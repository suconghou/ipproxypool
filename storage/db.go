package storage

import (
	"fmt"
	"sync"
)

var (
	ProxyItemListIn   chan ProxyItem = make(chan ProxyItem, 100)
	ProxyItemListGood chan ProxyItem = make(chan ProxyItem, 100)
	ProxyItemListBad  chan ProxyItem = make(chan ProxyItem, 1000)
)

type ProxyMap struct {
	Data map[string]ProxyItem
	Lock *sync.Mutex
}

type ProxyItem struct {
	Ip         string
	Port       uint16
	Latency    uint16
	Status     bool
	Http       bool
	Https      bool
	GoodStatus uint32
	BadStatus  uint32
}

var GlobalProxyMap = ProxyMap{
	map[string]ProxyItem{},
	new(sync.Mutex),
}

func init() {

}

func (proxyMap ProxyMap) Get(key string) (ProxyItem, bool) {
	proxyMap.Lock.Lock()
	val, ok := proxyMap.Data[key]
	proxyMap.Lock.Unlock()
	return val, ok
}

func (proxyMap ProxyMap) Set(key string, val ProxyItem) {
	proxyMap.Lock.Lock()
	proxyMap.Data[key] = val
	proxyMap.Lock.Unlock()
}

func (proxyMap ProxyMap) Len() int {
	return len(proxyMap.Data)
}

func NewProxyItem(ip string, port uint16) ProxyItem {
	return ProxyItem{Ip: ip, Port: port, Latency: 0, Status: false, Http: false, Https: false, GoodStatus: 0, BadStatus: 0}
}

// 爬虫爬入的入栈,与现有对比去重
func SaveProxyIn(ipList []ProxyItem) {
	for _, item := range ipList {
		proxy := fmt.Sprintf("%s:%d", item.Ip, item.Port)
		if _, ok := GlobalProxyMap.Get(proxy); !ok {
			GlobalProxyMap.Set(proxy, item)
			ProxyItemListIn <- item
		}
	}
}

// 从线程池中取出可用的用于服务
func GetOneProxy() ProxyItem {
	item := <-ProxyItemListGood
	return item
}
