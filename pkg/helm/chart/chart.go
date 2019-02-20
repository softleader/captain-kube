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

// Templates 等於 chart/template
type Templates map[string][]*Image

// Size 回傳 Templates 數量
func (t *Templates) Size() int {
	return len(reflect.ValueOf(t).MapKeys())
}

// LoadArchiveBytes 從 chart 壓縮檔的 bytes 載入 template
func LoadArchiveBytes(log *logrus.Logger, filename string, data []byte, set ...string) (tpls Templates, err error) {
	tmp, err := ioutil.TempDir(os.TempDir(), "load-bytes-")
	if err != nil {
		return
	}
	defer os.Remove(tmp)
	archive := filepath.Join(tmp, filename)
	if err := ioutil.WriteFile(archive, data, 0644); err != nil {
		return nil, err
	}
	return LoadArchive(log, archive, set...)
}

// LoadArchive 從 chart 壓縮檔的載入 template
func LoadArchive(log *logrus.Logger, archivePath string, set ...string) (tpls Templates, err error) {
	tmp, err := ioutil.TempDir(os.TempDir(), "load-archive-")
	if err != nil {
		return
	}
	defer os.RemoveAll(tmp)
	extractPath := filepath.Join(tmp, filepath.Base(archivePath))
	if err = arc.Extract(log, archivePath, extractPath); err != nil {
		return
	}
	chartDir, err := findFirstDir(extractPath)
	if err != nil {
		return
	}
	tplPath := filepath.Join(chartDir, templateDir)
	if err = helm.Template(log, chartDir, tplPath, set...); err != nil {
		return
	}
	tpls, err = LoadDir(log, tplPath)
	log.Debugf("%v template(s) loaded", len(tpls))
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

// LoadDir 從指定的 path 載入 template
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
			i = append(i, NewImage(c.Image))
		}
	}
	return i, nil
}

// TemplateYAML 定義了 template yaml 要讀取的格式
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
