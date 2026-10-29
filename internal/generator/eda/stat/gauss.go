package stat

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var sigmaLevels = []float64{3, 3.5, 4, 4.5, 5, 5.5, 6}

// sigmaToProb: one-sided tail at k sigma in ppm (KP 1), expected failing cells in an array
// (KP 2), and the cell sigma an array-yield target requires (KP 3).
func sigmaToProb(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		k := edakit.Pick([]float64{2.5, 3, 3.25, 3.5, 3.75, 4, 4.25, 4.5, 5})
		p := tail(k)
		ppm := p * 1e6
		pl := sigPlaces(ppm)
		ans := edakit.Num(ppm, pl)
		return edakit.Numeric(edakit.Sprintf("A cell fails when a Gaussian margin falls more than %s standard deviations below its mean. What is the failure probability, in parts per million? %s", f(k, 2), roundTo(pl)),
			ans, edakit.Sprintf("$P = Q(%s) = 1 - \\Phi(%s) = %s$, which is %s ppm.", f(k, 2), f(k, 2), sci(p), ans), map[string]string{"q": sci(p)})
	case 2:
		k := edakit.Pick(sigmaLevels)
		mb := edakit.Pick([]int{1, 2, 4, 8, 16, 32, 64, 128})
		n := float64(mb) * 1048576
		e := n * tail(k)
		pl := sigPlaces(e)
		if pl > 4 {
			pl = 4
		}
		ans := edakit.Num(e, pl)
		return edakit.Numeric(edakit.Sprintf("An SRAM macro holds %d Mb ($%d \\times 2^{20}$ cells). Each cell independently fails if its read margin is beyond %s sigma on the failing side. How many failing cells does one macro contain on average? %s", mb, mb, f(k, 1), roundTo(pl)),
			ans, edakit.Sprintf("Per-cell $p = Q(%s) = %s$. Expected failures $= N p = %s \\times %s = %s$.", f(k, 1), sci(tail(k)), f(n, 0), sci(tail(k)), ans), map[string]string{"p_cell": sci(tail(k))})
	default:
		mb := edakit.Pick([]int{1, 4, 16, 64, 256})
		y := edakit.Pick([]float64{0.9, 0.95, 0.99, 0.999})
		n := float64(mb) * 1048576
		p := -math.Expm1(math.Log(y) / n)
		k := tailInv(p)
		ans := edakit.Num(k, 2)
		return edakit.Numeric(edakit.Sprintf("A %d Mb SRAM ($%d \\times 2^{20}$ cells, no redundancy) must work with probability %s. Cells fail independently with a one-sided Gaussian tail. At least how many sigma of margin must each cell have? Round to 2 decimal places.", mb, mb, f(y, 3)),
			ans, edakit.Sprintf("All $N = %s$ cells must pass: $(1-p)^N = %s$, so $p = 1 - %s^{1/N} = %s$. Then $k = Q^{-1}(p) = %s$.", f(n, 0), f(y, 3), f(y, 3), sci(p), ans), map[string]string{"p_cell": sci(p)})
	}
}

