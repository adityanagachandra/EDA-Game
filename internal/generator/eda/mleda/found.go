package mleda

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var whereMLHelpsItems = []edakit.Item{
	{KP: 1, Q: "A global placer spends most of its iterations waiting for a full global route to tell it where congestion is. Where can a learned model help most directly?",
		A: "As a fast congestion estimate that the placer queries between expensive route runs", D: []string{"As a replacement for the legalizer", "As a replacement for the final global router", "As a way to skip design rule checking"},
		Why: "Learning pays where an exact evaluation is slow and a fast, approximate answer is enough to steer a search. The router still produces the real routes."},
	{KP: 1, Q: "Which task is the best fit for a learned predictor in a physical-design flow?",
		A: "Estimating post-route wire delay at placement time, where the true value needs a route that has not happened yet", D: []string{"Computing the exact HPWL of a placed net", "Checking that two cells do not overlap", "Topologically sorting the timing graph"},
		Why: "HPWL, overlap and topological order are cheap and exact. A predictor earns its place only when the true quantity is expensive or not yet available."},
	{KP: 1, Q: "A team must pick one of 200 flow configurations, and each full run takes a day. Why is a learned model attractive here?",
		A: "It can rank configurations from cheap early-stage features, so only a few need full runs", D: []string{"It guarantees the best configuration without any runs", "It makes each full run faster", "It removes the need for timing signoff"},
		Why: "The value is triage: spending the expensive evaluations on the candidates most likely to be good."},
	{KP: 1, Q: "Which property of a problem makes a learned heuristic most likely to help a classical algorithm?",
		A: "The algorithm faces many similar instances and makes repeated choices whose quality is hard to compute locally", D: []string{"The problem has a closed-form solution", "Each instance is solved once and never resembles another", "The exact answer is cheaper than running a network"},
		Why: "Learning amortizes experience across similar instances. With a closed form or a cheap exact answer there is nothing to amortize."},
	{KP: 2, Q: "A neural network predicts that a routed design has zero DRC violations. What may the flow conclude?",
		A: "Nothing for signoff: the design still has to pass the rule deck in a DRC engine", D: []string{"The design can be taped out", "DRC can be skipped on this block", "The router can be turned off for future blocks"},
		Why: "A prediction is a statistical estimate. Manufacturability is decided by checking the actual geometry against the actual rules."},
	{KP: 2, Q: "An LLM rewrites an RTL module to save area and reports that the behavior is unchanged. What establishes that claim?",
		A: "Equivalence checking of the new module against the original", D: []string{"The model's confidence score", "The code compiling without warnings", "A second LLM agreeing"},
		Why: "Functional equivalence is a formal property. A combinational or sequential equivalence checker proves it; a model's statement does not."},
	{KP: 2, Q: "A learned placer outputs coordinates with a few overlapping cells. Which component must still run?",
		A: "A legalizer that snaps cells to rows and sites without overlap", D: []string{"A second learned placer", "A congestion predictor", "Nothing, because overlaps are cosmetic"},
		Why: "Legality is a hard constraint enforced by an exact algorithm. A continuous or learned placement is never legal by construction."},
	{KP: 2, Q: "Why can a timing predictor not replace static timing analysis at signoff?",
		A: "Its error has no guaranteed bound on an unseen design, while STA computes a conservative bound by construction", D: []string{"Predictors cannot output negative numbers", "STA is faster than any predictor", "Timing does not depend on the netlist"},
		Why: "STA's guarantee comes from propagating worst-case delays over every path. A model's error is only measured on its test set."},
	{KP: 3, Q: "A flow uses a model to propose 50 buffer insertions and then runs incremental STA on each. What pattern is this?",
		A: "The model proposes and an exact engine decides", D: []string{"The model replaces the timer", "The timer trains the model online", "End-to-end learned signoff"},
		Why: "The safe pattern keeps learned outputs as proposals, each accepted only after an engine with known guarantees checks it."},
	{KP: 3, Q: "A congestion predictor is wrong on 10 percent of tiles. In which use is that error rate least harmful?",
		A: "Ordering which regions the placer spreads first, with a real global route afterwards", D: []string{"Deciding which tiles skip DRC", "Reporting final overflow to signoff", "Choosing which nets need no routing"},
		Why: "Errors in a heuristic ordering cost some runtime or quality, and the later exact step catches them. Errors in a decision that skips a check reach silicon."},
	{KP: 3, Q: "A paper reports a learned model that matches a classical heuristic's quality but takes three days of GPU training per design. What should a tool designer ask first?",
		A: "Whether the training cost is amortized across designs or paid again for each one", D: []string{"Whether the GPU was the newest model", "Whether the network had enough layers", "Whether the heuristic can be removed from the tool"},
		Why: "A learned method is useful in a tool only if its total cost, training included, beats the alternative over the designs it will see."},
	{KP: 3, Q: "Where should a learned component sit in a placer's architecture so that a bad model cannot produce an illegal result?",
		A: "In the strategy layer, feeding guidance to a core that enforces legality", D: []string{"Inside the legalizer, replacing its overlap check", "In the file writer, after legalization", "In the row and site definitions"},
		Why: "If the learned part only steers decisions and the core enforces invariants, the worst a bad model can do is a poor but legal result."},
}

