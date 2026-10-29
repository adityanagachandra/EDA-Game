package pv

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// maskPrep: shots for one rectangle (KP 1), fewest shots for an L-shape cut two ways (KP 2),
// vertices after OPC fragmentation (KP 3).
func maskPrep(ctx generator.GeneratorContext) generator.Problem {
	shots := func(w, h, m int) int { return ceilDiv(w, m) * ceilDiv(h, m) }
	switch kp(ctx) {
	case 1:
		m := edakit.Pick([]int{500, 1000, 2000})
		W, H := edakit.Scaled(ctx, 1500, 9000), edakit.Int(800, 6000)
		n := shots(W, H, m)
		q := edakit.Sprintf("A variable-shaped-beam mask writer exposes rectangles of at most %d by %d nm per shot. How many shots does it need for one %d by %d nm rectangle, split on a regular grid of maximal shots?", m, m, W, H)
		return edakit.Numeric(q, edakit.IntStr(n),
			edakit.Sprintf("ceil(%d / %d) = %d columns and ceil(%d / %d) = %d rows: %d shots. Write time grows with the shot count, not the area.", W, m, ceilDiv(W, m), H, m, ceilDiv(H, m), n),
			map[string]string{"shots": edakit.IntStr(n)})
	case 2:
		m := edakit.Pick([]int{500, 1000})
		X, Y := 100*edakit.Int(10, 60), 100*edakit.Int(10, 60)
		X1, Y1 := 100*edakit.Int(3, X/100-3), 100*edakit.Int(3, Y/100-3)
		hcut := shots(X, Y1, m) + shots(X1, Y-Y1, m)
		vcut := shots(X1, Y, m) + shots(X-X1, Y1, m)
		best := imin(hcut, vcut)
		q := edakit.Sprintf("An L-shaped polygon is the union of $[0, %d] \\times [0, %d]$ and $[0, %d] \\times [%d, %d]$ (nm). The fracturer cuts it into two rectangles, either with a horizontal cut at y = %d or a vertical cut at x = %d, and each rectangle is exposed with shots of at most %d by %d nm on a regular grid. What is the fewer shot count of the two cuts?", X, Y1, X1, Y1, Y, Y1, X1, m, m)
		return edakit.Numeric(q, edakit.IntStr(best),
			edakit.Sprintf("Horizontal cut: %d by %d and %d by %d give %d shots. Vertical cut: %d by %d and %d by %d give %d shots. The better cut needs %d; fracturing choices change write time even for the same polygon.", X, Y1, X1, Y-Y1, hcut, X1, Y, X-X1, Y1, vcut, best),
			map[string]string{"horizontal": edakit.IntStr(hcut), "vertical": edakit.IntStr(vcut)})
	default:
		f := edakit.Pick([]int{20, 40, 50, 80, 100})
		W, H := edakit.Scaled(ctx, 60, 1200), edakit.Int(40, 600)
		mw, mh := ceilDiv(W, f), ceilDiv(H, f)
		v := 4 + 2*(2*(mw-1)+2*(mh-1))
		q := edakit.Sprintf("OPC cuts every edge of a drawn %d by %d nm rectangle into fragments of at most %d nm (equal fragments, the fewest possible) and moves each fragment independently. If every pair of neighbouring fragments on a side ends up at a different offset, how many vertices does the corrected polygon have?", W, H, f)
		return edakit.Numeric(q, edakit.IntStr(v),
			edakit.Sprintf("The %d nm sides split into %d fragments and the %d nm sides into %d. Each internal break becomes a jog with 2 vertices, and the 4 corners stay: 4 + 2 x (2 x %d + 2 x %d) = %d. This vertex growth is why post-OPC mask data is far larger than the drawn layout.", W, mw, H, mh, mw-1, mh-1, v),
			map[string]string{"fragments_w": edakit.IntStr(mw), "fragments_h": edakit.IntStr(mh)})
	}
}

// overlappingPair draws two rectangles that overlap with positive area.
func overlappingPair(span, minS, maxS int) (rect, rect) {
	for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
		a, b := randRect(span, minS, maxS), randRect(span, minS, maxS)
		if interArea(a, b) > 0 && a != b {
			return a, b
		}
	}
	return rect{0, 0, minS + 2, minS + 2}, rect{1, 1, minS + 3, minS + 4}
}

