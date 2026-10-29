package cts

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func TestDomainContract(t *testing.T) {
	edakit.CheckDomain(t, Register, "clock_and_power")
}
