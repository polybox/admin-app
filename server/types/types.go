package types

import yaml "gopkg.in/yaml.v2"

type Installation struct {
	Id            string
	ApplicationId string
	Status        string
	Descriptor    []byte
}

type Process struct {
	Image string
	Ports []string
}

type Service struct {
	App Process
	Web Process
}

type AppDescriptor struct {
	Services Service
	Name     string
}

func (ad AppDescriptor) GetBytes() ([]byte, error) {
	return yaml.Marshal(ad)
}
