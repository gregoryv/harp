package warp

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// IPRange converts
// - 192.1.1.1 to one ip
// - 192.1.1.* to 255 range 192.1.1.1 to 192.1.1.255
// - 192.1.1.5-15 to 10 192.1.1.5 to 192.1.1.15
func IPRange(rangestr string) ([]net.IP, error) {
	if strings.HasSuffix(rangestr, "*") {
		res := make([]net.IP, 0, 255)
		prefix := rangestr[:len(rangestr)-1]
		for i := 1; i <= 255; i++ {
			ipstr := prefix + strconv.Itoa(i)
			ip := net.ParseIP(ipstr)
			if ip != nil {
				res = append(res, ip)
			}
		}
		return res, nil
	}
	if strings.Contains(rangestr, "-") {
		i := strings.LastIndex(rangestr, ".")
		j := strings.LastIndex(rangestr, "-")
		if j <= i || i == -1 || j == -1 {
			return nil, fmt.Errorf("%s: invalid rangestr", rangestr)
		}
		from, _ := strconv.Atoi(rangestr[i+1 : j])
		to, _ := strconv.Atoi(rangestr[j+1:])
		if from <= 0 || from > 255 || to <= 0 || to > 255 || from >= to {
			return nil, fmt.Errorf("%s: check range", rangestr)
		}
		prefix := rangestr[:i+1]
		res := make([]net.IP, 0, to-from+1)
		for i := from; i <= to; i++ {
			ipstr := prefix + strconv.Itoa(i)
			ip := net.ParseIP(ipstr)
			if ip != nil {
				res = append(res, ip)
			}
		}
		return res, nil
	}
	// assume only one

	ip := net.ParseIP(rangestr)
	if ip == nil {
		return nil, fmt.Errorf("%s: bad range", rangestr)
	}
	return []net.IP{ip}, nil
}
