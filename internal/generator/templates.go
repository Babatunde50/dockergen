package generator

// Go Dockerfile template
const goDockerfileTemplate = `# syntax=docker/dockerfile:1

{{if .UseMultiStage}}
# === Multi-stage build ===

# Build stage
FROM golang:{{.Version}}-alpine AS build
WORKDIR /app

# Copy go.mod and go.sum files first and download dependencies
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations for smaller binary size
RUN {{.BuildCmd}}

# Runtime stage with a minimal Alpine image
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata

# Create a non-root user to run the application
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Create app directory and set permissions
RUN mkdir -p /app && chown -R appuser:appgroup /app

# Set the working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build {{.Entrypoint}} /app/

# Switch to non-root user for security
USER appuser

{{if .Port}}
# Expose the application port
EXPOSE {{.Port}}
{{end}}

# Run the application
ENTRYPOINT ["/app/{{.BinaryName}}"]

{{else}}
# === Single-stage build ===

FROM golang:{{.Version}}-alpine
WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN {{.BuildCmd}}

{{if .Port}}
EXPOSE {{.Port}}
{{end}}

ENTRYPOINT ["{{.RunCmd}}"]

{{end}}`
