package csim

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// multirate covers multirate stepping, latency exploitation and partitioned solves.
func multirate(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		b := edakit.Scaled(ctx, 2, 5)
		hs := make([]float64, b)
		var parts []string
		hmin := math.Inf(1)
		for i := range hs {
			hs[i] = pickF(1, 2, 5, 10, 20, 50, 100)
			hmin = math.Min(hmin, hs[i])
			parts = append(parts, edakit.Sprintf("%s ps", f(hs[i], 0)))
		}
		sum := 0.0
		for _, h := range hs {
			sum += hmin / h
		}
		sp := float64(b) / sum
		return num(edakit.Sprintf("A fast SPICE engine partitions a circuit into %d blocks of similar cost whose accuracy needs steps of %s respectively. A single-rate simulator steps every block at the smallest step; a multirate one steps each block at its own. What speedup in block solves does multirate give over a fixed run time?", b, strings.Join(parts, ", ")), sp, 2,
			edakit.Sprintf("Single-rate work is $%d \\cdot T/h_{min}$ and multirate work is $\\sum T/h_i$. Their ratio is $%d / %s = %s$, ignoring the cost of interpolating signals between blocks.", b, b, f(sum, 4), f(sp, 2)), nil)
	case 2:
		n := edakit.Pick([]int{1000, 2000, 5000, 10000})
		active := edakit.Pick([]int{20, 50, 100, 200, 500})
		cost := pickF(0.01, 0.02, 0.05)
		sp := float64(n) / (float64(active) + float64(n-active)*cost)
		return num(edakit.Sprintf("Of %d partitions, on average %d are active at each timestep; the rest are latent. Fully evaluating a partition costs 1 unit, and checking that a latent partition's inputs have not changed costs %s units. What speedup does skipping latent partitions give over evaluating all of them?", n, active, f(cost, 2)), sp, 2,
			edakit.Sprintf("Work per step drops from %d to $%d + %d \\times %s = %s$ units, a speedup of %s.", n, active, n-active, f(cost, 2), f(float64(active)+float64(n-active)*cost, 2), f(sp, 2)), nil)
	}
	n := edakit.Pick([]int{1200, 2400, 3000, 6000})
	b := edakit.Pick([]int{4, 6, 10, 12})
	s := edakit.Pick([]int{50, 100, 200})
	full := math.Pow(float64(n), 3)
	part := float64(b)*math.Pow(float64(n)/float64(b), 3) + math.Pow(float64(s), 3)
	return num(edakit.Sprintf("Suppose each matrix were factored as dense, with cost proportional to size cubed. A %d-unknown circuit is split into %d equal blocks plus a %d-unknown interface system solved after them. What is the ratio of the unpartitioned cost to the partitioned cost?", n, b, s), full/part, 1,
		edakit.Sprintf("$%d^3$ against $%d \\cdot %d^3 + %d^3$: the ratio is %s. Real circuit matrices are sparse, so true gains are smaller, but the interface size still bounds the benefit.", n, b, n/b, s, f(full/part, 1)), nil)
}

