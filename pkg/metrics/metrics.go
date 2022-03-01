package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	urlUp         *prometheus.GaugeVec
	urlResponseMs *prometheus.GaugeVec
)

func init() {
	urlUp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "sample",
		Subsystem: "external",
		Name:      "url_up",
		Help:      "Boolean status of whether a URL is considered up or down.",
	}, []string{"url"})
	urlResponseMs = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "sample",
		Subsystem: "external",
		Name:      "url_response_ms",
		Help:      "Response time in milliseconds it took for the URL to respond.",
	}, []string{"url"})
	prometheus.MustRegister(urlUp, urlResponseMs)
	http.Handle("/metrics", promhttp.Handler())
}

func RegisterPingResult(url string, up bool, responseTime time.Duration) {
	urlUp.With(prometheus.Labels{
		"url": url,
	}).Set(func() float64 {
		if up {
			return 1
		} else {
			return 0
		}
	}())
	urlResponseMs.With(prometheus.Labels{
		"url": url,
	}).Set(float64(responseTime.Milliseconds()))
}