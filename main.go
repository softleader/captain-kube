package main

import (
	"fmt"
	"os"
	"github.com/softleader/captain-kube/app"
)

func main() {
	command := app.NewCaptainKubeCommand()

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
