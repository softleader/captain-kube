package chart

import (
	"github.com/softleader/captain-kube/pkg/logger"
	"text/template"
)

const saveScript = `
{{ $from := index . "from" }}
{{- range $path, $images := index . "tpls" }}
##---
# Source: {{ $path }}
{{- range $key, $image := $images }}
docker tag {{ $from }}/{{ $image.Name }} {{ $image.Host }}/{{ $image.Name }}
{{- end }}
{{- end }}
`

var saveTemplate = template.Must(template.New("").Parse(saveScript))

func (i *Templates) GenerateSaveScript(log *logger.Logger) error {
	data := make(map[string]interface{})
	data["tpls"] = i
	return saveTemplate.Execute(log.GetOutput(), data)
}
