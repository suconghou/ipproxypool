package route

import (
	"fmt"
	"ipproxypool/request"
	"ipproxypool/util"
	"net/http"
	"net/url"
)

type fetchcfg struct {
	Headers http.Header
	Method  string
	Timeout int
	Cache   int
	Proxy   string
	Urls    []*itemcfg
	Query   request.QueryConfig
}

type itemcfg struct {
	URL     string
	Method  string
	Body    string
	Headers http.Header
	Timeout int
	Proxy   string
	Retry   int
	Limit   int
}

func fetchurl(w http.ResponseWriter, r *http.Request, match []string) error {
	var data fetchcfg
	if err := parse(w, r, &data); err != nil {
		return err
	}
	if len(data.Urls) < 1 {
		err := fmt.Errorf("at least one url")
		util.JSONPut(w, resp{-4, err.Error(), nil})
		return err
	}
	var uitems = []*request.URLItem{}
	for _, item := range data.Urls {
		u, err := url.Parse(item.URL)
		if err != nil {
			util.JSONPut(w, resp{-5, err.Error(), nil})
			return err
		}
		uitems = append(uitems, &request.URLItem{
			URL:     u,
			Method:  item.Method,
			Body:    item.Body,
			Headers: item.Headers,
			Timeout: item.Timeout,
			Proxy:   item.Proxy,
			Retry:   item.Retry,
			Limit:   item.Limit,
		})
	}
	var cfg = &request.FetchConfig{
		Headers: data.Headers,
		Method:  data.Method,
		Timeout: data.Timeout,
		Proxy:   data.Proxy,
		Urls:    uitems,
	}
	ret, err := request.New(cfg).Do(match[1], data.Query)
	if err != nil {
		util.JSONPut(w, resp{-6, err.Error(), nil})
		return err
	}
	if data.Cache > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", data.Cache))
	}
	if bs, ok := ret.([]byte); ok {
		_, err = w.Write(bs)
	} else {
		_, err = util.JSONPut(w, resp{0, "ok", ret})
	}
	return err
}
