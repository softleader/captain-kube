package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func (s *CaptainServer) GenerateScript(req *proto.GenerateScriptRequest, stream proto.Captain_GenerateScriptServer) error {
	log := logrus.New()
	log.SetOutput(sio.NewStreamWriter(func(p []byte) error {
		return stream.Send(&proto.ChunkMessage{
			Msg: p,
		})
	}))
	if req.GetVerbose() {
		log.SetLevel(logrus.DebugLevel)
	}

	log.Debugf("chart: %s\n", req.GetChart().GetFileName())
	log.Debugf("retag: %+v\n", req.GetRetag())
	log.Debugf("pull: %+v\n", req.GetPull())
	log.Debugf("load: %+v\n", req.GetLoad())
	log.Debugf("save: %+v\n", req.GetSave())

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
		if err := tpls.GenerateReTagScript(log, from, to); err != nil {
			return err
		}
	}

	if req.Pull {
		if err := tpls.GeneratePullScript(log); err != nil {
			return err
		}
	}

	if req.Load {
		if err := tpls.GenerateLoadScript(log); err != nil {
			return err
		}
	}

	if req.Save {
		if err := tpls.GenerateSaveScript(log); err != nil {
			return err
		}
	}

	return nil
}
