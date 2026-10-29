package flow

import (
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// kp returns the knowledge point to serve: the requested one, or a uniform draw from 1..n when
// the request is unguided.
func kp(ctx generator.GeneratorContext, n int) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= n {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(n)
}

// comma formats an integer with thousands separators, for readable option text.
func comma(n int) string {
	neg := n < 0
	if neg {
		n = -n
	}
	s := strconv.Itoa(n)
	var b strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}
	if neg {
		return "-" + b.String()
	}
	return b.String()
}

// pickOthers returns n distinct strings from pool that differ from every excluded string.
func pickOthers(pool []string, n int, exclude ...string) []string {
	ex := map[string]bool{}
	for _, e := range exclude {
		ex[e] = true
	}
	var out []string
	for _, s := range edakit.Shuffle(pool) {
		if ex[s] {
			continue
		}
		ex[s] = true
		out = append(out, s)
		if len(out) == n {
			break
		}
	}
	return out
}

// lowerFirst lowercases the first letter of s for use mid-sentence, unless s starts with an
// acronym (a second capital letter), as in "RTL" or "FPGA".
func lowerFirst(s string) string {
	if len(s) < 2 || (s[1] >= 'A' && s[1] <= 'Z') {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

// nearInt reports whether x is within tol of an integer, so a floor() answer is not on a knife edge.
func nearInt(x, tol float64) bool { return math.Abs(x-math.Round(x)) < tol }

// f2 formats a float compactly for a stem.
func f2(x float64, places int) string { return edakit.Num(x, places) }

func joinInts(xs []int) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = strconv.Itoa(x)
	}
	return strings.Join(parts, ", ")
}
