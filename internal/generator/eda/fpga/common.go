package fpga

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// whichKP resolves the knowledge point: the requested one, or a uniform pick when unguided.
func whichKP(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

func num(q, a, why string, facts map[string]string) generator.Problem {
	return edakit.Numeric(q, a, why, facts)
}

func ceilDiv(a, b int) int {
	if b <= 0 {
		return 0
	}
	return (a + b - 1) / b
}

func maxInt(xs ...int) int {
	m := xs[0]
	for _, x := range xs[1:] {
		if x > m {
			m = x
		}
	}
	return m
}

func minInt(xs ...int) int {
	m := xs[0]
	for _, x := range xs[1:] {
		if x < m {
			m = x
		}
	}
	return m
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func joinInts(xs []int) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = edakit.IntStr(x)
	}
	return strings.Join(parts, ", ")
}

func pow2(k int) int { return 1 << uint(k) }

// setText renders a sorted set of names as {a, b, c}.
func setText(xs []string) string { return "{" + strings.Join(xs, ", ") + "}" }

// union returns the sorted union of string sets, keeping first-seen order stable by sorting.
func union(sets ...[]string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range sets {
		for _, x := range s {
			if !seen[x] {
				seen[x] = true
				out = append(out, x)
			}
		}
	}
	sortStrings(out)
	return out
}

func sortStrings(xs []string) {
	for i := 1; i < len(xs); i++ {
		for j := i; j > 0 && less(xs[j], xs[j-1]); j-- {
			xs[j], xs[j-1] = xs[j-1], xs[j]
		}
	}
}

// less orders names with a numeric suffix naturally (n2 before n10).
func less(a, b string) bool {
	pa, na := split(a)
	pb, nb := split(b)
	if pa != pb {
		return pa < pb
	}
	return na < nb
}

func split(s string) (string, int) {
	i := len(s)
	for i > 0 && s[i-1] >= '0' && s[i-1] <= '9' {
		i--
	}
	n := 0
	for _, r := range s[i:] {
		n = n*10 + int(r-'0')
	}
	return s[:i], n
}

// pickDistinct returns n distinct elements of xs in random order.
func pickDistinct[T any](xs []T, n int) []T {
	s := edakit.Shuffle(xs)
	if n > len(s) {
		n = len(s)
	}
	return s[:n]
}
