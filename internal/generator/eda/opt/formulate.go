package opt

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var variablesObjectiveItems = []edakit.Item{
	{KP: 1, Q: "A gate-sizing tool chooses a drive strength for each of 40,000 cells from a library with 6 sizes per function. What are the decision variables?",
		A: "The size chosen for each cell", D: []string{"The arrival time at each pin, chosen freely", "The 6 library sizes themselves", "The clock period"},
		Why: "Decision variables are what the optimizer is free to set. Arrival times are derived from the sizes through the timer, the library is data, and the period is a fixed constraint."},
	{KP: 1, Q: "In global placement, which quantity is a decision variable rather than data?",
		A: "The (x, y) location of each movable cell", D: []string{"The netlist connectivity", "The positions of fixed I/O pads", "The die outline"},
		Why: "The placer sets movable-cell coordinates. Connectivity, pad positions and the die are inputs the formulation must respect."},
	{KP: 1, Q: "An HLS scheduler assigns operations to clock cycles under a resource limit. Which encoding of the decision makes an ILP possible?",
		A: "A binary $x_{i,t}$ that is 1 when operation i starts in cycle t", D: []string{"A real number per operation equal to its delay", "One integer for the total latency only", "A binary per pair of resources"},
		Why: "Time-indexed binaries let both precedence and per-cycle resource limits be written as linear constraints."},
	{KP: 1, Q: "A floorplanner may rotate each soft macro by 90 degrees. How is that choice usually modeled?",
		A: "One binary rotation variable per macro that swaps its width and height", D: []string{"A continuous angle between 0 and 90 degrees", "A constraint on the die aspect ratio", "A change to the netlist"},
		Why: "A rotation is a discrete choice, so it becomes a binary that selects which dimension is the width."},
	{KP: 2, Q: "A buffering tool must meet every endpoint's required time and use as little buffer area as possible. How should it be formulated?",
		A: "Minimize buffer area subject to timing constraints at every endpoint", D: []string{"Minimize the sum of arrival times with no constraints", "Maximize buffer count subject to an area limit", "Minimize worst slack subject to an area limit"},
		Why: "Meeting timing is a hard requirement, so it is a constraint; area is what we trade, so it is the objective."},
	{KP: 2, Q: "A placer adds a density term to its wirelength objective with weight $\\lambda$ instead of a hard bin capacity. What is the density requirement now?",
		A: "A soft constraint: violations are allowed but penalized", D: []string{"A hard constraint enforced at every iteration", "A decision variable", "Irrelevant to the result"},
		Why: "Moving a constraint into the objective with a weight makes it soft; the optimizer may trade some overflow for wirelength unless the weight grows."},
	{KP: 2, Q: "A router minimizes total wirelength plus 1000 times total overflow. A run returns overflow 3 with much shorter wires than an overflow-free run. What does this reveal?",
		A: "Overflow is a weighted soft term, so the weight was too low to force legality", D: []string{"The router found a bug in the capacity model", "Wirelength is a constraint in this formulation", "Overflow-free routing is impossible"},
		Why: "With a finite penalty the optimizer is allowed to buy wirelength with overflow; legality needs a hard constraint or a larger weight."},
	{KP: 2, Q: "Timing-driven placement wants to reduce the worst negative slack. Written as an optimization, what is the cleanest objective form?",
		A: "Minimize t subject to t greater than or equal to each path's lateness", D: []string{"Minimize the product of all slacks", "Maximize the number of paths", "Minimize the clock period with no other constraint"},
		Why: "A max of terms is linearized with an auxiliary variable t bounded below by every term: an epigraph formulation."},
	{KP: 3, Q: "A formulation for clock-tree sizing has no feasible point. What is the first thing to check?",
		A: "Whether the constraints conflict, for example a skew bound tighter than wire delays allow", D: []string{"Whether the objective is convex", "Whether the step size is too large", "Whether the objective weights sum to 1"},
		Why: "Infeasibility comes from constraints alone; the objective plays no part in whether a feasible point exists."},
	{KP: 3, Q: "Two formulations of the same legalization problem have the same optimum, but one has an LP relaxation bound much closer to it. Why does that matter?",
		A: "A tighter relaxation prunes more of the branch-and-bound tree, so the solver is faster", D: []string{"It changes the optimal integer solution", "It makes the problem convex", "It removes the need for integer variables"},
		Why: "Equivalent integer models can differ greatly in relaxation strength, and solve time follows the bound quality."},
	{KP: 3, Q: "A sizing objective sums delay in ps and area in square microns with weights 1 and 1. Area values are around 10,000 while delays are around 200. What happens?",
		A: "Area dominates the objective, so the optimizer nearly ignores timing", D: []string{"The two terms are balanced automatically", "The optimizer rejects mixed units", "Delay dominates because it is listed first"},
		Why: "Weights must account for scale. Without normalization the larger-magnitude term controls the trade-off."},
	{KP: 3, Q: "A placer optimizes HPWL, but the final metric is routed wirelength after detailed routing. What is this a case of?",
		A: "Optimizing a tractable proxy of the true objective", D: []string{"A constraint violation", "A nonconvex constraint", "A decision variable error"},
		Why: "Most EDA objectives are proxies (HPWL for routed length, Elmore for signoff delay); the gap between proxy and signoff metric must be tracked."},
}

