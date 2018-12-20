package app

import (
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/proto"
	"os"
	"testing"
)

func TestRetrieveServer(t *testing.T) {
	if _, err := retrieveServer("non-exist"); err == nil {
		t.Error("error must exist")
	}
}

func TestGrpc(t *testing.T) {
	//command := app.CaptainCmd{
	//	Port:           env.LookupInt(captain.EnvPort, captain.DefaultPort),
	//	CapletPort:     env.LookupInt(caplet.EnvPort, caplet.DefaultPort),
	//	CapletHostname: env.Lookup(caplet.EnvHostname, caplet.DefaultHostname),
	//}
	//command.Run()

	out := os.Stdout
	req := &proto.PullImageRequest{}
	req.Images = append(req.Images, &proto.Image{
		Host: "softleader",
		Repo: "helm",
	})

	ep := caplet.Endpoint{
		Target:  "localhost",
		Port:    caplet.DefaultPort,
		Timeout: dur.DefaultDeadlineSecond,
	}

	if err := ep.PullImage(out, req); err != nil {
		t.Error(err)
	}
}
