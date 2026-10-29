package cmos

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// setupSweep is a characterization sweep of clock-to-q against setup skew.
type setupSweep struct {
	skews []int
	tcq   []int
}

func drawSweep(ctx generator.GeneratorContext) setupSweep {
	step := edakit.Pick([]int{5, 10})
	n := edakit.Scaled(ctx, 7, 10)
	smin := step * edakit.Int(1, 3)
	t0 := edakit.Int(35, 80)
	c := float64(edakit.Int(40, 120))
	h := float64(edakit.Pick([]int{6, 8, 10, 12}))
	var sw setupSweep
	for i := n - 1; i >= 0; i-- {
		s := smin + step*i
		sw.skews = append(sw.skews, s)
		sw.tcq = append(sw.tcq, t0+int(c*pow(2, -float64(s-smin)/h)+0.5))
	}
	return sw
}

func (s setupSweep) text() string {
	var parts []string
	for i := range s.skews {
		parts = append(parts, edakit.Sprintf("%d ps: %d", s.skews[i], s.tcq[i]))
	}
	return "Setup skew (data-to-clock) versus measured clock-to-q, in ps: " + strings.Join(parts, "; ") + "."
}

// setupHoldChar reads a setup sweep: push-out limit, setup time, and minimum D-to-Q.
func setupHoldChar(ctx generator.GeneratorContext) generator.Problem {
	sw := drawSweep(ctx)
	nom := sw.tcq[0]
	pct := edakit.Pick([]int{5, 10, 15, 20})
	limit := float64(nom) * (1 + float64(pct)/100)
	pre := "A characterization run sweeps how long before the clock edge the data input changes. " + sw.text() + " "
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(limit, 2)
		q := pre + edakit.Sprintf("Take the nominal clock-to-q as the value at the largest skew. The setup criterion allows a %d%% push-out. What is the largest acceptable clock-to-q, in ps? Round to 2 decimal places.", pct)
		return edakit.Numeric(q, ans, edakit.Sprintf("The nominal value is %d ps, so the limit is $%d \\times %s = %s$ ps.", nom, nom, edakit.Num(1+float64(pct)/100, 2), ans), map[string]string{"limit_ps": ans})
	case 2:
		best := sw.skews[0]
		for i := range sw.skews {
			if float64(sw.tcq[i]) <= limit+1e-9 {
				best = sw.skews[i]
			}
		}
		q := pre + edakit.Sprintf("Take the nominal clock-to-q as the value at the largest skew. The setup time is the smallest skew in the table whose clock-to-q is within %d%% of nominal. What setup time does the library get, in ps?", pct)
		return edakit.Numeric(q, edakit.IntStr(best),
			edakit.Sprintf("The limit is %s ps. Scanning from large skew to small, the last skew whose clock-to-q stays at or below the limit is %d ps.", edakit.Num(limit, 2), best),
			map[string]string{"limit_ps": edakit.Num(limit, 2), "setup_ps": edakit.IntStr(best)})
	}
	bestS, bestSum := 0, 1<<30
	for i := range sw.skews {
		if sum := sw.skews[i] + sw.tcq[i]; sum < bestSum || (sum == bestSum && sw.skews[i] > bestS) {
			bestS, bestSum = sw.skews[i], sum
		}
	}
	q := pre + "The data-to-output delay of the flop is skew plus clock-to-q. At which skew in the table is the data-to-output delay smallest? If several tie, give the largest such skew. Answer in ps."
	return edakit.Numeric(q, edakit.IntStr(bestS),
		edakit.Sprintf("Adding each skew to its clock-to-q, the minimum sum is %d ps at a skew of %d ps. Below that point the push-out grows faster than the skew shrinks.", bestSum, bestS),
		map[string]string{"min_dq_ps": edakit.IntStr(bestSum), "skew_ps": edakit.IntStr(bestS)})
}

// metastability: synchronizer MTBF from latch regeneration.
func metastability(ctx generator.GeneratorContext) generator.Problem {
	tau := float64(edakit.Int(8, 30))
	tw := float64(edakit.Int(10, 50))
	fc := edakit.Pick([]float64{0.5, 1, 1.5, 2})
	fd := float64(edakit.Pick([]int{10, 50, 100, 200, 500}))
	rate := tw * 1e-12 * fc * 1e9 * fd * 1e6
	switch kpOf(ctx) {
	case 1:
		tr := float64(edakit.Scaled(ctx, 400, 1200))
		lg := tr/(tau*ln(10)) - log10(rate)
		ans := edakit.Num(lg, 2)
		q := edakit.Sprintf("A synchronizer flop has regeneration time constant $\\tau = %s$ ps and metastability window $T_w = %s$ ps. It samples data toggling at %s MHz with a %s GHz clock and allows %s ps of resolution time. With $\\text{MTBF} = e^{t_r/\\tau}/(T_w f_c f_d)$, what is $\\log_{10}$ of the MTBF in seconds? Round to 2 decimal places.",
			edakit.Num(tau, 0), edakit.Num(tw, 0), edakit.Num(fd, 0), edakit.Num(fc, 1), edakit.Num(tr, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$\\log_{10}\\text{MTBF} = t_r/(\\tau \\ln 10) - \\log_{10}(T_w f_c f_d) = %s - %s = %s$.", edakit.Num(tr/(tau*ln(10)), 3), edakit.Num(log10(rate), 3), ans),
			map[string]string{"log_mtbf": ans})
	case 2:
		years := float64(edakit.Pick([]int{10, 100, 1000, 10000}))
		target := years * 3.15e7
		tr := tau * ln(target*rate)
		ans := edakit.Num(tr, 1)
		q := edakit.Sprintf("A synchronizer flop has $\\tau = %s$ ps and $T_w = %s$ ps, sampling data toggling at %s MHz with a %s GHz clock. Take one year as $3.15 \\times 10^7$ s. What resolution time $t_r$, in ps, gives an MTBF of %s years, using $\\text{MTBF} = e^{t_r/\\tau}/(T_w f_c f_d)$? Round to 1 decimal place.",
			edakit.Num(tau, 0), edakit.Num(tw, 0), edakit.Num(fd, 0), edakit.Num(fc, 1), edakit.Num(years, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Solve for $t_r = \\tau \\ln(\\text{MTBF} \\cdot T_w f_c f_d) = %s \\times \\ln(%s) = %s$ ps.", edakit.Num(tau, 0), sciTex(target*rate), ans),
			map[string]string{"tr_ps": ans})
	}
	period := 1000 / fc
	over := float64(edakit.Int(40, 120))
	dt := period - over
	dec := dt / (tau * ln(10))
	ans := edakit.Num(dec, 2)
	q := edakit.Sprintf("A two-flop synchronizer clocked at %s GHz is extended to three flops. The extra stage adds one clock period minus %s ps of flop overhead to the resolution time. The flops have $\\tau = %s$ ps. By how many decades (powers of ten) does the MTBF improve? Round to 2 decimal places.", edakit.Num(fc, 1), edakit.Num(over, 0), edakit.Num(tau, 0))
	return edakit.Numeric(q, ans,
		edakit.Sprintf("The added resolution time is %s - %s = %s ps, which multiplies the MTBF by $e^{\\Delta t/\\tau}$, that is $\\Delta t/(\\tau \\ln 10) = %s$ decades.", edakit.Num(period, 1), edakit.Num(over, 0), edakit.Num(dt, 1), ans),
		map[string]string{"dt_ps": edakit.Num(dt, 1), "decades": ans})
}
