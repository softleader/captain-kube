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

func newContext(args ...string) (*Context, error) {
	c := &Context{
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
