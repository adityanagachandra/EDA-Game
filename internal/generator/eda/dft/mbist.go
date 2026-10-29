package dft

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// marchElem is one March element: an address order and a list of operations such as "r0",
// "w1".
type marchElem struct {
	order string // "up", "down", "any"
	ops   []string
}

type marchAlg struct {
	name  string
	elems []marchElem
}

func me(order string, ops ...string) marchElem { return marchElem{order, ops} }

var marchAlgs = []marchAlg{
	{"MATS+", []marchElem{me("any", "w0"), me("up", "r0", "w1"), me("down", "r1", "w0")}},
	{"MATS++", []marchElem{me("any", "w0"), me("up", "r0", "w1"), me("down", "r1", "w0", "r0")}},
	{"March X", []marchElem{me("any", "w0"), me("up", "r0", "w1"), me("down", "r1", "w0"), me("any", "r0")}},
	{"March Y", []marchElem{me("any", "w0"), me("up", "r0", "w1", "r1"), me("down", "r1", "w0", "r0"), me("any", "r0")}},
	{"March C-", []marchElem{me("any", "w0"), me("up", "r0", "w1"), me("up", "r1", "w0"), me("down", "r0", "w1"), me("down", "r1", "w0"), me("any", "r0")}},
	{"March A", []marchElem{me("any", "w0"), me("up", "r0", "w1", "w0", "w1"), me("up", "r1", "w0", "w1"), me("down", "r1", "w0", "w1", "w0"), me("down", "r0", "w1", "w0")}},
	{"March B", []marchElem{me("any", "w0"), me("up", "r0", "w1", "r1", "w0", "r0", "w1"), me("up", "r1", "w0", "w1"), me("down", "r1", "w0", "w1", "w0"), me("down", "r0", "w1", "w0")}},
}

func (a marchAlg) text() string {
	arrow := map[string]string{"up": "\\Uparrow", "down": "\\Downarrow", "any": "\\Updownarrow"}
	var p []string
	for _, e := range a.elems {
		p = append(p, arrow[e.order]+"("+strings.Join(e.ops, ",")+")")
	}
	return "$\\{" + strings.Join(p, ";\\ ") + "\\}$"
}

func (a marchAlg) opsPerCell() int {
	n := 0
	for _, e := range a.elems {
		n += len(e.ops)
	}
	return n
}

// randMarch draws a consistent March algorithm: each read expects the value last written.
func randMarch(ctx generator.GeneratorContext) marchAlg {
	ne := edakit.Scaled(ctx, 3, 6)
	elems := []marchElem{me("any", "w0")}
	v := 0
	for i := 1; i < ne; i++ {
		order := edakit.Pick([]string{"up", "down"})
		ops := []string{edakit.Sprintf("r%d", v)}
		k := edakit.Int(1, 4)
		for j := 0; j < k; j++ {
			if edakit.Int(0, 2) == 0 {
				ops = append(ops, edakit.Sprintf("r%d", v))
			} else {
				v = 1 - v
				ops = append(ops, edakit.Sprintf("w%d", v))
			}
		}
		elems = append(elems, me(order, ops...))
	}
	if edakit.Int(0, 1) == 1 {
		elems = append(elems, me("any", edakit.Sprintf("r%d", v)))
	}
	return marchAlg{"", elems}
}

