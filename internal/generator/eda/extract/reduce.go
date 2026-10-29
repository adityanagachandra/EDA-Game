package extract

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// rcReductionCalc offers four reduced versions of a two-segment chain and asks which one
// preserves both total capacitance and the Elmore delay at the sink.
func rcReductionCalc(ctx generator.GeneratorContext) generator.Problem {
	r1, r2 := 250.0, 400.0 // known-good fallback with r1 != r2
	for try := 0; try < generator.MaxDrawAttempts; try++ {
		a, b := float64(50*edakit.Int(2, 12)), float64(50*edakit.Int(2, 12))
		if a != b {
			r1, r2 = a, b
			break
		}
	}
	c1, c2 := float64(edakit.Int(2, 12)), float64(edakit.Int(2, 12))
	rt := r1 + r2
	opt := func(cd, cs float64) string {
		return edakit.Sprintf("R = %s ohm, %s fF at the driver end, %s fF at the sink", f(rt, 0), f(cd, 2), f(cs, 2))
	}
	good := opt(c1*r2/rt, c2+c1*r1/rt)
	d := []string{opt(0, c1+c2), opt(0, c2), opt(c1/2, c2+c1/2), opt(c1*r1/rt, c2+c1*r2/rt)}
	q := edakit.Sprintf("An RC chain runs from a driver through %s ohm to internal node n1 (%s fF to ground), then through %s ohm to the sink (%s fF). A reducer eliminates n1 and keeps a single %s ohm resistor. Which reduced network preserves both the total capacitance and the Elmore delay at the sink?", f(r1, 0), f(c1, 0), f(r2, 0), f(c2, 0), f(rt, 0))
	el := r1*(c1+c2) + r2*c2
	return edakit.Choice(q, good, edakit.Sprintf("The original Elmore delay at the sink is %s x %s + %s x %s = %s ohm-fF. Splitting $C_1$ in proportion to the conductances, $C_1 R_2/(R_1+R_2)$ to the driver end and $C_1 R_1/(R_1+R_2)$ to the sink, keeps the total and gives (R1 + R2) x sink capacitance = the same %s ohm-fF. Capacitance at the driver end adds no wire delay.", f(r1, 0), f(c1+c2, 0), f(r2, 0), f(c2, 0), f(el, 0), f(el, 0)), d...)
}

