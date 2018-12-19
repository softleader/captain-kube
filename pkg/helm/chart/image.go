package chart

import "strings"

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
	Repo string // e.g. captain-kube
	Tag  string // latest
}

func newImage(img string) (i *Image) {
	i = &Image{
		Host: before(img, "/"),
		Name: after(img, "/"),
	}
	if n := i.Name; strings.ContainsAny(n, ":") {
		s := strings.Split(n, ":")
		i.Repo = s[0]
		i.Tag = s[1]
	} else {
		i.Repo = n
	}
	return
}

type Images map[string][]*Image

func (i *Image) ReTag(from, to string) {
	if from != "" && to != "" && i.Host == from {
		i.Host = to
	}
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