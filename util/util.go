package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	// Log to stdout
	Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
)

type resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// JSONPut resp json,如果v是byte类型,我们应直接使用,byte类型再json.Marshal就是base64字符串了,string类型经json.Marshal后转为对于的byte
func JSONPut(w http.ResponseWriter, v any) (int, error) {
	var (
		bs  []byte
		err error
	)
	if bb, ok := v.([]byte); !ok {
		if bs, err = json.Marshal(v); err != nil {
			return 0, err
		}
	} else {
		bs = bb
	}
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Max-Age", "864000")
	return w.Write(bs)
}

func JSON(w http.ResponseWriter, msg string, code int) (int, error) {
	return JSONPut(w, resp{code, msg, nil})
}

func JSONData(w http.ResponseWriter, data any) (int, error) {
	return JSONPut(w, resp{0, "ok", data})
}

// 文件或文件夹存在返回true,特别的如果是一个空文件报告为false,使得上层可以覆盖这个文件
func FileExists(filename string) bool {
	s, err := os.Stat(filename)
	if err == nil {
		if s.IsDir() {
			return true
		}
		if s.Size() < 1 {
			return false
		}
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

// 当返回错误时，已向客户端发送错误消息
func Parse(w http.ResponseWriter, r *http.Request, v any) error {
	bs, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 8192))
	if err == nil {
		if len(bs) <= 4 {
			err = fmt.Errorf("bad request")
		}
	}
	if err != nil {
		JSON(w, err.Error(), -2)
		return err
	}
	err = json.Unmarshal(bs, v)
	if err != nil {
		JSON(w, err.Error(), -3)
		return err
	}
	return nil
}
