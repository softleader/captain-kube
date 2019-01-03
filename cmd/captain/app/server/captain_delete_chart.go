package server

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/utils"
)

func (s *CaptainServer) DeleteChart(req *proto.DeleteChartRequest, stream proto.Captain_DeleteChartServer) error {
	log := logrus.New()
	log.SetOutput(sio.NewStreamWriter(func(p []byte) error {
		return stream.Send(&proto.ChunkMessage{
			Msg: p,
		})
	}))
	log.SetFormatter(&utils.PlainFormatter{})
	if req.GetVerbose() {
		log.SetLevel(logrus.DebugLevel)
	}
	d, err := chart.NewDeleter(s.K8s, req.GetTiller(), req.GetChartName(), req.GetChartVersion())
	if err != nil {
		return err
	}
	if err := d.Delete(log); err != nil {
		return err
	}
	return nil
}
