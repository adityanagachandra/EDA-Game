package stat

import (
	"math"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// monteCarlo: standard error of a yield estimate (KP 1), 95% interval half-width (KP 2), and
// the sample count for a target standard error of a mean (KP 3).
func monteCarlo(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1, 2:
		n := edakit.Pick([]int{200, 400, 500, 1000, 2000, 5000})
		fail := edakit.Int(2, n/8)
		p := float64(n-fail) / float64(n)
		se := math.Sqrt(p * (1 - p) / float64(n))
		stem := edakit.Sprintf("A Monte Carlo run simulates %d random instances of a circuit and %d of them fail a spec.", n, fail)
		if ctx.KPOrdinal == 1 || (ctx.KPOrdinal == 0 && edakit.Int(0, 1) == 0) {
			ans := edakit.Num(100*se, 3)
			return edakit.Numeric(stem+" What is the standard error of the estimated yield, in percentage points? Round to 3 decimal places.",
				ans, edakit.Sprintf("$\\hat{Y} = %d/%d = %s$. $SE = \\sqrt{\\hat{Y}(1-\\hat{Y})/N} = %s$, which is %s percentage points.", n-fail, n, f(p, 4), f(se, 5), ans),
				map[string]string{"yield": f(p, 4)})
		}
		ans := edakit.Num(100*1.96*se, 2)
		return edakit.Numeric(stem+" What is the half-width of the 95 percent normal-approximation confidence interval on the yield, in percentage points? Round to 2 decimal places.",
			ans, edakit.Sprintf("$\\hat{Y} = %s$, $SE = %s$, half-width $= 1.96 \\times SE = %s$ percentage points.", f(p, 4), f(se, 5), ans),
			map[string]string{"se": f(se, 5)})
	default:
		n0 := edakit.Pick([]int{100, 200, 250, 500})
		s := float64(edakit.Int(5, 40)) / 10
		target := float64(edakit.Int(2, 10)) / 100
		need := int(math.Ceil((s / target) * (s / target)))
		return edakit.Numeric(edakit.Sprintf("A pilot Monte Carlo run of %d samples estimates the standard deviation of a delay as %s ps. How many samples in total are needed for the standard error of the mean delay to be at most %s ps? Give the smallest integer.", n0, f(s, 1), f(target, 2)),
			edakit.IntStr(need), edakit.Sprintf("$SE = s/\\sqrt{N} \\le %s$ gives $N \\ge (s/SE)^2 = (%s/%s)^2$, so N = %d. The error falls only as $1/\\sqrt{N}$.", f(target, 2), f(s, 1), f(target, 2), need),
			map[string]string{"ratio": f(s/target, 3)})
	}
}

