package dur

import (
	"strconv"
	"time"
)

var DefaultDeadline = "30m"

func Parse(s string) time.Duration {
	if i, err := strconv.Atoi(s); err == nil {
		return time.Duration(i) * time.Second
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		d, _ = time.ParseDuration(DefaultDeadline)
	}
	return d
}
