package numer

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// circuit3 draws a 3x3 diagonally dominant nodal-like matrix, possibly unsymmetric (a
// controlled source), with integer entries.
func circuit3() [][]float64 {
	a := [][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i != j {
				a[i][j] = -float64(edakit.Int(1, 4))
			}
		}
	}
	if edakit.Int(0, 1) == 0 {
		a[1][0] = a[0][1]
		a[2][0] = a[0][2]
		a[2][1] = a[1][2]
	}
	for i := 0; i < 3; i++ {
		s := 0.0
		for j := 0; j < 3; j++ {
			if i != j {
				s += math.Abs(a[i][j])
			}
		}
		a[i][i] = s + float64(edakit.Int(1, 5))
	}
	return a
}

func lu(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		a := circuit3()
		i := edakit.Int(2, 3)
		l := a[i-1][0] / a[0][0]
		ans := f(l, 4)
		q := "Gaussian elimination without pivoting factors the nodal matrix " + pmat(a, 0) + sp(" as $A = LU$ with unit lower-triangular $L$. What is the multiplier $l_{%d1}$? Round to 4 decimal places.", i)
		return numeric(q, ans, sp("The multiplier is the entry being eliminated over the pivot: $l_{%d1} = a_{%d1}/a_{11} = %s/%s = %s$.", i, i, f(a[i-1][0], 0), f(a[0][0], 0), ans),
			map[string]string{"l": ans})
	case 2:
		a := circuit3()
		l21, l31 := a[1][0]/a[0][0], a[2][0]/a[0][0]
		u22 := a[1][1] - l21*a[0][1]
		u23 := a[1][2] - l21*a[0][2]
		a32 := a[2][1] - l31*a[0][1]
		l32 := a32 / u22
		u33 := a[2][2] - l31*a[0][2] - l32*u23
		ans := f(u33, 4)
		q := "Factor the nodal matrix " + pmat(a, 0) + " as $A = LU$ without pivoting. What is $u_{33}$, the last pivot? Round to 4 decimal places."
		return numeric(q, ans, sp("After step 1 the trailing block is the Schur complement: $u_{22} = %s$, $u_{23} = %s$, updated $a_{32} = %s$, so $l_{32} = %s$ and $u_{33} = a_{33} - l_{31}a_{13} - l_{32}u_{23} = %s$.", f(u22, 4), f(u23, 4), f(a32, 4), f(l32, 4), ans),
			map[string]string{"u22": f(u22, 4), "u33": ans})
	default:
		n := edakit.Scaled(ctx, 5, 7)
		g := randomGraph(n, 0.25)
		m := g.edgeCount()
		fill := fillForOrder(g, seq(1, n))
		ans := n + 2*m + 2*fill
		q := sp("The symbolic phase of a sparse LU analyses the pattern once before any numeric values are used. A structurally symmetric %d-by-%d matrix has a full diagonal and off-diagonal nonzero pairs at %s. With the natural elimination order 1, 2, ..., %d and no pivoting, how many entries do $L$ and $U$ store together? (Store the unit diagonal of $L$ implicitly, so count the diagonal once, in $U$.)",
			n, n, strings.Join(g.edges(), ", "), n)
		return numeric(q, itoa(ans), sp("The original pattern has %d + 2 · %d = %d nonzeros. Natural-order elimination adds %d fill edges, each a new entry in both $L$ and $U$: %d + 2 · %d = %d.", n, m, n+2*m, fill, n+2*m, fill, ans),
			map[string]string{"fill": itoa(fill), "nnz": itoa(n + 2*m)})
	}
}

