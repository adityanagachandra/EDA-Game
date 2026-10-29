package opt

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// half is one linear inequality a*x + b*y <= r.
type half struct{ a, b, r float64 }

type pt struct{ x, y float64 }

// vertices enumerates the vertices of {x, y >= 0} intersected with the half-planes.
func vertices(cons []half) []pt {
	all := append([]half{{-1, 0, 0}, {0, -1, 0}}, cons...)
	var out []pt
	for i := 0; i < len(all); i++ {
		for j := i + 1; j < len(all); j++ {
			p, q := all[i], all[j]
			det := p.a*q.b - p.b*q.a
			if math.Abs(det) < 1e-12 {
				continue
			}
			x := (p.r*q.b - p.b*q.r) / det
			y := (p.a*q.r - p.r*q.a) / det
			ok := true
			for _, h := range all {
				if h.a*x+h.b*y > h.r+1e-9 {
					ok = false
					break
				}
			}
			if !ok {
				continue
			}
			dup := false
			for _, v := range out {
				if math.Abs(v.x-x) < 1e-9 && math.Abs(v.y-y) < 1e-9 {
					dup = true
					break
				}
			}
			if !dup {
				out = append(out, pt{x, y})
			}
		}
	}
	return out
}

// lpMax solves max c1 x + c2 y over the bounded region, returning value and point.
func lpMax(c1, c2 float64, cons []half) (float64, pt) {
	best, arg := math.Inf(-1), pt{}
	for _, v := range vertices(cons) {
		if z := c1*v.x + c2*v.y; z > best+1e-12 {
			best, arg = z, v
		}
	}
	return best, arg
}

// twoConstraintLP draws max c1 x + c2 y s.t. two constraints that cross at an integer point
// (px, py) in the open quadrant, with constraint 1 steeper than constraint 2. When inCone is
// true the objective is a positive combination y1*(a1,b1) + y2*(a2,b2), so (px, py) is optimal
// and (y1, y2) are the optimal dual values.
type lp2 struct {
	a1, b1, r1, a2, b2, r2 int
	c1, c2                 int
	px, py                 int
	y1, y2                 int
}

func drawLP2(inCone bool) lp2 {
	for {
		var l lp2
		l.px, l.py = edakit.Int(1, 6), edakit.Int(1, 6)
		l.a1, l.b1 = edakit.Int(2, 6), edakit.Int(1, 3)
		l.a2, l.b2 = edakit.Int(1, 3), edakit.Int(2, 6)
		if l.a1*l.b2 <= l.a2*l.b1 {
			continue
		}
		l.r1 = l.a1*l.px + l.b1*l.py
		l.r2 = l.a2*l.px + l.b2*l.py
		if inCone {
			l.y1, l.y2 = edakit.Int(1, 3), edakit.Int(1, 3)
			l.c1 = l.y1*l.a1 + l.y2*l.a2
			l.c2 = l.y1*l.b1 + l.y2*l.b2
		} else {
			l.c1, l.c2 = edakit.Int(1, 8), edakit.Int(1, 8)
		}
		return l
	}
}

func (l lp2) cons() []half {
	return []half{{float64(l.a1), float64(l.b1), float64(l.r1)}, {float64(l.a2), float64(l.b2), float64(l.r2)}}
}

func (l lp2) stem() string {
	return edakit.Sprintf("A linearized budget LP allots x and y units of two resources to two blocks. Maximize $z = %s$ subject to (1) $%s \\le %d$, (2) $%s \\le %d$, and $x, y \\ge 0$.",
		lin2(l.c1, "x", l.c2, "y"), lin2(l.a1, "x", l.b1, "y"), l.r1, lin2(l.a2, "x", l.b2, "y"), l.r2)
}

