package discovery

import (
	"net"
	"testing"
)

func TestLocalIpGet(t *testing.T) {
	ifaces, err := net.Interfaces()
	if err != nil {
		t.Log(err)
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			v, ok := addr.(*net.IPNet)
			if ok {
				if v.IP.To4() != nil && !v.IP.IsLoopback() {
					t.Log(v.IP)
				}
			}
		}
	}
	// t.Log(ifaces)
}
