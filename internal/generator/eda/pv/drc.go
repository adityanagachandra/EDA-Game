package pv

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// widthSpacing: width violations (KP 1), Euclidean spacing violations (KP 2), width-dependent
// spacing (KP 3).
func widthSpacing(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		W := 10 * edakit.Int(4, 10)
		n := edakit.Scaled(ctx, 4, 9)
		var rs []rect
		bad := 0
		for i := 0; i < n; i++ {
			wd := imax(10, W+edakit.Int(-25, 40))
			ln := edakit.Int(3*W, 12*W)
			x := edakit.Int(0, 100)
			var r rect
			if edakit.Int(0, 1) == 0 {
				r = rect{x, 400 * i, x + ln, 400*i + wd}
			} else {
				r = rect{x, 400 * i, x + wd, 400*i + ln}
			}
			if wd < W {
				bad++
			}
			rs = append(rs, r)
		}
		q := edakit.Sprintf("The minimum width rule on metal 1 is %d nm. The layer holds the separate wires %s (nm). How many wires violate the width rule?", W, listRects("W", rs))
		return edakit.Numeric(q, edakit.IntStr(bad),
			edakit.Sprintf("A rectangle's width is its shorter side. %d of the wires have a shorter side below %d nm; their length does not matter.", bad, W),
			map[string]string{"violations": edakit.IntStr(bad)})
	case 2:
		S := edakit.Int(18, 36)
		n := edakit.Scaled(ctx, 4, 8)
		rs := separated(n, 160, 10, 50)
		bad := 0
		for i := range rs {
			for j := i + 1; j < len(rs); j++ {
				if dist2(rs[i], rs[j]) < S*S {
					bad++
				}
			}
		}
		q := edakit.Sprintf("A spacing rule says any two separate shapes on the layer must be at least %d nm apart, measured as Euclidean distance (corner to corner when the shapes are diagonal). The shapes are %s (nm), and no two touch. How many pairs violate the rule?", S, listRects("R", rs))
		return edakit.Numeric(q, edakit.IntStr(bad),
			edakit.Sprintf("For each pair take the x gap gx and y gap gy (zero where projections overlap); the distance is sqrt(gx^2 + gy^2). %d of the %d pairs fall below %d nm.", bad, n*(n-1)/2, S),
			map[string]string{"violations": edakit.IntStr(bad), "pairs": edakit.IntStr(n * (n - 1) / 2)})
	default:
		s0, s1, wt := 40, 80, 200
		n := edakit.Scaled(ctx, 4, 8)
		x := 0
		var rs []rect
		var ws []int
		for i := 0; i < n; i++ {
			w := edakit.Int(80, 150)
			if edakit.Int(0, 2) == 0 {
				w = edakit.Int(200, 300)
			}
			if i > 0 {
				x += edakit.Int(30, 110)
			}
			rs = append(rs, rect{x, 0, x + w, 2000})
			ws = append(ws, w)
			x += w
		}
		bad := 0
		for i := 0; i+1 < n; i++ {
			need := s0
			if ws[i] >= wt || ws[i+1] >= wt {
				need = s1
			}
			if rs[i+1].x1-rs[i].x2 < need {
				bad++
			}
		}
		q := edakit.Sprintf("A width-dependent spacing rule needs %d nm between two wires, rising to %d nm when either wire is at least %d nm wide. Parallel vertical wires are %s (nm). How many pairs of wires violate the rule?", s0, s1, wt, listRects("M", rs))
		return edakit.Numeric(q, edakit.IntStr(bad),
			edakit.Sprintf("Only neighbours can be closer than %d nm, since every wire is at least 80 nm wide. Checking each neighbour gap against %d nm, or %d nm when a wide wire is involved, gives %d violations.", s1, s0, s1, bad),
			map[string]string{"violations": edakit.IntStr(bad)})
	}
}

