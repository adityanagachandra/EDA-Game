package route

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// costPluginCalc evaluates a Tier 2 cost plug-in on three candidate paths.
func costPluginCalc(ctx generator.GeneratorContext) generator.Problem {
	for try := 0; ; try++ {
		alpha := edakit.Pick([]float64{0.5, 1, 2})
		beta := edakit.Pick([]float64{0, 0.5, 1})
		gamma := edakit.Pick([]float64{0, 2, 4})
		type edge struct {
			b     float64
			d, c  int
			h     float64
			crit  bool
			isVia bool
		}
		cost := func(e edge) float64 {
			ov := math.Max(0, float64(e.d+1-e.c))
			v := e.b*(1+alpha*ov*ov) + beta*e.h
			if e.crit && e.isVia {
				v += gamma
			}
			return v
		}
		names := []string{"A", "B", "C"}
		var parts []string
		var totals []float64
		for range names {
			n := edakit.Int(2, 4)
			tot := 0.0
			var es []string
			for i := 0; i < n; i++ {
				e := edge{b: 1, c: edakit.Int(2, 4), h: float64(edakit.Int(0, 3))}
				e.d = e.c + edakit.Int(-2, 1)
				e.isVia = edakit.Int(0, 3) == 0
				e.crit = true
				tot += cost(e)
				kind := "wire"
				if e.isVia {
					kind = "via"
				}
				es = append(es, edakit.Sprintf("%s (demand %d, cap %d, history %s)", kind, e.d, e.c, edakit.Num(e.h, 0)))
			}
			totals = append(totals, tot)
			parts = append(parts, strings.Join(es, ", "))
		}
		bi := 0
		for i := range totals {
			if totals[i] < totals[bi] {
				bi = i
			}
		}
		tie := false
		for i := range totals {
			if i != bi && math.Abs(totals[i]-totals[bi]) < 1e-9 {
				tie = true
			}
		}
		if tie && try < generator.MaxDrawAttempts {
			continue
		}
		var desc []string
		for i, p := range parts {
			desc = append(desc, edakit.Sprintf("path %s: %s", names[i], p))
		}
		q := edakit.Sprintf("A team registers a Tier 2 edge-cost plug-in in a router's cost seam: $cost(e) = b_e(1 + %s\\max(0,\\ d_e + 1 - c_e)^2) + %s h_e$, plus %s for every via on a timing-critical net. The router core calls it for each edge it expands and keeps its own legality checks. All edges have $b_e = 1$, and the net being routed is timing-critical. Candidate paths: %s. Which path does the router choose?",
			edakit.Num(alpha, 1), edakit.Num(beta, 1), edakit.Num(gamma, 0), strings.Join(desc, "; "))
		exp := edakit.Sprintf("Plug-in costs: A = %s, B = %s, C = %s. The core takes the cheapest, path %s. The plug-in changes only prices; it cannot make an illegal path legal, which stays the core's job.", edakit.Num(totals[0], 2), edakit.Num(totals[1], 2), edakit.Num(totals[2], 2), names[bi])
		var d []string
		for i, n := range names {
			if math.Abs(totals[i]-totals[bi]) > 1e-9 {
				d = append(d, "Path "+n)
			}
		}
		d = append(d, "No path: the plug-in must first reject the overflowed edges")
		return edakit.Choice(q, "Path "+names[bi], exp, d...)
	}
}

// drcGolden: wirelength delta (KP 1), unintended changed nets (KP 2), failing golden designs
// (KP 3).
func drcGolden(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		base := float64(edakit.Int(200000, 900000))
		delta := float64(edakit.Int(-3000, 3000)) / 1000
		nw := math.Round(base * (1 + delta/100))
		pct := (nw - base) / base * 100
		ans := edakit.Num(pct, 2)
		q := edakit.Sprintf("A router change is checked on a golden reference design. The baseline routing has %s um of wirelength; the new build routes it DRC-clean with %s um. What is the wirelength change, in percent of the baseline? Round to 2 decimal places (negative means shorter).", edakit.Num(base, 0), edakit.Num(nw, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("(%s - %s) / %s x 100 = %s%%. A QoR gate compares this against a tolerance only after the hard gates, zero DRC violations and zero opens, have passed.", edakit.Num(nw, 0), edakit.Num(base, 0), edakit.Num(base, 0), ans),
			map[string]string{"delta_pct": ans})
	case 2:
		n := edakit.Int(8, 12)
		var parts []string
		bad := 0
		for i := 0; i < n; i++ {
			clock := edakit.Int(0, 2) == 0
			b := edakit.Int(50, 400)
			nw := b
			if clock || edakit.Int(0, 3) == 0 {
				nw = b + edakit.Int(-20, 20)
			}
			if !clock && nw != b {
				bad++
			}
			kind := "signal"
			if clock {
				kind = "clock"
			}
			parts = append(parts, edakit.Sprintf("n%d (%s) %d to %d", i+1, kind, b, nw))
		}
		q := "A Tier 2 plug-in is meant to change routing of clock nets only. A differential test routes a reference design with and without it and lists per-net wirelength in um, baseline to new: " + strings.Join(parts, "; ") + ". How many nets changed that the plug-in was not supposed to touch?"
		return edakit.Numeric(q, edakit.IntStr(bad),
			edakit.Sprintf("Count signal nets whose wirelength moved: %d. Clock nets are allowed to change. Any unintended change fails the differential check even if total wirelength improved, because it shows the seam leaks behaviour outside its scope (or the router is not deterministic).", bad),
			map[string]string{"unintended": edakit.IntStr(bad)})
	default:
		tol := edakit.Pick([]float64{0.5, 1, 2})
		n := 5
		var parts []string
		fails := 0
		for i := 0; i < n; i++ {
			drc := 0
			if edakit.Int(0, 3) == 0 {
				drc = edakit.Int(1, 12)
			}
			opens := 0
			if edakit.Int(0, 5) == 0 {
				opens = edakit.Int(1, 3)
			}
			dw := float64(edakit.Int(-30, 30)) / 10
			if drc > 0 || opens > 0 || math.Abs(dw) > tol {
				fails++
			}
			parts = append(parts, edakit.Sprintf("design %d: %d DRC violations, %d opens, wirelength %s%%", i+1, drc, opens, signed(dw)))
		}
		q := edakit.Sprintf("A router's golden suite passes a design only if it routes with zero DRC violations, zero opens, and wirelength within %s%% of the stored baseline in either direction. Results: %s. How many designs fail?", edakit.Num(tol, 1), strings.Join(parts, "; "))
		return edakit.Numeric(q, edakit.IntStr(fails),
			edakit.Sprintf("Check the hard gates first, then the tolerance: %d designs fail. A change is promoted only when every golden design passes.", fails),
			map[string]string{"fails": edakit.IntStr(fails)})
	}
}

