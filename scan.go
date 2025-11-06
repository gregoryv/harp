package harp

import (
	"net"
	"net/http"
	"sync"
	"time"
)

func Scan(ips []net.IP) error {
	client := &http.Client{
		Timeout: 100 * time.Millisecond, // total time for the request
	}
	var wg sync.WaitGroup

	for _, ip := range ips {
		wg.Add(1)
		go func() {
			url := "http://" + ip.String()
			client.Head(url)
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}
