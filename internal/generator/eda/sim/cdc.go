package sim

import (
	"math"
	"math/bits"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func mtbf(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	for {
		tau := float64(edakit.Int(10, 30))    // ps
		t0 := float64(edakit.Int(2, 10) * 10) // ps
		fc := float64(edakit.Pick([]int{200, 250, 400, 500, 800, 1000}))
		fd := float64(edakit.Pick([]int{10, 20, 25, 50, 100}))
		den := t0 * 1e-12 * fc * 1e6 * fd * 1e6
		switch kp {
		case 1:
			tr := float64(edakit.Int(2, 8) * 50) // ps
			m := math.Exp(tr/tau) / den
			if m < 1 || m > 1e6 {
				continue
			}
			ans := edakit.Num(m, 1)
			q := edakit.Sprintf("A synchronizer flop has resolution time constant tau = %s ps and metastability window T0 = %s ps. It is clocked at %s MHz, the asynchronous data toggles at %s MHz, and %s ps are available for resolution. Using MTBF = e^(t_r/tau) / (T0 f_clk f_data), what is the MTBF in seconds? Round to 1 decimal place.",
				edakit.Num(tau, 0), edakit.Num(t0, 0), edakit.Num(fc, 0), edakit.Num(fd, 0), edakit.Num(tr, 0))
			return edakit.Numeric(q, ans,
				edakit.Sprintf("t_r/tau = %s, so e^(t_r/tau) = %s. The denominator T0 f_clk f_data = %s per second. MTBF = %s s, far too short for a product.", edakit.Num(tr/tau, 3), edakit.Num(math.Exp(tr/tau), 0), edakit.Num(den, 0), ans),
				map[string]string{"exponent": edakit.Num(tr/tau, 3)})
		case 2:
			period := 1e6 / fc // ps
			f := period / tau / math.Ln10
			ans := edakit.Num(f, 1)
			q := edakit.Sprintf("A two-flop synchronizer is clocked at %s MHz and its flops have tau = %s ps. Adding a third flop gives roughly one more clock period of resolution time. By how many orders of magnitude (log10 of the improvement factor) does the MTBF grow? Round to 1 decimal place.", edakit.Num(fc, 0), edakit.Num(tau, 0))
			return edakit.Numeric(q, ans,
				edakit.Sprintf("The extra stage multiplies MTBF by e^(T/tau) with T = %s ps, so log10 of the factor is T/(tau ln 10) = %s.", edakit.Num(period, 0), ans),
				map[string]string{"period_ps": edakit.Num(period, 0)})
		default:
			years := float64(edakit.Pick([]int{10, 100, 1000, 10000}))
			target := years * 365 * 24 * 3600
			tr := tau * math.Log(target*den)
			ans := edakit.Num(tr, 0)
			q := edakit.Sprintf("A synchronizer has tau = %s ps and T0 = %s ps, with f_clk = %s MHz and f_data = %s MHz. Using MTBF = e^(t_r/tau) / (T0 f_clk f_data) and a year of 365 days, what resolution time t_r in ps gives an MTBF of %s years? Round to the nearest ps.",
				edakit.Num(tau, 0), edakit.Num(t0, 0), edakit.Num(fc, 0), edakit.Num(fd, 0), edakit.Num(years, 0))
			return edakit.Numeric(q, ans,
				edakit.Sprintf("Solve t_r = tau ln(MTBF x T0 f_clk f_data) = %s x ln(%s x %s) = %s ps. Each factor of 10 in MTBF costs only tau ln 10 = %s ps more.", edakit.Num(tau, 0), edakit.Num(target, 0), edakit.Num(den, 0), ans, edakit.Num(tau*math.Ln10, 1)),
				map[string]string{"target_s": edakit.Num(target, 0)})
		}
	}
}

func multibit(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		b := edakit.Int(2, 1<<edakit.Scaled(ctx, 4, 8)-1)
		g := b ^ (b >> 1)
		q := edakit.Sprintf("A pointer crossing clock domains is converted to Gray code with g = b XOR (b >> 1). What is the Gray code of b = %d, written as a decimal number?", b)
		return edakit.Numeric(q, edakit.IntStr(g),
			edakit.Sprintf("b = %b in binary and b >> 1 = %b; their XOR is %b, which is %d.", b, b>>1, g, g),
			map[string]string{"gray_binary": edakit.Sprintf("%b", g)})
	case 2:
		w := edakit.Int(4, 8)
		k := edakit.Int(1, 1<<w-2)
		t := bits.OnesCount(uint(k ^ (k + 1)))
		q := edakit.Sprintf("A %d-bit binary counter in one clock domain is sampled by another clock domain through per-bit synchronizers. When it increments from %d to %d, how many bits change at once, each of which could be captured on either side of the change?", w, k, k+1)
		return edakit.Numeric(q, edakit.IntStr(t),
			edakit.Sprintf("%d XOR %d = %b, which has %d set bits. A Gray-coded counter would change exactly one bit, so any sample is either the old or the new value.", k, k+1, k^(k+1), t),
			map[string]string{"toggles": edakit.IntStr(t)})
	default:
		B := edakit.Int(4, 40) * 8
		fw := edakit.Pick([]int{100, 200, 250, 400, 500})
		fr := edakit.Pick([]int{50, 80, 100, 125, 150, 200})
		for fr >= fw {
			fr = edakit.Pick([]int{50, 80, 100, 125})
		}
		d := math.Ceil(float64(B)*(1-float64(fr)/float64(fw)) - 1e-9)
		ans := edakit.Num(d, 0)
		q := edakit.Sprintf("A burst of %d words is written into an asynchronous FIFO at one word per cycle at %d MHz, while the reader takes one word per cycle at %d MHz. Ignoring synchronizer latency, what is the minimum FIFO depth that avoids overflow?", B, fw, fr)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The burst lasts %d / %d MHz, during which the reader removes %d x %d/%d words. The backlog is B (1 - f_r/f_w) = %s, rounded up to %s.", B, fw, B, fr, fw, edakit.Num(float64(B)*(1-float64(fr)/float64(fw)), 2), ans),
			map[string]string{"depth": ans})
	}
}

