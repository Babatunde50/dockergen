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
