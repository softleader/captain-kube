package chart

import (
	"github.com/Sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestTemplates_GenerateLoadScript(t *testing.T) {
	path, err := ioutil.TempDir(os.TempDir(), "test-generate-script-")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(path)

	cmd := exec.Command("helm", "create", "foo")
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	cmd = exec.Command("helm", "package", "foo")
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&utils.PlainFormatter{})
	log.SetLevel(logrus.DebugLevel)

	log.Debugf("ddddddug")

	tpl, err := LoadDir(log, filepath.Join(path, "foo-0.1.0.tgz"))
	if err != nil {
		t.Error(err)
	}
	if err = tpl.GenerateLoadScript(log); err != nil {
		t.Error(err)
	}
}
