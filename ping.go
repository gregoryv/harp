package warp

import (
	"net"
	"os/exec"
	"runtime"
	"sync"
)

func pingAll(ips []net.IP) {
	var wg sync.WaitGroup
	for _, ip := range ips {
		wg.Add(1)
		go func() {
			ping(ip.String())
			wg.Done()
		}()
	}
	wg.Wait()
}

func ping(addr string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Windows ping default uses -n for count
		cmd = exec.Command("ping", "-n", "1", addr)
	} else {
		// Linux/macOS use -c for count
		cmd = exec.Command("ping", "-c", "1", addr)
	}
	debug.Println("ping", addr)
	if err := cmd.Run(); err != nil {
		debug.Print(err)
	}
}
