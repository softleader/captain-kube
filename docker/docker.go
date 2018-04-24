package docker

import (
	"github.com/softleader/captain-kube/sh"
	"io/ioutil"
	"path/filepath"
	"fmt"
	"os"
	"gopkg.in/yaml.v2"
	"github.com/softleader/captain-kube/chart"
)

func PullImage(opts *sh.Options, chart string) error {
	dir, err := ioutil.TempDir("/tmp", "extract-tgz-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir) // clean up

	// untar
	_, _, err = sh.C(opts, "tar zxvf", chart, "-C", dir, "--strip 1")
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}

	// Render chart templates locally
	rendered := filepath.Join(dir, "rendered")
	_, _, err = sh.C(opts, "mkdir -p", rendered, "&& helm template --output-dir", rendered, dir)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = filepath.Walk(rendered, func(path string, info os.FileInfo, err error) error {
		return pull(opts, path, info, err)
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func pull(opts *sh.Options, path string, f os.FileInfo, err error) error {
	if !f.IsDir() && filepath.Ext(path) == ".yaml" {
		fmt.Printf("pull: %s\n", path)
		in, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Println(string(in))
		t := chart.Template{}
		yaml.Unmarshal(in, &t)
		for _, c := range t.Spec.Template.Spec.Containers {
			sh.C(opts, "docker pull", c.Image)
		}
	}
	return nil
}
