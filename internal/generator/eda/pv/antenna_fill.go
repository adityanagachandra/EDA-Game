package pv

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// gateArea draws a gate W by L in um and returns its description and area.
func gateArea() (string, float64) {
	w := edakit.Pick([]float64{0.42, 0.5, 0.65, 1, 1.5, 2})
	l := edakit.Pick([]float64{0.1, 0.13, 0.15, 0.18, 0.25})
	return edakit.Sprintf("W = %s um, L = %s um", edakit.Num(w, 2), edakit.Num(l, 2)), w * l
}

func floorEps(x float64) int { return int(math.Floor(x + 1e-9)) }

// antennaRules: area ratio (KP 1), longest allowed wire (KP 2), side-wall ratio (KP 3).
func antennaRules(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		ng := edakit.Scaled(ctx, 1, 3)
		var gs []string
		G := 0.0
		for i := 0; i < ng; i++ {
			d, a := gateArea()
			gs = append(gs, d)
			G += a
		}
		ns := edakit.Scaled(ctx, 1, 4)
		var segs []string
		M := 0.0
		for i := 0; i < ns; i++ {
			w := edakit.Pick([]float64{0.14, 0.2, 0.28, 0.4})
			l := float64(edakit.Int(5, 120))
			M += w * l
			segs = append(segs, edakit.Sprintf("%s by %s um", edakit.Num(l, 0), edakit.Num(w, 2)))
		}
		ans := edakit.Num(M/G, 1)
		q := edakit.Sprintf("When metal 2 is etched, a net's metal 2 pieces %s are connected to the gates %s and to no diffusion. Using the area antenna ratio (connected metal area over total gate area), what is the net's metal 2 antenna ratio? Round to 1 decimal place.", strings.Join(segs, ", "), strings.Join(gs, "; "))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Metal area = %s square um; gate area = %s square um (W x L summed over the gates). Ratio = %s. All gates on the net share the collected charge, so their areas add.", edakit.Num(M, 3), edakit.Num(G, 4), ans),
			map[string]string{"metal_area": edakit.Num(M, 3), "gate_area": edakit.Num(G, 4)})
	case 2:
		R := 50 * edakit.Int(4, 20)
		ng := edakit.Int(1, 3)
		G := 0.0
		var gs []string
		for i := 0; i < ng; i++ {
			d, a := gateArea()
			gs = append(gs, d)
			G += a
		}
		w := edakit.Pick([]float64{0.14, 0.2, 0.28, 0.4})
		L := floorEps(float64(R) * G / w)
		q := edakit.Sprintf("A metal 1 wire %s um wide connects only to gate(s) %s. The metal 1 area antenna ratio limit is %d. What is the longest wire, in whole um, that stays within the limit?", edakit.Num(w, 2), strings.Join(gs, "; "), R)
		return edakit.Numeric(q, edakit.IntStr(L),
			edakit.Sprintf("Need w x L / G <= %d with G = %s square um, so L <= %d x %s / %s = %s um; the longest whole length is %d um.", R, edakit.Num(G, 4), R, edakit.Num(G, 4), edakit.Num(w, 2), edakit.Num(float64(R)*G/w, 2), L),
			map[string]string{"gate_area": edakit.Num(G, 4)})
	default:
		_, G := gateArea()
		gd := edakit.Sprintf("%s square um", edakit.Num(G, 4))
		t := edakit.Pick([]float64{0.2, 0.3, 0.36, 0.5})
		w := edakit.Pick([]float64{0.14, 0.2, 0.28})
		l := float64(edakit.Int(10, 200))
		per := 2 * (l + w)
		ratio := per * t / G
		ans := edakit.Num(ratio, 1)
		q := edakit.Sprintf("A deck measures the antenna ratio by side-wall area: the wire's perimeter times the metal thickness, over the gate area. A %s by %s um metal wire, %s um thick, connects to a gate of area %s. What is the side-wall antenna ratio? Round to 1 decimal place.", edakit.Num(l, 0), edakit.Num(w, 2), edakit.Num(t, 2), gd)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Perimeter = 2 x (%s + %s) = %s um; side-wall area = %s x %s = %s square um; ratio = %s. Plasma charging during etch collects through exposed side walls, so thick layers are often checked this way.", edakit.Num(l, 0), edakit.Num(w, 2), edakit.Num(per, 2), edakit.Num(per, 2), edakit.Num(t, 2), edakit.Num(per*t, 3), ans),
			map[string]string{"perimeter": edakit.Num(per, 2)})
	}
}

