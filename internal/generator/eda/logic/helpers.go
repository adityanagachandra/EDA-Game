// Package logic generates digital-logic-design problems: number representations, combinational
// and sequential building blocks, finite-state machines, arithmetic circuits, memories,
// microarchitecture and the agent-first frontier of RTL design.
package logic

import (
	"math"
	"math/bits"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds every digital_logic_design generator to reg.
func Register(reg *generator.Registry) {
	registerNumber(reg)
	registerComb(reg)
	registerSeq(reg)
	registerFSM(reg)
	registerArith(reg)
	registerMem(reg)
	registerArch(reg)
	registerFrontier(reg)
}

// variant returns the knowledge point a request asks for (1, 2 or 3), or a uniform pick for an
// unguided request.
func variant(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

// bin writes v as an n-bit binary string, most significant bit first.
func bin(v uint64, n int) string {
	var b strings.Builder
	for i := n - 1; i >= 0; i-- {
		if v>>uint(i)&1 == 1 {
			b.WriteByte('1')
		} else {
			b.WriteByte('0')
		}
	}
	return b.String()
}

func popcount(v uint64) int { return bits.OnesCount64(v) }

// clog2 is the ceiling of log base two of n, for n >= 1.
func clog2(n int) int {
	k := 0
	for (1 << uint(k)) < n {
		k++
	}
	return k
}

// signed decodes the low n bits of v as a two's-complement integer.
func signed(v uint64, n int) int64 {
	v &= (1 << uint(n)) - 1
	if v>>uint(n-1)&1 == 1 {
		return int64(v) - (1 << uint(n))
	}
	return int64(v)
}

// wrap reduces x into the n-bit two's-complement range.
func wrap(x int64, n int) int64 {
	m := int64(1) << uint(n)
	r := ((x % m) + m) % m
	if r >= m/2 {
		r -= m
	}
	return r
}

func i64(x int64) string { return edakit.Sprintf("%d", x) }

// joinInts formats a list of integers with separators for a stem.
func joinInts(xs []int, sep string) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = edakit.IntStr(x)
	}
	return strings.Join(parts, sep)
}

func bit(b bool) int {
	if b {
		return 1
	}
	return 0
}

// isTie reports whether x sits exactly halfway between two values at the given number of decimal
// places, where a learner rounding half up and the formatter rounding half to even would disagree.
func isTie(x float64, places int) bool {
	s := math.Abs(x) * math.Pow(10, float64(places))
	frac := s - math.Floor(s)
	return math.Abs(frac-0.5) < 1e-7
}

// plural attaches a unit word to a count, adding "s" unless the count is one.
func plural(n int, word string) string {
	if n == 1 {
		return edakit.Sprintf("%d %s", n, word)
	}
	return edakit.Sprintf("%d %ss", n, word)
}

var roundRe = regexp.MustCompile(`Round to (\d+) decimal place`)

// withExact records the unrounded answer so retryTies can reject draws that land on a rounding tie.
func withExact(p generator.Problem, x float64) generator.Problem {
	if p.Facts == nil {
		p.Facts = map[string]string{}
	}
	p.Facts["exact"] = strconv.FormatFloat(x, 'g', 17, 64)
	return p
}

// retryTies redraws a problem whose rounded answer would sit exactly on a tie.
func retryTies(f edakit.Func) edakit.Func {
	return func(ctx generator.GeneratorContext) generator.Problem {
		var p generator.Problem
		for i := 0; i < 50; i++ {
			p = f(ctx)
			m := roundRe.FindStringSubmatch(p.Question)
			ex, ok := p.Facts["exact"]
			if m == nil || !ok {
				return p
			}
			x, err := strconv.ParseFloat(ex, 64)
			places, _ := strconv.Atoi(m[1])
			if err != nil || !isTie(x, places) {
				return p
			}
		}
		return p
	}
}
