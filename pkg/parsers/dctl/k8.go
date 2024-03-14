package dctl

type Resources struct {
	Limits struct {
		Cpu    string `yaml:"cpu,omitempty"`
		Memory string `yaml:"memory,omitempty"`
	} `yaml:"limits,omitempty"`
	Requests struct {
		Cpu    string `yaml:"cpu,omitempty"`
		Memory string `yaml:"memory,omitempty"`
	} `yaml:"requests,omitempty"`
}

type Requirement struct {
	Name       string `yaml:"name,omitempty"`
	Version    string `yaml:"version,omitempty"`
	Repository string `yaml:"repository,omitempty"`
}

type Probe struct {
	Handler             `yaml:",inline"`
	InitialDelaySeconds int64 `yaml:"initialDelaySeconds,omitempty"`
	TimeoutSeconds      int64 `yaml:"timeoutSeconds,omitempty"`
}

type SELinuxOptions struct {
	// User is a SELinux user label that applies to the container.
	// More info: http://releases.k8s.io/HEAD/docs/user-guide/labels.md
	User string `yaml:"user,omitempty"`

	// Role is a SELinux role label that applies to the container.
	// More info: http://releases.k8s.io/HEAD/docs/user-guide/labels.md
	Role string `yaml:"role,omitempty"`

	// Type is a SELinux type label that applies to the container.
	// More info: http://releases.k8s.io/HEAD/docs/user-guide/labels.md
	Type string `yaml:"type,omitempty"`

	// Level is SELinux level label that applies to the container.
	// More info: http://releases.k8s.io/HEAD/docs/user-guide/labels.md
	Level string `yaml:"level,omitempty"`
}

type SecurityContext struct {
	Capabilities   Capabilities   `yaml:"capabilities,omitempty"`
	Privileged     bool           `yaml:"privileged,omitempty"`
	SELinuxOptions SELinuxOptions `yaml:"seLinuxOptions,omitempty"`
	RunAsUser      int64          `yaml:"runAsUser,omitempty"`
	RunAsNonRoot   bool           `yaml:"runAsNonRoot,omitempty"`
}

type DeploymentContainer struct {
	Resources *Resources `yaml:"resources,omitempty"`
	Ports     []string   `yaml:"ports,omitempty"`
	Pvc       []struct {
		Name      string `yaml:"name"`
		MountPath string `yaml:"mountPath"`
	} `yaml:"pvc"`
	Image    string `yaml:"image"`
	EmptyDir struct {
		MountPath string `yaml:"mountPath"`
		Enabled   bool   `yaml:"enabled" default:"false"`
	} `yaml:"emptyDir"`
	Lifecycle *struct {
		PostStart Handler `yaml:"postStart,omitempty"`
		PreStop   Handler `yaml:"preStop,omitempty"`
	} `yaml:"lifecycle,omitempty"`
	Env             map[string]map[string]string `yaml:"env"`
	Args            []string                     `yaml:"args,omitempty"`
	Command         []string                     `yaml:"command,omitempty"`
	WorkingDir      string                       `yaml:"workingDir,omitempty"`
	ReadinessProbe  *Probe                       `yaml:"readinessProbe,omitempty"`
	LivenessProbe   *Probe                       `yaml:"livenessProbe,omitempty"`
	ImagePullPolicy string                       `yaml:"imagePullPolicy,omitempty" default:"Always"`
	SecurityContext *SecurityContext             `yaml:"securityContext,omitempty"`
	Hooks           struct {
		PreRollback  []string `yaml:"pre-rollback"`
		PostRollback []string `yaml:"post-rollback"`
		PreInstall   []string `yaml:"pre-install"`
		PostInstall  []string `yaml:"post-install"`
		PreUpgrade   []string `yaml:"pre-upgrade"`
		PostUpgrade  []string `yaml:"post-upgrade"`
	} `yaml:"hooks,omitempty"`
}

type Capabilities struct {
	// Added capabilities
	Add []string `yaml:"add,omitempty"`
	// Removed capabilities
	Drop []string `yaml:"drop,omitempty"`
}