var tableModelItems = []edakit.Item{
	{KP: 1, Q: "A fast SPICE tool replaces BSIM evaluation with a table of drain current over a grid of terminal voltages. What does it gain?",
		A: "Much cheaper device evaluation per Newton iteration", D: []string{"Exact agreement with BSIM at every bias", "Smaller MNA matrices", "Removal of the need for timesteps"},
		Why: "Interpolating a precomputed table costs a few multiply-adds, whereas a compact model evaluates many equations per call."},
	{KP: 1, Q: "A table model uses piecewise-linear interpolation of current. What problem does that cause for Newton?",
		A: "The derivative jumps at grid lines, so the Jacobian is discontinuous and convergence can stall", D: []string{"The current becomes negative everywhere", "The matrix becomes singular", "None, since Newton needs only the current"},
		Why: "Newton uses conductances, the derivatives. A piecewise-linear table has piecewise-constant derivatives, which can make iterations oscillate across grid lines."},
	{KP: 1, Q: "Why do table models often use smooth splines rather than linear interpolation?",
		A: "Continuous first derivatives give continuous conductances, which Newton and small-signal analysis need", D: []string{"Splines need less memory", "Splines are exact for BSIM", "Linear interpolation cannot be stored"},
		Why: "A cubic spline is continuous with continuous first derivative, so the stamped conductances vary smoothly with bias."},
	{KP: 1, Q: "A table model is built from SPICE runs at one temperature. What must the fast simulator do for another temperature?",
		A: "Use another table (or an extra table dimension) built at that temperature", D: []string{"Scale every current by the absolute temperature ratio", "Nothing, since tables are temperature independent", "Reuse the table and rely on gmin"},
		Why: "Device behavior depends on temperature through mobility, threshold and more. A table is only valid where it was characterized."},
	{KP: 2, Q: "Fast SPICE tools relax Newton and timestep tolerances on digital-like partitions. Why is this often acceptable there?",
		A: "Digital nodes swing rail to rail, so small voltage errors rarely change logic timing or power much", D: []string{"Digital nodes have no capacitance", "Relaxed tolerances make results exact", "Digital partitions are never stiff"},
		Why: "The useful output of a digital partition is edge timing and current. Errors of a few millivolts are small compared with the swing."},
	{KP: 2, Q: "A memory array has 64K identical bit cells, most in the same state. How do fast SPICE engines avoid simulating each one?",
		A: "Isomorphism: identical subcircuits with identical inputs share one solution", D: []string{"They delete the cells", "They use AC analysis for the array", "They increase gmin on every cell"},
		Why: "Recognizing structurally identical partitions with the same boundary conditions lets the engine solve once and reuse the result."},
	{KP: 2, Q: "A post-layout netlist has millions of parasitic resistors and capacitors. What does a fast simulator typically do before simulation?",
		A: "Reduce the RC networks to smaller equivalents that preserve their low-frequency behavior", D: []string{"Remove all parasitics", "Convert them to inductors", "Simulate them with fixed steps of one femtosecond"},
		Why: "RC reduction, by node elimination or moment-matching projection, shrinks the parasitic network while keeping its delay behavior."},
	{KP: 2, Q: "A fast SPICE run and a SPICE run differ by 3% in the delay of a critical path. Which knob most directly trades this accuracy against runtime?",
		A: "The accuracy level or speed mode, which sets tolerances, table resolution and partitioning aggressiveness", D: []string{"The netlist node names", "The output file format", "The number of .print statements"},
		Why: "Fast SPICE tools expose modes that bundle tolerances and model choices. Moving toward accuracy shrinks the error at higher cost."},
	{KP: 3, Q: "A team wants the input offset of a comparator to within 0.1 mV. Should they use a fast SPICE mode?",
		A: "No: offset depends on small differences in device currents that table models and relaxed tolerances blur", D: []string{"Yes: fast SPICE is designed for offset", "Yes, if the clock is fast", "Only with gmin stepping"},
		Why: "Precision analog metrics need accurate device models and tight tolerances. Fast SPICE suits large digital and mixed blocks where small analog errors do not matter."},
	{KP: 3, Q: "Leakage current of a sleeping SRAM is the sign-off metric. What risk does a table-based fast model carry?",
		A: "Subthreshold currents are tiny and exponential, so table interpolation errors can be large in relative terms", D: []string{"Leakage cannot be computed at all", "Leakage is independent of the device model", "Tables always overestimate leakage by exactly 2x"},
		Why: "Tables tuned for on-current accuracy can be poorly resolved in the subthreshold region, where currents change by decades over a few hundred millivolts."},
	{KP: 3, Q: "A fast SPICE result is used for full-chip power, and three critical paths are re-run in accurate SPICE. What is this practice?",
		A: "Using the fast engine broadly and recursing into the accurate engine where the decision depends on precision", D: []string{"Redundant work with no purpose", "A replacement for all sign-off", "A way to avoid any accurate simulation"},
		Why: "The fast engine covers scale; the accurate engine confirms the places where the fast engine's error budget could change the answer."},
	{KP: 3, Q: "A fast SPICE engine's partitioner splits a tightly coupled analog feedback loop across two partitions with relaxed coupling. What can go wrong?",
		A: "The loop dynamics, such as its stability, are distorted because the coupling is treated as weak", D: []string{"Nothing, since partitions are always exact", "The DC point becomes exact", "The netlist cannot be parsed"},
		Why: "Partitioned methods work well when coupling between partitions is weak or unidirectional. Strong bidirectional feedback across a boundary needs a joint solve."},
}

