package route

import (
	"net/http"
	"regexp"
)

// 路由定义
type routeInfo struct {
	Reg     *regexp.Regexp
	Handler func(http.ResponseWriter, *http.Request, []string)
}

// Route export route list
var Route = []routeInfo{
	{regexp.MustCompile(`^/api/list$`), proxylist},
}

func proxylist(w http.ResponseWriter, r *http.Request, match []string) {

}
