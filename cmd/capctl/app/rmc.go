package app

import (
	"encoding/hex"
	"encoding/json"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	pb "github.com/softleader/captain-kube/pkg/proto"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
)

const (
	rmcHelp = `在每個 worker node 上刪除 helm-chart 中的 image

使用 '--endpoint' 指定刪除的 Captain Endpoint

	$ {{.}} rmc CHART... -e CAPTAIN_ENDPOINT

傳入 '--constraint' 指定 TAG Semver2 的範圍條件 (https://devhints.io/semver)

	$ {{.}} rmc CHART... -e CAPTAIN_ENDPOINT -c "<"
	$ {{.}} rmc CHART... -e CAPTAIN_ENDPOINT -c ^	
	$ {{.}} rmc CHART... -e CAPTAIN_ENDPOINT -c ~

如果需要在 rmc 前修改 values.yaml 中任何參數, 可以傳入 '--set key1=val1,key2=val2'

	$ {{.}} rmc CHART... -e CAPTAIN_ENDPOINT --set ingress.enabled=true

傳入 '--force' 就算當前還有開啟 Container, 都強制刪除

	$ {{.}} rmc CHART... -e CAPTAIN_ENDPOINT --force

傳入 '--dry-run' 可以模擬真實的 rmi, 但不會真的刪除, 通常可以用來檢視 TAG 的 Semver2 範圍條件是否符合預期

	$ {{.}} rmc CHART... -e CAPTAIN_ENDPOINT --dry-run
`
)

type rmcCmd struct {
	hex        bool
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
		Short: "在每個 worker node 上刪除 helm-chart 中的 image",
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
	f.BoolVar(&c.force, "force", false, "force removal of the image")
	f.StringVarP(&c.constraint, "constraint", "c", "", "tag semver2 constraint, more details: https://devhints.io/semver")
	f.BoolVar(&c.dryRun, "dry-run", false, `simulate an rmc "for real"`)
	f.StringArrayVar(&c.set, "set", []string{}, "set values (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	f.BoolVar(&c.hex, "hex", false, "convert and upload chart into hex string")
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

		req := &pb.RmcRequest{
			Timeout:    settings.Timeout,
			DryRun:     c.dryRun,
			Force:      c.force,
			Color:      settings.Color,
			Verbose:    settings.Verbose,
			Constraint: c.constraint,
			Chart: &pb.Chart{
				FileName: filepath.Base(abs),
				FileSize: int64(len(bytes)),
			},
		}

		if c.hex {
			req.Chart.ContentHex = hex.EncodeToString(bytes)
			if logrus.IsLevelEnabled(logrus.DebugLevel) {
				v, _ := json.Marshal(req)
				logrus.Println(string(v)) // 如果是 hex string 印出來才有意義
			}
		} else {
			req.Chart.Content = bytes
		}

		if err := captain.CallRmc(logrus.StandardLogger(), c.endpoint.String(), req, settings.TimeoutDuration()); err != nil {
			return err
		}
	}
	return nil
}
