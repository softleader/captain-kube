package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/utils"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

var (
	prefix   = "CAPUI_CTX"
	contexts = make(map[string][]string)
)

// Contexts 定義了 route 的相關 call back function
type Contexts struct {
	*capUICmd
}

// ActiveContext 包含了當前啟用的 context 資訊, for 頁面使用
type ActiveContext struct {
	*ctx.Context
	Name string
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
			args := strings.Split(strings.Join(s[1:], "="), " ")
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

// ContextInfo 定義了每個 context 要呈現在頁面的資訊
type ContextInfo struct {
	Active bool
	Name   string
	Args   string
}

// ListContext 回傳所有 contexts
func (s *Contexts) ListContext(c *gin.Context) {
	if ctx := c.Query("ctx"); len(ctx) > 0 {
		if context, err := switchContext(ctx); err != nil {
			c.Error(err)
		} else {
			s.Context = context
		}
	}
	var names []string
	for k := range contexts {
		names = append(names, k)
	}
	sort.Strings(names)
	if _, exists := c.GetQuery("json"); exists {
		c.JSON(http.StatusOK, names)
		return
	}
	var info []ContextInfo
	for _, name := range names {
		info = append(info, ContextInfo{
			Active: s.Context.Name == name,
			Name:   name,
			Args:   strings.Join(contexts[name], " "),
		})
	}
	c.HTML(http.StatusOK, "contexts.html", gin.H{
		"requestURI": c.Request.RequestURI,
		"config":     &s,
		"info":       info,
	})
}

func switchContext(ctx string) (*ActiveContext, error) {
	if ctx == "" {
		return nil, fmt.Errorf("can't switch to blank context: %q", ctx)
	}
	_, found := contexts[strings.ToLower(ctx)]
	if !found {
		return nil, fmt.Errorf("context %q not found", ctx)
	}
	return activateContext(logrus.StandardLogger(), ctx)
}

// ListContextVersions 回傳所有 context 及其版本
func (s *Contexts) ListContextVersions(c *gin.Context) {
	full := false
	color := false
	timeout := dur.DefaultDeadline
	if q := c.Query("full"); len(q) != 0 {
		full, _ = strconv.ParseBool(q)
	}
	if q := c.Query("color"); len(q) != 0 {
		color, _ = strconv.ParseBool(q)
	}
	if q := c.Query("timeout"); len(q) != 0 {
		timeout = q
	}
	if ctx := c.Query("ctx"); len(ctx) != 0 {
		c.JSON(http.StatusOK, contextVersions(ctx, full, color, timeout))
	} else {
		c.JSON(http.StatusOK, contextsVersions(full, color, timeout))
	}
}

func contextsVersions(full, color bool, timeout string) map[string][]string {
	contextsVersions := make(map[string][]string)
	for context := range contexts {
		contextsVersions[context] = contextVersions(context, full, color, timeout)
	}
	return contextsVersions
}

func contextVersions(context string, full, color bool, timeout string) []string {
	var versions []string
	log := logrus.New()
	log.SetOutput(sio.NewStreamWriter(func(p []byte) error {
		versions = append(versions, strings.TrimSuffix(string(p), fmt.Sprintln()))
		return nil
	}))
	log.SetFormatter(&utils.PlainFormatter{})
	if c, err := newContext(log, context); err != nil {
		log.Println(err)
	} else if err := captain.CallVersion(log, c.Endpoint.String(), full, color, dur.Parse(timeout)); err != nil {
		log.Println(err)
	}
	return versions
}

func activateContext(log *logrus.Logger, context string) (*ActiveContext, error) {
	log.Printf("activating context %s", context)
	ac, err := newContext(log, context)
	if err != nil {
		return nil, err
	}
	return &ActiveContext{
		Context: ac,
		Name:    context,
	}, nil
}
