package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/gregoryv/warp"
)

func main() {
	log.SetFlags(0)
	targetIP := flag.String(
		"ip", "", "arp IP range, e.g 192.1.1.3-128 or 192.1.1.*",
	)
	verbose := flag.Bool("verbose", false, "")
	flag.Parse()

	if *verbose {
		warp.SetDebugOutput(os.Stderr)
	}

	if *targetIP != "" {
		ips, err := warp.IPRange(*targetIP)
		if err != nil {
			log.Fatal(err)
		}
		if err := warp.Scan(ips); err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Duration(len(ips)) * time.Millisecond)
	}
	data, err := exec.Command("arp", "-a").Output()
	if err != nil {
		log.Fatal(err)
	}
	res := warp.ParseARPCache(bytes.NewReader(data))

	var ipList []string
	for ip := range res {
		switch {
		case strings.HasPrefix(ip, "224.0.0."):
		case strings.HasPrefix(ip, "239.192.152."):
		case strings.HasPrefix(ip, "239.255.255."):
		default:
			ipList = append(ipList, ip)
		}
	}
	slices.SortFunc(ipList, func(a, b string) int {
		ipA := net.ParseIP(a).To4()[3]
		ipB := net.ParseIP(b).To4()[3]
		switch {
		case ipA < ipB:
			return -1
		case ipA > ipB:
			return 1
		default:
			return 0
		}
	})
	for _, ip := range ipList {
		mac := strings.ToLower(res[ip])
		if mac == "ff:ff:ff:ff:ff:ff" {
			continue
		}
		fmt.Println(mac, ip)
	}
}
