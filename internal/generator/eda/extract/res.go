package extract

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// sheetResistance counts squares on a straight wire (KP 1), turns squares into ohms, possibly
// from resistivity and thickness (KP 2), and handles an L-shaped wire with a corner square
// (KP 3).
func sheetResistance(ctx generator.GeneratorContext) generator.Problem {
	w := edakit.Pick([]float64{0.02, 0.04, 0.05, 0.1, 0.2})
	switch kp(ctx) {
	case 1:
		n := 10 * edakit.Scaled(ctx, 3, 60)
		l := float64(n) * w
		return num(edakit.Sprintf("A straight wire on one metal layer is %s um long and %s um wide. How many squares does it contain?", f(l, 3), f(w, 3)),
			float64(n), 0, edakit.Sprintf("Squares are length over width: %s / %s = %d. Sheet resistance times squares gives the resistance, whatever the absolute size of each square.", f(l, 3), f(w, 3), n),
			map[string]string{"squares": edakit.IntStr(n)})
	case 2:
		rho := edakit.Pick([]float64{0.022, 0.025, 0.03, 0.035, 0.04})
		t := edakit.Pick([]float64{0.04, 0.05, 0.08, 0.1, 0.2})
		n := 10 * edakit.Scaled(ctx, 3, 60)
		l := float64(n) * w
		rs := rho / t
		r := rs * float64(n)
		return num(edakit.Sprintf("A metal layer has effective resistivity %s ohm-um and thickness %s um. A wire on it is %s um long and %s um wide. What is its resistance, in ohms?", f(rho, 3), f(t, 2), f(l, 3), f(w, 3)),
			r, 2, edakit.Sprintf("Sheet resistance is $\\rho/T$ = %s / %s = %s ohm/sq. The wire has %s / %s = %d squares, so R = %s x %d = %s ohm.", f(rho, 3), f(t, 2), f(rs, 4), f(l, 3), f(w, 3), n, f(rs, 4), n, f(r, 2)),
			map[string]string{"rs": f(rs, 4), "squares": edakit.IntStr(n)})
	}
	rs := edakit.Pick([]float64{0.2, 0.3, 0.4, 0.5, 0.8, 1.2})
	a := 10 * edakit.Scaled(ctx, 2, 40)
	b := 10 * edakit.Scaled(ctx, 2, 40)
	sq := float64(a) + float64(b) + 0.56
	r := rs * sq
	return num(edakit.Sprintf("An L-shaped wire of width %s um has a sheet resistance of %s ohm/sq. Excluding the corner square, its two arms are %s um and %s um long. Count the corner square as 0.56 squares. What is the resistance of the wire, in ohms?", f(w, 3), f(rs, 2), f(float64(a)*w, 3), f(float64(b)*w, 3)),
		r, 2, edakit.Sprintf("The arms hold %d and %d squares, and the corner adds 0.56, because current crowds around the inner corner and the outer corner carries little current. Total %s squares times %s ohm/sq = %s ohm.", a, b, f(sq, 2), f(rs, 2), f(r, 2)),
		map[string]string{"squares": f(sq, 2)})
}

