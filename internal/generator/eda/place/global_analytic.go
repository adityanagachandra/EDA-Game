package place

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type rect struct{ x, y, w, h int }

func overlap1(a0, a1, b0, b1 int) int {
	lo, hi := maxI(a0, b0), minI(a1, b1)
	if hi > lo {
		return hi - lo
	}
	return 0
}

func (r rect) overlapArea(o rect) int {
	return overlap1(r.x, r.x+r.w, o.x, o.x+o.w) * overlap1(r.y, r.y+r.h, o.y, o.y+o.h)
}

// density computes bin density, bin overflow and total overflow on a 2 x 2 bin grid.
func density(ctx generator.GeneratorContext) generator.Problem {
	B := edakit.Pick([]int{8, 10, 12})
	n := edakit.Scaled(ctx, 3, 6)
	cells := make([]rect, n)
	var desc []string
	total := 0
	for i := range cells {
		w, h := edakit.Int(2, 8), edakit.Int(2, 8)
		cells[i] = rect{edakit.Int(0, 2*B-w), edakit.Int(0, 2*B-h), w, h}
		total += w * h
		desc = append(desc, edakit.Sprintf("a %d x %d cell with lower-left corner (%d, %d)", w, h, cells[i].x, cells[i].y))
	}
	bins := []rect{{0, 0, B, B}, {B, 0, B, B}, {0, B, B, B}, {B, B, B, B}}
	names := []string{"lower-left", "lower-right", "upper-left", "upper-right"}
	area := make([]int, 4)
	for b := range bins {
		for _, c := range cells {
			area[b] += bins[b].overlapArea(c)
		}
	}
	t := edakit.Pick([]float64{0.5, 0.6, 0.7, 0.8, 0.9})
	head := edakit.Sprintf("A %d x %d um placement region is divided into four %d x %d um bins. Cells (width x height, um): %s. Cells may overlap each other in global placement. ", 2*B, 2*B, B, B, strings.Join(desc, "; "))
	switch ctx.KPOrdinal {
	case 1:
		b := edakit.Int(0, 3)
		dd := float64(area[b]) / float64(B*B)
		ans := edakit.Num(dd, 4)
		return edakit.Numeric(head+edakit.Sprintf("What is the density of the %s bin, the cell area inside it divided by the bin area? Round to 4 decimal places.", names[b]), ans,
			edakit.Sprintf("Clip each cell to the bin and add the overlap areas: %d um^2 over %d um^2 gives %s.", area[b], B*B, ans), map[string]string{"bin_area": edakit.IntStr(area[b])})
	case 2:
		b := edakit.Int(0, 3)
		of := math.Max(0, float64(area[b])-t*float64(B*B))
		ans := edakit.Num(of, 2)
		return edakit.Numeric(head+edakit.Sprintf("The target density is %s. What is the overflow of the %s bin, the cell area in excess of the target times the bin area, in um^2 (0 if none)? Round to 2 decimal places.", edakit.Num(t, 1), names[b]), ans,
			edakit.Sprintf("The bin holds %d um^2 against an allowance of %s x %d = %s um^2, so its overflow is %s um^2.", area[b], edakit.Num(t, 1), B*B, edakit.Num(t*float64(B*B), 1), ans), map[string]string{"bin_area": edakit.IntStr(area[b])})
	}
	sum := 0.0
	var parts []string
	for b := range bins {
		of := math.Max(0, float64(area[b])-t*float64(B*B))
		sum += of
		parts = append(parts, edakit.Sprintf("%s %d", names[b], area[b]))
	}
	ans := edakit.Num(sum/float64(total), 4)
	return edakit.Numeric(head+edakit.Sprintf("The target density is %s. The placer reports total overflow as the sum over bins of the excess cell area above target, divided by the total cell area. What is it? Round to 4 decimal places.", edakit.Num(t, 1)), ans,
		edakit.Sprintf("Bin cell areas (um^2): %s. The excess over %s each sums to %s um^2, and the total cell area is %d um^2, giving %s.", strings.Join(parts, ", "), edakit.Num(t*float64(B*B), 1), edakit.Num(sum, 2), total, ans),
		map[string]string{"excess": edakit.Num(sum, 2), "total": edakit.IntStr(total)})
}

