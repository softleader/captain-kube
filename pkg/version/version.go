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

func NewBuildMetadata(version, commit string) (b *BuildMetadata) {
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
	trunc := 7
	if len := len(b.GitCommit); len < 7 {
		trunc = len
	}
	return fmt.Sprintf("%s+%s", b.GitVersion, b.GitCommit[:trunc])
}

func (b *BuildMetadata) FullString() string {
	return fmt.Sprintf("%#v", b)
}
