// Package boolalg generates Boolean-reasoning problems: Boolean functions and their cofactors,
// cubes and covers, decomposition, bit-parallel and hash-consed representations, decision
// diagrams, lattices and fixpoints, CNF and SAT encodings, algorithm analysis and
// computational complexity, and the agent-first use of Boolean engines.
//
// This file is the small Boolean-function kit the generators share. A function of n <= 5
// variables is a truth table packed into a uint32: bit i holds f at minterm i, and minterm i
// writes i in binary with variable 0 (named a) as the most significant bit. That is the
// textbook minterm numbering, so m(5) over a, b, c is a = 1, b = 0, c = 1.
package boolalg

import (
	"math/bits"
	"math/rand"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var varNames = []string{"a", "b", "c", "d", "e"}

// tt is a truth table over n variables.
type tt struct {
	n int
	w uint32
}

func fullMask(n int) uint32 {
	if n >= 5 {
		return 0xFFFFFFFF
	}
	return (uint32(1) << (uint(1) << uint(n))) - 1
}

func rows(n int) int { return 1 << uint(n) }

// bitOf returns the value of variable v in minterm i.
func bitOf(i, v, n int) int { return (i >> uint(n-1-v)) & 1 }

func (f tt) at(i int) bool { return f.w>>uint(i)&1 == 1 }

func (f tt) ones() int { return bits.OnesCount32(f.w & fullMask(f.n)) }

func (f tt) not() tt { return tt{f.n, ^f.w & fullMask(f.n)} }

func (f tt) and(g tt) tt { return tt{f.n, f.w & g.w} }

func (f tt) or(g tt) tt { return tt{f.n, f.w | g.w} }

func (f tt) xor(g tt) tt { return tt{f.n, (f.w ^ g.w) & fullMask(f.n)} }

func (f tt) minterms() []int {
	var out []int
	for i := 0; i < rows(f.n); i++ {
		if f.at(i) {
			out = append(out, i)
		}
	}
	return out
}

func fromMinterms(n int, ms []int) tt {
	var w uint32
	for _, m := range ms {
		w |= 1 << uint(m)
	}
	return tt{n, w}
}

// varTT is the projection function of variable v.
func varTT(n, v int) tt {
	var w uint32
	for i := 0; i < rows(n); i++ {
		if bitOf(i, v, n) == 1 {
			w |= 1 << uint(i)
		}
	}
	return tt{n, w}
}

// cof is the cofactor f with variable v fixed to val, a function of the other n-1 variables
// in their original order.
func (f tt) cof(v, val int) tt {
	var w uint32
	k := 0
	for i := 0; i < rows(f.n); i++ {
		if bitOf(i, v, f.n) != val {
			continue
		}
		if f.at(i) {
			w |= 1 << uint(k)
		}
		k++
	}
	return tt{f.n - 1, w}
}

func (f tt) dependsOn(v int) bool { return f.cof(v, 0) != f.cof(v, 1) }

func (f tt) supportSize() int {
	s := 0
	for v := 0; v < f.n; v++ {
		if f.dependsOn(v) {
			s++
		}
	}
	return s
}

// permute returns g with g(y) = f(x) where variable v of f sits at position perm[v] of g.
func (f tt) permute(perm []int) tt {
	var w uint32
	for i := 0; i < rows(f.n); i++ {
		if !f.at(i) {
			continue
		}
		j := 0
		for v := 0; v < f.n; v++ {
			if bitOf(i, v, f.n) == 1 {
				j |= 1 << uint(f.n-1-perm[v])
			}
		}
		w |= 1 << uint(j)
	}
	return tt{f.n, w}
}

func (f tt) isConst() bool { return f.w&fullMask(f.n) == 0 || f.w&fullMask(f.n) == fullMask(f.n) }

// ---------- cubes

// cube holds one entry per variable: 0 (complemented literal), 1 (true literal) or 2 (absent).
type cube []int

func (c cube) contains(i, n int) bool {
	for v := 0; v < n; v++ {
		if c[v] != 2 && bitOf(i, v, n) != c[v] {
			return false
		}
	}
	return true
}

func (c cube) tt(n int) tt {
	var w uint32
	for i := 0; i < rows(n); i++ {
		if c.contains(i, n) {
			w |= 1 << uint(i)
		}
	}
	return tt{n, w}
}

func (c cube) lits() int {
	k := 0
	for _, x := range c {
		if x != 2 {
			k++
		}
	}
	return k
}

func (c cube) dashes() int { return len(c) - c.lits() }

func (c cube) size() int { return 1 << uint(c.dashes()) }

func (c cube) latex() string {
	var b strings.Builder
	for v, x := range c {
		switch x {
		case 1:
			b.WriteString(varNames[v])
		case 0:
			b.WriteString(`\overline{` + varNames[v] + `}`)
		}
	}
	if b.Len() == 0 {
		return "1"
	}
	return b.String()
}

func (c cube) positional() string {
	var b strings.Builder
	for _, x := range c {
		switch x {
		case 0:
			b.WriteByte('0')
		case 1:
			b.WriteByte('1')
		default:
			b.WriteByte('-')
		}
	}
	return b.String()
}

func (c cube) equal(d cube) bool {
	for i := range c {
		if c[i] != d[i] {
			return false
		}
	}
	return true
}

// within reports whether every minterm of c lies in d (c is contained in d).
func (c cube) within(d cube) bool {
	for i := range c {
		if d[i] != 2 && c[i] != d[i] {
			return false
		}
	}
	return true
}

func randCube(n, minLits, maxLits int) cube {
	k := edakit.Int(minLits, maxLits)
	c := make(cube, n)
	for i := range c {
		c[i] = 2
	}
	for _, v := range rand.Perm(n)[:k] {
		c[v] = rand.Intn(2)
	}
	return c
}

type sop []cube

func (s sop) tt(n int) tt {
	f := tt{n, 0}
	for _, c := range s {
		f = f.or(c.tt(n))
	}
	return f
}

func (s sop) latex() string {
	parts := make([]string, len(s))
	for i, c := range s {
		parts[i] = c.latex()
	}
	return strings.Join(parts, " + ")
}

func (s sop) mentions(v int) bool {
	for _, c := range s {
		if c[v] != 2 {
			return true
		}
	}
	return false
}

// randSOP draws k distinct cubes over n variables with 1 to n-1 literals each.
func randSOP(n, k int) sop {
	var s sop
	for tries := 0; len(s) < k && tries < 100; tries++ {
		c := randCube(n, 1, n-1)
		dup := false
		for _, d := range s {
			if d.equal(c) {
				dup = true
			}
		}
		if !dup {
			s = append(s, c)
		}
	}
	return s
}

func allCubes(n int) []cube {
	var out []cube
	total := 1
	for i := 0; i < n; i++ {
		total *= 3
	}
	for code := 0; code < total; code++ {
		c := make(cube, n)
		x := code
		for v := n - 1; v >= 0; v-- {
			c[v] = x % 3
			x /= 3
		}
		out = append(out, c)
	}
	return out
}

func isImplicant(c cube, f tt) bool {
	ct := c.tt(f.n)
	return ct.w&^f.w == 0
}

// primes returns the prime implicants of f, ordered by literal count and then position.
func primes(f tt) []cube {
	var imps []cube
	for _, c := range allCubes(f.n) {
		if isImplicant(c, f) {
			imps = append(imps, c)
		}
	}
	var out []cube
	for _, c := range imps {
		prime := true
		for v := range c {
			if c[v] == 2 {
				continue
			}
			d := append(cube(nil), c...)
			d[v] = 2
			if isImplicant(d, f) {
				prime = false
				break
			}
		}
		if prime {
			out = append(out, c)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].lits() < out[j].lits() })
	return out
}

