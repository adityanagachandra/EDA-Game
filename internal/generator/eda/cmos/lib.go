package cmos

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// table is a two-dimensional Liberty lookup table with integer entries.
type table struct {
	name   string
	idx1   []int // input transition (ps) or the constrained pin's transition
	idx2   []float64
	vals   [][]int
	i1name string
	i2name string
	i2unit string
}

func (t table) text() string {
	var b strings.Builder
	var i1, i2 []string
	for _, x := range t.idx1 {
		i1 = append(i1, edakit.IntStr(x))
	}
	for _, x := range t.idx2 {
		i2 = append(i2, edakit.Num(x, 2))
	}
	b.WriteString(edakit.Sprintf("Table %s: index_1 (%s, ps) = %s; index_2 (%s, %s) = %s.", t.name, t.i1name, strings.Join(i1, ", "), t.i2name, t.i2unit, strings.Join(i2, ", ")))
	for r, row := range t.vals {
		var cells []string
		for _, v := range row {
			cells = append(cells, edakit.IntStr(v))
		}
		b.WriteString(edakit.Sprintf(" Row %s = %d: %s.", t.i1name, t.idx1[r], strings.Join(cells, ", ")))
	}
	return b.String()
}

var slewGrids = [][]int{{10, 30, 90}, {5, 20, 80}, {8, 25, 75}, {15, 45, 135}, {10, 40, 160}, {5, 15, 45, 135}, {10, 20, 40, 80}}
var loadGrids = [][]float64{{1, 4, 16}, {0.5, 2, 8}, {2, 6, 18}, {1, 3, 9, 27}, {0.5, 1, 2, 4}, {1, 5, 25}}

// drawTable fills a delay-like table: v = a + b*s + c*l + e*s*l.
func drawTable(name string, slews []int, loads []float64, a, b, c, e float64) table {
	t := table{name: name, idx1: slews, idx2: loads, i1name: "input transition", i2name: "output load", i2unit: "fF"}
	for _, s := range slews {
		var row []int
		for _, l := range loads {
			row = append(row, int(a+b*float64(s)+c*l+e*float64(s)*l+0.5))
		}
		t.vals = append(t.vals, row)
	}
	return t
}

func randDelayTable(name string, slews []int, loads []float64) table {
	return drawTable(name, slews, loads, float64(edakit.Int(6, 30)), float64(edakit.Int(10, 40))/100, float64(edakit.Int(20, 90))/10, float64(edakit.Int(0, 3))/100)
}

func randTransTable(name string, slews []int, loads []float64) table {
	return drawTable(name, slews, loads, float64(edakit.Int(4, 15)), float64(edakit.Int(5, 25))/100, float64(edakit.Int(30, 120))/10, 0)
}

func gridPoint(t table) (int, int) { return edakit.Int(0, len(t.idx1)-1), edakit.Int(0, len(t.idx2)-1) }

