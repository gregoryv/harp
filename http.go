package warp

import (
	"net"
	"net/http"
	"sync"
	"time"
)

func headAll(ips []net.IP) {
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
}
