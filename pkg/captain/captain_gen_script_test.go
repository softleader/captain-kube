package captain

import (
	"github.com/sirupsen/logrus"
	chart2 "github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/softleader/captain-kube/pkg/utils/command"
	"github.com/softleader/captain-kube/pkg/utils/tcp"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGenerateScript(t *testing.T) {
	endpoint := "192.168.1.93"
	port := DefaultPort
	if !tcp.IsReachable(endpoint, port, 3) {
		t.Skipf("endpoint %s:%v is not reachable", endpoint, port)
	}

	helm := "helm"
	if !command.IsAvailable(helm) {
		t.Skipf("%q command does not exist", helm)
	}

	path, err := ioutil.TempDir(os.TempDir(), "test-generate-script-")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(path)

	cmd := exec.Command(helm, "create", "foo")
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	cmd = exec.Command(helm, "package", "foo")
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	log := logrus.New()
	log.SetFormatter(&utils.PlainFormatter{})

	tpls, err := chart2.LoadArchive(logrus.StandardLogger(), filepath.Join(path, "foo-0.1.0.tgz"))
	if err != nil {
		t.Error(err)
	}

	if b, err := tpls.GeneratePullScript(); err != nil {
		t.Error(err)
	} else {
		log.Out.Write(b)
	}

	if b, err := tpls.GenerateLoadScript(); err != nil {
		t.Error(err)
	} else {
		log.Out.Write(b)
	}

	if b, err := tpls.GenerateSaveScript(); err != nil {
		t.Error(err)
	} else {
		log.Out.Write(b)
	}

}
