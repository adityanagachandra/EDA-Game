package route

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// timingDriven compares a shared-trunk topology with a star for the delay to a critical sink.
func timingDriven(ctx generator.GeneratorContext) generator.Problem {
	rd := float64(50 * edakit.Int(2, 12))
	r := edakit.Pick([]float64{0.5, 1, 1.5, 2})
	c := edakit.Pick([]float64{0.1, 0.15, 0.2})
	l0 := float64(50 * edakit.Int(2, 8))
	l1 := float64(50 * edakit.Int(1, 6))
	l2 := float64(50 * edakit.Int(1, 6))
	c1, c2 := float64(edakit.Int(2, 10)), float64(edakit.Int(2, 10))
	// Topology A: trunk l0 to a Steiner point, branches l1 (sink 1) and l2 (sink 2).
	totA := c*(l0+l1+l2) + c1 + c2
	dA := rd*totA + r*l0*(c*l0/2+c*(l1+l2)+c1+c2) + r*l1*(c*l1/2+c1)
	// Topology B: separate wires of length l0+l1 and l0+l2 from the driver.
	d1, d2 := l0+l1, l0+l2
	totB := c*(d1+d2) + c1 + c2
	dB := rd*totB + r*d1*(c*d1/2+c1)
	model := edakit.Sprintf("A driver of %s ohm drives two sinks. Wire has %s ohm/um and %s fF/um, modelled as distributed RC, and delays are Elmore. Sink 1 (the critical one) has %s fF of pin capacitance and sink 2 has %s fF. Topology A runs a %s um trunk to a Steiner point, then %s um to sink 1 and %s um to sink 2. Topology B runs separate wires from the driver, %s um to sink 1 and %s um to sink 2. One ohm times one fF is one fs. ",
		edakit.Num(rd, 0), edakit.Num(r, 2), edakit.Num(c, 2), edakit.Num(c1, 0), edakit.Num(c2, 0), edakit.Num(l0, 0), edakit.Num(l1, 0), edakit.Num(l2, 0), edakit.Num(d1, 0), edakit.Num(d2, 0))
	switch kp(ctx) {
	case 1:
		ans := edakit.Num(dA/1000, 2)
		return edakit.Numeric(model+"What is the delay to sink 1 under topology A, in ps? Round to 2 decimal places.", ans,
			edakit.Sprintf("Driver: %s x %s fF. Trunk: $r l_0 (c l_0/2 + C_{below})$ where everything past the Steiner point hangs below it. Branch: $r l_1 (c l_1/2 + C_1)$. Total %s ohm-fF = %s ps.", edakit.Num(rd, 0), edakit.Num(totA, 1), edakit.Num(dA, 1), ans),
			map[string]string{"a_ps": ans})
	case 2:
		ans := edakit.Num(dB/1000, 2)
		return edakit.Numeric(model+"What is the delay to sink 1 under topology B, in ps? Round to 2 decimal places.", ans,
			edakit.Sprintf("The driver still charges all %s fF of both wires and pins, but sink 1's own wire carries only its own capacitance: $r d_1 (c d_1/2 + C_1)$. Total %s ohm-fF = %s ps.", edakit.Num(totB, 1), edakit.Num(dB, 1), ans),
			map[string]string{"b_ps": ans})
	default:
		best := math.Min(dA, dB)
		ans := edakit.Num(best/1000, 2)
		return edakit.Numeric(model+"A timing-driven router picks the topology with the smaller delay to the critical sink. What is that delay, in ps? Round to 2 decimal places.", ans,
			edakit.Sprintf("Topology A gives %s ps and B gives %s ps. B spends more wire and loads the driver more, but it isolates sink 1 from sink 2's branch; which effect wins depends on the driver resistance against the wire resistance.", edakit.Num(dA/1000, 2), edakit.Num(dB/1000, 2)),
			map[string]string{"a_ps": edakit.Num(dA/1000, 2), "b_ps": edakit.Num(dB/1000, 2)})
	}
}

// primDijkstra builds the Prim-Dijkstra tree for trade-off c; ok is false on a near tie.
func primDijkstra(ps []pt, c float64) (wl int, radius int, ok bool) {
	n := len(ps)
	in := make([]bool, n)
	pl := make([]int, n)
	in[0] = true
	ok = true
	for k := 1; k < n; k++ {
		best, second := math.Inf(1), math.Inf(1)
		bu, bv := -1, -1
		for u := 0; u < n; u++ {
			if !in[u] {
				continue
			}
			for v := 0; v < n; v++ {
				if in[v] {
					continue
				}
				key := c*float64(pl[u]) + float64(manh(ps[u], ps[v]))
				if key < best {
					second, best = best, key
					bu, bv = u, v
				} else if key < second {
					second = key
				}
			}
		}
		if second-best < 1e-9 {
			ok = false
		}
		in[bv] = true
		pl[bv] = pl[bu] + manh(ps[bu], ps[bv])
		wl += manh(ps[bu], ps[bv])
		radius = maxI(radius, pl[bv])
	}
	return wl, radius, ok
}

