package helm

import (
	"github.com/softleader/captain-kube/sh"
	"io/ioutil"
)

func Template(opts *sh.Options, chart string) (rendered string, err error) {
	rendered, err = ioutil.TempDir(chart, "")
	if err != nil {
		return
	}
	_, _, err = sh.C(opts, "helm template --output-dir", rendered, chart)
	return
}
