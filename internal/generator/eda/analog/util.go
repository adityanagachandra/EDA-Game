package analog

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Physical constants used across the package.
const (
	kBoltzmann = 1.380649e-23 // J/K
	eps0       = 8.854e-12    // F/m
)

// whichKP returns the requested knowledge point, or a uniform one when the request is unguided.
func whichKP(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

func num(q, a, why string, facts map[string]string) generator.Problem {
	return edakit.Numeric(q, a, why, facts)
}

// f formats a value for a stem or explanation with the given places, trailing zeros trimmed.
func f(x float64, places int) string { return edakit.Num(x, places) }

// pick returns a uniform element.
func pick[T any](xs ...T) T { return edakit.Pick(xs) }

// unif draws a value on [lo, hi] quantized to step.
func unif(lo, hi, step float64) float64 {
	n := int(math.Round((hi - lo) / step))
	return lo + step*float64(edakit.Int(0, n))
}

func db20(x float64) float64   { return 20 * math.Log10(x) }
func db10(x float64) float64   { return 10 * math.Log10(x) }
func undb10(x float64) float64 { return math.Pow(10, x/10) }
func undb20(x float64) float64 { return math.Pow(10, x/20) }

// phi is the standard normal cumulative distribution function.
func phi(x float64) float64 { return 0.5 * (1 + math.Erf(x/math.Sqrt2)) }

// pdf is the standard normal density.
func pdf(x float64) float64 { return math.Exp(-x*x/2) / math.Sqrt(2*math.Pi) }

func deg(rad float64) float64 { return rad * 180 / math.Pi }

func par(a, b float64) float64 { return a * b / (a + b) }

// choiceNums builds a computed multiple-choice problem whose options are formatted numbers with
// a unit suffix. Distractors that format equal to the answer or to each other are dropped and
// replaced by small perturbations so that at least three distractors remain.
func choiceNums(q string, ans float64, places int, unit, why string, wrong ...float64) generator.Problem {
	a := f(ans, places) + unit
	seen := map[string]bool{a: true}
	var ds []string
	add := func(x float64) {
		s := f(x, places) + unit
		if !seen[s] {
			seen[s] = true
			ds = append(ds, s)
		}
	}
	for _, w := range wrong {
		add(w)
	}
	for k := 1; len(ds) < 3 && k < 20; k++ {
		add(ans * (1 + 0.37*float64(k)))
	}
	return edakit.Choice(q, a, why, ds...)
}

// sci renders x for use inside $...$: a plain decimal for moderate magnitudes, otherwise
// mantissa times a power of ten.
func sci(x float64) string {
	ax := math.Abs(x)
	if ax == 0 || (ax >= 1e-2 && ax < 1e5) {
		return edakit.Sprintf("%.4g", x)
	}
	e := int(math.Floor(math.Log10(ax)))
	m := x / math.Pow(10, float64(e))
	return edakit.Sprintf("%s \\times 10^{%d}", f(m, 3), e)
}

// plural returns word with an s unless n is 1.
func plural(n int, word string) string {
	if n == 1 {
		return word
	}
	return word + "s"
}
