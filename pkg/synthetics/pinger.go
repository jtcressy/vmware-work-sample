package synthetics

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	// "strings"
	"time"

	"github.com/jtcressy/vmware-work-sample/pkg/metrics"
)

func NewPinger(opts *Options) (*Pinger, error) {

	opts.defaults()

	pingerContext := &PingerContext{
		Context:      context.Background(),
		URLs:         opts.TestUrls,
		PingInterval: opts.FetchInterval,
		BindAddr:     opts.BindAddress,
	}

	return &Pinger{
		ctx:         pingerContext,
		pingResults: make(map[string]*PingResult),
	}, nil
}

type Pinger struct {
	ctx         *PingerContext
	pingResults map[string]*PingResult
}

func (p *Pinger) GetContext() *PingerContext {
	return p.ctx
}

func (p *Pinger) Start(ctx context.Context) error {
	log.Printf("Starting pinger with status server on %s\n", p.ctx.BindAddr)
	http.Handle("/status", p)
	srv := &http.Server{
		Addr:    p.ctx.BindAddr,
		Handler: http.DefaultServeMux,
	}

	go p.startPing(ctx)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		} else {
			log.Printf("Server stopped: %s\n", err)
		}
	}()
	<-ctx.Done()
	log.Println("Shutting down pinger")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	<-ctx.Done()
	return nil
}

func (p *Pinger) startPing(ctx context.Context) {
	ticker := time.NewTicker(p.ctx.PingInterval)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				p.PingAll()
			}
		}
	}()

	<-ctx.Done()
	ticker.Stop()
	done <- true
}

func (p *Pinger) PingAll() {
	for _, u := range p.ctx.URLs {
		log.Printf("Running ping for %s\n", u.String())
		start := time.Now()
		r, err := http.Get(u.String())
		if err != nil {
			log.Println(err)
		}
		elapsed := time.Since(start)

		pr := &PingResult{
			Url:          u,
			Up:           bool(200 <= r.StatusCode && r.StatusCode <= 299),
			ResponseTime: elapsed,
		}
		p.pingResults[u.String()] = pr
		metrics.RegisterPingResult(
			pr.Url.String(),
			pr.Up,
			pr.ResponseTime,
		)
	}
}

func (p *Pinger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	j, err := func() ([]byte, error) {
		if r.URL.Query().Get("pretty") == "true" {
			return json.MarshalIndent(p.pingResults, "", "    ")
		} else {
			return json.Marshal(p.pingResults)
		}
	}()
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.Write(j)
	}
}

type PingResult struct {
	Url          *url.URL      `json:"-"`
	Up           bool          `json:"up"`
	ResponseTime time.Duration `json:"response_time_ms"`
}

func (pr *PingResult) MarshalJSON() ([]byte, error) {
	type Alias PingResult
	return json.Marshal(&struct {
		ResponseTime int64 `json:"response_time_ms"`
		*Alias
	}{
		ResponseTime: pr.ResponseTime.Milliseconds(),
		Alias:        (*Alias)(pr),
	})
}
