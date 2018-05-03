package ansible

import (
	"testing"
	"encoding/json"
	"github.com/softleader/captain-kube/ansible/playbook"
	"github.com/softleader/captain-kube/sh"
	"fmt"
	"github.com/softleader/captain-kube/charts"
)

func TestExtendsDefaultValues(t *testing.T) {
	dft := playbook.NewStaging()
	ExtendsDefaultValues("/Users/Matt/tmp/anible", dft)
	fmt.Printf("%+v", dft)
}

func TestCommandOf(t *testing.T) {
	j := `{"inventory":"hosts","tags":["icp"],"version":"abc","namespace":"gardenia","chart":"softleader-jasmine","verbose":true}`
	b := playbook.NewStaging()
	err := json.Unmarshal([]byte(j), &b)
	if err != nil {
		t.Error(err)
	}
	expected := `ansible-playbook -i hosts -t "icp" -e "version=abc chart=softleader-jasmine chart_path= script= script_path= namespace=gardenia" -v staging.yml`
	actual := commandOf(b)
	if actual != expected {
		t.Error("\nexpected:\n", expected, "\nactual was:\n", actual)
	}
}

func TestPlay(t *testing.T) {
	play := `ansible-playbook -i hosts -t "icp" -e "version=abc chart=softleader-jasmine chart_path= script= script_path= namespace=gardenia" -v staging.yml`

	opts := sh.Options{
		Pwd:     "/Users/Matt/go/src/github.com/softleader/captain-kube/docs/playbooks",
		Verbose: true,
	}
	sh.C(&opts, play)
}

func TestPrintAnsibleExtraVars(t *testing.T) {
	e := make(map[string]interface{})
	e["version"] = "aaa"
	e["chart"] = "bbb"
	e["chart_path"] = "xxx"
	e["namespace"] = "default"
	e["images"] = []charts.Image{
		{Registry: "hub.softleader.com.tw", Name: "a"}, {Registry: "hub.softleader.com.tw", Name: "b"}, {Registry: "hub.softleader.com.tw", Name: "c"},
	}
	b, _ := json.Marshal(e)
	fmt.Println(string(b))
}
