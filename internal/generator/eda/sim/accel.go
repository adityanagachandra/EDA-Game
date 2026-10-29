package sim

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func throughput(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	fs := float64(edakit.Pick([]int{50, 100, 200, 500, 1000, 2000})) // cycles per second
	fe := float64(edakit.Pick([]int{500, 1000, 2000, 4000})) * 1000  // cycles per second
	C := float64(edakit.Int(2, 12))                                  // hours
	switch kp {
	case 1:
		N := float64(edakit.Int(1, 50)) * 1e6
		h := N / fs / 3600
		ans := edakit.Num(h, 2)
		q := edakit.Sprintf("An SoC simulates at %s cycles per second. How many hours does a %s-million-cycle boot test take? Round to 2 decimal places.", edakit.Num(fs, 0), edakit.Num(N/1e6, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("%s cycles / %s cycles per second = %s s, which is %s hours.", edakit.Num(N, 0), edakit.Num(fs, 0), edakit.Num(N/fs, 0), ans),
			map[string]string{"seconds": edakit.Num(N/fs, 0)})
	case 2:
		N := float64(edakit.Int(1, 100)) * 1e9
		h := C + N/fe/3600
		ans := edakit.Num(h, 2)
		q := edakit.Sprintf("An emulator runs a design at %s MHz after a compile that takes %s hours. What is the total time, in hours, to compile and run %s billion cycles? Round to 2 decimal places.", edakit.Num(fe/1e6, 1), edakit.Num(C, 0), edakit.Num(N/1e9, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Running takes %s cycles / %s cycles per second = %s hours; adding the %s-hour compile gives %s hours.", edakit.Num(N, 0), edakit.Num(fe, 0), edakit.Num(N/fe/3600, 2), edakit.Num(C, 0), ans),
			map[string]string{"run_hours": edakit.Num(N/fe/3600, 2)})
	default:
		nStar := C * 3600 / (1/fs - 1/fe)
		ans := edakit.Num(nStar/1e6, 2)
		q := edakit.Sprintf("Simulation runs at %s cycles per second with no compile cost to speak of. An emulator runs at %s MHz but needs a %s-hour compile. Above how many million cycles is the emulator faster overall? Round to 2 decimal places.", edakit.Num(fs, 0), edakit.Num(fe/1e6, 1), edakit.Num(C, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Set N / f_s = C + N / f_e, so N* = C / (1/f_s - 1/f_e) = %s s / (%s - %s) = %s million cycles.", edakit.Num(C*3600, 0), edakit.Num(1/fs, 6), edakit.Num(1/fe, 9), ans),
			map[string]string{"breakeven_cycles": edakit.Num(nStar, 0)})
	}
}

func partitioning(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		nb := edakit.Scaled(ctx, 5, 8)
		side := make([]int, nb)
		for i := range side {
			side[i] = edakit.Int(0, 1)
		}
		side[0], side[1] = 0, 1
		nn := edakit.Scaled(ctx, 5, 9)
		var desc []string
		cut := 0
		for k := 0; k < nn; k++ {
			deg := edakit.Int(2, 3)
			seen := map[int]bool{}
			var pins []int
			for len(pins) < deg {
				p := edakit.Int(0, nb-1)
				// Nets are mostly local to one FPGA, as a partitioner would leave them.
				if len(pins) > 0 && side[p] != side[pins[0]] && edakit.Int(0, 99) < 60 {
					continue
				}
				if !seen[p] {
					seen[p] = true
					pins = append(pins, p)
				}
			}
			var names []string
			s0, s1 := false, false
			for _, p := range pins {
				names = append(names, edakit.Sprintf("B%d", p+1))
				if side[p] == 0 {
					s0 = true
				} else {
					s1 = true
				}
			}
			if s0 && s1 {
				cut++
			}
			desc = append(desc, edakit.Sprintf("n%d {%s}", k+1, strings.Join(names, ", ")))
		}
		var fa, fb []string
		for i, s := range side {
			if s == 0 {
				fa = append(fa, edakit.Sprintf("B%d", i+1))
			} else {
				fb = append(fb, edakit.Sprintf("B%d", i+1))
			}
		}
		q := edakit.Sprintf("A design is split across two FPGAs: FPGA A holds %s and FPGA B holds %s. Its nets, each listing the blocks it connects, are %s. How many nets must cross between the FPGAs?",
			strings.Join(fa, ", "), strings.Join(fb, ", "), strings.Join(desc, "; "))
		return edakit.Numeric(q, edakit.IntStr(cut),
			edakit.Sprintf("A net is cut when its blocks sit on both FPGAs. %d nets are cut, and each needs an inter-FPGA signal.", cut),
			map[string]string{"cut": edakit.IntStr(cut)})
	case 2:
		pins := edakit.Pick([]int{200, 400, 480, 600, 800})
		sig := pins * edakit.Int(2, 12)
		sig += edakit.Int(-pins/2, pins/2)
		r := int(math.Ceil(float64(sig) / float64(pins)))
		q := edakit.Sprintf("A partition cuts %d signals between two FPGAs, which have %d physical I/O pins available between them. With time-division multiplexing, what is the smallest whole multiplexing ratio (signals per pin) that carries every signal?", sig, pins)
		return edakit.Numeric(q, edakit.IntStr(r),
			edakit.Sprintf("Each pin carries ceil(%d / %d) = %d signals in turn.", sig, pins, r),
			map[string]string{"ratio": edakit.IntStr(r)})
	default:
		link := float64(edakit.Pick([]int{200, 400, 500, 800, 1000}))
		r := edakit.Int(4, 32)
		o := edakit.Int(2, 6)
		f := link / float64(r+o)
		ans := edakit.Num(f, 2)
		q := edakit.Sprintf("Inter-FPGA links toggle at %s MHz. Each system clock cycle must move every multiplexed signal across a pin with a ratio of %d, plus %d link cycles of overhead for alignment. What is the highest system clock in MHz? Round to 2 decimal places.", edakit.Num(link, 0), r, o)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("One system cycle needs %d + %d = %d link cycles, so the system clock is %s / %d = %s MHz. Halving the cut halves the ratio and nearly doubles the speed.", r, o, r+o, edakit.Num(link, 0), r+o, ans),
			map[string]string{"link_cycles": edakit.IntStr(r + o)})
	}
}

var emulationItems = []edakit.Item{
	{KP: 1, Q: "What does a hardware emulator do with the design under test?",
		A: "Compiles it onto specialized hardware (custom processors or FPGAs) that executes it far faster than a software simulator", D: []string{"Simulates it on a larger server farm in software", "Converts it to a transaction-level model", "Runs a static timing analysis"},
		Why: "Emulation maps the RTL onto hardware that evaluates the logic in parallel."},
	{KP: 1, Q: "Compared with an FPGA prototype, what is an emulator's main advantage?",
		A: "Much deeper debug visibility and faster turnaround for design changes", D: []string{"Higher clock speed", "Lower cost per unit", "Running at the real product frequency"},
		Why: "Emulators are built for debug: trace capture and quick recompiles. Prototypes trade visibility for speed."},
	{KP: 1, Q: "Compared with an emulator, what is an FPGA prototype's main advantage?",
		A: "Higher run speed, which suits long software workloads", D: []string{"Full visibility of every signal at every cycle", "No need to synthesize the design", "Faster compile for each change"},
		Why: "Prototypes run faster but take longer to bring up and expose fewer signals."},
	{KP: 1, Q: "Why must the design be synthesizable for emulation?",
		A: "It is mapped onto real hardware primitives, so constructs with no hardware meaning cannot be compiled", D: []string{"Emulators only accept gate-level SDF", "Synthesizable code simulates faster in software", "Emulators need analog models"},
		Why: "Emulation compiles RTL to hardware; behavioral testbench code needs another home."},
	{KP: 2, Q: "A UVM testbench runs on a simulator host while the DUT runs on an emulator. What keeps the emulator from waiting on the host every cycle?",
		A: "Transaction-level communication through transactors, so the host exchanges messages rather than signals", D: []string{"Running the testbench at a higher clock", "Removing the testbench", "Converting the DUT to four-state"},
		Why: "Co-modeling moves pin-level protocol work into synthesizable transactors and sends transactions across the link."},
	{KP: 2, Q: "What is a speed bridge in an emulation or prototyping setup?",
		A: "Hardware that adapts a real-speed external interface, such as a network link, to the slower emulated design", D: []string{"A faster host computer", "A cache for waveforms", "A timing constraint file"},
		Why: "Real devices run at full speed; the bridge buffers and rate-matches traffic."},
	{KP: 2, Q: "A test runs one billion cycles at full visibility on an emulator. Why does capturing every signal at every cycle hurt?",
		A: "Trace storage and upload bandwidth limit how much can be captured, so capture windows are usually triggered", D: []string{"Visibility changes the design's function", "Emulators cannot store signals", "It forces a recompile"},
		Why: "Debug data volume is the practical limit; triggers and windows keep it manageable."},
	{KP: 2, Q: "Why are behavioral delays like #5 removed or replaced before emulation?",
		A: "The hardware runs on clocks and has no way to wait an arbitrary simulated delay", D: []string{"Delays are always ignored by simulators too", "Delays create X values", "Delays are not legal in SystemVerilog"},
		Why: "Emulation evaluates cycle by cycle; timing has to be expressed with clocks."},
	{KP: 3, Q: "A block is early in development and changes daily. Which platform suits its verification best?",
		A: "Software simulation, for fast turnaround and full debug", D: []string{"An FPGA prototype", "A multi-FPGA emulator", "A silicon test chip"},
		Why: "Early RTL churn favors the platform with no compile wait and complete visibility."},
	{KP: 3, Q: "A team needs to boot an operating system on a full SoC before tape-out and debug driver issues. Which platform is the usual choice?",
		A: "Emulation, which runs fast enough to boot and still offers hardware debug visibility", D: []string{"Event-driven gate-level simulation", "Formal property checking", "Static timing analysis"},
		Why: "Billions of cycles are too slow for simulation; emulation balances speed and visibility."},
	{KP: 3, Q: "Firmware developers need many copies of the design running near real speed for weeks. What fits?",
		A: "FPGA prototypes, replicated for the software team", D: []string{"One emulator shared by everyone", "Simulation farms", "Transaction-level models only"},
		Why: "Prototypes are cheaper per seat and faster, which suits long software development."},
	{KP: 3, Q: "Why does emulation still need the verification environment from simulation, such as scoreboards and coverage?",
		A: "Running faster finds nothing unless the results are checked and the reached scenarios are measured", D: []string{"Emulators check themselves", "Coverage is meaningless on emulators", "The DUT changes on the emulator"},
		Why: "Speed only multiplies the value of good checking and coverage."},
}
