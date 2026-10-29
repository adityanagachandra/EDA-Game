package place

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// netDelay estimates a two-pin net's Elmore delay from placement.
func netDelay(ctx generator.GeneratorContext) generator.Problem {
	r := edakit.Pick([]float64{0.5, 1, 1.5, 2})
	c := edakit.Pick([]float64{0.1, 0.15, 0.2})
	a := pt{10 * edakit.Int(0, 100), 10 * edakit.Int(0, 100)}
	b := pt{10 * edakit.Int(0, 100), 10 * edakit.Int(0, 100)}
	for try := 0; manhattan(a, b) < 100 && try < generator.MaxDrawAttempts; try++ {
		b = pt{10 * edakit.Int(0, 100), 10 * edakit.Int(0, 100)}
	}
	if manhattan(a, b) < 100 {
		a, b = pt{100, 200}, pt{600, 500}
	}
	L := float64(manhattan(a, b))
	rd := float64(50 * edakit.Int(2, 20))
	cl := float64(edakit.Int(2, 20))
	head := edakit.Sprintf("A driver at %s and its only sink at %s (um) are joined by a wire estimated at their Manhattan distance, with %s ohm/um and %s fF/um. ", fmtPt(a), fmtPt(b), edakit.Num(r, 1), edakit.Num(c, 2))
	switch ctx.KPOrdinal {
	case 1:
		fs := r * c * L * L / 2
		ans := edakit.Num(fs/1000, 2)
		return edakit.Numeric(head+"Ignoring the driver and the load, what is the wire's distributed Elmore delay $rcL^2/2$, in ps? Round to 2 decimal places.", ans,
			edakit.Sprintf("L = %s um, so the delay is %s x %s x %s^2 / 2 = %s ohm-fF = %s ps.", edakit.Num(L, 0), edakit.Num(r, 1), edakit.Num(c, 2), edakit.Num(L, 0), edakit.Num(fs, 0), ans), map[string]string{"L": edakit.Num(L, 0)})
	case 3:
		f := edakit.Pick([]float64{1.5, 2, 2.5, 3})
		ratio := f * f
		ans := edakit.Num(ratio, 2)
		return edakit.Numeric(head+edakit.Sprintf("A congestion-driven move pushes the sink so that the wire becomes %s times longer. Ignoring driver and load, by what factor does the wire's Elmore delay grow? Round to 2 decimal places.", edakit.Num(f, 1)), ans,
			edakit.Sprintf("The distributed delay scales with $L^2$: %s^2 = %s. This quadratic growth is why long nets are buffered rather than tolerated.", edakit.Num(f, 1), ans), map[string]string{"ratio": ans})
	}
	fs := rd*(c*L+cl) + r*L*(c*L/2+cl)
	ans := edakit.Num(fs/1000, 2)
	return edakit.Numeric(head+edakit.Sprintf("The driver resistance is %s ohm and the sink pin is %s fF. Using $R_d(cL + C_L) + rL(cL/2 + C_L)$, what is the Elmore delay from driver to sink, in ps? Round to 2 decimal places.", edakit.Num(rd, 0), edakit.Num(cl, 0)), ans,
		edakit.Sprintf("L = %s um. The driver term is %s ohm-fF and the wire term %s ohm-fF, %s ps in total.", edakit.Num(L, 0), edakit.Num(rd*(c*L+cl), 0), edakit.Num(r*L*(c*L/2+cl), 0), ans), map[string]string{"L": edakit.Num(L, 0)})
}

