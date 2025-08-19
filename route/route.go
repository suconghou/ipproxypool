package route

import (
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
}
