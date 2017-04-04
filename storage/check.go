package storage

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

var IpReg *regexp.Regexp = regexp.MustCompile(`((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))))`)

var ProxyReg *regexp.Regexp = regexp.MustCompile(`((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))):([1-9]?\d{0,4})`)

var urls map[string]string = map[string]string{
	"https": "https://ipinfo.io",
	"http":  "http://ip.taobao.com/service/getIpInfo.php?ip=myip",
}

func GetByProxy(url_addr, proxy_addr string) ([]byte, error) {
	request, err := http.NewRequest("GET", url_addr, nil)
	if err != nil {
		return nil, err
	}
	proxy, err := url.Parse(proxy_addr)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	str, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return str, err
	}
	return str, nil
}

func ProxyStatus(item ProxyItem) ProxyItem {
	proxy := fmt.Sprintf("%s:%s", item.Ip, item.Port)
	for _, url := range urls {
		res, err := GetByProxy(url, proxy)
		if err != nil {
			item.Status = false
		} else {
			item.Status = true
			fmt.Println(res)
		}

	}
	return item
}

func FindAllProxy(str string) []ProxyItem {
	var ipList []ProxyItem
	matches := ProxyReg.FindAllStringSubmatch(str, -1)
	for _, item := range matches {
		portInt, err := strconv.Atoi(item[8])
		if err != nil {
			continue
		}
		ipList = append(ipList, NewProxyItem(item[1], uint16(portInt)))
	}
	return ipList
}
