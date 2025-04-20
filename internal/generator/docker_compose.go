package generator

import (
	"bytes"
	"fmt"
	"strings"
)

// GenerateDockerCompose creates a default docker-compose.yml file
func GenerateDockerCompose(projectName, port string) (string, error) {

	if projectName == "" || port == "" {
		return "", fmt.Errorf("project name and port are required")
	}

	composeTemplate := DockerComposeTemplate{
		Version: "3.8",
		Name:    projectName,
		Services: []Service{
			{
				Name:          "app",
				ContainerName: fmt.Sprintf("%s-app", projectName),
				Restart:       "unless-stopped",
				Build: Build{
					Context:    ".",
					Dockerfile: "Dockerfile",
				},
				Environment: map[string]string{
					"ENV": "development",
				},
			},
		},
	}

	return renderDockerCompose(composeTemplate)
}

// renderDockerCompose converts the template to YAML content
func renderDockerCompose(template DockerComposeTemplate) (string, error) {
	var buf bytes.Buffer

	// Helper to write quoted strings if necessary
	quote := func(s string) string {
		if strings.ContainsAny(s, ": \t\"'") {
			return fmt.Sprintf("'%s'", strings.ReplaceAll(s, "'", "''"))
		}
		return s
	}

	// Write version
	buf.WriteString(fmt.Sprintf("version: '%s'\n\n", template.Version))

	// Write name
	buf.WriteString(fmt.Sprintf("name: '%s'\n\n", template.Name))

	// Write services
	if len(template.Services) > 0 {
		buf.WriteString("services:\n")

		for _, service := range template.Services {
			buf.WriteString(fmt.Sprintf("  %s:\n", service.Name))

			// Container name
			if service.ContainerName != "" {
				buf.WriteString(fmt.Sprintf("    container_name: %s\n", quote(service.ContainerName)))
			}

			// Image
			if service.Image != "" {
				buf.WriteString(fmt.Sprintf("    image: %s\n", quote(service.Image)))
			}

			// Build
			if service.Build.Context != "" {
				buf.WriteString("    build:\n")
				buf.WriteString(fmt.Sprintf("      context: %s\n", quote(service.Build.Context)))
				if service.Build.Dockerfile != "" {
					buf.WriteString(fmt.Sprintf("      dockerfile: %s\n", quote(service.Build.Dockerfile)))
				}
				if len(service.Build.Args) > 0 {
					buf.WriteString("      args:\n")
					for key, value := range service.Build.Args {
						buf.WriteString(fmt.Sprintf("        %s: %s\n", quote(key), quote(value)))
					}
				}
				if service.Build.Target != "" {
					buf.WriteString(fmt.Sprintf("      target: %s\n", quote(service.Build.Target)))
				}
				if len(service.Build.CacheFrom) > 0 {
					buf.WriteString("      cache_from:\n")
					for _, cacheFrom := range service.Build.CacheFrom {
						buf.WriteString(fmt.Sprintf("        - %s\n", quote(cacheFrom)))
					}
				}
			}

			// Restart policy
			if service.Restart != "" {
				buf.WriteString(fmt.Sprintf("    restart: %s\n", quote(service.Restart)))
			}

			// Ports
			if len(service.Ports) > 0 {
				buf.WriteString("    ports:\n")
				for _, port := range service.Ports {
					buf.WriteString(fmt.Sprintf("      - %s\n", quote(port)))
				}
			}

			// Environment variables (as map)
			if len(service.Environment) > 0 {
				buf.WriteString("    environment:\n")
				for key, value := range service.Environment {
					buf.WriteString(fmt.Sprintf("      %s: %s\n", quote(key), quote(value)))
				}
			}

			// Env files
			if len(service.EnvFile) > 0 {
				buf.WriteString("    env_file:\n")
				for _, envFile := range service.EnvFile {
					buf.WriteString(fmt.Sprintf("      - %s\n", quote(envFile)))
				}
			}

			// Volumes
			if len(service.Volumes) > 0 {
				buf.WriteString("    volumes:\n")
				for _, volume := range service.Volumes {
					buf.WriteString(fmt.Sprintf("      - %s\n", quote(volume)))
				}
			}

			// Networks
			if len(service.Networks) > 0 {
				buf.WriteString("    networks:\n")
				for _, network := range service.Networks {
					buf.WriteString(fmt.Sprintf("      - %s\n", quote(network)))
				}
			}

			// Depends on
			if len(service.DependsOn) > 0 {
				buf.WriteString("    depends_on:\n")
				for _, dependency := range service.DependsOn {
					buf.WriteString(fmt.Sprintf("      - %s\n", quote(dependency)))
				}
			}

			// Health Check
			if service.HealthCheck.Test != nil {
				buf.WriteString("    healthcheck:\n")
				buf.WriteString("      test:\n")
				for _, testCmd := range service.HealthCheck.Test {
					buf.WriteString(fmt.Sprintf("        - %s\n", quote(testCmd)))
				}
				if service.HealthCheck.Interval != "" {
					buf.WriteString(fmt.Sprintf("      interval: %s\n", quote(service.HealthCheck.Interval)))
				}
				if service.HealthCheck.Timeout != "" {
					buf.WriteString(fmt.Sprintf("      timeout: %s\n", quote(service.HealthCheck.Timeout)))
				}
				if service.HealthCheck.Retries != 0 {
					buf.WriteString(fmt.Sprintf("      retries: %d\n", service.HealthCheck.Retries))
				}
				if service.HealthCheck.StartPeriod != "" {
					buf.WriteString(fmt.Sprintf("      start_period: %s\n", quote(service.HealthCheck.StartPeriod)))
				}
			}

			// Deploy configuration
			if service.Deploy.Mode != "" || service.Deploy.Replicas > 0 {
				buf.WriteString("    deploy:\n")
				if service.Deploy.Mode != "" {
					buf.WriteString(fmt.Sprintf("      mode: %s\n", quote(service.Deploy.Mode)))
				}
				if service.Deploy.Replicas > 0 {
					buf.WriteString(fmt.Sprintf("      replicas: %d\n", service.Deploy.Replicas))
				}
				// Add more deploy fields (resources, update_config, etc.) as needed
			}

			// Labels
			if len(service.Labels) > 0 {
				buf.WriteString("    labels:\n")
				for key, value := range service.Labels {
					buf.WriteString(fmt.Sprintf("      %s: %s\n", quote(key), quote(value)))
				}
			}

			// Command
			if service.Command != "" {
				buf.WriteString(fmt.Sprintf("    command: %s\n", quote(service.Command)))
			}

			// Entrypoint
			if service.Entrypoint != "" {
				buf.WriteString(fmt.Sprintf("    entrypoint: %s\n", quote(service.Entrypoint)))
			}

			// User
			if service.User != "" {
				buf.WriteString(fmt.Sprintf("    user: %s\n", quote(service.User)))
			}

			// Working directory
			if service.WorkingDir != "" {
				buf.WriteString(fmt.Sprintf("    working_dir: %s\n", quote(service.WorkingDir)))
			}

			// Read-only
			if service.ReadOnly {
				buf.WriteString("    read_only: true\n")
			}

			buf.WriteString("\n")
		}
	}

	// Write networks
	if len(template.Networks) > 0 {
		buf.WriteString("networks:\n")
		for _, network := range template.Networks {
			buf.WriteString(fmt.Sprintf("  %s:\n", network.Name))
			if network.Driver != "" {
				buf.WriteString(fmt.Sprintf("    driver: %s\n", quote(network.Driver)))
			}
			if network.External {
				buf.WriteString("    external: true\n")
			}
			if network.Attachable {
				buf.WriteString("    attachable: true\n")
			}
			if len(network.Labels) > 0 {
				buf.WriteString("    labels:\n")
				for key, value := range network.Labels {
					buf.WriteString(fmt.Sprintf("      %s: %s\n", quote(key), quote(value)))
				}
			}
			// IPAM configuration can be added here if needed
			buf.WriteString("\n")
		}
	}

	// Write volumes
	if len(template.Volumes) > 0 {
		buf.WriteString("volumes:\n")
		for _, volume := range template.Volumes {
			buf.WriteString(fmt.Sprintf("  %s:\n", volume.Name))
			if volume.Driver != "" {
				buf.WriteString(fmt.Sprintf("    driver: %s\n", quote(volume.Driver)))
			}
			if volume.External {
				buf.WriteString("    external: true\n")
			}
			if len(volume.Labels) > 0 {
				buf.WriteString("    labels:\n")
				for key, value := range volume.Labels {
					buf.WriteString(fmt.Sprintf("      %s: %s\n", quote(key), quote(value)))
				}
			}
			if len(volume.Options) > 0 {
				buf.WriteString("    options:\n")
				for key, value := range volume.Options {
					buf.WriteString(fmt.Sprintf("      %s: %s\n", quote(key), quote(value)))
				}
			}
			buf.WriteString("\n")
		}
	}

	// Write configs
	if len(template.Configs) > 0 {
		buf.WriteString("configs:\n")
		for _, config := range template.Configs {
			buf.WriteString(fmt.Sprintf("  %s:\n", config.Name))
			if config.External {
				buf.WriteString("    external: true\n")
			} else {
				buf.WriteString(fmt.Sprintf("    file: %s\n", quote(config.File)))
			}
			buf.WriteString("\n")
		}
	}

	// Write secrets
	if len(template.Secrets) > 0 {
		buf.WriteString("secrets:\n")
		for _, secret := range template.Secrets {
			buf.WriteString(fmt.Sprintf("  %s:\n", secret.Name))
			if secret.External {
				buf.WriteString("    external: true\n")
			} else {
				buf.WriteString(fmt.Sprintf("    file: %s\n", quote(secret.File)))
			}
			buf.WriteString("\n")
		}
	}

	return buf.String(), nil
}