var synchronizerItems = []edakit.Item{
	{KP: 1, Q: "Why is a single-bit control signal from another clock domain passed through two flip-flops before use?",
		A: "The first flop may go metastable; the second gives it a full cycle to resolve before the value is used", D: []string{"Two flops double the signal's speed", "The second flop removes glitches from combinational logic", "It converts the signal to Gray code"},
		Why: "The synchronizer does not prevent metastability; it makes the probability that it propagates very small."},
	{KP: 1, Q: "What does a two-flop synchronizer not guarantee?",
		A: "The cycle in which the destination sees the change: it may be one cycle later or earlier", D: []string{"That the output is a valid logic level most of the time", "That metastability is filtered with high probability", "That the destination flop is clocked by its own clock"},
		Why: "Latency through a synchronizer is uncertain by a cycle, which designs must tolerate."},
	{KP: 1, Q: "Why must the source of a synchronized signal come straight from a flip-flop?",
		A: "Combinational logic can glitch, and a glitch may be captured as a false pulse in the other domain", D: []string{"Flops are faster than gates", "Synchronizers only accept Gray code", "It saves area"},
		Why: "The receiving domain samples at arbitrary times, so any glitch can be caught."},
	{KP: 1, Q: "What does a structural CDC tool report on a design?",
		A: "Crossings between clock domains that lack an approved synchronization scheme", D: []string{"Setup violations within one domain", "Coverage holes", "Functional mismatches against a reference model"},
		Why: "Structural CDC analysis classifies every crossing and flags unsynchronized or badly synchronized paths."},
	{KP: 2, Q: "Two related control bits are synchronized by separate two-flop synchronizers and recombined in the destination. What is the risk?",
		A: "Reconvergence: the bits can arrive in different cycles, creating a combination the source never produced", D: []string{"No risk: each bit is safe", "The bits always arrive together", "It causes a setup violation in the source domain"},
		Why: "Independent synchronizers have independent latency, so the destination can see mixtures."},
	{KP: 2, Q: "A one-cycle pulse from a 500 MHz domain is synchronized into a 100 MHz domain. What can happen?",
		A: "The pulse can fall between destination clock edges and be missed entirely", D: []string{"It is always seen exactly once", "It is seen five times", "It becomes metastable forever"},
		Why: "Fast-to-slow pulses must be stretched, or a toggle or handshake used."},
	{KP: 2, Q: "Which scheme safely moves a multi-bit data word from one domain to another when it changes rarely?",
		A: "A handshake: hold the data stable and send a synchronized request, then wait for an acknowledgement", D: []string{"A two-flop synchronizer on each bit", "A single flop on each bit", "An inverter chain"},
		Why: "Data held stable while a single control bit is synchronized is sampled safely."},
	{KP: 2, Q: "A CDC tool flags a synchronizer where the synchronized signal fans out to two separate flops in the destination before the second stage. Why?",
		A: "The first-stage flop may be metastable, and two loads may resolve it differently", D: []string{"Fanout is never allowed in RTL", "It violates Gray coding", "It doubles the clock frequency"},
		Why: "Only the second stage should read the first stage; otherwise metastability can spread with divergent values."},
	{KP: 3, Q: "Why does ordinary zero-delay RTL simulation almost never reveal CDC bugs?",
		A: "It never models metastability or the extra cycle of uncertain latency at a synchronizer", D: []string{"It runs too slowly", "It ignores clocks", "CDC bugs only appear at gate level"},
		Why: "In RTL simulation a flop always captures the clean value, so latency-dependent bugs stay hidden."},
	{KP: 3, Q: "What does metastability injection in simulation do?",
		A: "It randomly adds or removes a cycle of delay at synchronizer outputs to emulate unpredictable resolution", D: []string{"It sets signals to X permanently", "It slows down the clock", "It removes the synchronizers"},
		Why: "Random latency exposes logic that depends on exact crossing timing, such as reconvergence bugs."},
	{KP: 3, Q: "A design passes CDC structural checks and all tests, then fails in silicon with a rare data corruption on a crossing. What is the likely gap?",
		A: "A functional protocol issue, such as data changing before the handshake completes, that structure alone cannot catch", D: []string{"A missing timing constraint inside one domain", "Too few random seeds in unit tests", "A simulator bug"},
		Why: "Structural CDC proves synchronizers exist; protocol correctness needs assertions or formal checks."},
	{KP: 3, Q: "Why are assertions placed on CDC handshake signals in simulation?",
		A: "They check protocol rules, for example that data is stable while the request is pending", D: []string{"They speed up the crossing", "They replace synchronizers", "They measure MTBF"},
		Why: "Assertions verify the functional contract that keeps a structurally correct crossing safe."},
}