// nldmLookup: one lookup, choosing rise or fall from unateness, and a two-stage chain.
func nldmLookup(ctx generator.GeneratorContext) generator.Problem {
	slews, loads := edakit.Pick(slewGrids), edakit.Pick(loadGrids)
	switch kpOf(ctx) {
	case 1:
		t := randDelayTable(edakit.Pick([]string{"cell_rise", "cell_fall", "rise_transition", "fall_transition"}), slews, loads)
		i, j := gridPoint(t)
		q := t.text() + edakit.Sprintf(" Values are in ps. Which value does the table give at an input transition of %d ps and an output load of %s fF?", t.idx1[i], edakit.Num(t.idx2[j], 2))
		return edakit.Numeric(q, edakit.IntStr(t.vals[i][j]),
			edakit.Sprintf("index_1 selects the row (input transition %d ps) and index_2 the column (load %s fF); the entry there is %d ps.", t.idx1[i], edakit.Num(t.idx2[j], 2), t.vals[i][j]),
			map[string]string{"value_ps": edakit.IntStr(t.vals[i][j])})
	case 2:
		rise := randDelayTable("cell_rise", slews, loads)
		fall := randDelayTable("cell_fall", slews, loads)
		i, j := gridPoint(rise)
		neg := edakit.Int(0, 1) == 0
		inRise := edakit.Int(0, 1) == 0
		cell, sense := "BUF", "positive_unate"
		if neg {
			cell, sense = "INV", "negative_unate"
		}
		outRise := inRise != neg
		tab, edge := fall, "falling"
		if outRise {
			tab, edge = rise, "rising"
		}
		in := "falls"
		if inRise {
			in = "rises"
		}
		q := edakit.Sprintf("%s cell's arc from A to Y is %s. ", capArticle(cell), sense) + rise.text() + " " + fall.text() +
			edakit.Sprintf(" Values are in ps. Input A %s with a %d ps transition while Y drives %s fF. What is the arc delay, in ps?", in, rise.idx1[i], edakit.Num(rise.idx2[j], 2))
		return edakit.Numeric(q, edakit.IntStr(tab.vals[i][j]),
			edakit.Sprintf("The tables are named for the output edge. A %s arc with input that %s gives a %s output, so read %s at that grid point: %d ps.", sense, in, edge, tab.name, tab.vals[i][j]),
			map[string]string{"delay_ps": edakit.IntStr(tab.vals[i][j])})
	}
	a := randDelayTable("U1 cell_fall", slews, loads)
	at := randTransTable("U1 fall_transition", slews, loads)
	i, j := gridPoint(a)
	tA := at.vals[i][j]
	set := map[int]bool{tA: true}
	for _, x := range []int{tA / 3, tA * 3, tA / 2, tA * 2} {
		if x > 0 && len(set) < 3 {
			set[x] = true
		}
	}
	var bSlews []int
	for x := range set {
		bSlews = append(bSlews, x)
	}
	sort.Ints(bSlews)
	bLoads := edakit.Pick(loadGrids)
	b := randDelayTable("U2 cell_rise", bSlews, bLoads)
	k := edakit.Int(0, len(bLoads)-1)
	bi := sort.SearchInts(bSlews, tA)
	tot := a.vals[i][j] + b.vals[bi][k]
	q := edakit.Sprintf("Inverter U1 receives a rising input with a %d ps transition and drives %s fF, which includes the input of inverter U2; U2 drives %s fF. ", a.idx1[i], edakit.Num(a.idx2[j], 2), edakit.Num(bLoads[k], 2)) +
		a.text() + " " + at.text() + " " + b.text() + " All values are in ps. What is the delay from U1's input to U2's output, in ps?"
	return edakit.Numeric(q, edakit.IntStr(tot),
		edakit.Sprintf("U1's output falls: cell_fall gives %d ps and fall_transition gives %d ps. That transition is U2's input slew, and U2's output rises: cell_rise at (%d ps, %s fF) gives %d ps. Total %d ps.", a.vals[i][j], tA, tA, edakit.Num(bLoads[k], 2), b.vals[bi][k], tot),
		map[string]string{"u1_ps": edakit.IntStr(a.vals[i][j]), "u1_slew_ps": edakit.IntStr(tA), "u2_ps": edakit.IntStr(b.vals[bi][k])})
}

