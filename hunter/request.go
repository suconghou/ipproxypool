package hunter

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"time"
)

var reqHeader map[string]string = map[string]string{
	"Connection":                "keep-alive",
	"Cache-Control":             "max-age=0",
	"Upgrade-Insecure-Requests": "1",
	"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	// "Accept-Encoding":           "gzip, deflate, sdch",
	"Accept-Language": "zh-CN,zh;q=0.8,en;q=0.6",
}

func reqWithHeader(req *http.Request) *http.Request {
	for key, value := range reqHeader {
		req.Header.Add(key, value)
	}
	return req
}

func initDocument(url_addr string) (*goquery.Document, error) {
	request, err := http.NewRequest("GET", url_addr, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	res, err := client.Do(reqWithHeader(request))
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
