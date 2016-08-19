package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// AgentHost is the remote agent hostname or IP address.
var AgentHost string

// AgentPort is the remote agent port number.
var AgentPort int

// AgentTimeout is the timeout used for calls to the remote agent.
var AgentTimeout int

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "inservice-cli",
	Short: "Provides CLI access to the InService Agent and Plugins.",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.inservice-cli.yaml)")

	RootCmd.PersistentFlags().StringVarP(&AgentHost, "host", "H", "127.0.0.1", "InService Agent Host Address")
	RootCmd.PersistentFlags().IntVarP(&AgentPort, "port", "P", 8080, "InService Agent Host Port")
	RootCmd.PersistentFlags().IntVarP(&AgentTimeout, "timeout", "t", 5, "InService Agent Request Timeout (Seconds)")
}

func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".inservice-cli") // name of config file (without extension)
	viper.AddConfigPath("$HOME")          // adding home directory as first search path
	viper.AutomaticEnv()                  // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func loadViperConfig() error {

	viper.SetConfigType("json")
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/inservice", AgentHost, AgentPort+1))
	if err != nil {
		return fmt.Errorf("Unable to get Config: %s", err)
	}
	defer resp.Body.Close()

	viper.ReadConfig(resp.Body)
	return nil
}
