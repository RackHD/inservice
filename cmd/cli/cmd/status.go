package cmd

import (
	"fmt"
	"io"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/RackHD/InService/agent/grpc/host"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Lists the status of all Plugins.",
	Long:  `inservice-cli plugins status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := grpc.Dial(
			fmt.Sprintf("%s:%d", AgentHost, AgentPort),
			grpc.WithInsecure(),
			grpc.WithTimeout(time.Duration(AgentTimeout)*time.Second),
		)
		if err != nil {
			return fmt.Errorf("Unable to connect to host: %s", err)
		}
		defer conn.Close()

		client := host.NewHostClient(conn)

		stream, err := client.Status(context.Background(), &host.StatusRequest{})
		if err != nil {
			return fmt.Errorf("Unable to list status: %s", err)
		}

		for {
			item, err := stream.Recv()
			if err == io.EOF {
				return nil
			}

			if err != nil {
				return fmt.Errorf("Unable to stream status: %s", err)
			}

			fmt.Printf("%+v\n", item)
		}
	},
}

func init() {
	pluginsCmd.AddCommand(statusCmd)
}
