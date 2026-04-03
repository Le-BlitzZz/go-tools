package directives

import (
	"reflect"

	"github.com/Le-BlitzZz/tools/go/analysis"
	"github.com/Le-BlitzZz/go-tools/analysis/lint"
)

func directives(pass *analysis.Pass) (any, error) {
	return lint.ParseDirectives(pass.Files, pass.Fset), nil
}

var Analyzer = &analysis.Analyzer{
	Name:             "directives",
	Doc:              "extracts linter directives",
	Run:              directives,
	RunDespiteErrors: true,
	ResultType:       reflect.TypeFor[[]lint.Directive](),
}