// searchSpaceSize asks for the size of a discrete decision space.
func searchSpaceSize(ctx generator.GeneratorContext) generator.Problem {
	switch edakit.Int(0, 2) {
	case 0:
		n := edakit.Int(5, 10)
		k := edakit.Int(2, 4)
		v := int(math.Pow(float64(k), float64(n)))
		q := edakit.Sprintf("A sizing tool picks one of %d drive strengths for each of %d gates independently. How many distinct sizings are there?", k, n)
		return edakit.Choice(q, edakit.IntStr(v),
			edakit.Sprintf("Each gate has %d choices, so there are $%d^{%d} = %d$ sizings. The space is discrete and grows exponentially with the gate count.", k, k, n, v),
			edakit.Distinct(v, 3, int(math.Pow(float64(n), float64(k))), n*k, fact(n))...)
	case 1:
		n := edakit.Int(4, 9)
		v := fact(n)
		q := edakit.Sprintf("A detailed placer assigns %d cells to %d legal slots, one cell per slot. How many assignments are possible?", n, n)
		return edakit.Choice(q, edakit.IntStr(v),
			edakit.Sprintf("An assignment is a permutation, so there are $%d! = %d$. This is a combinatorial problem: no gradient exists between permutations.", n, v),
			edakit.Distinct(v, 3, int(math.Pow(2, float64(n))), n*n, int(math.Pow(float64(n), float64(n))))...)
	default:
		n := 2 * edakit.Int(2, 6)
		v := binom(n, n/2) / 2
		q := edakit.Sprintf("A partitioner splits %d cells into two unlabeled halves of %d cells each. How many bisections are there?", n, n/2)
		return edakit.Choice(q, edakit.IntStr(v),
			edakit.Sprintf("Choose the %d cells of one half, $\\binom{%d}{%d} = %d$, then divide by 2 because swapping the halves gives the same bisection: %d.", n/2, n, n/2, binom(n, n/2), v),
			edakit.Distinct(v, 3, binom(n, n/2), int(math.Pow(2, float64(n))), int(math.Pow(2, float64(n-1))))...)
	}
}

func fact(n int) int {
	r := 1
	for i := 2; i <= n; i++ {
		r *= i
	}
	return r
}

