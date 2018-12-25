package chart

import (
	"github.com/sirupsen/logrus"
	"text/template"
)

const loadScript = `
{{- $tpls := index . "tpls" -}}
{{- $len := len $tpls -}} 
{{- if eq $len 0 -}}
# no sources found in template
{{- else -}}
{{- range $path, $images := $tpls -}}
##---
# Source: {{ $path }}
{{- $len = len $images -}} 
{{- if eq $len 0 -}}
# no images found in source
{{- else -}}
{{- range $key, $image := $images }}
docker load -i ./{{ $image.Name }}.tar
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
`

var loadTemplate = template.Must(template.New("").Parse(loadScript))

func (t *Templates) GenerateLoadScript(log *logrus.Logger) error {
	data := make(map[string]interface{})
	data["tpls"] = t
	out := log.Writer()
	defer out.Close()
	return loadTemplate.Execute(out, data)
}
