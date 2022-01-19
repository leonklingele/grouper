package consts

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// https://go.dev/ref/spec#Constant_declarations

type Const struct {
	Decl    *ast.GenDecl
	IsGroup bool
}

func Filepass(c *Config, p *analysis.Pass, f *ast.File) error {
	var consts []*Const
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if genDecl.Tok == token.CONST {
			consts = append(consts, &Const{
				Decl:    genDecl,
				IsGroup: genDecl.Lparen != 0,
			})
		}
	}

	numConsts := len(consts)
	if numConsts == 0 {
		// Bail out early
		return nil
	}

	if c.RequireSingleConst && numConsts > 1 {
		msg := fmt.Sprintf("should only use a single global 'const' declaration, %d found", numConsts)
		firstdup := consts[1]
		decl := firstdup.Decl

		p.Report(analysis.Diagnostic{ //nolint:exhaustivestruct // we do not need all fields
			Pos:     decl.Pos(),
			End:     decl.End(),
			Message: msg,
			Related: toRelated(consts[1:]),
			// TODO(leon): Suggest fix
		})
	}

	if c.RequireGrouping {
		var ungroupedConsts []*Const
		for _, c := range consts {
			if !c.IsGroup {
				ungroupedConsts = append(ungroupedConsts, c)
			}
		}

		if numUngroupedConsts := len(ungroupedConsts); numUngroupedConsts != 0 {
			msg := "should only use grouped global 'const' declarations"
			firstmatch := ungroupedConsts[0]
			decl := firstmatch.Decl

			report := analysis.Diagnostic{ //nolint:exhaustivestruct // we do not need all fields
				Pos:     decl.Pos(),
				End:     decl.End(),
				Message: msg,
				// TODO(leon): Suggest fix
			}

			if numUngroupedConsts > 1 {
				report.Related = toRelated(ungroupedConsts[1:])
			}

			p.Report(report)
		}
	}

	return nil
}

func toRelated(consts []*Const) []analysis.RelatedInformation {
	related := make([]analysis.RelatedInformation, 0, len(consts))
	for _, c := range consts {
		decl := c.Decl

		related = append(related, analysis.RelatedInformation{
			Pos:     decl.Pos(),
			End:     decl.End(),
			Message: "found here",
		})
	}

	return related
}
