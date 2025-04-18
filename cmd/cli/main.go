package main

import (
	"fmt"
	"os"

	"github.com/Babatunde50/dockergen/cmd/cli/app"
)

func main() {
	app := app.NewApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "dockergen: %s\n", err)
		os.Exit(1)
	}
}
