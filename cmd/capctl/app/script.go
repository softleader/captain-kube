package app

import (
	"errors"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
)

type scriptCmd struct {
	pull  bool
	retag bool
	save  bool
	load  bool

	sourceRegistry string
	registry       string
	charts         []string

	registryAuthUsername string // docker registry 的帳號
	registryAuthPassword string // docker registry 的密碼

	endpoint *captain.Endpoint // captain 的 endpoint ip
}

func newScriptCmd() *cobra.Command {
	c := scriptCmd{
		registryAuthUsername: env.Lookup(captain.EnvRegistryAuthUsername, captain.DefaultRegistryAuthUsername),
		registryAuthPassword: env.Lookup(captain.EnvRegistryAuthPassword, captain.DefaultRegistryAuthPassword),
	}

	cmd := &cobra.Command{
		Use:   "script [CHART...]",
		Short: "build script of helm chart",
		Long:  "build script of helm chart",
		RunE: func(cmd *cobra.Command, args []string) error {
			if c.charts = args; len(c.charts) == 0 {
				return errors.New("chart path is required")
			}
			if err := c.endpoint.Validate(); err != nil {
				return err
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&c.pull, "pull", "p", c.pull, "pull images in Chart")
	f.BoolVarP(&c.retag, "retag", "r", c.retag, "retag images in Chart")
	f.BoolVarP(&c.save, "save", "s", c.save, "save images in Chart")
	f.BoolVarP(&c.load, "load", "l", c.load, "load images in Chart")

	f.StringVarP(&c.sourceRegistry, "retag-from", "f", c.sourceRegistry, "specify the host of re-tag from, required when Sync")
	f.StringVarP(&c.registry, "retag-to", "t", c.registry, "specify the host of re-tag to, required when Sync")

	c.endpoint = captain.AddEndpointFlags(f)

	return cmd
}

func (c *scriptCmd) run() error {
	for _, chart := range c.charts {
		logrus.Println("### chart:", chart, "###")
		if err := runScript(c, chart); err != nil {
			return err
		}
	}
	return nil
}

func runScript(c *scriptCmd, path string) error {
	if expanded, err := homedir.Expand(path); err != nil {
		path = expanded
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(abs)
	if err != nil {
		return err
	}

	request := proto.GenerateScriptRequest{
		Chart: &proto.Chart{
			FileName: filepath.Base(abs),
			Content:  bytes,
			FileSize: int64(len(bytes)),
		},
		Pull: c.pull,
		Retag: &proto.ReTag{
			From: c.sourceRegistry,
			To:   c.registry,
		},
		Save: c.save,
		Load: c.load,
	}

	if err := captain.GenerateScript(logrus.StandardLogger(), c.endpoint.String(), &request, settings.timeout); err != nil {
		return err
	}

	return nil
}
