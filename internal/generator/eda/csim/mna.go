package csim

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// mnaStamps asks for an entry, or the sum of entries, of the nodal conductance matrix of a
// random resistor network.
func mnaStamps(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 2, 4)
	rs := randomNetwork(n, edakit.Scaled(ctx, 1, 3))
	g := stampG(n, rs)
	desc := describeResistors(rs)
	switch kp(ctx) {
	case 1:
		i := edakit.Int(1, n)
		return num(edakit.Sprintf("A resistor network has nodes 1 to %d and ground: %s. In the nodal conductance matrix $G$ (unknowns ordered by node number, ground row and column removed), what is the diagonal entry $G_{%d%d}$ in mS?", n, desc, i, i),
			g[i-1][i-1], 3,
			edakit.Sprintf("A diagonal entry sums the conductances of every resistor touching the node, including those to ground: $G_{%d%d} = %s$ mS.", i, i, f(g[i-1][i-1], 3)), nil)
	case 2:
		i := edakit.Int(1, n)
		j := edakit.Int(1, n)
		for j == i {
			j = edakit.Int(1, n)
		}
		return num(edakit.Sprintf("A resistor network has nodes 1 to %d and ground: %s. In the nodal conductance matrix $G$ (ground removed), what is the entry $G_{%d%d}$ in mS?", n, desc, i, j),
			g[i-1][j-1], 3,
			edakit.Sprintf("An off-diagonal entry is minus the total conductance connected directly between the two nodes (zero if none): $G_{%d%d} = %s$ mS.", i, j, f(g[i-1][j-1], 3)), nil)
	}
	sum := 0.0
	toGround := 0.0
	for _, r := range rs {
		if r.b == 0 {
			toGround += 1 / r.r
		}
	}
	for i := range g {
		for j := range g[i] {
			sum += g[i][j]
		}
	}
	return num(edakit.Sprintf("A resistor network has nodes 1 to %d and ground: %s. What is the sum of all entries of the nodal conductance matrix $G$ (ground removed), in mS?", n, desc),
		sum, 3,
		edakit.Sprintf("Each resistor between two non-ground nodes stamps $+g, +g, -g, -g$, which sums to zero. A resistor to ground stamps only $+g$ on one diagonal. So the sum is the total conductance to ground, %s mS.", f(toGround, 3)),
		map[string]string{"sum": f(sum, 3)})
}

type isrc struct {
	from, to int // current flows from node 'from' through the source into node 'to'
	i        float64
}

// currentSources asks for an entry of the right-hand side vector.
func currentSources(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 3, 5)
	m := edakit.Scaled(ctx, 2, 4)
	k := kp(ctx)
	var srcs []isrc
	var parts []string
	for s := 0; s < m; s++ {
		a := edakit.Int(0, n)
		b := edakit.Int(0, n)
		for b == a {
			b = edakit.Int(0, n)
		}
		val := stepF(0.5, 0.5, 7)
		if k == 1 {
			// a source pushing val into node a, drawn from node b
			srcs = append(srcs, isrc{b, a, val})
			switch {
			case b == 0:
				parts = append(parts, edakit.Sprintf("a source pushing %s mA from ground into node %d", f(val, 2), a))
			case a == 0:
				parts = append(parts, edakit.Sprintf("a source drawing %s mA out of node %d into ground", f(val, 2), b))
			default:
				parts = append(parts, edakit.Sprintf("a source pushing %s mA into node %d, drawn from node %d", f(val, 2), a, b))
			}
		} else {
			// SPICE card I n+ n- val: current flows from n+ through the source to n-.
			srcs = append(srcs, isrc{a, b, val})
			parts = append(parts, edakit.Sprintf("I%d %s %s %sm", s+1, nodeName(a), nodeName(b), f(val, 2)))
		}
	}
	rhs := make([]float64, n+1)
	for _, s := range srcs {
		rhs[s.to] += s.i
		rhs[s.from] -= s.i
	}
	switch k {
	case 1:
		node := edakit.Int(1, n)
		return num(edakit.Sprintf("A circuit with nodes 1 to %d and ground has resistors and these current sources: %s. With KCL written as currents leaving through elements equal to source currents injected, what is the right-hand-side entry $b_%d$, in mA?", n, strings.Join(parts, "; "), node),
			rhs[node], 2,
			edakit.Sprintf("Add every current pushed into node %d and subtract every current drawn out of it: $b_%d = %s$ mA.", node, node, f(rhs[node], 2)), nil)
	case 2:
		node := edakit.Int(1, n)
		return num(edakit.Sprintf("A SPICE deck for a circuit with nodes 1 to %d contains these current sources (positive current flows from the first node, through the source, to the second): %s. With $b$ holding currents injected into each node, what is $b_%d$ in mA?", n, strings.Join(parts, "; "), node),
			rhs[node], 2,
			edakit.Sprintf("A card I n+ n- I takes I out of n+ and delivers it into n-, so it stamps $-I$ at n+ and $+I$ at n-. Summing at node %d gives %s mA.", node, f(rhs[node], 2)), nil)
	}
	total := 0.0
	for i := 1; i <= n; i++ {
		total += rhs[i]
	}
	return num(edakit.Sprintf("A SPICE deck for a circuit with nodes 1 to %d contains these current sources (positive current flows from the first node, through the source, to the second; node 0 is ground): %s. After the ground row is dropped, what is the sum of all entries of the right-hand side $b$, in mA?", n, strings.Join(parts, "; ")),
		total, 2,
		edakit.Sprintf("A source between two non-ground nodes stamps $+I$ and $-I$ and cancels. Only sources touching ground survive, so the sum is the net current delivered from ground into the circuit: %s mA.", f(total, 2)),
		map[string]string{"sum": f(total, 2)})
}

