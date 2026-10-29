package fab

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// dpw is the classic dies-per-wafer estimate with an edge-loss term, floored.
func dpw(d, area float64) int {
	n := math.Pi*d*d/(4*area) - math.Pi*d/math.Sqrt(2*area)
	if n < 0 {
		return 0
	}
	return int(math.Floor(n))
}

// diesPerWafer: area ratio, the edge-loss formula, and edge exclusion with scribe lanes.
func diesPerWafer(ctx generator.GeneratorContext) generator.Problem {
	d := edakit.Pick([]float64{200, 300})
	w := step(3, 25, 0.5)
	h := step(3, 25, 0.5)
	A := w * h
	switch kpOf(ctx) {
	case 1:
		r := math.Pi * d * d / 4 / A
		ans := edakit.Num(r, 1)
		q := edakit.Sprintf("A die measures %s mm by %s mm. Dividing the area of a %.0f mm wafer by the die area, how many dies would fit if edges did not matter? Round to 1 decimal place.", edakit.Num(w, 1), edakit.Num(h, 1), d)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Wafer area is $\\pi (%.0f)^2/4 = %s$ mm^2; die area is %s mm^2; the ratio is %s. Real counts are lower because square dies do not tile a circle.", d, edakit.Num(math.Pi*d*d/4, 0), edakit.Num(A, 2), ans),
			map[string]string{"die_area": edakit.Num(A, 2)})
	case 2:
		n := dpw(d, A)
		ans := edakit.IntStr(n)
		q := edakit.Sprintf("Estimate dies per wafer with $DPW = \\pi d^2/(4A) - \\pi d/\\sqrt{2A}$ for a %.0f mm wafer and a %s mm by %s mm die. Round down to a whole number.", d, edakit.Num(w, 1), edakit.Num(h, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$A = %s$ mm^2. The area term is %s and the edge-loss term is %s, so DPW = %s.", edakit.Num(A, 2), edakit.Num(math.Pi*d*d/(4*A), 2), edakit.Num(math.Pi*d/math.Sqrt(2*A), 2), ans),
			map[string]string{"area_term": edakit.Num(math.Pi*d*d/(4*A), 2)})
	default:
		e := float64(edakit.Int(2, 5))
		s := step(0.05, 0.2, 0.05)
		de := d - 2*e
		Ae := (w + s) * (h + s)
		n := dpw(de, Ae)
		ans := edakit.IntStr(n)
		q := edakit.Sprintf("A %.0f mm wafer has a %.0f mm edge exclusion and %s mm scribe lanes between dies. The die is %s mm by %s mm. Using $DPW = \\pi d^2/(4A) - \\pi d/\\sqrt{2A}$ with the usable diameter and the die pitch area, how many dies fit? Round down.", d, e, edakit.Num(s, 2), edakit.Num(w, 1), edakit.Num(h, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Usable diameter is %.0f mm and each die occupies $(%s)(%s) = %s$ mm^2 with its scribe. The formula gives %s dies.", de, edakit.Num(w+s, 2), edakit.Num(h+s, 2), edakit.Num(Ae, 3), ans),
			map[string]string{"pitch_area": edakit.Num(Ae, 3)})
	}
}

// defectDensity: Poisson yield.
func defectDensity(ctx generator.GeneratorContext) generator.Problem {
	amm := float64(10 * edakit.Int(5, 80))
	A := amm / 100
	D := step(0.05, 0.5, 0.01)
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(A*D, 3)
		q := edakit.Sprintf("Killer defects land at random with density %s per cm^2. What is the expected number of killer defects on a %.0f mm^2 die? Round to 3 decimal places.", edakit.Num(D, 2), amm)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Convert: %.0f mm^2 = %s cm^2. The mean is $\\lambda = A D_0 = %s$.", amm, edakit.Num(A, 2), ans),
			map[string]string{"A_cm2": edakit.Num(A, 2)})
	case 2:
		y := math.Exp(-A*D) * 100
		ans := edakit.Num(y, 2)
		q := edakit.Sprintf("Using the Poisson model with defect density %s per cm^2, what is the yield of a %.0f mm^2 die, in percent? Round to 2 decimal places.", edakit.Num(D, 2), amm)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("A die works if it gets zero defects: $Y = e^{-AD_0} = e^{-%s} = %s$%%.", edakit.Num(A*D, 4), ans),
			map[string]string{"lambda": edakit.Num(A*D, 4)})
	default:
		y := step(30, 95, 1)
		d0 := -math.Log(y/100) / A
		ans := edakit.Num(d0, 3)
		q := edakit.Sprintf("A product with a %.0f mm^2 die yields %.0f%% at wafer sort. Assuming all loss is random defects following the Poisson model, what defect density does that imply, in defects per cm^2? Round to 3 decimal places.", amm, y)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Invert $Y = e^{-AD_0}$: $D_0 = -\\ln(%s)/%s = %s$ per cm^2.", edakit.Num(y/100, 2), edakit.Num(A, 2), ans),
			map[string]string{"A_cm2": edakit.Num(A, 2)})
	}
}

