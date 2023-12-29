package dctl

type GitlabStageStruct struct {
	Name   string `yaml:"name"`
	Docker struct {
		Image string `yaml:"image"`
		Build struct {
			Args map[string]string `yaml:"args"`
		}
	}
	Before       []string `yaml:"before"`
	Scripts      []string `yaml:"scripts"`
	After        []string `yaml:"after"`
	AllowFailure bool     `yaml:"allow_failure" default:"false"`
	Services     []string `yaml:"services"`
	Require      string   `yaml:"require"`
	Timeout      int      `yaml:"timeout"`
	Tags         string   `yaml:"tags"`
	Only         []string `yaml:"only"`
}

type DctlEntity struct {
	Version float32 `yaml:"version"`
	Name    string  `yaml:"name"`
	K8      struct {
		Enabled bool `yaml:"enabled" default:"false"`
	}
	Docker struct {
		Enabled  bool   `yaml:"enabled" default:"true"`
		Registry string `yaml:"registry"`
	} `yaml:"docker"`
	Containers map[string]*struct {
		Image       string            `yaml:"image"`
		Ports       []string          `yaml:"ports"`
		Volumes     []string          `yaml:"volumes"`
		Links       []string          `yaml:"links"`
		DependsOn   []string          `yaml:"depends_on"`
		Restart     string            `yaml:"restart"`
		Environment map[string]string `yaml:"environment"`
		Entrypoint  string            `yaml:"entrypoint"`
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
	Gitlab struct {
		Cache struct {
			Paths []string `yaml:"paths"`
		} `yaml:"cache"`
		OnlyWhen GitlabWorkflow      `yaml:"only_when" default:"merge_request"`
		Tests    []GitlabStageStruct `yaml:"tests"`
		Deploy   []GitlabStageStruct `yaml:"deploy"`
	} `yaml:"gitlab"`
}

type GitlabWorkflow string

const (
	MERGE_REQUEST        GitlabWorkflow = "merge_request"
	ALWAYS               GitlabWorkflow = "always"
	MERGE_REQUEST_MASTER GitlabWorkflow = "merge_request_master"
	NEVER                GitlabWorkflow = "never"
)
