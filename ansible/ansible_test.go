package ansible

import (
	"testing"
	"encoding/json"
	"github.com/softleader/captain-kube/ansible/playbook"
)

func TestCommandOf(t *testing.T) {
	j := `{"inventory":"hosts-93","tags":["icp"],"version":"abc","namespace":"gardenia","chart":"softleader-jasmine","verbose":true}`
	b := playbook.NewStaging()
	err := json.Unmarshal([]byte(j), &b)
	if err != nil {
		t.Error(err)
	}
	if commandOf(b) != `ansible-playbook -i hosts-93 -t "icp" -e "version=abc chart=softleader-jasmine namespace=gardenia" -v staging.yml` {
		t.Error()
	}
}
