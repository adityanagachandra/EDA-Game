package techmap

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// snode is a node of a NAND2/INV subject graph or of a pattern graph.
// kind: 'i' input (subject), 'x' free leaf (pattern), 'n' NAND2, 'v' INV.
type snode struct {
	kind byte
	name string
	kids []*snode
}

func pin(name string) *snode { return &snode{kind: 'i', name: name} }
func px() *snode             { return &snode{kind: 'x'} }
func pn(a, b *snode) *snode  { return &snode{kind: 'n', kids: []*snode{a, b}} }
func pv(a *snode) *snode     { return &snode{kind: 'v', kids: []*snode{a}} }
func countKind(s *snode, k byte) int {
	n := 0
	if s.kind == k {
		n = 1
	}
	for _, c := range s.kids {
		n += countKind(c, k)
	}
	return n
}

// expr is an AND/OR/NOT expression over named inputs, with explicit association.
type expr struct {
	op   byte // 'v' variable, '&', '|', '!'
	a, b *expr
	v    string
}

func (e *expr) String() string {
	switch e.op {
	case 'v':
		return e.v
	case '!':
		if e.a.op == 'v' {
			return "NOT " + e.a.String()
		}
		return "NOT (" + e.a.String() + ")"
	}
	w := " AND "
	if e.op == '|' {
		w = " OR "
	}
	return paren(e.a) + w + paren(e.b)
}

func paren(e *expr) string {
	if e.op == '&' || e.op == '|' {
		return "(" + e.String() + ")"
	}
	return e.String()
}

// decomposer turns an expression into a NAND2/INV tree and counts cancelled inverter pairs.
type decomposer struct{ cancels int }

func (d *decomposer) inv(x *snode) *snode {
	if x.kind == 'v' {
		d.cancels++
		return x.kids[0]
	}
	return pv(x)
}

func (d *decomposer) build(e *expr) *snode {
	switch e.op {
	case 'v':
		return pin(e.v)
	case '!':
		return d.inv(d.build(e.a))
	case '&':
		return d.inv(pn(d.build(e.a), d.build(e.b)))
	default:
		return pn(d.inv(d.build(e.a)), d.inv(d.build(e.b)))
	}
}

func countOps(e *expr) (and, or, not int) {
	if e == nil || e.op == 'v' {
		return
	}
	a1, o1, n1 := countOps(e.a)
	a2, o2, n2 := countOps(e.b)
	and, or, not = a1+a2, o1+o2, n1+n2
	switch e.op {
	case '&':
		and++
	case '|':
		or++
	case '!':
		not++
	}
	return
}

// randExpr draws an expression with nops binary operators over fresh inputs.
func randExpr(nops int, next *int) *expr {
	var e *expr
	if nops == 0 {
		e = &expr{op: 'v', v: piName(*next)}
		*next++
	} else {
		left := edakit.Int(0, nops-1)
		op := byte('&')
		if edakit.Int(0, 1) == 1 {
			op = '|'
		}
		e = &expr{op: op, a: randExpr(left, next), b: randExpr(nops-1-left, next)}
	}
	if edakit.Int(1, 10) <= 3 {
		e = &expr{op: '!', a: e}
	}
	return e
}

func decomposition(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	var e *expr
	var root *snode
	var d decomposer
	for {
		next := 0
		e = randExpr(edakit.Scaled(ctx, 2, 5), &next)
		d = decomposer{}
		root = d.build(e)
		if k != 3 || d.cancels > 0 {
			break
		}
	}
	nand, inv := countKind(root, 'n'), countKind(root, 'v')
	and, or, not := countOps(e)
	raw := and + 2*or + not
	rules := "Decompose it into a NAND2/INV subject tree by the standard rules: AND(x, y) becomes INV(NAND(x, y)), OR(x, y) becomes NAND(INV x, INV y), NOT x becomes INV x. Build a tree (no sharing) and then delete every pair of back-to-back inverters."
	q := "The function is f = " + e.String() + ". " + rules
	facts := map[string]string{"nand2": edakit.IntStr(nand), "inv": edakit.IntStr(inv), "raw_inv": edakit.IntStr(raw), "cancelled_pairs": edakit.IntStr(d.cancels)}
	why := edakit.Sprintf("There are %d AND, %d OR and %d NOT operators, so the raw decomposition has %d NAND2 nodes and %d + 2(%d) + %d = %d inverters. Cancelling back-to-back pairs removes %d pair(s), leaving %d inverters, so the tree has %d nodes in all.",
		and, or, not, nand, and, or, not, raw, d.cancels, inv, nand+inv)
	switch k {
	case 1:
		return edakit.Numeric(q+" How many nodes (NAND2 plus INV) does the final subject tree have?", edakit.IntStr(nand+inv), why, facts)
	case 2:
		return edakit.Numeric(q+" How many INV nodes does the final subject tree have?", edakit.IntStr(inv), why, facts)
	}
	return edakit.Numeric(q+" How many back-to-back inverter pairs are deleted?", edakit.IntStr(d.cancels), why, facts)
}

