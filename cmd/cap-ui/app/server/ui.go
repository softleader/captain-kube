package server

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/proto"
	"html/template"
	"io"
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

	r.GET("/install", func(c *gin.Context) {
		c.HTML(http.StatusOK, "install.html", gin.H{
			"config": &cfg,
		})
	})
	r.POST("/install", func(c *gin.Context) {
		fmt.Fprintln(c.Writer, "call: POST /install")

		var request InstallRequest
		if err := c.Bind(&request); err != nil {
			fmt.Fprintln(c.Writer, "binding form data error:", err)
			return
		} else {
			fmt.Fprintln(c.Writer, "form:", request)
		}

		if file, header, err := c.Request.FormFile("file"); err != nil {
			fmt.Fprintln(c.Writer, "loading form file error:", err)
			return
		} else {
			fmt.Fprintln(c.Writer, "file.header:", header)
			fmt.Fprintln(c.Writer, "file:", file)

			buf := bytes.NewBuffer(nil)
			if readed, err := io.Copy(buf, file); err != nil {
				fmt.Fprintln(c.Writer, "reading file failed:", err)
			} else {
				fmt.Fprintln(c.Writer, "readed ", readed, " bytes")
			}

			request := proto.InstallChartRequest{
				Chart: &proto.Chart{
					FileName: header.Filename,
					Content:  buf.Bytes(),
					FileSize: header.Size,
				},
			}
			if err := captain.InstallChart(c.Writer, cfg.DefaultValue.CaptainUrl, &request, 300); err != nil {
				fmt.Fprintln(c.Writer, "call captain InstallChart failed:", err)
			}
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
			return
		} else {
			fmt.Fprintln(c.Writer, "form:", request)
		}

		if file, header, err := c.Request.FormFile("file"); err != nil {
			fmt.Fprintln(c.Writer, "loading form file error:", err)
			return
		} else {
			fmt.Fprintln(c.Writer, "file.header:", header)
			fmt.Fprintln(c.Writer, "file:", file)

			buf := bytes.NewBuffer(nil)
			if readed, err := io.Copy(buf, file); err != nil {
				fmt.Fprintln(c.Writer, "reading file failed:", err)
			} else {
				fmt.Fprintln(c.Writer, "readed ", readed, " bytes")
			}

			request := proto.GenerateScriptRequest{
				Chart: &proto.Chart{
					FileName: header.Filename,
					Content:  buf.Bytes(),
					FileSize: header.Size,
				},
			}

			if err := captain.GenerateScript(c.Writer, cfg.DefaultValue.CaptainUrl, &request, 300); err != nil {
				fmt.Fprintln(c.Writer, "call captain InstallChart failed:", err)
			}
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
