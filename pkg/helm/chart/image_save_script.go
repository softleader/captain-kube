package chart

import (
	"io"
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

func (i *Templates) GenerateSaveScript(out io.Writer) error {
	data := make(map[string]interface{})
	data["tpls"] = i
	return saveTemplate.Execute(out, data)
}
