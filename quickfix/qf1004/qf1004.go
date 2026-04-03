package qf1004

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/Le-BlitzZz/go-tools/analysis/edit"
	"github.com/Le-BlitzZz/go-tools/analysis/lint"
	"github.com/Le-BlitzZz/go-tools/analysis/report"
	typeindexanalyzer "github.com/Le-BlitzZz/go-tools/internal/analysisinternal/typeindex"
	"github.com/Le-BlitzZz/go-tools/internal/typesinternal/typeindex"

	"github.com/Le-BlitzZz/tools/go/analysis"
)

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: &analysis.Analyzer{
		Name:     "QF1004",
		Run:      run,
		Requires: []*analysis.Analyzer{typeindexanalyzer.Analyzer},
	},
	Doc: &lint.RawDocumentation{
		Title:    `Use \'strings.ReplaceAll\' instead of \'strings.Replace\' with \'n == -1\'`,
		Since:    "2021.1",
		Severity: lint.SeverityHint,
	},
})

var Analyzer = SCAnalyzer.Analyzer

var fns = []struct {
	path        string
	name        string
	replacement string
}{
	{"strings", "Replace", "strings.ReplaceAll"},
	{"strings", "SplitN", "strings.Split"},
	{"strings", "SplitAfterN", "strings.SplitAfter"},
	{"bytes", "Replace", "bytes.ReplaceAll"},
	{"bytes", "SplitN", "bytes.Split"},
	{"bytes", "SplitAfterN", "bytes.SplitAfter"},
}

func run(pass *analysis.Pass) (any, error) {
	// XXX respect minimum Go version

	// FIXME(dh): create proper suggested fix for renamed import

	index := pass.ResultOf[typeindexanalyzer.Analyzer].(*typeindex.Index)
	for _, fn := range fns {
		for c := range index.Calls(index.Object(fn.path, fn.name)) {
			call := c.Node().(*ast.CallExpr)
			if op, ok := call.Args[len(call.Args)-1].(*ast.UnaryExpr); ok && op.Op == token.SUB {
				if lit, ok := op.X.(*ast.BasicLit); ok && lit.Value == "1" {
					report.Report(pass, call.Fun, fmt.Sprintf("could use %s instead", fn.replacement),
						report.Fixes(edit.Fix(fmt.Sprintf("Use %s instead", fn.replacement),
							edit.ReplaceWithString(call.Fun, fn.replacement),
							edit.Delete(op))))
				}
			}
		}
	}
	return nil, nil
}
