package cmd

import "github.com/spf13/cobra"

var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "A collection of commands related to Plugin management.",
}

func init() {
	RootCmd.AddCommand(pluginsCmd)
}
