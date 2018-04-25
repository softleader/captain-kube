package playbook

import "encoding/json"

type Release struct {
	Inventory      string   `json:"inventory"`
	Tags           []string `json:"tags"`
	Namespace      string   `json:"namespace"`
	Version        string   `json:"version,omitempty"`
	Chart          string   `json:"chart,omitempty"`
	ChartPath      string   `json:"-"`
	Verbose        bool     `json:"verbose"`
	Images         []string `json:"-"`
	SourceRegistry string   `json:"source_registry"`
	Registry       string   `json:"registry"`
}

func NewRelease() *Release {
	return &Release{
		Inventory: "hosts",
		Namespace: "default",
	}
}

func (b Release) Yaml() []string {
	return []string{"release.yml"}
}

func (b Release) I() string {
	return b.Inventory
}

func (b Release) T() []string {
	return b.Tags
}

func (b Release) E() string {
	e := make(map[string]interface{})
	e["version"] = b.Version
	e["chart"] = b.Chart
	e["chart_path"] = b.ChartPath
	e["namespace"] = b.Namespace
	e["images"] = b.Images
	e["registry"] = b.Registry
	e["source_registry"] = b.SourceRegistry
	bs, _ := json.Marshal(e)
	return string(bs)
}

func (b Release) V() bool {
	return b.Verbose
}
