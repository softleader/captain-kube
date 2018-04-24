package helm

import (
	"path/filepath"
	"github.com/softleader/captain-kube/sh"
)

const renderDir = "rendered"

func Template(opts *sh.Options, chart string) (rendered string, err error) {
	rendered = filepath.Join(chart, renderDir)
	_, _, err = sh.C(opts, "mkdir -p", rendered, "&& helm template --output-dir", rendered, chart)
	return
}
