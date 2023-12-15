package main

import (
	"dctl/pkg/initializers"
	"dctl/pkg/parsers/dctl"
	"dctl/pkg/transformers/compose"
	"dctl/pkg/transformers/gitlab"
	"dctl/pkg/transformers/k8"
	"dctl/pkg/transformers/sh"
	"dctl/pkg/version"
	"flag"
	"fmt"
	"os"
)

func main() {
	shouldUpdate := flag.Bool("update", false, "Should update")
	shouldInit := flag.String("init", "", "Should init project and which")
	shouldShowVersion := flag.Bool("version", false, "Print latest version")
	flag.Parse()

	if *shouldShowVersion {
		fmt.Println(version.Version)
		os.Exit(0)
	}
	if *shouldUpdate {
		version.UpdateVersion()
		fmt.Println("Updated")
		os.Exit(0)
	}
	if *shouldInit != "" {
		initializers.Initialize(*shouldInit)
		os.Exit(0)
	}

	isOutdated := version.CheckVersion()
	if isOutdated {
		fmt.Println("New version is out. Run dctl update to update your version.")
	}

	entity := dctl.ParseDctl()
	compose.Transform(&entity)
	gitlab.Transform(&entity)
	sh.Transform(&entity)
	k8.Transform(&entity)
}
