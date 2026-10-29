package mleda

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var learnedRoutingItems = []edakit.Item{
	{KP: 1, Q: "Which part of a negotiated-congestion global router is the safest place for a learned model?",
		A: "Predicting which nets to rip up and the order to reroute them, while the router still finds legal paths", D: []string{"Replacing the maze search that produces the path geometry", "Deciding whether design rules apply", "Skipping the final overflow check"},
		Why: "Ordering and selection are heuristic choices. Path legality and overflow accounting stay with the exact router."},
	{KP: 1, Q: "A learned model predicts a congestion cost map that the router adds to its edge costs. What does a bad prediction cost?",
		A: "Longer runtime or worse wirelength, while routes stay legal because the router enforces rules", D: []string{"Illegal routes that pass signoff", "Shorts between nets that no check can see", "A crash of the DRC engine"},
		Why: "A cost-map hint only biases the search. The router's legality model is unchanged, so errors cost quality, not correctness."},
	{KP: 1, Q: "Why is learning a heuristic for A* in a maze router delicate?",
		A: "A heuristic that overestimates the remaining cost can make A* return a longer path than the optimum", D: []string{"A* cannot use heuristics", "A learned heuristic always speeds A* up with no cost", "A* ignores the heuristic on routing grids"},
		Why: "A* is optimal only with an admissible heuristic. A learned estimate can overestimate, so optimality must be traded or bounded."},
	{KP: 1, Q: "A detailed router uses a model to predict which pin-access patterns will cause DRC violations. How should the prediction be used?",
		A: "To order or prune candidate access points, with the DRC-aware router checking every chosen one", D: []string{"To declare the selected patterns DRC-clean", "To remove the spacing rules from the router", "To skip the final DRC run"},
		Why: "The model narrows the search; the rule-aware engine verifies each decision."},
	{KP: 1, Q: "What representation suits learning to predict where a global router will overflow before routing?",
		A: "Gcell-grid feature maps such as pin density, RUDY and macro blockage, predicted as a 2-D map", D: []string{"The GDS polygons of the final layout", "The Liberty timing tables", "The list of SDC constraints"},
		Why: "Before routing, the information is spatial and on the gcell grid, so image-like encodings fit the task."},
	{KP: 1, Q: "Why are end-to-end learned routers, which output wire geometry directly, rare in production flows?",
		A: "Their outputs would need full DRC and connectivity checks anyway, and repairing violations needs a classical router", D: []string{"Routing has no data to learn from", "Neural networks cannot output coordinates", "Routers are already optimal"},
		Why: "Legality is non-negotiable. A generator of geometry still needs the exact machinery to verify and fix its output."},
	{KP: 3, Q: "A learned net-ordering policy cuts global-routing runtime on its training designs but increases overflow on a new design with many macros. What is the likely cause?",
		A: "Distribution shift: the policy learned orderings that suit designs unlike the new one", D: []string{"The router has a bug that only appears with macros", "Net ordering cannot affect overflow", "The policy was too small to run"},
		Why: "Heuristics learned from one design family can fail on a different floorplan style. Evaluation must include held-out design types."},
	{KP: 3, Q: "An agent uses a learned congestion predictor to decide whether to run the global router at all for a candidate placement. When is that safe?",
		A: "Only for rejecting clearly bad candidates; any candidate that is kept must still be routed", D: []string{"Always, if the predictor's average error is small", "Only on the final signoff run", "Never, because routers cannot be skipped in exploration"},
		Why: "Skipping evaluation on rejections only costs a missed candidate. Accepting without routing lets a wrong prediction through."},
	{KP: 3, Q: "A paper reports that a learned rip-up policy reduces DRC violations by a large factor. What baseline makes the claim meaningful?",
		A: "The same router with its standard rip-up heuristic, tuned, on the same designs and with the same runtime budget", D: []string{"A router with rip-up disabled", "A different router on different designs", "The policy's own early training checkpoint"},
		Why: "The comparison must isolate the policy. A weak or untuned baseline inflates the gain."},
	{KP: 3, Q: "A router's learned cost model makes A* expand far fewer nodes, but wirelength rises by 3 percent. What trade did it make?",
		A: "It used an inflated or inadmissible heuristic, trading path optimality for search speed", D: []string{"It found shorter paths that use more vias", "It reduced the routing grid resolution", "It changed the design rules"},
		Why: "Heuristics that overestimate focus the search but can return suboptimal paths. Weighted A* makes this trade explicit with a bound."},
	{KP: 3, Q: "Where should a learned routing heuristic sit in the router's architecture?",
		A: "Behind a cost or ordering seam in the strategy layer, so the core path search and legality checks are unchanged", D: []string{"Inside the design-rule model", "In the LEF parser", "In place of the connectivity check"},
		Why: "A seam lets the learned heuristic be swapped, tested and rolled back without touching the core."},
	{KP: 3, Q: "What should a differential test of a router with a new learned net-ordering policy check?",
		A: "Routes are legal and connected, and overflow, wirelength and runtime changes against the baseline are within agreed bounds", D: []string{"Only that runtime decreased", "Only that the policy's loss decreased", "That every net took the same route as before"},
		Why: "Correctness must be unchanged; quality and runtime deltas are measured against a baseline before promotion."},
}

