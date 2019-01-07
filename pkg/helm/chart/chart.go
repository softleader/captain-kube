package chart

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/arc"
	"github.com/softleader/captain-kube/pkg/helm"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const templateDir = "t"

type Templates map[string][]*Image

func (t *Templates) Size() int {
	return len(reflect.ValueOf(t).MapKeys())
}

func LoadBytes(log *logrus.Logger, content []byte) (tpls Templates, err error) {
	tf, err := ioutil.TempFile(os.TempDir(), "load-bytes-")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tf.Name())
	if _, err := tf.Write(content); err != nil {
		return nil, err
	}
	return LoadArchive(log, tf.Name())
}

func LoadArchive(log *logrus.Logger, archivePath string) (tpls Templates, err error) {
	path, err := ioutil.TempDir(os.TempDir(), "load-archive-")
	if err != nil {
		return
	}
	defer os.RemoveAll(path)
	extractPath := filepath.Join(path, filepath.Base(archivePath))
	if err = arc.Extract(log, archivePath, extractPath); err != nil {
		return
	}
	chartDir, err := findFirstDir(extractPath)
	if err != nil {
		return
	}
	tplPath := filepath.Join(chartDir, templateDir)
	if err = helm.Template(log, chartDir, tplPath); err != nil {
		return
	}
	tpls, err = LoadDir(log, tplPath)
	log.Debugf("%v template(s) loaded\n", len(tpls))
	return
}

func findFirstDir(path string) (string, error) {
	extractDir, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}
	for _, f := range extractDir {
		if f.IsDir() {
			return filepath.Join(path, f.Name()), err
		}
	}
	return "", fmt.Errorf("no dir found in %q", path)
}

func LoadDir(log *logrus.Logger, chartPath string) (tpls Templates, err error) {
	tpls = make(map[string][]*Image)
	log.Debugf("loading helm template: %s", chartPath)
	err = filepath.Walk(chartPath, func(path string, info os.FileInfo, err error) error {
		i, err := image(path, info)
		if len(i) > 0 {
			src := strings.Replace(path, chartPath+"/", "", -1)
			tpls[src] = i
		}
		return err
	})
	return
}

func image(path string, f os.FileInfo) ([]*Image, error) {
	var i []*Image
	if !f.IsDir() && filepath.Ext(path) == ".yaml" {
		in, err := ioutil.ReadFile(path)
		if err != nil {
			return i, err
		}
		t := TemplateYAML{}
		yaml.Unmarshal(in, &t)
		for _, c := range t.Spec.SpecTemplate.Spec.Containers {
			i = append(i, newImage(c.Image))
		}
	}
	return i, nil
}

type TemplateYAML struct {
	Spec struct {
		SpecTemplate struct {
			Spec struct {
				Containers []struct {
					Name  string `yaml:"name"`
					Image string `yaml:"image"`
				} `yaml:"containers"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}
