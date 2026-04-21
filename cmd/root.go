package cmd

import (
	"went/app"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "went",
	Short: "went creates project scaffolding",
	Long:  `A CLI tool to scaffold Go projects with best practices.`,
}

func NewRootCmd(
	use string,
	short string,
	long string,
	buildSpecs *app.BuildSpecs,
) *cobra.Command {
	return &cobra.Command{
		Use:     use,
		Short:   short,
		Long:    long,
		Version: buildSpecs.GetVersion() + " (" + buildSpecs.GetChannel() + ")",
	}
}

func Execute() error {
	return rootCmd.Execute()
}

func AddCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
