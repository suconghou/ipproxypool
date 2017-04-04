package hunter

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"ipproxypool/request"
	"ipproxypool/storage"
	"ipproxypool/util"
	"strconv"
	"strings"
	"time"
)

func StartHunter() {
	go func() {
		for {
			ipList, err := iP181()
			if err != nil {
				util.Debug(fmt.Sprintf("task has an error %s", err))
			} else {
				storage.SaveProxyIn(ipList)
			}
			ipList, err = xici()
			if err != nil {
				util.Debug(fmt.Sprintf("task has an error %s", err))
			} else {
				storage.SaveProxyIn(ipList)
			}
			ipList, err = ydl()
			if err != nil {
				util.Debug(fmt.Sprintf("task has an error %s", err))
			} else {
				storage.SaveProxyIn(ipList)
			}
			ipList, err = xdl()
			if err != nil {
				util.Debug(fmt.Sprintf("task has an error %s", err))
			} else {
				storage.SaveProxyIn(ipList)
			}
			time.Sleep(time.Second * 120)
		}
	}()

	go func() {
		for {
			item := <-storage.ProxyItemListIn
			item = storage.ProxyStatus(item)
			if item.Status {
				item.GoodStatus = item.GoodStatus + 1
				storage.ProxyItemListGood <- item
			} else {
				item.BadStatus = item.BadStatus + 1
				storage.ProxyItemListBad <- item
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * 60)
			item := <-storage.ProxyItemListBad
			item = storage.ProxyStatus(item)
			if item.Status {
				item.GoodStatus = item.GoodStatus + 1
				storage.ProxyItemListGood <- item
			} else if (item.GoodStatus < 1 && item.BadStatus > 10) || (item.BadStatus > 20) {
				item.BadStatus = item.BadStatus + 1
				storage.ProxyItemListBad <- item
			} else {
				util.Debug(fmt.Sprintf("drop bad proxy %s:%d", item.Ip, item.Port))
			}
		}
	}()

}

func iP181() ([]storage.ProxyItem, error) {
	var url string = "http://www.ip181.com/"
	var ipList []storage.ProxyItem
	doc, err := request.Document(url)
	if err != nil {
		return ipList, err
	}
	doc.Find(".panel-info .panel-body table  tbody  tr").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			ipAddr, err1 := s.Find("td:nth-child(1)").Html()
			portStr, err2 := s.Find("td:nth-child(2)").Html()
			if err1 != nil || err2 != nil {
				util.Debug(fmt.Sprintf("iP181 parse error %s %s", err1, err2))
			} else if ipAddr == "" || portStr == "" {
				util.Debug(fmt.Sprintf("iP181 matched empty %s:%s", ipAddr, portStr))
			} else {
				portInt, err3 := strconv.Atoi(portStr)
				if err3 != nil || !storage.IpReg.MatchString(ipAddr) {
					util.Debug(fmt.Sprintf("iP181 ip %s or port parse error %s", ipAddr, err3))
				} else {
					ipList = append(ipList, storage.NewProxyItem(ipAddr, uint16(portInt)))
				}
			}
		}
	})
	util.Debug(fmt.Sprintf("iP181 done. found %d", len(ipList)))
	return ipList, nil
}

func xici() ([]storage.ProxyItem, error) {
	var url string = "http://www.xicidaili.com/"
	var ipList []storage.ProxyItem
	doc, err := request.Document(url)
	if err != nil {
		return ipList, err
	}
	doc.Find("#ip_list tbody tr").Each(func(i int, s *goquery.Selection) {
		nodes := s.Find("td")
		if len(nodes.Nodes) == 8 {
			ipAddr, err1 := s.Find("td:nth-child(2)").Html()
			portStr, err2 := s.Find("td:nth-child(3)").Html()
			if err1 != nil || err2 != nil {
				util.Debug(fmt.Sprintf("xici parse error %s %s", err1, err2))
			} else if ipAddr == "" || portStr == "" {
				util.Debug(fmt.Sprintf("xici matched empty %s:%s", ipAddr, portStr))
			} else {
				portInt, err3 := strconv.Atoi(portStr)
				if err3 != nil || !storage.IpReg.MatchString(ipAddr) {
					util.Debug(fmt.Sprintf("xici ip %s or port parse error %s", ipAddr, err3))
				} else {
					ipList = append(ipList, storage.NewProxyItem(ipAddr, uint16(portInt)))
				}
			}
		}
	})
	util.Debug(fmt.Sprintf("xici done. found %d", len(ipList)))
	return ipList, nil
}

func ydl() ([]storage.ProxyItem, error) {
	var url string = "http://www.youdaili.net/Daili/http/"
	var ipList []storage.ProxyItem
	doc, err := request.Document(url)
	if err != nil {
		return ipList, err
	}
	url, exists := doc.Find(".chunlist ul li:nth-child(1) a").Attr("href")
	if !exists {
		return ipList, fmt.Errorf("ydl entry url not found")
	}
	doc, err = request.Document(url)
	if err != nil {
		return ipList, err
	}
	str, err := doc.Find(".conl .content").Html()
	if err != nil {
		return ipList, err
	}
	totalPage, err1 := doc.Find(".conl .pagebreak li:nth-last-child(2)").Find("a").Html()
	totalPageInt, err2 := strconv.Atoi(totalPage)
	if err1 != nil || err2 != nil {
		ipList = storage.FindAllProxy(str)
		util.Debug(fmt.Sprintf("ydl done. found %d", len(ipList)))
		return ipList, nil
	}
	for i := 2; i <= totalPageInt; i++ {
		currUrl := strings.Replace(url, ".html", fmt.Sprintf("_%d.html", i), 1)
		doc, err = request.Document(currUrl)
		if err != nil {
			continue
		}
		currStr, err := doc.Find(".conl .content").Html()
		if err != nil {
			continue
		}
		str = str + currStr
	}
	ipList = storage.FindAllProxy(str)
	util.Debug(fmt.Sprintf("ydl done. found %d from %d pages", len(ipList), totalPageInt))
	return ipList, nil
}

func xdl() ([]storage.ProxyItem, error) {
	var ipList []storage.ProxyItem
	var url string = "http://www.xdaili.cn/ipagent/freeip/getFreeIps?page=1&rows=10"
	str, err := request.HttpGet(url, map[string]string{"Accept": "application/json, text/javascript, */*; q=0.01"})
	if err != nil {
		return ipList, err
	}
	type res struct {
		Rows []struct {
			Ip   string
			Port string
		}
	}
	var info res
	err = json.Unmarshal(str, &info)
	if err != nil {
		return ipList, err
	}
	for _, item := range info.Rows {
		portInt, err := strconv.Atoi(item.Port)
		if err != nil {
			continue
		}
		ipList = append(ipList, storage.NewProxyItem(item.Ip, uint16(portInt)))
	}
	util.Debug(fmt.Sprintf("xdl done.found %d", len(ipList)))
	return ipList, nil
}
