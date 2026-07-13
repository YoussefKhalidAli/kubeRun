{{/*
Validate that target namespace matches values.namespace.
*/}}
{{- define "kuberun.validateNamespace" -}}
{{- if and .Release.Namespace .Values.namespace (ne .Release.Namespace .Values.namespace) -}}
{{- fail (printf "Namespace mismatch: target namespace (%s) must match .Values.namespace (%s). Please set both to the same value, e.g. --namespace %s --set namespace=%s" .Release.Namespace .Values.namespace .Values.namespace .Values.namespace) -}}
{{- end -}}
{{- end -}}

{{/*
Controller ServiceAccount Name
*/}}
{{- define "kuberun.controllerServiceAccountName" -}}
{{- if .Values.controller.serviceAccount.create -}}
    {{- default "kuberun-controller-sa" .Values.controller.serviceAccount.name -}}
{{- else -}}
    {{- default "default" .Values.controller.serviceAccount.name -}}
{{- end -}}
{{- end -}}
