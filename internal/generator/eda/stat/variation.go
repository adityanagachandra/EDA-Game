package stat

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// varianceSplitChoice draws a die-to-die and a within-die random sigma and asks one of three
// computed questions: the total sigma of one device, the correlation of two devices on the same
// die, or the sigma of their difference.
func varianceSplitChoice(ctx generator.GeneratorContext) generator.Problem {
	g := float64(edakit.Int(3, 30))
	l := float64(edakit.Int(3, 30))
	for l == g {
		l = float64(edakit.Int(3, 30))
	}
	param := edakit.Pick([]string{"threshold voltage", "gate delay"})
	unit := "mV"
	if param == "gate delay" {
		unit = "ps"
		g, l = g/5, l/5
	}
	stem := edakit.Sprintf("A transistor's %s has a die-to-die (global) component with $\\sigma_g = %s$ %s, shared by every device on a die, and an independent within-die random component with $\\sigma_l = %s$ %s per device.",
		param, f(g, 1), unit, f(l, 1), unit)
	tot := math.Hypot(g, l)
	switch edakit.Int(1, 3) {
	case 1:
		return numChoice(stem+" What is the total standard deviation of one device's value across many dies?", tot, 2, unit,
			edakit.Sprintf("Independent components add in variance: $\\sqrt{%s^2 + %s^2} = %s$ %s. Adding sigmas directly overstates the spread.", f(g, 1), f(l, 1), f(tot, 2), unit),
			g+l, math.Abs(g-l), (g+l)/2, tot/math.Sqrt2)
	case 2:
		rho := g * g / (g*g + l*l)
		return numChoice(stem+" What is the correlation coefficient between two devices on the same die?", rho, 3, "",
			edakit.Sprintf("Only the global part is shared, so the covariance is $\\sigma_g^2$ and the correlation is $\\sigma_g^2/(\\sigma_g^2+\\sigma_l^2) = %s$.", f(rho, 3)),
			g/(g+l), g/tot, l*l/(g*g+l*l))
	default:
		d := math.Sqrt2 * l
		return numChoice(stem+" What is the standard deviation of the difference between two devices on the same die (their mismatch)?", d, 2, unit,
			edakit.Sprintf("The shared global part cancels in the difference, leaving two independent random parts: $\\sqrt{2}\\,\\sigma_l = %s$ %s.", f(d, 2), unit),
			math.Sqrt2*tot, 2*l, l, math.Sqrt2*g)
	}
}

