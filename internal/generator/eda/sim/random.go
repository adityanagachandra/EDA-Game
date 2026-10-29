package sim

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type cons struct {
	text string
	f    func(x, y int) bool
}

func drawCons(xmax, ymax int) cons {
	c := edakit.Int(2, xmax)
	k := edakit.Int(1, 3)
	opts := []cons{
		{edakit.Sprintf("x + y < %d", c), func(x, y int) bool { return x+y < c }},
		{"x < y", func(x, y int) bool { return x < y }},
		{"x != y", func(x, y int) bool { return x != y }},
		{edakit.Sprintf("x + y == %d", c), func(x, y int) bool { return x+y == c }},
		{"x >= 2 * y", func(x, y int) bool { return x >= 2*y }},
		{edakit.Sprintf("x <= y + %d", k), func(x, y int) bool { return x <= y+k }},
		{edakit.Sprintf("y > %d", k), func(x, y int) bool { return y > k }},
	}
	return edakit.Pick(opts)
}

type space struct {
	wx, wy int
	cs     []cons
}

func (s space) legal() (int, map[int]int) {
	cnt := 0
	byX := map[int]int{}
	for x := 0; x < 1<<s.wx; x++ {
		for y := 0; y < 1<<s.wy; y++ {
			ok := true
			for _, c := range s.cs {
				if !c.f(x, y) {
					ok = false
					break
				}
			}
			if ok {
				cnt++
				byX[x]++
			}
		}
	}
	return cnt, byX
}

func (s space) text() string {
	var ts []string
	for _, c := range s.cs {
		ts = append(ts, c.text+";")
	}
	return edakit.Sprintf("rand bit [%d:0] x; rand bit [%d:0] y; constraint c { %s }", s.wx-1, s.wy-1, strings.Join(ts, " "))
}

func drawSpace(nCons int) (space, int, map[int]int) {
	for {
		s := space{wx: edakit.Int(2, 3), wy: edakit.Int(2, 3)}
		for i := 0; i < nCons; i++ {
			s.cs = append(s.cs, drawCons(1<<s.wx-1, 1<<s.wy-1))
		}
		if nCons == 2 && s.cs[0].text == s.cs[1].text {
			continue
		}
		n, byX := s.legal()
		if n >= 2 && n < (1<<(s.wx+s.wy)) {
			return s, n, byX
		}
	}
}

func constrained(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	nCons := 1
	if kp >= 2 && (kp == 2 || edakit.Int(0, 1) == 0) {
		nCons = 2
	}
	s, n, byX := drawSpace(nCons)
	setup := "A class declares " + s.text() + ". Arithmetic in the constraint does not overflow."
	facts := map[string]string{"solutions": edakit.IntStr(n)}
	if kp <= 2 {
		return edakit.Numeric(setup+" How many (x, y) pairs satisfy the constraints?", edakit.IntStr(n),
			edakit.Sprintf("Enumerating the %d x %d grid of values and keeping the pairs that satisfy every constraint leaves %d legal solutions.", 1<<s.wx, 1<<s.wy, n), facts)
	}
	// Probability of a value of x under uniform-over-solutions.
	var xs []int
	for x := range byX {
		xs = append(xs, x)
	}
	x := edakit.Pick(xs)
	p := float64(byX[x]) / float64(n)
	ans := edakit.Num(p, 3)
	return edakit.Numeric(setup+edakit.Sprintf(" With no dist or solve-before, randomize() picks uniformly among all legal (x, y) pairs. What is the probability that x == %d? Round to 3 decimal places.", x), ans,
		edakit.Sprintf("%d of the %d legal pairs have x = %d, so the probability is %d/%d = %s.%s", byX[x], n, x, byX[x], n, ans, uniformNote(byX, n)), facts)
}

