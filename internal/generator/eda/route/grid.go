package route

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// gcellsTracks: tracks inside one gcell (KP 1), tracks left after a blockage (KP 2), gcell-edge
// capacity after layer adjustment on two same-direction layers (KP 3).
func gcellsTracks(ctx generator.GeneratorContext) generator.Problem {
	pitches := []int{80, 100, 120, 140, 160, 200, 240, 280}
	switch kp(ctx) {
	case 1:
		p := edakit.Pick(pitches)
		G := p * edakit.Int(8, 20)
		G += 10 * edakit.Int(0, p/10-1)
		o := 10 * edakit.Int(0, p/10-1)
		i := edakit.Int(1, 30)
		x0, x1 := i*G, (i+1)*G
		n := 0
		for k := 0; o+k*p < x1; k++ {
			if o+k*p >= x0 {
				n++
			}
		}
		q := edakit.Sprintf("Vertical routing tracks on a metal layer sit at x = %d + %dk nm for k = 0, 1, 2, and so on. Gcells are %d nm wide, and gcell column %d covers the half-open range [%d, %d) nm. How many tracks of this layer run through that gcell column?", o, p, G, i, x0, x1)
		return edakit.Numeric(q, edakit.IntStr(n),
			edakit.Sprintf("Count the integers k with %d <= %d + %dk < %d. The first such track is at k = %d and the last at k = %d, so there are %d tracks.", x0, o, p, x1, ceilDiv(x0-o, p), (x1-o-1)/p, n),
			map[string]string{"tracks": edakit.IntStr(n)})
	case 2:
		p := edakit.Pick(pitches)
		G := p * edakit.Int(10, 20)
		o := p / 2
		m := edakit.Pick([]int{p / 4, p / 2, 3 * p / 4})
		b1 := 10 * edakit.Int(G/40, G/20)
		b2 := b1 + 10*edakit.Int(p/10, 4*p/10)
		total, lost := 0, 0
		for k := 0; o+k*p < G; k++ {
			x := o + k*p
			total++
			if x > b1-m && x < b2+m {
				lost++
			}
		}
		q := edakit.Sprintf("A gcell spans x in [0, %d) nm. Its vertical tracks sit at x = %d + %dk nm. A macro obstruction covers x from %d to %d nm over the full gcell height. A track is unusable if its centerline lies strictly within %d nm (half the wire width plus the minimum spacing) of the obstruction, that is, at %d < x < %d. How many usable tracks does the gcell keep?", G, o, p, b1, b2, m, b1-m, b2+m)
		return edakit.Numeric(q, edakit.IntStr(total-lost),
			edakit.Sprintf("The gcell holds %d tracks. Those with centerline in (%d, %d) are lost: %d of them. The usable capacity is %d.", total, b1-m, b2+m, lost, total-lost),
			map[string]string{"tracks": edakit.IntStr(total), "lost": edakit.IntStr(lost)})
	default:
		t1, t2 := edakit.Int(10, 30), edakit.Int(8, 24)
		a1 := edakit.Pick([]int{0, 10, 20, 30, 40, 50})
		a2 := edakit.Pick([]int{0, 10, 20, 30, 40, 50, 60})
		c1 := t1 * (100 - a1) / 100
		c2 := t2 * (100 - a2) / 100
		q := edakit.Sprintf("A global router builds the capacity of a horizontal gcell edge from two horizontal layers. Metal 2 has %d tracks crossing the edge and a layer adjustment of %d%%; metal 4 has %d tracks and an adjustment of %d%%. The adjustment removes that share of the layer's tracks, and each layer's remaining capacity is rounded down to an integer. What capacity does the router assign to the edge?", t1, a1, t2, a2)
		return edakit.Numeric(q, edakit.IntStr(c1+c2),
			edakit.Sprintf("Metal 2 keeps floor(%d x %s) = %d tracks and metal 4 keeps floor(%d x %s) = %d. The edge capacity is %d. Adjustments reserve room for pin access, power straps and detailed-routing detours that the gcell model cannot see.", t1, edakit.Num(float64(100-a1)/100, 2), c1, t2, edakit.Num(float64(100-a2)/100, 2), c2, c1+c2),
			map[string]string{"m2": edakit.IntStr(c1), "m4": edakit.IntStr(c2)})
	}
}

