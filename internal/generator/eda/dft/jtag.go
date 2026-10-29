package dft

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// The sixteen-state IEEE 1149.1 TAP controller: next[state] = {on TMS 0, on TMS 1}.
var tapNext = map[string][2]string{
	"Test-Logic-Reset": {"Run-Test/Idle", "Test-Logic-Reset"},
	"Run-Test/Idle":    {"Run-Test/Idle", "Select-DR-Scan"},
	"Select-DR-Scan":   {"Capture-DR", "Select-IR-Scan"},
	"Capture-DR":       {"Shift-DR", "Exit1-DR"},
	"Shift-DR":         {"Shift-DR", "Exit1-DR"},
	"Exit1-DR":         {"Pause-DR", "Update-DR"},
	"Pause-DR":         {"Pause-DR", "Exit2-DR"},
	"Exit2-DR":         {"Shift-DR", "Update-DR"},
	"Update-DR":        {"Run-Test/Idle", "Select-DR-Scan"},
	"Select-IR-Scan":   {"Capture-IR", "Test-Logic-Reset"},
	"Capture-IR":       {"Shift-IR", "Exit1-IR"},
	"Shift-IR":         {"Shift-IR", "Exit1-IR"},
	"Exit1-IR":         {"Pause-IR", "Update-IR"},
	"Pause-IR":         {"Pause-IR", "Exit2-IR"},
	"Exit2-IR":         {"Shift-IR", "Update-IR"},
	"Update-IR":        {"Run-Test/Idle", "Select-DR-Scan"},
}

var tapStates = []string{"Test-Logic-Reset", "Run-Test/Idle", "Select-DR-Scan", "Capture-DR", "Shift-DR", "Exit1-DR", "Pause-DR", "Exit2-DR", "Update-DR",
	"Select-IR-Scan", "Capture-IR", "Shift-IR", "Exit1-IR", "Pause-IR", "Exit2-IR", "Update-IR"}

func tapRun(s string, tms []int) string {
	for _, t := range tms {
		s = tapNext[s][t]
	}
	return s
}

func tapDist(a, b string) int {
	dist := map[string]int{a: 0}
	q := []string{a}
	for len(q) > 0 {
		x := q[0]
		q = q[1:]
		if x == b {
			return dist[x]
		}
		for _, y := range tapNext[x] {
			if _, ok := dist[y]; !ok {
				dist[y] = dist[x] + 1
				q = append(q, y)
			}
		}
	}
	return -1
}

func tmsText(t []int) string {
	var p []string
	for _, x := range t {
		p = append(p, edakit.IntStr(x))
	}
	return strings.Join(p, ", ")
}

// tapFSM: state after a TMS sequence (KP 1), shortest TMS sequence (KP 2), TCK cycles for a full
// scan (KP 3).
func tapFSM(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		s := edakit.Pick(tapStates)
		n := edakit.Scaled(ctx, 3, 8)
		tms := make([]int, n)
		for i := range tms {
			tms[i] = edakit.Int(0, 1)
		}
		ans := tapRun(s, tms)
		short := tapRun(s, tms[:n-1])
		flip := append([]int(nil), tms...)
		flip[n-1] ^= 1
		ds := uniqStates(ans, tapRun(s, flip), short, tapRun(s, append(tms, 0)), edakit.Pick(tapStates), edakit.Pick(tapStates), edakit.Pick(tapStates))
		return edakit.Choice(edakit.Sprintf("A TAP controller is in %s. TMS takes the values %s on successive rising edges of TCK. Which state is it in afterwards?", s, tmsText(tms)),
			ans, edakit.Sprintf("Following the 1149.1 state diagram edge by edge from %s ends in %s.", s, ans), ds...)
	case 2:
		ia := edakit.Int(0, len(tapStates)-1)
		a := tapStates[ia]
		b := tapStates[(ia+edakit.Int(1, len(tapStates)-1))%len(tapStates)]
		d := tapDist(a, b)
		return edakit.Choice(edakit.Sprintf("What is the fewest TCK cycles that move a TAP controller from %s to %s?", a, b),
			edakit.IntStr(d), edakit.Sprintf("A breadth-first search over the 16-state diagram gives %d.", d), edakit.Distinct(d, 3, d+1, d-1, d+2)...)
	default:
		k := edakit.Scaled(ctx, 4, 300)
		ir := edakit.Int(0, 1) == 1
		ans := k + 5
		reg := "data register"
		if ir {
			ans = k + 6
			reg = "instruction register"
		}
		return edakit.Choice(edakit.Sprintf("Starting in Run-Test/Idle, a host scans %d bits through the %s and returns to Run-Test/Idle by the shortest TMS sequence, shifting the last bit on the edge that leaves the Shift state. How many TCK cycles does that take?", k, reg),
			edakit.IntStr(ans), edakit.Sprintf("Reaching Shift takes %d edges, shifting takes %d (the last one exits to Exit1), and Exit1 to Update to Run-Test/Idle takes 2: %d.", ans-k-2, k, ans),
			edakit.Distinct(ans, 3, k+4, k+5, k+6, k+7, 2*k)...)
	}
}

