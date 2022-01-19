# grouper — a Go linter to analyze expression groups

## Installation

```sh
go get -u github.com/leonklingele/grouper/...
grouper -help
```

## Run analyzer

```sh
grouper -import-require-single-import -import-require-grouping ./...

# Example output:
GOPATH/src/github.com/leonklingele/grouper/pkg/analyzer/analyzer.go:8:1: should only use a single 'import' declaration, 2 found
GOPATH/src/github.com/leonklingele/grouper/pkg/analyzer/flags.go:3:1: should only use grouped 'import' declarations
```
