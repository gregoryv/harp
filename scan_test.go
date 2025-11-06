package harp

import "testing"

func TestScan(t *testing.T) {
	ips, _ := IPRange("127.0.0.1")
	if err := Scan(ips); err != nil {
		t.Error(err)
	}
}
