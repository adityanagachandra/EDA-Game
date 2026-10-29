package techmap

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// bindings returns every way pattern p matches rooted at subject node s, each as the list of
// subject nodes bound to the pattern's free leaves. NAND2 inputs are commutative.
func bindings(p, s *snode) [][]*snode {
	switch p.kind {
	case 'x':
		return [][]*snode{{s}}
	case 'v':
		if s.kind != 'v' {
			return nil
		}
		return bindings(p.kids[0], s.kids[0])
	}
	if s.kind != 'n' {
		return nil
	}
	var out [][]*snode
	for _, ord := range [][2]int{{0, 1}, {1, 0}} {
		l := bindings(p.kids[0], s.kids[ord[0]])
		r := bindings(p.kids[1], s.kids[ord[1]])
		for _, a := range l {
			for _, b := range r {
				out = append(out, append(append([]*snode{}, a...), b...))
			}
		}
	}
	return out
}

// randLibrary returns INV and NAND2 plus extra cells, with small random jitter on area and delay.
func randLibrary(extra int) []libCell {
	all := cellLibrary()
	lib := []libCell{all[0], all[1]}
	rest := edakit.Shuffle(all[2:])
	if extra > len(rest) {
		extra = len(rest)
	}
	lib = append(lib, rest[:extra]...)
	for i := range lib {
		lib[i].area += edakit.Int(0, 1)
		lib[i].delay += edakit.Int(0, 1)
	}
	return lib
}

type coverResult struct {
	cost  map[*snode]int
	arr   map[*snode]int
	area  map[*snode]int // area of the delay-optimal cover (ties broken by area)
	pick  map[*snode]string
	dpick map[*snode]string
}

// coverTree runs the area DP and the delay DP (ties by area) over a subject tree.
func coverTree(root *snode, lib []libCell, piArr map[string]int) coverResult {
	r := coverResult{cost: map[*snode]int{}, arr: map[*snode]int{}, area: map[*snode]int{}, pick: map[*snode]string{}, dpick: map[*snode]string{}}
	for _, s := range internalNodes(root) {
		best, bestArr, bestArea := 1<<30, 1<<30, 1<<30
		for _, c := range lib {
			for _, b := range bindings(c.pat, s) {
				cost, mx, ar := c.area, 0, c.area
				for _, l := range b {
					if l.kind == 'i' {
						if piArr[l.name] > mx {
							mx = piArr[l.name]
						}
						continue
					}
					cost += r.cost[l]
					ar += r.area[l]
					if r.arr[l] > mx {
						mx = r.arr[l]
					}
				}
				if cost < best {
					best = cost
					r.pick[s] = c.name
				}
				a := mx + c.delay
				if a < bestArr || (a == bestArr && ar < bestArea) {
					bestArr, bestArea = a, ar
					r.dpick[s] = c.name
				}
			}
		}
		r.cost[s], r.arr[s], r.area[s] = best, bestArr, bestArea
	}
	return r
}

// greedyArea covers top-down, taking at each node the match that covers the most subject nodes
// (ties: smaller area).
func greedyArea(s *snode, lib []libCell) int {
	if s.kind == 'i' {
		return 0
	}
	bestCover, bestArea := -1, 0
	var bestLeaves []*snode
	for _, c := range lib {
		for _, b := range bindings(c.pat, s) {
			cov := countKind(c.pat, 'n') + countKind(c.pat, 'v')
			if cov > bestCover || (cov == bestCover && c.area < bestArea) {
				bestCover, bestArea, bestLeaves = cov, c.area, b
			}
		}
	}
	total := bestArea
	for _, l := range bestLeaves {
		total += greedyArea(l, lib)
	}
	return total
}

func treeStem(root *snode) string {
	lines, names := listing(root)
	return "The subject tree is " + strings.Join(lines, "; ") + ". Its root is " + names[root] + "."
}

