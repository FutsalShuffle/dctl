package dctl

import (
	yaml "gopkg.in/yaml.v3"
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

	//Default environments are dev and prod
	if len(entity.K8.Environments) == 0 {
		envs := make([]string, 2)
		envs[0] = "dev"
		envs[1] = "prod"
		entity.K8.Environments = envs
	}

	if entity.K8.Namespace == "" {
		entity.K8.Namespace = entity.Name
	}

	return entity
}
