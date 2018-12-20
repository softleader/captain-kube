package chart

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestLoadArchive(t *testing.T) {
	out := os.Stdout

	path, err := ioutil.TempDir(os.TempDir(), "load-archive-test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(path)

	cmd := exec.Command("helm", "create", "foo")
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	cmd = exec.Command("helm", "package", "foo")
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	tpls, err := LoadArchive(out, filepath.Join(path, "foo-0.1.0.tgz"))
	if err != nil {
		t.Error(err)
	}
	if len(tpls) == 0 {
		t.Error("no template found")
	}
}