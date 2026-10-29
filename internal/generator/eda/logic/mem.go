package logic

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func registerMem(reg *generator.Registry) {
	reg.Register("logic.mem.sram", retryTies(sram))
	reg.Register("logic.mem.dram", retryTies(dram))
	reg.Register("logic.mem.register_file", retryTies(registerFile))
	reg.Register("logic.mem.cam", retryTies(cam))
	reg.Register("logic.mem.macros_nvm", &edakit.Bank{Items: macroItems})
}

func sram(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		W := 1 << uint(edakit.Scaled(ctx, 8, 14))
		B := edakit.Pick([]int{8, 16, 32, 64})
		return edakit.Numeric(edakit.Sprintf("An SRAM macro stores %d words of %d bits in 6T cells. How many transistors are in its cell array (excluding decoders, sense amplifiers and other periphery)?", W, B), edakit.IntStr(6*W*B),
			edakit.Sprintf("$6 \\times %d \\times %d = %d$.", W, B, 6*W*B), map[string]string{"bits": edakit.IntStr(W * B)})
	case 2:
		W := 1 << uint(edakit.Int(8, 14))
		B := edakit.Pick([]int{8, 16, 32, 64})
		m := edakit.Pick([]int{2, 4, 8, 16})
		rows := W / m
		switch edakit.Int(0, 3) {
		case 0:
			return edakit.Numeric(edakit.Sprintf("An SRAM of %d words by %d bits uses a column-mux factor of %d. How many physical rows does its array have?", W, B, m), edakit.IntStr(rows),
				edakit.Sprintf("Rows = W / m = %d / %d = %d.", W, m, rows), map[string]string{"columns": edakit.IntStr(B * m)})
		case 1:
			return edakit.Numeric(edakit.Sprintf("An SRAM of %d words by %d bits uses a column-mux factor of %d. How many physical columns does its array have?", W, B, m), edakit.IntStr(B*m),
				edakit.Sprintf("Columns = B times m = %d times %d = %d.", B, m, B*m), map[string]string{"rows": edakit.IntStr(rows)})
		case 2:
			return edakit.Numeric(edakit.Sprintf("An SRAM of %d words by %d bits uses a column-mux factor of %d. How many address bits drive its row decoder?", W, B, m), edakit.IntStr(clog2(rows)),
				edakit.Sprintf("There are %d rows, so $\\log_2 %d = %d$ row-address bits; the other %d bits drive the column mux.", rows, rows, clog2(rows), clog2(m)), map[string]string{"rows": edakit.IntStr(rows)})
		default:
			return edakit.Numeric(edakit.Sprintf("An SRAM of %d words by %d bits uses a column-mux factor of %d. How many address bits drive its column multiplexer?", W, B, m), edakit.IntStr(clog2(m)),
				edakit.Sprintf("$\\log_2 %d = %d$ bits select the word within a row.", m, clog2(m)), map[string]string{"rows": edakit.IntStr(rows)})
		}
	default:
		B := edakit.Pick([]int{8, 16, 32, 64})
		j := edakit.Int(1, 4)
		m := 1 << uint(j)
		W := B * m * m
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(edakit.Sprintf("An SRAM of %d words by %d bits is to be laid out as a square array by choosing its column-mux factor. How many rows will the square array have?", W, B), edakit.IntStr(W/m),
				edakit.Sprintf("Square means $W/m = Bm$, so $m = \\sqrt{%d/%d} = %d$ and rows = %d = columns.", W, B, m, W/m), map[string]string{"mux": edakit.IntStr(m)})
		}
		return edakit.Numeric(edakit.Sprintf("An SRAM of %d words by %d bits is to be laid out as a square array. What column-mux factor achieves that?", W, B), edakit.IntStr(m),
			edakit.Sprintf("$m = \\sqrt{W/B} = \\sqrt{%d} = %d$, giving %d rows by %d columns.", W/B, m, W/m, B*m), map[string]string{"rows": edakit.IntStr(W / m)})
	}
}

