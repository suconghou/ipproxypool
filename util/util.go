package util

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"
)

var (
	regproxyurl = regexp.MustCompile(`^/(?i:https?):/{1,2}[[:print:]]+$`)
	// Log to stdout
	Log = log.New(os.Stdout, "", 0)
)

// JSONPut resp json
func JSONPut(w http.ResponseWriter, v interface{}) (int, error) {
	bs, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	return w.Write(bs)
}

// PortOpen test port is reachable
func PortOpen(ipPort string) bool {
	_, err := net.DialTimeout("tcp", ipPort, time.Second)
	if err != nil {
		return false
	}
	return true
}

// IoCopy copy two stream
func IoCopy(c1, c2 io.ReadWriteCloser) error {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { io.Copy(c1, c2); wg.Done() }()
	go func() { io.Copy(c2, c1); wg.Done() }()
	var e1 = c1.Close()
	var e2 = c2.Close()
	wg.Wait()
	if e1 == nil {
		return e2
	}
	return e1
}

// FileExists check if file exist or dir exist , !info.IsDir()
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		// error , treat as exist, so stop create it
		return true
	}
}

// ValidMethod test http method ok
func ValidMethod(m string) bool {
	if m == http.MethodGet || m == http.MethodPost || m == http.MethodPut || m == http.MethodDelete {
		return true
	}
	return false
}

// ValidProxyURL do valid url
func ValidProxyURL(u string) bool {
	return regproxyurl.MatchString(u)
}
