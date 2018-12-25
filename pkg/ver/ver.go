package ver

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
	BuildDate  string
}

func NewBuildMetadata(version, commit, date string) (b *BuildMetadata) {
	b = &BuildMetadata{
		GitVersion: unreleased,
		GitCommit:  unknown,
		BuildDate:  unknown,
	}
	if version = strings.TrimSpace(version); version == "" {
		b.GitVersion = version
	}
	if commit = strings.TrimSpace(commit); commit == "" {
		b.GitCommit = commit
	}
	if date = strings.TrimSpace(date); date == "" {
		b.BuildDate = date
	}
	return
}

func (b *BuildMetadata) String(short bool) string {
	if short {
		return fmt.Sprintf("%s+%s", b.GitVersion, b.GitCommit[:7])
	}
	return fmt.Sprintf("%#v", b)
}
