package graphs

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func laplacianGen(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1, 2:
		n := edakit.Scaled(ctx, 4, 7)
		g := randWGraph(n, edakit.Scaled(ctx, 2, 5), 6)
		intro := "A weighted graph on v1 to v" + edakit.IntStr(n) + " has edges " + edgeText(g, true) + ". Its Laplacian is $L = D - A$. "
		if ctx.KPOrdinal == 1 || (ctx.KPOrdinal == 0 && edakit.Int(0, 1) == 0) {
			i := edakit.Int(0, n-1)
			d := 0
			for j := 0; j < n; j++ {
				d += g.w[i][j]
			}
			return edakit.Numeric(intro+edakit.Sprintf("What is the diagonal entry $L_{%d%d}$?", i+1, i+1), edakit.IntStr(d),
				edakit.Sprintf("The diagonal holds the weighted degree, the sum of the weights of edges at %s: %d.", vName(i), d),
				map[string]string{"degree": edakit.IntStr(d)})
		}
		i, j := edakit.Int(0, n-1), edakit.Int(0, n-1)
		for i == j {
			j = edakit.Int(0, n-1)
		}
		if i > j {
			i, j = j, i
		}
		return edakit.Numeric(intro+edakit.Sprintf("What is the off-diagonal entry $L_{%d%d}$?", i+1, j+1), edakit.IntStr(-g.w[i][j]),
			edakit.Sprintf("Off the diagonal L holds minus the edge weight, and 0 where there is no edge: %d.", -g.w[i][j]),
			map[string]string{"entry": edakit.IntStr(-g.w[i][j])})
	default:
		n := edakit.Scaled(ctx, 5, 8)
		h := randHyper(n, edakit.Scaled(ctx, 4, 7), 5)
		c := edakit.Int(0, n-1)
		d := 0.0
		var terms []string
		for _, net := range h.nets {
			for _, p := range net {
				if p == c {
					k := len(net)
					d += float64(k-1) * 2 / float64(k)
					terms = append(terms, edakit.Sprintf("%d·2/%d", k-1, k))
				}
			}
		}
		if len(terms) == 0 {
			terms = []string{"0"}
		}
		return edakit.Numeric("Nets: "+netText(h, false)+". Each k-pin net is expanded as a clique with edge weight $2/k$ and the resulting graph's Laplacian is formed. What is the diagonal entry for "+cellName(c)+"? Round to 2 decimal places.",
			edakit.Num(d, 2),
			"In a k-pin net, "+cellName(c)+" has k - 1 clique edges of weight 2/k. Summing over its nets: "+strings.Join(terms, " + ")+" = "+edakit.Num(d, 2)+".",
			map[string]string{"degree": edakit.Num(d, 2)})
	}
}

func quadForm(g wgraph, x []float64) float64 {
	s := 0.0
	for i := 0; i < g.n; i++ {
		for j := i + 1; j < g.n; j++ {
			if g.w[i][j] > 0 {
				d := x[i] - x[j]
				s += float64(g.w[i][j]) * d * d
			}
		}
	}
	return s
}