type vsrc struct {
	p, m int
	v    float64
}

// voltageSources asks for the augmented system size, an augmented entry, or a branch current.
func voltageSources(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 3, 8)
		nv := edakit.Int(1, 3)
		ni := edakit.Int(0, 3)
		nr := edakit.Int(n, 2*n)
		nc := edakit.Int(0, n)
		return edakit.Numeric(edakit.Sprintf("A circuit has %d non-ground nodes, %d resistors, %d capacitors (open at DC), %d independent current sources and %d independent voltage sources. How many unknowns does the DC modified nodal analysis system have?", n, nr, nc, ni, nv),
			edakit.IntStr(n+nv),
			edakit.Sprintf("Every node voltage is an unknown, and each voltage source adds its branch current, because its current cannot be written from its terminal voltages. Resistors, capacitors and current sources add none: %d + %d = %d.", n, nv, n+nv), nil)
	case 2:
		n := edakit.Scaled(ctx, 2, 3)
		rs := randomNetwork(n, 1)
		nv := edakit.Int(1, 2)
		var vs []vsrc
		var vparts []string
		for k := 0; k < nv; k++ {
			p := edakit.Int(1, n)
			m := 0
			if edakit.Int(0, 2) == 0 && n > 1 {
				m = edakit.Int(1, n)
				for m == p {
					m = edakit.Int(1, n)
				}
			}
			v := stepF(0.5, 0.5, 9)
			vs = append(vs, vsrc{p, m, v})
			vparts = append(vparts, edakit.Sprintf("V%d %s %s %s", k+1, nodeName(p), nodeName(m), f(v, 2)))
		}
		size := n + nv
		a := zeros(size)
		g := stampG(n, rs)
		for i := 0; i < n; i++ {
			copy(a[i], g[i])
		}
		for k, s := range vs {
			row := n + k
			if s.p > 0 {
				a[s.p-1][row] += 1
				a[row][s.p-1] += 1
			}
			if s.m > 0 {
				a[s.m-1][row] -= 1
				a[row][s.m-1] -= 1
			}
		}
		desc := describeResistors(rs) + "; and voltage sources " + strings.Join(vparts, ", ") + " (V name, + node, - node, volts)"
		order := edakit.Sprintf("Unknowns are ordered $v_1, \\ldots, v_%d$, then the branch currents of the voltage sources in card order; resistances give conductances in mS.", n)
		var r, c int
		if edakit.Int(0, 2) == 0 {
			r, c = edakit.Int(1, n), edakit.Int(1, n)
		} else {
			k := edakit.Int(0, nv-1)
			node := edakit.Int(1, n)
			if edakit.Int(0, 1) == 0 {
				r, c = node, n+k+1
			} else {
				r, c = n+k+1, node
			}
		}
		return num("A circuit has "+desc+". "+order+edakit.Sprintf(" What is the entry in row %d, column %d of the MNA matrix?", r, c),
			a[r-1][c-1], 3,
			edakit.Sprintf("The node block holds the conductance stamps; a voltage source puts $+1$ at its + node and $-1$ at its - node in its branch column and, transposed, in its branch row. Entry (%d, %d) is %s.", r, c, f(a[r-1][c-1], 3)), nil)
	}
	vsv := stepF(1, 0.5, 8)
	r1 := edakit.Pick(kohmChoices)
	r2 := edakit.Pick(kohmChoices)
	withR3 := edakit.Int(0, 1) == 0
	r3 := edakit.Pick(kohmChoices)
	i := vsv / (r1 + r2)
	stem := edakit.Sprintf("A deck has V1 1 0 %s, R1 1 2 %sk and R2 2 0 %sk", f(vsv, 2), f(r1, 2), f(r2, 2))
	if withR3 {
		i += vsv / r3
		stem += edakit.Sprintf(" and R3 1 0 %sk", f(r3, 2))
	}
	ans := -i
	return num(stem+". MNA gives V1 a branch-current unknown with the SPICE sign convention: positive current flows into the + terminal, through the source, and out of the - terminal. What value does the simulator report for that branch current, in mA?",
		ans, 3,
		edakit.Sprintf("The source delivers %s mA out of its + terminal into the resistors. Current into the + terminal is the negative of that, so the branch current is %s mA; a source that delivers power reports a negative current.", f(i, 3), f(ans, 3)),
		map[string]string{"delivered_ma": f(i, 3)})
}

