package tcp

import (
	"fmt"
	"net"
	"time"
)

func IsReachable(host string, port, timeoutSeconds int) bool {
	timeOut := time.Duration(timeoutSeconds) * time.Second
	_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%v", host, port), timeOut)
	return err == nil
}
