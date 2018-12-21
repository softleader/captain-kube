package script

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/comm"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/utils"
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
		sw := utils.SSEWriter{GinContext: c}

		var form Request
		if err := c.Bind(&form); err != nil {
			fmt.Fprintln(&sw, "binding form data error:", err)
			return
		}

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			fmt.Fprintln(&sw, "loading form file error:", err)
			return
		}

		fmt.Fprintln(&sw, "call: POST /script")
		fmt.Fprintln(&sw, "form:", form)
		fmt.Fprintln(&sw, "file:", file)

		buf := bytes.NewBuffer(nil)
		if readed, err := io.Copy(buf, file); err != nil {
			fmt.Fprintln(&sw, "reading file failed:", err)
			return
		} else {
			fmt.Fprintln(&sw, "readed ", readed, " bytes")
		}

		request := proto.GenerateScriptRequest{
			Chart: &proto.Chart{
				FileName: header.Filename,
				Content:  buf.Bytes(),
				FileSize: header.Size,
			},
			Pull: strutil.Contains(form.Tags, "p"),
			Retag: &proto.ReTag{
				From: form.SourceRegistry,
				To:   form.Registry,
			},
			Save: strutil.Contains(form.Tags, "s"),
			Load: strutil.Contains(form.Tags, "l"),
		}

		if err := captain.GenerateScript(&sw, cfg.DefaultValue.CaptainUrl, &request, 300); err != nil {
			fmt.Fprintln(&sw, "call captain GenerateScript failed:", err)
		} else {
			fmt.Fprintln(&sw, "GenerateScript finish")
		}

	})
}
