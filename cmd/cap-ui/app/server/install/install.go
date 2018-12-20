package install

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/comm"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/proto"
	"io"
	"net/http"
)

type Request struct {
	Platform       string   `form:"platform"`
	Namespace      string   `form:"namespace"`
	Tags           []string `form:"tags"`
	SourceRegistry string   `form:"sourceRegistry"`
	Registry       string   `form:"registry"`
	Verbose        string   `form:"verbose"`
}

func Serve(path string, r *gin.Engine, cfg *comm.Config) {
	r.GET(path, func(c *gin.Context) {
		c.HTML(http.StatusOK, "install.html", gin.H{
			"config": &cfg,
		})
	})
	r.POST(path, func(c *gin.Context) {
		fmt.Fprintln(c.Writer, "call: POST /install")

		var request Request
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
			fmt.Fprintln(c.Writer, "file:", file)

			buf := bytes.NewBuffer(nil)
			if readed, err := io.Copy(buf, file); err != nil {
				fmt.Fprintln(c.Writer, "reading file failed:", err)
				return
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
			if err := captain.InstallChart(c.Writer, cfg.DefaultValue.CaptainUrl, &request, 30*1000); err != nil {
				fmt.Fprintln(c.Writer, "call captain InstallChart failed:", err)
			}
		}
	})
}
