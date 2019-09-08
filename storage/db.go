package storage

import "time"

var (
	proxyItemListIn   = make(chan ProxyItem, 100)
	proxyItemListGood = make(chan ProxyItem, 100)
	proxyItemListBad  = make(chan ProxyItem, 100)
)

// ProxyItem is one proxy
type ProxyItem struct {
	IP      string `json:"ip"`
	Port    uint16 `json:"port"`
	Latency uint16 `json:"latency"`
	Status  bool   `json:"status"`
	Succeed uint32 `json:"succeed"`
	Failed  uint32 `json:"failed"`
}

func init() {
	go func() {
		for {
			item := <-proxyItemListIn
			now := time.Now()
			if proxyOpen(item) {
				item.Latency = uint16(time.Since(now).Seconds() * 1000)
				item.Status = true
				item.Succeed++
				proxyItemListGood <- item
			} else {
				item.Status = false
				item.Failed++
				proxyItemListBad <- item
			}
		}
	}()
}

// NewProxyItem create new proxy item
func NewProxyItem(ip string, port uint16) ProxyItem {
	return ProxyItem{IP: ip, Port: port, Latency: 0, Status: false, Succeed: 0, Failed: 0}
}

// SaveProxyIn add new proxy to check
func SaveProxyIn(ipList []ProxyItem) {
	for _, item := range ipList {
		proxyItemListIn <- item
	}
}

// GetOneProxy return one good proxy
func GetOneProxy() ProxyItem {
	item := <-proxyItemListGood
	return item
}
