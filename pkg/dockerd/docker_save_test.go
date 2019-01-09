package dockerd

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSave(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "docker-save-")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(tmp)
	out := filepath.Join(tmp, "saved.tar")

	images := []*chart.Image{
		{
			Repo: "busybox",
		},
		{
			Host: "registry.au-syd.bluemix.net/softleader",
			Repo: "aiclm-backend",
			Tag:  "experimental",
		},
	}

	if err := Save(logrus.StandardLogger(), images, out, true); err != nil {
		t.Error(err)
		t.Skipf("maybe just docker not exist")
	}
}