func distSolveBefore(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		w0, w1 := edakit.Int(1, 8), edakit.Int(1, 8)
		m := edakit.Int(3, 15)
		op := edakit.Pick([]string{":=", ":/"})
		which := edakit.Int(0, 1)
		var p float64
		if op == ":=" {
			tot := float64(w0 + m*w1)
			if which == 0 {
				p = float64(w0) / tot
			} else {
				p = float64(w1) / tot
			}
		} else {
			tot := float64(w0 + w1)
			if which == 0 {
				p = float64(w0) / tot
			} else {
				p = float64(w1) / float64(m) / tot
			}
		}
		ans := edakit.Num(p, 3)
		q := edakit.Sprintf("A variable is constrained by x dist { 0 := %d, [1:%d] %s %d }; and has no other constraints. What is the probability that x == %d? Round to 3 decimal places.", w0, m, op, w1, which)
		why := edakit.Sprintf("With :=, every value in the range gets weight %d, so the total weight is %d + %d x %d.", w1, w0, m, w1)
		if op == ":/" {
			why = edakit.Sprintf("With :/, the weight %d is shared across the %d values of the range, so each gets %d/%d and the total weight is %d + %d.", w1, m, w1, m, w0, w1)
		}
		return edakit.Numeric(q, ans, why+" The probability is "+ans+".", map[string]string{"p": ans})
	case 2:
		k := edakit.Int(2, 6)
		c := edakit.Int(1, 4)
		with := edakit.Int(0, 1) == 1
		askD := edakit.Int(0, 1) == 1
		D := float64(int(1) << k)
		var p float64
		cons := "s -> d == 0"
		if c > 1 {
			cons = edakit.Sprintf("s -> d < %d", c)
		}
		cc := float64(c)
		switch {
		case !with && !askD:
			p = cc / (D + cc)
		case with && !askD:
			p = 0.5
		case !with && askD:
			p = 2 / (D + cc)
		default:
			p = 0.5/cc + 0.5/D
		}
		ans := edakit.Num(p, 4)
		mode := "There is no solve-before."
		if with {
			mode = "The class also has constraint order { solve s before d; }."
		}
		target := "s == 1"
		if askD {
			target = "d == 0"
		}
		q := edakit.Sprintf("A class declares rand bit s; rand bit [%d:0] d; constraint c { %s; }. %s What is the probability that %s after randomize()? Round to 4 decimal places.", k-1, cons, mode, target)
		why := edakit.Sprintf("The legal solutions are s = 0 with any of %d values of d, and s = 1 with %d value(s) of d. ", 1<<k, c)
		if with {
			why += edakit.Sprintf("Solving s first picks s = 0 or 1 with probability 1/2 each, then d uniformly among the values legal for that s, giving %s.", ans)
		} else {
			why += edakit.Sprintf("Uniform over all %d solutions gives %s.", (1<<k)+c, ans)
		}
		return edakit.Numeric(q, ans, why, map[string]string{"solutions": edakit.IntStr((1 << k) + c)})
	default:
		wr, ww := edakit.Int(1, 9), edakit.Int(1, 9)
		a, b := edakit.Int(2, 8), edakit.Int(2, 16)
		p := float64(wr)/float64(wr+ww)/float64(a) + float64(ww)/float64(wr+ww)/float64(b)
		ans := edakit.Num(p, 3)
		q := edakit.Sprintf("A bus transaction has rand op_e op; rand int len; with constraint c { op dist { READ := %d, WRITE := %d }; (op == READ) -> len inside {[1:%d]}; (op == WRITE) -> len inside {[1:%d]}; solve op before len; }. What is the probability that len == 1? Round to 3 decimal places.", wr, ww, a, b)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("By total probability, P(len = 1) = P(READ)/%d + P(WRITE)/%d = (%d/%d)/%d + (%d/%d)/%d = %s.", a, b, wr, wr+ww, a, ww, wr+ww, b, ans),
			map[string]string{"p_read": edakit.Num(float64(wr)/float64(wr+ww), 3)})
	}
}

