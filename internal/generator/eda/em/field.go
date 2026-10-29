package em

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func registerField(reg *generator.Registry) {
	reg.Register("em.field.maxwell_equations", &edakit.Bank{Items: maxwellItems})
	reg.Register("em.field.quasi_static", edakit.Func(quasiStatic))
	reg.Register("em.field.electrostatics_capacitance", edakit.Func(electrostaticsCapacitance))
	reg.Register("em.field.fringe_capacitance", edakit.Func(fringeCapacitance))
	reg.Register("em.field.coupling_capacitance", edakit.Func(couplingCapacitance))
	reg.Register("em.field.wire_resistance", edakit.Func(wireResistance))
	reg.Register("em.field.inductance", edakit.Func(inductance))
	reg.Register("em.field.skin_effect", edakit.Func(skinEffect))
	reg.Register("em.field.dielectric_loss", edakit.Func(dielectricLoss))
}

var dielectrics = []struct {
	name string
	er   float64
}{
	{"silicon dioxide ($\\varepsilon_r = 3.9$)", 3.9},
	{"a low-k dielectric ($\\varepsilon_r = 2.7$)", 2.7},
	{"FR-4 ($\\varepsilon_r = 4.3$)", 4.3},
	{"a dielectric with $\\varepsilon_r = 3.0$", 3.0},
	{"polyimide ($\\varepsilon_r = 3.4$)", 3.4},
}

func quasiStatic(ctx generator.GeneratorContext) generator.Problem {
	d := edakit.Pick(dielectrics)
	switch kp(ctx) {
	case 1:
		f := pick(1, 2, 5, 10, 20, 28, 40, 60, 77, 100)
		lam := cLight / (f * 1e9 * math.Sqrt(d.er)) * 1e3
		return edakit.Numeric(
			sp("What is the wavelength of a %s GHz wave travelling in %s, in mm? Round to 2 decimal places.", num(f, 0), d.name),
			num(lam, 2),
			sp("$\\lambda = c/(f\\sqrt{\\varepsilon_r}) = 3.00 \\times 10^8/(%s \\times 10^9 \\cdot %s)$ m, which is %s mm.", num(f, 0), num(math.Sqrt(d.er), 4), num(lam, 2)),
			facts("lambda_mm", num(lam, 2)))
	case 2:
		if edakit.Int(0, 1) == 0 {
			f := pick(1, 2, 5, 10, 20, 40, 60, 100)
			lam := cLight / (f * 1e9 * math.Sqrt(d.er)) * 1e3
			return edakit.Numeric(
				sp("A model treats a structure as quasi-static (lumped) when its largest dimension is below one tenth of the wavelength. In %s at %s GHz, what is the largest dimension that passes, in mm? Round to 2 decimal places.", d.name, num(f, 0)),
				num(lam/10, 2),
				sp("$\\lambda = %s$ mm, so the limit is $\\lambda/10 = %s$ mm.", num(lam, 2), num(lam/10, 2)),
				facts("lambda_mm", num(lam, 2)))
		}
		l := pick(0.5, 1, 2, 5, 10, 20, 50)
		f := cLight / (10 * l * 1e-3 * math.Sqrt(d.er)) / 1e9
		return edakit.Numeric(
			sp("A model treats a structure as quasi-static (lumped) when its largest dimension is below one tenth of the wavelength. A %s mm structure is embedded in %s. Up to what frequency does it pass, in GHz? Round to 2 decimal places.", num(l, 1), d.name),
			num(f, 2),
			sp("Require $\\lambda = c/(f\\sqrt{\\varepsilon_r}) \\ge 10\\ell$, so $f \\le c/(10\\ell\\sqrt{\\varepsilon_r}) = %s$ GHz.", num(f, 2)),
			facts("f_ghz", num(f, 2)))
	default:
		rho := pick(0.01, 0.1, 1, 5, 10, 20, 50, 100)
		tau := eps0 * 11.7 * rho / 100 * 1e12
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(
				sp("A silicon substrate ($\\varepsilon_r = 11.7$) has resistivity %s ohm cm. What is its charge relaxation time $\\tau = \\varepsilon/\\sigma$, in picoseconds? Round to 3 decimal places.", num(rho, 2)),
				num(tau, 3),
				sp("$\\tau = \\varepsilon_0\\varepsilon_r\\rho = 8.854 \\times 10^{-12} \\cdot 11.7 \\cdot %s$ s (ohm cm divided by 100 is ohm m), which is %s ps.", num(rho/100, 4), num(tau, 3)),
				facts("tau_ps", num(tau, 3)))
		}
		f := 1 / (2 * math.Pi * tau * 1e-12) / 1e9
		return edakit.Numeric(
			sp("A silicon substrate ($\\varepsilon_r = 11.7$) has resistivity %s ohm cm. Above the frequency $1/(2\\pi\\tau)$, with $\\tau = \\varepsilon/\\sigma$, it behaves more like a dielectric than a conductor. What is that frequency, in GHz? Round to 2 decimal places.", num(rho, 2)),
			num(f, 2),
			sp("$\\tau = \\varepsilon_0 \\varepsilon_r \\rho = %s$ ps, and $1/(2\\pi\\tau) = %s$ GHz. Below it conduction current dominates displacement current.", num(tau, 3), num(f, 2)),
			facts("tau_ps", num(tau, 3)))
	}
}

