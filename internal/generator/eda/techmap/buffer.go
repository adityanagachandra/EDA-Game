package techmap

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// cand is a van Ginneken candidate: downstream load C (fF) and required time q (ps).
type cand struct{ c, q float64 }

func (x cand) String() string {
	return "(" + edakit.Num(x.c, 1) + " fF, " + edakit.Num(x.q, 1) + " ps)"
}

func candList(cs []cand) string {
	var s []string
	for _, c := range cs {
		s = append(s, c.String())
	}
	return strings.Join(s, ", ")
}

// prune keeps the non-dominated candidates: no other has C <= and q >= with one strict.
func prune(cs []cand) []cand {
	var out []cand
	for i, a := range cs {
		dom := false
		for j, b := range cs {
			if i == j {
				continue
			}
			if b.c <= a.c && b.q >= a.q && (b.c < a.c || b.q > a.q || j < i) {
				dom = true
				break
			}
		}
		if !dom {
			out = append(out, a)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].c < out[j].c })
	return out
}

func wire(x cand, r, cw float64) cand { return cand{x.c + cw, x.q - r*(cw/2+x.c)} }

type buf struct{ cb, rb, db float64 }

func (b buf) apply(x cand) cand { return cand{b.cb, x.q - b.db - b.rb*x.c} }

func randBuf() buf {
	return buf{float64(edakit.Int(3, 8)), float64(edakit.Int(5, 15)) / 10, float64(edakit.Int(10, 30))}
}

func (b buf) String() string {
	return edakit.Sprintf("input capacitance %s fF, drive resistance %s kOhm, intrinsic delay %s ps", edakit.Num(b.cb, 0), edakit.Num(b.rb, 1), edakit.Num(b.db, 0))
}

const vgUnits = "Resistance is in kOhm and capacitance in fF, so R x C is in ps."

// randFront draws a non-dominated list: C and q both increase.
func randFront(n int) []cand {
	var cs []cand
	c, q := float64(edakit.Int(5, 15)), float64(edakit.Int(300, 400))
	for i := 0; i < n; i++ {
		cs = append(cs, cand{c, q})
		c += float64(edakit.Int(5, 20))
		q += float64(edakit.Int(5, 30))
	}
	return cs
}

func bufferCandidates(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		x := cand{float64(edakit.Int(5, 60)), float64(edakit.Int(300, 600))}
		r, cw := float64(edakit.Int(1, 8))/10, float64(edakit.Int(10, 80))
		y := wire(x, r, cw)
		return edakit.Numeric(edakit.Sprintf("%s A van Ginneken candidate at the downstream end of a wire segment is (C, q) = %s. The segment has resistance %s kOhm and capacitance %s fF, modelled as a pi segment with Elmore delay. What is the required time q at the upstream end of the segment, in ps? Round to 1 decimal place.", vgUnits, x, edakit.Num(r, 1), edakit.Num(cw, 0)),
			edakit.Num(y.q, 1),
			edakit.Sprintf("q' = q - R (C_w / 2 + C) = %s - %s x (%s + %s) = %s ps, and the load grows to C' = C + C_w = %s fF.", edakit.Num(x.q, 0), edakit.Num(r, 1), edakit.Num(cw/2, 1), edakit.Num(x.c, 0), edakit.Num(y.q, 1), edakit.Num(y.c, 0)),
			map[string]string{"q": edakit.Num(y.q, 1), "c": edakit.Num(y.c, 0)})
	case 2:
		x := cand{float64(edakit.Int(20, 150)), float64(edakit.Int(300, 600))}
		b := randBuf()
		y := b.apply(x)
		return edakit.Numeric(edakit.Sprintf("%s At a buffer site the candidate is (C, q) = %s. A buffer with %s is inserted at the site. What is the new candidate's required time q, in ps? Round to 1 decimal place.", vgUnits, x, b),
			edakit.Num(y.q, 1),
			edakit.Sprintf("The buffer drives the downstream load: q' = q - d_b - R_b C = %s - %s - %s x %s = %s ps. Upstream it presents only C' = C_b = %s fF, which is the point of buffering.", edakit.Num(x.q, 0), edakit.Num(b.db, 0), edakit.Num(b.rb, 1), edakit.Num(x.c, 0), edakit.Num(y.q, 1), edakit.Num(b.cb, 0)),
			map[string]string{"q": edakit.Num(y.q, 1), "c": edakit.Num(y.c, 0)})
	}
	l1, l2 := randFront(edakit.Int(2, 3)), randFront(edakit.Int(2, 3))
	var merged []cand
	for _, a := range l1 {
		for _, b := range l2 {
			merged = append(merged, cand{a.c + b.c, minF(a.q, b.q)})
		}
	}
	kept := prune(merged)
	return edakit.Numeric(edakit.Sprintf("%s At a branch point, the left subtree's candidates are %s and the right subtree's are %s. Merging pairs one candidate from each side into (C_1 + C_2, min(q_1, q_2)). After merging all pairs and removing dominated candidates (another candidate with no more load and no less required time), how many candidates remain?", vgUnits, candList(l1), candList(l2)),
		edakit.IntStr(len(kept)),
		edakit.Sprintf("The %d merged pairs reduce to %d non-dominated candidates: %s. The linear-time merge that walks both sorted lists produces at most |L1| + |L2| - 1 of them.", len(merged), len(kept), candList(kept)),
		map[string]string{"merged": edakit.IntStr(len(merged)), "kept": edakit.IntStr(len(kept))})
}