func lpGeometry(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		l := drawLP2(false)
		z := l.c1*l.px + l.c2*l.py
		return num(l.stem()+" What is the objective value at the vertex where constraints (1) and (2) are both tight?"+roundNote(0),
			edakit.IntStr(z),
			edakit.Sprintf("Solve the two tight constraints together: $x = %d$, $y = %d$. Then $z = %d \\cdot %d + %d \\cdot %d = %d$.", l.px, l.py, l.c1, l.px, l.c2, l.py, z),
			"x", edakit.IntStr(l.px), "y", edakit.IntStr(l.py), "z", edakit.IntStr(z))
	case 2:
		l := drawLP2(edakit.Int(0, 1) == 1)
		z, v := lpMax(float64(l.c1), float64(l.c2), l.cons())
		var vs []string
		for _, p := range vertices(l.cons()) {
			vs = append(vs, edakit.Sprintf("(%s, %s) gives %s", f(p.x, 2), f(p.y, 2), f(float64(l.c1)*p.x+float64(l.c2)*p.y, 2)))
		}
		return num(l.stem()+" What is the optimal objective value?"+roundNote(2), f(z, 2),
			"An LP optimum sits at a vertex. The vertices and their objective values are "+strings.Join(vs, "; ")+edakit.Sprintf(". The best is %s at (%s, %s).", f(z, 2), f(v.x, 2), f(v.y, 2)),
			"z_opt", f(z, 2))
	default:
		l := drawLP2(false)
		cons := l.cons()
		var extra string
		if edakit.Int(0, 1) == 0 {
			u := edakit.Int(1, 8)
			cons = append(cons, half{1, 0, float64(u)})
			extra = edakit.Sprintf("$x \\le %d$", u)
		} else {
			u := edakit.Int(1, 8)
			cons = append(cons, half{0, 1, float64(u)})
			extra = edakit.Sprintf("$y \\le %d$", u)
		}
		vs := vertices(cons)
		var list []string
		for _, p := range vs {
			list = append(list, edakit.Sprintf("(%s, %s)", f(p.x, 2), f(p.y, 2)))
		}
		return num(l.stem()+" A third constraint (3) "+extra+" is added. How many vertices (corner points) does the feasible region now have?"+roundNote(0),
			edakit.IntStr(len(vs)),
			"Intersect every pair of boundary lines, including the axes, and keep the points that satisfy all constraints. The corners are "+strings.Join(list, ", ")+edakit.Sprintf(", so there are %d. A constraint that never touches the region adds no corner.", len(vs)),
			"vertices", edakit.IntStr(len(vs)))
	}
}

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

func simplexPivot(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n, m := edakit.Scaled(ctx, 2, 6), edakit.Int(2, 4)
		b := binom(n+m, m)
		return num(edakit.Sprintf("A routing-capacity LP has %d variables and %d inequality constraints of the form $Ax \\le b$, plus $x \\ge 0$. After adding one slack variable per inequality to reach standard form, what is the largest possible number of basic solutions (choices of the basic variables)?", n, m)+roundNote(0),
			edakit.IntStr(b),
			edakit.Sprintf("Standard form has %d + %d = %d variables and %d equality rows. A basic solution picks %d basic variables, so there are at most $\\binom{%d}{%d} = %d$ of them. Simplex visits only a few, which is why it is fast in practice.", n, m, n+m, m, m, n+m, m, b),
			"vars", edakit.IntStr(n+m), "bases", edakit.IntStr(b))
	case 2:
		rows := edakit.Scaled(ctx, 3, 5)
		var col, rhs []int
		var parts []string
		best, bi := math.Inf(1), -1
		for {
			col, rhs, parts = nil, nil, nil
			best, bi = math.Inf(1), -1
			for i := 0; i < rows; i++ {
				a := edakit.Int(-2, 5)
				b := edakit.Int(2, 20)
				col = append(col, a)
				rhs = append(rhs, b)
				parts = append(parts, edakit.Sprintf("row %d: coefficient %d, right-hand side %d", i+1, a, b))
				if a > 0 && float64(b)/float64(a) < best-1e-9 {
					best, bi = float64(b)/float64(a), i
				}
			}
			if bi >= 0 {
				// require a unique minimum ratio
				uniq := true
				for i := 0; i < rows; i++ {
					if i != bi && col[i] > 0 && math.Abs(float64(rhs[i])/float64(col[i])-best) < 1e-9 {
						uniq = false
					}
				}
				if uniq {
					break
				}
			}
		}
		return num("In a simplex tableau the entering variable's column and the right-hand sides are: "+strings.Join(parts, "; ")+". By the minimum-ratio test, what value does the entering variable take after the pivot?"+roundNote(2),
			f(best, 2),
			edakit.Sprintf("Only rows with a positive coefficient limit the step; rows with a zero or negative coefficient never become negative. The smallest ratio is row %d: %d / %d = %s, so that row's basic variable leaves.", bi+1, rhs[bi], col[bi], f(best, 2)),
			"ratio", f(best, 2), "leaving_row", edakit.IntStr(bi+1))
	default:
		z0 := edakit.Int(10, 60)
		d := edakit.Int(1, 6)
		a1, a2 := edakit.Int(1, 4), edakit.Int(1, 4)
		b1, b2 := edakit.Int(4, 24), edakit.Int(4, 24)
		theta := math.Min(float64(b1)/float64(a1), float64(b2)/float64(a2))
		z := float64(z0) + float64(d)*theta
		return num(edakit.Sprintf("A maximization tableau has current objective %d. Nonbasic variable $x_e$ has reduced cost %d (each unit of $x_e$ adds %d to the objective). Its column has coefficients %d and %d in rows with right-hand sides %d and %d. After pivoting $x_e$ in, what is the new objective value?", z0, d, d, a1, a2, b1, b2)+roundNote(2),
			f(z, 2),
			edakit.Sprintf("The ratio test gives the step $\\theta = \\min(%d/%d, %d/%d) = %s$. The objective rises by the reduced cost times the step: $%d + %d \\cdot %s = %s$.", b1, a1, b2, a2, f(theta, 2), z0, d, f(theta, 2), f(z, 2)),
			"theta", f(theta, 4), "z_new", f(z, 2))
	}
}

