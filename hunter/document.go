package hunter

import (
	"ipproxypool/request"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func document(u string) (*goquery.Document, error) {
	var (
		method  = "GET"
		headers = http.Header{}
		data    = ""
		proxy   = ""
		timeout = 10
		retry   = 3
	)
	url, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	resp, err := request.GetResponse(url, method, headers, data, proxy, timeout, retry)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	return doc, err
}
