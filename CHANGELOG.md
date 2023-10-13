<a name="unreleased"></a>
## [Unreleased]


<a name="v0.6.0"></a>
## [v0.6.0] - 2023-10-13
### Code Refactoring

- Refactor common graph approach both for Pipeline/PipelineRun
- Remove unused variables from function signature

### Testing

- Hide output in tests

### Routine

- We don't have plugins

### Documentation

- How to install with Homebrew


<a name="v0.5.0"></a>
## [v0.5.0] - 2023-10-12
### Code Refactoring

- Move common graph logic to separate package

### Testing

- Add fake kubeconfig

### Routine

- Add brew package delivery
- Update .gorelease.yaml fromatting
- Update changelog


<a name="v0.4.0"></a>
## [v0.4.0] - 2023-10-10
### Features

- Allow to get graph for specific Pipeline/PipelineRun provided by name

### Code Refactoring

- Do some minor sruff
- Use go-template for DOT format
- Use go-template for PlantUML
- Use go-template for Mermaid generation
- Use isRoot to define root task
- Move create logic to separate function


<a name="v0.3.0"></a>
## [v0.3.0] - 2023-10-08
### Features

- Add start node for PUML format
- Add start node for DOT format
- Add support for start node using Mermaid format

### Testing

- Add more tests

### Routine

- Remove deprecated code

### Documentation

- Update mmd example graph
- Update mmd format
- Update documentation with start/stop nodes support


<a name="v0.2.0"></a>
## [v0.2.0] - 2023-10-07
### Code Refactoring

- Add support for writing graphs to file
- Add output format validator
- Align code base to tektoncd/cli repo approach
- First refactor iteration

### Routine

- Add code coverage report for SonarCloud
- Add sonarqube check
- Add CHANGELEG.md file and approach to build it

### Documentation

- Update documentation
- Update documentation


<a name="v0.1.0"></a>
## v0.1.0 - 2023-10-03
### Features

- Add support for the TaskRef output
- Add Task Reference name to the graph

### Code Refactoring

- Refactor tests to get rid of flanky results
- Refactor build graph function
- Change project structure
- Refactor buildgraph function

### Routine

- Update gorelease options
- Add goreleaser
- Add go vet target to Makefile
- Add linting stuff
- Fix module
- Remove back-quote for mermaid format

### Documentation

- Add TOC for README file
- Update examples in documentation
- Update README.md file


[Unreleased]: https://github.com/sergk/tkn-graph/compare/v0.6.0...HEAD
[v0.6.0]: https://github.com/sergk/tkn-graph/compare/v0.5.0...v0.6.0
[v0.5.0]: https://github.com/sergk/tkn-graph/compare/v0.4.0...v0.5.0
[v0.4.0]: https://github.com/sergk/tkn-graph/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/sergk/tkn-graph/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/sergk/tkn-graph/compare/v0.1.0...v0.2.0
