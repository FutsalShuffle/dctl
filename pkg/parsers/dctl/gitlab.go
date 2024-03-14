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
