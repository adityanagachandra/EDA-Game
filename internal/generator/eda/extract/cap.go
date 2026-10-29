package extract

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// parallelPlate computes area capacitance, its per-length value, and a scaling ratio.
func parallelPlate(ctx generator.GeneratorContext) generator.Problem {
	er := edakit.Pick([]float64{2.5, 2.7, 3.0, 3.9})
	w := edakit.Pick([]float64{0.1, 0.2, 0.4, 0.5, 1})
	h := edakit.Pick([]float64{0.05, 0.08, 0.1, 0.2})
	l := float64(10 * edakit.Scaled(ctx, 2, 30))
	c := eps0 * er * w * l / h
	switch kp(ctx) {
	case 1:
		return num(edakit.Sprintf("A wire %s um wide and %s um long runs %s um above a wide grounded plane, in a dielectric of relative permittivity %s. Using the parallel-plate formula with $\\varepsilon_0 = 8.854 \\times 10^{-3}$ fF/um, what is its area capacitance, in fF?", f(w, 2), f(l, 0), f(h, 2), f(er, 1)),
			c, 2, edakit.Sprintf("$C = \\varepsilon_0 \\varepsilon_r W L / h$ = 0.008854 x %s x %s x %s / %s = %s fF.", f(er, 1), f(w, 2), f(l, 0), f(h, 2), f(c, 2)), nil)
	case 2:
		per := 1000 * eps0 * er * w / h
		return num(edakit.Sprintf("A wire %s um wide runs %s um above a grounded plane in a dielectric with $\\varepsilon_r = %s$. Take $\\varepsilon_0 = 8.854 \\times 10^{-3}$ fF/um. What is its parallel-plate capacitance per unit length, in aF/um?", f(w, 2), f(h, 2), f(er, 1)),
			per, 1, edakit.Sprintf("Per unit length, $C/L = \\varepsilon_0 \\varepsilon_r W/h$ = 0.008854 x %s x %s / %s = %s fF/um, which is %s aF/um.", f(er, 1), f(w, 2), f(h, 2), f(per/1000, 4), f(per, 1)), nil)
	}
	sw := edakit.Pick([]float64{0.5, 0.7, 0.8})
	sh := edakit.Pick([]float64{0.5, 0.7, 0.8, 1.25})
	ser := edakit.Pick([]float64{1, 0.9, 0.85})
	ratio := sw * ser / sh
	return num(edakit.Sprintf("A process shrink scales a wire's width by %s and its height above the plane below by %s, and a new dielectric scales $\\varepsilon_r$ by %s. Length is unchanged. By what factor does the parallel-plate area capacitance change (new over old)?", f(sw, 2), f(sh, 2), f(ser, 2)),
		ratio, 3, edakit.Sprintf("Area capacitance goes as $\\varepsilon_r W / h$: %s x %s / %s = %s. Thinning the dielectric works against narrowing the wire, which is one reason capacitance per length has stayed roughly flat across nodes.", f(ser, 2), f(sw, 2), f(sh, 2), f(ratio, 3)), nil)
}

// fringe adds per-edge fringe capacitance to area capacitance.
func fringe(ctx generator.GeneratorContext) generator.Problem {
	ca := float64(10 * edakit.Int(2, 12)) // aF/um^2
	cf := float64(5 * edakit.Int(4, 16))  // aF/um per edge
	w := edakit.Pick([]float64{0.05, 0.1, 0.2, 0.4, 0.8})
	l := float64(10 * edakit.Scaled(ctx, 2, 40))
	area := ca * w * l / 1000
	fr := 2 * cf * l / 1000
	switch kp(ctx) {
	case 1:
		return num(edakit.Sprintf("A wire %s um wide and %s um long has area capacitance %s aF/um^2 to the plane below and fringe capacitance %s aF/um along each of its two edges. What is its total capacitance to ground, in fF?", f(w, 2), f(l, 0), f(ca, 0), f(cf, 0)),
			area+fr, 3, edakit.Sprintf("Area: %s x %s x %s = %s fF. Fringe: 2 x %s x %s = %s fF. Total %s fF.", f(ca, 0), f(w, 2), f(l, 0), f(area, 3), f(cf, 0), f(l, 0), f(fr, 3), f(area+fr, 3)),
			map[string]string{"area": f(area, 3), "fringe": f(fr, 3)})
	case 2:
		pct := 100 * fr / (area + fr)
		return num(edakit.Sprintf("A wire %s um wide has area capacitance %s aF/um^2 and fringe capacitance %s aF/um on each of its two edges. What percentage of its ground capacitance is fringe?", f(w, 2), f(ca, 0), f(cf, 0)),
			pct, 1, edakit.Sprintf("Per micron of length: area %s x %s = %s aF/um, fringe 2 x %s = %s aF/um. Fringe share = %s%%. The length cancels.", f(ca, 0), f(w, 2), f(ca*w, 2), f(cf, 0), f(2*cf, 0), f(pct, 1)), nil)
	}
	wstar := 2 * cf / ca
	return num(edakit.Sprintf("A layer has area capacitance %s aF/um^2 and fringe capacitance %s aF/um per edge, both independent of width in this model. Below what wire width, in um, does fringe capacitance exceed area capacitance?", f(ca, 0), f(cf, 0)),
		wstar, 3, edakit.Sprintf("Equal when $c_a W = 2 c_f$, so $W^* = 2c_f/c_a$ = %s / %s = %s um. Minimum-width wires at advanced nodes sit far below this, so their ground capacitance is mostly fringe and a parallel-plate estimate is badly low.", f(2*cf, 0), f(ca, 0), f(wstar, 3)), nil)
}

