package mleda

import (
	"math"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// mdpFormulation: discounted return of a reward sequence (KP 1), the return of a sparse
// terminal reward (KP 2), and a one-step TD error (KP 3).
func mdpFormulation(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 3, 6)
		g := edakit.Pick([]float64{0.5, 0.8, 0.9, 0.95, 0.99})
		r := make([]float64, n)
		ret := 0.0
		var terms []string
		for t := range r {
			r[t] = float64(edakit.Int(-5, 5))
			ret += math.Pow(g, float64(t)) * r[t]
			terms = append(terms, edakit.Sprintf("%s^{%d}(%s)", f(g, 2), t, f(r[t], 0)))
		}
		ans := edakit.Num(ret, 4)
		return edakit.Numeric(edakit.Sprintf("An RL agent tuning a synthesis script receives rewards %s at steps 0 to %d (the reduction in node count, in hundreds, after each command). With discount factor $\\gamma = %s$, what is the return $G_0 = \\sum_t \\gamma^t r_t$? Round to 4 decimal places.", floats(r, 0), n-1, f(g, 2)),
			ans, edakit.Sprintf("$G_0 = %s = %s$.", strings.Join(terms, " + "), ans), map[string]string{"gamma": f(g, 2)})
	case 2:
		m := edakit.Scaled(ctx, 10, 200)
		g := edakit.Pick([]float64{0.9, 0.95, 0.99, 0.999, 1})
		c := float64(edakit.Int(50, 150)) / 100
		ret := -c * math.Pow(g, float64(m-1))
		ans := edakit.Num(ret, 4)
		return edakit.Numeric(edakit.Sprintf("A macro-placement episode places %d macros, one per step. Every step has reward 0 except the last, where the reward is the negative proxy cost $-%s$ of the finished placement. With discount $\\gamma = %s$, what is the return from step 0? Round to 4 decimal places.", m, f(c, 2), f(g, 3)),
			ans, edakit.Sprintf("Only the final reward counts, discounted by $\\gamma^{%d}$: $G_0 = %s^{%d} \\cdot (-%s) = %s$. With sparse terminal rewards, a discount below 1 shrinks the learning signal for early macros; placement formulations therefore often use $\\gamma = 1$.", m-1, f(g, 3), m-1, f(c, 2), ans),
			map[string]string{"weight": f(math.Pow(g, float64(m-1)), 6)})
	default:
		r := float64(edakit.Int(-10, 10)) / 2
		g := edakit.Pick([]float64{0.9, 0.95, 0.99})
		v := float64(edakit.Int(-20, 20)) / 2
		vn := float64(edakit.Int(-20, 20)) / 2
		td := r + g*vn - v
		ans := edakit.Num(td, 3)
		return edakit.Numeric(edakit.Sprintf("A critic estimates the value of the current partial placement as $V(s) = %s$. The agent places the next macro, receives an intermediate reward $r = %s$ from a density penalty, and the critic values the new state at $V(s') = %s$. With $\\gamma = %s$, what is the TD error $\\delta = r + \\gamma V(s') - V(s)$? Round to 3 decimal places.", f(v, 1), f(r, 1), f(vn, 1), f(g, 2)),
			ans, edakit.Sprintf("$\\delta = %s + %s \\cdot %s - %s = %s$. A positive value means the move turned out better than the critic expected.", f(r, 1), f(g, 2), pn(vn, 1), pn(v, 1), ans),
			map[string]string{"target": f(r+g*vn, 3)})
	}
}

