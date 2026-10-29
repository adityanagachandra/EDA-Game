package stat

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

const boltzmannEV = 8.617e-5 // eV per kelvin
const hoursPerYear = 8760.0

// weibull: survival at a time (KP 1), the time to a failure fraction (KP 2), and the mean time
// to failure (KP 3).
func weibull(ctx generator.GeneratorContext) generator.Problem {
	eta := float64(edakit.Int(5, 40))
	beta := float64(edakit.Int(5, 40)) / 10
	stem := edakit.Sprintf("Oxide breakdown times of a device population follow a Weibull distribution with scale $\\eta = %s$ years and shape $\\beta = %s$.", f(eta, 0), f(beta, 1))
	switch kp(ctx) {
	case 1:
		t := float64(edakit.Int(1, int(eta)))
		r := math.Exp(-math.Pow(t/eta, beta))
		ans := edakit.Num(100*r, 3)
		return edakit.Numeric(stem+edakit.Sprintf(" What percentage of devices survive %s years? Round to 3 decimal places.", f(t, 0)),
			ans, edakit.Sprintf("$R(t) = e^{-(t/\\eta)^\\beta} = e^{-(%s)^{%s}} = %s$ percent.", f(t/eta, 4), f(beta, 1), ans), map[string]string{"cum_hazard": f(math.Pow(t/eta, beta), 5)})
	case 2:
		q := edakit.Pick([]float64{0.001, 0.01, 0.05, 0.1, 0.5})
		t := eta * math.Pow(-math.Log1p(-q), 1/beta)
		ans := edakit.Num(t, 3)
		return edakit.Numeric(stem+edakit.Sprintf(" By what time, in years, will %s percent of devices have failed? Round to 3 decimal places.", f(100*q, 1)),
			ans, edakit.Sprintf("Solve $1 - e^{-(t/\\eta)^\\beta} = %s$: $t = \\eta(-\\ln(1-%s))^{1/\\beta} = %s$ years.", f(q, 3), f(q, 3), ans), map[string]string{"ln_term": f(-math.Log1p(-q), 6)})
	default:
		m := eta * math.Gamma(1+1/beta)
		ans := edakit.Num(m, 3)
		return edakit.Numeric(stem+" What is the mean time to failure, in years? Round to 3 decimal places.",
			ans, edakit.Sprintf("$MTTF = \\eta\\,\\Gamma(1 + 1/\\beta) = %s\\,\\Gamma(%s) = %s$ years.", f(eta, 0), f(1+1/beta, 4), ans), map[string]string{"gamma": f(math.Gamma(1+1/beta), 5)})
	}
}

// hazardFIT: Weibull hazard in FIT (KP 1), FIT from a life test (KP 2), and MTTF in years from
// a constant FIT (KP 3).
func hazardFIT(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		eta := float64(edakit.Int(1, 20)) * 1e6
		beta := float64(edakit.Int(5, 30)) / 10
		ty := float64(edakit.Int(1, 15))
		t := ty * hoursPerYear
		h := beta / eta * math.Pow(t/eta, beta-1)
		fit := h * 1e9
		pl := sigPlaces(fit)
		if pl > 4 {
			pl = 4
		}
		ans := edakit.Num(fit, pl)
		return edakit.Numeric(edakit.Sprintf("A failure mechanism follows a Weibull law with $\\eta = %s$ hours and $\\beta = %s$. What is the hazard rate after %s years of operation (8760 hours per year), in FIT (failures per $10^9$ device-hours)? %s", sci(eta), f(beta, 1), f(ty, 0), roundTo(pl)),
			ans, edakit.Sprintf("$h(t) = (\\beta/\\eta)(t/\\eta)^{\\beta-1}$ with $t = %s$ h gives %s per hour, which is %s FIT. With $\\beta %s 1$ the rate is %s.", f(t, 0), sciT(h), ans, map[bool]string{true: ">", false: "\\le"}[beta > 1], map[bool]string{true: "rising: wear-out", false: "flat or falling"}[beta > 1]),
			map[string]string{"h_per_hour": sci(h)})
	case 2:
		n := edakit.Pick([]int{77, 100, 231, 500, 1000})
		hrs := edakit.Pick([]int{168, 500, 1000, 2000})
		af := float64(edakit.Pick([]int{1, 20, 50, 100}))
		r := edakit.Int(1, 5)
		fit := float64(r) / (float64(n) * float64(hrs) * af) * 1e9
		ans := edakit.Num(fit, 2)
		return edakit.Numeric(edakit.Sprintf("A life test runs %d devices for %d hours at a stress with acceleration factor %s, and %d fail. What is the point estimate of the field failure rate, in FIT? Round to 2 decimal places.", n, hrs, f(af, 0), r),
			ans, edakit.Sprintf("Equivalent field device-hours: $%d \\times %d \\times %s = %s$. FIT $= %d / %s \\times 10^9 = %s$.", n, hrs, f(af, 0), sci(float64(n)*float64(hrs)*af), r, sci(float64(n)*float64(hrs)*af), ans),
			map[string]string{"device_hours": sci(float64(n) * float64(hrs) * af)})
	default:
		fit := float64(edakit.Int(1, 200)) * 5
		y := 1e9 / fit / hoursPerYear
		ans := edakit.Num(y, 1)
		return edakit.Numeric(edakit.Sprintf("A component has a constant failure rate of %s FIT. What is its mean time to failure, in years (8760 hours per year)? Round to 1 decimal place.", f(fit, 0)),
			ans, edakit.Sprintf("For a constant rate $MTTF = 1/\\lambda = 10^9/%s = %s$ hours, which is %s years.", f(fit, 0), sci(1e9/fit), ans), map[string]string{"hours": sci(1e9 / fit)})
	}
}