func uniqStates(ans string, cands ...string) []string {
	seen := map[string]bool{ans: true}
	var out []string
	for _, c := range cands {
		if !seen[c] && len(out) < 3 {
			seen[c] = true
			out = append(out, c)
		}
	}
	for _, c := range edakit.Shuffle(tapStates) {
		if !seen[c] && len(out) < 3 {
			seen[c] = true
			out = append(out, c)
		}
	}
	return out
}

// boundaryScan: DR length of a board chain (KP 1), interconnect test vectors (KP 2), EXTEST time
// (KP 3).
func boundaryScan(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 3, 6)
		total := 0
		var parts []string
		for i := 0; i < n; i++ {
			bsr := edakit.Int(20, 600)
			if edakit.Int(0, 2) == 0 {
				parts = append(parts, edakit.Sprintf("U%d (boundary register %d cells) in BYPASS", i+1, bsr))
				total++
			} else {
				parts = append(parts, edakit.Sprintf("U%d (boundary register %d cells) in EXTEST", i+1, bsr))
				total += bsr
			}
		}
		return edakit.Numeric("A board's JTAG chain links "+strings.Join(parts, ", ")+". How many bits long is the selected data-register path from TDI to TDO?",
			edakit.IntStr(total), edakit.Sprintf("EXTEST selects a device's boundary register; BYPASS selects its 1-bit bypass register. The sum is %d.", total), nil)
	case 2:
		nets := edakit.Scaled(ctx, 10, 2000)
		v := int(math.Ceil(math.Log2(float64(nets + 2))))
		return edakit.Numeric(edakit.Sprintf("A boundary-scan interconnect test must give every one of %d board nets a unique serial code, with no net seeing all 0s or all 1s (so that stuck nets and wired shorts are distinguishable). Each parallel test vector supplies one bit of every code. What is the fewest vectors?", nets),
			edakit.IntStr(v), edakit.Sprintf("Excluding the all-0 and all-1 codes leaves $2^v - 2 \\ge %d$, so $v = \\lceil \\log_2(%d + 2) \\rceil = %d$. This is the modified counting sequence.", nets, nets, v), nil)
	default:
		l := edakit.Scaled(ctx, 100, 3000)
		irl := edakit.Int(4, 32)
		vecs := edakit.Scaled(ctx, 8, 200)
		f := edakit.Pick([]int{5, 10, 20, 25, 50})
		cyc := irl + 6 + vecs*(l+5)
		us := float64(cyc) / float64(f)
		ans := edakit.Num(us, 2)
		return edakit.Numeric(edakit.Sprintf("An EXTEST session loads the instruction once (an IR scan of %d bits costs %d + 6 TCK cycles from Run-Test/Idle back to it) and then applies %d vectors, each a DR scan of the %d-bit boundary chain costing that length + 5 cycles. TCK runs at %d MHz. How long does the session take, in microseconds? Round to 2 decimal places.", irl, irl, vecs, l, f),
			ans, edakit.Sprintf("$%d + 6 + %d(%d + 5) = %d$ cycles at %d MHz is %s microseconds.", irl, vecs, l, cyc, f, ans), nil)
	}
}

