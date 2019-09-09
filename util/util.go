package util

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	// Logger to stdout
	Logger = log.New(os.Stdout, "", 0)
)

// JSONPut resp json
func JSONPut(w http.ResponseWriter, v interface{}) (int, error) {
	bs, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}
	h := w.Header()
	h.Set("Content-Type", "text/json; charset=utf-8")
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