var rcReductionItems = []edakit.Item{
	{KP: 1, Q: "A detailed extraction of a block has 40 million RC nodes. Why does the flow reduce it before timing?",
		A: "Delay calculation and SI analysis cost grows with network size, and most nodes add accuracy no timer can use", D: []string{"SPEF cannot store more than a million nodes", "Reduction makes the parasitics more accurate", "Timers ignore resistors anyway"},
		Why: "Reduction keeps the electrical behaviour seen at the pins (moments, Elmore delay, total capacitance) with far fewer nodes, so timing runs faster with negligible error."},
	{KP: 1, Q: "Which quantities should a timing-oriented RC reduction preserve at minimum?",
		A: "Total capacitance of the net and the Elmore delay (first moment) to each sink", D: []string{"The number of resistors", "The node names of every internal node", "The ground capacitance of each segment separately"},
		Why: "Total capacitance sets the driver load; the first moment to each sink sets the wire delay. Preserving both keeps the delay calculator's answers stable."},
	{KP: 1, Q: "What makes RC reduction different from simply deleting small capacitors?",
		A: "It moves charge-storage and resistance into the remaining nodes so the port behaviour is kept", D: []string{"It deletes only the large capacitors", "It replaces resistors by ideal wires", "It removes coupling capacitance first"},
		Why: "Deleting elements changes total load and delays. Node elimination redistributes the eliminated node's capacitance and builds equivalent resistors among its neighbours."},
	{KP: 1, Q: "How does Krylov-projection reduction, such as PRIMA, differ from node elimination methods such as TICER?",
		A: "It projects the whole network onto a small subspace matching several moments, rather than removing nodes one at a time", D: []string{"It only works for single-segment wires", "It cannot preserve passivity", "It works only in the time domain"},
		Why: "Projection methods match multiple moments at once and PRIMA guarantees passivity. Node elimination is local, keeps a physical RC netlist and is easy to apply incrementally."},
	{KP: 1, Q: "A reduced model is accurate to its first moment only. Which quantity can it get noticeably wrong?",
		A: "The waveform shape and slew at the far end, which depend on higher moments", D: []string{"The total capacitance", "The Elmore delay", "The DC resistance from driver to sink"},
		Why: "Matching only the first moment fixes the Elmore delay. Slew and the 50% delay of resistive-shielded nets depend on the second and higher moments."},
	{KP: 1, Q: "Why must a reduced RC network be passive?",
		A: "A non-passive model can generate energy and make a downstream simulation unstable", D: []string{"Passivity makes the file smaller", "Timers reject any network with resistors", "Passivity guarantees zero error"},
		Why: "RC interconnect is passive. A reduced model that is not can produce growing responses when combined with drivers, which is why PRIMA's passivity guarantee matters."},
	{KP: 3, Q: "TICER eliminates nodes whose time constant is small compared with the signal's rise time. Why that criterion?",
		A: "A node that settles much faster than the signal changes behaves quasi-statically, so removing it barely changes the waveform", D: []string{"Small time constants mean the node has no capacitance", "Large time constants are cheaper to eliminate", "The criterion preserves every moment exactly"},
		Why: "A fast node tracks its neighbours almost instantly; folding it into them is accurate. Slow nodes shape the response and are kept."},
	{KP: 3, Q: "Eliminating a node with six neighbours by a star-mesh transform creates how many new resistors among the neighbours, before merging with existing ones?",
		A: "15, one for every pair of neighbours", D: []string{"6, one per neighbour", "5, a chain through the neighbours", "36, one per ordered pair"},
		Why: "Star-mesh connects every pair of neighbours: 6 choose 2 = 15. High-degree nodes create fill-in, so reducers eliminate low-degree nodes first, as in sparse elimination ordering."},
	{KP: 3, Q: "A reducer eliminates a node on a coupled net that carries a coupling capacitor to an aggressor. What must happen to the coupling capacitor?",
		A: "It is split among the neighbouring nodes like a ground capacitor, keeping its other terminal on the aggressor", D: []string{"It is deleted", "It is converted to a ground capacitor", "It is doubled to account for Miller effect"},
		Why: "Coupling is redistributed in the same proportions as ground capacitance, so SI analysis still sees the full coupling to the right aggressor."},
	{KP: 3, Q: "After reduction, a sink's Elmore delay matches the detailed network but its 50% delay differs by 6%. What is the likely reason?",
		A: "The reduction preserved the first moment but not the higher moments that set the waveform", D: []string{"The total capacitance changed by 6%", "The reducer renamed the sink", "Elmore delay always equals 50% delay"},
		Why: "Elmore delay is only the first moment. Waveform-sensitive metrics need second and third moments matched, or a less aggressive threshold."},
	{KP: 3, Q: "Why do incremental flows prefer node elimination over global projection for reducing ECO-touched nets?",
		A: "Elimination is local, so only the changed nets need re-reduction and the result stays a readable RC netlist", D: []string{"Projection cannot handle more than one net", "Elimination preserves all moments exactly", "Projection requires a field solver"},
		Why: "Local methods reduce each net independently, which fits incremental extraction after an ECO, and the output is still R and C elements in SPEF."},
	{KP: 3, Q: "A reduced net preserves total capacitance but moves 30% of it from the far end to the driver end. Which timing number does this distort?",
		A: "The wire delay to the far sink falls, because less capacitance sits behind the wire resistance", D: []string{"The driver's load rises", "The total capacitance falls", "Nothing, because total capacitance is preserved"},
		Why: "Capacitance behind more resistance costs more delay. Preserving the total is necessary but not sufficient; the Elmore delays must also be kept."},
}

type rcNode struct {
	parent int
	r, c   float64
}

