package flow

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var implChain = []string{"RTL design", "Logic synthesis", "Scan insertion", "Floorplanning", "Placement",
	"Clock tree synthesis", "Routing", "Parasitic extraction", "Signoff", "Tape-out", "Fabrication",
	"Packaging", "Production test"}

// firstKnown pairs a quantity with the first stage of implChain at which it exists.
var firstKnown = []struct {
	what  string
	stage int
}{
	{"the number of library gates in the design", 1},
	{"scan chains that connect the flip-flops into shift registers", 2},
	{"the die outline and the positions of the memory macros", 3},
	{"a legal row-and-site location for every standard cell", 4},
	{"the real clock skew between two flip-flops, set by buffers and wires", 5},
	{"the exact metal geometry of every signal wire", 6},
	{"the coupling capacitance between two neighbouring routed wires", 7},
	{"a verdict on timing and design rules accepted by the foundry flow", 8},
	{"the final layout database released for mask making", 9},
	{"physical silicon dies on a wafer", 10},
	{"the solder balls or bond wires that connect the die to a board", 11},
}

// rtlToGDS: stage successors (KP 1), orderings (KP 2) and the stage at which information first
// exists (KP 3), all drawn from the canonical chain.
func rtlToGDS(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx, 3) {
	case 1:
		i := edakit.Int(0, len(implChain)-2)
		ans := implChain[i+1]
		var near []string
		if i > 0 {
			near = append(near, implChain[i-1])
		}
		if i+2 < len(implChain) {
			near = append(near, implChain[i+2])
		}
		ds := append(near, pickOthers(implChain, 4-len(near), append([]string{ans, implChain[i]}, near...)...)...)
		q := edakit.Sprintf("In a standard digital implementation flow, which stage comes directly after %s?", lowerFirst(implChain[i]))
		if i == 0 {
			q = "In a standard digital implementation flow, which stage comes directly after RTL design is complete and verified?"
		}
		return edakit.Choice(q, ans, edakit.Sprintf("The chain runs %s, so %s follows %s.", strings.Join(implChain, ", "), ans, implChain[i]), ds...)
	case 2:
		n := 4
		idx := edakit.Shuffle([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})[:n]
		sort.Ints(idx)
		names := make([]string, n)
		for k, v := range idx {
			names[k] = implChain[v]
		}
		correct := strings.Join(names, " → ")
		perm := func(a, b int) string {
			x := append([]string(nil), names...)
			x[a], x[b] = x[b], x[a]
			return strings.Join(x, " → ")
		}
		q := "Which ordering of these four stages matches the standard digital implementation flow?"
		return edakit.Choice(q, correct,
			"In the flow each stage needs the previous one's output: "+correct+".",
			perm(0, 1), perm(1, 2), perm(2, 3), perm(0, 3))
	default:
		it := firstKnown[edakit.Int(0, len(firstKnown)-1)]
		ans := implChain[it.stage]
		ds := []string{implChain[it.stage-1], implChain[it.stage+1]}
		ds = append(ds, pickOthers(implChain, 1, ans, ds[0], ds[1])...)
		q := edakit.Sprintf("At which stage of the flow does %s first exist?", it.what)
		return edakit.Choice(q, ans,
			edakit.Sprintf("Before %s this quantity can only be estimated; %s is the stage that creates it.", strings.ToLower(ans), ans), ds...)
	}
}

