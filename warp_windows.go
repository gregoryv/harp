package warp

import (
	"fmt"
	"net"
	"unsafe"

	"golang.org/x/sys/windows"
)

func sendARP(ip net.IP, _ net.Interface) error {
	ipv4 := ip.To4()
	if ipv4 == nil {
		return fmt.Errorf("sendARP(%q): not ipv4", ip.String())
	}

	// Convert byte slice to uint32 (big-endian to little-endian for Windows IPAddr struct)
	destIP := uint32(ipv4[0]) | uint32(ipv4[1])<<8 | uint32(ipv4[2])<<16 | uint32(ipv4[3])<<24

	// Note: You can optionally specify a source IP here instead of 0 to force a specific interface.
	// todo ip from interface name
	srcIP := uint32(0) // 0 lets Windows choose the best interface

	// 2. Prepare buffers for the response
	macAddr := [6]byte{}           // Buffer to store the 6-byte MAC address
	macLen := uint32(len(macAddr)) // Length of the MAC address buffer

	// 3. Call the SendARP API function
	// The call takes raw pointers via unsafe, as required for system calls
	phlSendARP.Call(
		uintptr(destIP),
		uintptr(srcIP),
		uintptr(unsafe.Pointer(&macAddr[0])),
		uintptr(unsafe.Pointer(&macLen)),
	)
	return nil
}

// Define the SendARP function signature for linking with iphlpapi.dll
// This signature matches the C function: DWORD SendARP(IPAddr DestIP,
// IPAddr SrcIP, PULONG pMacAddr, PULONG PhyAddrLen);
var (
	iphlpapi   = windows.NewLazySystemDLL("iphlpapi.dll")
	phlSendARP = iphlpapi.NewProc("SendARP")
)
