package cmd

import (
	"fmt"
	"github.com/RackHD/inservice/plugins/lldp/grpc/lldp"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// lldpNeighborsListCmd represents the list command
var lldpNeighborsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show information on the LLDP neighbors",
	Long: `Displays list of neighbors seen by LLDP Plugin or those seen from specific interface:
	        Usage:
		        Neighbor List - inservice-cli lldp neighbors list
			Interface Neighbors - inservice-cli lldp neighbors <name>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := loadViperConfig()
		if err != nil {
			return fmt.Errorf("Could not load viper config: %s", err)
		}
		LLDPPort := viper.GetInt("plugins.inservice-lldp.port")

		conn, err := grpc.Dial(
			fmt.Sprintf("%s:%d", AgentHost, LLDPPort),
			grpc.WithInsecure(),
			grpc.WithTimeout(time.Duration(AgentTimeout)*time.Second),
		)
		if err != nil {
			return fmt.Errorf("Unable to connect to host: %s", err)
		}
		defer conn.Close()

		client := lldp.NewLldpClient(conn)

		stream, err := client.ListNeighbors(context.Background(), &lldp.EmptyMessage{})
		if err != nil {
			return fmt.Errorf("Unable to list neighbors: %s", err)
		}

		for {
			item, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return fmt.Errorf("Unable to stream neighbors: %s", err)
			}
			fmt.Printf("%+v\n", item)
		}
	},
}

func init() {
	neighborsCmd.AddCommand(lldpNeighborsListCmd)
}
