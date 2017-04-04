package util

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var debuglog bool = true

func JsonPut(w http.ResponseWriter, bs []byte, httpCache bool, cacheTime uint32) {
	CrossShare(w)
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	if httpCache {
		UseHttpCache(w, cacheTime)
	}
	w.Write(bs)
}

func UseHttpCache(w http.ResponseWriter, cacheTime uint32) {
	w.Header().Set("Expires", time.Now().Add(time.Second*time.Duration(cacheTime)).Format(http.TimeFormat))
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", cacheTime))
}

func CrossShare(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, HEAD, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Content-Length, Accept, Accept-Encoding")
}

func Debug(args ...interface{}) {
	if debuglog {
		log.Println(args...)
	}
}

func Halt(args ...interface{}) {
	if debuglog {
		log.Println(args...)
	}
	os.Exit(1)
}

func SetDebug(debug bool) {
	debuglog = debug
}

func SetOutput(w io.Writer) {
	log.SetOutput(w)
}
