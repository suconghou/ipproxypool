package route

import (
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
	Mode    int8
}

// 下载队列信息
func taskinfo(w http.ResponseWriter, r *http.Request, match []string) error {
	var ret = stream.DefaultWorker.GetStatus()
	_, err := util.JSONPut(w, ret)
	return err
}

// 添加新下载任务
func taskadd(w http.ResponseWriter, r *http.Request, match []string) error {
	var (
		data     []taskItem
		taskList []*stream.TaskItem
	)
	if err := parse(w, r, &data); err != nil {
		return err
	}
	for _, item := range data {
		urlinfo, err := url.Parse(item.URL)
		if err != nil {
			util.JSONPut(w, resp{-4, err.Error(), nil})
			return err
		}
		task := &stream.TaskItem{
			URL:     urlinfo,
			Method:  item.Method,
			Timeout: item.Timeout,
			Proxy:   item.Proxy,
			Headers: item.Headers,
			Body:    item.Body,
			Retry:   item.Retry,
			Name:    item.Name,
			Path:    item.Path,
			Mode:    item.Mode,
		}
		taskList = append(taskList, task)
	}
	for _, task := range taskList {
		stream.DefaultWorker.Put(task)
	}
	_, err := util.JSONPut(w, resp{0, "ok", nil})
	return err
}
