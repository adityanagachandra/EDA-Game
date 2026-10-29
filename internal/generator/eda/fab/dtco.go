package fab

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var trackHeights = []float64{5, 5.5, 6, 6.5, 7, 7.5, 9, 10.5, 12}

// cellArchitecture: track height and cell height.
func cellArchitecture(ctx generator.GeneratorContext) generator.Problem {
	p := float64(2 * edakit.Int(10, 24))
	switch kpOf(ctx) {
	case 1:
		t := edakit.Pick(trackHeights)
		ans := edakit.Num(t*p, 1)
		q := edakit.Sprintf("A standard-cell library is a %s-track library on a %.0f nm horizontal metal pitch. What is the cell height, in nanometres? Round to 1 decimal place.", edakit.Num(t, 1), p)
		return edakit.Numeric(q, ans, edakit.Sprintf("Cell height is tracks times metal pitch: $%s \\times %.0f = %s$ nm.", edakit.Num(t, 1), p, ans), map[string]string{"h": ans})
	case 2:
		t := edakit.Pick(trackHeights)
		h := t * p
		ans := edakit.Num(t, 1)
		q := edakit.Sprintf("A cell row is %s nm tall and the horizontal routing pitch is %.0f nm. What is the library's track height? Round to 1 decimal place.", edakit.Num(h, 1), p)
		return edakit.Numeric(q, ans, edakit.Sprintf("Track height is cell height over pitch: %s / %.0f = %s tracks.", edakit.Num(h, 1), p, ans), map[string]string{"t": ans})
	default:
		i := edakit.Int(1, len(trackHeights)-1)
		t1 := trackHeights[i]
		t2 := trackHeights[edakit.Int(0, i-1)]
		red := (1 - t2/t1) * 100
		ans := edakit.Num(red, 1)
		q := edakit.Sprintf("A team moves a library from %s tracks to %s tracks at the same metal pitch and the same cell widths. By what percentage does cell area fall? Round to 1 decimal place.", edakit.Num(t1, 1), edakit.Num(t2, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Area scales with height at fixed width: $1 - %s/%s = %s$%%. Whether the block shrinks as much depends on routability with fewer tracks per cell.", edakit.Num(t2, 1), edakit.Num(t1, 1), ans),
			map[string]string{"reduction": ans})
	}
}

// cppWidth: cell width from contacted poly pitch, row width, and node area scaling.
func cppWidth(ctx generator.GeneratorContext) generator.Problem {
	cpp := float64(2 * edakit.Int(20, 32))
	switch kpOf(ctx) {
	case 1:
		g := edakit.Int(1, 8)
		ddb := edakit.Int(0, 1) == 1
		n := g + 1
		kind := "single diffusion break (one CPP of boundary per cell)"
		if ddb {
			n = g + 2
			kind = "double diffusion break (two CPP of boundary per cell)"
		}
		ans := edakit.Num(float64(n)*cpp, 0)
		q := edakit.Sprintf("A cell has %d gate (poly) tracks and uses a %s. The contacted poly pitch is %.0f nm. What is the cell width, in nanometres?", g, kind, cpp)
		return edakit.Numeric(q, ans, edakit.Sprintf("Width is $(%d + %d) \\times %.0f = %s$ nm.", g, n-g, cpp, ans), map[string]string{"cpp_count": edakit.IntStr(n)})
	case 2:
		k := edakit.Scaled(ctx, 3, 8)
		names := []string{"INV (1 gate)", "NAND2 (2 gates)", "NOR3 (3 gates)", "AOI22 (4 gates)", "DFF (12 gates)", "BUF (2 gates)", "XOR2 (6 gates)"}
		gates := []int{1, 2, 3, 4, 12, 2, 6}
		sum := 0
		var parts []string
		for i := 0; i < k; i++ {
			j := edakit.Int(0, len(names)-1)
			sum += gates[j] + 1
			parts = append(parts, names[j])
		}
		w := float64(sum) * cpp / 1000
		ans := edakit.Num(w, 3)
		q := edakit.Sprintf("A row holds, abutted, %s. Each cell is (gates + 1) CPP wide with single diffusion breaks, and CPP is %.0f nm. What is the total row length used, in micrometres? Round to 3 decimal places.", strings.Join(parts, ", "), cpp)
		return edakit.Numeric(q, ans, edakit.Sprintf("Total CPP count is %d, so length is $%d \\times %.0f$ nm $= %s$ um.", sum, sum, cpp, ans), map[string]string{"cpps": edakit.IntStr(sum)})
	default:
		mp := float64(2 * edakit.Int(12, 24))
		t := edakit.Pick(trackHeights[:6])
		cpp2 := cpp * step(0.75, 0.95, 0.01)
		mp2 := mp * step(0.7, 0.95, 0.01)
		t2 := t
		if edakit.Int(0, 1) == 1 && t > 5 {
			t2 = t - 0.5
		}
		r := (cpp2 * mp2 * t2) / (cpp * mp * t) * 100
		ans := edakit.Num(r, 1)
		q := edakit.Sprintf("Node A has CPP %.0f nm, metal pitch %.0f nm and a %s-track cell. Node B has CPP %s nm, metal pitch %s nm and a %s-track cell. Taking cell area as CPP times track height times metal pitch, what is a node B cell's area as a percentage of node A's? Round to 1 decimal place.", cpp, mp, edakit.Num(t, 1), edakit.Num(cpp2, 1), edakit.Num(mp2, 1), edakit.Num(t2, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Ratio $= (%s \\times %s \\times %s)/(%.0f \\times %s \\times %.0f) = %s$%%. Track reduction scales area as much as pitch does.", edakit.Num(cpp2, 1), edakit.Num(t2, 1), edakit.Num(mp2, 1), cpp, edakit.Num(t, 1), mp, ans),
			map[string]string{"ratio": ans})
	}
}

// finDepop: drive versus fin count, fin-pitch cell height, and drive density.
func finDepop(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n1 := edakit.Int(2, 4)
		n2 := edakit.Int(1, n1-1)
		ans := edakit.Num(float64(n2)/float64(n1)*100, 1)
		q := edakit.Sprintf("A library depopulates its transistors from %d fins to %d fins each. To first order, drive current is proportional to fin count. What percentage of the original drive current remains? Round to 1 decimal place.", n1, n2)
		return edakit.Numeric(q, ans, edakit.Sprintf("$%d/%d = %s$%%. Lower drive also means lower input capacitance, so delay into a light load falls less than this.", n2, n1, ans), map[string]string{"ratio": ans})
	case 2:
		fp := float64(2 * edakit.Int(12, 24))
		nf := edakit.Int(1, 4)
		g := edakit.Int(2, 6)
		h := float64(2*nf+g) * fp
		ans := edakit.Num(h, 0)
		q := edakit.Sprintf("In a FinFET cell, the height is set on the fin grid: %d active fins for NMOS, %d for PMOS, plus %d fin pitches for the n-to-p gap and the shared rails. The fin pitch is %.0f nm. What is the cell height, in nanometres?", nf, nf, g, fp)
		return edakit.Numeric(q, ans, edakit.Sprintf("Height is $(2 \\times %d + %d) \\times %.0f = %s$ nm. Each fin removed per device takes two fin pitches out of the cell.", nf, g, fp, ans), map[string]string{"pitches": edakit.IntStr(2*nf + g)})
	default:
		fp := float64(2 * edakit.Int(12, 24))
		n1 := edakit.Int(2, 4)
		n2 := n1 - 1
		g := edakit.Int(2, 6)
		h1 := float64(2*n1+g) * fp
		h2 := float64(2*n2+g) * fp
		r := (float64(n2) / h2) / (float64(n1) / h1)
		ans := edakit.Num(r, 3)
		q := edakit.Sprintf("A cell's height is $(2 n_{fin} + %d)$ fin pitches of %.0f nm. Depopulating from %d to %d fins per device, what is the ratio of drive current per unit cell area (new over old)? Assume equal widths and drive proportional to fins. Round to 3 decimal places.", g, fp, n1, n2)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Drive falls to $%d/%d$ and height to $%.0f/%.0f$ nm, so drive density changes by %s. Fixed overhead pitches keep area from shrinking as fast as drive.", n2, n1, h2, h1, ans),
			map[string]string{"h_old": edakit.Num(h1, 0), "h_new": edakit.Num(h2, 0)})
	}
}