// preferredDirection: horizontal capacity across a stack (KP 1), highest layer a demand needs
// (KP 2), the longest wrong-way jog cheaper than a via pair (KP 3).
func preferredDirection(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	switch k {
	case 1, 2:
		n := edakit.Int(6, 9)
		firstH := edakit.Int(0, 1) == 0
		tr := make([]int, n+1)
		var parts []string
		horiz := func(l int) bool { return (l%2 == 1) == firstH }
		for l := 2; l <= n; l++ {
			tr[l] = edakit.Int(4, 16) - l/2
			if tr[l] < 3 {
				tr[l] = 3
			}
			d := "vertical"
			if horiz(l) {
				d = "horizontal"
			}
			parts = append(parts, edakit.Sprintf("M%d %s with %d tracks per gcell", l, d, tr[l]))
		}
		m1 := "vertical"
		if firstH {
			m1 = "horizontal"
		}
		intro := edakit.Sprintf("A stack has metal layers M1 to M%d with alternating preferred directions, M1 being %s. M1 is reserved for cell pins, so routing uses M2 and above: %s. ", n, m1, strings.Join(parts, "; "))
		if k == 1 {
			sum := 0
			var used []string
			for l := 2; l <= n; l++ {
				if horiz(l) {
					sum += tr[l]
					used = append(used, edakit.Sprintf("M%d", l))
				}
			}
			q := intro + "Wires crossing the boundary between two horizontally adjacent gcells run east-west. What is the capacity of that boundary, in tracks, counting every routing layer that can carry such wires in its preferred direction?"
			return edakit.Numeric(q, edakit.IntStr(sum),
				edakit.Sprintf("An east-west crossing uses horizontal tracks, so only the horizontal layers %s count. Their tracks sum to %d. Vertical layers contribute nothing to this edge unless the router allows wrong-way wiring.", strings.Join(used, ", "), sum),
				map[string]string{"capacity": edakit.IntStr(sum)})
		}
		sum, top := 0, 0
		total := 0
		for l := 2; l <= n; l++ {
			if horiz(l) {
				total += tr[l]
			}
		}
		D := edakit.Int(total/3+1, total)
		for l := 2; l <= n && sum < D; l++ {
			if horiz(l) {
				sum += tr[l]
				top = l
			}
		}
		q := intro + edakit.Sprintf("A layer assigner places %d east-west wires across one gcell boundary and fills horizontal layers from the bottom up, using every track of a layer before moving to the next horizontal layer. Which is the highest metal layer it must use? Answer with the layer number.", D)
		return edakit.Numeric(q, edakit.IntStr(top),
			edakit.Sprintf("Accumulating horizontal-layer tracks from the bottom, the running total first reaches %d at M%d (cumulative %d tracks).", D, top, sum),
			map[string]string{"layer": edakit.IntStr(top)})
	default:
		cv := edakit.Int(2, 8)
		f := edakit.Pick([]float64{1.5, 2, 2.5, 3, 4})
		lim := 2 * float64(cv) / (f - 1)
		nmax := int(lim)
		if float64(nmax) >= lim-1e-9 {
			nmax--
		}
		q := edakit.Sprintf("A detailed router charges 1 per track step in a layer's preferred direction and %s per step against it (wrong-way), and %d per via. A net on a horizontal layer needs a short vertical jog of n steps. It can jog wrong-way on the same layer, or drop a via to the vertical layer below, run n preferred-direction steps and via back up. What is the largest n for which the wrong-way jog is strictly cheaper?", edakit.Num(f, 1), cv)
		return edakit.Numeric(q, edakit.IntStr(nmax),
			edakit.Sprintf("Wrong-way costs %sn; the via detour costs n + 2(%d). The jog wins when %sn < n + %d, that is n < %s, so the largest integer is %d.", edakit.Num(f, 1), cv, edakit.Num(f, 1), 2*cv, edakit.Num(lim, 3), nmax),
			map[string]string{"n_max": edakit.IntStr(nmax)})
	}
}

