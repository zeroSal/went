package main

import (
	"embed"
	"fmt"
	"os"
	"went/app"
	_ "went/cmd"
	"went/registry"

	"github.com/zeroSal/went-clio/clio"
	"github.com/zeroSal/went-command/command"

	"github.com/spf13/cobra"
)

var Version = ""
var Channel = ""
var BuildDate = ""

//go:embed res/*
var EmbedFS embed.FS

func main() {
	clio := clio.NewClio()

	data, err := EmbedFS.ReadFile("res/banner.template")
	if err != nil {
		clio.Error("Error loading the banner template.")
		os.Exit(3)
	}

	specs := app.NewSpecs(Version, Channel, BuildDate)
	clio.SetBanner(string(data), Version, Channel, BuildDate)

	kernel := app.NewKernel(EmbedFS, specs, clio)

	root := &cobra.Command{
		Version: fmt.Sprintf("%s-%s (%s)", Version, Channel, BuildDate),
		Use:     "went",
		Short:   "Go project wireframe following best practices",
		Long:    "Go project wireframe following best practices.",
	}

	run := func(command command.Interface) {
		if err := kernel.Run(command.Invoke()); err != nil {
			clio.Fatal("%s", err.Error())
			os.Exit(1)
		}
	}

	if err := command.Mount(registry.Command.All(), root, run).Execute(); err != nil {
		clio.Fatal("Error mounting commands: %s", err.Error())
		os.Exit(2)
	}
}
