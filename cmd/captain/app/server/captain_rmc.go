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

// Rmc 將上傳的 chart 中的 images 從 caplet 的 docker 中刪除
func (s *CaptainServer) Rmc(req *pb.RmcRequest, stream pb.Captain_RmcServer) error {
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

	endpoints, err := s.lookupCaplet(log, req.GetColor())
	if err != nil {
		return err
	}

	tmp, err := ioutil.TempDir(os.TempDir(), "rmc-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	chartPath := filepath.Join(tmp, req.GetChart().GetFileName())
	if err := saveChart(req.GetChart(), chartPath); err != nil {
		return err
	}
	tpls, err := chart.LoadArchive(log, chartPath, req.GetSet()...)
	if err != nil {
		return err
	}
	log.Debugf("%v template(s) loaded\n", len(tpls))
	log.SetNoLock()
	timeout := dur.Parse(req.GetTimeout())
	endpoints.Each(func(e *caplet.Endpoint) {
		if err := e.CallRmi(log, newRmiRequest(
			tpls,
			req.GetConstraint(),
			req.GetVerbose(),
			req.GetForce(),
			req.GetDryRun()),
			timeout,
		); err != nil {
			log.Error(err)
		}
	})
	return nil
}

func newRmiRequest(tpls chart.Templates, constraint string, verbose, force, dryRun bool) (req *pb.RmiRequest) {
	req = &pb.RmiRequest{
		Verbose: verbose,
		Force:   force,
		DryRun:  dryRun,
	}
	for _, tpl := range tpls {
		for _, img := range tpl {
			i := &pb.Image{
				Host: img.Host,
				Repo: img.Repo,
				Tag:  constraint + img.Tag,
			}
			req.Images = append(req.Images, i)
		}
	}
	return
}
