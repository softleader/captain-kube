package app

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/dockerd"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sse"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/softleader/captain-kube/pkg/utils/strutil"
	"io"
	"mime/multipart"
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
	Cmd *capUiCmd
}

func (s *Install) View(c *gin.Context) {
	c.HTML(http.StatusOK, "install.html", gin.H{
		"config": &s.Cmd,
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
	for _, file := range files {
		filename := file.Filename
		log.Println("### Chart:", filename, "###")
		if err := doInstall(log, s, &form, file); err != nil {
			log.Errorln("### [ERROR]", filename, err)
			logrus.Errorln(filename, err)
		}
		log.Println("### Finish:", filename, "###")
		log.Println("#")
		log.Println("#")
	}

}

func doInstall(log *logrus.Logger, s *Install, form *InstallRequest, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("open file stream failed: %s", err)
	}

	log.Debugln("call: POST /install")
	log.Debugln("form:", form)
	log.Debugln("file:", file)

	buf := bytes.NewBuffer(nil)
	if readed, err := io.Copy(buf, file); err != nil {
		return fmt.Errorf("reading file failed: %s", err)
	} else {
		log.Debugln("readed ", readed, " bytes")
	}

	// prepare rquest
	request := tw_com_softleader.InstallChartRequest{
		Chart: &tw_com_softleader.Chart{
			FileName: fileHeader.Filename,
			Content:  buf.Bytes(),
			FileSize: fileHeader.Size,
		},
		Sync: strutil.Contains(form.Tags, "s"),
		Retag: &tw_com_softleader.ReTag{
			From: form.SourceRegistry,
			To:   form.Registry,
		},
		RegistryAuth: &tw_com_softleader.RegistryAuth{
			Username: s.Cmd.RegistryAuth.Username,
			Password: s.Cmd.RegistryAuth.Password,
		},
		Tiller: &tw_com_softleader.Tiller{
			Endpoint:          s.Cmd.Tiller.Endpoint,
			Username:          s.Cmd.Tiller.Username,
			Password:          s.Cmd.Tiller.Password,
			Account:           s.Cmd.Tiller.Account,
			SkipSslValidation: s.Cmd.Tiller.SkipSslValidation,
		},
	}

	var tpls chart.Templates

	if strutil.Contains(form.Tags, "p") {
		if tpls == nil {
			if tpls, err = chart.LoadArchiveBytes(logrus.StandardLogger(), request.Chart.FileName, request.Chart.Content); err != nil {
				return err
			}
		}
		if err := dockerd.PullFromTemplates(logrus.StandardLogger(), tpls, request.RegistryAuth); err != nil {
			return err
		}
	}

	if len(request.Retag.From) > 0 && len(request.Retag.To) > 0 {
		if tpls == nil {
			if tpls, err = chart.LoadArchiveBytes(logrus.StandardLogger(), request.Chart.FileName, request.Chart.Content); err != nil {
				return err
			}
		}
		if err := dockerd.ReTagFromTemplates(logrus.StandardLogger(), tpls, request.Retag, request.RegistryAuth); err != nil {
			return err
		}
	}

	if err := captain.InstallChart(log, s.Cmd.Endpoint.String(), &request, dur.DefaultDeadlineSecond); err != nil {
		return fmt.Errorf("call captain InstallChart failed: %s", err)
	}

	log.Debugln("InstallChart finish")
	return nil
}
