package main

import (
	"bytes"
	"errors"
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
	available := listInterfaces()
	if len(available) == 0 {
		log.Fatal("no interfaces available")
	}
	iface := flag.String("iface", available[0].Name, "interface to scan")
	verbose := flag.Bool("verbose", false, "debug logs")
	flag.Parse()

	if *verbose {
		warp.SetDebugOutput(os.Stderr)
	}

	if *iface == "" {
		showInterfaces(listInterfaces())
		log.Fatal("missing -iface")
		os.Exit(1)
	}
	var selectedInterface net.Interface
	for _, ic := range available {
		if ic.Name == *iface {
			selectedInterface = ic
		}
	}

	if *targetIP != "" {
		ips, err := warp.IPRange(*targetIP)
		if err != nil {
			log.Fatal(err)
		}
		if err := warp.SendARP(ips, selectedInterface); err != nil {
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

func showInterfaces(available []net.Interface) {
	fmt.Println("# available interfaces")
	for _, iface := range available {
		addr, _ := findIP4(&iface)
		ip := ipstr(addr)
		fmt.Fprintf(os.Stderr, "%s %s\n", iface.Name, ip)
	}
}

// listInterfaces returns a filtered list of interfaces that may be
// scanned.
func listInterfaces() []net.Interface {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	var available []net.Interface
	for _, iface := range ifaces {
		addr, _ := findIP4(&iface)
		ip := ipstr(addr)
		if ip == "" {
			continue
		}
		if strings.HasPrefix(ip, "172.") {
			continue
		}
		available = append(available, iface)
	}
	return available
}

func ipstr(addr *net.IPNet) string {
	if addr == nil {
		return ""
	}
	v := addr.String()
	i := strings.Index(v, "/")
	if i < 0 {
		return v
	}
	return v[:i]
}

func findIP4(iface *net.Interface) (*net.IPNet, error) {
	var addr *net.IPNet
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				addr = &net.IPNet{
					IP:   ip4,
					Mask: ipnet.Mask[len(ipnet.Mask)-4:],
				}
				break
			}
		}
	}
	// Sanity-check that the interface has a good address.
	switch {
	case addr == nil:
		return nil, fmt.Errorf("%s: not found", iface.Name)

	case addr.IP[0] == 127:
		return nil, errors.New("is localhost")

	case addr.Mask[0] != 0x_ff || addr.Mask[1] != 0xff:
		return nil, errors.New("mask too large")
	}

	return addr, nil
}
