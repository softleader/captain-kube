package chart

import (
	"bytes"
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

func (i *Templates) GenerateSaveScript() (buf bytes.Buffer, err error) {
	data := make(map[string]interface{})
	data["tpls"] = i
	err = saveTemplate.Execute(&buf, data)
	return
}
