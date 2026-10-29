package cts

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

const ln9 = 2.1972245773362196

// skewLatency: KP1 latency of one sink as a sum of stage delays, KP2 global skew over a sink
// set, KP3 local skew over sequentially adjacent pairs only.
func skewLatency(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		src := 10 * edakit.Int(0, 15)
		n := edakit.Scaled(ctx, 3, 6)
		var st []int
		total := src
		for i := 0; i < n; i++ {
			d := edakit.Int(25, 70)
			st = append(st, d)
			total += d
		}
		wire := edakit.Int(3, 20)
		total += wire
		q := edakit.Sprintf("A clock leaves its source with %d ps of source latency. On its way to flop FF7 it passes %d buffers with delays %s ps, and the last wire to the clock pin adds %d ps. What is the clock latency (insertion delay) at FF7, in ps?", src, n, ints(st), wire)
		return edakit.Numeric(q, edakit.IntStr(total),
			edakit.Sprintf("Latency is the sum of everything from the source to the clock pin: %d + %s + %d = %d ps.", src, strings.ReplaceAll(ints(st), ", ", " + "), wire, total),
			map[string]string{"latency": edakit.IntStr(total)})
	case 2:
		n := edakit.Scaled(ctx, 4, 8)
		lat := make([]int, n)
		for i := range lat {
			lat[i] = edakit.Int(380, 470)
		}
		sk := maxInt(lat) - minInt(lat)
		q := edakit.Sprintf("After CTS, the %d clock sinks of a block report latencies of %s ps. What is the global skew of the tree, in ps?", n, ints(lat))
		return edakit.Numeric(q, edakit.IntStr(sk),
			edakit.Sprintf("Global skew is the largest latency minus the smallest: %d - %d = %d ps.", maxInt(lat), minInt(lat), sk),
			map[string]string{"global_skew": edakit.IntStr(sk)})
	}
	names := []string{"A", "B", "C", "D", "E"}
	n := 5
	lat := make([]int, n)
	for i := range lat {
		lat[i] = edakit.Int(380, 470)
	}
	// Choose three related pairs.
	type pr struct{ a, b int }
	var all []pr
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				all = append(all, pr{i, j})
			}
		}
	}
	all = edakit.Shuffle(all)
	pairs := all[:3]
	best := -1
	var desc []string
	for _, p := range pairs {
		desc = append(desc, names[p.a]+" to "+names[p.b])
		if s := absInt(lat[p.b] - lat[p.a]); s > best {
			best = s
		}
	}
	var ls []string
	for i := 0; i < n; i++ {
		ls = append(ls, edakit.Sprintf("%s = %d", names[i], lat[i]))
	}
	glob := maxInt(lat) - minInt(lat)
	q := edakit.Sprintf("Five flops have clock latencies (in ps) %s. Timing paths exist only from %s. What is the worst local skew, the largest absolute latency difference over flop pairs joined by a timing path, in ps?", strings.Join(ls, ", "), strings.Join(desc, ", "))
	return edakit.Numeric(q, edakit.IntStr(best),
		edakit.Sprintf("Only sequentially adjacent pairs matter for timing. The largest difference over the three connected pairs is %d ps, while the global skew over all five flops would be %d ps.", best, glob),
		map[string]string{"local_skew": edakit.IntStr(best), "global_skew": edakit.IntStr(glob)})
}

