package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

func ui(cfg *config) (err error) {
	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"Contains": contains,
		"NotContains": func(vs []string, s string) bool {
			return !contains(vs, s)
		},
	})

	r.Static("/static", "static")
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"config": &cfg,
		})
	})
	r.GET("/staging", func(c *gin.Context) {
		c.HTML(http.StatusOK, "staging.html", gin.H{
			"config": &cfg,
		})
	})

	r.Run(fmt.Sprintf(":%v", cfg.DefaultValue.UiPort))

	return
}

func contains(vs []string, s string) bool {
	for _, v := range vs {
		if v == s {
			return true
		}
	}
	return false
}