// sampleSize: samples for a relative error at a confidence (KP 1), zero-failure demonstration
// (KP 2), and naive Monte Carlo cost at a sigma level (KP 3).
func sampleSize(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		p := edakit.Pick([]float64{0.01, 0.005, 0.002, 0.001, 0.0005, 0.0001})
		eps := edakit.Pick([]float64{0.05, 0.1, 0.2})
		z := edakit.Pick([]float64{1.645, 1.96, 2.576})
		conf := map[float64]string{1.645: "90", 1.96: "95", 2.576: "99"}[z]
		n := int(math.Ceil(z * z * (1 - p) / (p * eps * eps)))
		return edakit.Numeric(edakit.Sprintf("A failure rate near $p = %s$ is to be estimated by plain Monte Carlo with a relative error of at most %s percent at %s percent confidence ($z = %s$). Using $N = z^2(1-p)/(p\\,\\varepsilon^2)$, how many samples are needed? Give the smallest integer.", f(p, 4), f(100*eps, 0), conf, f(z, 3)),
			edakit.IntStr(n), edakit.Sprintf("$N = %s^2 (1-%s)/(%s \\cdot %s^2) = %d$.", f(z, 3), f(p, 4), f(p, 4), f(eps, 2), n),
			map[string]string{"n": edakit.IntStr(n)})
	case 2:
		p0 := edakit.Pick([]float64{0.01, 0.001, 0.0001, 0.00001})
		c := edakit.Pick([]float64{0.9, 0.95, 0.99})
		n := int(math.Ceil(math.Log(1-c) / math.Log1p(-p0)))
		return edakit.Numeric(edakit.Sprintf("You will claim a failure rate below %s with %s percent confidence if a Monte Carlo run shows zero failures. What is the smallest number of samples that supports the claim? Give an integer.", f(p0, 5), f(100*c, 0)),
			edakit.IntStr(n), edakit.Sprintf("Zero failures in $N$ samples has probability $(1-p_0)^N$. Requiring $(1-p_0)^N \\le %s$ gives $N \\ge \\ln(%s)/\\ln(1-%s) = %d$. At 95 percent this is the rule of three, $N \\approx 3/p_0$.", f(1-c, 2), f(1-c, 2), f(p0, 5), n),
			map[string]string{"n": edakit.IntStr(n)})
	default:
		k := edakit.Pick([]float64{4, 4.5, 5, 5.5, 6})
		eps := edakit.Pick([]float64{0.1, 0.2})
		p := tail(k)
		n := 1.96 * 1.96 * (1 - p) / (p * eps * eps)
		unit, scale := "millions", 1e6
		if k >= 5.5 {
			unit, scale = "billions", 1e9
		}
		ans := edakit.Num(n/scale, 1)
		return edakit.Numeric(edakit.Sprintf("An SRAM bitcell fails at the %s-sigma tail of a Gaussian margin. How many plain Monte Carlo samples, in %s, are needed to estimate its failure rate to %s percent relative error at 95 percent confidence ($z = 1.96$)? Round to 1 decimal place.", f(k, 1), unit, f(100*eps, 0)),
			ans, edakit.Sprintf("$p = Q(%s) = %s$. $N = 1.96^2(1-p)/(p\\,\\varepsilon^2) = %s$, which is %s %s of simulations.", f(k, 1), sci(p), sci(n), ans, unit),
			map[string]string{"p": sci(p)})
	}
}

// radicalInverse returns the van der Corput radical inverse of n in base b.
func radicalInverse(n, b int) float64 {
	x, scale := 0.0, 1.0/float64(b)
	for n > 0 {
		x += float64(n%b) * scale
		n /= b
		scale /= float64(b)
	}
	return x
}

