package place

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// mincut counts cut nets, bisection levels and terminal-propagated cuts.
func mincut(ctx generator.GeneratorContext) generator.Problem {
	switch ctx.KPOrdinal {
	case 2:
		n := edakit.Pick([]int{1000, 5000, 20000, 100000, 250000, 1000000, 3000000})
		k := edakit.Pick([]int{1, 4, 8, 10, 16, 32})
		lv, regions := 0, 1
		for float64(n)/float64(regions) > float64(k) {
			lv++
			regions *= 2
		}
		return edakit.Numeric(edakit.Sprintf("A min-cut placer recursively bisects a netlist of %d cells, splitting every region into two halves of equal cell count at each level. It stops when every region holds at most %d cells. How many levels of bisection does it perform?", n, k), edakit.IntStr(lv),
			edakit.Sprintf("After L levels there are $2^L$ regions of %d/$2^L$ cells. The smallest L with %d / $2^L \\le$ %d is %d.", n, n, k, lv), map[string]string{"levels": edakit.IntStr(lv)})
	case 3:
		m := edakit.Int(4, 6)
		var desc []string
		cutL, cutR := 0, 0
		for i := 0; i < m; i++ {
			var mem []string
			hasL, hasR := false, false
			for j := 0; j < edakit.Int(1, 3); j++ {
				side := edakit.Pick([]string{"L", "R"})
				if edakit.Int(0, 3) == 0 {
					mem = append(mem, "a terminal propagated to the "+map[string]string{"L": "left", "R": "right"}[side])
				} else {
					mem = append(mem, edakit.Sprintf("a cell in %s", side))
				}
				if side == "L" {
					hasL = true
				} else {
					hasR = true
				}
			}
			if hasR {
				cutL++
			}
			if hasL {
				cutR++
			}
			list := strings.Join(mem, " and ")
			if len(mem) > 2 {
				list = strings.Join(mem[:len(mem)-1], ", ") + " and " + mem[len(mem)-1]
			}
			desc = append(desc, edakit.Sprintf("net n%d joins c with %s", i+1, list))
		}
		side := edakit.Pick([]string{"left", "right"})
		ans := cutL
		if side == "right" {
			ans = cutR
		}
		return edakit.Numeric(edakit.Sprintf("A min-cut placer is splitting a region into a left half L and a right half R. Every cell except c is already assigned, and connections leaving the region have been propagated to the nearer side as fixed terminals. Cell c has these nets: %s. If c goes to the %s half, how many of its nets are cut?", strings.Join(desc, "; "), side), edakit.IntStr(ans),
			edakit.Sprintf("A net is cut when it has members on both sides. Placing c left cuts %d nets and placing it right cuts %d, so the %s side gives %d. Without the propagated terminals the placer would ignore pulls from outside the region.", cutL, cutR, side, ans),
			map[string]string{"cut_left": edakit.IntStr(cutL), "cut_right": edakit.IntStr(cutR)})
	}
	n := edakit.Scaled(ctx, 6, 10)
	side := make([]string, n)
	var ls, rs []string
	for i := range side {
		side[i] = edakit.Pick([]string{"L", "R"})
		if side[i] == "L" {
			ls = append(ls, edakit.Sprintf("c%d", i+1))
		} else {
			rs = append(rs, edakit.Sprintf("c%d", i+1))
		}
	}
	if len(ls) == 0 || len(rs) == 0 {
		side[0], side[n-1] = "L", "R"
		ls, rs = nil, nil
		for i := range side {
			if side[i] == "L" {
				ls = append(ls, edakit.Sprintf("c%d", i+1))
			} else {
				rs = append(rs, edakit.Sprintf("c%d", i+1))
			}
		}
	}
	m := edakit.Scaled(ctx, 4, 8)
	cut := 0
	var nd []string
	for i := 0; i < m; i++ {
		mem := distinctInts(edakit.Int(2, 4), 0, n-1)
		names := make([]string, len(mem))
		l, r := false, false
		for j, c := range mem {
			names[j] = edakit.Sprintf("c%d", c+1)
			if side[c] == "L" {
				l = true
			} else {
				r = true
			}
		}
		if l && r {
			cut++
		}
		nd = append(nd, "{"+strings.Join(names, ", ")+"}")
	}
	return edakit.Numeric(edakit.Sprintf("A min-cut placer has bisected a region: cells %s are on the left and %s on the right. The nets are %s. How many nets cross the cut line?", strings.Join(ls, ", "), strings.Join(rs, ", "), strings.Join(nd, ", ")), edakit.IntStr(cut),
		edakit.Sprintf("Count the nets with at least one cell on each side: %d.", cut), map[string]string{"cut": edakit.IntStr(cut)})
}