// clockSlew uses the single-pole estimate t_10-90 = ln 9 * R * C.
func clockSlew(ctx generator.GeneratorContext) generator.Problem {
	r := float64(50 * edakit.Int(4, 16))
	switch kpOf(ctx) {
	case 1:
		c := float64(10 * edakit.Int(5, 30))
		s := ln9 * r * c / 1000
		q := edakit.Sprintf("A clock buffer with an output resistance of $%s\\ \\Omega$ drives a lumped load of %s fF. Treating the stage as a single RC pole, what is the 10-90%% output transition time, in ps? Round to 1 decimal place.", f(r, 0), f(c, 0))
		return edakit.Numeric(q, f(s, 1),
			edakit.Sprintf("A single pole crosses 10%% at $\\tau\\ln(10/9)$ and 90%% at $\\tau\\ln 10$, so the 10-90 time is $\\tau\\ln 9 = 2.197\\,RC$. With $RC = %s \\cdot %s$ ohm-fF $= %s$ ps the slew is %s ps.", f(r, 0), f(c, 0), f(r*c/1000, 2), f(s, 1)),
			map[string]string{"tau_ps": f(r*c/1000, 3)})
	case 2:
		lim := float64(10 * edakit.Int(6, 15))
		cmax := lim * 1000 / (ln9 * r)
		q := edakit.Sprintf("The clock transition limit is %s ps. A buffer with output resistance $%s\\ \\Omega$ drives a lumped load. Using the single-pole estimate (10-90%% time $= RC\\ln 9$), what is the largest load it may drive, in fF? Round to 1 decimal place.", f(lim, 0), f(r, 0))
		return edakit.Numeric(q, f(cmax, 1),
			edakit.Sprintf("Set $RC\\ln 9 = %s$ ps: $C = %s \\cdot 1000 / (2.197 \\cdot %s) = %s$ fF.", f(lim, 0), f(lim, 0), f(r, 0), f(cmax, 1)),
			map[string]string{"cmax_ff": f(cmax, 2)})
	}
	lim := float64(10 * edakit.Int(6, 15))
	cpin := pick(1.0, 3.0, 0.5)
	cw := pick(1.0, 4.0, 0.5)
	cmax := lim * 1000 / (ln9 * r)
	n := int(math.Floor(cmax / (cpin + cw)))
	if n < 1 { // unreachable for the drawn ranges; fall back to light sinks
		cpin, cw = 1, 1
		n = int(math.Floor(cmax / (cpin + cw)))
	}
	q := edakit.Sprintf("A clock buffer has output resistance $%s\\ \\Omega$ and the transition limit is %s ps (single-pole estimate, 10-90%% time $= RC\\ln 9$). Each sink adds %s fF of clock-pin capacitance and %s fF of wire. How many sinks can the buffer drive without violating the limit?", f(r, 0), f(lim, 0), f(cpin, 1), f(cw, 1))
	return edakit.Numeric(q, edakit.IntStr(n),
		edakit.Sprintf("The largest load is $%s/(2.197\\cdot%s)$ = %s fF. Each sink costs %s fF, so $\\lfloor %s / %s \\rfloor = %d$ sinks.", f(lim*1000, 0), f(r, 0), f(cmax, 2), f(cpin+cw, 1), f(cmax, 2), f(cpin+cw, 1), n),
		map[string]string{"cmax_ff": f(cmax, 2), "sinks": edakit.IntStr(n)})
}

// clockPower: P = C V^2 f for a clock net (activity 1), with pF x V^2 x GHz = mW.
func clockPower(ctx generator.GeneratorContext) generator.Problem {
	v := pick(0.6, 1.0, 0.05)
	fr := pick(0.5, 3.0, 0.25)
	switch kpOf(ctx) {
	case 1:
		c := float64(10 * edakit.Int(5, 40))
		p := c * v * v * fr
		q := edakit.Sprintf("A clock network has %s pF of total switched capacitance and runs at %s GHz from a %s V supply. What is its dynamic power, in mW? Round to 1 decimal place.", f(c, 0), f(fr, 2), f(v, 2))
		return edakit.Numeric(q, f(p, 1),
			edakit.Sprintf("A clock makes one rising transition every cycle, so its activity factor is 1: $P = C V^2 f = %s \\cdot %s^2 \\cdot %s = %s$ mW (pF times V squared times GHz gives mW).", f(c, 0), f(v, 2), f(fr, 2), f(p, 1)),
			map[string]string{"power_mw": f(p, 3)})
	case 2:
		wl := float64(5 * edakit.Int(10, 60))
		cw := pick(0.15, 0.25, 0.01)
		n := 1000 * edakit.Int(5, 40)
		cp := pick(0.5, 1.5, 0.1)
		nb := 10 * edakit.Int(5, 40)
		cb := pick(2, 6, 1)
		c := wl*cw + float64(n)*cp/1000 + float64(nb)*cb/1000
		p := c * v * v * fr
		q := edakit.Sprintf("A clock tree has %s mm of wire at %s pF/mm, drives %d flop clock pins of %s fF each, and contains %d buffers whose switched capacitance is %s fF each. At %s GHz and %s V, what is the clock power, in mW? Round to 1 decimal place.", f(wl, 0), f(cw, 2), n, f(cp, 1), nb, f(cb, 0), f(fr, 2), f(v, 2))
		return edakit.Numeric(q, f(p, 1),
			edakit.Sprintf("Total capacitance: %s pF of wire + %s pF of pins + %s pF of buffers = %s pF. Then $P = CV^2f = %s$ mW.", f(wl*cw, 2), f(float64(n)*cp/1000, 2), f(float64(nb)*cb/1000, 2), f(c, 3), f(p, 1)),
			map[string]string{"c_pf": f(c, 3), "power_mw": f(p, 3)})
	}
	cc := float64(10 * edakit.Int(5, 30))
	cd := float64(50 * edakit.Int(4, 20))
	a := pick(0.05, 0.25, 0.01)
	pc := cc * v * v * fr
	pd := a * cd * v * v * fr
	share := 100 * pc / (pc + pd)
	q := edakit.Sprintf("A design switches %s pF of clock capacitance (activity 1) and %s pF of data capacitance with an average activity factor of %s (0-to-1 transitions per cycle). What share of the dynamic power does the clock take, in percent? Round to 1 decimal place.", f(cc, 0), f(cd, 0), f(a, 2))
	return edakit.Numeric(q, f(share, 1),
		edakit.Sprintf("$V^2 f$ is common to both and cancels. Clock: %s pF. Data: $%s \\cdot %s = %s$ pF effective. Share = %s / %s = %s%%.", f(cc, 0), f(a, 2), f(cd, 0), f(a*cd, 2), f(cc, 0), f(cc+a*cd, 2), f(share, 1)),
		map[string]string{"share": f(share, 2)})
}

