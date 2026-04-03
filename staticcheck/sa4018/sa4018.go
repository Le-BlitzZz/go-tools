package sa4018

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"

	"github.com/Le-BlitzZz/go-tools/analysis/code"
	"github.com/Le-BlitzZz/go-tools/analysis/facts/generated"
	"github.com/Le-BlitzZz/go-tools/analysis/facts/purity"
	"github.com/Le-BlitzZz/go-tools/analysis/lint"
	"github.com/Le-BlitzZz/go-tools/analysis/report"

	"github.com/Le-BlitzZz/tools/go/analysis"
	"github.com/Le-BlitzZz/tools/go/analysis/passes/inspect"
)

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: &analysis.Analyzer{
		Name:     "SA4018",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer, generated.Analyzer, purity.Analyzer},
	},
	Doc: &lint.RawDocumentation{
		Title:    `Self-assignment of variables`,
		Since:    "2017.1",
		Severity: lint.SeverityWarning,
		MergeIf:  lint.MergeIfAny,
	},
})

var Analyzer = SCAnalyzer.Analyzer

func run(pass *analysis.Pass) (any, error) {
	pure := pass.ResultOf[purity.Analyzer].(purity.Result)

	fn := func(node ast.Node) {
		assign := node.(*ast.AssignStmt)
		if assign.Tok != token.ASSIGN || len(assign.Lhs) != len(assign.Rhs) {
			return
		}
		for i, lhs := range assign.Lhs {
			rhs := assign.Rhs[i]
			if reflect.TypeOf(lhs) != reflect.TypeOf(rhs) {
				continue
			}
			if code.MayHaveSideEffects(pass, lhs, pure) || code.MayHaveSideEffects(pass, rhs, pure) {
				continue
			}

			rlh := report.Render(pass, lhs)
			rrh := report.Render(pass, rhs)
			if rlh == rrh {
				report.Report(pass, assign, fmt.Sprintf("self-assignment of %s to %s", rrh, rlh), report.FilterGenerated())
			}
		}
	}
	code.Preorder(pass, fn, (*ast.AssignStmt)(nil))
	return nil, nil
}
