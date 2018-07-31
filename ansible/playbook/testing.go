package playbook

type Testing struct {
	Inventory string  `json:"inventory"`
	Verbose   Verbose `json:"verbose"`
}

func NewTesting() *Testing {
	return &Testing{}
}

func (b Testing) Yaml() []string {
	return []string{"testing.yml"}
}

func (b Testing) I() string {
	return b.Inventory
}

func (b Testing) T() []string {
	return []string{}
}

func (b Testing) E() string {
	return ""
}

func (b Testing) V() string {
	return b.Verbose.String()
}
