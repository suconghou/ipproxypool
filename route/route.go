package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"ipproxypool/storage"
	"ipproxypool/util"
	"net/http"
	"regexp"
)

// 路由定义
type routeInfo struct {
	Reg     *regexp.Regexp
	Handler func(http.ResponseWriter, *http.Request, []string) error
}

// Route export route list
var Route = []routeInfo{
	{regexp.MustCompile(`^/api/proxy/getone$`), proxygetone},
	{regexp.MustCompile(`^/api/proxy/add$`), proxyadd},
	{regexp.MustCompile(`^/api/task/info$`), taskinfo},
	{regexp.MustCompile(`^/api/task/add$`), taskadd},
}

type resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type proxyItem struct {
	IP   string
	Port uint16
}

func proxygetone(w http.ResponseWriter, r *http.Request, match []string) error {
	item := storage.GetOneProxy()
	defer storage.SaveProxyIn([]storage.ProxyItem{item})
	_, err := util.JSONPut(w, item)
	return err
}

func proxyadd(w http.ResponseWriter, r *http.Request, match []string) error {
	bs, err := ioutil.ReadAll(http.MaxBytesReader(w, r.Body, 8192))
	if err == nil {
		if len(bs) <= 4 {
			err = fmt.Errorf("bad request")
		}
	}
	if err != nil {
		util.JSONPut(w, resp{-2, err.Error()})
		return err
	}
	var data []proxyItem
	err = json.Unmarshal(bs, &data)
	if err != nil {
		util.JSONPut(w, resp{-3, err.Error()})
		return err
	}
	var items = []storage.ProxyItem{}
	for _, v := range data {
		items = append(items, storage.NewProxyItem(v.IP, v.Port))
	}
	storage.SaveProxyIn(items)
	_, err = util.JSONPut(w, resp{0, "ok"})
	return err
}