// enclosureExtension: minimum via enclosure (KP 1), two-sided enclosure rule (KP 2), gate
// extension repair (KP 3).
func enclosureExtension(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		v := 10 * edakit.Int(3, 8)
		el, er, eb, et := edakit.Int(0, 40), edakit.Int(0, 40), edakit.Int(0, 40), edakit.Int(0, 40)
		vx, vy := edakit.Int(100, 400), edakit.Int(100, 400)
		via := rect{vx, vy, vx + v, vy + v}
		m := rect{vx - el, vy - eb, vx + v + er, vy + v + et}
		mn := imin(imin(el, er), imin(eb, et))
		q := edakit.Sprintf("A via cut is %s and the metal above it is %s (nm). What is the via's enclosure by the metal, the smallest distance from a via edge to the corresponding metal edge, in nm?", via, m)
		return edakit.Numeric(q, edakit.IntStr(mn),
			edakit.Sprintf("The four margins are left %d, right %d, bottom %d, top %d nm; the enclosure is the minimum, %d nm.", el, er, eb, et, mn),
			map[string]string{"left": edakit.IntStr(el), "right": edakit.IntStr(er), "bottom": edakit.IntStr(eb), "top": edakit.IntStr(et)})
	case 2:
		a, b := edakit.Int(0, 10), 5*edakit.Int(4, 8)
		n := edakit.Scaled(ctx, 4, 8)
		var parts []string
		bad := 0
		for i := 0; i < n; i++ {
			e := [4]int{edakit.Int(0, b+10), edakit.Int(0, b+10), edakit.Int(0, b+10), edakit.Int(0, b+10)}
			hx, vy := imin(e[0], e[1]), imin(e[2], e[3])
			ok := (hx >= b && vy >= a) || (hx >= a && vy >= b)
			if !ok {
				bad++
			}
			parts = append(parts, edakit.Sprintf("V%d (%d, %d, %d, %d)", i+1, e[0], e[1], e[2], e[3]))
		}
		q := edakit.Sprintf("A via enclosure rule requires at least %d nm on both sides along one axis and at least %d nm on both sides along the other axis (either orientation is allowed). The vias' enclosures as (left, right, bottom, top), in nm, are %s. How many vias violate the rule?", b, a, strings.Join(parts, ", "))
		return edakit.Numeric(q, edakit.IntStr(bad),
			edakit.Sprintf("For each via take min(left, right) and min(bottom, top). It passes if one of them is at least %d and the other at least %d. %d vias fail both orientations.", b, a, bad),
			map[string]string{"violations": edakit.IntStr(bad)})
	default:
		E := 10 * edakit.Int(5, 12)
		n := edakit.Scaled(ctx, 2, 6)
		var parts []string
		total := 0
		for i := 0; i < n; i++ {
			d1 := edakit.Int(0, 50) * 10
			d2 := d1 + 10*edakit.Int(15, 60)
			top, bot := edakit.Int(E-40, E+30), edakit.Int(E-40, E+30)
			top, bot = imax(0, top), imax(0, bot)
			total += imax(0, E-top) + imax(0, E-bot)
			parts = append(parts, edakit.Sprintf("G%d: diffusion y from %d to %d, poly y from %d to %d", i+1, d1, d2, d1-bot, d2+top))
		}
		q := edakit.Sprintf("Poly must extend at least %d nm past diffusion at both ends of every gate (the endcap rule). Vertical poly gates cross horizontal diffusion as follows (nm): %s. What total length of poly, in nm, must be added to fix every short endcap?", E, strings.Join(parts, "; "))
		return edakit.Numeric(q, edakit.IntStr(total),
			edakit.Sprintf("Each end's extension is the poly end minus the diffusion edge. Summing max(0, %d - extension) over both ends of every gate gives %d nm.", E, total),
			map[string]string{"added": edakit.IntStr(total)})
	}
}

