package main

import (
	"flag"
	"log"

	"github.com/gregoryv/warp"
)

func main() {
	targetIP := flag.String("ip", "127.0.0.1", "arp IP range")
	flag.Parse()

	log.SetFlags(0)

	ips, err := warp.IPRange(*targetIP)
	if err != nil {
		log.Fatal(err)
	}
	for _, ip := range ips {
		warp.SendARP(ip)
	}
}
