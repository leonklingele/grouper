package consts

import (
	"go/ast"
	"go/token"

	"github.com/leonklingele/grouper/pkg/analyzer/globals"

	"golang.org/x/tools/go/analysis"
)

// https://go.dev/ref/spec#Constant_declarations

func Filepass(c *Config, p *analysis.Pass, f *ast.File) error {
	//nolint:wrapcheck // Wrapper func doesn't need to be wrapped
	return globals.Filepass(
		p, f,
		token.CONST, c.RequireSingleConst, c.RequireGrouping,
	)
}
