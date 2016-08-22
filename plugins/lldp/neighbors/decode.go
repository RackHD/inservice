package neighbors

import (
	"errors"
	"github.com/google/gopacket/layers"
	"net"
	"strconv"
	"strings"
)

func (s *Neighbors) translatePacket(input *Packet) (*layers.LinkLayerDiscovery, *layers.LinkLayerDiscoveryInfo, error) {
	d := input.Packet.Layer(layers.LayerTypeLinkLayerDiscovery)
	if d == nil {
		return &layers.LinkLayerDiscovery{}, &layers.LinkLayerDiscoveryInfo{}, errors.New("LLDP: Failed to translate packet")
	}
	i := input.Packet.Layer(layers.LayerTypeLinkLayerDiscoveryInfo)
	if i == nil {
		return &layers.LinkLayerDiscovery{}, &layers.LinkLayerDiscoveryInfo{}, errors.New("LLDP: Failed to translate packet")
	}
	return d.(*layers.LinkLayerDiscovery), i.(*layers.LinkLayerDiscoveryInfo), nil
}

func (s *Neighbors) decodeChassis(discovery *layers.LinkLayerDiscovery) string {
	subType := discovery.ChassisID.Subtype
	switch subType {
	case 4:
		return net.HardwareAddr(discovery.ChassisID.ID).String()
	default:
		return ""
	}
}

func (s *Neighbors) decodeMgmtAddress(info *layers.LinkLayerDiscoveryInfo) string {
	subType := info.MgmtAddress.Subtype
	switch subType {
	case 1:
		return net.IP(info.MgmtAddress.Address).String()
	default:
		return ""
	}
}

func (s *Neighbors) decodePortID(discovery *layers.LinkLayerDiscovery) string {
	subType := discovery.PortID.Subtype
	switch subType {
	case 1, 2, 5, 7:
		return string(discovery.PortID.ID)
	case 3:
		return net.HardwareAddr(discovery.PortID.ID).String()
	case 4:
		return net.IP(discovery.PortID.ID).String()
	default:
		return ""
	}
}

func (s *Neighbors) decodePortDescription(info *layers.LinkLayerDiscoveryInfo) string {
	return info.PortDescription
}

func (s *Neighbors) decodeVlan(info *layers.LinkLayerDiscoveryInfo) string {
	info8021, err := info.Decode8021()
	if err != nil {
		return ""
	}
	return strconv.Itoa(int(info8021.PVID))
}

func (s *Neighbors) determineSysType(discovery *layers.LinkLayerDiscovery, info *layers.LinkLayerDiscoveryInfo) string {
	if s.isVMwareOui(s.decodeChassis(discovery)) {
		return "vm"
	}
	if discovery.PortID.Subtype == 7 {
		return "nic"
	}
	if info.SysCapabilities.EnabledCap.Router {
		return "router"
	}
	if info.SysCapabilities.EnabledCap.Bridge {
		return "switch"
	}
	return "unknown"
}

var vmwareoui = [...]string{"00:1C:14", "00:0C:29", "00:50:56", "00:05:69"}

func (s *Neighbors) isVMwareOui(mac string) bool {
	var oui string
	if len(mac) > 0 {
		oui = strings.ToUpper(mac[:8])
	} else {
		return false
	}
	for _, addr := range vmwareoui {
		if oui == addr {
			return true
		}
	}
	return false
}