// vccsStamp adds a VCCS into g: current gm*(v(cp)-v(cm)) flows from np through the source to nm.
func vccsStamp(g [][]float64, np, nm, cp, cm int, gm float64) {
	add := func(r, c int, v float64) {
		if r > 0 && c > 0 {
			g[r-1][c-1] += v
		}
	}
	add(np, cp, gm)
	add(np, cm, -gm)
	add(nm, cp, -gm)
	add(nm, cm, gm)
}

// vccsStamps asks for entries of a matrix that contains a VCCS.
func vccsStamps(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 3, 4)
	nodes := edakit.Shuffle([]int{1, 2, 3, 4}[:n])
	np, cp := nodes[0], nodes[1]
	nm, cm := 0, 0
	if edakit.Int(0, 1) == 0 {
		nm = nodes[2]
	}
	if edakit.Int(0, 1) == 0 && n > 3 {
		cm = nodes[3]
	} else if edakit.Int(0, 2) == 0 {
		cm = nodes[2]
	}
	gm := stepF(0.5, 0.5, 9)
	card := edakit.Sprintf("G1 %s %s %s %s %sm", nodeName(np), nodeName(nm), nodeName(cp), nodeName(cm), f(gm, 2))
	intro := "A VCCS card G1 n+ n- nc+ nc- gm drives a current $g_m (v_{nc+} - v_{nc-})$ from n+ through the source into n-. "
	k := kp(ctx)
	g := zeros(n)
	var rs []resistor
	if k >= 2 {
		rs = randomNetwork(n, 1)
		g = stampG(n, rs)
	}
	vccsStamp(g, np, nm, cp, cm, gm)
	ctrl := []int{cp}
	if cm > 0 {
		ctrl = append(ctrl, cm)
	}
	out := []int{np}
	if nm > 0 {
		out = append(out, nm)
	}
	switch k {
	case 1:
		r, c := edakit.Pick(out), edakit.Pick(ctrl)
		return num(intro+edakit.Sprintf("A circuit with nodes 1 to %d contains %s (gm in mS, 0 is ground). Which entries does this card alone stamp into the nodal matrix entry $G_{%d%d}$, in mS?", n, card, r, c),
			g[r-1][c-1], 2,
			edakit.Sprintf("The VCCS stamps $+g_m$ at (n+, nc+), $-g_m$ at (n+, nc-), $-g_m$ at (n-, nc+) and $+g_m$ at (n-, nc-), dropping any ground row or column. Entry (%d, %d) is %s mS.", r, c, f(g[r-1][c-1], 2)), nil)
	case 2:
		r, c := edakit.Pick(out), edakit.Pick(ctrl)
		if edakit.Int(0, 2) == 0 {
			r, c = edakit.Int(1, n), edakit.Int(1, n)
		}
		return num(intro+edakit.Sprintf("A circuit with nodes 1 to %d has %s, and also %s (gm in mS). What is the nodal matrix entry $G_{%d%d}$ in mS?", n, describeResistors(rs), card, r, c),
			g[r-1][c-1], 3,
			edakit.Sprintf("Stamp the resistors, then add the VCCS entries on top. Entry (%d, %d) is %s mS.", r, c, f(g[r-1][c-1], 3)), nil)
	}
	a, b := np, cp
	if edakit.Int(0, 1) == 0 {
		a, b = cp, np
	}
	d := g[a-1][b-1] - g[b-1][a-1]
	return num(intro+edakit.Sprintf("A circuit with nodes 1 to %d has %s, and also %s (gm in mS). Resistors alone give a symmetric matrix. What is $G_{%d%d} - G_{%d%d}$ in mS?", n, describeResistors(rs), card, a, b, b, a),
		d, 3,
		edakit.Sprintf("Resistor stamps cancel in the difference, so only the VCCS remains: %s mS. Controlled sources make the MNA matrix unsymmetric, which is why SPICE uses a general sparse LU rather than Cholesky.", f(d, 3)), nil)
}

