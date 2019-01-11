package dockerd

import (
	"context"
	"github.com/fsouza/go-dockerclient"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/utils/tcp"
	"testing"
)

func TestEncode(t *testing.T) {
	auth := captainkube_v2.RegistryAuth{
		Username:         "a",
		Password:         "b",
		XXX_unrecognized: []byte("c"),
		XXX_sizecache:    1,
	}

	encoded, err := encode(&auth)
	if err != nil {
		t.Error(err)
	}

	expected := "eyJVc2VybmFtZSI6ImEiLCJQYXNzd29yZCI6ImIifQ=="
	if encoded != expected {
		t.Errorf("expected %q, but get %q", expected, encoded)
	}
}

func TestIsDockerUnauthorized(t *testing.T) {
	if !tcp.IsReachable("hub.softleader.com.tw", 8443, 3) {
		t.Skipf("hub.softleader.com.tw:8443 not reachable")
	}
	cli, err := docker.NewClientFromEnv()
	if err != nil {
		t.Skipf("can not get docker client")
	}
	options := docker.PushImageOptions{
		Context: context.Background(),
		Name:    "hub.softleader.com.tw",
		Tag:     "busybox",
	}

	err = cli.PushImage(options, docker.AuthConfiguration{})
	if !isDockerUnauthorized(err) {
		t.Errorf("err should be unauthorized, but got %s", err)
	}
}
