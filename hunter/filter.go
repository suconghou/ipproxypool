package hunter

import (
	"bytes"
	"fmt"
	"ipproxypool/request"
	"ipproxypool/storage"
	"ipproxypool/util"
	"os"
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

	ipProxy = make(chan *ipport, 100)
	ipSet   = map[string]uint32{}
)

func init() {
	var fetch = os.Getenv("PROXY_FETCH")
	if fetch == "" {
		return
	}
	go func() {
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
			case <-time.After(time.Second):
				go func() {
					if err := ip66(); err != nil {
						util.Logger.Print(err)
					}
				}()
				time.Sleep(time.Second * 5)
			}
		}
	}()
}

func ip66() error {
	var url = "http://www.66ip.cn/mo.php?sxb=&tqsl=1000&port=&export=&ktip=&sxa=&submit=%CC%E1++%C8%A1&textarea=http%3A%2F%2Fwww.66ip.cn%2F%3Fsxb%3D%26tqsl%3D10%26ports%255B%255D2%3D%26ktip%3D%26sxa%3D%26radio%3Dradio%26submit%3D%25CC%25E1%2B%2B%25C8%25A1"
	data, err := request.GetResponseData(url, 6, nil)
	if err != nil {
		return err
	}
	var matches = proxyReg.FindAll(data, -1)
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
