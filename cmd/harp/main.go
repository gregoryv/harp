package main

import (
	"flag"
	"fmt"
	"log"
	"time"

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
		if err := harp.Scan(ips); err != nil {
			log.Fatal(err)
		}

		time.Sleep(8 * time.Duration(len(ips)) * time.Millisecond)
	}

	result, err := harp.Cache()
	if err != nil {
		log.Fatal(err)
	}
	for _, hit := range result {
		fmt.Println(hit.MAC, hit.IP)
	}
}