// weightedAStarChoice: the suboptimality bound of weighted or epsilon-admissible A*.
func weightedAStarChoice(ctx generator.GeneratorContext) generator.Problem {
	c := float64(edakit.Int(20, 120))
	if edakit.Int(0, 1) == 0 {
		w := edakit.Pick([]float64{1.2, 1.5, 2, 2.5, 3})
		q := edakit.Sprintf("A maze router searches with weighted A*, $f(n) = g(n) + %s\\,h(n)$, where $h$ is the admissible Manhattan-distance heuristic scaled by a learned confidence weight. The optimal route for a net costs %s track units. What is the largest route cost the search is guaranteed not to exceed?", f(w, 1), f(c, 0))
		why := edakit.Sprintf("Weighted A* with an admissible $h$ returns a path of cost at most $w C^* = %s \\times %s = %s$.", f(w, 1), f(c, 0), f(w*c, 1))
		return numChoice(q, w*c, 1, "units", why, c, c/w, c+w, w*c+c)
	}
	e := float64(edakit.Int(2, 15))
	q := edakit.Sprintf("A router uses A* with a nonnegative learned heuristic that may overestimate the true remaining cost by at most %s track units at any node. The optimal route costs %s units. What is the largest route cost the search is guaranteed not to exceed?", f(e, 0), f(c, 0))
	why := edakit.Sprintf("If $h(n) \\le h^*(n) + \\epsilon$ everywhere, A* returns a path of cost at most $C^* + \\epsilon = %s$.", f(c+e, 0))
	return numChoice(q, c+e, 1, "units", why, c, c*(1+e/10), c+2*e, c-e)
}

// generativeDesign: useful designs from a generator (KP 1), cost per accepted design (KP 2),
// and the throughput ratio against random search (KP 3).
func generativeDesign(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 200, 5000)
	v := float64(edakit.Int(30, 95)) / 100
	s := float64(edakit.Int(5, 40)) / 100
	tg := float64(edakit.Int(1, 20)) / 100
	tsim := float64(edakit.Int(2, 30)) / 10
	acc := float64(n) * v * s
	stem := edakit.Sprintf("A generative model proposes %d analog circuit topologies with sizings. A netlist checker finds %s percent of them valid (connected, simulatable). SPICE shows %s percent of the valid ones meet the specification.", n, f(100*v, 0), f(100*s, 0))
	switch kp(ctx) {
	case 1:
		u := float64(edakit.Int(50, 100)) / 100
		use := acc * u
		ans := edakit.Num(use, 1)
		return edakit.Numeric(stem+edakit.Sprintf(" Of the designs that meet the specification, %s percent are not near-duplicates of training examples or of each other. How many novel, spec-meeting designs does the batch yield? Round to 1 decimal place.", f(100*u, 0)),
			ans, edakit.Sprintf("$%d \\times %s \\times %s \\times %s = %s$. Every stage of the funnel is checked by a tool or a comparison, not by the generator.", n, f(v, 2), f(s, 2), f(u, 2), ans),
			map[string]string{"accepted": f(acc, 1)})
	case 2:
		cost := float64(n)*tg + float64(n)*v*tsim
		per := cost / acc
		ans := edakit.Num(per, 3)
		return edakit.Numeric(stem+edakit.Sprintf(" Generating a design costs %s CPU-minutes and simulating a valid one costs %s CPU-minutes; invalid designs are not simulated. What is the cost per spec-meeting design, in CPU-minutes? Round to 3 decimal places.", f(tg, 2), f(tsim, 1)),
			ans, edakit.Sprintf("Cost $= %d \\times %s + %d \\times %s \\times %s = %s$ min. Accepted $= %s$. Per design: %s min.", n, f(tg, 2), n, f(v, 2), f(tsim, 1), f(cost, 2), f(acc, 1), ans),
			map[string]string{"cost": f(cost, 2)})
	default:
		ab := float64(edakit.Int(1, 20)) / 1000
		cost := float64(n)*tg + float64(n)*v*tsim
		rg := acc / cost
		rb := ab / tsim
		ratio := rg / rb
		ans := edakit.Num(ratio, 2)
		return edakit.Numeric(stem+edakit.Sprintf(" Generating costs %s CPU-minutes per design and simulating a valid one %s CPU-minutes. A random-sampling baseline simulates every sample at the same cost, and %s percent of its samples meet the spec. How many times more spec-meeting designs per CPU-minute does the generator deliver? Round to 2 decimal places.", f(tg, 2), f(tsim, 1), f(100*ab, 1)),
			ans, edakit.Sprintf("Generator: $%s / %s = %s$ per min. Baseline: $%s / %s = %s$ per min. Ratio: %s. A fair comparison also counts the generator's training cost.", f(acc, 1), f(cost, 2), f(rg, 5), f(ab, 3), f(tsim, 1), f(rb, 5), ans),
			map[string]string{"rate_gen": f(rg, 5)})
	}
}
