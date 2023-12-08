package main

import (
	"dctl/pkg/initializers"
	"dctl/pkg/parsers/dctl"
	"dctl/pkg/transformers/compose"
	"dctl/pkg/transformers/k8"
	"dctl/pkg/transformers/sh"
	"dctl/pkg/version"
	"flag"
	"fmt"
)

func main() {
	shouldUpdate := flag.Bool("update", false, "Should update")
	shouldInit := flag.String("init", "", "Should init project and which")
	flag.Parse()

	if *shouldUpdate {
		version.UpdateVersion()
		fmt.Printf("Updated")
		return
	}

	if *shouldInit != "" {
		initializers.Initialize(*shouldInit)
		return
	}

	isOutdated := version.CheckVersion()
	if isOutdated {
		fmt.Printf("New version is out. Run dctl update to update your version.")
	}

	entity := dctl.ParseDctl()
	compose.Transform(&entity)
	sh.Transform(&entity)
	k8.Transform(&entity)
}
