package warp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

func sendARP(ip net.IP, iface net.Interface) error {
	srcIP, err := getInterfaceIPv4(&iface)
	if err != nil {
		return fmt.Errorf("Could not get source IP for %s: %v. Check IP configuration.", iface.Name, err)
	}

	// 2. Build the Raw ARP Packet
	// buildARPPacket(iface.HardwareAddr, srcIP, ip)
	packetBytes := NewARPRequest(ip, &iface, srcIP)

	// 3. Create the Raw Socket
	// AF_PACKET: Address family for the device-level packet interface
	// SOCK_RAW: We supply the entire frame (including Ethernet header)
	// ETH_P_ARP: Filter for ARP protocol (0x0806 in host byte order)
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(unix.ETH_P_ARP))
	if err != nil {
		return fmt.Errorf("Failed to create raw socket. You need root privileges: %v", err)
	}
	defer syscall.Close(fd) // todo why Close after each send

	// 4. Define the target address structure (SockaddrLinklayer)
	// This tells the kernel which interface to send the raw data out of.
	addr := syscall.SockaddrLinklayer{
		Protocol: unix.ETH_P_ARP,
		Ifindex:  iface.Index,
		Hatype:   unix.ARPHRD_ETHER, // Hardware type: Ethernet
		Pkttype:  syscall.PACKET_HOST,
		Halen:    6, // MAC address length
	}

	return syscall.Sendto(fd, packetBytes, 0, &addr)
}

// --- Packet Structure Definitions ---

// EthernetHeader matches the first 14 bytes of an Ethernet frame.
type EthernetHeader struct {
	DestMAC   [6]byte // Broadcast FF:FF:FF:FF:FF:FF
	SrcMAC    [6]byte // Local interface MAC
	EtherType uint16  // 0x0806 for ARP
}

// ARPHeader matches the 28-byte payload of an ARP packet.
type ARPHeader struct {
	HType     uint16  // Hardware type (1 for Ethernet)
	PType     uint16  // Protocol type (0x0800 for IPv4)
	HLen      uint8   // Hardware address length (6 bytes for MAC)
	PLen      uint8   // Protocol address length (4 bytes for IPv4)
	Operation uint16  // Operation (1 for Request, 2 for Reply)
	SHA       [6]byte // Sender hardware address (Source MAC)
	SPA       [4]byte // Sender protocol address (Source IP)
	THA       [6]byte // Target hardware address (00:00:00:00:00:00 for Request)
	TPA       [4]byte // Target protocol address (Target IP)
}

// buildARPPacket manually serializes the Ethernet and ARP headers into a byte buffer.
func buildARPPacket(srcMAC net.HardwareAddr, srcIP, targetIP net.IP) []byte {
	// 1. Ethernet Header
	ethHdr := EthernetHeader{
		DestMAC:   [6]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, // Broadcast
		EtherType: 0x0806,                                      // ARP (needs to be in Network Byte Order, which is Big Endian)
	}
	copy(ethHdr.SrcMAC[:], srcMAC)

	// 2. ARP Header
	arpHdr := ARPHeader{
		HType:     1,         // Hardware type: Ethernet
		PType:     0x0800,    // Protocol type: IPv4
		HLen:      6,         // MAC length
		PLen:      4,         // IPv4 length
		Operation: 1,         // ARP Request
		THA:       [6]byte{}, // Target MAC (all zeros in request)
	}
	copy(arpHdr.SHA[:], srcMAC)
	copy(arpHdr.SPA[:], srcIP.To4())
	copy(arpHdr.TPA[:], targetIP.To4())

	// 3. Serialize Headers using Big Endian (Network Byte Order)
	buf := new(bytes.Buffer)

	// Note: The structure fields are written manually to ensure correct packing
	// and network byte order (Big Endian).

	// Write Ethernet Header
	buf.Write(ethHdr.DestMAC[:])
	buf.Write(ethHdr.SrcMAC[:])
	binary.Write(buf, binary.BigEndian, ethHdr.EtherType) // 0x0806

	// Write ARP Header
	binary.Write(buf, binary.BigEndian, arpHdr.HType) // 1
	binary.Write(buf, binary.BigEndian, arpHdr.PType) // 0x0800
	buf.WriteByte(arpHdr.HLen)
	buf.WriteByte(arpHdr.PLen)
	binary.Write(buf, binary.BigEndian, arpHdr.Operation) // 1 (Request)
	buf.Write(arpHdr.SHA[:])
	buf.Write(arpHdr.SPA[:])
	buf.Write(arpHdr.THA[:])
	buf.Write(arpHdr.TPA[:])

	// ARP packets must be at least 42 bytes (28 ARP + 14 Eth).
	// The Ethernet frame must be 60 bytes minimum. Padding is required.
	// The full length should be 60 bytes (excluding 4-byte CRC).
	paddingLen := 60 - buf.Len()
	if paddingLen > 0 {
		buf.Write(bytes.Repeat([]byte{0x00}, paddingLen))
	}

	return buf.Bytes()
}