// vcvsStamps covers the VCVS branch row, unknown counting and an ideal amplifier stage.
func vcvsStamps(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 3, 4)
		nodes := edakit.Shuffle([]int{1, 2, 3, 4}[:n])
		np, cp := nodes[0], nodes[1]
		nm, cm := 0, 0
		if edakit.Int(0, 1) == 0 {
			nm = nodes[2]
		}
		if edakit.Int(0, 1) == 0 && n > 3 {
			cm = nodes[3]
		}
		mu := pickF(2, 5, 10, 20, 0.5, 4)
		withV := edakit.Int(0, 1) == 0
		nb := 1
		pre := ""
		if withV {
			nb = 2
			pre = edakit.Sprintf("V1 %d 0 1, then ", cp)
		}
		size := n + nb
		a := zeros(size)
		row := size - 1
		put := func(r, c int, v float64) {
			if r >= 0 && c >= 0 {
				a[r][c] += v
			}
		}
		put(row, np-1, 1)
		if nm > 0 {
			put(row, nm-1, -1)
			put(nm-1, row, -1)
		}
		put(np-1, row, 1)
		put(row, cp-1, -mu)
		if cm > 0 {
			put(row, cm-1, mu)
		}
		var r, c int
		switch edakit.Int(0, 2) {
		case 0:
			r, c = row+1, cp
		case 1:
			r, c = row+1, edakit.Int(1, n)
		default:
			r, c = edakit.Int(1, n), row+1
		}
		card := edakit.Sprintf("E1 %s %s %s %s %s", nodeName(np), nodeName(nm), nodeName(cp), nodeName(cm), f(mu, 2))
		return num(edakit.Sprintf("A circuit with nodes 1 to %d contains %s%s (E1 n+ n- nc+ nc- gain). Unknowns are $v_1, \\ldots, v_%d$ followed by the branch currents in card order, so E1's branch current is unknown %d and its constraint is row %d. What is the entry in row %d, column %d contributed by these cards?", n, pre, card, n, size, size, r, c),
			a[r-1][c-1], 2,
			edakit.Sprintf("The E1 row encodes $v_{n+} - v_{n-} - \\mu (v_{nc+} - v_{nc-}) = 0$: $+1$ at n+, $-1$ at n-, $-\\mu$ at nc+, $+\\mu$ at nc-. Its branch column has $+1$ at n+ and $-1$ at n-. Entry (%d, %d) is %s.", r, c, f(a[r-1][c-1], 2)), nil)
	case 2:
		n := edakit.Scaled(ctx, 3, 7)
		counts := map[string]int{"V": edakit.Int(1, 2), "E": edakit.Int(1, 2), "G": edakit.Int(0, 2), "I": edakit.Int(0, 2), "R": edakit.Int(n, 2*n), "C": edakit.Int(0, 3), "L": edakit.Int(0, 2)}
		u := n + counts["V"] + counts["E"] + counts["L"]
		return edakit.Numeric(edakit.Sprintf("A circuit has %d non-ground nodes and these elements: %d resistors, %d capacitors, %d inductors, %d independent voltage sources, %d independent current sources, %d VCVS (E) and %d VCCS (G). How many unknowns does its DC MNA system have, if every inductor carries a branch current?", n, counts["R"], counts["C"], counts["L"], counts["V"], counts["I"], counts["E"], counts["G"]),
			edakit.IntStr(u),
			edakit.Sprintf("Node voltages %d, plus one branch current per voltage source (%d), VCVS (%d) and inductor (%d). VCCS, resistors, capacitors and current sources stamp without new unknowns: %d.", n, counts["V"], counts["E"], counts["L"], u), nil)
	}
	vin := stepF(0.1, 0.1, 9)
	mu := pickF(2, 4, 5, 10, 20)
	rl := edakit.Pick(kohmChoices)
	split := edakit.Int(0, 1) == 0
	r3 := edakit.Pick(kohmChoices)
	load := rl
	stem := edakit.Sprintf("A deck has V1 1 0 %s, E1 2 0 1 0 %s and R1 2 0 %sk", f(vin, 2), f(mu, 2), f(rl, 2))
	if split {
		stem = edakit.Sprintf("A deck has V1 1 0 %s, E1 2 0 1 0 %s, R1 2 3 %sk and R2 3 0 %sk", f(vin, 2), f(mu, 2), f(rl, 2), f(r3, 2))
		load = rl + r3
	}
	vout := mu * vin
	ans := -vout / load
	return num(stem+". Using the SPICE sign convention (positive branch current flows into the + terminal of E1 and out of its - terminal), what is the branch current of E1 in mA?",
		ans, 3,
		edakit.Sprintf("E1 forces $v_2 = %s \\times %s = %s$ V. It delivers %s mA into %s k$\\Omega$ of load, so the current into its + terminal is %s mA.", f(mu, 2), f(vin, 2), f(vout, 3), f(vout/load, 3), f(load, 2), f(ans, 3)),
		map[string]string{"vout": f(vout, 3)})
}

