package web

import (
	"fmt"
	"io/ioutil"
	"ipproxypool/storage"
	_ "ipproxypool/util"
	"net/http"
	"regexp"
	"strconv"
)

type SimpleProxy struct {
	Ip   string
	Port uint16
}

// 路由定义
type routeInfo struct {
	Reg     *regexp.Regexp
	Handler func(http.ResponseWriter, *http.Request, []string)
}

// 路由添加
var RoutePath = []routeInfo{
	{regexp.MustCompile(`^/api/add/(.*)$`), apiAdd},
	{regexp.MustCompile(`^/api/get/(.*)$`), apiGet},
}

func Start() {

}

func apiAdd(w http.ResponseWriter, r *http.Request, match []string) {
	fmt.Println(match)
	result, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}
	lists, ok := parseProxyIp(string(result))
	if ok {
		// 与数据库对比去重,然后检查并存入
		fmt.Println(lists)
	} else {

	}

}

func apiGet(w http.ResponseWriter, r *http.Request, match []string) {
	fmt.Println(match)
}

func parseProxyIp(str string) ([]SimpleProxy, bool) {
	var proxyList []SimpleProxy
	if storage.ProxyReg.MatchString(str) {
		matches := storage.ProxyReg.FindAllStringSubmatch(str, -1)
		for _, item := range matches {
			port, err := strconv.Atoi(item[8])
			if err != nil {
				continue
			}
			var proxy SimpleProxy = SimpleProxy{item[1], uint16(port)}
			proxyList = append(proxyList, proxy)
		}
		return proxyList, true
	} else {
		return proxyList, false
	}
}
