package cmd

import (
	"fmt"
	"github.com/RackHD/InService/plugins/lldp/grpc/lldp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"time"
)

// lldpInterfaceListCmd represents the list command
var lldpInterfaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "Information about LLDP interfaces",
	Long: `Displays list of interfaces in use by LLDP Plugin:
		Usage:
			Interface List - inservice-cli lldp interface list`,
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

		stream, err := client.ListInterfaces(context.Background(), &lldp.EmptyMessage{})
		if err != nil {
			return fmt.Errorf("Unable to list interfaces: %s", err)
		}

		for {
			item, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return fmt.Errorf("Unable to stream interfaces: %s", err)
			}
			fmt.Printf("%+v\n", item)
		}
	},
}

func init() {
	interfaceCmd.AddCommand(lldpInterfaceListCmd)
}