// timingWeights computes criticality-based net weights and their effect.
func timingWeights(ctx generator.GeneratorContext) generator.Problem {
	wns := -float64(edakit.Int(50, 300))
	alpha := edakit.Pick([]float64{1, 2, 4, 8})
	beta := edakit.Pick([]float64{1, 2, 3})
	crit := func(s float64) float64 { return math.Max(0, math.Min(1, -s/math.Abs(wns))) }
	weight := func(s float64) float64 { return 1 + alpha*math.Pow(crit(s), beta) }
	rule := edakit.Sprintf("A timing-driven placer weights each net by $w = 1 + %s c^{%s}$, where the criticality $c = \\min(1, \\max(0, -s/|WNS|))$ uses the net's worst slack s and the design's WNS of %s ps. ", edakit.Num(alpha, 0), edakit.Num(beta, 0), edakit.Num(wns, 0))
	switch ctx.KPOrdinal {
	case 2:
		n := edakit.Int(3, 4)
		total := 0.0
		var desc, parts []string
		for i := 0; i < n; i++ {
			s := float64(edakit.Int(int(wns), 100))
			h := float64(edakit.Int(10, 200))
			w := weight(s)
			total += w * h
			desc = append(desc, edakit.Sprintf("HPWL %s um with slack %s ps", edakit.Num(h, 0), edakit.Num(s, 0)))
			parts = append(parts, edakit.Sprintf("%s x %s", edakit.Num(w, 4), edakit.Num(h, 0)))
		}
		ans := edakit.Num(total, 2)
		return edakit.Numeric(rule+edakit.Sprintf("The nets are: %s. What weighted wirelength does the placer now minimise over these nets, in um? Round to 2 decimal places.", strings.Join(desc, "; ")), ans,
			edakit.Sprintf("Weights times HPWL: %s = %s um. Nets with non-negative slack keep weight 1.", strings.Join(parts, " + "), ans), map[string]string{"total": ans})
	case 3:
		sc := wns * edakit.Pick([]float64{0.5, 0.75, 1})
		wc := weight(sc)
		xa, xb := float64(edakit.Int(0, 30)), float64(edakit.Int(70, 100))
		x0 := (xa + xb) / 2
		x1 := (wc*xa + xb) / (wc + 1)
		ans := edakit.Num(x0-x1, 2)
		return edakit.Numeric(rule+edakit.Sprintf("A movable cell connects to a fixed pin at x = %s through a net with slack %s ps, and to a fixed pin at x = %s through a net with positive slack. In quadratic placement, how far does the cell move toward the critical pin once weights are applied, in um? Round to 2 decimal places.", edakit.Num(xa, 0), edakit.Num(sc, 1), edakit.Num(xb, 0)), ans,
			edakit.Sprintf("Unweighted, the cell sits at the midpoint %s. The critical net gets weight %s, so the cell moves to (%s x %s + %s)/(%s + 1) = %s, a shift of %s um.", edakit.Num(x0, 2), edakit.Num(wc, 4), edakit.Num(wc, 4), edakit.Num(xa, 0), edakit.Num(xb, 0), edakit.Num(wc, 4), edakit.Num(x1, 2), ans),
			map[string]string{"w": edakit.Num(wc, 4)})
	}
	s := float64(edakit.Int(int(wns), 50))
	w := weight(s)
	ans := edakit.Num(w, 4)
	return edakit.Numeric(rule+edakit.Sprintf("What weight does a net with slack %s ps receive? Round to 4 decimal places.", edakit.Num(s, 0)), ans,
		edakit.Sprintf("c = %s, so w = 1 + %s x %s^%s = %s.", edakit.Num(crit(s), 4), edakit.Num(alpha, 0), edakit.Num(crit(s), 4), edakit.Num(beta, 0), ans), map[string]string{"c": edakit.Num(crit(s), 4)})
}

