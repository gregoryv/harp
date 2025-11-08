package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gregoryv/harp"
)

func main() {
	flag.Usage = func() {
		fmt.Printf(`Usage: %s [OPTIONS] [IP-range]

Examples

  $ harp 192.168.1.3
  $ harp 192.168.1.3-9
  $ harp 192.168.1.*

without IP-range harp shows the arp -a cache only.

Options
`, os.Args[0])
		flag.PrintDefaults()
	}
	log.SetFlags(0)
	flag.Parse()

	version()

	rangestr := flag.Arg(0)
	if rangestr != "" {
		ips, err := harp.IPRange(rangestr)
		if err != nil {
			log.Fatal(err)
		}

		// remove existing ips if already cached
		cache, err := harp.Cache()
		if err != nil {
			// fail early here incase the command is not found
			log.Fatal(err)
		}
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
