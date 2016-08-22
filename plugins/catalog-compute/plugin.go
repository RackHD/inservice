package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"log"
	"os/exec"

	"golang.org/x/net/context"
)

// CatalogComputePlugin implements the Plugin interface.
type CatalogComputePlugin struct {
	address string
	port    int

	laddr    *net.TCPAddr
	listener *net.TCPListener
}

// NewCatalogComputePlugin initializes a new CatalogComputePlugin struct.
func NewCatalogComputePlugin(address string, port int) (*CatalogComputePlugin, error) {
	var err error

	plugin := &CatalogComputePlugin{
		address: address,
		port:    port,
	}

	plugin.laddr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, err
	}

	return plugin, nil
}

// Start handles plugin startup. This creates goroutines to handle the packet capture
// and serve the gRPC API or REST API
func (p *CatalogComputePlugin) Start(ctx context.Context) error {
        log.Println("Catalog-Compute Plugin Started.")
	if p.listener != nil {
		return fmt.Errorf("Plugin already started.")
	}

	var err error

	p.listener, err = net.ListenTCP("tcp", p.laddr)
	if err != nil {
		return err
	}

	http.HandleFunc("/lshw", p.HandleLshw)
	http.Serve(p.listener, nil)

	return nil
}

// Stop closes a channel that should stop all capture
func (p *CatalogComputePlugin) Stop(ctx context.Context) error {
	if p.listener == nil {
		return fmt.Errorf("Plugin not started.")
	}

	err := p.listener.Close()

	p.listener = nil

	return err
}

// Status is a dummy function for now until a better status mechanism is identified
func (p *CatalogComputePlugin) Status(ctx context.Context) error {
	if p.listener == nil {
		return fmt.Errorf("Plugin not started.")
	}

	return nil
}

// HandleLshw replies to HTTP requests for the lshw JSON catalog.
func (p *CatalogComputePlugin) HandleLshw(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	out, err := exec.Command("lshw", "-json").Output()
	if err != nil {
		io.WriteString(w, err.Error())
		w.WriteHeader(500)
	} else {
		io.WriteString(w, string(out))
	}
}
