package service

import (
	"bytes"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/comm"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/dockerctl"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sse"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/softleader/captain-kube/pkg/utils/strutil"
	"io"
	"net/http"
	"strconv"
)

type InstallRequest struct {
	Platform       string   `form:"platform"`
	Namespace      string   `form:"namespace"`
	Tags           []string `form:"tags"`
	SourceRegistry string   `form:"sourceRegistry"`
	Registry       string   `form:"registry"`
	Verbose        bool     `form:"verbose"`
}

type Install struct {
	Log *logrus.Logger // 這個是 server 自己的 log
	Cfg *comm.Config
}

func (s *Install) View(c *gin.Context) {
	c.HTML(http.StatusOK, "install.html", gin.H{
		"config": &s.Cfg,
	})
}

func (s *Install) Chart(c *gin.Context) {
	log := logrus.New() // 這個是這次請求要往前吐的 log
	log.SetFormatter(&utils.PlainFormatter{})
	log.SetOutput(sse.NewWriter(c))
	if v, _ := strconv.ParseBool(c.Request.FormValue("verbose")); v {
		log.SetLevel(logrus.DebugLevel)
	}

	var form InstallRequest
	if err := c.Bind(&form); err != nil {
		//sw.WriteStr(fmt.Sprint("binding form data error:", err))
		log.Errorln("binding form data error:", err)
		s.Log.Errorln("binding form data error:", err)
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		//sw.WriteStr(fmt.Sprint("loading form file error:", err))
		log.Errorln("loading form file error:", err)
		s.Log.Errorln("loading form file error:", err)
		return
	}

	// ps. 在讀完request body後才可以開始response, 否則body會close

	log.Debugln("call: POST /install")
	log.Debugln("form:", form)
	log.Debugln("file:", file)

	buf := bytes.NewBuffer(nil)
	if readed, err := io.Copy(buf, file); err != nil {
		log.Errorln("reading file failed:", err)
		s.Log.Errorln("reading file failed:", err)
		return
	} else {
		log.Debugln("readed ", readed, " bytes")
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
			Username: s.Cfg.RegistryAuth.Username,
			Password: s.Cfg.RegistryAuth.Password,
		},
		Tiller: &proto.Tiller{
			Endpoint:          s.Cfg.Tiller.Endpoint,
			Username:          s.Cfg.Tiller.Username,
			Password:          s.Cfg.Tiller.Password,
			Account:           s.Cfg.Tiller.Account,
			SkipSslValidation: s.Cfg.Tiller.SkipSslValidation,
		},
	}

	if err := dockerctl.PullAndSync(log, &request); err != nil {
		log.Errorln("Pull/Sync failed:", err)
		s.Log.Errorln("Pull/Sync failed:", err)
	}

	if err := captain.InstallChart(log, s.Cfg.DefaultValue.CaptainUrl, &request, 300); err != nil {
		log.Errorln("call captain InstallChart failed:", err)
		s.Log.Errorln("call captain InstallChart failed:", err)
	}
	log.Debugln("InstallChart finish")
}
