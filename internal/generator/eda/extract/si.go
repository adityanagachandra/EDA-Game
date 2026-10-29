package extract

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type aggState struct {
	name string
	mf   float64
}

var aggStates = []aggState{
	{"switches in the same direction at the same time", 0},
	{"stays quiet", 1},
	{"switches in the opposite direction at the same time", 2},
}

// millerFactor computes effective coupling and the charge a victim driver must supply.
func millerFactor(ctx generator.GeneratorContext) generator.Problem {
	cc := float64(edakit.Int(5, 60)) / 10
	st := edakit.Pick(aggStates)
	switch kp(ctx) {
	case 1:
		return num(edakit.Sprintf("A victim net rises while its neighbour, coupled through %s fF, %s. Modelling the coupling capacitor as a grounded capacitor scaled by a Miller factor, what effective capacitance does it add to the victim's load, in fF?", f(cc, 1), st.name),
			cc*st.mf, 2, edakit.Sprintf("The Miller factor is %s: the voltage change across the capacitor is %s times the victim's own swing. Effective coupling %s x %s = %s fF.", f(st.mf, 0), f(st.mf, 0), f(cc, 1), f(st.mf, 0), f(cc*st.mf, 2)), nil)
	case 2:
		vdd := edakit.Pick([]float64{0.7, 0.75, 0.8, 0.9, 1.0})
		q := cc * vdd * st.mf
		return num(edakit.Sprintf("A victim rises from 0 to %s V while an aggressor coupled through %s fF %s, swinging the full %s V. How much charge must the victim's driver push through the coupling capacitor during the transition, in fC?", f(vdd, 2), f(cc, 1), st.name, f(vdd, 2)),
			q, 3, edakit.Sprintf("Charge is $C_c$ times the change in voltage across it. The victim end moves +%s V and the aggressor end moves %s, so the change across the capacitor is %s V, giving %s fC.", f(vdd, 2), map[float64]string{0: "+" + f(vdd, 2) + " V", 1: "0 V", 2: "-" + f(vdd, 2) + " V"}[st.mf], f(vdd*st.mf, 2), f(q, 3)), nil)
	}
	n := edakit.Scaled(ctx, 2, 4)
	var cs []string
	sum := 0.0
	for i := 0; i < n; i++ {
		c := float64(edakit.Int(5, 40)) / 10
		sum += c
		cs = append(cs, f(c, 1))
	}
	return num(edakit.Sprintf("A victim has coupling capacitances of %s fF to %d aggressors whose switching times are unknown. Max-delay analysis assumes each aggressor switches against the victim; min-delay analysis assumes each switches with it. By how much does the effective coupling capacitance differ between the two analyses, in fF?", strings.Join(cs, ", "), n),
		2*sum, 2, edakit.Sprintf("Max delay uses Miller factor 2 on each coupling, min delay uses 0. The difference is $2\\sum C_c$ = 2 x %s = %s fF. This spread is the price of not knowing the aggressors' timing.", f(sum, 1), f(2*sum, 2)), nil)
}

// effectiveCap sums ground, pin and Miller-scaled coupling capacitance.
func effectiveCap(ctx generator.GeneratorContext) generator.Problem {
	cg := float64(edakit.Int(20, 120)) / 10
	cp := float64(edakit.Int(5, 40)) / 10
	n := edakit.Scaled(ctx, 2, 4)
	ccs := make([]float64, n)
	var sts []aggState
	var desc []string
	sum := 0.0
	for i := range ccs {
		ccs[i] = float64(edakit.Int(5, 40)) / 10
		sum += ccs[i]
		s := edakit.Pick(aggStates)
		sts = append(sts, s)
		desc = append(desc, edakit.Sprintf("aggressor %d (%s fF) %s", i+1, f(ccs[i], 1), s.name))
	}
	stem := edakit.Sprintf("A victim net has %s fF of wire capacitance to ground and %s fF of sink pin capacitance. ", f(cg, 1), f(cp, 1))
	switch kp(ctx) {
	case 1:
		eff := cg + cp
		for i := range ccs {
			eff += sts[i].mf * ccs[i]
		}
		return num(stem+"During its transition, "+strings.Join(desc, ", ")+". Using Miller factors 0, 1 and 2, what effective load capacitance does the victim's driver see, in fF?", eff, 2,
			edakit.Sprintf("Add ground and pin capacitance, then each coupling times its Miller factor: %s fF.", f(eff, 2)), nil)
	case 2:
		var cs []string
		for _, c := range ccs {
			cs = append(cs, f(c, 1))
		}
		worst := cg + cp + 2*sum
		return num(stem+edakit.Sprintf("It couples to %d aggressors through %s fF. For worst-case (setup) delay, with every aggressor assumed to switch against the victim, what is the effective load, in fF?", n, strings.Join(cs, ", ")), worst, 2,
			edakit.Sprintf("%s + %s + 2 x %s = %s fF.", f(cg, 1), f(cp, 1), f(sum, 1), f(worst, 2)), nil)
	}
	var cs []string
	for _, c := range ccs {
		cs = append(cs, f(c, 1))
	}
	worst := cg + cp + 2*sum
	best := cg + cp
	return num(stem+edakit.Sprintf("It couples to %d aggressors through %s fF. What is the ratio of the worst-case (all aggressors opposing) to the best-case (all switching with the victim) effective load?", n, strings.Join(cs, ", ")), worst/best, 3,
		edakit.Sprintf("Worst %s fF, best %s fF (coupling contributes nothing when both ends move together). Ratio %s. When coupling dominates a net's capacitance, this ratio, and so the delay spread, becomes large.", f(worst, 2), f(best, 2), f(worst/best, 3)), nil)
}

