package neighbors

import (
	"fmt"
	"github.com/RackHD/InService/plugins/lldp/grpc/lldp"
	"github.com/google/gopacket"
	"golang.org/x/net/context"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// Neighbors stores the local neighbors topology
type Neighbors struct {
	Rw         sync.RWMutex                          `json:"-"`
	NetDevList []net.Interface                       `json:"-"`
	Interfaces map[string]map[string]NeighborDetails `json:"connections,omitempty"`
	ChassisID  string                                `json:"host,omitempty"`
	SysName    string                                `json:"hostname,omitempty"`
	SysDesc    string                                `json:"description,omitempty"`
	Address    string                                `json:"address,omitempty"`
	Type       string                                `json:"type,omitempty"`
}

// NeighborDetails stores relevent information to identify a neighboring node
type NeighborDetails struct {
	PortID          string      `json:"portid,omitempty"`
	PortDescription string      `json:"portdesc,omitempty"`
	SysName         string      `json:"hostname,omitempty"`
	SysDesc         string      `json:"description,omitempty"`
	Address         string      `json:"address,omitempty"`
	Vlan            string      `json:"vlan,omitempty"`
	Type            string      `json:"type,omitempty"`
	TTL             *time.Timer `json:"-"`
}

// Packet encodes the ingress interface with the receives LLDP packet
type Packet struct {
	Iface  net.Interface
	Packet gopacket.Packet
}

// NewNeighbors creates the Neighbors data storage object and returns it to the LLDP Plugin
func NewNeighbors() (*Neighbors, error) {
	t := make(map[string]map[string]NeighborDetails)
	return &Neighbors{
		Interfaces: t,
	}, nil

}

// ListInterfaces provides a list of all active interfaces
func (n *Neighbors) ListInterfaces(in *lldp.EmptyMessage, stream lldp.Lldp_ListInterfacesServer) error {
	for _, iface := range n.NetDevList {
		err := stream.Send(&lldp.Interface{
			Name:         iface.Name,
			Hardwareaddr: iface.HardwareAddr.String(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// GetInterfaceDetails provides details indofmration about each interface
func (n *Neighbors) GetInterfaceDetails(ctx context.Context, in *lldp.Interface) (*lldp.Interface, error) {
	for _, iface := range n.NetDevList {
		if iface.Name == in.Name {
			return &lldp.Interface{
				Index:        int32(iface.Index),
				Mtu:          int32(iface.MTU),
				Name:         iface.Name,
				Hardwareaddr: iface.HardwareAddr.String(),
				Flags:        uint32(iface.Flags),
			}, nil
		}
	}
	return &lldp.Interface{}, nil
}

// ListInterfaceNeighbors shows the neighbors attached to a specific interface
func (n *Neighbors) ListInterfaceNeighbors(in *lldp.Interface, stream lldp.Lldp_ListInterfaceNeighborsServer) error {
	var targetInterface *net.Interface
	for _, iface := range n.NetDevList {
		if in.Name == iface.Name {
			targetInterface = &iface
			break
		}
	}
	if targetInterface == nil {
		return fmt.Errorf("Interface does not exist or is not active")
	}

	n.Rw.RLock()
	defer n.Rw.RUnlock()
	val, ok := n.Interfaces[targetInterface.Name]
	if !ok {
		return fmt.Errorf("Interface does not have any neighbors")
	}
	for neighbors := range val {
		data := val[neighbors]
		err := stream.Send(&lldp.Neighbor{
			Portid:          data.PortID,
			Portdescription: data.PortDescription,
			Sysname:         data.SysName,
			Sysdesc:         data.SysDesc,
			Address:         data.Address,
			Vlan:            data.Vlan,
			Type:            data.Type,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// ListNeighbors shows all neighbors attached to all active interfaces on the host
func (n *Neighbors) ListNeighbors(in *lldp.EmptyMessage, stream lldp.Lldp_ListNeighborsServer) error {
	n.Rw.RLock()
	defer n.Rw.RUnlock()
	for _, val := range n.Interfaces {
		for neighbors := range val {
			data := val[neighbors]
			err := stream.Send(&lldp.Neighbor{
				Portid:          data.PortID,
				Portdescription: data.PortDescription,
				Sysname:         data.SysName,
				Sysdesc:         data.SysDesc,
				Address:         data.Address,
				Vlan:            data.Vlan,
				Type:            data.Type,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ProcessPacket takes a cpatured packet
func (n *Neighbors) ProcessPacket(wg *sync.WaitGroup, in Packet) {
	//	log.Printf("%v\n", in.Packet.Data())
	lldpDiscovery, lldpInfo, err := n.translatePacket(&in)
	if err != nil {
		wg.Done()
		return
	}
	chassis := n.decodeChassis(lldpDiscovery)
	address := n.decodeMgmtAddress(lldpInfo)
	neighborPort := n.decodePortID(lldpDiscovery)
	neighborPortDescription := n.decodePortDescription(lldpInfo)
	vlan := n.decodeVlan(lldpInfo)
	sysType := n.determineSysType(lldpDiscovery, lldpInfo)
	if n.isLocalInterface(chassis) {
		n.ChassisID = chassis
		n.SysName = lldpInfo.SysName
		n.SysDesc = lldpInfo.SysDescription
		n.Address = address
		n.Type = sysType
	} else {
		if chassis == "" {
			wg.Done()
			return
		}
		n.Rw.Lock()
		defer n.Rw.Unlock()
		c, ok := n.Interfaces[in.Iface.Name]
		if !ok {
			c = make(map[string]NeighborDetails)
			n.Interfaces[in.Iface.Name] = c
		}
		t, ok := c[chassis]
		if ok {
			t.TTL.Stop()
		}
		f := func() {
			n.expireNeighbor(in.Iface.Name, chassis)
		}
		timeout := time.AfterFunc(time.Duration(lldpDiscovery.TTL)*time.Second, f)
		node := NeighborDetails{
			PortID:          neighborPort,
			PortDescription: neighborPortDescription,
			SysName:         lldpInfo.SysName,
			SysDesc:         lldpInfo.SysDescription,
			Address:         address,
			Vlan:            vlan,
			Type:            sysType,
			TTL:             timeout,
		}
		c[chassis] = node
	}
	wg.Done()
	return
}

func (n *Neighbors) isLocalInterface(input string) bool {
	for _, iface := range n.NetDevList {
		if strings.Compare(iface.HardwareAddr.String(), input) == 0 {
			return true
		}
	}
	return false
}

func (n *Neighbors) expireNeighbor(iface string, chassis string) {
	n.Rw.Lock()
	defer n.Rw.Unlock()
	neighbor, ok := n.Interfaces[iface]
	if !ok {
		return
	}
	delete(neighbor, chassis)
	log.Printf("LLDP: Removed expired neighbor %s\n", chassis)
	return
}
