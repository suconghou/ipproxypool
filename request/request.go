package request

import (
	"io"
	"io/ioutil"
	"ipproxypool/util"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	defaultHeader = http.Header{
		"User-Agent": []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36",
		},
	}
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
	Headers http.Header
	Method  string
	Timeout int
	Proxy   string
	Retry   int
	Limit   int
	Urls    []*URLItem
}

// URLItem define one fetch item config
type URLItem struct {
	URL     *url.URL
	Method  string
	Body    string
	Headers http.Header
	Timeout int
	Proxy   string
	Retry   int
	Limit   int
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

func newRequest(targetURL string, method string, reqHeader http.Header, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, targetURL, body)
	if err != nil {
		return req, err
	}
	if reqHeader != nil {
		req.Header = reqHeader
	} else {
		req.Header = defaultHeader
	}
	return req, nil
}

// New Fetch
func New(config *FetchConfig) *Fetcher {
	var (
		method  = config.Method
		timeout = config.Timeout
		retry   = config.Retry
		limit   = config.Limit
		proxy   = config.Proxy
	)
	if !util.ValidMethod(method) {
		method = "GET"
	}
	if timeout < 1 || timeout > 120 {
		timeout = 20
	}
	if retry < 1 || retry > 100 {
		retry = 3
	}
	if limit < 1 || limit > 8388608 {
		limit = 1048576
	}
	return &Fetcher{
		headers: config.Headers,
		method:  method,
		timeout: timeout,
		proxy:   proxy,
		retry:   retry,
		limit:   limit,
		client:  newClient(timeout, proxy),
		urls:    config.Urls,
	}
}

// Do the fetch, get resp and parse
func (f Fetcher) Do(action string, query QueryConfig) (interface{}, error) {
	respMap, err := f.doFetch()
	if err != nil {
		return nil, err
	}
	return process(respMap, action, query)
}

func (f Fetcher) doFetch() (map[string][]byte, error) {
	var tasks = []*task{}
	for _, item := range f.urls {
		var (
			method  = item.Method
			headers = item.Headers
			body    io.Reader
		)
		if !util.ValidMethod(method) {
			method = f.method
		}
		if headers == nil {
			headers = f.headers
		}
		if item.Body != "" {
			body = strings.NewReader(item.Body)
		}
		var (
			client       *http.Client
			request, err = newRequest(item.URL.String(), method, headers, body)
			retry        = item.Retry
			limit        = item.Limit
		)
		if err != nil {
			return nil, err
		}
		if (item.Proxy != "" && item.Proxy != f.proxy) || (item.Timeout > 0 && item.Timeout < 120 && item.Timeout != f.timeout) {
			client = newClient(item.Timeout, item.Proxy)
		} else {
			client = f.client
		}

		if retry <= 0 {
			retry = f.retry
		}
		if limit <= 0 {
			limit = f.limit
		}
		tasks = append(tasks, &task{
			client,
			request,
			retry,
			limit,
		})
	}
	return getTasksData(tasks)
}

// GetResponse for large http response
func GetResponse(url *url.URL, method string, headers http.Header, body io.Reader, proxy string, timeout int, retry int) (*http.Response, error) {
	if !util.ValidMethod(method) {
		method = "GET"
	}
	if timeout < 1 || timeout > 86400 {
		timeout = 86400
	}
	if retry < 1 || retry > 100 {
		retry = 3
	}
	var (
		client       = newClient(timeout, proxy)
		request, err = newRequest(url.String(), method, headers, body)
		limit        = 0
	)
	if err != nil {
		return nil, err
	}
	var taskItem = &task{
		client,
		request,
		retry,
		limit,
	}
	return getTaskResponse(taskItem)
}

// GetResponseData like GetResponse but only do GET request and return bytes for easy use
func GetResponseData(target string, timeout int, headers http.Header) ([]byte, error) {
	var (
		method = "GET"
		body   io.Reader
		proxy        = ""
		retry        = 2
		limit  int64 = 1048576
	)
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	resp, err := GetResponse(u, method, headers, body, proxy, timeout, retry)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(io.LimitReader(resp.Body, limit))
}

func getTaskResponse(taskItem *task) (*http.Response, error) {
	var (
		times = 0
		resp  *http.Response
		err   error
	)
	for ; times < taskItem.retry; times++ {
		resp, err = taskItem.client.Do(taskItem.request)
		if err == nil {
			return resp, err
		}
		time.Sleep(time.Millisecond)
	}
	return resp, err
}

func getTasksData(tasks []*task) (map[string][]byte, error) {
	var (
		ch       = make(chan *resItem)
		response = make(map[string][]byte)
	)
	for _, u := range tasks {
		go func(taskItem *task) {
			var (
				url       = taskItem.request.URL.String()
				resp, err = getTaskResponse(taskItem)
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
