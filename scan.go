package harp

import (
	"net"
	"sync"
)

// Scan tries to do a ipv4 connection to each ip on port 80
func Scan(ips []net.IP) error {
	var wg sync.WaitGroup

	for _, ip := range ips {
		wg.Add(1)
		go func() {
			net.Dial("tcp4", ip.String()+":80")
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}
