package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
	args, found := contexts[activeCtx]
	if !found {
		return nil, ctx.ErrNoActiveContextPresent
	}
	return ctx.NewContext(args...)
}

func initContext(envs []string) error {
	for _, env := range envs {
		if strings.HasPrefix(env, prefix) {
			s := strings.Split(env, "=")
			key := strings.Trim(s[0], prefix)
			args := strings.Split(s[1], " ")
			contexts[strings.ToLower(key)] = args
		}
	}
	if len(contexts) == 0 {
		return errors.New("can't initial any contexts")
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
