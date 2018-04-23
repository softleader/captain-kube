package docker

import (
	"github.com/softleader/captain-kube/sh"
	"testing"
	"path/filepath"
	"os"
)

func TestPullImage(t *testing.T) {
	opts := sh.Options{
		Pwd:     "/Users/Matt/tmp",
		Verbose: true,
	}
	PullImage(&opts, "/Users/Matt/tmp/softleader-gardenia-1.1.6.tgz")
}

func TestCollectYaml(t *testing.T) {
	opts := sh.Options{
		Pwd:     "/Users/Matt/tmp",
		Verbose: true,
	}
	err := filepath.Walk("/tmp/extract-tgz-879389998/rendered", func(path string, info os.FileInfo, err error) error {
		return pull(&opts, path, info, err)
	})
	if err != nil {
		t.Error(err)
	}
}
