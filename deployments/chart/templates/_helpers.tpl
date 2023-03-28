{{- define "prometheus-multi-tenant-proxy-server.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "prometheus-multi-tenant-proxy-server.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{- define "prometheus-multi-tenant-proxy-server.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "prometheus-multi-tenant-proxy-server.selectorLabels" -}}
app.kubernetes.io/name: {{ include "prometheus-multi-tenant-proxy-server.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "prometheus-multi-tenant-proxy-server.labels" -}}
helm.sh/chart: {{ include "prometheus-multi-tenant-proxy-server.chart" . }}
{{ include "prometheus-multi-tenant-proxy-server.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}