type fixedPin struct {
	w float64
	x int
}

func fmtFixed(fs []fixedPin) string {
	parts := make([]string, len(fs))
	for i, f := range fs {
		parts[i] = edakit.Sprintf("x = %d with weight %s", f.x, edakit.Num(f.w, 0))
	}
	return strings.Join(parts, ", ")
}

func randFixed(n int) []fixedPin {
	out := make([]fixedPin, n)
	for i := range out {
		out[i] = fixedPin{float64(edakit.Int(1, 4)), edakit.Int(0, 100)}
	}
	return out
}

// quadratic solves 1-, 2- and 3-cell quadratic placements in one dimension.
func quadratic(ctx generator.GeneratorContext) generator.Problem {
	head := "A quadratic placer minimises the sum over two-pin connections of $w (x_i - x_j)^2$. Coordinates are in um. "
	switch ctx.KPOrdinal {
	case 2:
		fa, fb := randFixed(edakit.Int(1, 2)), randFixed(edakit.Int(1, 2))
		wab := float64(edakit.Int(1, 4))
		A := [][]float64{{wab, -wab}, {-wab, wab}}
		b := []float64{0, 0}
		for _, f := range fa {
			A[0][0] += f.w
			b[0] += f.w * float64(f.x)
		}
		for _, f := range fb {
			A[1][1] += f.w
			b[1] += f.w * float64(f.x)
		}
		x := solve(A, b)
		who := edakit.Int(0, 1)
		ans := edakit.Num(x[who], 2)
		return edakit.Numeric(head+edakit.Sprintf("Movable cell a connects to fixed pins at %s. Movable cell b connects to fixed pins at %s. Cells a and b are connected with weight %s. What is the optimal x of cell %s? Round to 2 decimal places.", fmtFixed(fa), fmtFixed(fb), edakit.Num(wab, 0), []string{"a", "b"}[who]), ans,
			edakit.Sprintf("Setting both partial derivatives to zero gives the system [[%s, %s], [%s, %s]] x = [%s, %s], whose solution is $x_a$ = %s, $x_b$ = %s.", edakit.Num(A[0][0], 0), edakit.Num(A[0][1], 0), edakit.Num(A[1][0], 0), edakit.Num(A[1][1], 0), edakit.Num(b[0], 0), edakit.Num(b[1], 0), edakit.Num(x[0], 2), edakit.Num(x[1], 2)),
			map[string]string{"xa": edakit.Num(x[0], 4), "xb": edakit.Num(x[1], 4)})
	case 3:
		l := fixedPin{float64(edakit.Int(1, 4)), edakit.Int(0, 30)}
		r := fixedPin{float64(edakit.Int(1, 4)), edakit.Int(70, 100)}
		w1, w2 := float64(edakit.Int(1, 4)), float64(edakit.Int(1, 4))
		extra := fixedPin{float64(edakit.Int(1, 3)), edakit.Int(0, 100)}
		A := [][]float64{{l.w + w1, -w1, 0}, {-w1, w1 + w2 + extra.w, -w2}, {0, -w2, w2 + r.w}}
		b := []float64{l.w * float64(l.x), extra.w * float64(extra.x), r.w * float64(r.x)}
		x := solve(A, b)
		ans := edakit.Num(x[1], 2)
		return edakit.Numeric(head+edakit.Sprintf("Three movable cells form a chain: fixed pin at x = %d -(weight %s)- a -(weight %s)- b -(weight %s)- c -(weight %s)- fixed pin at x = %d. Cell b also connects to a fixed pin at x = %d with weight %s. What is the optimal x of cell b? Round to 2 decimal places.", l.x, edakit.Num(l.w, 0), edakit.Num(w1, 0), edakit.Num(w2, 0), edakit.Num(r.w, 0), r.x, extra.x, edakit.Num(extra.w, 0)), ans,
			edakit.Sprintf("The optimality conditions form a tridiagonal 3 x 3 system. Solving it gives $x_a$ = %s, $x_b$ = %s, $x_c$ = %s.", edakit.Num(x[0], 2), edakit.Num(x[1], 2), edakit.Num(x[2], 2)),
			map[string]string{"xa": edakit.Num(x[0], 4), "xb": edakit.Num(x[1], 4), "xc": edakit.Num(x[2], 4)})
	}
	fs := randFixed(edakit.Scaled(ctx, 2, 4))
	num, den := 0.0, 0.0
	for _, f := range fs {
		num += f.w * float64(f.x)
		den += f.w
	}
	ans := edakit.Num(num/den, 2)
	return edakit.Numeric(head+edakit.Sprintf("One movable cell connects to fixed pins at %s. What is its optimal x? Round to 2 decimal places.", fmtFixed(fs)), ans,
		edakit.Sprintf("Setting the derivative to zero gives the weighted mean: %s / %s = %s.", edakit.Num(num, 0), edakit.Num(den, 0), ans), map[string]string{"x": ans})
}