// pathR returns the resistance shared by the root paths of nodes i and k.
func pathR(t []rcNode, i, k int) float64 {
	anc := map[int]bool{}
	for x := i; x > 0; x = t[x].parent {
		anc[x] = true
	}
	s := 0.0
	for x := k; x > 0; x = t[x].parent {
		if anc[x] {
			s += t[x].r
		}
	}
	return s
}

func treeMoments(t []rcNode) ([]float64, []float64) {
	n := len(t)
	m1 := make([]float64, n)
	m2 := make([]float64, n)
	for i := 1; i < n; i++ {
		for k := 1; k < n; k++ {
			m1[i] += pathR(t, i, k) * t[k].c
		}
	}
	for i := 1; i < n; i++ {
		for k := 1; k < n; k++ {
			m2[i] += pathR(t, i, k) * t[k].c * m1[k]
		}
	}
	return m1, m2
}

func describeTree(t []rcNode) string {
	var parts []string
	for i := 1; i < len(t); i++ {
		from := "the driver"
		if t[i].parent > 0 {
			from = edakit.Sprintf("node %d", t[i].parent)
		}
		parts = append(parts, edakit.Sprintf("%s kohm from %s to node %d (%s fF)", f(t[i].r, 2), from, i, f(t[i].c, 0)))
	}
	return strings.Join(parts, ", ")
}

// moments computes the first and second moments of an RC tree.
func moments(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		t := []rcNode{{}, {0, float64(edakit.Int(1, 5)) / 10, float64(edakit.Int(2, 10))}}
		t = append(t, rcNode{1, float64(edakit.Int(1, 5)) / 10, float64(edakit.Int(2, 10))})
		t = append(t, rcNode{1, float64(edakit.Int(1, 5)) / 10, float64(edakit.Int(2, 10))})
		if edakit.Scaled(ctx, 0, 1) == 1 {
			t = append(t, rcNode{2, float64(edakit.Int(1, 5)) / 10, float64(edakit.Int(2, 10))})
		}
		m1, _ := treeMoments(t)
		s := edakit.Int(2, len(t)-1)
		return num(edakit.Sprintf("An RC tree has resistors %s, each node with the stated capacitance to ground. One kohm times one fF is one ps. What is the first moment of the impulse response at node %d (its Elmore delay), in ps?", describeTree(t), s),
			m1[s], 2, edakit.Sprintf("$m_1(i) = \\sum_k R_{ik} C_k$, where $R_{ik}$ is the resistance shared by the paths to i and k. At node %d this is %s ps.", s, f(m1[s], 2)), nil)
	case 2:
		n := edakit.Scaled(ctx, 2, 3)
		t := []rcNode{{}}
		for i := 1; i <= n; i++ {
			t = append(t, rcNode{i - 1, float64(edakit.Int(1, 5)) / 10, float64(edakit.Int(2, 10))})
		}
		m1, m2 := treeMoments(t)
		var m1s []string
		for i := 1; i <= n; i++ {
			m1s = append(m1s, f(m1[i], 2))
		}
		return num(edakit.Sprintf("An RC chain has resistors %s, each node with the stated capacitance to ground; kohm times fF is ps. The first moments at nodes 1 to %d are %s ps. Using $m_2(i) = \\sum_k R_{ik} C_k m_1(k)$, what is the second moment at node %d, in ps^2?", describeTree(t), n, strings.Join(m1s, ", "), n),
			m2[n], 3, edakit.Sprintf("Each capacitor contributes its shared resistance times its capacitance times its own first moment. The sum at node %d is %s ps^2. The second moment adds information about the spread of the response, which the Elmore delay alone lacks.", n, f(m2[n], 3)), nil)
	}
	n := edakit.Scaled(ctx, 2, 4)
	t := []rcNode{{}}
	rt, ct := 0.0, 0.0
	for i := 1; i <= n; i++ {
		t = append(t, rcNode{i - 1, float64(edakit.Int(1, 5)) / 10, float64(edakit.Int(2, 10))})
		rt += t[i].r
		ct += t[i].c
	}
	m1, _ := treeMoments(t)
	pct := 100 * (rt*ct - m1[n]) / m1[n]
	return num(edakit.Sprintf("An RC chain has resistors %s; kohm times fF is ps. A crude model lumps it into one resistor equal to the total resistance driving one capacitor equal to the total capacitance. By what percentage does that model overestimate the first moment at node %d?", describeTree(t), n),
		pct, 1, edakit.Sprintf("Exact: $m_1$ = %s ps. Lumped: %s kohm x %s fF = %s ps. Overestimate %s%%. Lumping puts every capacitor behind all the resistance, which only the last one actually sees.", f(m1[n], 2), f(rt, 2), f(ct, 0), f(rt*ct, 2), f(pct, 1)), nil)
}

