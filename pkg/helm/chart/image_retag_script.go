package chart

import (
	"github.com/Sirupsen/logrus"
	"text/template"
)

const retagScript = `
{{ $from := index . "from" }}
{{- range $path, $images := index . "tpls" }}
##---
# Source: {{ $path }}
{{- range $key, $image := $images }}
docker tag {{ $from }}/{{ $image.Name }} {{ $image.Host }}/{{ $image.Name }} && docker push {{ $image.Host }}/{{ $image.Name }}
{{- end }}
{{- end }}
`

var retagTemplate = template.Must(template.New("").Parse(retagScript))

func (i *Templates) GenerateReTagScript(log *logrus.Logger, from, to string) error {
	var retags map[string][]*Image
	for src, images := range *i {
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
	return retagTemplate.Execute(log.Writer(), data)
}
