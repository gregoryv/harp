package warp

import (
	"errors"
	"net"
)

func SendARP(ips []net.IP) error {
	var err error
	for _, ip := range ips {
		e := sendARP(ip)
		err = errors.Join(err, e)
	}
	return err
}
