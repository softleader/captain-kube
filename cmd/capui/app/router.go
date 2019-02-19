package app

import (
	"github.com/gin-gonic/gin"
	"github.com/softleader/captain-kube/pkg/utils/strutil"
	"html/template"
	"net/http"
)

func NewCapUIServer(cmd *capUICmd) (r *gin.Engine) {
	r = gin.Default()

	r.SetFuncMap(template.FuncMap{
		"Contains": strutil.Contains,
		"NotContains": func(vs []string, s string) bool {
			return !strutil.Contains(vs, s)
		},
	})

	// static and template
	r.Static("/static", "static")
	r.LoadHTMLGlob("templates/*.html")

	// index
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"requestURI": c.Request.RequestURI,
			"context":    &activeContext,
		})
	})

	installRoute := r.Group("/install")
	{
		install := &Install{
			cmd,
		}
		installRoute.GET("/", install.View)
		installRoute.POST("/", install.Chart)
	}

	contextsRoute := r.Group("/contexts")
	{
		ctxs := &Contexts{}
		contextsRoute.GET("/", ctxs.ListContext)
		contextsRoute.PUT("/:ctx", ctxs.SwitchContext)
		contextsRoute.GET("/versions", ctxs.ListContextVersions)
	}

	return
}
