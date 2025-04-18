package generator

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Babatunde50/dockergen/internal/detector"
)

// GenerateDockerfile creates a Dockerfile based on the detected project type
func GenerateDockerfile(project *detector.Project, useMultiStage bool) (string, error) {
	if project == nil {
		return "", fmt.Errorf("project cannot be nil")
	}

	tmpl := DockerfileTemplate{
		Port:          project.Port,
		UseMultiStage: useMultiStage,
		Version:       project.Version,
	}

	switch project.Type {
	case detector.Go:
		return generateGoDockerfile(project, tmpl)
	default:
		return "", fmt.Errorf("unsupported project type: %s", project.Type)
	}
}

// generateGoDockerfile creates a Dockerfile for Go projects
func generateGoDockerfile(project *detector.Project, tmpl DockerfileTemplate) (string, error) {

	// Extract binary name from entrypoint path
	binaryName := "app"
	if project.Entrypoint != "" {
		entrypoint := strings.TrimSuffix(project.Entrypoint, ".go")
		if entrypoint != "" {
			parts := strings.Split(entrypoint, "/")
			if len(parts) > 0 && parts[len(parts)-1] != "" {
				binaryName = parts[len(parts)-1]
			}
		}
	}

	tmpl.BinaryName = binaryName

	tmpl.Entrypoint = fmt.Sprintf("/app/%s", binaryName)

	tmpl.BuildCmd = fmt.Sprintf("CGO_ENABLED=0 go build -ldflags=\"-s -w\" -o /app/%s ./%s", binaryName, trimLastPart(project.Entrypoint))

	tmpl.RunCmd = fmt.Sprintf("/app/%s", binaryName)

	return renderDockerfile(goDockerfileTemplate, tmpl)
}

// renderDockerfile applies the template data to the specified template
func renderDockerfile(dockerfileTemplate string, tmpl DockerfileTemplate) (string, error) {
	t, err := template.New("dockerfile").Parse(dockerfileTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse dockerfile template: %v", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, tmpl); err != nil {
		return "", fmt.Errorf("failed to execute dockerfile template: %v", err)
	}

	return buf.String(), nil
}

func trimLastPart(path string) string {
	dir := filepath.Dir(path)
	return dir + string(filepath.Separator)
}
