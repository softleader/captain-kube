package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/utils"
)

func (s *CaptainServer) Version(req *proto.VersionRequest, stream proto.Captain_VersionServer) error {
	log := logrus.New()
	log.SetOutput(sio.NewStreamWriter(func(p []byte) error {
		return stream.Send(&proto.ChunkMessage{
			Msg: p,
		})
	}))
	log.SetFormatter(&utils.PlainFormatter{})

	log.Println(fmt.Sprintf(`Captain: %s`, s.Metadata.String(req.GetShort())))

	endpoints, err := s.lookupCaplet()
	if err != nil {
		return err
	}

	log.SetNoLock()
	endpoints.Each(func(e *caplet.Endpoint) {
		if r, err := e.Version(req.GetShort(), req.GetTimeout()); err != nil {
			log.Error(err)
		} else {
			msg := fmt.Sprintf(`Caplet %s: %s`, r.GetHostname(), r.GetVersion())
			if req.GetColor() {
				msg = e.Color(msg)
			}
			log.Println(msg)
		}
	})

	return nil
}