func smallPivot(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		k := edakit.Scaled(ctx, 3, 8)
		d := edakit.Int(1, 9)
		a := edakit.Int(1, 9)
		ans := d * int(math.Pow(10, float64(k)))
		q := sp("A nodal matrix has a nearly floating node: $A = \\begin{pmatrix} 10^{-%d} & 1 \\\\ %d & %d \\end{pmatrix}$. Eliminating without pivoting, what is the multiplier $l_{21}$?", k, d, a)
		return numeric(q, itoa(ans), sp("$l_{21} = a_{21}/a_{11} = %d / 10^{-%d} = %d$, a huge multiplier from a tiny pivot.", d, k, ans), nil)
	case 2:
		k := edakit.Scaled(ctx, 3, 7)
		d := edakit.Int(1, 9)
		a := edakit.Int(1, 9)
		eps := math.Pow(10, -float64(k))
		u22 := float64(a) - float64(d)/eps
		maxA := math.Max(1, math.Max(float64(d), float64(a)))
		rho := math.Abs(u22) / maxA
		ans := f(rho, 0)
		q := sp("Factor $A = \\begin{pmatrix} 10^{-%d} & 1 \\\\ %d & %d \\end{pmatrix}$ without pivoting. The growth factor is the largest entry magnitude in $U$ divided by the largest in $A$. What is it? Round to the nearest integer.", k, d, a)
		return numeric(q, ans, sp("$u_{22} = %d - %d \\cdot 10^{%d} = %s$, so the growth factor is $%s / %s = %s$.", a, d, k, f(u22, 0), f(math.Abs(u22), 0), f(maxA, 0), ans),
			map[string]string{"u22": f(u22, 0)})
	default:
		k := edakit.Int(2, 4)
		d := edakit.Int(2, 9)
		a := edakit.Int(1, 9)
		eps := math.Pow(10, -float64(k))
		u22 := 1 - eps/float64(d)*float64(a)
		ans := f(u22, 6)
		q := sp("Factor $A = \\begin{pmatrix} 10^{-%d} & 1 \\\\ %d & %d \\end{pmatrix}$ with partial pivoting, which swaps the rows so the pivot is %d. What is $u_{22}$ of the pivoted factorization? Round to 6 decimal places.", k, d, a, d)
		return numeric(q, ans, sp("After the swap the multiplier is $10^{-%d}/%d$, at most 1 in magnitude, and $u_{22} = 1 - (10^{-%d}/%d) \\cdot %d = %s$. No entry grows.", k, d, k, d, a, ans),
			map[string]string{"u22": ans})
	}
}

func thresholdPivoting(ctx generator.GeneratorContext) generator.Problem {
	m := edakit.Scaled(ctx, 4, 6)
	vals := make([]float64, m)
	rows := make([]int, m)
	for i := range vals {
		v := float64(edakit.Int(1, 99)) / 10
		if edakit.Int(0, 1) == 0 {
			v = -v
		}
		vals[i] = v
		rows[i] = edakit.Int(2, 9)
	}
	u := edakit.Pick([]float64{0.1, 0.25, 0.5, 0.01})
	maxAbs := 0.0
	for _, v := range vals {
		maxAbs = math.Max(maxAbs, math.Abs(v))
	}
	var desc []string
	for i := range vals {
		desc = append(desc, sp("row %d: value %s, %d nonzeros in its row", i+1, f(vals[i], 1), rows[i]))
	}
	stem := "The active pivot column of a sparse LU has these candidates: " + strings.Join(desc, "; ") + ". "
	switch kpOf(ctx) {
	case 1:
		cnt := 0
		for _, v := range vals {
			if math.Abs(v) >= u*maxAbs {
				cnt++
			}
		}
		q := stem + sp("With threshold $u = %s$, how many candidates pass the test $|a_{ij}| \\ge u \\max_k |a_{kj}|$?", f(u, 2))
		return numeric(q, itoa(cnt), sp("The largest magnitude is %s, so the cutoff is %s. %d candidates reach it.", f(maxAbs, 1), f(u*maxAbs, 3), cnt),
			map[string]string{"cutoff": f(u*maxAbs, 3)})
	case 2:
		best, br := -1, 1<<30
		for i, v := range vals {
			if math.Abs(v) >= u*maxAbs {
				if rows[i] < br || (rows[i] == br && math.Abs(v) > math.Abs(vals[best])) {
					best, br = i, rows[i]
				}
			}
		}
		q := stem + sp("With threshold $u = %s$, the solver picks the eligible candidate with the least Markowitz count (fewest row nonzeros, since all share this column), breaking ties by larger magnitude. Which row is chosen?", f(u, 2))
		return numeric(q, itoa(best+1), sp("Eligible candidates have magnitude at least %s. Among them, row %d has the fewest row nonzeros (%d).", f(u*maxAbs, 3), best+1, br),
			map[string]string{"row": itoa(best + 1)})
	default:
		if edakit.Int(0, 1) == 0 {
			i := edakit.Int(0, m-1)
			ans := f(math.Abs(vals[i])/maxAbs, 4)
			q := stem + sp("What is the largest threshold $u$ for which the row %d candidate is still eligible? Round to 4 decimal places.", i+1)
			return numeric(q, ans, sp("Row %d is eligible while $u \\le |a|/\\max|a| = %s/%s = %s$.", i+1, f(math.Abs(vals[i]), 1), f(maxAbs, 1), ans), nil)
		}
		k := edakit.Int(2, 4)
		g := math.Pow(1+1/u, float64(k))
		ans := f(g, 2)
		q := sp("Threshold pivoting with $u = %s$ bounds every multiplier by $1/u$, so each elimination step can grow entries by at most a factor $1 + 1/u$. What is the worst-case growth bound after %d steps? Round to 2 decimal places.", f(u, 2), k)
		return numeric(q, ans, sp("$(1 + 1/u)^{%d} = %s^{%d} = %s$. Smaller $u$ favours sparsity but weakens this bound.", k, f(1+1/u, 2), k, ans), nil)
	}
}

