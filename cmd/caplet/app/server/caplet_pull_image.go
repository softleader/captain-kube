package server

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/dockerd"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/utils"
)

// PullImage 執行 docker pull
func (s *CapletServer) PullImage(req *captainkube_v2.PullImageRequest, stream captainkube_v2.Caplet_PullImageServer) error {
	log := logrus.New()
	log.SetOutput(sio.NewStreamWriter(func(p []byte) error {
		return stream.Send(&captainkube_v2.ChunkMessage{
			Hostname: s.hostname,
			Msg:      p,
		})
	}))
	log.SetFormatter(&utils.PlainFormatter{})
	if req.GetVerbose() {
		log.SetLevel(logrus.DebugLevel)
	}
	for _, image := range req.GetImages() {
		if err := pull(log, image, req.GetRegistryAuth()); err != nil {
			return err
		}
	}
	return nil
}

func pull(log *logrus.Logger, image *captainkube_v2.Image, auth *captainkube_v2.RegistryAuth) error {
	if tag := image.GetTag(); len(tag) == 0 {
		image.Tag = chart.DefaultTag
	}
	err := dockerd.Pull(log, chart.Image{
		Host: image.Host,
		Repo: image.Repo,
		Tag:  image.Tag,
	}, auth)
	if err != nil {
		return err
	}
	return nil
}
