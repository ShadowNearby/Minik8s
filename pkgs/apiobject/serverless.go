package core

type Function struct {
	Kind       string `yaml:"kind" json:"kind"`
	APIVersion string `yaml:"apiVersion" json:"apiVersion"`
	Name       string `yaml:"name" json:"name"`
	Path       string `yaml:"path" json:"path"` // this path should be an absolute one
}

// only for controller
type TriggerRequest struct {
	Url    string `json:"url" yaml:"url"`
	Params []byte `json:"params" yaml:"params"`
}

type TriggerMessage struct {
	Kind   string `yaml:"kind" json:"kind"`
	Name   string `yaml:"name" json:"name"`
	Params string `yaml:"params" json:"params"`
	ID     string `yaml:"id,omitempty" json:"id,omitempty"`
}

type PingSource struct {
	ApiVersion string         `yaml:"apiVersion" json:"apiVersion"`
	Kind       string         `yaml:"kind" json:"kind"`
	MetaData   MetaData       `json:"metaData" yaml:"metaData"`
	ID         string         `json:"id,omitempty" yaml:"id,omitempty"`
	Spec       PingSourceSpec `yaml:"spec" json:"spec"`
}

/*
*  *  *  *  *
|  |  |  |  |
|  |  |  |  +--- 星期几 (0 - 6) (Sunday to Saturday)
|  |  |  +------ 月份 (1 - 12)
|  |  +--------- 日期 (1 - 31)
|  +------------ 小时 (0 - 23)
+--------------- 分钟 (0 - 59)
*/
// e.g. */1 * * * * 每分钟一次

type PingSourceSpec struct {
	Schedule string `json:"schedule" yaml:"schedule"`
	JsonData string `json:"jsonData" yaml:"jsonData"`
	Sink     Sink   `json:"sink" yaml:"sink"`
}

type Sink struct {
	Ref SinkRef `yaml:"ref" json:"ref"`
}

type SinkRef struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	Name       string `json:"name" yaml:"name"`
}

type TriggerResult struct {
	ID     string `json:"id"`
	Result string `json:"result"`
}
