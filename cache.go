package warp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
)

/*
	ie. output from $ arp -a

Interface: 192.168.1.71 --- 0x3

	Internet Address      Physical Address      Type
	192.168.1.1           ac-8b-a9-ab-b1-ad     dynamic
	192.168.1.2           00-01-c0-2d-07-39     dynamic
	192.168.1.58          ac-91-a1-d5-0c-aa     dynamic
	192.168.1.59          ac-cc-8e-c2-2f-8b     dynamic
	192.168.1.150         ac-cc-8e-c3-3c-2c     dynamic
	192.168.1.190         30-05-5c-a1-2a-77     dynamic
	192.168.1.255         ff-ff-ff-ff-ff-ff     static
	224.0.0.22            01-00-5e-00-00-16     static
	224.0.0.251           01-00-5e-00-00-fb     static
	224.0.0.252           01-00-5e-00-00-fc     static
	239.192.152.143       01-00-5e-40-98-8f     static
	239.255.255.250       01-00-5e-7f-ff-fa     static
	255.255.255.255       ff-ff-ff-ff-ff-ff     static
*/
func parseArpWindows(r io.Reader) Result {
	res := make(Result)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		line = whitespaces.ReplaceAllString(line, " ")
		line = strings.TrimSpace(line)
		parts := strings.Split(line, " ")

		var ip, mac string
		if len(parts) > 0 {
			ip = parts[0]
		}
		if len(parts) > 1 {
			mac = parts[1]
		}
		if v := net.ParseIP(ip); len(v) == 0 {
			continue
		}
		mac = strings.ReplaceAll(mac, "-", ":")
		mac = strings.ToUpper(mac)
		res[mac] = ip
	}
	return res
}

/*
? (192.168.1.1) at ac:8b:a9:ab:b1:ad on en0 ifscope [ethernet]
? (192.168.1.2) at 0:1:c0:2d:7:39 on en0 ifscope [ethernet]
? (192.168.1.46) at e8:68:e7:6f:da:b5 on en0 ifscope [ethernet]
? (192.168.1.47) at (incomplete) on en0 ifscope [ethernet]
? (192.168.1.61) at 40:cb:c0:e5:4b:bc on en0 ifscope [ethernet]
? (192.168.1.190) at 30:5:5c:a1:2a:77 on en0 ifscope [ethernet]
? (192.168.1.255) at ff:ff:ff:ff:ff:ff on en0 ifscope [ethernet]
mdns.mcast.net (224.0.0.251) at 1:0:5e:0:0:fb on en0 ifscope permanent [ethernet]
*/
func parseArpDarwin(r io.Reader) Result {
	res := make(Result)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		line = strings.TrimSpace(line)
		parts := strings.Split(line, " ")

		var ip, mac string
		if len(parts) > 1 {
			ip = parts[1]
		}
		if len(parts) > 3 {
			mac = parts[3]
		}
		ip = strings.Trim(ip, "()")
		if v := net.ParseIP(ip); len(v) == 0 {
			fmt.Println(ip)
			continue
		}
		if _, err := net.ParseMAC(mac); err != nil {
			continue
		}
		mac = strings.ToUpper(mac)
		res[mac] = ip
	}
	return res
}

/*
? (192.168.1.62) at f0:9f:c2:79:5c:ab [ether] on enp4s0
? (192.168.1.220) at 06:9e:6f:94:0c:2e [ether] on enp4s0
? (192.168.1.71) at c8:7f:54:03:a3:e3 [ether] on enp4s0
? (192.168.1.58) at ac:91:a1:d5:0c:aa [ether] on enp4s0
? (192.168.1.42) at d8:b3:70:b0:0a:7d [ether] on enp4s0
? (192.168.1.190) at 30:05:5c:a1:2a:77 [ether] on enp4s0
? (192.168.1.41) at e0:63:da:b6:5a:74 [ether] on enp4s0
? (192.168.1.20) at <incomplete> on enp4s0
? (192.168.1.55) at f4:fe:fb:2e:c7:bc [ether] on enp4s0
? (192.168.1.188) at e4:0d:36:fe:8c:f1 [ether] on enp4s0
? (192.168.1.213) at f0:9f:c2:60:2b:17 [ether] on enp4s0
*/
func parseArpLinux(r io.Reader) Result {
	return parseArpDarwin(r)
}

var whitespaces = regexp.MustCompile(`\s+`)

type Result map[string]string
