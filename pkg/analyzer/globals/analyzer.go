package globals

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

type Global struct {
	Decl    *ast.GenDecl
	IsGroup bool
}

func Filepass(
	p *analysis.Pass, f *ast.File,
	token token.Token, requireSingle, requireGrouping bool,
) error {
	var globals []*Global
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if genDecl.Tok == token {
			globals = append(globals, &Global{
				Decl:    genDecl,
				IsGroup: genDecl.Lparen != 0,
			})
		}
	}

	numGlobals := len(globals)
	if numGlobals == 0 {
		// Bail out early
		return nil
	}

	if requireSingle && numGlobals > 1 {
		msg := fmt.Sprintf("should only use a single global '%s' declaration, %d found", token.String(), numGlobals)
		firstdup := globals[1]
		decl := firstdup.Decl

		p.Report(analysis.Diagnostic{ //nolint:exhaustivestruct // we do not need all fields
			Pos:     decl.Pos(),
			End:     decl.End(),
			Message: msg,
			Related: toRelated(globals[1:]),
			// TODO(leon): Suggest fix
		})
	}

	if requireGrouping {
		var ungrouped []*Global
		for _, g := range globals {
			if !g.IsGroup {
				ungrouped = append(ungrouped, g)
			}
		}

		if numUngrouped := len(ungrouped); numUngrouped != 0 {
			msg := fmt.Sprintf("should only use grouped global '%s' declarations", token.String())
			firstmatch := ungrouped[0]
			decl := firstmatch.Decl

			report := analysis.Diagnostic{ //nolint:exhaustivestruct // we do not need all fields
				Pos:     decl.Pos(),
				End:     decl.End(),
				Message: msg,
				// TODO(leon): Suggest fix
			}

			if numUngrouped > 1 {
				report.Related = toRelated(ungrouped[1:])
			}

			p.Report(report)
		}
	}

	return nil
}

func toRelated(globals []*Global) []analysis.RelatedInformation {
	related := make([]analysis.RelatedInformation, 0, len(globals))
	for _, g := range globals {
		decl := g.Decl

		related = append(related, analysis.RelatedInformation{
			Pos:     decl.Pos(),
			End:     decl.End(),
			Message: "found here",
		})
	}

	return related
}