// primDijkstraGen: Prim tree length (KP 1), Dijkstra tree length (KP 2), radius at an
// intermediate trade-off (KP 3).
func primDijkstraGen(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	c := 0.0
	switch k {
	case 2:
		c = 1
	case 3:
		c = edakit.Pick([]float64{0.25, 0.5, 0.75})
	}
	for try := 0; ; try++ {
		ps := randPins(edakit.Scaled(ctx, 4, 6), 20)
		wl, rad, ok := primDijkstra(ps, c)
		w0, r0, ok0 := primDijkstra(ps, 0)
		if (!ok || !ok0) && try < generator.MaxDrawAttempts {
			continue
		}
		intro := edakit.Sprintf("The Prim-Dijkstra construction grows a routing tree from the source at %s over the sinks %s with Manhattan distances. At each step it adds the outside pin v and tree pin u minimising $c \\cdot pl(u) + d(u, v)$, where $pl(u)$ is the tree path length from the source to u. ", ptText(ps[0]), pinsText(ps[1:]))
		switch k {
		case 1:
			return edakit.Numeric(intro+"With c = 0 the rule is Prim's algorithm. What is the total wirelength of the tree?", edakit.IntStr(wl),
				edakit.Sprintf("With c = 0 only the edge length matters, so the tree is a rectilinear MST of length %d. Its radius, the longest source-to-sink path, is %d.", wl, rad),
				map[string]string{"wl": edakit.IntStr(wl), "radius": edakit.IntStr(rad)})
		case 2:
			return edakit.Numeric(intro+"With c = 1 the rule is Dijkstra's algorithm, giving a shortest-path tree. What is the total wirelength of the tree?", edakit.IntStr(wl),
				edakit.Sprintf("With c = 1 the key is the source path length through u, so every sink gets its shortest path, at a wirelength of %d against %d for the MST. Its radius is %d.", wl, w0, rad),
				map[string]string{"wl": edakit.IntStr(wl), "mst": edakit.IntStr(w0)})
		default:
			return edakit.Numeric(intro+edakit.Sprintf("With c = %s, what is the radius of the tree, the longest tree path from the source to a sink?", edakit.Num(c, 2)), edakit.IntStr(rad),
				edakit.Sprintf("Growing the tree with keys $%s \\cdot pl(u) + d(u, v)$ gives wirelength %d and radius %d. The MST (c = 0) has wirelength %d and radius %d: raising c trades wirelength for shorter source-sink paths.", edakit.Num(c, 2), wl, rad, w0, r0),
				map[string]string{"radius": edakit.IntStr(rad), "wl": edakit.IntStr(wl)})
		}
	}
}

// crosstalk: coupling capacitance (KP 1), the delay window from Miller factors (KP 2), the
// worst-case delay saved by doubling the spacing (KP 3).
func crosstalk(ctx generator.GeneratorContext) generator.Problem {
	c0 := edakit.Pick([]float64{0.04, 0.05, 0.06, 0.08, 0.1})
	L := float64(50 * edakit.Int(2, 40))
	m := edakit.Int(1, 3)
	rd := float64(50 * edakit.Int(4, 20))
	cc := c0 * L / float64(m)
	model := edakit.Sprintf("To first order, coupling capacitance between parallel wires scales with their parallel run length and inversely with their spacing: at the minimum spacing it is %s fF/um. ", edakit.Num(c0, 2))
	switch kp(ctx) {
	case 1:
		ans := edakit.Num(cc, 2)
		q := model + edakit.Sprintf("A victim runs parallel to an aggressor for %s um at %d times the minimum spacing. What is their coupling capacitance, in fF? Round to 2 decimal places.", edakit.Num(L, 0), m)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("%s x %s / %d = %s fF.", edakit.Num(c0, 2), edakit.Num(L, 0), m, ans),
			map[string]string{"cc": ans})
	case 2:
		w := rd * 2 * cc / 1000
		ans := edakit.Num(w, 2)
		q := model + edakit.Sprintf("A victim driven through %s ohm couples to an aggressor over %s um at %d times the minimum spacing. With a Miller factor of 0 when the aggressor switches with the victim and 2 when it switches against it, by how much can the victim's delay vary, in ps? Use delay = driver resistance times effective capacitance. Round to 2 decimal places.", edakit.Num(rd, 0), edakit.Num(L, 0), m)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$C_c$ = %s fF. The effective coupling runs from 0 to $2C_c$, so the delay varies by %s x %s fF = %s ps.", edakit.Num(cc, 2), edakit.Num(rd, 0), edakit.Num(2*cc, 2), ans),
			map[string]string{"cc": edakit.Num(cc, 2)})
	default:
		cc2 := cc / 2
		sv := rd * 2 * (cc - cc2) / 1000
		ans := edakit.Num(sv, 2)
		q := model + edakit.Sprintf("A crosstalk-aware router moves a critical victim, driven through %s ohm, from %d to %d times the minimum spacing from an aggressor over a %s um parallel run. Taking the worst-case Miller factor of 2 and delay = driver resistance times effective capacitance, how much worst-case delay does it save, in ps? Round to 2 decimal places.", edakit.Num(rd, 0), m, 2*m, edakit.Num(L, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Doubling the spacing halves $C_c$ from %s to %s fF. Worst-case delay falls by %s x 2 x %s fF = %s ps. Spacing costs tracks; a grounded shield costs a track too but also removes the switching aggressor.", edakit.Num(cc, 2), edakit.Num(cc2, 2), edakit.Num(rd, 0), edakit.Num(cc-cc2, 2), ans),
			map[string]string{"saved": ans})
	}
}

