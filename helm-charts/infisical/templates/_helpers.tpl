{{/*
Expand the name of the chart.
*/}}
{{- define "gsoc2.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "gsoc2.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create unified labels for gsoc2 components
*/}}
{{- define "gsoc2.common.matchLabels" -}}
app: {{ template "gsoc2.name" . }}
release: {{ .Release.Name }}
{{- end -}}

{{- define "gsoc2.common.metaLabels" -}}
chart: {{ template "gsoc2.chart" . }}
heritage: {{ .Release.Service }}
{{- end -}}

{{- define "gsoc2.common.labels" -}}
{{ include "gsoc2.common.matchLabels" . }}
{{ include "gsoc2.common.metaLabels" . }}
{{- end -}}


{{- define "gsoc2.backend.labels" -}}
{{ include "gsoc2.backend.matchLabels" . }}
{{ include "gsoc2.common.metaLabels" . }}
{{- end -}}

{{- define "gsoc2.backend.matchLabels" -}}
component: {{ .Values.backend.name | quote }}
{{ include "gsoc2.common.matchLabels" . }}
{{- end -}}

{{- define "gsoc2.mongodb.labels" -}}
{{ include "gsoc2.mongodb.matchLabels" . }}
{{ include "gsoc2.common.metaLabels" . }}
{{- end -}}

{{- define "gsoc2.mongodb.matchLabels" -}}
component: {{ .Values.mongodb.name | quote }}
{{ include "gsoc2.common.matchLabels" . }}
{{- end -}}

{{/*
Create a fully qualified backend name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "gsoc2.backend.fullname" -}}
{{- if .Values.backend.fullnameOverride -}}
{{- .Values.backend.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- printf "%s-%s" .Release.Name .Values.backend.name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-%s" .Release.Name $name .Values.backend.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}


{{/*
Create a fully qualified mongodb name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "gsoc2.mongodb.fullname" -}}
{{- if .Values.mongodb.fullnameOverride -}}
{{- .Values.mongodb.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- printf "%s-%s" .Release.Name .Values.mongodb.name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-%s" .Release.Name $name .Values.mongodb.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create the mongodb connection string.
*/}}
{{- define "gsoc2.mongodb.connectionString" -}}
{{- $host := include "gsoc2.mongodb.fullname" . -}}
{{- $port := 27017 -}}
{{- $user := first .Values.mongodb.auth.usernames | default "root" -}}
{{- $pass := first .Values.mongodb.auth.passwords | default "root" -}}
{{- $database := first .Values.mongodb.auth.databases | default "test" -}}
{{- $connectionString := printf "mongodb://%s:%s@%s:%d/%s" $user $pass $host $port $database -}}
{{/* Backward compatibility (< 0.1.16, deprecated) */}}
{{- if .Values.mongodbConnection.externalMongoDBConnectionString -}}
{{- $connectionString = .Values.mongodbConnection.externalMongoDBConnectionString -}}
{{- end -}}
{{- printf "%s" $connectionString -}}
{{- end -}}