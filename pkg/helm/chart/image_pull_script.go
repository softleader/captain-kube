package chart

import (
	"github.com/Sirupsen/logrus"
	"text/template"
)

const pullScript = `
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
docker pull {{ $image.String -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
`

var pullTemplate = template.Must(template.New("").Parse(pullScript))

func (t *Templates) GeneratePullScript(log *logrus.Logger) error {
	data := make(map[string]interface{})
	data["tpls"] = t
	out := log.Writer()
	defer out.Close()
	return pullTemplate.Execute(out, data)
}