// notchEOL: notch count (KP 1), end-of-line spacing (KP 2), parallel-run-length spacing (KP 3).
func notchEOL(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		N := 10 * edakit.Int(5, 9)
		n := edakit.Scaled(ctx, 4, 8)
		x := 0
		var prongs []string
		bad := 0
		for i := 0; i < n; i++ {
			if i > 0 {
				g := 10 * imax(2, edakit.Int(N/10-3, N/10+3))
				if g < N {
					bad++
				}
				x += g
			}
			w := 10 * edakit.Int(6, 15)
			prongs = append(prongs, edakit.Sprintf("[%d, %d]", x, x+w))
			x += w
		}
		q := edakit.Sprintf("A comb-shaped polygon is one horizontal spine with %d vertical prongs rising from it; the prongs span x-ranges %s (nm). The notch rule requires at least %d nm between two edges of the same polygon that face each other. How many notches violate it?", n, strings.Join(prongs, ", "), N)
		return edakit.Numeric(q, edakit.IntStr(bad),
			edakit.Sprintf("Each gap between neighbouring prongs is a notch: the edges belong to one polygon, so a merged-shape spacing check sees it. %d of the %d gaps are below %d nm.", bad, n-1, N),
			map[string]string{"violations": edakit.IntStr(bad)})
	case 2:
		S := 10 * edakit.Int(4, 6)
		Seol := S + 10*edakit.Int(2, 4)
		Weol := 10 * edakit.Int(6, 9)
		n := edakit.Scaled(ctx, 4, 8)
		var parts []string
		bad := 0
		for i := 0; i < n; i++ {
			w := 10 * edakit.Int(4, 12)
			d := 5 * edakit.Int(S/5-2, Seol/5+2)
			if (w < Weol && d < Seol) || d < S {
				bad++
			}
			parts = append(parts, edakit.Sprintf("E%d (width %d, gap %d)", i+1, w, d))
		}
		q := edakit.Sprintf("The normal spacing rule is %d nm. An end-of-line rule adds that a wire end narrower than %d nm needs at least %d nm to the shape in front of it. The wire ends, with the end's width and the gap to the facing shape in nm, are %s. How many wire ends violate a rule?", S, Weol, Seol, strings.Join(parts, ", "))
		return edakit.Numeric(q, edakit.IntStr(bad),
			edakit.Sprintf("A narrow end (width below %d) needs %d nm; a wide end only needs the normal %d nm. %d ends fall short.", Weol, Seol, S, bad),
			map[string]string{"violations": edakit.IntStr(bad)})
	default:
		S0 := 10 * edakit.Int(4, 6)
		S1 := S0 + 10*edakit.Int(2, 5)
		Wt := 10 * edakit.Int(15, 30)
		P := 100 * edakit.Int(3, 8)
		n := edakit.Scaled(ctx, 4, 8)
		var parts []string
		bad := 0
		for i := 0; i < n; i++ {
			w := 10 * edakit.Int(5, 40)
			prl := 50 * edakit.Int(2, 20)
			s := 5 * edakit.Int(S0/5-2, S1/5+2)
			need := S0
			if w >= Wt && prl > P {
				need = S1
			}
			if s < need {
				bad++
			}
			parts = append(parts, edakit.Sprintf("P%d (wider width %d, run %d, spacing %d)", i+1, w, prl, s))
		}
		q := edakit.Sprintf("A spacing table requires %d nm between parallel wires, but %d nm when the wider wire is at least %d nm wide and the two run in parallel for more than %d nm. The wire pairs are %s, all in nm. How many pairs violate the table?", S0, S1, Wt, P, strings.Join(parts, ", "))
		return edakit.Numeric(q, edakit.IntStr(bad),
			edakit.Sprintf("A pair needs %d nm only if both conditions hold (width at least %d and run above %d); otherwise %d nm. %d pairs fall short.", S1, Wt, P, S0, bad),
			map[string]string{"violations": edakit.IntStr(bad)})
	}
}

