package stream

import (
	"encoding/json"
	"fmt"
	"io"
	"ipproxypool/request"
	"ipproxypool/util"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	itemStatusWating   int8 = 1
	itemStatusStarted  int8 = 2
	itemStatusResolved int8 = 3
	itemStatusRejected int8 = 4
	itemStatusIgnored  int8 = 5

	itemModeIgnore    int8 = 1
	itemModeRename    int8 = 2
	itemModeOverWrite int8 = 3
	itemModeAppend    int8 = 4
)

// TaskItem desc the task
type TaskItem struct {
	URL     *url.URL    `json:"-"`
	Method  string      `json:"-"`
	Timeout int         `json:"-"`
	Proxy   string      `json:"-"`
	Headers http.Header `json:"-"`
	Body    string      `json:"-"`
	Retry   int         `json:"-"`
	Name    string      `json:"name"`
	Path    string      `json:"path"`
	Mode    int8        `json:"mode"`
	Status  int8        `json:"status"`
	Size    int64       `json:"size"`
	Start   int64       `json:"start"` // 任务开始时间
	End     int64       `json:"end"`   // 任务结束时间
}

// WorkerStatus intro work status
type WorkerStatus struct {
	Thread   int32                `json:"thread"`
	Runing   int32                `json:"runing"`
	Tasks    int                  `json:"tasks"`
	Pendings map[string]*TaskItem `json:"pendings"` // 正在处理的
	Items    []*TaskItem          `json:"items"`    // 已处理完毕的
}

// Worker do job
type Worker struct {
	thread   int32
	runing   int32
	receive  chan *TaskItem
	pendings map[string]*TaskItem
	items    []*TaskItem
	r        *sync.RWMutex
}

// Put taskitem to this worker
func (w *Worker) Put(t *TaskItem) {
	w.start()
	if t.Mode < itemModeIgnore || t.Mode > itemModeAppend {
		t.Mode = itemModeIgnore
	}
	t.Status = itemStatusWating
	w.receive <- t
}

// start this work use how many thread
func (w *Worker) start() {
	if atomic.LoadInt32(&w.runing) >= w.thread {
		return
	}
	go func() {
		defer func() {
			atomic.AddInt32(&w.runing, -1)
		}()
		for {
			select {
			case t := <-w.receive:
				t.before()
				// 忽略重复任务,根据name(url地址)字段判断
				w.r.Lock()
				// 如果进行中的任务没有这个则可以允许，后续根据Mode策略检测
				if _, exist := w.pendings[t.Name]; !exist {
					w.pendings[t.Name] = t
					w.r.Unlock()
					if err := t.after(t.start()); err != nil {
						util.Log.Print(err)
					}
					w.r.Lock()
					delete(w.pendings, t.Name)
					w.items = append(w.items, t)
					w.r.Unlock()
				} else {
					w.r.Unlock()
				}
			case <-time.After(time.Minute):
				return
			}
		}
	}()
	atomic.AddInt32(&w.runing, 1)
}

// GetStatus return worker status
func (w *Worker) GetStatus() ([]byte, error) {
	w.r.RLock()
	status := &WorkerStatus{
		Thread:   w.thread,
		Runing:   atomic.LoadInt32(&w.runing),
		Tasks:    len(w.receive),
		Items:    w.items,
		Pendings: w.pendings,
	}
	bs, err := json.Marshal(status)
	w.r.RUnlock()
	return bs, err
}

func (t *TaskItem) start() (int64, string, error) {
	t.Status = itemStatusStarted
	t.Start = time.Now().Unix()
	var (
		fpath = t.Path
		flag  = os.O_WRONLY | os.O_CREATE
		exist = util.FileExists(fpath)
	)
	if exist {
		switch t.Mode {
		case itemModeRename:
			fpath = fmt.Sprintf("%s.%d", t.Path, time.Now().Unix()) // 如果这个也重名，则文件被覆盖
		case itemModeAppend:
			flag = os.O_WRONLY | os.O_APPEND
		case itemModeOverWrite:
			flag = os.O_WRONLY | os.O_TRUNC
		case itemModeIgnore:
			return -1, fpath, nil
		}
	}
	err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
	if err != nil {
		return 0, fpath, err
	}
	file, err := os.OpenFile(fpath, flag, os.ModePerm)
	if err != nil {
		return 0, fpath, err
	}
	resp, err := request.GetResponse(t.URL, t.Method, t.Headers, strings.NewReader(t.Body), t.Proxy, t.Timeout, t.Retry)
	if err != nil {
		return 0, fpath, err
	}
	defer resp.Body.Close()
	defer file.Close()
	n, err := io.Copy(file, resp.Body)
	if err == nil {
		if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusIMUsed {
			// status not ok, we logger
			err = fmt.Errorf("%v %s", t.URL, resp.Status)
		}
	}
	return n, fpath, err
}

func (t *TaskItem) before() {
	t.Status = itemStatusWating
	if t.Name == "" {
		t.Name = resolveName(t.URL)
	}
	if t.Path == "" {
		t.Path = resolvePath(t.URL)
	}
}

func (t *TaskItem) after(n int64, fpath string, err error) error {
	if err != nil {
		t.Status = itemStatusRejected
	} else {
		if n == -1 {
			t.Status = itemStatusIgnored
		} else {
			t.Status = itemStatusResolved
		}
	}
	t.Path = fpath
	t.Size = n
	t.End = time.Now().Unix()
	return err
}

func resolveName(u *url.URL) string {
	return u.Host + u.Path
}

func resolvePath(u *url.URL) string {
	return filepath.Join(u.Host, strings.ReplaceAll(u.Path, "/", "_"))
}
