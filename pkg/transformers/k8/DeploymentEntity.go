package k8

import (
	"dctl/pkg/parsers/dctl"
)

type DeploymentEntity struct {
	Deployment  dctl.Deployment
	Name        string
	ProjectName string
	Namespace   string
	Environment string
	Secrets     map[string]string
}
