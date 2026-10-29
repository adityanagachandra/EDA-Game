package numer

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func jacobiGS(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Int(2, 4)
		var parts []string
		num, den := 0.0, 0.0
		for i := 0; i < n; i++ {
			g := float64(edakit.Int(1, 20))
			v := float64(edakit.Int(900, 1000)) / 1000
			num += g * v
			den += g
			parts = append(parts, sp("%s S to a node at %s V", f(g, 0), f(v, 3)))
		}
		load := float64(edakit.Int(1, 40)) / 10
		vi := (num - load) / den
		ans := f(vi, 4)
		q := "A power-grid node connects through " + strings.Join(parts, ", ") + sp(", and its cells draw %s A. Holding the neighbours fixed, what voltage does one Jacobi update assign to the node, in volts? Round to 4 decimal places.", f(load, 1))
		return numeric(q, ans, sp("KCL: $\\sum g_j(v_j - v_i) = I_{load}$, so $v_i = (\\sum g_j v_j - I_{load})/\\sum g_j = (%s - %s)/%s = %s$ V.", f(num, 4), f(load, 1), f(den, 0), ans),
			map[string]string{"sum_gv": f(num, 4), "sum_g": f(den, 0)})
	case 2:
		vl := 1.0
		vr := float64(edakit.Int(950, 1000)) / 1000
		g := float64(edakit.Int(5, 20))
		var loads [3]float64
		var init [3]float64
		for i := range loads {
			loads[i] = float64(edakit.Int(1, 20)) / 100
			init[i] = float64(edakit.Int(900, 1000)) / 1000
		}
		gs := edakit.Int(0, 1) == 0
		v := init
		if gs {
			for i := 0; i < 3; i++ {
				left, right := vl, vr
				if i > 0 {
					left = v[i-1]
				}
				if i < 2 {
					right = v[i+1]
				}
				v[i] = (left + right - loads[i]/g) / 2
			}
		} else {
			old := v
			for i := 0; i < 3; i++ {
				left, right := vl, vr
				if i > 0 {
					left = old[i-1]
				}
				if i < 2 {
					right = old[i+1]
				}
				v[i] = (left + right - loads[i]/g) / 2
			}
		}
		k := edakit.Int(2, 3)
		if !gs {
			k = edakit.Int(1, 3)
		}
		ans := f(v[k-1], 4)
		method := "Jacobi (all updates use the old values)"
		if gs {
			method = "Gauss-Seidel in the order 1, 2, 3 (each update uses the newest values)"
		}
		q := sp("A supply rail is a chain: pad (%s V), node 1, node 2, node 3, pad (%s V), with %s S between neighbours. Nodes 1 to 3 draw %s, %s and %s A and start at %s, %s and %s V. After one sweep of %s, what is the voltage of node %d? Round to 4 decimal places.",
			f(vl, 3), f(vr, 3), f(g, 0), f(loads[0], 2), f(loads[1], 2), f(loads[2], 2), f(init[0], 3), f(init[1], 3), f(init[2], 3), method, k)
		return numeric(q, ans, sp("Each update is $v_i = (v_{left} + v_{right} - I_i/g)/2$. The sweep gives $(%s, %s, %s)$ V.", f(v[0], 4), f(v[1], 4), f(v[2], 4)),
			map[string]string{"v1": f(v[0], 4), "v2": f(v[1], 4), "v3": f(v[2], 4)})
	default:
		n := edakit.Scaled(ctx, 10, 200)
		d := edakit.Int(2, 6)
		rho := math.Cos(math.Pi / float64(n+1))
		gs := edakit.Int(0, 1) == 0
		r := rho
		name := "Jacobi"
		if gs {
			r = rho * rho
			name = "Gauss-Seidel"
		}
		it := math.Ceil(float64(d) * math.Ln10 / -math.Log(r))
		ans := f(it, 0)
		q := sp("A rail of %d equal resistor segments has %d interior nodes between two fixed pads. The Jacobi iteration matrix has spectral radius $\\rho = \\cos(\\pi/(n+1))$ with $n = %d$, and Gauss-Seidel has $\\rho^2$. How many %s iterations reduce the error by a factor $10^{%d}$?", n+1, n, n, name, d)
		return numeric(q, ans, sp("$\\rho_{%s} = %s$. Iterations $= \\lceil %d \\ln 10 / (-\\ln \\rho) \\rceil = %s$. Long chains converge very slowly.", name, f(r, 6), d, ans),
			map[string]string{"rho": f(r, 6)})
	}
}