// nodeElimination applies the star-mesh transform used by TICER.
func nodeElimination(ctx generator.GeneratorContext) generator.Problem {
	k := 3
	if ctx.Difficulty > 0.8 {
		k = edakit.Int(3, 4)
	}
	rs := make([]float64, k)
	gs := 0.0
	var parts []string
	for i := range rs {
		rs[i] = float64(10 * edakit.Int(2, 20))
		gs += 1 / rs[i]
		parts = append(parts, edakit.Sprintf("%s ohm to node %c", f(rs[i], 0), 'a'+i))
	}
	c0 := float64(edakit.Int(2, 20)) / 2
	stem := edakit.Sprintf("Internal node x has %s fF to ground and resistors %s. A reducer eliminates x with the star-mesh transform, which connects every pair of former neighbours i, j by $R_{ij} = R_i R_j \\sum_k 1/R_k$ and splits x's capacitance among them in proportion to their conductances.", f(c0, 1), strings.Join(parts, ", "))
	switch kp(ctx) {
	case 1:
		i := edakit.Int(0, k-1)
		j := (i + edakit.Int(1, k-1)) % k
		rij := rs[i] * rs[j] * gs
		return num(stem+edakit.Sprintf(" What is the new resistance between nodes %c and %c, in ohms?", 'a'+i, 'a'+j), rij, 2,
			edakit.Sprintf("$\\sum 1/R_k$ = %s S. $R_{%c%c}$ = %s x %s x %s = %s ohm.", f(gs, 5), 'a'+i, 'a'+j, f(rs[i], 0), f(rs[j], 0), f(gs, 5), f(rij, 2)), nil)
	case 2:
		i := edakit.Int(0, k-1)
		share := c0 * (1 / rs[i]) / gs
		return num(stem+edakit.Sprintf(" How much capacitance is added to node %c, in fF?", 'a'+i), share, 3,
			edakit.Sprintf("Node %c's conductance share is (1/%s) / %s = %s, so it receives %s x %s = %s fF. The shares sum to one, so total capacitance is preserved.", 'a'+i, f(rs[i], 0), f(gs, 5), f(1/rs[i]/gs, 4), f(c0, 1), f(1/rs[i]/gs, 4), f(share, 3)), nil)
	}
	// TICER threshold: count nodes whose time constant C / sum(g) is below a threshold.
	m := edakit.Scaled(ctx, 4, 7)
	thr := float64(edakit.Int(2, 6)) // ps
	cnt := 0
	var ds []string
	for i := 0; i < m; i++ {
		c := float64(edakit.Int(1, 20)) / 2
		gsum := 0.0
		var rr []string
		deg := edakit.Int(2, 3)
		for d := 0; d < deg; d++ {
			r := float64(100 * edakit.Int(1, 20))
			gsum += 1 / r
			rr = append(rr, f(r, 0))
		}
		tau := c / gsum / 1000
		if tau < thr {
			cnt++
		}
		ds = append(ds, edakit.Sprintf("node %d: %s fF, resistors %s ohm", i+1, f(c, 1), strings.Join(rr, ", ")))
	}
	return num(edakit.Sprintf("TICER marks a node for elimination when its time constant $\\tau = C / \\sum_k (1/R_k)$ is below a threshold, here %s ps (ohm times fF is fs). The candidate nodes are %s. How many nodes are marked?", f(thr, 0), strings.Join(ds, "; ")),
		float64(cnt), 0, edakit.Sprintf("Compute each node's capacitance times the parallel combination of its resistors, convert fs to ps, and compare with %s ps: %d nodes fall below it. These are the nodes that settle fast relative to the signal, so removing them barely changes the waveform.", f(thr, 0), cnt), nil)
}