// march: ops per cell (KP 1), total test time (KP 2), word-oriented backgrounds (KP 3).
func march(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		a := randMarch(ctx)
		c := a.opsPerCell()
		return edakit.Numeric("A March test is written "+a.text()+", where $\\Uparrow$ and $\\Downarrow$ are ascending and descending address orders and $\\Updownarrow$ means either. How many operations does it apply per memory cell (its complexity is that number times n)?",
			edakit.IntStr(c), edakit.Sprintf("Each element applies all its operations to every address, so the complexity is the total operation count, %dn.", c), nil)
	case 2:
		a := edakit.Pick(marchAlgs)
		kw := edakit.Pick([]int{1, 4, 16, 64, 256, 1024})
		n := kw * 1024
		ns := edakit.Pick([]int{1, 2, 5, 10})
		us := float64(a.opsPerCell()*n*ns) / 1000
		ans := edakit.Num(us, 2)
		return edakit.Numeric(edakit.Sprintf("%s is %s. MBIST runs it on a %dK-word memory (1K = 1024) with one operation per %d ns cycle. How long does the test take, in microseconds? Round to 2 decimal places.", a.name, a.text(), kw, ns),
			ans, edakit.Sprintf("%s is %dn: $%d \\times %d \\times %d$ ns = %s microseconds.", a.name, a.opsPerCell(), a.opsPerCell(), n, ns, ans), nil)
	default:
		a := edakit.Pick(marchAlgs)
		b := edakit.Pick([]int{4, 8, 16, 32, 64})
		words := edakit.Pick([]int{256, 512, 1024, 2048, 4096})
		bg := int(math.Log2(float64(b))) + 1
		ans := a.opsPerCell() * words * bg
		return edakit.Numeric(edakit.Sprintf("%s (%dn) is applied to a word-oriented memory of %d words of %d bits. To catch intra-word coupling faults it is repeated once per data background, using $\\log_2 B + 1$ backgrounds for B-bit words. How many read and write operations are applied in total?", a.name, a.opsPerCell(), words, b),
			edakit.IntStr(ans), edakit.Sprintf("$%d \\times %d \\times (\\log_2 %d + 1) = %d \\times %d \\times %d = %d$.", a.opsPerCell(), words, b, a.opsPerCell(), words, bg, ans), nil)
	}
}

type memFault struct {
	kind    string // saf, tf, cfin, cfid
	v, a    int    // victim, aggressor
	val, up int    // stuck value or forced value; up = 1 for rising trigger or TF rising
	desc    string
}

func randMemFault(n int, allowCF bool) memFault {
	v := edakit.Int(0, n-1)
	kinds := []string{"saf", "tf"}
	if allowCF {
		kinds = append(kinds, "cfin", "cfid", "cfin", "cfid")
	}
	switch edakit.Pick(kinds) {
	case "saf":
		x := edakit.Int(0, 1)
		return memFault{kind: "saf", v: v, val: x, desc: edakit.Sprintf("cell %d is stuck-at-%d", v, x)}
	case "tf":
		up := edakit.Int(0, 1)
		d := "cannot make a 1 to 0 transition (down transition fault)"
		if up == 1 {
			d = "cannot make a 0 to 1 transition (up transition fault)"
		}
		return memFault{kind: "tf", v: v, up: up, desc: edakit.Sprintf("cell %d %s", v, d)}
	case "cfin":
		a := (v + edakit.Int(1, n-1)) % n
		up := edakit.Int(0, 1)
		t := "falling"
		if up == 1 {
			t = "rising"
		}
		return memFault{kind: "cfin", v: v, a: a, up: up, desc: edakit.Sprintf("an inversion coupling fault: a %s transition written into cell %d inverts cell %d", t, a, v)}
	}
	a := (v + edakit.Int(1, n-1)) % n
	up := edakit.Int(0, 1)
	x := edakit.Int(0, 1)
	t := "falling"
	if up == 1 {
		t = "rising"
	}
	return memFault{kind: "cfid", v: v, a: a, up: up, val: x, desc: edakit.Sprintf("an idempotent coupling fault: a %s transition written into cell %d forces cell %d to %d", t, a, v, x)}
}

// runMarch returns the 1-based index of the first failing read, or 0.
func runMarch(alg marchAlg, n int, f memFault) (int, string) {
	mem := make([]int, n)
	if f.kind == "saf" {
		mem[f.v] = f.val
	}
	idx := 0
	for ei, e := range alg.elems {
		addrs := seq(n)
		if e.order == "down" {
			for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
				addrs[i], addrs[j] = addrs[j], addrs[i]
			}
		}
		for _, ad := range addrs {
			for _, op := range e.ops {
				idx++
				d := int(op[1] - '0')
				if op[0] == 'r' {
					if mem[ad] != d {
						return idx, edakit.Sprintf("operation %d is %s at address %d in element %d and reads %d", idx, op, ad, ei+1, mem[ad])
					}
					continue
				}
				old := mem[ad]
				nv := d
				if f.kind == "saf" && ad == f.v {
					nv = f.val
				}
				if f.kind == "tf" && ad == f.v {
					if f.up == 1 && old == 0 && d == 1 {
						nv = 0
					}
					if f.up == 0 && old == 1 && d == 0 {
						nv = 1
					}
				}
				mem[ad] = nv
				if (f.kind == "cfin" || f.kind == "cfid") && ad == f.a && old != nv {
					rising := nv == 1
					if (f.up == 1) == rising {
						if f.kind == "cfin" {
							mem[f.v] ^= 1
						} else {
							mem[f.v] = f.val
						}
					}
				}
			}
		}
	}
	return 0, "no read ever mismatches, so this test does not detect the fault"
}

