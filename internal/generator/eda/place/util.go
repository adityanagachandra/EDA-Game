package place

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// pt is an integer point on the placement grid.
type pt struct{ x, y int }

func absI(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func minI(xs ...int) int {
	m := xs[0]
	for _, x := range xs[1:] {
		if x < m {
			m = x
		}
	}
	return m
}

func maxI(xs ...int) int {
	m := xs[0]
	for _, x := range xs[1:] {
		if x > m {
			m = x
		}
	}
	return m
}

func manhattan(a, b pt) int { return absI(a.x-b.x) + absI(a.y-b.y) }

// hpwl is the half-perimeter of the bounding box of pts.
func hpwl(pts []pt) int {
	if len(pts) == 0 {
		return 0
	}
	xs, ys := make([]int, len(pts)), make([]int, len(pts))
	for i, p := range pts {
		xs[i], ys[i] = p.x, p.y
	}
	return maxI(xs...) - minI(xs...) + maxI(ys...) - minI(ys...)
}

// span is max minus min of xs.
func span(xs []int) int { return maxI(xs...) - minI(xs...) }

func fmtPt(p pt) string { return edakit.Sprintf("(%d, %d)", p.x, p.y) }

func fmtPts(pts []pt) string {
	parts := make([]string, len(pts))
	for i, p := range pts {
		parts[i] = fmtPt(p)
	}
	return strings.Join(parts, ", ")
}

func fmtInts(xs []int) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = edakit.IntStr(x)
	}
	return strings.Join(parts, ", ")
}

func fmtFloats(xs []float64, places int) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = edakit.Num(x, places)
	}
	return strings.Join(parts, ", ")
}

// distinctInts draws n distinct integers in [lo, hi].
func distinctInts(n, lo, hi int) []int {
	seen := map[int]bool{}
	out := make([]int, 0, n)
	for try := 0; len(out) < n && try < generator.MaxDrawAttempts*n; try++ {
		v := edakit.Int(lo, hi)
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	for v := lo; len(out) < n; v++ {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

func randPts(n, lo, hi int) []pt {
	out := make([]pt, n)
	for i := range out {
		out[i] = pt{edakit.Int(lo, hi), edakit.Int(lo, hi)}
	}
	return out
}

// numChoices turns a correct float and candidate wrong floats into distinct option strings with a
// unit suffix, padding with nearby values when the candidates collide.
func numChoices(correct float64, places int, unit string, cands ...float64) (string, []string) {
	ans := edakit.Num(correct, places)
	seen := map[string]bool{ans: true}
	var d []string
	for _, c := range cands {
		s := edakit.Num(c, places)
		if !seen[s] {
			seen[s] = true
			d = append(d, s+unit)
		}
	}
	step := math.Max(math.Abs(correct)*0.1, math.Pow(10, -float64(places)))
	for k := 1; len(d) < 3; k++ {
		for _, c := range []float64{correct + float64(k)*step, correct - float64(k)*step} {
			s := edakit.Num(c, places)
			if !seen[s] && len(d) < 3 {
				seen[s] = true
				d = append(d, s+unit)
			}
		}
	}
	if len(d) > 4 {
		d = d[:4]
	}
	return ans + unit, d
}

// permutations returns every permutation of 0..n-1.
func permutations(n int) [][]int {
	var out [][]int
	var rec func(cur []int, used []bool)
	rec = func(cur []int, used []bool) {
		if len(cur) == n {
			out = append(out, append([]int(nil), cur...))
			return
		}
		for i := 0; i < n; i++ {
			if !used[i] {
				used[i] = true
				rec(append(cur, i), used)
				used[i] = false
			}
		}
	}
	rec(nil, make([]bool, n))
	return out
}

// solve solves a small dense linear system by Gaussian elimination with partial pivoting.
func solve(a [][]float64, b []float64) []float64 {
	n := len(b)
	m := make([][]float64, n)
	for i := range m {
		m[i] = append(append([]float64(nil), a[i]...), b[i])
	}
	for c := 0; c < n; c++ {
		p := c
		for r := c + 1; r < n; r++ {
			if math.Abs(m[r][c]) > math.Abs(m[p][c]) {
				p = r
			}
		}
		m[c], m[p] = m[p], m[c]
		for r := c + 1; r < n; r++ {
			f := m[r][c] / m[c][c]
			for k := c; k <= n; k++ {
				m[r][k] -= f * m[c][k]
			}
		}
	}
	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		s := m[i][n]
		for k := i + 1; k < n; k++ {
			s -= m[i][k] * x[k]
		}
		x[i] = s / m[i][i]
	}
	return x
}

func sortedCopy(xs []int) []int {
	out := append([]int(nil), xs...)
	sort.Ints(out)
	return out
}

func letter(i int) string { return string(rune('A' + i)) }
