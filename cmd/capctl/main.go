package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/cmd/capctl/app"
	ver "github.com/softleader/captain-kube/pkg/version"
	"os"
)

var (
	version string
	commit  string
)

func main() {
	metadata := ver.NewBuildMetadata(version, commit)
	if command, err := app.NewRootCmd(os.Args[1:], metadata); err != nil {
		logrus.Error(err)
		os.Exit(1)
	} else if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
