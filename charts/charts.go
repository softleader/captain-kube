package charts

import (
	"path/filepath"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"strings"
)

type Template struct {
	Spec struct {
		Template struct {
			Spec struct {
				Containers []struct {
					Name  string `yaml:"name"`
					Image string `yaml:"image"`
				} `yaml:"containers"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}

type Image struct {
	Registry string // hub.softleader.com.tw
	Name     string // captain-kube:latest
}

func CollectImages(chart string, filter func(string) bool) (images map[string][]Image, err error) {
	images = make(map[string][]Image)
	err = filepath.Walk(chart, func(path string, info os.FileInfo, err error) error {
		i, err := image(filter, path, info)
		if len(i) > 0 {
			images[path] = i
		}
		return err
	})
	return
}

func image(filter func(string) bool, path string, f os.FileInfo) ([]Image, error) {
	var i []Image
	if !f.IsDir() && filepath.Ext(path) == ".yaml" {
		in, err := ioutil.ReadFile(path)
		if err != nil {
			return i, err
		}
		t := Template{}
		yaml.Unmarshal(in, &t)
		for _, c := range t.Spec.Template.Spec.Containers {
			image := Image{
				Registry: before(c.Image, "/"),
				Name:     after(c.Image, "/"),
			}
			if filter(image.Registry) {
				i = append(i, image)
			}
		}
	}
	return i, nil
}

func before(value string, a string) string {
	// Get substring before a string.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
}

func after(value string, a string) string {
	// Get substring after a string.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}