package plugins

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/RackHD/InService/agent/grpc/plugin"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// This is a helper to generate a suitable unix socket dialer for gRPC.
func unixDialer(path string, timeout time.Duration) (net.Conn, error) {
	laddr, err := net.ResolveUnixAddr("unix", "")
	if err != nil {
		return nil, err
	}

	raddr, err := net.ResolveUnixAddr("unix", path)
	if err != nil {
		return nil, err
	}

	return net.DialUnix("unix", laddr, raddr)
}

// PluginProcess encapsulates plugin process lifecycle as
// managed by the Plugin Host.
type PluginProcess struct {
	Name    string
	restart bool
	running bool
	cmd     *exec.Cmd
}

// NewPluginProcess returns a new PluginProcess.
func NewPluginProcess(name string) *PluginProcess {
	return &PluginProcess{name, true, false, nil}
}

// Restartable returns a bool indicating whether the PluginProcess should be
// restarted by the host if the process ends.  It's assumed to restart unless
// an API call was made to the Plugin Host to set the PluginProcess to stop.
func (p *PluginProcess) Restartable() bool {
	return p.restart
}

// Start brings up the PluginProcess providing a channel for the plugin to notify
// when the process ends.  It redirects stdout/stderr to the host console and
// also contacts the PluginProcess over gRPC to le it know it can start.
func (p *PluginProcess) Start(notify chan<- *PluginProcess) error {
	if p.running {
		return fmt.Errorf("Pluging Already Running: %s\n", p.Name)
	}

	path, err := exec.LookPath(p.Name)
	if err != nil {
		dir, errDir := filepath.Abs(filepath.Dir(os.Args[0]))
		if errDir != nil {
			return errDir
		}
		path = fmt.Sprintf("%s/%s", dir, p.Name)

	}

	p.cmd = exec.Command(path)

	p.cmd.Stdout = os.Stdout
	p.cmd.Stderr = os.Stderr

	err = p.cmd.Start()
	if err != nil {
		return err
	}

	p.restart = true
	p.running = true

	go func() {
		e := p.cmd.Wait()
		if e != nil {
			log.Printf("Plugin Exited with Error: %s (%s)\n", p.Name, e)
		} else {
			log.Printf("Plugin Exited: %s\n", p.Name)
		}

		p.running = false
		p.cmd = nil

		notify <- p
	}()

	conn, err := grpc.Dial(
		fmt.Sprintf(pluginSockFormat, p.Name),
		grpc.WithDialer(unixDialer),
		grpc.WithInsecure(),
		grpc.WithTimeout(PluginResponseTimeout),
	)
	if err != nil {
		log.Printf("Error Dialing Plugin: %s\n", err)
		return p.Stop()
	}
	defer conn.Close()

	client := plugin.NewPluginClient(conn)

	_, err = client.Start(context.Background(), &plugin.StartRequest{})
	if err != nil {
		log.Printf("Error Starting Plugin: %s\n", err)
		return p.Stop()
	}

	return nil
}

// Stop contacts a running PluginProcess over gRPC and asks it to shut down
// then it kills the PluginProcess if it hasn't shutdown in a suitable amount
// of time.
func (p *PluginProcess) Stop() error {
	p.restart = false

	conn, err := grpc.Dial(
		fmt.Sprintf(pluginSockFormat, p.Name),
		grpc.WithInsecure(),
		grpc.WithDialer(unixDialer),
		grpc.WithTimeout(PluginResponseTimeout),
	)

	if err != nil {
		return err
	}

	defer conn.Close()

	client := plugin.NewPluginClient(conn)

	_, err = client.Stop(context.Background(), &plugin.StopRequest{})
	if err != nil {
		return err
	}

	// TODO: Implement a wait mechanism so the PluginProcess can gracefully exit.

	return p.cmd.Process.Kill()
}

// Status returns a bool indicating whether the PluginProcess is running.
func (p *PluginProcess) Status() bool {
	return p.running
}
