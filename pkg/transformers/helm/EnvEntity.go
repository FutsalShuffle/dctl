package helm

import "dctl/pkg/parsers/dctl"

type EnvEntity struct {
	Environment string
	Entity      *dctl.DctlEntity
}