var problemClassItems = []edakit.Item{
	{KP: 2, Q: "Minimize $\\sum w_{ij}(x_i - x_j)^2$ over real cell coordinates, with some coordinates fixed. What class of problem is this?",
		A: "An unconstrained convex quadratic program", D: []string{"An integer linear program", "A nonconvex QCQP", "A linear program"},
		Why: "The objective is a positive semidefinite quadratic in continuous variables with no inequality constraints; solving it is one sparse linear system."},
	{KP: 2, Q: "Choose 0-1 variables to cover every fault with the fewest test patterns. What class is this?",
		A: "An integer linear program (set covering)", D: []string{"A linear program", "A convex QP", "An unconstrained smooth problem"},
		Why: "Binary selection variables, a linear objective and linear covering constraints make a 0-1 ILP."},
	{KP: 2, Q: "Minimize maximum lateness with linear delay models and continuous buffer delays. Which class fits after introducing an epigraph variable?",
		A: "A linear program", D: []string{"A mixed-integer nonlinear program", "A quadratic program", "A combinatorial search with no relaxation"},
		Why: "A max of linear functions becomes linear constraints on an auxiliary variable, so the whole model is an LP."},
	{KP: 2, Q: "Which problem is a mixed-integer linear program?",
		A: "Floorplanning with continuous block coordinates and binary relative-position variables in linear non-overlap constraints", D: []string{"Quadratic placement with fixed pads", "Gradient descent on wirelength", "Computing an Elmore delay"},
		Why: "It mixes continuous coordinates with binaries and keeps every constraint linear through big-M terms."},
	{KP: 2, Q: "Gate sizing with continuous sizes and Elmore-style delays written as posynomials is which kind of problem?",
		A: "A geometric program, convex after a log change of variables", D: []string{"An integer program", "A linear program in the sizes", "A nonconvex problem with no convex form"},
		Why: "Posynomial objectives and constraints become convex in log variables, which is why continuous sizing (as in TILOS-era work) has global optima."},
	{KP: 2, Q: "A model constrains every pair of cells to satisfy $(x_i - x_j)^2 + (y_i - y_j)^2 \\ge d^2$. What does this make the problem?",
		A: "A nonconvex quadratically constrained program", D: []string{"A convex QCQP", "A linear program", "An unconstrained QP"},
		Why: "A 'greater than or equal to' bound on a convex quadratic describes the outside of a disk, a nonconvex set."},
	{KP: 3, Q: "Why do global placers solve a continuous problem first even though legal placement is discrete?",
		A: "The continuous relaxation is tractable with gradients and gives a good starting point for legalization", D: []string{"Legal placement is continuous by definition", "Discrete placement has no objective", "Rows and sites do not affect placement"},
		Why: "Relaxing site legality makes the problem smooth and solvable at millions of cells; legalization then repairs the small discrete violations."},
	{KP: 3, Q: "Detailed routing on a track grid, with discrete layer and track choices, is best described as which kind of problem?",
		A: "Combinatorial: a discrete choice of grid paths subject to design rules", D: []string{"Continuous convex", "An unconstrained QP", "A linear least-squares fit"},
		Why: "Wires live on discrete tracks and layers, so the feasible set is a finite collection of paths."},
	{KP: 3, Q: "A continuous optimizer returns sizes like 2.37x for a library offering 1x, 2x, 4x and 8x. What step remains?",
		A: "Discretization: snap to library sizes and repair the timing that rounding breaks", D: []string{"Nothing: tools accept fractional sizes", "Rerun the same continuous solve", "Change the netlist"},
		Why: "Continuous sizing is a relaxation of the discrete library choice; rounding can violate constraints, so a repair or local search follows."},
	{KP: 3, Q: "What is the key structural difference between continuous and combinatorial optimization for an algorithm designer?",
		A: "Continuous problems offer gradients and local models; combinatorial ones need enumeration, relaxation or search", D: []string{"Combinatorial problems are always smaller", "Continuous problems never have local optima", "Combinatorial problems cannot have constraints"},
		Why: "Gradients and Taylor models exist only on continuous spaces; discrete spaces need bounds, neighborhoods or relaxations instead."},
	{KP: 3, Q: "Technology mapping chooses library cells to cover a subject graph. Which class does it belong to?",
		A: "Combinatorial covering, solved exactly on trees by dynamic programming", D: []string{"Convex continuous optimization", "A linear program with real variables", "Unconstrained smooth minimization"},
		Why: "Each match is a discrete choice; on trees the cover decomposes into subproblems that DP solves exactly."},
	{KP: 3, Q: "A tool tunes 10 continuous flow knobs, but each evaluation is a 6-hour place-and-route run with no gradient. What class of problem is this?",
		A: "Expensive black-box optimization", D: []string{"A linear program", "A convex QP solved by Newton's method", "A combinatorial covering problem"},
		Why: "With no derivatives and costly evaluations, the problem calls for surrogate-based methods such as Bayesian optimization."},
}