// linearCombo: sigma of a sum of two correlated delays (KP 1), sigma of n equally correlated
// stages (KP 2), and the 3-sigma-over-mean ratio for that path (KP 3).
func linearCombo(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		s1, s2 := float64(edakit.Int(2, 15)), float64(edakit.Int(2, 15))
		rho := float64(edakit.Int(-5, 9)) / 10
		v := s1*s1 + s2*s2 + 2*rho*s1*s2
		ans := edakit.Num(math.Sqrt(v), 3)
		return edakit.Numeric(edakit.Sprintf("Two stage delays on a path have standard deviations %s ps and %s ps and correlation coefficient %s. What is the standard deviation of their sum, in ps? Round to 3 decimal places.", f(s1, 0), f(s2, 0), f(rho, 1)),
			ans, edakit.Sprintf("$\\mathrm{Var} = \\sigma_1^2 + \\sigma_2^2 + 2\\rho\\sigma_1\\sigma_2 = %s + %s %s = %s$, so $\\sigma = %s$ ps.", f(s1*s1, 0), f(s2*s2, 0), signed2(2*rho*s1*s2), f(v, 2), ans), map[string]string{"var": f(v, 2)})
	case 2:
		n := edakit.Scaled(ctx, 3, 20)
		s := float64(edakit.Int(2, 10))
		rho := float64(edakit.Int(1, 9)) / 10
		nf := float64(n)
		v := nf*s*s + nf*(nf-1)*rho*s*s
		ans := edakit.Num(math.Sqrt(v), 3)
		return edakit.Numeric(edakit.Sprintf("A path has %d identical stages. Each stage delay has standard deviation %s ps and every pair of stages has correlation %s. What is the standard deviation of the path delay, in ps? Round to 3 decimal places.", n, f(s, 0), f(rho, 1)),
			ans, edakit.Sprintf("$\\mathrm{Var} = n\\sigma^2 + n(n-1)\\rho\\sigma^2 = %s + %s = %s$, so $\\sigma = %s$ ps.", f(nf*s*s, 0), f(nf*(nf-1)*rho*s*s, 2), f(v, 2), ans), map[string]string{"var": f(v, 2)})
	default:
		n := edakit.Scaled(ctx, 4, 30)
		mu := float64(edakit.Int(10, 40))
		s := float64(edakit.Int(1, 6))
		corr := edakit.Pick([]bool{true, false})
		nf := float64(n)
		sig := math.Sqrt(nf) * s
		word := "independent"
		if corr {
			sig = nf * s
			word = "perfectly correlated"
		}
		r := 100 * 3 * sig / (nf * mu)
		ans := edakit.Num(r, 2)
		return edakit.Numeric(edakit.Sprintf("A path has %d identical stages, each with mean delay %s ps and standard deviation %s ps. The stage delays are %s. What is the ratio of the path's 3-sigma spread to its mean delay, in percent? Round to 2 decimal places.", n, f(mu, 0), f(s, 0), word),
			ans, edakit.Sprintf("The mean is $%d \\times %s = %s$ ps. The sigma is %s ps (%s stages add %s). The ratio is $3(%s)/%s = %s$ percent.", n, f(mu, 0), f(nf*mu, 0), f(sig, 3), word, map[bool]string{true: "linearly in sigma", false: "in variance"}[corr], f(sig, 3), f(nf*mu, 0), ans),
			map[string]string{"path_sigma": f(sig, 3)})
	}
}

// clark returns the tightness probability, mean and variance of max(A, B) for jointly
// Gaussian A and B, following Clark (1961).
func clark(m1, s1, m2, s2, rho float64) (t, mean, variance, theta, alpha float64) {
	theta = math.Sqrt(s1*s1 + s2*s2 - 2*rho*s1*s2)
	alpha = (m1 - m2) / theta
	t = cdf(alpha)
	mean = m1*t + m2*(1-t) + theta*phi(alpha)
	second := (m1*m1+s1*s1)*t + (m2*m2+s2*s2)*(1-t) + (m1+m2)*theta*phi(alpha)
	variance = second - mean*mean
	return
}

// clarkMax: tightness probability (KP 1), mean of the max (KP 2), sigma of the max (KP 3).
func clarkMax(ctx generator.GeneratorContext) generator.Problem {
	m1 := float64(edakit.Int(90, 130))
	m2 := float64(edakit.Int(90, 130))
	s1, s2 := float64(edakit.Int(3, 12)), float64(edakit.Int(3, 12))
	rho := float64(edakit.Int(0, 6)) / 10
	t, mean, v, theta, alpha := clark(m1, s1, m2, s2, rho)
	stem := edakit.Sprintf("Two arrival times at a gate input are jointly Gaussian: $A \\sim N(%s, %s^2)$ ps and $B \\sim N(%s, %s^2)$ ps with correlation %s.", f(m1, 0), f(s1, 0), f(m2, 0), f(s2, 0), f(rho, 1))
	facts := map[string]string{"theta": f(theta, 4), "alpha": f(alpha, 4), "tightness": f(t, 4)}
	work := edakit.Sprintf("$\\theta = \\sqrt{\\sigma_A^2 + \\sigma_B^2 - 2\\rho\\sigma_A\\sigma_B} = %s$ and $\\alpha = (\\mu_A - \\mu_B)/\\theta = %s$.", f(theta, 4), f(alpha, 4))
	switch kp(ctx) {
	case 1:
		ans := edakit.Num(t, 4)
		return edakit.Numeric(stem+" What is the probability that $A$ is the larger (the tightness probability of $A$)? Round to 4 decimal places.",
			ans, work+edakit.Sprintf(" The tightness probability is $\\Phi(\\alpha) = %s$.", ans), facts)
	case 2:
		ans := edakit.Num(mean, 3)
		return edakit.Numeric(stem+" Using Clark's formulas, what is the mean of $\\max(A, B)$, in ps? Round to 3 decimal places.",
			ans, work+edakit.Sprintf(" $E[\\max] = \\mu_A\\Phi(\\alpha) + \\mu_B\\Phi(-\\alpha) + \\theta\\varphi(\\alpha) = %s$ ps, above both means.", ans), facts)
	default:
		ans := edakit.Num(math.Sqrt(v), 3)
		return edakit.Numeric(stem+" Using Clark's second moment, what is the standard deviation of $\\max(A, B)$, in ps? Round to 3 decimal places.",
			ans, work+edakit.Sprintf(" $E[\\max^2] = (\\mu_A^2+\\sigma_A^2)\\Phi(\\alpha) + (\\mu_B^2+\\sigma_B^2)\\Phi(-\\alpha) + (\\mu_A+\\mu_B)\\theta\\varphi(\\alpha)$. With mean %s the variance is %s, so $\\sigma = %s$ ps.", f(mean, 4), f(v, 4), ans), facts)
	}
}

