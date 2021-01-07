package util

import (
	"log"
	"net"
)

// 获取本机IP
var ServerIp string

func init() {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ServerIp = ipnet.IP.String()
				return
			}

		}
	}

	log.Fatal("No network found")
}