// enumerationTime asks how long exhaustive enumeration would take.
func enumerationTime(ctx generator.GeneratorContext) generator.Problem {
	rate := edakit.Pick([]int{6, 7, 8, 9})
	var log10N float64
	var what string
	if edakit.Int(0, 1) == 0 {
		k := edakit.Int(2, 4)
		n := edakit.Int(30, 120)
		log10N = float64(n) * math.Log10(float64(k))
		what = edakit.Sprintf("all $%d^{%d}$ sizings of %d gates with %d sizes each", k, n, n, k)
	} else {
		n := edakit.Int(15, 40)
		lg, _ := math.Lgamma(float64(n + 1))
		log10N = lg / math.Ln10
		what = edakit.Sprintf("all $%d!$ orderings of %d cells in a row", n, n)
	}
	e := int(math.Floor(log10N - float64(rate)))
	q := edakit.Sprintf("An exhaustive optimizer evaluates $10^{%d}$ configurations per second. Roughly how long does it take to enumerate %s?", rate, what)
	ans := edakit.Sprintf("About $10^{%d}$ seconds", e)
	return edakit.Choice(q, ans,
		edakit.Sprintf("The space has about $10^{%s}$ points; dividing by $10^{%d}$ per second leaves about $10^{%d}$ seconds. For scale, the age of the universe is about $4 \\times 10^{17}$ seconds.", f(log10N, 1), rate, e),
		edakit.Sprintf("About $10^{%d}$ seconds", e-rate), edakit.Sprintf("About $10^{%d}$ seconds", e+3), edakit.Sprintf("About $10^{%d}$ seconds", int(math.Floor(log10N))), edakit.Sprintf("About $10^{%d}$ seconds", e/2))
}