func isPrime(c cube, f tt) bool {
	if !isImplicant(c, f) {
		return false
	}
	for v := range c {
		if c[v] == 2 {
			continue
		}
		d := append(cube(nil), c...)
		d[v] = 2
		if isImplicant(d, f) {
			return false
		}
	}
	return true
}

// essentials returns the primes of care (ON plus don't-care) that alone cover some ON minterm.
func essentials(on tt, ps []cube) []cube {
	var out []cube
	for pi, p := range ps {
		ess := false
		for _, m := range on.minterms() {
			if !p.contains(m, on.n) {
				continue
			}
			alone := true
			for qi, q := range ps {
				if qi != pi && q.contains(m, on.n) {
					alone = false
					break
				}
			}
			if alone {
				ess = true
				break
			}
		}
		if ess {
			out = append(out, p)
		}
	}
	return out
}

// minCover returns the minimum number of primes covering on, and among the covers of that
// size the smallest total literal count. ps must be primes of a function containing on.
func minCover(on tt, ps []cube) (int, int) {
	if on.ones() == 0 {
		return 0, 0
	}
	// Keep only primes that touch the ON-set.
	var cand []cube
	var masks []uint32
	for _, p := range ps {
		m := p.tt(on.n).w & on.w
		if m != 0 {
			cand = append(cand, p)
			masks = append(masks, m)
		}
	}
	best, bestLits := len(cand)+1, 1<<30
	var rec func(start int, covered uint32, k, lits int)
	rec = func(start int, covered uint32, k, lits int) {
		if covered == on.w {
			if k < best || (k == best && lits < bestLits) {
				best, bestLits = k, lits
			}
			return
		}
		if k >= best {
			return
		}
		// Branch on the lowest uncovered minterm: some chosen prime must cover it.
		rest := on.w &^ covered
		m := bits.TrailingZeros32(rest)
		for i := range cand {
			if masks[i]>>uint(m)&1 == 1 {
				rec(i+1, covered|masks[i], k+1, lits+cand[i].lits())
			}
		}
	}
	rec(0, 0, 0, 0)
	return best, bestLits
}

