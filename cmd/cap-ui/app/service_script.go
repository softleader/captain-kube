package app

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
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
	Log *logrus.Logger // 這個是 server 自己的 log
	Cmd *capUiCmd
}

func (s *Script) View(c *gin.Context) {
	c.HTML(http.StatusOK, "script.html", gin.H{
		"config": &s.Cmd,
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
		} else {
			buf = &bytes.Buffer{}
			//log.SetOutput(io.MultiWriter(&sseWriter, buf))
			log.SetOutput(buf)
		}
	}

	for _, file := range files {
		filename := file.Filename
		log.Println("### Chart:", filename, "###")
		if err := doScript(log, s, &form, file); err != nil {
			log.Errorln("### [ERROR]", filename, err)
			logrus.Errorln(filename, err)
		}
		log.Println("")
		log.Println("### Finish:", filename, "###")
		log.Println("#")
		log.Println("#")

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

	tpls, err := chart.LoadArchiveBytes(log, fileHeader.Filename, buf.Bytes())
	if err != nil {
		return err
	}
	log.Debugf("%v template(s) loaded\n", len(tpls))

	if from, to := strings.TrimSpace(form.SourceRegistry), strings.TrimSpace(form.Registry); from != "" && to != "" {
		if b, err := tpls.GenerateReTagScript(from, to); err != nil {
			return err
		} else {
			log.Out.Write(b)
		}
	}

	if strutil.Contains(form.Tags, "p") {
		if b, err := tpls.GeneratePullScript(); err != nil {
			return err
		} else {
			log.Out.Write(b)
		}
	}

	if strutil.Contains(form.Tags, "l") {
		if b, err := tpls.GenerateLoadScript(); err != nil {
			return err
		} else {
			log.Out.Write(b)
		}
	}

	if strutil.Contains(form.Tags, "s") {
		if b, err := tpls.GenerateSaveScript(); err != nil {
			return err
		} else {
			log.Out.Write(b)
		}
	}

	return nil
}
