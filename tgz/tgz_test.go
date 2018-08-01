package tgz

import (
	"testing"
	"github.com/softleader/captain-kube/sh"
)

func TestExtract(t *testing.T) {
	opts := sh.Options{}
	Extract(&opts, "/Users/Matt/GitHub/softleader-charts/softleader-product-1.1.5 (1).tgz", "/Users/Matt/tmp/a")
}