func electrostaticsCapacitance(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		v := edakit.Pick(supplies)
		dnm := float64(edakit.Int(1, 50))
		e := v / dnm * 1000
		return edakit.Numeric(
			sp("Two parallel plates %s nm apart are held %s V apart. What is the electric field between them, ignoring fringing, in MV/m? Round to 1 decimal place.", num(dnm, 0), num(v, 1)),
			num(e, 1),
			sp("The field between ideal plates is uniform, $E = V/d = %s/(%s \\times 10^{-9})$ V/m, which is %s MV/m.", num(v, 1), num(dnm, 0), num(e, 1)),
			facts("e_mv_per_m", num(e, 1)))
	case 2:
		d := edakit.Pick(dielectrics)
		w, l := float64(edakit.Int(1, 50)), float64(edakit.Int(1, 50))
		t := pick(0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1)
		c := 8.854e-3 * d.er * w * l / t
		return edakit.Numeric(
			sp("A %s um by %s um plate sits %s um above a ground plane, separated by %s. Ignoring fringing, what is its capacitance, in fF? Round to 3 decimal places.", num(w, 0), num(l, 0), num(t, 2), d.name),
			num(c, 3),
			sp("$C = \\varepsilon_0\\varepsilon_r A/d$ with $\\varepsilon_0 = 8.854 \\times 10^{-3}$ fF/um: $8.854 \\times 10^{-3} \\cdot %s \\cdot %s/%s = %s$ fF.", num(d.er, 1), num(w*l, 0), num(t, 2), num(c, 3)),
			facts("c_ff", num(c, 3)))
	default:
		w, l := float64(edakit.Int(1, 50)), float64(edakit.Int(1, 50))
		d1, d2 := pick(0.05, 0.1, 0.2, 0.3), pick(0.02, 0.05, 0.1, 0.2)
		e1, e2 := pick(2.7, 3.0, 3.9), pick(4.2, 5.5, 7.0)
		c := 8.854e-3 * w * l / (d1/e1 + d2/e2)
		return edakit.Numeric(
			sp("A %s um by %s um plate sits over ground through two stacked layers: %s um of $\\varepsilon_r = %s$ and %s um of $\\varepsilon_r = %s$. Ignoring fringing, what is its capacitance, in fF? Round to 3 decimal places.", num(w, 0), num(l, 0), num(d1, 2), num(e1, 1), num(d2, 2), num(e2, 1)),
			num(c, 3),
			sp("The layers act as capacitors in series: $C = \\varepsilon_0 A/(d_1/\\varepsilon_1 + d_2/\\varepsilon_2) = 8.854 \\times 10^{-3} \\cdot %s/(%s) = %s$ fF.", num(w*l, 0), num(d1/e1+d2/e2, 5), num(c, 3)),
			facts("c_ff", num(c, 3)))
	}
}

