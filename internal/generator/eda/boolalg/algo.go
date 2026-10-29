package boolalg

import (
	"math"
	"math/bits"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func analysis(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 8, 30)
		width := edakit.Pick([]int{32, 64})
		lw := 5
		if width == 64 {
			lw = 6
		}
		if n < lw+1 {
			n = lw + 1
		}
		ans := edakit.IntStr(pow2(n - lw))
		q := edakit.Sprintf("A tool stores the complete truth table of a %d-variable function bit-parallel in %d-bit machine words. How many words does it need?", n, width)
		return edakit.Numeric(q, ans, edakit.Sprintf("The table has 2^%d bits; dividing by 2^%d bits per word leaves 2^%d = %s words. Each extra variable doubles it, which is why explicit tables stop near 16 to 20 inputs.", n, lw, n-lw, ans), nil)
	case 2:
		n := edakit.Scaled(ctx, 3, 12)
		switch rand.Intn(3) {
		case 0:
			t := 1
			for i := 1; i <= n; i++ {
				t = 2*t + 1
			}
			ans := edakit.IntStr(t)
			q := edakit.Sprintf("A recursive tautology checker splits on one variable per level with no pruning, so its call count satisfies $T(0) = 1$ and $T(k) = 2T(k-1) + 1$. How many calls does it make on a function of %d variables?", n)
			return edakit.Numeric(q, ans, edakit.Sprintf("Unrolling gives T(k) = 2^(k+1) - 1, so T(%d) = %s: every node of the full Shannon tree is visited.", n, ans), nil)
		case 1:
			c := edakit.Int(2, 5)
			t := 1
			for i := 1; i <= n; i++ {
				t = 2*t + c
			}
			ans := edakit.IntStr(t)
			q := edakit.Sprintf("A recursive procedure on $k$ variables does %d units of work and then recurses on both cofactors: $T(0) = 1$, $T(k) = 2T(k-1) + %d$. What is $T(%d)$?", c, c, n)
			return edakit.Numeric(q, ans, "Iterate the recurrence from T(0): T("+edakit.IntStr(n)+") = "+ans+".", nil)
		default:
			a, b := 1, 1 // T(0), T(1)
			for i := 2; i <= n; i++ {
				a, b = b, a+b+1
			}
			ans := edakit.IntStr(b)
			q := edakit.Sprintf("A branching procedure satisfies $T(0) = T(1) = 1$ and $T(k) = T(k-1) + T(k-2) + 1$. What is $T(%d)$?", n)
			return edakit.Numeric(q, ans, "Iterate upward from the base cases: T("+edakit.IntStr(n)+") = "+ans+". The growth is about 1.618^k, slower than 2^k.", nil)
		}
	default:
		n := edakit.Int(20, 40)
		k := edakit.Int(1, 6)
		t := float64(edakit.Int(2, 30)) / 10
		if rand.Intn(2) == 0 {
			res := t * float64(pow2(k))
			ans := edakit.Num(res, 1)
			q := edakit.Sprintf("An exact algorithm runs in time proportional to $2^n$. It takes %s s on an instance with $n = %d$. How long, in seconds, on a similar instance with $n = %d$? Round to 1 decimal place.", edakit.Num(t, 1), n, n+k)
			return edakit.Numeric(q, ans, edakit.Sprintf("Adding %d variables multiplies the time by 2^%d = %d: %s s.", k, k, pow2(k), ans), nil)
		}
		m := n + edakit.Int(5, 40)
		res := t * math.Pow(float64(m)/float64(n), 3)
		ans := edakit.Num(res, 2)
		q := edakit.Sprintf("A polynomial algorithm runs in time proportional to $n^3$. It takes %s s at $n = %d$. How long, in seconds, at $n = %d$? Round to 2 decimal places.", edakit.Num(t, 1), n, m)
		return edakit.Numeric(q, ans, edakit.Sprintf("Scale by (%d / %d)^3: %s s. Polynomial growth multiplies by a ratio; exponential growth multiplies by a factor per added unit of size.", m, n, ans), nil)
	}
}

func counterFlips(N int) int { return 2*N - bits.OnesCount(uint(N)) }

func doublingCopies(N, c0 int) int {
	capacity, copies := c0, 0
	for size := 0; size < N; size++ {
		if size == capacity {
			copies += size
			capacity *= 2
		}
	}
	return copies
}

