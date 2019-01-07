package ctx

import (
	"bytes"
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
)

type Context struct {
	Global       *Global
	Endpoint     *Endpoint
	HelmTiller   *HelmTiller
	RegistryAuth *RegistryAuth
	ReTag        *ReTag
}

func NewContextFromEnv() (c *Context) {
	c = &Context{
		Global:       newGlobalFromEnv(),
		Endpoint:     newEndpointFromEnv(),
		HelmTiller:   newHelmTillerFromEnv(),
		RegistryAuth: newRegistryAuthFromEnv(),
		ReTag:        newReTagFromEnv(),
	}
	return
}

func FlagsString() (string, error) {
	c := &Context{
		Global:       &Global{},
		Endpoint:     &Endpoint{},
		HelmTiller:   &HelmTiller{},
		RegistryAuth: &RegistryAuth{},
		ReTag:        &ReTag{},
	}
	cmd := &cobra.Command{}
	cmd.SetUsageTemplate(`{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}`)
	f := cmd.Flags()
	addFlags(c, f)
	bb := new(bytes.Buffer)
	cmd.SetOutput(bb)
	if err := cmd.Usage(); err != nil {
		return "", err
	}
	return bb.String(), nil
}

func newContext(args ...string) (*Context, error) {
	c := &Context{
		Global:       &Global{},
		Endpoint:     &Endpoint{},
		HelmTiller:   &HelmTiller{},
		RegistryAuth: &RegistryAuth{},
		ReTag:        &ReTag{},
	}
	if len(args) == 0 {
		return c, nil
	}
	cmd := &cobra.Command{}
	f := cmd.Flags()
	addFlags(c, f)
	return c, cmd.ParseFlags(args)
}

func (ctx *Context) expandEnv() error {
	defaultCtx := NewContextFromEnv()
	return mergo.Merge(ctx, defaultCtx)
}
