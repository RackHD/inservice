package cmd

import (
	"github.com/spf13/cobra"
)

// lldpCmd represents the lldp command
var lldpCmd = &cobra.Command{
	Use:   "lldp",
	Short: "lldp plugin control commands",
}

func init() {
	RootCmd.AddCommand(lldpCmd)
}
