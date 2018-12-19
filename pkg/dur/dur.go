package dur

import "time"

var DefaultDeadlineSecond = 30

func Deadline(timeout int) time.Duration {
	if timeout <= 0 {
		timeout = DefaultDeadlineSecond
	}
	return time.Duration(timeout)
}
