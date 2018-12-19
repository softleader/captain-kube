package chart


import (
	"bytes"
	"text/template"
)

const loadScript = `
{{- range $source, $images := index . "images" }}
##---
# Source: {{ $source }}
{{- range $key, $image := $images }}
docker load -i ./{{ $image.Name }}.tar
{{- end }}
{{- end }}
`

var loadTemplate = template.Must(template.New("").Parse(loadScript))

func (i *Images) GenerateLoadScript() (buf bytes.Buffer, err error) {
	data := make(map[string]interface{})
	data["images"] = i
	err = loadTemplate.Execute(&buf, data)
	return
}
