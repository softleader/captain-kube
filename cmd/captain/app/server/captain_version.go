package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/dur"
	pb "github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
)

// Version 返回 captain 所有 caplets 的版本
func (s *CaptainServer) Version(req *pb.VersionRequest, stream pb.Captain_VersionServer) error {
	out := sio.NewStreamWriter(func(p []byte) error {
		return stream.Send(&pb.ChunkMessage{
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

	timeout := dur.Parse(req.GetTimeout())
	endpoints.Each(func(e *caplet.Endpoint) {
		if r, err := e.CallVersion(req.GetFull(), timeout); err != nil {
			out.Write([]byte(err.Error()))
		} else {
			msg := fmt.Sprintln("Caplet " + r.GetHostname() + ": " + r.GetVersion())
			out.Write(e.Color([]byte(msg)))
		}
	})

	return nil
}