func lpDuality(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		l := drawLP2(false)
		var y1, y2 int
		for {
			y1, y2 = edakit.Int(0, 5), edakit.Int(0, 5)
			if y1*l.a1+y2*l.a2 >= l.c1 && y1*l.b1+y2*l.b2 >= l.c2 {
				break
			}
		}
		ub := y1*l.r1 + y2*l.r2
		return num(l.stem()+edakit.Sprintf(" The dual is: minimize $%s$ subject to $%s \\ge %d$, $%s \\ge %d$, $y \\ge 0$. The point $y_1 = %d$, $y_2 = %d$ is dual feasible. What upper bound on the primal optimum does it certify?", lin2(l.r1, "y_1", l.r2, "y_2"), lin2(l.a1, "y_1", l.a2, "y_2"), l.c1, lin2(l.b1, "y_1", l.b2, "y_2"), l.c2, y1, y2)+roundNote(0),
			edakit.IntStr(ub),
			edakit.Sprintf("Weak duality: every dual-feasible point bounds the primal maximum from above. The bound is $%d \\cdot %d + %d \\cdot %d = %d$.", l.r1, y1, l.r2, y2, ub),
			"bound", edakit.IntStr(ub))
	case 2:
		l := drawLP2(edakit.Int(0, 1) == 1)
		z, _ := lpMax(float64(l.c1), float64(l.c2), l.cons())
		return num(l.stem()+" What is the optimal value of the dual LP, minimize $r_1 y_1 + r_2 y_2$ over $A^T y \\ge c$, $y \\ge 0$?"+roundNote(2),
			f(z, 2),
			edakit.Sprintf("Strong duality: when the primal has a finite optimum, the dual optimum equals it. Solve the primal by its vertices: the best is %s.", f(z, 2)),
			"z_opt", f(z, 2))
	default:
		l := drawLP2(true)
		k := edakit.Int(1, 2)
		y := l.y1
		if k == 2 {
			y = l.y2
		}
		return num(l.stem()+edakit.Sprintf(" The optimum is at $x = %d$, $y = %d$, where both constraints are tight. What is the optimal dual value $y_%d$ of constraint (%d)?", l.px, l.py, k, k)+roundNote(2),
			f(float64(y), 2),
			edakit.Sprintf("Both constraints are tight, so the dual solves $A^T y = c$: $%s = %d$ and $%s = %d$. That gives $y_1 = %d$ and $y_2 = %d$, both nonnegative, which confirms the vertex is optimal.", lin2(l.a1, "y_1", l.a2, "y_2"), l.c1, lin2(l.b1, "y_1", l.b2, "y_2"), l.c2, l.y1, l.y2),
			"y1", edakit.IntStr(l.y1), "y2", edakit.IntStr(l.y2))
	}
}

