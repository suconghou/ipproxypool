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
	Proxy   string
	Urls    []*itemcfg
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
	util.Logger.Print(data.Urls[0].URL)
	if len(data.Urls) < 1 {
		err := fmt.Errorf("at least one url")
		util.JSONPut(w, resp{-4, err.Error()})
		return err
	}
	var uitems = []*request.URLItem{}
	for _, item := range data.Urls {
		u, err := url.Parse(item.URL)
		if err != nil {
			util.JSONPut(w, resp{-5, err.Error()})
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

	request.New(cfg).Do()

	return nil
}
