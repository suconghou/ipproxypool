package proxy

import (
	"io"
	"ipproxypool/request"
	"ipproxypool/util"
	"net/http"
	"net/url"
	"strings"
)

var (
	fwdHeaders = []string{
		"User-Agent",
		"Accept",
		"Accept-Encoding",
		"Accept-Language",
		"If-Modified-Since",
		"If-None-Match",
		"Range",
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

func copyHeader(from http.Header, to http.Header, headers []string) http.Header {
	for _, k := range headers {
		if v := from.Get(k); v != "" {
			to.Set(k, v)
		}
	}
	return to
}

// URL proxy request to target
func URL(w http.ResponseWriter, r *http.Request) {
	var (
		u         = r.RequestURI
		reqHeader = http.Header{}
	)
	u = strings.Replace(strings.TrimPrefix(u, "/"), ":/", "://", 1)
	target, err := url.Parse(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	resp, err := request.GetResponse(target, r.Method, copyHeader(r.Header, reqHeader, fwdHeaders), r.Body, "", 600, 2)
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
