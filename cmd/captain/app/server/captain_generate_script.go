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

func (s *CaptainServer) GenerateScript(req *proto.GenerateScriptRequest, stream proto.Captain_GenerateScriptServer) error {
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
		if b, err := tpls.GenerateReTagScript(from, to); err != nil {
			return err
		} else {
			log.Out.Write(b)
		}
	}

	if req.Pull {
		if b, err := tpls.GeneratePullScript(); err != nil {
			return err
		} else {
			log.Out.Write(b)
		}
	}

	if req.Load {
		if b, err := tpls.GenerateLoadScript(); err != nil {
			return err
		} else {
			log.Out.Write(b)
		}
	}

	if req.Save {
		if b, err := tpls.GenerateSaveScript(); err != nil {
			return err
		} else {
			log.Out.Write(b)
		}
	}

	return nil
}