func fringeCapacitance(ctx generator.GeneratorContext) generator.Problem {
	er := pick(2.5, 2.7, 3.0, 3.9)
	h := pick(0.1, 0.2, 0.3, 0.5)
	w := h * pick(0.5, 1, 1.5, 2, 3, 5, 10)
	t := h * pick(0.5, 1, 1.5, 2, 3)
	e := 8.854 * er // aF/um
	area := e * w / h
	total := e * (1.15*(w/h) + 2.80*math.Pow(t/h, 0.222))
	stem := sp("A wire of width %s um and thickness %s um runs %s um above a ground plane in a dielectric with $\\varepsilon_r = %s$.", num(w, 2), num(t, 2), num(h, 2), num(er, 1))
	f := facts("area_af_um", num(area, 2), "total_af_um", num(total, 2))
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(stem+" Counting only the parallel-plate (area) term, what is its capacitance per unit length, in aF/um? Round to 2 decimal places.",
			num(area, 2),
			sp("$C/\\ell = \\varepsilon W/H = 8.854 \\cdot %s \\cdot %s/%s = %s$ aF/um ($\\varepsilon_0 = 8.854$ aF/um).", num(er, 1), num(w, 2), num(h, 2), num(area, 2)), f)
	case 2:
		return edakit.Numeric(stem+" Using the Sakurai-Tamaru fit $C/(\\varepsilon\\ell) = 1.15(W/H) + 2.80(T/H)^{0.222}$, what is its total capacitance per unit length, in aF/um? Round to 2 decimal places.",
			num(total, 2),
			sp("$W/H = %s$ and $T/H = %s$, so $C/(\\varepsilon\\ell) = 1.15 \\cdot %s + 2.80 \\cdot %s^{0.222} = %s$. Times $\\varepsilon = %s$ aF/um gives %s aF/um.", num(w/h, 2), num(t/h, 2), num(w/h, 2), num(t/h, 2), num(total/e, 4), num(e, 3), num(total, 2)), f)
	default:
		share := 100 * (total - area) / total
		return edakit.Numeric(stem+" Using the Sakurai-Tamaru fit $C/(\\varepsilon\\ell) = 1.15(W/H) + 2.80(T/H)^{0.222}$ for the total, what percentage of the total is not accounted for by the parallel-plate term $\\varepsilon W/H$? Round to 1 decimal place.",
			num(share, 1),
			sp("Parallel plate: %s aF/um. Total: %s aF/um. The fringe share is $(%s - %s)/%s = %s$ percent; it grows as $W/H$ shrinks.", num(area, 2), num(total, 2), num(total, 2), num(area, 2), num(total, 2), num(share, 1)), f)
	}
}

