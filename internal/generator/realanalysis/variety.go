package realanalysis

import (
	"fmt"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

// distinctDistractors ensures option candidates do not collide with the correct
// answer or with each other, providing fallbacks when collisions occur.
func distinctDistractors(correct string, candidates ...string) []string {
	seen := map[string]bool{correct: true}
	out := make([]string, 0, 3)
	for _, cand := range candidates {
		if cand != "" && !seen[cand] {
			seen[cand] = true
			out = append(out, cand)
			if len(out) == 3 {
				break
			}
		}
	}
	fallbacks := []string{
		"not enough information is provided",
		"only on compact subspaces",
		"none of the other choices holds",
		"holds only for finite sets",
		"diverges everywhere",
	}
	for _, fb := range fallbacks {
		if len(out) == 3 {
			break
		}
		if !seen[fb] {
			seen[fb] = true
			out = append(out, fb)
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// 1. Connectedness (ra.limit.connected)
// ---------------------------------------------------------------------------

type connectedGen struct{}

func (g *connectedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(5)
	b := a + 1 + rand.Intn(4)
	c := b + 1 + rand.Intn(4)
	d := c + 1 + rand.Intn(4)

	switch rand.Intn(5) {
	case 0:
		// Disjoint open intervals
		q := fmt.Sprintf("In the real line $\\mathbb{R}$ with the standard topology, the subset $S = (%d, %d) \\cup (%d, %d)$ is...", a, b, c, d)
		correct := "disconnected"
		dist := distinctDistractors(correct, "connected", "an open interval", "dense in $\\mathbb{R}$")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The gap between the two open intervals splits the set into two separated open components."
		return p
	case 1:
		// Closed intervals meeting at a point
		q := fmt.Sprintf("In $\\mathbb{R}$ with the standard topology, the union of closed intervals $[%d, %d] \\cup [%d, %d]$ is...", a, b, b, d)
		correct := "connected"
		dist := distinctDistractors(correct, "disconnected", "open in $\\mathbb{R}$", "a discrete space")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The intervals share the boundary point, forming the single connected interval [a, d]."
		return p
	case 2:
		// Identification of connected subsets of R
		q := fmt.Sprintf("Which of the following subsets of $\\mathbb{R}$ with the standard topology is connected?")
		correct := fmt.Sprintf("[%d, %d]", a, d)
		dist := distinctDistractors(correct,
			fmt.Sprintf("(%d, %d) \\cup (%d, %d)", a, b, c, d),
			fmt.Sprintf("\\mathbb{Q} \\cap [%d, %d]", a, d),
			fmt.Sprintf("\\{%d, %d, %d\\}", a, b, c),
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The connected subsets of the real line are exactly the intervals."
		return p
	case 3:
		// Components of punctured real line
		q := fmt.Sprintf("How many connected components does the punctured real line $\\mathbb{R} \\setminus \\{%d, %d\\}$ have?", a, c)
		correct := "3"
		dist := distinctDistractors(correct, "1", "2", "4")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Removing two distinct points partitions the real line into three disjoint open intervals."
		return p
	default:
		// Clopen characterization
		q := "A topological space $X$ is connected if and only if the only subsets that are both open and closed are..."
		correct := "the empty set and the entire space $X$"
		dist := distinctDistractors(correct, "all compact subsets", "all singleton subsets", "the interior and closure of $X$")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "A space is connected when it cannot be partitioned into two disjoint nonempty open sets."
		return p
	}
}

// ---------------------------------------------------------------------------
// 2. Classification of discontinuities (ra.limit.discontinuity)
// ---------------------------------------------------------------------------

type discontinuityGen struct{}

func (g *discontinuityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	c := 1 + rand.Intn(7)
	k := 2 + rand.Intn(6)

	switch rand.Intn(5) {
	case 0:
		// Removable discontinuity
		q := fmt.Sprintf("What type of discontinuity does the function $f(x) = \\frac{x^2 - %d}{x - %d}$ have at $x = %d$?", c*c, c, c)
		correct := "removable discontinuity"
		dist := distinctDistractors(correct, "jump discontinuity", "essential discontinuity", "infinite discontinuity")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The limit as x approaches c exists and equals 2c, but f(c) is undefined."
		return p
	case 1:
		// Jump discontinuity
		q := fmt.Sprintf("What type of discontinuity does the piecewise function $f(x) = \\begin{cases} x + %d & x \\ge %d \\\\ x - %d & x < %d \\end{cases}$ have at $x = %d$?", k, c, k, c, c)
		correct := "jump discontinuity"
		dist := distinctDistractors(correct, "removable discontinuity", "essential discontinuity", "infinite discontinuity")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Both one-sided limits exist at c but they are unequal, creating a jump."
		return p
	case 2:
		// Essential oscillatory discontinuity
		q := fmt.Sprintf("What type of discontinuity does $f(x) = \\sin\\left(\\frac{%d}{x - %d}\\right)$ (with $f(%d) = 0$) possess at $x = %d$?", k, c, c, c)
		correct := "essential discontinuity"
		dist := distinctDistractors(correct, "removable discontinuity", "jump discontinuity", "infinite discontinuity")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The function oscillates infinitely often between -1 and 1 as x approaches c; neither one-sided limit exists."
		return p
	case 3:
		// Monotone discontinuity constraint
		q := fmt.Sprintf("If $f: [%d, %d] \\to \\mathbb{R}$ is a monotonically increasing function, what type of discontinuities can $f$ possess?", c, c+10)
		correct := "only jump discontinuities"
		dist := distinctDistractors(correct, "only removable discontinuities", "essential discontinuities", "infinite discontinuities")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "For any monotone function, one-sided limits exist at every point, ruling out essential discontinuities."
		return p
	default:
		// Discontinuity count theorem
		q := "The set of points at which a monotone function $f: \\mathbb{R} \\to \\mathbb{R}$ is discontinuous is always..."
		correct := "at most countable"
		dist := distinctDistractors(correct, "finite", "uncountable", "empty")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Each jump discontinuity contains a disjoint rational interval, which bounds the set to at most countable."
		return p
	}
}

// ---------------------------------------------------------------------------
// 3. Monotone functions (ra.limit.monotone_fn)
// ---------------------------------------------------------------------------

type monotoneFnGen struct{}

func (g *monotoneFnGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(5)
	b := a + 2 + rand.Intn(6)

	switch rand.Intn(5) {
	case 0:
		q := fmt.Sprintf("For a bounded monotonically increasing function $f$ on $(%d, %d)$ and any point $c \\in (%d, %d)$, the one-sided limits...", a, b, a, b)
		correct := "both exist and satisfy $\\lim_{x \\to c^-} f(x) \\le f(c) \\le \\lim_{x \\to c^+} f(x)$"
		dist := distinctDistractors(correct,
			"must be strictly equal, proving continuity",
			"exist only if $f$ is continuous at $c$",
			"diverge whenever $f$ has a jump at $c$",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Monotonicity guarantees that the supremum from the left and infimum from the right bound f(c)."
		return p
	case 1:
		q := fmt.Sprintf("By Lebesgue's differentiation theorem, any monotone function $f: [%d, %d] \\to \\mathbb{R}$ is...", a, b)
		correct := "differentiable almost everywhere"
		dist := distinctDistractors(correct,
			"differentiable everywhere on the open interval",
			"infinitely differentiable",
			"nowhere differentiable",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Lebesgue proved that bounded variation and monotonicity imply differentiability almost everywhere."
		return p
	case 2:
		q := fmt.Sprintf("If a strictly increasing function $f: [%d, %d] \\to \\mathbb{R}$ has the intermediate value property, then $f$ must be...", a, b)
		correct := "continuous on $[a, b]$"
		dist := distinctDistractors(correct, "differentiable on $(a, b)$", "uniformly convex", "a linear polynomial")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "If a monotone function had a jump discontinuity, it would omit values, violating the intermediate value property."
		return p
	case 3:
		q := "The inverse of a continuous, strictly increasing function $f: [a, b] \\to [c, d]$ is..."
		correct := "continuous and strictly increasing on $[c, d]$"
		dist := distinctDistractors(correct,
			"strictly decreasing on $[c, d]$",
			"continuous only on the interior $(c, d)$",
			"not necessarily continuous",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Strict monotonicity preserves order in the inverse, and intervals map to intervals ensuring continuity."
		return p
	default:
		q := "A monotone function $f: [a, b] \\to \\mathbb{R}$ has bounded variation because..."
		correct := "its total variation equals $|f(b) - f(a)|$"
		dist := distinctDistractors(correct,
			"its derivative is uniformly bounded",
			"its graph is a rectifiable polygon",
			"it has no local extrema",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "For a monotone function, differences telescopically sum to the difference between endpoints."
		return p
	}
}

// ---------------------------------------------------------------------------
// 4. Darboux's theorem (ra.deriv.darboux)
// ---------------------------------------------------------------------------

type darbouxGen struct{}

func (g *darbouxGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(4)
	b := a + 2 + rand.Intn(5)
	fa := -5 + rand.Intn(4)
	fb := 2 + rand.Intn(5)

	switch rand.Intn(4) {
	case 0:
		q := fmt.Sprintf("Let $f$ be differentiable on $[%d, %d]$ with $f'(%d) = %d$ and $f'(%d) = %d$. For any value $\\lambda$ strictly between %d and %d, Darboux's theorem guarantees...", a, b, a, fa, b, fb, fa, fb)
		correct := fmt.Sprintf("there exists $c \\in (%d, %d)$ with $f'(c) = \\lambda$", a, b)
		dist := distinctDistractors(correct,
			fmt.Sprintf("f attains its maximum value at $c = (%d + %d)/2$", a, b),
			"f must be twice differentiable on the interval",
			"f' is continuous on the entire interval",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Darboux's theorem establishes that every derivative has the intermediate value property."
		return p
	case 1:
		q := "As a consequence of Darboux's theorem, a derivative function $f'$ on an interval cannot have..."
		correct := "simple jump discontinuities"
		dist := distinctDistractors(correct, "oscillatory discontinuities", "zero values", "unbounded values")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "A jump discontinuity misses intermediate values, which Darboux's theorem forbids for derivatives."
		return p
	case 2:
		k := 2 + rand.Intn(5)
		q := fmt.Sprintf("The function $f(x) = x^2 \\sin\\left(\\frac{%d}{x}\\right)$ for $x \\ne 0$ and $f(0) = 0$ illustrates Darboux's theorem because...", k)
		correct := "$f'(x)$ exists everywhere but is discontinuous at 0"
		dist := distinctDistractors(correct,
			"$f(x)$ is not differentiable at $x = 0$",
			"$f'(x)$ is uniformly continuous on $\\mathbb{R}$",
			"$f'(x)$ does not satisfy the intermediate value property",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The derivative exists at 0 with value 0, yet oscillates wildly near 0, showing derivatives need not be continuous."
		return p
	default:
		q := "A function $g: I \\to \\mathbb{R}$ is said to have the Darboux property on an interval $I$ when..."
		correct := "the image of every subinterval of $I$ is an interval"
		dist := distinctDistractors(correct,
			"it is Riemann integrable on every subinterval",
			"its derivative is continuous on $I$",
			"it is monotone on every compact subinterval",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "A set is connected in R if and only if it is an interval; the Darboux property means preserving connected sets."
		return p
	}
}

// ---------------------------------------------------------------------------
// 5. Uniform convergence (ra.deriv.uniform_conv)
// ---------------------------------------------------------------------------

type uniformConvGen struct{}

func (g *uniformConvGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	pExp := 2 + rand.Intn(4)
	k := 2 + rand.Intn(5)

	switch rand.Intn(4) {
	case 0:
		q := fmt.Sprintf("Consider the series of functions $\\sum_{n=1}^\\infty \\frac{\\cos(%d n x)}{n^%d}$ on $\\mathbb{R}$. By the Weierstrass M-test, this series...", k, pExp)
		correct := fmt.Sprintf("converges uniformly on $\\mathbb{R}$ because $\\left|\\frac{\\cos(%d n x)}{n^%d}\\right| \\le \\frac{1}{n^%d}$", k, pExp, pExp)
		dist := distinctDistractors(correct,
			"diverges at $x = 0$",
			"converges pointwise but not uniformly on $\\mathbb{R}$",
			"converges uniformly only on compact intervals $[0, 1]$",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The terms are bounded by 1/n^p with p > 1, so the Weierstrass M-test implies uniform convergence."
		return p
	case 1:
		q := "If a sequence of continuous functions $f_n: [a, b] \\to \\mathbb{R}$ converges uniformly to $f$ on $[a, b]$, then the limit function $f$ is..."
		correct := "continuous on $[a, b]$"
		dist := distinctDistractors(correct, "differentiable on $(a, b)$", "a polynomial function", "strictly monotonic")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The uniform limit of continuous functions is always continuous."
		return p
	case 2:
		q := "The sequence of functions $f_n(x) = x^n$ on the interval $[0, 1]$ converges pointwise to $f(x) = 0$ for $x < 1$ and $f(1) = 1$. The convergence is..."
		correct := "pointwise but not uniform on $[0, 1]$"
		dist := distinctDistractors(correct,
			"uniform on $[0, 1]$",
			"divergent at all points $x < 1$",
			"uniform on every open subinterval $(0, 1)$",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Each fn is continuous but the limit is discontinuous at 1, so convergence cannot be uniform."
		return p
	default:
		q := "Dini's Theorem guarantees uniform convergence of a sequence of functions $f_n$ on a compact interval $[a, b]$ when..."
		correct := "each $f_n$ is continuous, the sequence is monotone and the pointwise limit is continuous"
		dist := distinctDistractors(correct,
			"each $f_n$ is bounded and the limit is differentiable",
			"the sequence of derivatives $f_n'$ converges uniformly",
			"the functions $f_n$ are polynomials with bounded degrees",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Dini's theorem asserts that monotonicity and compactness elevate pointwise continuity to uniform convergence."
		return p
	}
}

// ---------------------------------------------------------------------------
// 6. Interchange of limits (ra.deriv.interchange)
// ---------------------------------------------------------------------------

type interchangeGen struct{}

func (g *interchangeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(4)
	b := a + 2 + rand.Intn(5)

	switch rand.Intn(4) {
	case 0:
		q := fmt.Sprintf("If continuous functions $f_n: [%d, %d] \\to \\mathbb{R}$ converge uniformly to $f$ on $[%d, %d]$, then $\\lim_{n \\to \\infty} \\int_{%d}^{%d} f_n(x)\\,dx$ equals...", a, b, a, b, a, b)
		correct := fmt.Sprintf("$\\int_{%d}^{%d} f(x)\\,dx$", a, b)
		dist := distinctDistractors(correct,
			fmt.Sprintf("$f(%d) - f(%d)$", b, a),
			"0",
			"the limit does not exist in general",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Uniform convergence justifies interchanging the limit with the Riemann integral on a compact interval."
		return p
	case 1:
		q := "To justify differentiating a sequence of differentiable functions term-by-term on $(a, b)$, it is sufficient that..."
		correct := "the sequence of derivatives $f_n'$ converges uniformly on $(a, b)$ and $f_n(x_0)$ converges at some point $x_0$"
		dist := distinctDistractors(correct,
			"$f_n(x)$ converges uniformly on $(a, b)$",
			"each $f_n$ is a polynomial",
			"$f_n'(x)$ is bounded on $(a, b)$",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Uniform convergence of the derivatives is the essential condition for interchanging limit and derivative."
		return p
	case 2:
		k := 2 + rand.Intn(5)
		q := fmt.Sprintf("The sequence $f_n(x) = \\frac{\\sin(%d n x)}{\\sqrt{n}}$ converges uniformly to 0 on $\\mathbb{R}$. The sequence of derivatives $f_n'(x)$...", k)
		correct := fmt.Sprintf("diverges because $f_n'(0) = %d \\sqrt{n} \\to \\infty$", k)
		dist := distinctDistractors(correct,
			"converges uniformly to 0",
			"converges uniformly to a nonzero constant",
			"converges pointwise everywhere on $\\mathbb{R}$",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Uniform convergence of fn does not imply uniform or even pointwise convergence of fn'."
		return p
	default:
		q := "If $f_n(x)$ is a sequence of Riemann integrable functions on $[a, b]$ converging pointwise to $f(x)$, the integral $\\int_a^b f_n(x)\\,dx$..."
		correct := "does not necessarily converge to $\\int_a^b f(x)\\,dx$ without uniform convergence or a dominated bound"
		dist := distinctDistractors(correct,
			"always converges to $\\int_a^b f(x)\\,dx$",
			"diverges to infinity",
			"equals $f(b) - f(a)$",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Pointwise convergence alone is insufficient to interchange limits with Riemann integration."
		return p
	}
}

// ---------------------------------------------------------------------------
// 7. Open and closed sets (ra.metric.open_closed)
// ---------------------------------------------------------------------------

type openClosedGen struct{}

func (g *openClosedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(5)
	b := a + 2 + rand.Intn(6)

	switch rand.Intn(5) {
	case 0:
		q := fmt.Sprintf("In the real line $\\mathbb{R}$ with the Euclidean metric, the open interval $(%d, %d)$ is...", a, b)
		correct := "open but not closed"
		dist := distinctDistractors(correct, "closed but not open", "both open and closed", "neither open nor closed")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Every point has an open ball around it inside (a, b), but the limit points a and b are not contained in the set."
		return p
	case 1:
		q := fmt.Sprintf("In $\\mathbb{R}$ with the Euclidean metric, the closed interval $[%d, %d]$ is...", a, b)
		correct := "closed but not open"
		dist := distinctDistractors(correct, "open but not closed", "both open and closed", "neither open nor closed")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Its complement is the union of open rays, so it is closed; the endpoints have no open ball inside the set."
		return p
	case 2:
		q := fmt.Sprintf("In $\\mathbb{R}$ with the Euclidean metric, the half-open interval $[%d, %d)$ is...", a, b)
		correct := "neither open nor closed"
		dist := distinctDistractors(correct, "open but not closed", "closed but not open", "both open and closed")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "It contains the point a with no open ball inside it, and omits the limit point b."
		return p
	case 3:
		q := "In any metric space $(X, d)$, any finite subset $\\{x_1, x_2, \\dots, x_k\\}$ is..."
		correct := "always closed"
		dist := distinctDistractors(correct, "always open", "open only if the space is discrete", "never closed")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Singletons in a metric space are closed, and finite unions of closed sets are closed."
		return p
	default:
		q := "In a discrete metric space $(X, d)$, every subset $S \\subset X$ is..."
		correct := "both open and closed (clopen)"
		dist := distinctDistractors(correct, "neither open nor closed", "closed only if finite", "open only if bounded")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The ball B(x, 1/2) is the singleton {x}, so singletons are open; arbitrary unions of singletons are thus open and closed."
		return p
	}
}

// ---------------------------------------------------------------------------
// 8. Closure and boundary (ra.metric.closure)
// ---------------------------------------------------------------------------

type closureGen struct{}

func (g *closureGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(5)
	b := a + 2 + rand.Intn(6)

	switch rand.Intn(5) {
	case 0:
		q := fmt.Sprintf("In $\\mathbb{R}$ with the standard metric, what is the closure of the interval $S = (%d, %d)$?", a, b)
		correct := fmt.Sprintf("[%d, %d]", a, b)
		dist := distinctDistractors(correct,
			fmt.Sprintf("(%d, %d)", a, b),
			fmt.Sprintf("\\{%d, %d\\}", a, b),
			fmt.Sprintf("[%d, %d)", a, b),
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The closure is the union of the set with its limit points, adding the endpoints."
		return p
	case 1:
		q := fmt.Sprintf("In $\\mathbb{R}$ with the standard metric, what is the interior of the interval $S = [%d, %d]$?", a, b)
		correct := fmt.Sprintf("(%d, %d)", a, b)
		dist := distinctDistractors(correct,
			fmt.Sprintf("[%d, %d]", a, b),
			fmt.Sprintf("\\{%d, %d\\}", a, b),
			"\\emptyset",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The interior is the largest open set contained within [a, b], which drops the endpoints."
		return p
	case 2:
		q := fmt.Sprintf("In $\\mathbb{R}$ with the standard metric, what is the boundary $\\partial S$ of the interval $S = (%d, %d)$?", a, b)
		correct := fmt.Sprintf("\\{%d, %d\\}", a, b)
		dist := distinctDistractors(correct,
			fmt.Sprintf("[%d, %d]", a, b),
			fmt.Sprintf("(%d, %d)", a, b),
			"\\emptyset",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The boundary consists of the closure minus the interior: [a, b] \\ (a, b) = {a, b}."
		return p
	case 3:
		q := "In $\\mathbb{R}$ with the standard metric, what is the closure of the set of rational numbers $\\mathbb{Q}$?"
		correct := "$\\mathbb{R}$"
		dist := distinctDistractors(correct, "$\\mathbb{Q}$", "$\\emptyset$", "the set of irrational numbers")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Rationals are dense in R: every real number is a limit point of rationals."
		return p
	default:
		c := 1 + rand.Intn(8)
		q := fmt.Sprintf("In $\\mathbb{R}$ with the standard metric, what is the closure of the set $S = \\{%d + \\frac{1}{n} : n \\in \\mathbb{N}\\}$?", c)
		correct := fmt.Sprintf("$S \\cup \\{%d\\}$", c)
		dist := distinctDistractors(correct,
			"the set $S$ itself",
			fmt.Sprintf("[%d, %d]", c, c+1),
			fmt.Sprintf("\\{%d\\}", c),
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The sequence c + 1/n converges to c, which is the unique limit point of S."
		return p
	}
}

// ---------------------------------------------------------------------------
// 9. Convergence in metric spaces (ra.metric.convergence)
// ---------------------------------------------------------------------------

type metricConvGen struct{}

func (g *metricConvGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	c := 1 + rand.Intn(9)
	k := 2 + rand.Intn(5)

	switch rand.Intn(5) {
	case 0:
		q := fmt.Sprintf("In $\\mathbb{R}$ with the standard metric, to what limit does the sequence $x_n = %d + \\frac{%d}{n^2}$ converge?", c, k)
		correct := fmt.Sprintf("%d", c)
		dist := distinctDistractors(correct, fmt.Sprintf("%d", c+k), fmt.Sprintf("%d", k), "0")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "As n approaches infinity, k/n^2 approaches 0, so xn approaches c."
		return p
	case 1:
		q := "In a discrete metric space $(X, d_{\\text{disc}})$, a sequence $(x_n)$ converges to $x$ if and only if..."
		correct := "there exists $N$ such that $x_n = x$ for all $n \\ge N$ (eventually constant)"
		dist := distinctDistractors(correct,
			"$d(x_n, x) < 1/n$ for all $n$",
			"$x_n$ is bounded",
			"$x_n$ has a convergent subsequence",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "For epsilon < 1, the ball B(x, epsilon) contains only x itself."
		return p
	case 2:
		y0 := 1 + rand.Intn(8)
		q := fmt.Sprintf("In $\\mathbb{R}^2$ with the Euclidean metric, to what point does $p_n = \\left(%d + \\frac{1}{n},\\, %d - \\frac{%d}{n}\\right)$ converge?", c, y0, k)
		correct := fmt.Sprintf("(%d, %d)", c, y0)
		dist := distinctDistractors(correct,
			fmt.Sprintf("(%d, %d)", c+1, y0-k),
			fmt.Sprintf("(%d, %d)", c, y0-k),
			"(0, 0)",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Convergence in Euclidean space is equivalent to coordinate-wise convergence."
		return p
	case 3:
		q := "In any metric space $(X, d)$, if a sequence $(x_n)$ converges, its limit is..."
		correct := "unique"
		dist := distinctDistractors(correct, "nonempty but not necessarily unique", "an isolated point", "a boundary point")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The triangle inequality implies that any sequence in a metric space can converge to at most one point."
		return p
	default:
		q := fmt.Sprintf("In $\\mathbb{R}$ with the standard metric, the sequence $x_n = %d + (-1)^n \\cdot %d$...", c, k)
		correct := "diverges by oscillation"
		dist := distinctDistractors(correct,
			fmt.Sprintf("converges to %d", c),
			fmt.Sprintf("converges to %d", c+k),
			"converges to 0",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The terms oscillate between c-k and c+k, having two distinct subsequential limits."
		return p
	}
}

// ---------------------------------------------------------------------------
// 10. Complete metric spaces (ra.metric.complete)
// ---------------------------------------------------------------------------

type completeGen struct{}

func (g *completeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(5)
	b := a + 2 + rand.Intn(6)

	switch rand.Intn(4) {
	case 0:
		q := "The metric space $(\\mathbb{Q}, d_{\\text{std}})$ of rational numbers with the Euclidean metric is..."
		correct := "incomplete, because Cauchy sequences of rationals can converge to irrationals"
		dist := distinctDistractors(correct,
			"complete, because it has countable cardinality",
			"complete by the Bolzano-Weierstrass theorem",
			"compact and therefore complete",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Rational sequences like (1+1/n)^n are Cauchy in Q but their limit e is not in Q."
		return p
	case 1:
		q := "A subset $A$ of a complete metric space $(X, d)$ is complete in the subspace metric if and only if..."
		correct := "$A$ is a closed subset of $X$"
		dist := distinctDistractors(correct,
			"$A$ is an open subset of $X$",
			"$A$ is bounded in $X$",
			"$A$ is dense in $X$",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Closed subsets of complete metric spaces contain all their limit points, hence all Cauchy limits."
		return p
	case 2:
		q := fmt.Sprintf("The open interval $(%d, %d)$ with the standard Euclidean metric is...", a, b)
		correct := "not complete, because sequences approaching the endpoints are Cauchy but have no limit in the set"
		dist := distinctDistractors(correct,
			"complete, because it is homeomorphic to $\\mathbb{R}$",
			"complete, because it is open in $\\mathbb{R}$",
			"compact and therefore complete",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The sequence a + 1/n is Cauchy in (a, b) but converges to a, which is outside the set."
		return p
	default:
		dim := 2 + rand.Intn(4)
		q := fmt.Sprintf("The Euclidean space $\\mathbb{R}^%d$ is complete because...", dim)
		correct := "every Cauchy sequence in $\\mathbb{R}^%d$ converges to a point in $\\mathbb{R}^%d$"
		dist := distinctDistractors(correct,
			"it is a compact metric space",
			"it contains no isolated points",
			"its Heine-Borel dimension is finite",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Completeness of R implies coordinate-wise completeness of Euclidean n-space."
		return p
	}
}

// ---------------------------------------------------------------------------
// 11. Compact metric spaces (ra.metric.compact)
// ---------------------------------------------------------------------------

type metricCompactGen struct{}

func (g *metricCompactGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(5)
	b := a + 2 + rand.Intn(6)

	switch rand.Intn(5) {
	case 0:
		dim := 1 + rand.Intn(3)
		q := fmt.Sprintf("By the Heine-Borel theorem, a subset of $\\mathbb{R}^%d$ with the Euclidean metric is compact if and only if it is...", dim)
		correct := "closed and bounded"
		dist := distinctDistractors(correct, "open and connected", "complete and countable", "nowhere dense")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "In finite-dimensional Euclidean space, compactness is equivalent to being closed and bounded."
		return p
	case 1:
		q := fmt.Sprintf("In $\\mathbb{R}$ with the Euclidean metric, which of the following subsets is compact?")
		correct := fmt.Sprintf("[%d, %d]", a, b)
		dist := distinctDistractors(correct,
			fmt.Sprintf("(%d, %d)", a, b),
			fmt.Sprintf("[%d, \\infty)", a),
			fmt.Sprintf("\\mathbb{Q} \\cap [%d, %d]", a, b),
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The closed interval [a, b] is both closed and bounded, so it is compact by Heine-Borel."
		return p
	case 2:
		c := 1 + rand.Intn(8)
		q := fmt.Sprintf("In $\\mathbb{R}$ with the standard metric, the set $S = \\{%d + \\frac{1}{n} : n \\in \\mathbb{N}\\} \\cup \\{%d\\}$ is...", c, c)
		correct := "compact, because it is closed and bounded"
		dist := distinctDistractors(correct,
			"not compact, because it contains infinitely many points",
			"open, because each point is isolated",
			"dense in $\\mathbb{R}$",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Including the limit point c makes the bounded set closed, hence compact."
		return p
	case 3:
		q := fmt.Sprintf("Why is the open interval $(%d, %d)$ not compact in $\\mathbb{R}$?", a, b)
		correct := "it is bounded but not closed"
		dist := distinctDistractors(correct,
			"it is closed but not bounded",
			"it has infinite Lebesgue measure",
			"it contains no rational numbers",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "An open interval is bounded, but lacks its boundary limit points, so it is not closed."
		return p
	default:
		q := "If $f: X \\to Y$ is a continuous mapping from a compact metric space $X$ into a metric space $Y$, then the image $f(X)$ is..."
		correct := "compact in $Y$"
		dist := distinctDistractors(correct, "open in $Y$", "complete only if $Y$ is complete", "an interval")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The continuous image of a compact space is always compact."
		return p
	}
}

// ---------------------------------------------------------------------------
// 12. Continuity in metric spaces (ra.metric.continuity)
// ---------------------------------------------------------------------------

type metricContinuityGen struct{}

func (g *metricContinuityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	k := 2 + rand.Intn(6)
	m := 1 + rand.Intn(5)

	switch rand.Intn(5) {
	case 0:
		q := "A mapping $f: (X, d_X) \\to (Y, d_Y)$ between metric spaces is continuous everywhere if and only if..."
		correct := "the preimage $f^{-1}(V)$ of every open set $V \\subset Y$ is open in $X$"
		dist := distinctDistractors(correct,
			"the image $f(U)$ of every open set $U \\subset X$ is open in $Y$",
			"the preimage $f^{-1}(K)$ of every compact set is compact",
			"it preserves the distance between any two points",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "The topological definition of continuity states that preimages of open sets are open."
		return p
	case 1:
		q := fmt.Sprintf("Which of the following functions is uniformly continuous on the entire real line $\\mathbb{R}$?")
		correct := fmt.Sprintf("$f(x) = %d x + %d$", k, m)
		dist := distinctDistractors(correct, "$f(x) = x^2$", "$f(x) = e^x$", "$f(x) = x^3$")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Linear functions are Lipschitz continuous with constant L = k, hence uniformly continuous on R."
		return p
	case 2:
		a := 1 + rand.Intn(4)
		b := a + 3 + rand.Intn(5)
		q := fmt.Sprintf("By the Heine-Cantor theorem, any continuous function $f: [%d, %d] \\to \\mathbb{R}$ on a compact interval is...", a, b)
		correct := "uniformly continuous"
		dist := distinctDistractors(correct, "differentiable on $(a, b)$", "a contraction mapping", "strictly monotonic")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Continuous functions on compact metric spaces are automatically uniformly continuous."
		return p
	case 3:
		L := 2 + rand.Intn(5)
		q := fmt.Sprintf("A function $f: \\mathbb{R} \\to \\mathbb{R}$ satisfying $|f(x) - f(y)| \\le %d |x - y|$ for all $x, y$ is...", L)
		correct := "Lipschitz continuous, which implies uniform continuity"
		dist := distinctDistractors(correct,
			"a contraction mapping on $\\mathbb{R}$",
			"differentiable with $|f'(x)| \\le 1$",
			"bounded on $\\mathbb{R}$",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "A Lipschitz condition with any finite constant L implies uniform continuity (take delta = epsilon/L)."
		return p
	default:
		q := "The epsilon-delta definition of continuity of $f: (X, d_X) \\to (Y, d_Y)$ at $x_0$ states that for every $\\epsilon > 0$..."
		correct := "there exists $\\delta > 0$ such that $d_X(x, x_0) < \\delta \\implies d_Y(f(x), f(x_0)) < \\epsilon$"
		dist := distinctDistractors(correct,
			"there exists $\\delta > 0$ such that $d_Y(f(x), f(x_0)) < \\delta \\implies d_X(x, x_0) < \\epsilon$",
			"for all $\\delta > 0$, $d_X(x, x_0) < \\epsilon \\implies d_Y(f(x), f(x_0)) < \\delta$",
			"there exists $\\delta > 0$ such that $f(B(x_0, \\epsilon)) \\subset B(f(x_0), \\delta)$",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Continuity means keeping outputs within epsilon whenever inputs are within delta."
		return p
	}
}

// ---------------------------------------------------------------------------
// 13. Baire category theorem (ra.metric.baire)
// ---------------------------------------------------------------------------

type baireGen struct{}

func (g *baireGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		q := "The Baire Category Theorem states that every complete metric space $X$..."
		correct := "cannot be expressed as a countable union of nowhere dense closed subsets"
		dist := distinctDistractors(correct,
			"is compact and connected",
			"contains a countable dense subset",
			"has diameter at most 1",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "In a complete metric space, any countable union of nowhere dense sets has empty interior."
		return p
	case 1:
		q := "In any complete metric space $(X, d)$, the intersection of a countable collection of dense open sets is..."
		correct := "dense in $X$"
		dist := distinctDistractors(correct, "empty", "finite", "nowhere dense in $X$")
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "This is the open dual formulation of the Baire Category Theorem."
		return p
	case 2:
		q := "A subset $E$ of a metric space $(X, d)$ is defined to be nowhere dense if..."
		correct := "the interior of its closure is empty: $\\text{int}(\\overline{E}) = \\emptyset$"
		dist := distinctDistractors(correct,
			"it contains no isolated points",
			"its closure is equal to $X$",
			"it is a countable set",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "A set is nowhere dense when its closure contains no nonempty open balls."
		return p
	case 3:
		a := 1 + rand.Intn(5)
		b := a + 3 + rand.Intn(5)
		q := fmt.Sprintf("In the real line $\\mathbb{R}$ with the standard Euclidean metric, which of the following subsets is nowhere dense?")
		correct := "the set of integers $\\mathbb{Z}$"
		dist := distinctDistractors(correct,
			"the set of rational numbers $\\mathbb{Q}$",
			fmt.Sprintf("the open interval (%d, %d)", a, b),
			fmt.Sprintf("the closed ray [%d, \\infty)", a),
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Z is closed and contains no open intervals, so its interior is empty."
		return p
	default:
		q := "The set of rational numbers $\\mathbb{Q}$ is a countable union of singletons $\\{q\\}$, each of which is nowhere dense in $\\mathbb{R}$. Therefore, $\\mathbb{Q}$ is..."
		correct := "of first category (meager) in $\\mathbb{R}$"
		dist := distinctDistractors(correct,
			"of second category (non-meager) in $\\mathbb{R}$",
			"a complete metric space",
			"an open subset of $\\mathbb{R}$",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "A countable union of nowhere dense sets is by definition of first category (meager)."
		return p
	}
}

// ---------------------------------------------------------------------------
// 14. Equivalent metrics (ra.metric.equivalent)
// ---------------------------------------------------------------------------

type equivalentGen struct{}

func (g *equivalentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	dim := 2 + rand.Intn(4)

	switch rand.Intn(5) {
	case 0:
		q := "Two metrics $d_1$ and $d_2$ on a set $X$ are strongly (Lipschitz) equivalent if there exist positive constants $c_1, c_2 > 0$ such that for all $x, y \\in X$..."
		correct := "$c_1 d_1(x, y) \\le d_2(x, y) \\le c_2 d_1(x, y)$"
		dist := distinctDistractors(correct,
			"$d_1(x, y) + d_2(x, y) = c_1$",
			"$d_1(x, y) \\cdot d_2(x, y) \\le c_1$",
			"$d_2(x, y) = c_1 d_1(x, y)^2$",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Strong equivalence requires each metric to be bounded by a constant multiple of the other."
		return p
	case 1:
		q := fmt.Sprintf("In the Euclidean space $\\mathbb{R}^%d$, the $\\ell_1$ taxicab metric, the $\\ell_2$ Euclidean metric and the $\\ell_\\infty$ maximum metric...", dim)
		correct := "are all strongly equivalent and induce the same topology"
		dist := distinctDistractors(correct,
			"induce three strictly distinct topologies",
			"are topologically equivalent but none are strongly equivalent",
			"are equivalent only on the unit ball",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "All norms on a finite-dimensional real vector space are strongly equivalent."
		return p
	case 2:
		q := "Two metrics $d$ and $d'$ on a set $X$ are topologically equivalent if and only if..."
		correct := "they generate the exact same collection of open sets"
		dist := distinctDistractors(correct,
			"they assign the exact same distance to every pair of points",
			"the space is compact under both metrics",
			"all Cauchy sequences are constant",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "Topological equivalence means inducing identical topologies, regardless of numerical distance scaling."
		return p
	case 3:
		q := "For any metric $d$ on $X$, the standard bounded metric defined by $d'(x, y) = \\frac{d(x, y)}{1 + d(x, y)}$..."
		correct := "is topologically equivalent to $d$ and satisfies $d'(x, y) < 1$"
		dist := distinctDistractors(correct,
			"is strongly equivalent to $d$ on unbounded spaces",
			"makes $X$ a compact metric space",
			"changes the convergent sequences in $X$",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "d' preserves small distances and open balls while bounding all distances below 1."
		return p
	default:
		q := "On a finite set $X = \\{x_1, x_2, \\dots, x_n\\}$, any two metrics $d_1$ and $d_2$ are..."
		correct := "strongly equivalent"
		dist := distinctDistractors(correct,
			"never equivalent unless $d_1 = d_2$",
			"topologically equivalent but not strongly equivalent",
			"equivalent only if $n \\le 2$",
		)
		p := generator.MultipleChoice(q, correct, dist...)
		p.Explanation = "On a finite set, the minimum nonzero distance is strictly positive, giving strong equivalence constants."
		return p
	}
}
