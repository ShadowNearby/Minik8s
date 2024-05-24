package config

import "time"

const PrometheusNodeFilePath = "prometheus/sd_node.json"

const PrometheusPodFilePath = "prometheus/sd_pod.json"

const PrometheusScrapeInterval = 10 * time.Second
