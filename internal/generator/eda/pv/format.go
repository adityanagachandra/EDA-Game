package pv

import (
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// padded is the byte length of a GDSII string field: odd lengths get one null pad byte.
func padded(n int) int { return n + n%2 }

var cellNames = []string{"INV", "NAND2", "VIA12", "BITCELL", "DFFR", "SRAM_COL", "TAP", "IO_PAD", "MUX4", "DECAP", "FILL1", "PLL_CORE", "ADC_BIAS"}

// gdsRecords: bytes of one element (KP 1), records in a structure (KP 2), bytes of a
// structure (KP 3).
func gdsRecords(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		if edakit.Int(0, 1) == 0 {
			n := 2 * edakit.Scaled(ctx, 2, 12)
			b := 32 + 8*n
			q := edakit.Sprintf("A GDSII BOUNDARY element is written as the records BOUNDARY, LAYER, DATATYPE, XY and ENDEL. Every record has a 4-byte header (2-byte length, record type, data type); LAYER and DATATYPE each carry one 2-byte integer, and XY carries 4-byte x and y integers for every point, with the first point repeated at the end to close the ring. How many bytes does a polygon with %d vertices occupy?", n)
			return edakit.Numeric(q, edakit.IntStr(b),
				edakit.Sprintf("BOUNDARY 4 + LAYER 6 + DATATYPE 6 + XY (4 + 8 x %d points) + ENDEL 4 = %d bytes. The XY record lists %d points because the ring is closed by repeating the first vertex.", n+1, b, n+1),
				map[string]string{"xy_bytes": edakit.IntStr(4 + 8*(n+1))})
		}
		m := edakit.Scaled(ctx, 2, 14)
		b := 32 + 8*m
		q := edakit.Sprintf("A GDSII PATH element (a wire given by its centerline) is written as PATH, LAYER, DATATYPE, WIDTH, XY and ENDEL, with no optional records. Every record has a 4-byte header; LAYER and DATATYPE carry one 2-byte integer each, WIDTH one 4-byte integer, and XY 8 bytes per centerline point (a path is not closed). How many bytes does a path with %d points occupy?", m)
		return edakit.Numeric(q, edakit.IntStr(b),
			edakit.Sprintf("PATH 4 + LAYER 6 + DATATYPE 6 + WIDTH 8 + XY (4 + 8 x %d) + ENDEL 4 = %d bytes.", m, b),
			map[string]string{"xy_bytes": edakit.IntStr(4 + 8*m)})
	case 2:
		b, s, t := edakit.Scaled(ctx, 2, 30), edakit.Int(0, 8), edakit.Int(0, 4)
		r := 3 + 5*b + 4*s + 6*t
		q := edakit.Sprintf("In GDSII a structure is BGNSTR, STRNAME, its elements, then ENDSTR. A BOUNDARY element is 5 records (BOUNDARY, LAYER, DATATYPE, XY, ENDEL), an SREF is 4 (SREF, SNAME, XY, ENDEL) and a TEXT label is 6 (TEXT, LAYER, TEXTTYPE, XY, STRING, ENDEL). A structure holds %d boundaries, %d SREFs and %d text labels, with no optional records. How many records does it contain?", b, s, t)
		return edakit.Numeric(q, edakit.IntStr(r),
			edakit.Sprintf("3 framing records + 5 x %d + 4 x %d + 6 x %d = %d records.", b, s, t, r),
			map[string]string{"records": edakit.IntStr(r)})
	default:
		name, sub := edakit.Pick(cellNames), edakit.Pick(cellNames)
		for sub == name {
			sub = edakit.Pick(cellNames)
		}
		b, s := edakit.Scaled(ctx, 2, 20), edakit.Int(1, 6)
		pn, ps := padded(len(name)), padded(len(sub))
		total := 28 + (4 + pn) + 64*b + s*(24+ps) + 4
		q := edakit.Sprintf("A GDSII structure named %s holds %d rectangles (each a 64-byte BOUNDARY element) and %d SREFs of the structure %s, with no optional records. BGNSTR is 28 bytes, STRNAME and SNAME are a 4-byte header plus the name padded to an even length, an SREF's own record and its ENDEL are 4 bytes each, its XY record is 12 bytes, and ENDSTR is 4 bytes. How many bytes does the structure occupy?", name, b, s, sub)
		return edakit.Numeric(q, edakit.IntStr(total),
			edakit.Sprintf("BGNSTR 28 + STRNAME %d + %d x 64 + %d x (4 + %d + 12 + 4) + ENDSTR 4 = %d bytes. Names of odd length gain one pad byte.", 4+pn, b, s, 4+ps, total),
			map[string]string{"sref_bytes": edakit.IntStr(24 + ps)})
	}
}

