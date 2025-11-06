package warp

import (
	"net"
	"slices"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// writeARP writes an ARP request for each address on our local
// network to the pcap handle.
func NewARPRequest(iface *net.Interface, addr *net.IPNet, ip net.IP) []byte {
	// Set up all the layers' fields we can.
	eth := layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(iface.HardwareAddr),
		SourceProtAddress: []byte(addr.IP),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
	}
	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	arp.DstProtAddress = []byte(ip)
	gopacket.SerializeLayers(buf, opts, &eth, &arp)
	return slices.Clone(buf.Bytes())
}
