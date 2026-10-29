package numer

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func autodiffForward(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		a := edakit.Int(-4, 5)
		c := edakit.Pick([]int{-3, -2, -1, 1, 2, 3, 4})
		d := edakit.Pick([]int{-3, -2, -1, 1, 2, 3, 4})
		// f(x) = (x^2 + c x)(x - d), dual seed x = a + 1 eps.
		u, du := a*a+c*a, 2*a+c
		w, dw := a-d, 1
		fv, df := u*w, du*w+u*dw
		q := sp("Forward-mode AD evaluates $f(x) = (x^2 + %dx)(x - %d)$ on the dual number $x = %d + 1\\varepsilon$ with $\\varepsilon^2 = 0$. The result is $f + f'\\varepsilon$. What is the $\\varepsilon$ coefficient?", c, d, a)
		q = fixSigns(q)
		return numeric(q, itoa(df), fixSigns(sp("$u = x^2 + %dx = %d + %d\\varepsilon$, $w = x - %d = %d + 1\\varepsilon$, and $(u_0 + u_1\\varepsilon)(w_0 + w_1\\varepsilon) = u_0w_0 + (u_0w_1 + u_1w_0)\\varepsilon = %d + %d\\varepsilon$.", c, u, du, d, w, fv, df)),
			map[string]string{"f": itoa(fv)})
	case 2:
		x := float64(edakit.Int(-20, 20)) / 10
		y := float64(edakit.Int(-20, 20)) / 10
		wrtX := edakit.Int(0, 1) == 0
		var d float64
		var seed, name string
		if wrtX {
			d = 2*x*y + y*math.Cos(x)
			seed, name = "(1, 0)", "x"
		} else {
			d = x*x + math.Sin(x)
			seed, name = "(0, 1)", "y"
		}
		ans := f(d, 4)
		q := sp("Forward-mode AD propagates tangents through $f(x, y) = x^2y + y\\sin x$ at $(x, y) = (%s, %s)$ with seed direction %s. What tangent emerges, i.e. $\\partial f/\\partial %s$? Round to 4 decimal places.", f(x, 1), f(y, 1), seed, name)
		return numeric(q, ans, sp("Each operation carries (value, tangent). The result is $\\partial f/\\partial %s = %s$ (using $\\partial f/\\partial x = 2xy + y\\cos x$, $\\partial f/\\partial y = x^2 + \\sin x$).", name, ans), nil)
	default:
		n := edakit.Int(2, 2000)
		m := edakit.Int(1, 50)
		w := edakit.Pick([]int{1, 4, 8, 16})
		passes := (n + w - 1) / w
		q := sp("A compact model maps %d parameters to %d outputs. Forward-mode AD computes a Jacobian-vector product per tangent, and this implementation carries %d tangent direction(s) per pass. How many forward passes give the full %d-by-%d Jacobian?", n, m, w, m, n)
		return numeric(q, itoa(passes), sp("Forward mode needs one tangent per input: %d tangents at %d per pass is %d passes. Reverse mode would need %d (one per output).", n, w, passes, m), nil)
	}
}

func fixSigns(s string) string {
	s = strings.ReplaceAll(s, "+ -", "- ")
	s = strings.ReplaceAll(s, "- -", "+ ")
	return strings.ReplaceAll(s, " 1x", " x")
}

func adjoint(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		x, y := edakit.Int(-4, 5), edakit.Int(-4, 5)
		forms := []struct {
			expr   string
			dx, dy int
		}{
			{"(x + y) \\cdot xy", (x*y + (x+y)*y), (x*y + (x+y)*x)},
			{"x^2 y + 3y", 2 * x * y, x*x + 3},
			{"(x - 2y)(x + y)", (x + y) + (x - 2*y), -2*(x+y) + (x - 2*y)},
		}
		fm := edakit.Pick(forms)
		wrtX := edakit.Int(0, 1) == 0
		ans, name := fm.dy, "y"
		if wrtX {
			ans, name = fm.dx, "x"
		}
		q := sp("A reverse-mode sweep records $f = %s$ at $(x, y) = (%d, %d)$, seeds $\\bar f = 1$ and propagates adjoints backward through each operation. What adjoint $\\bar %s = \\partial f/\\partial %s$ does it produce?", fm.expr, x, y, name, name)
		return numeric(q, itoa(ans), sp("Each node passes its adjoint times the local partial derivative to its inputs, and contributions to a shared input add: $\\bar %s = %d$.", name, ans), nil)
	case 2:
		a, b := spd2()
		_ = b
		k := edakit.Int(1, 2)
		c := []float64{0, 0}
		c[k-1] = 1
		lam := solve2(a, c)
		j := edakit.Int(1, 2)
		ans := f(lam[j-1], 4)
		q := "A linear circuit solves $Gx = b$ with $G = $ " + pmat(a, 0) + sp(" mS, $x$ in V and $b$ the injected currents in mA. The output is $y = x_%d$. The adjoint method solves $G^T\\lambda = c$ once. What is the sensitivity $\\partial y/\\partial b_%d$, in kΩ (V per mA)? Round to 4 decimal places.", k, j)
		return numeric(q, ans, sp("With $c = e_%d$, $\\lambda = G^{-T}c = (%s, %s)$, and $\\partial y/\\partial b = \\lambda$, so the answer is %s. One transposed solve gives the sensitivity to every injection.", k, f(lam[0], 4), f(lam[1], 4), ans),
			map[string]string{"lambda1": f(lam[0], 6), "lambda2": f(lam[1], 6)})
	default:
		n := edakit.Int(10, 5000)
		m := edakit.Int(1, 20)
		fwd := float64(n)
		rev := 3 * float64(m)
		ans := f(fwd/rev, 2)
		q := sp("A differentiable placer needs the full Jacobian of %d output metric(s) with respect to %d cell coordinates. A forward-mode pass costs about one function evaluation per input direction; a reverse-mode pass costs about 3 evaluations per output. What is the ratio of forward-mode cost to reverse-mode cost? Round to 2 decimal places.", m, n)
		return numeric(q, ans, sp("Forward: %d evaluations. Reverse: 3 · %d = %s. Ratio %s; reverse mode wins when inputs far outnumber outputs, at the price of storing the tape.", n, m, f(rev, 0), ans), nil)
	}
}