const partialRule = "The partial ratio of layer i counts only the layer-i metal already connected to the gate when layer i is etched, that is, connected through layers 1 to i. The cumulative ratio at layer i is the sum of the partial ratios of layers 1 to i."

// antennaCumulative: one partial ratio (KP 1), cumulative ratio (KP 2), largest partial ratio
// (KP 3).
func antennaCumulative(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Scaled(ctx, 3, 5)
	G := float64(edakit.Int(5, 40)) / 100
	p := make([]float64, k+1)
	var parts []string
	for i := 1; i <= k; i++ {
		p[i] = float64(edakit.Int(2, 60))
		parts = append(parts, edakit.Sprintf("metal %d: %s square um connected through layers up to %d", i, edakit.Num(p[i], 0), i))
	}
	dl := edakit.Int(1, k-1)
	via := edakit.Int(dl+1, k)
	b := float64(edakit.Int(10, 80))
	parts = append(parts, edakit.Sprintf("metal %d: another %s square um that reaches the gate only through metal %d", dl, edakit.Num(b, 0), via))
	stem := edakit.Sprintf("A gate of area %s square um connects to this metal: %s. %s", edakit.Num(G, 2), strings.Join(parts, "; "), partialRule)
	switch kp(ctx) {
	case 1:
		i := edakit.Int(1, k)
		ans := edakit.Num(p[i]/G, 1)
		q := stem + edakit.Sprintf(" What is the partial antenna ratio of metal %d? Round to 1 decimal place.", i)
		why := edakit.Sprintf("Metal %d counts %s square um", i, edakit.Num(p[i], 0))
		if i == dl {
			why += edakit.Sprintf("; the other %s square um on metal %d is still floating when metal %d is etched, because it joins the gate only at metal %d", edakit.Num(b, 0), dl, dl, via)
		}
		return edakit.Numeric(q, ans, why+edakit.Sprintf(". Ratio = %s / %s = %s.", edakit.Num(p[i], 0), edakit.Num(G, 2), ans),
			map[string]string{"area": edakit.Num(p[i], 0)})
	case 2:
		sum := 0.0
		for i := 1; i <= k; i++ {
			sum += p[i]
		}
		ans := edakit.Num(sum/G, 1)
		q := stem + edakit.Sprintf(" What is the cumulative antenna ratio at metal %d? Round to 1 decimal place.", k)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Sum the partial ratios: (%s) / %s = %s / %s = %s. The %s square um on metal %d never counts, because when metal %d was etched it was not yet connected to the gate.", joinNums(p[1:]), edakit.Num(G, 2), edakit.Num(sum, 0), edakit.Num(G, 2), ans, edakit.Num(b, 0), dl, dl),
			map[string]string{"sum_area": edakit.Num(sum, 0)})
	default:
		mx := 0.0
		for i := 1; i <= k; i++ {
			mx = math.Max(mx, p[i])
		}
		ans := edakit.Num(mx/G, 1)
		q := stem + " What is the largest partial antenna ratio over all layers? Round to 1 decimal place."
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The largest connected area on one layer is %s square um, so the largest partial ratio is %s / %s = %s.", edakit.Num(mx, 0), edakit.Num(mx, 0), edakit.Num(G, 2), ans),
			map[string]string{"max_area": edakit.Num(mx, 0)})
	}
}

func joinNums(xs []float64) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = edakit.Num(x, 0)
	}
	return strings.Join(parts, " + ")
}