// spd2 draws a 2x2 symmetric positive definite integer matrix and an integer right-hand side.
func spd2() ([][]float64, []float64) {
	a := float64(edakit.Int(2, 9))
	d := float64(edakit.Int(2, 9))
	c := float64(edakit.Int(-3, 3))
	for c*c >= a*d {
		c /= 2
		c = math.Trunc(c)
	}
	b := []float64{float64(edakit.Int(-5, 6)), float64(edakit.Int(1, 6))}
	return [][]float64{{a, c}, {c, d}}, b
}

func conjugateGradient(ctx generator.GeneratorContext) generator.Problem {
	a, b := spd2()
	stem := "Solve $Ax = b$ with $A = $ " + pmat(a, 0) + " and $b = $ " + pvec(b, 0) + " by conjugate gradients from $x_0 = 0$. "
	switch kpOf(ctx) {
	case 1:
		x := []float64{float64(edakit.Int(-3, 3)), float64(edakit.Int(-3, 3))}
		phi := 0.5*dot(x, matvec(a, x)) - dot(b, x)
		ans := f(phi, 3)
		q := "For $A = $ " + pmat(a, 0) + " and $b = $ " + pvec(b, 0) + sp(", CG minimizes the energy $\\phi(x) = \\tfrac12 x^TAx - b^Tx$. What is $\\phi$ at $x = (%s, %s)$? Round to 3 decimal places.", f(x[0], 0), f(x[1], 0))
		return numeric(q, ans, sp("$x^TAx = %s$ and $b^Tx = %s$, so $\\phi = %s$.", f(dot(x, matvec(a, x)), 3), f(dot(b, x), 3), ans), nil)
	case 2:
		r := b
		ar := matvec(a, r)
		alpha := dot(r, r) / dot(r, ar)
		x1 := []float64{alpha * r[0], alpha * r[1]}
		r1 := []float64{r[0] - alpha*ar[0], r[1] - alpha*ar[1]}
		switch edakit.Int(0, 2) {
		case 0:
			ans := f(alpha, 4)
			return numeric(stem+"What is the first step length $\\alpha_0 = r_0^Tr_0/(p_0^TAp_0)$? Round to 4 decimal places.", ans,
				sp("$r_0 = p_0 = b$, $r_0^Tr_0 = %s$, $p_0^TAp_0 = %s$, so $\\alpha_0 = %s$.", f(dot(r, r), 0), f(dot(r, ar), 0), ans), nil)
		case 1:
			i := edakit.Int(0, 1)
			ans := f(x1[i], 4)
			return numeric(stem+sp("What is component %d of the first iterate $x_1$? Round to 4 decimal places.", i+1), ans,
				sp("$\\alpha_0 = %s$ and $x_1 = \\alpha_0 b = (%s, %s)$.", f(alpha, 4), f(x1[0], 4), f(x1[1], 4)), map[string]string{"alpha": f(alpha, 6)})
		default:
			i := edakit.Int(0, 1)
			ans := f(r1[i], 4)
			return numeric(stem+sp("What is component %d of the new residual $r_1 = r_0 - \\alpha_0 Ap_0$? Round to 4 decimal places.", i+1), ans,
				sp("$Ap_0 = (%s, %s)$ and $\\alpha_0 = %s$, so $r_1 = (%s, %s)$, orthogonal to $r_0$.", f(ar[0], 0), f(ar[1], 0), f(alpha, 4), f(r1[0], 4), f(r1[1], 4)), map[string]string{"alpha": f(alpha, 6)})
		}
	default:
		n := edakit.Scaled(ctx, 6, 14)
		pool := []int{1, 2, 3, 4, 5, 6, 8}
		k := edakit.Int(2, 5)
		vals := edakit.Shuffle(pool)[:k]
		var eig []int
		for i := 0; i < n; i++ {
			eig = append(eig, vals[i%k])
		}
		eig = edakit.Shuffle(eig)
		q := sp("A symmetric positive definite %d-by-%d matrix (a regular grid with repeated symmetric blocks) has eigenvalues %s. In exact arithmetic, at most how many CG iterations reach the exact solution from any starting point?", n, n, ints(eig))
		return numeric(q, itoa(k), sp("CG terminates in at most as many iterations as there are distinct eigenvalues, here %d, because a degree-%d polynomial can vanish on all of them.", k, k), nil)
	}
}

