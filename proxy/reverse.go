package proxy

import (
	"cloud/util"
	"io"
	"ipproxypool/request"
	"net/http"
	"net/url"
	"strings"
)

var (
	xheaders = []string{
		"X-Forwarded-For",
		"X-Forwarded-Host",
		"X-Forwarded-Server",
		"X-Forwarded-Port",
		"X-Forwarded-Proto",
		"X-Client-Ip",
		"Cookie",
	}
	exposeHeaders = []string{
		"Accept-Ranges",
		"Content-Range",
		"Content-Length",
		"Content-Type",
		"Content-Encoding",
		"Date",
		"Expires",
		"Cache-Control",
	}
)

func cleanHeader(header http.Header, headers []string) http.Header {
	for _, k := range headers {
		header.Del(k)
	}
	return header
}

func copyHeader(from http.Header, to http.Header, headers []string) {
	for _, k := range headers {
		if v := from.Get(k); v != "" {
			to.Set(k, v)
		}
	}
}

// URL proxy request to target
func URL(w http.ResponseWriter, r *http.Request) {
	var u = r.RequestURI
	u = strings.Replace(strings.TrimPrefix(u, "/"), ":/", "://", 1)
	target, err := url.Parse(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	resp, err := request.GetResponse(target, r.Method, cleanHeader(r.Header, xheaders), r.Body, "", 600, 2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	to := w.Header()
	copyHeader(resp.Header, to, exposeHeaders)
	if resp.StatusCode == http.StatusOK {
		to.Set("Cache-Control", "public, max-age=15552000")
	}
	w.WriteHeader(resp.StatusCode)
	n, err := io.Copy(w, resp.Body)
	util.Log.Printf("%s %d %v", u, n, err)
}