// bpr: buried power rails free tracks and lower rail resistance.
func bpr(ctx generator.GeneratorContext) generator.Problem {
	p := float64(2 * edakit.Int(10, 20))
	sig := edakit.Pick([]float64{4, 4.5, 5, 5.5, 6})
	rail := edakit.Pick([]float64{1, 1.5, 2})
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num((sig)*p, 1)
		q := edakit.Sprintf("A cell is %s tracks tall on a %.0f nm pitch: %s signal tracks plus power rails that together occupy %s tracks. Burying the power rails below the devices removes the rail tracks without changing the signal tracks. What is the new cell height, in nanometres? Round to 1 decimal place.", edakit.Num(sig+rail, 1), p, edakit.Num(sig, 1), edakit.Num(rail, 1))
		return edakit.Numeric(q, ans, edakit.Sprintf("Only the %s signal tracks remain: $%s \\times %.0f = %s$ nm.", edakit.Num(sig, 1), edakit.Num(sig, 1), p, ans), map[string]string{"tracks": edakit.Num(sig, 1)})
	case 2:
		red := rail / (sig + rail) * 100
		ans := edakit.Num(red, 1)
		q := edakit.Sprintf("A %s-track cell spends %s tracks on its front-side power rails. Moving the rails into buried power rails frees those tracks. By what percentage can cell height shrink at the same signal-track count? Round to 1 decimal place.", edakit.Num(sig+rail, 1), edakit.Num(rail, 1))
		return edakit.Numeric(q, ans, edakit.Sprintf("$%s/%s = %s$%%.", edakit.Num(rail, 1), edakit.Num(sig+rail, 1), ans), map[string]string{"red": ans})
	default:
		w1 := float64(2 * edakit.Int(10, 20))
		h1 := float64(2 * edakit.Int(15, 30))
		w2 := float64(2 * edakit.Int(8, 15))
		h2 := float64(10 * edakit.Int(8, 20))
		r := (w2 * h2) / (w1 * h1)
		ans := edakit.Num(r, 2)
		q := edakit.Sprintf("A front-side M0 power rail is %.0f nm wide and %.0f nm thick. A buried power rail is %.0f nm wide but %.0f nm tall. With equal resistivity, how many times lower is the buried rail's resistance per unit length? Round to 2 decimal places.", w1, h1, w2, h2)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Resistance per length goes as one over cross-section: $(%.0f \\times %.0f)/(%.0f \\times %.0f) = %s$. A buried rail can be deep because it does not compete for routing tracks.", w2, h2, w1, h1, ans),
			map[string]string{"area_ratio": ans})
	}
}