func cornerProbability(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		w1, w2 := edakit.Int(3, 8), edakit.Int(2, 6)
		e := float64(int(1) << (w1 + w2))
		q := edakit.Sprintf("A transaction has an unconstrained %d-bit field addr and an unconstrained %d-bit field len, each uniform and independent. The corner of interest is addr at its maximum value and len == 0 in the same transaction. On average, how many transactions are drawn until the first one hits the corner?", w1, w2)
		return edakit.Numeric(q, edakit.Num(e, 0),
			edakit.Sprintf("The corner probability is p = 2^-%d x 2^-%d = 2^-%d. The number of draws to the first hit is geometric with mean 1/p = %s.", w1, w2, w1+w2, edakit.Num(e, 0)),
			map[string]string{"inv_p": edakit.Num(e, 0)})
	case 2:
		w := edakit.Int(6, 12)
		N := edakit.Int(1, 40) * 100
		p := math.Pow(2, -float64(w))
		pr := 100 * (1 - math.Pow(1-p, float64(N)))
		ans := edakit.Num(pr, 2)
		q := edakit.Sprintf("Each random test hits a corner with probability 2^-%d, independently. What is the probability, in percent, that %d tests hit it at least once? Round to 2 decimal places.", w, N)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("P(no hit) = (1 - 2^-%d)^%d, so P(at least one) = 1 - (1 - 2^-%d)^%d = %s%%.", w, N, w, N, ans),
			map[string]string{"p": edakit.Num(p, 6)})
	default:
		w := edakit.Int(6, 14)
		c := edakit.Pick([]int{90, 95, 99})
		p := math.Pow(2, -float64(w))
		n := math.Ceil(math.Log(1-float64(c)/100) / math.Log(1-p))
		ans := edakit.Num(n, 0)
		q := edakit.Sprintf("Each random test hits a corner with probability 2^-%d, independently. What is the smallest number of tests that hits it at least once with probability at least %d%%?", w, c)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Solve 1 - (1 - p)^N >= %.2f: N >= ln(%.2f) / ln(1 - 2^-%d), which rounds up to %s. That is about %s times the mean wait 1/p.", float64(c)/100, 1-float64(c)/100, w, ans, edakit.Num(n*p, 1)),
			map[string]string{"n": ans})
	}
}

// rejectionSampling asks for the expected number of raw draws per accepted solution.
func rejectionSampling(ctx generator.GeneratorContext) generator.Problem {
	s, n, _ := drawSpace(edakit.Int(1, 2))
	tot := 1 << (s.wx + s.wy)
	e := float64(tot) / float64(n)
	ans := edakit.Num(e, 2)
	q := edakit.Sprintf("A naive solver handles %s by drawing x and y uniformly and rejecting the pair if a constraint fails. On average, how many raw draws does it need per accepted solution?", s.text())
	why := edakit.Sprintf("%d of %d pairs are legal, so acceptance probability is %d/%d and the expected number of draws is %d/%d = %s. As constraints tighten, this ratio explodes, which is why real solvers do not just reject.", n, tot, n, tot, tot, n, ans)
	d := []string{edakit.Num(float64(n)/float64(tot), 2), edakit.IntStr(n), edakit.Num(e/2, 2)}
	if n < tot {
		d = append(d, edakit.Num(float64(tot)/float64(tot-n), 2))
	}
	return edakit.Choice(q, ans, why, d...)
}

