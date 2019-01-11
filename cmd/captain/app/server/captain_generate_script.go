package server

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func (s *CaptainServer) GenerateScript(req *captainkube_v2.GenerateScriptRequest, stream captainkube_v2.Captain_GenerateScriptServer) error {
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

	log.Debugf("chart: %s", req.GetChart().GetFileName())
	log.Debugf("retag: %+v", req.GetRetag())
	log.Debugf("pull: %+v", req.GetPull())
	log.Debugf("load: %+v", req.GetLoad())
	log.Debugf("save: %+v", req.GetSave())

	tmp, err := ioutil.TempDir(os.TempDir(), "generate-script-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	chartPath := filepath.Join(tmp, req.GetChart().GetFileName())
	if err := ioutil.WriteFile(chartPath, req.GetChart().GetContent(), 0644); err != nil {
		return err
	}

	tpls, err := chart.LoadArchive(log, chartPath)
	if err != nil {
		return err
	}
	log.Debugf("%v template(s) loaded\n", len(tpls))

	if from, to := strings.TrimSpace(req.GetRetag().GetFrom()), strings.TrimSpace(req.GetRetag().GetTo()); from != "" && to != "" {
		b, err := tpls.GenerateReTagScript(from, to)
		if err != nil {
			return err
		}
		log.Out.Write(b)
	}

	if req.Pull {
		b, err := tpls.GeneratePullScript()
		if err != nil {
			return err
		}
		log.Out.Write(b)
	}

	if req.Load {
		b, err := tpls.GenerateLoadScript()
		if err != nil {
			return err
		}
		log.Out.Write(b)
	}

	if req.Save {
		b, err := tpls.GenerateSaveScript()
		if err != nil {
			return err
		}
		log.Out.Write(b)
	}

	return nil
}
