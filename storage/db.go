package storage

import (
	"ipproxypool/util"
	"time"
)

var (
	// ProxyItemListIn for push in
	ProxyItemListIn = make(chan ProxyItem, 100)
	// ProxyItemListGood for get good item
	ProxyItemListGood = make(chan ProxyItem, 100)
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
	for i := 0; i < 1; i++ {
		go func() {
			for {
				select {
				case item := <-ProxyItemListIn:
					now := time.Now()
					if proxyOpen(item) {
						item.Latency = uint16(time.Since(now).Seconds() * 1000)
						item.Status = true
						item.Succeed++
						ProxyItemListGood <- item
					} else {
						item.Status = false
						item.Failed++
						if item.Succeed > item.Failed && item.Failed < 10 {
							select {
							case ProxyItemListIn <- item:
							default:
							}
						}
					}
				case <-time.After(time.Second * time.Duration(30)):
					// 检查队列不足了,需要补充一些新的IP,去抓取吧
					util.Logger.Print("抓取的不够哦")
				}
			}
		}()
	}
}

// NewProxyItem create new proxy item
func NewProxyItem(ip string, port uint16) {
	ProxyItemListIn <- ProxyItem{IP: ip, Port: port, Latency: 0, Status: false, Succeed: 0, Failed: 0}
}
