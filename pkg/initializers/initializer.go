package initializers

import (
	"dctl/pkg/initializers/bitrix"
	"dctl/pkg/initializers/django"
	"dctl/pkg/initializers/laravel"
	"dctl/pkg/initializers/next"
	"dctl/pkg/initializers/symfony"
	"log"
)

type ProjectInitializer interface {
	Init()
}

func Initialize(projectType string) {
	initializers := map[string]ProjectInitializer{
		"laravel": laravel.Initializer{},
		"symfony": symfony.Initializer{},
		"bitrix":  bitrix.Initializer{},
		"django":  django.Initializer{},
		"next":    next.Initializer{},
	}

	val, exists := initializers[projectType]
	if !exists {
		log.Fatalln("Project type " + projectType + " does not exists")
	}

	val.Init()
}