// libCell is a library gate with its NAND2/INV pattern.
type libCell struct {
	name  string
	fn    string
	pat   *snode
	area  int
	delay int
}

func cellLibrary() []libCell {
	return []libCell{
		{"INV", "NOT a", pv(px()), 2, 1},
		{"NAND2", "NOT (a AND b)", pn(px(), px()), 3, 2},
		{"AND2", "a AND b", pv(pn(px(), px())), 4, 3},
		{"NOR2", "NOT (a OR b)", pv(pn(pv(px()), pv(px()))), 3, 2},
		{"OR2", "a OR b", pn(pv(px()), pv(px())), 4, 3},
		{"NAND3", "NOT ((a AND b) AND c)", pn(px(), pv(pn(px(), px()))), 4, 3},
		{"AND3", "(a AND b) AND c", pv(pn(px(), pv(pn(px(), px())))), 5, 4},
		{"NAND4", "NOT ((a AND b) AND (c AND d))", pn(pv(pn(px(), px())), pv(pn(px(), px()))), 5, 4},
		{"AOI21", "NOT ((a AND b) OR c)", pv(pn(pn(px(), px()), pv(px()))), 4, 3},
		{"OAI21", "NOT ((a OR b) AND c)", pn(pn(pv(px()), pv(px())), px()), 4, 3},
		{"AOI22", "NOT ((a AND b) OR (c AND d))", pv(pn(pn(px(), px()), pn(px(), px()))), 5, 3},
		{"OAI22", "NOT ((a OR b) AND (c OR d))", pn(pn(pv(px()), pv(px())), pn(pv(px()), pv(px()))), 5, 3},
	}
}

// patternCells are the cells whose function is decomposed for the pattern-graph drill.
var patternCells = []struct {
	name string
	e    func() *expr
}{
	{"AND3", func() *expr { return eand(eand(ev("a"), ev("b")), ev("c")) }},
	{"OR3", func() *expr { return eor(eor(ev("a"), ev("b")), ev("c")) }},
	{"NAND3", func() *expr { return enot(eand(eand(ev("a"), ev("b")), ev("c"))) }},
	{"NOR3", func() *expr { return enot(eor(eor(ev("a"), ev("b")), ev("c"))) }},
	{"AND4", func() *expr { return eand(eand(ev("a"), ev("b")), eand(ev("c"), ev("d"))) }},
	{"OR4", func() *expr { return eor(eor(ev("a"), ev("b")), eor(ev("c"), ev("d"))) }},
	{"NOR4", func() *expr { return enot(eor(eor(ev("a"), ev("b")), eor(ev("c"), ev("d")))) }},
	{"AOI21", func() *expr { return enot(eor(eand(ev("a"), ev("b")), ev("c"))) }},
	{"OAI21", func() *expr { return enot(eand(eor(ev("a"), ev("b")), ev("c"))) }},
	{"AOI22", func() *expr { return enot(eor(eand(ev("a"), ev("b")), eand(ev("c"), ev("d")))) }},
	{"OAI22", func() *expr { return enot(eand(eor(ev("a"), ev("b")), eor(ev("c"), ev("d")))) }},
	{"AOI211", func() *expr { return enot(eor(eor(eand(ev("a"), ev("b")), ev("c")), ev("d"))) }},
	{"OAI211", func() *expr { return enot(eand(eand(eor(ev("a"), ev("b")), ev("c")), ev("d"))) }},
	{"AO21", func() *expr { return eor(eand(ev("a"), ev("b")), ev("c")) }},
	{"OA21", func() *expr { return eand(eor(ev("a"), ev("b")), ev("c")) }},
}

func ev(s string) *expr     { return &expr{op: 'v', v: s} }
func eand(a, b *expr) *expr { return &expr{op: '&', a: a, b: b} }
func eor(a, b *expr) *expr  { return &expr{op: '|', a: a, b: b} }
func enot(a *expr) *expr    { return &expr{op: '!', a: a} }

