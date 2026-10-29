package sta

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func ioDelay(ctx generator.GeneratorContext) generator.Problem {
	t := 50 * edakit.Int(10, 30)
	clk := edakit.Sprintf("create_clock -period %d sets the clock in ps, and clocks are ideal. ", t)
	switch ctx.KPOrdinal {
	case 1:
		x := 10 * edakit.Int(10, t/20)
		d := 10 * edakit.Int(10, t/20)
		su := edakit.Int(20, 60)
		s := t - su - (x + d)
		return edakit.Numeric(clk+edakit.Sprintf("Input port din has set_input_delay %d -clock clk. From din, logic of %d ps reaches a flop with setup %d ps. What is the setup slack of this input path, in ps?", x, d, su), edakit.IntStr(s),
			edakit.Sprintf("The input delay models the external launch: data reaches din at %d ps. Slack = %d - %d - (%d + %d) = %d ps.", x, t, su, x, d, s), map[string]string{"slack": edakit.IntStr(s)})
	case 3:
		x := 10 * edakit.Int(10, t/25)
		y := 10 * edakit.Int(10, t/25)
		b := t - x - y
		return edakit.Numeric(clk+edakit.Sprintf("Port a has set_input_delay %d -clock clk and port z has set_output_delay %d -clock clk. A purely combinational path runs from a to z. What is the largest delay that path may have and still meet setup, in ps?", x, y), edakit.IntStr(b),
			edakit.Sprintf("The external launch uses %d ps and the external capture needs %d ps of the %d ps cycle, leaving %d ps.", x, y, t, b), map[string]string{"budget": edakit.IntStr(b)})
	}
	y := 10 * edakit.Int(10, t/20)
	tcq, d := edakit.Int(30, 80), 10*edakit.Int(10, t/20)
	s := t - y - (tcq + d)
	return edakit.Numeric(clk+edakit.Sprintf("Output port dout has set_output_delay %d -clock clk. A flop with clock-to-Q %d ps drives dout through %d ps of logic. What is the setup slack at dout, in ps?", y, tcq, d), edakit.IntStr(s),
		edakit.Sprintf("The output delay is the external logic plus setup outside the block, so dout is required at %d - %d = %d ps. Slack = %d - (%d + %d) = %d ps.", t, y, t-y, t-y, tcq, d, s), map[string]string{"slack": edakit.IntStr(s)})
}

func multicycle(ctx generator.GeneratorContext) generator.Problem {
	t := 50 * edakit.Int(6, 16)
	n := edakit.Int(2, 4)
	tcq := edakit.Int(30, 70)
	d := 10 * edakit.Int(t/10, (n*t)/10-10)
	dmin := edakit.Int(20, 150)
	su, th := edakit.Int(20, 50), edakit.Int(10, 40)
	q := edakit.Sprintf("A path between two flops on clock clk (period %d ps, ideal) has clock-to-Q %d ps, longest logic delay %d ps and shortest %d ps; setup is %d ps and hold %d ps. The SDC has set_multicycle_path %d -setup -from [get_cells FF1] -to [get_cells FF2]",
		t, tcq, d, dmin, su, th, n)
	switch ctx.KPOrdinal {
	case 1:
		s := n*t - su - (tcq + d)
		return edakit.Numeric(q+". What is the setup slack, in ps?", edakit.IntStr(s),
			edakit.Sprintf("Setup now captures at edge %d, at %d ps. Slack = %d - %d - (%d + %d) = %d ps.", n, n*t, n*t, su, tcq, d, s), map[string]string{"setup_slack": edakit.IntStr(s)})
	case 3:
		m := edakit.Int(1, n-1)
		edge := (n - 1 - m) * t
		h := tcq + dmin - (edge + th)
		return edakit.Numeric(q+edakit.Sprintf(", plus set_multicycle_path %d -hold on the same path. What is the hold slack, in ps?", m), edakit.IntStr(h),
			edakit.Sprintf("Moving setup to edge %d moves the default hold check to %d ps. The -hold %d moves it %d cycles earlier, to %d ps. Slack = (%d + %d) - (%d + %d) = %d ps.", n, (n-1)*t, m, m, edge, tcq, dmin, edge, th, h), map[string]string{"hold_edge": edakit.IntStr(edge), "hold_slack": edakit.IntStr(h)})
	}
	edge := (n - 1) * t
	h := tcq + dmin - (edge + th)
	return edakit.Numeric(q+", and no hold multicycle. What is the hold slack, in ps?", edakit.IntStr(h),
		edakit.Sprintf("The default hold check sits one capture edge before the setup edge, at %d ps. Slack = (%d + %d) - (%d + %d) = %d ps, which is why a -hold %d normally accompanies -setup %d.", edge, tcq, dmin, edge, th, h, n-1, n), map[string]string{"hold_edge": edakit.IntStr(edge), "hold_slack": edakit.IntStr(h)})
}

