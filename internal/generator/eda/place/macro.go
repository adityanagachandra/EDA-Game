package place

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// macroChannel sizes the routing channel between two facing macros (computed multiple choice).
func macroChannel(ctx generator.GeneratorContext) generator.Problem {
	nets := edakit.Scaled(ctx, 40, 400)
	layers := edakit.Int(2, 4)
	pitch := edakit.Pick([]float64{0.14, 0.2, 0.28, 0.4})
	halo := edakit.Pick([]float64{2, 5, 10})
	tracks := int(math.Ceil(float64(nets) / float64(layers)))
	w := float64(tracks)*pitch + 2*halo
	q := edakit.Sprintf("Two macros face each other across a vertical channel. %d nets must cross the channel running vertically, and %d routing layers in the channel carry vertical wires at a track pitch of %s um. Each macro has a halo (keep-out ring) of %s um that standard cells and wires inside the channel must respect. What is the smallest channel width between the two macro edges that gives every net its own track?",
		nets, layers, edakit.Num(pitch, 2), edakit.Num(halo, 0))
	ans, d := numChoices(w, 2, " um",
		float64(tracks)*pitch,
		float64(nets)*pitch+2*halo,
		float64(tracks)*pitch+halo,
		float64(nets)/float64(layers)*pitch+2*halo)
	return edakit.Choice(q, ans, edakit.Sprintf("Each layer offers one wire per track, so %d nets over %d layers need ceil(%d/%d) = %d tracks, %s um of width. Adding a %s um halo on each side gives %s.", nets, layers, nets, layers, tracks, edakit.Num(float64(tracks)*pitch, 2), edakit.Num(halo, 0), ans), d...)
}

// pinAssignment assigns pins to boundary slots to minimise total Manhattan distance.
func pinAssignment(ctx generator.GeneratorContext) generator.Problem {
	n := 3
	if ctx.KPOrdinal != 1 {
		n = edakit.Scaled(ctx, 3, 4)
	}
	cells := make([]pt, n)
	slots := make([]pt, n)
	W := 100
	for i := range cells {
		cells[i] = pt{edakit.Int(10, 90), edakit.Int(10, 90)}
		switch edakit.Int(0, 3) {
		case 0:
			slots[i] = pt{0, edakit.Int(0, W)}
		case 1:
			slots[i] = pt{W, edakit.Int(0, W)}
		case 2:
			slots[i] = pt{edakit.Int(0, W), 0}
		default:
			slots[i] = pt{edakit.Int(0, W), W}
		}
	}
	var cd, sd []string
	for i := range cells {
		cd = append(cd, edakit.Sprintf("p%d connects to a cell at %s", i+1, fmtPt(cells[i])))
		sd = append(sd, edakit.Sprintf("s%d at %s", i+1, fmtPt(slots[i])))
	}
	cost := func(perm []int) int {
		t := 0
		for i, j := range perm {
			t += manhattan(cells[i], slots[j])
		}
		return t
	}
	id := make([]int, n)
	for i := range id {
		id[i] = i
	}
	best := cost(id)
	for _, p := range permutations(n) {
		if c := cost(p); c < best {
			best = c
		}
	}
	head := edakit.Sprintf("A 100 x 100 um block has %d pins to assign to %d legal boundary slots, one pin per slot. Pin %s. The slots are %s. The cost of an assignment is the total Manhattan distance from each pin's slot to the cell it connects to. ", n, n, strings.Join(cd, "; pin "), strings.Join(sd, ", "))
	switch ctx.KPOrdinal {
	case 1:
		c := cost(id)
		return edakit.Numeric(head+"What is the cost of the assignment p1 to s1, p2 to s2, p3 to s3, in um?", edakit.IntStr(c),
			edakit.Sprintf("Sum the three distances: the cost is %d um.", c), map[string]string{"cost": edakit.IntStr(c)})
	case 3:
		c := cost(id)
		return edakit.Numeric(head+"The current assignment puts each pin pi in slot si. By how many um does the optimal assignment reduce the cost?", edakit.IntStr(c-best),
			edakit.Sprintf("The current cost is %d um. Solving the assignment problem (here by trying every permutation; in a tool, by bipartite matching) gives %d um, a saving of %d um.", c, best, c-best),
			map[string]string{"current": edakit.IntStr(c), "best": edakit.IntStr(best)})
	}
	return edakit.Numeric(head+"What is the minimum possible cost, in um?", edakit.IntStr(best),
		edakit.Sprintf("This is a minimum-cost bipartite matching between pins and slots. Checking every assignment gives %d um.", best),
		map[string]string{"best": edakit.IntStr(best)})
}

