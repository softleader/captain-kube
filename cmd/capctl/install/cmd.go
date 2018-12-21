package install

import (
	"fmt"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/dockerctl"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"strconv"
)

type installCmd struct {
	out            io.Writer
	pull           bool
	sync           bool
	k8sVendor      string
	namespace      string
	sourceRegistry string
	registry       string
	chartPath      string
	verbose        bool
	timeout        int64

	registryAuthUsername string
	registryAuthPassword string

	tillerEndpoint          string
	tillerUsername          string
	tillerPassword          string
	tillerAccount           string
	tillerSkipSslValidation bool

	captainUrl string
}

func NewCmd(out io.Writer) *cobra.Command {
	c := installCmd{
		out:       out,
		k8sVendor: env.Lookup(captain.EnvK8sVendor, captain.DefaultK8sVendor),
		namespace: "default",
		chartPath: "./chart.tgz",
		timeout:   300,

		registryAuthUsername: env.Lookup(captain.EnvRegistryAuthUsername, captain.DefaultRegistryAuthUsername),
		registryAuthPassword: env.Lookup(captain.EnvRegistryAuthPassword, captain.DefaultRegistryAuthPassword),

		tillerEndpoint:          env.Lookup(captain.EnvTillerEndpoint, captain.DefaultTillerEndpoint),
		tillerUsername:          env.Lookup(captain.EnvTillerUsername, captain.DefaultTillerUsername),
		tillerPassword:          env.Lookup(captain.EnvTillerPassword, captain.DefaultTillerPassword),
		tillerAccount:           env.Lookup(captain.EnvTillerAccount, captain.DefaultTillerAccount),
		tillerSkipSslValidation: env.LookupBool(captain.EnvTillerSkipSslValidation, captain.DefaultTillerSkipSslValidation),

		captainUrl: fmt.Sprintf("localhost:%v", env.LookupInt(captain.EnvPort, captain.DefaultPort)),
	}

	cmd := &cobra.Command{
		Use:   "install",
		Short: "install helm chart",
		Long:  "install helm chart",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&c.pull, "pull", "p", c.pull, "Pull images in Chart")
	f.BoolVarP(&c.sync, "sync", "s", c.sync, "Retag images & Sync to install nodes")
	f.StringVar(&c.k8sVendor, "k8s-vendor", c.k8sVendor, "specify the vendor of k8s, override "+captain.EnvK8sVendor)
	f.StringVar(&c.namespace, "namespace", c.namespace, "specify the namespace of gcp, not available now")
	f.StringVar(&c.sourceRegistry, "src", c.sourceRegistry, "specify the host of Retage from, reqiured when Sync")
	f.StringVar(&c.registry, "tgt", c.registry, "specify the host of Retage to, reqiured when Sync")
	f.StringVarP(&c.chartPath, "chartPath", "f", c.chartPath, "specify the path of chart, must be a tgz file, default: "+c.chartPath)
	f.Int64VarP(&c.timeout, "timeout", "t", c.timeout, "seconds of captain run timeout, default: "+string(c.timeout))

	f.StringVar(&c.registryAuthUsername, "reguser", c.registryAuthUsername, "specify the registryAuthUsername, reqiured when Pull&Sync, default: "+c.registryAuthUsername)
	f.StringVar(&c.registryAuthPassword, "regpass", c.registryAuthPassword, "specify the registryAuthPassword, reqiured when Pull&Sync, default: "+c.registryAuthPassword)

	f.StringVar(&c.tillerEndpoint, "tiller", c.tillerEndpoint, "specify the tillerEndpoint, default: "+c.tillerEndpoint)
	f.StringVar(&c.tillerUsername, "tillerU", c.tillerUsername, "specify the tillerUsername, default: "+c.tillerUsername)
	f.StringVar(&c.tillerPassword, "tillerP", c.tillerPassword, "specify the tillerPassword, default: "+c.tillerPassword)
	f.StringVar(&c.tillerAccount, "tillerA", c.tillerAccount, "specify the tillerAccount, default: "+c.tillerAccount)
	f.BoolVar(&c.tillerSkipSslValidation, "tillerS", c.tillerSkipSslValidation, "specify the tillerSkipSslValidation, default: "+strconv.FormatBool(c.tillerSkipSslValidation))

	f.StringVar(&c.captainUrl, "captainUrl", c.captainUrl, "specify the captainUrl, default: "+c.captainUrl)

	return cmd
}

func (c *installCmd) run() error {
	bytes, err := ioutil.ReadFile(c.chartPath)
	if err != nil {
		return err
	}

	request := proto.InstallChartRequest{
		Chart: &proto.Chart{
			FileName: "chart.tgz", //FIXME 需要動態的取得filename, 否則到時候解壓時有可能會被附檔名誤導
			Content:  bytes,
			FileSize: int64(len(bytes)),
		},
		Pull: c.pull,
		Sync: c.sync,
		Retag: &proto.ReTag{
			From: c.sourceRegistry,
			To:   c.registry,
		},
		Tiller: &proto.Tiller{
			Endpoint:          c.tillerEndpoint,
			Username:          c.tillerUsername,
			Password:          c.tillerPassword,
			Account:           c.tillerAccount,
			SkipSslValidation: c.tillerSkipSslValidation,
		},
		RegistryAuth: &proto.RegistryAuth{
			Username: c.registryAuthUsername,
			Password: c.registryAuthPassword,
		},
		Timeout: c.timeout,
	}

	if err := dockerctl.PullAndSync(c.out, &request); err != nil {
		return err
	}

	if err := captain.InstallChart(c.out, c.captainUrl, &request, 300); err != nil {
		return err
	}

	return nil
}