// acceleration: Arrhenius factor (KP 1), equivalent field time (KP 2), and Black's equation
// ratio for electromigration (KP 3).
func acceleration(ctx generator.GeneratorContext) generator.Problem {
	ea := float64(edakit.Int(5, 10)) / 10
	tu := float64(edakit.Pick([]int{55, 70, 85, 105}))
	ts := float64(edakit.Pick([]int{125, 150, 175}))
	af := math.Exp(ea / boltzmannEV * (1/(tu+273.15) - 1/(ts+273.15)))
	switch kp(ctx) {
	case 1:
		ans := edakit.Num(af, 2)
		return edakit.Numeric(edakit.Sprintf("A failure mechanism has activation energy %s eV. What is the Arrhenius acceleration factor between a stress temperature of %s C and a use temperature of %s C? Use $k = 8.617 \\times 10^{-5}$ eV/K and add 273.15 to convert to kelvin. Round to 2 decimal places.", f(ea, 1), f(ts, 0), f(tu, 0)),
			ans, edakit.Sprintf("$AF = \\exp[\\frac{E_a}{k}(\\frac{1}{T_u} - \\frac{1}{T_s})] = \\exp[%s(%s - %s)] = %s$.", f(ea/boltzmannEV, 1), f(1/(tu+273.15), 7), f(1/(ts+273.15), 7), ans), map[string]string{"af": ans})
	case 2:
		hrs := float64(edakit.Pick([]int{168, 500, 1000}))
		y := hrs * af / hoursPerYear
		ans := edakit.Num(y, 2)
		return edakit.Numeric(edakit.Sprintf("A burn-in runs for %s hours at %s C. The mechanism has activation energy %s eV and the product runs at %s C. How many years of field use does the stress represent (8760 hours per year, $k = 8.617 \\times 10^{-5}$ eV/K, kelvin = C + 273.15)? Round to 2 decimal places.", f(hrs, 0), f(ts, 0), f(ea, 1), f(tu, 0)),
			ans, edakit.Sprintf("$AF = %s$, so the test equals $%s \\times %s = %s$ field hours, which is %s years.", f(af, 3), f(hrs, 0), f(af, 3), f(hrs*af, 0), ans), map[string]string{"af": f(af, 3)})
	default:
		n := edakit.Pick([]float64{1, 1.5, 2})
		j1 := float64(edakit.Int(5, 20)) / 10
		j2 := j1 * float64(edakit.Int(11, 30)) / 10
		ratio := math.Pow(j2/j1, n) * math.Exp(ea/boltzmannEV*(1/(tu+273.15)-1/(ts+273.15)))
		ans := edakit.Num(ratio, 2)
		return edakit.Numeric(edakit.Sprintf("Black's equation models electromigration life as $MTTF = A J^{-n} e^{E_a/kT}$ with $n = %s$ and $E_a = %s$ eV. A wire runs in use at %s MA/cm$^2$ and %s C; a stress test uses %s MA/cm$^2$ and %s C. By what factor does the use-condition MTTF exceed the stress MTTF ($k = 8.617 \\times 10^{-5}$ eV/K, kelvin = C + 273.15)? Round to 2 decimal places.", f(n, 1), f(ea, 1), f(j1, 1), f(tu, 0), f(j2, 2), f(ts, 0)),
			ans, edakit.Sprintf("Current factor $(%s/%s)^{%s} = %s$; thermal factor %s. Their product is %s.", f(j2, 2), f(j1, 1), f(n, 1), f(math.Pow(j2/j1, n), 4), f(af, 3), ans), map[string]string{"af_t": f(af, 3)})
	}
}

