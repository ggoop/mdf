package utils

import (
	"fmt"
	"net"
	"strings"
)

func GetIpAddrs(nets ...bool) []string {
	var ips []string
	if len(nets) > 0 && nets[0] {
		ping := DefaultConfig.GetValue("PING")
		if ping == "" {
			ping = "baidu.com:80"
		}
		if strings.Index(ping, ":") < 0 {
			ping = ping + ":80"
		}
		if conn, err := net.Dial("udp", ping); err != nil {

		} else {
			defer conn.Close()
			ips = append(ips, strings.Split(conn.LocalAddr().String(), ":")[0])
		}
	}
	if len(ips) == 0 {
		netInterfaces, err := net.Interfaces()
		if err != nil {
			fmt.Println("net.Interfaces failed, err:", err.Error())
			return ips
		}
		for i := 0; i < len(netInterfaces); i++ {
			if (netInterfaces[i].Flags & net.FlagUp) != 0 {
				addrs, _ := netInterfaces[i].Addrs()
				for _, address := range addrs {
					if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
						if ipnet.IP.To4() != nil {
							ips = append(ips, ipnet.IP.String())
						}
					}
				}
			}
		}
	}
	return ips
}
