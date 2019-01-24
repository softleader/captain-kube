package server

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/utils"
)

func (s *CaptainServer) DeleteChart(req *captainkube_v2.DeleteChartRequest, stream captainkube_v2.Captain_DeleteChartServer) error {
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
	k8s, err := s.k8s()
	if err != nil {
		return err
	}
	d, err := chart.NewDeleter(k8s, req.GetTiller(), req.GetChartName(), req.GetChartVersion())
	if err != nil {
		return err
	}
	if err := d.Delete(log); err != nil {
		return err
	}
	return nil
}