// quasiMC: base-2 radical inverse (KP 1), a Halton point coordinate in base 3 or 5 (KP 2), and
// sample counts under the two convergence rates (KP 3).
func quasiMC(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 3, 60)
		x := radicalInverse(n, 2)
		ans := edakit.Num(x, 6)
		return edakit.Numeric(edakit.Sprintf("What is point number %d of the base-2 van der Corput sequence (the radical inverse of %d in base 2)? Round to 6 decimal places.", n, n),
			ans, edakit.Sprintf("$%d = %s_2$. Mirroring the digits about the radix point gives $0.%s_2 = %s$.", n, strconv.FormatInt(int64(n), 2), reverse(strconv.FormatInt(int64(n), 2)), ans),
			map[string]string{"binary": strconv.FormatInt(int64(n), 2)})
	case 2:
		b := edakit.Pick([]int{3, 5})
		n := edakit.Scaled(ctx, 4, 80)
		x := radicalInverse(n, b)
		ans := edakit.Num(x, 4)
		dim := map[int]string{3: "second", 5: "third"}[b]
		return edakit.Numeric(edakit.Sprintf("A Halton sequence uses the prime bases 2, 3, 5 for its first three coordinates. What is the %s coordinate (base %d) of point number %d? Round to 4 decimal places.", dim, b, n),
			ans, edakit.Sprintf("$%d = %s_{%d}$. Reversing the digits after the radix point gives $0.%s_{%d} = %s$.", n, strconv.FormatInt(int64(n), b), b, reverse(strconv.FormatInt(int64(n), b)), b, ans),
			map[string]string{"digits": strconv.FormatInt(int64(n), b)})
	default:
		n := edakit.Pick([]int{256, 500, 1000, 1024, 2000, 4096})
		r := edakit.Int(2, 10)
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(edakit.Sprintf("A plain Monte Carlo estimate from %d samples has error that shrinks as $N^{-1/2}$. How many samples in total are needed to cut the error by a factor of %d?", n, r),
				edakit.IntStr(n*r*r), edakit.Sprintf("Error $\\propto N^{-1/2}$, so a factor %d needs $%d^2 = %d$ times the samples: %d.", r, r, r*r, n*r*r), map[string]string{"factor": edakit.IntStr(r * r)})
		}
		return edakit.Numeric(edakit.Sprintf("A quasi-Monte Carlo estimate from %d low-discrepancy points has error that, ignoring logarithmic factors, shrinks as $N^{-1}$. How many points in total are needed to cut the error by a factor of %d?", n, r),
			edakit.IntStr(n*r), edakit.Sprintf("Error $\\propto N^{-1}$, so a factor %d needs %d times the points: %d. Plain Monte Carlo would need %d.", r, r, n*r, n*r*r), map[string]string{"factor": edakit.IntStr(r)})
	}
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// importanceSampling: weight of one sample under a mean shift (KP 1), an estimate from a few
// weighted samples (KP 2), and a two-dimensional weight (KP 3).
func importanceSampling(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		mu := float64(edakit.Int(4, 7)) / 2
		x := mu + float64(edakit.Int(-4, 4))/4
		w := math.Exp(-mu*x + mu*mu/2)
		pl := sigPlaces(w)
		ans := edakit.Num(w, pl)
		return edakit.Numeric(edakit.Sprintf("A failure rate under $X \\sim N(0,1)$ is estimated by drawing from the shifted proposal $N(%s, 1)$. What is the likelihood-ratio weight of the sample $x = %s$? %s", f(mu, 1), f(x, 2), roundTo(pl)),
			ans, edakit.Sprintf("$w(x) = \\varphi(x)/\\varphi(x-\\mu) = e^{-\\mu x + \\mu^2/2} = e^{-%s + %s} = %s$.", f(mu*x, 4), f(mu*mu/2, 3), ans), map[string]string{"log_w": f(-mu*x+mu*mu/2, 4)})
	case 2:
		mu := float64(edakit.Int(6, 9)) / 2
		t := mu
		n := edakit.Scaled(ctx, 4, 6)
		xs := make([]float64, n)
		var parts []string
		sum := 0.0
		for i := range xs {
			xs[i] = mu + float64(edakit.Int(-6, 6))/4
			parts = append(parts, f(xs[i], 2))
			if xs[i] > t {
				sum += math.Exp(-mu*xs[i] + mu*mu/2)
			}
		}
		if sum == 0 {
			xs[0] = t + 0.25
			parts[0] = f(xs[0], 2)
			sum = math.Exp(-mu*xs[0] + mu*mu/2)
		}
		est := sum / float64(n) * 1e6
		pl := sigPlaces(est)
		ans := edakit.Num(est, pl)
		return edakit.Numeric(edakit.Sprintf("To estimate $P(X > %s)$ for $X \\sim N(0,1)$, you draw %d samples from $N(%s, 1)$: %s. What is the importance-sampling estimate, in parts per million? %s", f(t, 1), n, f(mu, 1), strings.Join(parts, ", "), roundTo(pl)),
			ans, edakit.Sprintf("Keep samples above %s, weight each by $e^{-\\mu x + \\mu^2/2}$, and divide the weighted sum %s by $N = %d$: %s ppm. The exact value is %s ppm.", f(t, 1), sciT(sum), n, ans, f(tail(t)*1e6, 3)),
			map[string]string{"weighted_sum": sci(sum)})
	default:
		m1, m2 := float64(edakit.Int(2, 5))/2, float64(edakit.Int(2, 5))/2
		x1, x2 := m1+float64(edakit.Int(-2, 2))/4, m2+float64(edakit.Int(-2, 2))/4
		lw := -(m1*x1 + m2*x2) + (m1*m1+m2*m2)/2
		w := math.Exp(lw)
		pl := sigPlaces(w)
		ans := edakit.Num(w, pl)
		return edakit.Numeric(edakit.Sprintf("Two independent standard-normal process parameters are sampled from a proposal whose mean is shifted to $(%s, %s)$ with unit variances. What is the importance weight of the sample $(%s, %s)$? %s", f(m1, 1), f(m2, 1), f(x1, 2), f(x2, 2), roundTo(pl)),
			ans, edakit.Sprintf("$w = e^{-\\mu\\cdot x + \\|\\mu\\|^2/2} = e^{%s} = %s$.", f(lw, 4), ans), map[string]string{"log_w": f(lw, 4)})
	}
}

