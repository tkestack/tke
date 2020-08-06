{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "tke-platform-controller.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "tke-platform-controller.fullname" -}}
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
{{- define "tke-platform-controller.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "tke-platform-controller.labels" -}}
helm.sh/chart: {{ include "tke-platform-controller.chart" . }}
{{ include "tke-platform-controller.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "tke-platform-controller.selectorLabels" -}}
app.kubernetes.io/name: {{ include "tke-platform-controller.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "tke-platform-controller.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "tke-platform-controller.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Generate certificates for server
*/}}
{{- define "tke-platform-controller.gen-certs" -}}
{{- $altNames := list (include "tke-platform-controller.fullname" .) ( printf "%s.%s" (include "tke-platform-controller.fullname" .) .Release.Namespace ) ( printf "%s.%s.svc" (include "tke-platform-controller.fullname" .) .Release.Namespace ) -}}
{{- $ca := genCA "tke-platform-controller-ca" 365 -}}
{{- $cert := genSignedCert ( include "tke-platform-controller.name" . ) nil $altNames 365 $ca -}}
ca.crt: {{ $ca.Cert | b64enc }}
ca.key: {{ $ca.Key | b64enc }}
server.crt: {{ $cert.Cert | b64enc }}
server.key: {{ $cert.Key | b64enc }}
{{- end -}}