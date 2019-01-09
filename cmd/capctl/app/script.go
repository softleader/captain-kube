package app

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/softleader/captain-kube/pkg/utils/strutil"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	scriptHelp = `依照一或多個 Helm Chart 內容產生 docker scripts

傳入 flag 產生對應的 docker script

	--pull    : docker pull IMAGE
	--re-tag  : docker tag IMAGE NEW_IMAGE && docker push NEW_IMAGE
	--save    : docker save IMAGE -o FILE
	--load    : docker load -i FILE

flags 可以自由的混搭使用, 你也可以使用 '>' 再將產生的 script 輸出成檔案

	$ {{.}} script CHART... -e CAPTAIN_ENDPOINT -prsl
	$ {{.}} script CHART... -e CAPTAIN_ENDPOINT -sl > save-and-load.sh

結合 '--diff' 可以只產生差異 image 的 script

	$ {{.}} script CHART ANOTHER_CHART -e CAPTAIN_ENDPOINT -prsld
`
)

type scriptCmd struct {
	pull bool
	rt   bool
	save bool
	load bool
	diff bool

	charts []string

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
		Short: "generate script of helm-chart",
		Long:  usage(scriptHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			if c.charts = args; len(c.charts) == 0 {
				return errors.New("chart path is required")
			}
			if err := c.endpoint.Validate(); err != nil {
				return err
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&c.pull, "pull", "p", c.pull, "pull images in Chart")
	f.BoolVarP(&c.rt, "re-tag", "r", c.rt, "re-tag images in Chart")
	f.BoolVarP(&c.save, "save", "s", c.save, "save images in Chart")
	f.BoolVarP(&c.load, "load", "l", c.load, "load images in Chart")
	f.BoolVarP(&c.diff, "diff", "d", c.diff, "show diff of two charts")

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
		} else {
			log = logrus.New() // 這個是這次請求要往前吐的 log
			buf = &bytes.Buffer{}
			//log.SetOutput(io.MultiWriter(&sseWriter, buf))
			log.SetOutput(buf)
			log.SetFormatter(&utils.PlainFormatter{})
		}
	} else {
		log = logrus.StandardLogger()
	}

	for _, chart := range c.charts {
		logrus.Println("### chart:", chart, "###")
		if err := runScript(log, c, chart); err != nil {
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

func runScript(log *logrus.Logger, c *scriptCmd, path string) error {
	if expanded, err := homedir.Expand(path); err != nil {
		path = expanded
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(abs)
	if err != nil {
		return err
	}

	request := proto.GenerateScriptRequest{
		Chart: &proto.Chart{
			FileName: filepath.Base(abs),
			Content:  bytes,
			FileSize: int64(len(bytes)),
		},
		Pull: c.pull,
		Retag: &proto.ReTag{
			From: c.retag.From,
			To:   c.retag.To,
		},
		Save: c.save,
		Load: c.load,
	}

	if err := captain.GenerateScript(log, c.endpoint.String(), &request, settings.Timeout); err != nil {
		return err
	}

	return nil
}
