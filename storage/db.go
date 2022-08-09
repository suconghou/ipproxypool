package storage

import (
	"encoding/json"
	"ipproxypool/util"
	"os"
	"sync/atomic"
	"time"
)

const (
	proxyfile = "proxylist.json"
)

var (
	// ProxyItemListIn 代表检查队列, 如果检测未通过,并且曾经是好的,并且小于3次检测不通过,则会稍后再次检验
	ProxyItemListIn = make(chan ProxyItem, 100)
	// ProxyItemListGood 代表已检查通过的,可能是socks5或https代理,不保留纯http代理
	ProxyItemListGood = make(chan ProxyItem, 100)

	// Runing current runing
	Runing int32
	// Thread max thread num
	Thread int32 = 10

	proxyList = []ProxyItem{}
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
					if item.Succeed > item.Failed && item.Failed < 3 {
						select {
						case ProxyItemListIn <- item:
						case <-time.After(time.Minute):
							return
						}
					}
				}
			case <-time.After(time.Minute):
				// 没什么可检查的,需要补充一些新的IP,去抓取吧
				util.Log.Print("需要抓取新IP")
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

func dump() error {
	var bs, err = json.Marshal(proxyList)
	if err != nil {
		return err
	}
	return os.WriteFile(proxyfile, bs, 0666)
}

func load() error {
	var bs, err = os.ReadFile(proxyfile)
	if err != nil {
		return err
	}
	var data []ProxyItem
	err = json.Unmarshal(bs, &data)
	if err != nil {
		return err
	}
	proxyList = data
	return nil
}
