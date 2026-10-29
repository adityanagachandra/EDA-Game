package fpga

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func seedRuns(center float64, n, spread int) ([]float64, float64) {
	var v []float64
	s := 0.0
	for i := 0; i < n; i++ {
		x := center + float64(edakit.Int(-spread, spread))/2
		v = append(v, x)
		s += x
	}
	return v, s / float64(n)
}

func fl(v []float64, p int) string {
	var s []string
	for _, x := range v {
		s = append(s, edakit.Num(x, p))
	}
	return strings.Join(s, ", ")
}

func optionTuning(ctx generator.GeneratorContext) generator.Problem {
	tool := edakit.Pick([]string{"nextpnr", "VPR"})
	switch whichKP(ctx) {
	case 1:
		n := edakit.Int(3, 6)
		v, m := seedRuns(float64(edakit.Int(80, 300)), n, 20)
		q := edakit.Sprintf("An agent tunes %s options for one design. Because placement is seeded, it runs one candidate configuration with %d placement seeds and records Fmax values of %s MHz. What mean Fmax does it log for the configuration? Round to 2 decimal places.", tool, n, fl(v, 1))
		return num(q, edakit.Num(m, 2),
			edakit.Sprintf("Mean = %s MHz. A single seed would confuse configuration effects with seed noise.", edakit.Num(m, 2)),
			map[string]string{"mean": edakit.Num(m, 2)})
	case 2:
		c := float64(edakit.Int(100, 250))
		b, bm := seedRuns(c, 3, 16)
		var lines []string
		best := math.Inf(-1)
		var ms []string
		for i := 0; i < 3; i++ {
			v, m := seedRuns(c+float64(edakit.Int(-6, 10)), 3, 16)
			lines = append(lines, edakit.Sprintf("configuration %d: %s", i+1, fl(v, 1)))
			ms = append(ms, edakit.Num(m, 2))
			if m > best {
				best = m
			}
		}
		d := best - bm
		q := edakit.Sprintf("An agent compares %s option sets on three seeds each. Baseline Fmax: %s MHz; %s MHz. By how much does the best configuration's mean Fmax exceed the baseline mean, in MHz? Round to 2 decimal places (negative if none is better).", tool, fl(b, 1), strings.Join(lines, "; "))
		return num(q, edakit.Num(d, 2),
			edakit.Sprintf("Baseline mean %s; configuration means %s. The best exceeds the baseline by %s MHz.", edakit.Num(bm, 2), strings.Join(ms, ", "), edakit.Num(d, 2)),
			map[string]string{"baseline": edakit.Num(bm, 2)})
	default:
		c := float64(edakit.Int(100, 250))
		b, bm := seedRuns(c, 4, 12)
		lo, hi := b[0], b[0]
		for _, x := range b {
			lo, hi = math.Min(lo, x), math.Max(hi, x)
		}
		rg := hi - lo
		k := 0
		var lines, ms []string
		for i := 0; i < 4; i++ {
			v, m := seedRuns(c+float64(edakit.Int(-4, 16)), 4, 12)
			lines = append(lines, edakit.Sprintf("configuration %d: %s", i+1, fl(v, 1)))
			ms = append(ms, edakit.Num(m, 2))
			if m-bm > rg {
				k++
			}
		}
		q := edakit.Sprintf("Before promoting a %s configuration, an agent requires its mean Fmax to beat the baseline mean by more than the baseline's seed-to-seed range (max minus min). Baseline over four seeds: %s MHz. Candidates over four seeds: %s MHz. How many candidates pass?", tool, fl(b, 1), strings.Join(lines, "; "))
		return num(q, edakit.IntStr(k),
			edakit.Sprintf("Baseline mean %s, range %s. Candidate means %s. %d exceed %s + %s. The noise gate keeps the agent from promoting a lucky seed.", edakit.Num(bm, 2), edakit.Num(rg, 1), strings.Join(ms, ", "), k, edakit.Num(bm, 2), edakit.Num(rg, 1)),
			map[string]string{"range": edakit.Num(rg, 1)})
	}
}

func fastTiming(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 8)
	var est, act []int
	var lines []string
	for i := 0; i < n; i++ {
		a := 10 * edakit.Int(-15, 40)
		e := a + 10*edakit.Int(-8, 8)
		est = append(est, e)
		act = append(act, a)
		lines = append(lines, edakit.Sprintf("P%d estimated %d, routed %d", i+1, e, a))
	}
	head := "A placement-stage timing estimator is compared with post-route static timing on the same paths (setup slack in ps): " + strings.Join(lines, "; ") + "."
	switch whichKP(ctx) {
	case 1:
		s := 0
		for i := range est {
			s += absInt(est[i] - act[i])
		}
		m := float64(s) / float64(n)
		return num(head+" What is the estimator's mean absolute error, in ps? Round to 1 decimal place.", edakit.Num(m, 1),
			edakit.Sprintf("Sum of |estimated - routed| = %d over %d paths, so MAE = %s ps.", s, n, edakit.Num(m, 1)),
			map[string]string{"abs_sum": edakit.IntStr(s)})
	case 2:
		g := 10 * edakit.Int(3, 10)
		k := 0
		for _, e := range est {
			if e < g {
				k++
			}
		}
		return num(head+edakit.Sprintf(" An agent trusts the fast estimate only when its slack is at least a guard band of %d ps, and sends every other path to the accurate post-route timer. Using only the estimates, how many paths does it send?", g), edakit.IntStr(k),
			edakit.Sprintf("Count estimated slacks below %d ps: %d. The guard band trades accurate-timer run time against the risk of trusting a wrong estimate.", g, k),
			map[string]string{"sent": edakit.IntStr(k)})
	default:
		k := 0
		for i := range est {
			if est[i] >= 0 && act[i] < 0 {
				k++
			}
		}
		return num(head+" How many paths does the estimator falsely pass, with estimated slack of 0 or more but negative routed slack?", edakit.IntStr(k),
			edakit.Sprintf("False passes: %d. These are the dangerous errors; a pessimistic error only costs extra work, while a false pass reaches signoff as a violation.", k),
			map[string]string{"false_pass": edakit.IntStr(k)})
	}
}
