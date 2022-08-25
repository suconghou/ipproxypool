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
		"Range",
		"Content-Length",
		"Content-Type",
	}
	exposeHeaders = []string{
		"Accept-Ranges",
		"Content-Range",
		"Content-Length",
		"Content-Type",
		"Content-Encoding",
		"Cache-Control",
	}
)

func copyHeader(from http.Header, to http.Header, headers []string) http.Header {
	for _, k := range headers {
		if v := from.Get(k); v != "" {
			to.Set(k, v)
		}
	}
	return to
}

// URL proxy request to target
func URL(w http.ResponseWriter, r *http.Request, match []string) error {
	var (
		u         = r.RequestURI
		reqHeader = http.Header{}
	)
	u = strings.Replace(strings.TrimPrefix(u, "/"), ":/", "://", 1)
	target, err := url.Parse(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	defer r.Body.Close()
	resp, err := request.GetResponse(target, r.Method, copyHeader(r.Header, reqHeader, fwdHeaders), r.Body, "", 600, 2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	defer resp.Body.Close()
	h := w.Header()
	copyHeader(resp.Header, h, exposeHeaders)
	if resp.StatusCode == http.StatusOK {
		h.Set("Cache-Control", "public, max-age=15552000")
	}
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Max-Age", "864000")
	w.WriteHeader(resp.StatusCode)
	if n, err := io.Copy(w, resp.Body); err == nil {
		util.Log.Printf("%s %d %d", u, n, resp.StatusCode)
	} else {
		util.Log.Printf("%s %d %d %v", u, n, resp.StatusCode, err)
	}
	return nil
}
