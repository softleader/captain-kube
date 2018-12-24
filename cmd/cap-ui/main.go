package main

import (
	"fmt"
	"github.com/softleader/captain-kube/cmd/cap-ui/app"
	"os"
)

func main() {
	command := app.NewCapUiCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