// elmorePreserving eliminates the first internal node of a chain.
func elmorePreserving(ctx generator.GeneratorContext) generator.Problem {
	r1 := float64(edakit.Int(1, 8)) / 10
	r2 := float64(edakit.Int(1, 8)) / 10
	c1 := float64(edakit.Int(2, 12))
	c2 := float64(edakit.Int(2, 12))
	down := edakit.Scaled(ctx, 0, 1) == 1
	r3, c3 := 0.0, 0.0
	desc := edakit.Sprintf("An RC chain runs from a driver through %s kohm to node n1 (%s fF), then %s kohm to node n2 (%s fF)", f(r1, 1), f(c1, 0), f(r2, 1), f(c2, 0))
	sink := "n2"
	if down {
		r3 = float64(edakit.Int(1, 8)) / 10
		c3 = float64(edakit.Int(2, 12))
		desc += edakit.Sprintf(", then %s kohm to the sink n3 (%s fF)", f(r3, 1), f(c3, 0))
		sink = "n3"
	}
	desc += ". kohm times fF is ps. A reducer eliminates n1: it replaces the two resistors around n1 with one of " + f(r1+r2, 1) + " kohm and splits n1's capacitance between the driver end and n2 in proportion to the conductances of the two removed resistors."
	before := r1*(c1+c2+c3) + r2*(c2+c3) + r3*c3
	switch kp(ctx) {
	case 1:
		add := c1 * r1 / (r1 + r2)
		return num(desc+" How much capacitance is added to n2, in fF?", add, 3,
			edakit.Sprintf("n2's share is $(1/R_2)/(1/R_1 + 1/R_2) = R_1/(R_1+R_2)$ = %s / %s, so n2 gains %s fF; the driver end gets the remaining %s fF.", f(r1, 1), f(r1+r2, 1), f(add, 3), f(c1-add, 3)), nil)
	case 2:
		return num(desc+edakit.Sprintf(" What is the Elmore delay at %s after the elimination, in ps?", sink), before, 3,
			edakit.Sprintf("After elimination, %s kohm drives n2's new capacitance plus everything downstream, and the driver-end share adds nothing. The result equals the original Elmore delay, %s ps: this split preserves the first moment exactly at every node downstream of n1.", f(r1+r2, 1), f(before, 3)), nil)
	}
	errv := r2 * c1
	return num(desc+edakit.Sprintf(" A naive reducer instead moves all of n1's capacitance onto n2. By how many ps does that overestimate the Elmore delay at %s?", sink), errv, 3,
		edakit.Sprintf("Naive lumping puts $C_1$ behind $R_1 + R_2$ instead of behind $R_1$ only, adding $R_2 C_1$ = %s x %s = %s ps of delay that does not exist.", f(r2, 1), f(c1, 0), f(errv, 3)), nil)
}

type ser [4]float64

func (a ser) mul(b ser) ser {
	var o ser
	for i := 0; i < 4; i++ {
		for j := 0; i+j < 4; j++ {
			o[i+j] += a[i] * b[j]
		}
	}
	return o
}

// inv inverts a series with a nonzero constant term.
func (a ser) inv() ser {
	var o ser
	o[0] = 1 / a[0]
	for n := 1; n < 4; n++ {
		s := 0.0
		for k := 1; k <= n; k++ {
			s += a[k] * o[n-k]
		}
		o[n] = -s / a[0]
	}
	return o
}

// piFromMoments returns the near capacitance, far capacitance and resistance of the
// O'Brien-Savarino pi model from driving-point admittance moments.
func piFromMoments(y1, y2, y3 float64) (float64, float64, float64) {
	cf := y2 * y2 / y3
	return y1 - cf, cf, -y3 * y3 / (y2 * y2 * y2)
}