var mixedSignalItems = []edakit.Item{
	{KP: 1, Q: "An analog node drives a digital gate input in a mixed-signal simulation. What converts the analog waveform into a logic value?",
		A: "A connect module (an A-to-D interface) with thresholds, often with hysteresis", D: []string{"The digital simulator samples the node at every analog timestep", "A DC operating point at each clock edge", "The analog solver rounds the voltage to 0 or 1"},
		Why: "Connect modules define when an analog voltage counts as 0, 1 or X, typically with separate rising and falling thresholds to avoid chatter."},
	{KP: 1, Q: "A digital output drives an analog node. What should the D-to-A interface present to the analog solver?",
		A: "A voltage or current source with finite rise time and output resistance", D: []string{"An ideal voltage step with zero rise time", "Nothing; the analog node keeps its old voltage", "A logic value stamped into the MNA matrix"},
		Why: "Ideal steps create discontinuities that the analog solver handles badly. A finite slope and impedance give a well-posed, realistic drive."},
	{KP: 1, Q: "A slowly rising analog input crosses a single threshold while noise rides on it. Without hysteresis in the A-to-D interface, what does the digital side see?",
		A: "Many spurious events as the signal crosses back and forth", D: []string{"One clean event", "No events at all", "An X state for the whole run"},
		Why: "Small wiggles around a single threshold produce repeated crossings. Hysteresis requires the signal to move past a second threshold before switching back."},
	{KP: 1, Q: "Why do mixed-signal simulators insert interface elements automatically based on port disciplines?",
		A: "Each net must be resolved consistently in one domain, and the interface defines the conversion at each boundary", D: []string{"Because digital simulators cannot read netlists", "To add parasitics", "To randomize the simulation"},
		Why: "Disciplines (electrical, logic) say what a port is. When a net joins different disciplines, a connect module is placed at the boundary to translate values."},
	{KP: 2, Q: "The digital engine schedules an event at 12.3 ns that changes an analog source. What must the analog solver do?",
		A: "Place a breakpoint at 12.3 ns so its time grid hits the event exactly", D: []string{"Ignore it until the next accepted step", "Restart from the DC operating point", "Halve all future timesteps"},
		Why: "Digital events act like source corners for the analog side. Treating them as breakpoints keeps the discontinuity aligned with a time point."},
	{KP: 2, Q: "The analog solver takes a step to 20 ns, but during that step an analog output crossed a threshold at 18.7 ns. What must happen?",
		A: "The crossing is located, the analog side backs up, and the digital engine is notified at 18.7 ns", D: []string{"The digital engine sees the event at 20 ns", "Nothing, since crossings are only checked at breakpoints", "The analog step is accepted and the event is dropped"},
		Why: "Accurate synchronization needs threshold crossings at the right time. The analog step is cut or interpolated to report the event when it happened."},
	{KP: 2, Q: "In lockstep synchronization, how do the analog and digital engines advance?",
		A: "Neither advances past the next point where it may need input from the other", D: []string{"Each runs the entire simulation independently and the results are merged", "The digital engine always waits for the analog to finish", "They exchange data only at the end"},
		Why: "Lockstep keeps both engines at consistent times, advancing to the next event or analog time point, which is safe but limits each engine's step choice."},
	{KP: 2, Q: "Waveform relaxation lets coupled subsystems simulate a whole window independently, exchange waveforms, and repeat. When does it converge quickly?",
		A: "When coupling between the subsystems is weak or mostly one-directional", D: []string{"When the subsystems are tightly coupled by feedback", "Always in one iteration", "Only when every subsystem is linear and identical"},
		Why: "Each relaxation sweep uses the other subsystem's previous waveform. Weak coupling means those waveforms change little between sweeps."},
	{KP: 3, Q: "A verification team co-simulates an RTL testbench in a digital simulator with an analog block in a SPICE engine. Where does most runtime typically go?",
		A: "The analog engine, which solves nonlinear systems at small timesteps while the digital side processes events cheaply", D: []string{"The digital engine, because RTL is slower than SPICE", "The interface modules", "Parsing the netlists"},
		Why: "Event-driven RTL simulation is fast per event. The analog part pays for Newton solves at every time point, so its size and accuracy dominate the cost."},
	{KP: 3, Q: "To speed up a chip-level mixed-signal regression, a team swaps a PLL's transistor netlist for a real-number model. What is a real-number model?",
		A: "A discrete-time model in the digital simulator that passes real values, such as voltages, between blocks as events", D: []string{"A SPICE model with real rather than complex numbers", "A Verilog-A model solved by Newton", "A model that uses only integer arithmetic"},
		Why: "Real-number modeling keeps analog behavior approximate but runs entirely in the event-driven simulator, avoiding the analog solver."},
	{KP: 3, Q: "ngspice includes XSPICE. What does XSPICE add to the simulator?",
		A: "Code models, including event-driven digital models, that run alongside the analog circuit", D: []string{"A harmonic-balance engine", "A GPU solver", "A layout extractor"},
		Why: "XSPICE extends SPICE with user-definable code models and an event-driven algorithm for digital and mixed-mode simulation."},
	{KP: 3, Q: "A co-simulation shows a digital signal toggling one event later than in an all-SPICE reference run. Which interface setting is the first suspect?",
		A: "The A-to-D threshold and hysteresis settings that decide when the analog waveform becomes a logic event", D: []string{"The analog engine's gmin", "The number of digital signals", "The waveform file format"},
		Why: "Event timing at the boundary depends on where the analog waveform is declared to cross. Different thresholds shift events relative to the reference."},
}
