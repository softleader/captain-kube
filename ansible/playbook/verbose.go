package playbook

import "strings"

type Verbose struct {
	Enabled bool
	Level   int
}

func (v *Verbose) String() (s string) {
	if v.Enabled && v.Level > 0 {
		s = "-" + strings.Repeat("v", v.Level)
	}
	return
}
