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
		dft, err := cmd.newDefaultValue()
		if err != nil {
			c.Error(err)
			return
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"config":       cmd,
			"defaultValue": dft,
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

	scriptRoute := r.Group("/script")
	{
		script := &Script{
			cmd,
		}
		scriptRoute.GET("/", script.View)
		scriptRoute.POST("/", script.Generate)
	}

	contextsRoute := r.Group("/contexts")
	{
		ctxs := &Contexts{
			cmd,
		}
		contextsRoute.GET("/", ctxs.ListContext)
		contextsRoute.PUT("/:ctx", ctxs.SwitchContext)
	}

	return
}
