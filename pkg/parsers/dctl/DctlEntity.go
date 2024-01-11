package dctl

type GitlabStageStruct struct {
	Name          string                 `yaml:"name,omitempty"`
	Variables     map[string]string      `yaml:"variables,omitempty"`
	Image         string                 `yaml:"image,omitempty"`
	BeforeScript  []string               `yaml:"before_script,omitempty"`
	Artifacts     []string               `yaml:"artifacts,omitempty"`
	Script        []string               `yaml:"script,omitempty"`
	AfterScript   []string               `yaml:"after_script,omitempty"`
	Cache         map[string]interface{} `yaml:"cache,omitempty"`
	AllowFailure  bool                   `yaml:"allow_failure" default:"false"`
	Services      map[string]interface{} `yaml:"services,omitempty"`
	Require       string                 `yaml:"require,omitempty"`
	Timeout       int                    `yaml:"timeout,omitempty"`
	Tags          []string               `yaml:"tags,omitempty"`
	Only          []string               `yaml:"only,omitempty"`
	Interruptible bool                   `yaml:"interruptible,omitempty" default:"false"`
	Environment   string                 `yaml:"environment,omitempty"`
	Retry         map[string]interface{} `yaml:"retry,omitempty"`
	Release       map[string]interface{} `yaml:"release,omitempty"`
	Needs         []interface{}          `yaml:"needs,omitempty"`
	When          string                 `yaml:"when,omitempty"`
	Secrets       map[string]interface{} `yaml:"secrets,omitempty"`
	Rules         []interface{}          `yaml:"rules,omitempty"`
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
	Restart string `yaml:"restart" default:"Always"`
	Pvc     []struct {
		Storage  string `yaml:"storage"`
		Name     string `yaml:"name"`
		HostPath string `yaml:"hostPath"`
	} `yaml:"pvc"`
	EmptyDir struct {
		SizeLimit string `yaml:"sizeLimit"`
		Enabled   bool   `yaml:"enabled" default:"false"`
	} `yaml:"emptyDir"`
	Replicas   int                            `yaml:"replicas"`
	Service    bool                           `yaml:"service" default:"true"`
	Containers map[string]DeploymentContainer `yaml:"containers"`
	Secrets    map[string]map[string]string   `yaml:"secrets"`
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
	Envs map[string]map[string]string `yaml:"envs"`
}

type ComposeContainer struct {
	CPU         int               `yaml:"cpu_shares,omitempty"`
	DNS         []string          `yaml:"dns,omitempty"`
	Domain      []string          `yaml:"dns_search,omitempty"`
	EnvFile     []string          `yaml:"env_file,omitempty"`
	Expose      []int             `yaml:"expose,omitempty"`
	Hostname    string            `yaml:"hostname,omitempty"`
	Memory      int               `yaml:"mem_limit,omitempty"`
	Network     []string          `yaml:"networks,omitempty"`
	NetworkMode string            `yaml:"network_mode,omitempty"`
	Pid         string            `yaml:"pid,omitempty"`
	Privileged  bool              `yaml:"privileged,omitempty"`
	User        string            `yaml:"user,omitempty"`
	WorkDir     string            `yaml:"working_dir,omitempty"`
	Image       string            `yaml:"image,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Links       []string          `yaml:"links,omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty"`
	Restart     string            `yaml:"restart,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Entrypoint  string            `yaml:"entrypoint,omitempty"`
	Command     []string          `yaml:"command,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Build       struct {
		Context    string            `yaml:"context,omitempty"`
		Dockerfile string            `yaml:"dockerfile,omitempty"`
		Args       map[string]string `yaml:"args,omitempty"`
	} `yaml:"build,omitempty"`
}

type DctlEntity struct {
	Version float32 `yaml:"version"`
	Name    string  `yaml:"name"`
	K8      struct {
		Enabled          bool     `yaml:"enabled" default:"false"`
		Namespace        string   `yaml:"namespace" default:"default"`
		Environments     []string `yaml:"environments"`
		UseSealedSecrets bool     `yaml:"useSealedSecrets" default:"false"`
	}
	Docker struct {
		Enabled  bool   `yaml:"enabled" default:"true"`
		Registry string `yaml:"registry"`
	} `yaml:"docker"`
	Containers  map[string]*ComposeContainer `yaml:"containers"`
	Deployments map[string]Deployment        `yaml:"deployments"`
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
