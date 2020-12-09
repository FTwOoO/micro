package server

import (
	"net"
	"strings"
)

func InternalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, inter := range ifaces {
		if inter.Flags&net.FlagUp != 0 && (strings.HasPrefix(inter.Name, "eth") || strings.HasPrefix(inter.Name, "en")) {
			addrs, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}

				}
			}
		}
	}
	return ""
}
