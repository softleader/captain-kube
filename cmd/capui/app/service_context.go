package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/ctx"
	"net/http"
	"strings"
)

var (
	prefix   = "CTX_"
	contexts = make(map[string][]string)
)

type Contexts struct {
	*capUICmd
}

func newActiveContext(activeCtx string) (*ctx.Context, error) {
	target := strings.ToLower(activeCtx)
	args, found := contexts[target]
	if !found {
		return nil, ctx.ErrNoActiveContextPresent
	}
	logrus.Debugf("loading context '%s' with its args: %s", target, strings.Join(args, " "))
	c, err := ctx.NewContext(args...)
	if err != nil {
		return nil, err
	}
	err = c.ExpandEnv()
	return c, err
}

func initContext(envs []string) error {
	for _, env := range envs {
		if strings.HasPrefix(env, prefix) {
			s := strings.Split(env, "=")
			key := strings.Replace(s[0], prefix, "", -1)
			args := strings.Split(s[1], " ")
			// to make sure args are alright
			if _, err := ctx.NewContext(args...); err != nil {
				return err
			}
			contexts[strings.ToLower(key)] = args
		}
	}
	if len(contexts) == 0 {
		return errors.New("can't found any contexts")
	}
	logrus.Printf("context loaded:")
	for k, v := range contexts {
		logrus.Printf("%s: %s", k, v)
	}
	return nil
}

func (s *Contexts) ListContext(c *gin.Context) {
	var names []string
	for k := range contexts {
		names = append(names, k)
	}
	c.JSON(http.StatusOK, names)
}

func (s *Contexts) SwitchContext(c *gin.Context) {
	ctx := c.Param("ctx")
	if ctx == "" {
		c.Error(fmt.Errorf("can't switch to blank context: %q", ctx))
		return
	}
	_, found := contexts[strings.ToLower(ctx)]
	if !found {
		c.Error(fmt.Errorf("context %q not found", ctx))
		return
	}
	s.ActiveCtx = ctx
	c.Status(http.StatusOK)
}
