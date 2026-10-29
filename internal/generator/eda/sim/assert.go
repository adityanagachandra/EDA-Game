package sim

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func bitsRow(v []int) string {
	s := make([]string, len(v))
	for i, x := range v {
		s[i] = edakit.IntStr(x)
	}
	return strings.Join(s, " ")
}

func implication(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	n := edakit.Scaled(ctx, 8, 14)
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		if edakit.Int(0, 99) < 40 {
			a[i] = 1
		}
		if edakit.Int(0, 99) < 60 {
			b[i] = 1
		}
	}
	k := 0
	prop := "a |-> b"
	switch kp {
	case 2:
		k = edakit.Int(1, 3)
		if k == 1 && edakit.Int(0, 1) == 0 {
			prop = "a |=> b"
		} else {
			prop = edakit.Sprintf("a |-> ##%d b", k)
		}
	case 3:
		k = edakit.Int(0, 2)
		if k > 0 {
			prop = edakit.Sprintf("a |-> ##%d b", k)
		}
	}
	m := n - k
	fail, vac, real := 0, 0, 0
	for i := 0; i < m; i++ {
		if a[i] == 0 {
			vac++
		} else if b[i+k] == 0 {
			fail++
		} else {
			real++
		}
	}
	setup := edakit.Sprintf("The property assert property (@(posedge clk) %s); is checked against sampled values over cycles 1 to %d: a = %s and b = %s. An attempt starts at every cycle; consider only attempts starting in cycles 1 to %d, whose checks fall inside the trace.", prop, n, bitsRow(a), bitsRow(b), m)
	facts := map[string]string{"failures": edakit.IntStr(fail), "vacuous": edakit.IntStr(vac), "real_passes": edakit.IntStr(real)}
	when := edakit.Sprintf("%d cycle(s) later", k)
	if k == 0 {
		when = "in the same cycle"
	}
	if kp <= 2 {
		return edakit.Numeric(setup+" How many attempts fail?", edakit.IntStr(fail),
			edakit.Sprintf("An attempt fails when a is 1 at its start cycle and b is 0 %s. That happens %d time(s); %d attempt(s) pass for real and %d pass vacuously.", when, fail, real, vac), facts)
	}
	if edakit.Int(0, 1) == 0 {
		return edakit.Numeric(setup+" How many attempts succeed vacuously?", edakit.IntStr(vac),
			edakit.Sprintf("An attempt whose antecedent a is 0 succeeds vacuously: nothing was checked. There are %d of them, against %d real pass(es) and %d failure(s).", vac, real, fail), facts)
	}
	return edakit.Numeric(setup+" How many attempts pass non-vacuously, with the consequent actually checked?", edakit.IntStr(real),
		edakit.Sprintf("Real passes need a = 1 at the start and b = 1 %s: %d attempt(s). The %d with a = 0 pass vacuously.", when, real, vac), facts)
}

var assertItems = []edakit.Item{
	{KP: 1, Q: "What distinguishes a concurrent assertion from an immediate assertion?",
		A: "A concurrent assertion is clocked and can check behavior across cycles; an immediate one checks a condition when the statement executes", D: []string{"A concurrent assertion runs in a separate operating-system thread", "An immediate assertion can span many cycles", "Concurrent assertions are only for formal tools"},
		Why: "Concurrent assertions describe temporal properties sampled on a clock; immediate ones are procedural checks."},
	{KP: 1, Q: "An immediate assertion inside an always_comb block fires on a glitch that settles within the same time step. What variant avoids the false failure?",
		A: "A deferred immediate assertion (assert #0 or assert final)", D: []string{"A concurrent assertion with no clock", "An assume property", "A cover property"},
		Why: "Deferred assertions report only after the values have settled for the time step, so transient glitches do not fire."},
	{KP: 1, Q: "Where can a concurrent assertion be placed?",
		A: "In a module, interface or program, outside procedural code, or bound into a design with bind", D: []string{"Only inside initial blocks", "Only inside functions", "Only in the testbench top"},
		Why: "Concurrent assertions are structural; bind attaches them to a design without editing its source."},
	{KP: 1, Q: "Which statement checks a property without asserting it, recording each time it is matched?",
		A: "cover property", D: []string{"assume property", "restrict property", "expect"},
		Why: "cover property counts matches and contributes to functional coverage."},
	{KP: 2, Q: "A concurrent assertion is clocked on posedge clk and checks q, which a flop updates at that same edge. Which value of q does it see?",
		A: "The value from just before the edge, sampled in the Preponed region", D: []string{"The new value after the NBA update", "Whichever value the simulator computes first", "An X for one cycle"},
		Why: "Sampling in Preponed makes assertions race-free: they see pre-edge values."},
	{KP: 2, Q: "In which region is the body of a concurrent assertion evaluated, and where do its action blocks run?",
		A: "It is evaluated in the Observed region and its action blocks run in the Reactive region", D: []string{"Both in the Active region", "Evaluated in NBA, actions in Postponed", "Evaluated in Preponed, actions in Active"},
		Why: "Observed evaluates properties on sampled values; Reactive then runs pass and fail actions."},
	{KP: 2, Q: "A waveform shows req rising at the same edge where an assertion reports req was 0. Is the simulator wrong?",
		A: "No: the assertion samples req just before the edge, when it was still 0", D: []string{"Yes: the assertion should see 1", "Yes: assertions are evaluated after NBA updates", "No: assertions always lag by two cycles"},
		Why: "Sampled-value semantics mean an assertion at an edge sees the value before that edge."},
	{KP: 2, Q: "Which SVA sampled-value function (named here without its leading dollar sign) is true when a signal went from 0 to 1 between the previous and current clock edges?",
		A: "rose", D: []string{"stable", "past", "countones"},
		Why: "The rose function compares the current sampled value with the previous sample; stable checks for no change and past returns an earlier sample."},
	{KP: 3, Q: "Why are the same SVA properties valuable in both simulation and formal verification?",
		A: "In simulation they check the tests that run; in formal the tool tries to prove them for all inputs", D: []string{"Formal tools need assertions to run faster", "Simulation cannot evaluate assertions", "They are only for documentation"},
		Why: "One property, two engines: simulation catches violations on the traces it sees, formal searches all traces."},
	{KP: 3, Q: "A team adds thousands of assertions to a large design and simulation slows sharply. Why?",
		A: "Every concurrent assertion starts an attempt on each clock and the simulator must track open threads", D: []string{"Assertions disable compiled simulation", "Assertions force four-state evaluation", "Assertions stop the clock"},
		Why: "Each attempt is evaluation work; long or unbounded sequences multiply the number of live attempts."},
	{KP: 3, Q: "An assertion checks that a FIFO never overflows, but the FIFO never fills during any test. What does its passing mean?",
		A: "Very little: the property was never meaningfully exercised", D: []string{"The FIFO is proven correct", "The assertion is wrong", "The FIFO cannot overflow"},
		Why: "Passing without reaching the condition is like a vacuous pass; pair assertions with cover properties."},
	{KP: 3, Q: "Why do verification teams attach assertions with bind instead of editing RTL files?",
		A: "The checker stays separate from the design source and can be added or removed without touching the RTL", D: []string{"bind makes assertions synthesizable", "bind is faster in simulation", "Assertions in RTL are illegal"},
		Why: "bind instantiates a checker module inside a design instance from outside, keeping ownership separate."},
}
