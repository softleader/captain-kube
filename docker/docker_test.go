package docker

import (
	"github.com/softleader/captain-kube/sh"
	"testing"
)

func TestPullImage(t *testing.T) {
	opts := sh.Options{
		Pwd:     "/Users/Matt/tmp",
		Verbose: true,
	}
	Pull(&opts, "/Users/Matt/tmp/", "softleader-gardenia-1.1.6.tgz")
}
