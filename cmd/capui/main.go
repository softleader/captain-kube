package main

import (
	"fmt"
	"github.com/softleader/captain-kube/cmd/capui/app"
	ver "github.com/softleader/captain-kube/pkg/version"
	"os"
)

var (
	version string
	commit  string
)

func main() {
	metadata := ver.NewBuildMetadata(version, commit)
	command := app.NewCapUICommand(metadata)
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}