package app

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/dockerd"
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
	Ctx            string   `form:"ctx"`
	Timeout        int64    `form:"timeout"`
}

type Install struct {
	*capUICmd
}

func (s *Install) View(c *gin.Context) {
	c.HTML(http.StatusOK, "install.html", gin.H{
		"config":  &s,
		"context": &activeContext,
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
		log.Errorln("binding form data error:", err)
		logrus.Errorln("binding form data error:", err)
		return
	}

	mForm, err := c.MultipartForm()
	if err != nil {
		log.Errorln("loading form files error:", err)
		logrus.Errorln("loading form files error:", err)
		return
	}

	// 而頁面到送進來的過程中, global activeCtx 有可能改變
	// 所以這邊必須以頁面的 ctx 為主
	selectedCtx, err := newContext(log, form.Ctx)
	if err != nil {
		log.Errorln(err)
		logrus.Errorln(err)
		return
	}

	// ps. 在讀完request body後才可以開始response, 否則body會close
	files := mForm.File["files"]
	for _, file := range files {
		filename := file.Filename
		log.Println("Installing chart:", filename)
		if err := s.install(log, selectedCtx, &form, file); err != nil {
			log.Errorln(err)
			logrus.Errorln(err)
		} else {
			log.Println("Successfully installed chart:", filename)
		}
		log.Println("")
	}

}

func (s *Install) install(log *logrus.Logger, activeCtx *ctx.Context, form *InstallRequest, fileHeader *multipart.FileHeader) error {
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

	// prepare rquest
	request := captainkube_v2.InstallChartRequest{
		Verbose: form.Verbose,
		Chart: &captainkube_v2.Chart{
			FileName: fileHeader.Filename,
			Content:  buf.Bytes(),
			FileSize: fileHeader.Size,
		},
		Sync: strutil.Contains(form.Tags, "s"),
		Retag: &captainkube_v2.ReTag{
			From: form.SourceRegistry,
			To:   form.Registry,
		},
		RegistryAuth: &captainkube_v2.RegistryAuth{
			Username: activeCtx.RegistryAuth.Username,
			Password: activeCtx.RegistryAuth.Password,
		},
		Tiller: &captainkube_v2.Tiller{
			Endpoint:          activeCtx.HelmTiller.Endpoint,
			Username:          activeCtx.HelmTiller.Username,
			Password:          activeCtx.HelmTiller.Password,
			Account:           activeCtx.HelmTiller.Account,
			SkipSslValidation: activeCtx.HelmTiller.SkipSslValidation,
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

	if err := captain.InstallChart(log, activeCtx.Endpoint.String(), &request, form.Timeout); err != nil {
		return fmt.Errorf("failed to call backend: %s", err)
	}
	return nil
}
