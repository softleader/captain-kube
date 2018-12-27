package tcp

import (
	"fmt"
	"net"
	"time"
)

func Reachable(host string, port, timeOutSec int) bool {
	timeOut := time.Duration(timeOutSec) * time.Second
	_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%v", host, port), timeOut)
	return err == nil
}
