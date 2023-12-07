package main

import (
	"dctl/pkg/parsers/dctl"
	"dctl/pkg/transformers/compose"
)

func main() {
	entity := dctl.ParseDctl()
	compose.Transform(entity)
}
