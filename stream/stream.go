package stream

import "sync"

// 下载器
// 负责下载图片,视频等资源文件,存储到指定地址

// NewWorker create new stream worker

var (
	// DefaultWorker for http api use
	DefaultWorker *Worker
)

func init() {
	DefaultWorker = NewWorker()
}

// NewWorker create new worker
func NewWorker() *Worker {
	return &Worker{
		thread:    20,
		runing:    0,
		receive:   make(chan *TaskItem, 100),
		statusMap: map[string]*TaskItem{},
		r:         &sync.RWMutex{},
	}
}
