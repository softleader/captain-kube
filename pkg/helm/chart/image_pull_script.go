package chart

import (
	"bytes"
	"text/template"
)

const pullScript = `
{{ range $source, $images := index . "images" }}
##---
# Source: {{ $source }}
{{- range $key, $image := $images }}
docker pull {{ $image.Host }}/{{ $image.Name }}
{{- end }}
{{- end }}
`

var pullTemplate = template.Must(template.New("").Parse(pullScript))

func (i *Images) GeneratePullScript() (buf bytes.Buffer, err error) {
	data := make(map[string]interface{})
	data["images"] = i
	err = pullTemplate.Execute(&buf, data)
	return
}
