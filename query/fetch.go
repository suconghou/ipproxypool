package query

import (
	"io"
	"ipproxypool/request"
	"ipproxypool/util"
	"net/http"
	"strings"
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

// New Fetch
func NewFetcher(config *FetchConfig) *Fetcher {
	var (
		method  = config.Method
		timeout = config.Timeout
		retry   = config.Retry
		limit   = config.Limit
		proxy   = config.Proxy
	)
	if !util.ValidMethod(method) {
		method = http.MethodGet
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
		client:  request.NewClient(timeout, proxy),
		urls:    config.Urls,
	}
}

// Do the fetch, get resp and parse
func (f Fetcher) Do(action string, query QueryConfig) (any, error) {
	respMap, err := f.doFetch()
	if err != nil {
		return nil, err
	}
	return process(respMap, action, query)
}

func (f Fetcher) doFetch() (map[string][]byte, error) {
	var tasks = []*request.Task{}
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
			client   *http.Client
			req, err = request.NewRequest(item.URL, method, headers, body)
			retry    = item.Retry
			limit    = item.Limit
		)
		if err != nil {
			return nil, err
		}
		if (item.Proxy != "" && item.Proxy != f.proxy) || (item.Timeout > 0 && item.Timeout < 120 && item.Timeout != f.timeout) {
			client = request.NewClient(item.Timeout, item.Proxy)
		} else {
			client = f.client
		}

		if retry <= 0 {
			retry = f.retry
		}
		if limit <= 0 {
			limit = f.limit
		}
		tasks = append(tasks, &request.Task{
			Client:    client,
			Request:   req,
			Retry:     retry,
			Limit:     limit,
			Transform: item.Transform,
		})
	}
	return request.GetTasksData(tasks)
}