// marchSim: first failing operation for a cell fault (KP 1, 2) or a coupling fault (KP 3).
func marchSim(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 8)
	k := kp(ctx)
	var alg marchAlg
	switch k {
	case 1:
		alg = marchAlgs[edakit.Int(0, 2)]
	default:
		alg = edakit.Pick(marchAlgs)
	}
	f := randMemFault(n, k == 3)
	idx, why := runMarch(alg, n, f)
	return edakit.Numeric(edakit.Sprintf("%s, %s, runs on a memory of %d one-bit cells at addresses 0 to %d ($\\Uparrow$ ascending; $\\Updownarrow$ also taken ascending). Operations are numbered from 1 in execution order: element by element, address by address, operation by operation. The memory powers up holding all 0s, and %s. Which operation number is the first read that returns the wrong value? Answer 0 if no read fails.", alg.name, alg.text(), n, n-1, f.desc),
		edakit.IntStr(idx), "Simulating the memory, "+why+".", map[string]string{"first_fail": edakit.IntStr(idx)})
}

// repair: must-repair lines (KP 1), minimum spare lines (KP 2), minimum columns given spare
// rows (KP 3).
func repair(ctx generator.GeneratorContext) generator.Problem {
	R, C := 8, 8
	nf := edakit.Scaled(ctx, 4, 9)
	cells := map[[2]int]bool{}
	r0 := edakit.Int(0, R-1)
	for i := 0; i < 3; i++ {
		cells[[2]int{r0, edakit.Int(0, C-1)}] = true
	}
	if edakit.Int(0, 1) == 1 {
		c0 := edakit.Int(0, C-1)
		for i := 0; i < 3; i++ {
			cells[[2]int{edakit.Int(0, R-1), c0}] = true
		}
	}
	for len(cells) < nf {
		cells[[2]int{edakit.Int(0, R-1), edakit.Int(0, C-1)}] = true
	}
	var list [][2]int
	for r := 0; r < R; r++ {
		for c := 0; c < C; c++ {
			if cells[[2]int{r, c}] {
				list = append(list, [2]int{r, c})
			}
		}
	}
	var lt []string
	for _, x := range list {
		lt = append(lt, edakit.Sprintf("(%d, %d)", x[0], x[1]))
	}
	stem := edakit.Sprintf("Memory BIST on an 8 by 8 array logs faulty cells at (row, column) %s.", strings.Join(lt, ", "))
	colsLeft := func(rowMask int) int {
		cs := map[int]bool{}
		for _, x := range list {
			if rowMask>>x[0]&1 == 0 {
				cs[x[1]] = true
			}
		}
		return len(cs)
	}
	switch kp(ctx) {
	case 1:
		sr, sc := edakit.Int(1, 2), edakit.Int(1, 2)
		rowc, colc := map[int]int{}, map[int]int{}
		for _, x := range list {
			rowc[x[0]]++
			colc[x[1]]++
		}
		mr, mc := 0, 0
		for _, v := range rowc {
			if v > sc {
				mr++
			}
		}
		for _, v := range colc {
			if v > sr {
				mc++
			}
		}
		return edakit.Numeric(stem+edakit.Sprintf(" The array has spare rows: %d, spare columns: %d. In the first must-repair pass, a row with more faulty cells than there are spare columns must take a spare row, and a column with more faulty cells than spare rows must take a spare column. How many lines are must-repair in that pass?", sr, sc),
			edakit.IntStr(mr+mc), edakit.Sprintf("Rows with more than %d faults: %d. Columns with more than %d faults: %d. Total %d.", sc, mr, sr, mc, mr+mc), nil)
	case 2:
		best := 1 << 30
		for m := 0; m < 1<<R; m++ {
			if c := popcount(m) + colsLeft(m); c < best {
				best = c
			}
		}
		return edakit.Numeric(stem+" With unlimited spare rows and columns, what is the smallest total number of spare lines (rows plus columns) that covers every faulty cell?",
			edakit.IntStr(best), edakit.Sprintf("This is a minimum vertex cover of the bipartite row-column fault graph, which by König's theorem equals its maximum matching: %d lines.", best), nil)
	default:
		sr := edakit.Int(1, 2)
		best := 1 << 30
		for m := 0; m < 1<<R; m++ {
			if popcount(m) <= sr {
				if c := colsLeft(m); c < best {
					best = c
				}
			}
		}
		return edakit.Numeric(stem+edakit.Sprintf(" Repair may use at most %d spare row(s). What is the fewest spare columns needed to cover the remaining faulty cells?", sr),
			edakit.IntStr(best), edakit.Sprintf("Trying every choice of at most %d rows and counting the distinct columns left gives %d columns at best.", sr, best), nil)
	}
}

