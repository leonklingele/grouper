package vars

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// https://go.dev/ref/spec#Variable_declarations

type Var struct {
	Decl    *ast.GenDecl
	IsGroup bool
}

func Filepass(c *Config, p *analysis.Pass, f *ast.File) error {
	var vars []*Var
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if genDecl.Tok == token.VAR {
			vars = append(vars, &Var{
				Decl:    genDecl,
				IsGroup: genDecl.Lparen != 0,
			})
		}
	}

	numVars := len(vars)
	if numVars == 0 {
		// Bail out early
		return nil
	}

	if c.RequireSingleVar && numVars > 1 {
		msg := fmt.Sprintf("should only use a single global 'var' declaration, %d found", numVars)
		firstdup := vars[1]
		decl := firstdup.Decl

		p.Report(analysis.Diagnostic{ //nolint:exhaustivestruct // we do not need all fields
			Pos:     decl.Pos(),
			End:     decl.End(),
			Message: msg,
			Related: toRelated(vars[1:]),
			// TODO(leon): Suggest fix
		})
	}

	if c.RequireGrouping {
		var ungroupedVars []*Var
		for _, v := range vars {
			if !v.IsGroup {
				ungroupedVars = append(ungroupedVars, v)
			}
		}

		if numUngroupedVars := len(ungroupedVars); numUngroupedVars != 0 {
			msg := "should only use grouped global 'var' declarations"
			firstmatch := ungroupedVars[0]
			decl := firstmatch.Decl

			report := analysis.Diagnostic{ //nolint:exhaustivestruct // we do not need all fields
				Pos:     decl.Pos(),
				End:     decl.End(),
				Message: msg,
				// TODO(leon): Suggest fix
			}

			if numUngroupedVars > 1 {
				report.Related = toRelated(ungroupedVars[1:])
			}

			p.Report(report)
		}
	}

	return nil
}

func toRelated(vars []*Var) []analysis.RelatedInformation {
	related := make([]analysis.RelatedInformation, 0, len(vars))
	for _, v := range vars {
		decl := v.Decl

		related = append(related, analysis.RelatedInformation{
			Pos:     decl.Pos(),
			End:     decl.End(),
			Message: "found here",
		})
	}

	return related
}
