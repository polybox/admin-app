package types

import yaml "gopkg.in/yaml.v2"

type Application struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	IsRunning   bool   `json:"is_running"`
	IconUrl     string `json:"icon_url"`
	Descriptor  []byte `json:"-"`
	Description string `json:"description"`
	RemoteUrl   string `json:"remote_url,omitempty"`
}

type Process struct {
	Image string
	Ports []string
	Ui    bool
	Sound bool
}

type Service struct {
	App    Process
	Remote Process
}

type AppDescriptor struct {
	Services    Service `yaml:"services"`
	Name        string  `yaml:"name"`
	IconUrl     string  `yaml:"icon_url"`
	RemoteUrl   string  `yaml:"remote_url"`
	Description string  `yaml:"description"`
}

func (ad AppDescriptor) GetBytes() ([]byte, error) {
	return yaml.Marshal(ad)
}
