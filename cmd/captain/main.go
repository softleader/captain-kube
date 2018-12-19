package main

import (
	"fmt"
	"github.com/softleader/captain-kube/cmd/captain/app"
	"os"
)

func main() {
	command := app.NewCaptainCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