// hTree: level-1 H has a horizontal bar D/2 and two vertical bars D/2 (1.5 D in all); each
// level halves the size and quadruples the count.
func hTree(ctx generator.GeneratorContext) generator.Problem {
	d := float64(1000 * edakit.Int(2, 12))
	geom := edakit.Sprintf("An H-tree serves a square die of side $D = %s$ µm. The level-1 H has a horizontal bar of length $D/2$ whose ends feed two vertical bars, each of length $D/2$. Each tip of a level-$k$ H is the centre of a level-$(k+1)$ H of half the size.", f(d, 0))
	switch kpOf(ctx) {
	case 1:
		k := edakit.Scaled(ctx, 2, 5)
		wl := 1.5 * d * math.Pow(2, float64(k-1))
		q := geom + edakit.Sprintf(" What is the total wirelength of all the level-%d H's together, in µm?", k)
		return edakit.Numeric(q, f(wl, 1),
			edakit.Sprintf("Level %d has $4^{%d} = %d$ H's, each $1.5D/2^{%d}$ long. The total is $1.5D\\cdot 2^{%d} = %s$ µm: wirelength doubles every level because the count quadruples while the size halves.", k, k-1, int(math.Pow(4, float64(k-1))), k-1, k-1, f(wl, 1)),
			map[string]string{"level_wl": f(wl, 1)})
	case 2:
		n := edakit.Scaled(ctx, 2, 5)
		wl := 1.5 * d * (math.Pow(2, float64(n)) - 1)
		q := geom + edakit.Sprintf(" The tree has %d levels and so drives $4^{%d} = %d$ sinks. What is its total wirelength, in µm?", n, n, int(math.Pow(4, float64(n))))
		return edakit.Numeric(q, f(wl, 1),
			edakit.Sprintf("Level $k$ contributes $1.5D\\cdot2^{k-1}$. Summing the geometric series over $k = 1..%d$ gives $1.5D(2^{%d}-1) = %s$ µm.", n, n, f(wl, 1)),
			map[string]string{"total_wl": f(wl, 1)})
	}
	n := edakit.Scaled(ctx, 2, 6)
	pl := d * (1 - math.Pow(2, -float64(n)))
	q := geom + edakit.Sprintf(" The tree has %d levels. What is the wire path length from the centre of the die to any one sink, in µm? Round to 1 decimal place.", n)
	return edakit.Numeric(q, f(pl, 1),
		edakit.Sprintf("Inside a level-$k$ H the path runs half a horizontal bar and half a vertical bar, $D/2^{k+1} + D/2^{k+1} = D/2^k$. Summing over %d levels gives $D(1 - 2^{-%d}) = %s$ µm, the same for every sink, which is why the H-tree has zero nominal skew.", n, n, f(pl, 1)),
		map[string]string{"path": f(pl, 2)})
}

