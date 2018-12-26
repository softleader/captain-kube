package dockerd

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"testing"
)

func TestEncode(t *testing.T) {
	auth := proto.RegistryAuth{
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

func TestPull(t *testing.T) {
	if err := Pull(logrus.StandardLogger(), chart.Image{
		Host: "library",
		Repo: "ubuntu",
		Tag:  "xenial",
		Name: "ubuntu:xenial",
	}, &proto.RegistryAuth{
		Username: "dev",
		Password: "sleader",
	}); err != nil {
		t.Error(err)
	}

	if err := Pull(logrus.StandardLogger(), chart.Image{
		Host: "hub.softleader.com.tw",
		Repo: "softleader-common-mail-rpc",
		Tag:  "latest",
		Name: "softleader-common-mail-rpc:latest",
	}, &proto.RegistryAuth{
		Username: "dev",
		Password: "sleader",
	}); err != nil {
		t.Error(err)
	}
}