var processVariationItems = []edakit.Item{
	{KP: 1, Q: "Every die on one wafer runs 4% slower than the dies on the next wafer, but devices inside a die track each other closely. Which class of variation dominates?",
		A: "Global (die-to-die or wafer-to-wafer) variation", D: []string{"Local random mismatch", "Within-die systematic variation from layout density", "Temporal variation from aging"},
		Why: "A shift shared by every device on the die and differing between wafers is global variation. Mismatch would make neighbours differ from each other."},
	{KP: 1, Q: "Two identical, adjacent, minimum-size transistors differ in threshold voltage by 18 mV, and repeated pairs show both signs equally often. What is the source?",
		A: "Local random variation such as random dopant fluctuation and line-edge roughness", D: []string{"Global variation in gate-oxide thickness", "Systematic lithography proximity effects", "Supply voltage droop"},
		Why: "Neighbours with a zero-mean difference of random sign show independent local variation; a global shift would move both devices together."},
	{KP: 1, Q: "Gate lengths are reliably 2 nm shorter in dense arrays than in isolated lines, on every die. How should a model treat this?",
		A: "As systematic variation, predictable from layout context and modeled deterministically", D: []string{"As independent random variation per device", "As global variation drawn once per wafer", "As measurement noise to be averaged out"},
		Why: "A repeatable dependence on layout context is systematic. It can be predicted and corrected (for example by lithography-aware extraction) rather than treated as random."},
	{KP: 1, Q: "Why does a statistical timer split variation into global and local parts instead of using one total sigma per gate?",
		A: "The global part is shared by all gates and adds coherently along a path; local parts average out", D: []string{"The global part is always larger than the local part", "Local variation does not affect delay", "Splitting makes every gate delay independent"},
		Why: "Correlated components add linearly in sigma along a path while independent ones add in variance. One lumped sigma cannot get both behaviours right."},
	{KP: 1, Q: "A design is checked at a slow corner and a fast corner. Which kind of variation do those corners represent well, and which badly?",
		A: "They capture global variation well and local mismatch badly", D: []string{"They capture local mismatch well and global variation badly", "They capture both equally", "They capture only aging"},
		Why: "A corner moves every device the same way, which is what global variation does. Mismatch moves neighbours in opposite directions, which no single corner represents."},
	{KP: 1, Q: "Which variation is not a manufacturing variation at all, yet still enters a variability budget?",
		A: "Supply and temperature variation during operation", D: []string{"Random dopant fluctuation", "Wafer-level oxide thickness drift", "Line-edge roughness"},
		Why: "Voltage and temperature vary in the field (environmental variation). They are budgeted alongside process variation but have different statistics and correlation."},
	{KP: 3, Q: "A path has 20 gates, each with independent local delay sigma 2 ps and a fully shared global sigma 2 ps. Which component dominates the path sigma?",
		A: "The global component, 40 ps against about 8.9 ps", D: []string{"The local component, 40 ps against 2 ps", "They contribute equally, 20 ps each", "Neither: the sigmas cancel"},
		Why: "Shared parts add linearly: 20 times 2 = 40 ps. Independent parts add in variance: the square root of 20 times 4 is about 8.9 ps."},
	{KP: 3, Q: "A differential pair's offset is what limits an analog design. Which variation component sets that offset?",
		A: "The local mismatch between the two devices", D: []string{"The global die-to-die shift", "The systematic within-wafer gradient only", "The wafer-to-wafer variation"},
		Why: "A shared global shift moves both devices together and cancels in a differential measurement. The mismatch remains."},
	{KP: 3, Q: "A process engineer reports one total Vt sigma of 30 mV. Why is that not enough for an SRAM designer?",
		A: "Cell stability depends on mismatch between devices in one cell, which needs the local part separately", D: []string{"SRAM cells are immune to global variation", "SRAM designers only need the mean", "A 30 mV sigma is always negligible"},
		Why: "Read and write stability depend on the relative strengths of neighbouring devices, which is driven by local mismatch. The global part shifts the whole array."},
	{KP: 3, Q: "Two clock sinks are 2 mm apart. Why is the skew variation between them larger than between two sinks 20 um apart?",
		A: "Correlation of within-die variation falls with distance, so distant sinks share less of it", D: []string{"Global variation grows with distance", "Random mismatch grows linearly with distance", "Wire resistance is a random variable"},
		Why: "Spatially correlated variation is nearly shared by close neighbours and cancels in their skew. Distant sinks see more independent values."},
	{KP: 3, Q: "Doubling the number of identical independent stages on a path changes the relative random sigma (sigma over mean) how?",
		A: "It falls by a factor of the square root of 2", D: []string{"It doubles", "It halves", "It is unchanged"},
		Why: "The mean doubles and the independent sigma grows by the square root of 2, so their ratio drops by the square root of 2."},
	{KP: 3, Q: "Why must a variation-aware flow keep systematic variation out of the random sigma it feeds to Monte Carlo?",
		A: "Predictable shifts treated as random get averaged away or double-counted instead of being corrected", D: []string{"Monte Carlo cannot sample systematic variation", "Systematic variation has zero mean by definition", "It makes the simulation run slower"},
		Why: "A deterministic layout-dependent shift belongs in the nominal model. Folding it into a random sigma gives the wrong mean and the wrong correlation."},
}

