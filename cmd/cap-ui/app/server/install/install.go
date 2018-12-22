package install

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/comm"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/dockerctl"
	"github.com/softleader/captain-kube/pkg/logger"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sse"
	"github.com/softleader/captain-kube/pkg/utils/strutil"
	"io"
	"net/http"
	"strconv"
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
		var form Request

		log := logger.New(sse.NewWriter(c))
		v, _ := strconv.ParseBool(c.Request.FormValue("verbose"))
		log.WithVerbose(v)

		if err := c.Bind(&form); err != nil {
			//sw.WriteStr(fmt.Sprint("binding form data error:", err))
			log.Println("binding form data error:", err)
			return
		}

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			//sw.WriteStr(fmt.Sprint("loading form file error:", err))
			log.Println("loading form file error:", err)
			return
		}

		// ps. 在讀完request body後才可以開始response, 否則body會close

		log.Println( "call: POST /install")

		log.Println( "form:", form)
		log.Println("file:", file)

		buf := bytes.NewBuffer(nil)
		if readed, err := io.Copy(buf, file); err != nil {
			log.Println( "reading file failed:", err)
			return
		} else {
			log.Println( "readed ", readed, " bytes")
		}

		// prepare rquest
		request := proto.InstallChartRequest{
			Chart: &proto.Chart{
				FileName: header.Filename,
				Content:  buf.Bytes(),
				FileSize: header.Size,
			},
			Pull: strutil.Contains(form.Tags, "p"),
			Sync: strutil.Contains(form.Tags, "r"),
			Retag: &proto.ReTag{
				From: form.SourceRegistry,
				To:   form.Registry,
			},
			RegistryAuth: &proto.RegistryAuth{
				Username: cfg.RegistryAuth.Username,
				Password: cfg.RegistryAuth.Password,
			},
			Tiller: &proto.Tiller{
				Endpoint:          cfg.Tiller.Endpoint,
				Username:          cfg.Tiller.Username,
				Password:          cfg.Tiller.Password,
				Account:           cfg.Tiller.Account,
				SkipSslValidation: cfg.Tiller.SkipSslValidation,
			},
		}

		if err := dockerctl.PullAndSync(log, &request); err != nil {
			log.Println( "Pull/Sync failed:", err)
		}

		if err := captain.InstallChart(log, cfg.DefaultValue.CaptainUrl, &request, 300); err != nil {
			fmt.Fprintln(c.Writer, "call captain InstallChart failed:", err)
		}
		fmt.Fprintln(c.Writer, "InstallChart finish")
	})
}
