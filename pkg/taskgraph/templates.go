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

// dotTemplate is the template used to generate the DOT graph
// The template is based on the DOT language: https://graphviz.org/doc/info/lang.html
// We replace "-" with "_" in the node names to avoid issues with the DOT language
const plantumlTemplate = `@startuml
hide empty description
title {{ .PipelineName }}
{{ range $name, $node := .Nodes }}
{{- $trName := replace $name "-" "_" }}
{{- if eq (len $node.Dependencies) 0 }}
   {{ $trName }} --> [*]
{{- end }}
{{- if $node.IsRoot }}
   [*] --> {{ $trName }}
{{- end }}
{{- range $dep := $node.Dependencies }}
   {{- $trDepName := replace $dep.Name "-" "_" }}
   {{ $trName }} -down-> {{ $trDepName }}
{{- end }}
{{ end }}
@enduml
`

// dotTemplate is the template used to generate the DOT graph
// The template is based on the DOT language: https://graphviz.org/doc/info/lang.html
// We replace "-" with "_" in the node names to avoid issues with the DOT language
const plantumlTemplateWithTaskRef = `@startuml
hide empty description
title {{ .PipelineName }}
{{ range $name, $node := .Nodes }}
{{- $trName := replace $name "-" "_" }}
   {{ $trName }}: {{ $node.TaskRefName }}
{{- if eq (len $node.Dependencies) 0 }}
   {{ $trName }} --> [*]
{{- end }}
{{- if $node.IsRoot }}
   [*] --> {{ $trName }}
{{- end }}
{{- range $dep := $node.Dependencies }}
   {{- $trDepName := replace $dep.Name "-" "_" }}
   {{ $trName }} -down-> {{ $trDepName }}
{{- end }}
{{ end }}
@enduml
`
