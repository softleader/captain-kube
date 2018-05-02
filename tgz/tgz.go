package tgz

import (
	"github.com/softleader/captain-kube/sh"
)

func Extract(opts *sh.Options, src, dest string) (err error) {
	_, _, err = sh.C(opts, "tar zxvf", src, "-C", dest, "--strip 1")
	return
}