var hardnessItems = []edakit.Item{
	{KP: 2, Q: "Min-cut bisection, floorplanning and detailed placement are NP-hard. What does that mean for a tool designer?",
		A: "No known algorithm solves every instance optimally in polynomial time, so tools rely on heuristics and bounds", D: []string{"Every instance takes exponential time", "The problems have no optimal solution", "Exact solvers can never be used on any instance"},
		Why: "NP-hardness is a worst-case statement. Small or structured instances may still be solved exactly, but scaling demands heuristics."},
	{KP: 2, Q: "A colleague says ILP global routing is fine because the LP relaxation solves in seconds. What is the flaw?",
		A: "The integer problem can need exponentially many branch-and-bound nodes even when each LP is fast", D: []string{"LPs cannot be solved in polynomial time", "Global routing has no integer variables", "LP relaxations are always integral"},
		Why: "Each node's LP is cheap, but the number of nodes is what explodes for hard integer programs."},
	{KP: 2, Q: "Why is 2-coloring a conflict graph easy while 3-coloring (as in triple patterning lithography decomposition) is hard?",
		A: "2-colorability is a bipartiteness check in linear time; 3-colorability is NP-complete", D: []string{"Both are NP-complete", "Both are linear time", "3-coloring is easy and 2-coloring is hard"},
		Why: "A breadth-first search decides 2-coloring; 3-coloring is a classic NP-complete problem, which is why TPL decomposition uses heuristics or ILP on small components."},
	{KP: 2, Q: "Which property usually makes a continuous optimization problem tractable to solve to global optimality?",
		A: "Convexity of the objective and the feasible set", D: []string{"Having many variables", "Being smooth but nonconvex", "Having integer variables"},
		Why: "For a convex problem every local optimum is global, and interior-point and first-order methods have polynomial guarantees."},
	{KP: 2, Q: "A problem is NP-hard in general but the instances a router sees are planar with small degree. What is the right conclusion?",
		A: "Structure may allow exact or near-exact methods on those instances; hardness is a worst-case statement", D: []string{"The instances are guaranteed to be hard", "NP-hardness disappears for planar graphs in all problems", "The router must use brute force"},
		Why: "Restricted classes can be easy (some problems become polynomial on planar or bounded-treewidth graphs), so structure should be exploited."},
	{KP: 2, Q: "Which statement about local optima is correct?",
		A: "A nonconvex continuous problem can have many local optima, and descent finds only the one whose basin it starts in", D: []string{"Gradient descent always finds the global optimum", "Discrete problems have no local optima", "Local optima exist only in linear programs"},
		Why: "Descent methods are local; without convexity the starting point decides which optimum is found."},
	{KP: 3, Q: "Given NP-hardness, what do production EDA engines typically promise about solution quality?",
		A: "Good solutions in bounded runtime, often with no optimality certificate", D: []string{"Provably optimal solutions on every design", "Random solutions", "Optimality whenever the design is small"},
		Why: "Engines trade optimality for scalability, and quality is measured by benchmarks and QoR baselines rather than proofs."},
	{KP: 3, Q: "Why does a MILP solver's reported gap matter in an EDA flow?",
		A: "It bounds how far the returned solution can be from optimal, so the flow can stop early with a known guarantee", D: []string{"It measures the runtime", "It is the number of integer variables", "It shows the LP is infeasible"},
		Why: "The gap between incumbent and best bound is a certificate; stopping at 1% gap saves enormous time for a known loss."},
	{KP: 3, Q: "A partitioner's runtime grows from 1 s at 10^4 cells to 100 s at 10^5 cells. What empirical complexity does that suggest?",
		A: "About quadratic", D: []string{"About linear", "About cubic", "Exponential"},
		Why: "Ten times the size and a hundred times the time is $n^2$ scaling, which is too slow at $10^7$ cells and motivates multilevel methods."},
	{KP: 3, Q: "Which technique makes a hard problem easier by dropping constraints to get a bound?",
		A: "Relaxation", D: []string{"Enumeration", "Rounding", "Restarting"},
		Why: "Relaxing integrality or hard constraints gives an easier problem whose optimum bounds the original."},
	{KP: 3, Q: "Why do multilevel partitioners and placers coarsen the netlist first?",
		A: "Solving a much smaller clustered problem gives a good global structure cheaply, refined level by level", D: []string{"Coarsening makes the problem convex", "Coarsening removes all local optima", "Coarsening is required for legality"},
		Why: "Multilevel methods decompose by scale: global decisions on few clusters, local refinement on many cells."},
	{KP: 3, Q: "An approximation algorithm for set cover returns a cover at most about $\\ln n$ times the optimum. What kind of guarantee is that?",
		A: "A worst-case approximation ratio", D: []string{"An average-case runtime", "An optimality certificate", "A convexity proof"},
		Why: "The greedy set-cover algorithm has a logarithmic approximation ratio: a bound on quality relative to the unknown optimum."},
}

