package storage

import (
	"fmt"
	"ipproxypool/util"
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
	_, err := net.DialTimeout("tcp", ipPort, time.Second*10)
	util.Logger.Print(ipPort, err)
	if err != nil {
		return false
	}
	return true
}
