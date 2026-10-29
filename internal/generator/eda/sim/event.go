package sim

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// drawToggle draws a netlist and an input toggle that produces some activity.
func drawToggle(ctx generator.GeneratorContext, minEvents int) (netlist, []int, int, simResult) {
	for {
		n := randomNetlist(edakit.Int(3, 4), edakit.Scaled(ctx, 4, 7))
		in := inputVector(n.nIn)
		which := edakit.Int(0, n.nIn-1)
		r := n.simulateToggle(in, which)
		if r.events >= minEvents {
			return n, in, which, r
		}
	}
}

func eventDriven(ctx generator.GeneratorContext) generator.Problem {
	n, in, which, r := drawToggle(ctx, 3)
	setup := edakit.Sprintf("A netlist has inputs %s and gates %s. Every gate has a delay of one time unit. The circuit is stable with %s. At time 0 input %s toggles. An event-driven simulator applies the net changes of each time step, evaluates every fanout gate of a changed net once in that step, and schedules an output change one unit later only if the new value differs.",
		strings.Join(n.names[:n.nIn], ", "), n.describe(), n.inputText(in), n.names[which])
	facts := map[string]string{"events": edakit.IntStr(r.events), "evaluations": edakit.IntStr(r.evals)}
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		return edakit.Numeric(setup+" How many events (net value changes) occur in total, counting the input toggle?",
			edakit.IntStr(r.events),
			edakit.Sprintf("Propagating the toggle step by step, nets change %d times in total, the input toggle included, over %d time steps.", r.events, r.steps), facts)
	case 2:
		return edakit.Numeric(setup+" How many gate evaluations does the simulator perform in total?",
			edakit.IntStr(r.evals),
			edakit.Sprintf("Each step evaluates the fanout gates of the nets that changed in it, each gate once per step. Summing over the %d steps gives %d evaluations; a gate whose output does not change still costs an evaluation.", r.steps, r.evals), facts)
	default:
		cnt := 0
		for _, ok := range r.evaluatedOK {
			if ok {
				cnt++
			}
		}
		pct := 100 * float64(cnt) / float64(len(n.gates))
		ans := edakit.Num(pct, 1)
		facts["gates_touched"] = edakit.IntStr(cnt)
		return edakit.Numeric(setup+" What percentage of the gates are evaluated at least once? Round to 1 decimal place.",
			ans,
			edakit.Sprintf("%d of the %d gates lie on an active path and are evaluated at least once, which is %s%%. The rest are never touched, which is where event-driven simulation saves work.", cnt, len(n.gates), ans), facts)
	}
}

