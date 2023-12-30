package k8

type K8DeploymentEntity struct {
	Name           string
	Ports          []string          `yaml:"ports"`
	Volumes        []string          `yaml:"volumes"`
	Environment    map[string]string `yaml:"environment"`
	ProjectName    string
	DockerRegistry string
	Containers     []string
	Ingress        struct {
		Paths []struct {
			Path string
			Port string
		} `yaml:"paths"`
		Enabled bool `yaml:"enabled" default:"false"`
	} `yaml:"ingress"`
	Restart   string `yaml:"restart" default:"Always"`
	Resources struct {
		Limits struct {
			Cpu    string `yaml:"cpu"`
			Memory string `yaml:"memory"`
		} `yaml:"limits"`
		Requests struct {
			Cpu    string `yaml:"cpu"`
			Memory string `yaml:"memory"`
		} `yaml:"requests"`
	} `yaml:"resources"`
}
