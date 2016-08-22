package main

import (
	"fmt"
	"net/http"
)

// HTTPServer serves the InService Agent configuration file as part of the SSDP
// advertisement process for consumers to identify the configuration of the Agent.
type HTTPServer struct {
	Address    string
	Port       int
	ConfigFile string
	URI        string
}

// Serve runs the HTTPServer in blocking mode.
func (s *HTTPServer) Serve() {
	http.HandleFunc(fmt.Sprintf("/%s", s.URI), s.handler)
	http.ListenAndServe(fmt.Sprintf("%s:%d", s.Address, s.Port), nil)
}

func (s *HTTPServer) handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, s.ConfigFile)
}
