package main

import (
	"github.com/softleader/captain-kube/cmd/capctl/app"
	ver "github.com/softleader/captain-kube/pkg/version"
	"os"
)

var (
	version string
	commit  string
	date    string
)

func main() {
	metadata := ver.NewBuildMetadata(version, commit, date)
	if command, err := app.NewRootCmd(os.Args[1:], metadata); err != nil {
		panic(err)
	} else {
		if err := command.Execute(); err != nil {
			os.Exit(1)
		}

	}

}
