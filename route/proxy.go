package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"ipproxypool/storage"
	"ipproxypool/util"
	"net/http"
	"time"
)

type proxyItem struct {
	IP   string
	Port uint16
}

func proxyone(w http.ResponseWriter, r *http.Request, match []string) error {
	item := <-storage.ProxyItemListGood
	_, err := util.JSONPut(w, item)
	go func() {
		select {
		case storage.ProxyItemListIn <- item:
		case <-time.After(time.Minute):
		}
	}()
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
		util.JSONPut(w, resp{-2, err.Error(), nil})
		return err
	}
	var data []proxyItem
	err = json.Unmarshal(bs, &data)
	if err != nil {
		util.JSONPut(w, resp{-3, err.Error(), nil})
		return err
	}
	for _, v := range data {
		storage.NewProxyItem(v.IP, v.Port)
	}
	_, err = util.JSONPut(w, resp{0, "ok", nil})
	return err
}

func proxyinfo(w http.ResponseWriter, r *http.Request, match []string) error {
	var data = map[string]interface{}{
		"queued":   len(storage.ProxyItemListGood),
		"checking": len(storage.ProxyItemListIn),
	}
	_, err := util.JSONPut(w, data)
	return err
}
