package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/gregoryv/warp"
)

func main() {
	targetIP := flag.String(
		"ip", "", "arp IP range, e.g 192.1.1.3-128 or 192.1.1.*",
	)
	flag.Parse()

	log.SetFlags(0)

	if *targetIP != "" {
		ips, err := warp.IPRange(*targetIP)
		if err != nil {
			log.Fatal(err)
		}
		if err := warp.SendARP(ips); err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Duration(len(ips)) * time.Millisecond)
	}
	data, err := exec.Command("arp", "-a").Output()
	if err != nil {
		log.Fatal(err)
	}
	res := warp.ParseARPCache(bytes.NewReader(data))
	for mac, ip := range res {
		fmt.Println(strings.ToLower(mac), ip)
	}
}
