package em

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func registerSolver(reg *generator.Registry) {
	reg.Register("em.solver.finite_difference", edakit.Func(finiteDifference))
	reg.Register("em.solver.method_of_moments", edakit.Func(methodOfMoments))
	reg.Register("em.solver.capacitance_matrix", edakit.Func(capacitanceMatrix))
	reg.Register("em.solver.fdtd", edakit.Func(fdtd))
	reg.Register("em.solver.method_selection", &edakit.Bank{Items: methodSelectionItems})
}

func finiteDifference(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n, s, e, w := rnd(0, 0.1, 10), rnd(0, 0.1, 10), rnd(0, 0.1, 10), rnd(0, 0.1, 10)
		v := (n + s + e + w) / 4
		return edakit.Numeric(
			sp("A charge-free region is discretized on a square grid with spacing $h$. An interior node's four neighbours currently hold %s, %s, %s and %s V. What value does one Jacobi update of Laplace's equation assign to the node, in volts? Round to 3 decimal places.", num(n, 1), num(s, 1), num(e, 1), num(w, 1)),
			num(v, 3),
			sp("The five-point stencil of $\\nabla^2\\phi = 0$ is $(\\phi_N + \\phi_S + \\phi_E + \\phi_W - 4\\phi_0)/h^2 = 0$, so $\\phi_0$ becomes the neighbour average: %s V.", num(v, 3)),
			facts("v", num(v, 3)))
	case 2:
		vl, vr := edakit.Pick(supplies), 0.0
		x := []float64{rnd(0, 0.1, 10), rnd(0, 0.1, 10), rnd(0, 0.1, 10)}
		gs := edakit.Int(0, 1) == 0
		y := append([]float64(nil), x...)
		if gs {
			y[0] = (vl + y[1]) / 2
			y[1] = (y[0] + y[2]) / 2
			y[2] = (y[1] + vr) / 2
		} else {
			y[0] = (vl + x[1]) / 2
			y[1] = (x[0] + x[2]) / 2
			y[2] = (x[1] + vr) / 2
		}
		k := edakit.Int(1, 3)
		method, how := "Jacobi", "using only the old values"
		if gs {
			method, how = "Gauss-Seidel", "using each new value as soon as it is computed"
		}
		return edakit.Numeric(
			sp("A 1D finite-difference model of $\\phi'' = 0$ has boundary nodes fixed at %s V (left) and 0 V (right) and three interior nodes $x_1, x_2, x_3$ (left to right) with current values %s, %s and %s V. After one %s sweep from left to right, what is $x_%d$, in volts? Round to 4 decimal places.", num(vl, 1), num(x[0], 1), num(x[1], 1), num(x[2], 1), method, k),
			num(y[k-1], 4),
			sp("Each interior node becomes the average of its two neighbours, %s. The sweep gives $x = (%s, %s, %s)$ V; repeated sweeps converge to the straight line between the boundaries.", how, num(y[0], 4), num(y[1], 4), num(y[2], 4)),
			facts("x1", num(y[0], 4), "x2", num(y[1], 4), "x3", num(y[2], 4)))
	default:
		v := edakit.Pick(supplies)
		er := pick(1, 2.7, 3.9, 4.3)
		m := edakit.Scaled(ctx, 6, 12)
		var ds []string
		sum := 0.0
		for i := 0; i < m; i++ {
			d := rnd(0.05, 0.01, 20) * v
			d = edakit.Round(d, 3)
			sum += d
			ds = append(ds, num(d, 3))
		}
		c := 8.854 * er * sum / v
		return edakit.Numeric(
			sp("A 2D finite-difference solution (square cells) has an inner conductor held at %s V inside a grounded outer conductor, with $\\varepsilon_r = %s$. A closed contour around the inner conductor crosses %d cell edges, and the potential drops across those edges (inside minus outside) are %s V. Using Gauss's law on the contour, what is the capacitance per unit length, in pF/m? Round to 2 decimal places.", num(v, 1), num(er, 1), m, strings.Join(ds, ", ")),
			num(c, 2),
			sp("Each crossed edge carries flux $\\varepsilon (\\Delta\\phi/h) \\cdot h = \\varepsilon\\Delta\\phi$ per unit length, so $Q' = \\varepsilon_0\\varepsilon_r \\sum \\Delta\\phi = 8.854 \\cdot %s \\cdot %s$ pC/m and $C' = Q'/V = %s$ pF/m.", num(er, 1), num(sum, 3), num(c, 2)),
			facts("sum_dphi", num(sum, 3)))
	}
}

