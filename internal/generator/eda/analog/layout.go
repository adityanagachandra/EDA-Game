package analog

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func pelgrom(ctx generator.GeneratorContext) generator.Problem {
	avt := unif(1, 5, 0.1)
	switch whichKP(ctx) {
	case 1:
		w := unif(0.5, 20, 0.5)
		l := unif(0.1, 4, 0.1)
		s := avt / math.Sqrt(w*l)
		q := edakit.Sprintf("A process has Pelgrom coefficient $A_{VT} = %s$ mV$\\cdot\\mu$m. What is the standard deviation of the threshold-voltage difference between two identical devices with $W = %s\\ \\mu$m and $L = %s\\ \\mu$m, in mV? Round to 2 decimal places.", f(avt, 1), f(w, 1), f(l, 1))
		return num(q, f(s, 2), edakit.Sprintf("$\\sigma(\\Delta V_T) = A_{VT}/\\sqrt{WL} = %s/\\sqrt{%s} = %s$ mV.", f(avt, 1), f(w*l, 2), f(s, 2)), map[string]string{"area": f(w*l, 2)})
	case 2:
		s := unif(0.2, 3, 0.05)
		a := (avt / s) * (avt / s)
		q := edakit.Sprintf("With $A_{VT} = %s$ mV$\\cdot\\mu$m, what gate area $WL$ does each device of a matched pair need so that $\\sigma(\\Delta V_T) = %s$ mV, in square micrometres? Round to 2 decimal places.", f(avt, 1), f(s, 2))
		return num(q, f(a, 2), edakit.Sprintf("Invert Pelgrom's law: $WL = (A_{VT}/\\sigma)^2 = (%s/%s)^2 = %s\\ \\mu\\mathrm{m}^2$. Halving sigma costs four times the area.", f(avt, 1), f(s, 2), f(a, 2)), nil)
	default:
		s1 := unif(0.5, 6, 0.1)
		k := pick(2.0, 3.0, 4.0, 6.0, 8.0, 9.0, 16.0)
		up := edakit.Int(0, 1) == 0
		var s2 float64
		var how string
		if up {
			s2 = s1 / math.Sqrt(k)
			how = "multiplied"
		} else {
			s2 = s1 * math.Sqrt(k)
			how = "divided"
		}
		q := edakit.Sprintf("A matched pair shows $\\sigma(\\Delta V_T) = %s$ mV. Both devices' gate areas are %s by %s, keeping the same shape. What is the new sigma, in mV? Round to 3 decimal places.", f(s1, 1), how, f(k, 0))
		return num(q, f(s2, 3), edakit.Sprintf("Sigma scales as $1/\\sqrt{WL}$, so it changes by a factor $\\sqrt{%s} = %s$: the new value is %s mV.", f(k, 0), f(math.Sqrt(k), 3), f(s2, 3)), nil)
	}
}

