{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "tke-monitor-api.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "tke-monitor-api.fullname" -}}
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

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "tke-monitor-api.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "tke-monitor-api.labels" -}}
helm.sh/chart: {{ include "tke-monitor-api.chart" . }}
{{ include "tke-monitor-api.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "tke-monitor-api.selectorLabels" -}}
app.kubernetes.io/name: {{ include "tke-monitor-api.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "tke-monitor-api.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "tke-monitor-api.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Generate certificates for server
*/}}
{{- define "tke-monitor-api.gen-certs" -}}
{{- $altNames := list (include "tke-monitor-api.fullname" .) ( printf "%s.%s" (include "tke-monitor-api.fullname" .) .Release.Namespace ) ( printf "%s.%s.svc" (include "tke-monitor-api.fullname" .) .Release.Namespace ) -}}
{{- $ca := genCA "tke-monitor-api-ca" 365 -}}
{{- $cert := genSignedCert ( include "tke-monitor-api.name" . ) nil $altNames 365 $ca -}}
ca.crt: {{ $ca.Cert | b64enc }}
ca.key: {{ $ca.Key | b64enc }}
server.crt: {{ $cert.Cert | b64enc }}
server.key: {{ $cert.Key | b64enc }}
{{- end -}}