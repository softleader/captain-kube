package app

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/dockerd"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	pb "github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/softleader/captain-kube/pkg/utils/strutil"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

// InstallRequest 代表頁面送進來的 form request
type InstallRequest struct {
	Platform       string   `form:"platform"`
	Namespace      string   `form:"namespace"`
	Tags           []string `form:"tags"`
	SourceRegistry string   `form:"sourceRegistry"`
	Registry       string   `form:"registry"`
	Verbose        bool     `form:"verbose"`
	Ctx            string   `form:"ctx"`
	Timeout        string   `form:"timeout"`
}

// Charts 定義了 route 的相關 call back function
type Charts struct {
	*capUICmd
}

// View 轉到 charts 頁面
func (cs *Charts) View(c *gin.Context) {
	c.HTML(http.StatusOK, "charts.html", gin.H{
		"requestURI": c.Request.RequestURI,
		"config":     &cs,
		"context":    &activeContext,
	})
}

// Install 接收上傳的 chart
func (cs *Charts) Install(c *gin.Context) {
	log := logrus.New() // 這個是這次請求要往前吐的 log
	log.SetFormatter(&utils.PlainFormatter{})
	log.SetOutput(sio.NewSSEventWriter(c))
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
		if err := cs.install(log, selectedCtx, &form, file); err != nil {
			log.Errorln(err)
			logrus.Errorln(err)
		} else {
			log.Println("Successfully handled chart:", filename)
		}
		log.Println("")
	}

}

func (cs *Charts) install(log *logrus.Logger, activeCtx *ctx.Context, form *InstallRequest, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open file stream: %s", err)
	}
	defer file.Close()
	log.Debugf("received form: %+v", form)

	buf := bytes.NewBuffer(nil)
	read, err := io.Copy(buf, file)
	if err != nil {
		return fmt.Errorf("failed to copy buffer to file: %s", err)
	}
	log.Debugf("received chart size: %v", read)

	// prepare request
	c := &pb.Chart{
		FileName: fileHeader.Filename,
		Content:  buf.Bytes(),
		FileSize: fileHeader.Size,
	}
	rt := &pb.ReTag{
		From: form.SourceRegistry,
		To:   form.Registry,
	}
	ra := &pb.RegistryAuth{
		Username: activeCtx.RegistryAuth.Username,
		Password: activeCtx.RegistryAuth.Password,
	}

	var tpls chart.Templates

	if strutil.Contains(form.Tags, "p") {
		if tpls == nil {
			if tpls, err = chart.LoadArchiveBytes(logrus.StandardLogger(), c.FileName, c.Content); err != nil {
				return err
			}
		}
		if err := dockerd.PullFromTemplates(logrus.StandardLogger(), tpls, ra); err != nil {
			return err
		}
	}

	if strutil.Contains(form.Tags, "r") {
		if len(rt.From) > 0 && len(rt.To) > 0 {
			if tpls == nil {
				if tpls, err = chart.LoadArchiveBytes(logrus.StandardLogger(), c.FileName, c.Content); err != nil {
					return err
				}
			}
			if err := dockerd.ReTagFromTemplates(logrus.StandardLogger(), tpls, rt, ra); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("requires both 'retag-from' and 'retag-to' fields to tag and push images")
		}
	}

	if strutil.Contains(form.Tags, "i") {
		req := &pb.InstallChartRequest{
			Verbose:      form.Verbose,
			Chart:        c,
			Timeout:      form.Timeout,
			Sync:         strutil.Contains(form.Tags, "s"),
			Retag:        rt,
			RegistryAuth: ra,
			Tiller: &pb.Tiller{
				Endpoint:          activeCtx.HelmTiller.Endpoint,
				Username:          activeCtx.HelmTiller.Username,
				Password:          activeCtx.HelmTiller.Password,
				Account:           activeCtx.HelmTiller.Account,
				SkipSslValidation: activeCtx.HelmTiller.SkipSslValidation,
			},
		}
		if err := captain.CallInstallChart(log, activeCtx.Endpoint.String(), req, dur.Parse(form.Timeout)); err != nil {
			return fmt.Errorf("failed to call backend: %s", err)
		}
	} else if strutil.Contains(form.Tags, "s") {
		req := &pb.SyncChartRequest{
			Verbose:      form.Verbose,
			Chart:        c,
			Timeout:      form.Timeout,
			Retag:        rt,
			RegistryAuth: ra,
		}
		if err := captain.CallSyncChart(log, activeCtx.Endpoint.String(), req, dur.Parse(form.Timeout)); err != nil {
			return fmt.Errorf("failed to call backend: %s", err)
		}
	}
	return nil
}
