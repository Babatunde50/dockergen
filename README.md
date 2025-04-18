# DockerGen

> Simple, smart Docker file generation for your projects

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Babatunde50/dockergen)](https://goreportcard.com/report/github.com/Babatunde50/dockergen)

## Overview

DockerGen is a CLI tool that automatically generates optimized Docker and docker-compose files for your projects. It intelligently analyzes your codebase to detect the project type, dependencies, and configuration, then creates tailored Docker files that follow best practices.

> **Note:** This project is still a work in progress.

## Features

- **Automatic Project Detection**: Automatically identifies Go, Node.js, and Python projects
- **Smart Configuration Detection**: Detects ports, entry points, and project structure
- **Version-Aware**: Uses the same language version defined in your project (e.g., Go version from go.mod)
- **Best Practices**: Generates optimized, secure multi-stage Dockerfiles
- **docker-compose Support**: Creates docker-compose.yml for local development

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/Babatunde50/dockergen.git
cd dockergen

# Build and install
make install
```

## Usage

Generate Docker files for your project:

```bash
# Basic usage - generates a Dockerfile
dockergen init

# Generate both Dockerfile and docker-compose.yml
dockergen init --compose

# Specify a custom port
dockergen init --port 8080

# Force overwrite existing files
dockergen init --force
```

## Examples

### Go Project

For a Go project with the following structure:

```
myapp/
├── cmd/
│   └── main.go
├── go.mod
└── go.sum
```

DockerGen will:
1. Detect the Go version from go.mod
2. Identify the main entry point
3. Generate an optimized multi-stage Dockerfile:

```dockerfile
# Build stage
FROM golang:1.21-alpine AS build
WORKDIR /app

# Copy go.mod and go.sum files first for better caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/myapp ./cmd/

# Runtime stage with a minimal Alpine image
FROM alpine:latest
...
```

## Project Status

This project is under active development. Currently supported:

- [x] Go project detection and Dockerfile generation
- [ ] Complete Node.js support
- [ ] Complete Python support
- [] docker-compose generation
- [ ] Container dependency detection

## Development

### Prerequisites

- Go 1.18+
- Make

### Setting Up Development Environment

```bash
# Clone the repository
git clone https://github.com/Babatunde50/dockergen.git
cd dockergen

# Run tests
make test

# Build the binary
make build

# Run the application
make run
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin feature/my-new-feature`
5. Submit a pull request 