func patternGraphs(ctx generator.GeneratorContext) generator.Problem {
	c := edakit.Pick(patternCells)
	e := c.e()
	var d decomposer
	root := d.build(e)
	nand, inv := countKind(root, 'n'), countKind(root, 'v')
	and, or, not := countOps(e)
	raw := and + 2*or + not
	q := "The library cell " + c.name + " computes y = " + e.String() + ". Its pattern graph is built by decomposing that function, as written, into NAND2 and INV with AND(x, y) = INV(NAND(x, y)), OR(x, y) = NAND(INV x, INV y), and back-to-back inverters deleted."
	why := edakit.Sprintf("Each of the %d binary operators becomes one NAND2. The raw decomposition has %d inverters; cancelling back-to-back pairs removes %d pair(s), so %d remain: %d nodes in all.", nand, raw, d.cancels, inv, nand+inv)
	switch kp(ctx) {
	case 1:
		return edakit.Choice(q+" How many NAND2 nodes does the pattern graph have?", edakit.IntStr(nand), why,
			edakit.Distinct(nand, 3, nand+1, nand-1, nand+inv)...)
	case 2:
		return edakit.Choice(q+" How many INV nodes does the pattern graph have?", edakit.IntStr(inv), why,
			edakit.Distinct(inv, 3, raw, inv+1, inv+2, inv-1)...)
	}
	return edakit.Choice(q+" How many nodes in total (NAND2 plus INV) does the pattern graph have?", edakit.IntStr(nand+inv), why,
		edakit.Distinct(nand+inv, 3, nand+raw, nand+inv+1, nand+inv-1)...)
}

// randSubjectTree draws a NAND2/INV tree with s NAND2 nodes, distinct inputs at the leaves,
// and no back-to-back inverters.
func randSubjectTree(s int, next *int) *snode {
	var n *snode
	if s == 0 {
		n = pin(piName(*next))
		*next++
	} else {
		left := edakit.Int(0, s-1)
		a := maybeInv(randSubjectTree(left, next))
		b := maybeInv(randSubjectTree(s-1-left, next))
		n = pn(a, b)
	}
	return n
}

func maybeInv(x *snode) *snode {
	if x.kind != 'v' && edakit.Int(1, 10) <= 4 {
		return pv(x)
	}
	return x
}

// listing names the internal nodes n1, n2, ... in post-order and prints one line per node.
func listing(root *snode) ([]string, map[*snode]string) {
	names := map[*snode]string{}
	var lines []string
	var walk func(s *snode)
	walk = func(s *snode) {
		if s.kind == 'i' {
			names[s] = s.name
			return
		}
		for _, c := range s.kids {
			walk(c)
		}
		names[s] = edakit.Sprintf("n%d", len(lines)+1)
		if s.kind == 'n' {
			lines = append(lines, names[s]+" = NAND("+names[s.kids[0]]+", "+names[s.kids[1]]+")")
		} else {
			lines = append(lines, names[s]+" = INV("+names[s.kids[0]]+")")
		}
	}
	walk(root)
	return lines, names
}

func internalNodes(root *snode) []*snode {
	var out []*snode
	var walk func(s *snode)
	walk = func(s *snode) {
		if s.kind == 'i' {
			return
		}
		for _, c := range s.kids {
			walk(c)
		}
		out = append(out, s)
	}
	walk(root)
	return out
}

func describeLib(lib []libCell, withArea, withDelay bool) string {
	var parts []string
	for _, c := range lib {
		s := c.name + " (" + patString(c.pat)
		if withArea {
			s += edakit.Sprintf(", area %d", c.area)
		}
		if withDelay {
			s += edakit.Sprintf(", delay %d", c.delay)
		}
		parts = append(parts, s+")")
	}
	return strings.Join(parts, "; ")
}

func patString(p *snode) string {
	switch p.kind {
	case 'x':
		return "x"
	case 'v':
		return "INV(" + patString(p.kids[0]) + ")"
	}
	return "NAND(" + patString(p.kids[0]) + ", " + patString(p.kids[1]) + ")"
}

