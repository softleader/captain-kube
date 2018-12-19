package chart

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
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
	Host string // e.g. hub.softleader.com.tw
	Name string // e.g. captain-kube:latest
}

type Images map[string][]Image

func (i *Image) ReTag(from, to string) {
	if from != "" && to != "" && i.Host == from {
		i.Host = to
	}
}

func CollectImages(chart string, filter func(Image) bool, comsumer func(Image) Image) (images Images, err error) {
	images = make(map[string][]Image)
	err = filepath.Walk(chart, func(path string, info os.FileInfo, err error) error {
		i, err := image(path, info, filter, comsumer)
		if len(i) > 0 {
			src := strings.Replace(path, chart+"/", "", -1)
			images[src] = i
		}
		return err
	})
	return
}

func image(path string, f os.FileInfo, filter func(Image) bool, comsumer func(Image) Image) ([]Image, error) {
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
				Host: before(c.Image, "/"),
				Name: after(c.Image, "/"),
			}
			if filter(image) {
				i = append(i, comsumer(image))
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