// pelgrom: sigma of Vt mismatch from area (KP 1), area for a target sigma (KP 2), and the
// single-device sigma from the pair sigma (KP 3).
func pelgrom(ctx generator.GeneratorContext) generator.Problem {
	a := float64(edakit.Int(10, 50)) / 10 // mV um
	switch kp(ctx) {
	case 1:
		w := float64(edakit.Int(2, 20)) / 10
		l := float64(edakit.Pick([]int{3, 4, 5, 6, 8, 10, 15, 20})) / 100
		s := a / math.Sqrt(w*l)
		ans := edakit.Num(s, 2)
		return edakit.Numeric(edakit.Sprintf("A process has a Pelgrom coefficient $A_{V_T} = %s$ mV$\\cdot\\mu$m. A matched pair uses devices with $W = %s\\ \\mu$m and $L = %s\\ \\mu$m. What is the standard deviation of the threshold-voltage difference between the two devices, in mV? Round to 2 decimal places.", f(a, 1), f(w, 1), f(l, 2)),
			ans, edakit.Sprintf("$\\sigma(\\Delta V_T) = A_{V_T}/\\sqrt{WL} = %s/\\sqrt{%s} = %s$ mV.", f(a, 1), f(w*l, 4), ans),
			map[string]string{"area_um2": f(w*l, 4)})
	case 2:
		target := float64(edakit.Int(5, 30)) / 10
		area := (a / target) * (a / target)
		ans := edakit.Num(area, 3)
		return edakit.Numeric(edakit.Sprintf("With $A_{V_T} = %s$ mV$\\cdot\\mu$m, what gate area $WL$ in square micrometres gives a pair mismatch of $\\sigma(\\Delta V_T) = %s$ mV? Round to 3 decimal places.", f(a, 1), f(target, 1)),
			ans, edakit.Sprintf("Invert $\\sigma = A/\\sqrt{WL}$: $WL = (A/\\sigma)^2 = (%s/%s)^2 = %s\\ \\mu$m$^2$. Halving the sigma costs four times the area.", f(a, 1), f(target, 1), ans),
			map[string]string{"area_um2": ans})
	default:
		w := float64(edakit.Int(2, 20)) / 10
		l := float64(edakit.Pick([]int{3, 4, 5, 6, 8, 10, 15, 20})) / 100
		pair := a / math.Sqrt(w*l)
		single := pair / math.Sqrt2
		ans := edakit.Num(single, 2)
		return edakit.Numeric(edakit.Sprintf("Pelgrom's law with $A_{V_T} = %s$ mV$\\cdot\\mu$m gives the sigma of the difference between two devices. For one device of $W = %s\\ \\mu$m and $L = %s\\ \\mu$m, what local Vt sigma should a Monte Carlo model assign to each device independently, in mV? Round to 2 decimal places.", f(a, 1), f(w, 1), f(l, 2)),
			ans, edakit.Sprintf("The pair sigma is $%s/\\sqrt{%s} = %s$ mV. The difference of two independent devices has variance $2\\sigma^2$, so each device gets $%s/\\sqrt{2} = %s$ mV.", f(a, 1), f(w*l, 4), f(pair, 3), f(pair, 3), ans),
			map[string]string{"pair_sigma": f(pair, 3)})
	}
}