func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func bufferPruning(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 5, 8)
		var cs []cand
		for i := 0; i < n; i++ {
			cs = append(cs, cand{float64(edakit.Int(1, 12) * 5), float64(edakit.Int(60, 90) * 5)})
		}
		kept := prune(cs)
		return edakit.Numeric(edakit.Sprintf("A van Ginneken node holds candidates (C, q): %s. A candidate is dominated if another has C no larger and q no smaller (exact duplicates collapse to one). How many candidates survive pruning?", candList(cs)),
			edakit.IntStr(len(kept)),
			edakit.Sprintf("The survivors are %s: sorted by C, their q values strictly increase, which is the shape every pruned list has.", candList(kept)),
			map[string]string{"kept": edakit.IntStr(len(kept))})
	case 2:
		cs := randFront(edakit.Int(3, 5))
		rd := float64(edakit.Int(5, 30)) / 10
		best, bc := -1e9, cand{}
		var vals []string
		for _, c := range cs {
			v := c.q - rd*c.c
			vals = append(vals, edakit.Num(v, 1))
			if v > best {
				best, bc = v, c
			}
		}
		return edakit.Numeric(edakit.Sprintf("%s The pruned candidate list at the source of a net is %s. The driver has output resistance %s kOhm, so each candidate gives a required time at the driver of q - R_d C. What is the best required time at the driver, in ps? Round to 1 decimal place.", vgUnits, candList(cs), edakit.Num(rd, 1)),
			edakit.Num(best, 1),
			edakit.Sprintf("The candidates give %s ps. The best is %s ps from %s; tracing back its buffer decisions gives the solution.", strings.Join(vals, ", "), edakit.Num(best, 1), bc),
			map[string]string{"best_q": edakit.Num(best, 1)})
	}
	cs := randFront(edakit.Int(3, 5))
	b := randBuf()
	all := append([]cand{}, cs...)
	var bs []cand
	for _, c := range cs {
		bs = append(bs, b.apply(c))
	}
	all = append(all, bs...)
	kept := prune(all)
	return edakit.Numeric(edakit.Sprintf("%s A buffer site holds the pruned candidates %s. Each candidate may also be buffered with a buffer of %s. After adding every buffered candidate and pruning the combined list, how many candidates remain?", vgUnits, candList(cs), b),
		edakit.IntStr(len(kept)),
		edakit.Sprintf("Every buffered candidate has C = %s fF, so only the one with the largest q, %s ps, can survive; it then competes with the unbuffered ones. The final list is %s.", edakit.Num(b.cb, 0), edakit.Num(bestBuffered(bs), 1), candList(kept)),
		map[string]string{"kept": edakit.IntStr(len(kept))})
}