func dram(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		win := edakit.Pick([]int{32, 64})
		R := edakit.Pick([]int{4096, 8192, 16384})
		trefi := float64(win) * 1000 / float64(R)
		if edakit.Int(0, 1) == 0 {
			ns := trefi * 1000
			return edakit.Numeric(edakit.Sprintf("A DRAM must receive %d refresh commands, evenly spaced, within every %d ms retention window. What is the interval between refresh commands, in nanoseconds? Give the exact value.", R, win), edakit.Num(ns, 3),
				edakit.Sprintf("$%d\\ \\text{ms} / %d = %s$ ns.", win, R, edakit.Num(ns, 3)), map[string]string{"trefi_us": edakit.Num(trefi, 6)})
		}
		rfc := 10 * edakit.Int(16, 55)
		pct := float64(rfc) / (trefi * 1000) * 100
		return withExact(edakit.Numeric(edakit.Sprintf("A DRAM receives %d refresh commands per %d ms window, evenly spaced, and each keeps the device busy for %d ns. What percentage of the time is spent refreshing? Round to 2 decimal places.", R, win, rfc), edakit.Num(pct, 2),
			edakit.Sprintf("The interval is %s ns, so the overhead is $%d / %s = %s$ percent.", edakit.Num(trefi*1000, 2), rfc, edakit.Num(trefi*1000, 2), edakit.Num(pct, 2)), map[string]string{"trefi_ns": edakit.Num(trefi*1000, 2)}), pct)
	case 2:
		w := edakit.Pick([]int{16, 32, 64, 128})
		r := edakit.Pick([]int{1600, 2400, 2666, 3200, 4266, 6400})
		bw := w * r / 8
		return edakit.Numeric(edakit.Sprintf("A DRAM channel is %d bits wide and runs at %d MT/s (mega-transfers per second). What is its peak bandwidth, in MB/s?", w, r), edakit.IntStr(bw),
			edakit.Sprintf("$%d \\times %d / 8 = %d$ MB/s.", w, r, bw), map[string]string{"bits_per_s": edakit.IntStr(w * r)})
	default:
		cl := edakit.Int(12, 22)
		rcd := edakit.Int(12, 22)
		rp := edakit.Int(12, 22)
		h := 5 * edakit.Int(8, 17)
		e := 5 * edakit.Int(0, (100-h)/5)
		c := 100 - h - e
		lat := (float64(h)*float64(cl) + float64(e)*float64(rcd+cl) + float64(c)*float64(rp+rcd+cl)) / 100
		return withExact(edakit.Numeric(edakit.Sprintf("A DRAM has $t_{CL} = %d$ ns, $t_{RCD} = %d$ ns and $t_{RP} = %d$ ns. Of its accesses, %d%% hit the open row, %d%% find the bank idle (no open row) and %d%% conflict with a different open row. What is the expected access latency, in ns? Round to 2 decimal places.", cl, rcd, rp, h, e, c), edakit.Num(lat, 2),
			edakit.Sprintf("Hit %d ns, idle %d ns, conflict %d ns; weighted: %s ns.", cl, rcd+cl, rp+rcd+cl, edakit.Num(lat, 2)), map[string]string{"conflict_ns": edakit.IntStr(rp + rcd + cl)}), lat)
	}
}

func registerFile(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		w := edakit.Int(1, 8)
		R := edakit.Pick([]int{16, 32, 64, 128, 256})
		switch edakit.Int(0, 2) {
		case 0:
			return edakit.Numeric(edakit.Sprintf("A core issues %d instructions per cycle, each reading up to two registers and writing one. How many register-file ports (read plus write) avoid structural conflicts?", w), edakit.IntStr(3*w),
				edakit.Sprintf("$2 \\times %d$ read plus %d write ports: %d.", w, w, 3*w), map[string]string{"read": edakit.IntStr(2 * w)})
		case 1:
			return edakit.Numeric(edakit.Sprintf("A core issues %d instructions per cycle, each reading up to two source registers. How many read ports does its register file need?", w), edakit.IntStr(2*w),
				edakit.Sprintf("Two per instruction: %d.", 2*w), map[string]string{"write": edakit.IntStr(w)})
		default:
			p := 3 * w
			return edakit.Numeric(edakit.Sprintf("A %d-entry register file has %d ports in total. How many address bits do all ports decode together each cycle?", R, p), edakit.IntStr(p*clog2(R)),
				edakit.Sprintf("Each port takes $\\log_2 %d = %d$ bits: $%d \\times %d = %d$.", R, clog2(R), p, clog2(R), p*clog2(R)), map[string]string{"bits_per_port": edakit.IntStr(clog2(R))})
		}
	case 2:
		a := edakit.Int(1, 4)
		p1 := edakit.Int(2, 6)
		p2 := p1 + edakit.Int(2, 12)
		r := float64((a+p2)*(a+p2)) / float64((a+p1)*(a+p1))
		return withExact(edakit.Numeric(edakit.Sprintf("A register-file cell needs %d wire tracks in each direction plus one per port, so its side grows as %d + p and its area as the square. By what factor does cell area grow going from %d ports to %d ports? Round to 2 decimal places.", a, a, p1, p2), edakit.Num(r, 2),
			edakit.Sprintf("$(%d / %d)^2 = %s$.", a+p2, a+p1, edakit.Num(r, 2)), map[string]string{"side_ratio": edakit.Num(float64(a+p2)/float64(a+p1), 3)}), r)
	default:
		w := edakit.Int(1, 6)
		k := edakit.Int(1, 4)
		return edakit.Numeric(edakit.Sprintf("A %d-wide core forwards results from %d pipeline stages (each stage holding up to %d results) to every ALU operand, in addition to the register-file value. How many sources does each operand mux choose among?", w, k, w), edakit.IntStr(1+w*k),
			edakit.Sprintf("The register file plus %d stages times %d results: $1 + %d \\cdot %d = %d$.", k, w, w, k, 1+w*k), map[string]string{"bypass": edakit.IntStr(w * k)})
	}
}

