package stat

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// yieldModel: two-sided parametric yield (KP 1), a product of independent yields (KP 2), and
// good dies per wafer (KP 3).
func yieldModel(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		mu := float64(edakit.Int(90, 110))
		s := float64(edakit.Int(2, 8))
		lo := mu - s*float64(edakit.Int(10, 35))/10
		hi := mu + s*float64(edakit.Int(10, 35))/10
		lo, hi = math.Round(lo), math.Round(hi)
		y := cdf((hi-mu)/s) - cdf((lo-mu)/s)
		ans := edakit.Num(100*y, 2)
		return edakit.Numeric(edakit.Sprintf("A ring-oscillator frequency is Gaussian with mean %s MHz and standard deviation %s MHz. A die passes if the frequency lies between %s and %s MHz. What is the parametric yield, in percent? Round to 2 decimal places.", f(mu, 0), f(s, 0), f(lo, 0), f(hi, 0)),
			ans, edakit.Sprintf("$Y = \\Phi(%s) - \\Phi(%s) = %s - %s$, which is %s percent.", f((hi-mu)/s, 3), f((lo-mu)/s, 3), f(cdf((hi-mu)/s), 4), f(cdf((lo-mu)/s), 4), ans),
			map[string]string{"z_hi": f((hi-mu)/s, 3), "z_lo": f((lo-mu)/s, 3)})
	case 2:
		yd := float64(edakit.Int(70, 95)) / 100
		n := edakit.Scaled(ctx, 2, 4)
		y := yd
		var parts []string
		for i := 0; i < n; i++ {
			v := float64(edakit.Int(900, 999)) / 1000
			y *= v
			parts = append(parts, f(v, 3))
		}
		ans := edakit.Num(100*y, 2)
		return edakit.Numeric(edakit.Sprintf("A die is good if it has no killer defect (probability %s) and independently passes each of %d parametric specs (pass probabilities %s). What is the overall yield, in percent? Round to 2 decimal places.", f(yd, 2), n, strings.Join(parts, ", ")),
			ans, edakit.Sprintf("Independent requirements multiply: $%s \\times %s = %s$, which is %s percent.", f(yd, 2), strings.Join(parts, " \\times "), f(y, 5), ans), map[string]string{"yield": f(y, 5)})
	default:
		d := float64(edakit.Pick([]int{200, 300}))
		a := float64(edakit.Int(20, 400))
		y := float64(edakit.Int(40, 95)) / 100
		dpw := math.Pi*d*d/(4*a) - math.Pi*d/math.Sqrt(2*a)
		good := dpw * y
		ans := edakit.Num(good, 0)
		return edakit.Numeric(edakit.Sprintf("A %s mm wafer carries dies of %s mm$^2$. Using $DPW = \\pi d^2/(4A) - \\pi d/\\sqrt{2A}$ for gross dies per wafer and a die yield of %s percent, how many good dies does a wafer give on average? Round to the nearest integer.", f(d, 0), f(a, 0), f(100*y, 0)),
			ans, edakit.Sprintf("$DPW = %s - %s = %s$ gross dies; times %s gives %s good dies.", f(math.Pi*d*d/(4*a), 2), f(math.Pi*d/math.Sqrt(2*a), 2), f(dpw, 2), f(y, 2), ans), map[string]string{"dpw": f(dpw, 2)})
	}
}

// poissonYield: yield from area and defect density (KP 1), defect density from yield (KP 2),
// and yield summed over layers (KP 3).
func poissonYield(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		a := float64(edakit.Int(20, 600))
		d0 := float64(edakit.Int(5, 50)) / 100
		y := math.Exp(-a / 100 * d0)
		ans := edakit.Num(100*y, 2)
		return edakit.Numeric(edakit.Sprintf("A die of %s mm$^2$ is made in a process with %s killer defects per cm$^2$, spread at random. Under the Poisson model, what is the die yield, in percent? Round to 2 decimal places.", f(a, 0), f(d0, 2)),
			ans, edakit.Sprintf("$A = %s$ cm$^2$, so $AD_0 = %s$ and $Y = e^{-AD_0} = %s$ percent.", f(a/100, 2), f(a/100*d0, 4), ans), map[string]string{"ad0": f(a/100*d0, 4)})
	case 2:
		a := float64(edakit.Int(50, 600))
		y := float64(edakit.Int(30, 95)) / 100
		d0 := -math.Log(y) / (a / 100)
		ans := edakit.Num(d0, 4)
		return edakit.Numeric(edakit.Sprintf("Dies of %s mm$^2$ yield %s percent. Assuming Poisson-distributed random defects, what defect density does that imply, in defects per cm$^2$? Round to 4 decimal places.", f(a, 0), f(100*y, 0)),
			ans, edakit.Sprintf("$D_0 = -\\ln Y / A = -\\ln(%s)/%s = %s$ per cm$^2$.", f(y, 2), f(a/100, 2), ans), map[string]string{"neg_ln_y": f(-math.Log(y), 4)})
	default:
		n := edakit.Scaled(ctx, 2, 4)
		names := []string{"metal 1 shorts", "via opens", "poly shorts", "metal 2 opens"}
		lam := 0.0
		var parts []string
		for i := 0; i < n; i++ {
			ac := float64(edakit.Int(10, 120))
			di := float64(edakit.Int(5, 40)) / 100
			lam += ac / 100 * di
			parts = append(parts, edakit.Sprintf("%s: critical area %s mm$^2$, %s defects per cm$^2$", names[i], f(ac, 0), f(di, 2)))
		}
		ans := edakit.Num(100*math.Exp(-lam), 2)
		return edakit.Numeric(edakit.Sprintf("A layout's defect-limited yield is computed from critical areas per mechanism, each with Poisson defects. %s. What is the yield, in percent? Round to 2 decimal places.", strings.Join(parts, "; ")),
			ans, edakit.Sprintf("Expected faults add over mechanisms: $\\lambda = \\sum A_{c,i} D_i = %s$, so $Y = e^{-\\lambda} = %s$ percent.", f(lam, 4), ans), map[string]string{"lambda": f(lam, 4)})
	}
}

