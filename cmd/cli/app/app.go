package app

import (
	"github.com/Babatunde50/dockergen/cmd/cli/commands/initialize"
	"github.com/urfave/cli/v2"
)

func init() {

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print the version",
	}
	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "Show help",
	}
}

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Dockergen"
	app.Usage = "A minimal CLI to auto-generate Dockerfiles and Docker Compose files based on project analysis. No container management – just smart scaffolding."

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Enable verbose/debug logging",
		},
		&cli.BoolFlag{
			Name:    "yes",
			Aliases: []string{"y"},
			Usage:   "Auto-confirm all prompts",
		},
		&cli.StringFlag{
			Name:  "log-file",
			Usage: "Write logs to `FILE`",
			Value: "",
		},
	}

	app.Commands = []*cli.Command{
		initialize.Command,
	}

	return app
}
