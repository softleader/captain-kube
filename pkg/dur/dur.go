package dur

import "time"

var DefaultDeadlineSecond = 300

func Deadline(timeout int) time.Duration {
	if timeout <= 0 {
		timeout = DefaultDeadlineSecond
	}
	return time.Duration(timeout)
}