// mpfp: reliability index of a linear boundary (KP 1), its failure probability (KP 2), and the
// index with unequal parameter sigmas (KP 3).
func mpfp(ctx generator.GeneratorContext) generator.Problem {
	a1, a2 := float64(edakit.Int(1, 9)), float64(edakit.Int(1, 9))
	norm := math.Hypot(a1, a2)
	switch kp(ctx) {
	case 1:
		b := math.Round(norm * float64(edakit.Int(20, 60)) / 10)
		beta := b / norm
		ans := edakit.Num(beta, 3)
		return edakit.Numeric(edakit.Sprintf("In standard-normal parameter space a cell fails when $%s x_1 + %s x_2 > %s$. How far from the origin, in sigmas, is the most probable failure point? Round to 3 decimal places.", f(a1, 0), f(a2, 0), f(b, 0)),
			ans, edakit.Sprintf("The closest point of the plane $a\\cdot x = b$ is at distance $\\beta = b/\\|a\\| = %s/%s = %s$, along the direction of $a$.", f(b, 0), f(norm, 4), ans), map[string]string{"norm_a": f(norm, 4)})
	case 2:
		beta := float64(edakit.Int(30, 55)) / 10
		b := beta * norm
		p := tail(b / norm)
		ppm := p * 1e6
		pl := sigPlaces(ppm)
		ans := edakit.Num(ppm, pl)
		return edakit.Numeric(edakit.Sprintf("In standard-normal space a failure region is $%s x_1 + %s x_2 > %s$. What is its exact probability, in parts per million? %s", f(a1, 0), f(a2, 0), f(b, 3), roundTo(pl)),
			ans, edakit.Sprintf("$a\\cdot x \\sim N(0, \\|a\\|^2)$, so $P = Q(b/\\|a\\|) = Q(%s) = %s$, which is %s ppm. For a linear boundary the first-order estimate is exact.", f(b/norm, 3), sci(p), ans), map[string]string{"beta": f(b/norm, 3)})
	default:
		s1, s2 := float64(edakit.Int(2, 12)), float64(edakit.Int(2, 12))
		b := float64(edakit.Int(20, 200))
		den := math.Hypot(a1*s1, a2*s2)
		beta := b / den
		ans := edakit.Num(beta, 3)
		return edakit.Numeric(edakit.Sprintf("Two independent zero-mean parameters have $\\sigma_1 = %s$ mV and $\\sigma_2 = %s$ mV. A cell fails when $%s \\Delta V_1 + %s \\Delta V_2 > %s$ mV. After whitening to standard-normal coordinates, what is the reliability index $\\beta$? Round to 3 decimal places.", f(s1, 0), f(s2, 0), f(a1, 0), f(a2, 0), f(b, 0)),
			ans, edakit.Sprintf("Substitute $\\Delta V_i = \\sigma_i x_i$: the boundary becomes $%s x_1 + %s x_2 > %s$, so $\\beta = %s/\\sqrt{%s^2+%s^2} = %s$.", f(a1*s1, 0), f(a2*s2, 0), f(b, 0), f(b, 0), f(a1*s1, 0), f(a2*s2, 0), ans), map[string]string{"whitened_norm": f(den, 4)})
	}
}

// rareEventBudgetChoice asks how many plain Monte Carlo samples a k-sigma failure needs for a
// 10 percent relative error at 95 percent confidence, with the classic wrong counts as options.
func rareEventBudgetChoice(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Pick([]float64{3.5, 4, 4.5, 5, 5.5, 6, 6.5})
	eps := edakit.Pick([]float64{0.1, 0.05})
	p := tail(k)
	n := 1.96 * 1.96 / (p * eps * eps)
	q := edakit.Sprintf("A bitcell fails at the %s-sigma tail ($p \\approx %s$). Roughly how many plain Monte Carlo simulations does a %s percent relative error at 95 percent confidence need?", f(k, 1), sci(p), f(100*eps, 0))
	return edakit.Choice(q, "About "+sciT(n),
		edakit.Sprintf("Relative standard error is $\\sqrt{(1-p)/(Np)}$, so $N \\approx z^2/(p\\,\\varepsilon^2) = %s$. That is why rare-event methods exist.", sci(n)),
		"About "+sciT(1/p), "About "+sciT(1.96/(p*eps)), "About "+sciT(1/(p*eps*eps)), "About "+sciT(1.96*1.96/p))
}

