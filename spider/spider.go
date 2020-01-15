package spider

import (
	"bytes"
	"fmt"
	"ipproxypool/request"
	"ipproxypool/storage"
	"ipproxypool/util"
	"regexp"
	"strconv"
	"time"
)

type ipport struct {
	ip   string
	port uint16
}

var (
	ipReg    = regexp.MustCompile(`((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))))`)
	proxyReg = regexp.MustCompile(`(((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?):[1-9]?\d{0,4})`)

	regIP                = regexp.MustCompile(`(?:(?:[0,1]?\d?\d|2[0-4]\d|25[0-5])\.){3}(?:[0,1]?\d?\d|2[0-4]\d|25[0-5])`)
	regProxy             = regexp.MustCompile(`(?:(?:[0,1]?\d?\d|2[0-4]\d|25[0-5])\.){3}(?:[0,1]?\d?\d|2[0-4]\d|25[0-5]):\d{0,5}`)
	regProxyWithoutColon = regexp.MustCompile(`(?:(?:[0,1]?\d?\d|2[0-4]\d|25[0-5])\.){3}(?:[0,1]?\d?\d|2[0-4]\d|25[0-5]) \d{0,5}`)

	ipProxy  = make(chan *ipport, 100)
	ipSet    = map[string]uint32{} // 用于判断是否已添加过,已添加过的,后面程序会自动校验是否可用,我们这里不向他们提交重复数据
	pageData = make(chan *[]byte, 200)
)

// Start 启动后,期待ipProxy来数据,若始终无数据,则超时发起新请求,抓取IP
// 抓取的IP若数量足够,将被ipProxy阻塞,最终由NewProxyItem消费,NewProxyItem决定最佳状态
func Start() {
	for {
		select {
		case item := <-ipProxy:
			var ipstr = fmt.Sprintf("%s:%d", item.ip, item.port)
			if v, ok := ipSet[ipstr]; ok {
				ipSet[ipstr] = v + 1
			} else {
				ipSet[ipstr] = 1
				storage.NewProxyItem(item.ip, item.port)
			}
		case page := <-pageData:
			parse(page)
		case <-time.After(time.Second):
			go func() {
				if err := ip66(); err != nil {
					util.Log.Print(err)
				}
			}()
			time.Sleep(time.Second * 5)
		}
	}
}

// 在网页中用正则查找IP,尝试当做代理去处理
func parse(page *[]byte) error {
	var matches = proxyReg.FindAll(*page, -1)
	if matches != nil {
		for _, v := range matches {
			var arr = bytes.Split(v, []byte(":"))
			var ip = string(arr[0])
			var port, err = strconv.Atoi(string(arr[1]))
			if err != nil {
				return err
			}
			ipProxy <- &ipport{ip, uint16(port)}
		}
	}
	return nil
}

func ip66() error {
	var url = "http://www.66ip.cn/mo.php?sxb=&tqsl=1000&port=&export=&ktip=&sxa=&submit=%CC%E1++%C8%A1&textarea=http%3A%2F%2Fwww.66ip.cn%2F%3Fsxb%3D%26tqsl%3D10%26ports%255B%255D2%3D%26ktip%3D%26sxa%3D%26radio%3Dradio%26submit%3D%25CC%25E1%2B%2B%25C8%25A1"
	data, err := request.GetResponseData(url, 6, nil)
	if err != nil {
		return err
	}
	pageData <- &data
	return nil
}
