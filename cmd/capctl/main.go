package main

import (
	"github.com/softleader/captain-kube/cmd/capctl/app"
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
	command := app.NewRootCmd(os.Args[1:], metadata)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