func couplingCapacitance(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		er := pick(2.5, 2.7, 3.0, 3.9)
		t := pick(0.1, 0.2, 0.3, 0.4, 0.5)
		s := pick(0.05, 0.1, 0.14, 0.2, 0.3)
		l := float64(50 * edakit.Int(1, 40))
		cc := 8.854e-3 * er * t * l / s
		return edakit.Numeric(
			sp("Two parallel wires of thickness %s um run side by side for %s um with %s um spacing, in a dielectric with $\\varepsilon_r = %s$. Treating the facing sidewalls as a parallel-plate capacitor, what is the coupling capacitance, in fF? Round to 3 decimal places.", num(t, 2), num(l, 0), num(s, 2), num(er, 1)),
			num(cc, 3),
			sp("$C_c = \\varepsilon_0\\varepsilon_r T\\ell/S = 8.854 \\times 10^{-3} \\cdot %s \\cdot %s \\cdot %s/%s = %s$ fF.", num(er, 1), num(t, 2), num(l, 0), num(s, 2), num(cc, 3)),
			facts("cc_ff", num(cc, 3)))
	case 2:
		cg, cc := float64(edakit.Int(5, 60)), float64(edakit.Int(2, 40))
		cases := []struct {
			txt string
			k   float64
		}{
			{"switches in the same direction at the same time as the victim", 0},
			{"stays quiet while the victim switches", 1},
			{"switches in the opposite direction at the same time as the victim", 2},
		}
		x := edakit.Pick(cases)
		eff := cg + x.k*cc
		return edakit.Numeric(
			sp("A victim wire has %s fF to ground and %s fF of coupling to one neighbour. The neighbour %s. Using the Miller-factor model, what effective capacitance does the victim's driver charge, in fF?", num(cg, 0), num(cc, 0), x.txt),
			num(eff, 0),
			sp("The coupling capacitor is charged by the voltage difference across it: no change when the neighbour moves with the victim (factor 0), one swing when it is quiet (1), two swings when it moves the other way (2). Here the factor is %s: $%s + %s \\cdot %s = %s$ fF.", num(x.k, 0), num(cg, 0), num(x.k, 0), num(cc, 0), num(eff, 0)),
			facts("miller", num(x.k, 0)))
	default:
		cg, cc := float64(edakit.Int(5, 60)), float64(edakit.Int(2, 40))
		v := edakit.Pick(supplies)
		dv := 1000 * v * cc / (cc + cg)
		return edakit.Numeric(
			sp("A precharged dynamic node floats (nothing drives it) with %s fF to ground. A neighbouring wire coupled to it by %s fF swings by %s V. How large is the glitch on the floating node, in mV? Round to 1 decimal place.", num(cg, 0), num(cc, 0), num(v, 1)),
			num(dv, 1),
			sp("Charge is conserved on the floating node, so it follows the capacitive divider: $\\Delta V = V \\cdot C_c/(C_c + C_g) = %s \\cdot %s/%s = %s$ mV.", num(1000*v, 0), num(cc, 0), num(cc+cg, 0), num(dv, 1)),
			facts("dv_mv", num(dv, 1)))
	}
}

var metals = []struct {
	name string
	rho  float64 // microohm cm
}{
	{"copper (effective resistivity 2.2 microohm cm)", 2.2},
	{"copper (resistivity 1.72 microohm cm)", 1.72},
	{"aluminium (resistivity 2.7 microohm cm)", 2.7},
	{"tungsten (resistivity 5.6 microohm cm)", 5.6},
	{"a thin copper line (effective resistivity 3.5 microohm cm)", 3.5},
}