// ccsLvf: CCS charge integration, LVF late delay, and the early-late spread.
func ccsLvf(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 3, 6)
		dt := float64(edakit.Pick([]int{2, 4, 5, 10}))
		c := float64(edakit.Int(2, 20))
		var parts []string
		q0 := 0.0
		for i := 0; i < n; i++ {
			cur := float64(edakit.Int(5, 80))
			q0 += cur * dt
			parts = append(parts, edakit.Num(cur, 0))
		}
		v := q0 / c
		ans := edakit.Num(v, 1)
		q := edakit.Sprintf("A CCS driver model gives the output current into a %s fF load as samples held constant over successive %s ps intervals: %s uA. Starting from 0 V, what is the output voltage after these intervals, in millivolts? Round to 1 decimal place.", edakit.Num(c, 0), edakit.Num(dt, 0), strings.Join(parts, ", "))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The delivered charge is $\\sum i \\Delta t = %s$ uA ps $= %s$ aC, and $V = Q/C$; aC per fF is mV, so %s mV.", edakit.Num(q0, 0), edakit.Num(q0, 0), ans),
			map[string]string{"charge_ac": edakit.Num(q0, 0), "v_mv": ans})
	case 2:
		slews, loads := edakit.Pick(slewGrids), edakit.Pick(loadGrids)
		mu := randDelayTable("cell_rise", slews, loads)
		sg := drawTable("ocv_sigma_cell_rise (late)", slews, loads, float64(edakit.Int(1, 3)), float64(edakit.Int(1, 5))/100, float64(edakit.Int(1, 8))/10, 0)
		i, j := gridPoint(mu)
		ns := edakit.Pick([]int{2, 3})
		late := mu.vals[i][j] + ns*sg.vals[i][j]
		q := mu.text() + " " + sg.text() + edakit.Sprintf(" Values are in ps. A POCV analysis uses %d sigma. What late rise delay does it use at an input transition of %d ps and a load of %s fF, in ps?", ns, mu.idx1[i], edakit.Num(mu.idx2[j], 2))
		return edakit.Numeric(q, edakit.IntStr(late),
			edakit.Sprintf("LVF stores the sigma on the same grid as the mean. Late delay is $\\mu + n\\sigma = %d + %d \\times %d = %d$ ps.", mu.vals[i][j], ns, sg.vals[i][j], late),
			map[string]string{"mean_ps": edakit.IntStr(mu.vals[i][j]), "sigma_ps": edakit.IntStr(sg.vals[i][j])})
	}
	mu := float64(edakit.Int(20, 120))
	se := float64(edakit.Int(10, 60)) / 10
	sl := se + float64(edakit.Int(0, 30))/10
	ns := float64(edakit.Pick([]int{2, 3}))
	if edakit.Int(0, 1) == 0 {
		ans := edakit.Num(mu-ns*se, 1)
		q := edakit.Sprintf("At one grid point an LVF arc has mean delay %s ps, early sigma %s ps and late sigma %s ps. Early and late sigmas differ because the delay distribution is skewed. What early delay does a %s-sigma analysis use, in ps? Round to 1 decimal place.", edakit.Num(mu, 0), edakit.Num(se, 1), edakit.Num(sl, 1), edakit.Num(ns, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("Early delay subtracts the early sigma: $%s - %s \\times %s = %s$ ps.", edakit.Num(mu, 0), edakit.Num(ns, 0), edakit.Num(se, 1), ans), nil)
	}
	ans := edakit.Num(ns*(se+sl), 1)
	q := edakit.Sprintf("At one grid point an LVF arc has mean delay %s ps, early sigma %s ps and late sigma %s ps. In a %s-sigma analysis, how far apart are the late and early delays, in ps? Round to 1 decimal place.", edakit.Num(mu, 0), edakit.Num(se, 1), edakit.Num(sl, 1), edakit.Num(ns, 0))
	return edakit.Numeric(q, ans, edakit.Sprintf("Late is $\\mu + n\\sigma_{late}$ and early is $\\mu - n\\sigma_{early}$, so the gap is $n(\\sigma_{early} + \\sigma_{late}) = %s$ ps; the mean cancels.", ans), nil)
}

