package request

import (
	"bytes"
	"ipproxypool/util"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// QueryConfig for query
type QueryConfig map[string]struct {
	Attrs   map[string]string
	Methods map[string]string
	Query   QueryConfig
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
		ret = append(ret, strings.TrimSpace(doc.Text())...)
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
			doc.Find(q).Each(func(i int, s *goquery.Selection) {
				var one = map[string]string{}
				for attr, key := range query.Attrs {
					one[key] = s.AttrOr(attr, "")
				}
				for m, key := range query.Methods {
					switch m {
					case "html":
						str, err := s.Html()
						one[key] = str
						if err != nil {
							util.Log.Print(err)
						}
					case "text":
						one[key] = s.Text()
					default:
						one[key] = s.Text()
					}
				}
				subQuery(s, query.Query, &one)
				oneret = append(oneret, one)
			})
			ret[q] = oneret
		}
	}
	return ret, nil
}

func subQuery(doc *goquery.Selection, queries QueryConfig, one *map[string]string) {
	if queries == nil {
		return
	}
	var oneitem = *one
	for q, query := range queries {
		doc.Find(q).Each(func(i int, s *goquery.Selection) {
			for attr, key := range query.Attrs {
				oneitem[key] = s.AttrOr(attr, "")
			}
			for m, key := range query.Methods {
				switch m {
				case "html":
					str, err := s.Html()
					oneitem[key] = str
					if err != nil {
						util.Log.Print(err)
					}
				case "text":
					oneitem[key] = s.Text()
				default:
					oneitem[key] = s.Text()
				}
			}
			subQuery(s, query.Query, one)
		})
	}
}
