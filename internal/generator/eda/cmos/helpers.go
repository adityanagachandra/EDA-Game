package cmos

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func sqrt(x float64) float64   { return math.Sqrt(x) }
func pow(x, y float64) float64 { return math.Pow(x, y) }
func ln(x float64) float64     { return math.Log(x) }
func log10(x float64) float64  { return math.Log10(x) }
func ceil(x float64) int       { return int(math.Ceil(x - 1e-9)) }
func floor(x float64) int      { return int(math.Floor(x + 1e-9)) }
func absf(x float64) float64   { return math.Abs(x) }

// article prefixes a cell or gate name with "a" or "an".
func article(name string) string {
	if name == "" {
		return name
	}
	switch name[0] {
	case 'A', 'E', 'I', 'O', 'U', 'a', 'e', 'i', 'o', 'u', 'X':
		return "an " + name
	}
	return "a " + name
}

// capArticle is article with a capital first letter, for the start of a sentence.
func capArticle(name string) string {
	s := article(name)
	return strings.ToUpper(s[:1]) + s[1:]
}

// sciTex writes x as m \times 10^{e} for use inside $...$.
func sciTex(x float64) string {
	e := floor(log10(x))
	m := x / pow(10, float64(e))
	return edakit.Sprintf("%s \\times 10^{%d}", edakit.Num(m, 3), e)
}

// signed writes "+ n" or "- n" for a term appended to an expression.
func signed(n int) string {
	if n < 0 {
		return edakit.Sprintf("- %d", -n)
	}
	return edakit.Sprintf("+ %d", n)
}
