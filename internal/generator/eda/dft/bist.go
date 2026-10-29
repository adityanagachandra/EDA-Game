package dft

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// lfsr is an external-XOR (Fibonacci) shift register: Q1 takes the XOR of the tapped stages and
// every other stage takes its left neighbour. Stage n is always tapped.
type lfsr struct {
	n    int
	taps []int // 1-based
}

func (l lfsr) step(q []int, in []int) []int {
	fb := 0
	for _, t := range l.taps {
		fb ^= q[t-1]
	}
	nq := make([]int, l.n)
	nq[0] = fb
	for i := 1; i < l.n; i++ {
		nq[i] = q[i-1]
	}
	if in != nil {
		for i := range nq {
			nq[i] ^= in[i]
		}
	}
	return nq
}

func randLFSR(n int) lfsr {
	taps := []int{n}
	for t := 1; t < n; t++ {
		if edakit.Int(0, 2) == 0 {
			taps = append(taps, t)
		}
	}
	if len(taps) == 1 {
		taps = append(taps, edakit.Int(1, n-1))
	}
	return lfsr{n, sortInts(taps)}
}

func sortInts(a []int) []int {
	for i := range a {
		for j := i + 1; j < len(a); j++ {
			if a[j] < a[i] {
				a[i], a[j] = a[j], a[i]
			}
		}
	}
	return a
}

func (l lfsr) tapText() string {
	var p []string
	for _, t := range l.taps {
		p = append(p, edakit.Sprintf("Q_%d", t))
	}
	return "$" + strings.Join(p, " \\oplus ") + "$"
}

// poly renders the characteristic polynomial x^n + sum_{t in T} x^{n-t}.
func (l lfsr) poly() string {
	var p []string
	p = append(p, edakit.Sprintf("x^{%d}", l.n))
	for _, t := range l.taps {
		e := l.n - t
		switch e {
		case 0:
			p = append(p, "1")
		case 1:
			p = append(p, "x")
		default:
			p = append(p, edakit.Sprintf("x^{%d}", e))
		}
	}
	return "$" + strings.Join(p, " + ") + "$"
}

func (l lfsr) period(seed []int) int {
	q := append([]int(nil), seed...)
	for k := 1; k <= 1<<l.n; k++ {
		q = l.step(q, nil)
		if bits(q) == bits(seed) {
			return k
		}
	}
	return 0
}

func randBits(n int, nonzero bool) []int {
	for attempt := 0; ; attempt++ {
		q := make([]int, n)
		any := false
		for i := range q {
			q[i] = edakit.Int(0, 1)
			if q[i] == 1 {
				any = true
			}
		}
		if any || !nonzero {
			return q
		}
		if attempt >= generator.MaxDrawAttempts {
			q[0] = 1
			return q
		}
	}
}

// lfsrMISR: one LFSR step (KP 1), a MISR signature (KP 2), the state after k steps (KP 3).
func lfsrMISR(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 8)
	l := randLFSR(n)
	stem := edakit.Sprintf("An external-XOR LFSR with %d stages shifts $Q_1 \\to Q_2 \\to \\ldots \\to Q_%d$ on each clock while $Q_1$ takes the feedback %s.", n, n, l.tapText())
	switch kp(ctx) {
	case 1:
		q := randBits(n, true)
		nq := l.step(q, nil)
		ans := bits(nq)
		return edakit.Numeric(stem+edakit.Sprintf(" The state $Q_1 \\ldots Q_%d$ is %s. What is the state after one clock? Answer with the %d bits in order.", n, bits(q), n),
			ans, edakit.Sprintf("Feedback = %d, which enters $Q_1$; the rest shift right: %s.", nq[0], ans), map[string]string{"state": ans})
	case 2:
		k := edakit.Scaled(ctx, 2, 4)
		q := make([]int, n)
		var words []string
		var trace []string
		for i := 0; i < k; i++ {
			d := randBits(n, false)
			words = append(words, bits(d))
			q = l.step(q, d)
			trace = append(trace, bits(q))
		}
		ans := bits(q)
		return edakit.Numeric(stem+edakit.Sprintf(" It is used as a MISR: on each clock stage $i$ also XORs in response bit $d_i$. Starting from all zeros, it absorbs the response words $d_1 \\ldots d_%d$ = %s in that order. What is the signature? Answer with the %d bits $Q_1 \\ldots Q_%d$.", n, strings.Join(words, ", "), n, n),
			ans, "Each clock: shift with feedback, then XOR the word in. States: "+strings.Join(trace, ", ")+".", map[string]string{"signature": ans})
	default:
		k := edakit.Scaled(ctx, 3, 6)
		q := randBits(n, true)
		start := bits(q)
		var trace []string
		for i := 0; i < k; i++ {
			q = l.step(q, nil)
			trace = append(trace, bits(q))
		}
		ans := bits(q)
		return edakit.Numeric(stem+edakit.Sprintf(" Starting from %s, what is the state after %d clocks? Answer with the %d bits $Q_1 \\ldots Q_%d$.", start, k, n, n),
			ans, "States: "+strings.Join(trace, ", ")+".", map[string]string{"state": ans})
	}
}