func signed(x float64) string {
	if x > 0 {
		return "+" + edakit.Num(x, 1)
	}
	return edakit.Num(x, 1)
}

type bbox struct{ x1, y1, x2, y2 int }

// fastCongestion: RUDY of one net (KP 1), RUDY summed at a tile (KP 2), estimator misses against
// the accurate router (KP 3).
func fastCongestion(ctx generator.GeneratorContext) generator.Problem {
	pitch := edakit.Pick([]float64{0.1, 0.14, 0.2})
	def := edakit.Sprintf("RUDY estimates routing demand by spreading each net's wire area uniformly over its bounding box: for a box W by H um and wire pitch %s um, the density inside the box is $p(W + H)/(WH)$, the fraction of the area the net's half-perimeter wire would occupy. ", edakit.Num(pitch, 2))
	switch kp(ctx) {
	case 1:
		W, H := float64(edakit.Int(2, 40)), float64(edakit.Int(2, 40))
		v := pitch * (W + H) / (W * H)
		ans := edakit.Num(v, 3)
		return edakit.Numeric(def+edakit.Sprintf("What is the RUDY density inside the box of a net spanning %s um by %s um? Round to 3 decimal places.", edakit.Num(W, 0), edakit.Num(H, 0)), ans,
			edakit.Sprintf("%s x (%s + %s) / (%s x %s) = %s. Small boxes are dense; long, thin nets spread thin.", edakit.Num(pitch, 2), edakit.Num(W, 0), edakit.Num(H, 0), edakit.Num(W, 0), edakit.Num(H, 0), ans),
			map[string]string{"rudy": ans})
	case 2:
		n := edakit.Scaled(ctx, 4, 6)
		tx, ty := edakit.Int(5, 15), edakit.Int(5, 15)
		var parts []string
		sum := 0.0
		hit := 0
		for i := 0; i < n; i++ {
			b := bbox{edakit.Int(0, 15), edakit.Int(0, 15), 0, 0}
			b.x2, b.y2 = b.x1+edakit.Int(2, 10), b.y1+edakit.Int(2, 10)
			if i < 2 {
				b = bbox{tx - edakit.Int(1, 5), ty - edakit.Int(1, 5), tx + edakit.Int(2, 6), ty + edakit.Int(2, 6)}
			}
			parts = append(parts, edakit.Sprintf("net %d box x %d to %d, y %d to %d", i+1, b.x1, b.x2, b.y1, b.y2))
			if b.x1 <= tx && tx+1 <= b.x2 && b.y1 <= ty && ty+1 <= b.y2 {
				W, H := float64(b.x2-b.x1), float64(b.y2-b.y1)
				sum += pitch * (W + H) / (W * H)
				hit++
			}
		}
		ans := edakit.Num(sum, 3)
		q := def + edakit.Sprintf("Net bounding boxes, in um: %s. What is the RUDY density of the 1 by 1 um tile from (%d, %d) to (%d, %d), summing every net whose box contains the whole tile? Round to 3 decimal places.", strings.Join(parts, "; "), tx, ty, tx+1, ty+1)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("%d of the boxes contain the tile; adding their densities gives %s. RUDY needs no routing at all, so a placer can evaluate it inside its own loop.", hit, ans),
			map[string]string{"nets": edakit.IntStr(hit), "rudy": ans})
	default:
		n := edakit.Int(8, 12)
		thr := edakit.Pick([]float64{0.8, 0.9, 1.0})
		var parts []string
		miss := 0
		for i := 0; i < n; i++ {
			est := float64(edakit.Int(60, 120)) / 100
			acc := est + float64(edakit.Int(-20, 45))/100
			if est < thr && acc > 1.0 {
				miss++
			}
			parts = append(parts, edakit.Sprintf("t%d %s / %s", i+1, edakit.Num(est, 2), edakit.Num(acc, 2)))
		}
		q := edakit.Sprintf("An agent's modeling loop flags a tile for the accurate global router only if the fast estimator's utilisation is at least %s. Later, the accurate router reports utilisation per tile; above 1.00 the tile overflows. Fast and accurate values: %s. How many overflowing tiles did the fast estimator fail to flag?", edakit.Num(thr, 2), strings.Join(parts, "; "))
		return edakit.Numeric(q, edakit.IntStr(miss),
			edakit.Sprintf("Count tiles with fast value below %s and accurate value above 1.00: %d. These misses set how low the threshold must sit, and when the loop must recurse into the accurate engine instead of trusting the fast one.", edakit.Num(thr, 2), miss),
			map[string]string{"misses": edakit.IntStr(miss)})
	}
}