func gpu(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := 1000 * edakit.Int(100, 5000)
		d := edakit.Int(3, 12)
		nnz := n * d
		bytes := 12*float64(nnz) + 4*float64(n+1) + 16*float64(n)
		ai := 2 * float64(nnz) / bytes
		ans := f(ai, 4)
		q := sp("A CSR SpMV on n = %d rows with %d nonzeros per row does 2 flops per nonzero and moves $12\\,\\text{nnz} + 4(n + 1) + 16n$ bytes (values, indices, row pointers, $x$ and $y$). What is its arithmetic intensity, in flop/byte? Round to 4 decimal places.", n, d)
		return numeric(q, ans, sp("nnz = %d, flops = %d, bytes = %s, intensity %s flop/byte, far below the balance point of any processor.", nnz, 2*nnz, f(bytes, 0), ans), nil)
	case 2:
		peak := float64(edakit.Pick([]int{10, 20, 40, 60, 80})) // TFLOP/s
		bw := float64(edakit.Pick([]int{1, 2, 3, 5}))*0.8 + 0.2 // TB/s
		ai := edakit.Pick([]float64{0.14, 0.16, 0.25, 0.5, 2, 8, 30, 60, 120})
		att := math.Min(peak*1000, ai*bw*1000)
		ans := f(att, 1)
		q := sp("A GPU has peak %s TFLOP/s and memory bandwidth %s TB/s. By the roofline model, what performance can a kernel with arithmetic intensity %s flop/byte attain, in GFLOP/s? Round to 1 decimal place.", f(peak, 0), f(bw, 1), f(ai, 2))
		bound := "memory-bound"
		if att == peak*1000 {
			bound = "compute-bound"
		}
		return numeric(q, ans, sp("$\\min(\\text{peak}, I \\cdot BW) = \\min(%s, %s \\cdot %s) = %s$ GFLOP/s: %s.", f(peak*1000, 0), f(ai, 2), f(bw*1000, 0), ans, bound), nil)
	default:
		n := 1000 * edakit.Int(100, 5000)
		d := edakit.Int(3, 12)
		nnz := n * d
		bytes := 12*float64(nnz) + 4*float64(n+1) + 16*float64(n)
		bw := float64(edakit.Pick([]int{1, 2, 3, 5}))*0.8 + 0.2
		t := bytes / (bw * 1e12) * 1e6
		ans := f(t, 2)
		q := sp("A CG solver for an IR-drop grid runs on a GPU with %s TB/s of memory bandwidth. Its SpMV has n = %d rows, %d nonzeros per row, and moves $12\\,\\text{nnz} + 4(n + 1) + 16n$ bytes. Treating SpMV as purely bandwidth-bound, how long does one take, in µs? Round to 2 decimal places.", f(bw, 1), n, d)
		return numeric(q, ans, sp("%s bytes at %s TB/s is %s µs. A CPU with a fifth of the bandwidth would take about five times as long, whatever its flop rate.", f(bytes, 0), f(bw, 1), ans), nil)
	}
}

