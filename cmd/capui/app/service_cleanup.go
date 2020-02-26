package app

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/dur"
	pb "github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/utils"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

// CleanUpRequest 代表頁面送進來的 form request
type CleanUpRequest struct {
	Prune      bool   `form:"prune"`
	Set        string `form:"set"`
	Verbose    bool   `form:"verbose"`
	Force      bool   `form:"force"`
	Ctx        string `form:"ctx"`
	Timeout    string `form:"timeout"`
	DryRun     bool   `form:"dryRun"`
	Constraint string `form:"constraint"`
}

// CleanUp 定義了 route 的相關 call back function
type CleanUp struct {
	*capUICmd
}

// View 轉到 CleanUp 頁面
func (s *CleanUp) View(c *gin.Context) {
	c.HTML(http.StatusOK, "cleanup.html", gin.H{
		"requestURI": c.Request.RequestURI,
		"config":     &s,
	})
}

// Clean 接收 clean up 的資訊
func (s *CleanUp) Clean(c *gin.Context) {
	log := logrus.New() // 這個是這次請求要往前吐的 log
	log.SetFormatter(&utils.PlainFormatter{})
	log.SetOutput(sio.NewSSEventWriter(c))
	if v, _ := strconv.ParseBool(c.Request.FormValue("verbose")); v {
		log.SetLevel(logrus.DebugLevel)
	}

	var form CleanUpRequest
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

	if form.Prune {
		if err := captain.CallPrune(log, selectedCtx.Endpoint.String(), form.Verbose, false, dur.Parse(form.Timeout)); err != nil {
			log.Errorln(err)
			logrus.Errorln(err)
			return
		}
	}

	if len(mForm.File) > 0 {
		if form.DryRun {
			log.Warnln("running in dry-run mode, specify the '-v' flag if you want to turn on verbose output")
		}

		// ps. 在讀完request body後才可以開始response, 否則body會close
		files := mForm.File["files"]
		for _, file := range files {
			if err := s.rmc(log, selectedCtx, &form, file); err != nil {
				log.Errorln(err)
				logrus.Errorln(err)
				return
			}
		}
	}
}

func (s *CleanUp) rmc(log *logrus.Logger, activeCtx *ctx.Context, form *CleanUpRequest, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		log.Errorln(err)
		logrus.Errorln(err)
		return err
	}
	defer file.Close()
	log.Debugf("received form: %+v", form)

	buf := bytes.NewBuffer(nil)
	read, err := io.Copy(buf, file)
	if err != nil {
		return err
	}
	log.Debugf("received chart size: %v", read)

	req := &pb.RmcRequest{
		Chart: &pb.Chart{
			FileName: fileHeader.Filename,
			Content:  buf.Bytes(),
			FileSize: fileHeader.Size,
		},
		Timeout:    form.Timeout,
		DryRun:     form.DryRun,
		Force:      form.Force,
		Set:        []string{form.Set},
		Color:      false,
		Verbose:    form.Verbose,
		Constraint: form.Constraint,
	}

	log.Println("cleaning up:", fileHeader.Filename)
	return captain.CallRmc(log, activeCtx.Endpoint.String(), req, dur.Parse(form.Timeout))
}
