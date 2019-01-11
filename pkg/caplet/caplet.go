package caplet

import (
	"bytes"
	"fmt"
	"github.com/softleader/captain-kube/pkg/color"
	"github.com/softleader/captain-kube/pkg/proto"
	"sync"
)

const (
	EnvPort         = "CAPLET_PORT"
	EnvHostname     = "CAPLET_HOSTNAME"
	DefaultPort     = 50051
	DefaultHostname = "caplet"
)

type Endpoint struct {
	Target string
	Port   int
	Color  func([]byte) []byte // output 塗色
}

func NewEndpoint(target string, port int) *Endpoint {
	return &Endpoint{
		Target: target,
		Port:   port,
		Color:  color.Plain,
	}
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%v", e.Target, e.Port)
}

type Endpoints []*Endpoint

func (endpoints Endpoints) Each(consumer func(e *Endpoint)) {
	var wg sync.WaitGroup
	for _, ep := range endpoints {
		wg.Add(1)
		go func(ep *Endpoint) {
			defer wg.Done()
			consumer(ep)
		}(ep)
	}
	wg.Wait()
}

func format(last, chunk *captainkube_v2.ChunkMessage) []byte {
	msg := chunk.GetMsg()
	if msg == nil || len(msg) == 0 {
		return msg
	}
	if last != nil { // 如果上一次的訊息沒有斷行符號, 這次就不要加上 hostname 吧
		if lastMsg := last.GetMsg(); lastMsg != nil && !bytes.ContainsAny(lastMsg, "\n") {
			return msg
		}
	}
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s | ", chunk.GetHostname()))
	buf.Write(msg)
	return buf.Bytes()
}