func offsetBudget(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		gmid := float64(edakit.Int(4, 24))
		svt := unif(0.5, 5, 0.1)
		s := gmid * svt / 10
		q := edakit.Sprintf("Two mirror devices biased at $g_m/I_D = %s\\ \\mathrm{V}^{-1}$ have $\\sigma(\\Delta V_T) = %s$ mV. Ignoring current-factor mismatch, what is the standard deviation of their relative current mismatch $\\Delta I/I$, in percent? Round to 2 decimal places.", f(gmid, 0), f(svt, 1))
		return num(q, f(s, 2), edakit.Sprintf("A threshold shift acts like a gate-voltage error, so $\\sigma(\\Delta I/I) = (g_m/I_D)\\sigma(\\Delta V_T) = %s \\cdot %s \\times 10^{-3} = %s$%%. Mirrors want low $g_m/I_D$.", f(gmid, 0), f(svt, 1), f(s, 2)), nil)
	case 2:
		gmid := float64(edakit.Int(4, 20))
		svt := unif(0.5, 4, 0.1)
		sb := unif(0.2, 2, 0.1)
		vt := gmid * svt / 10
		s := math.Sqrt(vt*vt + sb*sb)
		q := edakit.Sprintf("A current-mirror pair runs at $g_m/I_D = %s\\ \\mathrm{V}^{-1}$ with $\\sigma(\\Delta V_T) = %s$ mV and an independent current-factor mismatch $\\sigma(\\Delta\\beta/\\beta) = %s$%%. What is $\\sigma(\\Delta I/I)$ in percent? Round to 2 decimal places.", f(gmid, 0), f(svt, 1), f(sb, 1))
		return num(q, f(s, 2), edakit.Sprintf("The threshold term is $%s \\cdot %s \\times 10^{-3} = %s$%%. Independent terms add in variance: $\\sqrt{%s^2 + %s^2} = %s$%%.", f(gmid, 0), f(svt, 1), f(vt, 2), f(vt, 2), f(sb, 1), f(s, 2)), nil)
	default:
		s1 := unif(0.5, 4, 0.1)
		s3 := unif(0.5, 4, 0.1)
		g1 := float64(50 * edakit.Int(2, 20))
		g3 := float64(50 * edakit.Int(1, 16))
		v := math.Sqrt(s1*s1 + (g3/g1)*(g3/g1)*s3*s3)
		q := edakit.Sprintf("In a five-transistor OTA the input pair has $\\sigma(\\Delta V_{T1}) = %s$ mV and $g_{m1} = %s\\ \\mu$S; the mirror load has $\\sigma(\\Delta V_{T3}) = %s$ mV and $g_{m3} = %s\\ \\mu$S. What is the input-referred offset sigma, in mV? Round to 2 decimal places.", f(s1, 1), f(g1, 0), f(s3, 1), f(g3, 0))
		return num(q, f(v, 2), edakit.Sprintf("Load mismatch produces a current error $g_{m3}\\Delta V_{T3}$, which the input pair cancels with $\\Delta V_{T3} g_{m3}/g_{m1}$ at its gate. $\\sigma(V_{os}) = \\sqrt{%s^2 + (%s/%s)^2 %s^2} = %s$ mV.", f(s1, 1), f(g3, 0), f(g1, 0), f(s3, 1), f(v, 2)), nil)
	}
}

// ccPattern draws a two-row arrangement of unit devices A and B with equal counts. Half the time
// it is point-symmetric (each device's centroid at the array centre); otherwise it is shuffled.
func ccPattern() [2][]byte {
	k := edakit.Int(2, 4)
	n := 2 * k
	var g [2][]byte
	if edakit.Int(0, 1) == 0 {
		row := make([]byte, 0, n)
		for i := 0; i < k; i++ {
			row = append(row, 'A', 'B')
		}
		row = edakit.Shuffle(row)
		g[0] = row
		g[1] = make([]byte, n)
		for j := 0; j < n; j++ {
			g[1][j] = row[n-1-j]
		}
		return g
	}
	cells := make([]byte, 0, 2*n)
	for i := 0; i < n; i++ {
		cells = append(cells, 'A', 'B')
	}
	cells = edakit.Shuffle(cells)
	g[0] = cells[:n]
	g[1] = cells[n:]
	return g
}

func centroid(g [2][]byte, d byte) (float64, float64) {
	var sx, sy, c float64
	for r := 0; r < 2; r++ {
		for j, x := range g[r] {
			if x == d {
				sx += float64(j + 1)
				sy += float64(r + 1)
				c++
			}
		}
	}
	return sx / c, sy / c
}

func ccText(g [2][]byte) string {
	return edakit.Sprintf("row 1 is %s and row 2 is %s", string(g[0]), string(g[1]))
}

