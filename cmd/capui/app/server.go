package app

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/capui"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/softleader/captain-kube/pkg/release"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

type capUICmd struct {
	Metadata         *release.Metadata
	port             int
	ActiveCtx        string
	defaultPlatform  string
	defaultNamespace string
	BaseUrl          string
}

type DefaultValue struct {
	Platform  string // 平台(Google/ICP)
	Namespace string
	*ctx.Context
}

func (c *capUICmd) newDefaultValue() (*DefaultValue, error) {
	ac, err := newActiveContext(logrus.StandardLogger(), c.ActiveCtx)
	if err != nil {
		return nil, err
	}
	return &DefaultValue{
		Platform:  c.defaultPlatform,
		Namespace: c.defaultNamespace,
		Context:   ac,
	}, nil
}

func NewCapUICommand(metadata *release.Metadata) (cmd *cobra.Command) {
	var verbose bool
	c := capUICmd{
		Metadata:         metadata,
		defaultPlatform:  env.Lookup(capui.EnvPlatform, capui.DefaultPlatform),
		defaultNamespace: env.Lookup(capui.EnvNamespace, capui.DefaultNamespace),
	}

	cmd = &cobra.Command{
		Use:  "capui",
		Long: "capui is a web interface for captain-kube",
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.SetOutput(colorable.NewColorableStdout()) // for windows color output
			if verbose {
				logrus.SetLevel(logrus.DebugLevel)
			}
			if !strings.HasSuffix(c.BaseUrl, "/") {
				c.BaseUrl += "/"
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	f.IntVarP(&c.port, "port", "p", 8080, "port of web ui serve port")
	f.StringVarP(&c.defaultPlatform, "platform", "k", c.defaultPlatform, "default value of k8s platform")
	f.StringVarP(&c.defaultNamespace, "namespace", "n", c.defaultNamespace, "default value of the namespace of gcp, not available now")
	f.StringVar(&c.ActiveCtx, "active-ctx", "", "active ctx")
	f.StringVar(&c.BaseUrl, "base-url", "/", "specify base url, more details: https://www.w3schools.com/tags/tag_base.asp")

	cmd.MarkFlagRequired("active-ctx")

	return
}

func (c *capUICmd) run() error {
	if err := initContext(os.Environ()); err != nil {
		return err
	}
	logrus.Printf("activated default context: %s", c.ActiveCtx)
	server := NewCapUIServer(c)
	addr := fmt.Sprintf(":%v", c.port)
	logrus.Printf("listening and serving CapUI on %v", addr)
	return server.Run(addr)
}