var resetItems = []edakit.Item{
	{KP: 1, Q: "What does a reset synchronizer do?",
		A: "Asserts reset asynchronously and releases it synchronously to the local clock", D: []string{"Asserts and releases reset only on clock edges", "Delays assertion by two cycles", "Removes the need for a reset"},
		Why: "Immediate assertion works without a clock; synchronous release avoids recovery and removal violations."},
	{KP: 1, Q: "Why is asynchronous reset deassertion near a clock edge dangerous?",
		A: "Flops can violate recovery or removal time and come out of reset in different cycles or metastable", D: []string{"It causes a setup violation on the data path only", "It is harmless because reset is asynchronous", "It changes the clock frequency"},
		Why: "Recovery and removal are the reset-pin equivalents of setup and hold."},
	{KP: 1, Q: "A reset synchronizer is built from two flops with their asynchronous reset tied to the external reset. What is fed into the first flop's D input?",
		A: "A constant 1 (the deasserted level), so release ripples through the two flops", D: []string{"The external reset itself", "The clock", "The second flop's output"},
		Why: "With D tied to the inactive level, the release propagates on clock edges after reset goes away."},
	{KP: 1, Q: "Each clock domain has its own reset synchronizer. Why?",
		A: "Release must be synchronous to each domain's own clock", D: []string{"One synchronizer cannot drive many flops", "Resets must be Gray coded", "To save power"},
		Why: "Synchronous release is per clock; a single synchronizer serves only one domain."},
	{KP: 2, Q: "What is a reset-domain crossing?",
		A: "A path from a flop reset by one reset to a flop reset by a different reset, or not reset", D: []string{"A path between two clock domains", "A path from reset to the clock pin", "A reset applied to a memory"},
		Why: "RDC concerns asynchronous resets in different domains, even when the clock is the same."},
	{KP: 2, Q: "A source flop is asynchronously reset while its destination flop, on the same clock but not reset, keeps capturing. What can happen?",
		A: "The source output changes asynchronously and the destination can capture it metastable", D: []string{"Nothing, because the clock is shared", "The destination resets too", "A setup check catches it in STA"},
		Why: "An asynchronous reset changes the source at an arbitrary time relative to the destination clock."},
	{KP: 2, Q: "Which fix is standard for a reset-domain crossing where the source can be reset while the destination runs?",
		A: "Gate or isolate the destination's capture (for example hold its enable) while the source reset is asserted", D: []string{"Add a third synchronizer stage to the clock", "Remove reset from the source", "Increase the clock frequency"},
		Why: "Isolating the destination during the source reset prevents capturing the asynchronous change."},
	{KP: 2, Q: "Why do timing tools not catch reset-domain crossing problems by default?",
		A: "Asynchronous reset paths are usually not timed against the destination clock as data paths", D: []string{"Timing tools ignore resets completely", "RDC only matters at gate level", "STA cannot handle two flops"},
		Why: "RDC is a structural-and-functional check, like CDC, outside ordinary STA."},
	{KP: 3, Q: "A block's soft reset can fire while the rest of the chip runs. What must verification check?",
		A: "That every path from the reset block into running logic is isolated or tolerant during and after the reset", D: []string{"Only that the block's registers reach their reset values", "Only the reset pulse width", "That the reset is Gray coded"},
		Why: "Partial resets are the classic source of RDC bugs; isolation must be verified on every crossing."},
	{KP: 3, Q: "Why can a simulation show a clean reset while silicon fails to come out of reset reliably?",
		A: "Simulation does not model recovery and removal timing or metastability on release", D: []string{"Simulation always releases reset early", "Silicon has no reset", "The testbench applied reset too long"},
		Why: "Zero-delay simulation treats release as clean; release timing issues need synchronizers and checks."},
	{KP: 3, Q: "Resets in two domains must release in a specific order. How is that ensured?",
		A: "Sequence the releases deliberately, for example by gating one domain's reset release on the other's", D: []string{"Use one reset synchronizer for both domains", "Release them at the same instant", "Remove the reset from one domain"},
		Why: "Ordering dependencies between domains are designed in, then checked with assertions."},
	{KP: 3, Q: "Why is X-propagation in simulation useful when checking reset?",
		A: "Flops that are not reset stay X, revealing logic that depends on them before it is initialized", D: []string{"X speeds up the simulation", "X proves the reset is synchronous", "X replaces metastability modeling exactly"},
		Why: "Four-state simulation exposes uninitialized state that a two-state simulator would hide."},
}