var solverItems = []edakit.Item{
	{KP: 1, Q: "A BDD-based stimulus solver builds a BDD of the constraints. How does it then pick a uniformly random solution?",
		A: "It counts solutions below each node and walks down choosing branches in proportion to those counts", D: []string{"It picks a random path ignoring counts", "It enumerates every solution into a list", "It runs a SAT solver on the BDD"},
		Why: "Path counts give exact branch probabilities, so a single top-down walk is an exact uniform sample."},
	{KP: 1, Q: "Why do BDD-based constraint solvers struggle with constraints like a * b == c on 32-bit variables?",
		A: "BDDs for multiplication grow exponentially with bit width", D: []string{"BDDs cannot represent equality", "Multiplication is not allowed in constraints", "BDDs only handle one variable"},
		Why: "Multiplier BDDs are exponential for every variable order, so memory blows up."},
	{KP: 1, Q: "What does a SAT- or SMT-based stimulus solver find easily, and what does it not give for free?",
		A: "It finds a legal solution quickly, but uniform sampling over all solutions is not automatic", D: []string{"It gives uniform samples, but finding any solution is slow", "It counts all solutions exactly", "It only works on unconstrained variables"},
		Why: "SAT search returns some solution biased by its heuristics; near-uniform sampling needs extra techniques such as random hashing constraints."},
	{KP: 1, Q: "How do hashing-based samplers make a SAT solver's output closer to uniform?",
		A: "They add random XOR (parity) constraints that split the solution space into small cells, then pick a solution from a random cell", D: []string{"They restart the solver with a different seed only", "They sort solutions alphabetically", "They remove constraints until all values are legal"},
		Why: "Random XORs partition solutions roughly evenly, so a random cell yields a near-uniform sample."},
	{KP: 1, Q: "Which is a natural encoding of a constraint x + y < 10 on 4-bit variables for a SAT solver?",
		A: "Bit-blast the adder and comparator into CNF clauses over the bits of x and y", D: []string{"One clause per possible value of x", "A single clause x OR y", "A linear program over real numbers"},
		Why: "Word-level arithmetic is bit-blasted into gates and then Tseitin-encoded as CNF."},
	{KP: 1, Q: "Why is solving stimulus constraints NP-hard in general?",
		A: "Constraints can encode Boolean satisfiability, so a general solver solves SAT", D: []string{"Because random numbers are hard to generate", "Because constraints are evaluated in simulation time", "Because constraints must be linear"},
		Why: "Any CNF formula can be written as a constraint block, so the problem inherits SAT's hardness."},
	{KP: 3, Q: "A constraint block makes randomize() take seconds per call. Which change usually helps most?",
		A: "Splitting independent random variables into separate constraint sets so they are solved separately", D: []string{"Adding more rand variables", "Calling randomize() twice", "Converting all variables to 64 bits"},
		Why: "Solvers partition variables connected by constraints; breaking unnecessary links shrinks each problem."},
	{KP: 3, Q: "randomize() returns 0 for a transaction class. What happened?",
		A: "The constraints are contradictory for the current state, so no solution exists", D: []string{"The seed is zero", "All variables were set to zero", "The solver timed out and chose defaults"},
		Why: "randomize() returns 0 when the constraint set is unsatisfiable, and the variables keep their old values."},
	{KP: 3, Q: "A constraint uses a function call: len == calc(mode). How does the solver treat it?",
		A: "It solves mode first, calls the function, then treats the result as a fixed value for len", D: []string{"It inverts the function to solve for mode", "It ignores the constraint", "It samples the function at random inputs"},
		Why: "Function calls in constraints impose an implicit ordering: arguments are solved before the function is evaluated."},
	{KP: 3, Q: "A test needs one specific address pattern 0.01% of the time but random constraints never produce it. What is the standard fix?",
		A: "Add a dist or a weighted constraint mode that biases toward the pattern", D: []string{"Run ten times more seeds without changes", "Remove all constraints", "Make the address a non-random variable"},
		Why: "Biasing the distribution raises the probability of the rare case directly."},
	{KP: 3, Q: "Why do constraint solvers prefer to fix non-random state variables before solving?",
		A: "State variables act as constants, which shrinks the problem the solver must search", D: []string{"State variables cannot appear in constraints", "It makes the distribution uniform", "It is required for randc"},
		Why: "Only rand variables are solved; everything else is a constant that simplifies the constraints."},
	{KP: 3, Q: "A randc variable of 4 bits is randomized 16 times. What is guaranteed?",
		A: "Each of the 16 values appears exactly once before any repeats", D: []string{"The values are uniformly random with repeats allowed", "The values appear in increasing order", "Only even values appear"},
		Why: "randc cycles through a random permutation of its range before repeating."},
}