// antennaDiode: metal to remove (KP 1), diode cells needed (KP 2), diode cells on a net with
// several gates (KP 3).
func antennaDiode(ctx generator.GeneratorContext) generator.Problem {
	R0 := 50 * edakit.Int(4, 10)
	k := kp(ctx)
	switch k {
	case 1:
		_, G := gateArea()
		ratio := float64(R0) * (1 + float64(edakit.Int(10, 150))/100)
		M := ratio * G
		M = math.Round(M*10) / 10
		ex := M - float64(R0)*G
		ans := edakit.Num(ex, 2)
		q := edakit.Sprintf("A net has %s square um of metal 3 connected to a single gate of area %s square um. The metal 3 area ratio limit without protection is %d. How much metal 3 area, in square um, would have to be moved off the net's metal 3 to meet the limit? Round to 2 decimal places.", edakit.Num(M, 1), edakit.Num(G, 4), R0)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The allowed area is %d x %s = %s square um, so %s - %s = %s square um must go (in practice by a layer jump rather than by deleting wire).", R0, edakit.Num(G, 4), edakit.Num(float64(R0)*G, 3), edakit.Num(M, 1), edakit.Num(float64(R0)*G, 3), ans),
			map[string]string{"allowed": edakit.Num(float64(R0)*G, 3)})
	default:
		ng := 1
		if k == 3 {
			ng = edakit.Int(2, 4)
		}
		G := 0.0
		var gs []string
		for i := 0; i < ng; i++ {
			d, a := gateArea()
			G += a
			gs = append(gs, d)
		}
		kf := 100 * edakit.Int(2, 8)
		ac := edakit.Pick([]float64{0.1, 0.2, 0.25, 0.5})
		ratio := float64(R0) + float64(kf)*ac*float64(edakit.Int(1, 5)) - float64(edakit.Int(1, 20))
		M := math.Round(ratio*G*10) / 10
		ratio = M / G
		n := int(math.Ceil((ratio-float64(R0))/(float64(kf)*ac) - 1e-9))
		if n < 0 {
			n = 0
		}
		q := edakit.Sprintf("A net carries %s square um of metal connected to the gate(s) %s. The deck's limit is %d without protection, and it rises by %d for every square um of reverse-biased diode attached to the net (a linearized protection table). Each antenna diode cell adds %s square um of diode. How many diode cells must the router insert on the net?", edakit.Num(M, 1), strings.Join(gs, "; "), R0, kf, edakit.Num(ac, 2))
		return edakit.Numeric(q, edakit.IntStr(n),
			edakit.Sprintf("Gate area %s square um gives a ratio of %s. The excess over %d is %s, and each cell buys %d x %s = %s, so ceil(%s / %s) = %d cells. One diode protects every gate on the net, because they share the node.", edakit.Num(G, 4), edakit.Num(ratio, 1), R0, edakit.Num(ratio-float64(R0), 1), kf, edakit.Num(ac, 2), edakit.Num(float64(kf)*ac, 0), edakit.Num(ratio-float64(R0), 1), edakit.Num(float64(kf)*ac, 0), n),
			map[string]string{"ratio": edakit.Num(ratio, 1)})
	}
}

// antennaJumper: ratio after a jump (KP 1), farthest jump point (KP 2), nets fixable by
// jumping (KP 3).
func antennaJumper(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		_, G := gateArea()
		w := edakit.Pick([]float64{0.14, 0.2, 0.28})
		L := float64(edakit.Int(200, 900))
		s := float64(edakit.Int(2, 30))
		before, after := w*L/G, w*s/G
		ans := edakit.Num(after, 1)
		q := edakit.Sprintf("A %s um long metal 1 route, %s um wide, connects to a gate of area %s square um. The router breaks it %s um from the gate and carries the rest over metal 2, dropping back to metal 1 only at the far end. What is the metal 1 partial antenna ratio seen by the gate after the jump? Round to 1 decimal place.", edakit.Num(L, 0), edakit.Num(w, 2), edakit.Num(G, 4), edakit.Num(s, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Before the jump the gate saw %s; after it, only the %s um stub is connected when metal 1 is etched: %s x %s / %s = %s. The far metal 1 piece joins only through metal 2, too late to charge the gate.", edakit.Num(before, 1), edakit.Num(s, 0), edakit.Num(w, 2), edakit.Num(s, 0), edakit.Num(G, 4), ans),
			map[string]string{"before": edakit.Num(before, 1)})
	case 2:
		_, G := gateArea()
		R := 50 * edakit.Int(4, 16)
		w := edakit.Pick([]float64{0.14, 0.2, 0.28})
		pin := float64(edakit.Int(1, 8)) / 10
		smax := floorEps((float64(R)*G - pin) / w)
		q := edakit.Sprintf("A gate of area %s square um has a metal 1 pin shape of %s square um. A metal 1 route %s um wide leaves the pin, and the router will jump to metal 2 at some distance from the pin. The metal 1 ratio limit is %d. What is the farthest whole distance, in um, the route may run on metal 1 before jumping?", edakit.Num(G, 4), edakit.Num(pin, 1), edakit.Num(w, 2), R)
		return edakit.Numeric(q, edakit.IntStr(smax),
			edakit.Sprintf("Need (%s + %s s) / %s <= %d, so s <= (%d x %s - %s) / %s = %s um; the answer is %d um.", edakit.Num(pin, 1), edakit.Num(w, 2), edakit.Num(G, 4), R, R, edakit.Num(G, 4), edakit.Num(pin, 1), edakit.Num(w, 2), edakit.Num((float64(R)*G-pin)/w, 2), smax),
			map[string]string{"limit_area": edakit.Num(float64(R)*G, 3)})
	default:
		top := edakit.Int(5, 7)
		R := 50 * edakit.Int(6, 10)
		n := edakit.Scaled(ctx, 4, 8)
		var parts []string
		fix := 0
		for i := 0; i < n; i++ {
			layer := edakit.Int(2, top)
			up := edakit.Int(R/4, R)
			add := 10 * edakit.Int(2, 30)
			if layer < top && up+add <= R {
				fix++
			}
			if layer == top {
				parts = append(parts, edakit.Sprintf("N%d (violates on metal %d)", i+1, layer))
			} else {
				parts = append(parts, edakit.Sprintf("N%d (violates on metal %d; metal %d partial ratio now %d, jump adds %d)", i+1, layer, layer+1, up, add))
			}
		}
		q := edakit.Sprintf("An antenna repair step must choose between a layer jump and a diode for each violating net. A jump moves the violating length up one layer, which is possible only below the top layer (metal %d) and succeeds only if the upper layer's partial ratio stays within the limit of %d. The violating nets are %s. How many nets can be fixed by a jump alone?", top, R, strings.Join(parts, ", "))
		return edakit.Numeric(q, edakit.IntStr(fix),
			edakit.Sprintf("A net on metal %d has no layer above it and needs a diode; elsewhere the jump works when the upper ratio plus the added amount is at most %d. %d nets qualify.", top, R, fix),
			map[string]string{"jumpable": edakit.IntStr(fix)})
	}
}

