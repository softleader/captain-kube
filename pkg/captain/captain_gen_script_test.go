package captain

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/proto"
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

	chart, err := ioutil.ReadFile(filepath.Join(path, "foo-0.1.0.tgz"))
	if err != nil {
		t.Error(err)
	}

	log := logrus.New()
	log.SetFormatter(&utils.PlainFormatter{})
	err = CallGenerateScript(log, fmt.Sprintf("%v:%v", endpoint, port), &captainkube_v2.GenerateScriptRequest{
		Chart: &captainkube_v2.Chart{
			Content:  chart,
			FileName: "foo-0.1.0.tgz",
		},
		Pull:    true,
		Retag:   &captainkube_v2.ReTag{},
		Save:    true,
		Load:    true,
		Verbose: true,
	}, dur.Parse(dur.DefaultDeadline))
	if err != nil {
		t.Error(err)
	}
}