var rlPlacementItems = []edakit.Item{
	{KP: 1, Q: "In the sequential-decision formulation of macro placement, what is one action?",
		A: "Choosing the grid cell where the next macro is placed", D: []string{"Moving every standard cell by one step", "Choosing the next routing layer", "Selecting a clock period"},
		Why: "The agent places macros one at a time onto a coarse grid; the state is the partial placement and the netlist."},
	{KP: 1, Q: "Why does the macro-placement policy mask some actions at each step?",
		A: "Grid cells where the macro would overlap placed macros or leave the canvas are illegal", D: []string{"To speed up the GNN encoder", "Because every action has the same reward", "To force macros onto the die edge"},
		Why: "Masking enforces hard constraints so that the policy can only choose legal positions."},
	{KP: 1, Q: "In the published RL macro placer, how are the standard cells placed?",
		A: "By a separate conventional placer, such as a force-directed method, after or between macro decisions", D: []string{"By the same policy, one cell per step", "They are not placed at all", "By a second RL agent for each cell"},
		Why: "Millions of cells would make the episode far too long. The agent handles macros; a classical method places the cells for evaluating the reward."},
	{KP: 1, Q: "What role does the graph neural network play in the RL placement agent?",
		A: "It encodes the netlist and the partial placement into features for the policy and value networks", D: []string{"It legalizes the final placement", "It computes the exact routed wirelength", "It replaces the reward"},
		Why: "The GNN gives a representation of the netlist that the policy can condition on, and that can transfer across netlists in pretraining."},
	{KP: 1, Q: "Why is the reward computed from a proxy cost rather than from a full route and signoff timing?",
		A: "A full route per episode would make training far too slow, so a fast estimate stands in", D: []string{"Routing results are not numbers", "Proxy costs are exactly equal to routed results", "Signoff timing is not affected by placement"},
		Why: "RL needs many episodes. A cheap proxy of wirelength, congestion and density keeps each episode affordable."},
	{KP: 1, Q: "What makes macro placement a hard RL problem compared with board games?",
		A: "The reward arrives only after all macros are placed and is a proxy for the real objective", D: []string{"There are fewer actions per step than in a board game", "The state is fully random", "The episode length is always one"},
		Why: "Sparse, delayed reward and a proxy that may not match the signoff objective make credit assignment and evaluation hard."},
	{KP: 3, Q: "An RL placer achieves a lower proxy cost than a baseline, but after routing the baseline has better timing. What does this reveal?",
		A: "The proxy does not capture what signoff measures, so proxy wins need routed confirmation", D: []string{"The router must be broken", "The baseline was trained on the test design", "Timing never depends on macro placement"},
		Why: "Optimizing a proxy hard tends to exploit its blind spots. Only post-route metrics show whether the placement is better."},
	{KP: 3, Q: "Which comparison isolates the effect of learning in an RL macro placer?",
		A: "The same proxy cost and evaluation flow optimized by a strong non-learned method such as simulated annealing with a matched compute budget", D: []string{"A comparison against random macro positions only", "A comparison of training curves across seeds", "A comparison of GPU types"},
		Why: "If a classical optimizer reaches the same proxy cost with the same budget, the gain is from the formulation, not from learning."},
	{KP: 3, Q: "A pretrained RL policy is fine-tuned on a new block and wins. What must the report state for the result to be interpretable?",
		A: "The pretraining data, the fine-tuning compute, and the compute given to each baseline", D: []string{"Only the final proxy cost", "Only the GPU vendor", "Only the number of macros"},
		Why: "Without the compute and data accounting, a win may reflect more search effort rather than better learning."},
	{KP: 3, Q: "The policy starts from an initial placement of standard-cell clusters produced by a commercial tool. Why does that matter for evaluation?",
		A: "The starting point carries information from another optimizer, so results may depend on it", D: []string{"It does not, since the agent ignores the input", "It makes the result illegal", "It removes the need for a reward"},
		Why: "An informative initial placement can explain part of the quality. Ablating it is needed to credit the learned policy."},
	{KP: 3, Q: "Which evaluation practice would most help a reader trust an RL placement result?",
		A: "Results on open designs with open flows, several seeds, and post-route PPA alongside the proxy cost", D: []string{"A single best run on a proprietary design", "Only the proxy cost on the training designs", "Only wall-clock time to the first placement"},
		Why: "Reproducible inputs, variance across seeds and real downstream metrics make the comparison checkable."},
	{KP: 3, Q: "Where does an RL macro placer sit in the four-layer engine model?",
		A: "In the strategy layer, proposing macro positions that the core legalizes and the flow evaluates", D: []string{"In the immutable core, replacing legalization", "In the observation hooks only", "In the file format layer"},
		Why: "A learned policy decides how to place; legality and evaluation stay in the core and downstream engines."},
}

// proxyCostChoice: the weighted proxy cost of a macro placement.
func proxyCostChoice(ctx generator.GeneratorContext) generator.Problem {
	wl := float64(edakit.Int(50, 150)) / 100
	c := float64(edakit.Int(80, 200)) / 100
	d := float64(edakit.Int(50, 150)) / 100
	lc := edakit.Pick([]float64{0.5, 1, 2})
	ld := edakit.Pick([]float64{0.5, 1, 2})
	cost := wl + lc*c + ld*d
	q := edakit.Sprintf("An RL macro placer's reward is the negative of a proxy cost $\\mathrm{WL} + \\lambda_c C + \\lambda_d D$, with normalized wirelength WL, congestion C and density D. A finished placement has $\\mathrm{WL} = %s$, $C = %s$ and $D = %s$, with $\\lambda_c = %s$ and $\\lambda_d = %s$. What proxy cost does it receive?", f(wl, 2), f(c, 2), f(d, 2), f(lc, 1), f(ld, 1))
	why := edakit.Sprintf("$%s + %s \\cdot %s + %s \\cdot %s = %s$. The weights decide which failure the agent avoids most, and none of the terms is routed timing.", f(wl, 2), f(lc, 1), f(c, 2), f(ld, 1), f(d, 2), f(cost, 4))
	return numChoice(q, cost, 4, "", why, wl+c+d, wl+ld*c+lc*d, wl, lc*c+ld*d)
}