// negBinomial: clustered-defect yield.
func negBinomial(ctx generator.GeneratorContext) generator.Problem {
	amm := float64(10 * edakit.Int(5, 80))
	A := amm / 100
	D := step(0.05, 0.6, 0.01)
	alpha := edakit.Pick([]float64{0.5, 1, 1.5, 2, 3, 4, 5})
	ynb := math.Pow(1+A*D/alpha, -alpha) * 100
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(ynb, 2)
		q := edakit.Sprintf("Using the negative binomial yield model $Y = (1 + A D_0/\\alpha)^{-\\alpha}$ with clustering parameter $\\alpha = %s$, what is the yield of a %.0f mm^2 die at %s defects per cm^2, in percent? Round to 2 decimal places.", edakit.Num(alpha, 1), amm, edakit.Num(D, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$A D_0 = %s$, so $Y = (1 + %s/%s)^{-%s} = %s$%%.", edakit.Num(A*D, 4), edakit.Num(A*D, 4), edakit.Num(alpha, 1), edakit.Num(alpha, 1), ans),
			map[string]string{"AD": edakit.Num(A*D, 4)})
	case 2:
		yp := math.Exp(-A*D) * 100
		ans := edakit.Num(ynb-yp, 2)
		q := edakit.Sprintf("A %.0f mm^2 die sees %s defects per cm^2. How many percentage points higher is the negative binomial yield with $\\alpha = %s$ than the Poisson yield? Round to 2 decimal places.", amm, edakit.Num(D, 2), edakit.Num(alpha, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Negative binomial gives %s%% and Poisson %s%%; the difference is %s points. Clustering puts several defects on the same die, leaving more dies clean.", edakit.Num(ynb, 2), edakit.Num(yp, 2), ans),
			map[string]string{"y_poisson": edakit.Num(yp, 2)})
	default:
		y := step(40, 95, 1)
		d0 := alpha * (math.Pow(y/100, -1/alpha) - 1) / A
		ans := edakit.Num(d0, 3)
		q := edakit.Sprintf("A %.0f mm^2 die yields %.0f%%. Assuming the negative binomial model with $\\alpha = %s$, what defect density does that imply, in defects per cm^2? Round to 3 decimal places.", amm, y, edakit.Num(alpha, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Invert: $D_0 = \\alpha (Y^{-1/\\alpha} - 1)/A = %s$ per cm^2, compared with %s from the Poisson model.", ans, edakit.Num(-math.Log(y/100)/A, 3)),
			map[string]string{"d_poisson": edakit.Num(-math.Log(y/100)/A, 3)})
	}
}

// criticalArea: short and open critical area for parallel wires.
func criticalArea(ctx generator.GeneratorContext) generator.Problem {
	L := float64(50 * edakit.Int(2, 40))
	s := float64(edakit.Int(16, 60))
	w := float64(edakit.Int(16, 60))
	x := float64(edakit.Int(10, 120))
	short := math.Max(0, x-s) * L / 1000
	open := math.Max(0, x-w) * L / 1000
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(short, 3)
		q := edakit.Sprintf("Two parallel wires run side by side for %.0f um with a %.0f nm space between them. A circular conducting defect of diameter %.0f nm shorts them if it touches both. Ignoring end effects, what is the critical area for shorts, in um^2? Round to 3 decimal places.", L, s, x)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The defect centre must lie in a band of width $\\max(0, x - s) = %.0f$ nm along the gap, so $A_c = L \\max(0, x - s) = %s$ um^2.", math.Max(0, x-s), ans),
			map[string]string{"band_nm": edakit.Num(math.Max(0, x-s), 0)})
	case 2:
		ans := edakit.Num(open, 3)
		q := edakit.Sprintf("A wire %.0f nm wide runs for %.0f um. A circular missing-material defect of diameter %.0f nm opens it if it spans the full width. Ignoring end effects, what is the critical area for opens, in um^2? Round to 3 decimal places.", w, L, x)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The centre must sit in a band of width $\\max(0, x - w) = %.0f$ nm, so $A_c = %s$ um^2. A defect smaller than the width cannot open the wire.", math.Max(0, x-w), ans),
			map[string]string{"band_nm": edakit.Num(math.Max(0, x-w), 0)})
	default:
		n := edakit.Scaled(ctx, 3, 12)
		tot := float64(n-1)*short + float64(n)*open
		ans := edakit.Num(tot, 3)
		q := edakit.Sprintf("A bus of %d parallel wires, each %.0f nm wide with %.0f nm spaces, runs for %.0f um. For defects of diameter %.0f nm, what is the total critical area for shorts plus opens, in um^2? Ignore end effects and defects that bridge three wires. Round to 3 decimal places.", n, w, s, L, x)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("There are %d gaps, each with short critical area %s um^2, and %d wires, each with open critical area %s um^2. Total %s um^2.", n-1, edakit.Num(short, 3), n, edakit.Num(open, 3), ans),
			map[string]string{"short_each": edakit.Num(short, 3), "open_each": edakit.Num(open, 3)})
	}
}

