package initializers

import (
	"dctl/pkg/initializers/laravel"
	"dctl/pkg/initializers/symfony"
	"log"
	"os"
)

type ProjectInitializer interface {
	Init()
}

func Initialize(projectType string) {
	initializers := map[string]ProjectInitializer{
		"laravel": laravel.Initializer{},
		"symfony": symfony.Initializer{},
	}

	val, exists := initializers[projectType]
	if !exists {
		log.Fatalln("Project type " + projectType + " does not exists")
	}

	val.Init()
	os.Exit(0)
}
