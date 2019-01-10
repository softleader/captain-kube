package server

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/dockerd"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/utils"
)

func (s *CapletServer) Prune(req *tw_com_softleader.PruneRequest, stream tw_com_softleader.Caplet_PruneServer) error {
	log := logrus.New()
	log.SetOutput(sio.NewStreamWriter(func(p []byte) error {
		return stream.Send(&tw_com_softleader.ChunkMessage{
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