func commonCentroid(ctx generator.GeneratorContext) generator.Problem {
	g := ccPattern()
	ax, ay := centroid(g, 'A')
	bx, by := centroid(g, 'B')
	setup := "Two matched devices A and B are each split into equal unit cells in a two-row array: " + ccText(g) + ". Columns are numbered from 1 at the left and rows from 1 at the top, one pitch apart."
	switch whichKP(ctx) {
	case 1:
		d := pick[byte]('A', 'B')
		x, _ := centroid(g, d)
		q := setup + edakit.Sprintf(" What is the column coordinate of device %s's centroid? Round to 3 decimal places.", string(d))
		return num(q, f(x, 3), edakit.Sprintf("Average the column numbers of %s's cells: %s. (A is at column %s, B at %s.)", string(d), f(x, 3), f(ax, 3), f(bx, 3)), nil)
	case 2:
		dist := math.Hypot(ax-bx, ay-by)
		q := setup + " What is the distance between the centroids of A and B, in pitches? Answer 0 if the arrangement is common-centroid. Round to 3 decimal places."
		return num(q, f(dist, 3), edakit.Sprintf("A's centroid is (%s, %s) and B's is (%s, %s) in (column, row). The distance is %s pitches; only a distance of 0 is common-centroid.", f(ax, 3), f(ay, 3), f(bx, 3), f(by, 3), f(dist, 3)), nil)
	default:
		gx := unif(0.05, 1, 0.05)
		gy := unif(0.05, 1, 0.05)
		m := math.Abs(gx*(ax-bx) + gy*(ay-by))
		q := setup + edakit.Sprintf(" Threshold voltage has a linear gradient of %s mV per column and %s mV per row, increasing to the right and downward. What is the magnitude of the resulting mean threshold difference between A and B, in mV? Round to 3 decimal places.", f(gx, 2), f(gy, 2))
		return num(q, f(m, 3), edakit.Sprintf("Under a linear gradient each device's mean shift is the gradient evaluated at its centroid. $|%s(%s - %s) + %s(%s - %s)| = %s$ mV.", f(gx, 2), f(ax, 3), f(bx, 3), f(gy, 2), f(ay, 3), f(by, 3), f(m, 3)), nil)
	}
}

func row1D() string {
	pats := []string{"AABB", "ABAB", "ABBA", "AAABBB", "ABABAB", "ABBAAB", "AABBBBAA", "ABABBABA", "ABBAABBA", "AABBAABB", "ABABABAB", "AAAABBBB", "BAAB", "BABA", "ABBABAAB", "AABABB"}
	return pick(pats...)
}

func pos1D(p string, d byte) []float64 {
	n := len(p)
	c := float64(n+1) / 2
	var xs []float64
	for i := 0; i < n; i++ {
		if p[i] == d {
			xs = append(xs, float64(i+1)-c)
		}
	}
	return xs
}

func mean(xs []float64) float64 {
	s := 0.0
	for _, x := range xs {
		s += x
	}
	return s / float64(len(xs))
}

func interdig(ctx generator.GeneratorContext) generator.Problem {
	k := whichKP(ctx)
	p := row1D()
	if k == 3 {
		p = pick("ABBA", "BAAB", "AABBBBAA", "ABBAABBA", "BAABBAAB", "ABABBABA", "AABBAABB", "BBAAAABB", "ABBBBAAA", "AAABBABB")
	}
	a, b := pos1D(p, 'A'), pos1D(p, 'B')
	dc := math.Abs(mean(a) - mean(b))
	row := "A pair of matched transistors A and B is laid out as one row of equal fingers, left to right: " + strings.Join(strings.Split(p, ""), " ") + ", with dummy fingers at both ends."
	switch k {
	case 1:
		q := row + " What is the distance between the centroids of A and B, in finger pitches? Round to 3 decimal places."
		return num(q, f(dc, 3), edakit.Sprintf("Measure positions from the row centre. A's mean position is %s and B's is %s, so the centroid separation is %s finger pitches.", f(mean(a), 3), f(mean(b), 3), f(dc, 3)), nil)
	case 2:
		g := unif(0.05, 2, 0.05)
		m := g * dc
		q := row + edakit.Sprintf(" Oxide thickness produces a linear threshold gradient of %s mV per finger pitch along the row. What is the systematic threshold mismatch between A and B, in mV? Round to 3 decimal places.", f(g, 2))
		return num(q, f(m, 3), edakit.Sprintf("A linear gradient shifts each device by the gradient at its centroid, so the mismatch is $%s \\cdot %s = %s$ mV.", f(g, 2), f(dc, 3), f(m, 3)), nil)
	default:
		c := unif(0.01, 0.2, 0.01)
		var sa, sb float64
		for _, x := range a {
			sa += x * x
		}
		for _, x := range b {
			sb += x * x
		}
		m := c * math.Abs(sa/float64(len(a))-sb/float64(len(b)))
		q := row + edakit.Sprintf(" Suppose the threshold varies quadratically along the row, $\\Delta V_T(x) = %s x^2$ mV, where $x$ is the finger's offset from the row centre in pitches. What is the systematic mismatch between A's and B's mean thresholds, in mV? Round to 3 decimal places.", f(c, 2))
		return num(q, f(m, 3), edakit.Sprintf("Average $x^2$ over each device's fingers: %s for A and %s for B. The mismatch is $%s \\cdot |%s - %s| = %s$ mV. A pattern that cancels a linear gradient need not cancel a quadratic one.", f(sa/float64(len(a)), 3), f(sb/float64(len(b)), 3), f(c, 2), f(sa/float64(len(a)), 3), f(sb/float64(len(b)), 3), f(m, 3)), nil)
	}
}