// clip returns the parts of rs inside window w.
func clip(rs []rect, w rect) []rect {
	var out []rect
	for _, r := range rs {
		if c, ok := intersect(r, w); ok {
			out = append(out, c)
		}
	}
	return out
}

// fillDensity: density of one window (KP 1), windows below the minimum (KP 2), density from
// routing pitch around a blockage (KP 3).
func fillDensity(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		W := 10 * edakit.Int(5, 10)
		win := rect{0, 0, W, W}
		n := edakit.Scaled(ctx, 3, 6)
		var rs []rect
		for i := 0; i < n; i++ {
			r := randRect(W+20, 8, W/2)
			rs = append(rs, r.shift(-10, -10))
		}
		a := unionArea(clip(rs, win))
		ans := edakit.Num(100*float64(a)/float64(W*W), 1)
		q := edakit.Sprintf("A density window is $[0, %d] \\times [0, %d]$ um. The metal shapes near it are %s (um); some overlap each other and some extend past the window. What is the metal density of the window, in percent? Round to 1 decimal place.", W, W, listRects("R", rs))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Clip every shape to the window (AND), merge the clipped shapes (OR) so overlaps count once, and take the union area: %d square um out of %d, %s percent.", a, W*W, ans),
			map[string]string{"metal": edakit.IntStr(a)})
	case 2:
		W := 10 * edakit.Int(4, 8)
		D := 5 * edakit.Int(3, 8)
		n := edakit.Scaled(ctx, 4, 8)
		var rs []rect
		for i := 0; i < n; i++ {
			rs = append(rs, randRect(2*W, 5, W))
		}
		low := 0
		var ds []string
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				w := rect{i * W, j * W, (i + 1) * W, (j + 1) * W}
				a := unionArea(clip(rs, w))
				d := 100 * float64(a) / float64(W*W)
				ds = append(ds, edakit.Num(d, 1))
				if d < float64(D) {
					low++
				}
			}
		}
		q := edakit.Sprintf("A $%d \\times %d$ um region is checked as four $%d \\times %d$ um windows (a 2 by 2 tiling with corners at multiples of %d). The minimum density is %d percent. The metal shapes are %s (um). How many of the four windows are below the minimum?", 2*W, 2*W, W, W, W, D, listRects("R", rs))
		return edakit.Numeric(q, edakit.IntStr(low),
			edakit.Sprintf("The window densities, clipped and merged, are %s percent; %d are below %d percent.", strings.Join(ds, ", "), low, D),
			map[string]string{"low": edakit.IntStr(low)})
	default:
		W := 10 * edakit.Int(5, 20)
		bw, bh := 10*edakit.Int(1, W/20), 10*edakit.Int(1, W/20)
		w := edakit.Pick([]float64{0.1, 0.14, 0.2, 0.28})
		p := w * float64(edakit.Int(2, 5))
		d := 100 * (float64(W*W - bw*bh)) * (w / p) / float64(W*W)
		ans := edakit.Num(d, 1)
		q := edakit.Sprintf("A $%d \\times %d$ um window contains a $%d \\times %d$ um macro blockage with no metal on this layer. Everywhere else the layer is fully routed with parallel wires %s um wide at a %s um pitch. What is the window's metal density, in percent? Round to 1 decimal place.", W, W, bw, bh, edakit.Num(w, 2), edakit.Num(p, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Routed area %d - %d = %d square um at density %s / %s = %s, over %d square um: %s percent.", W*W, bw*bh, W*W-bw*bh, edakit.Num(w, 2), edakit.Num(p, 2), edakit.Num(w/p, 3), W*W, ans),
			map[string]string{"line_density": edakit.Num(w/p, 3)})
	}
}

