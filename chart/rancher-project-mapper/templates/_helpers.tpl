{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "rancher-project-mapper.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "rancher-project-mapper.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "rancher-project-mapper.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "rancher-project-mapper.gen-certs" -}}
{{- $altNames := list ( printf "%s.%s" (include "rancher-project-mapper.fullname" .) .Release.Namespace ) ( printf "%s.%s.svc" (include "rancher-project-mapper.fullname" .) .Release.Namespace ) -}}
{{- $ca := genCA "rancher-project-mapper-ca" 365 -}}
{{- $cert := genSignedCert ( include "rancher-project-mapper.name" . ) nil $altNames 365 $ca -}}
apiVersion: v1
kind: Secret
type: kubernetes.io/tls
metadata:
  name: {{ include "rancher-project-mapper.fullname" . }}
  labels:
    app: {{ include "rancher-project-mapper.fullname" . }}
    chart: {{ include "rancher-project-mapper.chart" . }}
    heritage: {{ .Release.Service }}
    release: {{ .Release.Name }}
  annotations:
    "helm.sh/hook": "pre-install"
    "helm.sh/hook-delete-policy": "before-hook-creation"
data:
    tls.crt: {{ $cert.Cert | b64enc }}
    tls.key: {{ $cert.Key | b64enc }}
    ca.crt: {{ $ca.Cert | b64enc }}
---
apiVersion: admissionregistration.k8s.io/v1beta1
kind: MutatingWebhookConfiguration
metadata:
  name: "projectmapper.rancher.1eb100.net"
webhooks:
  - name: "projectmapper.rancher.1eb100.net"
    rules:
      - apiGroups: [""]
        apiVersions: ["v1"]
        operations: ["CREATE"]
        resources: ["namespaces"]
        scope: "Cluster"
    clientConfig:
      caBundle: {{ $ca.Cert | b64enc }}
      service:
        namespace: {{ .Release.Namespace }}
        name: {{ include "rancher-project-mapper.fullname" . }}
        path: "/namespace"
    sideEffects: None
    timeoutSeconds: 10
{{- end -}}