func parasitics(ctx generator.GeneratorContext) generator.Problem {
	gm6 := float64(100 * edakit.Int(5, 40))
	cl := unif(0.5, 5, 0.1)
	cp := unif(0.1, 2, 0.05)
	gbw := unif(5, 60, 1)
	switch whichKP(ctx) {
	case 1:
		p2 := gm6 / (2 * math.Pi * (cl + cp))
		q := edakit.Sprintf("A two-stage op-amp's non-dominant pole is $g_{m6}/(2\\pi C_{out})$ with $g_{m6} = %s\\ \\mu$S. The schematic load is %s pF, and extraction adds %s pF of wiring and junction capacitance at the output. What is the post-layout pole frequency in MHz? Round to 2 decimal places.", f(gm6, 0), f(cl, 1), f(cp, 2))
		return num(q, f(p2, 2), edakit.Sprintf("$C_{out} = %s + %s = %s$ pF, so $f_{p2} = %s/(2\\pi \\cdot %s) = %s$ MHz, down from %s MHz before layout.", f(cl, 1), f(cp, 2), f(cl+cp, 2), f(gm6, 0), f(cl+cp, 2), f(p2, 2), f(gm6/(2*math.Pi*cl), 2)), nil)
	case 2:
		p2 := gm6 / (2 * math.Pi * (cl + cp))
		pm := 90 - deg(math.Atan(gbw/p2))
		q := edakit.Sprintf("An op-amp crosses over at GBW = %s MHz with a dominant pole contributing $-90$ degrees. Its non-dominant pole is $g_{m6}/(2\\pi C_{out})$, $g_{m6} = %s\\ \\mu$S, with %s pF of schematic load plus %s pF of extracted parasitics. What is the post-layout phase margin in degrees? Round to 1 decimal place.", f(gbw, 0), f(gm6, 0), f(cl, 1), f(cp, 2))
		return num(q, f(pm, 1), edakit.Sprintf("$f_{p2} = %s$ MHz, so $\\mathrm{PM} = 90 - \\arctan(%s/%s) = %s$ degrees (schematic: %s degrees).", f(p2, 3), f(gbw, 0), f(p2, 3), f(pm, 1), f(90-deg(math.Atan(gbw*2*math.Pi*cl/gm6)), 1)), nil)
	default:
		pmT := float64(5 * edakit.Int(11, 13))
		cpm := unif(0.1, 0.3, 0.01)
		gm := gm6
		// Need p2 >= gbw / tan(90 - pm); cmax = gm/(2 pi p2).
		p2 := gbw / math.Tan((90-pmT)*math.Pi/180)
		cmax := gm / (2 * math.Pi * p2)
		for cmax <= cl+0.1 {
			gm += 500
			cmax = gm / (2 * math.Pi * p2)
		}
		l := (cmax - cl) * 1000 / cpm
		q := edakit.Sprintf("An op-amp has GBW = %s MHz, a dominant pole contributing $-90$ degrees, and a non-dominant pole $g_{m6}/(2\\pi C_{out})$ with $g_{m6} = %s\\ \\mu$S and %s pF of device load. The output route adds %s fF per micrometre. What is the longest output route, in micrometres, that keeps the phase margin at least %s degrees? Round down to a whole micrometre.", f(gbw, 0), f(gm, 0), f(cl, 1), f(cpm, 2), f(pmT, 0))
		return num(q, f(math.Floor(l), 0), edakit.Sprintf("The margin needs $f_{p2} \\ge %s/\\tan(%s) = %s$ MHz, so $C_{out} \\le %s/(2\\pi \\cdot %s) = %s$ pF. That leaves %s pF for wiring, or %s micrometres, rounded down to %s. This is the kind of constraint a router must honour.", f(gbw, 0), f(90-pmT, 0), f(p2, 3), f(gm, 0), f(p2, 3), f(cmax, 4), f(cmax-cl, 4), f(l, 1), f(math.Floor(l), 0)), nil)
	}
}
