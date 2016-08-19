package main

import (
	"fmt"
	"log"
	"time"

	"github.com/king-jam/gossdp"
)

type ssdpHandler struct{}

// Response is the callback to process inbound SSDP messages.
func (h *ssdpHandler) Response(message gossdp.ResponseMessage) {
	fmt.Printf("%+v\n", message)
}

func main() {
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

	time.Sleep(60 * time.Second)
}
