package synthetics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	// "strings"
	"time"

	"github.com/jtcressy/vmware-work-sample/pkg/metrics"
)

func NewPinger(opts Options) (*Pinger, error) {

	opts.defaults()

	pingerContext := &PingerContext{
		URLs:         opts.TestUrls,
		PingInterval: opts.FetchInterval,
		BindAddr:     opts.BindAddress,
	}

	return &Pinger{
		ctx:         pingerContext,
		pingResults: make(map[string]PingResult),
	}, nil
}

type Pinger struct {
	ctx         *PingerContext
	pingResults map[string]PingResult
}

func (p *Pinger) Start(ctx context.Context) error {
	http.Handle("/status", p)

	go p.startPing(ctx)
	http.ListenAndServe(p.ctx.BindAddr, nil)

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
				// do url fetches
				fmt.Printf("Running Pings at %v\n", time.Now())
				for _, u := range p.ctx.URLs {
					start := time.Now()
					r, err := http.Get(u.String())
					if err != nil {
						fmt.Println(err)
					}
					elapsed := time.Since(start)

					pr := PingResult{
						Url:          u,
						Up:           bool(200 <= r.StatusCode && r.StatusCode <= 299),
						ResponseTime: ResponseTimeMs(elapsed),
					}
					p.pingResults[u.String()] = pr
					metrics.RegisterPingResult(
						pr.Url.String(),
						pr.Up,
						time.Duration(pr.ResponseTime),
					)
				}
			}
		}
	}()

	<-ctx.Done()
	ticker.Stop()
	done <- true
}

func (p *Pinger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(p.pingResults, "", "    ")
	if err != nil {
		fmt.Println(err)
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
	ResponseTime ResponseTimeMs `json:"responseTimeMs"`
}

type ResponseTimeMs time.Duration

func (rtm *ResponseTimeMs) MarshalJSON() ([]byte, error) {
	fmt.Println("marshalling")
	rt := time.Duration(*rtm)
	return json.Marshal(rt.String())
}
// func (pr *PingResult) UnmarshalJSON(j []byte) error {
// 	var rawStrings map[string]string
// 	err := json.Unmarshal(j, &rawStrings)
// 	if err != nil {
// 		return err
// 	}

// 	for k, v := range rawStrings {
// 		if strings.ToLower(k) == "url" {
// 			pr.Url, err = url.Parse(v)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		if strings.ToLower(k) == "up" {
// 			pr.Up = bool(v == "1" || v == "true")
// 		}
// 		if k == "responseTimeMs" {
// 			pr.ResponseTime, err = time.ParseDuration(v)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// func (pr *PingResult) MarshalJSON() ([]byte, error) {
// 	basicPingResult := struct {
// 		Url          string `json:"url"`
// 		Up           bool   `json:"up"`
// 		ResponseTime int    `json:"responseTimeMs"`
// 	}{
// 		Url:          pr.Url.RequestURI(),
// 		Up:           pr.Up,
// 		ResponseTime: int(pr.ResponseTime.Milliseconds()),
// 	}
// 	return json.Marshal(basicPingResult)
// }
