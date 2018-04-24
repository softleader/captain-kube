package playbook

type Release struct {
	Inventory  string   `json:"inventory"`
	Tags       []string `json:"tags"`
	Namespace  string   `json:"namespace"`
	Version    string   `json:"version,omitempty"`
	Chart      string   `json:"chart,omitempty"`
	ChartPath  string   `json:"-"`
	Script     string   `json:"-"`
	ScriptPath string   `json:"-"`
	Verbose    bool     `json:"verbose"`
	Registry   string   `json:"registry"`
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

func (b Release) E() []string {
	return []string{"version=" + b.Version, "chart=" + b.Chart, "chart_path=" + b.ChartPath, "script=" + b.Script, "script_path=" + b.ScriptPath, "namespace=" + b.Namespace}
}

func (b Release) V() bool {
	return b.Verbose
}