// canonicalForm: sigma of a canonical delay (KP 1), correlation of two canonical delays (KP 2),
// and one sensitivity of their statistical max by tightness weighting (KP 3).
func canonicalForm(ctx generator.GeneratorContext) generator.Problem {
	a := []float64{float64(edakit.Int(1, 6)), nonzero(-3, 5), float64(edakit.Int(1, 4))}
	b := []float64{float64(edakit.Int(1, 6)), nonzero(-3, 5), float64(edakit.Int(1, 4))}
	a0, b0 := float64(edakit.Int(100, 140)), float64(edakit.Int(100, 140))
	form := func(m float64, c []float64) string {
		s := f(m, 0)
		for i, v := range c[:2] {
			sign := "+"
			if v < 0 {
				sign = "-"
			}
			s += edakit.Sprintf(" %s %s\\,\\Delta X_%d", sign, f(math.Abs(v), 0), i+1)
		}
		return s + edakit.Sprintf(" + %s\\,\\Delta R", f(c[2], 0))
	}
	stem := edakit.Sprintf("In first-order canonical form, with $\\Delta X_1, \\Delta X_2$ shared principal components and each $\\Delta R$ an independent random term (all standard normal), two arrival times are $A = %s$ and $B = %s$ ps.", form(a0, a), form(b0, b))
	va := a[0]*a[0] + a[1]*a[1] + a[2]*a[2]
	vb := b[0]*b[0] + b[1]*b[1] + b[2]*b[2]
	cov := a[0]*b[0] + a[1]*b[1]
	switch kp(ctx) {
	case 1:
		ans := edakit.Num(math.Sqrt(va), 3)
		return edakit.Numeric(stem+" What is the standard deviation of $A$, in ps? Round to 3 decimal places.",
			ans, edakit.Sprintf("$\\sigma_A^2 = %s^2 + %s^2 + %s^2 = %s$, so $\\sigma_A = %s$ ps.", pn(a[0], 0), pn(a[1], 0), pn(a[2], 0), f(va, 0), ans), map[string]string{"var_a": f(va, 0)})
	case 2:
		r := cov / math.Sqrt(va*vb)
		ans := edakit.Num(r, 3)
		return edakit.Numeric(stem+" What is the correlation coefficient of $A$ and $B$? Round to 3 decimal places.",
			ans, edakit.Sprintf("Only the shared components covary: $\\mathrm{Cov} = %s \\cdot %s + %s \\cdot %s = %s$. Dividing by $\\sigma_A\\sigma_B = \\sqrt{%s \\cdot %s}$ gives %s.", pn(a[0], 0), pn(b[0], 0), pn(a[1], 0), pn(b[1], 0), f(cov, 0), f(va, 0), f(vb, 0), ans), map[string]string{"cov": f(cov, 0)})
	default:
		rho := cov / math.Sqrt(va*vb)
		t, _, _, theta, alpha := clark(a0, math.Sqrt(va), b0, math.Sqrt(vb), rho)
		c1 := t*a[0] + (1-t)*b[0]
		ans := edakit.Num(c1, 3)
		return edakit.Numeric(stem+" A block-based statistical timer re-expresses $\\max(A,B)$ in canonical form by weighting sensitivities with the tightness probability. What is the coefficient of $\\Delta X_1$ in the result? Round to 3 decimal places.",
			ans, edakit.Sprintf("$\\theta = \\sqrt{\\sigma_A^2+\\sigma_B^2-2\\,\\mathrm{Cov}} = %s$, $\\alpha = (%s-%s)/\\theta = %s$, $T = \\Phi(\\alpha) = %s$. The coefficient is $T a_1 + (1-T) b_1 = %s$.", f(theta, 4), f(a0, 0), f(b0, 0), f(alpha, 4), f(t, 4), ans),
			map[string]string{"tightness": f(t, 4)})
	}
}

