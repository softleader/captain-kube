package caplet

import (
	"bytes"
	"fmt"
	"github.com/softleader/captain-kube/pkg/color"
	pb "github.com/softleader/captain-kube/pkg/proto"
	"sync"
)

const (
	// EnvPort key to specify caplet port
	EnvPort = "CAPLET_PORT"
	// EnvHostname key to specify hostname to lookup
	EnvHostname = "CAPLET_HOSTNAME"
	// DefaultPort caplet default port
	DefaultPort = 50051
	// DefaultHostname caplet default hostname to lookup
	DefaultHostname = "caplet"
)

// Endpoint 封裝了 caplet 的連線資訊
type Endpoint struct {
	Target string
	Port   int
	Color  func([]byte) []byte // output 塗色
}

// NewEndpoint 建立 Endpoint 物件
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

// Endpoints 代表多個 Endpoint
type Endpoints []*Endpoint

// Each 非同步的針對每個 Endpoint 執行傳入的 consumer function
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

func format(last, chunk *pb.ChunkMessage) []byte {
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
