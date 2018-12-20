package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
	"io"
)

type capCmd struct {
	out     io.Writer	
	offline bool
	verbose bool
	token   string
}

func main() {
	c := capCmd{}
	cmd := &cobra.Command{
		Use:   "cap",
		Short: "the cap plugin",
		Long:  "The cap plugin",
		RunE: func(cmd *cobra.Command, args []string) error {
			c.token = os.ExpandEnv(c.token)
			return c.run()
		},
	}
	
	c.out = cmd.OutOrStdout()
	c.offline, _ = strconv.ParseBool(os.Getenv("SL_OFFLINE"))
	c.verbose, _ = strconv.ParseBool(os.Getenv("SL_VERBOSE"))

	f := cmd.Flags()
	f.BoolVarP(&c.offline, "offline", "o", c.offline, "work offline, Overrides $SL_OFFLINE")
	f.BoolVarP(&c.verbose, "verbose", "v", c.verbose, "enable verbose output, Overrides $SL_VERBOSE")
	f.StringVar(&c.token, "token", "$SL_TOKEN", "github access token. Overrides $SL_TOKEN")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func (c *capCmd) run() error {
	// use os.LookupEnv to retrieve the specific value of the environment (e.g. os.LookupEnv("SL_TOKEN"))
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "SL_") {
			fmt.Println(env)
		}
	}
	fmt.Printf("%+v\n", c)
	return nil
}