// verification: (state, input) space size, cross-coverage percentage, and late-bug cost.
func verification(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx, 3) {
	case 1:
		n, r := edakit.Int(2, 8), edakit.Scaled(ctx, 2, 12)
		ans := 1 << (n + r)
		q := edakit.Sprintf("A small controller has %d input bits and %d bits of state. An exhaustive test must apply every input combination in every state. How many (state, input) pairs is that?", n, r)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("There are $2^{%d}$ states and $2^{%d}$ input combinations, so $2^{%d} = %d$ pairs, and reaching each state may itself take many cycles.", r, n, n+r, ans),
			map[string]string{"pairs": edakit.IntStr(ans)})
	case 2:
		a, b := edakit.Int(3, 10), edakit.Int(2, 6)
		total := a * b
		dims := edakit.Sprintf("a coverpoint with %d bins crossed with a coverpoint with %d bins", a, b)
		if ctx.Difficulty > 0.6 {
			c := edakit.Int(2, 3)
			total *= c
			dims = edakit.Sprintf("coverpoints with %d, %d and %d bins crossed together", a, b, c)
		}
		hit := edakit.Int(total/3, total-1)
		ans := edakit.Num(100*float64(hit)/float64(total), 2)
		q := edakit.Sprintf("A functional coverage model has %s. Regression tests have hit %d of the cross bins. What is the cross coverage, in percent? Round to 2 decimal places.", dims, hit)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The cross has %d bins; %d/%d is %s%%.", total, hit, total, ans),
			map[string]string{"bins": edakit.IntStr(total)})
	default:
		c := edakit.Int(1, 5)
		m := edakit.Pick([]int{5, 8, 10})
		k := edakit.Int(1, 3)
		stages := []string{"the gate-level netlist", "layout", "silicon"}
		ans := c
		for i := 0; i < k; i++ {
			ans *= m
		}
		q := edakit.Sprintf("A bug costs %d engineer-days to fix when found in RTL. Assume the fix cost multiplies by %d for each later stage it escapes to, in the order gate-level netlist, layout, silicon. How many engineer-days does it cost if it is first found in %s?", c, m, stages[k-1])
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("It escapes %d stage(s), so the cost is $%d \\cdot %d^{%d} = %d$ engineer-days.", k, c, m, k, ans),
			map[string]string{"stages_escaped": edakit.IntStr(k)})
	}
}

// scanTest: chain length, test cycles and test time.
func scanTest(ctx generator.GeneratorContext) generator.Problem {
	f := 100 * edakit.Scaled(ctx, 10, 500)
	c := edakit.Pick([]int{4, 8, 16, 24, 32, 64, 100})
	l := (f + c - 1) / c
	switch kp(ctx, 3) {
	case 1:
		q := edakit.Sprintf("A design has %s scan flip-flops, divided as evenly as possible into %d scan chains. How many flip-flops are in the longest chain?", comma(f), c)
		return edakit.Numeric(q, edakit.IntStr(l),
			edakit.Sprintf("$\\lceil %d / %d \\rceil = %d$.", f, c, l), map[string]string{"chain_length": edakit.IntStr(l)})
	case 2:
		n := 100 * edakit.Int(5, 50)
		cyc := (n+1)*l + n
		q := edakit.Sprintf("The longest scan chain has %d flip-flops and the test applies %d patterns. Shifting a pattern in takes one cycle per flip-flop, each pattern then needs one capture cycle, and the shift-out of each response overlaps the shift-in of the next pattern, with one final shift-out. How many test clock cycles does the whole test take?", l, n)
		return edakit.Numeric(q, edakit.IntStr(cyc),
			edakit.Sprintf("$(N + 1)L + N = %d \\cdot %d + %d = %d$ cycles.", n+1, l, n, cyc), map[string]string{"cycles": edakit.IntStr(cyc)})
	default:
		n := 100 * edakit.Int(5, 50)
		mhz := edakit.Pick([]int{10, 20, 25, 50, 100})
		cyc := (n+1)*l + n
		ms := float64(cyc) / (float64(mhz) * 1000)
		ans := edakit.Num(ms, 2)
		q := edakit.Sprintf("A design has %s scan flip-flops balanced into %d chains, and the test applies %d patterns with one capture cycle each and overlapped shift-in and shift-out. The scan shift clock runs at %d MHz. How long does the test take, in milliseconds? Round to 2 decimal places.", comma(f), c, n, mhz)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Chain length $L = \\lceil %d/%d \\rceil = %d$. Cycles $= (N+1)L + N = %d$. At %d MHz that is %s ms.", f, c, l, cyc, mhz, ans),
			map[string]string{"chain_length": edakit.IntStr(l), "cycles": edakit.IntStr(cyc)})
	}
}

