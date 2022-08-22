package route

import (
	"encoding/json"
	"fmt"
	"io"
	"ipproxypool/util"
	"net/http"
	"regexp"
)

// 路由定义
type routeInfo struct {
	Reg     *regexp.Regexp
	Handler func(http.ResponseWriter, *http.Request, []string) error
}

// Route export route list
var Route = []routeInfo{
	{regexp.MustCompile(`^/api/proxy/one$`), proxyone},
	{regexp.MustCompile(`^/api/proxy/add$`), proxyadd},
	{regexp.MustCompile(`^/api/proxy/info$`), proxyinfo},
	{regexp.MustCompile(`^/api/task/info$`), taskinfo},
	{regexp.MustCompile(`^/api/task/add$`), taskadd},
	{regexp.MustCompile(`^/api/fetch/(\w{1,10})$`), fetchurl},
}

type resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func parse(w http.ResponseWriter, r *http.Request, v any) error {
	bs, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 8192))
	if err == nil {
		if len(bs) <= 4 {
			err = fmt.Errorf("bad request")
		}
	}
	if err != nil {
		util.JSONPut(w, resp{-2, err.Error(), nil})
		return err
	}
	err = json.Unmarshal(bs, v)
	if err != nil {
		util.JSONPut(w, resp{-3, err.Error(), nil})
		return err
	}
	return nil
}