// gridGraph: edges of the 2D gcell graph (KP 1), of the 3D graph with preferred directions and
// vias (KP 2), edges touched by an interior macro (KP 3).
func gridGraph(ctx generator.GeneratorContext) generator.Problem {
	X, Y := edakit.Scaled(ctx, 10, 200), edakit.Scaled(ctx, 10, 200)
	switch kp(ctx) {
	case 1:
		e := (X-1)*Y + X*(Y-1)
		q := edakit.Sprintf("A global router tiles the die into %d columns by %d rows of gcells. Its 2D routing graph has one vertex per gcell and one edge between each pair of horizontally or vertically adjacent gcells. How many edges does it have?", X, Y)
		return edakit.Numeric(q, edakit.IntStr(e),
			edakit.Sprintf("There are (%d - 1) x %d = %d horizontal edges and %d x (%d - 1) = %d vertical edges, %d in all. Each edge carries a capacity and a demand.", X, Y, (X-1)*Y, X, Y, X*(Y-1), e),
			map[string]string{"h_edges": edakit.IntStr((X - 1) * Y), "v_edges": edakit.IntStr(X * (Y - 1))})
	case 2:
		L := edakit.Int(3, 8)
		firstH := edakit.Int(0, 1) == 0
		nh, nv := 0, 0
		for l := 1; l <= L; l++ {
			if (l%2 == 1) == firstH {
				nh++
			} else {
				nv++
			}
		}
		wires := nh*(X-1)*Y + nv*X*(Y-1)
		vias := (L - 1) * X * Y
		d := "vertical"
		if firstH {
			d = "horizontal"
		}
		q := edakit.Sprintf("A 3D global-routing graph covers %d by %d gcells on %d routing layers. The lowest routing layer is %s and directions alternate upward. Each layer has wire edges only between gcells adjacent along its preferred direction, and a via edge joins each gcell to the same gcell on the next layer up. How many edges does the graph have?", X, Y, L, d)
		return edakit.Numeric(q, edakit.IntStr(wires+vias),
			edakit.Sprintf("%d horizontal layers give %d x %d = %d wire edges, %d vertical layers give %d x %d = %d, and %d via levels give %d x %d = %d via edges. The total is %d.", nh, nh, (X-1)*Y, nh*(X-1)*Y, nv, nv, X*(Y-1), nv*X*(Y-1), L-1, L-1, X*Y, vias, wires+vias),
			map[string]string{"wire_edges": edakit.IntStr(wires), "via_edges": edakit.IntStr(vias)})
	default:
		a, b := edakit.Int(2, 12), edakit.Int(2, 12)
		x0, y0 := edakit.Int(2, X-a-1), edakit.Int(2, Y-b-1)
		if X-a-1 < 2 || Y-b-1 < 2 {
			x0, y0 = 2, 2
			X, Y = maxI(X, a+4), maxI(Y, b+4)
		}
		in := (a-1)*b + a*(b-1)
		bd := 2*a + 2*b
		q := edakit.Sprintf("On a %d by %d gcell grid, a macro covers the gcells in columns %d to %d and rows %d to %d, away from the die edge. The router must lower the capacity of every 2D edge that has at least one endpoint inside the macro. How many edges is that?", X, Y, x0, x0+a-1, y0, y0+b-1)
		return edakit.Numeric(q, edakit.IntStr(in+bd),
			edakit.Sprintf("The macro is %d by %d gcells. Edges with both ends inside: (%d - 1) x %d + %d x (%d - 1) = %d. Edges crossing its boundary: 2(%d) + 2(%d) = %d. Total %d.", a, b, a, b, a, b, in, a, b, bd, in+bd),
			map[string]string{"inside": edakit.IntStr(in), "boundary": edakit.IntStr(bd)})
	}
}