func treeMatching(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	next := 0
	root := maybeInv(randSubjectTree(edakit.Scaled(ctx, 2, 5), &next))
	lib := randLibrary(edakit.Int(5, 8))
	_, names := listing(root)
	q := treeStem(root) + " The library patterns are " + describeLib(lib, false, false) + ". A pattern leaf x may bind to any subject node, and NAND inputs may be matched in either order."
	switch k {
	case 1, 2:
		target := root
		if k == 2 {
			in := internalNodes(root)
			target = in[edakit.Int(0, len(in)-1)]
		}
		var hits []string
		for _, c := range lib {
			if len(bindings(c.pat, target)) > 0 {
				hits = append(hits, c.name)
			}
		}
		return edakit.Numeric(q+" How many of these library cells have a match rooted at node "+names[target]+"?",
			edakit.IntStr(len(hits)),
			edakit.Sprintf("Rooted at %s the cells that match are %s: %d in all.", names[target], strings.Join(hits, ", "), len(hits)),
			map[string]string{"matches": edakit.IntStr(len(hits))})
	}
	total := 0
	for _, s := range internalNodes(root) {
		for _, c := range lib {
			if len(bindings(c.pat, s)) > 0 {
				total++
			}
		}
	}
	return edakit.Numeric(q+" Counting each (node, cell) pair once, how many matches exist over all nodes of the tree?",
		edakit.IntStr(total),
		edakit.Sprintf("Every node matches INV or NAND2 trivially by its own type, and larger cells match where the structure below allows. Summing over the nodes gives %d (node, cell) matches.", total),
		map[string]string{"matches": edakit.IntStr(total)})
}

func dpCovering(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	size := edakit.Scaled(ctx, 3, 6)
	if k == 1 {
		size = edakit.Int(2, 4)
	}
	for {
		next := 0
		root := maybeInv(randSubjectTree(size, &next))
		lib := randLibrary(edakit.Int(3, 6))
		res := coverTree(root, lib, map[string]int{})
		q := treeStem(root) + " The library is " + describeLib(lib, true, false) + "."
		if k == 3 {
			g := greedyArea(root, lib)
			if g == res.cost[root] && edakit.Int(1, 3) > 1 {
				continue
			}
			return edakit.Numeric(q+" A greedy mapper works top-down from the root and at each node takes the match that covers the most subject nodes (ties: smaller area), then recurses on that match's inputs. What total area does the greedy cover have?",
				edakit.IntStr(g),
				edakit.Sprintf("Greedy commits to the biggest pattern without asking what its inputs will cost, and reaches area %d. The dynamic program, which minimizes cell area plus the optimal cost of every input, reaches %d.", g, res.cost[root]),
				map[string]string{"greedy_area": edakit.IntStr(g), "optimal_area": edakit.IntStr(res.cost[root])})
		}
		var steps []string
		for _, s := range internalNodes(root) {
			_, names := listing(root)
			steps = append(steps, edakit.Sprintf("%s: %d (%s)", names[s], res.cost[s], res.pick[s]))
		}
		return edakit.Numeric(q+" Primary inputs cost nothing. What is the minimum total area of a cover of the tree?",
			edakit.IntStr(res.cost[root]),
			"Process the nodes bottom-up. The cost of a node is the minimum over its matches of the cell area plus the cost of each match input. Best costs: "+strings.Join(steps, ", ")+".",
			map[string]string{"min_area": edakit.IntStr(res.cost[root])})
	}
}

func delayCovering(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	size := edakit.Scaled(ctx, 3, 6)
	if k == 1 {
		size = edakit.Int(2, 4)
	}
	next := 0
	root := maybeInv(randSubjectTree(size, &next))
	lib := randLibrary(edakit.Int(3, 6))
	arr := map[string]int{}
	var ins []string
	for i := 0; i < next; i++ {
		a := 0
		if k != 1 {
			a = edakit.Int(0, 4)
		}
		arr[piName(i)] = a
		ins = append(ins, edakit.Sprintf("%s = %d", piName(i), a))
	}
	res := coverTree(root, lib, arr)
	_, names := listing(root)
	q := treeStem(root) + " The library, with load-independent delays in ps, is " + describeLib(lib, k == 3, true) + ". Input arrival times in ps are " + strings.Join(ins, ", ") + "."
	var steps []string
	for _, s := range internalNodes(root) {
		steps = append(steps, edakit.Sprintf("%s: %d (%s)", names[s], res.arr[s], res.dpick[s]))
	}
	if k == 3 {
		return edakit.Numeric(q+" Run the delay-optimal tree covering, keeping at each node the match with the smallest arrival time and breaking ties by the smaller total area. What is the total area of the resulting cover?",
			edakit.IntStr(res.area[root]),
			edakit.Sprintf("The delay DP gives root arrival %d ps with area %d, while the area-optimal cover has area %d. Best arrivals: %s.", res.arr[root], res.area[root], res.cost[root], strings.Join(steps, ", ")),
			map[string]string{"arrival": edakit.IntStr(res.arr[root]), "delay_cover_area": edakit.IntStr(res.area[root]), "min_area": edakit.IntStr(res.cost[root])})
	}
	return edakit.Numeric(q+" What is the earliest possible arrival time at the root over all covers, in ps?",
		edakit.IntStr(res.arr[root]),
		"Bottom-up, the arrival of a node is the minimum over its matches of the cell delay plus the latest arrival among the match inputs. Best arrivals: "+strings.Join(steps, ", ")+".",
		map[string]string{"arrival": edakit.IntStr(res.arr[root])})
}

