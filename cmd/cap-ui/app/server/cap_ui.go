package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/comm"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/service"
	"github.com/softleader/captain-kube/pkg/utils/strutil"
	"html/template"
	"net/http"
)

func NewCapUiServer(log *logrus.Logger, cfg *comm.Config) (r *gin.Engine) {
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
			"config": cfg,
		})
	})

	installRoute := r.Group("/install")
	{
		install := &service.Install{
			Log: log,
			Cfg: cfg,
		}
		installRoute.GET("/", install.View)
		installRoute.POST("/", install.Chart)
	}

	scriptRoute := r.Group("/script")
	{
		script := &service.Script{
			Log: log,
			Cfg: cfg,
		}
		scriptRoute.GET("/", script.View)
		scriptRoute.POST("/", script.Generate)
	}

	return
}
