package ctx

import (
	"fmt"
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
	ErrMountEnvNotExist = errors.New(fmt.Sprintf(`%q not exist
for more details: https://github.com/softleader/slctl/wiki/Plugins-Guide#environment-variables
`, EnvMountVolume))
	ErrNoActiveContextPresent = errors.New("no active context present") // 代表當前沒有 active 的 context
)

var PlainContexts = new(Contexts)

type Context struct {
	Endpoint     *Endpoint
	HelmTiller   *HelmTiller
	RegistryAuth *RegistryAuth
	ReTag        *ReTag
	addAllFlags  func(f *pflag.FlagSet)
}

type Contexts struct {
	log      *logrus.Logger
	path     string
	Contexts map[string]*Context
	Active   string // 當前
	Previous string // 上一個
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
		return nil, ErrMountEnvNotExist
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

func (ctx *Context) MergeFromEnv() (*Context, error) {
	envCtx := newContext(true)
	data, err := yaml.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	return envCtx, yaml.Unmarshal(data, envCtx)
}

func (c *Contexts) GetActive() (*Context, error) {
	if c.Active == "" {
		return newContext(true), ErrNoActiveContextPresent
	}
	if ctx, found := c.Contexts[c.Active]; !found {
		return nil, fmt.Errorf("no active context exists with name %q", c.Active)
	} else {
		return ctx, nil
	}
}

func (c *Contexts) Add(name string, args []string) error {
	if _, found := c.Contexts[name]; found {
		return fmt.Errorf("context %q already exists", name)
	}
	cmd := &cobra.Command{}
	f := cmd.Flags()
	ctx := newContext(false)
	ctx.addAllFlags(f)
	c.Contexts[name] = ctx
	cmd.ParseFlags(args)
	return c.save()
}

func (c *Contexts) Delete(name string) error {
	if _, found := c.Contexts[name]; !found {
		return fmt.Errorf("no context exists with name %q", name)
	}
	delete(c.Contexts, name)
	if c.Active == name {
		c.Active = ""
	}
	return c.save()
}

func (c *Contexts) Switch(name string) error {
	if name == "-" {
		return c.switchToPrevious()
	}
	if _, found := c.Contexts[name]; !found {
		return fmt.Errorf("no context exists with name %q", name)
	}
	c.Previous = c.Active
	c.Active = name
	c.log.Printf("Active context is %q.\n", c.Active)
	return c.save()
}

func (c *Contexts) switchToPrevious() error {
	last := c.Previous
	c.Previous = c.Active
	c.Active = last
	c.log.Printf("Active context is %q.\n", c.Active)
	return c.save()
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
