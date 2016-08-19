package plugins

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/RackHD/InService/agent/grpc/plugin"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

// This is a helper to normalize the plugin process name for use in generating
// a unix socket file name.
func executableName() string {
	_, file := filepath.Split(os.Args[0])
	return file
}

// pluginSockFormat provides a format string for the Plugin unix socket path.
const pluginSockFormat string = "/tmp/%s.sock"

// Service is an interface plugins need to implement in order
// to interact with the PluginHost.
type Service interface {
	Start(context.Context) error
	Stop(context.Context) error
	Status(context.Context) error
}

// Plugin is the main structure for a Plugin hosted in an executable.
type Plugin struct {
	Service Service
}

// NewPlugin returns a new Plugin object.
func NewPlugin(service Service) (*Plugin, error) {
	if service == nil {
		return nil, fmt.Errorf("Service Interface Required.")
	}

	return &Plugin{
		Service: service,
	}, nil
}

// Serve should be called by the Plugin executable and will block allowing the plugin
// to execute via the Service interface.
func (p *Plugin) Serve() {
	sock := fmt.Sprintf(pluginSockFormat, executableName())

	if err := os.Remove(sock); err != nil {
		log.Printf("Unable to remove sock %s: %s", sock, err)
	}

	listen, err := net.Listen("unix", sock)
	if err != nil {
		log.Fatalf("Failed to Listen: %v\n", err)
	}

	server := grpc.NewServer()

	plugin.RegisterPluginServer(server, p)

	server.Serve(listen)
}

// Start implements the gRPC Start call and proxies it to the Service interface.
func (p *Plugin) Start(ctx context.Context, in *plugin.StartRequest) (*plugin.StartResponse, error) {
	err := p.Service.Start(ctx)
	if err != nil {
		return nil, err
	}

	return &plugin.StartResponse{}, nil
}

// Stop implements the gRPC Stop call and proxies it to the Service interface.
func (p *Plugin) Stop(ctx context.Context, in *plugin.StopRequest) (*plugin.StopResponse, error) {
	err := p.Service.Stop(ctx)
	if err != nil {
		return nil, err
	}

	return &plugin.StopResponse{}, nil
}

// Status implements the gRCP Status call and proxies it to the Service interface.
func (p *Plugin) Status(ctx context.Context, in *plugin.StatusRequest) (*plugin.StatusResponse, error) {
	err := p.Service.Status(ctx)
	if err != nil {
		return nil, err
	}

	return &plugin.StatusResponse{}, nil
}
