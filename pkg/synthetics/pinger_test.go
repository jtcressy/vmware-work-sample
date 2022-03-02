package synthetics

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestPinger(t *testing.T) {
	rawUrls := []string{"https://httpstat.us/200", "https://httpstat.us/503"}
	pinger, err := NewPinger(&Options{
		FetchInterval: 5 * time.Second,
		TestUrls: func() (urls MultipleUrls) {
			for _, u := range rawUrls {
				parsedUrl, err := url.Parse(u)
				if err != nil {
					t.Fatal(err)
				}
				urls = append(urls, parsedUrl)
			}
			return urls
		}(),
	})
	if err != nil {
		t.Fatal(err)
	}
	pinger.PingAll()

	t.Run("StatusGet", func(t *testing.T) {
		wr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/status", nil)
		pinger.ServeHTTP(wr, req)
		if wr.Code != http.StatusOK {
			t.Errorf("got HTTP status code %d, expected 200", wr.Code)
		}

		for _, rawUrl := range rawUrls {
			if !strings.Contains(wr.Body.String(), rawUrl) {
				t.Errorf(`status response body "%s" does not contain "%s"`, wr.Body.String(), rawUrl)
			}
		}
	})

	t.Run("MetricsGet", func(t *testing.T) {
		wr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		
		promhttp.Handler().ServeHTTP(wr, req)
		
		if wr.Code != http.StatusOK {
			t.Errorf("got HTTP status code %d, expected 200", wr.Code)
		}

		for _, rawUrl := range rawUrls {
			if !strings.Contains(wr.Body.String(), rawUrl) {
				t.Errorf(`metrics response body "%s" does not contain "%s"`, wr.Body.String(), rawUrl)
			}
		}

		metricNames := []string{
			"sample_external_url_up",
			"sample_external_url_response_ms",
		}
		for _,metricName := range metricNames {
			if !strings.Contains(wr.Body.String(), metricName) {
				t.Errorf(`metrics response body "%s" does not contain "%s"`, wr.Body.String(), metricName)
			}
		}
	})
}