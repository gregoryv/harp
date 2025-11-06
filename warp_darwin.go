package warp

import (
	"net"
)

// sendARP uses ping, difficult to send raw ARP request
func sendARP(ips []net.IP, _ net.Interface) error {
	pingAll(ips)
	return nil
}
