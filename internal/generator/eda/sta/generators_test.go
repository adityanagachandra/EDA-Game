package sta

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func TestDomainContract(t *testing.T) {
	edakit.CheckDomain(t, Register, "timing_analysis")
}
