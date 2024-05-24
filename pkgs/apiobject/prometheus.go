package core

type PrometheusSdConfig struct {
	Targets []string          `json:"targets" yaml:"targets"`
	Labels  map[string]string `json:"labels" yaml:"labels"`
}