// ---------- rendering

func mintermList(ms []int) string {
	parts := make([]string, len(ms))
	for i, m := range ms {
		parts[i] = edakit.IntStr(m)
	}
	return strings.Join(parts, ", ")
}

func varList(n int) string { return strings.Join(varNames[:n], ", ") }

// sigma renders f as a minterm list, for example f(a, b, c) = sum m(1, 4, 6).
func sigma(name string, f tt) string {
	return "$" + name + "(" + varList(f.n) + `) = \sum m(` + mintermList(f.minterms()) + ")$"
}

// drawFunc draws a non-constant function of n variables and renders it either as a
// sum-of-products expression or as a minterm list.
func drawFunc(n int, exprProb float64) (tt, string) {
	for {
		if rand.Float64() < exprProb {
			s := randSOP(n, edakit.Int(2, n))
			f := s.tt(n)
			if f.isConst() {
				continue
			}
			return f, "$f(" + varList(n) + ") = " + s.latex() + "$"
		}
		f := tt{n, rand.Uint32() & fullMask(n)}
		if f.isConst() || f.ones() < 2 || f.ones() > rows(n)-2 {
			continue
		}
		return f, sigma("f", f)
	}
}

// drawSOP draws a non-constant SOP function with an expression.
func drawSOP(n int) (tt, sop) {
	for {
		s := randSOP(n, edakit.Int(2, n))
		f := s.tt(n)
		if !f.isConst() {
			return f, s
		}
	}
}

func kpOf(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

// nVars draws 3 or 4 variables, 4 more often at higher difficulty.
func nVars(ctx generator.GeneratorContext) int { return edakit.Scaled(ctx, 3, 4) }

func lit(v, val int) string {
	if val == 1 {
		return varNames[v]
	}
	return `\overline{` + varNames[v] + `}`
}

func pow2(k int) int { return 1 << uint(k) }

func binom(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	r := 1
	for i := 1; i <= k; i++ {
		r = r * (n - k + i) / i
	}
	return r
}

// otherVars names the variables of f other than v, for "over b, c, d" phrases.
func otherVars(n int, skip ...int) string {
	var out []string
	for v := 0; v < n; v++ {
		keep := true
		for _, s := range skip {
			if s == v {
				keep = false
			}
		}
		if keep {
			out = append(out, varNames[v])
		}
	}
	return strings.Join(out, ", ")
}