func wireResistance(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		m := edakit.Pick(metals)
		t := float64(10 * edakit.Int(3, 100))
		rs := 10 * m.rho / t
		return edakit.Numeric(
			sp("A metal layer of %s is %s nm thick. What is its sheet resistance, in ohm per square? Round to 4 decimal places.", m.name, num(t, 0)),
			num(rs, 4),
			sp("$R_s = \\rho/t = %s \\times 10^{-8}\\text{ ohm m}/(%s \\times 10^{-9}\\text{ m}) = %s$ ohm per square.", num(m.rho, 2), num(t, 0), num(rs, 4)),
			facts("rs", num(rs, 4)))
	case 2:
		rs := pick(0.02, 0.04, 0.05, 0.08, 0.1, 0.2)
		w := pick(0.05, 0.1, 0.14, 0.2, 0.4, 1)
		l := float64(10 * edakit.Int(1, 200))
		r := rs * l / w
		return edakit.Numeric(
			sp("A wire on a layer with sheet resistance %s ohm per square is %s um wide and %s um long. What is its resistance, in ohms? Round to 2 decimal places.", num(rs, 2), num(w, 2), num(l, 0)),
			num(r, 2),
			sp("The wire is $L/W = %s/%s = %s$ squares long, so $R = %s \\cdot %s = %s$ ohm.", num(l, 0), num(w, 2), num(l/w, 1), num(rs, 2), num(l/w, 1), num(r, 2)),
			facts("squares", num(l/w, 1)))
	default:
		rs1, rs2 := pick(0.08, 0.1, 0.15), pick(0.02, 0.03, 0.04)
		w1, w2 := pick(0.05, 0.07, 0.1), pick(0.2, 0.4)
		l1, l2 := float64(10*edakit.Int(1, 50)), float64(10*edakit.Int(1, 100))
		rv := pick(1, 2, 5, 10)
		n := edakit.Int(1, 4)
		r := rs1*l1/w1 + rv/float64(n) + rs2*l2/w2
		return edakit.Numeric(
			sp("A route runs %s um on a lower layer ($R_s = %s$ ohm per square, width %s um), climbs through a via array of %d parallel cuts of %s ohm each, and runs %s um on an upper layer ($R_s = %s$ ohm per square, width %s um). What is the total resistance, in ohms? Round to 2 decimal places.", num(l1, 0), num(rs1, 2), num(w1, 2), n, num(rv, 0), num(l2, 0), num(rs2, 2), num(w2, 2)),
			num(r, 2),
			sp("Lower segment: $%s \\cdot %s/%s = %s$ ohm. Via array: $%s/%d = %s$ ohm. Upper segment: %s ohm. Series total: %s ohm.", num(rs1, 2), num(l1, 0), num(w1, 2), num(rs1*l1/w1, 2), num(rv, 0), n, num(rv/float64(n), 2), num(rs2*l2/w2, 2), num(r, 2)),
			facts("via", num(rv/float64(n), 2)))
	}
}

func inductance(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		a := pick(0.1, 0.2, 0.25, 0.5, 1)
		b := a * pick(2, 2.3, 3, 3.5, 4, 5, 10)
		l := 0.2 * math.Log(b/a)
		return edakit.Numeric(
			sp("A coaxial line has an inner conductor of radius %s mm and an outer conductor of inner radius %s mm. Assuming the current flows on the facing surfaces, what is its inductance per unit length, in nH/mm? Round to 4 decimal places.", num(a, 2), num(b, 3)),
			num(l, 4),
			sp("Ampere's law gives $B = \\mu_0 I/(2\\pi r)$ between the conductors; the flux per unit length is $\\frac{\\mu_0 I}{2\\pi}\\ln(b/a)$, so $L' = 0.2\\ln(%s) = %s$ nH/mm ($\\mu_0/2\\pi = 0.2$ nH/mm).", num(b/a, 3), num(l, 4)),
			facts("ratio", num(b/a, 3)))
	case 2:
		a := pick(0.05, 0.1, 0.2, 0.25)
		if edakit.Int(0, 1) == 0 {
			d := a * pick(3, 4, 5, 8, 10, 20)
			l := 0.4 * math.Acosh(d/(2*a))
			return edakit.Numeric(
				sp("Two parallel round wires of radius %s mm, with centres %s mm apart, carry a current out and back. What is the loop inductance per unit length, in nH/mm? Round to 4 decimal places.", num(a, 2), num(d, 2)),
				num(l, 4),
				sp("$L' = (\\mu_0/\\pi)\\,\\mathrm{acosh}(D/2a) = 0.4\\,\\mathrm{acosh}(%s) = %s$ nH/mm, close to $0.4\\ln(D/a) = %s$ when $D \\gg a$.", num(d/(2*a), 3), num(l, 4), num(0.4*math.Log(d/a), 4)),
				facts("arg", num(d/(2*a), 3)))
		}
		h := a * pick(2, 3, 4, 5, 10, 20)
		l := 0.2 * math.Acosh(h/a)
		return edakit.Numeric(
			sp("A round wire of radius %s mm runs with its centre %s mm above a perfectly conducting return plane. What is the loop inductance per unit length, in nH/mm? Round to 4 decimal places.", num(a, 2), num(h, 2)),
			num(l, 4),
			sp("By image theory the plane acts like a return wire at depth $h$ below it, and the field above the plane is half of the two-wire case: $L' = (\\mu_0/2\\pi)\\,\\mathrm{acosh}(h/a) = 0.2\\,\\mathrm{acosh}(%s) = %s$ nH/mm.", num(h/a, 2), num(l, 4)),
			facts("arg", num(h/a, 2)))
	default:
		d := edakit.Pick(dielectrics)
		cp := float64(10 * edakit.Int(5, 20))
		l := d.er / (cLight * cLight * cp * 1e-12) * 1e9
		return edakit.Numeric(
			sp("A uniform two-conductor line is embedded entirely in %s and has a capacitance of %s pF/m. What is its inductance per unit length, in nH/m? Round to 1 decimal place.", d.name, num(cp, 0)),
			num(l, 1),
			sp("For a TEM line in a homogeneous dielectric $L'C' = \\mu_0\\varepsilon_0\\varepsilon_r = \\varepsilon_r/c^2$, so $L' = %s/((3.00 \\times 10^8)^2 \\cdot %s \\times 10^{-12})$ H/m, which is %s nH/m.", num(d.er, 1), num(cp, 0), num(l, 1)),
			facts("l_nh_m", num(l, 1)))
	}
}