// mesh: a grid of m horizontal and m vertical wires spanning a D x D region.
func mesh(ctx generator.GeneratorContext) generator.Problem {
	d := pick(1, 4, 0.5)
	m := 4 * edakit.Int(4, 16)
	switch kpOf(ctx) {
	case 1:
		wl := 2 * float64(m) * d
		q := edakit.Sprintf("A clock mesh covers a %s mm by %s mm region with %d horizontal and %d vertical wires, each spanning the full region. What is the mesh wirelength, in mm?", f(d, 1), f(d, 1), m, m)
		return edakit.Numeric(q, f(wl, 1),
			edakit.Sprintf("Each direction has %d wires of %s mm: $2 \\cdot %d \\cdot %s = %s$ mm.", m, f(d, 1), m, f(d, 1), f(wl, 1)),
			map[string]string{"mesh_wl": f(wl, 2)})
	case 2:
		cw := pick(0.15, 0.3, 0.01)
		n := 1000 * edakit.Int(5, 30)
		cp := pick(0.5, 1.5, 0.1)
		v := pick(0.7, 1.0, 0.05)
		fr := pick(1, 3, 0.5)
		c := 2*float64(m)*d*cw + float64(n)*cp/1000
		p := c * v * v * fr
		q := edakit.Sprintf("A clock mesh of %d horizontal and %d vertical wires spans a %s mm square; the wire has %s pF/mm. It drives %d clock pins of %s fF. At %s GHz and %s V, what is the power to switch the mesh and its pins, in mW? Round to 1 decimal place.", m, m, f(d, 1), f(cw, 2), n, f(cp, 1), f(fr, 1), f(v, 2))
		return edakit.Numeric(q, f(p, 1),
			edakit.Sprintf("Mesh wire: $2\\cdot%d\\cdot%s\\cdot%s = %s$ pF; pins: %s pF; total %s pF. $P = CV^2f = %s$ mW.", m, f(d, 1), f(cw, 2), f(2*float64(m)*d*cw, 2), f(float64(n)*cp/1000, 2), f(c, 2), f(p, 1)),
			map[string]string{"c_pf": f(c, 3)})
	}
	n := edakit.Int(3, 4)
	for 2*float64(m) < 1.5*(math.Pow(2, float64(n))-1)+4 {
		m += 4
	}
	cw := pick(0.15, 0.3, 0.01)
	v := pick(0.7, 1.0, 0.05)
	fr := pick(1, 3, 0.5)
	meshWL := 2 * float64(m) * d
	htWL := 1.5 * d * (math.Pow(2, float64(n)) - 1)
	extra := (meshWL - htWL) * cw * v * v * fr
	q := edakit.Sprintf("On a %s mm square region, compare a clock mesh of %d horizontal and %d vertical full-span wires with a %d-level H-tree (level-1 H of total length $1.5D$, wirelength doubling each level). With %s pF/mm of wire, %s GHz and %s V, how much more wire power does the mesh burn than the H-tree, in mW? Round to 2 decimal places.", f(d, 1), m, m, n, f(cw, 2), f(fr, 1), f(v, 2))
	return edakit.Numeric(q, f(extra, 2),
		edakit.Sprintf("Mesh: %s mm. H-tree: $1.5D(2^{%d}-1) = %s$ mm. The difference, %s mm, costs $%s \\cdot %s \\cdot %s^2 \\cdot %s = %s$ mW.", f(meshWL, 2), n, f(htWL, 2), f(meshWL-htWL, 2), f(meshWL-htWL, 2), f(cw, 2), f(v, 2), f(fr, 1), f(extra, 2)),
		map[string]string{"mesh_wl": f(meshWL, 2), "htree_wl": f(htWL, 2)})
}

