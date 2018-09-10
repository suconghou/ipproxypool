package storage

import (
	"fmt"
	"ipproxypool/request"

	"os"
	"regexp"
)

var (
	IpReg = regexp.MustCompile(`((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))))`)

	ProxyReg *regexp.Regexp = regexp.MustCompile(`((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))):([1-9]?\d{0,4})`)
)

var urls = map[string]string{
	"https": "https://ipinfo.io",
	"http":  "http://ip.taobao.com/service/getIpInfo.php?ip=myip",
}

// 保存更新后的状态,同时判断响应是否是正确的
func ProxyStatus(item ProxyItem) ProxyItem {
	proxy := fmt.Sprintf("%s:%d", item.Ip, item.Port)
	proxyHttp := fmt.Sprintf("http://%s", proxy)
	// proxyHttpS:=fmt.Sprintf("https://%s",proxy)
	for _, url := range urls {
		res, err := request.GetByProxy(url, proxyHttp)
		if err != nil {
			item.Status = false
		} else {
			item.Status = true
			os.Stderr.Write(res)
		}
		GlobalProxyMap.Set(proxy, item)
	}
	return item
}
