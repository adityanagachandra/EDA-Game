package synth

import (
	"math/bits"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Two-level helpers. A function of n <= 5 variables is a truth table held in a uint64: bit m
// is f(m). Variable i (a, b, c, ...) is minterm bit n-1-i, so a is the most significant bit
// and the minterm index reads as the binary word abcd, the K-map convention.

var varNames = []string{"a", "b", "c", "d", "e", "f", "g", "h"}

// cube is a product term: mask marks the specified minterm bits, val holds their values.
type cube struct{ mask, val uint32 }

func (c cube) contains(m int) bool { return uint32(m)&c.mask == c.val }

func (c cube) lits() int { return bits.OnesCount32(c.mask) }

func (c cube) size(n int) int { return 1 << (n - c.lits()) }

// minterms returns the truth-table bits the cube covers.
func (c cube) table(n int) uint64 {
	var t uint64
	for m := 0; m < 1<<n; m++ {
		if c.contains(m) {
			t |= 1 << m
		}
	}
	return t
}

// tex renders the cube as a product of literals, $a\bar b c$ style (no dollar signs).
func (c cube) tex(n int) string {
	if c.mask == 0 {
		return "1"
	}
	var b strings.Builder
	for i := 0; i < n; i++ {
		p := uint32(1) << (n - 1 - i)
		if c.mask&p == 0 {
			continue
		}
		if c.val&p != 0 {
			b.WriteString(varNames[i])
		} else {
			b.WriteString("\\bar{" + varNames[i] + "}")
		}
	}
	return b.String()
}

// pos renders positional cube notation, one of 0, 1, - per variable in order a, b, c, ...
func (c cube) pos(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		p := uint32(1) << (n - 1 - i)
		switch {
		case c.mask&p == 0:
			b.WriteByte('-')
		case c.val&p != 0:
			b.WriteByte('1')
		default:
			b.WriteByte('0')
		}
	}
	return b.String()
}

// allCubes enumerates the 3^n cubes over n variables.
func allCubes(n int) []cube {
	out := []cube{{}}
	for i := 0; i < n; i++ {
		p := uint32(1) << (n - 1 - i)
		var next []cube
		for _, c := range out {
			next = append(next, c, cube{c.mask | p, c.val}, cube{c.mask | p, c.val | p})
		}
		out = next
	}
	return out
}

func randCube(n, minLits int) cube {
	for {
		var c cube
		for i := 0; i < n; i++ {
			p := uint32(1) << i
			switch edakit.Int(0, 2) {
			case 1:
				c.mask |= p
			case 2:
				c.mask |= p
				c.val |= p
			}
		}
		if c.lits() >= minLits {
			return c
		}
	}
}

// primes returns the prime implicants of ON with don't cares DC that cover at least one ON
// minterm, ordered by size then position.
func primes(n int, on, dc uint64) []cube {
	care := on | dc
	impl := func(c cube) bool { return c.table(n)&^care == 0 }
	var out []cube
	for _, c := range allCubes(n) {
		if !impl(c) || c.table(n)&on == 0 {
			continue
		}
		prime := true
		for i := 0; i < n && prime; i++ {
			p := uint32(1) << i
			if c.mask&p != 0 && impl(cube{c.mask &^ p, c.val &^ p}) {
				prime = false
			}
		}
		if prime {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].lits() != out[j].lits() {
			return out[i].lits() < out[j].lits()
		}
		return out[i].pos(n) < out[j].pos(n)
	})
	return out
}

// minCover returns the minimum cover of on by the given cubes: fewest cubes, then fewest
// literals. It brute-forces subsets by increasing size, which is fine for the small tables here.
func minCover(n int, on uint64, cs []cube) []cube {
	k := len(cs)
	tabs := make([]uint64, k)
	for i, c := range cs {
		tabs[i] = c.table(n) & on
	}
	for size := 0; size <= k; size++ {
		var best []int
		bestLits := 1 << 30
		var rec func(start int, chosen []int, cov uint64)
		rec = func(start int, chosen []int, cov uint64) {
			if len(chosen) == size {
				if cov == on {
					l := 0
					for _, i := range chosen {
						l += cs[i].lits()
					}
					if l < bestLits {
						bestLits = l
						best = append([]int(nil), chosen...)
					}
				}
				return
			}
			for i := start; i < k; i++ {
				rec(i+1, append(chosen, i), cov|tabs[i])
			}
		}
		rec(0, nil, 0)
		if best != nil {
			out := make([]cube, len(best))
			for i, j := range best {
				out[i] = cs[j]
			}
			return out
		}
	}
	return nil
}

func coverLits(cs []cube) int {
	l := 0
	for _, c := range cs {
		l += c.lits()
	}
	return l
}

func sopTex(n int, cs []cube) string {
	if len(cs) == 0 {
		return "0"
	}
	parts := make([]string, len(cs))
	for i, c := range cs {
		parts[i] = c.tex(n)
	}
	return strings.Join(parts, " + ")
}

func mintermList(t uint64, n int) []int {
	var out []int
	for m := 0; m < 1<<n; m++ {
		if t>>m&1 == 1 {
			out = append(out, m)
		}
	}
	return out
}

func intsCSV(xs []int) string {
	s := make([]string, len(xs))
	for i, x := range xs {
		s[i] = edakit.IntStr(x)
	}
	return strings.Join(s, ", ")
}

// sigma renders $f = \sum m(...)$ with an optional don't-care term.
func sigma(name string, on, dc uint64, n int) string {
	s := "$" + name + "(" + strings.Join(varNames[:n], ",") + ") = \\sum m(" + intsCSV(mintermList(on, n)) + ")"
	if dc != 0 {
		s += " + d(" + intsCSV(mintermList(dc, n)) + ")"
	}
	return s + "$"
}

// randFunc draws a nonconstant ON set with between lo and hi minterms.
func randFunc(n, lo, hi int) uint64 {
	full := 1 << n
	if hi > full-1 {
		hi = full - 1
	}
	k := edakit.Int(lo, hi)
	perm := edakit.Shuffle(seq(full))
	var t uint64
	for _, m := range perm[:k] {
		t |= 1 << m
	}
	return t
}

func seq(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i
	}
	return out
}

func popcount(t uint64) int { return bits.OnesCount64(t) }

// essentials returns the primes that alone cover some ON minterm.
func essentials(n int, on uint64, ps []cube) []cube {
	var out []cube
	for i, p := range ps {
		mine := p.table(n) & on
		for j, q := range ps {
			if j != i {
				mine &^= q.table(n)
			}
		}
		if mine != 0 {
			out = append(out, p)
		}
	}
	return out
}
