package server

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/utils"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (s *CaptainServer) InstallChart(req *proto.InstallChartRequest, stream proto.Captain_InstallChartServer) error {
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

	tmp, err := ioutil.TempDir(os.TempDir(), "install-chart-icp-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	chartPath := filepath.Join(tmp, req.GetChart().GetFileName())
	if err := ioutil.WriteFile(chartPath, req.GetChart().GetContent(), 0644); err != nil {
		return err
	}
	i, err := chart.NewInstaller(s.K8s, req.GetTiller(), chartPath)
	if err != nil {
		return err
	}

	if err := i.Install(log); err != nil {
		return err
	}

	if req.GetSync() {
		endpoints, err := s.lookupCaplet(log, req.GetColor())
		if err != nil {
			return err
		}
		tpls, err := chart.LoadArchive(log, chartPath)
		if err != nil {
			return err
		}
		log.Debugf("%v template(s) loaded\n", len(tpls))
		log.SetNoLock()
		endpoints.Each(func(e *caplet.Endpoint) {
			if err := e.PullImage(log, newPullImageRequest(tpls, req.GetRetag(), req.GetRegistryAuth()), req.GetTimeout()); err != nil {
				log.Error(err)
			}
		})
	}
	return nil
}
