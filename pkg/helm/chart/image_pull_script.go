package chart

import (
	"bytes"
	"text/template"
)

const pullScript = `
{{ range $path, $images := index . "tpls" }}
##---
# Source: {{ $path }}
{{- range $key, $image := $images }}
docker pull {{ $image.Host }}/{{ $image.Name }}
{{- end }}
{{- end }}
`

var pullTemplate = template.Must(template.New("").Parse(pullScript))

func (i *Templates) GeneratePullScript() (buf bytes.Buffer, err error) {
	data := make(map[string]interface{})
	data["tpls"] = i
	err = pullTemplate.Execute(&buf, data)
	return
}
