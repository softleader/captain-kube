package server

import (
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (s *CaptainServer) InstallChart(req *proto.InstallChartRequest, stream proto.Captain_InstallChartServer) error {
	tmp, err := ioutil.TempDir(os.TempDir(), "install-chart-icp-")
	if err != nil {
		return err
	}
	chartPath := filepath.Join(tmp, req.GetChart().GetFileName())
	if err := ioutil.WriteFile(chartPath, req.GetChart().GetContent(), 0644); err != nil {
		return err
	}
	i, err := chart.NewInstaller(s.K8s, req.GetTiller(), chartPath)
	if err != nil {
		return err
	}

	sout := sio.NewStreamWriter(func(p []byte) error {
		return stream.Send(&proto.ChunkMessage{
			Msg: p,
		})
	})

	if err := i.Install(sout); err != nil {
		return err
	}

	if req.GetSync() {
		endpoints, err := s.lookupCaplets()
		if err != nil {
			return err
		}
		tpls, err := chart.LoadArchive(s.Out, chartPath)
		if err != nil {
			return err
		}
		caplet.PullImage(s.Out, endpoints, newPullImageRequest(tpls, req.GetRegistryAuth()), req.GetTimeout())
	}
	return nil
}
