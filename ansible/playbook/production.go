package playbook

import (
	"encoding/json"
	"github.com/softleader/captain-kube/charts"
)

type Production struct {
	Inventory      string         `json:"inventory,omitempty"`
	Tags           []string       `json:"tags"`
	Namespace      string         `json:"namespace"`
	Version        string         `json:"version"`
	Chart          string         `json:"chart,omitempty"`
	ChartPath      string         `json:"-"`
	Verbose        bool           `json:"verbose"`
	Images         []charts.Image `json:"-"`
	SourceRegistry string         `json:"sourceRegistry" yaml:"sourceRegistry" `
	Registry       string         `json:"registry"`
}

func NewProduction() *Production {
	return &Production{}
}

func (b Production) Yaml() []string {
	return []string{"production.yml"}
}

func (b Production) I() string {
	return b.Inventory
}

func (b Production) T() []string {
	return b.Tags
}

func (b Production) E() string {
	e := make(map[string]interface{})
	e["version"] = b.Version
	e["chart"] = b.Chart
	e["chart_path"] = b.ChartPath
	e["namespace"] = b.Namespace
	e["images"] = b.Images
	e["registry"] = b.Registry
	bs, _ := json.Marshal(e)
	return string(bs)
}

func (b Production) V() bool {
	return b.Verbose
}
