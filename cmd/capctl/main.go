package main

import (
	"github.com/softleader/captain-kube/cmd/capctl/app"
	"os"
)

func main() {
	command := app.NewRootCmd(os.Args[1:])
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