// rudy estimates routing demand with RUDY and compares it with capacity.
func rudy(ctx generator.GeneratorContext) generator.Problem {
	B := edakit.Pick([]int{10, 20})
	n := edakit.Int(2, 3)
	nets := make([]rect, n)
	var desc []string
	for i := range nets {
		w, h := edakit.Int(4, 3*B), edakit.Int(4, 3*B)
		nets[i] = rect{edakit.Int(0, 2*B), edakit.Int(0, 2*B), w, h}
		desc = append(desc, edakit.Sprintf("a net whose pin bounding box spans x %d to %d and y %d to %d", nets[i].x, nets[i].x+w, nets[i].y, nets[i].y+h))
	}
	anchor := nets[edakit.Int(0, n-1)]
	bin := rect{anchor.x / B * B, anchor.y / B * B, B, B}
	dens := func(r rect) float64 { return float64(r.w+r.h) / float64(r.w*r.h) }
	head := "RUDY spreads each net's HPWL uniformly over its bounding box, so a net with box w x h has wire density $(w + h)/(wh)$ inside the box. "
	switch ctx.KPOrdinal {
	case 1:
		r := nets[0]
		ans := edakit.Num(dens(r), 4)
		return edakit.Numeric(head+edakit.Sprintf("Coordinates are in um. What is the RUDY density, in um of wire per um^2, of %s? Round to 4 decimal places.", desc[0]), ans,
			edakit.Sprintf("The box is %d x %d, so the density is (%d + %d)/(%d x %d) = %s. Small boxes get high density, which flags local pin crowding.", r.w, r.h, r.w, r.h, r.w, r.h, ans), map[string]string{"d": ans})
	}
	dem := 0.0
	var parts []string
	for _, r := range nets {
		ov := bin.overlapArea(r)
		dem += dens(r) * float64(ov)
		parts = append(parts, edakit.Sprintf("%s x %d", edakit.Num(dens(r), 4), ov))
	}
	q := head + edakit.Sprintf("Coordinates are in um. The design has %s. Consider the %d x %d um bin spanning x %d to %d and y %d to %d. ", strings.Join(desc, "; "), B, B, bin.x, bin.x+B, bin.y, bin.y+B)
	if ctx.KPOrdinal == 3 {
		layers := edakit.Int(2, 6)
		pitch := edakit.Pick([]float64{0.5, 1, 2})
		capa := float64(layers) * float64(B*B) / pitch
		d0 := math.Round(capa * float64(edakit.Int(40, 120)) / 100)
		ratio := (d0 + dem) / capa
		ans := edakit.Num(ratio, 4)
		return edakit.Numeric(q+edakit.Sprintf("The many other nets crossing this bin already contribute %s um of RUDY demand. The bin offers %d routing layers with a track pitch of %s um, so its capacity is layers x bin area / pitch um of wire. What is the bin's total demand-to-capacity ratio? Round to 4 decimal places.", edakit.Num(d0, 0), layers, edakit.Num(pitch, 1)), ans,
			edakit.Sprintf("The listed nets add density times overlap area: %s = %s um, for a total of %s um. Capacity is %d x %d / %s = %s um, so the ratio is %s; above 1 the bin is predicted to overflow.", strings.Join(parts, " + "), edakit.Num(dem, 3), edakit.Num(d0+dem, 3), layers, B*B, edakit.Num(pitch, 1), edakit.Num(capa, 1), ans),
			map[string]string{"demand": edakit.Num(dem, 3), "capacity": edakit.Num(capa, 1)})
	}
	ans := edakit.Num(dem, 2)
	return edakit.Numeric(q+"What is the RUDY routing demand in that bin, in um of wire? Round to 2 decimal places.", ans,
		edakit.Sprintf("Each net contributes its density times the area of its box inside the bin: %s = %s um.", strings.Join(parts, " + "), ans),
		map[string]string{"demand": ans})
}

// inflation bloats cells in congested bins and checks the area budget.
func inflation(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Pick([]float64{1, 1.5, 2})
	fmax := edakit.Pick([]float64{1.5, 2, 2.5})
	rule := edakit.Sprintf("A routability-driven placer inflates each cell by the factor $f = \\min(%s, r^{%s})$ when its bin's congestion ratio r (demand over capacity) exceeds 1, and leaves it alone otherwise. ", edakit.Num(fmax, 1), edakit.Num(k, 1))
	f := func(r float64) float64 {
		if r <= 1 {
			return 1
		}
		return math.Min(fmax, math.Pow(r, k))
	}
	n := edakit.Scaled(ctx, 3, 5)
	areas, rs := make([]float64, n), make([]float64, n)
	var desc []string
	tot, inf := 0.0, 0.0
	for i := range areas {
		areas[i] = float64(edakit.Int(1, 12))
		rs[i] = float64(edakit.Int(70, 160)) / 100
		tot += areas[i]
		inf += areas[i] * f(rs[i])
		desc = append(desc, edakit.Sprintf("%s um^2 in a bin with r = %s", edakit.Num(areas[i], 0), edakit.Num(rs[i], 2)))
	}
	switch ctx.KPOrdinal {
	case 1:
		a := areas[0]
		r := 1 + float64(edakit.Int(5, 80))/100
		ans := edakit.Num(a*f(r), 3)
		return edakit.Numeric(rule+edakit.Sprintf("A cell of area %s um^2 sits in a bin with r = %s. What is its inflated area, in um^2? Round to 3 decimal places.", edakit.Num(a, 0), edakit.Num(r, 2)), ans,
			edakit.Sprintf("f = min(%s, %s^%s) = %s, so the area becomes %s um^2.", edakit.Num(fmax, 1), edakit.Num(r, 2), edakit.Num(k, 1), edakit.Num(f(r), 4), ans), map[string]string{"f": edakit.Num(f(r), 4)})
	case 3:
		region := math.Ceil(tot / edakit.Pick([]float64{0.6, 0.7, 0.8}))
		u := 100 * inf / region
		ans := edakit.Num(u, 1)
		return edakit.Numeric(rule+edakit.Sprintf("A region of %s um^2 holds cells of %s. After inflation, what percentage of the region do the cells occupy? Round to 1 decimal place.", edakit.Num(region, 0), strings.Join(desc, ", ")), ans,
			edakit.Sprintf("The cells grow from %s to %s um^2, so utilization rises from %s%% to %s%%. If it passes the target density, the placer must cap inflation or spread cells further.", edakit.Num(tot, 0), edakit.Num(inf, 3), edakit.Num(100*tot/region, 1), ans),
			map[string]string{"inflated": edakit.Num(inf, 3)})
	}
	ans := edakit.Num(inf, 3)
	return edakit.Numeric(rule+edakit.Sprintf("The cells are %s. What is their total area after inflation, in um^2? Round to 3 decimal places.", strings.Join(desc, ", ")), ans,
		edakit.Sprintf("Only cells with r > 1 grow; summing area times factor gives %s um^2 against %s um^2 before.", ans, edakit.Num(tot, 0)),
		map[string]string{"before": edakit.Num(tot, 0)})
}