func condition(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		g1, g2, g3 := float64(edakit.Int(1, 20)), float64(edakit.Int(1, 20)), float64(edakit.Int(1, 20))
		if edakit.Int(0, 1) == 0 {
			g3 = float64(edakit.Int(1, 9)) / 100
		}
		a := [][]float64{{g1 + g2, -g2}, {-g2, g2 + g3}}
		det := a[0][0]*a[1][1] - a[0][1]*a[1][0]
		inv := [][]float64{{a[1][1] / det, -a[0][1] / det}, {-a[1][0] / det, a[0][0] / det}}
		nA := math.Max(math.Abs(a[0][0])+math.Abs(a[0][1]), math.Abs(a[1][0])+math.Abs(a[1][1]))
		nI := math.Max(math.Abs(inv[0][0])+math.Abs(inv[0][1]), math.Abs(inv[1][0])+math.Abs(inv[1][1]))
		k := nA * nI
		ans := f(k, 2)
		q := sp("Node 1 connects to ground through %s mS and to node 2 through %s mS; node 2 connects to ground through %s mS. Its conductance matrix is ", f(g1, 2), f(g2, 2), f(g3, 2)) + pmat(a, 2) + " mS. What is the condition number $\\kappa_\\infty(G) = \\|G\\|_\\infty \\|G^{-1}\\|_\\infty$? Round to 2 decimal places."
		return numeric(q, ans, sp("$\\|G\\|_\\infty = %s$ (largest row sum), $\\det G = %s$, $\\|G^{-1}\\|_\\infty = %s$, so $\\kappa_\\infty = %s$.", f(nA, 3), f(det, 4), f(nI, 6), ans),
			map[string]string{"normA": f(nA, 3), "normInv": f(nI, 6)})
	case 2:
		km := float64(edakit.Int(11, 99)) / 10
		ke := edakit.Int(3, 7)
		rm := float64(edakit.Int(11, 99)) / 10
		re := -(ke + edakit.Int(2, 5))
		ks, kv := sci(km, ke)
		rs, rv := sci(rm, re)
		bound := kv * rv
		ans := f(bound*1e6, 2)
		q := sp("A nodal solve has $\\kappa(G) = $ %s, and the computed solution has relative residual $\\|b - G\\hat x\\|/\\|b\\| = $ %s. What upper bound does this give on the relative error $\\|x - \\hat x\\|/\\|x\\|$, in parts per million? Round to 2 decimal places.", ks, rs)
		return numeric(q, ans, sp("The relative error is at most $\\kappa$ times the relative residual: %s ppm.", ans),
			map[string]string{"bound": f(bound, 12)})
	default:
		bm := float64(edakit.Int(1, 9))
		be := edakit.Int(0, 3)
		sm := float64(edakit.Int(1, 9))
		se := -edakit.Int(6, 12)
		bs, bv := sci(bm, be)
		ss, sv := sci(sm, se)
		k := bv / sv
		digits := 16 - math.Log10(k)
		ans := f(digits, 1)
		q := sp("A circuit matrix is, after scaling, close to diagonal with its largest conductance %s S (a switch) and its smallest %s S (a gmin leak), so $\\kappa \\approx g_{max}/g_{min}$. Taking double precision as about 16 significant digits and the solve as losing $\\log_{10}\\kappa$ of them, about how many significant digits remain in the solution? Round to 1 decimal place.", bs, ss)
		return numeric(q, ans, sp("$\\kappa \\approx %s$, $\\log_{10}\\kappa = %s$, so about $16 - %s = %s$ digits remain.", f(k, 0), f(math.Log10(k), 2), f(math.Log10(k), 2), ans),
			map[string]string{"log10_kappa": f(math.Log10(k), 3)})
	}
}

