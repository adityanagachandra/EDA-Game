package synth

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// rgraph is a Leiserson-Saxe retiming graph: node delays and register-weighted edges.
type rgraph struct {
	d     []int
	edges [][3]int // u, v, w
}

func (g *rgraph) desc() string {
	var ds, es []string
	for i, d := range g.d {
		ds = append(ds, edakit.Sprintf("d(v%d) = %d", i+1, d))
	}
	for _, e := range g.edges {
		es = append(es, edakit.Sprintf("v%d to v%d with %d register%s", e[0]+1, e[1]+1, e[2], map[bool]string{true: "", false: "s"}[e[2] == 1]))
	}
	return "Node delays (ns): " + strings.Join(ds, ", ") + ". Edges: " + strings.Join(es, "; ") + "."
}

// period is the longest path delay through zero-register edges, or -1 on a zero-weight cycle.
func (g *rgraph) period(w []int) int {
	n := len(g.d)
	memo := make([]int, n)
	state := make([]int, n)
	bad := false
	var longest func(v int) int
	longest = func(v int) int {
		if state[v] == 2 {
			return memo[v]
		}
		if state[v] == 1 {
			bad = true
			return 0
		}
		state[v] = 1
		best := 0
		for i, e := range g.edges {
			if e[0] == v && w[i] == 0 {
				if x := longest(e[1]); x > best {
					best = x
				}
			}
		}
		state[v] = 2
		memo[v] = g.d[v] + best
		return memo[v]
	}
	p := 0
	for v := 0; v < n; v++ {
		if x := longest(v); x > p {
			p = x
		}
	}
	if bad {
		return -1
	}
	return p
}

func (g *rgraph) weights() []int {
	w := make([]int, len(g.edges))
	for i, e := range g.edges {
		w[i] = e[2]
	}
	return w
}

func (g *rgraph) retimed(r []int) ([]int, bool) {
	w := make([]int, len(g.edges))
	for i, e := range g.edges {
		w[i] = e[2] + r[e[1]] - r[e[0]]
		if w[i] < 0 {
			return nil, false
		}
	}
	return w, true
}

func sum(xs []int) int {
	s := 0
	for _, x := range xs {
		s += x
	}
	return s
}

// randRGraph draws a ring with chords whose every cycle carries a register.
func randRGraph(n int) *rgraph {
	for {
		g := &rgraph{d: make([]int, n)}
		for i := range g.d {
			g.d[i] = edakit.Int(1, 6)
		}
		for i := 0; i < n; i++ {
			g.edges = append(g.edges, [3]int{i, (i + 1) % n, edakit.Int(0, 2)})
		}
		for c := 0; c < edakit.Int(1, 2); c++ {
			u, v := edakit.Int(0, n-1), edakit.Int(0, n-1)
			if u == v || (u+1)%n == v {
				continue
			}
			g.edges = append(g.edges, [3]int{u, v, edakit.Int(0, 2)})
		}
		if sum(g.weights()) >= 2 && g.period(g.weights()) > 0 {
			return g
		}
	}
}

// retiming: registers on an edge, total registers, clock period after a given retiming.
func retiming(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 5)
	for {
		g := randRGraph(n)
		r := make([]int, n)
		for i := range r {
			r[i] = edakit.Int(-1, 1)
		}
		w, ok := g.retimed(r)
		if !ok || g.period(w) < 0 {
			continue
		}
		var rs []string
		for i, x := range r {
			rs = append(rs, edakit.Sprintf("r(v%d) = %d", i+1, x))
		}
		stem := g.desc() + " Apply the retiming " + strings.Join(rs, ", ") + ", where an edge $e: u \\to v$ gets $w_r(e) = w(e) + r(v) - r(u)$ registers."
		switch kp(ctx) {
		case 1:
			i := edakit.Int(0, len(g.edges)-1)
			e := g.edges[i]
			return edakit.Numeric(stem+edakit.Sprintf(" How many registers sit on the edge from v%d to v%d after retiming?", e[0]+1, e[1]+1),
				edakit.IntStr(w[i]), edakit.Sprintf("$w_r = %d + %d - (%d) = %d$.", e[2], r[e[1]], r[e[0]], w[i]),
				map[string]string{"w_r": edakit.IntStr(w[i])})
		case 2:
			t := sum(w)
			return edakit.Numeric(stem+" What is the total number of registers after retiming?",
				edakit.IntStr(t), edakit.Sprintf("Apply the formula on every edge and sum: %d registers, against %d before. The total changes by $\\sum_v r(v)(\\mathrm{in}(v) - \\mathrm{out}(v))$, so only nodes with unequal fanin and fanout edge counts can change it.", t, sum(g.weights())),
				map[string]string{"registers": edakit.IntStr(t)})
		}
		p := g.period(w)
		return edakit.Numeric(stem+" What is the clock period of the retimed circuit, in ns (the largest total node delay along a path of zero-register edges)?",
			edakit.IntStr(p), edakit.Sprintf("Before retiming the period is %d ns. After, the longest register-free path sums to %d ns.", g.period(g.weights()), p),
			map[string]string{"period": edakit.IntStr(p)})
	}
}

