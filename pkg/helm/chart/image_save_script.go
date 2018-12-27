package chart

import (
	"bytes"
	"strings"
	"text/template"
)

const saveScript = `
{{- $tpls := index . "tpls" -}}
{{- $len := len $tpls -}} 
{{- if eq $len 0 -}}
# no sources found in template
{{- else -}}
{{- range $path, $images := $tpls -}}
##---
# Source: {{ $path }}
{{- $len = len $images -}} 
{{- if eq $len 0 -}}
# no images found in source
{{- else -}}
{{- range $key, $image := $images }}
docker save -o ./{{ replace $image.Name ":" "_" -1 }}.tar {{ $image.String }}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
`

var saveFuncMap = template.FuncMap{
	"replace": strings.Replace,
}

var saveTemplate = template.Must(template.New("").Funcs(saveFuncMap).Parse(saveScript))

func (t *Templates) GenerateSaveScript() ([]byte, error) {
	data := make(map[string]interface{})
	data["tpls"] = t
	var buf bytes.Buffer
	if err := saveTemplate.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
