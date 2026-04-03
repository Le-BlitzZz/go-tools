package sa4028

import (
	"github.com/Le-BlitzZz/go-tools/analysis/code"
	"github.com/Le-BlitzZz/go-tools/analysis/lint"
	"github.com/Le-BlitzZz/go-tools/analysis/report"
	"github.com/Le-BlitzZz/go-tools/pattern"

	"github.com/Le-BlitzZz/tools/go/analysis"
)

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: &analysis.Analyzer{
		Name:     "SA4028",
		Run:      run,
		Requires: code.RequiredAnalyzers,
	},
	Doc: &lint.RawDocumentation{
		Title:    `\'x % 1\' is always zero`,
		Since:    "2022.1",
		Severity: lint.SeverityWarning,
		MergeIf:  lint.MergeIfAny, // MergeIfAny if we only flag literals, not named constants
	},
})

var Analyzer = SCAnalyzer.Analyzer

var moduloOneQ = pattern.MustParse(`(BinaryExpr _ "%" (IntegerLiteral "1"))`)

func run(pass *analysis.Pass) (any, error) {
	for node := range code.Matches(pass, moduloOneQ) {
		report.Report(pass, node, "x % 1 is always zero")
	}
	return nil, nil
}