// crosstalkDelay applies a lumped driver model to effective-capacitance changes.
func crosstalkDelay(ctx generator.GeneratorContext) generator.Problem {
	rd := float64(edakit.Int(5, 40)) / 10 // kohm
	cg := float64(edakit.Int(20, 100)) / 10
	cc := float64(edakit.Int(5, 50)) / 10
	model := edakit.Sprintf("Model the victim's driver as %s kohm driving a lumped load, with 50%% delay $t = 0.69 R C$ (kohm times fF is ps). The victim has %s fF to ground (including pins) and %s fF of coupling to one aggressor.", f(rd, 1), f(cg, 1), f(cc, 1))
	switch kp(ctx) {
	case 1:
		d := 0.69 * rd * (cg + cc)
		return num(model+" What is the victim's delay with the aggressor quiet, in ps?", d, 2,
			edakit.Sprintf("A quiet aggressor acts as ground, Miller factor 1: 0.69 x %s x (%s + %s) = %s ps.", f(rd, 1), f(cg, 1), f(cc, 1), f(d, 2)), nil)
	case 2:
		dd := 0.69 * rd * cc
		up := edakit.Int(0, 1) == 0
		if up {
			return num(model+" By how many ps does the delay increase when the aggressor switches in the opposite direction at the same time, relative to the quiet case?", dd, 2,
				edakit.Sprintf("The Miller factor goes from 1 to 2, adding $C_c$ to the effective load: 0.69 x %s x %s = %s ps.", f(rd, 1), f(cc, 1), f(dd, 2)), nil)
		}
		return num(model+" By how many ps does the delay decrease when the aggressor switches in the same direction at the same time, relative to the quiet case?", dd, 2,
			edakit.Sprintf("The Miller factor goes from 1 to 0, removing $C_c$ from the effective load: 0.69 x %s x %s = %s ps faster. The speed-up matters for hold checks.", f(rd, 1), f(cc, 1), f(dd, 2)), nil)
	}
	slack := float64(edakit.Int(-5, 40))
	hold := edakit.Int(0, 1) == 0
	dd := 0.69 * rd * cc
	if hold {
		return num(model+edakit.Sprintf(" A hold check through this victim, timed with the aggressor quiet, has %s ps of hold slack. If the aggressor can switch in the same direction at the same time, what is the hold slack, in ps?", f(slack, 0)),
			slack-dd, 2, edakit.Sprintf("Same-direction switching speeds the victim by %s ps, so data arrives earlier and hold slack falls to %s ps. Hold analysis uses the crosstalk speed-up; setup uses the slow-down.", f(dd, 2), f(slack-dd, 2)), nil)
	}
	return num(model+edakit.Sprintf(" A setup path through this victim, timed with the aggressor quiet, has %s ps of setup slack. If the aggressor can switch in the opposite direction at the same time, what is the setup slack, in ps?", f(slack, 0)),
		slack-dd, 2, edakit.Sprintf("Opposite switching slows the victim by %s ps, so setup slack falls to %s ps.", f(dd, 2), f(slack-dd, 2)), nil)
}

