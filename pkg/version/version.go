package version

import (
	"fmt"
	"strings"
)

const (
	unreleased = "unreleased"
	unknown    = "unknown"
)

type BuildMetadata struct {
	GitVersion string
	GitCommit  string
}

func NewBuildMetadata(version, commit, date string) (b *BuildMetadata) {
	b = &BuildMetadata{
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

func (b *BuildMetadata) String() string {
	return fmt.Sprintf("%s+%s", b.GitVersion, b.GitCommit[:7])
}

func (b *BuildMetadata) FullString() string {
	return fmt.Sprintf("%#v", b)
}
