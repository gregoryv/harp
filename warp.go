package warp

import (
	"net"
)

func SendARP(ips []net.IP, iface net.Interface) error {
	headAll(ips)
	return nil
}