// lseMax is gamma * ln(sum exp(x/gamma)), computed stably.
func lseMax(xs []float64, g float64) float64 {
	m := xs[0]
	for _, x := range xs {
		m = math.Max(m, x)
	}
	s := 0.0
	for _, x := range xs {
		s += math.Exp((x - m) / g)
	}
	return m + g*math.Log(s)
}

func neg(xs []float64) []float64 {
	out := make([]float64, len(xs))
	for i, x := range xs {
		out[i] = -x
	}
	return out
}

func lseWL(xs []float64, g float64) float64 { return lseMax(xs, g) + lseMax(neg(xs), g) }

func softmax(xs []float64, g float64) []float64 {
	m := xs[0]
	for _, x := range xs {
		m = math.Max(m, x)
	}
	s := 0.0
	out := make([]float64, len(xs))
	for i, x := range xs {
		out[i] = math.Exp((x - m) / g)
		s += out[i]
	}
	for i := range out {
		out[i] /= s
	}
	return out
}

func waMax(xs []float64, g float64) float64 {
	w := softmax(xs, g)
	s := 0.0
	for i, x := range xs {
		s += w[i] * x
	}
	return s
}

func waWL(xs []float64, g float64) float64 { return waMax(xs, g) + waMax(neg(xs), g) }

func drawCoords(ctx generator.GeneratorContext) ([]float64, []int, float64) {
	k := edakit.Scaled(ctx, 3, 5)
	ints := make([]int, k)
	fl := make([]float64, k)
	for i := range ints {
		ints[i] = edakit.Int(0, 20)
		fl[i] = float64(ints[i])
	}
	if span(ints) == 0 {
		ints[0], fl[0] = ints[0]+3, fl[0]+3
	}
	g := edakit.Pick([]float64{1, 2, 3, 4, 5})
	return fl, ints, g
}

// smoothWL evaluates the log-sum-exp smooth maximum and wirelength.
func smoothWL(ctx generator.GeneratorContext) generator.Problem {
	xs, ints, g := drawCoords(ctx)
	head := edakit.Sprintf("A net's pins are at x = %s (um) and the smoothing parameter is $\\gamma$ = %s um. ", fmtInts(ints), edakit.Num(g, 0))
	hp := float64(span(ints))
	switch ctx.KPOrdinal {
	case 1:
		v := lseMax(xs, g)
		ans := edakit.Num(v, 3)
		return edakit.Numeric(head+"The smooth maximum is $\\gamma \\ln \\sum_i e^{x_i/\\gamma}$. What is its value? Round to 3 decimal places.", ans,
			edakit.Sprintf("Factor out the largest term: %s + %s ln(sum of $e^{(x_i - %s)/\\gamma}$) = %s, a little above the true maximum %d.", edakit.Num(float64(maxI(ints...)), 0), edakit.Num(g, 0), edakit.Num(float64(maxI(ints...)), 0), ans, maxI(ints...)),
			map[string]string{"smax": ans})
	case 3:
		e := lseWL(xs, g) - hp
		ans := edakit.Num(e, 3)
		return edakit.Numeric(head+"The log-sum-exp wirelength in x is $\\gamma \\ln \\sum e^{x_i/\\gamma} + \\gamma \\ln \\sum e^{-x_i/\\gamma}$. By how much does it overestimate the net's exact x span? Round to 3 decimal places.", ans,
			edakit.Sprintf("LSE gives %s and the span is %s, so the error is %s. Each of the two terms errs by at most $\\gamma \\ln n$ = %s, so the total error is at most %s.", edakit.Num(lseWL(xs, g), 3), edakit.Num(hp, 0), ans, edakit.Num(g*math.Log(float64(len(xs))), 3), edakit.Num(2*g*math.Log(float64(len(xs))), 3)),
			map[string]string{"lse": edakit.Num(lseWL(xs, g), 3)})
	}
	v := lseWL(xs, g)
	ans := edakit.Num(v, 3)
	return edakit.Numeric(head+"The log-sum-exp wirelength in x is $\\gamma \\ln \\sum e^{x_i/\\gamma} + \\gamma \\ln \\sum e^{-x_i/\\gamma}$. What is its value? Round to 3 decimal places.", ans,
		edakit.Sprintf("The smooth maximum is %s and the smooth negated minimum is %s; their sum %s slightly exceeds the exact span %s.", edakit.Num(lseMax(xs, g), 3), edakit.Num(lseMax(neg(xs), g), 3), ans, edakit.Num(hp, 0)),
		map[string]string{"lse": ans})
}

