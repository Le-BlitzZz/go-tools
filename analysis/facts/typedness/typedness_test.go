package typedness

import (
	"testing"

	"github.com/Le-BlitzZz/tools/go/analysis/analysistest"
)

func TestTypedness(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analysis, "example.com/Typedness")
}
