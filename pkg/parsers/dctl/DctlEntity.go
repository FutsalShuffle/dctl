package dctl

type Handler struct {
	Exec struct {
		Command []string `yaml:"command,omitempty"`
	} `yaml:"exec,omitempty"`
	HttpGet struct {
		Host        string `yaml:"host,omitempty"`
		HttpHeaders []struct {
			Name  string `yaml:"name,omitempty"`
			Value string `yaml:"value,omitempty"`
		} `yaml:"httpHeaders,omitempty"`
		Path   string `yaml:"path,omitempty"`
		Port   string `yaml:"port,omitempty"`
		Scheme string `yaml:"scheme,omitempty" default:"http"`
	} `yaml:"httpGet,omitempty"`
	TCPSocket struct {
		Port string `yaml:"port"`
	} `yaml:"tcpSocket,omitempty"`
}

type Deployment struct {
	Ingress struct {
		Paths []struct {
			Path string
			Port string
		} `yaml:"paths"`
		Enabled bool     `yaml:"enabled" default:"false"`
		Type    string   `yaml:"type" default:"http"`
		Hosts   []string `yaml:"hosts,omitempty"`
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
		Enabled          bool          `yaml:"enabled" default:"false"`
		Namespace        string        `yaml:"namespace" default:"default"`
		Environments     []string      `yaml:"environments"`
		UseSealedSecrets bool          `yaml:"useSealedSecrets" default:"false"`
		Requirements     []Requirement `yaml:"requirements,omitempty"`
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
