package request

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var reqHeader map[string]string = map[string]string{
	"Connection":    "keep-alive",
	"Cache-Control": "max-age=0",
	"User-Agent":    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
	"Accept":        "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	// "Accept-Encoding":           "gzip, deflate, sdch",
	"Accept-Language": "zh-CN,zh;q=0.8,en;q=0.6",
}

func reqWithHeader(req *http.Request, extra map[string]string) *http.Request {
	for key, value := range reqHeader {
		req.Header.Set(key, value)
	}
	for key, value := range extra {
		req.Header.Set(key, value)
	}
	return req
}

func Document(url_addr string) (*goquery.Document, error) {
	request, err := http.NewRequest("GET", url_addr, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	res, err := client.Do(reqWithHeader(request, nil))
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func HttpGet(url string, extra map[string]string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	res, err := client.Do(reqWithHeader(request, extra))
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	str, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return str, err
	}
	return str, nil
}

func GetByProxy(url_addr, proxy_addr string) ([]byte, error) {
	request, err := http.NewRequest("GET", url_addr, nil)
	if err != nil {
		return nil, err
	}
	proxy, err := url.Parse(proxy_addr)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}
	res, err := client.Do(reqWithHeader(request, nil))
	if err != nil {
		return nil, err
	}
	str, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return str, err
	}
	return str, nil
}
