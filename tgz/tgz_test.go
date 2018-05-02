package tgz

import (
	"testing"
	"github.com/softleader/captain-kube/sh"
)

func TestExtract(t *testing.T) {
	opts := sh.Options{}
	Extract(&opts, "/Users/Matt/GitHub/softleader-charts/softleader-gardenia/softleader-gardenia-1.2.15.tgz", "/Users/Matt/tmp/a")
}