// coupling computes sidewall coupling, its spacing scaling, and coupling over a partial
// parallel run.
func coupling(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		er := edakit.Pick([]float64{2.5, 2.7, 3.0})
		t := edakit.Pick([]float64{0.04, 0.05, 0.08, 0.1})
		s := edakit.Pick([]float64{0.02, 0.03, 0.04, 0.05})
		l := float64(10 * edakit.Scaled(ctx, 2, 40))
		c := eps0 * er * t * l / s
		return num(edakit.Sprintf("Two parallel wires of thickness %s um face each other across a spacing of %s um for %s um, in a dielectric with $\\varepsilon_r = %s$. Treat the facing sidewalls as a parallel-plate capacitor with $\\varepsilon_0 = 8.854 \\times 10^{-3}$ fF/um. What is the coupling capacitance, in fF?", f(t, 2), f(s, 2), f(l, 0), f(er, 1)),
			c, 2, edakit.Sprintf("$C_c = \\varepsilon_0 \\varepsilon_r T L / s$ = 0.008854 x %s x %s x %s / %s = %s fF. Fringe fields add to this, so a field solver would report somewhat more.", f(er, 1), f(t, 2), f(l, 0), f(s, 2), f(c, 2)), nil)
	case 2:
		c0 := float64(5 * edakit.Int(10, 30)) // aF/um
		s0 := edakit.Pick([]float64{0.04, 0.05, 0.06})
		s1 := s0 * edakit.Pick([]float64{0.5, 0.75, 1.5, 2})
		l := float64(10 * edakit.Scaled(ctx, 3, 40))
		c := c0 * (s0 / s1) * l / 1000
		return num(edakit.Sprintf("A characterization table gives %s aF/um of coupling between two minimum-width wires at spacing %s um. Assume coupling per length scales inversely with spacing. Two such wires run in parallel for %s um at spacing %s um. What is their coupling capacitance, in fF?", f(c0, 0), f(s0, 3), f(l, 0), f(s1, 3)),
			c, 3, edakit.Sprintf("Per length: %s x %s / %s = %s aF/um. Over %s um: %s fF.", f(c0, 0), f(s0, 3), f(s1, 3), f(c0*s0/s1, 2), f(l, 0), f(c, 3)), nil)
	}
	cc := float64(5 * edakit.Int(10, 30))
	vs := 0
	ve, as, ae := 200, 40, 150 // known-good fallback instance
	for try := 0; try < generator.MaxDrawAttempts; try++ {
		v := 10 * edakit.Scaled(ctx, 8, 30)
		a := 10 * edakit.Int(-5, v/10)
		e := a + 10*edakit.Int(3, 25)
		if e > vs+10 && a < v-10 {
			ve, as, ae = v, a, e
			break
		}
	}
	lo, hi := as, ae
	if lo < vs {
		lo = vs
	}
	if hi > ve {
		hi = ve
	}
	ov := float64(hi - lo)
	c := cc * ov / 1000
	return num(edakit.Sprintf("A victim wire runs along a track from x = %d um to x = %d um. On the adjacent track at minimum spacing, an aggressor runs from x = %d um to x = %d um. Coupling at that spacing is %s aF/um of parallel run, and end effects are ignored. What is the coupling capacitance between them, in fF?", vs, ve, as, ae, f(cc, 0)),
		c, 3, edakit.Sprintf("Only the overlap couples: from x = %d to x = %d, which is %s um. %s aF/um x %s um = %s fF. Extractors compute exactly this parallel run length from the geometry before looking up a per-length value.", lo, hi, f(ov, 0), f(cc, 0), f(ov, 0), f(c, 3)),
		map[string]string{"overlap": f(ov, 0)})
}

