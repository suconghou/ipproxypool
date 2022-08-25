package request

import (
	"io"
	"ipproxypool/encoding"
	"ipproxypool/util"
	"net/http"
	"net/url"
	"time"
)

var (
	defaultHeader = http.Header{
		"User-Agent": []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
		},
	}
)

type resItem struct {
	bytes []byte
	url   string
	err   error
}

type Task struct {
	Client    *http.Client
	Request   *http.Request
	Retry     int
	Limit     int
	Transform bool
}

func NewClient(timeout int, urlproxy string) *http.Client {
	var client = &http.Client{
		Timeout:   timeoutConfig(timeout),
		Transport: transportConfig(urlproxy),
	}
	return client
}

func NewRequest(targetURL string, method string, reqHeader http.Header, body io.Reader) (*http.Request, error) {
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

// 并发执行多个http调用
func GetTasksData(tasks []*Task) (map[string][]byte, error) {
	var (
		ch       = make(chan *resItem)
		response = make(map[string][]byte)
	)
	for _, u := range tasks {
		go func(taskItem *Task) {
			var (
				url       = taskItem.Request.URL.String()
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
			limitedr := http.MaxBytesReader(nil, resp.Body, int64(taskItem.Limit))
			var bytes []byte
			if taskItem.Transform {
				bytes, err = encoding.GbkReaderToUtf8(limitedr)
			} else {
				bytes, err = io.ReadAll(limitedr)
			}

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

// 最终执行http请求的方法，多个方法最终均复用此函数
func getTaskResponse(taskItem *Task) (*http.Response, error) {
	var (
		times = 0
		resp  *http.Response
		err   error
	)
	for ; times < taskItem.Retry; times++ {
		resp, err = taskItem.Client.Do(taskItem.Request)
		if err == nil {
			return resp, err
		}
		time.Sleep(time.Millisecond)
	}
	return resp, err
}

// GetResponse for large http response
func GetResponse(url *url.URL, method string, headers http.Header, body io.Reader, proxy string, timeout int, retry int) (*http.Response, error) {
	if !util.ValidMethod(method) {
		method = http.MethodGet
	}
	if timeout < 1 || timeout > 86400 {
		timeout = 86400
	}
	if retry < 1 || retry > 100 {
		retry = 3
	}
	var (
		client       = NewClient(timeout, proxy)
		request, err = NewRequest(url.String(), method, headers, body)
		limit        = 0
	)
	if err != nil {
		return nil, err
	}
	var taskItem = &Task{
		client,
		request,
		retry,
		limit,
		false,
	}
	return getTaskResponse(taskItem)
}

// GetResponseData 和 GetResponse 类似，但是仅GET请求，限制响应1MB以内
func GetResponseData(target string, timeout int, headers http.Header) ([]byte, error) {
	var (
		body  io.Reader
		proxy       = ""
		retry       = 2
		limit int64 = 1048576
	)
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	resp, err := GetResponse(u, http.MethodGet, headers, body, proxy, timeout, retry)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(http.MaxBytesReader(nil, resp.Body, limit))
}
