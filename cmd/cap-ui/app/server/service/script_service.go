package service

import (
	"bytes"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/comm"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sse"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/softleader/captain-kube/pkg/utils/strutil"
	"io"
	"net/http"
	"strconv"
)

type ScriptRequest struct {
	Tags           []string `form:"tags"`
	SourceRegistry string   `form:"sourceRegistry"`
	Registry       string   `form:"registry"`
	Verbose        bool     `form:"verbose"`
}

type Script struct {
	Log *logrus.Logger // 這個是 server 自己的 log
	Cfg *comm.Config
}

func (s *Script) View(c *gin.Context) {
	c.HTML(http.StatusOK, "script.html", gin.H{
		"config": &s.Cfg,
	})
}

func (s *Script) Generate(c *gin.Context) {
	log := logrus.New() // 這個是這次請求要往前吐的 log
	log.SetOutput(sse.NewWriter(c))
	log.SetFormatter(&utils.PlainFormatter{})
	if v, _ := strconv.ParseBool(c.Request.FormValue("verbose")); v {
		log.SetLevel(logrus.DebugLevel)
	}

	var form ScriptRequest
	if err := c.Bind(&form); err != nil {
		log.Errorln("binding form data error:", err)
		s.Log.Errorln("binding form data error:", err)
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Errorln("loading form file error:", err)
		s.Log.Errorln("loading form file error:", err)
		return
	}

	log.Debugln("call: POST /script")
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

	if err := captain.GenerateScript(log, s.Cfg.DefaultValue.CaptainUrl, &request, 300); err != nil {
		log.Errorln("call captain GenerateScript failed:", err)
		s.Log.Errorln("call captain GenerateScript failed:", err)
	} else {
		log.Debugln("GenerateScript finish")
	}
}