// areaVia: minimum-area extension (KP 1), via array count (KP 2), minimum-cut rule (KP 3).
func areaVia(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		w := 10 * edakit.Int(4, 10)
		amin := 1000 * edakit.Int(8, 30)
		l := edakit.Int(40, amin/w-10)
		ext := ceilDiv(amin, w) - l
		q := edakit.Sprintf("An isolated metal stub is %d nm wide and %d nm long. The minimum area rule is %d square nm. By how many nm must the stub be lengthened, keeping its width, to meet the rule? Give the smallest whole number of nm.", w, l, amin)
		return edakit.Numeric(q, edakit.IntStr(ext),
			edakit.Sprintf("The stub has area %d x %d = %d. It needs length ceil(%d / %d) = %d nm, so %d nm more.", w, l, w*l, amin, w, ceilDiv(amin, w), ext),
			map[string]string{"area": edakit.IntStr(w * l)})
	case 2:
		v := 10 * edakit.Int(3, 6)
		s := 10 * edakit.Int(3, 7)
		e := 5 * edakit.Int(1, 6)
		W := 2*e + v + edakit.Scaled(ctx, 0, 8)*(v+s) + edakit.Int(0, v+s-1)
		H := 2*e + v + edakit.Int(0, 5)*(v+s) + edakit.Int(0, v+s-1)
		nx := (W - 2*e + s) / (v + s)
		ny := (H - 2*e + s) / (v + s)
		q := edakit.Sprintf("Two wide wires overlap in a %d by %d nm region. Via cuts are %d nm squares, at least %d nm apart, and each must be enclosed by at least %d nm of metal on every side. What is the largest number of cuts that fit in a regular array in the overlap?", W, H, v, s, e)
		return edakit.Numeric(q, edakit.IntStr(nx*ny),
			edakit.Sprintf("Along x the usable length is %d - 2 x %d = %d, and n cuts need n x %d + (n - 1) x %d, so n = floor((%d + %d) / %d) = %d. Along y, %d. Total %d x %d = %d.", W, e, W-2*e, v, s, W-2*e, s, v+s, nx, ny, nx, ny, nx*ny),
			map[string]string{"nx": edakit.IntStr(nx), "ny": edakit.IntStr(ny)})
	default:
		w2 := 10 * edakit.Int(20, 40)
		w4 := w2 + 10*edakit.Int(20, 60)
		n := edakit.Scaled(ctx, 4, 9)
		var parts []string
		bad := 0
		for i := 0; i < n; i++ {
			w := 10 * edakit.Int(8, w4/10+30)
			c := edakit.Pick([]int{1, 1, 2, 2, 3, 4, 4})
			need := 1
			if w >= w4 {
				need = 4
			} else if w >= w2 {
				need = 2
			}
			if c < need {
				bad++
			}
			parts = append(parts, edakit.Sprintf("C%d (width %d, %d cuts)", i+1, w, c))
		}
		q := edakit.Sprintf("A minimum-cut rule says a via landing on a wire at least %d nm wide needs at least 2 cuts, and on a wire at least %d nm wide at least 4 cuts. The connections, with wire width in nm and cut count, are %s. How many violate the rule?", w2, w4, strings.Join(parts, ", "))
		return edakit.Numeric(q, edakit.IntStr(bad),
			edakit.Sprintf("Compare each connection's cut count with 1, 2 or 4 according to its width band: %d fall short.", bad),
			map[string]string{"violations": edakit.IntStr(bad)})
	}
}

// engineCandidates: pairs that survive a bounding-box filter.
func engineCandidates(ctx generator.GeneratorContext) generator.Problem {
	S := edakit.Int(15, 30)
	n := edakit.Scaled(ctx, 5, 8)
	rs := separated(n, 150, 10, 40)
	cand, viol := 0, 0
	for i := range rs {
		for j := i + 1; j < n; j++ {
			gx, gy := gaps(rs[i], rs[j])
			if gx < S && gy < S {
				cand++
			}
			if gx*gx+gy*gy < S*S {
				viol++
			}
		}
	}
	q := edakit.Sprintf("A DRC engine checks a %d nm Euclidean spacing rule. Before computing exact distances it keeps only pairs whose bounding boxes, one of them expanded by %d nm on every side, overlap with positive area. The shapes are %s (nm). How many pairs reach the exact distance test?", S, S, listRects("R", rs))
	why := edakit.Sprintf("A pair survives when both its x gap and its y gap are below %d nm. That keeps %d of the %d pairs. The exact test then confirms %d of them as violations: a diagonal pair can pass the box filter yet be %d nm or more apart.", S, cand, n*(n-1)/2, viol, S)
	return numChoice(q, float64(cand), 0, "pairs", why, float64(n*(n-1)/2), float64(viol), float64(cand+1))
}