// lseGradient computes partial derivatives of the LSE wirelength.
func lseGradient(ctx generator.GeneratorContext) generator.Problem {
	xs, ints, g := drawCoords(ctx)
	i := edakit.Int(0, len(xs)-1)
	head := edakit.Sprintf("A net's pins are at x = %s (um), and its LSE wirelength in x is $W = \\gamma \\ln \\sum e^{x_j/\\gamma} + \\gamma \\ln \\sum e^{-x_j/\\gamma}$ with $\\gamma$ = %s. ", fmtInts(ints), edakit.Num(g, 0))
	sp, sm := softmax(xs, g), softmax(neg(xs), g)
	switch ctx.KPOrdinal {
	case 1:
		ans := edakit.Num(sp[i], 4)
		return edakit.Numeric(head+edakit.Sprintf("What is the derivative of the first (smooth-maximum) term with respect to the pin at x = %d, listed in position %d? Round to 4 decimal places.", ints[i], i+1), ans,
			edakit.Sprintf("The derivative of $\\gamma \\ln \\sum e^{x_j/\\gamma}$ is the softmax weight $e^{x_i/\\gamma}/\\sum_j e^{x_j/\\gamma}$ = %s. Pins near the maximum get most of the weight.", ans), map[string]string{"softmax": ans})
	case 3:
		xs2, ints2, _ := drawCoords(ctx)
		j := edakit.Int(0, len(xs2)-1)
		xs2[j], ints2[j] = xs[i], ints[i]
		sp2, sm2 := softmax(xs2, g), softmax(neg(xs2), g)
		gr := sp[i] - sm[i] + sp2[j] - sm2[j]
		ans := edakit.Num(gr, 4)
		return edakit.Numeric(head+edakit.Sprintf("The cell whose pin is listed in position %d also drives a second net with pins at x = %s, where it is listed in position %d (same $\\gamma$). Treating the cell as one point for both pins, what is $\\partial (W_1 + W_2)/\\partial x$ for this cell? Round to 4 decimal places.", i+1, fmtInts(ints2), j+1), ans,
			edakit.Sprintf("Net 1 contributes %s - %s and net 2 contributes %s - %s; the gradient of a sum of nets is the sum of their gradients: %s.", edakit.Num(sp[i], 4), edakit.Num(sm[i], 4), edakit.Num(sp2[j], 4), edakit.Num(sm2[j], 4), ans),
			map[string]string{"g1": edakit.Num(sp[i]-sm[i], 4), "g2": edakit.Num(sp2[j]-sm2[j], 4)})
	}
	gr := sp[i] - sm[i]
	ans := edakit.Num(gr, 4)
	return edakit.Numeric(head+edakit.Sprintf("What is $\\partial W/\\partial x_i$ for the pin at x = %d, listed in position %d? Round to 4 decimal places.", ints[i], i+1), ans,
		edakit.Sprintf("The gradient is $e^{x_i/\\gamma}/\\sum e^{x_j/\\gamma} - e^{-x_i/\\gamma}/\\sum e^{-x_j/\\gamma}$ = %s - %s = %s. Over all pins the gradients sum to zero, because shifting the whole net does not change its length.", edakit.Num(sp[i], 4), edakit.Num(sm[i], 4), ans),
		map[string]string{"g": ans})
}