// floorplan: core area from utilization, core width from aspect ratio, die area with margins.
func floorplan(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx, 3) {
	case 1:
		a := 1000 * edakit.Scaled(ctx, 200, 5000)
		u := float64(edakit.Int(10, 16)) * 0.05
		core := float64(a) / u
		ans := edakit.Num(core, 0)
		q := edakit.Sprintf("Synthesis reports %s square microns of standard cells and macros. The floorplan targets a utilization of %s. What core area does that require, in square microns? Round to the nearest whole number.", comma(a), f2(u, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Core area = cell area / utilization = %d / %s = %s square microns.", a, f2(u, 2), ans), map[string]string{"core_area": ans})
	case 2:
		a := 10000 * edakit.Int(20, 500)
		r := edakit.Pick([]float64{0.5, 0.75, 1, 1.5, 2})
		w := math.Sqrt(float64(a) / r)
		h := r * w
		if edakit.Pick([]bool{true, false}) {
			ans := edakit.Num(w, 1)
			q := edakit.Sprintf("A core must have an area of %s square microns with aspect ratio $H/W = %s$. What is its width, in microns? Round to 1 decimal place.", comma(a), f2(r, 2))
			return edakit.Numeric(q, ans, edakit.Sprintf("$W = \\sqrt{A/r} = \\sqrt{%d/%s} = %s$ microns.", a, f2(r, 2), ans), map[string]string{"width": ans})
		}
		ans := edakit.Num(h, 1)
		q := edakit.Sprintf("A core must have an area of %s square microns with aspect ratio $H/W = %s$. What is its height, in microns? Round to 1 decimal place.", comma(a), f2(r, 2))
		return edakit.Numeric(q, ans, edakit.Sprintf("$W = \\sqrt{%d/%s} = %s$ microns and $H = rW = %s$ microns.", a, f2(r, 2), f2(w, 2), ans), map[string]string{"height": ans})
	default:
		w, h := 10*edakit.Int(50, 300), 10*edakit.Int(50, 300)
		m := edakit.Pick([]int{50, 80, 100, 120, 150, 200})
		die := float64((w+2*m)*(h+2*m)) / 1e6
		ans := edakit.Num(die, 3)
		q := edakit.Sprintf("A core is %d microns wide and %d microns tall. The I/O ring and seal add a margin of %d microns on every side. What is the die area, in square millimetres? Round to 3 decimal places.", w, h, m)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The die is $(%d + %d)(%d + %d) = %d$ square microns, which is %s square millimetres.", w, 2*m, h, 2*m, (w+2*m)*(h+2*m), ans), map[string]string{"die_um2": edakit.IntStr((w + 2*m) * (h + 2*m))})
	}
}

type pt struct{ x, y int }

func hpwl(ps []pt) int {
	minx, maxx, miny, maxy := ps[0].x, ps[0].x, ps[0].y, ps[0].y
	for _, p := range ps {
		minx, maxx = min(minx, p.x), max(maxx, p.x)
		miny, maxy = min(miny, p.y), max(maxy, p.y)
	}
	return (maxx - minx) + (maxy - miny)
}

func ptsStr(ps []pt) string {
	parts := make([]string, len(ps))
	for i, p := range ps {
		parts[i] = edakit.Sprintf("(%d, %d)", p.x, p.y)
	}
	return strings.Join(parts, ", ")
}

func randPts(n, span int) []pt {
	seen := map[pt]bool{}
	var ps []pt
	for len(ps) < n {
		p := pt{edakit.Int(0, span), edakit.Int(0, span)}
		if !seen[p] {
			seen[p] = true
			ps = append(ps, p)
		}
	}
	return ps
}

// placement: HPWL of one net, total HPWL, and the HPWL change of a cell move.
func placement(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx, 3) {
	case 1:
		ps := randPts(edakit.Scaled(ctx, 3, 7), 30)
		h := hpwl(ps)
		q := "A net connects pins at " + ptsStr(ps) + ", in microns. What is its half-perimeter wirelength, in microns?"
		return edakit.Numeric(q, edakit.IntStr(h),
			edakit.Sprintf("The bounding box width plus height is %d microns.", h), map[string]string{"hpwl": edakit.IntStr(h)})
	case 2:
		n := edakit.Scaled(ctx, 3, 5)
		var lines []string
		total := 0
		for i := 0; i < n; i++ {
			ps := randPts(edakit.Int(2, 4), 20)
			total += hpwl(ps)
			lines = append(lines, edakit.Sprintf("net %d has pins at %s", i+1, ptsStr(ps)))
		}
		q := "In a placement (coordinates in microns), " + strings.Join(lines, "; ") + ". What is the total half-perimeter wirelength of these nets, in microns?"
		return edakit.Numeric(q, edakit.IntStr(total),
			edakit.Sprintf("Summing each net's bounding-box width plus height gives %d microns.", total), map[string]string{"total_hpwl": edakit.IntStr(total)})
	default:
		for {
			old := pt{edakit.Int(0, 20), edakit.Int(0, 20)}
			nw := pt{edakit.Int(0, 20), edakit.Int(0, 20)}
			if old == nw {
				continue
			}
			k := edakit.Scaled(ctx, 2, 3)
			delta := 0
			var lines []string
			for i := 0; i < k; i++ {
				others := randPts(edakit.Int(1, 3), 20)
				before := hpwl(append([]pt{old}, others...))
				after := hpwl(append([]pt{nw}, others...))
				delta += after - before
				lines = append(lines, edakit.Sprintf("net %d also connects pins at %s", i+1, ptsStr(others)))
			}
			if delta == 0 {
				continue
			}
			q := edakit.Sprintf("Cell X sits at (%d, %d) and connects to %d nets: %s. The placer moves X to (%d, %d). By how much does the total HPWL of these nets change, in microns? Give new minus old, so a reduction is negative.",
				old.x, old.y, k, strings.Join(lines, "; "), nw.x, nw.y)
			return edakit.Numeric(q, edakit.IntStr(delta),
				edakit.Sprintf("Only X's nets change. Recomputing each bounding box with X at its new position and subtracting gives %d microns.", delta),
				map[string]string{"delta": edakit.IntStr(delta)})
		}
	}
}