// ecc: SEC-DED check bits (KP 1), Hamming syndrome (KP 2), uncorrectable-word probability (KP 3).
func ecc(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		m := edakit.Pick([]int{8, 16, 26, 32, 57, 64, 100, 120, 128, 256})
		if edakit.Int(0, 2) == 0 {
			m = edakit.Int(5, 200)
		}
		r := 1
		for (1 << r) < m+r+1 {
			r++
		}
		return edakit.Numeric(edakit.Sprintf("A memory word holds %d data bits. How many check bits does a SEC-DED code (single-error-correcting, double-error-detecting extended Hamming code) need?", m),
			edakit.IntStr(r+1), edakit.Sprintf("SEC needs the smallest r with $2^r \\ge m + r + 1$: r = %d for m = %d. One overall parity bit adds double-error detection: %d.", r, m, r+1), nil)
	case 2:
		n := 7
		if ctx.Difficulty > 0.6 {
			n = 15
		}
		w := make([]int, n+1)
		for i := 1; i <= n; i++ {
			if i&(i-1) != 0 {
				w[i] = edakit.Int(0, 1)
			}
		}
		for p := 1; p <= n; p <<= 1 {
			x := 0
			for i := 1; i <= n; i++ {
				if i&p != 0 && i != p {
					x ^= w[i]
				}
			}
			w[p] = x
		}
		e := edakit.Int(0, n)
		if e > 0 {
			w[e] ^= 1
		}
		s := 0
		for i := 1; i <= n; i++ {
			if w[i] == 1 {
				s ^= i
			}
		}
		return edakit.Numeric(edakit.Sprintf("A Hamming(%d, %d) code puts check bits at positions 1, 2, 4%s with even parity, each check bit covering the positions whose binary index contains it. The word read from memory, positions 1 to %d, is %s. What is the syndrome, the position of the single flipped bit (0 if none)?", n, n-int(math.Log2(float64(n+1))), map[bool]string{true: ", 8", false: ""}[n == 15], n, bits(w[1:])),
			edakit.IntStr(s), edakit.Sprintf("The syndrome is the XOR of the indices of all 1 bits, which is %d.", s), map[string]string{"syndrome": edakit.IntStr(s)})
	default:
		n := edakit.Pick([]int{39, 72, 137})
		p := float64(edakit.Int(1, 50)) * 1e-5
		p2 := 1 - math.Pow(1-p, float64(n)) - float64(n)*p*math.Pow(1-p, float64(n-1))
		ans := edakit.Num(p2*1e6, 3)
		return edakit.Numeric(edakit.Sprintf("Each bit of a %d-bit SEC-DED codeword flips independently with probability %s before it is read. SEC corrects one flip. What is the probability that a word has two or more flips and cannot be corrected, in parts per million? Round to 3 decimal places.", n, edakit.Num(p, 5)),
			ans, edakit.Sprintf("$1 - (1-p)^{%d} - %d p (1-p)^{%d} = %s$ ppm. Periodic scrubbing keeps single errors from accumulating into doubles.", n, n, n-1, ans), nil)
	}
}