// glitch computes the charge-sharing glitch peak on a floating victim.
func glitch(ctx generator.GeneratorContext) generator.Problem {
	vdd := edakit.Pick([]float64{0.7, 0.75, 0.8, 0.9, 1.0})
	cg := float64(edakit.Int(10, 80)) / 10
	switch kp(ctx) {
	case 1:
		cc := float64(edakit.Int(5, 40)) / 10
		v := 1000 * vdd * cc / (cc + cg)
		return num(edakit.Sprintf("A quiet victim with %s fF to ground is coupled through %s fF to an aggressor that switches from 0 to %s V. Treating the victim as floating (its driver too weak to respond during the fast edge), what glitch peak appears on the victim, in mV?", f(cg, 1), f(cc, 1), f(vdd, 2)),
			v, 1, edakit.Sprintf("The capacitive divider gives $V_p = V_{DD} C_c/(C_c + C_g)$ = %s x %s / %s = %s mV. This is the upper bound on the glitch; a holding driver reduces it.", f(vdd, 2), f(cc, 1), f(cc+cg, 1), f(v, 1)), nil)
	case 2:
		c1 := float64(edakit.Int(5, 30)) / 10
		c2 := float64(edakit.Int(5, 30)) / 10
		both := edakit.Int(0, 1) == 0
		var v float64
		var q string
		if both {
			v = 1000 * vdd * (c1 + c2) / (c1 + c2 + cg)
			q = edakit.Sprintf("A quiet, floating victim with %s fF to ground couples through %s fF to aggressor A and %s fF to aggressor B. Both aggressors rise from 0 to %s V at the same instant. What is the glitch peak on the victim, in mV?", f(cg, 1), f(c1, 1), f(c2, 1), f(vdd, 2))
		} else {
			v = 1000 * vdd * c1 / (c1 + c2 + cg)
			q = edakit.Sprintf("A quiet, floating victim with %s fF to ground couples through %s fF to aggressor A and %s fF to aggressor B. Aggressor A rises from 0 to %s V while B stays quiet. What is the glitch peak on the victim, in mV?", f(cg, 1), f(c1, 1), f(c2, 1), f(vdd, 2))
		}
		return num(q, v, 1, edakit.Sprintf("Switching couplings form the top of the divider; quiet couplings act as ground and join the bottom. The total capacitance is always %s fF, and the peak is %s mV.", f(c1+c2+cg, 1), f(v, 1)), nil)
	}
	cc := float64(edakit.Int(10, 40)) / 10
	k := edakit.Pick([]float64{0.5, 0.6, 0.7, 0.8})
	before := 1000 * vdd * cc / (cc + cg)
	// A grounded shield track takes fraction k of the old coupling and converts it to ground.
	after := 1000 * vdd * (cc * (1 - k)) / (cc + cg)
	return num(edakit.Sprintf("A floating victim has %s fF to ground and %s fF of coupling to an aggressor that swings %s V. Inserting a grounded shield wire between them converts %s%% of that coupling into capacitance to the grounded shield and leaves total capacitance unchanged. By how many mV does the glitch peak fall?", f(cg, 1), f(cc, 1), f(vdd, 2), f(100*k, 0)),
		before-after, 1, edakit.Sprintf("Before: %s mV. After: the coupling is %s fF and the total stays %s fF, giving %s mV. Reduction %s mV. Shielding costs a routing track, which is why it is reserved for clocks and sensitive nets.", f(before, 1), f(cc*(1-k), 2), f(cc+cg, 1), f(after, 1), f(before-after, 1)), nil)
}