// randDAG draws a NAND2/INV network over npi inputs with ng gates. fanins[i] lists signal
// indices (inputs 0..npi-1, gates npi..).
func randDAG(npi, ng int) [][]int {
	fanins := make([][]int, npi+ng)
	for g := npi; g < npi+ng; g++ {
		pick := func(avoid int) int {
			for {
				var x int
				if g > npi && edakit.Int(1, 10) <= 6 {
					lo := g - 3
					if lo < npi {
						lo = npi
					}
					x = edakit.Int(lo, g-1)
				} else {
					x = edakit.Int(0, g-1)
				}
				if x != avoid {
					return x
				}
			}
		}
		if edakit.Int(1, 10) <= 3 {
			fanins[g] = []int{pick(-1)}
		} else {
			a := pick(-1)
			fanins[g] = []int{a, pick(a)}
		}
	}
	return fanins
}

func sigName(i, npi int) string {
	if i < npi {
		return piName(i)
	}
	return edakit.Sprintf("g%d", i-npi+1)
}

func treePartition(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	npi, ng := edakit.Int(3, 4), edakit.Scaled(ctx, 6, 10)
	fanins := randDAG(npi, ng)
	fo := make([]int, npi+ng)
	for g := npi; g < npi+ng; g++ {
		for _, f := range fanins[g] {
			fo[f]++
		}
	}
	isOut := map[int]bool{}
	for g := npi; g < npi+ng; g++ {
		if fo[g] == 0 {
			isOut[g] = true
		}
	}
	// Occasionally an internal signal is also a primary output.
	if edakit.Int(0, 1) == 1 {
		g := edakit.Int(npi, npi+ng-1)
		isOut[g] = true
	}
	root := map[int]bool{}
	for g := npi; g < npi+ng; g++ {
		if isOut[g] || fo[g] >= 2 {
			root[g] = true
		}
	}
	var lines, outs []string
	for g := npi; g < npi+ng; g++ {
		var args []string
		for _, f := range fanins[g] {
			args = append(args, sigName(f, npi))
		}
		op := "NAND"
		if len(args) == 1 {
			op = "INV"
		}
		lines = append(lines, sigName(g, npi)+" = "+op+"("+strings.Join(args, ", ")+")")
		if isOut[g] {
			outs = append(outs, sigName(g, npi))
		}
	}
	// tree sizes
	var size func(g int) int
	size = func(g int) int {
		n := 1
		for _, f := range fanins[g] {
			if f >= npi && !root[f] {
				n += size(f)
			}
		}
		return n
	}
	roots := sortedKeys(root)
	var rnames []string
	largest, largestName := 0, ""
	for _, r := range roots {
		rnames = append(rnames, sigName(r, npi))
		if s := size(r); s > largest {
			largest, largestName = s, sigName(r, npi)
		}
	}
	q := "A NAND2/INV subject graph over inputs " + strings.Join(names(npi), ", ") + " is " + strings.Join(lines, "; ") + ". The primary outputs are " + strings.Join(outs, ", ") + ". A tree mapper partitions it into trees by cutting at every gate that is a primary output or drives two or more gate inputs."
	facts := map[string]string{"trees": edakit.IntStr(len(roots)), "largest": edakit.IntStr(largest)}
	why := "The tree roots are " + strings.Join(rnames, ", ") + edakit.Sprintf(": the primary outputs plus every gate with fanout two or more. ")
	switch k {
	case 1:
		return edakit.Numeric(q+" How many trees does the partition produce?", edakit.IntStr(len(roots)), why+edakit.Sprintf("That is %d trees.", len(roots)), facts)
	case 2:
		var big []int
		for _, x := range roots {
			if size(x) >= 2 {
				big = append(big, x)
			}
		}
		if len(big) == 0 {
			big = roots
		}
		r := edakit.Pick(big)
		s := size(r)
		return edakit.Numeric(q+" How many gates belong to the tree rooted at "+sigName(r, npi)+"?", edakit.IntStr(s),
			why+edakit.Sprintf("Walking down from %s and stopping at inputs and at other roots collects %d gate(s).", sigName(r, npi), s), facts)
	}
	return edakit.Numeric(q+" How many gates does the largest tree contain?", edakit.IntStr(largest),
		why+edakit.Sprintf("Walking down from each root and stopping at inputs and other roots, the largest tree is rooted at %s with %d gates.", largestName, largest), facts)
}

