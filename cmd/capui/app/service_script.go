package app

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sse"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/softleader/captain-kube/pkg/utils/strutil"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

type ScriptRequest struct {
	Tags           []string `form:"tags"`
	SourceRegistry string   `form:"sourceRegistry"`
	Registry       string   `form:"registry"`
	Verbose        bool     `form:"verbose"`
}

type Script struct {
	*capUICmd
}

func (s *Script) View(c *gin.Context) {
	dft, err := s.newDefaultValue()
	if err != nil {
		c.Error(err)
		return
	}
	c.HTML(http.StatusOK, "script.html", gin.H{
		"config":       &s,
		"defaultValue": dft,
	})
}

func (s *Script) Generate(c *gin.Context) {
	sseWriter := *sse.NewWriter(c)
	log := logrus.New() // 這個是這次請求要往前吐的 log
	log.SetOutput(&sseWriter)
	log.SetFormatter(&utils.PlainFormatter{})
	if v, _ := strconv.ParseBool(c.Request.FormValue("verbose")); v {
		log.SetLevel(logrus.DebugLevel)
	}

	var form ScriptRequest
	if err := c.Bind(&form); err != nil {
		log.Errorln("binding form data error:", err)
		logrus.Errorln("binding form data error:", err)
		return
	}

	mForm, err := c.MultipartForm()
	if err != nil {
		//sw.WriteStr(fmt.Sprint("loading form file error:", err))
		log.Errorln("loading form files error:", err)
		logrus.Errorln("loading form files error:", err)
		return
	}

	activeCtx, err := newActiveContext(log, s.ActiveCtx)
	if err != nil {
		log.Errorln(err)
		logrus.Errorln(err)
		return
	}

	// do some validation check
	if err := activeCtx.Endpoint.Validate(); err != nil {
		log.Errorln(err)
		logrus.Errorln(err)
		return
	}

	// ps. 在讀完request body後才可以開始response, 否則body會close
	files := mForm.File["files"]

	// diff validate
	var buf *bytes.Buffer
	var scripts []string
	if strutil.Contains(form.Tags, "d") {
		if len(files) != 2 {
			log.Errorln("diff mode must have two files")
			logrus.Errorln("diff mode must have two files")
			return
		}
		buf = &bytes.Buffer{}
		//log.SetOutput(io.MultiWriter(&sseWriter, buf))
		log.SetOutput(buf)
	}

	for _, file := range files {
		filename := file.Filename
		log.Println("### Chart:", filename, "###")
		if err := s.script(log, activeCtx, &form, file); err != nil {
			log.Errorln(err)
			logrus.Errorln(err)
		}
		log.Println("")

		// 如果buf裡面有存東西，則append到暫存裡面
		if buf != nil {
			scripts = append(scripts, buf.String())
			buf.Reset()
		}
	}

	if len(scripts) == 2 {
		log.SetOutput(&sseWriter)
		log.Println("### Diffs: ###")
		lines := strutil.DiffNewLines(scripts[0], scripts[1])
		log.Println(strings.Join(lines, "\n"))
	}

}

func (s *Script) script(log *logrus.Logger, activeCtx *ctx.Context, form *ScriptRequest, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open file stream: %s", err)
	}

	log.Debugf("received form: %+v", form)

	buf := bytes.NewBuffer(nil)
	read, err := io.Copy(buf, file)
	if err != nil {
		return fmt.Errorf("failed to copy buffer to file: %s", err)
	}
	log.Debugf("received chart size: %v", read)

	request := captainkube_v2.GenerateScriptRequest{
		Verbose: form.Verbose,
		Chart: &captainkube_v2.Chart{
			FileName: fileHeader.Filename,
			Content:  buf.Bytes(),
			FileSize: fileHeader.Size,
		},
		Pull: strutil.Contains(form.Tags, "p"),
		Retag: &captainkube_v2.ReTag{
			From: form.SourceRegistry,
			To:   form.Registry,
		},
		Save: strutil.Contains(form.Tags, "s"),
		Load: strutil.Contains(form.Tags, "l"),
	}

	if err := captain.GenerateScript(log, activeCtx.Endpoint.String(), &request, 300); err != nil {
		return fmt.Errorf("failed to call backend: %s", err)
	}
	return nil
}
