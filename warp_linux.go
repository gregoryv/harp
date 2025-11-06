package warp

import (
	"net"
)

func sendARP(ips []net.IP, _ net.Interface) error {
	pingAll(ips)
	return nil
}
