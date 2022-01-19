package imports

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// https://go.dev/ref/spec#Import_declarations

type Import struct {
	Decl    *ast.GenDecl
	IsGroup bool
}

func Filepass(c *Config, p *analysis.Pass, f *ast.File) error {
	var imports []*Import
	ast.Inspect(f, func(n ast.Node) bool {
		if decl, ok := n.(*ast.GenDecl); ok {
			if decl.Tok == token.IMPORT {
				imports = append(imports, &Import{
					Decl:    decl,
					IsGroup: decl.Lparen != 0,
				})
			}
		}

		return true
	})

	numImports := len(imports)
	if numImports == 0 {
		// Bail out early
		return nil
	}

	if c.RequireSingleImport && numImports > 1 {
		msg := fmt.Sprintf("should only use a single 'import' declaration, %d found", numImports)
		firstdup := imports[1]
		decl := firstdup.Decl

		p.Report(analysis.Diagnostic{ //nolint:exhaustivestruct // we do not need all fields
			Pos:     decl.Pos(),
			End:     decl.End(),
			Message: msg,
			Related: toRelated(imports[1:]),
			// TODO(leon): Suggest fix
		})
	}

	if c.RequireGrouping {
		var ungroupedImports []*Import
		for _, imp := range imports {
			if !imp.IsGroup {
				ungroupedImports = append(ungroupedImports, imp)
			}
		}

		if numUngroupedImports := len(ungroupedImports); numUngroupedImports != 0 {
			msg := "should only use grouped 'import' declarations"
			firstmatch := ungroupedImports[0]
			decl := firstmatch.Decl

			report := analysis.Diagnostic{ //nolint:exhaustivestruct // we do not need all fields
				Pos:     decl.Pos(),
				End:     decl.End(),
				Message: msg,
				// TODO(leon): Suggest fix
			}

			if numUngroupedImports > 1 {
				report.Related = toRelated(ungroupedImports[1:])
			}

			p.Report(report)
		}
	}

	return nil
}

func toRelated(imports []*Import) []analysis.RelatedInformation {
	related := make([]analysis.RelatedInformation, 0, len(imports))
	for _, imp := range imports {
		decl := imp.Decl

		related = append(related, analysis.RelatedInformation{
			Pos:     decl.Pos(),
			End:     decl.End(),
			Message: "found here",
		})
	}

	return related
}