var rareEventItems = []edakit.Item{
	{KP: 2, Q: "Which method shifts the sampling distribution toward the failure region and corrects the bias with likelihood-ratio weights?",
		A: "Importance sampling", D: []string{"Latin hypercube sampling", "Quasi-Monte Carlo", "Statistical blockade"},
		Why: "Importance sampling draws from a proposal that hits failures often and multiplies each sample by the ratio of the true to the proposal density."},
	{KP: 2, Q: "A method trains a cheap classifier to skip simulating samples that are clearly far from the tail, simulates only the rest, and fits a generalized Pareto distribution to the exceedances. What is it?",
		A: "Statistical blockade", D: []string{"Scaled-sigma sampling", "Subset simulation", "Mean-shift importance sampling"},
		Why: "Statistical blockade combines a blocking classifier with extreme value theory to model the tail beyond a threshold."},
	{KP: 2, Q: "A flow inflates every parameter sigma by factors 2, 2.5 and 3, measures the (now frequent) failure rate at each, and extrapolates a model of failure rate versus scale back to 1. What is it?",
		A: "Scaled-sigma sampling", D: []string{"Importance sampling with a mixture proposal", "Subset simulation", "Quasi-Monte Carlo"},
		Why: "Scaled-sigma sampling makes failures common by widening the distribution and relies on a fitted model of how the failure probability depends on the scale."},
	{KP: 2, Q: "A method writes a tiny probability as a product of larger conditional probabilities of nested intermediate failure events and estimates each with Markov chain sampling. What is it?",
		A: "Subset simulation (splitting)", D: []string{"Statistical blockade", "Scaled-sigma sampling", "Plain Monte Carlo with more seeds"},
		Why: "Subset simulation turns one rare event into a chain of conditional events, each with probability around 0.1, which is cheap to estimate."},
	{KP: 2, Q: "Before importance sampling an SRAM cell, an optimizer finds the point on the failure boundary nearest the origin in standard-normal space. Why centre the proposal there?",
		A: "It is the most probable failure point, so samples near it carry most of the failure probability", D: []string{"It is the mean of the failure region", "It makes all weights equal to one", "It is where the circuit simulator converges fastest"},
		Why: "The Gaussian density on the failure region peaks at its closest point to the origin. A proposal centred there produces failures whose weights are not vanishingly small."},
	{KP: 3, Q: "A bitcell can fail by read upset when one device is weak or by write failure when another is strong, in opposite directions of parameter space. A single mean-shift importance sampler targets the read region. What goes wrong?",
		A: "The write-failure region is almost never sampled, so the total failure rate is underestimated", D: []string{"The estimate is unbiased and only slower", "The weights of read failures become negative", "The estimate double-counts read failures"},
		Why: "Importance sampling is unbiased only if the proposal covers every failure region. A proposal centred on one region leaves the other with nearly zero samples."},
	{KP: 3, Q: "An engineer counts the failing importance samples and divides by N, forgetting the weights. What is the effect?",
		A: "The estimate is wildly too high, because it reports the proposal's failure rate", D: []string{"The estimate is slightly too low", "No effect, since the weights average to one", "The estimate becomes zero"},
		Why: "Under the shifted proposal failures are common by design. Only the weights bring each failure back to its true, tiny probability."},
	{KP: 3, Q: "Importance sampling in 200 dimensions gives weights where one sample holds 99 percent of the total. What does that indicate?",
		A: "Weight degeneracy: the effective sample size has collapsed and the estimate is unreliable", D: []string{"An excellent proposal with low variance", "The failure rate is about 0.99", "The estimate has converged"},
		Why: "Effective sample size, $(\\sum w)^2/\\sum w^2$, near one means the answer rests on a single sample, a common failure of naive proposals in high dimension."},
	{KP: 3, Q: "An importance-sampling run reports a narrow confidence interval. Why is that not proof the estimate is right?",
		A: "The sample variance is computed from the samples drawn, so a missed failure region leaves no trace in it", D: []string{"Confidence intervals do not apply to weighted samples", "Narrow intervals always mean too few samples", "The interval is always exact for Gaussians"},
		Why: "The variance estimate only sees regions the proposal visited. Golden checks against known cases and multiple-region diagnostics are needed."},
	{KP: 3, Q: "In statistical blockade the classifier threshold is set to block the 97 percent of samples least likely to be in the tail, though the tail of interest starts at the 99th percentile. Why the gap?",
		A: "A relaxed threshold keeps the classifier's errors from silently discarding true tail samples", D: []string{"To make the generalized Pareto fit exact", "To reduce the number of tail samples", "Because the classifier must block exactly 97 percent"},
		Why: "Misclassified tail samples are lost forever. Classifying against a looser threshold than the one used for the tail model buys a safety margin."},
}

