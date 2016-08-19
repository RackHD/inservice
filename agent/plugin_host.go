package plugins

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	"github.com/RackHD/InService/agent/grpc/host"
)

// PluginResponseTimeout sets the maximum duration allowed to respond to a request
// for a Plugin Service action (start/stop/status).
const PluginResponseTimeout time.Duration = 10 * time.Second

// PluginHost encapsulates a collection of PluginProcesses and manages their
// lifecycles.
type PluginHost struct {
	address string
	port    int

	processes     []*PluginProcess
	notifications chan *PluginProcess
	plugins       []string
}

// NewPluginHost returns a new PluginHost.
func NewPluginHost(address string, port int, plugins []string) (*PluginHost, error) {
	return &PluginHost{
		address: address,
		port:    port,
		plugins: plugins,
	}, nil
}

// Serve loads and manages all Plugins.
func (h *PluginHost) Serve() {
	go h.serveGRPC()
	h.servePlugins()
}

func (h *PluginHost) servePlugins() {
	for _, plugin := range h.plugins {
		h.processes = append(
			h.processes,
			NewPluginProcess(
				plugin,
			),
		)
	}

	h.notifications = make(chan *PluginProcess, len(h.plugins))

	for _, plugin := range h.processes {
		err := plugin.Start(h.notifications)
		if err != nil {
			log.Printf("Error Starting Plugin: %s, %+v\n", err, plugin)
		}
	}

	// Monitor
	for plugin := range h.notifications {
		if plugin.Restartable() {
			err := plugin.Start(h.notifications)
			if err != nil {
				log.Printf("Plugin Unable to Restart: %s (%s)\n", plugin.Name, err)
			}

			log.Printf("Plugin Restarted: %s\n", plugin.Name)
		} else {
			log.Printf("Plugin Not Set for Restart: %s\n", plugin.Name)
		}
	}
}

func (h *PluginHost) serveGRPC() {
	listen, err := net.Listen(
		"tcp",
		fmt.Sprintf("%s:%d", h.address, h.port),
	)
	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}

	server := grpc.NewServer()

	host.RegisterHostServer(server, h)

	server.Serve(listen)
}

// Start implements the gRPC plugin host API and starts the requested plugin.
func (h *PluginHost) Start(ctx context.Context, in *host.StartRequest) (*host.StartResponse, error) {
	for _, plugin := range h.processes {
		if plugin.Name == in.Name {
			err := plugin.Start(h.notifications)
			if err != nil {
				return nil, err
			}

			return &host.StartResponse{}, nil
		}
	}

	return nil, fmt.Errorf("Unable to locate plugin %s.", in.Name)
}

// Stop implements the gRPC plugin host API and stops the requested plugin.
func (h *PluginHost) Stop(ctx context.Context, in *host.StopRequest) (*host.StopResponse, error) {
	for _, plugin := range h.processes {
		if plugin.Name == in.Name {
			err := plugin.Stop()
			if err != nil {
				return nil, err
			}

			return &host.StopResponse{}, nil
		}
	}

	return nil, fmt.Errorf("Unable to locate plugin %s.", in.Name)
}

// Status implements the gRPC plugin host API and gets status on all plugins.
func (h *PluginHost) Status(in *host.StatusRequest, stream host.Host_StatusServer) error {
	for _, plugin := range h.processes {
		var state = "Not Running"

		if plugin.Status() {
			state = "Running"
		}

		stream.Send(&host.StatusResponse{
			Name:   plugin.Name,
			Status: state,
		})
	}

	stream.Send(nil)

	return nil
}
