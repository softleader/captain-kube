package dur

import "time"

var DefaultDeadline = "30m"

func Parse(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		d, _ = time.ParseDuration(DefaultDeadline)
	}
	return d
}
