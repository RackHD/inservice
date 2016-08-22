package main

import (
	"log"

	"github.com/RackHD/inservice/agent"
)

// PluginServer provides PluginHost configuration and management.
type PluginServer struct {
	plugins []string

	address string
	port    int

	host *plugins.PluginHost
}

// NewPluginServer returns a new PluginServer instance.
func NewPluginServer(address string, port int, plugins []string) (*PluginServer, error) {
	return &PluginServer{
		address: address,
		port:    port,
		plugins: plugins,
	}, nil
}

// Serve runs the PluginServer in blocking mode.
func (s *PluginServer) Serve() {
	var err error

	s.host, err = plugins.NewPluginHost(s.address, s.port, s.plugins)
	if err != nil {
		log.Fatalf("Error creating Plugin Host: %s\n", err)
	}

	s.host.Serve()
}