// spatialCorrelation: exponential correlation at a distance (KP 1), total correlation of two
// gate delays (KP 2), and the sigma of their difference (KP 3).
func spatialCorrelation(ctx generator.GeneratorContext) generator.Problem {
	lam := float64(edakit.Pick([]int{200, 300, 400, 500, 800, 1000}))
	d := float64(edakit.Int(1, 20) * 50)
	rho := math.Exp(-d / lam)
	switch kp(ctx) {
	case 1:
		ans := edakit.Num(rho, 3)
		return edakit.Numeric(edakit.Sprintf("Within-die variation of gate length is modeled with the exponential correlation $\\rho(d) = e^{-d/\\lambda}$ and correlation length $\\lambda = %s\\ \\mu$m. What is the correlation between two gates $%s\\ \\mu$m apart? Round to 3 decimal places.", f(lam, 0), f(d, 0)),
			ans, edakit.Sprintf("$\\rho = e^{-%s/%s} = %s$.", f(d, 0), f(lam, 0), ans), map[string]string{"rho_spatial": ans})
	case 2:
		sg, ss, sr := float64(edakit.Int(1, 6)), float64(edakit.Int(1, 6)), float64(edakit.Int(1, 6))
		c := (sg*sg + ss*ss*rho) / (sg*sg + ss*ss + sr*sr)
		ans := edakit.Num(c, 3)
		return edakit.Numeric(edakit.Sprintf("A gate delay varies as $D = \\mu + G + S + R$: a global part with $\\sigma_G = %s$ ps shared by the whole die, a spatial part with $\\sigma_S = %s$ ps whose correlation between two gates is $e^{-d/\\lambda}$ with $\\lambda = %s\\ \\mu$m, and an independent random part with $\\sigma_R = %s$ ps. Two gates sit $%s\\ \\mu$m apart. What is the correlation coefficient of their delays? Round to 3 decimal places.", f(sg, 0), f(ss, 0), f(lam, 0), f(sr, 0), f(d, 0)),
			ans, edakit.Sprintf("Spatial correlation is $e^{-%s/%s} = %s$. Covariance $= \\sigma_G^2 + \\sigma_S^2\\rho = %s$; variance $= %s$; correlation $= %s$.", f(d, 0), f(lam, 0), f(rho, 4), f(sg*sg+ss*ss*rho, 4), f(sg*sg+ss*ss+sr*sr, 0), ans),
			map[string]string{"rho_spatial": f(rho, 4)})
	default:
		sg, ss, sr := float64(edakit.Int(1, 6)), float64(edakit.Int(2, 8)), float64(edakit.Int(1, 4))
		v := 2 * (ss*ss*(1-rho) + sr*sr)
		ans := edakit.Num(math.Sqrt(v), 2)
		return edakit.Numeric(edakit.Sprintf("Two clock-buffer delays each vary as $\\mu + G + S + R$ with $\\sigma_G = %s$ ps (shared die-wide), $\\sigma_S = %s$ ps (spatial, correlation $e^{-d/\\lambda}$ with $\\lambda = %s\\ \\mu$m) and $\\sigma_R = %s$ ps (independent). The buffers are $%s\\ \\mu$m apart. What is the standard deviation of the skew, the difference of the two delays, in ps? Round to 2 decimal places.", f(sg, 0), f(ss, 0), f(lam, 0), f(sr, 0), f(d, 0)),
			ans, edakit.Sprintf("The global part cancels. $\\mathrm{Var} = 2\\sigma_S^2(1-\\rho) + 2\\sigma_R^2 = 2(%s)(1-%s) + 2(%s) = %s$, so the skew sigma is %s ps.", f(ss*ss, 0), f(rho, 4), f(sr*sr, 0), f(v, 3), ans),
			map[string]string{"rho_spatial": f(rho, 4)})
	}
}

// pcaVariation: largest eigenvalue of a 2x2 covariance (KP 1), share explained by the first
// component (KP 2), and the first-component share of an equicorrelated grid (KP 3).
func pcaVariation(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1, 2:
		a := float64(edakit.Int(2, 20))
		b := float64(edakit.Int(2, 20))
		c := float64(edakit.Int(1, int(math.Floor(math.Sqrt(a*b)))))
		if c*c >= a*b {
			c = math.Floor(math.Sqrt(a*b)) - 1
		}
		if c < 1 {
			c = 1
		}
		m := (a + b) / 2
		r := math.Sqrt((a-b)*(a-b)/4 + c*c)
		l1, l2 := m+r, m-r
		stem := edakit.Sprintf("The gate-length variations of two regions of a die have covariance matrix $\\begin{pmatrix} %s & %s \\\\ %s & %s \\end{pmatrix}$ in nm$^2$.", f(a, 0), f(c, 0), f(c, 0), f(b, 0))
		if ctx.KPOrdinal == 1 || (ctx.KPOrdinal == 0 && edakit.Int(0, 1) == 0) {
			ans := edakit.Num(l1, 3)
			return edakit.Numeric(stem+" What is the variance carried by the first principal component, in nm$^2$? Round to 3 decimal places.",
				ans, edakit.Sprintf("Eigenvalues of a symmetric 2x2: $\\frac{a+b}{2} \\pm \\sqrt{(\\frac{a-b}{2})^2 + c^2} = %s \\pm %s$. The larger is %s.", f(m, 1), f(r, 4), ans),
				map[string]string{"lambda1": ans, "lambda2": edakit.Num(l2, 3)})
		}
		pct := 100 * l1 / (a + b)
		ans := edakit.Num(pct, 1)
		return edakit.Numeric(stem+" What percentage of the total variance does the first principal component account for? Round to 1 decimal place.",
			ans, edakit.Sprintf("$\\lambda_1 = %s$ and the trace is %s, so the share is $%s/%s = %s$ percent.", f(l1, 4), f(a+b, 0), f(l1, 4), f(a+b, 0), ans),
			map[string]string{"lambda1": edakit.Num(l1, 4)})
	default:
		n := edakit.Scaled(ctx, 4, 64)
		rho := float64(edakit.Int(1, 9)) / 10
		share := 100 * (1 + float64(n-1)*rho) / float64(n)
		ans := edakit.Num(share, 1)
		return edakit.Numeric(edakit.Sprintf("A die is divided into %d grid cells. The variation in each cell has unit variance and every pair of cells has correlation %s. What percentage of the total variance does the first principal component account for? Round to 1 decimal place.", n, f(rho, 1)),
			ans, edakit.Sprintf("An equicorrelated matrix has one eigenvalue $1 + (n-1)\\rho = %s$ (the all-ones direction, a die-wide shift) and $n-1$ eigenvalues $1-\\rho$. The trace is %d, so the share is %s percent.", f(1+float64(n-1)*rho, 2), n, ans),
			map[string]string{"lambda1": f(1+float64(n-1)*rho, 2)})
	}
}

