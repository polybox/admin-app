package types

import yaml "gopkg.in/yaml.v2"

type Application struct {
	Id          string        `json:"id"`
	Name        string        `json:"name"`
	IsRunning   bool          `json:"is_running"`
	IconUrl     string        `json:"icon_url"`
	Descriptor  AppDescriptor `json:"-"`
	Description string        `json:"description"`
	RemotePath  string        `json:"remote_path,omitempty"`
	RemotePort  string        `json:"remote_port"`
}

type AppDescriptor struct {
	Services    Service `yaml:"services"`
	Name        string  `yaml:"name"`
	IconUrl     string  `yaml:"icon_url"`
	RemotePath  string  `yaml:"remote_path"`
	Description string  `yaml:"description"`
}

func (ad AppDescriptor) GetBytes() ([]byte, error) {
	return yaml.Marshal(ad)
}

type Process struct {
	Command []string
	Image   string
	Ports   []string
	Ui      bool
	Sound   bool
}

type Service struct {
	App    Process
	Remote Process
}