// powerViews: state-dependent leakage, internal power from a table, and total pin power.
func powerViews(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		pa := float64(edakit.Int(1, 9)) / 10
		pb := float64(edakit.Int(1, 9)) / 10
		l := []float64{float64(edakit.Int(1, 8)), float64(edakit.Int(5, 20)), float64(edakit.Int(5, 20)), float64(edakit.Int(15, 60))}
		pr := []float64{(1 - pa) * (1 - pb), (1 - pa) * pb, pa * (1 - pb), pa * pb}
		avg := 0.0
		for i := range l {
			avg += l[i] * pr[i]
		}
		ans := edakit.Num(avg, 3)
		q := edakit.Sprintf("A NAND2's Liberty leakage_power groups give, by input state: when \"!A & !B\" %s nW, when \"!A & B\" %s nW, when \"A & !B\" %s nW, when \"A & B\" %s nW. Inputs are independent with $P(A = 1) = %s$ and $P(B = 1) = %s$. What is the expected leakage in nW? Round to 3 decimal places.",
			edakit.Num(l[0], 0), edakit.Num(l[1], 0), edakit.Num(l[2], 0), edakit.Num(l[3], 0), edakit.Num(pa, 1), edakit.Num(pb, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Weight each state's leakage by its probability: the state probabilities are %s, %s, %s and %s, giving %s nW.", edakit.Num(pr[0], 2), edakit.Num(pr[1], 2), edakit.Num(pr[2], 2), edakit.Num(pr[3], 2), ans),
			map[string]string{"leak_nw": ans})
	case 2:
		slews, loads := edakit.Pick(slewGrids), edakit.Pick(loadGrids)
		t := drawTable("rise_power", slews, loads, float64(edakit.Int(1, 4)), float64(edakit.Int(1, 4))/100, float64(edakit.Int(1, 6))/10, 0)
		t.name = "internal_power rise_power (fJ per transition)"
		i, j := gridPoint(t)
		rate := float64(edakit.Pick([]int{50, 100, 200, 250, 400, 500}))
		p := float64(t.vals[i][j]) * rate / 1000
		ans := edakit.Num(p, 3)
		q := t.text() + edakit.Sprintf(" The output rises %s million times per second, each time with a %d ps input transition and a %s fF load. What internal power do the rising transitions contribute, in microwatts? Round to 3 decimal places.", edakit.Num(rate, 0), t.idx1[i], edakit.Num(t.idx2[j], 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The table gives %d fJ per rising transition at that point. Times %s million per second is %s uW.", t.vals[i][j], edakit.Num(rate, 0), ans),
			map[string]string{"energy_fj": edakit.IntStr(t.vals[i][j])})
	}
	eint := float64(edakit.Int(1, 12))
	c := float64(edakit.Int(2, 30))
	v := edakit.Pick([]float64{0.7, 0.8, 0.9, 1.0})
	rate := float64(edakit.Pick([]int{50, 100, 200, 400}))
	sw := 0.5 * c * v * v
	p := (eint + sw) * rate / 1000
	ans := edakit.Num(p, 3)
	q := edakit.Sprintf("An output pin toggles %s million times per second (counting rising and falling edges). Its internal_power table gives %s fJ per toggle at the operating point; that energy excludes charging the external load. The net load is %s fF at $V_{DD} = %s$ V, and each toggle of the net costs $\\frac{1}{2} C V^2$ of switching energy. What is the total dynamic power of the pin, in microwatts? Round to 3 decimal places.", edakit.Num(rate, 0), edakit.Num(eint, 0), edakit.Num(c, 0), edakit.Num(v, 1))
	return edakit.Numeric(q, ans,
		edakit.Sprintf("Switching energy per toggle is $\\frac{1}{2} \\times %s \\times %s^2 = %s$ fJ. Adding the internal %s fJ and multiplying by %s million per second gives %s uW.", edakit.Num(c, 0), edakit.Num(v, 1), edakit.Num(sw, 3), edakit.Num(eint, 0), edakit.Num(rate, 0), ans),
		map[string]string{"switching_fj": edakit.Num(sw, 3), "p_uw": ans})
}

// characterization: simulation counts, threshold measurements and full-swing extrapolation.
func characterization(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		cells := 10 * edakit.Scaled(ctx, 5, 100)
		arcs := edakit.Int(2, 8)
		n1, n2 := edakit.Pick([]int{5, 7, 8}), edakit.Pick([]int{5, 7, 8})
		corners := edakit.Int(1, 12)
		tot := cells * arcs * 2 * n1 * n2 * corners
		q := edakit.Sprintf("A library of %d cells averages %d delay arcs per cell. Each arc is characterized for a rising and a falling output on a %d by %d grid of input transition and load, at %d PVT corners, one SPICE transient per point. How many transient simulations is that?", cells, arcs, n1, n2, corners)
		return edakit.Numeric(q, edakit.IntStr(tot),
			edakit.Sprintf("Multiply: %d cells x %d arcs x 2 edges x %d grid points x %d corners = %d.", cells, arcs, n1*n2, corners, tot),
			map[string]string{"sims": edakit.IntStr(tot)})
	case 2:
		ts := float64(edakit.Int(0, 50))
		tin := float64(edakit.Int(10, 120))
		to := ts + float64(edakit.Int(5, 60))
		tout := float64(edakit.Int(10, 150))
		if edakit.Int(0, 1) == 0 {
			d := (to + tout/2) - (ts + tin/2)
			ans := edakit.Num(d, 1)
			q := edakit.Sprintf("In a characterization transient the input ramps linearly from 0 to $V_{DD}$ between %s ps and %s ps, and the output falls linearly from $V_{DD}$ to 0 between %s ps and %s ps. Delay is measured between the 50%% crossings. What is the delay, in ps? Round to 1 decimal place.", edakit.Num(ts, 0), edakit.Num(ts+tin, 0), edakit.Num(to, 0), edakit.Num(to+tout, 0))
			return edakit.Numeric(q, ans,
				edakit.Sprintf("The input crosses 50%% at %s ps and the output at %s ps, so the delay is %s ps.", edakit.Num(ts+tin/2, 1), edakit.Num(to+tout/2, 1), ans),
				map[string]string{"delay_ps": ans})
		}
		lo := edakit.Pick([]int{10, 20, 30})
		hi := 100 - lo
		tr := tout * float64(hi-lo) / 100
		ans := edakit.Num(tr, 1)
		q := edakit.Sprintf("In a characterization transient the output falls linearly from $V_{DD}$ to 0 between %s ps and %s ps. The library measures transition between the %d%% and %d%% points. What transition is measured, in ps? Round to 1 decimal place.", edakit.Num(to, 0), edakit.Num(to+tout, 0), hi, lo)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("A linear ramp spends %d%% of its %s ps swing between those thresholds: %s ps.", hi-lo, edakit.Num(tout, 0), ans),
			map[string]string{"transition_ps": ans})
	}
	lo := edakit.Pick([]int{10, 20, 30})
	hi := 100 - lo
	m := float64(edakit.Int(8, 120))
	full := m * 100 / float64(hi-lo)
	ans := edakit.Num(full, 2)
	q := edakit.Sprintf("A library measures output transitions between the %d%% and %d%% points but stores them as the equivalent full-swing (0 to 100%%) linear ramp time. A simulation measures %s ps between those thresholds. What value is stored, in ps? Round to 2 decimal places.", lo, hi, edakit.Num(m, 0))
	return edakit.Numeric(q, ans,
		edakit.Sprintf("The thresholds span %d%% of the swing, so the full ramp is $%s / %s = %s$ ps. A tool reading the library must know which convention was used, or every slew is off by this factor.", hi-lo, edakit.Num(m, 0), edakit.Num(float64(hi-lo)/100, 2), ans),
		map[string]string{"full_ps": ans})
}

