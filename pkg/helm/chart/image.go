package chart

import (
	"bytes"
	"fmt"
	"github.com/blang/semver"
	"strings"
	"text/template"
)

// DefaultTag 預設的 tag
const DefaultTag = "latest"

// Image 封裝了 docker image 的相關資訊
type Image struct {
	Host string // e.g. hub.softleader.com.tw
	Repo string // e.g. captain-kube
	Tag  string // e.g. latest
}

// HostRepo 回傳 image 的 host/repot
func (i *Image) HostRepo() string {
	var buf bytes.Buffer
	if i.Host != "" {
		buf.WriteString(fmt.Sprintf("%s/", i.Host))
	}
	buf.WriteString(i.Repo)
	return buf.String()
}

// Name 回傳 image repo:tag
func (i *Image) Name() string {
	var buf bytes.Buffer
	buf.WriteString(i.Repo)
	if i.Tag != "" {
		buf.WriteString(fmt.Sprintf(":%s", i.Tag))
	}
	return buf.String()
}

// String 回傳 image 的完整名稱: host/repo:tag
func (i *Image) String() string {
	var buf bytes.Buffer
	if i.Host != "" {
		buf.WriteString(fmt.Sprintf("%s/", i.Host))
	}
	buf.WriteString(i.Name())
	return buf.String()
}

// NewImage 建立 image 物件
func NewImage(img string) (i *Image) {
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

// ReTag 比對 image 的 host, 若符合則將 host 從 from 更換成 to
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

var templateFuncs = template.FuncMap{
	"replace": strings.Replace,
}

// CheckTag 檢查 tag 是否符合傳入的 constraint
func (i *Image) CheckTag(other string) (bool, error) {
	r, err := semver.ParseRange(other)
	if err != nil {
		return false, err
	}
	v, err := semver.Parse(i.Tag)
	if err != nil {
		return false, err
	}
	return r(v), nil
}