func cutQuadraticGen(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 6)
		g := randWGraph(n, edakit.Scaled(ctx, 1, 4), 4)
		x := make([]float64, n)
		var xi []int
		for i := range x {
			v := edakit.Int(-3, 4)
			x[i] = float64(v)
			xi = append(xi, v)
		}
		q := quadForm(g, x)
		return edakit.Numeric("A weighted graph has edges "+edgeText(g, true)+". For the vector x = ("+intsText(xi)+"), indexed by v1 onward, what is $x^T L x$?",
			edakit.Num(q, 0),
			"Use the identity $x^T L x = \\sum_{(i,j)} w_{ij}(x_i - x_j)^2$ over edges rather than forming L: the sum is "+edakit.Num(q, 0)+".",
			map[string]string{"quadratic": edakit.Num(q, 0)})
	case 2:
		n := edakit.Scaled(ctx, 5, 8)
		g := randWGraph(n, edakit.Scaled(ctx, 2, 6), 5)
		side := randSides(n)
		x := make([]float64, n)
		var xs []string
		for i := range x {
			x[i] = 1
			xs = append(xs, "1")
			if side[i] == 1 {
				x[i] = -1
				xs[i] = "-1"
			}
		}
		q := quadForm(g, x)
		return edakit.Numeric("A weighted graph has edges "+edgeText(g, true)+". A bisection is encoded as x = ("+strings.Join(xs, ", ")+"), +1 for side A and -1 for side B. What is $x^T L x$?",
			edakit.Num(q, 0),
			edakit.Sprintf("Uncut edges contribute 0 and each cut edge contributes w · 2² = 4w. The cut weight is %d, so $x^T L x = 4 \\times %d = %s$.", graphCut(g, side), graphCut(g, side), edakit.Num(q, 0)),
			map[string]string{"cut": edakit.IntStr(graphCut(g, side))})
	default:
		n := edakit.Scaled(ctx, 5, 7)
		h := randHyper(n, edakit.Scaled(ctx, 3, 6), 4)
		x := make([]int, n)
		for i := range x {
			x[i] = edakit.Int(0, 8)
		}
		total := 0.0
		for _, net := range h.nets {
			k := len(net)
			for a := 0; a < k; a++ {
				for b := a + 1; b < k; b++ {
					d := float64(x[net[a]] - x[net[b]])
					total += d * d / float64(k-1)
				}
			}
		}
		var xt []string
		for i, v := range x {
			xt = append(xt, edakit.Sprintf("%s at x = %d", cellName(i), v))
		}
		return edakit.Numeric("Nets: "+netText(h, false)+". Cells are placed on a line: "+strings.Join(xt, ", ")+". With each k-pin net expanded as a clique of weight $1/(k-1)$, the quadratic wirelength is $x^T L x$. What is it? Round to 2 decimal places.",
			edakit.Num(total, 2),
			"Sum over each net's pin pairs of (x_i - x_j)²/(k - 1): "+edakit.Num(total, 2)+". This is the objective a quadratic placer minimizes, and it is the same quadratic form whose ±1 version measures cut size.",
			map[string]string{"quadratic": edakit.Num(total, 2)})
	}
}

func ratioCutGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 6, 9)
	g := randWGraph(n, edakit.Scaled(ctx, 3, 7), 4)
	intro := "A weighted graph has edges " + edgeText(g, true) + ". "
	size := func(side []int) (int, int) {
		a := 0
		for _, s := range side {
			if s == 0 {
				a++
			}
		}
		return a, len(side) - a
	}
	switch kpOf(ctx) {
	case 1:
		side := randSides(n)
		a, b := size(side)
		c := graphCut(g, side)
		v := float64(c) / float64(a*b)
		return edakit.Numeric(intro+"For the partition A = "+vSideText(side, 0)+", B = "+vSideText(side, 1)+", what is the ratio cut $\\text{cut}(A,B)/(|A|\\,|B|)$? Round to 3 decimal places.",
			edakit.Num(v, 3),
			edakit.Sprintf("The cut weight is %d and |A||B| = %d × %d = %d, so the ratio cut is %s.", c, a, b, a*b, edakit.Num(v, 3)),
			map[string]string{"cut": edakit.IntStr(c), "ratio": edakit.Num(v, 3)})
	case 2:
		for {
			s1, s2 := randSides(n), randSides(n)
			a1, b1 := size(s1)
			a2, b2 := size(s2)
			r1 := float64(graphCut(g, s1)) / float64(a1*b1)
			r2 := float64(graphCut(g, s2)) / float64(a2*b2)
			if math.Abs(r1-r2) < 0.002 {
				continue
			}
			best := math.Min(r1, r2)
			return edakit.Numeric(intro+"Partition P has A = "+vSideText(s1, 0)+"; partition Q has A = "+vSideText(s2, 0)+" (B is the rest in each). Which has the smaller ratio cut $\\text{cut}/(|A|\\,|B|)$? Report that smaller ratio cut, rounded to 3 decimal places.",
				edakit.Num(best, 3),
				edakit.Sprintf("P: %d/(%d·%d) = %s. Q: %d/(%d·%d) = %s. The ratio cut rewards balance without a hard constraint.", graphCut(g, s1), a1, b1, edakit.Num(r1, 3), graphCut(g, s2), a2, b2, edakit.Num(r2, 3)),
				map[string]string{"p": edakit.Num(r1, 3), "q": edakit.Num(r2, 3)})
		}
	default:
		side := randSides(n)
		a, b := size(side)
		c := graphCut(g, side)
		v := float64(c)/float64(a) + float64(c)/float64(b)
		return edakit.Numeric(intro+"For the partition A = "+vSideText(side, 0)+", B = "+vSideText(side, 1)+", what is the spectral-clustering RatioCut $\\text{cut}/|A| + \\text{cut}/|B|$? Round to 2 decimal places.",
			edakit.Num(v, 2),
			edakit.Sprintf("%d/%d + %d/%d = %s, which equals n times the Wei-Cheng ratio cut, n·cut/(|A||B|). Its relaxation over real vectors orthogonal to the all-ones vector is solved by the Fiedler vector.", c, a, c, b, edakit.Num(v, 2)),
			map[string]string{"ratiocut": edakit.Num(v, 2)})
	}
}

