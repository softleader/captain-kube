package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/utils"
	"net/http"
	"strconv"
	"strings"
)

var (
	prefix        = "CAPUI_CTX"
	contexts      = make(map[string][]string)
	activeContext *ctx.Context // for 頁面的 default value 呈現使用
)

type Contexts struct {
}

func newContext(log *logrus.Logger, context string) (*ctx.Context, error) {
	target := strings.ToLower(context)
	args, found := contexts[target]
	if !found {
		return nil, ctx.ErrNoActiveContextPresent
	}
	log.Debugf("loading context '%s' with its args: %s", target, strings.Join(args, " "))
	c, err := ctx.NewContext(args...)
	if err != nil {
		return nil, err
	}
	err = c.ExpandEnv()
	// do some validation check
	if err := c.Endpoint.Validate(); err != nil {
		return nil, fmt.Errorf("failed validating context %q: %s", context, err)
	}
	// apply some default value
	if te := strings.TrimSpace(c.HelmTiller.Endpoint); len(te) == 0 {
		c.HelmTiller.Endpoint = c.Endpoint.Host
	}
	return c, err
}

func initContext(envs []string) error {
	for _, env := range envs {
		if strings.HasPrefix(env, prefix) {
			s := strings.Split(env, "=")
			key := s[0][len(prefix)+1:]
			args := strings.Split(s[1], " ")
			// to make sure args are alright
			_, err := ctx.NewContext(args...)
			if err != nil {
				return fmt.Errorf("failed loading context %q: %s", key, err)
			}
			contexts[strings.ToLower(key)] = args
		}
	}
	if len(contexts) == 0 {
		return errors.New("can't found any contexts")
	}
	logrus.Printf("loading context:")
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
	if err := activateContext(logrus.StandardLogger(), ctx); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *Contexts) ListContextVersions(c *gin.Context) {
	full := false
	color := false
	timeout := int64(5)
	if q := c.Query("full"); len(q) != 0 {
		full, _ = strconv.ParseBool(q)
	}
	if q := c.Query("color"); len(q) != 0 {
		color, _ = strconv.ParseBool(q)
	}
	if q := c.Query("timeout"); len(q) != 0 {
		if i, err := strconv.Atoi(q); err != nil {
			timeout = int64(i)
		}
	}
	contextsVersions := make(map[string][]string)
	for context := range contexts {
		var versions []string
		log := logrus.New()
		log.SetOutput(sio.NewStreamWriter(func(p []byte) error {
			versions = append(versions, strings.TrimSuffix(string(p), fmt.Sprintln()))
			return nil
		}))
		log.SetFormatter(&utils.PlainFormatter{})
		if c, err := newContext(log, context); err != nil {
			log.Println(err)
		} else if err := captain.Version(log, c.Endpoint.String(), full, color, timeout); err != nil {
			log.Println(err)
		}
		contextsVersions[context] = versions
	}
	c.JSON(http.StatusOK, contextsVersions)
}

func activateContext(log *logrus.Logger, context string) error {
	log.Printf("activating context %s", context)
	ac, err := newContext(log, context)
	if err != nil {
		return err
	}
	activeContext = ac
	return nil
}