// solveSmall asks for a node voltage of a two-node network.
func solveSmall(ctx generator.GeneratorContext) generator.Problem {
	r1, r12, r2 := edakit.Pick(kohmChoices), edakit.Pick(kohmChoices), edakit.Pick(kohmChoices)
	g1, g12, g2 := 1/r1, 1/r12, 1/r2
	k := kp(ctx)
	if k <= 2 {
		i := stepF(0.5, 0.5, 7)
		det := (g1+g12)*(g2+g12) - g12*g12
		v1 := i * (g2 + g12) / det
		v2 := i * g12 / det
		stem := edakit.Sprintf("A current source pushes %s mA into node 1. R1 = %s k$\\Omega$ connects node 1 to ground, R2 = %s k$\\Omega$ connects nodes 1 and 2, and R3 = %s k$\\Omega$ connects node 2 to ground. Writing $G v = b$ in mS, mA and V,", f(i, 2), f(r1, 2), f(r12, 2), f(r2, 2))
		facts := map[string]string{"det": f(det, 4), "v1": f(v1, 4), "v2": f(v2, 4)}
		expl := edakit.Sprintf("$G = \\begin{pmatrix} %s & -%s \\\\ -%s & %s \\end{pmatrix}$, $b = (%s, 0)$, $\\det G = %s$. By Cramer's rule $v_1 = %s$ V and $v_2 = %s$ V.", f(g1+g12, 4), f(g12, 4), f(g12, 4), f(g2+g12, 4), f(i, 2), f(det, 4), f(v1, 4), f(v2, 4))
		if k == 1 {
			return num(stem+" what is $v_1$ in volts?", v1, 3, expl, facts)
		}
		return num(stem+" what is $v_2$ in volts?", v2, 3, expl, facts)
	}
	vs := stepF(1, 0.5, 6)
	ix := stepF(-1, 0.5, 4)
	v2 := (vs*g12 + ix) / (g12 + g2)
	iv := -(vs*g12 - v2*g12)
	return num(edakit.Sprintf("V1 holds node 1 at %s V. R2 = %s k$\\Omega$ connects nodes 1 and 2, R3 = %s k$\\Omega$ connects node 2 to ground, and a current source pushes %s mA into node 2. The MNA unknowns are $v_1$, $v_2$ and the branch current of V1. What is $v_2$ in volts?", f(vs, 2), f(r12, 2), f(r2, 2), f(ix, 2)),
		v2, 3,
		edakit.Sprintf("The branch row fixes $v_1 = %s$. The node-2 row is $-g_2 v_1 + (g_2 + g_3) v_2 = I$, so $v_2 = (%s \\cdot %s + %s)/%s = %s$ V. (The V1 branch current is %s mA.)", f(vs, 2), f(g12, 4), f(vs, 2), f(ix, 2), f(g12+g2, 4), f(v2, 4), f(iv, 4)),
		map[string]string{"v2": f(v2, 4)})
}

