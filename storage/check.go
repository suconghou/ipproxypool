package storage

import (
	"fmt"
	"net"
	"time"
)

// proxyOpen test proxy is reachable
func proxyOpen(item ProxyItem) bool {
	var ipPort = fmt.Sprintf("%s:%d", item.IP, item.Port)
	_, err := net.DialTimeout("tcp", ipPort, time.Second)
	return err == nil
}