func methodOfMoments(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		q1, q2 := rnd(-1, 0.1, 30), rnd(-1, 0.1, 30)
		r1, r2 := pick(0.5, 1, 2, 2.5, 4, 5), pick(0.5, 1, 2, 2.5, 4, 5, 10)
		v := 8.988 * (q1/r1 + q2/r2)
		return edakit.Numeric(
			sp("In a free-space method-of-moments model, two panels are approximated as point charges of %s pC and %s pC. A third panel's centre is %s mm from the first and %s mm from the second. What potential do they produce at that centre, in volts? Round to 3 decimal places.", num(q1, 1), num(q2, 1), num(r1, 1), num(r2, 1)),
			num(v, 3),
			sp("The free-space Green's function is $1/(4\\pi\\varepsilon_0 r)$, which is 8.988 V per pC at 1 mm. So $V = 8.988(%s/%s + %s/%s) = %s$ V.", num(q1, 1), num(r1, 1), num(q2, 1), num(r2, 1), num(v, 3)),
			facts("v", num(v, 3)))
	case 2:
		a := float64(edakit.Int(10, 60))
		b := float64(edakit.Int(1, int(a)/2))
		det := a*a - b*b
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(
				sp("A two-conductor method-of-moments model (one panel per conductor) has potential coefficients $P_{11} = P_{22} = %s$ V/pC and $P_{12} = P_{21} = %s$ V/pC. Conductor 1 is held at 1 V and conductor 2 at 0 V. What charge sits on conductor 1, in pC, which equals the Maxwell matrix entry $C_{11}$ in pF? Round to 5 decimal places.", num(a, 0), num(b, 0)),
				num(a/det, 5),
				sp("Solve $Pq = v$ with $v = (1, 0)$: $q = P^{-1}v$, and $P^{-1} = \\frac{1}{%s}\\begin{pmatrix} %s & -%s \\\\ -%s & %s\\end{pmatrix}$. So $q_1 = %s/%s = %s$ pC.", num(det, 0), num(a, 0), num(b, 0), num(b, 0), num(a, 0), num(a, 0), num(det, 0), num(a/det, 5)),
				facts("c11", num(a/det, 5), "c12", num(-b/det, 5)))
		}
		return edakit.Numeric(
			sp("A two-conductor method-of-moments model (one panel per conductor) has potential coefficients $P_{11} = P_{22} = %s$ V/pC and $P_{12} = P_{21} = %s$ V/pC. Conductor 1 is held at 1 V and conductor 2 at 0 V. What charge is induced on conductor 2, in pC, which equals the Maxwell matrix entry $C_{21}$ in pF? Round to 5 decimal places.", num(a, 0), num(b, 0)),
			num(-b/det, 5),
			sp("$q = P^{-1}(1, 0)^T$ with $P^{-1} = \\frac{1}{%s}\\begin{pmatrix} %s & -%s \\\\ -%s & %s\\end{pmatrix}$, so $q_2 = -%s/%s = %s$ pC. Off-diagonal Maxwell entries are negative.", num(det, 0), num(a, 0), num(b, 0), num(b, 0), num(a, 0), num(b, 0), num(det, 0), num(-b/det, 5)),
			facts("c11", num(a/det, 5), "c12", num(-b/det, 5)))
	default:
		n := 1000 * edakit.Scaled(ctx, 1, 50)
		if edakit.Int(0, 1) == 0 {
			mb := 8 * float64(n) * float64(n) / 1e6
			return edakit.Numeric(
				sp("A method-of-moments extraction discretizes the conductor surfaces into %d panels and stores the dense potential-coefficient matrix in double precision (8 bytes per entry). How much memory does the matrix take, in MB ($10^6$ bytes)? Round to 1 decimal place.", n),
				num(mb, 1),
				sp("Every panel interacts with every other, so the matrix has $N^2 = %s$ entries, $8N^2 = %s$ MB. Fast multipole or hierarchical compression avoids storing it.", num(float64(n)*float64(n), 0), num(mb, 1)),
				facts("entries", num(float64(n)*float64(n), 0)))
		}
		t := pick(0.5, 1, 2, 5, 10)
		k := edakit.Int(2, 5)
		return edakit.Numeric(
			sp("Direct LU factorization of a dense %d by %d method-of-moments matrix takes %s s. Refining the mesh multiplies the panel count by %d. Assuming the cubic cost of dense LU, how long does the factorization take now, in seconds? Round to 1 decimal place.", n, n, num(t, 1), k),
			num(t*float64(k*k*k), 1),
			sp("Dense LU costs about $\\tfrac23 N^3$ operations, so time scales by $%d^3 = %d$: %s s.", k, k*k*k, num(t*float64(k*k*k), 1)),
			facts("scale", edakit.IntStr(k*k*k)))
	}
}