func tcamMatch(entry, key string) bool {
	for i := range entry {
		if entry[i] != 'X' && entry[i] != key[i] {
			return false
		}
	}
	return true
}

func cam(ctx generator.GeneratorContext) generator.Problem {
	W := edakit.Scaled(ctx, 4, 6)
	E := edakit.Scaled(ctx, 5, 8)
	key := bitString(randBits(W))
	mk := func() string {
		b := []byte(key)
		for i := range b {
			switch edakit.Int(0, 9) {
			case 0, 1, 2:
				b[i] = 'X'
			case 3, 4:
				b[i] = '0' + ('1' - b[i])
			}
		}
		return string(b)
	}
	switch variant(ctx) {
	case 1:
		var ents []string
		cnt := 0
		for i := 0; i < E; i++ {
			e := mk()
			ents = append(ents, edakit.Sprintf("%d: %s", i, e))
			if tcamMatch(e, key) {
				cnt++
			}
		}
		return edakit.Numeric(edakit.Sprintf("A TCAM holds the entries %s (X is don't care). How many entries match the search key %s?", strings.Join(ents, ", "), key), edakit.IntStr(cnt),
			edakit.Sprintf("An entry matches when every non-X bit equals the key bit; the matching count is %d.", cnt), map[string]string{"matches": edakit.IntStr(cnt)})
	case 2:
		var ents []string
		first := -1
		for i := 0; i < E-1; i++ {
			e := mk()
			ents = append(ents, edakit.Sprintf("%d: %s", i, e))
			if first < 0 && tcamMatch(e, key) {
				first = i
			}
		}
		ents = append(ents, edakit.Sprintf("%d: %s", E-1, strings.Repeat("X", W)))
		if first < 0 {
			first = E - 1
		}
		return edakit.Numeric(edakit.Sprintf("A TCAM holds entries in priority order, %s, and its priority encoder returns the lowest-indexed match. What index does a search for %s return?", strings.Join(ents, ", "), key), edakit.IntStr(first),
			edakit.Sprintf("Scanning from entry 0, the first entry whose specified bits all equal the key is entry %d.", first), map[string]string{"index": edakit.IntStr(first)})
	default:
		E2 := edakit.Pick([]int{64, 128, 256, 512, 1024, 2048})
		W2 := edakit.Pick([]int{16, 32, 48, 64, 128})
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(edakit.Sprintf("A TCAM has %d entries of %d ternary bits, each ternary cell storing two bits. How many storage bits does it hold?", E2, W2), edakit.IntStr(2*E2*W2),
				edakit.Sprintf("$2 \\times %d \\times %d = %d$.", E2, W2, 2*E2*W2), map[string]string{"cells": edakit.IntStr(E2 * W2)})
		}
		return edakit.Numeric(edakit.Sprintf("A binary CAM has %d entries of %d bits. How many bit comparisons does one search perform?", E2, W2), edakit.IntStr(E2*W2),
			edakit.Sprintf("Every bit of every entry is compared in parallel: $%d \\times %d = %d$.", E2, W2, E2*W2), map[string]string{"entries": edakit.IntStr(E2)})
	}
}

