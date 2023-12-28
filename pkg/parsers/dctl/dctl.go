package dctl

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func ParseDctl() DctlEntity {
	var entity DctlEntity
	b, _ := os.ReadFile("./dctl.yaml")
	data := string(b)

	err := yaml.Unmarshal([]byte(data), &entity)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return entity
}