var conductors = []struct {
	name string
	rho  float64 // ohm m
}{
	{"copper ($\\rho = 1.72 \\times 10^{-8}$ ohm m)", 1.72e-8},
	{"aluminium ($\\rho = 2.65 \\times 10^{-8}$ ohm m)", 2.65e-8},
	{"gold ($\\rho = 2.44 \\times 10^{-8}$ ohm m)", 2.44e-8},
}

func skinDepth(rho, f float64) float64 { return math.Sqrt(rho / (math.Pi * f * mu0)) }

func skinEffect(ctx generator.GeneratorContext) generator.Problem {
	m := edakit.Pick(conductors)
	switch kp(ctx) {
	case 1:
		f := pick(0.01, 0.1, 0.5, 1, 2, 5, 10, 20, 50)
		d := skinDepth(m.rho, f*1e9) * 1e6
		return edakit.Numeric(
			sp("What is the skin depth in %s at %s GHz, in um? Use $\\mu = \\mu_0$. Round to 3 decimal places.", m.name, num(f, 2)),
			num(d, 3),
			sp("$\\delta = \\sqrt{\\rho/(\\pi f \\mu_0)} = \\sqrt{%s/(\\pi \\cdot %s \\times 10^9 \\cdot 4\\pi \\times 10^{-7})}$ m, which is %s um.", num(m.rho*1e8, 2)+" \\times 10^{-8}", num(f, 2), num(d, 3)),
			facts("delta_um", num(d, 3)))
	case 2:
		f := pick(1, 2, 5, 10, 20)
		d := skinDepth(m.rho, f*1e9) * 1e6
		a := float64(edakit.Int(10, 200))
		ratio := a / (2 * d)
		return edakit.Numeric(
			sp("A round wire of %s has radius %s um. At %s GHz the skin depth is much smaller than the radius. What is the ratio of its AC resistance to its DC resistance? Round to 2 decimal places.", m.name, num(a, 0), num(f, 0)),
			num(ratio, 2),
			sp("$\\delta = %s$ um. The current crowds into a shell of thickness $\\delta$, so $R_{ac} \\approx \\rho/(2\\pi a\\delta)$ per length against $R_{dc} = \\rho/(\\pi a^2)$: the ratio is $a/(2\\delta) = %s$.", num(d, 3), num(ratio, 2)),
			facts("delta_um", num(d, 3)))
	default:
		r1 := pick(0.5, 1, 2, 3, 5)
		f1 := pick(1, 2, 5)
		f2 := f1 * pick(4, 9, 16, 2, 10, 25)
		r2 := r1 * math.Sqrt(f2/f1)
		return edakit.Numeric(
			sp("A package trace deep in the skin-effect regime has an AC resistance of %s ohm at %s GHz. What is its AC resistance at %s GHz, in ohms? Round to 3 decimal places.", num(r1, 1), num(f1, 0), num(f2, 0)),
			num(r2, 3),
			sp("Skin depth scales as $1/\\sqrt f$, so resistance scales as $\\sqrt f$: $%s\\sqrt{%s/%s} = %s$ ohm.", num(r1, 1), num(f2, 0), num(f1, 0), num(r2, 3)),
			facts("ratio", num(math.Sqrt(f2/f1), 4)))
	}
}

