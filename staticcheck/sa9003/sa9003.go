package sa9003

import (
	"go/ast"

	"github.com/Le-BlitzZz/tools/go/analysis"
	"github.com/Le-BlitzZz/go-tools/analysis/facts/generated"
	"github.com/Le-BlitzZz/go-tools/analysis/lint"
	"github.com/Le-BlitzZz/go-tools/analysis/report"
	"github.com/Le-BlitzZz/go-tools/go/ir/irutil"
	"github.com/Le-BlitzZz/go-tools/internal/passes/buildir"
)

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: &analysis.Analyzer{
		Name:     "SA9003",
		Run:      run,
		Requires: []*analysis.Analyzer{buildir.Analyzer, generated.Analyzer},
	},
	Doc: &lint.RawDocumentation{
		Title:      `Empty body in an if or else branch`,
		Since:      "2017.1",
		NonDefault: true,
		Severity:   lint.SeverityWarning,
		MergeIf:    lint.MergeIfAny,
	},
})

var Analyzer = SCAnalyzer.Analyzer

func run(pass *analysis.Pass) (any, error) {
	for _, fn := range pass.ResultOf[buildir.Analyzer].(*buildir.IR).SrcFuncs {
		if fn.Source() == nil {
			continue
		}
		if irutil.IsExample(fn) {
			continue
		}
		cb := func(node ast.Node) bool {
			ifstmt, ok := node.(*ast.IfStmt)
			if !ok {
				return true
			}
			if ifstmt.Else != nil {
				b, ok := ifstmt.Else.(*ast.BlockStmt)
				if !ok || len(b.List) != 0 {
					return true
				}
				report.Report(pass, ifstmt.Else, "empty branch", report.FilterGenerated(), report.ShortRange())
			}
			if len(ifstmt.Body.List) != 0 {
				return true
			}
			report.Report(pass, ifstmt, "empty branch", report.FilterGenerated(), report.ShortRange())
			return true
		}
		if source := fn.Source(); source != nil {
			ast.Inspect(source, cb)
		}
	}
	return nil, nil
}