// noiseImmunity computes the glitch peak with a resistive holding driver and its margin.
func noiseImmunity(ctx generator.GeneratorContext) generator.Problem {
	vdd := edakit.Pick([]float64{0.75, 0.8, 0.9, 1.0})
	cg := float64(edakit.Int(10, 60)) / 10
	cc := float64(edakit.Int(5, 40)) / 10
	rh := float64(edakit.Int(5, 40)) / 10 // kohm
	tr := float64(5 * edakit.Int(2, 16))  // ps
	tau := rh * (cc + cg)
	vp := 1000 * vdd * (rh * cc / tr) * (1 - math.Exp(-tr/tau))
	cs := 1000 * vdd * cc / (cc + cg)
	stem := edakit.Sprintf("A quiet victim is held low by a driver of resistance %s kohm and has %s fF to ground. An aggressor coupled through %s fF ramps from 0 to %s V in %s ps. For a linear ramp, the victim voltage peaks at the end of the ramp at $V_p = V_{DD} \\frac{R C_c}{t_r}\\left(1 - e^{-t_r/\\tau}\\right)$ with $\\tau = R (C_c + C_g)$; kohm times fF is ps.", f(rh, 1), f(cg, 1), f(cc, 1), f(vdd, 2), f(tr, 0))
	switch kp(ctx) {
	case 1:
		return num(stem+" What is the glitch peak, in mV?", vp, 1,
			edakit.Sprintf("$\\tau$ = %s x %s = %s ps. $V_p$ = %s x (%s x %s / %s) x (1 - exp(-%s/%s)) = %s mV.", f(rh, 1), f(cc+cg, 1), f(tau, 2), f(1000*vdd, 0), f(rh, 1), f(cc, 1), f(tr, 0), f(tr, 0), f(tau, 2), f(vp, 1)),
			map[string]string{"tau": f(tau, 2)})
	case 2:
		th := edakit.Pick([]float64{0.25, 0.3, 0.35})
		m := 1000*vdd*th - vp
		return num(stem+edakit.Sprintf(" The receiving gate's noise immunity, for a pulse this short, is a flat threshold of %s%% of $V_{DD}$. What is the noise margin (threshold minus peak), in mV? A negative value means a glitch violation.", f(100*th, 0)), m, 1,
			edakit.Sprintf("Peak %s mV against a threshold of %s mV gives a margin of %s mV. Real noise-immunity curves let narrow pulses exceed the DC threshold, because a gate filters pulses shorter than its own delay.", f(vp, 1), f(1000*vdd*th, 0), f(m, 1)), nil)
	}
	pct := 100 * vp / cs
	return num(stem+" What is this peak as a percentage of the floating-victim charge-sharing bound $V_{DD} C_c/(C_c + C_g)$?", pct, 1,
		edakit.Sprintf("Bound %s mV, peak %s mV, ratio %s%%. The holding driver bleeds charge during the ramp; the slower the ramp relative to $\\tau$, the smaller the glitch.", f(cs, 1), f(vp, 1), f(pct, 1)), nil)
}

// timingWindows filters aggressors by switching-window overlap.
func timingWindows(ctx generator.GeneratorContext) generator.Problem {
	rd := float64(edakit.Int(5, 30)) / 10
	vs := float64(10 * edakit.Int(5, 30))
	ve := vs + float64(10*edakit.Int(2, 8))
	n := edakit.Scaled(ctx, 3, 5)
	var desc []string
	all, ov := 0.0, 0.0
	cnt := 0
	for i := 0; i < n; i++ {
		cc := float64(edakit.Int(5, 30)) / 10
		var s, e float64
		if edakit.Int(0, 1) == 0 {
			s = vs + float64(10*edakit.Int(-6, 6))
		} else {
			s = vs + float64(10*edakit.Int(-15, 15))
		}
		e = s + float64(10*edakit.Int(1, 5))
		all += cc
		if s <= ve && e >= vs {
			ov += cc
			cnt++
		}
		desc = append(desc, edakit.Sprintf("aggressor %d: %s fF, switching window %s to %s ps", i+1, f(cc, 1), f(s, 0), f(e, 0)))
	}
	stem := edakit.Sprintf("A victim switches somewhere in the window %s to %s ps. Its aggressors are %s. An aggressor can affect the victim's delay only if its window overlaps the victim's (windows that touch count as overlapping). Non-overlapping aggressors are treated as quiet.", f(vs, 0), f(ve, 0), strings.Join(desc, "; "))
	switch kp(ctx) {
	case 1:
		return num(stem+" How many aggressors survive the timing-window filter?", float64(cnt), 0,
			edakit.Sprintf("Check each window against [%s, %s]: %d overlap.", f(vs, 0), f(ve, 0), cnt), nil)
	case 2:
		d := 0.69 * rd * (all - ov)
		return num(stem+edakit.Sprintf(" The victim's driver is %s kohm, with delay $t = 0.69RC$ (kohm times fF is ps), and an overlapping aggressor is assumed to oppose the victim. How many ps of crosstalk pessimism does the window filter remove, compared with assuming every aggressor opposes?", f(rd, 1)), d, 2,
			edakit.Sprintf("Without windows the delta delay is 0.69 x %s x %s; with windows it is 0.69 x %s x %s. The difference is %s ps.", f(rd, 1), f(all, 1), f(rd, 1), f(ov, 1), f(d, 2)), nil)
	}
	d := 0.69 * rd * ov
	return num(stem+edakit.Sprintf(" The victim's driver is %s kohm, with delay $t = 0.69RC$ (kohm times fF is ps). Each surviving aggressor is assumed to oppose the victim. What crosstalk delay increase does the analysis report, relative to all aggressors quiet, in ps?", f(rd, 1)), d, 2,
		edakit.Sprintf("Only overlapping aggressors move from Miller factor 1 to 2, adding %s fF: 0.69 x %s x %s = %s ps. Windows come from the timer, and the delay change moves the windows, so SI and STA iterate until the windows converge.", f(ov, 1), f(rd, 1), f(ov, 1), f(d, 2)), nil)
}
