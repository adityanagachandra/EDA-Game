package cts

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// pdnGrid: strap resistance from sheet resistance, routing share, parallel straps.
func pdnGrid(ctx generator.GeneratorContext) generator.Problem {
	rs := pick(0.02, 0.12, 0.01)
	switch kpOf(ctx) {
	case 1:
		L := float64(10 * edakit.Int(2, 30))
		w := pick(0.5, 6, 0.5)
		r := rs * L / w
		return edakit.Numeric(edakit.Sprintf("A power strap on a metal layer with sheet resistance %s $\\Omega$/sq is %s µm wide. What is the resistance of a %s µm segment between two via stacks, in ohms? Round to 3 decimal places.", f(rs, 2), f(w, 1), f(L, 0)), f(r, 3),
			edakit.Sprintf("The segment is $L/W = %s/%s = %s$ squares long, so $R = R_s L/W = %s$ ohm.", f(L, 0), f(w, 1), f(L/w, 2), f(r, 3)),
			map[string]string{"squares": f(L/w, 3)})
	case 2:
		w := pick(1, 6, 0.5)
		p := float64(5 * edakit.Int(4, 24))
		if 2*w >= p {
			p = 2*w + 20
		}
		share := 100 * 2 * w / p
		return edakit.Numeric(edakit.Sprintf("On one routing layer, a VDD strap and a VSS strap, each %s µm wide, repeat every %s µm. What percentage of the layer's width do the power straps occupy? Round to 1 decimal place.", f(w, 1), f(p, 0)), f(share, 1),
			edakit.Sprintf("Each %s µm pitch holds $2 \\cdot %s$ µm of strap: $%s/%s = %s\\%%$. Every percent here is routing capacity the signal router loses, which is the basic tension of PDN design.", f(p, 0), f(w, 1), f(2*w, 1), f(p, 0), f(share, 1)),
			map[string]string{"share": f(share, 2)})
	}
	L := float64(100 * edakit.Int(2, 20))
	w := pick(1, 6, 0.5)
	n := edakit.Int(4, 40)
	r := 1000 * rs * L / (float64(n) * w)
	return edakit.Numeric(edakit.Sprintf("A block is crossed by %d parallel VDD straps, each %s µm wide and %s µm long, on a layer of %s $\\Omega$/sq. If they are tied together only at their two ends, what is the end-to-end resistance of the strap bundle, in milliohms? Round to 2 decimal places.", n, f(w, 1), f(L, 0), f(rs, 2)), f(r, 2),
		edakit.Sprintf("One strap: $R_s L/W = %s$ ohm. %d in parallel: $%s/%d = %s$ ohm, which is %s mohm.", f(rs*L/w, 3), n, f(rs*L/w, 4), n, f(r/1000, 5), f(r, 2)),
		map[string]string{"r_mohm": f(r, 3)})
}

// ladder solves the drops on a rail with n equal taps of current i (mA) and segment
// resistance r (ohm), fed from the left end only or from both ends.
func ladder(n int, r, i float64, both bool) []float64 {
	g := make([][]float64, n)
	b := make([]float64, n)
	for k := 0; k < n; k++ {
		g[k] = make([]float64, n)
		b[k] = i
		if k > 0 {
			g[k][k] += 1 / r
			g[k][k-1] -= 1 / r
		} else {
			g[k][k] += 1 / r // segment to the left feed
		}
		if k < n-1 {
			g[k][k] += 1 / r
			g[k][k+1] -= 1 / r
		} else if both {
			g[k][k] += 1 / r // segment to the right feed
		}
	}
	return solve(g, b)
}

