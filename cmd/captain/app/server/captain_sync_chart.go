package server

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	pb "github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/utils"
	"io/ioutil"
	"os"
	"path/filepath"
)

// SyncChart 將 image 同步到所有 caplet 上, 以避免直接啟動 chart 時拉 images 到 timeout
func (s *CaptainServer) SyncChart(req *pb.SyncChartRequest, stream pb.Captain_SyncChartServer) error {
	log := logrus.New()
	log.SetOutput(sio.NewStreamWriter(func(p []byte) error {
		return stream.Send(&pb.ChunkMessage{
			Msg: p,
		})
	}))
	log.SetFormatter(&utils.PlainFormatter{})
	if req.GetVerbose() {
		log.SetLevel(logrus.DebugLevel)
	}

	tmp, err := ioutil.TempDir(os.TempDir(), "sync-chart-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	chartPath := filepath.Join(tmp, req.GetChart().GetFileName())
	if err := saveChart(req.GetChart(), chartPath); err != nil {
		return err
	}
	log.Printf("Syncing images to all kubernetes worker nodes..")
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
	timeout := dur.Parse(req.GetTimeout())
	endpoints.Each(func(e *caplet.Endpoint) {
		if err := e.CallPullImage(log, newPullImageRequest(tpls, req.GetRetag(), req.GetRegistryAuth()), timeout); err != nil {
			log.Error(err)
		}
	})
	return nil
}