// labelsFromTools: wall-clock days to generate a labeled dataset (KP 1), labeled tiles per
// dataset (KP 2), and the number of reuses that pay back the labeling cost (KP 3).
func labelsFromTools(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		d := edakit.Scaled(ctx, 10, 60)
		v := edakit.Pick([]int{4, 8, 10, 16, 20})
		h := float64(edakit.Int(2, 12))
		m := edakit.Pick([]int{8, 16, 32, 50, 64})
		days := float64(d*v) * h / float64(m) / 24
		ans := edakit.Num(days, 2)
		return edakit.Numeric(edakit.Sprintf("To label a congestion dataset, a team runs %d designs through the flow with %d parameter variants each. One placement-and-global-route run takes %s hours on one machine, and %d machines run jobs in parallel. How many days of wall-clock time does labeling take? Round to 2 decimal places.", d, v, f(h, 0), m),
			ans, edakit.Sprintf("Runs $= %d \\times %d = %d$. Machine-hours $= %d \\times %s = %s$. Wall-clock $= %s / %d / 24 = %s$ days.", d, v, d*v, d*v, f(h, 0), f(float64(d*v)*h, 0), f(float64(d*v)*h, 0), m, ans),
			map[string]string{"runs": edakit.IntStr(d * v)})
	case 2:
		w := edakit.Pick([]int{200, 300, 400, 500, 600})
		hh := edakit.Pick([]int{200, 300, 400, 500})
		g := edakit.Pick([]int{2, 4, 5, 10})
		runs := edakit.Scaled(ctx, 20, 200)
		tiles := (w / g) * (hh / g)
		total := tiles * runs
		return edakit.Numeric(edakit.Sprintf("A congestion dataset labels every gcell of each routed run. A die is %d by %d um and a gcell is %d by %d um. The dataset holds %d runs. How many labeled gcells does it contain?", w, hh, g, g, runs),
			edakit.IntStr(total), edakit.Sprintf("Per run: $(%d/%d)(%d/%d) = %d \\times %d = %d$ gcells. Times %d runs: %d. These labels are far from independent: gcells of one run share a design and a placement.", w, g, hh, g, w/g, hh/g, tiles, runs, total),
			map[string]string{"tiles_per_run": edakit.IntStr(tiles)})
	default:
		label := float64(edakit.Int(200, 2000))
		train := float64(edakit.Int(10, 200))
		tool := float64(edakit.Int(2, 20))
		infer := float64(edakit.Int(1, 10)) / 100
		n := math.Floor((label+train)/(tool-infer)) + 1
		return edakit.Numeric(edakit.Sprintf("Building a predictor costs %s machine-hours of tool runs for labels plus %s hours of training. Afterwards each query takes %s hours instead of the %s-hour tool run it replaces in an exploration loop. What is the smallest number of queries for which the predictor saves machine-hours overall? Give an integer.", f(label, 0), f(train, 0), f(infer, 2), f(tool, 0)),
			edakit.IntStr(int(n)), edakit.Sprintf("Savings per query $= %s - %s = %s$ hours. Fixed cost $= %s + %s = %s$ hours. The smallest $N$ with $N(%s) > %s$ is $\\lfloor %s / %s \\rfloor + 1 = %d$. This assumes no query still needs a confirming tool run.", f(tool, 0), f(infer, 2), f(tool-infer, 2), f(label, 0), f(train, 0), f(label+train, 0), f(tool-infer, 2), f(label+train, 0), f(label+train, 0), f(tool-infer, 2), int(n)),
			map[string]string{"fixed": f(label+train, 0)})
	}
}