func shadowPrices(ctx generator.GeneratorContext) generator.Problem {
	l := drawLP2(true)
	var third half
	var thirdS string
	if edakit.Int(0, 1) == 0 {
		u := l.px + edakit.Int(1, 4)
		third = half{1, 0, float64(u)}
		thirdS = edakit.Sprintf("(3) $x \\le %d$", u)
	} else {
		s := l.px + l.py + edakit.Int(1, 4)
		third = half{1, 1, float64(s)}
		thirdS = edakit.Sprintf("(3) $x + y \\le %d$", s)
	}
	cons := append(l.cons(), third)
	stem := strings.TrimSuffix(l.stem(), ", and $x, y \\ge 0$.") + ", " + thirdS + ", and $x, y \\ge 0$. The optimum is at $x = " + edakit.IntStr(l.px) + "$, $y = " + edakit.IntStr(l.py) + "$."
	slack3 := third.r - third.a*float64(l.px) - third.b*float64(l.py)
	switch kp(ctx) {
	case 1:
		return num(stem+" What is the slack of constraint (3) at the optimum?"+roundNote(0),
			f(slack3, 0),
			edakit.Sprintf("Slack is right-hand side minus left-hand side at the optimum: %s. A constraint with positive slack is not binding.", f(slack3, 0)),
			"slack3", f(slack3, 0))
	case 2:
		k := edakit.Int(1, 3)
		y := map[int]int{1: l.y1, 2: l.y2, 3: 0}[k]
		return num(stem+edakit.Sprintf(" Using complementary slackness, what is the optimal dual value (shadow price) $y_%d$ of constraint (%d)?", k, k)+roundNote(2),
			f(float64(y), 2),
			edakit.Sprintf("Constraint (3) has slack %s, so complementary slackness forces $y_3 = 0$. The tight constraints (1) and (2) carry the duals, from $%s = %d$ and $%s = %d$: $y_1 = %d$, $y_2 = %d$.", f(slack3, 0), lin2(l.a1, "y_1", l.a2, "y_2"), l.c1, lin2(l.b1, "y_1", l.b2, "y_2"), l.c2, l.y1, l.y2),
			"y1", edakit.IntStr(l.y1), "y2", edakit.IntStr(l.y2), "y3", "0")
	default:
		for tries := 0; ; tries++ {
			k := edakit.Int(1, 2)
			delta := edakit.Int(1, 2)
			z0, _ := lpMax(float64(l.c1), float64(l.c2), cons)
			mod := append([]half(nil), cons...)
			mod[k-1].r += float64(delta)
			z1, _ := lpMax(float64(l.c1), float64(l.c2), mod)
			y := l.y1
			if k == 2 {
				y = l.y2
			}
			consistent := math.Abs((z1-z0)-float64(y*delta)) < 1e-6
			if consistent || tries > 20 {
				expl := edakit.Sprintf("The shadow price of constraint (%d) is its dual value $y_%d = %d$, the marginal objective per unit of right-hand side while the same constraints stay tight. The gain is $%d \\cdot %d = %s$.", k, k, y, y, delta, f(z1-z0, 2))
				if !consistent {
					expl = edakit.Sprintf("The shadow price $y_%d = %d$ holds only while the same constraints stay tight. Here the shift is large enough that another constraint becomes binding, so re-solving the LP gives a gain of %s rather than %d.", k, y, f(z1-z0, 2), y*delta)
				}
				return num(stem+edakit.Sprintf(" The right-hand side of constraint (%d) is raised by %d (more routing capacity, say). By how much does the optimal objective increase?", k, delta)+roundNote(2),
					f(z1-z0, 2),
					expl, "y", edakit.IntStr(y), "gain", f(z1-z0, 2))
			}
		}
	}
}

