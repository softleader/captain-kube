package tgz

import (
	"github.com/softleader/captain-kube/sh"
	"os"
	"path"
)

func Extract(opts *sh.Options, src, dest string) (err error) {
	if _, err := os.Stat(path.Join(dest, "Chart.yaml")); os.IsNotExist(err) {
		_, _, err = sh.C(opts, "tar zxvf", src, "-C", dest, "--strip 1")
	}
	return
}