// sequentialArcs: setup constraint lookup, setup slack and hold slack.
func sequentialArcs(ctx generator.GeneratorContext) generator.Problem {
	dsl := edakit.Pick(slewGrids)
	csl := []float64{10, 40, 160}
	if edakit.Int(0, 1) == 0 {
		csl = []float64{20, 60, 180}
	}
	kp := kpOf(ctx)
	name := "rise_constraint (timing_type setup_rising)"
	if kp == 3 {
		name = "rise_constraint (timing_type hold_rising)"
	}
	t := table{name: name, idx1: dsl, idx2: csl, i1name: "data transition", i2name: "clock transition", i2unit: "ps"}
	a := float64(edakit.Int(10, 40))
	bd := float64(edakit.Int(10, 30)) / 100
	bc := -float64(edakit.Int(5, 20)) / 100
	if kp == 3 {
		a = float64(edakit.Int(-10, 20))
		bd = -float64(edakit.Int(5, 20)) / 100
		bc = float64(edakit.Int(10, 30)) / 100
	}
	for _, s := range dsl {
		var row []int
		for _, c := range csl {
			v := a + bd*float64(s) + bc*c
			if v < 0 {
				row = append(row, int(v-0.5))
			} else {
				row = append(row, int(v+0.5))
			}
		}
		t.vals = append(t.vals, row)
	}
	i, j := edakit.Int(0, len(dsl)-1), edakit.Int(0, len(csl)-1)
	cv := t.vals[i][j]
	pre := "A flop's Liberty constraint arc from CK (related_pin) to D. " + t.text() + " Values are in ps. "
	switch kp {
	case 1:
		q := pre + edakit.Sprintf("What setup time applies when D rises with a %d ps transition and CK a %s ps transition?", dsl[i], edakit.Num(csl[j], 0))
		return edakit.Numeric(q, edakit.IntStr(cv), edakit.Sprintf("index_1 is the data transition and index_2 the clock transition here, so the entry at (%d, %s) is %d ps.", dsl[i], edakit.Num(csl[j], 0), cv), nil)
	case 2:
		period := 50 * edakit.Int(6, 30)
		clk := edakit.Int(20, 120)
		arr := period + clk - cv - edakit.Int(-60, 60)
		slack := period + clk - cv - arr
		q := pre + edakit.Sprintf("The clock period is %d ps and the capture clock reaches CK %d ps after the ideal edge, with a %s ps transition. Rising data arrives at D at %d ps (measured from the launching edge) with a %d ps transition. What is the setup slack, in ps?", period, clk, edakit.Num(csl[j], 0), arr, dsl[i])
		return edakit.Numeric(q, edakit.IntStr(slack),
			edakit.Sprintf("The setup time from the table is %d ps. Required time is $%d + %d - %d = %d$ ps, so slack is $%d - %d = %d$ ps.", cv, period, clk, cv, period+clk-cv, period+clk-cv, arr, slack),
			map[string]string{"setup_ps": edakit.IntStr(cv), "slack_ps": edakit.IntStr(slack)})
	}
	clk := edakit.Int(10, 80)
	arr := clk + cv + edakit.Int(-30, 40)
	slack := arr - clk - cv
	q := pre + edakit.Sprintf("The clock reaches CK %d ps after the edge with a %s ps transition. The earliest rising data change at D arrives %d ps after the same edge, with a %d ps transition. What is the hold slack, in ps?", clk, edakit.Num(csl[j], 0), arr, dsl[i])
	return edakit.Numeric(q, edakit.IntStr(slack),
		edakit.Sprintf("The hold time from the table is %d ps. Data must not change before $%d %s = %d$ ps, so hold slack is $%d - %d = %d$ ps.", cv, clk, signed(cv), clk+cv, arr, clk+cv, slack),
		map[string]string{"hold_ps": edakit.IntStr(cv), "slack_ps": edakit.IntStr(slack)})
}