func capacitanceMatrix(ctx generator.GeneratorContext) generator.Problem {
	cg := []int{edakit.Int(2, 30), edakit.Int(2, 30), edakit.Int(2, 30)}
	c12, c13, c23 := edakit.Int(1, 20), edakit.Int(0, 10), edakit.Int(1, 20)
	cc := [3][3]int{{0, c12, c13}, {c12, 0, c23}, {c13, c23, 0}}
	var m [3][3]int
	for i := 0; i < 3; i++ {
		m[i][i] = cg[i]
		for j := 0; j < 3; j++ {
			if i != j {
				m[i][i] += cc[i][j]
				m[i][j] = -cc[i][j]
			}
		}
	}
	mat := sp("$\\begin{pmatrix} %d & %d & %d \\\\ %d & %d & %d \\\\ %d & %d & %d \\end{pmatrix}$", m[0][0], m[0][1], m[0][2], m[1][0], m[1][1], m[1][2], m[2][0], m[2][1], m[2][2])
	switch kp(ctx) {
	case 1:
		i, j := edakit.Int(1, 3), edakit.Int(1, 3)
		return edakit.Numeric(
			sp("Three conductors have ground capacitances of %d, %d and %d fF, and coupling capacitances $C_{12} = %d$, $C_{13} = %d$ and $C_{23} = %d$ fF. What is entry $(%d, %d)$ of their Maxwell capacitance matrix, in fF?", cg[0], cg[1], cg[2], c12, c13, c23, i, j),
			edakit.IntStr(m[i-1][j-1]),
			sp("A diagonal entry is the conductor's total capacitance, ground plus all couplings; an off-diagonal entry is minus the coupling. The matrix is %s fF.", mat),
			facts("entry", edakit.IntStr(m[i-1][j-1])))
	case 2:
		i := edakit.Int(1, 3)
		return edakit.Numeric(
			sp("A field solver reports the Maxwell capacitance matrix of three conductors as %s fF. What is the capacitance of conductor %d to ground, in fF?", mat, i),
			edakit.IntStr(cg[i-1]),
			sp("Raising every conductor to 1 V leaves only ground capacitance charged, so the ground capacitance is the row sum: %d fF. The couplings are the negated off-diagonals.", cg[i-1]),
			facts("cg", edakit.IntStr(cg[i-1])))
	default:
		v := []float64{rnd(0, 0.1, 12), rnd(0, 0.1, 12), rnd(0, 0.1, 12)}
		i := edakit.Int(1, 3)
		q := 0.0
		for j := 0; j < 3; j++ {
			q += float64(m[i-1][j]) * v[j]
		}
		return edakit.Numeric(
			sp("Three conductors have the Maxwell capacitance matrix %s fF and sit at %s, %s and %s V. What charge is on conductor %d, in fC? Round to 2 decimal places.", mat, num(v[0], 1), num(v[1], 1), num(v[2], 1), i),
			num(q, 2),
			sp("$Q = CV$, so $Q_%d = \\sum_j C_{%dj}V_j = %s$ fC.", i, i, num(q, 2)),
			facts("q", num(q, 2)))
	}
}

