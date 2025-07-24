package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	requestCounter        prometheus.Counter
	responseCounter       *prometheus.CounterVec
	histogramResponseTime *prometheus.HistogramVec
}

var metrics *Metrics

func Init(namespace, appName string) {
	metrics = &Metrics{
		requestCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "http",
			Name:      appName + "_requests_total",
			Help:      "Number of requests",
		}),
		responseCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "http",
			Name:      appName + "_responses_total",
			Help:      "Number of responses",
		}, []string{"status", "request"}),
		histogramResponseTime: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "http",
			Name:      appName + "_histogram_response_time_seconds",
			Help:      "Response time",
			Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
		}, []string{"status"}),
	}
}

func IncRequestCounter() {
	metrics.requestCounter.Inc()
}

func IncResponseCounter(status string, method string) {
	metrics.responseCounter.WithLabelValues(status, method).Inc()
}

func HistogramResponseTimeObserve(status string, time float64) {
	metrics.histogramResponseTime.WithLabelValues(status).Observe(time)
}