// bspdnCompute: IR drop through a front-side via stack versus a backside connection.
func bspdnCompute(ctx generator.GeneratorContext) generator.Problem {
	I := step(0.5, 3.0, 0.1) // mA
	levels := edakit.Int(8, 14)
	rv := float64(edakit.Int(2, 8))
	k := edakit.Int(2, 6)
	rf := float64(levels) * rv / float64(k)
	rb := step(1.0, 4.0, 0.5)
	df := I * rf
	db := I * rb
	red := df - db
	if edakit.Int(0, 1) == 0 {
		q := edakit.Sprintf("A cell cluster draws %s mA. On the front side, current descends %d via levels, each %.0f ohm per via with %d vias in parallel per level. What is the IR drop across that via stack?", edakit.Num(I, 1), levels, rv, k)
		c, d := numOptions(df, 1, " mV", I*float64(levels)*rv, I*rv/float64(k), I*float64(levels)*rv*float64(k))
		return edakit.Choice(q, c,
			edakit.Sprintf("Stack resistance is $%d \\times %.0f / %d = %s$ ohm, and $V = I R$ gives %s. Every front-side power path pays this before reaching a cell.", levels, rv, k, edakit.Num(rf, 2), c), d...)
	}
	q := edakit.Sprintf("A cell cluster draws %s mA. The front-side power path is %d via levels of %.0f ohm per via with %d vias per level. A backside power network reaches the same cells through a %s ohm nano-through-silicon via path. How much lower is the IR drop with backside power?", edakit.Num(I, 1), levels, rv, k, edakit.Num(rb, 1))
	c, d := numOptions(red, 1, " mV", df, db, I*(float64(levels)*rv-rb))
	return edakit.Choice(q, c,
		edakit.Sprintf("Front side: $%s \\times %s = %s$ mV. Back side: $%s \\times %s = %s$ mV. The saving is %s.", edakit.Num(I, 1), edakit.Num(rf, 2), edakit.Num(df, 1), edakit.Num(I, 1), edakit.Num(rb, 1), edakit.Num(db, 1), c), d...)
}