func cgConvergence(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		km := float64(edakit.Int(1, 9))
		ke := edakit.Int(1, 6)
		ks, kv := sci(km, ke)
		rho := (math.Sqrt(kv) - 1) / (math.Sqrt(kv) + 1)
		ans := f(rho, 4)
		q := "The CG error bound is $\\|e_k\\|_A \\le 2\\rho^k\\|e_0\\|_A$ with $\\rho = (\\sqrt\\kappa - 1)/(\\sqrt\\kappa + 1)$. For $\\kappa = $ " + ks + ", what is $\\rho$? Round to 4 decimal places."
		return numeric(q, ans, sp("$\\sqrt\\kappa = %s$, so $\\rho = %s$. Compare $(\\kappa - 1)/(\\kappa + 1)$ for steepest descent, much closer to 1.", f(math.Sqrt(kv), 3), ans), nil)
	case 2:
		km := float64(edakit.Int(1, 9))
		ke := edakit.Int(2, 6)
		ks, kv := sci(km, ke)
		d := edakit.Int(3, 10)
		rho := (math.Sqrt(kv) - 1) / (math.Sqrt(kv) + 1)
		k := math.Ceil(math.Log(math.Pow(10, -float64(d))/2) / math.Log(rho))
		ans := f(k, 0)
		q := "An IR-drop matrix has $\\kappa = $ " + ks + sp(". Using $\\|e_k\\|_A \\le 2\\rho^k\\|e_0\\|_A$, $\\rho = (\\sqrt\\kappa - 1)/(\\sqrt\\kappa + 1)$, what is the smallest $k$ that guarantees a reduction by $10^{-%d}$?", d)
		return numeric(q, ans, sp("$\\rho = %s$; need $2\\rho^k \\le 10^{-%d}$, so $k \\ge \\ln(10^{-%d}/2)/\\ln\\rho$, giving %s.", f(rho, 6), d, d, ans),
			map[string]string{"rho": f(rho, 6)})
	default:
		k := edakit.Int(1, 10) * 100
		it := edakit.Int(5, 40) * 10
		r := edakit.Int(2, 4)
		ans := it * r
		q := sp("CG on the nodal matrix of a %d-by-%d power grid converges in %d iterations. The grid is refined by a factor %d in each direction (same physical size, %d times the nodes). The condition number grows like $h^{-2}$ and CG iterations like $\\sqrt\\kappa$. About how many iterations will the refined grid need?", k, k, it, r, r*r)
		return numeric(q, itoa(ans), sp("$\\kappa$ grows by $%d^2$, so $\\sqrt\\kappa$ and the iteration count grow by %d: about %d iterations, each on %d times as many unknowns.", r, r, ans, r*r), nil)
	}
}

// nonsym2 draws a 2x2 nonsymmetric integer matrix with a dominant diagonal.
func nonsym2() [][]float64 {
	return [][]float64{{float64(edakit.Int(3, 9)), float64(edakit.Int(-3, 3))}, {float64(edakit.Int(-3, 3)) + 0.5*float64(edakit.Int(1, 3)), float64(edakit.Int(2, 8))}}
}