func cholesky(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		g12, g1, g2 := float64(edakit.Int(1, 9)), float64(edakit.Int(1, 9)), float64(edakit.Int(1, 9))
		v1, v2 := float64(edakit.Int(-20, 20))/10, float64(edakit.Int(-20, 20))/10
		p := g12*(v1-v2)*(v1-v2) + g1*v1*v1 + g2*v2*v2
		ans := f(p, 3)
		a := [][]float64{{g12 + g1, -g12}, {-g12, g12 + g2}}
		q := sp("Two nodes are joined by %s mS; node 1 has %s mS and node 2 has %s mS to ground, so $G = $ ", f(g12, 0), f(g1, 0), f(g2, 0)) + pmat(a, 0) + sp(" mS. For $v = (%s, %s)$ V, compute $v^T G v$ in mW (it equals the power the resistors dissipate). Round to 3 decimal places.", f(v1, 1), f(v2, 1))
		return numeric(q, ans, sp("$v^T G v = g_{12}(v_1 - v_2)^2 + g_1 v_1^2 + g_2 v_2^2 = %s$ mW, positive for every nonzero $v$, which is why $G$ is positive definite.", ans),
			map[string]string{"p": ans})
	case 2:
		if edakit.Int(0, 1) == 0 {
			a21 := -float64(edakit.Int(1, 5))
			a11 := math.Abs(a21) + float64(edakit.Int(1, 8))
			a22 := math.Abs(a21) + float64(edakit.Int(1, 8))
			l11 := math.Sqrt(a11)
			l21 := a21 / l11
			l22 := math.Sqrt(a22 - l21*l21)
			ans := f(l22, 4)
			a := [][]float64{{a11, a21}, {a21, a22}}
			q := "Compute the Cholesky factor $G = LL^T$ of " + pmat(a, 0) + ". What is $l_{22}$? Round to 4 decimal places."
			return numeric(q, ans, sp("$l_{11} = \\sqrt{%s} = %s$, $l_{21} = %s/l_{11} = %s$, $l_{22} = \\sqrt{%s - l_{21}^2} = %s$.", f(a11, 0), f(l11, 4), f(a21, 0), f(l21, 4), f(a22, 0), ans),
				map[string]string{"l11": f(l11, 4), "l21": f(l21, 4)})
		}
		a := circuit3()
		a[1][0], a[2][0], a[2][1] = a[0][1], a[0][2], a[1][2]
		for i := 0; i < 3; i++ {
			a[i][i] = math.Abs(a[i][(i+1)%3]) + math.Abs(a[i][(i+2)%3]) + float64(edakit.Int(1, 5))
		}
		l11 := math.Sqrt(a[0][0])
		l21, l31 := a[1][0]/l11, a[2][0]/l11
		l22 := math.Sqrt(a[1][1] - l21*l21)
		l32 := (a[2][1] - l31*l21) / l22
		l33 := math.Sqrt(a[2][2] - l31*l31 - l32*l32)
		ans := f(l33, 4)
		q := "Compute the Cholesky factor $G = LL^T$ of the conductance matrix " + pmat(a, 0) + ". What is $l_{33}$? Round to 4 decimal places."
		return numeric(q, ans, sp("$l_{11} = %s$, $l_{21} = %s$, $l_{31} = %s$, $l_{22} = %s$, $l_{32} = %s$, and $l_{33} = \\sqrt{a_{33} - l_{31}^2 - l_{32}^2} = %s$.", f(l11, 4), f(l21, 4), f(l31, 4), f(l22, 4), f(l32, 4), ans),
			map[string]string{"l22": f(l22, 4)})
	default:
		n := 100 * edakit.Scaled(ctx, 5, 80)
		rate := float64(edakit.Pick([]int{50, 100, 200, 500, 1000}))
		flops := math.Pow(float64(n), 3) / 3
		ms := flops / (rate * 1e9) * 1000
		ans := f(ms, 2)
		q := sp("A dense block of a thermal model has n = %d unknowns. Dense Cholesky costs about $n^3/3$ flops (LU about $2n^3/3$). On a processor sustaining %s GFLOP/s, how long does the Cholesky factorization take, in milliseconds? Round to 2 decimal places.", n, f(rate, 0))
		return numeric(q, ans, sp("$n^3/3 = %s$ flops at %s GFLOP/s is %s ms; LU would take twice as long.", f(flops, 0), f(rate, 0), ans), nil)
	}
}
