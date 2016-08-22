package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

// lshwCmd represents the lshw command
var lshwCmd = &cobra.Command{
	Use:   "lshw",
	Short: "Retrieves the HW catalog from the Catalog Compute plugin.",
	Long:  `inservice-cli catalog lshw`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := loadViperConfig()
		if err != nil {
			return fmt.Errorf("Could not load viper config: %s", err)
		}
		catalogPort := viper.GetInt("plugins.inservice-catalog-compute.port")

		resp, err := http.Get(fmt.Sprintf("http://%s:%d/lshw", AgentHost, catalogPort))

		if err != nil {
			return fmt.Errorf("Unable to get HW Catalog: %s", err)
		}
		defer resp.Body.Close()

		hardware, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Unable to read HW Catalog: %s", err)
		}

		fmt.Printf("%s", hardware)
		return err
	},
}

func init() {
	catalogCmd.AddCommand(lshwCmd)
}