// xorCheck: XOR of two rectangles (KP 1), of two shape sets (KP 2), of a shifted shape (KP 3).
func xorCheck(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		a, b := overlappingPair(30, 5, 20)
		x := a.area() + b.area() - 2*interArea(a, b)
		q := edakit.Sprintf("A tape-out XOR compares a shape in the signed-off layout, %s, with the same shape in the final stream, %s (um). What area, in square um, does the XOR report?", a, b)
		return edakit.Numeric(q, edakit.IntStr(x),
			edakit.Sprintf("XOR = A + B - 2 (A AND B) = %d + %d - 2 x %d = %d square um. Any nonzero XOR on a layer that should be unchanged stops the tape-out until it is explained.", a.area(), b.area(), interArea(a, b), x),
			map[string]string{"and": edakit.IntStr(interArea(a, b))})
	case 2:
		n := edakit.Scaled(ctx, 2, 3)
		var A, B []rect
		for i := 0; i < n; i++ {
			A = append(A, randRect(30, 4, 14))
		}
		for _, r := range A {
			if edakit.Int(0, 2) == 0 {
				B = append(B, r)
			} else {
				B = append(B, r.shift(edakit.Int(-3, 3), edakit.Int(-3, 3)))
			}
		}
		if edakit.Int(0, 1) == 0 {
			B = append(B, randRect(30, 3, 8))
		}
		x := xorArea(A, B)
		q := edakit.Sprintf("Version 1 of a layer is %s and version 2 is %s (um). Shapes within a version may overlap and are merged first. What is the XOR area between the two versions, in square um?", listRects("A", A), listRects("B", B))
		return edakit.Numeric(q, edakit.IntStr(x),
			edakit.Sprintf("Merge each version, then sum the area covered by exactly one of the two unions: %d square um. A shape moved by a few units shows up as two thin slivers, not as a whole shape.", x),
			map[string]string{"xor": edakit.IntStr(x)})
	default:
		w, h := edakit.Int(10, 200), edakit.Int(10, 200)
		dx, dy := edakit.Int(-9, 9), edakit.Int(-9, 9)
		if dx == 0 && dy == 0 {
			dx = 1
		}
		x := 2 * (w*h - (w-iabs(dx))*(h-iabs(dy)))
		q := edakit.Sprintf("An ECO accidentally shifts a %d by %d nm rectangle by (%d, %d) nm. What XOR area, in square nm, does the comparison against the previous version report for that shape?", w, h, dx, dy)
		return edakit.Numeric(q, edakit.IntStr(x),
			edakit.Sprintf("The two copies overlap in a (%d - %d) by (%d - %d) box, so XOR = 2 x (%d - %d) = %d. A uniform shift of a whole cell gives the same signature on every shape, which helps a reviewer recognize it.", w, iabs(dx), h, iabs(dy), w*h, (w-iabs(dx))*(h-iabs(dy)), x),
			map[string]string{"overlap": edakit.IntStr((w - iabs(dx)) * (h - iabs(dy)))})
	}
}

// fixRouting: how many violations an agent may send back to the router.
func fixRouting(ctx generator.GeneratorContext) generator.Problem {
	sp, sh, ma, en := edakit.Scaled(ctx, 3, 80), edakit.Int(0, 10), edakit.Int(0, 30), edakit.Int(0, 20)
	ant, den := edakit.Int(0, 15), edakit.Int(0, 6)
	inCell := edakit.Int(0, sp/3)
	router := sp - inCell + sh + ma + en
	q := edakit.Sprintf("A signoff DRC report lists %d metal spacing, %d short, %d minimum-area and %d via-enclosure violations on routed metal, %d antenna violations and %d density-window violations. Of the spacing violations, %d lie entirely inside standard-cell layouts. The agent's policy: violations on routed metal go back to the detailed router as a Tier 1 rerun of the affected region; antenna violations go to the antenna-repair step; density goes to the fill step; anything inside a cell layout is escalated to a human. How many violations does the agent send to the router?", sp, sh, ma, en, ant, den, inCell)
	why := edakit.Sprintf("Routed-metal violations: (%d - %d) + %d + %d + %d = %d. The %d in-cell spacing violations are not the router's to fix: the router treats cell layouts as fixed, so changing them is a core library change.", sp, inCell, sh, ma, en, router, inCell)
	return numChoice(q, float64(router), 0, "", why, float64(sp+sh+ma+en), float64(router+ant), float64(sp+sh+ma+en+ant+den))
}