// clockTree: global skew, buffer count of a fanout-limited tree, and skew in timing checks.
func clockTree(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx, 3) {
	case 1:
		n := edakit.Scaled(ctx, 4, 8)
		lat := make([]int, n)
		for i := range lat {
			lat[i] = edakit.Int(350, 460)
		}
		mx, mn := lat[0], lat[0]
		for _, v := range lat {
			mx, mn = max(mx, v), min(mn, v)
		}
		if mx == mn {
			lat[0]++
			mx++
		}
		q := "After clock tree synthesis, the clock latencies of a group of flip-flops are " + joinInts(lat) + " ps. What is the global skew of the group, in picoseconds?"
		return edakit.Numeric(q, edakit.IntStr(mx-mn),
			edakit.Sprintf("Global skew is the maximum minus the minimum latency: $%d - %d = %d$ ps.", mx, mn, mx-mn), map[string]string{"skew": edakit.IntStr(mx - mn)})
	case 2:
		n := edakit.Pick([]int{100, 250, 500, 1000, 2000, 5000, 10000, 40000, 100000})
		k := edakit.Pick([]int{4, 6, 8, 10, 12, 16})
		total, cur := 0, n
		var levels []string
		for {
			cur = (cur + k - 1) / k
			total += cur
			levels = append(levels, edakit.IntStr(cur))
			if cur == 1 {
				break
			}
		}
		q := edakit.Sprintf("A clock must reach %s flip-flops. Each buffer may drive at most %d loads. The leaf level uses the fewest buffers that can drive all sinks, each level above uses the fewest buffers that can drive the level below, up to a single root buffer. How many buffers does the tree contain?", comma(n), k)
		return edakit.Numeric(q, edakit.IntStr(total),
			"Level by level from the leaves: "+strings.Join(levels, ", ")+edakit.Sprintf(" buffers, a total of %d.", total),
			map[string]string{"levels": edakit.IntStr(len(levels))})
	default:
		t := 50 * edakit.Int(16, 30)
		tcq, tsu := edakit.Int(40, 80), edakit.Int(20, 50)
		d := edakit.Int(t*6/10, t-tcq-tsu+40)
		skew := edakit.Int(-40, 40)
		if skew == 0 {
			skew = 15
		}
		slack := t + skew - (tcq + d + tsu)
		dir := "later"
		mag := skew
		if skew < 0 {
			dir, mag = "earlier", -skew
		}
		q := edakit.Sprintf("A path launches at flip-flop A and is captured at flip-flop B. The clock period is %d ps, clock-to-output is %d ps, the logic and wire delay is %d ps and setup is %d ps. After CTS, B's clock arrives %d ps %s than A's. What is the setup slack, in picoseconds?", t, tcq, d, tsu, mag, dir)
		return edakit.Numeric(q, edakit.IntStr(slack),
			edakit.Sprintf("Capture-minus-launch skew is %d ps, so slack $= %d + (%d) - (%d + %d + %d) = %d$ ps.", skew, t, skew, tcq, d, tsu, slack),
			map[string]string{"skew": edakit.IntStr(skew)})
	}
}

