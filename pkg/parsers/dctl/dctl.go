package dctl

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func ParseDctl() DctlEntity {
	var entity DctlEntity
	b, _ := os.ReadFile("dctl.yaml")
	data := string(b)

	err := yaml.Unmarshal([]byte(data), &entity)

	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println(prettyPrint(entity))

	return entity
}

func prettyPrint(i DctlEntity) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