// extremeValue: GPD exceedance (KP 1), the exponential (xi = 0) tail (KP 2), and a return
// level (KP 3).
func extremeValue(ctx generator.GeneratorContext) generator.Problem {
	pu := edakit.Pick([]float64{0.01, 0.02, 0.05})
	u := float64(edakit.Int(40, 80))
	sg := float64(edakit.Int(3, 10))
	switch kp(ctx) {
	case 1:
		xi := float64(edakit.Int(1, 4)) / 10
		y := float64(edakit.Int(5, 40))
		p := pu * math.Pow(1+xi*y/sg, -1/xi)
		ppm := p * 1e6
		pl := sigPlaces(ppm)
		ans := edakit.Num(ppm, pl)
		return edakit.Numeric(edakit.Sprintf("Simulated read delays above a threshold $u = %s$ ps occur with probability %s. The excesses over $u$ fit a generalized Pareto distribution with scale $\\sigma = %s$ ps and shape $\\xi = %s$. What is the probability that a delay exceeds %s ps, in parts per million? %s", f(u, 0), f(pu, 2), f(sg, 0), f(xi, 1), f(u+y, 0), roundTo(pl)),
			ans, edakit.Sprintf("$P(X > u + y) = P(X>u)(1 + \\xi y/\\sigma)^{-1/\\xi} = %s(1 + %s)^{-%s} = %s$, which is %s ppm.", f(pu, 2), f(xi*y/sg, 4), f(1/xi, 3), sci(p), ans), map[string]string{"p": sci(p)})
	case 2:
		y := float64(edakit.Int(5, 40))
		p := pu * math.Exp(-y/sg)
		ppm := p * 1e6
		pl := sigPlaces(ppm)
		ans := edakit.Num(ppm, pl)
		return edakit.Numeric(edakit.Sprintf("Exceedances of a delay over $u = %s$ ps (probability %s) fit a generalized Pareto distribution with shape $\\xi = 0$, the exponential tail, and scale %s ps. What is the probability the delay exceeds %s ps, in parts per million? %s", f(u, 0), f(pu, 2), f(sg, 0), f(u+y, 0), roundTo(pl)),
			ans, edakit.Sprintf("With $\\xi = 0$: $P = P(X>u)\\,e^{-y/\\sigma} = %s e^{-%s} = %s$, which is %s ppm. Gaussian tails fall in this domain of attraction.", f(pu, 2), f(y/sg, 4), sci(p), ans), map[string]string{"p": sci(p)})
	default:
		xi := edakit.Pick([]float64{0.05, 0.1, 0.15, 0.2})
		pt := edakit.Pick([]float64{1e-5, 1e-6, 1e-7})
		x := u + sg/xi*(math.Pow(pu/pt, xi)-1)
		ans := edakit.Num(x, 1)
		return edakit.Numeric(edakit.Sprintf("A tail model says delays exceed $u = %s$ ps with probability %s, and the excesses follow a generalized Pareto distribution with $\\sigma = %s$ ps and $\\xi = %s$. What delay, in ps, is exceeded with probability %s? Round to 1 decimal place.", f(u, 0), f(pu, 2), f(sg, 0), f(xi, 2), sciT(pt)),
			ans, edakit.Sprintf("Invert the tail: $x_p = u + \\frac{\\sigma}{\\xi}[(P(X>u)/p)^{\\xi} - 1] = %s + %s[(%s)^{%s} - 1] = %s$ ps.", f(u, 0), f(sg/xi, 2), sci(pu/pt), f(xi, 2), ans), map[string]string{"ratio": sci(pu / pt)})
	}
}