func gmres(ctx generator.GeneratorContext) generator.Problem {
	a := nonsym2()
	b := []float64{float64(edakit.Int(1, 5)), float64(edakit.Int(-4, 4))}
	ab := matvec(a, b)
	stem := "GMRES solves $Ax = b$ with the nonsymmetric $A = $ " + pmat(a, 1) + " and $b = $ " + pvec(b, 0) + " from $x_0 = 0$. "
	switch kpOf(ctx) {
	case 1:
		alpha := float64(edakit.Int(1, 9)) / 10
		r := []float64{b[0] - alpha*ab[0], b[1] - alpha*ab[1]}
		ans := f(norm2(r), 4)
		q := stem + sp("The first Krylov subspace is span$\\{b\\}$, so $x_1 = \\alpha b$. For the trial value $\\alpha = %s$, what is the residual norm $\\|b - \\alpha Ab\\|_2$? Round to 4 decimal places.", f(alpha, 1))
		return numeric(q, ans, sp("$Ab = (%s, %s)$, so $b - \\alpha Ab = (%s, %s)$ with norm %s.", f(ab[0], 2), f(ab[1], 2), f(r[0], 4), f(r[1], 4), ans), nil)
	case 2:
		alpha := dot(b, ab) / dot(ab, ab)
		res := math.Sqrt(math.Max(0, dot(b, b)-dot(b, ab)*dot(b, ab)/dot(ab, ab)))
		if edakit.Int(0, 1) == 0 {
			ans := f(alpha, 4)
			return numeric(stem+"GMRES picks $x_1 = \\alpha b$ minimizing $\\|b - \\alpha Ab\\|_2$. What is that $\\alpha$? Round to 4 decimal places.", ans,
				sp("Least squares in one variable: $\\alpha = b^TAb/\\|Ab\\|^2 = %s/%s = %s$.", f(dot(b, ab), 3), f(dot(ab, ab), 3), ans), nil)
		}
		ans := f(res, 4)
		return numeric(stem+"After the first GMRES step (the best $x_1 = \\alpha b$), what is the residual norm $\\|b - Ax_1\\|_2$? Round to 4 decimal places.", ans,
			sp("$\\alpha = %s$; the minimal residual is $\\sqrt{\\|b\\|^2 - (b^TAb)^2/\\|Ab\\|^2} = %s$, below $\\|b\\| = %s$.", f(alpha, 4), ans, f(norm2(b), 4)),
			map[string]string{"alpha": f(alpha, 6)})
	default:
		n := 1000 * edakit.Scaled(ctx, 100, 5000)
		m := edakit.Pick([]int{10, 20, 30, 50, 80, 100})
		mb := float64(m+1) * float64(n) * 8 / 1e6
		ans := f(mb, 1)
		q := sp("Restarted GMRES(%d) runs on a circuit matrix with n = %d unknowns and keeps %d + 1 Krylov basis vectors of 8-byte doubles. How much memory do the basis vectors take, in megabytes ($10^6$ bytes)? Round to 1 decimal place.", m, n, m)
		return numeric(q, ans, sp("$(m + 1) \\cdot n \\cdot 8 = %d \\cdot %d \\cdot 8$ bytes $= %s$ MB. Orthogonalization work also grows with the basis, which is why GMRES restarts.", m+1, n, ans), nil)
	}
}

