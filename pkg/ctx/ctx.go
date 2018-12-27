package ctx

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	ContextsFile   = "contexts.yaml"
	EnvMountVolume = "SL_PLUGIN_MOUNT"
)

var (
	ErrMountVolumeNotExist = errors.New(`mount volume not found
It looks like you are running the command outside slctl (https://github.com/softleader/slctl)
For more details: https://github.com/softleader/slctl/wiki/Plugins-Guide#mount-volume
`)
	ErrNoActiveContextPresent = errors.New("no active context present") // 代表當前沒有 active 的 context
)

var PlainContexts = new(Contexts)

type Context struct {
	Endpoint     *Endpoint
	HelmTiller   *HelmTiller
	RegistryAuth *RegistryAuth
	ReTag        *ReTag
	addAllFlags  func(f *pflag.FlagSet)
	clone        func() *Context
}

type Contexts struct {
	log      *logrus.Logger
	path     string
	Contexts map[string]*Context
	Active   string // 當前
	Previous string // 上一個
}

func NewContextFromEnv() (c *Context) {
	return newContext(true)
}

func newContext(expandEnv bool) (c *Context) {
	c = &Context{
		Endpoint:     &Endpoint{},
		HelmTiller:   &HelmTiller{},
		RegistryAuth: &RegistryAuth{},
		ReTag:        &ReTag{},
	}
	if expandEnv {
		c.Endpoint.ExpandEnv()
		c.RegistryAuth.ExpandEnv()
		c.HelmTiller.ExpandEnv()
		c.ReTag.ExpandEnv()
	}
	c.addAllFlags = func(f *pflag.FlagSet) {
		c.Endpoint.AddFlags(f)
		c.RegistryAuth.AddFlags(f)
		c.HelmTiller.AddFlags(f)
		c.ReTag.AddFlags(f)
	}
	return
}

func LoadContextsFromEnv(log *logrus.Logger) (*Contexts, error) {
	mount, found := os.LookupEnv(EnvMountVolume)
	if !found {
		return nil, ErrMountVolumeNotExist
	}
	return LoadContexts(log, filepath.Join(mount, ContextsFile))
}

func LoadContexts(log *logrus.Logger, path string) (*Contexts, error) {
	log.Debugf("loading ctx from: %s\n", path)
	ctx := &Contexts{
		log:  log,
		path: path,
	}
	data, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	} else if os.IsNotExist(err) {
		ctx.Contexts = make(map[string]*Context)
		return ctx, nil
	}
	return ctx, yaml.Unmarshal(data, ctx)
}

func (ctx *Context) expandEnv() error {
	defaultCtx := newContext(true)
	return mergo.Merge(ctx, defaultCtx)
}

func (c *Contexts) GetActiveExpandEnv() (*Context, error) {
	if c.Active == "" {
		return nil, ErrNoActiveContextPresent
	}
	if ctx, found := c.Contexts[c.Active]; !found {
		return nil, fmt.Errorf("no active context exists with name %q", c.Active)
	} else {
		return ctx, ctx.expandEnv()
	}
}

func (c *Contexts) Add(name string, args []string) (err error) {
	if _, found := c.Contexts[name]; found {
		return fmt.Errorf("context %q already exists", name)
	}
	cmd := &cobra.Command{}
	f := cmd.Flags()
	ctx := newContext(false)
	ctx.addAllFlags(f)
	c.Contexts[name] = ctx
	cmd.ParseFlags(args)
	if err = c.save(); err == nil {
		c.log.Printf("Context %q added.\n", name)
	}
	return
}

func (c *Contexts) Delete(names ...string) (err error) {
	for _, name := range names {
		if name == "." {
			name = c.Active
		}
		if _, found := c.Contexts[name]; !found {
			return fmt.Errorf("no context exists with name %q", name)
		}
		delete(c.Contexts, name)
		if c.Active == name {
			c.Active = ""
		}
		if err = c.save(); err == nil {
			c.log.Printf("Context %q deleted.\n", name)
		}
	}
	return
}

func (c *Contexts) Switch(name string) (err error) {
	if name == "-" {
		return c.switchToPrevious()
	}
	if name == "x" {
		return c.switchOff()
	}
	if _, found := c.Contexts[name]; !found {
		return fmt.Errorf("no context exists with name %q", name)
	}
	c.Previous = c.Active
	c.Active = name
	if err = c.save(); err == nil {
		c.log.Printf("Switched to context %q.\n", c.Active)
	}
	return
}

func (c *Contexts) switchOff() (err error) {
	c.Previous = c.Active
	c.Active = ""
	if err = c.save(); err == nil {
		c.log.Print("Switched off the context.\n")
	}
	return
}

func (c *Contexts) switchToPrevious() (err error) {
	last := c.Previous
	c.Previous = c.Active
	c.Active = last
	if err = c.save(); err == nil {
		c.log.Printf("Switched to context %q.\n", c.Active)
	}
	return
}

func (c *Contexts) save() error {
	if c == PlainContexts {
		return errors.New("plain contexts is not able to save")
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.path, data, 0644)
}
