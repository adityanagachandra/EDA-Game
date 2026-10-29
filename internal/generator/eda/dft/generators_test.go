package dft

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func TestDomainContract(t *testing.T) {
	edakit.CheckDomain(t, Register, "test_and_dft")
}
