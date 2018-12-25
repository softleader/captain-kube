package dockerd

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"io/ioutil"
	"os"
)

func PullAndSync(log *logrus.Logger, request *proto.InstallChartRequest) error {
	var tpls chart.Templates
	if request.Pull || request.Sync {
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
		tpls, err = chart.LoadArchive(log, tmpFile.Name())
		if err != nil {
			return err
		}

	}

	if request.Pull {
		// pull all image from chart
		for _, tpl := range tpls {
			for _, image := range tpl {
				log.Println("pulling ", image)
				err := Pull(log, *image, request.RegistryAuth)
				if err != nil {
					log.Println("pull image failed: ", image, ", error: ", err)
				}
			}
		}
	}

	if request.Sync {
		if len(request.Retag.From) > 0 && len(request.Retag.To) > 0 {
			// retag and push all image from chart
			for _, tpl := range tpls {
				for _, image := range tpl {
					if image.Host == request.Retag.From {
						log.Println("syncing ", image)
						err := ReTag(log, *image, chart.Image{
							Host: request.Retag.To,
							Repo: image.Repo,
							Tag:  image.Tag,
						}, request.RegistryAuth)
						if err != nil {
							log.Println("sync image failed: ", image, ", error: ", err)
						}
					}
				}
			}
		} else {
			return errors.New("Registry and SourceRegistry is required when retag mode")
		}
	}

	return nil
}
