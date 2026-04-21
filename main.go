package main

import (
	"embed"
	"went/app"
	"went/cmd"
)

var Version = ""
var Channel = ""
var BuildDate = ""

//go:embed embed/*
var EmbedFS embed.FS

func main() {
	buildSpecs := app.NewBuildSpecs(Version, Channel, BuildDate)

	rootCmd := cmd.NewRootCmd(
		"went",
		"The Go scaffolding CLI",
		"A CLI tool to scaffold Go projects with best practices.",
		buildSpecs,
	)

	initCmd := cmd.NewInitCmd(EmbedFS, buildSpecs).Command()

	rootCmd.AddCommand(initCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
