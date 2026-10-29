package extract

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// capSection draws the CAP section of a D_NET: ground capacitors on n internal nodes and k
// coupling capacitors to nodes of other nets, values in fF.
func capSection(net string) (string, float64, float64) {
	n := edakit.Int(2, 4)
	k := edakit.Int(1, 3)
	var caps []string
	gsum, csum := 0.0, 0.0
	idx := 1
	for i := 1; i <= n; i++ {
		v := float64(edakit.Int(3, 30)) / 10
		gsum += v
		caps = append(caps, edakit.Sprintf("%d %s:%d %s", idx, net, i, f(v, 1)))
		idx++
	}
	for i := 0; i < k; i++ {
		v := float64(edakit.Int(1, 9)) / 10
		csum += v
		caps = append(caps, edakit.Sprintf("%d %s:%d a%d:%d %s", idx, net, edakit.Int(1, n), edakit.Int(1, 9), edakit.Int(1, 5), f(v, 1)))
		idx++
	}
	return strings.Join(caps, "; "), gsum, csum
}

// spefCalc reads the CAP section of a printed D_NET and asks for its total capacitance,
// which is the number the D_NET header carries.
func spefCalc(ctx generator.GeneratorContext) generator.Problem {
	net := edakit.Sprintf("n%d", edakit.Int(10, 99))
	caps, gsum, csum := capSection(net)
	tot := gsum + csum
	q := edakit.Sprintf("A SPEF file declares C_UNIT 1 FF. The CAP section of net %s lists: %s. What total capacitance belongs in the D_NET header line of %s?", net, caps, net)
	return numChoice(q, tot, 1, " fF",
		edakit.Sprintf("The D_NET total is the sum of every CAP entry, ground and coupling alike: %s + %s = %s fF. Entries with two node names are coupling capacitors to other nets.", f(gsum, 1), f(csum, 1), f(tot, 1)),
		gsum, csum, tot+csum)
}

var spefItems = []edakit.Item{
	{KP: 1, Q: "In a SPEF D_NET section, what does a CAP entry with two node names, one of them on another net, represent?",
		A: "A coupling capacitor between the two nets", D: []string{"A ground capacitor at the first node", "A resistor between the two nodes", "A pin-to-pin delay"},
		Why: "A ground capacitor lists one node; a coupling capacitor lists both terminals, so an SI tool can find the aggressor net."},
	{KP: 1, Q: "What does the CONN section of a D_NET list?",
		A: "The ports and instance pins the net connects, with their directions", D: []string{"Every internal RC node of the net", "The coupling capacitors to neighbours", "The timing arcs of the driving cell"},
		Why: "CONN names the net's external connection points: ports (P entries) and instance pins (I entries) with I, O or B direction. Internal nodes appear only in the CAP and RES sections."},
	{KP: 1, Q: "Why does a SPEF file carry a NAME_MAP section?",
		A: "To replace long hierarchical names with short indices, shrinking the file and speeding parsing", D: []string{"To map logical nets to physical layers", "To record which cells are spare cells", "To store the netlist's connectivity"},
		Why: "Hierarchical names are repeated in every entry. Writing each once in a map and referring to it by index cuts file size substantially."},
	{KP: 1, Q: "Which information does a standard SPEF D_NET not record?",
		A: "The wire's layer geometry and coordinates, unless optional extensions are used", D: []string{"Resistances between nodes", "Ground capacitances", "Coupling capacitances to other nets"},
		Why: "SPEF is an electrical netlist of R and C elements. Geometry is abstracted away, which is why a tool cannot re-extract from SPEF alone."},
	{KP: 1, Q: "What do the header lines T_UNIT, C_UNIT and R_UNIT do?",
		A: "They scale every time, capacitance and resistance value in the file", D: []string{"They set the analysis corner", "They declare the netlist's hierarchy separator", "They list the maximum allowed values"},
		Why: "Values in the body are bare numbers; the header units give them meaning. Misreading C_UNIT 1 PF as femtofarads is a thousandfold error."},
	{KP: 1, Q: "A timer reads a reduced SPEF in which each net is only a total capacitance and no resistors. What can it no longer compute well?",
		A: "Wire delay and slew degradation along the net, which need the RC network", D: []string{"The driving cell's intrinsic delay", "Clock period checks", "Logic levels of the netlist"},
		Why: "A lumped capacitance gives the driver's load but no interconnect delay to each sink. Detailed or moment-preserving reduced RC is needed for that."},
	{KP: 3, Q: "One SPEF file is produced per extraction corner. Why?",
		A: "Resistance and capacitance values depend on the process and temperature corner, so each corner has its own parasitics", D: []string{"SPEF can only hold one net per file", "Different corners use different netlists", "Corners differ only in their name maps"},
		Why: "Cworst, Cbest, RCworst, RCbest and typical give different wire dimensions and dielectric values, and temperature changes resistance, so each corner is a separate parasitic database."},
	{KP: 3, Q: "A timer reports that 2% of nets in a SPEF have a D_NET total that differs from the sum of their CAP entries. What is the most likely cause?",
		A: "A writer bug or truncated file, since the header total should equal the sum of the entries", D: []string{"Normal rounding of coupling to ground", "Miller factors applied in the file", "Different temperature corners per net"},
		Why: "The D_NET total is a checksum-like summary. A mismatch means the file is inconsistent and should be rejected or regenerated."},
	{KP: 3, Q: "Why must the instance and pin names in the SPEF match the netlist the timer loaded?",
		A: "The timer annotates parasitics by name, so unmatched names leave nets unannotated and timed with estimates", D: []string{"SPEF names are only comments", "The timer re-derives the names from geometry", "Mismatched names only affect power analysis"},
		Why: "Back-annotation is a name match. A renamed hierarchy after an ECO, without a matching SPEF, silently falls back to wire-load estimates on the affected nets."},
	{KP: 3, Q: "A flow grounds all coupling capacitors when writing SPEF to save space. What analysis does this disable?",
		A: "Crosstalk delay and glitch analysis, which need to know each aggressor", D: []string{"Basic setup and hold timing", "Power-grid IR drop", "Clock skew reporting"},
		Why: "Grounded coupling preserves total load but loses which net is on the other side, which SI analysis needs."},
	{KP: 3, Q: "A coupling capacitor between nets A and B may be written in the D_NET of A, of B, or of both. What must a reader building a coupled RC network do?",
		A: "Merge the entries so that each physical coupling capacitor is instantiated once between the two nets", D: []string{"Add every listed entry, doubling couplings that appear twice", "Ground every coupling capacitor that appears twice", "Keep only the entries written in the victim's section and discard the rest"},
		Why: "The same element can be reported from either side. Instantiating it twice doubles the coupling and inflates crosstalk; dropping one side can lose it entirely."},
	{KP: 3, Q: "Why is DSPF still used alongside SPEF in some flows?",
		A: "It is a SPICE-like detailed format that circuit simulators read directly for transistor-level analysis", D: []string{"It is the only format with coupling capacitors", "It stores layout geometry", "It is required by STA tools"},
		Why: "DSPF writes parasitics in SPICE syntax, convenient for simulating a custom block. SPEF is the compact standard for timing and SI tools."},
}