// railLadder: IR drop along a rail with equal taps.
func railLadder(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 12)
	r := pick(0.1, 1.0, 0.05)
	i := pick(0.5, 5, 0.25)
	desc := edakit.Sprintf("A power rail has %d cells tapped along it, each drawing %s mA. Consecutive taps are joined by rail segments of $%s\\ \\Omega$, and a segment of the same resistance joins the first tap to the supply.", n, f(i, 2), f(r, 2))
	switch kpOf(ctx) {
	case 1:
		d := r * i * float64(n*(n+1)) / 2
		return edakit.Numeric(desc+" The rail is fed from that end only. What is the IR drop at the far tap, in mV? Round to 2 decimal places.", f(d, 2),
			edakit.Sprintf("Segment $j$ carries the current of every tap beyond it, $(%d - j + 1) \\cdot %s$ mA. Summing the segment drops: $R I \\cdot %d \\cdot %d / 2 = %s$ mV. The drop grows with the square of the rail length.", n, f(i, 2), n, n+1, f(d, 2)),
			map[string]string{"far_drop": f(d, 3)})
	case 2:
		k := edakit.Int(2, n-1)
		d := r * i * (float64(k*n) - float64(k*(k-1))/2)
		return edakit.Numeric(desc+edakit.Sprintf(" The rail is fed from that end only. What is the IR drop at tap %d (counting from the supply), in mV? Round to 2 decimal places.", k), f(d, 2),
			edakit.Sprintf("The first %d segments carry $%d, %d, \\ldots, %d$ tap currents. Drop $= R I \\sum_{j=1}^{%d}(%d - j + 1) = %s \\cdot %s \\cdot %s = %s$ mV.", k, n, n-1, n-k+1, k, n, f(r, 2), f(i, 2), f(float64(k*n)-float64(k*(k-1))/2, 0), f(d, 2)),
			map[string]string{"drop": f(d, 3)})
	}
	ds := ladder(n, r, i, true)
	worst := 0.0
	for _, d := range ds {
		worst = math.Max(worst, d)
	}
	single := r * i * float64(n*(n+1)) / 2
	return edakit.Numeric(desc+" A second supply connection, through another such segment, now feeds the far end too. What is the worst IR drop on the rail, in mV? Round to 2 decimal places.", f(worst, 2),
		edakit.Sprintf("Solving the tridiagonal system $G d = i$ for the %d taps gives drops (mV) %s. The worst is %s mV in the middle, against %s mV with a single feed: feeding both ends cuts the worst drop by roughly a factor of four.", n, floats(ds, 2), f(worst, 2), f(single, 2)),
		map[string]string{"worst": f(worst, 3)})
}

type pgrid struct {
	n    int
	r    map[[2]int]float64 // edge resistances, key i<j (0-based)
	pad  []float64          // pad resistance per node, 0 = none
	load []float64          // mA
}

func drawGrid() pgrid {
	g := pgrid{n: 4, r: map[[2]int]float64{}, pad: make([]float64, 4), load: make([]float64, 4)}
	rs := []float64{0.5, 1, 2, 2.5, 4, 5}
	for _, e := range [][2]int{{0, 1}, {1, 3}, {2, 3}, {0, 2}} {
		g.r[e] = edakit.Pick(rs)
	}
	g.pad[0] = edakit.Pick([]float64{0.25, 0.5, 1})
	if edakit.Int(0, 1) == 1 {
		g.pad[3] = edakit.Pick([]float64{0.25, 0.5, 1})
	}
	for i := range g.load {
		g.load[i] = float64(edakit.Int(1, 10))
	}
	return g
}

func (g pgrid) matrix() [][]float64 {
	m := make([][]float64, g.n)
	for i := range m {
		m[i] = make([]float64, g.n)
		if g.pad[i] > 0 {
			m[i][i] += 1 / g.pad[i]
		}
	}
	for e, r := range g.r {
		m[e[0]][e[0]] += 1 / r
		m[e[1]][e[1]] += 1 / r
		m[e[0]][e[1]] -= 1 / r
		m[e[1]][e[0]] -= 1 / r
	}
	return m
}

func (g pgrid) describe() string {
	nm := []string{"n1", "n2", "n3", "n4"}
	var es []string
	for _, e := range [][2]int{{0, 1}, {1, 3}, {2, 3}, {0, 2}} {
		es = append(es, edakit.Sprintf("%s-%s $%s\\ \\Omega$", nm[e[0]], nm[e[1]], f(g.r[e], 2)))
	}
	var pads []string
	for i, p := range g.pad {
		if p > 0 {
			pads = append(pads, edakit.Sprintf("%s through $%s\\ \\Omega$", nm[i], f(p, 2)))
		}
	}
	return edakit.Sprintf("A 2 by 2 power grid has nodes n1 (top left), n2 (top right), n3 (bottom left) and n4 (bottom right), joined by %s. The ideal supply connects to %s. Loads draw %s mA at n1 to n4.",
		strings.Join(es, ", "), strings.Join(pads, " and "), floats(g.load, 0))
}

