package customisor

import (
	"gopkg.in/yaml.v2"
	"log"
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
