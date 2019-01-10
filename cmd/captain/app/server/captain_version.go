package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
)

func (s *CaptainServer) Version(req *tw_com_softleader_captainkube.VersionRequest, stream tw_com_softleader_captainkube.Captain_VersionServer) error {
	out := sio.NewStreamWriter(func(p []byte) error {
		return stream.Send(&tw_com_softleader_captainkube.ChunkMessage{
			Msg: p,
		})
	})

	var v string
	if req.GetFull() {
		v = s.Metadata.FullString()
	} else {
		v = s.Metadata.String()
	}
	out.Write([]byte(fmt.Sprintln("Captain: " + v)))

	endpoints, err := s.lookupCaplet(logrus.StandardLogger(), req.GetColor())
	if err != nil {
		return err
	}

	endpoints.Each(func(e *caplet.Endpoint) {
		if r, err := e.Version(req.GetFull(), req.GetTimeout()); err != nil {
			out.Write([]byte(err.Error()))
		} else {
			msg := fmt.Sprintln("Caplet " + r.GetHostname() + ": " + r.GetVersion())
			out.Write(e.Color([]byte(msg)))
		}
	})

	return nil
}