// antenna: partial antenna ratio (KP 1), the longest legal run (KP 2), cumulative ratio over
// several layers (KP 3).
func antenna(ctx generator.GeneratorContext) generator.Problem {
	gw := edakit.Pick([]float64{0.2, 0.3, 0.5, 1})
	gl := edakit.Pick([]float64{0.03, 0.05, 0.1})
	ng := edakit.Int(1, 3)
	ga := gw * gl * float64(ng)
	w := edakit.Pick([]float64{0.05, 0.07, 0.1})
	gates := edakit.Sprintf("A net drives %d identical gates of width %s um and length %s um. ", ng, edakit.Num(gw, 2), edakit.Num(gl, 2))
	def := "During etch of a metal layer, the metal of that layer connected to the gates but not yet to any diffusion collects charge; its partial antenna ratio is that metal's top area divided by the total gate area. "
	switch kp(ctx) {
	case 1:
		L := float64(10 * edakit.Int(2, 60))
		par := L * w / ga
		ans := edakit.Num(par, 1)
		q := gates + def + edakit.Sprintf("On metal 2, the gate-connected piece is %s um long and %s um wide, and no diffusion is connected until a higher layer. What is its partial antenna ratio? Round to 1 decimal place.", edakit.Num(L, 0), edakit.Num(w, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Metal area %s x %s = %s um^2; gate area %d x %s x %s = %s um^2. Ratio %s.", edakit.Num(L, 0), edakit.Num(w, 2), edakit.Num(L*w, 3), ng, edakit.Num(gw, 2), edakit.Num(gl, 2), edakit.Num(ga, 4), ans),
			map[string]string{"metal_area": edakit.Num(L*w, 3), "gate_area": edakit.Num(ga, 4)})
	case 2:
		R := float64(50 * edakit.Int(4, 20))
		lmax := R * ga / w
		ans := edakit.Num(lmax, 1)
		q := gates + def + edakit.Sprintf("The rule allows a partial ratio of at most %s on metal 3. A metal-3 wire %s um wide connects to the gates before any diffusion is attached. What is the longest such wire, in um, that meets the rule? Round to 1 decimal place.", edakit.Num(R, 0), edakit.Num(w, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Length times width over gate area must stay at or below %s: L <= %s x %s / %s = %s um. Beyond that, the router must hop to a higher layer near the gate (a jumper) or a diode must be added.", edakit.Num(R, 0), edakit.Num(R, 0), edakit.Num(ga, 4), edakit.Num(w, 2), ans),
			map[string]string{"l_max": ans})
	default:
		ls := []float64{float64(5 * edakit.Int(2, 20)), float64(5 * edakit.Int(4, 40)), float64(5 * edakit.Int(4, 60))}
		car := 0.0
		for _, l := range ls {
			car += l * w / ga
		}
		ans := edakit.Num(car, 1)
		q := gates + def + edakit.Sprintf("The cumulative antenna ratio sums the partial ratios of every layer up to the current one. On metals 1, 2 and 3 the gate-connected pieces are %s, %s and %s um long, all %s um wide, and diffusion is connected only on metal 4. What is the cumulative ratio at metal 3? Round to 1 decimal place.", edakit.Num(ls[0], 0), edakit.Num(ls[1], 0), edakit.Num(ls[2], 0), edakit.Num(w, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Total metal area %s um^2 over gate area %s um^2 gives %s. A jumper fixes it by moving the gate-side connection up a layer, which disconnects long lower-layer metal from the gate while those layers are etched.", edakit.Num((ls[0]+ls[1]+ls[2])*w, 3), edakit.Num(ga, 4), ans),
			map[string]string{"car": ans})
	}
}
