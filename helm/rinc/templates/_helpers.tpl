{{- define "rinc.chart" -}}
	{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "namespace" }}
	{{- default .Release.Namespace .Values.namespace }}
{{- end }}

{{- define "deployment.name" -}}
  {{- if .Values.web.fullnameOverride }}
    {{- .Values.web.fullnameOverride | trunc 63 | trimSuffix "-" }}
  {{- else if .Values.web.nameOverride }}
    {{- printf "%s-%s" .Chart.Name .Values.web.nameOverride | trunc 63 | trimSuffix "-" }}
  {{- else }}
    {{- printf "%s-web" .Chart.Name | trunc 63 | trimSuffix "-" }}
  {{- end }}
{{- end }}

{{- define "deployment.selectorLabels" -}}
app.kubernetes.io/name: {{ include "deployment.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "deployment.labels" -}}
helm.sh/chart: {{ include "rinc.chart" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{ if .Chart.AppVersion -}}
  app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
{{ include "deployment.selectorLabels" . }}
{{- end }}

{{- define "configMap.name" -}}
  {{- if .Values.config.configMap.fullnameOverride }}
    {{- .Values.config.configMap.fullnameOverride | trunc 63 | trimSuffix "-" }}
  {{- else if .Values.config.configMap.nameOverride }}
    {{- printf "%s-%s" .Chart.Name .Values.config.configMap.nameOverride | trunc 63 | trimSuffix "-" }}
  {{- else }}
    {{- printf "%s-config" .Chart.Name | trunc 63 | trimSuffix "-" }}
  {{- end }}
{{- end }}

{{- define "cronjob.name" -}}
  {{- if .Values.reportingCronJob.fullnameOverride }}
    {{- .Values.reportingCronJob.fullnameOverride | trunc 63 | trimSuffix "-" }}
  {{- else if .Values.reportingCronJob.nameOverride }}
    {{- printf "%s-%s" .Chart.Name .Values.reportingCronJob.nameOverride | trunc 63 | trimSuffix "-" }}
  {{- else }}
    {{- printf "%s-reporting-cronjob" .Chart.Name | trunc 63 | trimSuffix "-" }}
  {{- end }}
{{- end }}

{{- define "cronjob.labels" -}}
helm.sh/chart: {{ include "rinc.chart" . }}
app.kubernetes.io/name: {{ include "cronjob.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{ if .Chart.AppVersion -}}
  app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
{{- end }}

{{- define "serviceAccount.name" -}}
  {{- if .Values.rbac.serviceAccount.fullnameOverride }}
    {{- .Values.rbac.serviceAccount.fullnameOverride | trunc 63 | trimSuffix "-" }}
  {{- else if .Values.rbac.serviceAccount.nameOverride }}
    {{- printf "%s-%s" .Chart.Name .Values.rbac.serviceAccount.nameOverride | trunc 63 | trimSuffix "-" }}
  {{- else }}
    {{- .Chart.Name | trunc 63 | trimSuffix "-" }}
  {{- end }}
{{- end }}

{{- define "clusterRoleBinding.name" -}}
  {{- if .Values.rbac.clusterRoleBinding.fullnameOverride }}
    {{- .Values.rbac.clusterRoleBinding.fullnameOverride | trunc 63 | trimSuffix "-" }}
  {{- else if .Values.rbac.clusterRoleBinding.nameOverride }}
    {{- printf "%s-%s" .Chart.Name .Values.rbac.clusterRoleBinding.nameOverride | trunc 63 | trimSuffix "-" }}
  {{- else }}
    {{- .Chart.Name | trunc 63 | trimSuffix "-" }}
  {{- end }}
{{- end }}
