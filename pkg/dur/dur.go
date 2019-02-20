package dur

import (
	"strconv"
	"time"
)

// DefaultDeadline 預設的 deadline
var DefaultDeadline = "30m"

// Parse 將傳入的 string 轉換成 duration, 若不給單位預設為秒
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