func bestBuffered(bs []cand) float64 {
	b := -1e9
	for _, c := range bs {
		if c.q > b {
			b = c.q
		}
	}
	return b
}

func vanGinneken(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	n := map[int]int{1: 1, 2: 2, 3: 3}[k]
	if k == 3 && ctx.Difficulty < 0.4 {
		n = 2
	}
	for {
		b := randBuf()
		rd := float64(edakit.Int(10, 30)) / 10
		sink := cand{float64(edakit.Int(5, 20)), float64(edakit.Int(400, 600))}
		type seg struct{ r, c float64 }
		segs := make([]seg, n+1)
		same := k != 3
		s0 := seg{float64(edakit.Int(2, 8)) / 10, float64(edakit.Int(20, 80))}
		for i := range segs {
			if same {
				segs[i] = s0
			} else {
				segs[i] = seg{float64(edakit.Int(2, 8)) / 10, float64(edakit.Int(20, 80))}
			}
		}
		// Enumerate placements; site i sits between segs[i] (upstream) and segs[i+1].
		eval := func(mask int) float64 {
			x := sink
			for i := n; i >= 0; i-- {
				x = wire(x, segs[i].r, segs[i].c)
				if i > 0 && mask&(1<<uint(i-1)) != 0 {
					x = b.apply(x)
				}
			}
			return x.q - rd*x.c
		}
		best, bestMask := -1e9, 0
		for m := 0; m < 1<<uint(n); m++ {
			if v := eval(m); v > best+1e-9 {
				best, bestMask = v, m
			}
		}
		none := eval(0)
		if k <= 2 && bestMask == 0 && edakit.Int(0, 2) > 0 {
			continue
		}
		var layout []string
		layout = append(layout, "driver")
		for i := 0; i <= n; i++ {
			if same {
				layout = append(layout, "segment")
			} else {
				layout = append(layout, edakit.Sprintf("segment (%s kOhm, %s fF)", edakit.Num(segs[i].r, 1), edakit.Num(segs[i].c, 0)))
			}
			if i < n {
				layout = append(layout, edakit.Sprintf("site %d", i+1))
			}
		}
		layout = append(layout, "sink")
		segText := ""
		if same {
			segText = edakit.Sprintf(" Every segment has %s kOhm and %s fF (Elmore, pi model).", edakit.Num(s0.r, 1), edakit.Num(s0.c, 0))
		} else {
			segText = " Segments use the Elmore pi model."
		}
		var used []string
		for i := 0; i < n; i++ {
			if bestMask&(1<<uint(i)) != 0 {
				used = append(used, edakit.Sprintf("site %d", i+1))
			}
		}
		choice := "no buffer"
		if len(used) > 0 {
			choice = "buffers at " + strings.Join(used, ", ")
		}
		q := edakit.Sprintf("%s A two-pin net runs %s.%s The sink has load %s fF and required time %s ps. Each site may take one buffer with %s. The driver has output resistance %s kOhm, and the required time at the driver is q - R_d C for the candidate seen there. Run van Ginneken's algorithm. What is the best required time at the driver, in ps? Round to 1 decimal place.",
			vgUnits, strings.Join(layout, " -> "), segText, edakit.Num(sink.c, 0), edakit.Num(sink.q, 0), b, edakit.Num(rd, 1))
		return edakit.Numeric(q, edakit.Num(best, 1),
			edakit.Sprintf("Propagate candidates from the sink: through each segment apply q' = q - R (C_w/2 + C), C' = C + C_w; at each site keep both the unbuffered candidate and the buffered one (C_b, q - d_b - R_b C), pruning dominated ones. The best driver required time is %s ps with %s; with no buffers it would be %s ps.", edakit.Num(best, 1), choice, edakit.Num(none, 1)),
			map[string]string{"best_q": edakit.Num(best, 1), "unbuffered_q": edakit.Num(none, 1)})
	}
}