func amortized(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		N := edakit.Scaled(ctx, 10, 2000)
		ans := edakit.IntStr(counterFlips(N))
		q := edakit.Sprintf("A binary counter starts at 0 and is incremented %d times. Each increment flips the trailing 1s to 0 and one 0 to 1. How many bit flips happen in total?", N)
		return edakit.Numeric(q, ans, edakit.Sprintf("Bit j flips once every 2^j increments, so the total is 2N minus the number of 1s in N: 2(%d) - %d = %s, under 2 flips per increment.", N, bits.OnesCount(uint(N)), ans), nil)
	case 2:
		N := edakit.Scaled(ctx, 10, 5000)
		c0 := edakit.Pick([]int{1, 2, 4, 8})
		ans := edakit.IntStr(doublingCopies(N, c0))
		q := edakit.Sprintf("A unique table starts with capacity %d and doubles whenever an insert finds it full, copying every stored entry into the new array. After %d inserts, how many entry copies have been made in total?", c0, N)
		return edakit.Numeric(q, ans, "The table grows at sizes "+edakit.IntStr(c0)+", "+edakit.IntStr(2*c0)+", ... below "+edakit.IntStr(N)+", copying that many entries each time. The copies sum to "+ans+", less than 2N.", nil)
	default:
		N := edakit.Scaled(ctx, 10, 5000)
		c0 := edakit.Pick([]int{1, 2, 4})
		cp := doublingCopies(N, c0)
		res := float64(N+cp) / float64(N)
		ans := edakit.Num(res, 2)
		q := edakit.Sprintf("A doubling table starts at capacity %d. Each insert costs 1 unit plus 1 unit per entry copied during a resize. What is the amortized cost per insert over %d inserts? Round to 2 decimal places.", c0, N)
		return edakit.Numeric(q, ans, edakit.Sprintf("Total cost is %d inserts plus %d copies = %d units; divided by %d inserts that is %s, bounded by 3 for any N.", N, cp, N+cp, N, ans), nil)
	}
}

func randomized(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 6, 20)
	d := edakit.Int(1, 8) * pow2(edakit.Int(0, n/2))
	if d > pow2(n)/4 {
		d = pow2(n) / 4
	}
	p := float64(d) / float64(pow2(n))
	switch kpOf(ctx) {
	case 1:
		t := int(math.Round(edakit.Pick([]float64{0.3, 0.5, 1, 1.5, 2, 3}) / p))
		if t < 2 {
			t = 2
		}
		res := math.Pow(1-p, float64(t))
		ans := edakit.Num(res, 4)
		q := edakit.Sprintf("Two %d-input functions differ on exactly %d of the $2^{%d}$ input vectors. A simulator applies %d independent, uniformly random vectors. What is the probability that none of them exposes the difference? Round to 4 decimal places.", n, d, n, t)
		return edakit.Numeric(q, ans, edakit.Sprintf("Each vector misses with probability 1 - %d/%d. Independence gives (1 - p)^%d = %s.", d, pow2(n), t, ans), nil)
	case 2:
		eps := edakit.Pick([]float64{0.1, 0.01, 0.001})
		t := int(math.Ceil(math.Log(eps) / math.Log(1-p)))
		ans := edakit.IntStr(t)
		q := edakit.Sprintf("Two %d-input functions differ on %d of the $2^{%d}$ input vectors. What is the smallest number of independent uniformly random vectors for which the probability of missing the difference is at most %s?", n, d, n, edakit.Num(eps, 3))
		return edakit.Numeric(q, ans, edakit.Sprintf("Solve (1 - p)^t <= %s with p = %d/%d: t >= ln(%s) / ln(1 - p), so t = %s.", edakit.Num(eps, 3), d, pow2(n), edakit.Num(eps, 3), ans), nil)
	default:
		res := 1 / p
		ans := edakit.Num(res, 2)
		q := edakit.Sprintf("Two %d-input functions differ on %d of the $2^{%d}$ input vectors. Random simulation draws independent uniform vectors until one exposes the difference. What is the expected number of vectors drawn, including the one that succeeds? Round to 2 decimal places.", n, d, n)
		return edakit.Numeric(q, ans, edakit.Sprintf("The count is geometric with success probability p = %d/%d, so its mean is 1/p = %s.", d, pow2(n), ans), nil)
	}
}