// routing: total overflow, gcell edge capacity and routed detour percentage.
func routing(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx, 3) {
	case 1:
		n := edakit.Scaled(ctx, 4, 8)
		var lines []string
		total := 0
		for i := 0; i < n; i++ {
			c := edakit.Int(8, 20)
			d := c + edakit.Int(-6, 5)
			if i == 0 && d <= c {
				d = c + edakit.Int(1, 4)
			}
			total += max(0, d-c)
			lines = append(lines, edakit.Sprintf("edge %d: demand %d, capacity %d", i+1, d, c))
		}
		q := "A global router reports these gcell edges: " + strings.Join(lines, "; ") + ". What is the total overflow?"
		return edakit.Numeric(q, edakit.IntStr(total),
			edakit.Sprintf("Overflow sums $\\max(0, \\text{demand} - \\text{capacity})$ over the edges, giving %d.", total), map[string]string{"overflow": edakit.IntStr(total)})
	case 2:
		pitch := edakit.Pick([]float64{0.1, 0.125, 0.16, 0.2, 0.25, 0.4, 0.5})
		tracks := edakit.Int(20, 100)
		length := float64(tracks) * pitch
		layers := edakit.Int(2, 5)
		ans := tracks * layers
		q := edakit.Sprintf("A gcell boundary is %s microns long. Each routing layer that crosses it has a track pitch of %s microns, and %d layers run in the crossing direction. How many routing tracks cross the boundary?", f2(length, 3), f2(pitch, 3), layers)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("Each layer has $%s/%s = %d$ tracks, so %d layers give %d.", f2(length, 3), f2(pitch, 3), tracks, layers, ans), map[string]string{"tracks_per_layer": edakit.IntStr(tracks)})
	default:
		h := edakit.Int(40, 400)
		r := h + edakit.Int(1, h/3+1)
		ans := edakit.Num(100*float64(r-h)/float64(h), 1)
		q := edakit.Sprintf("A net has a half-perimeter wirelength of %d microns after placement. After detailed routing its wire is %d microns long. By what percentage does the routed length exceed the HPWL? Round to 1 decimal place.", h, r)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$(%d - %d)/%d = %s\\%%$.", r, h, h, ans), map[string]string{"detour": edakit.IntStr(r - h)})
	}
}

