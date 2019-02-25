package server

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/dockerd"
	pb "github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/utils"
)

// Prune 執行 docker system prune
func (s *CapletServer) Prune(req *pb.PruneRequest, stream pb.Caplet_PruneServer) error {
	log := logrus.New()
	log.SetOutput(sio.NewStreamWriter(func(p []byte) error {
		return stream.Send(&pb.ChunkMessage{
			Hostname: s.hostname,
			Msg:      p,
		})
	}))
	log.SetFormatter(&utils.PlainFormatter{})
	if req.GetVerbose() {
		log.SetLevel(logrus.DebugLevel)
	}
	return dockerd.Prune(log)
}
