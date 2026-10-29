package numer

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func TestDomainContract(t *testing.T) {
	edakit.CheckDomain(t, Register, "numerical_methods_for_eda")
}