// stco: monolithic die versus chiplets.
func stco(ctx generator.GeneratorContext) generator.Problem {
	cw := float64(1000 * edakit.Int(8, 20))
	A := float64(10 * edakit.Int(30, 80)) // mm^2
	D := step(0.08, 0.3, 0.01)
	k := edakit.Pick([]int{2, 4})
	oh := step(0.05, 0.15, 0.01)
	ya := step(0.95, 0.99, 0.01)
	nm := dpw(300, A)
	ym := math.Exp(-A / 100 * D)
	cm := cw / (float64(nm) * ym)
	ac := A / float64(k) * (1 + oh)
	nc := dpw(300, ac)
	yc := math.Exp(-ac / 100 * D)
	cc := float64(k) * cw / (float64(nc) * yc) / ya
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(cm, 2)
		q := edakit.Sprintf("A monolithic %.0f mm^2 die is made on 300 mm wafers costing %.0f US dollars, with %d dies per wafer and Poisson yield at %s defects per cm^2. What is the cost per good die, in dollars? Round to 2 decimal places.", A, cw, nm, edakit.Num(D, 2))
		return edakit.Numeric(q, ans, edakit.Sprintf("Yield $e^{-%s \\times %s} = %s$; cost $= %.0f/(%d \\times %s) = %s$ dollars.", edakit.Num(A/100, 2), edakit.Num(D, 2), edakit.Num(ym, 4), cw, nm, edakit.Num(ym, 4), ans), map[string]string{"yield": edakit.Num(ym, 4)})
	case 2:
		ans := edakit.Num(cc, 2)
		q := edakit.Sprintf("The same %.0f mm^2 design is split into %d equal chiplets, each %s%% larger than its share for die-to-die interfaces, so each chiplet is %s mm^2 with %d per wafer. Wafers cost %.0f US dollars and defect density is %s per cm^2 (Poisson). Chiplets are tested before assembly, and assembly yield is %s. What is the silicon cost per good assembled system, in dollars? Round to 2 decimal places.", A, k, edakit.Num(oh*100, 0), edakit.Num(ac, 2), nc, cw, edakit.Num(D, 2), edakit.Num(ya, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Chiplet yield $e^{-%s \\times %s} = %s$; cost per good chiplet %s dollars; times %d chiplets, divided by assembly yield %s, gives %s dollars.", edakit.Num(ac/100, 4), edakit.Num(D, 2), edakit.Num(yc, 4), edakit.Num(cw/(float64(nc)*yc), 2), k, edakit.Num(ya, 2), ans),
			map[string]string{"chiplet_yield": edakit.Num(yc, 4)})
	default:
		ans := edakit.Num(cm/cc, 3)
		q := edakit.Sprintf("Monolithic: a %.0f mm^2 die costs %s dollars per good die. Chiplets: %d tested chiplets of %s mm^2 each, %d per wafer, with Poisson yield at %s defects per cm^2, %.0f-dollar wafers and assembly yield %s. What is the monolithic cost divided by the chiplet system cost? Round to 3 decimal places.", A, edakit.Num(cm, 2), k, edakit.Num(ac, 2), nc, edakit.Num(D, 2), cw, edakit.Num(ya, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The chiplet system costs %s dollars, so the ratio is %s. Above 1, splitting wins on silicon cost; packaging cost and interface power must then be weighed.", edakit.Num(cc, 2), ans),
			map[string]string{"chiplet_cost": edakit.Num(cc, 2)})
	}
}