func deltaCycles(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	if kp == 1 {
		// A chain of zero-delay continuous assignments.
		k := edakit.Scaled(ctx, 3, 8)
		var lines []string
		prev := "a"
		for i := 1; i <= k; i++ {
			lines = append(lines, edakit.Sprintf("assign w%d = ~%s;", i, prev))
			prev = edakit.Sprintf("w%d", i)
		}
		q := edakit.Sprintf("A module contains the zero-delay continuous assignments %s Input a is changed by the testbench at time 10 ns. Counting the delta cycle in which a itself changes as delta 1, how many delta cycles at time 10 ns contain at least one net change?",
			strings.Join(lines, " "))
		ans := edakit.IntStr(k + 1)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Each assignment reads a net that changed in the previous delta and updates in the next one. The change of a is delta 1, w1 changes in delta 2, and w%d changes in delta %d. No simulated time passes.", k, k+1),
			map[string]string{"deltas": ans})
	}
	n, in, which, r := drawToggle(ctx, 3)
	setup := edakit.Sprintf("A netlist has inputs %s and zero-delay gates %s. It is stable with %s, and input %s toggles at time 50 ns. In each delta cycle the simulator applies pending changes, evaluates the fanout gates of changed nets once, and schedules each changed output for the next delta cycle.",
		strings.Join(n.names[:n.nIn], ", "), n.describe(), n.inputText(in), n.names[which])
	facts := map[string]string{"deltas": edakit.IntStr(r.steps), "events": edakit.IntStr(r.events)}
	if kp == 2 {
		return edakit.Numeric(setup+" Counting the delta cycle that applies the input change as delta 1, how many delta cycles contain at least one net change?",
			edakit.IntStr(r.steps),
			edakit.Sprintf("Tracing the changes delta by delta, the last net change happens in delta %d. That is the longest chain of gates whose outputs actually change, not the depth of the whole netlist.", r.steps), facts)
	}
	// KP 3: glitches in zero time.
	maxNet, maxCh := -1, 0
	for i := n.nIn; i < len(r.changes); i++ {
		if r.changes[i] > maxCh {
			maxNet, maxCh = i, r.changes[i]
		}
	}
	if maxNet < 0 {
		maxNet = n.nIn + len(n.gates) - 1
	}
	// Ask for a randomly chosen gate output, preferring one that changed.
	target := maxNet
	if edakit.Int(0, 1) == 0 {
		target = n.nIn + edakit.Int(0, len(n.gates)-1)
	}
	ans := edakit.IntStr(r.changes[target])
	return edakit.Numeric(setup+edakit.Sprintf(" How many times does net %s change value during time 50 ns, across all delta cycles?", n.names[target]),
		ans,
		edakit.Sprintf("Net %s changes %s time(s). When reconvergent paths deliver their changes in different delta cycles, a net can change, then change back, all at one simulated time: a zero-time glitch that costs real events.", n.names[target], ans), facts)
}

func eventWheel(ctx generator.GeneratorContext) generator.Problem {
	size := edakit.Pick([]int{8, 16, 32, 64, 100, 128})
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		now := edakit.Int(size, 20*size)
		d := edakit.Int(1, size-1)
		slot := (now + d) % size
		q := edakit.Sprintf("A simulator keeps near-future events in a timing wheel of %d slots, one per time unit; an event for absolute time t goes into slot t mod %d. The current time is %d. A gate with delay %d schedules its output change. Which slot index receives the event?", size, size, now, d)
		return edakit.Numeric(q, edakit.IntStr(slot),
			edakit.Sprintf("The event is due at %d + %d = %d, and %d mod %d = %d.", now, d, now+d, now+d, size, slot),
			map[string]string{"due": edakit.IntStr(now + d)})
	case 2:
		base := edakit.Int(1, 50) * size
		m := edakit.Scaled(ctx, 5, 9)
		var ts []string
		over := 0
		for i := 0; i < m; i++ {
			t := base + edakit.Int(0, 2*size)
			if t >= base+size {
				over++
			}
			ts = append(ts, edakit.IntStr(t))
		}
		q := edakit.Sprintf("A timing wheel has %d slots and currently covers absolute times %d through %d. Events beyond the window go to a sorted overflow list. Events are scheduled for times %s. How many of them go to the overflow list?",
			size, base, base+size-1, strings.Join(ts, ", "))
		return edakit.Numeric(q, edakit.IntStr(over),
			edakit.Sprintf("Any event at time %d or later is outside the window. Counting those gives %d; they are moved into the wheel as time advances.", base+size, over),
			map[string]string{"overflow": edakit.IntStr(over)})
	default:
		m := edakit.Scaled(ctx, 6, 12)
		seen := map[int]bool{}
		var ts []int
		for i := 0; i < m; i++ {
			t := edakit.Int(1, 8) * 5
			if len(ts) > 0 && edakit.Int(0, 2) == 0 {
				t = ts[edakit.Int(0, len(ts)-1)]
			}
			ts = append(ts, t)
			seen[t] = true
		}
		strs := make([]string, len(ts))
		for i, t := range ts {
			strs[i] = edakit.IntStr(t)
		}
		q := edakit.Sprintf("The event queue holds %d events with timestamps (ns) %s. The kernel processes all events of one time before advancing to the next nonempty time. How many distinct simulation times does it visit?", m, strings.Join(strs, ", "))
		keys := make([]int, 0, len(seen))
		for k := range seen {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		return edakit.Numeric(q, edakit.IntStr(len(seen)),
			edakit.Sprintf("Events that share a timestamp are processed together. The distinct times are %s ns, so the kernel visits %d times; empty times in between are skipped.", intsText(keys), len(seen)),
			map[string]string{"distinct": edakit.IntStr(len(seen))})
	}
}

