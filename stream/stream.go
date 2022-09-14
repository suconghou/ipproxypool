package stream

import "sync"

// 下载器
// 负责下载图片,视频等资源文件,存储到指定地址

// NewWorker create new stream worker

var (
	// DefaultWorker for http api use
	DefaultWorker = NewWorker()
)

// NewWorker create new worker
func NewWorker() *Worker {
	return &Worker{
		thread:   20,
		receive:  make(chan *TaskItem, 100),
		pendings: map[string]*TaskItem{},
		items:    []*TaskItem{},
		r:        &sync.RWMutex{},
	}
}
