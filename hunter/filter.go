package hunter

import (
	"regexp"
	"time"
)

var (
	ipReg    = regexp.MustCompile(`((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))))`)
	proxyReg = regexp.MustCompile(`((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))):([1-9]?\d{0,4})`)
)

// Start for get ips
func Start() {
	go func() {
		for {
			time.Sleep(time.Second)
		}
	}()
}
