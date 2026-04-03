package qf1006

import (
	"go/ast"
	"go/token"

	"github.com/Le-BlitzZz/go-tools/analysis/code"
	"github.com/Le-BlitzZz/go-tools/analysis/edit"
	"github.com/Le-BlitzZz/go-tools/analysis/lint"
	"github.com/Le-BlitzZz/go-tools/analysis/report"
	"github.com/Le-BlitzZz/go-tools/go/ast/astutil"
	"github.com/Le-BlitzZz/go-tools/pattern"

	"github.com/Le-BlitzZz/tools/go/analysis"
)

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: &analysis.Analyzer{
		Name:     "QF1006",
		Run:      run,
		Requires: code.RequiredAnalyzers,
	},
	Doc: &lint.RawDocumentation{
		Title: `Lift \'if\'+\'break\' into loop condition`,
		Before: `
for {
    if done {
        break
    }
    ...
}`,

		After: `
for !done {
    ...
}`,
		Since:    "2021.1",
		Severity: lint.SeverityHint,
	},
})

var Analyzer = SCAnalyzer.Analyzer

var checkForLoopIfBreak = pattern.MustParse(`(ForStmt nil nil nil if@(IfStmt nil cond (BranchStmt "BREAK" nil) nil):_)`)

func run(pass *analysis.Pass) (any, error) {
	for node, m := range code.Matches(pass, checkForLoopIfBreak) {
		pos := node.Pos() + token.Pos(len("for"))
		r := astutil.NegateDeMorgan(m.State["cond"].(ast.Expr), false)

		// FIXME(dh): we're leaving behind an empty line when we
		// delete the old if statement. However, we can't just delete
		// an additional character, in case there closing curly brace
		// is followed by a comment, or Windows newlines.
		report.Report(pass, m.State["if"].(ast.Node), "could lift into loop condition",
			report.Fixes(edit.Fix("Lift into loop condition",
				edit.ReplaceWithString(edit.Range{pos, pos}, " "+report.Render(pass, r)),
				edit.Delete(m.State["if"].(ast.Node)))))
	}
	return nil, nil
}
