package main

import (
	"fmt"
	"github.com/softleader/captain-kube/cmd/capctl/app"
	"github.com/softleader/captain-kube/pkg/release"
	"os"
)

var (
	version string
	commit  string
)

func main() {
	metadata := release.NewMetadata(version, commit)
	if command, err := app.NewRootCmd(os.Args[1:], metadata); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
