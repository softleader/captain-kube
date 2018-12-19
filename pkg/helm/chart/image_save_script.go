package chart

import (
	"bytes"
	"text/template"
)

const saveScript = `
{{ $from := index . "from" }}
{{- range $source, $images := index . "images" }}
##---
# Source: {{ $source }}
{{- range $key, $image := $images }}
docker tag {{ $from }}/{{ $image.Name }} {{ $image.Host }}/{{ $image.Name }}
{{- end }}
{{- end }}
`

var saveTemplate = template.Must(template.New("").Parse(saveScript))

func (i *Images) GenerateSaveScript() (buf bytes.Buffer, err error) {
	data := make(map[string]interface{})
	data["images"] = i
	err = saveTemplate.Execute(&buf, data)
	return
}
