package warp

import (
	"errors"
	"net"
)

func SendARP(ips []net.IP, ic net.Interface) error {
	var err error
	for _, ip := range ips {
		e := sendARP(ip, ic)
		err = errors.Join(err, e)
	}
	return err
}
