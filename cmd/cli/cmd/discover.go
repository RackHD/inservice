// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"

	"github.com/king-jam/gossdp"
	"github.com/spf13/cobra"
)

type ssdpHandler struct{}

// Response is the callback to process inbound SSDP messages.
func (h *ssdpHandler) Response(message gossdp.ResponseMessage) {
	fmt.Printf("%+v\n", message)
}

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		handler := ssdpHandler{}

		client, err := gossdp.NewSsdpClient(&handler)
		if err != nil {
			log.Println("Failed to start client: ", err)
			return
		}
		// call stop  when we are done
		defer client.Stop()
		// run! this will block until stop is called. so open it in a goroutine here
		go client.Start()
		// send a request for the server type we are listening for.
		err = client.ListenFor("urn:skunkworxs:inservice:agent:0")
		if err != nil {
			log.Println("Error ", err)
		}

		fmt.Println("Press a key to exit discovery...")

		var input string

		fmt.Scanln(&input)
	},
}

func init() {
	RootCmd.AddCommand(discoverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// discoverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// discoverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
