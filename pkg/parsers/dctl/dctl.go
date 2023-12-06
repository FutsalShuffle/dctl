package dctl

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func ParseDctl() {
	var entity DctlEntity
	b, _ := os.ReadFile("example.yaml")
	data := string(b)

	err := yaml.Unmarshal([]byte(data), &entity)

	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println(prettyPrint(entity))
}

func prettyPrint(i DctlEntity) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
