package request

import (
	"bytes"
	"fmt"
	"ipproxypool/util"

	"github.com/PuerkitoBio/goquery"
)

// QueryConfig for query
type QueryConfig map[string]struct {
	Attrs   []string
	Methods []string
}

func process(data map[string][]byte, action string, query QueryConfig) (interface{}, error) {
	switch action {
	case "html":
		return html(data), nil
	case "text":
		return text(data)
	case "json":
		return data, nil
	default:
		return findAttr(data, query)
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

func findAttr(data map[string][]byte, queries QueryConfig) (map[string][]map[string]string, error) {
	var ret = map[string][]map[string]string{}
	for _, item := range data {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(item))
		if err != nil {
			return ret, err
		}
		for q, query := range queries {
			var oneret = []map[string]string{}
			fmt.Println(query.Methods, query.Attrs)
			doc.Find(q).Each(func(i int, s *goquery.Selection) {
				var one = map[string]string{}
				for _, attrName := range query.Attrs {
					if attrName != "" {
						v, _ := s.Attr(attrName)
						one[attrName] = v
					}
				}
				for _, m := range query.Methods {
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
				oneret = append(oneret, one)
			})
			ret[q] = oneret
		}
	}
	return ret, nil
}
