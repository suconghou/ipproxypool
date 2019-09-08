package request

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

func tslConfig() *tls.Config {
	return &tls.Config{InsecureSkipVerify: true}
}

func dialContextConfig() func(ctx context.Context, network, addr string) (net.Conn, error) {
	return (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext
}

func httpProxyConfig(str string) func(*http.Request) (*url.URL, error) {
	if str == "" {
		return http.ProxyFromEnvironment
	}
	urlproxy, err := url.Parse(str)
	if err != nil {
		return http.ProxyFromEnvironment
	}
	return http.ProxyURL(urlproxy)
}

func transportConfig(urlproxy string) *http.Transport {
	return &http.Transport{
		Proxy:                 httpProxyConfig(urlproxy),
		DialContext:           dialContextConfig(),
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       tslConfig(),
	}
}
