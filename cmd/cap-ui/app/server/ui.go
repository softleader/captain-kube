package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/comm"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/install"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/script"
	"github.com/softleader/captain-kube/pkg/utils/strutil"
	"html/template"
	"net/http"
)

func Ui(cfg *comm.Config, port int) (err error) {
	r := gin.Default()

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
			"config": &cfg,
		})
	})

	// services
	install.Serve("install", r, cfg)
	script.Serve("script", r, cfg)

	r.Run(fmt.Sprintf(":%v", port))

	return
}
