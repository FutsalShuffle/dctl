package k8

type SecretsEntity struct {
	Deployment map[string]struct {
		Secrets    map[string]string `yaml:"secrets"`
		Containers map[string]struct {
			Env map[string]string `yaml:"env"`
		} `yaml:"containers"`
	} `yaml:"deployments"`
}
