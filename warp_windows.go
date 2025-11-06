package warp

import (
	"fmt"
	"net"
	"unsafe"

	"golang.org/x/sys/windows"
)

func SendARP(targetIPStr string) error {
	// 1. Convert the target IP string to a 32-bit unsigned integer
	ip := net.ParseIP(targetIPStr)
	if ip == nil {
		return fmt.Errorf("invalid IP address format: %s", targetIPStr)
	}
	// We only work with IPv4 addresses for SendARP
	ipv4 := ip.To4()
	if ipv4 == nil {
		return fmt.Errorf("IP address is not a valid IPv4 address: %s", targetIPStr)
	}

	// Convert byte slice to uint32 (big-endian to little-endian for Windows IPAddr struct)
	destIP := uint32(ipv4[0]) | uint32(ipv4[1])<<8 | uint32(ipv4[2])<<16 | uint32(ipv4[3])<<24

	// Note: You can optionally specify a source IP here instead of 0 to force a specific interface.
	srcIP := uint32(0) // 0 lets Windows choose the best interface

	// 2. Prepare buffers for the response
	macAddr := [6]byte{}           // Buffer to store the 6-byte MAC address
	macLen := uint32(len(macAddr)) // Length of the MAC address buffer

	// 3. Call the SendARP API function
	// The call takes raw pointers via unsafe, as required for system calls
	_, _, err := sendARP.Call(
		uintptr(destIP),
		uintptr(srcIP),
		uintptr(unsafe.Pointer(&macAddr[0])),
		uintptr(unsafe.Pointer(&macLen)),
	)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

// Define the SendARP function signature for linking with iphlpapi.dll
// This signature matches the C function: DWORD SendARP(IPAddr DestIP,
// IPAddr SrcIP, PULONG pMacAddr, PULONG PhyAddrLen);
var (
	iphlpapi = windows.NewLazySystemDLL("iphlpapi.dll")
	sendARP  = iphlpapi.NewProc("SendARP")
)
