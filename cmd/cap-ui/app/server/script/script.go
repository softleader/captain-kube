package script

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/comm"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/utils/strutil"
	"io"
	"net/http"
)

type Request struct {
	Tags           []string `form:"tags"`
	SourceRegistry string   `form:"sourceRegistry"`
	Registry       string   `form:"registry"`
}

func Serve(path string, r *gin.Engine, cfg *comm.Config) {
	r.GET(path, func(c *gin.Context) {
		c.HTML(http.StatusOK, "script.html", gin.H{
			"config": &cfg,
		})
	})
	r.POST(path, func(c *gin.Context) {
		fmt.Fprintln(c.Writer, "call: POST /script")

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

			request := proto.GenerateScriptRequest{
				Chart: &proto.Chart{
					FileName: header.Filename,
					Content:  buf.Bytes(),
					FileSize: header.Size,
				},
				Pull: strutil.Contains(request.Tags, "p"),
				Retag: &proto.ReTag{
					From: request.SourceRegistry,
					To:   request.Registry,
				},
				Save: strutil.Contains(request.Tags, "s"),
				Load: strutil.Contains(request.Tags, "l"),
			}

			if err := captain.GenerateScript(c.Writer, cfg.DefaultValue.CaptainUrl, &request, 300); err != nil {
				fmt.Fprintln(c.Writer, "call captain GenerateScript failed:", err)
			}
			fmt.Fprintln(c.Writer, "GenerateScript finish")
		}

	})
}
