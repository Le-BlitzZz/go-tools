package deprecated

import (
	"testing"

	"github.com/Le-BlitzZz/tools/go/analysis/analysistest"
)

func TestDeprecated(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "example.com/Deprecated")
}