// avgCriticalArea: averaging the short critical area over a 1/x^3 size distribution.
func avgCriticalArea(ctx generator.GeneratorContext) generator.Problem {
	x0 := float64(edakit.Int(10, 30))
	switch kpOf(ctx) {
	case 1:
		x := x0 * step(1.5, 6, 0.5)
		f := x0 * x0 / (x * x) * 100
		ans := edakit.Num(f, 2)
		q := edakit.Sprintf("Defect diameters follow $f(x) = 2 x_0^2 / x^3$ for $x \\ge x_0$ with $x_0 = %.0f$ nm. What percentage of defects are larger than %s nm? Round to 2 decimal places.", x0, edakit.Num(x, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$P(X > x) = \\int_x^\\infty 2x_0^2/t^3\\,dt = (x_0/x)^2 = %s$%%.", ans),
			map[string]string{"p": ans})
	case 2:
		L := float64(10 * edakit.Int(5, 200))
		s := x0 + float64(edakit.Int(0, 40))
		a := L * x0 * x0 / s / 1000
		ans := edakit.Num(a, 3)
		q := edakit.Sprintf("Two parallel wires run %.0f um with a %.0f nm space. Defect diameters follow $f(x) = 2x_0^2/x^3$ for $x \\ge x_0 = %.0f$ nm, and the short critical area at size $x$ is $L(x - s)$ for $x > s$. What is the average short critical area, in um^2? Round to 3 decimal places.", L, s, x0)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$\\bar A = \\int_s^\\infty L(x-s)\\,2x_0^2/x^3\\,dx = L x_0^2/s = %.0f \\times %.0f^2 / %.0f$ nm um $= %s$ um^2.", L, x0, s, ans),
			map[string]string{"avg": ans})
	default:
		M := float64(10 * edakit.Int(5, 60)) // metres of wire
		s := x0 + float64(edakit.Int(0, 30))
		D0 := step(0.1, 1.0, 0.05)
		aCm2 := M * 1e6 * x0 * x0 / s / 1000 * 1e-8
		y := math.Exp(-D0*aCm2) * 100
		ans := edakit.Num(y, 2)
		q := edakit.Sprintf("A die has %.0f m of minimum-spaced parallel wire on one layer, with %.0f nm spaces. Defects follow $f(x) = 2x_0^2/x^3$ with $x_0 = %.0f$ nm and total density %s per cm^2. Using average critical area $L x_0^2/s$ and Poisson statistics, what is the yield limited by shorts on this layer, in percent? Round to 2 decimal places.", M, s, x0, edakit.Num(D0, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$\\bar A = %.0f$ m $\\times %.0f^2/%.0f$ nm $= %s$ cm^2, so $\\lambda = %s$ and $Y = e^{-\\lambda} = %s$%%.", M, x0, s, edakit.Num(aCm2, 5), edakit.Num(D0*aCm2, 5), ans),
			map[string]string{"avg_cm2": edakit.Num(aCm2, 5)})
	}
}

// costPerDie: wafer cost over good dies.
func costPerDie(ctx generator.GeneratorContext) generator.Problem {
	cw := float64(1000 * edakit.Int(3, 20))
	w := step(4, 20, 0.5)
	h := step(4, 20, 0.5)
	A := w * h
	n := dpw(300, A)
	D := step(0.05, 0.4, 0.01)
	y := math.Exp(-A / 100 * D)
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(cw/float64(n), 2)
		q := edakit.Sprintf("A processed 300 mm wafer costs %.0f US dollars and holds %d candidate dies. What is the cost per candidate die, in dollars, before yield? Round to 2 decimal places.", cw, n)
		return edakit.Numeric(q, ans, edakit.Sprintf("%.0f / %d = %s dollars.", cw, n, ans), map[string]string{"dpw": edakit.IntStr(n)})
	case 2:
		c := cw / (float64(n) * y)
		ans := edakit.Num(c, 2)
		q := edakit.Sprintf("A 300 mm wafer costs %.0f US dollars. The die is %s mm by %s mm, giving %d dies per wafer, and defect density is %s per cm^2 (Poisson). What is the cost per good die, in dollars? Round to 2 decimal places.", cw, edakit.Num(w, 1), edakit.Num(h, 1), n, edakit.Num(D, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Yield is $e^{-%s \\times %s} = %s$, so good dies are %s and cost is %s dollars.", edakit.Num(A/100, 4), edakit.Num(D, 2), edakit.Num(y, 4), edakit.Num(float64(n)*y, 1), ans),
			map[string]string{"yield": edakit.Num(y, 4)})
	default:
		f := step(0.5, 0.85, 0.05)
		A2 := A * f
		n2 := dpw(300, A2)
		y2 := math.Exp(-A2 / 100 * D)
		c1 := cw / (float64(n) * y)
		c2 := cw * step(1.0, 1.4, 0.05) / (float64(n2) * y2)
		ratio := c2 / c1
		mult := c2 * float64(n2) * y2 / cw
		ans := edakit.Num(ratio, 3)
		q := edakit.Sprintf("A %s mm^2 die on a 300 mm wafer costs %s dollars per good die (%d dies per wafer, defect density %s per cm^2, Poisson, wafer cost %.0f dollars). A shrink to %s of the area on a newer node raises wafer cost by a factor of %s, with the same defect density, and gives %d dies per wafer. What is the new cost per good die divided by the old? Round to 3 decimal places.",
			edakit.Num(A, 2), edakit.Num(c1, 2), n, edakit.Num(D, 2), cw, edakit.Num(f, 2), edakit.Num(mult, 2), n2)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("New yield $e^{-%s \\times %s} = %s$ on %d dies; new cost %s dollars; ratio %s. A shrink pays only if density and yield gains beat the wafer-cost increase.", edakit.Num(A2/100, 4), edakit.Num(D, 2), edakit.Num(y2, 4), n2, edakit.Num(c2, 2), ans),
			map[string]string{"c_new": edakit.Num(c2, 2)})
	}
}

