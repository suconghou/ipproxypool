package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"ipproxypool/hunter"
	"ipproxypool/pool"
	"ipproxypool/util"
	"ipproxypool/web"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	port      string
	doc       string
	startTime = time.Now()
)

var sysStatus struct {
	Uptime       string
	GoVersion    string
	Hostname     string
	MemAllocated uint64
	MemTotal     uint64
	MemSys       uint64
	NumGoroutine int
	CpuNum       int
	Pid          int
}

func init() {
	flag.StringVar(&port, "port", "7090", "give me a port number")
	flag.StringVar(&doc, "doc", "./", "document root dir")
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !filepath.IsAbs(doc) {
		doc = filepath.Join(pwd, doc)
	}
	if f, err := os.Stat(doc); err == nil {
		if !f.Mode().IsDir() {
			fmt.Println(doc + " is not directory")
			os.Exit(3)
		}
	} else {
		fmt.Println(doc + " not exists")
		os.Exit(2)
	}

}

func main() {
	hunter.StartHunter()
	pool.ServeConnection()

	flag.Parse()
	http.HandleFunc("/status", status)
	http.HandleFunc("/post", post)
	http.HandleFunc("/get", get)
	http.HandleFunc("/", routeMatch)
	fmt.Println("Starting up on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func status(w http.ResponseWriter, r *http.Request) {
	memStat := new(runtime.MemStats)
	runtime.ReadMemStats(memStat)
	sysStatus.Uptime = time.Since(startTime).String()
	sysStatus.NumGoroutine = runtime.NumGoroutine()
	sysStatus.MemAllocated = memStat.Alloc
	sysStatus.MemTotal = memStat.TotalAlloc
	sysStatus.MemSys = memStat.Sys
	sysStatus.CpuNum = runtime.NumCPU()
	sysStatus.GoVersion = runtime.Version()
	sysStatus.Hostname, _ = os.Hostname()
	sysStatus.Pid = os.Getpid()
	if bs, err := json.Marshal(&sysStatus); err != nil {
		http.Error(w, fmt.Sprintf("%s", err), 500)
	} else {
		util.JsonPut(w, bs, true, 60)
	}
}

func routeMatch(w http.ResponseWriter, r *http.Request) {
	found := false
	for _, p := range web.RoutePath {
		if p.Reg.MatchString(r.URL.Path) {
			found = true
			p.Handler(w, r, p.Reg.FindStringSubmatch(r.URL.Path))
			break
		}
	}
	if !found {
		fallback(w, r)
	}
}

func fallback(w http.ResponseWriter, r *http.Request) {
	var files []string
	if r.URL.Path == "/" {
		files = []string{"index.html"}
	} else {
		files = []string{r.URL.Path, filepath.Join(r.URL.Path, "index.html")}
	}
	if !tryFiles(files, w, r) {
		http.NotFound(w, r)
	}
}

func tryFiles(files []string, w http.ResponseWriter, r *http.Request) bool {
	for _, file := range files {
		var realpath string = filepath.Join(doc, file)
		if f, err := os.Stat(realpath); err == nil {
			if f.Mode().IsRegular() {
				http.ServeFile(w, r, realpath)
				return true
			}
		}
	}
	return false
}

func post(w http.ResponseWriter, r *http.Request) {

}

func get(w http.ResponseWriter, r *http.Request) {

}