// laplacianSystem asks for an entry of the quadratic placement system Ax = b.
func laplacianSystem(ctx generator.GeneratorContext) generator.Problem {
	n := 3
	W := make([][]int, n)
	for i := range W {
		W[i] = make([]int, n)
	}
	var desc []string
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if edakit.Int(0, 3) > 0 {
				W[i][j] = edakit.Int(1, 4)
				W[j][i] = W[i][j]
				desc = append(desc, edakit.Sprintf("m%d-m%d weight %d", i+1, j+1, W[i][j]))
			}
		}
	}
	diag := make([]int, n)
	bb := make([]int, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			diag[i] += W[i][j]
		}
		for k := 0; k < edakit.Int(1, 2); k++ {
			w, x := edakit.Int(1, 4), edakit.Int(0, 50)
			diag[i] += w
			bb[i] += w * x
			desc = append(desc, edakit.Sprintf("m%d-fixed pin at x = %d weight %d", i+1, x, w))
		}
	}
	desc = edakit.Shuffle(desc)
	head := edakit.Sprintf("Quadratic placement in x minimises $\\frac{1}{2}x^T A x - b^T x$ plus a constant for movable cells m1, m2, m3. The two-pin connections are: %s. ", strings.Join(desc, "; "))
	i := edakit.Int(0, n-1)
	switch ctx.KPOrdinal {
	case 2:
		j := (i + edakit.Int(1, 2)) % n
		v := -W[i][j]
		return edakit.Numeric(head+edakit.Sprintf("What is the entry $A_{%d%d}$?", i+1, j+1), edakit.IntStr(v),
			edakit.Sprintf("Off-diagonal entries are minus the connection weight between the two movable cells: %d.", v), map[string]string{"a_ij": edakit.IntStr(v)})
	case 3:
		return edakit.Numeric(head+edakit.Sprintf("What is the entry $b_{%d}$ of the right-hand side?", i+1), edakit.IntStr(bb[i]),
			edakit.Sprintf("The right-hand side collects weight times fixed-pin coordinate over m%d's fixed connections: %d.", i+1, bb[i]), map[string]string{"b_i": edakit.IntStr(bb[i])})
	}
	return edakit.Numeric(head+edakit.Sprintf("What is the diagonal entry $A_{%d%d}$?", i+1, i+1), edakit.IntStr(diag[i]),
		edakit.Sprintf("A diagonal entry is the total weight of every connection at m%d, to movable cells and to fixed pins alike: %d. The fixed connections make A strictly diagonally dominant in that row, which is what makes the system positive definite.", i+1, diag[i]),
		map[string]string{"a_ii": edakit.IntStr(diag[i])})
}

