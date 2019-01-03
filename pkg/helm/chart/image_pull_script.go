package chart

import (
	"bytes"
	"text/template"
)

const pullScript = `
{{- $tpls := index . "tpls" -}}
{{- $len := len $tpls -}} 
{{- if eq $len 0 }}
# no sources found in template
{{- else -}}
{{- range $path, $images := $tpls }}
##---
# Source: {{ $path }}
{{- $len = len $images -}} 
{{- if eq $len 0 }}
# no images found in source
{{- else -}}
{{- range $key, $image := $images }}
docker pull {{ $image.String }}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
`

var pullTemplate = template.Must(template.New("").Parse(pullScript))

func (t *Templates) GeneratePullScript() ([]byte, error) {
	data := make(map[string]interface{})
	data["tpls"] = t
	var buf bytes.Buffer
	if err := pullTemplate.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
