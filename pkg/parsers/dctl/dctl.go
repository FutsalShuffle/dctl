package dctl

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func ParseDctl() DctlEntity {
	var entity DctlEntity
	b, err := os.ReadFile("./dctl.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	data := string(b)

	err = yaml.Unmarshal([]byte(data), &entity)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return entity
}