// signoff: scenario counts, setup slack with extracted wires, and the effect of wire growth.
func signoff(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx, 3) {
	case 1:
		m, p, rc := edakit.Int(2, 5), edakit.Int(2, 6), edakit.Int(2, 5)
		ans := m * p * rc
		extra := ""
		if ctx.Difficulty > 0.6 {
			ans *= 2
			extra = " Setup and hold are analyzed as separate scenarios."
		}
		q := edakit.Sprintf("A chip is signed off in %d functional modes, at %d process-voltage-temperature corners and %d interconnect RC corners, with every combination analyzed.%s How many timing scenarios must be run?", m, p, rc, extra)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("Scenarios multiply: the product is %d.", ans), map[string]string{"scenarios": edakit.IntStr(ans)})
	case 2:
		t := 50 * edakit.Int(16, 30)
		tcq, tsu := edakit.Int(40, 80), edakit.Int(20, 50)
		wire := edakit.Int(50, 250)
		logic := edakit.Int(t/2, t-tcq-tsu-wire+30)
		slack := t - tcq - logic - wire - tsu
		q := edakit.Sprintf("With extracted parasitics, a register-to-register path has a %d ps clock period, %d ps clock-to-output delay, %d ps of cell delay, %d ps of wire delay and a %d ps setup time. Ignoring skew and uncertainty, what is its setup slack, in picoseconds?", t, tcq, logic, wire, tsu)
		return edakit.Numeric(q, edakit.IntStr(slack),
			edakit.Sprintf("$%d - %d - %d - %d - %d = %d$ ps.", t, tcq, logic, wire, tsu, slack), map[string]string{"slack": edakit.IntStr(slack)})
	default:
		t := 50 * edakit.Int(16, 30)
		tcq, tsu := edakit.Int(40, 80), edakit.Int(20, 50)
		wire := 10 * edakit.Int(5, 25)
		logic := edakit.Int(t/2, t-tcq-tsu-wire)
		pct := edakit.Pick([]int{10, 20, 25, 30, 40, 50, 60})
		nw := wire * (100 + pct) / 100
		before := t - tcq - logic - wire - tsu
		after := t - tcq - logic - nw - tsu
		q := edakit.Sprintf("Before routing, a path was estimated at %d ps clock-to-output, %d ps cell delay, %d ps wire delay and %d ps setup against a %d ps clock. Extraction shows the wire delay is %d%% higher than estimated. What is the setup slack now, in picoseconds?", tcq, logic, wire, tsu, t, pct)
		return edakit.Numeric(q, edakit.IntStr(after),
			edakit.Sprintf("Estimated slack was %d ps. Wire delay becomes %d ps, so slack is $%d - %d - %d - %d - %d = %d$ ps.", before, nw, t, tcq, logic, nw, tsu, after),
			map[string]string{"slack_before": edakit.IntStr(before)})
	}
}

var propChain = []string{"synthesis", "scan insertion", "floorplanning", "placement", "CTS", "routing", "extraction", "signoff"}

// errorPropagation: stages to re-run, hours to re-run, and expected hours with failing attempts.
func errorPropagation(ctx generator.GeneratorContext) generator.Problem {
	hours := make([]int, len(propChain))
	var parts []string
	for i := range hours {
		hours[i] = edakit.Int(1, 12)
		parts = append(parts, edakit.Sprintf("%s %d h", propChain[i], hours[i]))
	}
	k := edakit.Int(0, len(propChain)-2)
	sum := 0
	for i := k; i < len(propChain); i++ {
		sum += hours[i]
	}
	chain := "A flow runs " + strings.Join(propChain, ", ") + " in that order, and a change to a stage invalidates that stage and every later one."
	switch kp(ctx, 3) {
	case 1:
		n := len(propChain) - k
		q := chain + edakit.Sprintf(" An engineer changes the %s step. How many stages must re-run, counting the changed one?", propChain[k])
		return edakit.Numeric(q, edakit.IntStr(n),
			edakit.Sprintf("%s is stage %d of %d, so stages %d through %d re-run: %d stages.", propChain[k], k+1, len(propChain), k+1, len(propChain), n),
			map[string]string{"stages": edakit.IntStr(n)})
	case 2:
		q := chain + " Stage run times are: " + strings.Join(parts, ", ") + edakit.Sprintf(". How many hours of compute does a change to %s cost before the new result is known?", propChain[k])
		return edakit.Numeric(q, edakit.IntStr(sum),
			edakit.Sprintf("Sum the run times from %s through signoff: %d h.", propChain[k], sum), map[string]string{"hours": edakit.IntStr(sum)})
	default:
		qprob := edakit.Pick([]float64{0.25, 0.3, 0.4, 0.5, 0.6, 0.75, 0.8})
		ans := edakit.Num(float64(sum)/qprob, 1)
		q := chain + " Stage run times are: " + strings.Join(parts, ", ") + edakit.Sprintf(". After a change to %s, each re-run attempt closes timing with probability %s, independently. What is the expected number of hours until closure? Round to 1 decimal place.", propChain[k], f2(qprob, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("One attempt costs %d h. Attempts are geometric with mean $1/%s$, so the expected cost is $%d/%s = %s$ h.", sum, f2(qprob, 2), sum, f2(qprob, 2), ans),
			map[string]string{"attempt_hours": edakit.IntStr(sum)})
	}
}
