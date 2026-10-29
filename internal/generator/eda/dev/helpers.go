// Package dev generates semiconductor-device problems: carrier physics, junctions, MOSFET
// operation, short-channel and leakage behavior, advanced device architectures, variability,
// reliability, scaling, compact models and the agent-first future of device modeling.
package dev

import (
	"math"
	"regexp"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Physical constants, in the units the stems state.
const (
	kB    = 8.617e-5          // Boltzmann constant, eV/K
	qC    = 1.602e-19         // elementary charge, C
	eps0  = 8.854e-14         // vacuum permittivity, F/cm
	epsSi = 11.7 * eps0       // silicon permittivity, F/cm
	epsOx = 3.9 * eps0        // SiO2 permittivity, F/cm
	vT300 = 0.02585           // thermal voltage at 300 K, V
	ln10  = 2.302585092994046 // natural log of 10
)

// kpOf returns the requested knowledge point, or a uniform one when the request is unguided.
func kpOf(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

// sigPlaces is the number of decimal places that shows x to sig significant figures.
func sigPlaces(x float64, sig int) int {
	if x == 0 {
		return sig - 1
	}
	p := sig - 1 - int(math.Floor(math.Log10(math.Abs(x))))
	if p < 0 {
		p = 0
	}
	return p
}

// roundText is the rounding instruction for a stem.
func roundText(places int) string {
	switch places {
	case 0:
		return "Round to the nearest whole number."
	case 1:
		return "Round to 1 decimal place."
	default:
		return edakit.Sprintf("Round to %d decimal places.", places)
	}
}

// sigAnswer rounds x to sig significant figures and returns the answer and its rounding text.
func sigAnswer(x float64, sig int) (string, string) {
	p := sigPlaces(x, sig)
	return edakit.Num(x, p), roundText(p)
}

// sci renders m x 10^e in LaTeX (without surrounding dollars).
func sci(m float64, e int) string {
	if e == 0 {
		return edakit.Num(m, 2)
	}
	if m == 1 {
		return edakit.Sprintf("10^{%d}", e)
	}
	return edakit.Sprintf("%s \\times 10^{%d}", edakit.Num(m, 2), e)
}

// sciOf renders a positive value in scientific LaTeX with 3 significant figures.
func sciOf(x float64) string {
	e := int(math.Floor(math.Log10(x)))
	m := x / math.Pow(10, float64(e))
	m = edakit.Round(m, 2)
	if m >= 10 {
		m /= 10
		e++
	}
	return sci(m, e)
}

// mantissa splits a positive value into m x 10^e with 1 <= m < 10.
func mantissa(x float64) (float64, int) {
	e := int(math.Floor(math.Log10(x)))
	return x / math.Pow(10, float64(e)), e
}

// pickF returns a uniform element of the given values.
func pickF(xs ...float64) float64 { return edakit.Pick(xs) }

// fmt3 formats an option value with a fixed number of places, keeping trailing zeros so that
// computed options line up in the option list.
func fixed(x float64, places int) string {
	s := edakit.Sprintf("%.*f", places, x)
	if s == "-0" || strings.HasPrefix(s, "-0.") && strings.Trim(s[1:], "0.") == "" {
		s = s[1:]
	}
	return s
}

// currentUnit expresses a current in amperes in the unit that puts it between 1 and 1000.
func currentUnit(i float64) (float64, string) {
	units := []struct {
		scale float64
		name  string
	}{{1, "A"}, {1e-3, "mA"}, {1e-6, "µA"}, {1e-9, "nA"}, {1e-12, "pA"}, {1e-15, "fA"}}
	for _, u := range units {
		if math.Abs(i) >= u.scale {
			return i / u.scale, u.name
		}
	}
	return i / 1e-15, "fA"
}

// mix builds the common three-part mix: authored bank on KP 1, computed drill on KP 2, bank on KP 3.
func mix(items []edakit.Item, computed edakit.Func) generator.Generator {
	b := &edakit.Bank{Items: items}
	return &edakit.Mix{Parts: []generator.Generator{b, fn(computed)}}
}

var goExp = regexp.MustCompile(`(\d(?:\.\d+)?)e([+-])0*(\d+)`)

// fn wraps a computational generator so that Go's e-notation in its explanation (written inside
// math spans by %g verbs) renders as LaTeX scientific notation.
func fn(f edakit.Func) edakit.Func {
	return func(ctx generator.GeneratorContext) generator.Problem {
		p := f(ctx)
		p.Explanation = goExp.ReplaceAllStringFunc(p.Explanation, func(m string) string {
			g := goExp.FindStringSubmatch(m)
			sign := ""
			if g[2] == "-" {
				sign = "-"
			}
			return g[1] + " \\times 10^{" + sign + g[3] + "}"
		})
		return p
	}
}
