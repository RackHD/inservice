package main

import (
	"log"

	"github.com/king-jam/gossdp"
)

// SSDPServer provides SSDP advertisment for the InService Agent and it's plugins.
type SSDPServer struct {
	ServiceType string
	DeviceUUID  string
	Location    string
	MaxAge      int

	ssdp *gossdp.Ssdp
}

// Serve runs the SSDPServer in blocking mode.
func (s *SSDPServer) Serve() {
	var err error

	s.ssdp, err = gossdp.NewSsdp(nil)
	if err != nil {
		log.Fatalf("Error creating SSDP Server: %s\n", err)
	}

	defer s.ssdp.Stop()

	definition := gossdp.AdvertisableServer{
		ServiceType: s.ServiceType,
		DeviceUuid:  s.DeviceUUID,
		Location:    s.Location,
		MaxAge:      s.MaxAge,
	}

	s.ssdp.AdvertiseServer(definition)
	s.ssdp.Start()
}
