package tmpl

import (
	"testing"
	"fmt"
	"github.com/softleader/captain-kube/charts"
	"encoding/json"
)

func TestCompile(t *testing.T) {
	text := `#!/usr/bin/env bash
{{ range $source, $images := index . "images" }}
##---
# Source: {{ $source }}
{{- range $key, $image := $images }}
docker pull {{ $image.Registry }}/{{ $image.Name }}
{{- end }}
{{- end }}

exit 0`

	data := make(map[string]interface{})
	images := make(map[string]interface{})
	images["base-deployment.yaml"] = []charts.Image{{
		Registry: "hub.softleader.com.tw",
		Name:     "a",
	}, {
		Registry: "hub.softleader.com.tw",
		Name:     "b",
	}, {
		Registry: "hub.softleader.com.tw",
		Name:     "c",
	}}
	images["gateway-deployment.yaml"] = []charts.Image{{
		Registry: "hub.softleader.com.tw",
		Name:     "d",
	}, {
		Registry: "hub.softleader.com.tw",
		Name:     "e",
	}}
	data["registry"] = "softleader"
	data["images"] = images

	b, err := json.Marshal(data)
	fmt.Println(string(b))
	buf, err := compile(text, data)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(buf.String())
}

func TestCompile2(t *testing.T) {
	data := make(map[string]interface{})
	images := make(map[string]interface{})
	images["base-deployment.yaml"] = []charts.Image{{
		Registry: "hub.softleader.com.tw",
		Name:     "a",
	}, {
		Registry: "hub.softleader.com.tw",
		Name:     "b",
	}, {
		Registry: "hub.softleader.com.tw",
		Name:     "c",
	}}
	images["gateway-deployment.yaml"] = []charts.Image{{
		Registry: "hub.softleader.com.tw",
		Name:     "d",
	}, {
		Registry: "hub.softleader.com.tw",
		Name:     "e",
	}}
	data["registry"] = "softleader"
	data["images"] = images

	text := `#!/usr/bin/env bash
{{ $registry := index . "registry" }}
{{- range $source, $images := index . "images" }}
##---
# Source: {{ $source }}
{{- range $key, $image := $images }}
docker tag {{ $image.Registry }}/{{ $image.Name }} {{ $registry }}/{{ $image.Name }} && docker push {{ $registry }}/{{ $image.Name }}
{{- end }}
{{- end }}

exit 0`
	buf, err := compile(text, data)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(buf.String())
}