var sdcItems = []edakit.Item{
	{KP: 1, Q: "Which SDC command defines a clock with a period and waveform on a port?",
		A: "create_clock", D: []string{"create_generated_clock", "set_clock_latency", "set_input_delay"},
		Why: "create_clock names a clock source, its period and its edge times. Generated clocks are derived from a master with create_generated_clock."},
	{KP: 1, Q: "A divide-by-2 flop's output clocks half the design. How should it be constrained?",
		A: "create_generated_clock on the divider output with -source and -divide_by 2", D: []string{"A second create_clock with twice the period on the divider output", "set_multicycle_path 2 on every path it clocks", "set_clock_latency equal to the divider delay"},
		Why: "A generated clock keeps its phase relationship to the master and lets the timer include the path through the divider in its latency. An independent create_clock would lose that relationship."},
	{KP: 1, Q: "What does set_clock_latency -source 200 [get_clocks clk] model?",
		A: "Delay from the clock's origin, such as an off-chip oscillator or PLL, to its definition point", D: []string{"The delay through the on-chip clock tree", "Extra jitter on the clock", "A delay that applies only to hold checks"},
		Why: "Source latency is the delay before the clock definition point. Network latency is the on-chip tree, given ideally with set_clock_latency or computed with set_propagated_clock."},
	{KP: 1, Q: "Before clock tree synthesis, how should network latency normally be modeled?",
		A: "As an ideal latency estimate, switching to set_propagated_clock after CTS", D: []string{"As propagated clocks through the unbuffered clock net", "As zero, with no uncertainty", "With set_input_delay on every flop"},
		Why: "Pre-CTS the clock net is one huge unbuffered net; propagating through it gives meaningless numbers. An ideal estimate plus uncertainty stands in until the tree exists."},
	{KP: 2, Q: "set_input_delay 300 -clock clk [get_ports din] tells the timer what?",
		A: "Data at din arrives 300 time units after the clk edge, from logic outside the block", D: []string{"din must be stable 300 units before the clk edge", "Paths from din have 300 units of internal delay", "din is a false path for 300 units"},
		Why: "Input delay is the arrival at the port relative to a clock edge, standing in for an external launch flop plus its logic."},
	{KP: 2, Q: "set_output_delay 400 -clock clk [get_ports dout] tells the timer what?",
		A: "External logic and setup consume 400 units, so dout must arrive 400 units before the capturing edge", D: []string{"dout arrives 400 units after the clock edge", "dout has a 400 unit hold requirement", "The path to dout is a multicycle path of 400 units"},
		Why: "Output delay is the required time margin outside the block: the port's required time is the capturing edge minus the output delay."},
	{KP: 2, Q: "Which commands give an input port a realistic transition time instead of an ideal step?",
		A: "set_driving_cell or set_input_transition", D: []string{"set_load", "set_input_delay", "set_clock_transition"},
		Why: "set_driving_cell models the driver's slew and drive strength; set_input_transition gives a fixed slew. set_load is for the capacitance on output ports."},
	{KP: 2, Q: "An output port drives a board trace and a receiver. Which command models the capacitance it sees?",
		A: "set_load", D: []string{"set_driving_cell", "set_output_delay", "set_max_transition"},
		Why: "set_load adds capacitance to a port's net, which raises the delay and slew of the cell driving it."},
	{KP: 3, Q: "Why do SDC commands select objects with get_ports, get_pins, get_cells and get_clocks?",
		A: "The same name can denote a port, a pin, a cell or a clock, so each command must say which object type it means", D: []string{"They make the commands run faster", "They are needed only for hierarchical designs", "They convert names to the Liberty format"},
		Why: "SDC is Tcl applied to a design database. Object queries resolve names unambiguously, and a misspelled query that matches nothing silently drops a constraint."},
	{KP: 3, Q: "A script sets set_clock_uncertainty -setup 50 and -hold 20 on clk. What do the two numbers do?",
		A: "Setup checks lose 50 units of margin and hold checks gain a 20 unit larger requirement", D: []string{"Both add to every setup check", "The hold value is subtracted from setup", "They change the clock period"},
		Why: "Setup uncertainty pulls the required time earlier, and hold uncertainty raises the hold requirement. They model jitter and, before CTS, expected skew."},
	{KP: 3, Q: "In which units are SDC time values read?",
		A: "The time unit set by the design's Liberty or set_units, often ns", D: []string{"Always ps", "Always seconds", "The unit of the SPEF file"},
		Why: "SDC numbers are unitless; the timer interprets them in the library or session time unit. A ps-versus-ns mismatch is a classic constraint bug."},
	{KP: 3, Q: "A block's SDC has no constraint on an input port that feeds flops. How does a timer treat those paths?",
		A: "They are unconstrained and never reported as violations", D: []string{"They are timed against a full clock period from time zero", "They are treated as false paths with a warning only when they fail", "They inherit the output delay of the block"},
		Why: "Without an input delay there is no launch reference for the port, so the paths are unconstrained. Constraint checks such as check_timing flag them."},
}