// conductanceMatrix: diagonal entry, off-diagonal entry, and row sum of G.
func conductanceMatrix(ctx generator.GeneratorContext) generator.Problem {
	g := drawGrid()
	m := g.matrix()
	nm := []string{"n1", "n2", "n3", "n4"}
	desc := g.describe() + " Write the drop equations as $G d = i$, with $d$ the drops below the supply and $i$ the load currents."
	switch kpOf(ctx) {
	case 1:
		i := edakit.Int(0, 3)
		return edakit.Numeric(desc+edakit.Sprintf(" What is the diagonal entry of $G$ for %s, in siemens? Round to 2 decimal places.", nm[i]), f(m[i][i], 2),
			edakit.Sprintf("A diagonal entry is the sum of all conductances touching the node, including a supply connection: %s S.", f(m[i][i], 3)),
			map[string]string{"gii": f(m[i][i], 3)})
	case 2:
		pairs := [][2]int{{0, 1}, {1, 3}, {2, 3}, {0, 2}, {0, 3}, {1, 2}}
		p := edakit.Pick(pairs)
		v := m[p[0]][p[1]]
		why := edakit.Sprintf("Off-diagonal entries are minus the conductance of the branch joining the two nodes: $-1/%s = %s$ S.", f(g.r[p], 2), f(v, 2))
		if v == 0 {
			why = "The two nodes are diagonal neighbours with no branch between them, so the entry is 0. Each row has non-zeros only for a node's grid neighbours, which is why $G$ is sparse."
		}
		return edakit.Numeric(desc+edakit.Sprintf(" What is the entry of $G$ in row %s, column %s, in siemens? Round to 2 decimal places.", nm[p[0]], nm[p[1]]), f(v, 2), why,
			map[string]string{"gij": f(v, 3)})
	}
	i := edakit.Int(0, 3)
	s := 0.0
	for j := 0; j < 4; j++ {
		s += m[i][j]
	}
	why := "Each internal branch adds $+g$ on the diagonal and $-g$ off it, so it cancels in the row sum. Only a supply connection survives: "
	if g.pad[i] > 0 {
		why += edakit.Sprintf("this node's is $1/%s = %s$ S.", f(g.pad[i], 2), f(s, 2))
	} else {
		why += "this node has none, so the row sums to 0. Without any supply connection every row would sum to zero and $G$ would be singular; the pads make it positive definite."
	}
	return edakit.Numeric(desc+edakit.Sprintf(" What is the sum of the entries in row %s of $G$, in siemens? Round to 2 decimal places.", nm[i]), f(s, 2), why,
		map[string]string{"row_sum": f(s, 3)})
}

