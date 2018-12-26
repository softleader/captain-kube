package service

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/comm"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sse"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/softleader/captain-kube/pkg/utils/strutil"
	"io"
	"mime/multipart"
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

	mForm, err := c.MultipartForm()
	if err != nil {
		//sw.WriteStr(fmt.Sprint("loading form file error:", err))
		log.Errorln("loading form files error:", err)
		logrus.Errorln("loading form files error:", err)
		return
	}

	// ps. 在讀完request body後才可以開始response, 否則body會close
	files := mForm.File["files"]
	for _, file := range files {
		filename := file.Filename
		log.Println("### Chart:", filename, "###")
		if err := doScript(log, s, &form, file); err != nil {
			log.Errorln("### [ERROR]", filename, err)
			logrus.Errorln(filename, err)
		}
		log.Println("### Finish:", filename, "###")
		log.Println("#")
		log.Println("#")
	}

}

func doScript(log *logrus.Logger, s *Script, form *ScriptRequest, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("open file stream failed: %s", err)
	}

	log.Debugln("call: POST /script")
	log.Debugln("form:", form)
	log.Debugln("file:", file)

	buf := bytes.NewBuffer(nil)
	if readed, err := io.Copy(buf, file); err != nil {
		return fmt.Errorf("call captain GenerateScript failed: %s", err)
	} else {
		log.Debugln("readed ", readed, " bytes")
	}

	request := proto.GenerateScriptRequest{
		Chart: &proto.Chart{
			FileName: fileHeader.Filename,
			Content:  buf.Bytes(),
			FileSize: fileHeader.Size,
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
		return fmt.Errorf("call captain GenerateScript failed: %s", err)
	} else {
		log.Debugln("GenerateScript finish")
	}

	return nil
}
