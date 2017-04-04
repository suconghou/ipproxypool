package storage

import (
	_ "database/sql"
	"fmt"
	_ "log"
	// _"github.com/mattn/go-sqlite3"
)

var (
	ProxyItemListIn  chan ProxyItem = make(chan ProxyItem, 100)
	ProxyItemListOut chan ProxyItem = make(chan ProxyItem, 100)
)

type ProxyItem struct {
	Ip      string
	Port    uint16
	Latency uint16
	Status  bool
	Http    bool
	Https   bool
}

// var db *sql.DB

func init() {
	// fmt.Println("init db")
	// currDb, err := sql.Open("sqlite3", "./foo.db")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(currDb)
	// db = currDb

}

func NewProxyItem(ip string, port uint16) ProxyItem {
	return ProxyItem{
		ip,
		port,
		0,
		false,
		false,
		false,
	}
}

// 爬虫爬入的检查状态,存入数据库
func saveIn() {

	for {
		item := <-ProxyItemListIn
		fmt.Println(item)
	}

}

// 补充队列,从数据库中取出,并检查状态,状态好的进入线程池
func getFromDb() {
	for {
		item, ok := GetOneProxyFromDb()
		if ok {
			item = ProxyStatus(item)
			ProxyItemListOut <- item
		} else {
			fmt.Println("db is empty")
		}
	}
}

func GetOneProxyFromDb() (ProxyItem, bool) {
	var item ProxyItem
	return item, true
}

// 从线程池中取出可用的用于服务
func GetOneProxy() ProxyItem {
	item := <-ProxyItemListOut
	return item
}

func FlushListOut() {
	length := len(ProxyItemListOut)
	var i int = 0
	for ; i < length; i++ {
		GetOneProxy()
	}
}
