package request

import "time"

func timeoutConfig(timeout int) time.Duration {
	return time.Duration(timeout) * time.Second
}