var coveringProblemItems = []edakit.Item{
	{KP: 1, Q: "A technology mapper receives a Boolean network that is already optimized. What is its job, stated as a graph problem?",
		A:   "Cover the subject graph with library pattern graphs so every node is produced by some chosen gate, at minimum cost",
		D:   []string{"Find a minimum vertex cover of the netlist hypergraph", "Partition the network into equal-size blocks, one per library cell", "Color the subject graph so that adjacent nodes use different cells"},
		Why: "Mapping chooses a set of library matches whose patterns together cover the subject graph, with each match's inputs available as outputs of other matches or primary inputs."},
	{KP: 1, Q: "Why does a mapper first convert both the network and the library into the same primitive base, such as NAND2 and INV?",
		A:   "So that a library cell can be found in the network by structural comparison of two graphs over the same node types",
		D:   []string{"Because NAND2 and INV are always the fastest cells in the library", "Because the final netlist must contain only NAND2 and INV cells", "Because Boolean functions cannot be compared without a NAND2 basis"},
		Why: "With a common base, a cell is a small pattern graph and mapping becomes pattern matching plus a covering decision."},
	{KP: 1, Q: "In a valid cover, what must be true of the inputs of every chosen match?",
		A:   "Each input is a primary input or the output of another chosen match",
		D:   []string{"Each input is a primary input", "Each input feeds exactly one chosen match", "Each input is an internal node of the same match"},
		Why: "A gate needs its input signals to exist in the mapped netlist, so every match input must be a primary input or the root of another chosen match."},
	{KP: 1, Q: "Which information from a cell library does a pure area-oriented structural mapper actually use?",
		A:   "Each cell's function, as pattern graphs, and its area",
		D:   []string{"Each cell's layout polygons", "Each cell's NLDM delay tables only", "Each cell's transistor sizes and threshold voltages"},
		Why: "For area mapping the mapper needs to know what a cell computes (to match it) and what it costs. Timing tables enter only in delay-oriented mapping."},
	{KP: 2, Q: "A 2-input XOR cell has no tree-shaped NAND2/INV decomposition without repeating an input. What does this mean for a tree-covering mapper?",
		A:   "Its pattern is not a tree, so a plain tree matcher cannot find it and needs special handling",
		D:   []string{"The XOR cell can be matched anywhere because it has only two inputs", "The mapper must decompose the subject graph into XOR gates first", "Tree covering will always prefer the XOR cell"},
		Why: "A tree pattern has distinct leaves. XOR needs each input twice, which makes a DAG pattern that a tree matcher does not handle without extensions."},
	{KP: 2, Q: "The same function, a 4-input NAND, can be decomposed into NAND2/INV as a balanced shape or as a chain. Why do mappers often store both shapes in the pattern library?",
		A:   "A structural matcher only finds a cell whose pattern shape matches the subject graph, so each shape catches different subject structures",
		D:   []string{"The two shapes have different Boolean functions", "Storing both shapes halves the number of matches to evaluate", "One shape is for area mode and the other for power mode"},
		Why: "Structural matching is sensitive to the shape of the decomposition. Several pattern shapes per cell reduce this structural bias."},
	{KP: 2, Q: "What is the main weakness of structural (graph) matching compared with Boolean matching?",
		A:   "It misses matches where the function is the same but the subject graph has a different structure",
		D:   []string{"It cannot handle inverters", "It produces covers that are not logically equivalent to the network", "It is slower than any Boolean matching method"},
		Why: "Structural matching compares shapes. A Boolean matcher compares functions, so it finds a cell regardless of how the subject graph happens to be decomposed."},
	{KP: 2, Q: "A designer removes the AOI21 and OAI21 cells from the library and remaps. What happens to the mapped area in a typical design?",
		A:   "It usually grows, because functions those complex gates covered now need several simpler gates",
		D:   []string{"It is unchanged, because NAND2 and INV can cover anything", "It always shrinks, because simpler gates are smaller", "Mapping fails, because AOI21 is required for completeness"},
		Why: "NAND2 and INV keep the library complete, so mapping succeeds, but complex gates are usually cheaper than the several simple gates that replace them."},
	{KP: 3, Q: "The choice of subject-graph decomposition is made before mapping starts. Why does it matter for the final result?",
		A:   "The mapper only finds covers that exist in that particular subject graph, so a different decomposition can allow a better cover",
		D:   []string{"It does not matter: all decompositions give the same optimal cover", "It only changes the run time, not the result", "It only matters for sequential elements"},
		Why: "Covering is optimal relative to one subject graph. Mappers in ABC therefore use structural choices, several decompositions of the same logic, to widen the search."},
	{KP: 3, Q: "After mapping, an equivalence checker finds that the mapped netlist differs from the input network. Which bug is most likely?",
		A:   "A match was accepted whose pattern does not compute the cell's function, for example a pattern library entry with a wrong inverter",
		D:   []string{"The cover used too many NAND2 gates", "The mapper chose a slower cell than necessary", "The cover duplicated a node with fanout"},
		Why: "Covers built from correct patterns are equivalent by construction. A wrong pattern-to-function pairing breaks that guarantee. Duplication and poor cost choices hurt quality, not correctness."},
	{KP: 3, Q: "Why do mappers work on a technology-independent network rather than mapping RTL operators directly to cells?",
		A:   "Technology-independent optimization is cheaper to do on a uniform representation, and mapping then only has to choose cells for already-good logic",
		D:   []string{"Cells cannot implement RTL operators", "Liberty files describe only AND-inverter graphs", "Mapping RTL directly is NP-hard but mapping AIGs is linear"},
		Why: "The flow separates logic optimization (on AIGs or Boolean networks) from the library-specific decision of which cells implement the logic."},
	{KP: 3, Q: "Two runs map the same AIG with the same library and cost function and get different areas. The only difference is the order in which the AIG nodes were created. What explains it?",
		A:   "The subject graph structure differs, and structural matching and tie-breaking depend on it",
		D:   []string{"Liberty data changed between runs", "Area covering is randomized by definition", "The mapper used a different library each time"},
		Why: "The structure of the subject graph, and ties among equal-cost matches, depend on how the AIG was built. This structural sensitivity is why mappers add structural choices."},
}