// fillInsertion: squares needed (KP 1), legal fill slots (KP 2), largest reachable density (KP 3).
func fillInsertion(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	if k == 1 {
		W := 10 * edakit.Int(5, 20)
		D := 5 * edakit.Int(4, 10)
		Am := edakit.Int(W*W*D/400, W*W*D/100-10)
		f := edakit.Pick([]int{1, 2, 3, 4, 5})
		need := D*W*W/100 - Am
		n := ceilDiv(need, f*f)
		q := edakit.Sprintf("A $%d \\times %d$ um window holds %d square um of metal, and the minimum density is %d percent. Fill squares are %d by %d um and are placed only in empty space. How many fill squares are needed?", W, W, Am, D, f, f)
		return edakit.Numeric(q, edakit.IntStr(n),
			edakit.Sprintf("Target area %d percent of %d = %s square um; deficit %s; ceil(deficit / %d) = %d squares.", D, W*W, edakit.Num(float64(D*W*W)/100, 2), edakit.Num(float64(D*W*W)/100-float64(Am), 2), f*f, n),
			map[string]string{"deficit": edakit.Num(float64(D*W*W)/100-float64(Am), 2)})
	}
	f := edakit.Pick([]int{2, 3})
	qp := f + edakit.Int(2, 4)
	s := edakit.Int(1, 2)
	ns := edakit.Scaled(ctx, 5, 8)
	W := ns * qp
	var metal []rect
	for i := 0; i < edakit.Int(2, 4); i++ {
		metal = append(metal, randRect(W, 3, W/2))
	}
	slots := 0
	for i := 0; i < ns; i++ {
		for j := 0; j < ns; j++ {
			sq := rect{i * qp, j * qp, i*qp + f, j*qp + f}
			ok := true
			for _, m := range metal {
				if interArea(sq.grow(s), m) > 0 {
					ok = false
				}
			}
			if ok {
				slots++
			}
		}
	}
	stem := edakit.Sprintf("A $%d \\times %d$ um window is offered fill squares %d by %d um on a %d um grid, with lower-left corners at multiples of %d. A fill square is legal only if, grown by its %d um spacing to metal, it does not overlap any metal shape. The metal shapes are %s (um).", W, W, f, f, qp, qp, s, listRects("M", metal))
	if k == 2 {
		q := stem + " How many fill squares are legal?"
		return edakit.Numeric(q, edakit.IntStr(slots),
			edakit.Sprintf("Test each of the %d candidates: the free region is the window NOT (metal sized by %d), and %d squares fit inside it.", ns*ns, s, slots),
			map[string]string{"slots": edakit.IntStr(slots)})
	}
	am := unionArea(clip(metal, rect{0, 0, W, W}))
	d := 100 * float64(am+slots*f*f) / float64(W*W)
	ans := edakit.Num(d, 1)
	q := stem + " If every legal square is filled, what density does the window reach, in percent? Round to 1 decimal place."
	return edakit.Numeric(q, ans,
		edakit.Sprintf("Metal union inside the window is %d square um, and %d legal squares add %d. (%d + %d) / %d = %s percent. If this is still below the minimum, the window cannot be fixed by fill and the design itself must change.", am, slots, slots*f*f, am, slots*f*f, W*W, ans),
		map[string]string{"metal": edakit.IntStr(am), "slots": edakit.IntStr(slots)})
}

