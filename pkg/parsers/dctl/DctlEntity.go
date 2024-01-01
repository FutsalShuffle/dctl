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

type LifecycleHandler struct {
	Exec struct {
		Command []string `yaml:"command"`
	} `yaml:"exec"`
	HttpGet struct {
		Host        string `yaml:"host"`
		HttpHeaders []struct {
			Name  string `yaml:"name"`
			Value string `yaml:"value"`
		} `yaml:"httpHeaders"`
		Path   string `yaml:"path"`
		Port   string `yaml:"port"`
		Scheme string `yaml:"scheme" default:"http"`
	} `yaml:"httpGet"`
}

type Deployment struct {
	Ingress struct {
		Paths []struct {
			Path string
			Port string
		} `yaml:"paths"`
		Enabled bool `yaml:"enabled" default:"false"`
	} `yaml:"ingress"`
	Secret  bool   `yaml:"secret" default:"false"`
	Restart string `yaml:"restart" default:"Always"`
	Pvc     []struct {
		Storage string `yaml:"storage"`
		Name    string `yaml:"name"`
	} `yaml:"pvc"`
	EmptyDir struct {
		SizeLimit string `yaml:"sizeLimit"`
		Enabled   bool   `yaml:"enabled" default:"false"`
	} `yaml:"emptyDir"`
	Replicas   int                            `yaml:"replicas"`
	Service    bool                           `yaml:"service" default:"true"`
	Containers map[string]DeploymentContainer `yaml:"containers"`
}

type DeploymentContainer struct {
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
	Ports []string `yaml:"ports"`
	Pvc   []struct {
		Name      string `yaml:"name"`
		MountPath string `yaml:"mountPath"`
	} `yaml:"pvc"`
	Image    string `yaml:"image"`
	EmptyDir struct {
		MountPath string `yaml:"mountPath"`
		Enabled   bool   `yaml:"enabled" default:"false"`
	} `yaml:"emptyDir"`
	Env       map[string]string `yaml:"env"`
	Lifecycle struct {
		PostStart LifecycleHandler `yaml:"postStart"`
		PreStop   LifecycleHandler `yaml:"preStop"`
	}
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
	Deployments map[string]Deployment `yaml:"deployments"`
	Commands    struct {
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