// instructionsCalc (KP 2): the length of the selected path for an instruction mix.
func instructionsCalc(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 3, 6)
	j := edakit.Int(1, n)
	var irs []string
	irTot := 0
	for i := 0; i < n; i++ {
		x := edakit.Int(3, 12)
		irTot += x
		irs = append(irs, edakit.IntStr(x))
	}
	switch edakit.Int(0, 2) {
	case 0:
		ans := 32 + n - 1
		return edakit.Choice(edakit.Sprintf("A chain of %d devices is set up so that U%d holds IDCODE and every other device holds BYPASS. How many bits long is the data-register path from TDI to TDO?", n, j),
			edakit.IntStr(ans), edakit.Sprintf("IDCODE selects the 32-bit identification register; each of the other %d devices contributes its 1-bit bypass register: %d.", n-1, ans),
			edakit.Distinct(ans, 3, 32, 32*n, 32+n, n)...)
	case 1:
		return edakit.Choice(edakit.Sprintf("A chain of %d devices has instruction registers of %s bits. How many bits must an IR scan shift to load one instruction into every device?", n, strings.Join(irs, ", ")),
			edakit.IntStr(irTot), edakit.Sprintf("Instruction registers are always all in the path during an IR scan, so their lengths add: %d. There is no bypass for the IR.", irTot),
			edakit.Distinct(irTot, 3, irTot-n+1, n, irTot+n)...)
	default:
		bsr := edakit.Int(40, 500)
		ans := bsr + n - 1
		return edakit.Choice(edakit.Sprintf("A chain of %d devices is set up so that U%d (a %d-cell boundary register) holds SAMPLE/PRELOAD and every other device holds BYPASS. How many bits long is the data-register path?", n, j, bsr),
			edakit.IntStr(ans), edakit.Sprintf("SAMPLE/PRELOAD selects the boundary register of U%d; the others add 1 bit each: %d + %d = %d.", j, bsr, n-1, ans),
			edakit.Distinct(ans, 3, bsr, bsr+n, bsr*n)...)
	}
}

// ijtagCalc (KP 2): the active scan path length of a SIB network.
func ijtagCalc(ctx generator.GeneratorContext) generator.Problem {
	m := edakit.Scaled(ctx, 2, 4)
	total := 0
	var desc []string
	var why []string
	for i := 0; i < m; i++ {
		open := edakit.Int(0, 1) == 1
		st := "closed"
		if open {
			st = "open"
		}
		if edakit.Int(0, 2) == 0 {
			// nested SIBs
			k := edakit.Int(2, 3)
			var sub []string
			subLen := 0
			for j := 0; j < k; j++ {
				o2 := edakit.Int(0, 1) == 1
				tdr := edakit.Int(4, 64)
				s2 := "closed"
				if o2 {
					s2 = "open"
					subLen += 1 + tdr
				} else {
					subLen++
				}
				sub = append(sub, edakit.Sprintf("SIB%d.%d (%s) over a %d-bit instrument register", i+1, j+1, s2, tdr))
			}
			desc = append(desc, edakit.Sprintf("SIB%d (%s) whose segment holds %s", i+1, st, strings.Join(sub, " and ")))
			if open {
				total += 1 + subLen
				why = append(why, edakit.Sprintf("SIB%d open: 1 + %d", i+1, subLen))
			} else {
				total++
				why = append(why, edakit.Sprintf("SIB%d closed: 1", i+1))
			}
			continue
		}
		tdr := edakit.Int(4, 128)
		desc = append(desc, edakit.Sprintf("SIB%d (%s) over a %d-bit instrument register", i+1, st, tdr))
		if open {
			total += 1 + tdr
			why = append(why, edakit.Sprintf("SIB%d open: 1 + %d", i+1, tdr))
		} else {
			total++
			why = append(why, edakit.Sprintf("SIB%d closed: 1", i+1))
		}
	}
	return edakit.Choice("An IEEE 1687 network places, in series on one TDR path: "+strings.Join(desc, "; ")+". Each SIB is a 1-bit cell; an open SIB inserts its segment into the path, a closed one bypasses it. How many bits long is the active scan path?",
		edakit.IntStr(total), strings.Join(why, "; ")+edakit.Sprintf("; total %d. Nested segments count only when every SIB above them is open.", total),
		edakit.Distinct(total, 3, total-m, total+m, total+1)...)
}

