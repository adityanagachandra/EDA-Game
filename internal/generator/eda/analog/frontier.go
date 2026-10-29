package analog

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// agentLogChoice: an agent's sizing iterations, judged by SPICE at three corners; which
// iteration is the first the harness may accept?
func agentLogChoice(ctx generator.GeneratorContext) generator.Problem {
	spec := float64(5 * edakit.Int(11, 13))
	// Known-good fallback log (spec 60): iteration 2 passes at TT only, iteration 3 is first.
	fallback := [4][3]float64{{58, 55, 60}, {63, 57, 64}, {66, 61, 63}, {68, 64, 66}}
	for attempt := 0; ; attempt++ {
		if attempt >= generator.MaxDrawAttempts {
			spec = 60
		}
		var rows []string
		first := -1
		ttOnly := -1
		for it := 0; it < 4; it++ {
			tt := spec + unif(-4, 12, 1)
			ss := tt - unif(0, 10, 1)
			ff := tt + unif(-6, 3, 1)
			if attempt >= generator.MaxDrawAttempts {
				tt, ss, ff = fallback[it][0], fallback[it][1], fallback[it][2]
			}
			pass := tt >= spec && ss >= spec && ff >= spec
			if pass && first < 0 {
				first = it
			}
			if tt >= spec && !pass && ttOnly < 0 && first < 0 {
				ttOnly = it
			}
			rows = append(rows, edakit.Sprintf("iteration %d: TT %s, SS %s, FF %s", it+1, f(tt, 0), f(ss, 0), f(ff, 0)))
		}
		if first < 1 || ttOnly < 0 {
			continue
		}
		q := edakit.Sprintf("A sizing agent proposes a new op-amp sizing each iteration; the harness runs SPICE at three corners and the acceptance rule is phase margin at least %s degrees at every corner. SPICE phase margins (degrees): %s. Which is the first iteration the harness may accept?", f(spec, 0), strings.Join(rows, "; "))
		var ds []string
		for i := 0; i < 4; i++ {
			if i != first {
				ds = append(ds, edakit.Sprintf("Iteration %d", i+1))
			}
		}
		return edakit.Choice(q, edakit.Sprintf("Iteration %d", first+1),
			edakit.Sprintf("Iteration %d is the first whose worst corner meets the spec. Iteration %d passes at TT only; a nominal pass is not acceptance.", first+1, ttOnly+1), ds...)
	}
}

func surrogateSpice(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		a := unif(10, 300, 0.5)
		p := a * unif(0.75, 1.25, 0.01)
		p = edakit.Round(p, 1)
		if p == a {
			p += 2.5
		}
		e := math.Abs(p-a) / a * 100
		q := edakit.Sprintf("A neural surrogate predicts an op-amp's GBW as %s MHz for a candidate sizing; SPICE, the golden judge, reports %s MHz. What is the surrogate's absolute percentage error? Round to 2 decimal places.", f(p, 1), f(a, 1))
		return num(q, f(e, 2), edakit.Sprintf("$|%s - %s|/%s \\times 100 = %s$%%.", f(p, 1), f(a, 1), f(a, 1), f(e, 2)), nil)
	case 2:
		n := 100 * edakit.Int(5, 100)
		k := 5 * edakit.Int(1, 20)
		c := pick(5, 9, 15, 27)
		ts := unif(0.001, 0.05, 0.001)
		tS := float64(edakit.Int(5, 120))
		all := float64(n*c) * tS
		mixed := float64(n)*ts + float64(k*c)*tS
		sp := all / mixed
		q := edakit.Sprintf("An agent must screen %d candidate sizings. Full SPICE costs %s s per corner over %d corners. Instead it scores every candidate with a surrogate (%s s each) and sends only the top %d to SPICE at all corners. What is the speed-up over running SPICE on everything? Round to 2 decimal places.", n, f(tS, 0), c, f(ts, 3), k)
		return num(q, f(sp, 2), edakit.Sprintf("All-SPICE: $%d \\cdot %d \\cdot %s = %s$ s. Surrogate plus verification: $%d \\cdot %s + %d \\cdot %d \\cdot %s = %s$ s. Ratio %s.", n, c, f(tS, 0), f(all, 0), n, f(ts, 3), k, c, f(tS, 0), f(mixed, 2), f(sp, 2)), nil)
	default:
		spec := float64(5 * edakit.Int(11, 13))
		gb := float64(edakit.Int(0, 4))
		var rows []string
		fa := 0
		for i := 0; i < 8; i++ {
			s := spec + unif(-6, 10, 1)
			t := s + unif(-8, 4, 1)
			if s >= spec+gb && t < spec {
				fa++
			}
			rows = append(rows, edakit.Sprintf("C%d surrogate %s, SPICE %s", i+1, f(s, 0), f(t, 0)))
		}
		q := edakit.Sprintf("An agent's surrogate predicts phase margin, and the agent forwards a candidate as passing when the prediction is at least %s degrees (the %s-degree spec plus a %s-degree guard band). Predicted and SPICE phase margins in degrees: %s. How many forwarded candidates does SPICE reject?", f(spec+gb, 0), f(spec, 0), f(gb, 0), strings.Join(rows, "; "))
		return num(q, edakit.IntStr(fa), edakit.Sprintf("Count candidates with prediction at least %s but SPICE below %s: %d. These false accepts are why SPICE stays the judge; a wider guard band trades false accepts for false rejects.", f(spec+gb, 0), f(spec, 0), fa), nil)
	}
}

