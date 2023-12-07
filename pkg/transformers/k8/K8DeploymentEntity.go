package k8

type K8DeploymentEntity struct {
	Name        string
	Ports       []string          `yaml:"ports"`
	Volumes     []string          `yaml:"volumes"`
	Restart     string            `yaml:"restart"`
	Environment map[string]string `yaml:"environment"`
	ProjectName string
}
