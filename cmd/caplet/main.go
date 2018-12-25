package main

import (
	"fmt"
	"github.com/softleader/captain-kube/cmd/caplet/app"
	"github.com/softleader/captain-kube/pkg/ver"
	"os"
)

var (
	version string
	commit  string
	date    string
)

func main() {
	metadata := ver.NewBuildMetadata(version, commit, date)
	command := app.NewCapletCommand(metadata)
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
