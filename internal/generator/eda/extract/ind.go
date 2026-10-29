package extract

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// partialInductance combines partial self and mutual inductances into loop inductance.
func partialInductance(ctx generator.GeneratorContext) generator.Problem {
	ls := float64(10 * edakit.Scaled(ctx, 20, 120))
	switch kp(ctx) {
	case 1:
		lr := float64(10 * edakit.Int(20, 120))
		m := float64(10 * edakit.Int(5, int(math.Min(ls, lr)/10)-2))
		loop := ls + lr - 2*m
		return num(edakit.Sprintf("A signal wire has partial self inductance %s pH, its return wire has partial self inductance %s pH, and their partial mutual inductance is %s pH. Current flows out on the signal and back on the return. What is the loop inductance, in pH?", f(ls, 0), f(lr, 0), f(m, 0)),
			loop, 0, edakit.Sprintf("With opposite currents the mutual term subtracts: $L_{loop} = L_s + L_r - 2M$ = %s + %s - 2 x %s = %s pH.", f(ls, 0), f(lr, 0), f(m, 0), f(loop, 0)), nil)
	case 2:
		lr := float64(10 * edakit.Int(20, 100))
		m1 := float64(10 * edakit.Int(int(0.5*math.Min(ls, lr)/10), int(0.8*math.Min(ls, lr)/10)))
		m2 := float64(10 * edakit.Int(2, int(0.4*math.Min(ls, lr)/10)))
		l1 := ls + lr - 2*m1
		l2 := ls + lr - 2*m2
		return num(edakit.Sprintf("A signal wire (partial self inductance %s pH) has two candidate return paths of the same partial self inductance, %s pH: an adjacent ground wire with mutual inductance %s pH to the signal, and a distant power strap with mutual inductance %s pH. At high frequency the return current concentrates in the path that minimizes loop inductance. What is that loop inductance, in pH?", f(ls, 0), f(lr, 0), f(m1, 0), f(m2, 0)),
			math.Min(l1, l2), 0, edakit.Sprintf("Adjacent return: %s + %s - 2 x %s = %s pH. Distant strap: %s + %s - 2 x %s = %s pH. The closer return, with larger mutual inductance, gives the smaller loop, %s pH. This is why inductance depends on where the return current flows, which the signal net alone does not determine.", f(ls, 0), f(lr, 0), f(m1, 0), f(l1, 0), f(ls, 0), f(lr, 0), f(m2, 0), f(l2, 0), f(math.Min(l1, l2), 0)), nil)
	}
	la := float64(10 * edakit.Int(20, 80))
	lb := float64(10 * edakit.Int(20, 80))
	m := float64(10 * edakit.Int(2, int(math.Min(la, lb)/30)+2))
	tot := la + lb + 2*m
	return num(edakit.Sprintf("An extractor fractures a straight wire into two collinear segments with partial self inductances %s pH and %s pH, and a partial mutual inductance of %s pH between them. Current flows the same way in both. What is the partial self inductance of the whole wire, in pH?", f(la, 0), f(lb, 0), f(m, 0)),
		tot, 0, edakit.Sprintf("For segments in series with the same current direction, $L = L_a + L_b + 2M$ = %s + %s + 2 x %s = %s pH. Partial inductance is not additive over length, so an extractor must keep the mutual terms between segments, which makes the inductance matrix dense.", f(la, 0), f(lb, 0), f(m, 0), f(tot, 0)), nil)
}

// significance decides whether inductance matters for a wire: knee frequency, the ratio of
// inductive to resistive impedance, and time of flight.
func significance(ctx generator.GeneratorContext) generator.Problem {
	tr := float64(5 * edakit.Int(4, 40))
	switch kp(ctx) {
	case 1:
		fk := 0.35 / tr * 1000
		return num(edakit.Sprintf("A clock driver produces edges with a 10 to 90 percent rise time of %s ps. Using the knee-frequency rule $f \\approx 0.35/t_r$, what is the highest significant frequency in the signal, in GHz?", f(tr, 0)),
			fk, 2, edakit.Sprintf("0.35 / %s ps = %s GHz. Above this frequency the edge carries little energy, so this is the frequency at which to compare inductive and resistive impedance.", f(tr, 0), f(fk, 2)), nil)
	case 2:
		lp := edakit.Pick([]float64{0.4, 0.5, 0.6, 0.8, 1})
		rp := edakit.Pick([]float64{0.01, 0.02, 0.05, 0.1, 0.2})
		l := float64(100 * edakit.Scaled(ctx, 2, 30))
		fk := 0.35 / (tr * 1e-12)
		wl := 2 * math.Pi * fk * lp * 1e-12 * l
		r := rp * l
		ratio := wl / r
		return num(edakit.Sprintf("A wide global wire %s um long has resistance %s ohm/um and loop inductance %s pH/um. The signal on it has a rise time of %s ps, so take its significant frequency as $f = 0.35/t_r$. What is the ratio of inductive to resistive impedance, $\\omega L / R$, at that frequency?", f(l, 0), f(rp, 2), f(lp, 1), f(tr, 0)),
			ratio, 2, edakit.Sprintf("f = %s GHz, so $\\omega L$ = 2 pi x %s GHz x %s pH = %s ohm, against R = %s ohm. Ratio %s. A ratio well above 1 means an RC model misses ringing and overshoot; thin local wires, with high resistance, stay RC.", f(fk/1e9, 2), f(fk/1e9, 2), f(lp*l, 0), f(wl, 2), f(r, 1), f(ratio, 2)), nil)
	}
	lp := edakit.Pick([]float64{0.4, 0.5, 0.6, 0.8, 1})
	cp := edakit.Pick([]float64{0.15, 0.2, 0.25, 0.3})
	l := float64(500 * edakit.Scaled(ctx, 2, 10))
	if edakit.Int(0, 1) == 0 {
		tof := l * math.Sqrt(lp*cp) * math.Sqrt(1e-27) / 1e-12
		return num(edakit.Sprintf("A global wire %s um long has inductance %s pH/um and capacitance %s fF/um. What is its time of flight $l\\sqrt{L'C'}$, in ps? A rise time shorter than about twice this flags the wire for an RLC model.", f(l, 0), f(lp, 1), f(cp, 2)),
			tof, 2, edakit.Sprintf("$\\sqrt{L'C'}$ = sqrt(%s pH/um x %s fF/um) = %s ps/um. Times %s um gives %s ps.", f(lp, 1), f(cp, 2), f(math.Sqrt(lp*cp)*0.0316228, 5), f(l, 0), f(tof, 2)), nil)
	}
	z0 := math.Sqrt(lp*1e-12/(cp*1e-15)) * 1
	return num(edakit.Sprintf("A global wire has inductance %s pH/um and capacitance %s fF/um. What is its characteristic impedance $Z_0 = \\sqrt{L'/C'}$, in ohms? A wire whose total resistance exceeds about $2Z_0$ is damped enough that an RC model suffices.", f(lp, 1), f(cp, 2)),
		z0, 2, edakit.Sprintf("$\\sqrt{%s \\times 10^{-12} / (%s \\times 10^{-15})}$ = %s ohm.", f(lp, 1), f(cp, 2), f(z0, 2)), nil)
}
