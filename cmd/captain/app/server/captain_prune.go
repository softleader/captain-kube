package server

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/utils"
)

func (s *CaptainServer) Prune(req *proto.PruneRequest, stream proto.Captain_PruneServer) error {
	log := logrus.New()
	log.SetOutput(sio.NewStreamWriter(func(p []byte) error {
		return stream.Send(&proto.ChunkMessage{
			Msg: p,
		})
	}))
	log.SetFormatter(&utils.PlainFormatter{})
	if req.GetVerbose() {
		log.SetLevel(logrus.DebugLevel)
	}

	endpoints, err := s.lookupInstances()
	if err != nil {
		return err
	}

	log.SetNoLock()
	endpoints.Each(func(e *caplet.Endpoint) {
		if err := e.Prune(log, req, req.GetTimeout()); err != nil {
			log.Error(err)
		}
	})
	return nil
}
