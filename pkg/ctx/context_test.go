package ctx

import (
	"fmt"
	"github.com/softleader/captain-kube/pkg/captain"
	"os"
	"strconv"
	"testing"
)

func TestExpandEnv(t *testing.T) {
	ctx, err := newContext()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	ctx.HelmTiller.Endpoint = "192.168.1.93"
	ctx.HelmTiller.Account = "hello-tiller"
	ctx.HelmTiller.Password = "secret"
	os.Setenv(captain.DefaultTillerPassword, "surprised")
	os.Setenv(captain.EnvPort, strconv.Itoa(captain.DefaultPort))

	if err := ctx.expandEnv(); err != nil {
		t.Error(err)
	}
	if e := ctx.HelmTiller.Endpoint; e != "192.168.1.93" {
		t.Errorf("new endpoint should be 192.168.1.93, but got %s", e)
	}
	if a := ctx.HelmTiller.Account; a != "hello-tiller" {
		t.Errorf("new account should be hello-tiller, but got %s", a)
	}
	if p := ctx.HelmTiller.Password; p != "secret" {
		t.Errorf("new endpoint should be secret, but got %s", p)
	}
	if p := ctx.Endpoint.Port; p != captain.DefaultPort {
		t.Errorf("new endpoint port should be %v, but got %v", captain.DefaultPort, p)
	}
}

func TestUsageString(t *testing.T) {
	u, err := FlagsString()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(u)
}
