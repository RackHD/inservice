package cmd

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/RackHD/InService/agent/grpc/host"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the specified Plugin.",
	Long:  `inservice-cli plugins stop <plugin>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("Plugin Name Required")
		}

		conn, err := grpc.Dial(
			fmt.Sprintf("%s:%d", AgentHost, AgentPort),
			grpc.WithInsecure(),
			grpc.WithTimeout(time.Duration(AgentTimeout)*time.Second),
		)
		if err != nil {
			return err
		}
		defer conn.Close()

		client := host.NewHostClient(conn)

		_, err = client.Stop(
			context.Background(),
			&host.StopRequest{
				Name: args[0],
			},
		)

		return err
	},
}

func init() {
	pluginsCmd.AddCommand(stopCmd)
}