// clusteredYield: negative binomial yield (KP 1), its gap to Poisson (KP 2), and the yield of a
// larger die (KP 3).
func clusteredYield(ctx generator.GeneratorContext) generator.Problem {
	a := float64(edakit.Int(50, 600))
	d0 := float64(edakit.Int(5, 50)) / 100
	al := edakit.Pick([]float64{0.5, 1, 2, 3, 4, 5})
	lam := a / 100 * d0
	nb := math.Pow(1+lam/al, -al)
	stem := edakit.Sprintf("A die of %s mm$^2$ sees %s defects per cm$^2$, clustered with negative binomial parameter $\\alpha = %s$.", f(a, 0), f(d0, 2), f(al, 1))
	switch kp(ctx) {
	case 1:
		ans := edakit.Num(100*nb, 2)
		return edakit.Numeric(stem+" What is the die yield, in percent? Round to 2 decimal places.",
			ans, edakit.Sprintf("$AD_0 = %s$, so $Y = (1 + AD_0/\\alpha)^{-\\alpha} = %s$ percent.", f(lam, 4), ans), map[string]string{"ad0": f(lam, 4)})
	case 2:
		gap := 100 * (nb - math.Exp(-lam))
		ans := edakit.Num(gap, 2)
		return edakit.Numeric(stem+" By how many percentage points does this yield exceed the Poisson prediction for the same defect density? Round to 2 decimal places.",
			ans, edakit.Sprintf("Negative binomial: %s percent. Poisson: $e^{-%s} = %s$ percent. Clustering piles defects onto fewer dies, so more dies are clean: the gap is %s points.", f(100*nb, 3), f(lam, 4), f(100*math.Exp(-lam), 3), ans), map[string]string{"poisson": f(100*math.Exp(-lam), 3)})
	default:
		m := float64(edakit.Pick([]int{2, 3, 4}))
		y2 := math.Pow(1+m*lam/al, -al)
		ans := edakit.Num(100*y2, 2)
		return edakit.Numeric(stem+edakit.Sprintf(" A next-generation die is %s times larger in the same process. What is its yield, in percent? Round to 2 decimal places.", f(m, 0)),
			ans, edakit.Sprintf("$AD_0$ grows to %s, so $Y = (1 + %s/%s)^{-%s} = %s$ percent, against %s percent for the original die.", f(m*lam, 4), f(m*lam, 4), f(al, 1), f(al, 1), ans, f(100*nb, 2)), map[string]string{"ad0": f(m*lam, 4)})
	}
}

func poissonCDF(k int, lam float64) float64 {
	s, term := 0.0, math.Exp(-lam)
	for i := 0; i <= k; i++ {
		if i > 0 {
			term *= lam / float64(i)
		}
		s += term
	}
	return s
}