type marker struct {
	rule string
	x, y int
}

var ruleNames = []string{"M1.S.1", "M2.S.1", "M1.W.1", "V1.EN.1", "M2.A.1", "M3.S.2"}

// goldenClean: new markers (KP 1), fixed markers (KP 2), unmatched markers under a position
// tolerance (KP 3).
func goldenClean(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	n := edakit.Scaled(ctx, 4, 8)
	var base, next []marker
	for i := 0; i < n; i++ {
		base = append(base, marker{edakit.Pick(ruleNames), 100 * (i + 1), 100 * edakit.Int(1, 9)})
	}
	tol := 0
	if k == 3 {
		tol = edakit.Int(2, 5)
	}
	newCount := 0
	for _, m := range base {
		switch edakit.Int(0, 3) {
		case 0: // fixed: absent from the new run
		case 1: // same rule, different place far away: a new marker
			next = append(next, marker{m.rule, m.x + 50, m.y + 50})
			newCount++
		default:
			dx, dy := 0, 0
			if tol > 0 {
				dx = edakit.Int(-tol, tol)
				dy = edakit.Int(-(tol - iabs(dx)), tol-iabs(dx))
				if edakit.Int(0, 3) == 0 {
					dx = tol + edakit.Int(1, 3)
					newCount++
				}
			}
			next = append(next, marker{m.rule, m.x + dx, m.y + dy})
		}
	}
	extra := edakit.Int(0, 2)
	for i := 0; i < extra; i++ {
		next = append(next, marker{edakit.Pick(ruleNames), 100*(n+i+1) + 30, 100 * edakit.Int(1, 9)})
		newCount++
	}
	desc := func(ms []marker) string {
		if len(ms) == 0 {
			return "none"
		}
		var p []string
		for _, m := range ms {
			p = append(p, edakit.Sprintf("%s at (%d, %d)", m.rule, m.x, m.y))
		}
		return strings.Join(p, ", ")
	}
	next = edakit.Shuffle(next)
	stem := edakit.Sprintf("An agent's change is checked by a differential DRC run. The baseline markers are %s. The new run's markers are %s.", desc(base), desc(next))
	match := func(a, b marker) bool {
		return a.rule == b.rule && iabs(a.x-b.x)+iabs(a.y-b.y) <= tol
	}
	switch k {
	case 1, 3:
		cnt := 0
		for _, m := range next {
			found := false
			for _, b := range base {
				if match(m, b) {
					found = true
				}
			}
			if !found {
				cnt++
			}
		}
		_ = newCount
		if k == 1 {
			q := stem + " A marker matches only if the rule and the position are identical. How many markers did the change introduce?"
			return edakit.Numeric(q, edakit.IntStr(cnt),
				edakit.Sprintf("Introduced = new markers with no identical baseline marker: %d. A golden test for a fix requires this to be zero, whatever happened to the old markers.", cnt),
				map[string]string{"introduced": edakit.IntStr(cnt)})
		}
		q := stem + edakit.Sprintf(" Markers match if the rule is the same and the Manhattan distance between positions is at most %d nm, since a rerun can shift a marker slightly. How many new-run markers have no match in the baseline?", tol)
		return edakit.Numeric(q, edakit.IntStr(cnt),
			edakit.Sprintf("Match each new marker to a baseline marker of the same rule within %d nm: %d remain unmatched and count as introduced. Matching with tolerance stops jitter from looking like regressions, but too large a tolerance hides real ones.", tol, cnt),
			map[string]string{"unmatched": edakit.IntStr(cnt)})
	default:
		cnt := 0
		for _, b := range base {
			found := false
			for _, m := range next {
				if match(m, b) {
					found = true
				}
			}
			if !found {
				cnt++
			}
		}
		q := stem + " A marker matches only if the rule and the position are identical. How many baseline markers did the change fix (they no longer appear)?"
		return edakit.Numeric(q, edakit.IntStr(cnt),
			edakit.Sprintf("Fixed = baseline markers with no identical marker in the new run: %d. A marker of the same rule elsewhere is a new violation, not the old one moved.", cnt),
			map[string]string{"fixed": edakit.IntStr(cnt)})
	}
}