var seedItems = []edakit.Item{
	{KP: 1, Q: "A nightly regression fails on one seed out of 5,000. What is needed to reproduce the failure?",
		A: "The same seed with the same code and tool versions", D: []string{"Any seed and the same code", "The same seed with any later code version", "Only the waveform file"},
		Why: "The seed fixes the random sequence only for the same program; changing code can change how random calls consume it."},
	{KP: 1, Q: "Why should every regression log record the seed of each test?",
		A: "So any failure can be rerun exactly for debugging", D: []string{"So the simulator runs faster", "So coverage merges correctly", "So that seeds are never reused"},
		Why: "Without the seed, a random failure may never be seen again."},
	{KP: 1, Q: "What does it mean that a constrained-random test is reproducible?",
		A: "Running it again with the same seed produces the same stimulus and results", D: []string{"It passes on every seed", "It always takes the same wall-clock time", "It produces the same results on any seed"},
		Why: "Reproducibility means determinism given the seed; it says nothing about passing."},
	{KP: 1, Q: "A regression uses the time of day as its seed and does not log it. What is the consequence?",
		A: "Failures cannot be reproduced", D: []string{"Coverage cannot be collected", "The test becomes directed", "Every run produces identical stimulus"},
		Why: "An unrecorded seed loses the only key to the random sequence."},
	{KP: 2, Q: "In SystemVerilog, each object and thread has its own random generator seeded from its parent. What property does this give?",
		A: "Random stability: adding randomization in one object does not change the values another object draws", D: []string{"All objects draw the same values", "Randomization becomes uniform over all objects", "Seeds can be omitted"},
		Why: "Hierarchical seeding isolates random streams, so local edits do not perturb unrelated stimulus."},
	{KP: 2, Q: "An engineer adds a new random field to one sequence item class. Which stimulus is most likely to change for the same seed?",
		A: "The values drawn by that item class", D: []string{"Every object in the testbench", "Nothing at all", "Only the DUT's internal state"},
		Why: "Random stability limits the change to the object whose randomization changed, though downstream behavior can still differ."},
	{KP: 2, Q: "Which change most likely breaks random stability across the whole testbench for a given seed?",
		A: "Changing the order in which components and threads are created", D: []string{"Renaming a coverage bin", "Changing a print message", "Adding a comment"},
		Why: "Child generators are seeded from the parent in creation order, so reordering creation shifts their seeds."},
	{KP: 2, Q: "A module's initial blocks call the urandom system function. Why do they become a stability concern when new initial blocks are added?",
		A: "Each new process takes a seed from its parent, which can change the seeds of processes created after it", D: []string{"urandom ignores the seed", "urandom is never reproducible", "urandom only works inside classes"},
		Why: "Thread seeds are derived in creation order; inserting processes changes the order."},
	{KP: 3, Q: "A regression ran 1,000 seeds and reached 92% coverage. Why rank and keep a subset of seeds?",
		A: "A smaller set that preserves the same coverage reruns faster for each code change", D: []string{"Seeds wear out if reused", "Coverage only counts the last seed", "Ranking makes tests directed"},
		Why: "Test ranking keeps the seeds that add coverage, so short regressions keep the same reach."},
	{KP: 3, Q: "A test that failed on seed 42 passes on seed 42 after an unrelated RTL change. What can you conclude?",
		A: "Not that the bug is fixed: the stimulus may have changed, so the failing scenario may not have been reached", D: []string{"The bug is fixed", "The simulator is nondeterministic", "Seed 42 is broken"},
		Why: "Changing code can change the random sequence; the scenario must be checked, for example with a directed test or coverage."},
	{KP: 3, Q: "Why do teams run new random seeds every night instead of rerunning a fixed set?",
		A: "New seeds explore new stimulus and can find bugs the fixed set never reaches", D: []string{"Fixed seeds stop working after a week", "Simulators require it", "New seeds always increase code coverage to 100%"},
		Why: "Fixed seeds give regression stability; fresh seeds give exploration. Both have a role."},
	{KP: 3, Q: "A failure reproduces with seed 7 in one simulator but not in another with the same seed. Why is that expected?",
		A: "Different simulators use different random generators and orderings, so the same seed gives different stimulus", D: []string{"One simulator ignores constraints", "The design has no bugs", "Seeds are only valid in four-state mode"},
		Why: "The LRM does not fix the generator algorithm, so seeds are only meaningful within one tool and version."},
}

// uniformNote remarks that uniform-over-solutions is not uniform per variable, when the drawn
// instance shows it.
func uniformNote(byX map[int]int, n int) string {
	for _, c := range byX {
		if c*len(byX) != n {
			return edakit.Sprintf(" The values of x are not equally likely (there are %d of them): values with more legal partners for y are drawn more often.", len(byX))
		}
	}
	return " Here every legal x has the same number of partners, so x happens to be uniform over its legal values."
}
