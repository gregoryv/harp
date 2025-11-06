package warp

import (
	"net"
)

func SendARP(ips []net.IP, iface net.Interface) error {
	return sendARP(ips, iface)
}