// dfmCompute: yield gain from redundant vias.
func dfmCompute(ctx generator.GeneratorContext) generator.Problem {
	me := edakit.Pick([]int{6, 7})
	N := float64(edakit.Int(1, 9)) * math.Pow(10, float64(me))
	pm := edakit.Int(1, 9)
	p := float64(pm) * 1e-8
	ys := math.Exp(-N*p) * 100
	yd := math.Exp(-N*p*p) * 100
	if edakit.Int(0, 1) == 0 {
		q := edakit.Sprintf("A design has $%s \\times 10^{%d}$ single vias, each failing open independently with probability $%d \\times 10^{-8}$. What is the via-limited yield?", edakit.Num(N/math.Pow(10, float64(me)), 0), me, pm)
		c, d := numOptions(ys, 2, "%", (1-N*p)*100, (1-p)*100, yd)
		if (1-N*p)*100 < 0 {
			c, d = numOptions(ys, 2, "%", (1-p)*100, yd, ys/2)
		}
		return edakit.Choice(q, c,
			edakit.Sprintf("$Y = (1-p)^N \\approx e^{-Np} = e^{-%s}$, which is %s. Linearizing as $1 - Np$ overstates loss when $Np$ is not small.", edakit.Num(N*p, 3), c), d...)
	}
	q := edakit.Sprintf("A design has $%s \\times 10^{%d}$ via locations, each via failing with probability $%d \\times 10^{-8}$. If every location gets a redundant second via and a location fails only when both fail, what is the via-limited yield?", edakit.Num(N/math.Pow(10, float64(me)), 0), me, pm)
	c, d := numOptions(yd, 2, "%", ys, math.Exp(-N*p/2)*100, math.Exp(-2*N*p)*100)
	return edakit.Choice(q, c,
		edakit.Sprintf("A doubled location fails with $p^2$, so $Y \\approx e^{-N p^2}$, which rounds to %s, against %s with single vias.", c, edakit.Num(ys, 2)+"%"), d...)
}