func parallelScaling(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		p := float64(edakit.Int(50, 99)) / 100
		n := edakit.Pick([]int{4, 8, 16, 32, 64, 128, 1024})
		s := 1 / ((1 - p) + p/float64(n))
		ans := f(s, 2)
		q := sp("A transient simulation spends %s percent of its single-thread time in device evaluation and the sparse solve, which parallelize perfectly, and the rest in serial bookkeeping. By Amdahl's law, what speedup do %d threads give? Round to 2 decimal places.", f(100*p, 0), n)
		return numeric(q, ans, sp("$S = 1/((1 - p) + p/N) = 1/(%s + %s/%d) = %s$; the limit for infinite threads is %s.", f(1-p, 2), f(p, 2), n, ans, f(1/(1-p), 2)), nil)
	case 2:
		p := float64(edakit.Int(50, 99)) / 100
		n := edakit.Pick([]int{4, 8, 16, 32, 64, 128, 1024})
		s := (1 - p) + p*float64(n)
		ans := f(s, 2)
		q := sp("A power-grid analysis is weak-scaled: the grid grows with the number of nodes so each of %d processors keeps the same share. In the parallel run, %s percent of the time is parallel work. By Gustafson's law, what is the scaled speedup? Round to 2 decimal places.", n, f(100*p, 0))
		return numeric(q, ans, sp("$S = (1 - p) + pN = %s + %s \\cdot %d = %s$.", f(1-p, 2), f(p, 2), n, ans), nil)
	default:
		n := edakit.Scaled(ctx, 7, 14)
		parent := make([]int, n+1)
		for j := 1; j < n; j++ {
			parent[j] = edakit.Int(j+1, minInt(n, j+edakit.Int(1, 4)))
		}
		depth := make([]int, n+1)
		height := 0
		for j := n; j >= 1; j-- {
			if j == n {
				depth[j] = 1
			} else {
				depth[j] = depth[parent[j]] + 1
			}
			if depth[j] > height {
				height = depth[j]
			}
		}
		var ps []int
		for j := 1; j < n; j++ {
			ps = append(ps, parent[j])
		}
		q := sp("The elimination tree of a sparse Cholesky factor on columns 1 to %d has parent(j) = %s for j = 1, ..., %d (column %d is the root). Each column takes one unit of time and a column can be factored once all its children are done. What is the shortest possible factorization time with unlimited processors (the tree height, counting columns)?", n, ints(ps), n-1, n)
		return numeric(q, itoa(height), sp("The critical path is the longest leaf-to-root chain: %d columns, so the speedup is at most %d/%d.", height, n, height), nil)
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func reductionOrder(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		k := edakit.Int(24, 26)
		big := float32(math.Pow(2, float64(k)))
		s1, s2 := float32(edakit.Int(1, 3)), float32(edakit.Int(1, 3))
		var r float32
		left := edakit.Int(0, 1) == 0
		if left {
			r = big + s1
			r = r + s2
		} else {
			t := s1 + s2
			r = big + t
		}
		ans := f(float64(r), 0)
		order := sp("$(2^{%d} + %s) + %s$", k, f(float64(s1), 0), f(float64(s2), 0))
		if !left {
			order = sp("$2^{%d} + (%s + %s)$", k, f(float64(s1), 0), f(float64(s2), 0))
		}
		q := sp("In IEEE single precision (24-bit significand, round to nearest even), compute %s exactly as the hardware would, rounding after each addition. What is the result? ($2^{%d} = %s$.)", order, k, f(math.Pow(2, float64(k)), 0))
		return numeric(q, ans, sp("Near $2^{%d}$ the spacing of float32 numbers is %s, so small additions round to the nearest representable value at each step; the result is %s.", k, f(math.Pow(2, float64(k-23)), 0), ans), nil)
	case 2:
		big := float32(math.Pow(2, float64(edakit.Int(24, 25))))
		n := edakit.Int(3, 5)
		xs := []float32{big}
		for i := 0; i < n; i++ {
			xs = append(xs, float32(edakit.Int(1, 3)))
		}
		xs = edakit.Shuffle(xs)
		var s float32
		for _, x := range xs {
			s = s + x
		}
		var shown []int
		for _, x := range xs {
			shown = append(shown, int(x))
		}
		q := sp("Partial sums from %d threads arrive in the order %s and are added left to right in float32 (24-bit significand, round to nearest even). What is the final sum?", len(xs), ints(shown))
		return numeric(q, f(float64(s), 0), sp("Once the running sum reaches %s, each small term is at most one ulp and may be rounded away; the left-to-right float32 sum is %s, while the exact sum is %s.", f(float64(big), 0), f(float64(s), 0), f(exactSum(xs), 0)), nil)
	default:
		n := edakit.Pick([]int{1000, 4096, 10000, 65536, 100000, 1000000})
		lg := math.Ceil(math.Log2(float64(n)))
		ratio := float64(n-1) / lg
		ans := f(ratio, 1)
		q := sp("Summing n = %d terms, the first-order error bound is $(n - 1)u\\sum|x_i|$ for left-to-right summation and $\\lceil\\log_2 n\\rceil u\\sum|x_i|$ for pairwise (tree) summation. By what factor is the pairwise bound smaller? Round to 1 decimal place.", n)
		return numeric(q, ans, sp("$(n - 1)/\\lceil\\log_2 n\\rceil = %d/%s = %s$. A fixed reduction tree is both more accurate and reproducible.", n-1, f(lg, 0), ans), nil)
	}
}

func exactSum(xs []float32) float64 {
	s := 0.0
	for _, x := range xs {
		s += float64(x)
	}
	return s
}