var placementDebateItems = []edakit.Item{
	{KP: 1, Q: "What did the 2021 Nature paper from Google on chip floorplanning claim, in outline?",
		A: "That an RL agent with a graph-based encoder, pretrained across blocks, could place macros with quality comparable to or better than human experts in far less time", D: []string{"That RL could replace the entire place-and-route flow including signoff", "That RL found provably optimal macro placements", "That RL removed the need for standard-cell placement"},
		Why: "The claim concerned macro placement quality and time against human and baseline methods, judged by proxy and downstream metrics on Google's blocks."},
	{KP: 1, Q: "What did Google release that later evaluations relied on?",
		A: "An open-source implementation of the RL macro placer, known as Circuit Training", D: []string{"The production TPU netlists used in the paper", "A signoff-quality commercial router", "A new foundry PDK"},
		Why: "Independent groups ran the released code on open designs because the original blocks were proprietary."},
	{KP: 1, Q: "Why was the original result hard for outsiders to reproduce directly?",
		A: "The evaluation designs and parts of the flow were proprietary, so others had to use different designs", D: []string{"The paper used no neural network", "The algorithm was patented and could not be run", "The results were reported only as images"},
		Why: "Reproduction requires the same inputs and evaluation. With proprietary blocks, outsiders could only test the method on substitutes."},
	{KP: 1, Q: "Which ingredient of the published method made comparisons against classical placers subtle?",
		A: "The agent optimized a proxy cost, and baselines could be scored on the same proxy or on post-route results", D: []string{"The agent placed cells on a continuous canvas with no grid", "The method needed no netlist", "The baselines were all random"},
		Why: "Which metric is compared, proxy or routed PPA, changes the conclusion, so the choice must be explicit."},
	{KP: 2, Q: "An academic assessment ran the released RL code against simulated annealing and other placers on open designs. What kind of question does that comparison answer?",
		A: "Whether the learned policy beats strong non-learned optimizers given comparable inputs and evaluation", D: []string{"Whether RL can be implemented in Python", "Whether GPUs are faster than CPUs", "Whether macros exist in open designs"},
		Why: "A strong classical baseline on the same proxy and flow tests whether learning is what produces the quality."},
	{KP: 2, Q: "Critics examined how well the proxy cost tracked post-route metrics. Why was that measurement central to the debate?",
		A: "If the proxy correlates weakly with routed timing and power, a lower proxy cost does not imply a better chip", D: []string{"Because the proxy defines the legal macro positions", "Because post-route metrics are not affected by macros", "Because the proxy is computed by the router"},
		Why: "The method's training signal was the proxy. Its value for real design depends on that correlation."},
	{KP: 2, Q: "One critique concerned the initial placement information fed to the agent. What was the concern?",
		A: "Placement information from a commercial physical-synthesis step may give the method a head start that should be ablated", D: []string{"The initial placement was illegal and crashed the tool", "The agent had no input at all", "The initial placement was random noise"},
		Why: "If an input already encodes a good placement, part of the reported quality may come from that input rather than the policy."},
	{KP: 2, Q: "Proponents and critics disagreed about pretraining and compute in the reproductions. Why does that matter?",
		A: "A method run without its pretraining or with less compute is not the method as published, while matched compute is needed for a fair baseline comparison", D: []string{"It does not matter for RL methods", "Pretraining only affects runtime, never quality", "Compute matters only for classical methods"},
		Why: "Both sides' arguments turn on accounting: what was run, for how long, with what prior training, for each method."},
	{KP: 3, Q: "What is the main methodological lesson of the macro-placement debate for ML-in-EDA papers?",
		A: "Results need open inputs, strong compute-matched baselines and downstream metrics to be checkable", D: []string{"Learned methods should never be published", "Proxy costs should never be used", "Only industrial designs are valid benchmarks"},
		Why: "The disagreement persisted because the evidence could not be compared on shared ground. Shared ground is the fix."},
	{KP: 3, Q: "A new RL placer paper reports a win over a commercial tool on one proprietary design, with one seed. What is the weakest part of the evidence?",
		A: "One design and one seed say nothing about variance or generality", D: []string{"Commercial tools cannot be compared", "RL results are always reproducible", "One design is enough if it is large"},
		Why: "Stochastic training and design-to-design variation require several designs and seeds before a win means anything."},
	{KP: 3, Q: "Why is simulated annealing a natural baseline for an RL macro placer?",
		A: "It optimizes the same kind of objective over the same discrete grid without any learning", D: []string{"It always finds the global optimum", "It is the method used by every commercial tool", "It needs a neural network"},
		Why: "Same objective, same action space, no learning: any gap is attributable to the policy and its training."},
	{KP: 3, Q: "A tool team wants to adopt a learned macro placer. What evidence should they gather on their own designs first?",
		A: "Post-route PPA and runtime against their current macro placement method across several blocks and seeds", D: []string{"The proxy cost from the paper's designs", "The number of GPU hours in the paper", "The model's parameter count"},
		Why: "Adoption depends on results in their own flow and node, measured on the metrics they sign off with."},
}