func migration(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		gm := float64(50 * edakit.Int(4, 60))
		gmid := float64(edakit.Int(8, 20))
		da := unif(3, 30, 0.1)
		dbn := da * unif(1.2, 3, 0.1)
		dbn = edakit.Round(dbn, 2)
		wa := gm / gmid / da
		wb := gm / gmid / dbn
		q := edakit.Sprintf("A migration agent ports an amplifier by holding each device's $g_m/I_D$ and $g_m$ fixed. One device needs $g_m = %s\\ \\mu$S at $g_m/I_D = %s\\ \\mathrm{V}^{-1}$. Its lookup tables give a current density at that $g_m/I_D$ of %s $\\mu$A/$\\mu$m in the old process and %s $\\mu$A/$\\mu$m in the new one. What width does the device get in the new process, in micrometres? Round to 3 decimal places.", f(gm, 0), f(gmid, 0), f(da, 1), f(dbn, 2))
		return num(q, f(wb, 3), edakit.Sprintf("$I_D = %s/%s = %s\\ \\mu$A is unchanged. $W = I_D/(I_D/W) = %s/%s = %s\\ \\mu$m (it was %s).", f(gm, 0), f(gmid, 0), f(gm/gmid, 3), f(gm/gmid, 3), f(dbn, 2), f(wb, 3), f(wa, 3)), nil)
	case 2:
		ga := unif(20, 80, 1)
		gb := ga * unif(0.3, 0.9, 0.05)
		gb = edakit.Round(gb, 1)
		d := db20(gb / ga)
		q := edakit.Sprintf("At the chosen $g_m/I_D$ and channel length, a device's intrinsic gain $g_m/g_{ds}$ is %s in the old process and %s in the new one. A gain stage limited by this device's $r_o$ is ported with $g_m/I_D$ held. By how many dB does its gain change? Give a negative number for a loss. Round to 2 decimal places.", f(ga, 0), f(gb, 1))
		return num(q, f(d, 2), edakit.Sprintf("The stage gain scales with $g_m/g_{ds}$: $20\\log_{10}(%s/%s) = %s$ dB. The agent must add a cascode or lengthen L to recover it, then re-verify in SPICE.", f(gb, 1), f(ga, 0), f(d, 2)), nil)
	default:
		vdd := pick(0.8, 0.9, 1.0, 1.2)
		k := edakit.Int(2, 4)
		sw := unif(0.2, 0.5, 0.05)
		for vdd-sw <= 0.1*float64(k) {
			sw -= 0.05
		}
		g := 2 * float64(k) / (vdd - sw)
		q := edakit.Sprintf("After migration the supply drops to %s V. An output branch stacks %d saturated devices, each needing $V_{DS,sat} \\approx 2/(g_m/I_D)$, and the output must still swing %s V. Assuming the stacked devices share one $g_m/I_D$, what is the minimum $g_m/I_D$ they need, in $\\mathrm{V}^{-1}$? Round to 2 decimal places.", f(vdd, 1), k, f(sw, 2))
		return num(q, f(g, 2), edakit.Sprintf("Headroom requires $k \\cdot 2/(g_m/I_D) \\le V_{DD} - V_{swing} = %s$ V, so $g_m/I_D \\ge 2 \\cdot %d/%s = %s$. Low-voltage ports push devices toward moderate inversion.", f(vdd-sw, 2), k, f(vdd-sw, 2), f(g, 2)), nil)
	}
}

type specSet struct{ gain, gbw, pm float64 }

