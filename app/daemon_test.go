package app

import (
	"testing"
	"io/ioutil"
	"os"
	"path"
	"fmt"
	"github.com/softleader/captain-kube/ansible/playbook"
	"gopkg.in/yaml.v2"
)

const daemon = `
defaultValue:
  inventory: hosts
  tags:
  - icp
  namespace: default 
  version: 
  verbose: false
  sourceRegistry: hub.softleader.com.tw
  registry:
kubernetes:
  url: https://192.168.1.93:8443/
`

func TestLoadDaemonYAML(t *testing.T) {
	d := Daemon{}
	yaml.Unmarshal([]byte(daemon), &d)
	fmt.Printf("%+v", d)
}

func TestExtendsTo(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Error(err.Error())
	}
	err = ioutil.WriteFile(path.Join(tmp, daemonYaml), []byte(daemon), os.ModePerm)
	if err != nil {
		t.Error(err.Error())
	}
	d := GetDaemon(tmp)
	fmt.Printf("%+v\n", d)

	s := playbook.NewStaging()
	// d.DefaultValue.DeepCopyTo(&s)
	fmt.Printf("%+v", s)
}