// fastHotspot: speedup of screen-then-verify (KP 1), escapes at a score threshold (KP 2),
// real hotspots behind a proxy count (KP 3).
func fastHotspot(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		N := 1000 * edakit.Scaled(ctx, 10, 500)
		tm := edakit.Pick([]float64{0.5, 1, 2, 5})
		ts := float64(100 * edakit.Int(5, 60))
		f := float64(edakit.Int(1, 15)) / 100
		k := int(float64(N) * f)
		sp := float64(N) * ts / (float64(N)*tm + float64(k)*ts)
		ans := edakit.Num(sp, 1)
		q := edakit.Sprintf("A full lithography simulation takes %s ms per clip and a learned hotspot predictor %s ms per clip. A layer has %d clips. The fast flow runs the predictor on every clip and simulates only the %d clips it flags. What is the speedup of the fast flow over simulating every clip? Round to 1 decimal place.", edakit.Num(ts, 0), edakit.Num(tm, 1), N, k)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Full: %d x %s ms. Fast: %d x %s + %d x %s ms. Ratio %s. The speedup is bounded by the flag rate: if the predictor flags more, the accurate engine dominates again.", N, edakit.Num(ts, 0), N, edakit.Num(tm, 1), k, edakit.Num(ts, 0), ans),
			map[string]string{"flagged": edakit.IntStr(k)})
	case 2:
		n := edakit.Scaled(ctx, 8, 14)
		t := float64(edakit.Int(3, 7)) / 10
		var parts []string
		miss, hot := 0, 0
		for i := 0; i < n; i++ {
			s := float64(edakit.Int(1, 99)) / 100
			h := edakit.Int(0, 2) == 0
			lab := "clean"
			if h {
				lab = "hotspot"
				hot++
				if s < t {
					miss++
				}
			}
			parts = append(parts, edakit.Sprintf("%s (%s)", edakit.Num(s, 2), lab))
		}
		if hot == 0 {
			parts[0] = edakit.Sprintf("%s (hotspot)", edakit.Num(t/2, 2))
			miss++
		}
		q := edakit.Sprintf("A predictor scores clips, and an agent sends to signoff simulation only clips scoring at least %s. Scores with true labels from a later full run are: %s. How many true hotspots would escape if the agent trusted the screen and skipped the full run?", edakit.Num(t, 1), strings.Join(parts, ", "))
		return edakit.Numeric(q, edakit.IntStr(miss),
			edakit.Sprintf("Count hotspots scoring below %s: %d. These are the false negatives; only a final accurate check catches them, which is why the fast engine screens and the signoff engine decides.", edakit.Num(t, 1), miss),
			map[string]string{"escapes": edakit.IntStr(miss)})
	default:
		P := edakit.Scaled(ctx, 10, 400)
		prec := float64(edakit.Int(50, 95)) / 100
		rec := float64(edakit.Int(40, 95)) / 100
		real := float64(P) * prec / rec
		ans := edakit.Num(real, 0)
		q := edakit.Sprintf("An agent is rewarded for lowering the hotspot count reported by a fast predictor. After its edits the predictor reports %d hotspots. On layouts like this one, the predictor's precision is %s and its recall %s. How many real hotspots should the agent expect the layout to contain? Round to the nearest whole number.", P, edakit.Num(prec, 2), edakit.Num(rec, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("True positives = %d x %s = %s, and these are a fraction %s of the real hotspots, so real = %s / %s = %s. A reward on the proxy count can be lowered by edits that move layouts into the predictor's blind spots; the reward must be checked against signoff.", P, edakit.Num(prec, 2), edakit.Num(float64(P)*prec, 1), edakit.Num(rec, 2), edakit.Num(float64(P)*prec, 1), edakit.Num(rec, 2), ans),
			map[string]string{"tp": edakit.Num(float64(P)*prec, 1)})
	}
}