// fillWindows: number of window positions (KP 1), minimum window density from tile sums
// (KP 2), largest density gradient between neighbouring windows (KP 3).
func fillWindows(ctx generator.GeneratorContext) generator.Problem {
	switch kk := kp(ctx); kk {
	case 1:
		w := edakit.Pick([]int{20, 50, 100, 200})
		t := w / edakit.Pick([]int{1, 2, 4})
		nx, ny := edakit.Scaled(ctx, 5, 60), edakit.Int(5, 60)
		X, Y := w+(nx-1)*t, w+(ny-1)*t
		q := edakit.Sprintf("A density rule uses a %d by %d um window stepped by %d um in x and y over a %d by %d um die. How many window positions are checked?", w, w, t, X, Y)
		return edakit.Numeric(q, edakit.IntStr(nx*ny),
			edakit.Sprintf("Per axis (D - w) / t + 1: (%d - %d) / %d + 1 = %d in x and %d in y, so %d positions.", X, w, t, nx, ny, nx*ny),
			map[string]string{"nx": edakit.IntStr(nx), "ny": edakit.IntStr(ny)})
	default:
		n := edakit.Scaled(ctx, 4, 6)
		k := edakit.Pick([]int{2, 2, 3})
		if k >= n {
			k = 2
		}
		g := make([][]int, n)
		var rows []string
		for i := range g {
			g[i] = make([]int, n)
			var cells []string
			for j := range g[i] {
				g[i][j] = edakit.Int(5, 70)
				cells = append(cells, edakit.IntStr(g[i][j]))
			}
			rows = append(rows, "["+strings.Join(cells, " ")+"]")
		}
		m := n - k + 1
		d := make([][]float64, m)
		for i := 0; i < m; i++ {
			d[i] = make([]float64, m)
			for j := 0; j < m; j++ {
				s := 0
				for a := 0; a < k; a++ {
					for b := 0; b < k; b++ {
						s += g[i+a][j+b]
					}
				}
				d[i][j] = float64(s) / float64(k*k)
			}
		}
		stem := edakit.Sprintf("A layer is divided into a %d by %d grid of equal tiles, with metal densities in percent, row by row: %s. A density window covers %d by %d tiles and steps one tile at a time.", n, n, strings.Join(rows, ", "), k, k)
		if kk == 2 {
			mn := d[0][0]
			for i := range d {
				for j := range d[i] {
					mn = math.Min(mn, d[i][j])
				}
			}
			ans := edakit.Num(mn, 1)
			q := stem + " What is the lowest window density, in percent? Round to 1 decimal place."
			return edakit.Numeric(q, ans,
				edakit.Sprintf("Each window's density is the mean of its %d tiles. Over the %d window positions the minimum is %s percent. A 2D prefix sum makes each window an O(1) lookup.", k*k, m*m, ans),
				map[string]string{"windows": edakit.IntStr(m * m)})
		}
		mx := 0.0
		for i := 0; i < m; i++ {
			for j := 0; j < m; j++ {
				if i+1 < m {
					mx = math.Max(mx, math.Abs(d[i][j]-d[i+1][j]))
				}
				if j+1 < m {
					mx = math.Max(mx, math.Abs(d[i][j]-d[i][j+1]))
				}
			}
		}
		ans := edakit.Num(mx, 1)
		q := stem + " A gradient rule limits the density difference between windows one step apart. What is the largest difference between two such neighbouring windows, in percentage points? Round to 1 decimal place."
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Compute all %d window densities, then the absolute difference for every horizontally or vertically adjacent pair; the largest is %s points. Neighbouring windows share tiles, so the difference is set only by the tiles that enter and leave.", m*m, ans),
			map[string]string{"windows": edakit.IntStr(m * m)})
	}
}

// fillCoupling: sidewall coupling added by fill next to a signal wire.
func fillCoupling(ctx generator.GeneratorContext) generator.Problem {
	s := 10 * edakit.Int(4, 20)
	sf := 10 * edakit.Int(3, s/10)
	pct := 100 * float64(s) / float64(sf)
	q := edakit.Sprintf("A signal wire has a neighbour %d nm away on one side and empty space on the other. Fill is inserted on the empty side, %d nm from the wire. Using a sidewall model in which each side's coupling capacitance is proportional to 1 / spacing (same run length and thickness), by what percent does the wire's total sidewall capacitance increase?", s, sf)
	why := edakit.Sprintf("Before: k / %d. The fill adds k / %d, which is %d / %d = %s percent of the original. Fill-aware flows keep fill away from critical nets or extract it before signoff timing.", s, sf, s, sf, edakit.Num(pct, 0))
	return numChoice(q, pct, 0, "percent", why, 100*float64(sf)/float64(s), 100*float64(s-sf)/float64(s), pct/2)
}