var exceptionItems = []edakit.Item{
	{KP: 1, Q: "A configuration register is written once at boot and then held constant while the datapath runs. Its paths into the datapath fail setup. What is the right constraint?",
		A: "set_false_path from the register, after confirming it is static in operation", D: []string{"set_multicycle_path 2", "set_clock_uncertainty 0", "set_max_delay equal to one period"},
		Why: "A path that never propagates a transition in functional operation is false. A multicycle path would still time it against a later edge."},
	{KP: 1, Q: "What is the main danger of set_false_path?",
		A: "A path that is in fact real becomes invisible to every report", D: []string{"It doubles run time", "It removes hold checks only", "It changes the clock period"},
		Why: "False paths delete timing checks. A wrong one hides a real failure until silicon, so every exception needs a reason that can be reviewed."},
	{KP: 1, Q: "Two clocks come from unrelated oscillators, and all crossings go through synchronizers. How should they be constrained?",
		A: "set_clock_groups -asynchronous between them", D: []string{"set_clock_groups -logically_exclusive", "set_multicycle_path 2 on every crossing", "A shared create_clock definition"},
		Why: "Asynchronous groups tell the timer that no phase relationship exists, so crossings are not timed. The synchronizers handle metastability."},
	{KP: 1, Q: "A mux selects either clkA or clkB onto one clock net, never both at once. Which constraint fits?",
		A: "set_clock_groups -logically_exclusive between clkA and clkB", D: []string{"set_clock_groups -asynchronous", "set_false_path from clkA to clkA", "A single clock with the faster period"},
		Why: "Exclusive clocks cannot coexist, so paths between them are not real. Asynchronous would also stop timing but misdescribes the design and changes crosstalk analysis."},
	{KP: 2, Q: "A multiplier's result is sampled only every second cycle, controlled by an enable. What exception fits?",
		A: "set_multicycle_path 2 -setup, with set_multicycle_path 1 -hold", D: []string{"set_false_path", "set_multicycle_path 2 -setup only", "set_max_delay of half a period"},
		Why: "The data has two cycles to settle, so setup moves to the second edge. Without the -hold 1, the hold check moves forward a cycle too and fails."},
	{KP: 2, Q: "What does set_max_delay 2.0 -from A -to B do that a multicycle path cannot?",
		A: "It sets an absolute delay limit, independent of the clock edges", D: []string{"It removes the hold check", "It applies only to clock paths", "It makes the path asynchronous"},
		Why: "Max delay replaces the edge-derived required time with an absolute number. It is common on asynchronous crossings that still need a bounded delay."},
	{KP: 2, Q: "A false path and a multicycle path both match the same path. Which wins?",
		A: "The false path, because it has the highest precedence", D: []string{"The multicycle path, because it is more specific", "Whichever was written last", "Neither: the timer reports an error and times it normally"},
		Why: "Exceptions are ranked by type: false path above max and min delay, which are above multicycle."},
	{KP: 2, Q: "A multicycle path is written with -through [get_pins U5/A]. Which paths does it affect?",
		A: "Only paths that pass through pin U5/A", D: []string{"All paths into U5", "All paths from the clock that drives U5", "All paths in the fanout cone of U5 regardless of pin"},
		Why: "-through filters on the pins a path visits. Writing exceptions at the narrowest point avoids catching unrelated paths."},
	{KP: 3, Q: "Why are false paths on functional modes not removed in test mode?",
		A: "Exceptions are mode-specific: a path false in mission mode may be real during scan shift", D: []string{"Test mode ignores SDC entirely", "Test mode runs only hold checks", "False paths are global by definition"},
		Why: "Each mode has its own SDC. Scan shift exercises paths the functional mode never uses, so exceptions must be justified per mode."},
	{KP: 3, Q: "A single-bit signal crosses from a 100 MHz domain to an unrelated 75 MHz domain through a two-flop synchronizer. What should the timer still check?",
		A: "Nothing across the crossing, but the synchronizer flops should be placed close together and the crossing bounded if the design needs it", D: []string{"A normal setup check against the 75 MHz edge", "A multicycle path of 2 cycles", "Recovery and removal on the first flop"},
		Why: "Asynchronous crossings have no valid single-cycle check. Designers declare them asynchronous and, where needed, bound the delay with set_max_delay."},
	{KP: 3, Q: "An engineer adds set_multicycle_path 3 -setup to fix a failing path but does not prove the logic waits three cycles. What is the risk?",
		A: "The chip can fail in silicon when data is actually sampled after one cycle", D: []string{"The timer reports more hold violations but silicon is fine", "Run time triples", "Nothing: multicycle paths are always safe"},
		Why: "Exceptions are claims about the design's behavior. If the claim is false the timing report is optimistic and the failure moves to silicon."},
	{KP: 3, Q: "Why do signoff teams prefer clock groups to many pairwise false paths between two unrelated clocks?",
		A: "One clock-group statement covers every crossing in both directions and states the intent directly", D: []string{"Clock groups are faster to time but less accurate", "False paths cannot be written between clocks", "Clock groups also fix metastability"},
		Why: "Pairwise false paths are easy to miss in one direction. set_clock_groups captures the relationship once, and timers also use it in crosstalk analysis."},
}