// capMatrix reads coupling and ground capacitance out of a Maxwell capacitance matrix, or
// builds a diagonal entry.
func capMatrix(ctx generator.GeneratorContext) generator.Problem {
	cg := []float64{float64(edakit.Int(2, 12)), float64(edakit.Int(2, 12)), float64(edakit.Int(2, 12))}
	c12, c13, c23 := float64(edakit.Int(1, 8)), float64(edakit.Int(0, 6)), float64(edakit.Int(1, 8))
	cc := [3][3]float64{{0, c12, c13}, {c12, 0, c23}, {c13, c23, 0}}
	var m [3][3]float64
	for i := 0; i < 3; i++ {
		m[i][i] = cg[i]
		for j := 0; j < 3; j++ {
			if i != j {
				m[i][i] += cc[i][j]
				m[i][j] = -cc[i][j]
			}
		}
	}
	var rows []string
	for i := 0; i < 3; i++ {
		rows = append(rows, edakit.Sprintf("row %d: [%s, %s, %s]", i+1, f(m[i][0], 0), f(m[i][1], 0), f(m[i][2], 0)))
	}
	mat := "A field solver returns the Maxwell capacitance matrix of three conductors over a ground plane, in fF: " + strings.Join(rows, "; ") + ". "
	i := edakit.Int(0, 2)
	j := (i + edakit.Int(1, 2)) % 3
	switch kp(ctx) {
	case 1:
		return num(mat+edakit.Sprintf("What is the coupling capacitance between conductors %d and %d, in fF?", i+1, j+1),
			cc[i][j], 0, edakit.Sprintf("Off-diagonal entries are the negatives of the coupling capacitances: $C_{c} = -C_{%d%d}$ = %s fF.", i+1, j+1, f(cc[i][j], 0)), nil)
	case 2:
		return num(mat+edakit.Sprintf("What is the capacitance from conductor %d to the ground plane, in fF?", i+1),
			cg[i], 0, edakit.Sprintf("The ground capacitance is the row sum: %s + (%s) + (%s) = %s fF. The diagonal is the total capacitance, so subtracting the couplings leaves the part to ground.", f(m[i][0], 0), f(m[i][1], 0), f(m[i][2], 0), f(cg[i], 0)), nil)
	}
	q := edakit.Sprintf("Three conductors over a ground plane have ground capacitances %s, %s and %s fF. The couplings are %s fF between 1 and 2, %s fF between 1 and 3, and %s fF between 2 and 3. What is the diagonal entry $C_{%d%d}$ of the Maxwell capacitance matrix, in fF?", f(cg[0], 0), f(cg[1], 0), f(cg[2], 0), f(c12, 0), f(c13, 0), f(c23, 0), i+1, i+1)
	return num(q, m[i][i], 0, edakit.Sprintf("$C_{ii}$ is the charge on conductor i when it is at 1 V and all others at 0 V: its ground capacitance plus every coupling, %s fF.", f(m[i][i], 0)), nil)
}

// totalNet sums ground, coupling and pin capacitance of one net.
func totalNet(ctx generator.GeneratorContext) generator.Problem {
	l := float64(10 * edakit.Scaled(ctx, 3, 50))
	cg := float64(edakit.Int(40, 120)) // aF/um to ground (area + fringe)
	cl := float64(edakit.Int(20, 80))  // aF/um coupling left
	cr := float64(edakit.Int(20, 80))
	np := edakit.Scaled(ctx, 1, 4)
	pins := make([]float64, np)
	var ps []string
	psum := 0.0
	for k := range pins {
		pins[k] = float64(edakit.Int(5, 30)) / 10
		psum += pins[k]
		ps = append(ps, f(pins[k], 1))
	}
	wire := l * (cg + cl + cr) / 1000
	tot := wire + psum
	stem := edakit.Sprintf("A net is %s um of wire with %s aF/um to ground (area plus fringe), %s aF/um of coupling to the neighbour on one side and %s aF/um to the neighbour on the other. It drives %s of %s fF.", f(l, 0), f(cg, 0), f(cl, 0), f(cr, 0), plural(np, "sink pin"), strings.Join(ps, ", "))
	switch kp(ctx) {
	case 1:
		return num(stem+" What is the net's total capacitance with all coupling counted to quiet neighbours, in fF?", tot, 2,
			edakit.Sprintf("Wire: %s x (%s + %s + %s) aF = %s fF. Pins: %s fF. Total %s fF.", f(l, 0), f(cg, 0), f(cl, 0), f(cr, 0), f(wire, 2), f(psum, 1), f(tot, 2)),
			map[string]string{"wire": f(wire, 2), "pins": f(psum, 1)})
	case 2:
		cc := l * (cl + cr) / 1000
		pct := 100 * cc / tot
		return num(stem+" What percentage of the net's total capacitance is coupling to the neighbours?", pct, 1,
			edakit.Sprintf("Coupling: %s x %s aF = %s fF. Total: %s fF. Share %s%%. A high share means the net's delay depends strongly on what its neighbours are doing.", f(l, 0), f(cl+cr, 0), f(cc, 2), f(tot, 2), f(pct, 1)), nil)
	}
	gnd := l*cg/1000 + psum
	return num(stem+" A timer that ignores signal integrity grounds every coupling capacitor. A careless flow instead drops coupling entirely and keeps only the ground and pin capacitance. How many fF does the careless flow lose?", tot-gnd, 2,
		edakit.Sprintf("Grounded coupling keeps %s fF; dropping it keeps %s fF. The difference is the coupling, %s fF. Coupling must stay in the netlist either as grounded capacitance or as coupling elements for SI analysis.", f(tot, 2), f(gnd, 2), f(tot-gnd, 2)), nil)
}
