package ctx

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Context struct {
	// 結構再說
	Endpoint string
}

type Contexts struct {
	log  *logrus.Logger
	path string

	Ctxs     map[string]Context
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
		return ctx, nil
	}
	return ctx, yaml.Unmarshal(data, ctx)
}

func (c *Contexts) Add(name string, args []string) error {

	return c.save()
}

func (c *Contexts) Delete(name string) error {
	if _, found := c.Ctxs[name]; !found {
		return fmt.Errorf("no context exists with name %q", name)
	}
	delete(c.Ctxs, name)
	return c.save()
}

func (c *Contexts) SwitchTo(name string) error {
	if name == "-" {
		return c.switchPrevious()
	}
	if _, found := c.Ctxs[name]; !found {
		return fmt.Errorf("no context exists with name %q", name)
	}
	c.Previous = c.Active
	c.Active = name
	return c.save()
}

func (c *Contexts) switchPrevious() error {
	last := c.Previous
	c.Previous = c.Active
	c.Active = last
	c.log.Printf("Active context is %s.\n", c.Active)
	return c.save()
}

func (c *Contexts) save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.path, data, 0644)
}