// mnaStructure counts structural nonzeros and zero diagonals.
func mnaStructure(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 3, 6)
	rs := randomNetwork(n, edakit.Scaled(ctx, 1, 4))
	pos := map[[2]int]bool{}
	for _, r := range rs {
		pos[[2]int{r.a, r.a}] = true
		if r.b > 0 {
			pos[[2]int{r.b, r.b}] = true
			pos[[2]int{r.a, r.b}] = true
			pos[[2]int{r.b, r.a}] = true
		}
	}
	gnz := len(pos)
	desc := describeResistors(rs)
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(edakit.Sprintf("A circuit on nodes 1 to %d has %s. How many structurally nonzero entries does its nodal conductance matrix have (ground removed)?", n, desc),
			edakit.IntStr(gnz),
			edakit.Sprintf("Every node has a nonzero diagonal, and each distinct pair of non-ground nodes joined by a resistor adds two symmetric off-diagonal entries. Parallel resistors share a position. Total %d.", gnz), nil)
	case 2:
		nv := edakit.Int(1, 3)
		var parts []string
		extra := 0
		for k := 0; k < nv; k++ {
			p := edakit.Int(1, n)
			m := 0
			if edakit.Int(0, 1) == 0 {
				m = edakit.Int(1, n)
				if m == p {
					m = 0
				}
			}
			if m == 0 {
				extra += 2
			} else {
				extra += 4
			}
			parts = append(parts, edakit.Sprintf("V%d %s %s", k+1, nodeName(p), nodeName(m)))
		}
		tot := gnz + extra
		return edakit.Numeric(edakit.Sprintf("A circuit on nodes 1 to %d has %s, plus voltage sources %s (+ node, - node). How many structurally nonzero entries does the MNA matrix have?", n, desc, strings.Join(parts, ", ")),
			edakit.IntStr(tot),
			edakit.Sprintf("The node block has %d nonzeros. Each grounded voltage source adds a $+1$ in its column and in its row (2 entries); a floating one adds 4. Its own diagonal stays zero. Total %d.", gnz, tot), nil)
	}
	nv, nl, ne := edakit.Int(1, 3), edakit.Int(0, 3), edakit.Int(0, 2)
	nc := edakit.Int(0, 4)
	return edakit.Numeric(edakit.Sprintf("A circuit on nodes 1 to %d has %s, plus %d voltage sources, %d inductors, %d capacitors and %d VCVS. For the DC operating point SPICE treats each inductor as a branch current with zero voltage drop. How many zero entries lie on the diagonal of the MNA matrix?", n, desc, nv, nl, nc, ne),
		edakit.IntStr(nv+nl+ne),
		edakit.Sprintf("Node rows have positive conductance on the diagonal. Every branch row (voltage source, inductor at DC, VCVS) is a voltage constraint with no term in its own current, so it has a zero diagonal: %d. These zeros are why SPICE's LU must pivot or order with care.", nv+nl+ne), nil)
}
