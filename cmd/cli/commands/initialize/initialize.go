package initialize

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Babatunde50/dockergen/internal/detector"
	"github.com/Babatunde50/dockergen/internal/generator"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:  "init",
	Usage: "Initialize a Dockerfile for your project",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "compose",
			Aliases: []string{"c"},
			Usage:   "Generate docker-compose.yml",
		},
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "Overwrite existing files",
		},
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Usage:   "Specify app port (default: auto-detect)",
		},
		&cli.BoolFlag{
			Name:    "multi-stage",
			Aliases: []string{"m"},
			Usage:   "Use multi-stage build for Go projects",
			Value:   true,
		},
	},
	Action: func(cCtx *cli.Context) error {

		workDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %v", err)
		}

		project, err := detector.DetectProject(workDir)
		if err != nil {
			return fmt.Errorf("failed to detect project: %v", err)
		}

		if cCtx.IsSet("port") {
			project.Port = cCtx.Int("port")
		}

		dockerfilePath := filepath.Join(workDir, "Dockerfile")
		dockerComposeFilePath := filepath.Join(workDir, "docker-compose.yml")

		if !cCtx.Bool("force") {
			// Check if Dockerfile exists
			if _, err := os.Stat(dockerfilePath); err == nil {
				return fmt.Errorf("docker file already exists. use --force to overwrite")
			}

			// Check if docker-compose.yml exists (if generation requested)
			if cCtx.Bool("compose") {
				if _, err := os.Stat(dockerComposeFilePath); err == nil {
					return fmt.Errorf("docker-compose.yml already exists. Use --force to overwrite")
				}
			}
		}

		// Generate Dockerfile
		dockerfileContent, err := generator.GenerateDockerfile(project, cCtx.Bool("multi-stage"))

		if err != nil {
			return fmt.Errorf("failed to generate Dockerfile: %v", err)
		}

		// Write Dockerfile
		err = os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644)
		if err != nil {
			return fmt.Errorf("failed to write Dockerfile: %v", err)
		}
		fmt.Printf("✅ Generated Dockerfile for %s project\n", project.Type)

		// Generate docker-compose.yml
		if cCtx.Bool("compose") {
			dockerComposeContent, err := generator.GenerateDockerCompose(getProjectName(project), fmt.Sprintf("%d", project.Port))

			if err != nil {
				return fmt.Errorf("failed to write docker-compose.yml: %v", err)
			}

			err = os.WriteFile(dockerComposeFilePath, []byte(dockerComposeContent), 0644)
			if err != nil {
				return fmt.Errorf("failed to write docker-compose.yml: %v", err)
			}

			fmt.Printf("✅ Generated docker-compose.yml for %s project\n", project.Type)
		}

		fmt.Println("🚀 Dockerization complete!")
		return nil
	},
}

func getProjectName(project *detector.Project) string {

	baseName := filepath.Base(project.WorkDir)

	baseName = strings.ToLower(baseName)
	baseName = strings.ReplaceAll(baseName, " ", "-")
	baseName = strings.ReplaceAll(baseName, "_", "-")

	return baseName
}
