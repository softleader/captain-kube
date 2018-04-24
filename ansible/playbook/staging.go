package playbook

type Staging struct {
	Inventory  string   `json:"inventory"`
	Tags       []string `json:"tags"`
	Namespace  string   `json:"namespace"`
	Version    string   `json:"version,omitempty"`
	Chart      string   `json:"chart,omitempty"`
	ChartPath  string   `json:"-"`
	Script     string   `json:"-"`
	ScriptPath string   `json:"-"`
	Verbose    bool     `json:"verbose"`
	DockerPull bool     `json:"dockerPull"`
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

func (b Staging) E() []string {
	return []string{"version=" + b.Version, "chart=" + b.Chart, "chart_path=" + b.ChartPath, "script=" + b.Script, "script_path=" + b.ScriptPath, "namespace=" + b.Namespace}
}

func (b Staging) V() bool {
	return b.Verbose
}
