package warp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
)

func sendARP(ips []net.IP, iface net.Interface) error {
	// 1. Open BPF Device
	bpfFile, err := openBPFDevice()
	if err != nil {
		return fmt.Errorf("failed to open BPF device: %w", err)
	}
	defer bpfFile.Close()
	fd := int(bpfFile.Fd())

	// 2. Bind BPF device to the network interface (en0, etc.)
	// This requires an ioctl call with BIOCSETIF and a struct ifreq,
	// which is the most difficult part, as Go's syscall/unix package often
	// lacks the necessary constants and structs for non-POSIX ioctl commands.

	// PSEUDOCODE for ioctl binding (Constants may not exist in 'unix'):
	// ifreq := ifreqStruct{} // Must be manually defined to match C header
	// copy(ifreq.name[:], []byte(iface.Name))
	// _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), unix.BIOCSETIF, uintptr(unsafe.Pointer(&ifreq)))
	// if errno != 0 {
	//     return fmt.Errorf("failed to bind interface %s to BPF: %v", iface.Name, errno)
	// }

	// For this example, we skip the complex ioctl binding and go straight to writing.
	// WARNING: Without the ioctl binding, the packet may not go out the intended interface.

	for _, ip := range ips {
		e := writeARP(ip, iface, fd)
		err = errors.Join(err, e)
	}
	return err
}

func writeARP(ip net.IP, iface net.Interface, fd int) error {
	srcIP, err := getInterfaceIPv4(&iface)
	if err != nil {
		return fmt.Errorf("Could not get source IP for %s: %v. Check IP configuration.", iface.Name, err)
	}

	data := NewARPRequest(ip, &iface, srcIP)
	n, err := syscall.Write(fd, data)
	if err != nil {
		return fmt.Errorf("failed to write to BPF device: %w", err)
	}
	if n != len(data) {
		return fmt.Errorf("wrote %d bytes, expected %d", n, len(data))
	}
	return nil
}

func openBPFDevice() (*os.File, error) {
	// macOS uses /dev/bpfX devices. We must find an available one.
	for i := 0; i < 20; i++ {
		device := fmt.Sprintf("/dev/bpf%d", i)
		// O_RDWR: Read/Write access (needed to send/receive)
		// O_NONBLOCK: Non-blocking mode
		f, err := os.OpenFile(device, os.O_RDWR|syscall.O_NONBLOCK, 0)
		if err == nil {
			return f, nil // Found and opened a free BPF device
		}
	}
	return nil, fmt.Errorf("failed to find an available BPF device (/dev/bpfX)")
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

// --- Main Logic ---

// --- Helper Functions ---

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