// jacobiEigen returns the eigenvalues (ascending) and the matching eigenvectors (columns) of a
// symmetric matrix by cyclic Jacobi rotations.
func jacobiEigen(a [][]float64) ([]float64, [][]float64) {
	n := len(a)
	m := make([][]float64, n)
	v := make([][]float64, n)
	for i := range m {
		m[i] = append([]float64(nil), a[i]...)
		v[i] = make([]float64, n)
		v[i][i] = 1
	}
	for sweep := 0; sweep < 100; sweep++ {
		off := 0.0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				off += m[i][j] * m[i][j]
			}
		}
		if off < 1e-20 {
			break
		}
		for p := 0; p < n; p++ {
			for q := p + 1; q < n; q++ {
				if math.Abs(m[p][q]) < 1e-15 {
					continue
				}
				theta := (m[q][q] - m[p][p]) / (2 * m[p][q])
				t := math.Copysign(1, theta) / (math.Abs(theta) + math.Sqrt(theta*theta+1))
				c := 1 / math.Sqrt(t*t+1)
				s := t * c
				for k := 0; k < n; k++ {
					mkp, mkq := m[k][p], m[k][q]
					m[k][p] = c*mkp - s*mkq
					m[k][q] = s*mkp + c*mkq
				}
				for k := 0; k < n; k++ {
					mpk, mqk := m[p][k], m[q][k]
					m[p][k] = c*mpk - s*mqk
					m[q][k] = s*mpk + c*mqk
				}
				for k := 0; k < n; k++ {
					vkp, vkq := v[k][p], v[k][q]
					v[k][p] = c*vkp - s*vkq
					v[k][q] = s*vkp + c*vkq
				}
			}
		}
	}
	idx := seq(n)
	sort.Slice(idx, func(i, j int) bool { return m[idx[i]][idx[i]] < m[idx[j]][idx[j]] })
	vals := make([]float64, n)
	vecs := make([][]float64, n)
	for r, i := range idx {
		vals[r] = m[i][i]
		vecs[r] = make([]float64, n)
		for k := 0; k < n; k++ {
			vecs[r][k] = v[k][i]
		}
	}
	return vals, vecs
}

func lapMatrix(g wgraph) [][]float64 {
	L := make([][]float64, g.n)
	for i := range L {
		L[i] = make([]float64, g.n)
		for j := 0; j < g.n; j++ {
			if i != j && g.w[i][j] > 0 {
				L[i][j] = -float64(g.w[i][j])
				L[i][i] += float64(g.w[i][j])
			}
		}
	}
	return L
}

