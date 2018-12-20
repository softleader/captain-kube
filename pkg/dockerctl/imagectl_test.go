package dockerctl

import (
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
