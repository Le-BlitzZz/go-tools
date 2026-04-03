package nilness

import (
	"testing"

	"github.com/Le-BlitzZz/tools/go/analysis/analysistest"
)

func TestNilness(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analysis, "example.com/Nilness")
}
