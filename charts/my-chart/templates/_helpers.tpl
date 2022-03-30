{{- define "applyVersionOverrides" -}}
{{- $overrides := dict -}}
{{- range $override := .Values.versionOverrides -}}
{{- if semverCompare $override.constraint $.Capabilities.KubeVersion.Version -}}
{{- $_ := mergeOverwrite $overrides $override.values -}}
{{- end -}}
{{- end -}}
{{- $_ := mergeOverwrite .Values $overrides -}}
{{- end -}}
