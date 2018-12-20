package main

import (
	"strings"
)

const (
	unreleased  = "unreleased"
)

var version string

func ver() string {
	if v := strings.TrimSpace(version); v != "" {
		return v
	} else {
		return unreleased
	}
}
