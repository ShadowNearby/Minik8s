package handler

import "github.com/prometheus/client_golang/prometheus"

var Collectors = [...]prometheus.Collector{
	prometheus.NewCounter(prometheus.CounterOpts{
		Name: "custom_counter",
		Help: "This is a custom counter",
	}),
}