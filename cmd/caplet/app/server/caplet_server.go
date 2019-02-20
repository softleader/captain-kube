package server

import (
	"github.com/softleader/captain-kube/pkg/release"
	"os"
)

// CapletServer 封裝了啟動 server 的相關設定
type CapletServer struct {
	metadata *release.Metadata
	hostname string
}

// NewCapletServer 建議 CapletServer 物件
func NewCapletServer(metadata *release.Metadata) (s *CapletServer) {
	s = &CapletServer{
		metadata: metadata,
	}
	s.hostname, _ = os.Hostname()
	return
}
