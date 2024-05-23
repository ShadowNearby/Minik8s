package handler

import "github.com/prometheus/client_golang/prometheus"

func PrometheusRegister() {
	for _, collector := range Collectors {
		prometheus.MustRegister(collector)
	}
}