func names(npi int) []string {
	var out []string
	for i := 0; i < npi; i++ {
		out = append(out, piName(i))
	}
	return out
}

func dagCovering(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	m := edakit.Scaled(ctx, 2, 4)
	ai, an2, an3 := edakit.Int(1, 3), edakit.Int(2, 4), 0
	an3 = an2 + edakit.Int(0, 3)
	di, dn2 := edakit.Int(1, 3), edakit.Int(2, 4)
	dn3 := dn2 + edakit.Int(1, 3)
	var cons []string
	for i := 1; i <= m; i++ {
		cons = append(cons, edakit.Sprintf("y%d = NAND(INV(s), c%d)", i, i))
	}
	q := edakit.Sprintf("In a subject graph, s = NAND(a, b) feeds %d outputs: %s. The library has INV (pattern INV(x), area %d, delay %d ps), NAND2 (pattern NAND(x, x), area %d, delay %d ps) and NAND3 (pattern NAND(x, INV(NAND(x, x))), area %d, delay %d ps).",
		m, strings.Join(cons, ", "), ai, di, an2, dn2, an3, dn3)
	tree := an2 + m*(ai+an2)
	dup := m * an3
	best := tree
	if dup < best {
		best = dup
	}
	facts := map[string]string{"tree_area": edakit.IntStr(tree), "dup_area": edakit.IntStr(dup), "best_area": edakit.IntStr(best)}
	switch k {
	case 1:
		return edakit.Numeric(q+" A tree mapper cuts the graph at the multi-fanout node s. What is the total area of its best cover?", edakit.IntStr(tree),
			edakit.Sprintf("With s a tree root it must be a gate output: NAND2, area %d. Each output tree is NAND(INV(s), c), which with s as a leaf only matches INV plus NAND2, area %d + %d. Total %d + %d(%d) = %d.", an2, ai, an2, an2, m, ai+an2, tree), facts)
	case 2:
		return edakit.Numeric(q+" A DAG mapper may duplicate s into the outputs' matches. What is the minimum total area over all covers it can choose from?", edakit.IntStr(best),
			edakit.Sprintf("Duplicating s lets each output match NAND3 over a, b and its c input, for %d x %d = %d. Keeping s as a shared NAND2 costs %d. Mixing does not help, because keeping s for even one output pays its full NAND2 area. The minimum is %d.", m, an3, dup, tree, best), facts)
	}
	dTree := dn2 + di + dn2
	return edakit.Numeric(q+" All inputs arrive at 0 ps. By how many ps does duplicating s into NAND3 matches reduce the output arrival time compared with the tree cover?",
		edakit.IntStr(dTree-dn3),
		edakit.Sprintf("The tree cover's path is NAND2, INV, NAND2: %d + %d + %d = %d ps. The duplicated cover is one NAND3 of %d ps. The reduction is %d ps.", dn2, di, dn2, dTree, dn3, dTree-dn3),
		map[string]string{"tree_delay": edakit.IntStr(dTree), "dup_delay": edakit.IntStr(dn3)})
}

