package app

import (
	"testing"
	"io/ioutil"
	"os"
	"path"
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
`

func TestGetDaemon(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Error(err.Error())
	}
	defer os.RemoveAll(tmp)
	err = ioutil.WriteFile(path.Join(tmp, daemonYaml), []byte(daemon), os.ModePerm)
	if err != nil {
		t.Error(err.Error())
	}
	d := GetDaemon(tmp)
	//fmt.Printf("%+v\n", d)
	if i := d.DefaultValue.Inventory; i != "hosts" {
		t.Errorf("Expected DefaultValue.Inventory is 'hosts', but got '%s'", i)
	}
	if n := d.DefaultValue.Namespace; n != "default" {
		t.Errorf("Expected DefaultValue.Namespace is 'default', but got '%s'", n)
	}
	if v := d.DefaultValue.Verbose; v {
		t.Errorf("Expected DefaultValue.Verbose true, but got '%v'", v)
	}
	if r := d.DefaultValue.SourceRegistry; r != "hub.softleader.com.tw" {
		t.Errorf("Expected DefaultValue.SourceRegistry is 'hub.softleader.com.tw', but got '%s'", r)
	}
}