func fdtd(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		f := pick(5, 10, 20, 40, 50, 100)
		er := pick(1, 2.2, 3.5, 4.3, 9.8)
		n := edakit.Pick([]int{10, 15, 20, 30})
		lam := cLight / (f * 1e9 * math.Sqrt(er)) * 1e3
		return edakit.Numeric(
			sp("An FDTD model must resolve frequencies up to %s GHz in a dielectric with $\\varepsilon_r = %s$, using at least %d cells per shortest wavelength. What is the largest allowed cell size, in mm? Round to 3 decimal places.", num(f, 0), num(er, 1), n),
			num(lam/float64(n), 3),
			sp("The shortest wavelength is $c/(f_{max}\\sqrt{\\varepsilon_r}) = %s$ mm, and dividing by %d gives %s mm.", num(lam, 3), n, num(lam/float64(n), 3)),
			facts("lambda_mm", num(lam, 3)))
	case 2:
		dx := pick(0.01, 0.02, 0.05, 0.1, 0.2, 0.5)
		if edakit.Int(0, 1) == 0 {
			dt := dx * 1e-3 / (cLight * math.Sqrt(3)) * 1e15
			return edakit.Numeric(
				sp("A 3D FDTD grid with cubic cells of %s mm contains vacuum regions. What is the largest stable time step, in femtoseconds? Round to 2 decimal places.", num(dx, 2)),
				num(dt, 2),
				sp("The Courant limit for cubic cells in 3D is $\\Delta t \\le \\Delta x/(c\\sqrt3) = %s \\times 10^{-3}/(3.00 \\times 10^8 \\cdot 1.732)$ s, which is %s fs.", num(dx, 2), num(dt, 2)),
				facts("dt_fs", num(dt, 2)))
		}
		dt := dx * 1e-3 / (cLight * math.Sqrt(2)) * 1e15
		return edakit.Numeric(
			sp("A 2D FDTD grid with square cells of %s mm contains vacuum regions. What is the largest stable time step, in femtoseconds? Round to 2 decimal places.", num(dx, 2)),
			num(dt, 2),
			sp("In 2D the Courant limit for square cells is $\\Delta t \\le \\Delta x/(c\\sqrt2)$, which is %s fs.", num(dt, 2)),
			facts("dt_fs", num(dt, 2)))
	default:
		dx := pick(0.01, 0.02, 0.05, 0.1)
		s := pick(0.9, 0.95, 0.99)
		t := pick(0.1, 0.2, 0.5, 1, 2)
		dt := s * dx * 1e-3 / (cLight * math.Sqrt(3))
		steps := int(math.Ceil(t * 1e-9 / dt))
		return edakit.Numeric(
			sp("A 3D FDTD run uses cubic cells of %s mm and a time step of %s times the Courant limit. How many time steps does it take to simulate %s ns?", num(dx, 2), num(s, 2), num(t, 1)),
			edakit.IntStr(steps),
			sp("$\\Delta t = %s \\cdot \\Delta x/(c\\sqrt3) = %s$ fs, so $\\lceil %s\\text{ ns}/\\Delta t \\rceil = %d$ steps. Halving the cell size doubles the step count and multiplies the cell count by 8.", num(s, 2), num(dt*1e15, 3), num(t, 1), steps),
			facts("dt_fs", num(dt*1e15, 3)))
	}
}
