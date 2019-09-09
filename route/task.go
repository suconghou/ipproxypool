package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"ipproxypool/stream"
	"ipproxypool/util"
	"net/http"
	"net/url"
)

type taskItem struct {
	URL     string
	Method  string
	Timeout int
	Proxy   string
	Headers http.Header
	Body    string
	Retry   int
	Name    string
	Path    string
}

// 下载队列信息
func taskinfo(w http.ResponseWriter, r *http.Request, match []string) error {
	var ret = stream.DefaultWorker.GetStatus()
	_, err := util.JSONPut(w, ret)
	return err
}

// 添加新下载任务
func taskadd(w http.ResponseWriter, r *http.Request, match []string) error {
	bs, err := ioutil.ReadAll(http.MaxBytesReader(w, r.Body, 8192))
	if err == nil {
		if len(bs) <= 64 {
			err = fmt.Errorf("bad request")
		}
	}
	if err != nil {
		util.JSONPut(w, resp{-2, err.Error()})
		return err
	}
	var data taskItem
	err = json.Unmarshal(bs, &data)
	if err != nil {
		util.JSONPut(w, resp{-3, err.Error()})
		return err
	}
	urlinfo, err := url.Parse(data.URL)
	if err != nil {
		util.JSONPut(w, resp{-4, err.Error()})
		return err
	}
	var task = &stream.TaskItem{
		URL:     urlinfo,
		Method:  data.Method,
		Timeout: data.Timeout,
		Proxy:   data.Proxy,
		Headers: data.Headers,
		Body:    data.Body,
		Retry:   data.Retry,
		Name:    data.Name,
		Path:    data.Path,
	}
	return stream.DefaultWorker.Add(task)
}