// waWirelength evaluates the weighted-average wirelength.
func waWirelength(ctx generator.GeneratorContext) generator.Problem {
	xs, ints, g := drawCoords(ctx)
	head := edakit.Sprintf("A net's pins are at x = %s (um) and $\\gamma$ = %s. The weighted-average model estimates the maximum as $\\sum x_i e^{x_i/\\gamma} / \\sum e^{x_i/\\gamma}$ and the minimum as $\\sum x_i e^{-x_i/\\gamma} / \\sum e^{-x_i/\\gamma}$, and the x wirelength as their difference. ", fmtInts(ints), edakit.Num(g, 0))
	hp := float64(span(ints))
	switch ctx.KPOrdinal {
	case 1:
		v := waMax(xs, g)
		ans := edakit.Num(v, 3)
		return edakit.Numeric(head+"What is the weighted-average estimate of the maximum? Round to 3 decimal places.", ans,
			edakit.Sprintf("Weight each pin by its softmax weight and average: %s, just below the true maximum %d.", ans, maxI(ints...)), map[string]string{"wamax": ans})
	case 3:
		e := hp - waWL(xs, g)
		ans := edakit.Num(e, 3)
		return edakit.Numeric(head+"By how much does the weighted-average wirelength underestimate the exact x span? Round to 3 decimal places.", ans,
			edakit.Sprintf("WA gives %s against a span of %s, an underestimate of %s. For comparison the LSE model overestimates by %s at the same $\\gamma$.", edakit.Num(waWL(xs, g), 3), edakit.Num(hp, 0), ans, edakit.Num(lseWL(xs, g)-hp, 3)),
			map[string]string{"wa": edakit.Num(waWL(xs, g), 3)})
	}
	v := waWL(xs, g)
	ans := edakit.Num(v, 3)
	return edakit.Numeric(head+"What is the weighted-average x wirelength? Round to 3 decimal places.", ans,
		edakit.Sprintf("The smooth maximum is %s and the smooth minimum %s, so the estimate is %s, a little under the exact span %s.", edakit.Num(waMax(xs, g), 3), edakit.Num(-waMax(neg(xs), g), 3), ans, edakit.Num(hp, 0)),
		map[string]string{"wa": ans})
}

