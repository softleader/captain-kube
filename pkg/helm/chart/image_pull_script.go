package chart

import (
	"github.com/softleader/captain-kube/pkg/logger"
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

func (i *Templates) GeneratePullScript(log *logger.Logger) error {
	data := make(map[string]interface{})
	data["tpls"] = i
	return pullTemplate.Execute(log.GetOutput(), data)
}
