package neighbors_test

import (
	. "github.com/RackHD/inservice/plugins/lldp/neighbors"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	"sync"
)

var _ = Describe("Neighbors", func() {

	var (
		neighbors *Neighbors
		data      []byte
		iface     net.Interface
		packet    Packet
		wg        *sync.WaitGroup
	)

	BeforeEach(func() {
		data = []byte{1, 128, 194, 0, 0, 14, 0, 28, 115, 104, 23, 47, 136, 204, 2, 7,
			4, 0, 28, 115, 104, 23, 35, 4, 11, 3, 69, 116, 104, 101, 114, 110, 101, 116,
			49, 50, 6, 2, 0, 2, 10, 11, 68, 67, 79, 79, 66, 45, 55, 48, 53, 48, 83, 12,
			78, 65, 114, 105, 115, 116, 97, 32, 78, 101, 116, 119, 111, 114, 107, 115, 32,
			69, 79, 83, 32, 118, 101, 114, 115, 105, 111, 110, 32, 52, 46, 49, 52, 46, 54,
			77, 32, 114, 117, 110, 110, 105, 110, 103, 32, 111, 110, 32, 97, 110, 32, 65,
			114, 105, 115, 116, 97, 32, 78, 101, 116, 119, 111, 114, 107, 115, 32, 68, 67,
			83, 45, 55, 48, 53, 48, 83, 45, 53, 50, 14, 4, 0, 20, 0, 20, 16, 12, 5, 1, 10,
			240, 16, 9, 2, 0, 30, 132, 129, 0, 254, 6, 0, 128, 194, 1, 0, 1, 254, 9, 0, 18,
			15, 3, 1, 0, 0, 0, 0, 254, 6, 0, 18, 15, 4, 36, 20, 0, 0}
		gopacket := gopacket.NewPacket(
			data,
			layers.LinkTypeEthernet,
			gopacket.DecodeOptions{
				Lazy:               true,
				NoCopy:             true,
				SkipDecodeRecovery: false},
		)
		iface = net.Interface{Index: 1, MTU: 1500, Name: "eth0", HardwareAddr: []byte{00, 50, 56, 99, 49, 190}, Flags: 19}
		packet = Packet{Iface: iface, Packet: gopacket}

		cdp = []byte{1, 0, 12, 204, 204, 204, 0, 80, 86, 153, 78, 253, 0, 235, 170, 170, 3, 0,
			0, 12, 32, 0, 2, 180, 18, 138, 0, 1, 0, 31, 115, 119, 105, 116, 99, 104, 40, 49, 50, 51,
			49, 57, 52, 57, 53, 48, 51, 49, 52, 55, 50, 53, 57, 48, 48, 53, 41, 0, 2, 0, 17, 0, 0,
			0, 1, 1, 1, 204, 0, 4, 10, 240, 16, 50, 0, 3, 0, 12, 99, 111, 110, 116, 114, 111, 108,
			48, 0, 4, 0, 8, 0, 0, 0, 40, 0, 5, 0, 73, 67, 105, 115, 99, 111, 32, 78, 101, 120, 117,
			115, 32, 79, 112, 101, 114, 97, 116, 105, 110, 103, 32, 83, 121, 115, 116, 101, 109, 32,
			40, 78, 88, 45, 79, 83, 41, 32, 83, 111, 102, 116, 119, 97, 114, 101, 44, 32, 86, 101,
			114, 115, 105, 111, 110, 32, 53, 46, 50, 40, 49, 41, 83, 86, 51, 40, 50, 46, 49, 41, 0, 6,
			0, 14, 78, 101, 120, 117, 115, 49, 48, 48, 48, 86, 0, 18, 0, 5, 0, 0, 19, 0, 5, 0, 0, 11,
			0, 5, 1, 0, 17, 0, 8, 0, 0, 5, 220, 0, 20, 0, 10, 115, 119, 105, 116, 99, 104, 0, 21, 0,
			18, 6, 12, 43, 6, 1, 4, 1, 9, 12, 3, 1, 3, 134, 72, 0, 22, 0, 17, 0, 0, 0, 1, 1, 1, 204,
			0, 4, 10, 240, 16, 50}
		cdpGopacket := gopacket.NewPacket(
			cdpData,
			layers.LinkTypeEthernet,
			cdpGopacket.DecodeOptions{
				Lazy:               true,
				NoCopy:             true,
				SkipDecodeRecovery: false},
		)
		cdpIface = net.Interface{Index: 1, MTU: 1500, Name: "eth0", HardwareAddr: []byte{00, 50, 56, 99, 49, 190}, Flags: 19}
		cdpPacket = Packet{Iface: cdpIface, Packet: cdpGopacket}

		wg = &sync.WaitGroup{}

	})

	JustBeforeEach(func() {
		neighbors, _ = NewNeighbors()
	})

	Describe("NewNeighbors", func() {
		It("should return a Neighbors struct", func() {
			neighbors, err := NewNeighbors()
			Expect(err).To(Succeed())
			Expect(neighbors).To(BeAssignableToTypeOf(&Neighbors{}))
		})
	})

	Describe("ProcessPacket", func() {
		Context("when a good packet is received", func() {
			It("should decode a fake packet", func() {
				wg.Add(1)
				neighbors.ProcessPacket(wg, packet)
				neighbors.Rw.Lock()
				defer neighbors.Rw.Unlock()
				Expect(len(neighbors.Interfaces["eth0"])).Should(Equal(1))
			})
		})
		Context("when a good cdpPacket is received", func() {
			It("should decode a fake packet", func() {
				wg.Add(1)
				neighbors.ProcessPacket(wg, cdpPacket)
				neighbors.Rw.Lock()
				defer neighbors.Rw.Unlock()
				Expect(len(neighbors.Interfaces["eth0"])).Should(Equal(1))
			})
		})

		Context("when a bad packet is received", func() {
			BeforeEach(func() {
				data := []byte{}
				gopacket := gopacket.NewPacket(
					data,
					layers.LinkTypeEthernet,
					gopacket.DecodeOptions{
						Lazy:               true,
						NoCopy:             true,
						SkipDecodeRecovery: false},
				)
				packet = Packet{Iface: iface, Packet: gopacket}
			})
			It("should handle empty packet data", func() {
				wg.Add(1)
				neighbors.ProcessPacket(wg, packet)
				neighbors.Rw.Lock()
				defer neighbors.Rw.Unlock()
				Expect(len(neighbors.Interfaces["eth0"])).Should(Equal(0))
			})
		})
	})

	Describe("expireNeighbors", func() {
		It("show expire timed out neighbors", func() {
			wg.Add(1)
			neighbors.ProcessPacket(wg, packet)
			Eventually(func() int {
				neighbors.Rw.Lock()
				defer neighbors.Rw.Unlock()
				return len(neighbors.Interfaces["eth0"])
			}, "10s", "1s").Should(Equal(0))
		})
	})
})
