package chart

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/utils/tcp"
	"testing"
)

func TestIcpInstaller_Install(t *testing.T) {
	endpoint := "192.168.1.93"
	port := captain.DefaultPort
	if !tcp.IsReachable(endpoint, port, 3) {
		t.Skipf("endpoint %s:%v is not reachable", endpoint, port)
	}

	buf := bytes.NewBuffer(nil)
	log := logrus.New()
	log.SetOutput(buf)
	logrus.SetLevel(logrus.DebugLevel)
	if err := loginBxPr(logrus.StandardLogger(), endpoint, "admin", "admin", "mycluster Account", true); err != nil {
		t.Error(err)
	}
}