// buffering compares an unbuffered long wire with a midpoint repeater.
func buffering(ctx generator.GeneratorContext) generator.Problem {
	r := edakit.Pick([]float64{0.5, 1, 1.5, 2})
	c := edakit.Pick([]float64{0.1, 0.15, 0.2})
	L := float64(100 * edakit.Scaled(ctx, 10, 40))
	rd := float64(50 * edakit.Int(2, 8))
	cl := float64(edakit.Int(2, 10))
	rb := float64(50 * edakit.Int(2, 8))
	cb := float64(edakit.Int(2, 6))
	tb := float64(edakit.Int(10, 30))
	h := L / 2
	unbuf := (rd*(c*L+cl) + r*L*(c*L/2+cl)) / 1000
	s1 := (rd*(c*h+cb) + r*h*(c*h/2+cb)) / 1000
	s2 := tb + (rb*(c*h+cl)+r*h*(c*h/2+cl))/1000
	buf := s1 + s2
	head := edakit.Sprintf("A driver with resistance %s ohm drives a %s um wire (%s ohm/um, %s fF/um) into a %s fF sink. A repeater can be placed at the midpoint: input capacitance %s fF, output resistance %s ohm, intrinsic delay %s ps. Use the Elmore model $R(cL + C) + rL(cL/2 + C)$ for each driven segment, with 1 ohm x 1 fF = 0.001 ps. ",
		edakit.Num(rd, 0), edakit.Num(L, 0), edakit.Num(r, 1), edakit.Num(c, 2), edakit.Num(cl, 0), edakit.Num(cb, 0), edakit.Num(rb, 0), edakit.Num(tb, 0))
	switch ctx.KPOrdinal {
	case 1:
		ans := edakit.Num(unbuf, 2)
		return edakit.Numeric(head+"What is the unbuffered delay from driver to sink, in ps? Round to 2 decimal places.", ans,
			edakit.Sprintf("Driver term %s ps plus wire term %s ps gives %s ps.", edakit.Num(rd*(c*L+cl)/1000, 2), edakit.Num(r*L*(c*L/2+cl)/1000, 2), ans), map[string]string{"unbuf": ans})
	case 3:
		ans := edakit.Num(unbuf-buf, 2)
		return edakit.Numeric(head+"How many ps does the midpoint repeater save (negative if it makes the path slower)? Round to 2 decimal places.", ans,
			edakit.Sprintf("Unbuffered: %s ps. Buffered: %s ps for the first half plus %s ps for the repeater stage, %s ps. Saving: %s ps. The quadratic wire term halves twice over, which is what makes repeaters pay on long wires.", edakit.Num(unbuf, 2), edakit.Num(s1, 2), edakit.Num(s2, 2), edakit.Num(buf, 2), ans),
			map[string]string{"unbuf": edakit.Num(unbuf, 2), "buf": edakit.Num(buf, 2)})
	}
	ans := edakit.Num(buf, 2)
	return edakit.Numeric(head+"With the repeater inserted, what is the delay from driver to sink, in ps? Round to 2 decimal places.", ans,
		edakit.Sprintf("Stage 1: the driver charges half the wire and the repeater input, %s ps. Stage 2: the repeater's intrinsic delay plus its drive of the second half and the sink, %s ps. Total %s ps.", edakit.Num(s1, 2), edakit.Num(s2, 2), ans),
		map[string]string{"s1": edakit.Num(s1, 2), "s2": edakit.Num(s2, 2)})
}
