package main

import (
	"dctl/pkg/parsers/dctl"
	"dctl/pkg/transformers/compose"
	"dctl/pkg/transformers/k8"
	"dctl/pkg/transformers/sh"
)

func main() {
	entity := dctl.ParseDctl()
	compose.Transform(&entity)
	sh.Transform(&entity)
	k8.Transform(&entity)
}
