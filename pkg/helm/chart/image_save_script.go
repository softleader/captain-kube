package chart

import (
	"github.com/Sirupsen/logrus"
	"text/template"
)

const saveScript = `
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
docker save -o ./{{ $image.Name }}.tar {{ $image.String }}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
`

var saveTemplate = template.Must(template.New("").Parse(saveScript))

func (t *Templates) GenerateSaveScript(log *logrus.Logger) error {
	data := make(map[string]interface{})
	data["tpls"] = t
	out := log.Writer()
	defer out.Close()
	return saveTemplate.Execute(out, data)
}