// irDrop: solve small G d = i systems.
func irDrop(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		r1, r2 := pick(0.2, 2, 0.1), pick(0.2, 2, 0.1)
		i1, i2 := float64(edakit.Int(2, 20)), float64(edakit.Int(2, 20))
		d := r1*(i1+i2) + r2*i2
		return edakit.Numeric(edakit.Sprintf("A supply pad feeds node a through $%s\\ \\Omega$, and node a feeds node b through $%s\\ \\Omega$. Node a draws %s mA and node b draws %s mA. What is the IR drop at node b, in mV? Round to 2 decimal places.", f(r1, 1), f(r2, 1), f(i1, 0), f(i2, 0)), f(d, 2),
			edakit.Sprintf("The first resistor carries both currents and the second only node b's: $%s(%s + %s) + %s \\cdot %s = %s$ mV (ohms times mA give mV).", f(r1, 1), f(i1, 0), f(i2, 0), f(r2, 1), f(i2, 0), f(d, 2)),
			map[string]string{"drop": f(d, 3)})
	case 2:
		rs := []float64{pick(0.2, 2, 0.1), pick(0.2, 2, 0.1), pick(0.2, 2, 0.1), pick(0.2, 2, 0.1)}
		is := []float64{float64(edakit.Int(2, 20)), float64(edakit.Int(2, 20)), float64(edakit.Int(2, 20))}
		g := [][]float64{
			{1/rs[0] + 1/rs[1], -1 / rs[1], 0},
			{-1 / rs[1], 1/rs[1] + 1/rs[2], -1 / rs[2]},
			{0, -1 / rs[2], 1/rs[2] + 1/rs[3]},
		}
		d := solve(g, is)
		k := edakit.Int(0, 2)
		return edakit.Numeric(edakit.Sprintf("A rail runs pad - $%s\\ \\Omega$ - n1 - $%s\\ \\Omega$ - n2 - $%s\\ \\Omega$ - n3 - $%s\\ \\Omega$ - pad, with both pads at the ideal supply. The nodes draw %s mA. Solve $G d = i$ for the drops. What is the IR drop at n%d, in mV? Round to 2 decimal places.", f(rs[0], 1), f(rs[1], 1), f(rs[2], 1), f(rs[3], 1), floats(is, 0), k+1), f(d[k], 2),
			edakit.Sprintf("$G$ is tridiagonal with diagonal (%s) S. Solving gives drops (%s) mV, so n%d sits %s mV below the supply.", floats([]float64{g[0][0], g[1][1], g[2][2]}, 3), floats(d, 3), k+1, f(d[k], 2)),
			map[string]string{"drop": f(d[k], 3)})
	}
	g := drawGrid()
	d := solve(g.matrix(), g.load)
	k := edakit.Int(0, 3)
	worst := 0
	for i := range d {
		if d[i] > d[worst] {
			worst = i
		}
	}
	if edakit.Int(0, 1) == 1 {
		k = worst
	}
	return edakit.Numeric(g.describe()+edakit.Sprintf(" Solve $G d = i$ for the drops below the supply. What is the IR drop at n%d, in mV? Round to 2 decimal places.", k+1), f(d[k], 2),
		edakit.Sprintf("Stamping every branch and supply conductance gives a 4 by 4 symmetric positive-definite $G$. Its solution is $d = $ (%s) mV, so n%d drops %s mV; the worst node is n%d.", floats(d, 3), k+1, f(d[k], 2), worst+1),
		map[string]string{"drop": f(d[k], 3)})
}

// iterativeSolve: Jacobi and Gauss-Seidel sweeps on a 3-node rail, and stencil sparsity.
func iterativeSolve(ctx generator.GeneratorContext) generator.Problem {
	k := kpOf(ctx)
	if k == 3 {
		n := 10 * edakit.Int(5, 200)
		nnz := n*n + 4*n*(n-1)
		return edakit.Numeric(edakit.Sprintf("A power grid layer is modelled as an $%d \\times %d$ array of nodes, each joined to its up, down, left and right neighbours. How many non-zero entries does its conductance matrix $G$ have?", n, n), edakit.IntStr(nnz),
			edakit.Sprintf("Every node has a diagonal entry ($%d^2 = %d$). There are $2 \\cdot %d \\cdot %d = %d$ branches and each puts two symmetric off-diagonal entries in $G$: total %d. That is about 5 per row out of %d columns, which is why grid solvers store $G$ sparsely and use iterative or sparse direct methods.", n, n*n, n, n-1, 2*n*(n-1), nnz, n*n),
			map[string]string{"nnz": edakit.IntStr(nnz)})
	}
	rs := []float64{pick(0.5, 2, 0.5), pick(0.5, 2, 0.5), pick(0.5, 2, 0.5), pick(0.5, 2, 0.5)}
	is := []float64{float64(edakit.Int(2, 20)), float64(edakit.Int(2, 20)), float64(edakit.Int(2, 20))}
	g := [][]float64{
		{1/rs[0] + 1/rs[1], -1 / rs[1], 0},
		{-1 / rs[1], 1/rs[1] + 1/rs[2], -1 / rs[2]},
		{0, -1 / rs[2], 1/rs[2] + 1/rs[3]},
	}
	desc := edakit.Sprintf("A rail runs pad - $%s\\ \\Omega$ - n1 - $%s\\ \\Omega$ - n2 - $%s\\ \\Omega$ - n3 - $%s\\ \\Omega$ - pad, with both pads at the ideal supply, and the nodes draw %s mA. The drops satisfy $G d = i$. Start from $d = 0$.", f(rs[0], 1), f(rs[1], 1), f(rs[2], 1), f(rs[3], 1), floats(is, 0))
	sweep := func(d []float64, gs bool) []float64 {
		out := append([]float64(nil), d...)
		src := d
		if gs {
			src = out
		}
		for i := 0; i < 3; i++ {
			s := is[i]
			for j := 0; j < 3; j++ {
				if j != i {
					s -= g[i][j] * src[j]
				}
			}
			out[i] = s / g[i][i]
		}
		return out
	}
	if k == 1 {
		d1 := sweep([]float64{0, 0, 0}, false)
		d2 := sweep(d1, false)
		return edakit.Numeric(desc+" After two Jacobi sweeps, what is the estimate of the drop at n2, in mV? Round to 3 decimal places.", f(d2[1], 3),
			edakit.Sprintf("Jacobi updates every node from the previous sweep: $d_i \\leftarrow (i_i - \\sum_{j \\ne i} G_{ij} d_j)/G_{ii}$. Sweep 1 gives (%s), sweep 2 gives (%s) mV.", floats(d1, 3), floats(d2, 3)),
			map[string]string{"d2": f(d2[1], 4)})
	}
	d1 := sweep([]float64{0, 0, 0}, true)
	return edakit.Numeric(desc+" After one Gauss-Seidel sweep in the order n1, n2, n3, what is the estimate of the drop at n3, in mV? Round to 3 decimal places.", f(d1[2], 3),
		edakit.Sprintf("Gauss-Seidel uses each new value as soon as it exists: $d_1 = %s$, then $d_2 = %s$ using the new $d_1$, then $d_3 = %s$ mV using the new $d_2$.", f(d1[0], 3), f(d1[1], 3), f(d1[2], 3)),
		map[string]string{"d3": f(d1[2], 4)})
}

