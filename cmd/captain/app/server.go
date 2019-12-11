package app

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/cmd/captain/app/server"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/softleader/captain-kube/pkg/kubectl"
	pb "github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/release"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
)

type captainCmd struct {
	metadata       *release.Metadata
	serve          string
	endpoints      []string
	port           int
	capletHostname string
	capletPort     int
	k8sVendor      string
}

// NewCaptainCommand 建立 Captain root command
func NewCaptainCommand(metadata *release.Metadata) (cmd *cobra.Command) {
	var verbose bool
	c := captainCmd{
		metadata:       metadata,
		port:           env.LookupInt(captain.EnvPort, captain.DefaultPort),
		capletPort:     env.LookupInt(caplet.EnvPort, caplet.DefaultPort),
		capletHostname: env.Lookup(caplet.EnvHostname, caplet.DefaultHostname),
	}
	cmd = &cobra.Command{
		Use:  "captain",
		Long: "captain is the brain of captain-kube system",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			logrus.SetOutput(colorable.NewColorableStdout()) // for windows color output
			if verbose {
				logrus.SetLevel(logrus.DebugLevel)
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	f.IntVarP(&c.port, "port", "p", c.port, "specify the port of captain, override "+captain.EnvPort)
	f.StringVar(&c.capletHostname, "caplet-hostname", c.capletHostname, "specify the hostname of caplet daemon to lookup, override "+caplet.EnvHostname)
	f.IntVar(&c.capletPort, "caplet-port", c.capletPort, "specify the port of caplet daemon to connect, override "+caplet.EnvPort)
	f.StringArrayVarP(&c.endpoints, "caplet-endpoint", "e", []string{}, "specify the endpoint of caplet daemon to connect, override '--caplet-hostname'")
	f.StringVar(&c.k8sVendor, "k8s-vendor", c.k8sVendor, "specify the vendor of k8s")

	return
}

func (c *captainCmd) run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", c.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	srv := &server.CaptainServer{
		Metadata:  c.metadata,
		Hostname:  c.capletHostname,
		Endpoints: c.endpoints,
		Port:      c.capletPort,
	}
	if len(c.k8sVendor) > 0 {
		logrus.Printf("server has specified k8s vendor to %q, skipping auto detection", c.k8sVendor)
		srv.K8s = kubectl.NewKubeVersion(c.k8sVendor)
	} else {
		srv.K8s, err = kubectl.Version()
		if err != nil {
			return err
		}
		logrus.Printf("detected k8s vendor: %s", srv.K8s.Server.GitVersion)
	}
	s := grpc.NewServer()
	pb.RegisterCaptainServer(s, srv)
	logrus.Printf("registered captain server\n")

	healthpb.RegisterHealthServer(s, health.NewServer())
	logrus.Printf("registered health probe\n")

	reflection.Register(s)
	logrus.Printf("listening and serving GRPC on %v", lis.Addr().String())
	return s.Serve(lis)
}
