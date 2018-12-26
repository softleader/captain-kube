package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
)

func (s *CaptainServer) Version(req *proto.VersionRequest, stream proto.Captain_VersionServer) error {
	out := sio.NewStreamWriter(func(p []byte) error {
		return stream.Send(&proto.ChunkMessage{
			Msg: p,
		})
	})

	out.Write([]byte(fmt.Sprintln("Captain: " + s.Metadata.String(req.GetShort()))))

	endpoints, err := s.lookupCaplet(logrus.StandardLogger(), req.GetColor())
	if err != nil {
		return err
	}

	endpoints.Each(func(e *caplet.Endpoint) {
		if r, err := e.Version(req.GetShort(), req.GetTimeout()); err != nil {
			out.Write([]byte(err.Error()))
		} else {
			msg := fmt.Sprintln("Caplet " + r.GetHostname() + ": " + r.GetVersion())
			out.Write(e.Color([]byte(msg)))
		}
	})

	return nil
}
