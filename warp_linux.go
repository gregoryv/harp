package warp

import (
	"errors"
	"fmt"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

func sendARP(ips []net.IP, iface net.Interface) error {
	// Create the Raw Socket
	// AF_PACKET: Address family for the device-level packet interface
	// SOCK_RAW: We supply the entire frame (including Ethernet header)
	// ETH_P_ARP: Filter for ARP protocol (0x0806 in host byte order)
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(unix.ETH_P_ARP))
	if err != nil {
		return fmt.Errorf("Failed to create raw socket. You need root privileges: %v", err)
	}
	defer syscall.Close(fd)

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
	debug.Println("writeARP to", ip.String(), len(data), "bytes")
	return syscall.Sendto(fd, data, 0, &addr)
}

// getInterfaceIPv4 finds the primary IPv4 address of a given interface.
func getInterfaceIPv4(iface *net.Interface) (net.IP, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.To4(), nil
			}
		}
	}
	return nil, fmt.Errorf("interface %s has no IPv4 address", iface.Name)
}