func differenceConstraints(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 7)
		type e struct{ i, j, d int }
		var es []e
		for j := 1; j < n; j++ {
			es = append(es, e{j - 1, j, edakit.Int(1, 6)})
			if j >= 2 && edakit.Int(0, 1) == 1 {
				i := edakit.Int(0, j-2)
				es = append(es, e{i, j, edakit.Int(3, 12)})
			}
		}
		dist := make([]int, n)
		for j := 1; j < n; j++ {
			dist[j] = -1 << 30
		}
		for _, x := range es { // edges are generated in topological order of j
			if dist[x.i]+x.d > dist[x.j] {
				dist[x.j] = dist[x.i] + x.d
			}
		}
		var parts []string
		for _, x := range es {
			parts = append(parts, edakit.Sprintf("$x_%d \\ge x_%d + %d$", x.j, x.i, x.d))
		}
		return num(edakit.Sprintf("A one-dimensional layout compactor places %d edges left to right with $x_0 = 0$ and minimum-spacing constraints ", n)+strings.Join(parts, ", ")+edakit.Sprintf(" (in nm). What is the smallest feasible value of $x_%d$?", n-1)+roundNote(0),
			edakit.IntStr(dist[n-1]),
			edakit.Sprintf("Each constraint is an edge of weight d in the constraint graph, and the minimum position of a node is its longest-path distance from $x_0$. Relaxing the edges in topological order gives positions %s, so $x_%d = %d$.", ints(dist, ", "), n-1, dist[n-1]),
			"x_last", edakit.IntStr(dist[n-1]))
	case 2:
		for {
			n := edakit.Scaled(ctx, 3, 5)
			type e struct{ i, j, w int }
			var es []e
			m := n + edakit.Int(1, n)
			for k := 0; k < m; k++ {
				i, j := edakit.Int(1, n), edakit.Int(1, n)
				if i == j {
					continue
				}
				es = append(es, e{i, j, edakit.Int(-4, 6)})
			}
			if len(es) < n {
				continue
			}
			// shortest paths from a virtual source with 0-weight edges to every variable
			d := make([]int, n+1)
			for it := 0; it < n+1; it++ {
				for _, x := range es {
					if d[x.i]+x.w < d[x.j] {
						d[x.j] = d[x.i] + x.w
					}
				}
			}
			neg := false
			for _, x := range es {
				if d[x.i]+x.w < d[x.j] {
					neg = true
				}
			}
			if neg {
				continue
			}
			var parts []string
			for _, x := range es {
				parts = append(parts, edakit.Sprintf("$x_%d - x_%d \\le %d$", x.j, x.i, x.w))
			}
			k := edakit.Int(1, n)
			return num("A clock-skew or retiming feasibility check gives the difference constraints "+strings.Join(parts, ", ")+edakit.Sprintf(". Add a source s with $x_v - s \\le 0$ for every variable and take $x_v$ as the shortest-path distance from s (Bellman-Ford). What value does this give $x_%d$?", k)+roundNote(0),
				edakit.IntStr(d[k]),
				edakit.Sprintf("Each constraint $x_j - x_i \\le w$ is an edge i to j of weight w, and shortest-path distances satisfy every edge's triangle inequality, so they are a feasible solution. Starting from 0 everywhere and relaxing to a fixed point gives %s for $x_1$ to $x_%d$.", ints(d[1:], ", "), n),
				"x_k", edakit.IntStr(d[k]))
		}
	default:
		k := edakit.Scaled(ctx, 2, 5)
		ds := make([]int, k)
		sum := 0
		for i := range ds {
			ds[i] = 10 * edakit.Int(20, 90)
			sum += ds[i]
		}
		setup := 10 * edakit.Int(0, 5)
		var parts []string
		for i := 0; i < k; i++ {
			parts = append(parts, edakit.Sprintf("FF%d to FF%d is %d ps", i+1, (i+1)%k+1, ds[i]))
		}
		tmin := float64(sum)/float64(k) + float64(setup)
		return num(edakit.Sprintf("Useful clock skew lets each flop i have its own clock arrival $s_i$. %d flops form a loop: ", k)+strings.Join(parts, ", ")+edakit.Sprintf(". Setup at each capture flop is %d ps and each stage needs $s_i + d_i + t_{setup} \\le s_{i+1} + T$. What is the smallest clock period T for which skews exist? Answer in ps.", setup)+roundNote(2),
			f(tmin, 2),
			edakit.Sprintf("Each constraint is $s_i - s_{i+1} \\le T - d_i - t_{setup}$. Summing around the loop cancels the skews, so feasibility (no negative cycle) needs $kT \\ge \\sum d_i + k\\,t_{setup}$: $T \\ge %d/%d + %d = %s$ ps. Skew can balance a loop but never beat its average delay.", sum, k, setup, f(tmin, 2)),
			"t_min", f(tmin, 2))
	}
}

// lin2 renders a*u + b*v with unit coefficients written bare.
func lin2(a int, u string, b int, v string) string {
	t := func(c int, w string) string {
		if c == 1 {
			return w
		}
		return edakit.IntStr(c) + w
	}
	return t(a, u) + " + " + t(b, v)
}
