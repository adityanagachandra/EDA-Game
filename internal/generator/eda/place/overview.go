package place

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// siteKinds are (site width, row height) pairs in nanometres, typical of public libraries.
var siteKinds = [][2]int{{190, 1400}, {460, 2720}, {200, 1600}, {380, 2800}, {54, 270}, {170, 1200}}

func um(nm int) string { return edakit.Num(float64(nm)/1000, 3) }

// rowsSites counts rows, sites and site utilization in a core region.
func rowsSites(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Pick(siteKinds)
	site, rowH := k[0], k[1]
	rows := edakit.Scaled(ctx, 20, 180)
	perRow := edakit.Scaled(ctx, 100, 900)
	coreH := rows*rowH + edakit.Int(0, rowH-1)
	coreW := perRow*site + edakit.Int(0, site-1)
	head := edakit.Sprintf("A standard-cell library has a site %s um wide and rows %s um tall. The core region is %s um wide and %s um tall, and rows are stacked from the bottom edge of the core with no gaps. ",
		um(site), um(rowH), um(coreW), um(coreH))
	switch ctx.KPOrdinal {
	case 1:
		return edakit.Numeric(head+"How many complete placement rows fit in the core?", edakit.IntStr(rows),
			edakit.Sprintf("Rows = floor(%s / %s) = %d. The leftover %s um at the top cannot hold a full row.", um(coreH), um(rowH), rows, um(coreH-rows*rowH)),
			map[string]string{"rows": edakit.IntStr(rows)})
	case 3:
		w := edakit.Int(2, 12)
		total := rows * perRow
		n := int(float64(total) / float64(w) * float64(edakit.Int(45, 85)) / 100)
		u := 100 * float64(n*w) / float64(total)
		ans := edakit.Num(u, 1)
		return edakit.Numeric(head+edakit.Sprintf("The netlist holds %d cells, each %d sites wide. What percentage of the core's sites do they occupy? Round to 1 decimal place.", n, w), ans,
			edakit.Sprintf("Rows = %d and sites per row = floor(%s / %s) = %d, so the core has %d sites. The cells need %d x %d = %d sites, which is %s%%.", rows, um(coreW), um(site), perRow, total, n, w, n*w, ans),
			map[string]string{"rows": edakit.IntStr(rows), "sites_per_row": edakit.IntStr(perRow), "total_sites": edakit.IntStr(total)})
	}
	total := rows * perRow
	return edakit.Numeric(head+"Each row is filled with whole sites from the left edge. How many sites does the whole core hold?", edakit.IntStr(total),
		edakit.Sprintf("Rows = floor(%s / %s) = %d. Sites per row = floor(%s / %s) = %d. Total = %d x %d = %d.", um(coreH), um(rowH), rows, um(coreW), um(site), perRow, rows, perRow, total),
		map[string]string{"rows": edakit.IntStr(rows), "sites_per_row": edakit.IntStr(perRow)})
}

// utilization sizes the core from cell area, target utilization and aspect ratio.
func utilization(ctx generator.GeneratorContext) generator.Problem {
	area := float64(1000 * edakit.Scaled(ctx, 20, 400))
	u := edakit.Pick([]float64{0.5, 0.55, 0.6, 0.65, 0.7, 0.75, 0.8})
	ar := edakit.Pick([]float64{0.5, 0.75, 1, 1.25, 1.5, 2})
	core := area / u
	w := math.Sqrt(core / ar)
	h := ar * w
	head := edakit.Sprintf("The standard cells of a block have a total area of %s um^2. The floorplan targets a core utilization of %s%%. ", edakit.Num(area, 0), edakit.Num(100*u, 0))
	switch ctx.KPOrdinal {
	case 1:
		ans := edakit.Num(core, 1)
		return edakit.Numeric(head+"What core area is needed, in um^2? Round to 1 decimal place.", ans,
			edakit.Sprintf("Utilization is cell area over core area, so the core is %s / %s = %s um^2.", edakit.Num(area, 0), edakit.Num(u, 2), ans), map[string]string{"core_area": ans})
	case 3:
		m := float64(edakit.Int(5, 40))
		die := (w + 2*m) * (h + 2*m)
		ans := edakit.Num(die, 1)
		return edakit.Numeric(head+edakit.Sprintf("The core has aspect ratio (height over width) %s, and the die adds an I/O margin of %s um on every side of the core. What is the die area, in um^2? Round to 1 decimal place.", edakit.Num(ar, 2), edakit.Num(m, 0)), ans,
			edakit.Sprintf("Core area %s um^2 with H/W = %s gives W = %s um and H = %s um. The die is (W + %s)(H + %s) = %s um^2.", edakit.Num(core, 1), edakit.Num(ar, 2), edakit.Num(w, 3), edakit.Num(h, 3), edakit.Num(2*m, 0), edakit.Num(2*m, 0), ans),
			map[string]string{"core_w": edakit.Num(w, 3), "core_h": edakit.Num(h, 3)})
	}
	ans := edakit.Num(w, 2)
	return edakit.Numeric(head+edakit.Sprintf("The core's aspect ratio (height over width) is %s. What is the core width, in um? Round to 2 decimal places.", edakit.Num(ar, 2)), ans,
		edakit.Sprintf("Core area = %s / %s = %s um^2. With H = %s W, the area is %s W^2, so W = sqrt(%s / %s) = %s um.", edakit.Num(area, 0), edakit.Num(u, 2), edakit.Num(core, 1), edakit.Num(ar, 2), edakit.Num(ar, 2), edakit.Num(core, 1), edakit.Num(ar, 2), ans),
		map[string]string{"core_area": edakit.Num(core, 1)})
}

