package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gregoryv/harp"
)

func main() {
	log.SetFlags(0)
	targetIP := flag.String(
		"ip", "", "IP range to scan, e.g 192.1.1.3-128 or 192.1.1.*",
	)
	flag.Parse()

	if *targetIP != "" {
		ips, err := harp.IPRange(*targetIP)
		if err != nil {
			log.Fatal(err)
		}
		if err := harp.Scan(ips); err != nil {
			log.Fatal(err)
		}

		time.Sleep(8 * time.Duration(len(ips)) * time.Millisecond)
	}

	for _, hit := range harp.Cache() {
		fmt.Println(hit.MAC, hit.IP)
	}
}