// engineSizedMarker: area of the overlap of two shapes each grown by half the spacing.
func engineSizedMarker(ctx generator.GeneratorContext) generator.Problem {
	S := 2 * edakit.Int(10, 30)
	g := edakit.Int(2, S-2)
	h1, h2 := edakit.Int(30, 120), edakit.Int(30, 120)
	w1, w2 := edakit.Int(20, 80), edakit.Int(20, 80)
	a := rect{0, 0, w1, h1}
	y := edakit.Int(-h2+10, h1-10)
	b := rect{w1 + g, y, w1 + g + w2, y + h2}
	oy := overlap1(a.y1, a.y2, b.y1, b.y2)
	area := (S - g) * (oy + S)
	q := edakit.Sprintf("A polygon-based spacing check grows every shape by half the rule distance, %d nm, and reports the AND of grown shapes from different polygons. Shapes A = %s and B = %s (nm) are checked against a %d nm rule. What is the area of the reported marker, in square nm?", S/2, a, b, S)
	why := edakit.Sprintf("Grown by %d, the shapes overlap %d - %d = %d nm in x and %d + %d = %d nm in y (the y overlap of %d plus %d nm of growth on each end). Area %d. The marker is wider than the facing edges because sizing also grows the ends.", S/2, S, g, S-g, oy, S, oy+S, oy, S/2, area)
	return numChoice(q, float64(area), 0, "square nm", why, float64((S-g)*oy), float64((S-g)*(oy+S/2)), float64(g*(oy+S)))
}

// planeSweep: maximum active set (KP 1), y-interval query at one insertion (KP 2), pair tests
// saved by the sweep (KP 3).
func planeSweep(ctx generator.GeneratorContext) generator.Problem {
	S := edakit.Int(10, 30)
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 5, 10)
		var rs []rect
		for i := 0; i < n; i++ {
			rs = append(rs, randRect(200, 10, 60))
		}
		best := 0
		var xs []int
		for _, r := range rs {
			xs = append(xs, r.x1, r.x2+S)
		}
		for _, x := range xs {
			c := 0
			for _, r := range rs {
				if r.x1 <= x && x <= r.x2+S {
					c++
				}
			}
			best = imax(best, c)
		}
		q := edakit.Sprintf("A sweep-line spacing check for a %d nm rule moves a vertical line left to right. A shape enters the active set at its left edge and leaves at its right edge plus %d nm, so it stays available for comparison with shapes starting within the rule distance (closed intervals). The shapes are %s (nm). What is the largest number of shapes active at once?", S, S, listRects("R", rs))
		return edakit.Numeric(q, edakit.IntStr(best),
			edakit.Sprintf("Each shape is active on [x1, x2 + %d]. The maximum number of overlapping intervals, found by scanning the event points, is %d. The active set, not the whole layer, bounds the memory and the per-event work.", S, best),
			map[string]string{"max_active": edakit.IntStr(best)})
	case 2:
		n := edakit.Scaled(ctx, 5, 10)
		y1 := edakit.Int(40, 140)
		y2 := y1 + edakit.Int(10, 50)
		var parts []string
		hit := 0
		for i := 0; i < n; i++ {
			a := edakit.Int(0, 220)
			b := a + edakit.Int(5, 40)
			if a < y2+S && y1-S < b {
				hit++
			}
			parts = append(parts, edakit.Sprintf("[%d, %d]", a, b))
		}
		q := edakit.Sprintf("During a sweep for a %d nm spacing rule, the active shapes have y-intervals %s (nm). A new shape with y-interval [%d, %d] is inserted, and the engine queries its interval tree for active intervals that come closer than %d nm in y, that is, that overlap the open interval (%d, %d). How many active shapes does the query return?", S, strings.Join(parts, ", "), y1, y2, S, y1-S, y2+S)
		return edakit.Numeric(q, edakit.IntStr(hit),
			edakit.Sprintf("An interval [a, b] is returned when a < %d and b > %d. %d of the %d active intervals qualify; only those are compared exactly.", y2+S, y1-S, hit, n),
			map[string]string{"returned": edakit.IntStr(hit)})
	default:
		n := edakit.Scaled(ctx, 6, 10)
		rs := separated(n, 170, 10, 40)
		cand := 0
		for i := range rs {
			for j := i + 1; j < n; j++ {
				gx, gy := gaps(rs[i], rs[j])
				if gx < S && gy < S {
					cand++
				}
			}
		}
		saved := n*(n-1)/2 - cand
		q := edakit.Sprintf("For a %d nm spacing rule, a sweep with an interval tree compares only pairs whose x gap and y gap are both below %d nm. A brute-force checker compares every pair. The shapes are %s (nm). How many pair comparisons does the sweep avoid?", S, S, listRects("R", rs))
		return edakit.Numeric(q, edakit.IntStr(saved),
			edakit.Sprintf("Brute force makes %d x %d / 2 = %d comparisons. %d pairs are near in both axes, so the sweep avoids %d. On a real layer with millions of shapes, the near pairs grow roughly linearly while all pairs grow quadratically.", n, n-1, n*(n-1)/2, cand, saved),
			map[string]string{"all_pairs": edakit.IntStr(n * (n - 1) / 2), "near_pairs": edakit.IntStr(cand)})
	}
}