var memFaultItems = []edakit.Item{
	{KP: 1, Q: "A cell can be written to 1 but a write of 0 after a 1 leaves it at 1. Which memory fault is this?",
		A: "A down transition fault", D: []string{"A stuck-at-1 fault", "An inversion coupling fault", "An address decoder fault"},
		Why: "The cell can hold either value from initialization but cannot make the 1 to 0 transition; a stuck-at-1 cell could never be 0 at all."},
	{KP: 1, Q: "Writing address 5 also writes address 13. Which functional fault class is this?",
		A: "An address decoder fault", D: []string{"A stuck-at fault in cell 5", "A transition fault", "A data retention fault"},
		Why: "The decoder selects two cells for one address; March tests detect such faults with ascending and descending address orders."},
	{KP: 1, Q: "Why are memory faults modeled functionally (cells, decoders, read/write logic) rather than as gate-level stuck-at faults?",
		A: "Memory arrays are dense analog structures whose defects show up as cell behaviours, not as gate-level line faults", D: []string{"Memories contain no transistors", "Gate-level models are too optimistic for logic", "Stuck-at faults cannot occur in memory"},
		Why: "A bit cell is a cross-coupled pair plus access devices; its defects are best described by what reads and writes do."},
	{KP: 1, Q: "A cell loses its stored 1 after about 50 ms with no access. Which fault is this, and what does testing it need?",
		A: "A data retention fault, which needs a pause before the read", D: []string{"A transition fault, detected by any March test", "A stuck-at-0 fault", "An address decoder fault"},
		Why: "The cell works for immediate reads; only waiting long enough exposes the leakage, so retention tests add delay elements."},
	{KP: 2, Q: "A rising write to cell A inverts cell B, whatever B held. Which fault is this?",
		A: "An inversion coupling fault", D: []string{"An idempotent coupling fault", "A state coupling fault", "A stuck-open fault"},
		Why: "Inversion coupling flips the victim on the aggressor's transition; idempotent coupling forces it to a fixed value."},
	{KP: 2, Q: "Cell B reads 0 whenever cell A holds 1, regardless of transitions. Which fault is this?",
		A: "A state coupling fault", D: []string{"An inversion coupling fault", "A transition fault", "A data retention fault"},
		Why: "State coupling depends on the aggressor's value, not on a write event, which distinguishes it from inversion and idempotent coupling."},
	{KP: 2, Q: "Why do coupling faults matter more as SRAM cells shrink?",
		A: "Adjacent cells share closer wires and wells, so a write to one can disturb its neighbours", D: []string{"Smaller cells cannot hold charge at all", "Decoders become larger", "Reads become destructive in SRAM"},
		Why: "Physical proximity drives capacitive and leakage coupling between neighbours; scrambled address maps make physical neighbours non-adjacent logically."},
	{KP: 2, Q: "A neighbourhood pattern-sensitive fault makes a cell fail only when its four physical neighbours hold a specific pattern. Why is it expensive to test?",
		A: "The test must create many neighbourhood patterns for every cell, which needs the physical address scrambling and many more operations", D: []string{"It needs only a single March element", "It is detected by stuck-at tests", "It cannot occur in DRAM"},
		Why: "Pattern-sensitive tests depend on physical topology and need far more than linear-time March tests to cover all patterns."},
	{KP: 3, Q: "Reading a cell returns the correct value but flips the stored value. Which fault is this?",
		A: "A read destructive fault", D: []string{"A stuck-at fault", "A write disturb on a neighbour", "An address decoder fault"},
		Why: "The first read passes; a second read of the same cell is needed to see the damage, which is why some March tests read twice."},
	{KP: 3, Q: "Why do March tests apply both an ascending and a descending address sweep?",
		A: "Aggressor-victim coupling and decoder faults are only exposed when the victim is read after or before the aggressor is written, depending on order", D: []string{"To halve the test time", "To test the sense amplifier only", "Because the memory must be read twice"},
		Why: "With the victim at a lower address than the aggressor, only one sweep direction writes the aggressor before reading the victim."},
	{KP: 3, Q: "An engineer adds a fault model for a static open in a word line. Which March behaviour reveals it?",
		A: "Every cell on that word line fails, so writes and reads to those addresses mismatch", D: []string{"Only one bit flips after a delay", "Only the column decoder reports an error", "No March test can detect it"},
		Why: "A word-line open disconnects a whole row of cells from access, a large, easily detected functional failure."},
	{KP: 3, Q: "Why are DRAM and SRAM fault lists not identical?",
		A: "DRAM stores charge on a capacitor that leaks and needs refresh, adding retention and refresh-related faults", D: []string{"SRAM has no cells", "DRAM has no decoders", "SRAM cells cannot be coupled"},
		Why: "Charge storage brings leakage, refresh and sense-margin faults that a static latch-based SRAM cell does not have."},
}
