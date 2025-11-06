package harp

import (
	"strings"
	"testing"
)

func TestCache(t *testing.T) {
	Cache()
}

func Test_parseArpWindows(t *testing.T) {
	// arp -a
	arpOutput := `

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
`
	r := strings.NewReader(arpOutput)
	res := parseArpWindows(r)
	if v := len(res); v == 0 {
		t.Error("empty result")
	}
	for k, v := range res {
		t.Log(k, v)
	}
}

func Test_parseArpDarwin(t *testing.T) {
	// arp -a
	arpOutput := `
? (192.168.1.1) at ac:8b:a9:ab:b1:ad on en0 ifscope [ethernet]
? (192.168.1.2) at 0:1:c0:2d:7:39 on en0 ifscope [ethernet]
? (192.168.1.46) at e8:68:e7:6f:da:b5 on en0 ifscope [ethernet]
? (192.168.1.47) at (incomplete) on en0 ifscope [ethernet]
? (192.168.1.61) at 40:cb:c0:e5:4b:bc on en0 ifscope [ethernet]
? (192.168.1.190) at 30:5:5c:a1:2a:77 on en0 ifscope [ethernet]
? (192.168.1.255) at ff:ff:ff:ff:ff:ff on en0 ifscope [ethernet]
mdns.mcast.net (224.0.0.251) at 1:0:5e:0:0:fb on en0 ifscope permanent [ethernet]
`
	r := strings.NewReader(arpOutput)
	res := parseArpDarwin(r)
	if v := len(res); v == 0 {
		t.Error("empty result")
	}
	for k, v := range res {
		t.Log(k, v)
	}
}

func Test_parseArpLinux(t *testing.T) {
	// arp -a
	arpOutput := `
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
_gateway (192.168.1.1) at ac:8b:a9:ab:b1:ad [ether] on enp4s0
`
	r := strings.NewReader(arpOutput)
	res := parseArpLinux(r)
	if v := len(res); v == 0 {
		t.Error("empty result")
	}
	for k, v := range res {
		t.Log(k, v)
	}
}

// ip neigh show
const ipOutputLinux = `
192.168.1.62 dev enp4s0 lladdr f0:9f:c2:79:5c:ab REACHABLE 
192.168.1.220 dev enp4s0 lladdr 06:9e:6f:94:0c:2e REACHABLE 
192.168.1.71 dev enp4s0 lladdr c8:7f:54:03:a3:e3 REACHABLE 
192.168.1.58 dev enp4s0 lladdr ac:91:a1:d5:0c:aa STALE 
192.168.1.42 dev enp4s0 lladdr d8:b3:70:b0:0a:7d DELAY 
192.168.1.190 dev enp4s0 lladdr 30:05:5c:a1:2a:77 STALE 
192.168.1.41 dev enp4s0 lladdr e0:63:da:b6:5a:74 REACHABLE 
192.168.1.20 dev enp4s0 FAILED 
192.168.1.55 dev enp4s0 lladdr f4:fe:fb:2e:c7:bc STALE 
192.168.1.188 dev enp4s0 lladdr e4:0d:36:fe:8c:f1 STALE 
192.168.1.213 dev enp4s0 lladdr f0:9f:c2:60:2b:17 REACHABLE 
192.168.1.1 dev enp4s0 lladdr ac:8b:a9:ab:b1:ad REACHABLE 
fe80::c94f:77d0:938:fe58 dev enp4s0 lladdr c8:7f:54:03:a3:e3 STALE
`