// nonlinear covers one penalty-gradient step, the lambda schedule and lambda initialisation.
func nonlinear(ctx generator.GeneratorContext) generator.Problem {
	switch ctx.KPOrdinal {
	case 2:
		l0 := edakit.Pick([]float64{0.0001, 0.001, 0.01})
		mu := edakit.Pick([]float64{1.05, 1.1, 1.2, 1.5, 2})
		target := l0 * float64(edakit.Pick([]int{10, 50, 100, 1000}))
		k, l := 0, l0
		for l < target {
			l *= mu
			k++
		}
		return edakit.Numeric(edakit.Sprintf("A nonlinear placer minimises $W(x) + \\lambda D(x)$. It starts with $\\lambda_0$ = %s and multiplies $\\lambda$ by %s after each outer iteration. After how many outer iterations does $\\lambda$ first reach at least %s?", edakit.Num(l0, 4), edakit.Num(mu, 2), edakit.Num(target, 4)), edakit.IntStr(k),
			edakit.Sprintf("We need the smallest k with %s x %s^k >= %s: k = ceil(ln(%s)/ln(%s)) = %d.", edakit.Num(l0, 4), edakit.Num(mu, 2), edakit.Num(target, 4), edakit.Num(target/l0, 0), edakit.Num(mu, 2), k),
			map[string]string{"k": edakit.IntStr(k)})
	case 3:
		n := 4
		gw, gd := make([]float64, n), make([]float64, n)
		sw, sd := 0.0, 0.0
		for i := 0; i < n; i++ {
			gw[i] = float64(edakit.Int(-30, 30))
			gd[i] = float64(edakit.Int(-400, 400))
			sw += math.Abs(gw[i])
			sd += math.Abs(gd[i])
		}
		if sd == 0 {
			gd[0], sd = 100, 100
		}
		l := sw / sd
		ans := edakit.Num(l, 4)
		return edakit.Numeric(edakit.Sprintf("At the start of nonlinear placement, the x components of the wirelength gradient for four cells are %s and those of the density gradient are %s. The placer initialises $\\lambda_0 = \\|\\nabla W\\|_1 / \\|\\nabla D\\|_1$ so the two forces start comparable. What is $\\lambda_0$? Round to 4 decimal places.", fmtFloats(gw, 0), fmtFloats(gd, 0)), ans,
			edakit.Sprintf("The L1 norms are %s and %s, so $\\lambda_0$ = %s.", edakit.Num(sw, 0), edakit.Num(sd, 0), ans), map[string]string{"l0": ans})
	}
	x := float64(edakit.Int(10, 90))
	gw := float64(edakit.Int(-20, 20))
	gd := float64(edakit.Int(-50, 50))
	l := edakit.Pick([]float64{0.1, 0.2, 0.5, 1, 2})
	eta := edakit.Pick([]float64{0.05, 0.1, 0.2, 0.5})
	xn := x - eta*(gw+l*gd)
	ans := edakit.Num(xn, 3)
	return edakit.Numeric(edakit.Sprintf("A nonlinear placer minimises $W(x) + \\lambda D(x)$ with $\\lambda$ = %s. A cell sits at x = %s um, where $\\partial W/\\partial x$ = %s and $\\partial D/\\partial x$ = %s. With step size %s, where does one gradient-descent step move it? Round to 3 decimal places.", edakit.Num(l, 1), edakit.Num(x, 0), edakit.Num(gw, 0), edakit.Num(gd, 0), edakit.Num(eta, 2)), ans,
		edakit.Sprintf("The combined gradient is %s + %s x %s = %s, so x becomes %s - %s x %s = %s.", edakit.Num(gw, 0), edakit.Num(l, 1), edakit.Num(gd, 0), edakit.Num(gw+l*gd, 2), edakit.Num(x, 0), edakit.Num(eta, 2), edakit.Num(gw+l*gd, 2), ans),
		map[string]string{"grad": edakit.Num(gw+l*gd, 2)})
}

