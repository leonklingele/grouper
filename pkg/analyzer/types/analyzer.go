package types

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// https://go.dev/ref/spec#Type_declarations

type Type struct {
	Decl    *ast.GenDecl
	IsGroup bool
}

func Filepass(c *Config, p *analysis.Pass, f *ast.File) error {
	var types []*Type
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if genDecl.Tok == token.TYPE {
			types = append(types, &Type{
				Decl:    genDecl,
				IsGroup: genDecl.Lparen != 0,
			})
		}
	}

	numTypes := len(types)
	if numTypes == 0 {
		// Bail out early
		return nil
	}

	if c.RequireSingleType && numTypes > 1 {
		msg := fmt.Sprintf("should only use a single global 'type' declaration, %d found", numTypes)
		firstdup := types[1]
		decl := firstdup.Decl

		p.Report(analysis.Diagnostic{ //nolint:exhaustivestruct // we do not need all fields
			Pos:     decl.Pos(),
			End:     decl.End(),
			Message: msg,
			Related: toRelated(types[1:]),
			// TODO(leon): Suggest fix
		})
	}

	if c.RequireGrouping {
		var ungroupedTypes []*Type
		for _, t := range types {
			if !t.IsGroup {
				ungroupedTypes = append(ungroupedTypes, t)
			}
		}

		if numUngroupedTypes := len(ungroupedTypes); numUngroupedTypes != 0 {
			msg := "should only use grouped global 'type' declarations"
			firstmatch := ungroupedTypes[0]
			decl := firstmatch.Decl

			report := analysis.Diagnostic{ //nolint:exhaustivestruct // we do not need all fields
				Pos:     decl.Pos(),
				End:     decl.End(),
				Message: msg,
				// TODO(leon): Suggest fix
			}

			if numUngroupedTypes > 1 {
				report.Related = toRelated(ungroupedTypes[1:])
			}

			p.Report(report)
		}
	}

	return nil
}

func toRelated(types []*Type) []analysis.RelatedInformation {
	related := make([]analysis.RelatedInformation, 0, len(types))
	for _, t := range types {
		decl := t.Decl

		related = append(related, analysis.RelatedInformation{
			Pos:     decl.Pos(),
			End:     decl.End(),
			Message: "found here",
		})
	}

	return related
}
