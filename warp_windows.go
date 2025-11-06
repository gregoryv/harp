package warp

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

func sendARP(ips []net.IP, iface net.Interface) error {
	var wg sync.WaitGroup
	for _, ip := range ips {
		wg.Add(1)
		go func() {
			ping(ip.String())
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

var writeInterval = 4 * time.Millisecond

// on windows we might need to ping the devices as this will probably
// use the cache.
func writeARP(ip net.IP, iface net.Interface) error {
	ipv4 := ip.To4()
	if ipv4 == nil {
		return fmt.Errorf("sendARP(%q): not ipv4", ip.String())
	}

	// Convert byte slice to uint32 (big-endian to little-endian for Windows IPAddr struct)
	destIP := uint32(ipv4[0]) | uint32(ipv4[1])<<8 | uint32(ipv4[2])<<16 | uint32(ipv4[3])<<24

	// Note: You can optionally specify a source IP here instead of 0 to force a specific interface.
	// todo ip from interface name
	// srcIP := uint32(0) // 0 lets Windows choose the best interface
	srcIP4, _ := getInterfaceIPv4(&iface)
	srcIP, _ := IPToUint32(srcIP4)
	// 2. Prepare buffers for the response
	macAddr := [6]byte{}           // Buffer to store the 6-byte MAC address
	macLen := uint32(len(macAddr)) // Length of the MAC address buffer

	// 3. Call the SendARP API function
	// The call takes raw pointers via unsafe, as required for system calls
	_, _, err := phlSendARP.Call(
		uintptr(destIP),
		uintptr(srcIP),
		uintptr(unsafe.Pointer(&macAddr[0])),
		uintptr(unsafe.Pointer(&macLen)),
	)
	if err != nil && !strings.Contains(err.Error(), "completed successfully") {
		debug.Println("writeARP to", ip.String(), err)
	} else {
		debug.Println("writeARP to", ip.String())
	}
	return nil
}

// Define the SendARP function signature for linking with iphlpapi.dll
// This signature matches the C function: DWORD SendARP(IPAddr DestIP,
// IPAddr SrcIP, PULONG pMacAddr, PULONG PhyAddrLen);
var (
	iphlpapi   = windows.NewLazySystemDLL("iphlpapi.dll")
	phlSendARP = iphlpapi.NewProc("SendARP")
)

// IPToUint32 converts a net.IP (assumed to be IPv4) into a uint32.
func IPToUint32(ip net.IP) (uint32, error) {
	// 1. Ensure the IP is a 4-byte IPv4 address.
	ip4 := ip.To4()
	if ip4 == nil {
		return 0, fmt.Errorf("IP address is not a valid IPv4 address: %s", ip.String())
	}

	// 2. Use binary.BigEndian.Uint32 to convert the 4-byte slice to a uint32.
	// IPv4 addresses are always represented in Big Endian (Network Byte Order).
	ipInt := binary.BigEndian.Uint32(ip4)

	return ipInt, nil
}
