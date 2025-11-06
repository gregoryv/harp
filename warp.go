package warp

import (
	"fmt"
	"net"
)

func SendARP(ips []net.IP, iface net.Interface) error {
	return sendARP(ips, iface)
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
