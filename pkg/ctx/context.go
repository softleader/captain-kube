package ctx

import (
	"github.com/imdario/mergo"
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
