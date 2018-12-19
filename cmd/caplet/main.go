package main

import (
	"fmt"
	"github.com/softleader/captain-kube/cmd/caplet/app"
	"os"
)

func main() {
	command := app.NewCapletCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
