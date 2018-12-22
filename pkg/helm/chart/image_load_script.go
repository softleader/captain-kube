package chart

import (
	"github.com/softleader/captain-kube/pkg/logger"
	"text/template"
)

const loadScript = `
{{- range $path, $images := index . "tpls" }}
##---
# Source: {{ $path }}
{{- range $key, $image := $images }}
docker load -i ./{{ $image.Name }}.tar
{{- end }}
{{- end }}
`

var loadTemplate = template.Must(template.New("").Parse(loadScript))

func (i *Templates) GenerateLoadScript(log *logger.Logger) error {
	data := make(map[string]interface{})
	data["tpls"] = i
	return loadTemplate.Execute(log.GetOutput(), data)
}
