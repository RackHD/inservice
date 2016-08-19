package cmd

import (
	"fmt"
	"github.com/RackHD/InService/plugins/lldp/grpc/lldp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

// interfaceCmd represents the interface command
var interfaceCmd = &cobra.Command{
	Use:   "interface",
	Short: "Information about LLDP interfaces",
	Long: `Displays list of interfaces in use by LLDP Plugin or details about a
	specific interface:

	Usage:
		Interface Details - inservice-cli lldp interface <name>
		Interface List - inservice-cli lldp interface list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := loadViperConfig()
		if err != nil {
			return fmt.Errorf("Could not load viper config: %s", err)
		}
		LLDPPort := viper.GetInt("plugins.inservice-lldp.port")

		if len(args) != 1 {
			return fmt.Errorf("Interface Name or [command] Required")
		}

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

		details, err := client.GetInterfaceDetails(
			context.Background(),
			&lldp.Interface{
				Name: args[0],
			},
		)
		if err != nil {
			return fmt.Errorf("Unable to show details: %s", err)
		}
		fmt.Printf("%+v\n", details)
		return nil
	},
}

func init() {
	lldpCmd.AddCommand(interfaceCmd)
}
