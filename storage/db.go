package storage

import (
	"sync/atomic"
	"time"
)

var (
	// ProxyItemListIn for push in
	ProxyItemListIn = make(chan ProxyItem, 100)
	// ProxyItemListGood for get good item
	ProxyItemListGood = make(chan ProxyItem, 100)

	// Runing current runing
	Runing int32
	// Thread max thread num
	Thread int32 = 10
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

func start() {
	if atomic.LoadInt32(&Runing) >= Thread {
		return
	}
	go func() {
		defer func() {
			atomic.AddInt32(&Runing, -1)
		}()
		for {
			select {
			case item := <-ProxyItemListIn:
				now := time.Now()
				if proxyOpen(item) {
					item.Latency = uint16(time.Since(now).Seconds() * 1000)
					item.Status = true
					item.Succeed++
					select {
					case ProxyItemListGood <- item:
					case <-time.After(time.Minute):
						return
					}

				} else {
					item.Status = false
					item.Failed++
					if item.Succeed > item.Failed && item.Failed < 10 {
						select {
						case ProxyItemListIn <- item:
						case <-time.After(time.Minute):
							return
						}
					}
				}
			case <-time.After(time.Minute):
				// 没什么可检查的,需要补充一些新的IP,去抓取吧
				return
			}
		}

	}()
	atomic.AddInt32(&Runing, 1)

}

// NewProxyItem create new proxy item
func NewProxyItem(ip string, port uint16) {
	ProxyItemListIn <- ProxyItem{IP: ip, Port: port, Latency: 0, Status: false, Succeed: 0, Failed: 0}
	start()
}
