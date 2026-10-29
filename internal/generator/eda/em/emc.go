package em

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func registerEMC(reg *generator.Registry) {
	reg.Register("em.emc.return_path", edakit.Func(returnPath))
	reg.Register("em.emc.ground_bounce", edakit.Func(groundBounce))
	reg.Register("em.emc.shielding_grounding", &edakit.Bank{Items: shieldingItems})
}

func returnPath(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		r1, r2 := float64(edakit.Int(2, 60)), float64(edakit.Int(2, 60))
		fr := 100 * r2 / (r1 + r2)
		return edakit.Numeric(
			sp("A trace's return current can come back through two paths: the plane strip directly beneath the trace ($R_1 = %s$ milliohm, $L_1$ small) or a wide detour through the board ($R_2 = %s$ milliohm, $L_2$ large). At DC, what percentage of the return current takes the path under the trace? Ignore mutual coupling. Round to 1 decimal place.", num(r1, 0), num(r2, 0)),
			num(fr, 1),
			sp("At DC only resistance matters, and parallel paths divide current inversely to resistance: $R_2/(R_1 + R_2) = %s/%s = %s$ percent.", num(r2, 0), num(r1+r2, 0), num(fr, 1)),
			facts("fraction_pct", num(fr, 1)))
	case 2:
		l1, l2 := pick(0.5, 1, 1.5, 2, 3), pick(4, 5, 8, 10, 20)
		if edakit.Int(0, 1) == 0 {
			l1, l2 = l2, l1
		}
		fr := 100 * l2 / (l1 + l2)
		return edakit.Numeric(
			sp("A trace's return current can come back through path 1 (loop inductance %s nH) or path 2 (loop inductance %s nH). Far above the frequency where reactance exceeds resistance in both, what percentage of the return current takes path 1? Ignore mutual coupling between the paths. Round to 1 decimal place.", num(l1, 1), num(l2, 1)),
			num(fr, 1),
			sp("At high frequency the paths' impedances are $j\\omega L$, so current divides inversely to inductance: $L_2/(L_1 + L_2) = %s/%s = %s$ percent. The low-inductance path is the one that keeps the loop area smallest.", num(l2, 1), num(l1+l2, 1), num(fr, 1)),
			facts("fraction_pct", num(fr, 1)))
	default:
		r := float64(edakit.Int(1, 50))
		l := pick(0.5, 1, 2, 5, 10, 20)
		f := r * 1e-3 / (2 * math.Pi * l * 1e-9) / 1e3
		return edakit.Numeric(
			sp("A return path has %s milliohm of resistance and %s nH of loop inductance. At what frequency does its reactance equal its resistance, in kHz? Round to 2 decimal places.", num(r, 0), num(l, 1)),
			num(f, 2),
			sp("$\\omega L = R$ at $f = R/(2\\pi L) = %s \\times 10^{-3}/(2\\pi \\cdot %s \\times 10^{-9})$ Hz, which is %s kHz. Above this the return current follows inductance rather than resistance.", num(r, 0), num(l, 1), num(f, 2)),
			facts("f_khz", num(f, 2)))
	}
}

func groundBounce(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		l := pick(0.5, 1, 2, 3, 5)
		di := float64(5 * edakit.Int(2, 20))
		dt := pick(0.1, 0.2, 0.25, 0.5, 1)
		v := l * di / dt
		return edakit.Numeric(
			sp("An output driver's ground current rises by %s mA in %s ns through a ground pin and bond wire with %s nH of inductance. How far does the on-die ground move, in mV? Round to 1 decimal place.", num(di, 0), num(dt, 2), num(l, 1)),
			num(v, 1),
			sp("$V = L\\,di/dt = %s\\text{ nH} \\cdot %s\\text{ mA}/%s\\text{ ns} = %s$ mV.", num(l, 1), num(di, 0), num(dt, 2), num(v, 1)),
			facts("v_mv", num(v, 1)))
	case 2:
		n := edakit.Int(4, 64)
		p := edakit.Int(1, 8)
		l := pick(1, 2, 3, 5)
		di := float64(5 * edakit.Int(1, 10))
		dt := pick(0.2, 0.25, 0.5, 1)
		v := float64(n) * di / dt * l / float64(p)
		return edakit.Numeric(
			sp("%d output drivers switch at once, each pulling %s mA more ground current within %s ns. They share %d ground pins in parallel, each with %s nH of inductance. Ignoring mutual inductance, what is the ground bounce, in mV? Round to 1 decimal place.", n, num(di, 0), num(dt, 2), p, num(l, 0)),
			num(v, 1),
			sp("The pins act as $L/P = %s$ nH, and the total $di/dt$ is $%d \\cdot %s/%s = %s$ mA/ns, so $V = %s$ mV.", num(l/float64(p), 3), n, num(di, 0), num(dt, 2), num(float64(n)*di/dt, 1), num(v, 1)),
			facts("leff_nh", num(l/float64(p), 3)))
	default:
		i := float64(10 * edakit.Int(1, 50))
		dt := pick(0.5, 1, 2, 5, 10)
		dv := float64(10 * edakit.Int(1, 10))
		c := i * dt / dv
		return edakit.Numeric(
			sp("A block draws an extra %s mA for %s ns before the package can respond. The local supply may droop by at most %s mV. How much on-die decoupling capacitance must supply that charge, in nF? Round to 3 decimal places.", num(i, 0), num(dt, 1), num(dv, 0)),
			num(c, 3),
			sp("The decap delivers $Q = I\\,\\Delta t = %s$ pC while drooping $\\Delta V = Q/C$, so $C = %s\\text{ pC}/%s\\text{ mV} = %s$ nF.", num(i*dt, 1), num(i*dt, 1), num(dv, 0), num(c, 3)),
			facts("q_pc", num(i*dt, 1)))
	}
}
