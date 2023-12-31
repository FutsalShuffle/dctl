package k8

import (
	"dctl/pkg/parsers/dctl"
)

type K8DeploymentEntity struct {
	Deployment  dctl.Deployment
	Name        string
	ProjectName string
}
