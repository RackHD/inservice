package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var binaryName, buildDate, buildUser, commitHash, goVersion, osArch, releaseVersion string

// versioinCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints metadata version information about this CLI tool.",
	Long:  `inservice-cli version`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(binaryName)
		log.Println("  Release version: " + releaseVersion)
		log.Println("  Built On: " + buildDate)
		log.Println("  Build By: " + buildUser)
		log.Println("  Commit Hash: " + commitHash)
		log.Println("  Go version: " + goVersion)
		log.Println("  OS/Arch: " + osArch)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