// electrostatic computes the electrostatic energy, the force on a cell and the filler count.
func electrostatic(ctx generator.GeneratorContext) generator.Problem {
	switch ctx.KPOrdinal {
	case 2:
		q := float64(edakit.Int(1, 12))
		ex, ey := float64(edakit.Int(-20, 20))/10, float64(edakit.Int(-20, 20))/10
		f := q * math.Hypot(ex, ey)
		ans := edakit.Num(f, 3)
		return edakit.Numeric(edakit.Sprintf("In electrostatics-based placement a cell of area %s um^2 is a charge of %s. The electric field at its centre, from the density solve, is $E = (%s, %s)$. What is the magnitude of the density force on the cell? Round to 3 decimal places.", edakit.Num(q, 0), edakit.Num(q, 0), edakit.Num(ex, 1), edakit.Num(ey, 1)), ans,
			edakit.Sprintf("The force is $qE$, so its magnitude is %s x sqrt((%s)^2 + (%s)^2) = %s. Large cells are pushed harder, which is why macros and standard cells can be spread in one model.", edakit.Num(q, 0), edakit.Num(ex, 1), edakit.Num(ey, 1), ans), map[string]string{"f": ans})
	case 3:
		region := float64(100 * edakit.Int(50, 400))
		fixed := float64(100 * edakit.Int(0, 100))
		if fixed >= region/2 {
			fixed = region / 4
		}
		ws := region - fixed
		rho := edakit.Pick([]float64{0.7, 0.8, 0.9, 1})
		mov := math.Floor(ws * edakit.Pick([]float64{0.4, 0.5, 0.6}))
		fa := edakit.Pick([]float64{2, 4, 5, 10})
		cnt := int(math.Floor((rho*ws-mov)/fa + 1e-9))
		return edakit.Numeric(edakit.Sprintf("A placement region has area %s um^2, of which fixed macros cover %s um^2. The movable cells total %s um^2 and the target density is %s. The placer adds filler cells of %s um^2 each so that movable plus filler area equals the target density times the free area. How many whole filler cells does it add?", edakit.Num(region, 0), edakit.Num(fixed, 0), edakit.Num(mov, 0), edakit.Num(rho, 1), edakit.Num(fa, 0)), edakit.IntStr(cnt),
			edakit.Sprintf("The free area is %s um^2, so the fillers need %s x %s - %s = %s um^2, which is floor(%s / %s) = %d fillers. They occupy whitespace so the density force stops pushing real cells into it.", edakit.Num(ws, 0), edakit.Num(rho, 1), edakit.Num(ws, 0), edakit.Num(mov, 0), edakit.Num(rho*ws-mov, 1), edakit.Num(rho*ws-mov, 1), edakit.Num(fa, 0), cnt),
			map[string]string{"filler_area": edakit.Num(rho*ws-mov, 1)})
	}
	n := edakit.Scaled(ctx, 3, 5)
	qs, ps := make([]float64, n), make([]float64, n)
	s := 0.0
	for i := range qs {
		qs[i] = float64(edakit.Int(1, 10))
		ps[i] = float64(edakit.Int(-20, 60)) / 10
		s += qs[i] * ps[i]
	}
	ans := edakit.Num(s/2, 3)
	return edakit.Numeric(edakit.Sprintf("In electrostatics-based placement each cell is a charge equal to its area, and the density penalty is the system's potential energy $N = \\frac{1}{2}\\sum_i q_i \\psi_i$. Cells have charges %s and sit at points where the potential is %s. What is N? Round to 3 decimal places.", fmtFloats(qs, 0), fmtFloats(ps, 1)), ans,
		edakit.Sprintf("Sum charge times potential, %s, and halve it: %s.", edakit.Num(s, 2), ans), map[string]string{"sum": edakit.Num(s, 2)})
}