func violation(m, s specSet) float64 {
	v := 0.0
	for _, p := range [][2]float64{{m.gain, s.gain}, {m.gbw, s.gbw}, {m.pm, s.pm}} {
		if p[0] < p[1] {
			v += (p[1] - p[0]) / p[1]
		}
	}
	return v
}

func rewardFor(m specSet, pw float64, s specSet, pref float64) float64 {
	if v := violation(m, s); v > 0 {
		return -v
	}
	return 1 - pw/pref
}

func reward(ctx generator.GeneratorContext) generator.Problem {
	s := specSet{float64(edakit.Int(60, 80)), float64(10 * edakit.Int(2, 10)), 60}
	pref := unif(1, 5, 0.5)
	rule := edakit.Sprintf("A tool-grounded reward for a sizing policy uses specs gain at least %s dB, GBW at least %s MHz and PM at least %s degrees. If any spec fails, the reward is minus the sum of the normalized shortfalls $(\\mathrm{spec} - \\mathrm{measured})/\\mathrm{spec}$ over the failing specs; if all pass, it is $1 - P/P_{ref}$ with $P_{ref} = %s$ mW.", f(s.gain, 0), f(s.gbw, 0), f(s.pm, 0), f(pref, 1))
	switch whichKP(ctx) {
	case 1:
		var m specSet
		for attempt := 0; ; attempt++ {
			m = specSet{s.gain + unif(-8, 4, 1), s.gbw * unif(0.6, 1.2, 0.05), unif(45, 70, 1)}
			if attempt >= generator.MaxDrawAttempts {
				m = specSet{s.gain - 3, s.gbw * 0.8, 55}
			}
			if violation(m, s) > 0 {
				break
			}
		}
		pw := unif(0.2, pref, 0.05)
		r := rewardFor(m, pw, s, pref)
		q := rule + edakit.Sprintf(" SPICE measures gain %s dB, GBW %s MHz, PM %s degrees and power %s mW. What is the reward? Round to 4 decimal places.", f(m.gain, 0), f(m.gbw, 1), f(m.pm, 0), f(pw, 2))
		return num(q, f(r, 4), edakit.Sprintf("At least one spec fails, so power is ignored and the reward is -%s.", f(-r, 4)), nil)
	case 2:
		m := specSet{s.gain + unif(0, 8, 1), s.gbw * unif(1, 1.5, 0.05), unif(60, 80, 1)}
		pw := unif(0.2, pref, 0.05)
		r := rewardFor(m, pw, s, pref)
		q := rule + edakit.Sprintf(" SPICE measures gain %s dB, GBW %s MHz, PM %s degrees and power %s mW. What is the reward? Round to 4 decimal places.", f(m.gain, 0), f(m.gbw, 1), f(m.pm, 0), f(pw, 2))
		return num(q, f(r, 4), edakit.Sprintf("Every spec passes, so the reward is $1 - %s/%s = %s$.", f(pw, 2), f(pref, 1), f(r, 4)), nil)
	default:
		pw := unif(0.2, pref, 0.05)
		sur := specSet{s.gain + unif(0, 6, 1), s.gbw * unif(1, 1.4, 0.05), unif(60, 75, 1)}
		var sp specSet
		for attempt := 0; ; attempt++ {
			sp = specSet{sur.gain - unif(0, 6, 1), edakit.Round(sur.gbw*unif(0.7, 1, 0.05), 2), sur.pm - unif(0, 12, 1)}
			if attempt >= generator.MaxDrawAttempts {
				sp = specSet{sur.gain - 2, edakit.Round(sur.gbw*0.7, 2), sur.pm - 1}
			}
			if violation(sp, s) > 0 {
				break
			}
		}
		rs := rewardFor(sur, pw, s, pref)
		rt := rewardFor(sp, pw, s, pref)
		q := rule + edakit.Sprintf(" During training the policy was scored by a fast surrogate, which reported gain %s dB, GBW %s MHz, PM %s degrees at %s mW for its favourite design. SPICE reports gain %s dB, GBW %s MHz, PM %s degrees at the same power. By how much does the surrogate reward exceed the SPICE reward? Round to 4 decimal places.", f(sur.gain, 0), f(sur.gbw, 2), f(sur.pm, 0), f(pw, 2), f(sp.gain, 0), f(sp.gbw, 2), f(sp.pm, 0))
		return num(q, f(rs-rt, 4), edakit.Sprintf("Surrogate reward %s (all specs pass), SPICE reward %s. The gap of %s is reward the policy collected by exploiting the surrogate's error, the signature of reward hacking.", f(rs, 4), f(rt, 4), f(rs-rt, 4)), nil)
	}
}
