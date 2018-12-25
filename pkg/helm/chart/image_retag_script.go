package chart

import (
	"github.com/sirupsen/logrus"
	"text/template"
)

const retagScript = `
{{- $from := index . "from" -}}
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
docker tag {{ $from }}/{{ $image.Name }} {{ $image.Host }}/{{ $image.Name }} && docker push {{ $image.Host }}/{{ $image.Name }}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
`

var retagTemplate = template.Must(template.New("").Parse(retagScript))

func (t *Templates) GenerateReTagScript(log *logrus.Logger, from, to string) error {
	retags := make(map[string][]*Image)
	for src, images := range *t {
		for _, image := range images {
			if image.Host == from {
				image.ReTag(from, to)
				retags[src] = append(retags[src], image)
			}
		}
	}
	data := make(map[string]interface{})
	data["from"] = from
	data["tpls"] = retags
	out := log.Writer()
	defer out.Close()
	return retagTemplate.Execute(out, data)
}
