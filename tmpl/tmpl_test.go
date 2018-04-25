package tmpl

import (
	"testing"
	"fmt"
	"github.com/softleader/captain-kube/chart"
)

func TestCompile(t *testing.T) {
	text := `
{{ range $key, $value := .}}
{{ $value }}
{{ end }}
`
	buf, err := compile(text, []string{"a", "b", "c"})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(buf.String())
}

func TestCompile2(t *testing.T) {
	retag := Retag{
		Registry: "hub.softleader.com.tw",
		Images:   []charts.Image{{Name: "a"}, {Name: "b"}, {Name: "c"}},
	}
	text := `{{ $registry := .Registry }}
{{ range $key, $value := .Images}}
{{ $registry }}/{{ $value.Name }}
{{ end }}
`
	buf, err := compile(text, retag)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(buf.String())
}

type Retag struct {
	Registry string
	Images   []charts.Image
}
