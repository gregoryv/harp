package main

import (
	"flag"

	"github.com/gregoryv/warp"
)

func main() {
	// Example usage: Try to find the MAC address of a local gateway or specific device
	targetIP := flag.String("ip", "127.0.0.1", "IP you want to test")
	flag.Parse()

	warp.SendARP(*targetIP)
}
