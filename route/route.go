package route

import (
	"ipproxypool/proxy"
	"ipproxypool/query"
	"ipproxypool/tasks"
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
	{regexp.MustCompile(`^/api/task/info$`), tasks.Info},
	{regexp.MustCompile(`^/api/task/add$`), tasks.Add},
	{regexp.MustCompile(`^/api/fetch/(\w{1,10})$`), query.GoQuery},
	{regexp.MustCompile(`^/(?i:https?):/{1,2}[[:print:]]+$`), proxy.URL},
}
