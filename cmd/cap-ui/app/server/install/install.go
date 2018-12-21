package install

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/comm"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/dockerctl"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/softleader/captain-kube/pkg/utils/strutil"
	"io"
	"net/http"
)

type Request struct {
	Platform       string   `form:"platform"`
	Namespace      string   `form:"namespace"`
	Tags           []string `form:"tags"`
	SourceRegistry string   `form:"sourceRegistry"`
	Registry       string   `form:"registry"`
	Verbose        bool     `form:"verbose"`
}

func Serve(path string, r *gin.Engine, cfg *comm.Config) {
	r.GET(path, func(c *gin.Context) {
		c.HTML(http.StatusOK, "install.html", gin.H{
			"config": &cfg,
		})
	})
	r.POST(path, func(c *gin.Context) {
		sw := utils.SSEWriter{GinContext: c}

		var form Request
		if err := c.Bind(&form); err != nil {
			//sw.WriteStr(fmt.Sprint("binding form data error:", err))
			fmt.Fprintln(&sw, "binding form data error:", err)
			return
		}

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			//sw.WriteStr(fmt.Sprint("loading form file error:", err))
			fmt.Fprintln(&sw, "loading form file error:", err)
			return
		}

		// ps. 在讀完request body後才可以開始response, 否則body會close

		fmt.Fprintln(&sw, "call: POST /install")

		fmt.Fprintln(&sw, "form:", form)
		fmt.Fprintln(&sw, "file:", file)

		buf := bytes.NewBuffer(nil)
		if readed, err := io.Copy(buf, file); err != nil {
			fmt.Fprintln(&sw, "reading file failed:", err)
			return
		} else {
			fmt.Fprintln(&sw, "readed ", readed, " bytes")
		}

		// prepare rquest
		request := proto.InstallChartRequest{
			Chart: &proto.Chart{
				FileName: header.Filename,
				Content:  buf.Bytes(),
				FileSize: header.Size,
			},
			Pull:           strutil.Contains(form.Tags, "p"),
			Sync:           strutil.Contains(form.Tags, "r"),
			SourceRegistry: form.SourceRegistry,
			Registry:       form.Registry,
			RegistryAuth: &proto.RegistryAuth{
				Username: cfg.RegistryAuth.Username,
				Password: cfg.RegistryAuth.Password,
			},
		}

		if err := dockerctl.PullAndSync(&sw, &request); err != nil {
			fmt.Fprintln(&sw, "Pull/Sync failed:", err)
		}

		if err := captain.InstallChart(c.Writer, cfg.DefaultValue.CaptainUrl, &request, 30*1000); err != nil {
			fmt.Fprintln(c.Writer, "call captain InstallChart failed:", err)
		}
		fmt.Fprintln(c.Writer, "InstallChart finish")
	})
}
