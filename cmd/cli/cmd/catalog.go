package cmd

import (
	"github.com/spf13/cobra"
)

// catalogCmd represents the catalog command
var catalogCmd = &cobra.Command{
	Use:   "catalog",
	Short: "Catalog Compute plugin control commands",
}

func init() {
	RootCmd.AddCommand(catalogCmd)
}
