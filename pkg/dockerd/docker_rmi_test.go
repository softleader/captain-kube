package dockerd

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"testing"
)

func TestRmi(t *testing.T) {
	image := chart.NewImage("library/busybox:1.23.2")
	if err := Pull(logrus.StandardLogger(), *image, nil); err != nil {
		t.Skipf("maybe just can not connect to docker or internet: %s", err)
	}
	err := Rmi(logrus.StandardLogger(), true, false, image)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
