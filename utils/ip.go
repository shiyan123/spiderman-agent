package utils

import (
	"net"
)

func GetNodeIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic("net is error")
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	panic("net is error")
}