func preconditioning(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		k1 := math.Pow(10, float64(edakit.Int(4, 8)))
		k2 := math.Pow(10, float64(edakit.Int(1, 3)))
		n1 := edakit.Int(20, 300) * 10
		n2 := float64(n1) * math.Sqrt(k2/k1)
		ans := f(n2, 0)
		q := sp("Unpreconditioned CG on a grid matrix with $\\kappa = 10^{%d}$ takes %d iterations. A preconditioner brings the effective condition number to $10^{%d}$. If iterations scale like $\\sqrt\\kappa$, about how many iterations remain? Round to the nearest integer.", int(math.Log10(k1)), n1, int(math.Log10(k2)))
		return numeric(q, ans, sp("$%d \\cdot \\sqrt{10^{%d}/10^{%d}} = %s$.", n1, int(math.Log10(k2)), int(math.Log10(k1)), ans), nil)
	case 2:
		a := math.Pow(10, float64(edakit.Int(1, 4))) * float64(edakit.Int(1, 9))
		d := float64(edakit.Int(1, 9))
		c := -math.Sqrt(a*d) * float64(edakit.Int(1, 9)) / 10
		c = edakit.Round(c, 2)
		rho := math.Abs(c) / math.Sqrt(a*d)
		k := (1 + rho) / (1 - rho)
		ans := f(k, 3)
		m := [][]float64{{a, c}, {c, d}}
		q := "Symmetric Jacobi scaling replaces $A$ by $D^{-1/2}AD^{-1/2}$ with $D = \\mathrm{diag}(A)$. For $A = $ " + pmat(m, 2) + ", what is the 2-norm condition number of the scaled matrix? Round to 3 decimal places."
		return numeric(q, ans, sp("The scaled matrix is $\\begin{pmatrix} 1 & \\rho \\\\ \\rho & 1 \\end{pmatrix}$ with $\\rho = c/\\sqrt{ad} = %s$, eigenvalues $1 \\pm \\rho$, so $\\kappa = (1 + |\\rho|)/(1 - |\\rho|) = %s$.", f(-rho, 4), ans),
			map[string]string{"rho": f(rho, 6)})
	default:
		n1 := edakit.Int(50, 500) * 10
		n2 := edakit.Int(10, 100) * 2
		tmv := float64(edakit.Int(5, 50)) / 10
		tpc := float64(edakit.Int(5, 100)) / 10
		setup := float64(edakit.Int(10, 300))
		t1 := float64(n1) * tmv
		t2 := setup + float64(n2)*(tmv+tpc)
		ans := f(t1/t2, 2)
		q := sp("Without a preconditioner CG needs %d iterations at %s ms per matrix-vector product. An incomplete Cholesky preconditioner costs %s ms to build and %s ms per application, and cuts the count to %d iterations. Ignoring vector operations, what is the overall speedup (old time over new)? Round to 2 decimal places.", n1, f(tmv, 1), f(setup, 0), f(tpc, 1), n2)
		return numeric(q, ans, sp("Old: %s ms. New: $%s + %d(%s + %s) = %s$ ms. Speedup %s.", f(t1, 1), f(setup, 0), n2, f(tmv, 1), f(tpc, 1), f(t2, 1), ans), nil)
	}
}

func ilu(ctx generator.GeneratorContext) generator.Problem {
	a := [][]float64{
		{float64(edakit.Int(4, 9)), -float64(edakit.Int(1, 3)), -float64(edakit.Int(1, 3))},
		{-float64(edakit.Int(1, 3)), float64(edakit.Int(4, 9)), 0},
		{-float64(edakit.Int(1, 3)), 0, float64(edakit.Int(4, 9))},
	}
	l21, l31 := a[1][0]/a[0][0], a[2][0]/a[0][0]
	u22 := a[1][1] - l21*a[0][1]
	fill23 := -l21 * a[0][2]
	fill32 := -l31 * a[0][1]
	stem := "The matrix " + pmat(a, 0) + " has zeros at positions (2,3) and (3,2). "
	switch kpOf(ctx) {
	case 1:
		ans := f(fill23, 4)
		q := stem + "Exact Gaussian elimination of column 1 makes position (2,3) nonzero. What value appears there? Round to 4 decimal places."
		return numeric(q, ans, sp("$a_{23} - l_{21}a_{13} = 0 - (%s/%s)(%s) = %s$. ILU(0) discards this fill because (2,3) is outside the pattern of $A$.", f(a[1][0], 0), f(a[0][0], 0), f(a[0][2], 0), ans), nil)
	case 2:
		uILU := a[2][2] - l31*a[0][2]
		l32 := fill32 / u22
		uExact := a[2][2] - l31*a[0][2] - l32*fill23
		if edakit.Int(0, 1) == 0 {
			ans := f(uILU, 4)
			q := stem + "In ILU(0), every update that would land outside the pattern of $A$ is dropped. What is the last pivot $u_{33}$ of the ILU(0) factorization? Round to 4 decimal places."
			return numeric(q, ans, sp("The fills at (2,3) and (3,2) are dropped, so $l_{32} = 0$ and $u_{33} = a_{33} - l_{31}a_{13} = %s$. The exact LU gives %s.", ans, f(uExact, 4)),
				map[string]string{"exact": f(uExact, 4)})
		}
		ans := f(uExact, 4)
		q := stem + "In the exact LU factorization (no dropping), what is the last pivot $u_{33}$? Round to 4 decimal places."
		return numeric(q, ans, sp("Fill $u_{23} = %s$ and updated $a_{32} = %s$, so $l_{32} = %s$ and $u_{33} = a_{33} - l_{31}a_{13} - l_{32}u_{23} = %s$. ILU(0) would give %s.", f(fill23, 4), f(fill32, 4), f(l32, 4), ans, f(a[2][2]-l31*a[0][2], 4)), nil)
	default:
		n := edakit.Int(5, 7)
		row := make([]float64, n)
		var cand []float64
		for i := range row {
			row[i] = float64(edakit.Int(-90, 90)) / 10
			if row[i] == 0 {
				row[i] = 1.5
			}
		}
		for i := 0; i < edakit.Int(3, 5); i++ {
			cand = append(cand, float64(edakit.Int(-300, 300))/1000)
		}
		tau := edakit.Pick([]float64{0.001, 0.005, 0.01, 0.02})
		nrm := norm2(row)
		kept := 0
		for _, c := range cand {
			if math.Abs(c) >= tau*nrm {
				kept++
			}
		}
		q := sp("ILUT drops a computed fill entry when its magnitude is below $\\tau\\|a_{i*}\\|_2$, using the 2-norm of the original row. Row $i$ of $A$ has entries %s and elimination produces fill candidates %s in it. With $\\tau = %s$, how many fill entries are kept?", floatsStr(row), floatsStr(cand), f(tau, 3))
		return numeric(q, itoa(kept), sp("$\\|a_{i*}\\|_2 = %s$, so the cutoff is %s; %d candidates reach it.", f(nrm, 4), f(tau*nrm, 5), kept), nil)
	}
}

