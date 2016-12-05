package types

import (
	"github.com/twinj/uuid"
	yaml "gopkg.in/yaml.v2"
)

func init() {
	uuid.SwitchFormat(uuid.FormatHex)
}

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
	Services    Service `yaml:"services" json:"-"`
	Name        string  `yaml:"name" json:"name"`
	IconUrl     string  `yaml:"icon_url"json:"icon_url"`
	RemotePath  string  `yaml:"remote_path"json:"remote_path"`
	Description string  `yaml:"description"json:"description"`
}

func (ad AppDescriptor) GetBytes() ([]byte, error) {
	return yaml.Marshal(ad)
}

func (ad AppDescriptor) GetId() string {
	return uuid.NewV5(uuid.NameSpaceURL, uuid.Name(ad.Name)).String()
}

type Process struct {
	Command []string
	Image   string
	Ports   []string
	Ui      bool
	Sound   bool
	Input   bool
	Volumes []string
}

type Service struct {
	App    Process
	Remote Process
}
