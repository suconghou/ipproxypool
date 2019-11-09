package route

import (
	"fmt"
	"ipproxypool/request"
	"ipproxypool/util"
	"net/http"
)

type fetchcfg struct {
	*request.FetchConfig
	Query request.QueryConfig
}

func fetchurl(w http.ResponseWriter, r *http.Request, match []string) error {
	var data fetchcfg
	if err := parse(w, r, &data); err != nil {
		return err
	}
	if len(data.Urls) < 1 {
		err := fmt.Errorf("at least one url")
		util.JSONPut(w, resp{-4, err.Error(), nil})
		return err
	}
	ret, err := request.New(data.FetchConfig).Do(match[1], data.Query)
	if err != nil {
		util.JSONPut(w, resp{-6, err.Error(), nil})
		return err
	}
	if data.Cache > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", data.Cache))
	}
	if bs, ok := ret.([]byte); ok {
		_, err = w.Write(bs)
	} else {
		_, err = util.JSONPut(w, resp{0, "ok", ret})
	}
	return err
}