// timingYield: yield at a clock period (KP 1), period for a target yield (KP 2), and yield of
// two independent critical paths (KP 3).
func timingYield(ctx generator.GeneratorContext) generator.Problem {
	mu := float64(edakit.Int(400, 900))
	s := float64(edakit.Int(10, 40))
	switch kp(ctx) {
	case 1:
		t := mu + s*float64(edakit.Int(-10, 30))/10
		y := 100 * cdf((t-mu)/s)
		ans := edakit.Num(y, 2)
		return edakit.Numeric(edakit.Sprintf("Statistical timing gives the circuit delay as Gaussian with mean %s ps and standard deviation %s ps. What percentage of dies meet a clock period of %s ps? Round to 2 decimal places.", f(mu, 0), f(s, 0), f(t, 0)),
			ans, edakit.Sprintf("$z = (%s - %s)/%s = %s$ and $\\Phi(z) = %s$ percent.", f(t, 0), f(mu, 0), f(s, 0), f((t-mu)/s, 3), ans), map[string]string{"z": f((t-mu)/s, 3)})
	case 2:
		y := edakit.Pick([]float64{0.9, 0.95, 0.975, 0.99, 0.995, 0.999})
		z := tailInv(1 - y)
		t := mu + z*s
		ans := edakit.Num(t, 1)
		return edakit.Numeric(edakit.Sprintf("Circuit delay is Gaussian with mean %s ps and standard deviation %s ps. What is the shortest clock period, in ps, that gives a timing yield of %s percent? Round to 1 decimal place.", f(mu, 0), f(s, 0), f(100*y, 1)),
			ans, edakit.Sprintf("$z = \\Phi^{-1}(%s) = %s$, so $T = \\mu + z\\sigma = %s + %s(%s) = %s$ ps.", f(y, 3), f(z, 4), f(mu, 0), f(z, 4), f(s, 0), ans), map[string]string{"z": f(z, 4)})
	default:
		m2 := mu + float64(edakit.Int(-30, 30))
		s2 := float64(edakit.Int(10, 40))
		t := math.Max(mu, m2) + float64(edakit.Int(5, 25))*2
		y1, y2 := cdf((t-mu)/s), cdf((t-m2)/s2)
		ans := edakit.Num(100*y1*y2, 2)
		return edakit.Numeric(edakit.Sprintf("Two critical paths have independent Gaussian delays, $N(%s, %s^2)$ ps and $N(%s, %s^2)$ ps. The chip works only if both meet a clock period of %s ps. What is the timing yield, in percent? Round to 2 decimal places.", f(mu, 0), f(s, 0), f(m2, 0), f(s2, 0), f(t, 0)),
			ans, edakit.Sprintf("Path yields are $\\Phi(%s) = %s$ and $\\Phi(%s) = %s$. Independence lets them multiply: %s percent, below either path alone.", f((t-mu)/s, 3), f(y1, 4), f((t-m2)/s2, 3), f(y2, 4), ans), map[string]string{"y1": f(y1, 4), "y2": f(y2, 4)})
	}
}

// nonzero draws a nonzero integer in [lo, hi].
func nonzero(lo, hi int) float64 {
	v := edakit.Int(lo, hi)
	for v == 0 {
		v = edakit.Int(lo, hi)
	}
	return float64(v)
}

// signed2 renders a term with its sign and two decimals, "+ 3.5" or "- 3.5".
func signed2(v float64) string {
	if v < 0 {
		return "- " + f(-v, 2)
	}
	return "+ " + f(v, 2)
}
