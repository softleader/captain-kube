package chart

import (
	"github.com/softleader/captain-kube/pkg/logger"
	"text/template"
)

const saveScript = `
{{- range $path, $images := index . "images" }}
##---
# Source: {{ $path }}
{{- range $key, $image := $images }}
docker save -o ./{{ $image.Name }}.tar {{ $image.String }}
{{- end }}
{{- end }}
`

var saveTemplate = template.Must(template.New("").Parse(saveScript))

func (i *Templates) GenerateSaveScript(log *logger.Logger) error {
	data := make(map[string]interface{})
	data["tpls"] = i
	return saveTemplate.Execute(log.GetOutput(), data)
}
