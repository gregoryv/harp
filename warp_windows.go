package warp

import (
	"net"
)

// sendARP doesn't use arp as SendARP doesn't actually send ARP as
// expected.
func sendARP(ips []net.IP, _ net.Interface) error {
	pingAll(ips)
	return nil
}
