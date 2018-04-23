package playbook

type Release struct {
	inventory string   `json:"inventory"`
	tags      []string `json:"tags"`
	verbose   bool     `json:"verbose"`
	version   string   `json:"version,omitempty"`
	chart     string   `json:"chart,omitempty"`
}

func NewRelease() *Release {
	return &Release{
		inventory: "hosts",
	}
}

func (b Release) Yaml() []string {
	return []string{"release.yml"}
}

func (b Release) Inventory() string {
	return b.inventory
}

func (b Release) Tags() []string {
	return b.tags
}

func (b Release) ExtraVars() []string {
	return []string{"version=" + b.version, "chart=" + b.chart}
}

func (b Release) Verbose() bool {
	return b.verbose
}
