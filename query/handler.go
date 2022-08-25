package query

import (
	"fmt"
	"ipproxypool/util"
	"net/http"
)

// 每个请求URL的配置，当某个项未配置时继承本任务的全局配置
type URLItem struct {
	URL       string
	Transform bool
	Method    string
	Body      string
	Headers   http.Header
	Timeout   int
	Proxy     string
	Retry     int
	Limit     int
}

// 本次任务的全局配置
type FetchConfig struct {
	Headers http.Header
	Method  string
	Timeout int
	Cache   int
	Proxy   string
	Retry   int
	Limit   int
	Urls    []*URLItem
}

type fetchcfg struct {
	*FetchConfig
	Query QueryConfig
}

func GoQuery(w http.ResponseWriter, r *http.Request, match []string) error {
	var data fetchcfg
	if err := util.Parse(w, r, &data); err != nil {
		return err
	}
	if len(data.Urls) < 1 {
		err := fmt.Errorf("at least one url")
		_, err = util.JSON(w, err.Error(), -4)
		return err
	}
	ret, err := NewFetcher(data.FetchConfig).Do(match[1], data.Query)
	if err != nil {
		_, err = util.JSON(w, err.Error(), -6)
		return err
	}
	if data.Cache > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", data.Cache))
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if bs, ok := ret.([]byte); ok {
		_, err = w.Write(bs)
	} else {
		_, err = util.JSONData(w, ret)
	}
	return err
}
