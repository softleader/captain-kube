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
		fmt.Fprintln(c.Writer, "call: POST /install")

		var form Request
		if err := c.Bind(&form); err != nil {
			fmt.Fprintln(c.Writer, "binding form data error:", err)
			return
		} else {
			fmt.Fprintln(c.Writer, "form:", form)
		}

		if file, header, err := c.Request.FormFile("file"); err != nil {
			fmt.Fprintln(c.Writer, "loading form file error:", err)
			return
		} else {
			fmt.Fprintln(c.Writer, "file:", file)

			buf := bytes.NewBuffer(nil)
			if readed, err := io.Copy(buf, file); err != nil {
				fmt.Fprintln(c.Writer, "reading file failed:", err)
				return
			} else {
				fmt.Fprintln(c.Writer, "readed ", readed, " bytes")
			}

			request := proto.InstallChartRequest{
				Chart: &proto.Chart{
					FileName: header.Filename,
					Content:  buf.Bytes(),
					FileSize: header.Size,
				},
			}

			if err := pullAndSync(c.Writer, form, &request); err != nil {
				fmt.Fprintln(c.Writer, "Pull/Sync failed:", err)
			}

			if err := captain.InstallChart(c.Writer, cfg.DefaultValue.CaptainUrl, &request, form.Verbose, 30*1000); err != nil {
				fmt.Fprintln(c.Writer, "call captain InstallChart failed:", err)
			}
			fmt.Fprintln(c.Writer, "InstallChart finish")
		}
	})
}

func pullAndSync(out io.Writer, form Request, request *proto.InstallChartRequest) error {
	var tpls chart.Templates
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
	}

	if strutil.Contains(form.Tags, "p") {
		// pull all image from chart
		for _, tpl := range tpls {
			for _, image := range tpl {
				fmt.Fprintln(out, "pulling ", image)
				result, err := dockerctl.Pull(*image)
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
						result, err := dockerctl.Retage(*image, chart.Image{
							Host: form.Registry,
							Repo: image.Repo,
							Tag:  image.Tag,
						})
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