// spefCounts counts entries for a described net (KP 1), totals a printed CAP section (KP 2),
// and computes the Elmore delay from a printed chain with header units (KP 3).
func spefCounts(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		k := edakit.Scaled(ctx, 1, 5)
		m := edakit.Scaled(ctx, 2, 12)
		c := edakit.Int(0, 6)
		stem := edakit.Sprintf("An extractor writes the D_NET of a net with one driver pin and %s. Its resistor network is a tree spanning the pins and %d internal nodes. It puts one ground capacitor on every internal node and none on the pins, and the net has %s to other nets.", plural(k, "sink pin"), m, plural(c, "coupling capacitor"))
		switch edakit.Int(1, 3) {
		case 1:
			return num(stem+" How many entries does the RES section have?", float64(k+m), 0,
				edakit.Sprintf("A tree on %d + %d + 1 = %d nodes has one fewer edge: %d resistors.", k, m, k+m+1, k+m), nil)
		case 2:
			return num(stem+" How many entries does the CAP section have?", float64(m+c), 0,
				edakit.Sprintf("%d ground capacitors plus %d coupling capacitors: %d entries.", m, c, m+c), nil)
		}
		return num(stem+" How many entries does the CONN section have?", float64(k+1), 0,
			edakit.Sprintf("CONN lists the connection points only: the driver and %d sinks, %d entries. Internal nodes are not listed there.", k, k+1), nil)
	case 2:
		unit := edakit.Pick([]string{"FF", "PF"})
		if unit == "FF" {
			net := edakit.Sprintf("n%d", edakit.Int(10, 99))
			caps, gsum, csum := capSection(net)
			return num(edakit.Sprintf("A SPEF file declares C_UNIT 1 FF. The CAP section of net %s lists: %s. What total capacitance, in fF, belongs in the D_NET header line of %s?", net, caps, net),
				gsum+csum, 1, edakit.Sprintf("Sum every CAP entry: ground %s fF plus coupling %s fF gives %s fF.", f(gsum, 1), f(csum, 1), f(gsum+csum, 1)), nil)
		}
		// The same question in picofarads: rebuild with values scaled.
		net := edakit.Sprintf("n%d", edakit.Int(10, 99))
		n := edakit.Int(2, 4)
		var caps []string
		sum := 0.0
		for i := 1; i <= n; i++ {
			v := float64(edakit.Int(3, 30)) / 10000
			sum += v
			caps = append(caps, edakit.Sprintf("%d %s:%d %s", i, net, i, f(v, 4)))
		}
		cv := float64(edakit.Int(1, 9)) / 10000
		sum += cv
		caps = append(caps, edakit.Sprintf("%d %s:%d a%d:2 %s", n+1, net, edakit.Int(1, n), edakit.Int(1, 9), f(cv, 4)))
		return num(edakit.Sprintf("A SPEF file declares C_UNIT 1 PF. The CAP section of net %s lists: %s. What is the net's total capacitance, in fF?", net, strings.Join(caps, "; ")),
			sum*1000, 1, edakit.Sprintf("The entries sum to %s pF, which is %s fF. The header unit scales every value.", f(sum, 4), f(sum*1000, 1)), nil)
	}
	net := edakit.Sprintf("n%d", edakit.Int(10, 99))
	n := edakit.Scaled(ctx, 2, 4)
	rs := make([]float64, n+1)
	cs := make([]float64, n+1)
	var rl, cl []string
	nodes := []string{"u1:Z"}
	for i := 1; i <= n; i++ {
		nodes = append(nodes, edakit.Sprintf("%s:%d", net, i))
	}
	nodes = append(nodes, "u2:A")
	for i := 0; i <= n; i++ {
		rs[i] = float64(edakit.Int(1, 20)) / 100
		rl = append(rl, edakit.Sprintf("%d %s %s %s", i+1, nodes[i], nodes[i+1], f(rs[i], 2)))
	}
	for i := 1; i <= n+1; i++ {
		cs[i-1] = float64(edakit.Int(5, 40)) / 10
		cl = append(cl, edakit.Sprintf("%d %s %s", i, nodes[i], f(cs[i-1], 1)))
	}
	// Elmore at u2:A: sum over resistors of R_i times capacitance downstream of it.
	el := 0.0
	for i := 0; i <= n; i++ {
		down := 0.0
		for j := i; j <= n; j++ {
			down += cs[j]
		}
		el += rs[i] * down
	}
	return num(edakit.Sprintf("A SPEF file declares R_UNIT 1 KOHM and C_UNIT 1 FF. Net %s is driven by u1:Z and drives u2:A. Its CAP section lists %s, and its RES section lists %s. What is the Elmore delay from u1:Z to u2:A, in ps?", net, strings.Join(cl, "; "), strings.Join(rl, "; ")),
		el, 3, edakit.Sprintf("The nodes form a chain. Each resistor charges all capacitance downstream of it; one kohm times one fF is one ps. The sum is %s ps.", f(el, 3)), nil)
}

