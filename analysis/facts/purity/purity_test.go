package purity

import (
	"testing"

	"github.com/Le-BlitzZz/tools/go/analysis/analysistest"
)

func TestPurity(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "example.com/Purity")
}
