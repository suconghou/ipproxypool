package stream

import (
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
	Status  int8        `json:"status"`
	Size    int64       `json:"size"`
}

// WorkerStatus intro work status
type WorkerStatus struct {
	Thread int32                `json:"thread"`
	Runing int32                `json:"runing"`
	Tasks  int                  `json:"tasks"`
	Items  map[string]*TaskItem `json:"items"`
}

// Worker do job
type Worker struct {
	thread    int32
	runing    int32
	receive   chan *TaskItem
	statusMap map[string]*TaskItem
	r         *sync.RWMutex
}

// Put taskitem to this worker
func (w *Worker) Put(t *TaskItem) error {
	w.start()
	t.Status = itemStatusWating
	w.receive <- t
	return nil
}

// start this work use how many thread
func (w *Worker) start() {
	if atomic.LoadInt32(&w.runing) < w.thread {
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
					if _, ok := w.statusMap[t.Name]; !ok {
						w.statusMap[t.Name] = t
						w.r.Unlock()
						if err := t.after(t.start()); err != nil {
							util.Logger.Print(err)
						}
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
}

// GetStatus return worker status
func (w *Worker) GetStatus() *WorkerStatus {
	w.r.RLock()
	status := &WorkerStatus{
		Thread: w.thread,
		Runing: atomic.LoadInt32(&w.runing),
		Tasks:  len(w.receive),
		Items:  w.statusMap,
	}
	w.r.RUnlock()
	return status
}

func (t *TaskItem) start() (int64, string, error) {
	t.Status = itemStatusStarted
	var savepath string
	if _, err := os.Stat(t.Path); os.IsNotExist(err) {
		savepath = t.Path
	} else {
		savepath = fmt.Sprintf("%s.%d", t.Path, time.Now().Unix())
	}
	err := os.MkdirAll(filepath.Dir(savepath), os.ModePerm)
	if err != nil {
		return 0, "", err
	}
	file, err := os.OpenFile(savepath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return 0, "", err
	}
	resp, err := request.GetResponse(t.URL, t.Method, t.Headers, t.Body, t.Proxy, t.Timeout, t.Retry)
	if err != nil {
		return 0, "", err
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
	return n, savepath, err
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

func (t *TaskItem) after(n int64, savepath string, err error) error {
	if err != nil {
		t.Status = itemStatusRejected
	} else {
		t.Status = itemStatusResolved
	}
	t.Path = savepath
	t.Size = n
	return err
}

func resolveName(u *url.URL) string {
	return u.Host + u.Path
}

func resolvePath(u *url.URL) string {
	return filepath.Join(u.Host, strings.ReplaceAll(u.Path, "/", "_"))
}