// redundancy: repairable yield with spares (KP 1), the fewest spares for a target (KP 2), and a
// k-of-n block yield (KP 3).
func redundancy(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		lam := float64(edakit.Int(5, 40)) / 10
		r := edakit.Int(1, 5)
		y := poissonCDF(r, lam)
		ans := edakit.Num(100*y, 2)
		return edakit.Numeric(edakit.Sprintf("An SRAM array expects $\\lambda = %s$ random faults (Poisson). It has %d spare rows, each able to repair one fault, and the spares themselves are fault-free. What is the repaired yield, in percent? Round to 2 decimal places.", f(lam, 1), r),
			ans, edakit.Sprintf("The array is good with at most %d faults: $\\sum_{k=0}^{%d} e^{-\\lambda}\\lambda^k/k! = %s$ percent. With no spares it would be %s percent.", r, r, ans, f(100*math.Exp(-lam), 2)), map[string]string{"y0": f(100*math.Exp(-lam), 2)})
	case 2:
		lam := float64(edakit.Int(5, 50)) / 10
		target := edakit.Pick([]float64{0.9, 0.95, 0.99, 0.999})
		r := 0
		for poissonCDF(r, lam) < target {
			r++
		}
		return edakit.Numeric(edakit.Sprintf("A memory expects $\\lambda = %s$ Poisson faults, and each spare repairs one fault. What is the fewest spares that give a repaired yield of at least %s percent?", f(lam, 1), f(100*target, 1)),
			edakit.IntStr(r), edakit.Sprintf("Accumulate $P(K \\le r)$: with %d spares it reaches %s percent, the first value at or above the target.", r, f(100*poissonCDF(r, lam), 3)), map[string]string{"spares": edakit.IntStr(r)})
	default:
		n := edakit.Pick([]int{4, 8, 16})
		s := edakit.Int(1, 3)
		p := float64(edakit.Int(1, 15)) / 100
		tot := n + s
		y := 0.0
		for k := 0; k <= s; k++ {
			y += float64(binom(tot, k)) * math.Pow(p, float64(k)) * math.Pow(1-p, float64(tot-k))
		}
		ans := edakit.Num(100*y, 2)
		return edakit.Numeric(edakit.Sprintf("A cache needs %d working sub-arrays and is built with %d, so %d are spares. Each sub-array is independently defective with probability %s. What is the yield, in percent? Round to 2 decimal places.", n, tot, s, f(p, 2)),
			ans, edakit.Sprintf("Good if at most %d of the %d are defective: $\\sum_{k=0}^{%d}\\binom{%d}{k}p^k(1-p)^{%d-k} = %s$ percent. Without spares it would be $(1-p)^{%d} = %s$ percent.", s, tot, s, tot, tot, ans, n, f(100*math.Pow(1-p, float64(n)), 2)), map[string]string{"no_spare": f(100*math.Pow(1-p, float64(n)), 2)})
	}
}

// designCentering: yield gained by centering (KP 1), sigma required for a target yield (KP 2),
// and the yield after a sensitivity-guided shift (KP 3).
func designCentering(ctx generator.GeneratorContext) generator.Problem {
	lo := float64(edakit.Int(80, 95))
	hi := lo + float64(edakit.Int(10, 30))
	c := (lo + hi) / 2
	s := float64(edakit.Int(2, 6))
	yf := func(m float64) float64 { return cdf((hi-m)/s) - cdf((lo-m)/s) }
	switch kp(ctx) {
	case 1:
		mu := c + float64(edakit.Int(2, 8))*edakit.Pick([]float64{-1, 1})
		gain := 100 * (yf(c) - yf(mu))
		ans := edakit.Num(gain, 2)
		return edakit.Numeric(edakit.Sprintf("A performance spec runs from %s to %s. The performance is Gaussian with standard deviation %s and its mean currently sits at %s. How many percentage points of yield does moving the mean to the centre of the spec gain? Round to 2 decimal places.", f(lo, 0), f(hi, 0), f(s, 0), f(mu, 0)),
			ans, edakit.Sprintf("At mean %s the yield is %s percent; at the centre %s it is %s percent. The gain is %s points.", f(mu, 0), f(100*yf(mu), 3), f(c, 1), f(100*yf(c), 3), ans), map[string]string{"y_now": f(100*yf(mu), 3)})
	case 2:
		target := edakit.Pick([]float64{0.95, 0.99, 0.997, 0.999})
		z := tailInv((1 - target) / 2)
		sig := (hi - lo) / (2 * z)
		ans := edakit.Num(sig, 3)
		return edakit.Numeric(edakit.Sprintf("A centred design has a spec window from %s to %s. What is the largest standard deviation that still gives a yield of %s percent? Round to 3 decimal places.", f(lo, 0), f(hi, 0), f(100*target, 1)),
			ans, edakit.Sprintf("Each tail may hold $(1-Y)/2$, so the half-width must be $z = Q^{-1}(%s) = %s$ sigmas: $\\sigma = (U-L)/(2z) = %s$.", f((1-target)/2, 4), f(z, 4), ans), map[string]string{"z": f(z, 4)})
	default:
		mu := c + float64(edakit.Int(3, 8))*edakit.Pick([]float64{-1, 1})
		a := float64(edakit.Int(2, 10)) / 2
		dx := float64(edakit.Int(1, 6)) / 2
		if mu > c {
			dx = -dx
		}
		nm := mu + a*dx
		y := yf(nm)
		ans := edakit.Num(100*y, 2)
		return edakit.Numeric(edakit.Sprintf("A performance spec runs from %s to %s, with Gaussian spread of standard deviation %s and current mean %s. The performance changes by %s units per unit of a sizing parameter. The designer moves the sizing by %s units. What is the new yield, in percent? Round to 2 decimal places.", f(lo, 0), f(hi, 0), f(s, 0), f(mu, 0), f(a, 1), f(dx, 1)),
			ans, edakit.Sprintf("The mean moves by $%s \\times %s = %s$ to %s. Yield is $\\Phi(%s) - \\Phi(%s) = %s$ percent, from %s percent before.", f(a, 1), f(dx, 1), f(a*dx, 2), f(nm, 2), f((hi-nm)/s, 3), f((lo-nm)/s, 3), ans, f(100*yf(mu), 2)), map[string]string{"new_mean": f(nm, 2)})
	}
}