// poissonFFT solves one cosine mode of the density Poisson equation and compares costs.
func poissonFFT(ctx generator.GeneratorContext) generator.Problem {
	if ctx.KPOrdinal == 3 {
		k := edakit.Int(5, 11)
		m := 1 << k
		ratio := m * m / (2 * k)
		if (m*m)%(2*k) != 0 {
			r := float64(m*m) / float64(2*k)
			ans := edakit.Num(r, 1)
			return edakit.Numeric(edakit.Sprintf("A density grid has m x m bins with m = %d. Count a direct cosine transform as $m^4$ multiply-adds and an FFT-based transform as $m^2 \\log_2(m^2)$. How many times fewer operations does the FFT need? Round to 1 decimal place.", m), ans,
				edakit.Sprintf("The ratio is $m^4 / (m^2 \\log_2 m^2) = m^2/(2\\log_2 m)$ = %d / %d = %s.", m*m, 2*k, ans), map[string]string{"ratio": ans})
		}
		return edakit.Numeric(edakit.Sprintf("A density grid has m x m bins with m = %d. Count a direct cosine transform as $m^4$ multiply-adds and an FFT-based transform as $m^2 \\log_2(m^2)$. How many times fewer operations does the FFT need?", m), edakit.IntStr(ratio),
			edakit.Sprintf("The ratio is $m^4 / (m^2 \\log_2 m^2) = m^2/(2\\log_2 m)$ = %d / %d = %d.", m*m, 2*k, ratio), map[string]string{"ratio": edakit.IntStr(ratio)})
	}
	m := edakit.Pick([]int{16, 32, 64})
	u, v := edakit.Int(1, 6), edakit.Int(0, 6)
	a := float64(edakit.Int(1, 9)) / 10
	wu, wv := math.Pi*float64(u)/float64(m), math.Pi*float64(v)/float64(m)
	den := wu*wu + wv*wv
	head := edakit.Sprintf("A placer solves $\\nabla^2 \\psi = -\\rho$ on an m x m bin grid (m = %d, unit bins) with a cosine expansion, $\\omega_u = \\pi u / m$. The zero-mean density contains the mode $%s \\cos(\\omega_%d x)\\cos(\\omega_%d y)$. ", m, edakit.Num(a, 1), u, v)
	if ctx.KPOrdinal == 2 {
		x, y := float64(edakit.Int(0, m-1))+0.5, float64(edakit.Int(0, m-1))+0.5
		ex := a * wu / den * math.Sin(wu*x) * math.Cos(wv*y)
		ans := edakit.Num(ex, 4)
		return edakit.Numeric(head+edakit.Sprintf("The field is $E = -\\nabla \\psi$. What is this mode's contribution to $E_x$ at the bin centre (%s, %s)? Round to 4 decimal places.", edakit.Num(x, 1), edakit.Num(y, 1)), ans,
			edakit.Sprintf("The potential mode is $\\frac{a}{\\omega_u^2+\\omega_v^2}\\cos\\cos$, so $E_x = \\frac{a\\,\\omega_u}{\\omega_u^2+\\omega_v^2}\\sin(\\omega_u x)\\cos(\\omega_v y)$. With $\\omega_u$ = %s and $\\omega_v$ = %s this is %s.", edakit.Num(wu, 5), edakit.Num(wv, 5), ans),
			map[string]string{"wu": edakit.Num(wu, 5), "wv": edakit.Num(wv, 5)})
	}
	c := a / den
	ans := edakit.Num(c, 4)
	return edakit.Numeric(head+"What is the coefficient of the matching mode $\\cos(\\omega_u x)\\cos(\\omega_v y)$ in the potential $\\psi$? Round to 4 decimal places.", ans,
		edakit.Sprintf("Each cosine mode is an eigenfunction of the Laplacian with eigenvalue $-(\\omega_u^2 + \\omega_v^2)$, so the coefficient is divided by %s: %s / %s = %s.", edakit.Num(den, 6), edakit.Num(a, 1), edakit.Num(den, 6), ans),
		map[string]string{"den": edakit.Num(den, 6)})
}

