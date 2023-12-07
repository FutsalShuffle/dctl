package main

import (
	"dctl/pkg/parsers/dctl"
	"dctl/pkg/transformers/compose"
	"dctl/pkg/transformers/sh"
)

func main() {
	entity := dctl.ParseDctl()
	compose.Transform(&entity)
	sh.Transform(&entity)
}
