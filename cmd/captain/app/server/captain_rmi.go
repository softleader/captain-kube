package server

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/utils"
)

func (s *CaptainServer) Rmi(req *captainkube_v2.RmiRequest, stream captainkube_v2.Captain_RmiServer) error {
	log := logrus.New()
	log.SetOutput(sio.NewStreamWriter(func(p []byte) error {
		return stream.Send(&captainkube_v2.ChunkMessage{
			Msg: p,
		})
	}))
	log.SetFormatter(&utils.PlainFormatter{})
	if req.GetVerbose() {
		log.SetLevel(logrus.DebugLevel)
	}

	endpoints, err := s.lookupCaplet(log, req.GetColor())
	if err != nil {
		return err
	}

	log.SetNoLock()
	endpoints.Each(func(e *caplet.Endpoint) {
		if err := e.Rmi(log, req, req.GetTimeout()); err != nil {
			log.Error(err)
		}
	})
	return nil
}