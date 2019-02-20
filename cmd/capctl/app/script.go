package app

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/softleader/captain-kube/pkg/utils/strutil"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	scriptHelp = `依照 helm-chart 內容產生 shell script

傳入 flag 產生對應的 docker script

	--pull-script    : docker pull IMAGE
	--re-tag-script  : docker tag IMAGE NEW_IMAGE && docker push NEW_IMAGE
	--save-script    : docker save IMAGE -o FILE
	--load-script    : docker load -i FILE

flags 可以自由的混搭使用, 你也可以使用 '>' 再將產生的 script 輸出成檔案

	$ {{.}} script CHART... -e CAPTAIN_ENDPOINT -prsl
	$ {{.}} script CHART... -e CAPTAIN_ENDPOINT -sl > save-and-load.sh

結合 '--diff' 可以只產生差異 image 的 script

	$ {{.}} script CHART ANOTHER_CHART -e CAPTAIN_ENDPOINT -prsld

如果需要在產生 script 前修改 values.yaml 中任何參數, 可以傳入 '--set key1=val1,key2=val2'

	$ {{.}} script CHART... -e CAPTAIN_ENDPOINT --set ingress.enabled=true

若 '--endpoint' 不指定則可在當前執行環境下產生 script, 而不是交給 Captain 執行
但請注意當前環境也必須要有產生 script 的必要 package

	$ {{.}} script CHART...

可以結合 '{{.}} ctx' 指令: 清空 context, 執行 script 後, 再切回原 context

	$ {{.}} ctx --off && {{.}} script CHART... && {{.}} ctx -
`
)

type scriptCmd struct {
	hex  bool
	pull bool
	rt   bool
	save bool
	load bool
	diff bool

	charts []string
	set    []string

	retag        *ctx.ReTag
	endpoint     *ctx.Endpoint // captain 的 endpoint ip
	registryAuth *ctx.RegistryAuth
}

func newScriptCmd() *cobra.Command {
	c := scriptCmd{
		endpoint:     activeCtx.Endpoint,
		registryAuth: activeCtx.RegistryAuth,
		retag:        activeCtx.ReTag,
	}

	cmd := &cobra.Command{
		Use:   "script CHART...",
		Short: "依照 helm-chart 內容產生 shell script",
		Long:  usage(scriptHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			if c.charts = args; len(c.charts) == 0 {
				return errors.New("chart path is required")
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	f.StringArrayVar(&c.set, "set", []string{}, "set values (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	f.BoolVarP(&c.pull, "pull-script", "p", c.pull, "generate pull images script")
	f.BoolVarP(&c.rt, "re-tag-script", "r", c.rt, "generate re-tag images script")
	f.BoolVarP(&c.save, "save-script", "s", c.save, "generate save images script")
	f.BoolVarP(&c.load, "load-script", "l", c.load, "generate load images script")
	f.BoolVarP(&c.diff, "diff", "d", c.diff, "show diff of two charts")
	f.BoolVar(&c.hex, "hex", false, "convert and upload chart into hex string")

	c.endpoint.AddFlags(f)
	c.registryAuth.AddFlags(f)
	c.retag.AddFlags(f)

	return cmd
}

func (c *scriptCmd) run() error {

	var buf *bytes.Buffer
	var scripts []string
	var log *logrus.Logger
	if c.diff {
		if l := len(c.charts); l != 2 {
			return fmt.Errorf("required two charts in diff mode, but received %v", l)
		}
		log = logrus.New() // 這個是這次請求要往前吐的 log
		buf = &bytes.Buffer{}
		//log.SetOutput(io.MultiWriter(&sseWriter, buf))
		log.SetOutput(buf)
		log.SetFormatter(&utils.PlainFormatter{})
	} else {
		log = logrus.StandardLogger()
	}

	for _, chart := range c.charts {
		logrus.Println("### chart:", chart, "###")
		if err := c.runScript(log, chart); err != nil {
			return err
		}
		logrus.Println("")
		logrus.Println("")

		// 如果buf裡面有存東西，則append到暫存裡面
		if buf != nil {
			scripts = append(scripts, buf.String())
			buf.Reset()
		}
	}

	// 如果暫存結果有存在，則進行差異比較
	if len(scripts) == 2 {
		logrus.Println("### Diffs: ###")
		lines := strutil.DiffNewLines(scripts[0], scripts[1])
		logrus.Println(strings.Join(lines, "\n"))
	}

	return nil
}

func (c *scriptCmd) runScript(log *logrus.Logger, path string) error {
	if expanded, err := homedir.Expand(path); err != nil {
		path = expanded
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if c.endpoint.Specified() {
		bytes, err := ioutil.ReadFile(abs)
		if err != nil {
			return err
		}
		request := captainkube_v2.GenerateScriptRequest{
			Chart: &captainkube_v2.Chart{
				FileName: filepath.Base(abs),
				FileSize: int64(len(bytes)),
			},
			Pull: c.pull,
			Retag: &captainkube_v2.ReTag{
				From: c.retag.From,
				To:   c.retag.To,
			},
			Save: c.save,
			Load: c.load,
		}

		if c.hex {
			request.Chart.ContentHex = hex.EncodeToString(bytes)
			if logrus.IsLevelEnabled(logrus.DebugLevel) {
				v, _ := json.Marshal(request)
				logrus.Println(string(v)) // 如果是 hex string 印出來才有意義
			}
		} else {
			request.Chart.Content = bytes
		}

		return captain.CallGenerateScript(log, c.endpoint.String(), &request, settings.TimeoutDuration())
	}

	return c.runScriptOnClient(log, abs)
}

func (c *scriptCmd) runScriptOnClient(log *logrus.Logger, path string) error {
	tpls, err := chart.LoadArchive(log, path, c.set...)
	if err != nil {
		return err
	}
	log.Debugf("%v template(s) loaded\n", len(tpls))

	if from, to := strings.TrimSpace(c.retag.From), strings.TrimSpace(c.retag.To); from != "" && to != "" {
		b, err := tpls.GenerateReTagScript(from, to)
		if err != nil {
			return err
		}
		log.Out.Write(b)
	}

	if c.pull {
		b, err := tpls.GeneratePullScript()
		if err != nil {
			return err
		}
		log.Out.Write(b)
	}

	if c.load {
		b, err := tpls.GenerateLoadScript()
		if err != nil {
			return err
		}
		log.Out.Write(b)
	}

	if c.save {
		b, err := tpls.GenerateSaveScript()
		if err != nil {
			return err
		}
		log.Out.Write(b)
	}
	return nil
}