var strategyItems = []edakit.Item{
	{KP: 1, Q: "Quadratic global placement ignores cell overlap and site legality, then spreads cells. Which strategy is it?",
		A: "Relax: drop hard constraints, solve the easier problem, then repair", D: []string{"Decompose by Benders cuts", "Exhaustive search", "Dynamic programming"},
		Why: "The overlap and site constraints are relaxed away and restored later by spreading and legalization."},
	{KP: 1, Q: "Lagrangian-relaxation gate sizing moves the timing constraints into the objective with multipliers. What is gained?",
		A: "The relaxed problem separates into easy per-gate subproblems", D: []string{"The problem becomes integer", "Timing is ignored entirely", "The solution is always infeasible"},
		Why: "Dualizing coupling constraints decouples the gates; multipliers coordinate them through iteration."},
	{KP: 1, Q: "A global router solves the LP relaxation of its ILP and then rounds the fractional route choices. Which strategy pair is this?",
		A: "Relaxation followed by approximation (rounding)", D: []string{"Decomposition followed by enumeration", "Search followed by dynamic programming", "Pure exhaustive search"},
		Why: "The LP relaxes integrality; randomized or deterministic rounding approximates an integer solution from it."},
	{KP: 1, Q: "Which is a relaxation of a 0-1 ILP?",
		A: "Replacing $x \\in \\{0, 1\\}$ with $0 \\le x \\le 1$", D: []string{"Adding a cutting plane", "Fixing some variables to 1", "Branching on a fractional variable"},
		Why: "A relaxation enlarges the feasible set, so its optimum bounds the original from the optimistic side."},
	{KP: 2, Q: "Min-cut placement recursively bisects the netlist and the die. Which strategy is this?",
		A: "Decompose: split one big problem into many small ones", D: []string{"Relax integrality", "Penalty method", "Bayesian search"},
		Why: "Recursive bisection is divide and conquer: each region becomes an independent smaller placement problem."},
	{KP: 2, Q: "A hierarchical flow optimizes each block separately and then the top level with block abstracts. What is the main risk of this decomposition?",
		A: "Interactions across block boundaries, such as cross-block timing, can be suboptimal or missed", D: []string{"Each block becomes NP-hard", "The top level has no variables", "Decomposition always gives the global optimum"},
		Why: "Decomposition is exact only when subproblems are independent; coupling must be handled by budgets or iteration."},
	{KP: 2, Q: "Benders decomposition splits an ILP into an integer master and an LP subproblem. What links them?",
		A: "Cuts generated from the subproblem's dual, added to the master", D: []string{"Random restarts", "Shared gradients", "A cooling schedule"},
		Why: "Each subproblem solve returns a feasibility or optimality cut that tells the master about the continuous part."},
	{KP: 2, Q: "ADMM is used to split a placement objective across partitions that share boundary cells. Which strategy is that?",
		A: "Decomposition with coordination through dual variables", D: []string{"Pure relaxation with no coordination", "Exhaustive enumeration", "Simulated annealing"},
		Why: "Each partition solves its own subproblem and the dual (consensus) variables drive the shared copies to agree."},
	{KP: 3, Q: "Simulated annealing floorplanning perturbs a Polish expression and accepts or rejects. Which strategy is it?",
		A: "Search: explore the discrete space with moves and an acceptance rule", D: []string{"Relaxation with an exact bound", "Benders decomposition", "Convex optimization"},
		Why: "SA searches the space directly without bounds or relaxation, trading certificates for flexibility."},
	{KP: 3, Q: "Which approach gives a lower bound on the optimum of a minimization problem?",
		A: "Solving a relaxation", D: []string{"Running a heuristic search", "Rounding a fractional solution", "Evaluating a random feasible point"},
		Why: "Any feasible point gives an upper bound; a relaxation's optimum gives a lower bound. Together they bracket the optimum."},
	{KP: 3, Q: "A team must solve 10^6-cell detailed placement in minutes and also wants a quality certificate. Which pairing fits?",
		A: "Local search for the solution plus a relaxation bound to measure the gap", D: []string{"Exhaustive enumeration", "An ILP over all cells at once", "Random placement with no bound"},
		Why: "Search scales; a relaxation provides the bound that says how much quality might be left on the table."},
	{KP: 3, Q: "Which EDA problem is solved exactly and efficiently by dynamic programming rather than relaxation or search?",
		A: "Optimal buffer insertion on a routing tree with a fixed buffer library (van Ginneken)", D: []string{"General 3-coloring of a conflict graph", "Full-chip detailed routing", "Global placement with overlap constraints"},
		Why: "The tree structure gives optimal substructure, so a bottom-up DP with candidate pruning is exact."},
}
