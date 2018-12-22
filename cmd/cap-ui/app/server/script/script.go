package script

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/comm"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/logger"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sse"
	"github.com/softleader/captain-kube/pkg/utils/strutil"
	"io"
	"net/http"
	"strconv"
)

type Request struct {
	Tags           []string `form:"tags"`
	SourceRegistry string   `form:"sourceRegistry"`
	Registry       string   `form:"registry"`
	Verbose        bool     `form:"verbose"`
}

func Serve(path string, r *gin.Engine, cfg *comm.Config) {
	r.GET(path, func(c *gin.Context) {
		c.HTML(http.StatusOK, "script.html", gin.H{
			"config": &cfg,
		})
	})
	r.POST(path, func(c *gin.Context) {
		log := logger.New(sse.NewWriter(c))
		v, _ := strconv.ParseBool(c.Request.FormValue("verbose"))
		log.WithVerbose(v)

		var form Request
		if err := c.Bind(&form); err != nil {
			log.Println("binding form data error:", err)
			return
		}

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			log.Println("loading form file error:", err)
			return
		}

		log.Println("call: POST /script")
		log.Println("form:", form)
		log.Println("file:", file)

		buf := bytes.NewBuffer(nil)
		if readed, err := io.Copy(buf, file); err != nil {
			log.Println("reading file failed:", err)
			return
		} else {
			log.Println("readed ", readed, " bytes")
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

		if err := captain.GenerateScript(log, cfg.DefaultValue.CaptainUrl, &request, 300); err != nil {
			log.Println("call captain GenerateScript failed:", err)
		} else {
			log.Println("GenerateScript finish")
		}

	})
}
