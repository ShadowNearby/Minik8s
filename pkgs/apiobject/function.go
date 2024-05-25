package core

import "encoding/json"

type Function struct {
	Kind       string       `yaml:"kind" json:"kind,omitempty"`
	APIVersion string       `yaml:"apiVersion" json:"apiVersion,omitempty"`
	Status     VersionLabel `yaml:"status" json:"status,omitempty"`
	Name       string       `yaml:"name" json:"name"`
	Path       string       `yaml:"path" json:"path"`
}

func (r *Function) MarshalJSON() ([]byte, error) {
	type Alias Function
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}

func (r *Function) UnMarshalJSON(data []byte) error {
	type Alias Function
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}
