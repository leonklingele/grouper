package analyzer

import (
	"fmt"
	"go/ast"

	"github.com/leonklingele/grouper/pkg/analyzer/imports"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const (
	Name = "grouper"
	Doc  = `expression group analyzer: require 'import', 'const' and/or 'var' declaration groups`
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{ //nolint:exhaustivestruct // we do not need all fields
		Name:     Name,
		Doc:      Doc,
		Flags:    Flags(),
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func run(p *analysis.Pass) (interface{}, error) {
	flagLookupBool := func(name string) bool {
		return p.Analyzer.Flags.Lookup(name).Value.String() == "true"
	}

	c := &Config{
		ImportsConfig: &imports.Config{
			RequireSingleImport: flagLookupBool(FlagNameImportRequireSingleImport),
			RequireGrouping:     flagLookupBool(FlagNameImportRequireGrouping),
		},
	}

	return nil, pass(c, p)
}

func pass(c *Config, p *analysis.Pass) error {
	for _, f := range p.Files {
		if err := filepass(c, p, f); err != nil {
			return err
		}
	}

	return nil
}

func filepass(c *Config, p *analysis.Pass, f *ast.File) error {
	if err := imports.Filepass(c.ImportsConfig, p, f); err != nil {
		return fmt.Errorf("failed to imports.Filepass: %w", err)
	}

	return nil
}
