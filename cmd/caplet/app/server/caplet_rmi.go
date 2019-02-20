package server

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/dockerd"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/utils"
)

// Rmi 執行 docker rmi
func (s *CapletServer) Rmi(req *captainkube_v2.RmiRequest, stream captainkube_v2.Caplet_RmiServer) error {
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
		i := chart.Image{Host: image.Host, Repo: image.Repo, Tag: image.Tag}
		rm, err := dockerd.ImagesWithTagConstraint(log, i.HostRepo(), i.Tag)
		if err != nil {
			return err
		}
		if err := dockerd.Rmi(log, req.GetForce(), req.GetDryRun(), rm...); err != nil {
			return err
		}
	}
	return nil
}
