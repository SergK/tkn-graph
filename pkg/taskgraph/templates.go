package taskgraph

// mermaidTemplate is the template used to generate the mermaid graph
// The template is based on the mermaid flowchart syntax: https://mermaid-js.github.io/mermaid/#/flowchart
// The template uses the following variables:
//   - PipelineName: Name of the pipeline
//   - Nodes: Map of nodes in the graph
const mermaidTemplate = `---
title: {{ .PipelineName }}
---
flowchart TD
{{- range $name, $node := .Nodes }}
{{- if eq (len $node.Dependencies) 0 }}
   {{ $name }} --> stop([fa:fa-circle])
{{- end }}
{{- if $node.IsRoot }}
   start([fa:fa-circle]) --> {{ $name }}
{{- end }}
{{- range $dep := $node.Dependencies }}
   {{ $name }} --> {{ $dep.Name }}
{{- end }}
{{- end }}
`

// mermaidTemplateWithTaskRef is the template used to generate the mermaid graph with taskRefName
// The template is based on the mermaid flowchart syntax: https://mermaid-js.github.io/mermaid/#/flowchart
// taskRefName is added to the node name in the graph in format: (taskRefName). For example:
//
//	---------------
//	|  taskName    |
//	|(taskRefName) |
//	---------------
const mermaidTemplateWithTaskRef = `---
title: {{ .PipelineName }}
---
flowchart TD
{{- range $name, $node := .Nodes }}
{{- if eq (len $node.Dependencies) 0 }}
   {{ $name }}("{{ $node.Name }}
   ({{ $node.TaskRefName }})") --> stop([fa:fa-circle])
{{- end }}
{{- if $node.IsRoot }}
   start([fa:fa-circle]) --> {{ $name }}("{{ $node.Name }}
   ({{ $node.TaskRefName }})")
{{- end }}
{{- range $dep := $node.Dependencies }}
   {{ $name }}("{{ $node.Name }}
   ({{ $node.TaskRefName }})") --> {{ $dep.Name }}("{{ $dep.Name }}
   ({{ $dep.TaskRefName }})")
{{- end }}
{{- end }}
`
