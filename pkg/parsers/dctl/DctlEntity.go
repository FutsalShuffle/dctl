package dctl

type DctlEntity struct {
	Version    float32           `yaml:"version"`
	Name       string            `yaml:"name"`
	Docker     EnabledOnlyEntity `yaml:"docker" default:"true"`
	K8         EnabledOnlyEntity `yaml:"k8" default:"true"`
	Containers map[string]*struct {
		Image       string            `yaml:"image"`
		Ports       []string          `yaml:"ports"`
		Volumes     []string          `yaml:"volumes"`
		Links       []string          `yaml:"links"`
		DependsOn   []string          `yaml:"depends_on"`
		Restart     string            `yaml:"restart"`
		Environment map[string]string `yaml:"environment"`
		Command     []string          `yaml:"command"`
		Build       struct {
			Context    string            `yaml:"context"`
			Dockerfile string            `yaml:"dockerfile"`
			Args       map[string]string `yaml:"args"`
		} `yaml:"build"`
	} `yaml:"containers"`
	Deployments map[string]struct {
		//Resources string `yaml:"resources"`
		Ingress map[string]struct {
			Paths []string `yaml:"paths"`
		} `yaml:"ingress"`
	} `yaml:"deployments"`
	Commands struct {
		Db struct {
			Vendor    string `default:"mysql" yaml:"vendor"`
			Container string `default:"mysql" yaml:"container"`
		} `yaml:"db"`
		Run struct {
			Container string `default:"php" yaml:"container"`
		} `yaml:"run"`
		Extra []struct {
			Name    string `yaml:"name"`
			Command string `yaml:"command"`
		} `yaml:"extra"`
	} `yaml:"commands"`
}
