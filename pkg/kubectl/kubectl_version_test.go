package kubectl

import (
	"fmt"
	"github.com/softleader/captain-kube/pkg/utils/command"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestVersion(t *testing.T) {
	if !command.IsAvailable(kubectl) {
		t.Skipf("%q command does not exist", kubectl)
	}
	v, err := versionClientOnly()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("%+v\n", v)
}

func TestVersionUnmarshal(t *testing.T) {
	ver := `clientVersion:
  buildDate: 2018-10-05T16:46:06Z
  compiler: gc
  gitCommit: 4ed3216f3ec431b140b1d899130a69fc671678f4
  gitTreeState: clean
  gitVersion: v1.12.1
  goVersion: go1.10.4
  major: "1"
  minor: "12"
  platform: darwin/amd64
serverVersion:
  buildDate: 2018-02-23T07:20:41Z
  compiler: gc
  gitCommit: d97ba3f083461e0ae0a8881550e83350af4c8f57
  gitTreeState: clean
  gitVersion: v1.9.1+icp-ee
  goVersion: go1.9.2
  major: "1"
  minor: "9"
  platform: linux/amd64`
	kv := &KubeVersion{}
	if err := yaml.Unmarshal([]byte(ver), kv); err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("%+v", kv)
}