// powerPlanning counts stripes, their routing cost and the rail IR drop.
func powerPlanning(ctx generator.GeneratorContext) generator.Problem {
	switch ctx.KPOrdinal {
	case 2:
		s := edakit.Pick([]float64{1, 1.5, 2, 3, 4})
		p := edakit.Pick([]float64{20, 30, 40, 50, 60, 80})
		f := 100 * 2 * s / p
		ans := edakit.Num(f, 1)
		return edakit.Numeric(edakit.Sprintf("A power grid uses vertical stripes on metal 5. Every %s um pitch carries one VDD stripe and one VSS stripe, each %s um wide. What percentage of metal 5's routing width do the stripes consume? Round to 1 decimal place.", edakit.Num(p, 0), edakit.Num(s, 1)), ans,
			edakit.Sprintf("Per pitch, 2 x %s = %s um of %s um is power, which is %s%%. That routing resource is gone before the router starts.", edakit.Num(s, 1), edakit.Num(2*s, 1), edakit.Num(p, 0), ans), map[string]string{"fraction": ans})
	case 3:
		L := float64(10 * edakit.Int(4, 30))
		r := edakit.Pick([]float64{0.05, 0.1, 0.2, 0.4})
		I := float64(edakit.Int(1, 20))
		R := r * L
		drop := I * R / 8
		ans := edakit.Num(drop, 3)
		return edakit.Numeric(edakit.Sprintf("A metal-1 power rail runs %s um between two power stripes that hold it at VDD at both ends. The rail resistance is %s ohm/um, and the cells along it draw a total of %s mA spread uniformly along its length. What is the largest IR drop along the rail, in mV? Round to 3 decimal places.", edakit.Num(L, 0), edakit.Num(r, 2), edakit.Num(I, 0)), ans,
			edakit.Sprintf("A uniformly loaded line fed from both ends drops most at its midpoint, by $IR/8$. Here R = %s ohm, so the drop is %s x %s / 8 = %s mV. Halving the stripe pitch quarters it.", edakit.Num(R, 2), edakit.Num(I, 0), edakit.Num(R, 2), ans),
			map[string]string{"R": edakit.Num(R, 2)})
	}
	W := float64(10 * edakit.Scaled(ctx, 20, 200))
	o := float64(edakit.Int(5, 30))
	p := float64(edakit.Pick([]int{20, 25, 40, 50, 60}))
	if math.Mod(W-o, p) == 0 {
		o++
	}
	n := int(math.Floor((W-o)/p)) + 1
	return edakit.Numeric(edakit.Sprintf("A core is %s um wide. The power planner places vertical stripe pairs starting %s um from the left core edge and repeating every %s um, and it keeps every stripe pair whose start lies inside the core. How many stripe pairs does it place?", edakit.Num(W, 0), edakit.Num(o, 0), edakit.Num(p, 0)), edakit.IntStr(n),
		edakit.Sprintf("Starts are at %s, %s, ... up to %s um: floor((%s - %s)/%s) + 1 = %d.", edakit.Num(o, 0), edakit.Num(o+p, 0), edakit.Num(o+float64(n-1)*p, 0), edakit.Num(W, 0), edakit.Num(o, 0), edakit.Num(p, 0), n),
		map[string]string{"pairs": edakit.IntStr(n)})
}
