package util

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
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