// wireResistance computes a segment from resistivity and dimensions (KP 1), a two-layer route
// with a via array (KP 2), and the share of a route's resistance in its vias (KP 3).
func wireResistance(ctx generator.GeneratorContext) generator.Problem {
	rho := edakit.Pick([]float64{0.022, 0.025, 0.03, 0.035, 0.04})
	switch kp(ctx) {
	case 1:
		w := edakit.Pick([]float64{0.02, 0.03, 0.04, 0.05, 0.1})
		t := edakit.Pick([]float64{0.04, 0.05, 0.06, 0.08, 0.1})
		l := float64(5 * edakit.Scaled(ctx, 2, 40))
		r := rho * l / (w * t)
		return num(edakit.Sprintf("A copper wire has effective resistivity %s ohm-um, width %s um, thickness %s um and length %s um. What is its resistance, in ohms?", f(rho, 3), f(w, 3), f(t, 3), f(l, 0)),
			r, 1, edakit.Sprintf("$R = \\rho L/(WT)$ = %s x %s / (%s x %s) = %s ohm.", f(rho, 3), f(l, 0), f(w, 3), f(t, 3), f(r, 1)), nil)
	case 2:
		rs1 := edakit.Pick([]float64{0.5, 0.8, 1, 1.5, 2})
		rs2 := edakit.Pick([]float64{0.2, 0.3, 0.4, 0.5})
		w1, w2 := edakit.Pick([]float64{0.02, 0.04, 0.05}), edakit.Pick([]float64{0.04, 0.05, 0.1})
		n1, n2 := 10*edakit.Scaled(ctx, 2, 30), 10*edakit.Scaled(ctx, 2, 40)
		rv := float64(5 * edakit.Int(2, 10))
		cuts := edakit.Pick([]int{1, 2, 4})
		r1, r2 := rs1*float64(n1), rs2*float64(n2)
		rvia := rv / float64(cuts)
		tot := r1 + rvia + r2
		return num(edakit.Sprintf("A net runs %s um on M2 (sheet resistance %s ohm/sq, width %s um), climbs one via level through an array of %d identical cuts of %s ohm each, then runs %s um on M3 (sheet resistance %s ohm/sq, width %s um). What is the total resistance from end to end, in ohms?",
			f(float64(n1)*w1, 3), f(rs1, 2), f(w1, 3), cuts, f(rv, 0), f(float64(n2)*w2, 3), f(rs2, 2), f(w2, 3)),
			tot, 2, edakit.Sprintf("M2: %d squares x %s = %s ohm. Via array: %d cuts in parallel give %s / %d = %s ohm. M3: %d squares x %s = %s ohm. In series the total is %s ohm.", n1, f(rs1, 2), f(r1, 2), cuts, f(rv, 0), cuts, f(rvia, 2), n2, f(rs2, 2), f(r2, 2), f(tot, 2)),
			map[string]string{"r_m2": f(r1, 2), "r_via": f(rvia, 2), "r_m3": f(r2, 2)})
	}
	// A short route that climbs k via levels up and k down, each a single cut.
	rs := edakit.Pick([]float64{0.4, 0.5, 0.8, 1})
	w := edakit.Pick([]float64{0.02, 0.04, 0.05})
	n := 10 * edakit.Scaled(ctx, 2, 30)
	k := edakit.Int(1, 3)
	rv := float64(5 * edakit.Int(4, 16))
	rw := rs * float64(n)
	rvt := 2 * float64(k) * rv
	pct := 100 * rvt / (rw + rvt)
	return num(edakit.Sprintf("A short connection leaves a pin on M1, climbs %d via levels to a routing layer, runs %s um there (sheet resistance %s ohm/sq, width %s um), and descends %d via levels to the sink pin. Every via level is a single cut of %s ohm, and the M1 pin stubs are negligible. What percentage of the connection's resistance is in the vias?",
		k, f(float64(n)*w, 3), f(rs, 2), f(w, 3), k, f(rv, 0)),
		pct, 1, edakit.Sprintf("The wire is %d squares x %s = %s ohm. The %d single-cut vias add %s ohm. Vias are %s / %s = %s%% of the total. On short nets at advanced nodes the vias, not the wire, often dominate resistance, which is why routers prefer redundant vias and few layer changes.", n, f(rs, 2), f(rw, 1), 2*k, f(rvt, 0), f(rvt, 0), f(rw+rvt, 1), f(pct, 1)),
		map[string]string{"r_wire": f(rw, 1), "r_vias": f(rvt, 0)})
}

// temperature scales resistance with a linear temperature coefficient.
func temperature(ctx generator.GeneratorContext) generator.Problem {
	alpha := edakit.Pick([]float64{0.003, 0.0032, 0.0035, 0.0038, 0.004})
	t0 := 25.0
	switch kp(ctx) {
	case 1:
		r0 := float64(10 * edakit.Scaled(ctx, 5, 80))
		t := edakit.Pick([]float64{85, 105, 110, 125, 150})
		r := r0 * (1 + alpha*(t-t0))
		return num(edakit.Sprintf("A wire measures %s ohm at 25 C. Its temperature coefficient of resistance is %s per C. What is its resistance at %s C, in ohms?", f(r0, 0), f(alpha, 4), f(t, 0)),
			r, 2, edakit.Sprintf("$R(T) = R_0[1 + \\alpha (T - T_0)]$ = %s x (1 + %s x %s) = %s ohm.", f(r0, 0), f(alpha, 4), f(t-t0, 0), f(r, 2)), nil)
	case 2:
		r0 := float64(10 * edakit.Scaled(ctx, 5, 80))
		t := edakit.Pick([]float64{-40, 0})
		r := r0 * (1 + alpha*(t-t0))
		hot := 125.0
		rh := r0 * (1 + alpha*(hot-t0))
		return num(edakit.Sprintf("An extractor reports a net resistance of %s ohm at 125 C. The layer's temperature coefficient is %s per C, referenced to 25 C. What resistance should the extractor report for the same net at %s C, in ohms?", f(rh, 2), f(alpha, 4), f(t, 0)),
			r, 2, edakit.Sprintf("First recover the 25 C value: %s / (1 + %s x 100) = %s ohm. Then scale to %s C: %s x (1 + %s x (%s)) = %s ohm.", f(rh, 2), f(alpha, 4), f(r0, 2), f(t, 0), f(r0, 2), f(alpha, 4), f(t-t0, 0), f(r, 2)),
			map[string]string{"r25": f(r0, 2)})
	}
	dd := float64(edakit.Int(10, 40))
	dw := float64(edakit.Int(5, 40))
	t := edakit.Pick([]float64{105, 125})
	hot := dd + dw*(1+alpha*(t-t0))
	return num(edakit.Sprintf("At 25 C a net's Elmore delay splits into a driver term of %s ps (driver resistance times load) and a wire term of %s ps (wire resistance times downstream capacitance). Only the wire resistance is temperature dependent here, with coefficient %s per C; the driver term is held fixed. What is the net delay at %s C, in ps?", f(dd, 0), f(dw, 0), f(alpha, 4), f(t, 0)),
		hot, 2, edakit.Sprintf("Only the wire term scales: %s x (1 + %s x %s) = %s ps. Adding the fixed %s ps driver term gives %s ps. Long, wire-dominated nets are the ones whose delay is most temperature sensitive.", f(dw, 0), f(alpha, 4), f(t-t0, 0), f(hot-dd, 2), f(dd, 0), f(hot, 2)), nil)
}

