package ansible

import (
	"testing"
	"encoding/json"
	"github.com/softleader/captain-kube/ansible/playbook"
	"github.com/softleader/captain-kube/charts"
)

//const daemon = `inventory: hosts
//tags:
//  - icp
//    #  - retag
//namespace: default
//version:
//verbose: true
//sourceRegistry: hub.softleader.com.tw
//registry:`
//
//func TestExtendsDefaultValues(t *testing.T) {
//	tmp, err := ioutil.TempDir(os.TempDir(), "")
//	if err != nil {
//		t.Error(err.Error())
//	}
//	err = ioutil.WriteFile(path.Join(tmp, daemonYaml), []byte(daemon), os.ModePerm)
//	if err != nil {
//		t.Error(err.Error())
//	}
//	dft := playbook.NewStaging()
//	ExtendsDefaultValues(tmp, dft)
//	if i := dft.Inventory; i != "hosts" {
//		t.Errorf("Inventory should be hosts, but was %v", i)
//	}
//	if l := len(dft.Tags); l != 1 {
//		t.Errorf("Tags length should be 1, but was %v", l)
//	}
//}

func TestCommandOf(t *testing.T) {
	j := `{"inventory":"hosts","tags":["icp"],"version":"abc","namespace":"gardenia","chart":"softleader-jasmine","verbose":true}`
	b := playbook.NewStaging()
	err := json.Unmarshal([]byte(j), &b)
	if err != nil {
		t.Error(err)
	}
	expected := `ansible-playbook -i hosts -t "icp" -e '{"chart":"softleader-jasmine","chart_path":"","images":null,"namespace":"gardenia","version":"abc"}' -v staging.yml`
	actual := commandOf(b)
	if actual != expected {
		t.Error("\nexpected:\n", expected, "\nactual was:\n", actual)
	}
}

//func TestPlay(t *testing.T) {
//	play := `ansible-playbook -i hosts -t "icp" -e "version=abc chart=softleader-jasmine chart_path= script= script_path= namespace=gardenia" -v staging.yml`
//	opts := sh.Options{
//		Pwd:     "/Users/Matt/go/src/github.com/softleader/captain-kube/docs/playbooks",
//		Verbose: true,
//	}
//	sh.C(&opts, play)
//}

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
	expected := `{"chart":"bbb","chart_path":"xxx","images":[{"Registry":"hub.softleader.com.tw","Name":"a"},{"Registry":"hub.softleader.com.tw","Name":"b"},{"Registry":"hub.softleader.com.tw","Name":"c"}],"namespace":"default","version":"aaa"}`
	actual := string(b)
	if actual != expected {
		t.Error("\nexpected:\n", expected, "\nactual was:\n", actual)
	}
}