// strapPitch: a rail between two ideal straps with uniform current per length.
func strapPitch(ctx generator.GeneratorContext) generator.Problem {
	j := pick(0.002, 0.02, 0.001)
	r := pick(0.2, 1, 0.1)
	desc := edakit.Sprintf("A standard-cell rail with resistance %s $\\Omega$/µm runs between two power straps that hold it at the ideal supply. The cells along it draw a uniform %s mA per µm of rail.", f(r, 1), f(j, 3))
	switch kpOf(ctx) {
	case 1:
		p := float64(5 * edakit.Int(4, 30))
		d := j * r * p * p / 8
		return edakit.Numeric(desc+edakit.Sprintf(" The straps are %s µm apart. What is the worst IR drop on the rail, in mV? Round to 2 decimal places.", f(p, 0)), f(d, 2),
			edakit.Sprintf("The drop obeys $d'' = -jr$ with $d = 0$ at both straps, so $d(x) = \\tfrac{jr}{2}x(P - x)$. The worst is at the midpoint, $jrP^2/8 = %s \\cdot %s \\cdot %s^2/8 = %s$ mV.", f(j, 3), f(r, 1), f(p, 0), f(d, 2)),
			map[string]string{"drop": f(d, 3)})
	case 2:
		b := pick(1, 10, 0.5)
		p := math.Sqrt(8 * b / (j * r))
		return edakit.Numeric(desc+edakit.Sprintf(" The rail's share of the IR budget is %s mV. What is the largest strap pitch that meets it, in µm? Round to 1 decimal place.", f(b, 1)), f(p, 1),
			edakit.Sprintf("Set $jrP^2/8 = %s$: $P = \\sqrt{8 \\cdot %s/(%s \\cdot %s)} = %s$ µm. Halving the pitch quarters the rail drop, at the cost of twice the strap tracks.", f(b, 1), f(b, 1), f(j, 3), f(r, 1), f(p, 1)),
			map[string]string{"pitch": f(p, 2)})
	}
	p := float64(10 * edakit.Int(4, 15))
	x := float64(5 * edakit.Int(1, int(p/5)-1))
	d := j * r * x * (p - x) / 2
	return edakit.Numeric(desc+edakit.Sprintf(" The straps are %s µm apart. What is the IR drop %s µm from the left strap, in mV? Round to 3 decimal places.", f(p, 0), f(x, 0)), f(d, 3),
		edakit.Sprintf("$d(x) = \\tfrac{jr}{2}x(P - x) = \\tfrac{%s \\cdot %s}{2} \\cdot %s \\cdot %s = %s$ mV, a parabola that is zero at both straps.", f(j, 3), f(r, 1), f(x, 0), f(p-x, 0), f(d, 3)),
		map[string]string{"drop": f(d, 4)})
}
