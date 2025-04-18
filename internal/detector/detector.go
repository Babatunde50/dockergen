package detector

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type ProjectType string

const (
	Go     ProjectType = "go"
	NodeJS ProjectType = "nodejs"
	Python ProjectType = "python"
)

type Project struct {
	Type       ProjectType
	Entrypoint string
	Port       int
	WorkDir    string
	Version    string
}

func DetectProject(rootDir string) (*Project, error) {
	// Check if directory exists
	fileInfo, err := os.Stat(rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to access directory %s: %v", rootDir, err)
	}

	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", rootDir)
	}

	// Initialize project
	project := &Project{
		WorkDir: rootDir,
		Port:    detectPort(rootDir),
	}

	// Check for project type
	if isGoProject(rootDir) {
		project.Type = Go
		project.Entrypoint = findGoEntrypoint(rootDir)
		project.Version = detectVersion(rootDir, Go)
	} else if isNodeJSProject(rootDir) {
		project.Type = NodeJS
		project.Entrypoint = findNodeJSEntrypoint(rootDir)
		project.Version = detectVersion(rootDir, NodeJS)
	} else if isPythonProject(rootDir) {
		project.Type = Python
		project.Entrypoint = findPythonEntrypoint(rootDir)
		project.Version = detectVersion(rootDir, Python)
	} else {
		return nil, fmt.Errorf("unable to determine project type in %s", rootDir)
	}

	return project, nil
}

// isGoProject checks if the directory contains Go project indicators
func isGoProject(dir string) bool {
	// Look for go.mod, go.sum
	if fileExists(filepath.Join(dir, "go.mod")) ||
		fileExists(filepath.Join(dir, "go.sum")) {
		return true
	}

	// Check for .go files
	goFiles, _ := filepath.Glob(filepath.Join(dir, "*.go"))
	return len(goFiles) > 0
}

// isNodeJSProject checks if the directory contains NodeJS project indicators
func isNodeJSProject(dir string) bool {
	// Look for package.json
	if fileExists(filepath.Join(dir, "package.json")) {
		return true
	}

	// Look for node_modules
	if dirExists(filepath.Join(dir, "node_modules")) {
		return true
	}

	// Check for .js files
	jsFiles, _ := filepath.Glob(filepath.Join(dir, "*.js"))

	return len(jsFiles) > 0
}

// isPythonProject checks if the directory contains Python project indicators
func isPythonProject(dir string) bool {
	// Look for requirements.txt, setup.py, or Pipfile
	if fileExists(filepath.Join(dir, "requirements.txt")) ||
		fileExists(filepath.Join(dir, "setup.py")) ||
		fileExists(filepath.Join(dir, "Pipfile")) {
		return true
	}

	// Check for .py files
	pyFiles, _ := filepath.Glob(filepath.Join(dir, "*.py"))
	if len(pyFiles) > 0 {
		return true
	}

	// Look for venv or .venv directories
	if dirExists(filepath.Join(dir, "venv")) || dirExists(filepath.Join(dir, ".venv")) {
		return true
	}

	return false
}

// findGoEntrypoint attempts to locate the main Go file
func findGoEntrypoint(dir string) string {
	// Common patterns for Go entrypoints
	candidates := []string{
		filepath.Join(dir, "main.go"),
		filepath.Join(dir, "cmd", "main.go"),
		filepath.Join(dir, "cmd", filepath.Base(dir), "main.go"),
	}

	for _, candidate := range candidates {
		if fileExists(candidate) {
			relPath, _ := filepath.Rel(dir, candidate)
			return relPath
		}
	}

	// Look for any file with a main function
	var mainFile string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		if strings.Contains(string(content), "func main()") {
			mainFile = path
			return filepath.SkipDir // Stop walking once found
		}

		return nil
	})

	if mainFile != "" {
		relPath, _ := filepath.Rel(dir, mainFile)
		return relPath
	}

	return ""
}

// findNodeJSEntrypoint attempts to locate the main NodeJS file
func findNodeJSEntrypoint(dir string) string {
	// Check package.json for main field
	if fileExists(filepath.Join(dir, "package.json")) {
		data, err := os.ReadFile(filepath.Join(dir, "package.json"))
		if err == nil {
			content := string(data)
			// Simple regex to extract the main field
			re := regexp.MustCompile(`"main"\s*:\s*"([^"]+)"`)
			matches := re.FindStringSubmatch(content)
			if len(matches) > 1 {
				return matches[1]
			}
		}
	}

	// Common patterns for Node.js entrypoints
	candidates := []string{
		filepath.Join(dir, "index.js"),
		filepath.Join(dir, "server.js"),
		filepath.Join(dir, "app.js"),
		filepath.Join(dir, "main.js"),
		filepath.Join(dir, "src", "index.js"),
	}

	for _, candidate := range candidates {
		if fileExists(candidate) {
			relPath, _ := filepath.Rel(dir, candidate)
			return relPath
		}
	}

	return ""
}

// findPythonEntrypoint attempts to locate the main Python file
func findPythonEntrypoint(dir string) string {
	// Common patterns for Python entrypoints
	candidates := []string{
		filepath.Join(dir, "app.py"),
		filepath.Join(dir, "main.py"),
		filepath.Join(dir, "run.py"),
		filepath.Join(dir, filepath.Base(dir)+".py"),
	}

	for _, candidate := range candidates {
		if fileExists(candidate) {
			relPath, _ := filepath.Rel(dir, candidate)
			return relPath
		}
	}

	// Check for files with if __name__ == "__main__": pattern
	var mainFile string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".py") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		if strings.Contains(string(content), `if __name__ == "__main__"`) ||
			strings.Contains(string(content), `if __name__ == '__main__'`) {
			mainFile = path
			return filepath.SkipDir // Stop walking once found
		}

		return nil
	})

	if mainFile != "" {
		relPath, _ := filepath.Rel(dir, mainFile)
		return relPath
	}

	return ""
}

// detectPort tries to find the port number used in the project looking for the PORT variable in the environment files
func detectPort(dir string) int {

	// First check environment files
	envFiles := []string{
		filepath.Join(dir, ".env"),
		filepath.Join(dir, ".env.development"),
		filepath.Join(dir, ".env.local"),
	}

	for _, envFile := range envFiles {
		if fileExists(envFile) {
			content, err := os.ReadFile(envFile)
			if err == nil {
				portEnvPattern := regexp.MustCompile(`PORT\s*=\s*(\d+)`)
				matches := portEnvPattern.FindStringSubmatch(string(content))
				if len(matches) > 1 {
					port, _ := strconv.Atoi(matches[1])
					return port
				}
			}
		}
	}

	return 3000
}

// Helper functions
func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func dirExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// detectVersion attempts to extract the version from the project
func detectVersion(dir string, projectType ProjectType) string {
	switch projectType {
	case Go:
		return detectGoVersion(dir)
	case NodeJS:
		// TODO: Implement Node.js version detection
		return ""
	case Python:
		// TODO: Implement Python version detection
		return ""
	default:
		return ""
	}
}

func detectGoVersion(dir string) string {
	// Default Go version if we can't detect it
	defaultVersion := "1.22"

	// Try to read go.mod file
	goModPath := filepath.Join(dir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return defaultVersion
	}

	// Parse go.mod for the Go version using regex
	re := regexp.MustCompile(`go\s+(\d+\.\d+)`)
	matches := re.FindStringSubmatch(string(content))
	if len(matches) > 1 {
		return matches[1]
	}

	return defaultVersion
}