// nameMap computes bytes saved by a SPEF name map.
func nameMap(ctx generator.GeneratorContext) generator.Problem {
	l := edakit.Int(30, 90)
	d := edakit.Int(3, 7)
	rule := edakit.Sprintf("In a SPEF file, a net's full hierarchical name is %d characters long. With a name map, each reference is replaced by an index token of %d characters (the asterisk and digits), and the map adds one line per net: the token, a space, the full name and a newline.", l, d)
	entry := d + 1 + l + 1
	switch kp(ctx) {
	case 1:
		r := edakit.Int(4, 20)
		saved := r*(l-d) - entry
		return num(rule+edakit.Sprintf(" The net's name appears %d times in the body. How many bytes does the name map save for this net?", r), float64(saved), 0,
			edakit.Sprintf("Each reference saves %d - %d = %d bytes: %d x %d = %d. The map entry costs %d + 1 + %d + 1 = %d bytes. Net saving %d bytes.", l, d, l-d, r, l-d, r*(l-d), d, l, entry, saved), nil)
	case 2:
		r := edakit.Int(4, 20)
		nets := edakit.Pick([]int{1, 2, 5, 10, 20}) * 1000000
		saved := float64(nets) * float64(r*(l-d)-entry) / 1e6
		return num(rule+edakit.Sprintf(" A design has %d million nets, each named %d times in the body on average, and all names have this length. How many MB does the name map save (1 MB is one million bytes)?", nets/1000000, r), saved, 0,
			edakit.Sprintf("Per net: %d x %d - %d = %d bytes. Times %d nets: %s MB. Parse time falls with it, since the reader hashes each long name once.", r, l-d, entry, r*(l-d)-entry, nets, f(saved, 0)), nil)
	}
	rmin := int(math.Floor(float64(entry)/float64(l-d))) + 1
	return num(rule+" What is the smallest number of references to the net in the body for which the name map saves bytes overall?", float64(rmin), 0,
		edakit.Sprintf("Need $r(%d - %d) > %d$, so $r > %s$; the smallest integer is %d. Every net in a SPEF appears at least in its D_NET line and in its node names, so in practice the map always pays.", l, d, entry, f(float64(entry)/float64(l-d), 3), rmin), nil)
}
