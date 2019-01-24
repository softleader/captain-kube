package kubectl

import (
	"fmt"
	"github.com/softleader/captain-kube/pkg/utils/command"
	"testing"
)

func TestVersion(t *testing.T) {
	if !command.IsAvailable(kubectl) {
		t.Skipf("%q command does not exist", kubectl)
	}
	v, err := Version()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("%+v\n", v)
}