// spc: X-bar control limits (KP 1), process capability Cpk (KP 2), and the average run length
// after a mean shift (KP 3).
func spc(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		mu := float64(edakit.Int(200, 400)) / 10
		s := float64(edakit.Int(5, 30)) / 10
		n := edakit.Pick([]int{4, 5, 9, 16, 25})
		ucl := mu + 3*s/math.Sqrt(float64(n))
		ans := edakit.Num(ucl, 3)
		return edakit.Numeric(edakit.Sprintf("A fab monitors gate critical dimension with an $\\bar{X}$ chart. The in-control process has mean %s nm and standard deviation %s nm, and each plotted point averages %d sites. What is the upper control limit, in nm? Round to 3 decimal places.", f(mu, 1), f(s, 1), n),
			ans, edakit.Sprintf("$UCL = \\mu + 3\\sigma/\\sqrt{n} = %s + 3(%s)/\\sqrt{%d} = %s$ nm.", f(mu, 1), f(s, 1), n, ans), map[string]string{"se": f(s/math.Sqrt(float64(n)), 4)})
	case 2:
		lsl := float64(edakit.Int(180, 220)) / 10
		usl := lsl + float64(edakit.Int(40, 120))/10
		s := float64(edakit.Int(4, 15)) / 10
		mu := lsl + (usl-lsl)*float64(edakit.Int(30, 70))/100
		cpk := math.Min(usl-mu, mu-lsl) / (3 * s)
		ans := edakit.Num(cpk, 3)
		return edakit.Numeric(edakit.Sprintf("A gate length spec runs from %s nm to %s nm. The process mean is %s nm and its standard deviation is %s nm. What is the process capability index $C_{pk}$? Round to 3 decimal places.", f(lsl, 1), f(usl, 1), f(mu, 2), f(s, 1)),
			ans, edakit.Sprintf("$C_{pk} = \\min(USL-\\mu,\\ \\mu-LSL)/(3\\sigma) = %s/%s = %s$.", f(math.Min(usl-mu, mu-lsl), 2), f(3*s, 1), ans), map[string]string{"cpk": ans})
	default:
		delta := float64(edakit.Int(5, 20)) / 10
		n := edakit.Pick([]int{1, 2, 4, 5, 9})
		z := delta * math.Sqrt(float64(n))
		p := tail(3-z) + cdf(-3-z)
		arl := 1 / p
		ans := edakit.Num(arl, 1)
		return edakit.Numeric(edakit.Sprintf("An $\\bar{X}$ chart with 3-sigma limits plots the mean of %d samples. The process mean shifts up by %s process standard deviations. On average, how many plotted points pass before the chart signals (the average run length)? Round to 1 decimal place.", n, f(delta, 1)),
			ans, edakit.Sprintf("The plotted mean shifts by $\\delta\\sqrt{n} = %s$ of its own sigma. $P(\\text{signal}) = Q(3-%s) + \\Phi(-3-%s) = %s$, and $ARL = 1/P = %s$. With no shift the ARL is about 370.", f(z, 3), f(z, 3), f(z, 3), f(p, 5), ans),
			map[string]string{"p_signal": f(p, 6)})
	}
}
