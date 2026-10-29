package cmos

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// cellLayout: track height, cell area in contacted poly pitches, and tracks from height.
func cellLayout(ctx generator.GeneratorContext) generator.Problem {
	tracks := edakit.Pick([]float64{6, 6.5, 7.5, 9, 10.5, 12})
	pitch := float64(edakit.Pick([]int{24, 28, 30, 32, 36, 40, 48, 64, 80}))
	h := tracks * pitch
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(h, 1)
		q := edakit.Sprintf("A library is built on a %s-track cell height, where a track is the horizontal metal routing pitch of %s nm. What is the cell height in nm?", edakit.Num(tracks, 1), edakit.Num(pitch, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("Cell height is tracks times pitch: $%s \\times %s = %s$ nm.", edakit.Num(tracks, 1), edakit.Num(pitch, 0), ans), map[string]string{"height_nm": ans})
	case 2:
		cpp := float64(edakit.Pick([]int{42, 45, 48, 51, 54, 57, 64, 78, 90}))
		w := edakit.Scaled(ctx, 3, 24)
		area := float64(w) * cpp * h / 1e6
		ans := edakit.Num(area, 4)
		q := edakit.Sprintf("A cell is %d contacted poly pitches wide, with a CPP of %s nm. The library is %s tracks tall on a %s nm metal pitch. What is the cell area in square micrometres? Round to 4 decimal places.", w, edakit.Num(cpp, 0), edakit.Num(tracks, 1), edakit.Num(pitch, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Width is $%d \\times %s = %s$ nm and height is $%s \\times %s = %s$ nm, so the area is %s square micrometres.", w, edakit.Num(cpp, 0), edakit.Num(float64(w)*cpp, 0), edakit.Num(tracks, 1), edakit.Num(pitch, 0), edakit.Num(h, 1), ans),
			map[string]string{"width_nm": edakit.Num(float64(w)*cpp, 0), "height_nm": edakit.Num(h, 1)})
	}
	ans := edakit.Num(tracks, 1)
	q := edakit.Sprintf("A cell library's row height is %s nm and its horizontal routing pitch is %s nm. How many tracks tall is the library? Round to 1 decimal place.", edakit.Num(h, 1), edakit.Num(pitch, 0))
	return edakit.Numeric(q, ans, edakit.Sprintf("Tracks are height over pitch: $%s / %s = %s$.", edakit.Num(h, 1), edakit.Num(pitch, 0), ans), nil)
}

// rowArithmetic: rows in a core, placement sites, utilization and shared rails.
func rowArithmetic(ctx generator.GeneratorContext) generator.Problem {
	rowH := edakit.Pick([]int{216, 240, 270, 288, 320, 480, 540})
	site := edakit.Pick([]int{48, 51, 54, 57, 64, 90})
	coreH := rowH*edakit.Scaled(ctx, 20, 400) + edakit.Int(0, rowH-1)
	coreW := site*edakit.Scaled(ctx, 200, 4000) + edakit.Int(0, site-1)
	rows := coreH / rowH
	sitesPer := coreW / site
	switch kpOf(ctx) {
	case 1:
		q := edakit.Sprintf("A core area is %d nm tall and the library's row height is %d nm. Rows are stacked from the bottom edge with no gaps. How many complete placement rows fit?", coreH, rowH)
		return edakit.Numeric(q, edakit.IntStr(rows), edakit.Sprintf("$\\lfloor %d / %d \\rfloor = %d$ rows.", coreH, rowH, rows), map[string]string{"rows": edakit.IntStr(rows)})
	case 2:
		tot := rows * sitesPer
		q := edakit.Sprintf("A core is %d nm tall and %d nm wide. Rows are %d nm tall and each row is divided into placement sites %d nm wide, starting at the left edge. How many placement sites does the core hold in total?", coreH, coreW, rowH, site)
		return edakit.Numeric(q, edakit.IntStr(tot),
			edakit.Sprintf("There are $\\lfloor %d/%d \\rfloor = %d$ rows of $\\lfloor %d/%d \\rfloor = %d$ sites, so %d sites.", coreH, rowH, rows, coreW, site, sitesPer, tot),
			map[string]string{"rows": edakit.IntStr(rows), "sites_per_row": edakit.IntStr(sitesPer)})
	}
	if edakit.Int(0, 2) == 0 {
		q := edakit.Sprintf("A core holds %d placement rows. Alternate rows are flipped so that neighbouring rows share a power or ground rail along their common edge. How many horizontal rails (VDD plus VSS) run across the core?", rows)
		return edakit.Numeric(q, edakit.IntStr(rows+1),
			edakit.Sprintf("Every boundary between rows carries one shared rail, and the top and bottom edges carry one each: %d + 1 = %d rails.", rows, rows+1), nil)
	}
	used := rows * sitesPer * edakit.Int(45, 85) / 100
	util := 100 * float64(used) / float64(rows*sitesPer)
	ans := edakit.Num(util, 2)
	q := edakit.Sprintf("A core has %d rows of %d sites each. The placed cells occupy %d sites in total. What is the placement utilization, in percent? Round to 2 decimal places.", rows, sitesPer, used)
	return edakit.Numeric(q, ans, edakit.Sprintf("Utilization is used sites over available sites: $%d / %d$, which is %s%%.", used, rows*sitesPer, ans), map[string]string{"util_pct": ans})
}

var driveSizes = []int{1, 2, 4, 8, 16}

// driveStrength: delay of an Xk cell, the smallest size meeting a target, and the best size
// once the previous stage's load is counted.
func driveStrength(ctx generator.GeneratorContext) generator.Problem {
	r1 := float64(edakit.Int(4, 12))
	cp := float64(edakit.Int(1, 3))
	ci := float64(edakit.Int(1, 3))
	cl := float64(edakit.Scaled(ctx, 10, 200))
	delay := func(k int) float64 { return r1*cp + r1*cl/float64(k) }
	pre := edakit.Sprintf("In a cell family, the X1 size has output resistance %s kOhm, output parasitic capacitance %s fF and input capacitance %s fF. Size Xk has $1/k$ the resistance and $k$ times both capacitances. ", edakit.Num(r1, 0), edakit.Num(cp, 0), edakit.Num(ci, 0))
	switch kpOf(ctx) {
	case 1:
		k := edakit.Pick(driveSizes)
		ans := edakit.Num(delay(k), 2)
		q := pre + edakit.Sprintf("What is the RC delay of the X%d size driving %s fF, in ps? Round to 2 decimal places.", k, edakit.Num(cl, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$t = (R_1/k)(k C_p + C_L) = R_1 C_p + R_1 C_L / k = %s + %s = %s$ ps. The parasitic term does not shrink with size.", edakit.Num(r1*cp, 0), edakit.Num(r1*cl/float64(k), 3), ans),
			map[string]string{"delay_ps": ans})
	case 2:
		target := r1*cp + r1*cl/float64(edakit.Pick([]int{2, 4, 8, 16}))*float64(edakit.Int(100, 150))/100
		target = edakit.Round(target, 2)
		best := 0
		for _, k := range driveSizes {
			if delay(k) <= target+1e-9 {
				best = k
				break
			}
		}
		if best == 0 {
			best = 16
			target = delay(16)
		}
		q := pre + edakit.Sprintf("The cell drives %s fF and its RC delay must be at most %s ps. Sizes X1, X2, X4, X8 and X16 exist. What is the smallest k that meets the target?", edakit.Num(cl, 0), edakit.Num(target, 2))
		return edakit.Numeric(q, edakit.IntStr(best),
			edakit.Sprintf("Delay is $%s + %s/k$ ps. X%d gives %s ps, the first size at or under the target.", edakit.Num(r1*cp, 0), edakit.Num(r1*cl, 0), best, edakit.Num(delay(best), 2)),
			map[string]string{"k": edakit.IntStr(best)})
	}
	r0 := float64(edakit.Int(2, 12))
	cw := float64(edakit.Int(0, 10))
	for tries := 0; ; tries++ {
		total := func(k int) float64 { return r0*(cw+float64(k)*ci) + delay(k) }
		best, bestT, second := 0, 1e18, 1e18
		for _, k := range driveSizes {
			t := total(k)
			if t < bestT {
				second, best, bestT = bestT, k, t
			} else if t < second {
				second = t
			}
		}
		if second-bestT > 0.5 || tries > 20 {
			q := pre + edakit.Sprintf("A previous stage of resistance %s kOhm drives the cell's input through %s fF of wire, and the cell drives %s fF. Counting both stages' RC delays, which size among X1, X2, X4, X8 and X16 minimizes the total? Give k.", edakit.Num(r0, 0), edakit.Num(cw, 0), edakit.Num(cl, 0))
			return edakit.Numeric(q, edakit.IntStr(best),
				edakit.Sprintf("Total delay is $R_0(C_w + k C_{in}) + R_1 C_p + R_1 C_L/k$. Upsizing speeds the cell but loads its driver. Evaluating each size, X%d gives the minimum, %s ps.", best, edakit.Num(bestT, 2)),
				map[string]string{"k": edakit.IntStr(best), "total_ps": edakit.Num(bestT, 2)})
		}
		cl = float64(edakit.Scaled(ctx, 10, 200))
		r0 = float64(edakit.Int(2, 12))
	}
}
