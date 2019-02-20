package tcp

import (
	"fmt"
	"net"
	"time"
)

// IsReachable 檢查某個 host:port tcp 是否有通
func IsReachable(host string, port, timeoutSeconds int) bool {
	timeOut := time.Duration(timeoutSeconds) * time.Second
	_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%v", host, port), timeOut)
	return err == nil
}
