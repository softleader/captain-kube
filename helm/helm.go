package helm

import (
	"github.com/softleader/captain-kube/sh"
	"os"
	"path"
)

const template = "t"

func Template(opts *sh.Options, chart string) (rendered string, err error) {
	rendered = path.Join(chart, template)
	if _, err := os.Stat(rendered); os.IsNotExist(err) {
		err = os.MkdirAll(rendered, os.ModePerm)
		if err == nil {
			_, _, err = sh.C(opts, "helm template --output-dir", rendered, chart)
		}

	}
	return
}
