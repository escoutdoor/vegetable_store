package metrics

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	requestCounter        prometheus.Counter
	responseCounter       *prometheus.CounterVec
	histogramResponseTime *prometheus.HistogramVec
}

var m *metrics

func Init() {
	m = &metrics{
		requestCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "",
		}),
	}
}