// oasisVarint: bytes of an OASIS unsigned integer.
func oasisVarint(ctx generator.GeneratorContext) generator.Problem {
	bits := edakit.Int(3, 30)
	if ctx.Difficulty < 0.5 {
		bits = edakit.Int(3, 16)
	}
	v := 1<<(bits-1) + edakit.Int(0, 1<<(bits-1)-1)
	nb := (bits + 6) / 7
	q := edakit.Sprintf("OASIS writes an unsigned integer in 7-bit groups, one group per byte, with the top bit of each byte marking whether another byte follows. GDSII writes every coordinate as a fixed 4-byte integer. How many bytes does OASIS need for the value %d (for example a coordinate delta in DBU)?", v)
	why := edakit.Sprintf("%d needs %d bits (it lies in [2^%d, 2^%d)), and each byte carries 7 payload bits, so ceil(%d / 7) = %d bytes. Small deltas between neighbouring vertices therefore cost one or two bytes instead of four.", v, bits, bits-1, bits, bits, nb)
	return numChoice(q, float64(nb), 0, "bytes", why, 4, float64((bits+7)/8), float64(nb+1), float64(imax(1, nb-1)))
}

// gdsFlatBytes: bytes of a regular array stored flat, per instance, or as one AREF.
func gdsFlatBytes(ctx generator.GeneratorContext) generator.Problem {
	c, r := edakit.Scaled(ctx, 2, 40), edakit.Int(2, 30)
	name := edakit.Pick(cellNames)
	p := padded(len(name))
	cell := 28 + 4 + p + 64 + 4
	aref := 48 + p
	hier := cell + aref
	srefs := cell + c*r*(24+p)
	flat := 64 * c * r
	q := edakit.Sprintf("A %d by %d array of identical via rectangles is stored in GDSII as one structure %s holding a single 64-byte rectangle, plus one AREF placing it. The structure costs BGNSTR 28 + STRNAME (4 + name padded to even) + 64 + ENDSTR 4 bytes; the AREF costs 4 (AREF) + SNAME (4 + padded name) + COLROW 8 + XY 28 + ENDEL 4 bytes. How many bytes do the structure and the AREF take together?", c, r, name)
	why := edakit.Sprintf("Structure: 28 + %d + 64 + 4 = %d bytes. AREF: 4 + %d + 8 + 28 + 4 = %d bytes. Together %d bytes, against %d bytes for %d flat rectangles and %d bytes for one SREF per instance.", 4+p, cell, 4+p, aref, hier, flat, c*r, srefs)
	return numChoice(q, float64(hier), 0, "bytes", why, float64(flat), float64(srefs), float64(aref))
}

// unitsGrid: stored integer for a coordinate (KP 1), coordinate range of 32-bit integers
// (KP 2), off-grid vertices (KP 3).
func unitsGrid(ctx generator.GeneratorContext) generator.Problem {
	type unit struct {
		perUM int
		nm    string
	}
	units := []unit{{1000, "1"}, {2000, "0.5"}, {4000, "0.25"}, {200, "5"}, {10000, "0.1"}}
	switch kp(ctx) {
	case 1:
		u := edakit.Pick(units)
		n := edakit.Scaled(ctx, 2000, 900000)
		x := float64(n) / float64(u.perUM)
		xs := strconv.FormatFloat(x, 'f', -1, 64)
		q := edakit.Sprintf("A layout's UNITS record sets the database unit to %s nm, and every coordinate is stored as an integer count of database units. A vertex sits at x = %s um. What integer is stored for its x coordinate?", u.nm, xs)
		return edakit.Numeric(q, edakit.IntStr(n),
			edakit.Sprintf("One um is %d database units of %s nm, so %s um is %s x %d = %d.", u.perUM, u.nm, xs, xs, u.perUM, n),
			map[string]string{"dbu_per_um": edakit.IntStr(u.perUM)})
	case 2:
		u := edakit.Pick(units)
		mm := 2147483647.0 / float64(u.perUM) / 1000
		ans := edakit.Num(mm, 2)
		q := edakit.Sprintf("GDSII stores coordinates as signed 32-bit integers in database units. With a database unit of %s nm, what is the largest coordinate that can be represented, in mm? Round to 2 decimal places.", u.nm)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The largest value is 2^31 - 1 = 2147483647 database units; times %s nm that is %s mm. A finer database unit shrinks the reachable die, which is one reason the unit is chosen once per technology.", u.nm, ans),
			map[string]string{"max_dbu": "2147483647"})
	default:
		g := edakit.Pick([]int{2, 5, 10})
		n := edakit.Scaled(ctx, 6, 12)
		var parts []string
		off := 0
		for i := 0; i < n; i++ {
			x, y := g*edakit.Int(10, 400), g*edakit.Int(10, 400)
			if edakit.Int(0, 3) == 0 {
				x += edakit.Int(1, g-1)
			}
			if edakit.Int(0, 4) == 0 {
				y += edakit.Int(1, g-1)
			}
			if x%g != 0 || y%g != 0 {
				off++
			}
			parts = append(parts, edakit.Sprintf("(%d, %d)", x, y))
		}
		q := edakit.Sprintf("The database unit is 1 nm and the foundry's manufacturing grid is %d nm. A layer's vertices, in nm, are %s. How many vertices does an off-grid check flag?", g, joinStrs(parts))
		return edakit.Numeric(q, edakit.IntStr(off),
			edakit.Sprintf("A vertex is on grid only if both coordinates are multiples of %d. Checking each pair, %d vertices have at least one coordinate that is not.", g, off),
			map[string]string{"off_grid": edakit.IntStr(off)})
	}
}

