package dctl

type DctlEntity struct {
	Version    float32           `yaml:"version"`
	Name       string            `yaml:"name"`
	Docker     EnabledOnlyEntity `yaml:"docker"`
	K8         EnabledOnlyEntity `yaml:"k8"`
	Containers map[string]struct {
		Image       string            `yaml:"image"`
		Ports       map[int]string    `yaml:"ports"`
		Volumes     map[int]string    `yaml:"volumes"`
		Links       map[int]string    `yaml:"links"`
		DependsOn   map[int]string    `yaml:"depends_on"`
		Restart     string            `yaml:"restart"`
		Environment map[string]string `yaml:"environment"`
		Build       map[string]struct {
			Context    string         `yaml:"context"`
			Dockerfile string         `yaml:"dockerfile"`
			Args       map[int]string `yaml:"args"`
		} `yaml:"build"`
	} `yaml:"containers"`
	Deployments map[string]struct {
		Resources string `yaml:"resources"`
		Ingress   map[string]struct {
			Paths map[int]string `yaml:"paths"`
		} `yaml:"ingress"`
	} `yaml:"deployments"`
	Commands string `yaml:"commands"`
}
