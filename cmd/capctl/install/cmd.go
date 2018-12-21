package install

import (
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/dockerctl"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/spf13/cobra"
	"io"
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
}

func NewCmd(out io.Writer) *cobra.Command {
	c := installCmd{
		out:       out,
		k8sVendor: env.Lookup(captain.EnvK8sVendor, captain.DefaultK8sVendor),
		namespace: "default",
		chartPath: "./chart.tgz",
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
	f.StringVarP(&c.chartPath, "chartPath", "f", c.chartPath, "specify the path of char, must be a tgz file, default: "+c.chartPath)

	return cmd
}

func (c *installCmd) run() error {

	dockerctl.PullAndSync(c.out, &proto.InstallChartRequest{})

	return nil
}
