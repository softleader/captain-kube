package captain

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/utils/tcp"
	"testing"
)

func TestVersion(t *testing.T) {
	endpoint := "192.168.1.93"
	port := DefaultPort
	if !tcp.IsReachable(endpoint, port, 3) {
		t.Skipf("endpoint %s:%v is not reachable", endpoint, port)
	}
	addr := fmt.Sprintf("%v:%v", endpoint, port)
	if err := Version(logrus.StandardLogger(), addr, false, false, dur.Parse(dur.DefaultDeadline)); err != nil {
		t.Error(err)
	}
}
