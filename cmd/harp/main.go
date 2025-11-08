package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/gregoryv/harp"
)

func main() {
	flag.Usage = func() {
		fmt.Println(`Usage: harp [IP]

Examples

  $ harp 192.168.1.3
  $ harp 192.168.1.3-9
  $ harp 192.168.1.*

without IP harp shows the arp -a cache only.`)
	}
	log.SetFlags(0)
	flag.Parse()

	targetIP := flag.Arg(0)
	if targetIP != "" {
		ips, err := harp.IPRange(targetIP)
		if err != nil {
			log.Fatal(err)
		}

		// remove existing ips if already cached
		cache, _ := harp.Cache()
		existingIP := make(map[string]struct{})
		for _, hit := range cache {
			existingIP[hit.IP] = struct{}{}
		}

		filtered := make([]net.IP, 0, len(ips))
		for _, ip := range ips {
			if _, found := existingIP[ip.String()]; !found {
				filtered = append(filtered, ip)
			}
		}

		if err := harp.Scan(filtered); err != nil {
			log.Fatal(err)
		}
	}

	result, err := harp.Cache()
	if err != nil {
		log.Fatal(err)
	}
	for _, hit := range result {
		fmt.Println(hit.MAC, hit.IP)
	}
}
