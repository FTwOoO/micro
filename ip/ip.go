package ip

import (
	"net"
	"os"
	"strings"
)

func InternalIP() string {
	intranet_ip := os.Getenv("INTRANET_IP")
	if intranet_ip != "" {
		ip := net.ParseIP(intranet_ip)
		if ip != nil {
			return ip.String()
		}
	}

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
