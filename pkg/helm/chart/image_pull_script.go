package chart

import (
	"github.com/Sirupsen/logrus"
	"text/template"
)

const pullScript = `
{{ range $path, $images := index . "tpls" }}
##---
# Source: {{ $path }}
{{- range $key, $image := $images }}
docker pull {{ $image.String }}
{{- end }}
{{- end }}
`

var pullTemplate = template.Must(template.New("").Parse(pullScript))

func (i *Templates) GeneratePullScript(log *logrus.Logger) error {
	data := make(map[string]interface{})
	data["tpls"] = i
	return pullTemplate.Execute(log.WriterLevel(logrus.DebugLevel), data)
}