// treeSkew: Elmore delays in a small unbuffered clock tree and the dummy load that zeroes skew.
func treeSkew(ctx generator.GeneratorContext) generator.Problem {
	k := kpOf(ctx)
	if k == 3 {
		r0, ct := 50*edakit.Int(2, 8), 5*edakit.Int(1, 6)
		ra, ca := 50*edakit.Int(1, 8), 5*edakit.Int(2, 10)
		rb, cb := 50*edakit.Int(1, 8), 5*edakit.Int(2, 10)
		da := r0*(ct+ca+cb) + ra*ca
		db := r0*(ct+ca+cb) + rb*cb
		if da == db { // break the tie: a heavier sink b changes the skew by rb * 5 fF
			cb += 5
			da = r0*(ct+ca+cb) + ra*ca
			db = r0*(ct+ca+cb) + rb*cb
		}
		fast, rfast := "b", rb
		if da < db {
			fast, rfast = "a", ra
		}
		sk := absInt(da - db)
		dc := float64(sk) / float64(rfast)
		q := edakit.Sprintf("A clock driver feeds a trunk of $%d\\ \\Omega$ into node t ($C_t = %d$ fF). Branch a ($%d\\ \\Omega$) ends at sink a with %d fF; branch b ($%d\\ \\Omega$) ends at sink b with %d fF. Using Elmore delay, how much dummy capacitance must be added at the faster sink to make the skew zero, in fF? Round to 2 decimal places.", r0, ct, ra, ca, rb, cb)
		return edakit.Numeric(q, f(dc, 2),
			edakit.Sprintf("Elmore delays: sink a %s ps, sink b %s ps, so sink %s is faster by %s ps. Dummy capacitance $\\Delta C$ at sink %s raises both delays by $%d\\,\\Delta C$ through the shared trunk, and sink %s's delay by a further $%d\\,\\Delta C$ through its own branch. The skew therefore closes by $%d\\,\\Delta C$, so $\\Delta C = %d/%d = %s$ fF.", f(float64(da)/1000, 3), f(float64(db)/1000, 3), fast, f(float64(sk)/1000, 3), fast, r0, fast, rfast, rfast, sk, rfast, f(dc, 2)),
			map[string]string{"skew_ps": f(float64(sk)/1000, 3), "dummy_ff": f(dc, 3)})
	}
	// Two-level tree: root -> u (R1,C1), root -> w (R2,C2); u -> s1, s2; w -> s3 (and s4 at higher difficulty).
	nw := 1
	if ctx.Difficulty > 0.6 {
		nw = 2
	}
	R := map[string]int{}
	C := map[string]int{}
	for _, n := range []string{"u", "w"} {
		R[n], C[n] = 50*edakit.Int(1, 6), 5*edakit.Int(1, 4)
	}
	sinks := []string{"s1", "s2", "s3"}
	if nw == 2 {
		sinks = append(sinks, "s4")
	}
	parent := map[string]string{"s1": "u", "s2": "u", "s3": "w", "s4": "w"}
	for _, s := range sinks {
		R[s], C[s] = 50*edakit.Int(1, 8), 5*edakit.Int(2, 8)
	}
	down := map[string]int{}
	for _, s := range sinks {
		down[s] = C[s]
	}
	down["u"] = C["u"] + down["s1"] + down["s2"]
	down["w"] = C["w"] + down["s3"]
	if nw == 2 {
		down["w"] += down["s4"]
	}
	t := map[string]int{}
	t["u"] = R["u"] * down["u"]
	t["w"] = R["w"] * down["w"]
	var ds []int
	for _, s := range sinks {
		t[s] = t[parent[s]] + R[s]*C[s]
		ds = append(ds, t[s])
	}
	var parts []string
	parts = append(parts, edakit.Sprintf("root to u: $%d\\ \\Omega$, $C_u = %d$ fF", R["u"], C["u"]))
	parts = append(parts, edakit.Sprintf("root to w: $%d\\ \\Omega$, $C_w = %d$ fF", R["w"], C["w"]))
	for _, s := range sinks {
		parts = append(parts, edakit.Sprintf("%s to %s: $%d\\ \\Omega$, sink %d fF", parent[s], s, R[s], C[s]))
	}
	desc := "An ideal clock root drives an RC tree with edges " + strings.Join(parts, "; ") + "."
	if k == 1 {
		s := edakit.Pick(sinks)
		ans := float64(t[s]) / 1000
		return edakit.Numeric(desc+edakit.Sprintf(" What is the Elmore delay from the root to sink %s, in ps? Round to 3 decimal places.", s), f(ans, 3),
			edakit.Sprintf("Downstream capacitance: u carries %d fF and w carries %d fF. Delay to %s: $R_{%s} C_{down}(%s) + R_{%s} C_{%s}$ = %s ps.", down["u"], down["w"], s, parent[s], parent[s], s, s, f(ans, 3)),
			map[string]string{"delay": f(ans, 3)})
	}
	sk := float64(maxInt(ds)-minInt(ds)) / 1000
	var dl []string
	for i, s := range sinks {
		dl = append(dl, s+" "+f(float64(ds[i])/1000, 3))
	}
	return edakit.Numeric(desc+" What is the skew of this tree, the largest minus the smallest root-to-sink Elmore delay, in ps? Round to 3 decimal places.", f(sk, 3),
		edakit.Sprintf("Downstream capacitance at u is %d fF and at w is %d fF. Sink delays (ps): %s. Skew = %s ps.", down["u"], down["w"], strings.Join(dl, ", "), f(sk, 3)),
		map[string]string{"skew": f(sk, 3)})
}