// hierErrorCount: flat versus hierarchical marker counts.
func hierErrorCount(ctx generator.GeneratorContext) generator.Problem {
	e, m := edakit.Int(1, 6), edakit.Scaled(ctx, 4, 400)
	f, k := edakit.Int(1, 5), edakit.Int(2, 60)
	t := edakit.Int(0, 9)
	flat, hier := e*m+f*k+t, e+f+t
	if edakit.Int(0, 1) == 0 {
		q := edakit.Sprintf("Cell A has %d internal spacing violations and is placed %d times; cell B has %d internal violations and is placed %d times; the top cell adds %d violations of its own between instances. No instance's context changes the internal results. How many violations does a flat DRC run report?", e, m, f, k, t)
		why := edakit.Sprintf("Flat checking sees every copy: %d x %d + %d x %d + %d = %d. A hierarchical run reports each cell's internal errors once, %d in total.", e, m, f, k, t, flat, hier)
		return numChoice(q, float64(flat), 0, "", why, float64(hier), float64(e*m+f*k), float64((e+f)*(m+k)+t))
	}
	q := edakit.Sprintf("Cell A has %d internal spacing violations and is placed %d times; cell B has %d internal violations and is placed %d times; the top cell adds %d violations between instances. No instance's context changes the internal results. How many violations does a hierarchical DRC run report, listing each cell's errors once in the cell's own coordinates?", e, m, f, k, t)
	why := edakit.Sprintf("Each cell's internal errors are reported once: %d + %d, plus the %d top-level ones, %d. A flat run would report %d.", e, f, t, hier, flat)
	return numChoice(q, float64(hier), 0, "", why, float64(flat), float64(e+f), float64(e*f+t))
}

// hierInteractions: interaction regions in a placed array.
func hierInteractions(ctx generator.GeneratorContext) generator.Problem {
	S := 10 * edakit.Int(4, 8)
	c, r := edakit.Scaled(ctx, 2, 12), edakit.Int(2, 8)
	w, h := 10*edakit.Int(50, 200), 10*edakit.Int(50, 200)
	gx, gy := 10*edakit.Int(1, S/10+3), 10*edakit.Int(1, S/10+3)
	cnt := 0
	if gx < S {
		cnt += (c - 1) * r
	}
	if gy < S {
		cnt += c * (r - 1)
	}
	if gx*gx+gy*gy < S*S {
		cnt += 2 * (c - 1) * (r - 1)
	}
	q := edakit.Sprintf("A cell %d by %d nm, with metal running along all four edges of its boundary, is placed in a %d column by %d row array with pitch %d nm in x and %d nm in y. The spacing rule is %d nm (Euclidean). A hierarchical engine must check every pair of neighbouring instances whose shapes come closer than the rule. How many instance pairs form such interaction regions?", w, h, c, r, w+gx, h+gy, S)
	why := edakit.Sprintf("The gaps are %d nm in x and %d nm in y, and the diagonal gap is sqrt(%d^2 + %d^2). Horizontal neighbours (%d of them) interact if %d < %d, vertical neighbours (%d) if %d < %d, and the 2 x %d diagonal neighbours if the diagonal gap is below %d. Total %d.", gx, gy, gx, gy, (c-1)*r, gx, S, c*(r-1), gy, S, (c-1)*(r-1), S, cnt)
	return numChoice(q, float64(cnt), 0, "", why, float64((c-1)*r+c*(r-1)), float64(c*r), float64((c-1)*(r-1)))
}

