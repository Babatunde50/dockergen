package generator

type DockerfileTemplate struct {
	BuildCmd      string
	RunCmd        string
	Entrypoint    string
	Port          int
	UseMultiStage bool
	BinaryName    string
	Version       string
}

type DockerComposeTemplate struct {
	Version  string
	Name     string
	Services []Service
	Networks []Network
	Volumes  []Volume
	Configs  []Config
	Secrets  []Secret
}

type Service struct {
	Name          string
	Image         string
	Build         Build
	Ports         []string
	Environment   map[string]string
	EnvFile       []string
	DependsOn     []string
	Volumes       []string
	Networks      []string
	Restart       string
	HealthCheck   HealthCheck
	Deploy        Deploy
	Labels        map[string]string
	Command       string
	Entrypoint    string
	User          string
	WorkingDir    string
	ContainerName string
	ReadOnly      bool
}

type Build struct {
	Context    string
	Dockerfile string
	Args       map[string]string
	Target     string // For multi-stage builds
	CacheFrom  []string
}

type Network struct {
	Name       string
	Driver     string // e.g., "bridge", "overlay"
	External   bool
	Attachable bool
	IPAM       IPAM
	Labels     map[string]string
}

type IPAM struct {
	Driver string
	Config []IPAMConfig
}

type IPAMConfig struct {
	Subnet     string
	Gateway    string
	IPRange    string
	AuxAddress map[string]string
}

type Volume struct {
	Name     string
	Driver   string
	External bool
	Labels   map[string]string
	Options  map[string]string
}

type Config struct {
	Name     string
	File     string
	External bool
}

type Secret struct {
	Name     string
	File     string
	External bool
}

type HealthCheck struct {
	Test        []string
	Interval    string
	Timeout     string
	Retries     int
	StartPeriod string
}

type Deploy struct {
	Mode         string
	Replicas     int
	Resources    Resources
	UpdateConfig UpdateConfig
	Placement    Placement
}

type Resources struct {
	Limits       ResourceSpec
	Reservations ResourceSpec
}

type ResourceSpec struct {
	CPUs    string
	Memory  string
	Devices []DeviceSpec
}

type DeviceSpec struct {
	Capabilities []string
	Count        int
	Device       string
	Driver       string
}

type UpdateConfig struct {
	Parallelism     int
	Delay           string
	FailureAction   string
	MaxFailureRatio string
	Order           string
}

type Placement struct {
	Constraints []string
	Preferences []PlacementPreference
}

type PlacementPreference struct {
	Spread string
}