func fiedlerGen(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		w1, w2 := edakit.Int(1, 9), edakit.Int(1, 9)
		l2 := float64(w1+w2) - math.Sqrt(float64(w1*w1-w1*w2+w2*w2))
		return edakit.Numeric(edakit.Sprintf("Three cells form a path: c1-c2 with weight %d and c2-c3 with weight %d. What is the algebraic connectivity $\\lambda_2$ of its Laplacian? Round to 2 decimal places.", w1, w2),
			edakit.Num(l2, 2),
			edakit.Sprintf("The characteristic polynomial is $\\lambda(\\lambda^2 - 2(w_1+w_2)\\lambda + 3w_1w_2)$, so $\\lambda_2 = (w_1+w_2) - \\sqrt{w_1^2 - w_1w_2 + w_2^2} = %d - \\sqrt{%d} = %s$.", w1+w2, w1*w1-w1*w2+w2*w2, edakit.Num(l2, 2)),
			map[string]string{"lambda2": edakit.Num(l2, 2)})
	case 2:
		for {
			n := 2 * edakit.Scaled(ctx, 3, 4)
			g := newWGraph(n)
			perm := edakit.Shuffle(seq(n))
			half := n / 2
			for i := 0; i < n; i++ {
				for j := i + 1; j < n; j++ {
					same := (i < half) == (j < half)
					a, b := perm[i], perm[j]
					if same && edakit.Int(0, 9) < 7 {
						g.w[a][b] = edakit.Int(2, 5)
					} else if !same && edakit.Int(0, 9) < 2 {
						g.w[a][b] = 1
					}
					g.w[b][a] = g.w[a][b]
				}
			}
			vals, vecs := jacobiEigen(lapMatrix(g))
			if vals[1] < 1e-6 || vals[2]-vals[1] < 0.2 {
				continue
			}
			f := vecs[1]
			if f[0] < 0 {
				for i := range f {
					f[i] = -f[i]
				}
			}
			r := make([]float64, n)
			for i := range f {
				r[i] = edakit.Round(f[i], 2)
			}
			sorted := append([]float64(nil), r...)
			sort.Float64s(sorted)
			if sorted[half]-sorted[half-1] < 0.03 {
				continue
			}
			med := (sorted[half] + sorted[half-1]) / 2
			side := make([]int, n)
			for i := range r {
				if r[i] > med {
					side[i] = 1
				}
			}
			var ft []string
			for i := range r {
				ft = append(ft, edakit.Num(r[i], 2))
			}
			c := graphCut(g, side)
			return edakit.Numeric("A weighted graph has edges "+edgeText(g, true)+". Its Fiedler vector (eigenvector of $\\lambda_2$), indexed by v1 onward, is ("+strings.Join(ft, ", ")+"). Spectral bisection puts the half of the vertices with the smallest entries on one side and the rest on the other. What is the weight of the resulting cut?",
				edakit.IntStr(c),
				"Splitting at the median gives "+vSideText(side, 0)+" and "+vSideText(side, 1)+", which cuts weight "+edakit.IntStr(c)+". Vertices with close entries are strongly connected, so the median split falls between the two clusters.",
				map[string]string{"cut": edakit.IntStr(c), "lambda2": edakit.Num(vals[1], 3)})
		}
	default:
		n := edakit.Scaled(ctx, 4, 6)
		g := randWGraph(n, edakit.Scaled(ctx, 1, 4), 4)
		for {
			x := make([]float64, n)
			var xi []int
			sum := 0
			for i := 0; i < n-1; i++ {
				v := edakit.Int(-3, 3)
				xi = append(xi, v)
				sum += v
			}
			if sum < -3 || sum > 3 {
				continue
			}
			xi = append(xi, -sum)
			nz := 0.0
			for i, v := range xi {
				x[i] = float64(v)
				nz += float64(v * v)
			}
			if nz == 0 {
				continue
			}
			q := quadForm(g, x)
			rq := q / nz
			vals, _ := jacobiEigen(lapMatrix(g))
			return edakit.Numeric("A weighted graph has edges "+edgeText(g, true)+". The vector x = ("+intsText(xi)+") sums to zero, so it is orthogonal to the all-ones vector. What is its Rayleigh quotient $x^T L x / x^T x$? Round to 2 decimal places.",
				edakit.Num(rq, 2),
				edakit.Sprintf("$x^T L x = %s$ and $x^T x = %s$, so the quotient is %s. Any such x gives an upper bound on $\\lambda_2$, here %s, which is the minimum over all of them.", edakit.Num(q, 0), edakit.Num(nz, 0), edakit.Num(rq, 2), edakit.Num(vals[1], 3)),
				map[string]string{"rayleigh": edakit.Num(rq, 2)})
		}
	}
}