// unatenessPart is the computed KP 1 of liberty_structure.
func unatenessPart(ctx generator.GeneratorContext) generator.Problem {
	type arc struct{ desc, sense, why string }
	arcs := []arc{
		{"an AND2 from input A to Y", "positive_unate", "raising A can only raise or hold Y"},
		{"an OR3 from input B to Y", "positive_unate", "raising B can only raise or hold Y"},
		{"a NAND2 from input B to Y", "negative_unate", "raising B can only lower or hold Y"},
		{"a NOR2 from input A to Y", "negative_unate", "raising A can only lower or hold Y"},
		{"an inverter from A to Y", "negative_unate", "the output always moves opposite to the input"},
		{"a buffer from A to Y", "positive_unate", "the output always follows the input"},
		{"an XOR2 from input A to Y", "non_unate", "a rising A can raise or lower Y depending on B"},
		{"an XNOR2 from input B to Y", "non_unate", "the direction of Y's change depends on the other input"},
		{"a 2:1 mux from the select pin S to Y", "non_unate", "a rising S can raise or lower Y depending on which data input is 1"},
		{"a 2:1 mux from data input I0 to Y", "positive_unate", "with S fixed, Y can only follow I0"},
		{"an AOI21 from input C to Y", "negative_unate", "every input of an AND-OR-invert gate can only lower Y when it rises"},
		{"an OAI22 from input A to Y", "negative_unate", "every input of an OR-AND-invert gate can only lower Y when it rises"},
		{"a full adder from input A to the carry-out CO", "positive_unate", "the majority function is monotone increasing in each input"},
		{"a full adder from input A to the sum S", "non_unate", "the sum is an XOR of all three inputs"},
		{"an AO22 from input D to Y", "positive_unate", "AND-OR without inversion is monotone increasing"},
	}
	a := edakit.Pick(arcs)
	q := "What timing_sense should the Liberty timing group carry for the arc of " + a.desc + "?"
	return edakit.Choice(q, a.sense, "It is "+a.sense+" because "+a.why+".",
		[]string{"positive_unate", "negative_unate", "non_unate", "combinational"}[0:4]...)
}