// barrierScattering shrinks the copper cross-section by a liner and adds a size-effect term.
func barrierScattering(ctx generator.GeneratorContext) generator.Problem {
	w := edakit.Pick([]float64{0.02, 0.024, 0.03, 0.04, 0.05})
	t := edakit.Pick([]float64{0.04, 0.05, 0.06, 0.08})
	b := edakit.Pick([]float64{0.002, 0.003, 0.004})
	eff := (w - 2*b) * (t - b)
	frac := 100 * eff / (w * t)
	switch kp(ctx) {
	case 1:
		return num(edakit.Sprintf("A copper wire is drawn %s um wide and %s um thick. A resistive barrier liner of %s um coats both sidewalls and the bottom, and it carries no current. What percentage of the drawn cross-section is conducting copper?", f(w, 3), f(t, 3), f(b, 3)),
			frac, 1, edakit.Sprintf("The copper is $(W - 2b)(T - b)$ = %s x %s = %s um^2 out of %s um^2, which is %s%%. The liner does not shrink with the wire, so its share grows at every node.", f(w-2*b, 3), f(t-b, 3), f(eff, 6), f(w*t, 6), f(frac, 1)), nil)
	case 2:
		inc := 100 * (w*t/eff - 1)
		return num(edakit.Sprintf("A wire is drawn %s um wide and %s um thick, with a non-conducting barrier of %s um on both sidewalls and the bottom. By what percentage does the barrier increase the wire's resistance over a wire whose whole drawn cross-section were copper of the same resistivity?", f(w, 3), f(t, 3), f(b, 3)),
			inc, 1, edakit.Sprintf("Resistance scales as one over the conducting area. The ratio is %s / %s = %s, an increase of %s%%.", f(w*t, 6), f(eff, 6), f(w*t/eff, 4), f(inc, 1)), nil)
	}
	rho0 := 0.018
	k := edakit.Pick([]float64{0.01, 0.015, 0.02})
	l := float64(5 * edakit.Scaled(ctx, 2, 20))
	wc := w - 2*b
	rhoEff := rho0 * (1 + k/wc)
	r := rhoEff * l / eff
	return num(edakit.Sprintf("Use a simple size-effect model: copper resistivity is $\\rho = %s(1 + %s/W_c)$ ohm-um, where $W_c$ is the copper width in um. A wire is drawn %s um wide, %s um thick and %s um long, with a non-conducting barrier of %s um on both sidewalls and the bottom. What is its resistance, in ohms?", f(rho0, 3), f(k, 3), f(w, 3), f(t, 3), f(l, 0), f(b, 3)),
		r, 1, edakit.Sprintf("The copper is %s um wide and %s um thick. Surface and grain-boundary scattering raise resistivity to %s x (1 + %s / %s) = %s ohm-um. Then R = %s x %s / %s = %s ohm.", f(wc, 3), f(t-b, 3), f(rho0, 3), f(k, 3), f(wc, 3), f(rhoEff, 4), f(rhoEff, 4), f(l, 0), f(eff, 6), f(r, 1)),
		map[string]string{"rho_eff": f(rhoEff, 4)})
}
