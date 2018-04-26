package playbook

import (
	"encoding/json"
	"github.com/softleader/captain-kube/charts"
)

type Staging struct {
	Inventory  string         `json:"inventory"`
	Tags       []string       `json:"tags"`
	Namespace  string         `json:"namespace"`
	Version    string         `json:"version"`
	Chart      string         `json:"chart,omitempty"`
	ChartPath  string         `json:"-"`
	Images     []charts.Image `json:"-"`
	Verbose    bool           `json:"verbose"`
}

func NewStaging() *Staging {
	return &Staging{
		Inventory: "hosts",
		Namespace: "default",
	}
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
	bs, _ := json.Marshal(e)
	return string(bs)
}

func (b Staging) V() bool {
	return b.Verbose
}
