package realanalysis

// The second real-analysis batch: the rest of limits and continuity, the
// derivative and the integral made rigorous, uniform convergence, and the first
// look at metric spaces.
//
// The first batch built the real line from completeness and proved the theorems
// about sequences, series and the first two-thirds of continuity. This batch
// finishes continuity with connectedness, the classification of discontinuities
// and monotone functions; proves the mean value theorem, Taylor's theorem and
// the Riemann integral with the fundamental theorem, integration by parts and
// improper integrals; treats uniform convergence and the interchange of limits;
// and opens the metric-space picture with open and closed sets, closure,
// completeness, compactness, continuity, the contraction mapping principle,
// Baire's theorem and equivalent metrics.
//
// Ten concepts are numeric, so the definitions can be exercised, and fifteen are
// multiple choice, because the domain is largely about statements and proofs.
// Every numeric generator publishes the parameters its answer follows from, and
// the tests recompute the answer rather than trusting the author's arithmetic.

import (
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func raIntPow(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}

// The 14 multiple-choice generators (connectedGen through equivalentGen)
// are defined with elevated parameterization and variety in variety.go.

// ---------------------------------------------------------------------------
// The mean value theorem (numeric)
// ---------------------------------------------------------------------------

type mvtGen struct{}

func (g *mvtGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []int{2, 4, 6, 8, 10, 3}
	b := claims[rand.Intn(len(claims))]
	answer := decimal6(float64(b) / 2)
	return generator.Problem{
		Question: "For f(x) = x² on the interval from 0 to " + strconv.Itoa(b) +
			", the mean value theorem gives a point c with f'(c) equal to the average slope. What is c? (enter a decimal)",
		Answer:      answer,
		Explanation: "The average slope is (f(b) − f(0))/b = b²/b = b, and f'(c) = 2c, so 2c = b and c = b/2.",
		Facts: map[string]string{
			"b": strconv.Itoa(b), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// The chain rule (numeric)
// ---------------------------------------------------------------------------

type chainGen struct{}

func (g *chainGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		a, b, n, x0 int
	}
	claims := []claim{
		{a: 2, b: 1, n: 3, x0: 1},
		{a: 1, b: 2, n: 4, x0: 0},
		{a: 3, b: 1, n: 2, x0: 2},
		{a: 2, b: 3, n: 3, x0: 0},
		{a: 1, b: 1, n: 5, x0: 1},
		{a: 4, b: 2, n: 2, x0: 1},
	}
	c := claims[rand.Intn(len(claims))]
	inner := c.a*c.x0 + c.b
	answer := decimal6(float64(c.n * c.a * raIntPow(inner, c.n-1)))
	return generator.Problem{
		Question: "By the chain rule, the derivative of f(x) = (ax + b)^n is f'(x) = n·a·(ax + b)^(n−1). For a = " +
			strconv.Itoa(c.a) + ", b = " + strconv.Itoa(c.b) + ", n = " + strconv.Itoa(c.n) +
			", what is f'(x) at x = " + strconv.Itoa(c.x0) + "? (enter a decimal)",
		Answer: answer,
		Explanation: "The inner value is a·x + b = " + strconv.Itoa(inner) +
			", so f'(x) = n·a·(inner)^(n−1) = " + strconv.Itoa(c.n) + "·" + strconv.Itoa(c.a) + "·" +
			strconv.Itoa(raIntPow(inner, c.n-1)) + " = " + answer + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(c.a), "b": strconv.Itoa(c.b), "n": strconv.Itoa(c.n),
			"x0": strconv.Itoa(c.x0), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// L'Hopital's rule (numeric)
// ---------------------------------------------------------------------------

type lhopitalGen struct{}

func (g *lhopitalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []int{2, 3, 4, 5, 6, 1}
	a := claims[rand.Intn(len(claims))]
	answer := decimal6(float64(a*a) / 2)
	return generator.Problem{
		Question: "The limit of (1 − cos(ax))/x² as x tends to 0 is a²/2, by L'Hopital's rule applied twice. For a = " +
			strconv.Itoa(a) + ", what is the limit? (enter a decimal)",
		Answer:      answer,
		Explanation: "Both numerator and denominator tend to zero; differentiating twice gives a²cos(ax)/2, and cos(0) = 1, so the limit is a²/2 = " + answer + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(a), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// The Riemann integral (numeric)
// ---------------------------------------------------------------------------

type riemannGen struct{}

func (g *riemannGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		b, n int
	}
	claims := []claim{
		{b: 1, n: 4},
		{b: 2, n: 4},
		{b: 3, n: 5},
		{b: 2, n: 10},
		{b: 1, n: 5},
		{b: 4, n: 8},
	}
	c := claims[rand.Intn(len(claims))]
	sum := float64(c.b*c.b) * float64(c.n+1) / (2 * float64(c.n))
	answer := decimal6(sum)
	return generator.Problem{
		Question: "Use the right-endpoint Riemann sum with n = " + strconv.Itoa(c.n) + " equal subintervals to approximate the integral of f(x) = x on the interval from 0 to " +
			strconv.Itoa(c.b) + ". What is the value of the sum? (enter a decimal)",
		Answer:      answer,
		Explanation: "With step b/n the sum is Σ (i·b/n)(b/n) for i from 1 to n, which equals b²·n(n+1)/2/n² = b²(n+1)/(2n) = " + answer + ".",
		Facts: map[string]string{
			"b": strconv.Itoa(c.b), "n": strconv.Itoa(c.n), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// The fundamental theorem of calculus (numeric)
// ---------------------------------------------------------------------------

type ftcGen struct{}

func (g *ftcGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(6)
	c := rand.Intn(6)
	b := 2 + rand.Intn(8)
	value := float64(a)*float64(b*b)/2 + float64(c*b)
	answer := decimal6(value)
	return generator.Problem{
		Question: "By the fundamental theorem of calculus, the integral from 0 to b of (a·t + c) with respect to t is a·b²/2 + c·b. For a = " +
			strconv.Itoa(a) + ", c = " + strconv.Itoa(c) + " and b = " + strconv.Itoa(b) + ", what is the value? (enter a decimal)",
		Answer:      answer,
		Explanation: "An antiderivative of a·t + c is a·t²/2 + c·t, and evaluating it at b and at 0 gives a·b²/2 + c·b = " + answer + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(a), "c": strconv.Itoa(c), "b": strconv.Itoa(b), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Integration by parts (numeric)
// ---------------------------------------------------------------------------

type partsGen struct{}

func (g *partsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Six upper limits instead of three. The value is computed from the closed
	// form rather than looked up, so each limit is a different problem.
	claims := []int{1, 2, 3, 4, 5, 6}
	b := claims[rand.Intn(len(claims))]
	value := float64(b-1)*math.Exp(float64(b)) + 1
	answer := decimal4(value)
	return generator.Problem{
		Question: "Use integration by parts to evaluate the integral from 0 to b of x·e^x with respect to x; its value is (b − 1)e^b + 1. For b = " +
			strconv.Itoa(b) + ", what is the value? (enter a decimal)",
		Answer:      answer,
		Explanation: "Integrating by parts with u = x and dv = e^x dx gives [x·e^x − e^x] from 0 to b, which is (b − 1)e^b + 1 = " + answer + ".",
		Facts: map[string]string{
			"b": strconv.Itoa(b), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Improper integrals (numeric)
// ---------------------------------------------------------------------------

type improperGen struct{}

func (g *improperGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Every exponent above one is convergent, so the pool is the range rather
	// than a short list of selected values.
	p := 2 + rand.Intn(38)
	answer := decimal6(1 / float64(p-1))
	return generator.Problem{
		Question: "The improper integral from 1 to infinity of x^(−p) converges exactly when p > 1, to 1/(p − 1). For p = " +
			strconv.Itoa(p) + ", what is its value? (enter a decimal)",
		Answer:      answer,
		Explanation: "An antiderivative is x^(1−p)/(1−p), which tends to 0 at infinity when p > 1, leaving 1/(p − 1) = " + answer + ".",
		Facts: map[string]string{
			"p": strconv.Itoa(p), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Power series (numeric)
// ---------------------------------------------------------------------------

type powerSeriesGen struct{}

func (g *powerSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []int{2, 3, 4, 5, 8, 10}
	c := claims[rand.Intn(len(claims))]
	answer := decimal6(1 / float64(c))
	return generator.Problem{
		Question: "Find the radius of convergence \\(R\\) of the power series \\(\\sum_{n=1}^\\infty " +
			strconv.Itoa(c) + "^n x^n\\). (enter as a decimal)",
		Answer:      answer,
		Explanation: "By the root test (or ratio test), \\(\\limsup_{n\\to\\infty} |" + strconv.Itoa(c) + "^n|^{1/n} = " + strconv.Itoa(c) + "\\), so the radius of convergence is \\(R = \\frac{1}{" + strconv.Itoa(c) + "} = " + answer + "\\).",
		Facts: map[string]string{
			"c": strconv.Itoa(c), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Taylor's theorem (numeric)
// ---------------------------------------------------------------------------

type taylorGen struct{}

func (g *taylorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []int{1, 2, 3, 4, 5, 6}
	a := claims[rand.Intn(len(claims))]
	answer := decimal6(float64(a*a*a) / 6)
	return generator.Problem{
		Question: "The Taylor series at 0 of e^(a·x) has x³-coefficient a³/3! = a³/6. For a = " +
			strconv.Itoa(a) + ", what is that coefficient? (enter a decimal)",
		Answer:      answer,
		Explanation: "The n-th derivative of e^(a·x) at 0 is a^n, so the x³ coefficient is a³/6 = " + answer + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(a), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// The contraction mapping principle (numeric)
// ---------------------------------------------------------------------------

type contractionGen struct{}

func (g *contractionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		pn, pd, b int
	}
	claims := []claim{
		{pn: 1, pd: 2, b: 1},
		{pn: 1, pd: 3, b: 2},
		{pn: 2, pd: 3, b: 1},
		{pn: 3, pd: 4, b: 2},
		{pn: 1, pd: 2, b: 3},
		{pn: 1, pd: 4, b: 1},
	}
	c := claims[rand.Intn(len(claims))]
	a := float64(c.pn) / float64(c.pd)
	fixed := float64(c.b) / (1 - a)
	answer := decimal6(fixed)
	return generator.Problem{
		Question: "The contraction T(x) = a·x + b with a = " + strconv.Itoa(c.pn) + "/" + strconv.Itoa(c.pd) +
			" and b = " + strconv.Itoa(c.b) + " has a unique fixed point, the limit of its iterates from any start. What is the fixed point? (enter a decimal)",
		Answer: answer,
		Explanation: "The fixed point solves x = a·x + b, so x = b/(1 − a) = " + strconv.Itoa(c.b) + "/(1 − " +
			strconv.Itoa(c.pn) + "/" + strconv.Itoa(c.pd) + ") = " + answer + ".",
		Facts: map[string]string{
			"pn": strconv.Itoa(c.pn), "pd": strconv.Itoa(c.pd), "b": strconv.Itoa(c.b), "answer": answer,
		},
	}
}
