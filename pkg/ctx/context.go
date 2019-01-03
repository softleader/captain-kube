package ctx

import (
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
)

type Context struct {
	Endpoint     *Endpoint
	HelmTiller   *HelmTiller
	RegistryAuth *RegistryAuth
	ReTag        *ReTag
}

func NewContextFromEnv() (c *Context) {
	c = &Context{
		Endpoint:     newEndpointFromEnv(),
		HelmTiller:   newHelmTillerFromEnv(),
		RegistryAuth: newRegistryAuthFromEnv(),
		ReTag:        newReTagFromEnv(),
	}
	return
}

func newContextFromArgs(args []string) (*Context, error) {
	cmd := &cobra.Command{}
	f := cmd.Flags()
	ctx := newContext()
	addFlags(ctx, f)
	return ctx, cmd.ParseFlags(args)
}

func newContext() (c *Context) {
	c = &Context{
		Endpoint:     &Endpoint{},
		HelmTiller:   &HelmTiller{},
		RegistryAuth: &RegistryAuth{},
		ReTag:        &ReTag{},
	}
	return
}
func (ctx *Context) expandEnv() error {
	defaultCtx := NewContextFromEnv()
	return mergo.Merge(ctx, defaultCtx)
}