// synthesisRecipes: size of the recipe space (KP 1), a Q-learning update (KP 2), and a
// shaped area-delay reward (KP 3).
func synthesisRecipes(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		k := edakit.Int(4, 8)
		l := edakit.Scaled(ctx, 4, 12)
		n := int64(1)
		for i := 0; i < l; i++ {
			n *= int64(k)
		}
		return edakit.Numeric(edakit.Sprintf("A synthesis recipe is a sequence of %d logic-optimization commands, each chosen from %d ABC transforms (such as balance, rewrite, refactor and resub), with repeats allowed. How many distinct recipes exist?", l, k),
			strconv.FormatInt(n, 10), edakit.Sprintf("Each of %d positions has %d choices: $%d^{%d} = %d$. Exhaustive evaluation is out of reach, which motivates learned or bandit search over recipes.", l, k, k, l, n),
			map[string]string{"k": edakit.IntStr(k)})
	case 2:
		q := float64(edakit.Int(-20, 20)) / 10
		a := edakit.Pick([]float64{0.1, 0.2, 0.5})
		g := edakit.Pick([]float64{0.9, 0.95, 1})
		r := float64(edakit.Int(-10, 30)) / 10
		qs := []float64{float64(edakit.Int(-20, 30)) / 10, float64(edakit.Int(-20, 30)) / 10, float64(edakit.Int(-20, 30)) / 10}
		mx := math.Max(qs[0], math.Max(qs[1], qs[2]))
		nq := q + a*(r+g*mx-q)
		ans := edakit.Num(nq, 4)
		return edakit.Numeric(edakit.Sprintf("A Q-learning agent picks the next ABC command for an AIG. Currently $Q(s, \\text{rewrite}) = %s$. It applies rewrite, gets reward $r = %s$ (node reduction in percent), and in the new state the three commands have Q-values %s. With learning rate $\\alpha = %s$ and $\\gamma = %s$, what is the updated $Q(s, \\text{rewrite})$? Round to 4 decimal places.", f(q, 1), f(r, 1), floats(qs, 1), f(a, 1), f(g, 2)),
			ans, edakit.Sprintf("Target $= r + \\gamma \\max_a Q(s', a) = %s + %s \\cdot %s = %s$. Update: $%s + %s(%s - %s) = %s$.", f(r, 1), f(g, 2), pn(mx, 1), f(r+g*mx, 4), f(q, 1), f(a, 1), f(r+g*mx, 4), pn(q, 1), ans),
			map[string]string{"target": f(r+g*mx, 4)})
	default:
		a0 := float64(edakit.Int(800, 3000))
		d0 := float64(edakit.Int(20, 60))
		ap := a0 * float64(edakit.Int(85, 100)) / 100
		dp := d0 * float64(edakit.Int(85, 100)) / 100
		an := ap * float64(edakit.Int(90, 102)) / 100
		dn := dp * float64(edakit.Int(90, 104)) / 100
		ap, dp, an, dn = math.Round(ap), math.Round(dp), math.Round(an), math.Round(dn)
		wa := edakit.Pick([]float64{0.3, 0.5, 0.7})
		r := wa*(ap-an)/a0 + (1-wa)*(dp-dn)/d0
		ans := edakit.Num(r, 4)
		return edakit.Numeric(edakit.Sprintf("An RL recipe tuner scores each command by $r = w\\,(A_{\\mathrm{prev}} - A)/A_0 + (1 - w)(D_{\\mathrm{prev}} - D)/D_0$, where $A$ is AIG node count and $D$ is logic depth, both from the tool. The original circuit has $A_0 = %s$ and $D_0 = %s$. Before the command, $A = %s$ and $D = %s$; after it, $A = %s$ and $D = %s$. With $w = %s$, what is the reward? Round to 4 decimal places.", f(a0, 0), f(d0, 0), f(ap, 0), f(dp, 0), f(an, 0), f(dn, 0), f(wa, 1)),
			ans, edakit.Sprintf("$%s(%s - %s)/%s + %s(%s - %s)/%s = %s$. Node count and depth are proxies: the mapped area and delay after technology mapping can rank recipes differently.", f(wa, 1), f(ap, 0), f(an, 0), f(a0, 0), f(1-wa, 1), f(dp, 0), f(dn, 0), f(d0, 0), ans),
			map[string]string{"area_term": f((ap-an)/a0, 4)})
	}
}
