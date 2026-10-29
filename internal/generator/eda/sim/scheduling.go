package sim

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// blockingNonblocking asks for a register value after one clock edge of an always block that
// mixes blocking and nonblocking assignments.
func blockingNonblocking(ctx generator.GeneratorContext) generator.Problem {
	vals := edakit.Shuffle([]int{2, 3, 4, 5, 6, 7, 8, 9})
	A, B, C := vals[0], vals[1], vals[2]
	type tmpl struct {
		code string
		v    string
		res  int
		why  string
	}
	ts := []tmpl{
		{"a = b; b = a;", "b", B, "Blocking assignments update immediately: a takes b first, then b reads the new a, so both end up equal to the old b."},
		{"a = b; b = a;", "a", B, "The first blocking assignment copies the old b into a at once."},
		{"a <= b; b <= a;", "a", B, "Nonblocking right-hand sides are all read in the Active region before any update in the NBA region, so a and b swap."},
		{"a <= b; b <= a;", "b", A, "Both right-hand sides are sampled before the NBA region updates, so b receives the old a: the values swap."},
		{"b = a; c = b;", "c", A, "With blocking assignments b becomes a immediately and c then reads the new b, so the old a falls straight through to c."},
		{"b <= a; c <= b;", "c", B, "Nonblocking assignments read old values, so c receives the old b: this is a two-stage shift register."},
		{"a = a + 1; b = a;", "b", A + 1, "The blocking increment updates a immediately, so b reads the incremented value."},
		{"a <= a + 1; b <= a;", "b", A, "Nonblocking assignments sample a before the increment lands, so b receives the old a."},
		{"a <= b; b = a;", "b", A, "The blocking assignment runs in the Active region before the NBA update of a, so b reads the old a."},
		{"c = a + b; a <= c;", "a", A + B, "The blocking sum updates c at once, and the nonblocking assignment then samples the new c, so a receives a + b."},
	}
	t := edakit.Pick(ts)
	q := edakit.Sprintf("Registers start at a = %d, b = %d and c = %d. On one rising clock edge the block always @(posedge clk) begin %s end runs once. What is the value of %s after the time step completes?", A, B, C, t.code, t.v)
	d := edakit.Distinct(t.res, 3, A, B, C, A+1, A+B)
	return edakit.Choice(q, edakit.IntStr(t.res), t.why, d...)
}

var schedulingItems = []edakit.Item{
	{KP: 1, Q: "Within one simulation time step of an IEEE 1800 simulator, which order do these regions run in?",
		A: "Active, Inactive, NBA, Observed, Reactive, Postponed", D: []string{"Active, NBA, Inactive, Reactive, Observed, Postponed", "Preponed, NBA, Active, Observed, Postponed, Reactive", "Inactive, Active, NBA, Reactive, Observed, Postponed"},
		Why: "Active runs processes and blocking assignments, Inactive holds #0 events, NBA applies nonblocking updates, Observed evaluates assertions, Reactive runs program and assertion action code, and Postponed is the read-only end of the step."},
	{KP: 1, Q: "A process executes #0 x = 1;. In which region does the assignment to x take effect?",
		A: "The Inactive region of the current time step", D: []string{"The NBA region of the current time step", "The Active region of the next time step", "The Postponed region"},
		Why: "A zero delay suspends the process into the Inactive region, which runs after the Active region empties but in the same time step."},
	{KP: 1, Q: "Where does a nonblocking assignment q <= d evaluate d and where does it update q?",
		A: "It evaluates d in the Active region and updates q in the NBA region", D: []string{"It evaluates and updates both in the NBA region", "It evaluates d in the Preponed region and updates q in the Active region", "It evaluates d in the Observed region and updates q in Postponed"},
		Why: "The right-hand side is sampled when the statement executes (Active); the left-hand update is queued for the NBA region."},
	{KP: 1, Q: "The monitor system task prints a signal that is updated by nonblocking assignments. Why does it always print the final value of the time step?",
		A: "The monitor and strobe tasks report in the Postponed region, after all updates", D: []string{"The monitor task runs in the Preponed region", "The monitor task forces the NBA region to run first", "The monitor task delays printing by one time step"},
		Why: "Postponed is the last region and is read-only, so it sees the settled values of the step."},
	{KP: 1, Q: "After the NBA region applies updates, those updates wake further always blocks. What does the scheduler do?",
		A: "It returns to the Active region in the same time step, which is another delta cycle", D: []string{"It advances simulation time by one unit", "It discards the new events until the next clock edge", "It runs them in the Postponed region"},
		Why: "Regions iterate: new Active events after NBA cause another pass through the region loop at the same simulated time."},
	{KP: 1, Q: "Concurrent assertions must not see the half-updated values of a time step. Which region samples their operands?",
		A: "Preponed, before anything in the time step changes", D: []string{"NBA, after register updates", "Observed, at the moment they are evaluated", "Postponed, after the time step settles"},
		Why: "Operands are sampled in Preponed, so an assertion clocked on an edge sees the values from just before that edge."},
	{KP: 3, Q: "Two always @(posedge clk) blocks share a clock: one does q1 = d; the other does q2 = q1;. What is wrong?",
		A: "It is a race: the result depends on which block the simulator runs first", D: []string{"Nothing: blocking assignments always execute in source order across blocks", "q2 always gets the new q1", "The simulator reports a compile error"},
		Why: "The standard leaves the order of processes in the Active region unspecified. Nonblocking assignments remove the race."},
	{KP: 3, Q: "A testbench drives a DUT input with a blocking assignment exactly at the rising clock edge the DUT samples. What can go wrong?",
		A: "The DUT may sample either the old or the new value, depending on process ordering", D: []string{"The DUT always samples the new value", "The DUT always samples the old value", "The assignment is deferred to the next time step"},
		Why: "Driving and sampling in the same Active region is a race. Driving after the edge, with nonblocking assignments or a clocking block, makes it deterministic."},
	{KP: 3, Q: "Code uses #0 delays to 'fix' an ordering problem between two processes. Why is that fragile?",
		A: "Another #0 elsewhere can reintroduce the same race one region later", D: []string{"#0 is illegal in synthesizable code and fails simulation", "#0 advances time by one unit", "#0 forces the NBA region to run first"},
		Why: "#0 only moves a process into the Inactive region; two processes both using #0 race again there."},
	{KP: 3, Q: "Two simulators give different waveforms for the same RTL and the same seed. The RTL uses blocking assignments to registers in several clocked blocks. What is the likely cause?",
		A: "A race whose outcome depends on each simulator's process order", D: []string{"One simulator has a broken random generator", "A floating-point rounding difference", "One simulator ignores the clock"},
		Why: "Both simulators are LRM-compliant; the code has a race, so different legal orders give different results."},
	{KP: 3, Q: "A combinational always @* block reads a signal that it also writes later in the same block with a blocking assignment. What risk does this pose?",
		A: "The block can retrigger itself and loop through delta cycles, or read a stale value", D: []string{"The block becomes sequential and adds a clock cycle", "The nonblocking region runs twice", "Nothing: always @* never retriggers"},
		Why: "Writing a signal on the sensitivity list can wake the block again in the next delta, which is a zero-time loop risk."},
	{KP: 3, Q: "Why do clocking blocks in a SystemVerilog testbench remove most drive-sample races?",
		A: "They sample inputs at a skew before the edge and drive outputs at a skew after it", D: []string{"They run the testbench in the Preponed region only", "They convert every assignment into a blocking one", "They disable delta cycles in the DUT"},
		Why: "Input skew samples settled values and output skew drives after the edge, so the DUT and testbench never touch a signal in the same region."},
}
