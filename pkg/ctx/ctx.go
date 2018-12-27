package ctx

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// 結構再說
type Context struct {
	Endpoint *captain.Endpoint
}

type Contexts struct {
	log  *logrus.Logger
	path string

	Ctxs     map[string]*Context
	Active   string // 當前
	Previous string // 上一個
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
		ctx.Ctxs = make(map[string]*Context)
		return ctx, nil
	}
	return ctx, yaml.Unmarshal(data, ctx)
}

func AddContextFlags(f *pflag.FlagSet) (ctx *Context) {
	ctx = &Context{}
	ctx.Endpoint = captain.AddEndpointFlags(f)
	return
}

func (c *Contexts) Add(name string, args []string) error {
	if _, found := c.Ctxs[name]; found {
		return fmt.Errorf("context %q already exists", name)
	}
	cmd := &cobra.Command{}
	f := cmd.Flags()
	c.Ctxs[name] = AddContextFlags(f)
	cmd.ParseFlags(args)
	return c.save()
}

func (c *Contexts) Delete(name string) error {
	if _, found := c.Ctxs[name]; !found {
		return fmt.Errorf("no context exists with name %q", name)
	}
	delete(c.Ctxs, name)
	if c.Active == name {
		c.Active = ""
	}
	return c.save()
}

func (c *Contexts) Switch(name string) error {
	if name == "-" {
		return c.switchToPrevious()
	}
	if _, found := c.Ctxs[name]; !found {
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
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.path, data, 0644)
}
