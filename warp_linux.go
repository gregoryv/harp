package warp

import (
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

	// Create the Raw Socket
	// AF_PACKET: Address family for the device-level packet interface
	// SOCK_RAW: We supply the entire frame (including Ethernet header)
	// ETH_P_ARP: Filter for ARP protocol (0x0806 in host byte order)
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(unix.ETH_P_ARP))
	if err != nil {
		return fmt.Errorf("Failed to create raw socket. You need root privileges: %v", err)
	}
	defer syscall.Close(fd) // todo why Close after each send

	// Define the target address structure (SockaddrLinklayer) This
	// tells the kernel which interface to send the raw data out of.
	addr := syscall.SockaddrLinklayer{
		Protocol: unix.ETH_P_ARP,
		Ifindex:  iface.Index,
		Hatype:   unix.ARPHRD_ETHER, // Hardware type: Ethernet
		Pkttype:  syscall.PACKET_HOST,
		Halen:    6, // MAC address length
	}

	data := NewARPRequest(ip, &iface, srcIP)
	debug.Println("sendARP to", ip.String(), len(data), "bytes")
	return syscall.Sendto(fd, data, 0, &addr)
}