// piModel builds the moment-matched pi model of a uniform line or a two-segment chain.
func piModel(ctx generator.GeneratorContext) generator.Problem {
	if kp(ctx) == 1 {
		r := float64(50 * edakit.Int(2, 40))
		c := float64(6 * edakit.Int(2, 20))
		which := edakit.Int(1, 3)
		q := edakit.Sprintf("A uniform distributed RC line has total resistance %s ohm and total capacitance %s fF. Its driving-point admittance has moments $y_1 = C$, $y_2 = -RC^2/3$ and $y_3 = 2R^2C^3/15$. The pi model matching them has $C_{far} = y_2^2/y_3$, $C_{near} = y_1 - C_{far}$ and $R_{\\pi} = -y_3^2/y_2^3$. ", f(r, 0), f(c, 0))
		switch which {
		case 1:
			return num(q+"What is $C_{near}$, in fF?", c/6, 2, edakit.Sprintf("$C_{far} = (R^2C^4/9)/(2R^2C^3/15) = 5C/6$, so $C_{near} = C/6$ = %s fF.", f(c/6, 2)), nil)
		case 2:
			return num(q+"What is $C_{far}$, in fF?", 5*c/6, 2, edakit.Sprintf("$C_{far} = (R^2C^4/9)/(2R^2C^3/15) = 5C/6$ = %s fF.", f(5*c/6, 2)), nil)
		}
		return num(q+"What is $R_{\\pi}$, in ohms?", 12*r/25, 2, edakit.Sprintf("$R_\\pi = (4R^4C^6/225)/(R^3C^6/27) = 12R/25$ = %s ohm.", f(12*r/25, 2)), nil)
	}
	c0 := float64(edakit.Int(1, 10))
	r1, c1 := float64(edakit.Int(1, 8))/10, float64(edakit.Int(2, 20))
	r2, c2 := float64(edakit.Int(1, 8))/10, float64(edakit.Int(2, 20))
	// Y2 = s C2; Y1 = s C1 + Y2/(1 + R2 Y2); Y0 = s C0 + Y1/(1 + R1 Y1).
	stage := func(y ser, r float64) ser {
		den := ser{1, r * y[1], r * y[2], r * y[3]}
		return y.mul(den.inv())
	}
	y := ser{0, c2, 0, 0}
	y = stage(y, r2)
	y[1] += c1
	y = stage(y, r1)
	y[1] += c0
	cn, cf, rp := piFromMoments(y[1], y[2], y[3])
	q := edakit.Sprintf("A driver sees %s fF at its output pin, then %s kohm to node 1 (%s fF), then %s kohm to node 2 (%s fF). Expanding its driving-point admittance as $Y(s) = y_1 s + y_2 s^2 + y_3 s^3 + \\dots$ gives $y_1$ = %s fF, $y_2$ = %s kohm fF^2 and $y_3$ = %s kohm^2 fF^3. The O'Brien-Savarino pi model has $C_{far} = y_2^2/y_3$, $C_{near} = y_1 - C_{far}$ and $R_{\\pi} = -y_3^2/y_2^3$. ",
		f(c0, 0), f(r1, 1), f(c1, 0), f(r2, 1), f(c2, 0), f(y[1], 0), f(y[2], 4), f(y[3], 4))
	if kp(ctx) == 2 {
		return num(q+"What is $C_{far}$, in fF?", cf, 2, edakit.Sprintf("$C_{far}$ = (%s)^2 / %s = %s fF, and $C_{near}$ = %s - %s = %s fF.", f(y[2], 4), f(y[3], 4), f(cf, 2), f(y[1], 0), f(cf, 2), f(cn, 2)), nil)
	}
	return num(q+"What is $R_{\\pi}$, in kohm?", rp, 3, edakit.Sprintf("$R_\\pi = -y_3^2/y_2^3$ = -(%s)^2 / (%s)^3 = %s kohm. The driver delay calculator uses this pi load, with $C_{near}$ = %s fF and $C_{far}$ = %s fF, to compute an effective capacitance.", f(y[3], 4), f(y[2], 4), f(rp, 3), f(cn, 2), f(cf, 2)), nil)
}