var treeLimitsItems = []edakit.Item{
	{KP: 1, Q: "In what sense is DAGON-style tree covering optimal?",
		A:   "It finds a minimum-cost cover of each tree of the partitioned subject graph, for the given patterns",
		D:   []string{"It finds a minimum-area netlist for the Boolean function", "It finds a minimum-cost cover of the whole DAG", "It finds the minimum-delay netlist over all decompositions"},
		Why: "The dynamic program is exact per tree and relative to a fixed subject graph and pattern set. Cutting at multi-fanout points and fixing the decomposition both lose global optimality."},
	{KP: 1, Q: "Why is the dynamic program on a tree exact?",
		A:   "In a tree the subtrees under a match's inputs are disjoint, so their optimal costs can be chosen independently and added",
		D:   []string{"Because every tree has a unique cover", "Because trees have no inverters", "Because the library contains only tree patterns of depth one"},
		Why: "Optimal substructure holds when the subproblems do not share nodes. In a tree they never do."},
	{KP: 1, Q: "What is the time complexity of tree covering in the size of the subject tree, for a fixed library?",
		A:   "Linear, because each node tries a bounded number of patterns, each matched in bounded time",
		D:   []string{"Exponential in the number of nodes", "Quadratic, because every pair of nodes is compared", "Proportional to the number of possible covers"},
		Why: "With a fixed library, matching at a node is constant work, so the whole bottom-up pass is linear in the tree size."},
	{KP: 1, Q: "A tree mapper reports an optimal area of 120 for each of three trees and a total of 360. Is 360 a lower bound on any mapping of the network?",
		A:   "No: a cover that crosses the tree boundaries, for example by duplicating logic, may be cheaper",
		D:   []string{"Yes, tree covering is globally optimal", "Yes, because each tree's cost cannot be reduced", "No, but only if the library changes"},
		Why: "The partition forces each multi-fanout node to be a gate output. Covers that ignore that constraint are outside the tree mapper's search space."},
	{KP: 2, Q: "What is known about minimum-area covering of a general subject DAG with a library?",
		A:   "It is NP-hard",
		D:   []string{"It is solvable in linear time by the same DP", "It is polynomial when the library has at least three cells", "It is undecidable"},
		Why: "Area-optimal DAG covering is NP-hard in general, which is why practical mappers use tree partitioning or cut-based heuristics with area recovery."},
	{KP: 2, Q: "Under a load-independent delay model, what is known about delay-optimal mapping of a DAG?",
		A:   "It can be solved optimally in polynomial time, because a node's best arrival depends only on its transitive fanin",
		D:   []string{"It is NP-hard like area mapping", "It requires the network to be a tree", "It is only solvable for LUT targets"},
		Why: "Arrival times propagate forward, and the best arrival of a node does not depend on how other cones are covered, so a labeling pass gives the optimum (FlowMap is the LUT instance)."},
	{KP: 2, Q: "Why does area-optimal DAG covering resist the dynamic program that works on trees?",
		A:   "Cones of different outputs overlap, so the cost of a shared node cannot be charged independently to each user",
		D:   []string{"DAGs have no topological order", "Patterns cannot match across a fanout point", "The number of matches per node becomes infinite"},
		Why: "In a DAG a shared node's implementation serves several consumers. Charging it fully to each overcounts, and charging nothing undercounts: the subproblems are not independent."},
	{KP: 2, Q: "A mapper duplicates a 2-input gate into four fanout cones to save one level of logic on each. What does it trade away?",
		A:   "Area and possibly power, since the shared logic now exists several times",
		D:   []string{"Logical correctness", "Testability of the shared gate", "Nothing, duplication is free"},
		Why: "Duplication buys delay by recomputing logic in several places. The cost is extra area and switching capacitance."},
	{KP: 3, Q: "Some tree mappers insert an INV-INV pair on the edges around multi-fanout nodes and add a zero-cost pattern INV(INV(x)) to the library. What does this achieve?",
		A:   "Each consumer can pick the signal polarity it needs, absorbing an inverter into its match where that is cheaper",
		D:   []string{"It doubles the speed of the circuit", "It removes the need for tree partitioning", "It makes the subject graph acyclic"},
		Why: "With the pairs inserted, a match may absorb one inverter of a pair, so the cover chooses polarities freely. A pair that no match uses is covered by the zero-cost INV(INV(x)) pattern and disappears."},
	{KP: 3, Q: "A network has very few multi-fanout nodes. How close is tree covering likely to be to the best DAG cover?",
		A:   "Close, because the partition cuts few edges and most of the structure is already tree-like",
		D:   []string{"Far, because tree covering ignores area", "Identical in every case", "Far, because single-fanout logic cannot be mapped by trees"},
		Why: "The loss from tree covering comes from the forced cuts at fanout points. With few of them the tree optimum is near the DAG optimum."},
	{KP: 3, Q: "Why did cut-based mappers largely replace explicit tree covering in tools such as ABC?",
		A:   "Cuts let a match cross fanout points, so the mapper works on the DAG directly and then recovers area with heuristics",
		D:   []string{"Cut enumeration is exact for area on DAGs", "Tree covering cannot handle inverters", "Cuts remove the need for a library"},
		Why: "Cut-based mapping evaluates matches over k-input cones anywhere in the DAG, which subsumes duplication. Area is then recovered heuristically, with area flow and exact local area."},
	{KP: 3, Q: "A delay-optimal DAG mapping duplicates many nodes. Which step usually follows in a practical mapper?",
		A:   "Area recovery that re-maps nodes off the critical path with cheaper matches while keeping the required times",
		D:   []string{"Re-running tree covering from scratch", "Removing all duplicated nodes regardless of delay", "Upsizing every gate on the critical path"},
		Why: "Delay-optimal labeling is computed first, then non-critical nodes, which have slack, are re-covered to recover area without increasing the depth."},
}
