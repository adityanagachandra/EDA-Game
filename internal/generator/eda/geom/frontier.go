package geom

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// parallelTiling: tile count (KP 1), tiles that must see one shape given a halo (KP 2), halo
// overhead (KP 3).
func parallelTiling(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		T := edakit.Pick([]int{50, 100, 200, 250})
		W, H := edakit.Int(3*T/2, 40*T), edakit.Int(3*T/2, 40*T)
		nx, ny := (W+T-1)/T, (H+T-1)/T
		q := edakit.Sprintf("A parallel DRC engine cuts a %d by %d um die into %d by %d um tiles (the last row and column may be partial) and hands one tile to each worker. How many tiles are there?", W, H, T, T)
		return edakit.Numeric(q, edakit.IntStr(nx*ny),
			edakit.Sprintf("ceil(%d / %d) = %d columns and ceil(%d / %d) = %d rows give %d tiles.", W, T, nx, H, T, ny, nx*ny),
			map[string]string{"nx": edakit.IntStr(nx), "ny": edakit.IntStr(ny)})
	case 2:
		T := edakit.Pick([]int{20, 25, 40, 50})
		h := edakit.Int(1, T/4)
		n := 10
		x1, y1 := edakit.Int(T, 7*T), edakit.Int(T, 7*T)
		r := rect{x1, y1, x1 + edakit.Int(1, 2*T), y1 + edakit.Int(1, 2*T)}
		count := func(a1, a2 int) int {
			c := 0
			for i := 0; i < n; i++ {
				if a1 < (i+1)*T+h && a2 > i*T-h {
					c++
				}
			}
			return c
		}
		cx, cy := count(r.x1, r.x2), count(r.y1, r.y2)
		q := edakit.Sprintf("Tiles are half-open %d by %d um squares on a 10 by 10 grid starting at the origin. Each worker reads its tile expanded by a halo of %d um, the largest interaction distance. How many workers must read the shape $[%d, %d) \\times [%d, %d)$?", T, T, h, r.x1, r.x2, r.y1, r.y2)
		return edakit.Numeric(q, edakit.IntStr(cx*cy),
			edakit.Sprintf("Tile i covers [%di - %d, %d(i + 1) + %d) once expanded. The shape meets %d expanded columns and %d expanded rows: %d workers.", T, h, T, h, cx, cy, cx*cy),
			map[string]string{"cols": edakit.IntStr(cx), "rows": edakit.IntStr(cy)})
	default:
		T := 10 * edakit.Int(5, 50)
		h := edakit.Int(1, T/10)
		ov := 100 * (float64((T+2*h)*(T+2*h)) - float64(T*T)) / float64(T*T)
		ans := edakit.Num(ov, 1)
		q := edakit.Sprintf("Tiles are %d by %d um and each worker processes its tile plus a %d um halo on every side. By what percent does the processed area exceed the tile area? Round to 1 decimal place.", T, T, h)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("((%d + 2(%d))^2 - %d^2) / %d^2 = %s percent of redundant work. Smaller tiles balance load better but pay more halo; the halo is fixed by the rule deck, not by the tile size.", T, h, T, T, ans),
			map[string]string{"overhead_pct": ans})
	}
}

// differentialCheck: XOR area of two rectangles (KP 1), XOR area of two shape sets (KP 2), the XOR
// signature of a uniform shift (KP 3).
func differentialCheck(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		a, b := overlappingPair(20, 4, 12)
		x := a.area() + b.area() - 2*interArea(a, b)
		q := edakit.Sprintf("A regression compares a golden output shape %s with a new engine's output %s. What is the area of their XOR, the region the differential check reports?", a, b)
		return edakit.Numeric(q, edakit.IntStr(x),
			edakit.Sprintf("XOR = %d + %d - 2(%d) = %d. A correct change to a geometry kernel must leave this at zero unless the change was meant to move shapes.", a.area(), b.area(), interArea(a, b), x),
			map[string]string{"xor": edakit.IntStr(x)})
	case 2:
		k := edakit.Scaled(ctx, 2, 4)
		gold := randRects(k, 24, 3, 10)
		neu := append([]rect(nil), gold...)
		i := edakit.Int(0, k-1)
		r := neu[i]
		switch edakit.Int(0, 2) {
		case 0:
			r.x2 += edakit.Int(1, 3)
		case 1:
			r.y1 -= edakit.Int(1, 3)
		default:
			dd := edakit.Int(1, 2)
			r = rect{r.x1 + dd, r.y1, r.x2 + dd, r.y2}
		}
		neu[i] = r
		x := xorArea(gold, neu)
		q := edakit.Sprintf("The golden output layer is the union of %s. A modified Boolean engine outputs the union of %s. What is the XOR area of the two layers?", listRects("G", gold), listRects("N", neu))
		return edakit.Numeric(q, edakit.IntStr(x),
			edakit.Sprintf("Compare the unions cell by cell on the compressed grid: area covered by exactly one layer is %d. Overlaps with the other shapes can hide part of a change, so a raw per-shape diff and the XOR can disagree.", x),
			map[string]string{"xor": edakit.IntStr(x)})
	default:
		k := edakit.Int(2, 4)
		dlt := edakit.Int(1, 3)
		var rs, sh []rect
		x := 0
		for i := 0; i < k; i++ {
			w, h := edakit.Int(dlt+2, 12), edakit.Int(2, 12)
			r := rect{20 * i, 0, 20*i + w, h}
			rs = append(rs, r)
			sh = append(sh, rect{r.x1 + dlt, r.y1, r.x2 + dlt, r.y2})
			x += 2 * dlt * h
		}
		q := edakit.Sprintf("A new engine version shifts every output shape to the right by a distance of %d. The golden shapes are %s. What XOR area does the differential check report?", dlt, listRects("G", rs))
		return edakit.Numeric(q, edakit.IntStr(x),
			edakit.Sprintf("A rectangle of height h shifted by d less than its width leaves an XOR of two d-by-h slivers, 2dh. Summing over the shapes gives %d. A small shift produces a large, structured XOR, which is how a differential check catches an off-by-one in a coordinate transform.", x),
			map[string]string{"xor": edakit.IntStr(x)})
	}
}
