package server

import (
	"github.com/docker/docker/cli/command"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/softleader/captain-kube/pkg/dockerctl"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/logger"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
)

type CaptainServer struct {
	log *logger.Logger
}

func (s *CaptainServer) PullImage(req *proto.PullImageRequest, stream proto.Caplet_PullImageServer) error {
	log := logger.New(sio.NewStreamWriter(func(p []byte) error {
		return stream.Send(&proto.ChunkMessage{
			Msg: p,
		})
	})).WithVerbose(req.GetVerbose())
	for _, image := range req.Images {
		if err := pull(log, image, req.GetRegistryAuth()); err != nil {
			return err
		}
	}
	return nil
}

func pull(log *logger.Logger, image *proto.Image, auth *proto.RegistryAuth) error {
	if tag := image.GetTag(); len(tag) == 0 {
		image.Tag = "latest"
	}
	rc, err := dockerctl.Pull(log, chart.Image{
		Host: image.Host,
		Repo: image.Repo,
		Tag:  image.Tag,
	}, auth)
	if err != nil {
		return err
	}
	defer rc.Close()
	return jsonmessage.DisplayJSONMessagesToStream(rc, command.NewOutStream(log.GetOutput()), nil)
}
