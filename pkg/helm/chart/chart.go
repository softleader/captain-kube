package chart

import (
	"github.com/softleader/captain-kube/pkg/arc"
	"github.com/softleader/captain-kube/pkg/helm"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const template = "t"

type Templates map[string][]*Image

func LoadArchive(out io.Writer, archivePath string) (tpls Templates, err error) {
	path, err := ioutil.TempDir(os.TempDir(), "load-archive")
	if err != nil {
		return
	}
	defer os.RemoveAll(path)
	extractPath := filepath.Join(path, filepath.Base(archivePath))
	if err = arc.Extract(out, archivePath, extractPath); err != nil {
		return
	}
	tplPath := filepath.Join(archivePath, template)
	if err = helm.Template(out, extractPath, tplPath); err != nil {
		return
	}
	tpls, err = LoadDir(out, tplPath)
	return
}

func LoadDir(_ io.Writer, chartPath string) (tpls Templates, err error) {
	tpls = make(map[string][]*Image)
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
