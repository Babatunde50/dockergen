package detector

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestDir creates a temporary directory with files for testing
func setupTestDir(t *testing.T) string {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "detector-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	return tempDir
}

// createFile creates a file with the given content
func createFile(t *testing.T, path string, content string) {
	dir := filepath.Dir(path)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create file %s: %v", path, err)
	}
}

func TestDetectProject(t *testing.T) {
	tests := []struct {
		name               string
		setupFunc          func(string) // Function to set up test files
		expectedType       ProjectType
		expectedEntrypoint string
	}{
		{
			name: "Go Project",
			setupFunc: func(dir string) {
				// Create go.mod file
				createFile(t, filepath.Join(dir, "go.mod"), `module github.com/example/myproject
go 1.18
`)
				// Create main.go
				createFile(t, filepath.Join(dir, "main.go"), `package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	http.ListenAndServe(":8080", nil)
}
`)
			},
			expectedType:       Go,
			expectedEntrypoint: "main.go",
		},
		{
			name: "NodeJS Project",
			setupFunc: func(dir string) {
				// Create package.json
				createFile(t, filepath.Join(dir, "package.json"), `{
  "name": "myproject",
  "version": "1.0.0",
  "main": "server.js",
  "scripts": {
    "start": "node server.js",
    "build": "webpack"
  }
}
`)
				// Create server.js
				createFile(t, filepath.Join(dir, "server.js"), `
const express = require('express');
const app = express();
const PORT = 4000;

app.get('/', (req, res) => {
  res.send('Hello World!');
});

app.listen(PORT, () => {
  console.log('Server started on port ' + PORT);
});
`)
			},
			expectedType:       NodeJS,
			expectedEntrypoint: "server.js",
		},
		{
			name: "Python Project",
			setupFunc: func(dir string) {
				// Create requirements.txt
				createFile(t, filepath.Join(dir, "requirements.txt"), `
Flask==2.0.1
requests==2.26.0
`)
				// Create app.py
				createFile(t, filepath.Join(dir, "app.py"), `
from flask import Flask

app = Flask(__name__)

@app.route('/')
def hello():
    return 'Hello, World!'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
`)
			},
			expectedType:       Python,
			expectedEntrypoint: "app.py",
		},
		{
			name: "Complex Go Project Structure",
			setupFunc: func(dir string) {
				// Create go.mod file
				createFile(t, filepath.Join(dir, "go.mod"), `module github.com/example/complex-project
go 1.18
`)
				// Create cmd/myapp/main.go
				createFile(t, filepath.Join(dir, "cmd", "myapp", "main.go"), `package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	http.ListenAndServe(":8080", nil)
}
`)
			},
			expectedType:       Go,
			expectedEntrypoint: filepath.Join("cmd", "myapp", "main.go"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary directory for this test case
			tempDir := setupTestDir(t)

			// Set up the test files
			tc.setupFunc(tempDir)

			// Run detection
			project, err := DetectProject(tempDir)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			// Check project type
			if project.Type != tc.expectedType {
				t.Errorf("Expected project type %s, got %s", tc.expectedType, project.Type)
			}

			// Check entrypoint
			if project.Entrypoint != tc.expectedEntrypoint {
				t.Errorf("Expected entrypoint %s, got %s", tc.expectedEntrypoint, project.Entrypoint)
			}
		})
	}
}

func TestPort(t *testing.T) {
	// Create a test directory
	tempDir := setupTestDir(t)

	// Test .env file PORT detection
	createFile(t, filepath.Join(tempDir, ".env"), `
PORT=5678
DEBUG=true
`)

	// Test port detection
	port := detectPort(tempDir)
	if port != 5678 {
		t.Errorf("Expected port 5678 from .env file, got %d", port)
	}
}

func TestDetectionInNonExistentDirectory(t *testing.T) {
	// Try to detect in a non-existent directory
	_, err := DetectProject("/path/that/does/not/exist")
	if err == nil {
		t.Errorf("Expected an error when detecting in non-existent directory, but got nil")
	}
}

func TestDetectionInEmptyDirectory(t *testing.T) {
	// Create an empty directory
	tempDir := setupTestDir(t)

	// Try to detect the project
	_, err := DetectProject(tempDir)
	if err == nil {
		t.Errorf("Expected an error when detecting in empty directory, but got nil")
	}
}

func TestFileAndDirHelpers(t *testing.T) {
	// Create a temporary directory
	tempDir := setupTestDir(t)

	// Test non-existent file
	nonExistentPath := filepath.Join(tempDir, "non-existent.txt")
	if fileExists(nonExistentPath) {
		t.Errorf("Expected non-existent file to return false")
	}

	// Create a file and test fileExists
	filePath := filepath.Join(tempDir, "test.txt")
	createFile(t, filePath, "test content")
	if !fileExists(filePath) {
		t.Errorf("Expected existing file to return true")
	}

	// Create a directory and test dirExists
	dirPath := filepath.Join(tempDir, "testdir")
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	if !dirExists(dirPath) {
		t.Errorf("Expected existing directory to return true")
	}

	// Test directory with fileExists (should be false)
	if fileExists(dirPath) {
		t.Errorf("Expected fileExists to return false for directory")
	}

	// Test file with dirExists (should be false)
	if dirExists(filePath) {
		t.Errorf("Expected dirExists to return false for file")
	}
}

func TestDetectGoVersion(t *testing.T) {
	// Create a temporary directory
	tempDir := setupTestDir(t)

	// Create a go.mod file with a specific version
	createFile(t, filepath.Join(tempDir, "go.mod"), `module example.com/myproject

go 1.21.0

require (
	github.com/example/pkg v1.0.0
)
`)

	// Test version detection
	version := detectGoVersion(tempDir)
	if version != "1.21" {
		t.Errorf("Expected Go version 1.21, got %s", version)
	}

	// Test with a different go.mod format
	otherDir := setupTestDir(t)
	createFile(t, filepath.Join(otherDir, "go.mod"), `module test
go 1.19
`)

	version = detectGoVersion(otherDir)
	if version != "1.19" {
		t.Errorf("Expected Go version 1.19, got %s", version)
	}

	// Test with go.mod that has a patch version
	patchDir := setupTestDir(t)
	createFile(t, filepath.Join(patchDir, "go.mod"), `module example.com/patch
go 1.20.5
`)

	version = detectGoVersion(patchDir)
	if version != "1.20" {
		t.Errorf("Expected Go version 1.20 (ignoring patch), got %s", version)
	}

	// Test default version when go.mod doesn't exist
	emptyDir := setupTestDir(t)
	version = detectGoVersion(emptyDir)
	if version != "1.22" {
		t.Errorf("Expected default Go version 1.22, got %s", version)
	}

	// Test when project is detected through DetectProject
	projectDir := setupTestDir(t)
	createFile(t, filepath.Join(projectDir, "go.mod"), `module example.com/test
go 1.18.3
`)
	createFile(t, filepath.Join(projectDir, "main.go"), `package main
func main() {
	println("Hello world")
}`)

	project, err := DetectProject(projectDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if project.Type != Go {
		t.Errorf("Expected Go project, got %s", project.Type)
	}

	if project.Version != "1.18" {
		t.Errorf("Expected Go version 1.18, got %s", project.Version)
	}
}