// weibullFit: the Weibull-plot ordinate (KP 1), the shape from two plotted points (KP 2), and
// Benard's median rank (KP 3).
func weibullFit(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		fr := float64(edakit.Int(1, 95)) / 100
		y := math.Log(-math.Log(1 - fr))
		ans := edakit.Num(y, 4)
		return edakit.Numeric(edakit.Sprintf("On a Weibull plot the vertical coordinate is $\\ln(-\\ln(1-F))$. What is it for a cumulative failure fraction $F = %s$? Round to 4 decimal places.", f(fr, 2)),
			ans, edakit.Sprintf("$-\\ln(1-%s) = %s$ and its natural log is %s. Weibull data plot as a straight line with slope $\\beta$ against $\\ln t$.", f(fr, 2), f(-math.Log(1-fr), 5), ans), map[string]string{"inner": f(-math.Log(1-fr), 5)})
	case 2:
		t1 := float64(edakit.Int(20, 150))
		t2 := t1 * float64(edakit.Int(15, 50)) / 10
		f1 := float64(edakit.Int(2, 30)) / 100
		f2 := math.Min(0.95, f1+float64(edakit.Int(10, 60))/100)
		b := (math.Log(-math.Log(1-f2)) - math.Log(-math.Log(1-f1))) / (math.Log(t2) - math.Log(t1))
		ans := edakit.Num(b, 3)
		return edakit.Numeric(edakit.Sprintf("Accelerated electromigration tests give cumulative failure fractions $F = %s$ at %s hours and $F = %s$ at %s hours. Assuming a Weibull law, what shape parameter $\\beta$ does the line through these two points imply? Round to 3 decimal places.", f(f1, 4), f(t1, 0), f(f2, 4), f(t2, 1)),
			ans, edakit.Sprintf("Ordinates: $\\ln(-\\ln(1-F))$ = %s and %s. Slope $= \\Delta y/\\Delta \\ln t = %s/%s = %s$.", f(math.Log(-math.Log(1-f1)), 4), f(math.Log(-math.Log(1-f2)), 4), f(math.Log(-math.Log(1-f2))-math.Log(-math.Log(1-f1)), 4), f(math.Log(t2/t1), 4), ans),
			map[string]string{"dlnt": f(math.Log(t2/t1), 4)})
	default:
		n := edakit.Scaled(ctx, 5, 30)
		i := edakit.Int(1, n)
		mr := (float64(i) - 0.3) / (float64(n) + 0.4)
		ans := edakit.Num(mr, 4)
		return edakit.Numeric(edakit.Sprintf("%d devices are stressed until all fail, and the failure times are sorted. Using Benard's median-rank approximation, what cumulative failure fraction is plotted for failure number %d? Round to 4 decimal places.", n, i),
			ans, edakit.Sprintf("$F_i \\approx (i - 0.3)/(n + 0.4) = %s/%s = %s$. Using $i/n$ would put the last failure at $F = 1$, which cannot be plotted.", f(float64(i)-0.3, 1), f(float64(n)+0.4, 1), ans), map[string]string{"rank": ans})
	}
}

// systemReliability: a series system of Weibull parts (KP 1), triple modular redundancy (KP 2),
// and parallel redundancy (KP 3).
func systemReliability(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 2, 50)
		eta := float64(edakit.Int(20, 200))
		beta := float64(edakit.Int(10, 30)) / 10
		t := float64(edakit.Int(5, 15))
		r := math.Exp(-float64(n) * math.Pow(t/eta, beta))
		ans := edakit.Num(100*r, 3)
		return edakit.Numeric(edakit.Sprintf("A chip contains %d identical, independent wear-out sites, and it fails when any one fails. Each site has a Weibull life with $\\eta = %s$ years and $\\beta = %s$. What is the probability, in percent, that the chip survives %s years? Round to 3 decimal places.", n, f(eta, 0), f(beta, 1), f(t, 0)),
			ans, edakit.Sprintf("Series reliabilities multiply: $R = e^{-n(t/\\eta)^\\beta} = e^{-%s} = %s$ percent. A series of $n$ Weibull parts is Weibull with scale $\\eta n^{-1/\\beta}$.", f(float64(n)*math.Pow(t/eta, beta), 5), ans), map[string]string{"cum_hazard": f(float64(n)*math.Pow(t/eta, beta), 5)})
	case 2:
		r := float64(edakit.Int(50, 999)) / 1000
		tmr := 3*r*r - 2*r*r*r
		ans := edakit.Num(tmr, 5)
		return edakit.Numeric(edakit.Sprintf("A module has reliability %s over a mission. Three copies vote through a perfect majority voter (triple modular redundancy). What is the reliability of the voted system? Round to 5 decimal places.", f(r, 3)),
			ans, edakit.Sprintf("At least two of three must work: $3R^2(1-R) + R^3 = 3R^2 - 2R^3 = %s$. TMR helps only when $R > 0.5$.", ans), map[string]string{"r": f(r, 3)})
	default:
		r := float64(edakit.Int(50, 98)) / 100
		m := edakit.Int(2, 4)
		p := 1 - math.Pow(1-r, float64(m))
		ans := edakit.Num(p, 6)
		return edakit.Numeric(edakit.Sprintf("A link has %d independent redundant paths, each working with probability %s, and works if any path works. What is the link reliability? Round to 6 decimal places.", m, f(r, 2)),
			ans, edakit.Sprintf("It fails only if all paths fail: $1 - (1-%s)^{%d} = %s$.", f(r, 2), m, ans), map[string]string{"q": f(1-r, 2)})
	}
}