var macroItems = []edakit.Item{
	{KP: 1, Q: "A placer needs a memory macro's size, pin locations and routing blockages, but not its transistors. Which view supplies them?",
		A: "The abstract physical view (LEF)", D: []string{"The Liberty timing model", "The behavioral Verilog model", "The GDSII layout read transistor by transistor"},
		Why: "LEF abstracts the macro to its outline, pins and obstructions, which is all placement and routing need."},
	{KP: 1, Q: "A static timer analyzes a path ending at an SRAM macro's address pins. Where do the macro's setup requirement and clock-to-output delay come from?",
		A: "The macro's Liberty model, characterized per corner", D: []string{"Timing the macro's internal transistors with the gate-level timer", "The LEF file's pin geometry", "The behavioral model's delay statements"},
		Why: "The macro is a black box to STA; its arcs and checks are characterized by the compiler and delivered in Liberty."},
	{KP: 1, Q: "A logic simulator runs a testbench through a design containing a compiled SRAM. What does it simulate for the memory?",
		A: "A behavioral model of the macro's read and write protocol", D: []string{"The extracted transistor netlist", "The LEF abstract", "Nothing: memories are skipped in simulation"},
		Why: "Memory compilers emit a behavioral Verilog model so RTL and gate-level simulation can exercise the macro at speed."},
	{KP: 1, Q: "Why are SRAMs delivered by a memory compiler rather than synthesized from flip-flops?",
		A: "A custom bitcell array is many times denser and lower power than flip-flops", D: []string{"Synthesis tools cannot infer memories at all", "Flip-flops cannot be read and written in the same cycle", "Memory compilers produce faster flip-flops"},
		Why: "A 6T cell is far smaller than a flip-flop plus its mux; flops are used only for very small or heavily ported storage."},
	{KP: 2, Q: "What does a march test such as March C- do in memory BIST?",
		A: "Writes and reads every address in ascending and descending order to expose stuck-at, transition and coupling faults", D: []string{"Measures the access time of the slowest cell", "Writes random data once and compares a checksum", "Refreshes every row to test retention"},
		Why: "March algorithms are sequences of read and write elements over the address space in both directions, designed for known memory fault models."},
	{KP: 2, Q: "MBIST finds one failing row in a large SRAM that has spare rows. What raises yield in this case?",
		A: "Remapping the failing row to a spare through repair logic or fuses", D: []string{"Running the memory at a lower clock", "Discarding the bit with error correction alone", "Re-running synthesis on the memory"},
		Why: "Redundant rows and columns plus built-in repair turn a defective die into a good one."},
	{KP: 2, Q: "Why do memories get dedicated test logic rather than relying on scan chains?",
		A: "Their cells are not scan flip-flops, and their dense layout has fault types scan cannot target", D: []string{"Scan chains cannot reach the chip's top level", "Memories have no clock", "Scan insertion would change the memory's contents"},
		Why: "Bitcells are not part of scan, and coupling or address-decoder faults need algorithmic patterns that MBIST applies at speed."},
	{KP: 2, Q: "An MBIST controller is inserted around a macro. Which part of the flow must account for it?",
		A: "Timing and area, since the test muxes sit on the memory's functional inputs", D: []string{"Only the behavioral model", "Nothing, because test logic is removed after manufacturing", "Only the memory compiler's layout"},
		Why: "The test collar adds muxes in front of address and data pins, which costs area and adds delay on functional paths."},
	{KP: 3, Q: "Which property most distinguishes embedded flash from SRAM?",
		A: "It retains data without power but has slow, high-voltage writes and limited endurance", D: []string{"It is faster to write than SRAM", "It needs refresh every 64 ms", "It needs no additional process steps"},
		Why: "Flash trades write speed, endurance and process cost for non-volatility."},
	{KP: 3, Q: "In MRAM, what physical state stores a bit?",
		A: "The relative magnetization of the layers of a magnetic tunnel junction", D: []string{"Charge on a floating gate", "The crystalline or amorphous phase of a chalcogenide", "A conductive filament in an oxide"},
		Why: "MRAM reads the tunnel junction's resistance, which depends on whether the free and fixed layers are magnetized in parallel or antiparallel."},
	{KP: 3, Q: "A designer needs non-volatile storage for a few kilobits of calibration data written once at test time. Which choice is typical?",
		A: "One-time-programmable fuse or antifuse memory", D: []string{"A large DRAM array", "A TCAM", "A multi-ported register file"},
		Why: "OTP memory is small, cheap in process terms, and suits data written once and read thereafter."},
	{KP: 3, Q: "Why are emerging memories such as ReRAM and MRAM attractive for embedding beside logic?",
		A: "They can be integrated in the back-end metal layers and write faster than flash", D: []string{"They are denser than DRAM in every process", "They have no variability", "They need no sense amplifiers"},
		Why: "Back-end integration avoids many front-end process steps, though variability and sensing margin remain challenges."},
	{KP: 3, Q: "What is a mask ROM's main limitation?",
		A: "Its contents are fixed at fabrication and cannot be changed", D: []string{"It loses data at power-down", "It needs refresh", "It has very poor read speed"},
		Why: "Mask ROM is programmed by the layout itself, so any change needs a new mask."},
}