// wd computes W and D for every ordered pair by simple-path enumeration.
func (g *rgraph) wd() (W, D [][]int) {
	n := len(g.d)
	W = make([][]int, n)
	D = make([][]int, n)
	for u := 0; u < n; u++ {
		W[u] = make([]int, n)
		D[u] = make([]int, n)
		for v := range W[u] {
			W[u][v] = 1 << 20
		}
		visited := make([]bool, n)
		var dfs func(x, w, d int)
		dfs = func(x, w, d int) {
			if w < W[u][x] || (w == W[u][x] && d > D[u][x]) {
				W[u][x], D[u][x] = w, d
			}
			visited[x] = true
			for _, e := range g.edges {
				if e[0] == x && !visited[e[1]] {
					dfs(e[1], w+e[2], d+g.d[e[1]])
				}
			}
			visited[x] = false
		}
		dfs(u, 0, g.d[u])
	}
	return W, D
}

// wdMatrices: W(u,v), D(u,v), number of period constraints.
func wdMatrices(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 5)
	g := randRGraph(n)
	W, D := g.wd()
	stem := g.desc()
	u, v := edakit.Int(0, n-1), edakit.Int(0, n-2)
	if v >= u {
		v++
	}
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(stem+edakit.Sprintf(" What is W(v%d, v%d), the minimum number of registers on any path from v%d to v%d?", u+1, v+1, u+1, v+1),
			edakit.IntStr(W[u][v]), edakit.Sprintf("Minimize the register count over the paths from v%d to v%d: %d.", u+1, v+1, W[u][v]),
			map[string]string{"W": edakit.IntStr(W[u][v])})
	case 2:
		return edakit.Numeric(stem+edakit.Sprintf(" What is D(v%d, v%d) in ns: the largest total delay, endpoints included, over the paths from v%d to v%d that carry exactly W(v%d, v%d) registers?", u+1, v+1, u+1, v+1, u+1, v+1),
			edakit.IntStr(D[u][v]), edakit.Sprintf("W(v%d, v%d) = %d. Among the paths with that many registers, the slowest has delay %d ns.", u+1, v+1, W[u][v], D[u][v]),
			map[string]string{"W": edakit.IntStr(W[u][v]), "D": edakit.IntStr(D[u][v])})
	}
	maxd := 0
	for _, x := range g.d {
		if x > maxd {
			maxd = x
		}
	}
	c := edakit.Int(maxd, maxd+6)
	cnt := 0
	for a := 0; a < n; a++ {
		for b := 0; b < n; b++ {
			if a != b && D[a][b] > c {
				cnt++
			}
		}
	}
	return edakit.Numeric(stem+edakit.Sprintf(" To test the clock period c = %d ns, Leiserson-Saxe adds the constraint $r(u) - r(v) \\le W(u,v) - 1$ for every ordered pair of distinct nodes with $D(u,v) > c$. How many such constraints are there?", c),
		edakit.IntStr(cnt), edakit.Sprintf("Compute D for all %d ordered pairs and count those above %d ns: %d constraints, alongside one $r(u) - r(v) \\le w(e)$ per edge.", n*(n-1), c, cnt),
		map[string]string{"constraints": edakit.IntStr(cnt)})
}

// bestRetiming searches all legal retimings with r(v1) = 0 and returns the minimum period
// and the fewest registers at that period.
func (g *rgraph) bestRetiming() (minP, minRegs int) {
	n := len(g.d)
	W, _ := g.wd()
	lo := make([]int, n)
	hi := make([]int, n)
	for v := 1; v < n; v++ {
		lo[v], hi[v] = -W[0][v], W[v][0]
	}
	minP, minRegs = 1<<30, 1<<30
	r := make([]int, n)
	var rec func(v int)
	rec = func(v int) {
		if v == n {
			w, ok := g.retimed(r)
			if !ok {
				return
			}
			p := g.period(w)
			if p < 0 {
				return
			}
			s := sum(w)
			if p < minP || (p == minP && s < minRegs) {
				minP, minRegs = p, s
			}
			return
		}
		for x := lo[v]; x <= hi[v]; x++ {
			r[v] = x
			rec(v + 1)
		}
	}
	rec(1)
	return
}

// minPeriod: original period, minimum period, fewest registers at minimum period.
func minPeriod(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 3, 5)
	g := randRGraph(n)
	p0 := g.period(g.weights())
	mp, mr := g.bestRetiming()
	stem := g.desc()
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(stem+" What is the clock period of this circuit before retiming, in ns?",
			edakit.IntStr(p0), edakit.Sprintf("The longest path through register-free edges, adding node delays, is %d ns.", p0),
			map[string]string{"period": edakit.IntStr(p0)})
	case 2:
		return edakit.Numeric(stem+" What is the minimum clock period, in ns, that any legal retiming can achieve?",
			edakit.IntStr(mp), edakit.Sprintf("The original period is %d ns. Testing candidate periods from the D matrix with the W-D constraints (or searching every legal lag vector) gives a minimum of %d ns.", p0, mp),
			map[string]string{"period": edakit.IntStr(p0), "min_period": edakit.IntStr(mp)})
	}
	return edakit.Numeric(stem+" Among all legal retimings that achieve the minimum clock period, what is the fewest total registers?",
		edakit.IntStr(mr), edakit.Sprintf("The minimum period is %d ns, and the cheapest retiming meeting it uses %d registers (the circuit starts with %d).", mp, mr, sum(g.weights())),
		map[string]string{"min_period": edakit.IntStr(mp), "registers": edakit.IntStr(mr)})
}