// cgSolve takes the first conjugate-gradient step on a 2 x 2 placement system.
func cgSolve(ctx generator.GeneratorContext) generator.Problem {
	c := float64(edakit.Int(1, 4))
	a := c + float64(edakit.Int(1, 6))
	d := c + float64(edakit.Int(1, 6))
	b := []float64{float64(edakit.Int(5, 60)), float64(edakit.Int(5, 60))}
	Ab := []float64{a*b[0] - c*b[1], -c*b[0] + d*b[1]}
	rr := b[0]*b[0] + b[1]*b[1]
	pAp := b[0]*Ab[0] + b[1]*Ab[1]
	alpha := rr / pAp
	head := edakit.Sprintf("A two-cell placement system is $A x = b$ with $A = \\begin{pmatrix} %s & %s \\\\ %s & %s \\end{pmatrix}$ and $b = (%s, %s)$. ", edakit.Num(a, 0), edakit.Num(-c, 0), edakit.Num(-c, 0), edakit.Num(d, 0), edakit.Num(b[0], 0), edakit.Num(b[1], 0))
	switch ctx.KPOrdinal {
	case 1:
		ans := edakit.Num(alpha, 4)
		return edakit.Numeric(head+"Conjugate gradient starts from $x_0 = 0$, so $r_0 = p_0 = b$. What is the first step length $\\alpha_0 = r_0^T r_0 / p_0^T A p_0$? Round to 4 decimal places.", ans,
			edakit.Sprintf("$r_0^T r_0$ = %s and $A p_0$ = (%s, %s), so $p_0^T A p_0$ = %s and $\\alpha_0$ = %s.", edakit.Num(rr, 0), edakit.Num(Ab[0], 0), edakit.Num(Ab[1], 0), edakit.Num(pAp, 0), ans), map[string]string{"alpha": ans})
	case 3:
		mid := (a + d) / 2
		rad := math.Sqrt((a-d)*(a-d)/4 + c*c)
		k := (mid + rad) / (mid - rad)
		f := (math.Sqrt(k) - 1) / (math.Sqrt(k) + 1)
		ans := edakit.Num(f, 4)
		return edakit.Numeric(head+"The standard convergence bound for conjugate gradient is $\\|e_k\\|_A \\le 2\\rho^k \\|e_0\\|_A$ with $\\rho = (\\sqrt{\\kappa}-1)/(\\sqrt{\\kappa}+1)$, where $\\kappa$ is the condition number of A. What is $\\rho$ for this system? Round to 4 decimal places.", ans,
			edakit.Sprintf("The eigenvalues of the symmetric A are %s and %s, so $\\kappa$ = %s and $\\rho$ = %s. Weak connections to fixed pins make $\\lambda_{min}$ small and $\\kappa$ large, which is why placers precondition.", edakit.Num(mid+rad, 4), edakit.Num(mid-rad, 4), edakit.Num(k, 4), ans),
			map[string]string{"kappa": edakit.Num(k, 4)})
	}
	who := edakit.Int(0, 1)
	x1 := alpha * b[who]
	ans := edakit.Num(x1, 4)
	return edakit.Numeric(head+edakit.Sprintf("Starting from $x_0 = 0$ (so $r_0 = p_0 = b$), conjugate gradient takes the step $x_1 = x_0 + \\alpha_0 p_0$ with $\\alpha_0 = r_0^T r_0 / p_0^T A p_0$. What is component %d of $x_1$? Round to 4 decimal places.", who+1), ans,
		edakit.Sprintf("$\\alpha_0$ = %s / %s = %s, so $x_1$ = %s b = (%s, %s).", edakit.Num(rr, 0), edakit.Num(pAp, 0), edakit.Num(alpha, 6), edakit.Num(alpha, 6), edakit.Num(alpha*b[0], 4), edakit.Num(alpha*b[1], 4)),
		map[string]string{"alpha": edakit.Num(alpha, 6)})
}

