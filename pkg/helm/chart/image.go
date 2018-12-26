package chart

import (
	"bytes"
	"fmt"
	"strings"
)

type Image struct {
	Host string // e.g. hub.softleader.com.tw
	Repo string // e.g. captain-kube
	Tag  string // e.g. latest
}

func (i *Image) HostRepo() string {
	var buf bytes.Buffer
	if i.Host != "" {
		buf.WriteString(fmt.Sprintf("%s/", i.Host))
	}
	buf.WriteString(i.Repo)
	return buf.String()
}

func (i *Image) Name() string {
	var buf bytes.Buffer
	buf.WriteString(i.Repo)
	if i.Tag != "" {
		buf.WriteString(fmt.Sprintf("/%s", i.Tag))
	}
	return buf.String()
}

func (i *Image) String() string {
	var buf bytes.Buffer
	if i.Host != "" {
		buf.WriteString(fmt.Sprintf("%s/", i.Host))
	}
	buf.WriteString(i.Name())
	return buf.String()
}

func newImage(img string) (i *Image) {
	var name string
	i = &Image{}
	if strings.ContainsAny(img, "/") {
		i.Host = before(img, "/")
		name = after(img, "/")
	} else {
		name = img
	}
	if strings.ContainsAny(name, ":") {
		i.Repo = before(name, ":")
		i.Tag = after(name, ":")
	} else {
		i.Repo = name
	}
	return
}

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
	return value[adjustedPos:]
}