var dfmItems = []edakit.Item{
	{KP: 1, Q: "Which change is a typical DFM fix applied after routing?",
		A:   "Doubling single vias wherever there is room",
		D:   []string{"Lowering the clock frequency", "Removing all dummy fill", "Shrinking every wire to minimum width"},
		Why: "Redundant vias cut the probability that one bad via opens a net, with little area cost."},
	{KP: 1, Q: "Wire spreading moves wires apart where space allows. What does it reduce?",
		A:   "Critical area for shorts between neighbours",
		D:   []string{"Via resistance", "Gate leakage", "The number of metal layers"},
		Why: "Short critical area grows as spacing shrinks; extra space makes more defect sizes harmless."},
	{KP: 1, Q: "Wire widening on long wires reduces which failure mode?",
		A:   "Opens from missing-material defects",
		D:   []string{"Shorts to neighbours", "Hold violations", "Antenna charging"},
		Why: "A wider wire needs a larger defect to sever it, which lowers open critical area, though it can raise short critical area."},
	{KP: 1, Q: "Why are DFM fixes usually applied only where they do not cost timing or area?",
		A:   "They are optional yield improvements, traded against PPA",
		D:   []string{"They are illegal elsewhere", "They break LVS", "They always make timing better"},
		Why: "Recommended-rule fixes help yield but add capacitance or area, so tools apply them opportunistically."},
	{KP: 3, Q: "Random defect yield and systematic yield are modelled separately. What is a systematic yield loss?",
		A:   "Failures caused by specific layout patterns that print badly, not by random particles",
		D:   []string{"Failures from particles in the air", "Failures that happen only during test", "Failures caused by the package"},
		Why: "Systematic losses repeat wherever the pattern repeats, such as litho hotspots or CMP-sensitive regions; random losses follow defect statistics."},
	{KP: 3, Q: "Total yield is often written as $Y = Y_{sys} \\cdot Y_{rand}$. What does the product assume?",
		A:   "That systematic and random failure mechanisms act independently",
		D:   []string{"That the two yields are equal", "That random yield is always 1", "That systematic failures follow a Poisson law"},
		Why: "Independent mechanisms multiply their survival probabilities."},
	{KP: 3, Q: "Early in a node's life, systematic losses dominate. What tool flow responds to them?",
		A:   "Litho hotspot detection and pattern-based fixing in the layout",
		D:   []string{"Increasing wafer diameter", "Lowering defect density with cleaner rooms only", "Changing the timing constraints"},
		Why: "Pattern-driven failures are found by simulation or pattern matching and fixed by changing the shapes."},
	{KP: 3, Q: "A DFM score combines critical area, via redundancy and recommended-rule compliance. Why should a router use it as a soft cost rather than a hard rule?",
		A:   "Forcing every recommendation would cost area and timing that the design cannot afford",
		D:   []string{"Soft costs are ignored by routers", "DFM scores are random", "Hard rules are always faster"},
		Why: "Soft costs let the router buy yield where it is cheap and keep PPA elsewhere."},
}