// nesterov covers the Lipschitz step estimate, the momentum coefficient and one full update.
func nesterov(ctx generator.GeneratorContext) generator.Problem {
	switch ctx.KPOrdinal {
	case 1:
		v0 := []float64{float64(edakit.Int(0, 50)), float64(edakit.Int(0, 50))}
		v1 := []float64{v0[0] + float64(edakit.Int(-8, 8)), v0[1] + float64(edakit.Int(1, 8))}
		g0 := []float64{float64(edakit.Int(-20, 20)), float64(edakit.Int(-20, 20))}
		g1 := []float64{g0[0] + float64(edakit.Int(-15, 15)), g0[1] + float64(edakit.Int(1, 15))}
		dv := math.Hypot(v1[0]-v0[0], v1[1]-v0[1])
		dg := math.Hypot(g1[0]-g0[0], g1[1]-g0[1])
		ans := edakit.Num(dv/dg, 4)
		return edakit.Numeric(edakit.Sprintf("ePlace-style Nesterov placement estimates the step length as the inverse of a local Lipschitz constant, $\\alpha_k = \\|v_k - v_{k-1}\\| / \\|\\nabla f(v_k) - \\nabla f(v_{k-1})\\|$. For a one-cell toy, $v_{k-1} = (%s, %s)$, $v_k = (%s, %s)$, $\\nabla f(v_{k-1}) = (%s, %s)$ and $\\nabla f(v_k) = (%s, %s)$. What is $\\alpha_k$? Round to 4 decimal places.",
			edakit.Num(v0[0], 0), edakit.Num(v0[1], 0), edakit.Num(v1[0], 0), edakit.Num(v1[1], 0), edakit.Num(g0[0], 0), edakit.Num(g0[1], 0), edakit.Num(g1[0], 0), edakit.Num(g1[1], 0)), ans,
			edakit.Sprintf("The position changed by %s and the gradient by %s, so the local Lipschitz estimate is %s and the step is %s.", edakit.Num(dv, 4), edakit.Num(dg, 4), edakit.Num(dg/dv, 4), ans),
			map[string]string{"dv": edakit.Num(dv, 4), "dg": edakit.Num(dg, 4)})
	case 2:
		a := 1.0
		k := edakit.Int(0, 4)
		for i := 0; i < k; i++ {
			a = (1 + math.Sqrt(4*a*a+1)) / 2
		}
		an := (1 + math.Sqrt(4*a*a+1)) / 2
		beta := (a - 1) / an
		ans := edakit.Num(beta, 4)
		return edakit.Numeric(edakit.Sprintf("Nesterov's method keeps a sequence with $a_{k+1} = (1 + \\sqrt{4a_k^2 + 1})/2$ and uses momentum coefficient $(a_k - 1)/a_{k+1}$. If $a_k$ = %s, what is the momentum coefficient? Round to 4 decimal places.", edakit.Num(a, 4)), ans,
			edakit.Sprintf("$a_{k+1}$ = (1 + sqrt(4 x %s^2 + 1))/2 = %s, so the coefficient is (%s - 1)/%s = %s. It starts at 0 and climbs toward 1.", edakit.Num(a, 4), edakit.Num(an, 4), edakit.Num(a, 4), edakit.Num(an, 4), ans),
			map[string]string{"a_next": edakit.Num(an, 4)})
	}
	uk := float64(edakit.Int(10, 90))
	vk := uk + float64(edakit.Int(-5, 5))
	g := float64(edakit.Int(-20, 20))
	alpha := edakit.Pick([]float64{0.1, 0.2, 0.25, 0.5})
	a := edakit.Pick([]float64{1, 1.618, 2.1935, 2.7486})
	an := (1 + math.Sqrt(4*a*a+1)) / 2
	u1 := vk - alpha*g
	v1 := u1 + (a-1)/an*(u1-uk)
	ans := edakit.Num(v1, 3)
	return edakit.Numeric(edakit.Sprintf("One Nesterov step updates $u_{k+1} = v_k - \\alpha \\nabla f(v_k)$ and $v_{k+1} = u_{k+1} + \\frac{a_k - 1}{a_{k+1}}(u_{k+1} - u_k)$ with $a_{k+1} = (1 + \\sqrt{4a_k^2+1})/2$. For one coordinate, $u_k$ = %s, $v_k$ = %s, $\\nabla f(v_k)$ = %s, $\\alpha$ = %s and $a_k$ = %s. What is $v_{k+1}$? Round to 3 decimal places.", edakit.Num(uk, 0), edakit.Num(vk, 0), edakit.Num(g, 0), edakit.Num(alpha, 2), edakit.Num(a, 4)), ans,
		edakit.Sprintf("$u_{k+1}$ = %s - %s x %s = %s. $a_{k+1}$ = %s, so the momentum coefficient is %s and $v_{k+1}$ = %s.", edakit.Num(vk, 0), edakit.Num(alpha, 2), edakit.Num(g, 0), edakit.Num(u1, 3), edakit.Num(an, 4), edakit.Num((a-1)/an, 4), ans),
		map[string]string{"u_next": edakit.Num(u1, 3)})
}
