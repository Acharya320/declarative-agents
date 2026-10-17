{{/* Copyright (c) 2026 Nokia -- SPDX-License-Identifier: BSD-3-Clause */}}
{{- define "agent-operator.labels" -}}
app.kubernetes.io/name: agent-operator
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}
{{- define "agent-operator.image" -}}
{{- required "image.tag is required: the platform install sets it to the revision it built" .Values.image.tag | printf "%s:%s" .Values.image.repository -}}
{{- end -}}
