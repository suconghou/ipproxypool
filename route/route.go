package route

import (
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
}

type resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