func joinStrs(xs []string) string {
	out := ""
	for i, x := range xs {
		if i > 0 {
			out += ", "
		}
		out += x
	}
	return out
}

// hierarchy: flattened shape count (KP 1), leaf placements in the flat view (KP 2), bytes
// saved by an AREF (KP 3).
func hierarchy(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	switch k {
	case 1, 2:
		s1, s2, s3 := edakit.Int(3, 40), edakit.Int(2, 30), edakit.Int(0, 20)
		n1, n2 := edakit.Scaled(ctx, 2, 16), edakit.Scaled(ctx, 2, 24)
		c, r := edakit.Int(2, 32), edakit.Int(2, 16)
		stem := edakit.Sprintf("Cell LEAF holds %d shapes. Cell MID holds %d shapes and %d SREFs of LEAF. Cell TOP holds %d shapes, %d SREFs of MID and one AREF of LEAF with %d columns and %d rows.", s1, s2, n1, s3, n2, c, r)
		if k == 1 {
			flat := s3 + n2*(s2+n1*s1) + c*r*s1
			q := stem + " How many shapes does TOP contain once fully flattened?"
			return edakit.Numeric(q, edakit.IntStr(flat),
				edakit.Sprintf("Flatten bottom-up: MID expands to %d + %d x %d = %d shapes. TOP = %d + %d x %d + %d x %d x %d = %d. The hierarchy stores only %d shapes.", s2, n1, s1, s2+n1*s1, s3, n2, s2+n1*s1, c, r, s1, flat, s1+s2+s3),
				map[string]string{"mid_flat": edakit.IntStr(s2 + n1*s1), "stored": edakit.IntStr(s1 + s2 + s3)})
		}
		pl := n2*n1 + c*r
		q := stem + " How many placements of LEAF appear in the flattened view of TOP?"
		return edakit.Numeric(q, edakit.IntStr(pl),
			edakit.Sprintf("Each placement is one path from TOP to LEAF in the cell DAG, weighted by instance counts: through MID there are %d x %d = %d, and the AREF adds %d x %d = %d, for %d.", n2, n1, n2*n1, c, r, c*r, pl),
			map[string]string{"via_mid": edakit.IntStr(n2 * n1), "aref": edakit.IntStr(c * r)})
	default:
		name := edakit.Pick(cellNames)
		c, r := edakit.Scaled(ctx, 2, 64), edakit.Int(2, 32)
		p := padded(len(name))
		saved := c*r*(24+p) - (48 + p)
		q := edakit.Sprintf("A layout places cell %s in a regular %d by %d grid. Written as one SREF per instance, each SREF costs 4 (SREF) + SNAME (4 + name padded to even) + XY 12 + ENDEL 4 bytes. Written as one AREF, it costs 4 (AREF) + SNAME + COLROW 8 + XY 28 + ENDEL 4 bytes. How many bytes does the AREF save?", name, c, r)
		return edakit.Numeric(q, edakit.IntStr(saved),
			edakit.Sprintf("SREFs: %d x %d = %d bytes. AREF: %d bytes. Saving: %d bytes.", c*r, 24+p, c*r*(24+p), 48+p, saved),
			map[string]string{"sref_each": edakit.IntStr(24 + p), "aref": edakit.IntStr(48 + p)})
	}
}
