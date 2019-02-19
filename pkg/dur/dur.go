package dur

import "time"

var DefaultDeadlineSecond int64 = 1800

func Deadline(timeout int64) time.Duration {
	if timeout <= 0 {
		timeout = DefaultDeadlineSecond
	}
	return time.Duration(timeout) * time.Second
}