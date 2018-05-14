package playbook

import (
	"encoding/json"
	"github.com/softleader/captain-kube/charts"
)

type Staging struct {
	Inventory      string         `json:"inventory,omitempty"`
	Tags           []string       `json:"tags"`
	Namespace      string         `json:"namespace"`
	Version        string         `json:"version"`
	Chart          string         `json:"chart,omitempty"`
	ChartPath      string         `json:"-"`
	Images         []charts.Image `json:"-"` // chart 中所有的 image
	RetagImages    []charts.Image `json:"-"` // 需要被 retag 的 image
	Verbose        bool           `json:"verbose"`
	SourceRegistry string         `json:"sourceRegistry" yaml:"sourceRegistry" `
	Registry       string         `json:"registry"`
}

func NewStaging() *Staging {
	return &Staging{}
}

func (b Staging) Yaml() []string {
	return []string{"staging.yml"}
}

func (b Staging) I() string {
	return b.Inventory
}

func (b Staging) T() []string {
	return b.Tags
}

func (b Staging) E() string {
	e := make(map[string]interface{})
	e["version"] = b.Version
	e["chart"] = b.Chart
	e["chart_path"] = b.ChartPath
	e["namespace"] = b.Namespace
	e["images"] = b.Images
	e["retag_images"] = b.RetagImages
	e["registry"] = b.Registry
	bs, _ := json.Marshal(e)
	return string(bs)
}

func (b Staging) V() bool {
	return b.Verbose
}
