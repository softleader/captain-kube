package install

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/docker/docker/cli/command"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/gin-gonic/gin"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/comm"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/dockerctl"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/softleader/captain-kube/pkg/utils/strutil"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
		sw := utils.SSEWriter{GinContext: c}

		var form Request
		if err := c.Bind(&form); err != nil {
			//sw.WriteStr(fmt.Sprint("binding form data error:", err))
			fmt.Fprintln(&sw, "binding form data error:", err)
			return
		}

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			//sw.WriteStr(fmt.Sprint("loading form file error:", err))
			fmt.Fprintln(&sw, "loading form file error:", err)
			return
		}

		// 在讀完request body後才可以開始response, 否則body會close

		fmt.Fprintln(&sw, "call: POST /install")

		fmt.Fprintln(&sw, "form:", form)
		fmt.Fprintln(&sw, "file:", file)

		buf := bytes.NewBuffer(nil)
		if readed, err := io.Copy(buf, file); err != nil {
			fmt.Fprintln(&sw, "reading file failed:", err)
			return
		} else {
			fmt.Fprintln(&sw, "readed ", readed, " bytes")
		}

		request := proto.InstallChartRequest{
			Chart: &proto.Chart{
				FileName: header.Filename,
				Content:  buf.Bytes(),
				FileSize: header.Size,
			},
		}

		if err := pullAndSync(&sw, form, &request, cfg); err != nil {
			fmt.Fprintln(&sw, "Pull/Sync failed:", err)
		}

		if err := captain.InstallChart(&sw, cfg.DefaultValue.CaptainUrl, &request, form.Verbose, 30*1000); err != nil {
			fmt.Fprintln(&sw, "call captain InstallChart failed:", err)
		} else {
			fmt.Fprintln(&sw, "InstallChart finish")
		}
	})
}

func pullAndSync(out io.Writer, form Request, request *proto.InstallChartRequest, cfg *comm.Config) error {
	var tpls chart.Templates
	var auth proto.RegistryAuth
	if len(form.Tags) > 0 {
		// mk temp file
		tmpFile, err := ioutil.TempFile(os.TempDir(), "capui-*.tgz")
		if err != nil {
			return err
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write(request.Chart.Content); err != nil {
			return err
		}

		// load chart templates
		tpls, err = chart.LoadArchive(out, tmpFile.Name())
		if err != nil {
			return err
		}

		auth = proto.RegistryAuth{
			Username: cfg.RegistryAuth.Username,
			Password: cfg.RegistryAuth.Password,
		}
	}

	if strutil.Contains(form.Tags, "p") {
		// pull all image from chart
		for _, tpl := range tpls {
			for _, image := range tpl {
				fmt.Fprintln(out, "pulling ", image)
				result, err := dockerctl.Pull(out, *image, &auth)
				if err != nil {
					fmt.Fprintln(out, "pull image failed: ", image, ", error: ", err)
				}
				jsonmessage.DisplayJSONMessagesToStream(result, command.NewOutStream(out), nil)
			}
		}
	}

	if strutil.Contains(form.Tags, "r") {
		if len(form.SourceRegistry) > 0 && len(form.Registry) > 0 {
			// retag and push all image from chart
			for _, tpl := range tpls {
				for _, image := range tpl {
					if image.Host == form.SourceRegistry {
						fmt.Fprintln(out, "syncing ", image)
						result, err := dockerctl.ReTag(out, *image, chart.Image{
							Host: form.Registry,
							Repo: image.Repo,
							Tag:  image.Tag,
						}, &auth)
						if err != nil {
							fmt.Fprintln(out, "sync image failed: ", image, ", error: ", err)
						}
						jsonmessage.DisplayJSONMessagesToStream(result, command.NewOutStream(out), nil)
					}
				}
			}
		} else {
			return errors.New("Registry and SourceRegistry is required when retag mode")
		}
	}

	return nil
}
