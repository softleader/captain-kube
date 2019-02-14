package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/sse"
	"github.com/softleader/captain-kube/pkg/utils"
	"net/http"
	"strconv"
	"strings"
)

var (
	prefix               = "CAPUI_CTX"
	contexts             = make(map[string][]string)
	activeContext        *ctx.Context
	activeContextVersion []string
)

type Contexts struct {
	*capUICmd
}

func newActiveContext(log *logrus.Logger, activeCtx string) (*ctx.Context, error) {
	target := strings.ToLower(activeCtx)
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
		return nil, fmt.Errorf("failed validating context %q: %s", activeCtx, err)
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
	s.ActiveCtx = ctx

	log := logrus.New() // 這個是這次請求要往前吐的 log
	log.SetFormatter(&utils.PlainFormatter{})
	log.SetOutput(sse.NewWriter(c))
	if v, _ := strconv.ParseBool(c.Request.FormValue("verbose")); v {
		log.SetLevel(logrus.DebugLevel)
	}
	if err := activateContext(log, s.ActiveCtx); err != nil {
		log.Errorln("error activating context:", err)
		logrus.Errorln("error activating context:", err)
		return
	}

	c.Status(http.StatusOK)
}

func activateContext(log *logrus.Logger, context string) error {
	log.Printf("activating default context: %s", context)
	ac, err := newActiveContext(log, context)
	if err != nil {
		return err
	}
	activeContext = ac
	activeContextVersion = activeContextVersion[:0]
	streamLog := logrus.New()
	streamLog.SetOutput(sio.NewStreamWriter(func(p []byte) error {
		activeContextVersion = append(activeContextVersion, string(p))
		return nil
	}))
	streamLog.SetFormatter(&utils.PlainFormatter{})
	if err := captain.Version(streamLog, ac.Endpoint.String(), false, false, 5); err != nil {
		streamLog.Println(err)
	}
	return nil
}