// orientPin maps a pin offset inside a macro of size w x h through a DEF orientation.
func orientPin(o string, w, h, px, py int) (int, int) {
	switch o {
	case "S":
		return w - px, h - py
	case "FN":
		return w - px, py
	case "FS":
		return px, h - py
	}
	return px, py
}

// orientation places a LEF macro with a DEF orientation and asks for absolute pin coordinates.
func orientation(ctx generator.GeneratorContext) generator.Problem {
	dbu := edakit.Pick([]int{1000, 2000})
	scale := dbu / 1000 // DBU per nm-thousandth: coordinates below are drawn in nm, then scaled
	w := 10 * edakit.Int(40, 300)
	h := edakit.Pick([]int{1400, 1600, 2720, 2800})
	px := 10 * edakit.Int(2, w/10-2)
	py := 10 * edakit.Int(2, h/10-2)
	X := 10 * edakit.Int(100, 5000)
	Y := h * edakit.Int(0, 40)
	o := "N"
	switch ctx.KPOrdinal {
	case 2:
		o = edakit.Pick([]string{"FN", "FS"})
	case 3:
		o = edakit.Pick([]string{"S", "FN", "FS"})
	default:
		if ctx.KPOrdinal == 0 {
			o = edakit.Pick([]string{"N", "S", "FN", "FS"})
		}
	}
	ox, oy := orientPin(o, w, h, px, py)
	ax, ay := X+ox, Y+oy
	head := edakit.Sprintf("A LEF macro has SIZE %s BY %s (um) and a pin whose shape is centred at (%s, %s) um in the macro's own coordinates. The DEF file says UNITS DISTANCE MICRONS %d and places the instance with PLACED ( %d %d ) %s. In DEF the placement point is the lower-left corner of the instance after orientation; FN mirrors about the vertical axis, FS mirrors about the horizontal axis and S rotates by 180 degrees. ",
		um(w), um(h), um(px), um(py), dbu, X*scale, Y*scale, o)
	if ctx.KPOrdinal == 3 {
		other := pt{10 * edakit.Int(100, 8000), h*edakit.Int(0, 40) + 10*edakit.Int(0, h/10)}
		d := absI(ax-other.x) + absI(ay-other.y)
		ans := um(d)
		return edakit.Numeric(head+edakit.Sprintf("Another pin sits at (%s, %s) um in die coordinates. What is the Manhattan distance between the two pins, in um? Round to 3 decimal places.", um(other.x), um(other.y)), ans,
			edakit.Sprintf("The origin is (%s, %s) um. Under %s the offset (%s, %s) becomes (%s, %s), so the pin is at (%s, %s) um. Its distance to the other pin is %s um.", um(X), um(Y), o, um(px), um(py), um(ox), um(oy), um(ax), um(ay), ans),
			map[string]string{"pin_x": um(ax), "pin_y": um(ay)})
	}
	axis := edakit.Pick([]string{"x", "y"})
	if o == "FN" {
		axis = "x"
	} else if o == "FS" {
		axis = "y"
	}
	val := ax
	if axis == "y" {
		val = ay
	}
	ans := um(val)
	return edakit.Numeric(head+edakit.Sprintf("What is the pin's absolute %s coordinate on the die, in um? Round to 3 decimal places.", axis), ans,
		edakit.Sprintf("Divide by %d to get the origin (%s, %s) um. Under %s the pin offset (%s, %s) becomes (%s, %s). The absolute pin is at (%s, %s) um.", dbu, um(X), um(Y), o, um(px), um(py), um(ox), um(oy), um(ax), um(ay)),
		map[string]string{"pin_x": um(ax), "pin_y": um(ay)})
}