var instructionItems = []edakit.Item{
	{KP: 1, Q: "Which IEEE 1149.1 instruction connects a device's boundary register between TDI and TDO and drives the captured values onto its output pins?",
		A: "EXTEST", D: []string{"BYPASS", "SAMPLE/PRELOAD", "IDCODE"},
		Why: "EXTEST takes the pins away from the core and lets the boundary cells drive and capture board interconnect."},
	{KP: 1, Q: "A board is running in-system and the engineer wants to see the values on a device's pins without disturbing operation. Which instruction?",
		A: "SAMPLE/PRELOAD", D: []string{"EXTEST", "INTEST", "CLAMP"},
		Why: "SAMPLE captures the pin values into the boundary register while the core keeps control of the pins."},
	{KP: 1, Q: "What does the BYPASS instruction select, and why does it exist?",
		A: "A 1-bit register, so devices not under test add only one bit to the chain", D: []string{"The full boundary register, for safety", "The instruction register, for reloading", "No register: TDI connects to TDO by a wire"},
		Why: "Bypass shortens long board chains; the single flop also keeps TDO registered for timing."},
	{KP: 1, Q: "Which register does a compliant device select in Test-Logic-Reset if it implements IDCODE?",
		A: "The 32-bit device identification register", D: []string{"The bypass register", "The boundary register", "The instruction register"},
		Why: "After reset the IDCODE instruction (or BYPASS when there is no IDCODE) is loaded, so a blind DR scan can read the chain's device IDs."},
	{KP: 1, Q: "What does the CLAMP instruction do?",
		A: "Drives pins from the boundary register while selecting the bypass register as the data path", D: []string{"Puts every output pin in high impedance", "Resets the TAP controller", "Captures the core's internal state"},
		Why: "CLAMP holds outputs at preloaded values and keeps the scan path short; HIGHZ instead floats the outputs."},
	{KP: 1, Q: "Which instruction isolates a device by placing all its outputs in high impedance?",
		A: "HIGHZ", D: []string{"CLAMP", "BYPASS", "SAMPLE/PRELOAD"},
		Why: "HIGHZ disables output drivers so the device cannot fight a neighbour while other parts are tested."},
	{KP: 3, Q: "Why must an engineer PRELOAD the boundary register before switching to EXTEST?",
		A: "EXTEST drives the boundary cells onto the pins immediately on Update-IR, so they must already hold safe values", D: []string{"EXTEST cannot shift data", "PRELOAD resets the TAP", "Because IDCODE overwrites the boundary register"},
		Why: "Without a preload, the pins would drive whatever random values the cells held, which can cause bus contention on the board."},
	{KP: 3, Q: "During Capture-IR, what does every 1149.1 instruction register load into its two least significant bits?",
		A: "The fixed pattern 01, which lets a host check chain integrity and count devices", D: []string{"The current instruction", "All zeros", "The device's IDCODE"},
		Why: "The mandatory 01 capture value is a known signature for each IR, so a shifted-out IR stream reveals breaks and lengths."},
	{KP: 3, Q: "Why does the IR scan always include every device's instruction register while a DR scan can use bypass?",
		A: "The instruction registers are what select each device's data register, so they cannot themselves be bypassed", D: []string{"IR scans are shorter by design", "Bypass applies only to IR scans", "The standard forbids DR scans through more than one device"},
		Why: "BYPASS is an instruction; it can only take effect after it is loaded through the full IR chain."},
	{KP: 3, Q: "A board chain reads back 0xFFFFFFFF for every IDCODE. What is the most likely cause?",
		A: "TDO is stuck high or the chain is broken, so only the pull-up is seen", D: []string{"Every device uses the same manufacturer code", "The devices are in BYPASS correctly", "IDCODE is 64 bits long"},
		Why: "IDCODE bit 0 is always 1 and a valid code is never all ones; reading all ones points to a broken TDI-TDO path."},
	{KP: 3, Q: "INTEST is optional in 1149.1. What does it do?",
		A: "Applies test values to the core's inputs through the boundary cells and captures its outputs, isolating it from the board", D: []string{"Tests board interconnect between devices", "Reads the device ID", "Resets the core"},
		Why: "INTEST turns the boundary register into a wrapper around the core, the idea IEEE 1500 later generalizes to embedded cores."},
	{KP: 3, Q: "Why are user-defined private instructions common in real 1149.1 devices?",
		A: "They give access to internal scan chains, BIST controllers and debug logic through the same TAP", D: []string{"The standard forbids them", "They replace BYPASS", "They are required for EXTEST"},
		Why: "The TAP is the chip's universal test port; vendors add private instructions to reach internal instruments."},
}

