package sa4021

import (
	"github.com/Le-BlitzZz/go-tools/analysis/code"
	"github.com/Le-BlitzZz/go-tools/analysis/facts/generated"
	"github.com/Le-BlitzZz/go-tools/analysis/lint"
	"github.com/Le-BlitzZz/go-tools/analysis/report"
	"github.com/Le-BlitzZz/go-tools/pattern"

	"github.com/Le-BlitzZz/tools/go/analysis"
)

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: &analysis.Analyzer{
		Name:     "SA4021",
		Run:      run,
		Requires: append([]*analysis.Analyzer{generated.Analyzer}, code.RequiredAnalyzers...),
	},
	Doc: &lint.RawDocumentation{
		Title:    `\"x = append(y)\" is equivalent to \"x = y\"`,
		Since:    "2019.2",
		Severity: lint.SeverityWarning,
		MergeIf:  lint.MergeIfAny,
	},
})

var Analyzer = SCAnalyzer.Analyzer

var checkSingleArgAppendQ = pattern.MustParse(`(CallExpr (Builtin "append") [_])`)

func run(pass *analysis.Pass) (any, error) {
	for node := range code.Matches(pass, checkSingleArgAppendQ) {
		report.Report(pass, node, "x = append(y) is equivalent to x = y", report.FilterGenerated())
	}
	return nil, nil
}
