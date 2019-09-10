package request

import (
	"bytes"
	"fmt"
	"ipproxypool/util"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func process(data map[string][]byte, action string, params string) (interface{}, error) {
	switch action {
	case "html":
		return html(data), nil
	case "text":
		return text(data)
	case "json":
		return data, nil
	default:
		arr := strings.Split(params, "|")
		if len(arr) == 2 {
			return findAttr(data, action, strings.Split(arr[0], ","), strings.Split(arr[1], ","))
		} else if len(arr) == 1 {
			return findAttr(data, action, strings.Split(arr[0], ","), []string{})
		}
		return nil, fmt.Errorf("error action & params %s %s", action, params)
	}
}

func html(data map[string][]byte) []byte {
	var ret []byte
	for _, item := range data {
		ret = append(ret, item...)
	}
	return ret
}

func text(data map[string][]byte) ([]byte, error) {
	var ret []byte
	for _, item := range data {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(item))
		if err != nil {
			return ret, err
		}
		ret = append(ret, doc.Text()...)
	}
	return ret, nil
}

func findAttr(data map[string][]byte, selector string, attrNames []string, methods []string) ([]map[string]string, error) {
	var ret = []map[string]string{}
	for _, item := range data {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(item))
		if err != nil {
			return ret, err
		}
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			var one = map[string]string{}
			for _, attrName := range attrNames {
				if attrName != "" {
					v, _ := s.Attr(attrName)
					one[attrName] = v
				}
			}
			for _, m := range methods {
				switch m {
				case "html":
					str, err := s.Html()
					one[m] = str
					if err != nil {
						util.Logger.Print(err)
					}
				case "text":
					one[m] = s.Text()
				default:
					one["text"] = s.Text()
				}
			}
			ret = append(ret, one)
		})
	}
	return ret, nil
}
