package request

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Fetcher http fetch
type Fetcher struct {
	headers http.Header
	method  string
	timeout int
	proxy   string
	retry   int
	limit   int
	client  *http.Client
	urls    []*URLItem
}

// FetchConfig http request config for post payload parse
type FetchConfig struct {
	headers http.Header
	method  string
	timeout int
	proxy   string
	urls    []*URLItem
}

// URLItem define one fetch item config
type URLItem struct {
	url     *url.URL
	method  string
	body    string
	headers http.Header
	timeout int
	proxy   string
	retry   int
	limit   int
}

type resItem struct {
	bytes []byte
	url   string
	err   error
}

type task struct {
	client  *http.Client
	request *http.Request
	retry   int
	limit   int
}

func newClient(timeout int, urlproxy string) *http.Client {
	var client = &http.Client{
		Timeout:   timeoutConfig(timeout),
		Transport: transportConfig(urlproxy),
	}
	return client
}

func newRequest(targetURL string, method string, reqHeader http.Header, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, targetURL, body)
	if err != nil {
		panic(err)
	}
	if reqHeader != nil {
		req.Header = reqHeader
	}
	return req
}

// New Fetch
func New(config *FetchConfig) *Fetcher {
	var (
		method  = config.method
		timeout = config.timeout
	)
	if !isValidMethod(method) {
		method = "GET"
	}
	if timeout < 1 || timeout > 120 {
		timeout = 20
	}
	return &Fetcher{
		headers: config.headers,
		method:  method,
		timeout: timeout,
		proxy:   config.proxy,
		client:  newClient(config.timeout, config.proxy),
		urls:    config.urls,
	}
}

// Do the fetch, get resp and parse
func (f Fetcher) Do() {
	resp, err := f.doFetch()
	fmt.Println(resp, err)
}

func (f Fetcher) doFetch() (map[string][]byte, error) {
	var tasks = []*task{}
	for i, item := range f.urls {
		var (
			method  = item.method
			headers = item.headers
			body    io.Reader
		)
		if !isValidMethod(method) {
			method = f.method
		}
		if headers == nil {
			headers = f.headers
		}
		if item.body != "" {
			body = strings.NewReader(item.body)
		}
		var (
			client  *http.Client
			request = newRequest(item.url.String(), method, headers, body)
			retry   = item.retry
			limit   = item.limit
		)
		if (item.proxy != "" && item.proxy != f.proxy) || (item.timeout > 0 && item.timeout < 120 && item.timeout != f.timeout) {
			client = newClient(item.timeout, item.proxy)
		} else {
			client = f.client
		}

		if retry <= 0 {
			retry = f.retry
		}
		if limit <= 0 {
			limit = f.limit
		}
		tasks[i] = &task{
			client,
			request,
			retry,
			limit,
		}
	}
	return getURLBody(tasks)
}

// GetResponse for large http response
func GetResponse(url *url.URL, method string, headers http.Header, data string, proxy string, timeout int, retry int) (*http.Response, error) {
	var (
		body io.Reader
	)
	if !isValidMethod(method) {
		method = "GET"
	}
	if headers == nil {
		headers = http.Header{}
	}
	if data != "" {
		body = strings.NewReader(data)
	}
	if timeout < 1 || timeout > 7200 {
		timeout = 7200
	}
	if retry < 1 || retry > 100 {
		retry = 10
	}
	var (
		client  = newClient(timeout, proxy)
		request = newRequest(url.String(), method, headers, body)
		limit   = 0
	)
	var taskItem = &task{
		client,
		request,
		retry,
		limit,
	}
	return getURLResponse(taskItem)
}

func getURLResponse(taskItem *task) (*http.Response, error) {
	var (
		times = 0
		resp  *http.Response
		err   error
	)
	for ; times < taskItem.retry; times++ {
		resp, err = taskItem.client.Do(taskItem.request)
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond)
	}
	return resp, err
}

func getURLBody(tasks []*task) (map[string][]byte, error) {
	var (
		ch       = make(chan *resItem)
		response = make(map[string][]byte)
	)
	for _, u := range tasks {
		go func(taskItem *task) {
			var (
				url       = taskItem.request.URL.String()
				resp, err = getURLResponse(taskItem)
			)
			if err != nil {
				ch <- &resItem{
					nil,
					url,
					err,
				}
				return
			}
			defer resp.Body.Close()
			bytes, err := ioutil.ReadAll(io.LimitReader(resp.Body, int64(taskItem.limit)))
			ch <- &resItem{
				bytes,
				url,
				err,
			}
		}(u)
	}
	for range tasks {
		item := <-ch
		if item.err != nil {
			return response, item.err
		}
		response[item.url] = item.bytes
	}
	return response, nil
}

func isValidMethod(m string) bool {
	if m == "GET" || m == "POST" || m == "PUT" || m == "DELETE" {
		return true
	}
	return false
}