func multigrid(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 7, 63)
		k := edakit.Int(1, n)
		w := edakit.Pick([]float64{2.0 / 3, 0.5, 1})
		s := edakit.Int(1, 3)
		th := float64(k) * math.Pi / float64(n+1)
		g := math.Pow(1-w*(1-math.Cos(th)), float64(s))
		ans := f(g, 4)
		wname := map[float64]string{2.0 / 3: "2/3", 0.5: "1/2", 1: "1"}[w]
		q := sp("Damped Jacobi with weight $\\omega = %s$ on the 1D Laplacian with %d interior points multiplies the Fourier error mode $k$ (frequency $\\theta = k\\pi/(n+1)$) by $1 - \\omega(1 - \\cos\\theta)$ per sweep. By what factor is mode $k = %d$ multiplied after %d sweep(s)? Round to 4 decimal places.", wname, n, k, s)
		return numeric(q, ans, sp("$\\theta = %s$ rad, one sweep gives %s, and %d sweep(s) give %s. High modes shrink fast; low modes (small $\\theta$) barely change.", f(th, 4), f(1-w*(1-math.Cos(th)), 4), s, ans),
			map[string]string{"theta": f(th, 4)})
	case 2:
		L := edakit.Scaled(ctx, 3, 9)
		total := 0
		for l := 1; l <= L; l++ {
			m := (1 << l) - 1
			total += m * m
		}
		fine := (1<<L - 1) * (1<<L - 1)
		q := sp("A power-grid multigrid hierarchy starts from a %d-by-%d grid of unknowns and coarsens by keeping every other point in each direction ($m \\to (m - 1)/2$) down to a single unknown. How many unknowns are there over all levels together?", 1<<L-1, 1<<L-1)
		return numeric(q, itoa(total), sp("Levels have $(2^l - 1)^2$ unknowns for $l = %d, \\dots, 1$. The sum is %d, about 4/3 of the finest level's %d.", L, total, fine), nil)
	default:
		nu := edakit.Int(2, 4)
		rho := float64(edakit.Int(5, 25)) / 100
		d := edakit.Int(4, 10)
		cycles := math.Ceil(float64(d) * math.Ln10 / -math.Log(rho))
		wu := cycles * float64(nu) * 4 / 3
		ans := f(wu, 1)
		q := sp("A 2D V-cycle does %d smoothing sweeps per level in total and converges by a factor %s per cycle, independent of grid size. Count one fine-grid sweep as one work unit; coarser levels add a factor 4/3. How many work units reduce the error by $10^{%d}$? Round to 1 decimal place.", nu, f(rho, 2), d)
		return numeric(q, ans, sp("Cycles $= \\lceil %d\\ln 10/(-\\ln %s) \\rceil = %s$. Work $= %s \\cdot %d \\cdot 4/3 = %s$ work units, the same on any grid size: $O(n)$ total.", d, f(rho, 2), f(cycles, 0), f(cycles, 0), nu, ans), nil)
	}
}