var ijtagItems = []edakit.Item{
	{KP: 1, Q: "What does an IEEE 1500 wrapper add around an embedded core?",
		A: "A wrapper boundary register on the core's terminals, a wrapper instruction register and a bypass, controlled through a standard serial interface", D: []string{"A TAP controller with TMS and TCK pins on the package", "An on-chip tester with pattern memory", "A second copy of the core"},
		Why: "IEEE 1500 brings boundary-scan ideas inside the chip so each core can be tested in isolation and its interconnect tested separately."},
	{KP: 1, Q: "Why does IEEE 1500 not mandate a TAP state machine?",
		A: "Cores sit inside a chip, so the wrapper is driven by on-chip control signals rather than package pins", D: []string{"It uses the 1149.1 TAP of every core", "Cores have no test logic", "It only supports parallel access"},
		Why: "The wrapper serial port is controlled directly by chip-level test logic, often itself driven from the chip's 1149.1 TAP."},
	{KP: 1, Q: "What is the purpose of the wrapper parallel port in IEEE 1500?",
		A: "To feed many internal scan chains at once for fast core test, beyond the one-bit serial port", D: []string{"To replace the wrapper instruction register", "To connect the core to the package pins directly", "To implement the TAP state machine"},
		Why: "Serial access is enough for configuration; production scan test of a core needs width, which the optional parallel port provides."},
	{KP: 1, Q: "A core vendor ships a core with an IEEE 1500 wrapper and its test patterns. What does the integrator gain?",
		A: "The patterns can be reused at chip level by mapping them through the wrapper, without regenerating them", D: []string{"The core no longer needs testing", "The chip needs no scan chains", "ATPG must be run on the whole chip anyway"},
		Why: "Test reuse is the point: the wrapper defines a known boundary at which core-level patterns are valid."},
	{KP: 1, Q: "Which wrapper mode tests the wires between two wrapped cores?",
		A: "External test mode, with the wrapper cells driving and capturing core terminals", D: []string{"Internal test mode", "Bypass mode", "Functional mode"},
		Why: "Like EXTEST in 1149.1, external mode isolates the core and exercises the interconnect around it."},
	{KP: 3, Q: "What does a SIB (segment insertion bit) in an IEEE 1687 network do?",
		A: "It is a 1-bit scan cell that, when set, splices a sub-network into the active scan path", D: []string{"It stores the instrument's data", "It resets every instrument", "It replaces the TAP controller"},
		Why: "SIBs make the scan path reconfigurable so only the instruments of interest are in the path, keeping access short."},
	{KP: 3, Q: "What do ICL and PDL describe in IEEE 1687?",
		A: "ICL describes the network and instrument interfaces; PDL describes operations on instruments at their own terminals", D: []string{"ICL holds test patterns; PDL holds layout", "Both are netlist formats for synthesis", "ICL is a tester language and PDL a timing format"},
		Why: "Separating structure (ICL) from procedure (PDL) lets a tool retarget an instrument-level procedure to chip pins through any network."},
	{KP: 3, Q: "What does PDL retargeting compute?",
		A: "The sequence of scan operations at the chip's TAP that reaches an instrument through the current SIB configuration", D: []string{"A new placement of the instruments", "The timing of the instrument's internal paths", "A compressed version of the ATPG patterns"},
		Why: "The tool knows the network from ICL and translates an iWrite or iRead at the instrument into the path setup and shifts needed at the top."},
	{KP: 3, Q: "Why is IEEE 1687 well suited to hundreds of on-chip instruments such as sensors and BIST controllers?",
		A: "The reconfigurable network keeps the scan path short by opening only the needed segments", D: []string{"It gives each instrument its own package pin", "It requires no scan path", "It loads every instrument on every access"},
		Why: "A flat chain of all instrument registers would be very long; SIB-gated hierarchy makes access time scale with what is used."},
	{KP: 3, Q: "An instrument sits behind two nested SIBs, both closed. How many scan operations are needed to read it?",
		A: "At least three: open the outer SIB, then the inner SIB, then shift the instrument's register", D: []string{"One: the instrument is always in the path", "Two: one per SIB at once", "None: SIBs open automatically"},
		Why: "Each closed SIB hides what is under it until a scan sets it, so nested access costs one scan per level plus the data access."},
	{KP: 3, Q: "How do IEEE 1149.1, 1500 and 1687 typically fit together on one chip?",
		A: "The 1149.1 TAP is the chip-level port; 1500 wrappers isolate cores; a 1687 network reaches embedded instruments through the TAP", D: []string{"They are alternatives and never coexist", "1687 replaces the package TAP", "1500 is only for board test"},
		Why: "They address different scopes, board, core and instrument, and are designed to compose."},
	{KP: 3, Q: "What does IEEE 1500's core test language (CTL) carry that plain patterns do not?",
		A: "How the core's test modes, wrapper and patterns are configured, so the integrator can schedule and reuse them", D: []string{"The core's RTL", "The layout of the wrapper", "The power grid of the core"},
		Why: "CTL, built on STIL, describes test information about the core so chip-level tools can integrate it."},
}
