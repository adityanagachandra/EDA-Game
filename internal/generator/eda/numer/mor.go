package numer

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// chainMoments returns m0, m1, m2 at both nodes of a two-node RC chain driven by a unit
// voltage source through R1 (node 1, C1) and R2 (node 2, C2). Units kohm, pF, ns.
func chainMoments(r1, c1, r2, c2 float64) (m1, m2 [2]float64) {
	g := [][]float64{{1/r1 + 1/r2, -1 / r2}, {-1 / r2, 1 / r2}}
	m0 := []float64{1, 1}
	x := solve2(g, []float64{c1 * m0[0], c2 * m0[1]})
	m1 = [2]float64{-x[0], -x[1]}
	y := solve2(g, []float64{c1 * m1[0], c2 * m1[1]})
	m2 = [2]float64{-y[0], -y[1]}
	return
}

func moments(ctx generator.GeneratorContext) generator.Problem {
	r1, c1 := float64(edakit.Int(1, 10)), float64(edakit.Int(1, 10))/10
	r2, c2 := float64(edakit.Int(1, 10)), float64(edakit.Int(1, 10))/10
	chain := sp("A unit voltage source drives $R_1 = %s$ kΩ into node 1 ($C_1 = %s$ pF to ground), then $R_2 = %s$ kΩ into node 2 ($C_2 = %s$ pF to ground). ", f(r1, 0), f(c1, 1), f(r2, 0), f(c2, 1))
	m1, m2 := chainMoments(r1, c1, r2, c2)
	switch kpOf(ctx) {
	case 1:
		k := edakit.Int(1, 3)
		tau := r1 * c1
		m := math.Pow(-tau, float64(k))
		ans := f(m, 4)
		q := sp("A single RC section has $H(s) = 1/(1 + s\\tau)$ with $R = %s$ kΩ and $C = %s$ pF. Expanding $H(s) = \\sum_k m_k s^k$, what is the moment $m_%d$, in ns$^%d$? Round to 4 decimal places.", f(r1, 0), f(c1, 1), k, k)
		return numeric(q, ans, sp("$1/(1 + s\\tau) = \\sum (-\\tau)^k s^k$, so $m_%d = (-%s)^%d = %s$.", k, f(tau, 1), k, ans), nil)
	case 2:
		node := edakit.Int(1, 2)
		ans := f(m1[node-1], 4)
		q := chain + sp("The moments satisfy $Gm_0 = b$ and $Gm_k = -Cm_{k-1}$. What is the first moment $m_1$ of the voltage at node %d, in ns? Round to 4 decimal places.", node)
		return numeric(q, ans, sp("$m_0 = (1, 1)$; solving $Gm_1 = -Cm_0$ gives $m_1 = (%s, %s)$. Its negative is the Elmore delay: $R_1(C_1 + C_2)%s$.", f(m1[0], 4), f(m1[1], 4), map[int]string{1: "", 2: " + R_2C_2"}[node]),
			map[string]string{"m1_1": f(m1[0], 4), "m1_2": f(m1[1], 4)})
	default:
		tau := -m2[1] / m1[1]
		ans := f(tau, 4)
		q := chain + sp("At node 2 the moments are $m_0 = 1$, $m_1 = %s$ ns and $m_2 = %s$ ns$^2$. Fit a one-pole model $H(s) \\approx 1/(1 + s\\tau)$ by matching $m_1$ and $m_2$ (so $\\tau$ comes from their ratio). What is $\\tau$, in ns? Round to 4 decimal places.", f(m1[1], 4), f(m2[1], 4))
		return numeric(q, ans, sp("For one pole $m_1 = -\\tau$ and $m_2 = \\tau^2$, so $\\tau = -m_2/m_1 = %s$ ns. Matching $m_0, m_1$ alone would give the Elmore value %s ns instead.", ans, f(-m1[1], 4)),
			map[string]string{"m1": f(m1[1], 4), "m2": f(m2[1], 4)})
	}
}