func eulerPhi(n int) int {
	r := n
	for p := 2; p*p <= n; p++ {
		if n%p == 0 {
			for n%p == 0 {
				n /= p
			}
			r -= r / p
		}
	}
	if n > 1 {
		r -= r / n
	}
	return r
}

// lfsrPeriod: maximal period (KP 1), period by simulation (KP 2), primitive polynomial count
// (KP 3).
func lfsrPeriod(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 32)
		ans := (1 << n) - 1
		return edakit.Numeric(edakit.Sprintf("An LFSR with %d stages has a primitive characteristic polynomial and starts from a nonzero state. How many clocks pass before its state first repeats?", n),
			edakit.IntStr(ans), edakit.Sprintf("A primitive polynomial makes the LFSR visit every nonzero state once: $2^{%d} - 1 = %d$. The all-zero state is excluded because it maps to itself.", n, ans), nil)
	case 2:
		n := edakit.Scaled(ctx, 3, 5)
		l := randLFSR(n)
		seed := make([]int, n)
		seed[0] = 1
		p := l.period(seed)
		return edakit.Numeric(edakit.Sprintf("An external-XOR LFSR with %d stages shifts $Q_1 \\to \\ldots \\to Q_%d$ with $Q_1$ taking the feedback %s (characteristic polynomial %s). From the seed %s, how many clocks until the seed state recurs?", n, n, l.tapText(), l.poly(), bits(seed)),
			edakit.IntStr(p), edakit.Sprintf("Stepping the register returns to %s after %d clocks. A primitive polynomial would give $2^{%d}-1 = %d$; a shorter cycle means the polynomial is not primitive.", bits(seed), p, n, (1<<n)-1), map[string]string{"period": edakit.IntStr(p)})
	default:
		n := edakit.Scaled(ctx, 3, 16)
		m := (1 << n) - 1
		ph := eulerPhi(m)
		ans := ph / n
		return edakit.Numeric(edakit.Sprintf("How many primitive polynomials of degree %d are there over GF(2)?", n),
			edakit.IntStr(ans), edakit.Sprintf("The count is $\\varphi(2^{%d}-1)/%d = \\varphi(%d)/%d = %d/%d = %d$: each primitive element of GF($2^{%d}$) has %d conjugates sharing one minimal polynomial.", n, n, m, n, ph, n, ans, n, n), nil)
	}
}

