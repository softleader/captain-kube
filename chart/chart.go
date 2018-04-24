package chart

import (
	"path/filepath"
	"os"
	"fmt"
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
	Name       string // hub.softleader.com.tw/captain-kube:latest
	RemoteName string // captain-kube:latest
	Registry   string
}

type collect struct {
	Images []Image
}

func CollectImages(chart string) ([]Image, error) {
	collect := collect{}
	err := filepath.Walk(chart, func(path string, info os.FileInfo, err error) error {
		return image(&collect, path, info, err)
	})
	if err != nil {
		return nil, err
	}
	return collect.Images, nil
}

func image(collect *collect, path string, f os.FileInfo, err error) error {
	if !f.IsDir() && filepath.Ext(path) == ".yaml" {
		fmt.Printf("pull: %s\n", path)
		in, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Println(string(in))
		t := Template{}
		yaml.Unmarshal(in, &t)
		for _, c := range t.Spec.Template.Spec.Containers {
			image := Image{
				Name:       c.Image,
				RemoteName: after(c.Image, "/"),
			}
			collect.Images = append(collect.Images, image)
		}
	}
	return nil
}

//func before(value string, a string) string {
//	// Get substring before a string.
//	pos := strings.Index(value, a)
//	if pos == -1 {
//		return ""
//	}
//	return value[0:pos]
//}

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