func dielectricLoss(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		c := float64(edakit.Int(1, 50))
		f := pick(0.1, 0.5, 1, 2, 5, 10)
		tl := pick(0.001, 0.002, 0.005, 0.01, 0.02)
		g := 2 * math.Pi * f * 1e9 * c * 1e-12 * tl * 1e6
		return edakit.Numeric(
			sp("A %s pF capacitor uses a dielectric with loss tangent %s. At %s GHz, what parallel conductance models its loss, in microsiemens? Round to 2 decimal places.", num(c, 0), num(tl, 3), num(f, 1)),
			num(g, 2),
			sp("With $\\varepsilon = \\varepsilon'(1 - j\\tan\\delta)$, the admittance $j\\omega C$ gains a real part $G = \\omega C\\tan\\delta = 2\\pi \\cdot %s \\times 10^9 \\cdot %s \\times 10^{-12} \\cdot %s$ S, which is %s uS.", num(f, 1), num(c, 0), num(tl, 3), num(g, 2)),
			facts("g_us", num(g, 2)))
	case 2:
		d := edakit.Pick(dielectrics)
		f := pick(1, 2, 5, 10, 16, 28)
		tl := pick(0.002, 0.005, 0.01, 0.015, 0.02)
		l := float64(edakit.Int(1, 50))
		alpha := math.Pi * f * 1e9 * math.Sqrt(d.er) * tl / cLight // Np/m
		db := 8.686 * alpha * l / 100
		return edakit.Numeric(
			sp("A %s cm stripline is fully embedded in %s with loss tangent %s. Considering dielectric loss only, how much does a %s GHz signal attenuate along it, in dB? Round to 3 decimal places.", num(l, 0), d.name, num(tl, 3), num(f, 0)),
			num(db, 3),
			sp("$\\alpha_d = \\pi f\\sqrt{\\varepsilon_r}\\tan\\delta/c = %s$ Np/m. Times 8.686 dB per neper and %s m gives %s dB.", num(alpha, 4), num(l/100, 2), num(db, 3)),
			facts("alpha_np_m", num(alpha, 4)))
	default:
		f0 := pick(1, 2, 5)
		ac := pick(0.5, 1, 1.5, 2, 3)
		ad := ac * pick(0.25, 0.4, 0.5, 0.8, 1, 1.25)
		fx := f0 * (ac / ad) * (ac / ad)
		return edakit.Numeric(
			sp("At %s GHz a channel loses %s dB to conductor (skin-effect) loss and %s dB to dielectric loss. Conductor loss scales as $\\sqrt f$ and dielectric loss as $f$. At what frequency are the two equal, in GHz? Round to 2 decimal places.", num(f0, 0), num(ac, 1), num(ad, 2)),
			num(fx, 2),
			sp("Set $%s\\sqrt{f/%s} = %s\\,(f/%s)$: then $\\sqrt{f/%s} = %s/%s$, so $f = %s \\cdot (%s)^2 = %s$ GHz. Above it dielectric loss dominates.", num(ac, 1), num(f0, 0), num(ad, 2), num(f0, 0), num(f0, 0), num(ac, 1), num(ad, 2), num(f0, 0), num(ac/ad, 3), num(fx, 2)),
			facts("f_ghz", num(fx, 2)))
	}
}
