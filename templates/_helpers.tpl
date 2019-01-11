{{- define "captain-kube.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace " " ""  | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "captain-kube.captain" -}}
{{- printf "captain" -}}
{{- end -}}

{{- define "captain-kube.caplet" -}}
{{- printf "caplet" -}}
{{- end -}}

{{- define "captain-kube.config" -}}
{{- printf "captain-kube-config" -}}
{{- end -}}

{{- define "captain-kube.capui" -}}
{{- printf "capui" -}}
{{- end -}}