var krylovItems = []edakit.Item{
	{KP: 3, Q: "AWE computes the moment vectors $m_k$ explicitly and fits a Pade approximant. Why does it become unreliable beyond about eight poles?",
		A: "The moment vectors align with the dominant eigenvector, like the power method, so the moment matrix is numerically singular", D: []string{"Higher moments are always zero for RC circuits", "Pade approximants cannot have more than eight poles", "The LU factorization of G changes with each moment"},
		Why: "Repeated multiplication by the same matrix drives vectors toward the dominant direction; Arnoldi orthogonalizes the basis to keep it well conditioned."},
	{KP: 3, Q: "PRIMA projects $G$ and $C$ by congruence, $G_r = V^TGV$ and $C_r = V^TCV$. What does this preserve that AWE does not?",
		A: "Positive semidefiniteness of G and C, so the reduced RLC model stays passive", D: []string{"The exact poles of the full model", "The sparsity of G and C", "All moments to infinite order"},
		Why: "A congruence transform keeps semidefinite matrices semidefinite, which guarantees passivity of the reduced model."},
	{KP: 3, Q: "A reduced interconnect model is accurate below 1 GHz but wrong at 20 GHz. Which change targets the problem?",
		A: "Add Krylov vectors at a second expansion point near the high frequency (rational Krylov)", D: []string{"Add more vectors at s = 0 only", "Use a one-sided instead of a two-sided projection", "Drop the capacitance matrix"},
		Why: "Moments at s = 0 describe low-frequency behaviour; expansion points across the band give broadband accuracy."},
	{KP: 3, Q: "Why must a non-passive reduced model never be handed to a transient simulator, even if it matches the frequency response well?",
		A: "Connected to other passive parts it can generate energy and make the overall simulation unstable", D: []string{"The simulator cannot read reduced models", "It will have too many ports", "Non-passive models are always less accurate"},
		Why: "Passivity guarantees stability for any passive termination; an active-looking reduced model can oscillate or blow up in a larger circuit."},
	{KP: 3, Q: "What is the main cost of building a PRIMA model with q vectors for an n-node RC network?",
		A: "One sparse factorization of G, then about q sparse solves and orthogonalizations", D: []string{"A dense eigendecomposition of the n-by-n system", "q separate factorizations of G", "An n-by-n singular value decomposition"},
		Why: "Each Krylov vector needs one solve with the same factored G, so the cost is near linear in n for fixed q."},
	{KP: 3, Q: "A multiport interconnect has p = 8 ports, and block Arnoldi builds $k$ block iterations. How large is the reduced model?",
		A: "About kp = 8k states, since each block iteration adds one vector per port", D: []string{"k states regardless of ports", "p states regardless of k", "k^p states"},
		Why: "Block Krylov methods advance all port excitations together, so reduced size grows with the port count; many-port models need care."},
}

func krylovVector(ctx generator.GeneratorContext) generator.Problem {
	a := [][]float64{{float64(edakit.Int(-4, 4)), float64(edakit.Int(-4, 4))}, {float64(edakit.Int(-4, 4)), float64(edakit.Int(-4, 4))}}
	if a[0][1] == a[1][0] {
		a[0][1]++
	}
	r := []float64{float64(edakit.Int(-3, 3)), float64(edakit.Int(1, 3))}
	ar := matvec(a, r)
	at := [][]float64{{a[0][0], a[1][0]}, {a[0][1], a[1][1]}}
	atr := matvec(at, r)
	vs := func(v []float64) string { return "$(" + f(v[0], 0) + ", " + f(v[1], 0) + ")$" }
	q := "A Krylov basis for reduction starts from $r = $ " + pvec(r, 0) + " with $A = $ " + pmat(a, 0) + ". The subspace $K_2(A, r) = \\mathrm{span}\\{r, Ar\\}$. What is the second Krylov vector $Ar$ (before orthogonalization)?"
	return edakit.Choice(q, vs(ar), sp("Multiply row by row: $Ar = (%s, %s)$.", f(ar[0], 0), f(ar[1], 0)),
		vs(atr), vs([]float64{-ar[0], -ar[1]}), vs([]float64{ar[1], ar[0]}), vs([]float64{ar[0] + r[0], ar[1] + r[1]}))
}