// spreading adds an anchor pseudo-net toward a spread target.
func spreading(ctx generator.GeneratorContext) generator.Problem {
	fs := randFixed(edakit.Int(2, 3))
	S, Wt := 0.0, 0.0
	for _, f := range fs {
		S += f.w * float64(f.x)
		Wt += f.w
	}
	xs := S / Wt
	t := float64(edakit.Int(0, 100))
	for try := 0; math.Abs(t-xs) < 10 && try < generator.MaxDrawAttempts; try++ {
		t = float64(edakit.Int(0, 100))
	}
	if math.Abs(t-xs) < 10 {
		t = 0
		if xs < 50 {
			t = 100
		}
	}
	head := edakit.Sprintf("A movable cell connects to fixed pins at %s (um), so the plain quadratic solution puts it at x = %s. A spreading step computed a target position x = %s for it in a less crowded area. ", fmtFixed(fs), edakit.Num(xs, 2), edakit.Num(t, 0))
	switch ctx.KPOrdinal {
	case 2:
		a0 := edakit.Pick([]float64{0.5, 1, 2})
		k := edakit.Int(2, 6)
		prev := xs
		alpha := a0 * float64(k) / math.Abs(prev-t)
		x := (S + alpha*t) / (Wt + alpha)
		ans := edakit.Num(x, 2)
		return edakit.Numeric(head+edakit.Sprintf("The placer adds an anchor: a pseudo-net from the cell to its target with weight $\\alpha = %s k / |x_{prev} - x_{target}|$ at iteration k, where $x_{prev}$ is the cell's previous position, here the plain quadratic solution. At iteration k = %d, where does the cell land? Round to 2 decimal places.", edakit.Num(a0, 1), k), ans,
			edakit.Sprintf("The weight is %s x %d / %s = %s. The weighted mean of the fixed pins and the target is (%s + %s x %s) / (%s + %s) = %s.", edakit.Num(a0, 1), k, edakit.Num(math.Abs(prev-t), 2), edakit.Num(alpha, 4), edakit.Num(S, 0), edakit.Num(alpha, 4), edakit.Num(t, 0), edakit.Num(Wt, 0), edakit.Num(alpha, 4), ans),
			map[string]string{"alpha": edakit.Num(alpha, 4)})
	case 3:
		f := edakit.Pick([]float64{0.25, 0.5, 0.75})
		p := xs + f*(t-xs)
		alpha := (S - p*Wt) / (p - t)
		ans := edakit.Num(alpha, 2)
		return edakit.Numeric(head+edakit.Sprintf("The placer adds an anchor pseudo-net from the cell to its target with weight $\\alpha$. What value of $\\alpha$ makes the cell land at x = %s? Round to 2 decimal places.", edakit.Num(p, 2)), ans,
			edakit.Sprintf("The cell lands at (%s + %s t)/(%s + %s). Solving for a landing point %s gives $\\alpha$ = (%s - %s x %s)/(%s - %s) = %s.", edakit.Num(S, 0), "$\\alpha$", edakit.Num(Wt, 0), "$\\alpha$", edakit.Num(p, 2), edakit.Num(S, 0), edakit.Num(p, 2), edakit.Num(Wt, 0), edakit.Num(p, 2), edakit.Num(t, 0), ans),
			map[string]string{"alpha": ans})
	}
	alpha := float64(edakit.Int(1, 6))
	x := (S + alpha*t) / (Wt + alpha)
	ans := edakit.Num(x, 2)
	return edakit.Numeric(head+edakit.Sprintf("The placer adds an anchor: a pseudo-net of weight %s from the cell to a fixed point at the target. Where does the cell land now? Round to 2 decimal places.", edakit.Num(alpha, 0)), ans,
		edakit.Sprintf("The anchor is one more fixed connection: (%s + %s x %s) / (%s + %s) = %s.", edakit.Num(S, 0), edakit.Num(alpha, 0), edakit.Num(t, 0), edakit.Num(Wt, 0), edakit.Num(alpha, 0), ans),
		map[string]string{"x": ans})
}

