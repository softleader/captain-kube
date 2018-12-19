package chart

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func CollectImages(chart string) (images Images, err error) {
	images = make(map[string][]*Image)
	err = filepath.Walk(chart, func(path string, info os.FileInfo, err error) error {
		i, err := image(path, info)
		if len(i) > 0 {
			src := strings.Replace(path, chart+"/", "", -1)
			images[src] = i
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
		t := Template{}
		yaml.Unmarshal(in, &t)
		for _, c := range t.Spec.Template.Spec.Containers {
			i = append(i, newImage(c.Image))
		}
	}
	return i, nil
}
