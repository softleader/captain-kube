package server

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/utils"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (s *CaptainServer) Rmc(req *captainkube_v2.RmcRequest, stream captainkube_v2.Captain_RmcServer) error {
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
		if err := e.Rmi(log, newRmiRequest(
			tpls,
			req.GetRetag(),
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

func newRmiRequest(tpls chart.Templates, retag *captainkube_v2.ReTag, constraint string, verbose, force, dryRun bool) (req *captainkube_v2.RmiRequest) {
	req = &captainkube_v2.RmiRequest{
		Verbose: verbose,
		Force:   force,
		DryRun:  dryRun,
	}
	for _, tpl := range tpls {
		for _, img := range tpl {
			i := &captainkube_v2.Image{
				Host: img.Host,
				Repo: img.Repo,
				Tag:  constraint + img.Tag,
			}
			if i.Host == retag.From {
				i.Host = retag.To
			}
			req.Images = append(req.Images, i)
		}
	}
	return
}