// aliasing: 2^-n in ppm (KP 1), signature length for a target (KP 2), expected aliased parts
// (KP 3).
func aliasing(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 8, 24)
		p := math.Pow(2, -float64(n)) * 1e6
		ans := edakit.Num(p, 3)
		return edakit.Numeric(edakit.Sprintf("A long response stream is compressed into a MISR of %d bits. Assuming every faulty response is equally likely, what is the aliasing probability, in parts per million? Round to 3 decimal places.", n),
			ans, edakit.Sprintf("Aliasing approaches $2^{-%d}$ for long streams: %s ppm.", n, ans), nil)
	case 2:
		e := edakit.Int(4, 12)
		m := edakit.Pick([]int{1, 2, 5})
		target := float64(m) * math.Pow(10, -float64(e))
		n := int(math.Ceil(math.Log2(1 / target)))
		return edakit.Numeric(edakit.Sprintf("What is the smallest MISR length n for which the asymptotic aliasing probability $2^{-n}$ is at most $%d \\times 10^{-%d}$?", m, e),
			edakit.IntStr(n), edakit.Sprintf("Need $n \\ge \\log_2(1/(%d \\times 10^{-%d})) = %s$, so n = %d.", m, e, edakit.Num(math.Log2(1/target), 3), n), nil)
	default:
		n := edakit.Scaled(ctx, 8, 20)
		mdef := edakit.Scaled(ctx, 10000, 5000000)
		ans := edakit.Num(float64(mdef)*math.Pow(2, -float64(n)), 2)
		return edakit.Numeric(edakit.Sprintf("Logic BIST with a MISR of %d bits screens %d defective dies whose faults the patterns do reach the MISR. How many of them are expected to pass because their signature aliases to the good one? Round to 2 decimal places.", n, mdef),
			ans, edakit.Sprintf("Each aliases with probability $2^{-%d}$: $%d \\times 2^{-%d} = %s$.", n, mdef, n, ans), nil)
	}
}

// lbist: STUMPS session time (KP 1), observation-point benefit (KP 2), weighted random patterns
// (KP 3).
func lbist(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		p := edakit.Pick([]int{1000, 4096, 10000, 32768, 65536, 100000})
		l := edakit.Scaled(ctx, 50, 600)
		f := edakit.Pick([]int{50, 100, 200, 250})
		ms := float64(p*(l+1)) / (float64(f) * 1000)
		ans := edakit.Num(ms, 2)
		return edakit.Numeric(edakit.Sprintf("A STUMPS logic BIST controller applies %d pseudorandom patterns. Each pattern shifts %d cycles to fill the longest internal chain (the unload of the previous response into the MISR overlaps) and takes one capture cycle. At %d MHz, how long is the session, in milliseconds? Round to 2 decimal places.", p, l, f),
			ans, edakit.Sprintf("$%d \\times (%d + 1)$ cycles at %d MHz is %s ms.", p, l, f, ans), nil)
	case 2:
		e := edakit.Int(2, 8)
		pr := edakit.Int(3, 10)
		pe := math.Pow(2, -float64(e))
		pp := math.Pow(2, -float64(pr))
		n1 := int(math.Ceil(math.Log(0.05) / math.Log(1-pe*pp)))
		n2 := int(math.Ceil(math.Log(0.05) / math.Log(1-pe)))
		return edakit.Numeric(edakit.Sprintf("A fault is excited by a random pattern with probability $2^{-%d}$ and, once excited, reaches an observed point with probability $2^{-%d}$. An observation point inserted at the fault site makes propagation certain. How many fewer random patterns are needed to detect the fault with 95 percent probability?", e, pr),
			edakit.IntStr(n1-n2), edakit.Sprintf("Before: $p = 2^{-%d}$, $N = \\lceil \\ln 0.05/\\ln(1-p) \\rceil = %d$. After: $p = 2^{-%d}$, N = %d. Saved %d patterns.", e+pr, n1, e, n2, n1-n2), nil)
	default:
		k := edakit.Scaled(ctx, 4, 16)
		w := float64(edakit.Pick([]int{625, 750, 800, 875, 900, 9375})) / 1000
		if w > 1 {
			w /= 10
		}
		pw := math.Pow(w, float64(k))
		ans := edakit.Num(pw/math.Pow(0.5, float64(k)), 2)
		return edakit.Numeric(edakit.Sprintf("The output stuck-at-0 fault of a %d-input AND gate needs all inputs at 1. Weighted random pattern generation drives each input to 1 with probability %s instead of 0.5. By what factor does the per-pattern detection probability rise? Round to 2 decimal places.", k, edakit.Num(w, 4)),
			ans, edakit.Sprintf("$(%s/0.5)^{%d} = %s$. The same weighting makes the input stuck-at-1 faults, which need one input at 0, harder, so weight sets are tuned per region.", edakit.Num(w, 4), k, ans), nil)
	}
}
