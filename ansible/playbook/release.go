package playbook

type Release struct {
	Inventory string   `json:"inventory"`
	Tags      []string `json:"tags"`
	Namespace string   `json:"namespace"`
	Version   string   `json:"version,omitempty"`
	Chart     string   `json:"chart,omitempty"`
	Verbose   bool     `json:"verbose"`
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
	return []string{"version=" + b.Version, "chart=" + b.Chart, "namespace=" + b.Namespace}
}

func (b Release) V() bool {
	return b.Verbose
}