// inertialOut simulates a single-input buffer of delay d driven by input changes at times ts
// (alternating 0->1, 1->0, starting from 0). Each input change schedules the output at t+d and
// cancels any pending output event; times are even and d is odd, so nothing coincides.
func inertialOut(ts []int, d int) []int {
	out := 0
	var trans []int
	pending := false
	pt, pv := 0, 0
	for i, t := range ts {
		v := (i + 1) % 2
		if pending && pt < t {
			if pv != out {
				out = pv
				trans = append(trans, pt)
			}
			pending = false
		}
		pending, pt, pv = true, t+d, v
	}
	if pending && pv != out {
		trans = append(trans, pt)
	}
	return trans
}

func pulseTrain(ctx generator.GeneratorContext, n int) []int {
	ts := []int{2 * edakit.Int(1, 5)}
	for i := 1; i < n; i++ {
		ts = append(ts, ts[i-1]+2*edakit.Int(1, 5))
	}
	return ts
}

func delayModels(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	d := 2*edakit.Int(1, 3) + 1
	if kp == 1 {
		start := 2 * edakit.Int(5, 20)
		w := 2 * edakit.Int(1, 5)
		model := edakit.Pick([]string{"inertial", "transport"})
		trans := 2
		if model == "inertial" && w < d {
			trans = 0
		}
		q := edakit.Sprintf("A buffer with an %s delay of %d ns sees its input rise at %d ns and fall at %d ns, a pulse %d ns wide. How many transitions appear on its output?", model, d, start, start+w, w)
		why := "A transport delay shifts every input change by the delay, so both edges appear."
		if model == "inertial" {
			if w < d {
				why = edakit.Sprintf("Under inertial delay the falling edge at %d ns cancels the pending rise due at %d ns, because the pulse is narrower than the %d ns delay. Nothing reaches the output.", start+w, start+d, d)
			} else {
				why = edakit.Sprintf("The rise lands at %d ns before the fall arrives, because the pulse is at least as wide as the %d ns delay, so both edges pass.", start+d, d)
			}
		}
		return edakit.Numeric(q, edakit.IntStr(trans), why, map[string]string{"width": edakit.IntStr(w)})
	}
	m := 2 * edakit.Scaled(ctx, 2, 4)
	ts := pulseTrain(ctx, m)
	trans := inertialOut(ts, d)
	strs := make([]string, len(ts))
	for i, t := range ts {
		strs[i] = edakit.IntStr(t)
	}
	setup := edakit.Sprintf("A buffer has an inertial delay of %d ns: every input change schedules the output to follow %d ns later and cancels any output change still pending. The input starts at 0 and toggles at %s ns.", d, d, strings.Join(strs, ", "))
	tstr := make([]string, len(trans))
	for i, t := range trans {
		tstr[i] = edakit.IntStr(t)
	}
	trace := "none"
	if len(trans) > 0 {
		trace = strings.Join(tstr, ", ") + " ns"
	}
	if kp == 2 || len(trans) == 0 {
		return edakit.Numeric(setup+" How many transitions appear on the output?", edakit.IntStr(len(trans)),
			edakit.Sprintf("Any input segment shorter than %d ns is swallowed, because the next change cancels it before it matures. The surviving output transitions are at %s, so there are %d.", d, trace, len(trans)),
			map[string]string{"transitions": edakit.IntStr(len(trans))})
	}
	last := trans[len(trans)-1]
	return edakit.Numeric(setup+" At what time, in ns, does the last output transition occur?", edakit.IntStr(last),
		edakit.Sprintf("Segments shorter than %d ns are filtered. The output transitions occur at %s, so the last is at %d ns.", d, trace, last),
		map[string]string{"last": edakit.IntStr(last)})
}
