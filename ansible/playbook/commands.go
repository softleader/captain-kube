package playbook

type Commands struct {
	Inventory string   `json:"inventory"`
	Commands  []string `json:"commands"`
	Verbose   Verbose  `json:"verbose"`
}

func MewCommands() *Commands {
	return &Commands{}
}

func (b Commands) Yaml() []string {
	return []string{"commands.yml"}
}

func (b Commands) I() string {
	return b.Inventory
}

func (b Commands) T() []string {
	return b.Commands
}

func (b Commands) E() string {
	return ""
}

func (b Commands) V() string {
	return b.Verbose.String()
}
