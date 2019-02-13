package app

import (
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
)

const (
	rmcHelp = `刪除所有 node 上一個或多個 Chart 中的所有 image

使用 '--endpoint' 指定刪除的 Captain Endpoint

	$ {{.}} rmc CHART... -e CAPTAIN_ENDPOINT

傳入 '--range' 指定 TAG Semver2 的範圍條件 (https://devhints.io/semver)

	$ {{.}} rmc CHART... -e CAPTAIN_ENDPOINT -r "<"
	$ {{.}} rmc CHART... -e CAPTAIN_ENDPOINT -r ^	
	$ {{.}} rmc CHART... -e CAPTAIN_ENDPOINT -r ~

如果需要在 rmc 前修改 values.yaml 中任何參數, 可以傳入 '--set key1=val1,key2=val2'

	$ {{.}} rmc CHART... -e CAPTAIN_ENDPOINT --set ingress.enabled=true

傳入 '--force' 就算當前還有開啟 Container, 都強制刪除

	$ {{.}} rmc CHART... -e CAPTAIN_ENDPOINT -f

傳入 '--dry-run' 可以模擬真實的 rmi, 但不會真的刪除, 通常可以用來檢視 TAG 的 Semver2 範圍條件是否符合預期

	$ {{.}} rmc CHART... -e CAPTAIN_ENDPOINT --dry-run
`
)

type rmcCmd struct {
	force      bool
	charts     []string
	constraint string
	endpoint   *ctx.Endpoint // captain 的 endpoint ip
	dryRun     bool
	set        []string
}

func newRmcCmd() *cobra.Command {
	c := rmcCmd{
		endpoint: activeCtx.Endpoint,
	}

	cmd := &cobra.Command{
		Use:   "rmc CHARTS...",
		Short: "remove image of charts",
		Long:  usage(rmcHelp),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// do some validation check
			if err := c.endpoint.Validate(); err != nil {
				return err
			}
			c.charts = args
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&c.force, "force", "f", false, "force removal of the image")
	f.StringVarP(&c.constraint, "constraint", "c", "", "tag semver2 constraint, more details: https://devhints.io/semver")
	f.BoolVar(&c.dryRun, "dry-run", false, `simulate an rmc "for real"`)
	f.StringArrayVar(&c.set, "set", []string{}, "set values (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	c.endpoint.AddFlags(f)

	return cmd
}

func (c *rmcCmd) run() error {
	if c.dryRun {
		logrus.Warnln("running in dry-run mode, specify the '-v' flag if you want to turn on verbose output")
	}

	for _, chart := range c.charts {
		expanded, err := homedir.Expand(chart)
		if err != nil {
			expanded = chart
		}
		abs, err := filepath.Abs(expanded)
		if err != nil {
			return err
		}
		bytes, err := ioutil.ReadFile(abs)
		if err != nil {
			return err
		}

		req := &captainkube_v2.RmcRequest{
			Timeout:    settings.Timeout,
			DryRun:     c.dryRun,
			Force:      c.force,
			Color:      settings.Color,
			Verbose:    settings.Verbose,
			Constraint: c.constraint,
			Chart: &captainkube_v2.Chart{
				FileName: filepath.Base(abs),
				Content:  bytes,
				FileSize: int64(len(bytes)),
			},
		}
		if err := captain.Rmc(logrus.StandardLogger(), c.endpoint.String(), req, settings.Timeout); err != nil {
			return err
		}
	}
	return nil
}
