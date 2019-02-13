package dockerd

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"testing"
)

func TestImages(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	image := chart.NewImage("busybox:1.23.2")
	if err := Pull(logrus.StandardLogger(), *image, nil); err != nil {
		t.Skipf("maybe just can not connect to docker or internet: %s", err)
	}
	images, err := Images(logrus.StandardLogger(), image.HostRepo())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(images) == 0 {
		t.Errorf("should found more than 1 images")
	}
}
