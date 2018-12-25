package server

import (
	"fmt"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/proto"
)

func (s *CaptainServer) Version(req *proto.VersionRequest, stream proto.Captain_VersionServer) error {
	stream.Send(&proto.ChunkMessage{
		Msg: []byte(fmt.Sprintf(`Captain: %s`, s.Metadata.String(req.GetShort()))),
	})
	endpoints, err := s.lookupCaplets()
	if err != nil {
		return nil
	}

	return endpoints.Each(func(e *caplet.Endpoint) error {
		r, err := e.Version(req.GetShort(), req.GetTimeout())
		if err != nil {
			return err
		}
		return stream.Send(&proto.ChunkMessage{
			Msg: []byte(fmt.Sprintf(`Caplet %s: %s`, r.GetHostname(), r.GetVersion())),
		})
	})
}