// forceDirected computes zero-force targets and spring forces.
func forceDirected(ctx generator.GeneratorContext) generator.Problem {
	for attempt := 0; ; attempt++ {
		n := edakit.Scaled(ctx, 2, 4)
		nb := randPts(n, 0, 40)
		cs := make([]int, n)
		last := attempt >= generator.MaxDrawAttempts-1
		if last {
			n, nb, cs = 2, []pt{{13, 34}, {18, 40}}, make([]int, 2)
		}
		var desc []string
		sx, sy, sc := 0, 0, 0
		for i := range nb {
			cs[i] = edakit.Int(1, 4)
			if last {
				cs[i] = 3 + i
			}
			sx += cs[i] * nb[i].x
			sy += cs[i] * nb[i].y
			sc += cs[i]
			desc = append(desc, edakit.Sprintf("%s with connection weight %d", fmtPt(nb[i]), cs[i]))
		}
		cur := pt{edakit.Int(0, 40), edakit.Int(0, 40)}
		head := edakit.Sprintf("In force-directed placement each connection acts as a spring pulling with force equal to its weight times the distance. Cell c sits at %s and connects to cells at %s. ", fmtPt(cur), strings.Join(desc, ", "))
		tx, ty := float64(sx)/float64(sc), float64(sy)/float64(sc)
		switch ctx.KPOrdinal {
		case 2:
			f := sx - sc*cur.x
			return edakit.Numeric(head+"What is the x component of the net spring force on c (positive means toward larger x)?", edakit.IntStr(f),
				edakit.Sprintf("$F_x = \\sum_j c_j (x_j - x_c)$ = %d - %d x %d = %d.", sx, sc, cur.x, f), map[string]string{"fx": edakit.IntStr(f)})
		case 3:
			if !last && (math.Abs(tx-math.Floor(tx)-0.5) < 1e-9 || math.Abs(ty-math.Floor(ty)-0.5) < 1e-9) {
				continue
			}
			rx, ry := int(math.Round(tx)), int(math.Round(ty))
			d := absI(rx-cur.x) + absI(ry-cur.y)
			return edakit.Numeric(head+"The placer moves c to the integer grid point nearest its zero-force target (rounding each coordinate). How far does c move, in Manhattan distance?", edakit.IntStr(d),
				edakit.Sprintf("The zero-force target is (%s, %s), which rounds to (%d, %d). From %s that is %d.", edakit.Num(tx, 3), edakit.Num(ty, 3), rx, ry, fmtPt(cur), d),
				map[string]string{"tx": edakit.Num(tx, 3), "ty": edakit.Num(ty, 3)})
		}
		ans := edakit.Num(tx, 2)
		return edakit.Numeric(head+"What is the x coordinate of c's zero-force target location? Round to 2 decimal places.", ans,
			edakit.Sprintf("The forces balance at the weighted centroid: x = %d / %d = %s.", sx, sc, ans), map[string]string{"tx": ans})
	}
}
