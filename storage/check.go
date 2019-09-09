package storage

import (
	"fmt"
	"net"
	"time"
)

var urls = map[string]string{
	"https": "https://ipinfo.io",
	"http":  "http://ip.taobao.com/service/getIpInfo.php?ip=myip",
}

// proxyOpen test proxy is reachable
func proxyOpen(item ProxyItem) bool {
	var ipPort = fmt.Sprintf("%s:%d", item.IP, item.Port)
	_, err := net.DialTimeout("tcp", ipPort, time.Second)
	if err != nil {
		return false
	}
	return true
}
