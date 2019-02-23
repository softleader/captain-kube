package main

import (
	"fmt"
	"github.com/softleader/captain-kube/cmd/caplet/app"
	ver "github.com/softleader/captain-kube/pkg/release"
	// importing automaxprocs automatically adjusts GOMAXPROCS.
	_ "go.uber.org/automaxprocs"
	"os"
)

var (
	version string
	commit  string
)

func main() {
	metadata := ver.NewMetadata(version, commit)
	command := app.NewCapletCommand(metadata)
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
