package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/dockerd"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/spf13/cobra"
)

const (
	rmiHelp = `刪除一個或多個 docker image

使用 '--endpoint' 指定刪除的 Captain Endpoint

	$ {{.}} rmi IMAGE:TAG... -e CAPTAIN_ENDPOINT

IMAGE 中可以使用 * 做模糊查詢, 如刪除包含 gateway 字眼的 image

	$ {{.}} rmi hub.softleader.com.tw/*gateway*:TAG -e CAPTAIN_ENDPOINT

TAG 必須要指定, 可以是絕對條件或是 Semver2 的範圍條件 (https://devhints.io/semver)

	$ {{.}} rmi IMAGE:2.1.3 -e CAPTAIN_ENDPOINT
	$ {{.}} rmi IMAGE:^2.1.3 -e CAPTAIN_ENDPOINT
	$ {{.}} rmi IMAGE:~2.1.3 -e CAPTAIN_ENDPOINT

傳入 '--force' 就算當前還有開啟 Container, 都強制刪除

	$ {{.}} rmi IMAGE:TAG -e CAPTAIN_ENDPOINT -f

傳入 '--dry-run' 可以模擬真實的 rmi, 但不會真的刪除, 通常可以用來檢視 TAG 的 Semver2 範圍條件是否符合預期

	$ {{.}} rmi IMAGE:<=2.0.0 -e CAPTAIN_ENDPOINT --dry-run

若 '--endpoint' 不指定則可在當前執行環境下執行 rmi, 而不是 Kubernetes cluster

	$ {{.}} rmi IMAGE:TAG...

可以結合 '{{.}} ctx' 指令: 清空 context, 執行 rmi 後, 再切回原 context

	$ {{.}} ctx --off && {{.}} rmi IMAGE:TAG... && {{.}} ctx -
`
)

type rmiCmd struct {
	force    bool
	images   []*chart.Image
	endpoint *ctx.Endpoint // captain 的 endpoint ip
	dryRun   bool
}

func newRmiCmd() *cobra.Command {
	c := rmiCmd{
		endpoint: activeCtx.Endpoint,
	}

	cmd := &cobra.Command{
		Use:   "rmi IMAGES...",
		Short: "remove images",
		Long:  usage(rmiHelp),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				image := chart.NewImage(arg)
				if len(image.Tag) == 0 {
					return fmt.Errorf("%q must specify tag", arg)
				}
				c.images = append(c.images, image)
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&c.force, "force", "f", false, "force removal of the image")
	f.BoolVar(&c.dryRun, "dry-run", false, `simulate an rmi "for real"`)
	c.endpoint.AddFlags(f)

	return cmd
}

func (c *rmiCmd) run() error {
	if c.dryRun {
		logrus.Warnln("running in dry-run mode, specify the '-v' flag if you want to turn on verbose output")
	}

	if c.endpoint.Specified() {
		req := &captainkube_v2.RmiRequest{
			Timeout: settings.Timeout,
			DryRun:  c.dryRun,
			Force:   c.force,
			Color:   settings.Color,
			Verbose: settings.Verbose,
		}
		for _, i := range c.images {
			req.Images = append(req.Images, &captainkube_v2.Image{
				Host: i.Host,
				Repo: i.Repo,
				Tag:  i.Tag,
			})
		}
		return captain.Rmi(logrus.StandardLogger(), c.endpoint.String(), req, settings.Timeout)
	}

	for _, image := range c.images {
		rm, err := dockerd.ImagesWithTagConstraint(logrus.StandardLogger(), image.HostRepo(), image.Tag)
		if err != nil {
			return err
		}
		if err := dockerd.Rmi(logrus.StandardLogger(), c.force, c.dryRun, rm...); err != nil {
			return err
		}
	}
	return nil
}
