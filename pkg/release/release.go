package release

import (
	"fmt"
	"strings"
)

const (
	unreleased = "unreleased"
	unknown    = "unknown"
)

type Metadata struct {
	GitVersion string
	GitCommit  string
}

func NewMetadata(version, commit string) (b *Metadata) {
	b = &Metadata{
		GitVersion: unreleased,
		GitCommit:  unknown,
	}
	if version = strings.TrimSpace(version); version != "" {
		b.GitVersion = version
	}
	if commit = strings.TrimSpace(commit); commit != "" {
		b.GitCommit = commit
	}
	return
}

func (b *Metadata) String() string {
	trunc := 7
	if len := len(b.GitCommit); len < 7 {
		trunc = len
	}
	return fmt.Sprintf("%s+%s", b.GitVersion, b.GitCommit[:trunc])
}

func (b *Metadata) FullString() string {
	return fmt.Sprintf("%#v", b)
}
