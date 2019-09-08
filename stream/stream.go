package stream

// 下载器
// 负责下载图片,视频等资源文件,存储到指定地址
// 包含before 和 after 处理. 并且确保了一路并发
// 存内存操作,提供API,支持m3u8
// 包含取消操作

// NewWorker create new stream worker

var (
	// DefaultWorker for http api use
	DefaultWorker *Worker
)

func init() {
	DefaultWorker = NewWorker()
	DefaultWorker.Start(4)
}

// NewWorker create new worker
func NewWorker() *Worker {
	return &Worker{
		receive:   make(chan *TaskItem, 4),
		statusMap: map[string]*TaskItem{},
	}
}