// deckBoolean: area of a derived layer.
func deckBoolean(ctx generator.GeneratorContext) generator.Problem {
	poly, act := rect{50, 0, 65, 100}, rect{20, 30, 120, 70}
	for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
		p := rect{edakit.Int(40, 80), edakit.Int(0, 20), 0, 0}
		p.x2 = p.x1 + edakit.Int(5, 20)
		p.y2 = p.y1 + edakit.Int(80, 140)
		a := rect{edakit.Int(0, 40), edakit.Int(25, 50), 0, 0}
		a.x2 = a.x1 + edakit.Int(60, 120)
		a.y2 = a.y1 + edakit.Int(20, 60)
		if interArea(p, a) > 0 {
			poly, act = p, a
			break
		}
	}
	and := interArea(poly, act)
	switch edakit.Int(0, 2) {
	case 0:
		q := edakit.Sprintf("A rule deck derives GATE = POLY AND ACTIVE. POLY is %s and ACTIVE is %s (nm). What is the area of GATE, in square nm?", poly, act)
		return numChoice(q, float64(and), 0, "square nm", edakit.Sprintf("AND keeps the region covered by both layers: the intersection, area %d.", and),
			float64(poly.area()), float64(act.area()), float64(poly.area()+act.area()-and))
	case 1:
		v := poly.area() - and
		q := edakit.Sprintf("A rule deck derives FIELDPOLY = POLY NOT ACTIVE, the poly outside the transistor channels. POLY is %s and ACTIVE is %s (nm). What is the area of FIELDPOLY, in square nm?", poly, act)
		return numChoice(q, float64(v), 0, "square nm", edakit.Sprintf("NOT removes the overlap: %d - %d = %d.", poly.area(), and, v),
			float64(and), float64(act.area()-and), float64(poly.area()))
	default:
		v := act.area() - and
		q := edakit.Sprintf("A rule deck derives SD = ACTIVE NOT POLY, the source and drain diffusion. POLY is %s and ACTIVE is %s (nm). What is the area of SD, in square nm?", poly, act)
		return numChoice(q, float64(v), 0, "square nm", edakit.Sprintf("NOT removes the channel from the active area: %d - %d = %d.", act.area(), and, v),
			float64(and), float64(poly.area()-and), float64(act.area()))
	}
}

// deckSized: area of (A sized by d) AND B.
func deckSized(ctx generator.GeneratorContext) generator.Problem {
	d := 5 * edakit.Int(1, 6)
	a := randRect(120, 20, 60)
	b := a.shift(10, 10)
	for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
		a2, b2 := randRect(120, 20, 60), randRect(120, 20, 60)
		if interArea(a2.grow(d), b2) > 0 {
			a, b = a2, b2
			break
		}
	}
	v := interArea(a.grow(d), b)
	q := edakit.Sprintf("A deck computes HALO = (M1 SIZED BY %d) AND M2: grow every M1 shape by %d nm on all sides, then intersect with M2. M1 is %s and M2 is %s (nm). What is the area of HALO, in square nm?", d, d, a, b)
	why := edakit.Sprintf("Grown M1 is %s. Its intersection with M2 has area %d. Without the sizing the overlap would be %d.", a.grow(d), v, interArea(a, b))
	return numChoice(q, float64(v), 0, "square nm", why, float64(interArea(a, b)), float64(a.grow(d).area()), float64(interArea(a.grow(2*d), b)))
}
