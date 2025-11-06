package warp

import (
	"net"
)

func Scan(ips []net.IP) error {
	headAll(ips)
	return nil
}