func krylovMoments(ctx generator.GeneratorContext) generator.Problem {
	q0 := edakit.Int(2, 12)
	switch edakit.Int(0, 2) {
	case 0:
		q := sp("An Arnoldi-based (one-sided, PRIMA-style) reduction of a single-input single-output RC network uses an orthonormal basis of $K_{%d}(A, r)$ at the expansion point $s = 0$. How many moments of the transfer function does the reduced model match?", q0)
		return edakit.Choice(q, itoa(q0), sp("A one-sided projection onto a q-dimensional Krylov space matches q moments: %d.", q0),
			edakit.Distinct(q0, 3, 2*q0, q0-1, q0+1, 2*q0-1)...)
	case 1:
		q := sp("A two-sided (Lanczos, Pade-via-Lanczos) reduction builds left and right Krylov bases of dimension %d each at $s = 0$ for a single-input single-output model. How many moments does it match?", q0)
		return edakit.Choice(q, itoa(2*q0), sp("Two-sided projection gives a Pade approximant matching 2q moments: %d.", 2*q0),
			edakit.Distinct(2*q0, 3, q0, q0*q0, 2*q0+1, q0+1)...)
	default:
		a, b := edakit.Int(2, 6), edakit.Int(2, 6)
		q := sp("A rational Krylov reduction (one-sided) uses %d vectors at expansion point $s_1$ and %d vectors at $s_2$ for a single-input single-output model. How many moments does it match in total, counting moments about both points?", a, b)
		return edakit.Choice(q, itoa(a+b), sp("One-sided projection matches as many moments at each point as it has vectors there: %d + %d = %d.", a, b, a+b),
			edakit.Distinct(a+b, 3, 2*(a+b), a*b, a+b-1, maxInt(a, b))...)
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func svdTruncation(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 8)
	sv := make([]float64, n)
	v := float64(edakit.Int(50, 200)) / 10
	for i := range sv {
		sv[i] = edakit.Round(v, 3)
		v *= float64(edakit.Int(15, 70)) / 100
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(sv)))
	var ss []string
	for _, x := range sv {
		ss = append(ss, f(x, 3))
	}
	stem := "A snapshot matrix of a thermal model has singular values " + strings.Join(ss, ", ") + ". "
	switch kpOf(ctx) {
	case 1:
		k := edakit.Int(1, n-2)
		tot, top := 0.0, 0.0
		for i, x := range sv {
			tot += x * x
			if i < k {
				top += x * x
			}
		}
		ans := f(100*top/tot, 2)
		q := stem + sp("What percentage of the total energy $\\sum\\sigma_i^2$ do the first %d singular directions capture? Round to 2 decimal places.", k)
		return numeric(q, ans, sp("$\\sum_{i \\le %d}\\sigma_i^2 = %s$ of %s: %s percent.", k, f(top, 4), f(tot, 4), ans), nil)
	case 2:
		k := edakit.Int(1, n-2)
		if edakit.Int(0, 1) == 0 {
			ans := f(sv[k]/sv[0], 4)
			q := stem + sp("Truncating to rank %d, what is the relative 2-norm error $\\|A - A_%d\\|_2/\\|A\\|_2$? Round to 4 decimal places.", k, k)
			return numeric(q, ans, sp("Eckart-Young: the error is $\\sigma_{%d} = %s$ and $\\|A\\|_2 = \\sigma_1 = %s$, ratio %s.", k+1, f(sv[k], 3), f(sv[0], 3), ans), nil)
		}
		s := 0.0
		for _, x := range sv[k:] {
			s += x * x
		}
		ans := f(math.Sqrt(s), 4)
		q := stem + sp("Truncating to rank %d, what is the Frobenius-norm error $\\|A - A_%d\\|_F$? Round to 4 decimal places.", k, k)
		return numeric(q, ans, sp("$\\sqrt{\\sum_{i > %d}\\sigma_i^2} = %s$.", k, ans), nil)
	default:
		var tail []float64
		for i := range sv {
			t := 0.0
			for _, x := range sv[i:] {
				t += x
			}
			tail = append(tail, t)
		}
		k := edakit.Int(1, n-1)
		tol := 2*tail[k] + 0.001
		tol = edakit.Round(tol*float64(edakit.Int(100, 130))/100, 3)
		best := n
		for kk := 0; kk <= n; kk++ {
			t := 0.0
			if kk < n {
				t = tail[kk]
			}
			if 2*t <= tol {
				best = kk
				break
			}
		}
		q := "Balanced truncation of a stable thermal model has Hankel singular values " + strings.Join(ss, ", ") + sp(". Keeping $k$ states guarantees an $H_\\infty$ error at most $2\\sum_{i > k}\\sigma_i$. What is the smallest $k$ that guarantees an error of at most %s?", f(tol, 3))
		return numeric(q, itoa(best), sp("Tail sums $2\\sum_{i>k}\\sigma_i$ for k = 0, 1, ...: the first at or below %s is at $k = %d$.", f(tol, 3), best), nil)
	}
}
