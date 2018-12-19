package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

func Ui(cfg *config) (err error) {
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
	r.POST("/staging", func(c *gin.Context) {
		fmt.Fprintln(c.Writer, "call: POST /staging")

		var request StagingRequest
		if err := c.Bind(&request); err != nil {
			fmt.Fprintln(c.Writer, "binding form data error:", err)
		} else {
			fmt.Fprintln(c.Writer, "form:", request)
		}

		if file, header, err := c.Request.FormFile("file"); err != nil {
			fmt.Fprintln(c.Writer, "loading form file error:", err)
		} else {
			fmt.Fprintln(c.Writer, "file.header:", header)
			fmt.Fprintln(c.Writer, "file:", file)
		}
	})

	r.GET("/script", func(c *gin.Context) {
		c.HTML(http.StatusOK, "script.html", gin.H{
			"config": &cfg,
		})
	})
	r.POST("/script", func(c *gin.Context) {
		fmt.Fprintln(c.Writer, "call: POST /script")

		var request ScriptRequest
		if err := c.Bind(&request); err != nil {
			fmt.Fprintln(c.Writer, "binding form data error:", err)
		} else {
			fmt.Fprintln(c.Writer, "form:", request)
		}

		if file, header, err := c.Request.FormFile("file"); err != nil {
			fmt.Fprintln(c.Writer, "loading form file error:", err)
		} else {
			fmt.Fprintln(c.Writer, "file.header:", header)
			fmt.Fprintln(c.Writer, "file:", file)
		}

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
