package cmd

import (
	"fmt"
	"time"

	"github.com/RackHD/inservice/agent/grpc/host"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the specified Plugin.",
	Long:  `inservice-cli plugins start <plugin>`,
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

		_, err = client.Start(
			context.Background(),
			&host.StartRequest{
				Name: args[0],
			},
		)

		return err
	},
}

func init() {
	pluginsCmd.AddCommand(startCmd)
}
