package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"github.com/go-ini/ini"
)

type Customisor interface {
	ApplyCustomisations()
}


type BasicYamlData struct {
	Enabled bool	`yaml:"enabled"`
}

func (basic *BasicYamlData) BasicYamlLoader(data []byte) {
	err := yaml.Unmarshal(data, basic)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if basic.Enabled {
		log.Println("'Basic' enabled")
	} else {
		log.Println("'Basic' disabled")
	}
}

func Customise(content []byte, section *ini.Section, configurationFileName string) (bool) {
	if configurationFileName == "configuration-basic.yaml" {
		log.Println("Process as basic/yaml")
		var basic BasicYamlData
		basic.BasicYamlLoader(content)
		return true
	}
	return false
}
