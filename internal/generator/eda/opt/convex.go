package opt

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func quadStem(a, b, c, d int) string {
	return "$f(x, y) = " + poly([]int{a, b, c, d}, []string{"x^2", "xy", "y^2", "x"}) + "$"
}

func sgn(x int) string {
	if x < 0 {
		return "-"
	}
	return "+"
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func convexityClass(ctx generator.GeneratorContext) generator.Problem {
	a := edakit.Int(-3, 4)
	if a == 0 {
		a = 1
	}
	c := edakit.Int(-3, 4)
	if c == 0 {
		c = 2
	}
	b := edakit.Int(-6, 6)
	// sometimes force a singular PSD Hessian: 4ac = b^2 with a, c > 0
	if edakit.Int(0, 3) == 0 {
		a, c = edakit.Pick([]int{1, 4, 9}), 1
		b = int(math.Round(2 * math.Sqrt(float64(a*c))))
		if edakit.Int(0, 1) == 0 {
			b = -b
		}
	}
	d := edakit.Int(-5, 5)
	det := 4*a*c - b*b
	var ans string
	switch {
	case a > 0 && det > 0:
		ans = "Strictly convex"
	case a < 0 && det > 0:
		ans = "Strictly concave"
	case det == 0 && a > 0:
		ans = "Convex but not strictly convex"
	case det == 0 && a < 0:
		ans = "Concave but not strictly concave"
	default:
		ans = "Neither convex nor concave"
	}
	opts := []string{"Strictly convex", "Strictly concave", "Convex but not strictly convex", "Neither convex nor concave"}
	if ans == "Concave but not strictly concave" {
		opts[2] = ans
	}
	var dis []string
	for _, o := range opts {
		if o != ans {
			dis = append(dis, o)
		}
	}
	q := "A placement cost term is " + quadStem(a, b, c, d) + ". Its Hessian is constant. How is f classified on all of the plane?"
	return edakit.Choice(q, ans,
		edakit.Sprintf("The Hessian is $\\begin{pmatrix} %d & %d \\\\ %d & %d \\end{pmatrix}$ with determinant $4ac - b^2 = %d$ and leading entry %d. Positive determinant with a positive leading entry is positive definite; zero determinant is semidefinite; a negative determinant means one positive and one negative eigenvalue, a saddle.", 2*a, b, b, 2*c, det, 2*a),
		dis...)
}

func hessianDet(ctx generator.GeneratorContext) generator.Problem {
	a, c := edakit.Int(1, 6), edakit.Int(1, 6)
	b := edakit.Int(-8, 8)
	d := edakit.Int(-5, 5)
	det := 4*a*c - b*b
	q := "For " + quadStem(a, b, c, d) + ", what is the determinant of the Hessian matrix?"
	return edakit.Choice(q, edakit.IntStr(det),
		edakit.Sprintf("The Hessian is $\\begin{pmatrix} 2a & b \\\\ b & 2c \\end{pmatrix} = \\begin{pmatrix} %d & %d \\\\ %d & %d \\end{pmatrix}$, so its determinant is $4ac - b^2 = %d - %d = %d$. With $2a > 0$, f is convex exactly when this is nonnegative.", 2*a, b, b, 2*c, 4*a*c, b*b, det),
		edakit.Distinct(det, 3, a*c-b*b, 4*a*c+b*b, 2*a*c-b*b, 4*a*c-2*b*b)...)
}

func kkt(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	p, q := edakit.Pick([]int{1, 2, 4}), edakit.Pick([]int{1, 2, 4})
	a, b := edakit.Int(3, 12), edakit.Int(3, 12)
	var c int
	var s, lam, x, y float64
	for {
		if k == 3 && edakit.Int(0, 2) == 0 {
			c = a + b + edakit.Int(0, 4) // inactive or just touching
		} else {
			c = edakit.Int(2, a+b-1)
		}
		s = 1/(2*float64(p)) + 1/(2*float64(q))
		lam = math.Max(0, float64(a+b-c)/s)
		x = float64(a) - lam/(2*float64(p))
		y = float64(b) - lam/(2*float64(q))
		if x > 0 && y > 0 {
			break
		}
	}
	stem := edakit.Sprintf("A sizing step wants gate sizes near $a = %d$ and $b = %d$ but shares an area budget. It solves $\\min\\ %s(x - %d)^2 + %s(y - %d)^2$ subject to $x + y \\le %d$.", a, b, coefStr(p), a, coefStr(q), b, c)
	expl := edakit.Sprintf("Stationarity gives $2 \\cdot %d(x - %d) + \\lambda = 0$ and $2 \\cdot %d(y - %d) + \\lambda = 0$, so $x = %d - \\lambda/%d$ and $y = %d - \\lambda/%d$. ", p, a, q, b, a, 2*p, b, 2*q)
	if a+b > c {
		expl += edakit.Sprintf("The unconstrained point violates the budget, so the constraint is active: $x + y = %d$ gives $\\lambda = (%d + %d - %d)/%s = %s$, $x = %s$, $y = %s$.", c, a, b, c, f(s, 3), f(lam, 3), f(x, 3), f(y, 3))
	} else {
		expl += edakit.Sprintf("The unconstrained point (%d, %d) already satisfies $x + y \\le %d$, so complementary slackness gives $\\lambda = 0$.", a, b, c)
	}
	switch k {
	case 1:
		return num(stem+" What is the optimal x?"+roundNote(2), f(x, 2), expl, "lambda", f(lam, 3), "x", f(x, 3), "y", f(y, 3))
	default:
		return num(stem+" What is the KKT multiplier $\\lambda$ of the budget constraint at the optimum?"+roundNote(2), f(lam, 2), expl, "lambda", f(lam, 3), "x", f(x, 3), "y", f(y, 3))
	}
}

func quadraticPlacement(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 2, 5)
		sw, swp := 0, 0
		var parts []string
		for i := 0; i < n; i++ {
			w := edakit.Int(1, 4)
			p := edakit.Int(0, 100)
			sw += w
			swp += w * p
			parts = append(parts, edakit.Sprintf("a pin at x = %d with weight %d", p, w))
		}
		x := float64(swp) / float64(sw)
		return num("One movable cell connects by two-pin nets to fixed pins: "+strings.Join(parts, ", ")+". Quadratic placement minimizes $\\sum w_k (x - p_k)^2$. Where does the cell land, in um?"+roundNote(2),
			f(x, 2),
			edakit.Sprintf("Setting the derivative $2\\sum w_k (x - p_k)$ to zero gives the weighted mean $x = %d/%d = %s$.", swp, sw, f(x, 2)),
			"x", f(x, 2))
	case 2:
		w1, w2, w12 := edakit.Int(1, 4), edakit.Int(1, 4), edakit.Int(1, 4)
		p1, p2 := edakit.Int(0, 30), edakit.Int(60, 100)
		A11, A22, A12 := float64(w1+w12), float64(w2+w12), float64(-w12)
		b1, b2 := float64(w1*p1), float64(w2*p2)
		det := A11*A22 - A12*A12
		x1 := (b1*A22 - A12*b2) / det
		x2 := (A11*b2 - A12*b1) / det
		k := edakit.Int(1, 2)
		ans := x1
		if k == 2 {
			ans = x2
		}
		return num(edakit.Sprintf("Cells 1 and 2 are movable. Cell 1 connects to a pad at x = %d with weight %d, cell 2 to a pad at x = %d with weight %d, and the cells connect to each other with weight %d. Minimizing quadratic wirelength, what is the position of cell %d, in um?", p1, w1, p2, w2, w12, k)+roundNote(2),
			f(ans, 2),
			edakit.Sprintf("Zero gradient gives the linear system $%dx_1 - %dx_2 = %d$ and $-%dx_1 + %dx_2 = %d$. Its matrix is the weighted Laplacian plus pad weights on the diagonal, positive definite, so the solution $x_1 = %s$, $x_2 = %s$ is the unique minimum.", w1+w12, w12, w1*p1, w12, w2+w12, w2*p2, f(x1, 2), f(x2, 2)),
			"x1", f(x1, 3), "x2", f(x2, 3))
	default:
		w12, w13, w23 := edakit.Int(1, 5), edakit.Int(0, 5), edakit.Int(1, 5)
		pads := []int{edakit.Int(0, 3), edakit.Int(0, 3), edakit.Int(0, 3)}
		padx := []int{edakit.Int(0, 100), edakit.Int(0, 100), edakit.Int(0, 100)}
		i := edakit.Int(0, 2)
		W := [3][3]int{{0, w12, w13}, {w12, 0, w23}, {w13, w23, 0}}
		diag := pads[i]
		for j := 0; j < 3; j++ {
			diag += W[i][j]
		}
		var desc []string
		desc = append(desc, edakit.Sprintf("cell1-cell2 weight %d, cell1-cell3 weight %d, cell2-cell3 weight %d", w12, w13, w23))
		for j := 0; j < 3; j++ {
			desc = append(desc, edakit.Sprintf("cell %d to a pad at x = %d with weight %d", j+1, padx[j], pads[j]))
		}
		which := edakit.Int(0, 2)
		stem := "Three movable cells have nets " + strings.Join(desc, "; ") + ". Quadratic placement solves $Ax = b$."
		switch which {
		case 0:
			return num(stem+edakit.Sprintf(" What is the diagonal entry $A_{%d%d}$?", i+1, i+1)+roundNote(0), edakit.IntStr(diag),
				edakit.Sprintf("A diagonal entry is the total weight of every net at that cell, including its pad: %d.", diag), "a_ii", edakit.IntStr(diag))
		case 1:
			j := (i + 1) % 3
			v := -W[i][j]
			return num(stem+edakit.Sprintf(" What is the off-diagonal entry $A_{%d%d}$?", i+1, j+1)+roundNote(0), edakit.IntStr(v),
				edakit.Sprintf("An off-diagonal entry is minus the weight between the two movable cells: %d. Pads add only to the diagonal and to b.", v), "a_ij", edakit.IntStr(v))
		default:
			v := pads[i] * padx[i]
			return num(stem+edakit.Sprintf(" What is the right-hand side entry $b_%d$?", i+1)+roundNote(0), edakit.IntStr(v),
				edakit.Sprintf("Fixed pins move to the right-hand side: $b_%d$ is the pad weight times the pad position, %d times %d = %d.", i+1, pads[i], padx[i], v), "b_i", edakit.IntStr(v))
		}
	}
}

func geometricProgramming(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	switch k {
	case 1, 2:
		a := float64(edakit.Int(1, 9)) / 2
		b := float64(edakit.Int(4, 80))
		x := math.Sqrt(b / a)
		v := 2 * math.Sqrt(a*b)
		stem := edakit.Sprintf("Sizing one gate by a factor x adds input load to its driver and divides its own drive resistance, giving a path delay $D(x) = %sx + %s/x$ ps for $x > 0$.", coefStr(a), f(b, 0))
		expl := edakit.Sprintf("$D'(x) = %s - %s/x^2 = 0$ gives $x^* = \\sqrt{%s/%s} = %s$, and then both terms are equal, $D^* = 2\\sqrt{%s \\cdot %s} = %s$ ps. In $u = \\ln x$ the posynomial becomes a convex sum of exponentials, so this stationary point is global.", f(a, 1), f(b, 0), f(b, 0), f(a, 1), f(x, 3), f(a, 1), f(b, 0), f(v, 3))
		if k == 1 {
			return num(stem+" What size $x^*$ minimizes the delay?"+roundNote(2), f(x, 2), expl, "x_opt", f(x, 3), "d_opt", f(v, 3))
		}
		return num(stem+" What is the minimum delay, in ps?"+roundNote(2), f(v, 2), expl, "x_opt", f(x, 3), "d_opt", f(v, 3))
	default:
		a, b, c := edakit.Int(1, 6), edakit.Int(1, 6), edakit.Int(5, 60)
		t := math.Cbrt(float64(a * b * c))
		v := 3 * t
		return num(edakit.Sprintf("A two-gate sizing model gives the posynomial $D(x, y) = %sx + %sy + %d/(xy)$ ps for sizes $x, y > 0$. What is the minimum of D, in ps?", coefStr(a), coefStr(b), c)+roundNote(2),
			f(v, 2),
			edakit.Sprintf("At the optimum the three terms are equal (set both partial derivatives to zero): $%dx = %dy = %d/(xy) = t$. Multiplying the three gives $t^3 = %d \\cdot %d \\cdot %d$, so $t = %s$ and $D^* = 3t = %s$ ps.", a, b, c, a, b, c, f(t, 3), f(v, 3)),
			"t", f(t, 3), "d_opt", f(v, 3))
	}
}

func basinOfAttraction(ctx generator.GeneratorContext) generator.Problem {
	rs := distinctInts(3, -6, 6)
	r1, r2, r3 := rs[0], rs[1], rs[2]
	if r1 > r2 {
		r1, r2 = r2, r1
	}
	if r2 > r3 {
		r2, r3 = r3, r2
	}
	if r1 > r2 {
		r1, r2 = r2, r1
	}
	var x0 float64
	for {
		x0 = float64(edakit.Int(-70, 70)) / 10
		if x0 != float64(r1) && x0 != float64(r2) && x0 != float64(r3) {
			break
		}
	}
	ans := edakit.Sprintf("x = %d", r1)
	if x0 > float64(r2) {
		ans = edakit.Sprintf("x = %d", r3)
	}
	opts := []string{edakit.Sprintf("x = %d", r1), edakit.Sprintf("x = %d", r2), edakit.Sprintf("x = %d", r3), "It diverges to infinity"}
	var dis []string
	for _, o := range opts {
		if o != ans {
			dis = append(dis, o)
		}
	}
	q := edakit.Sprintf("A one-dimensional density-penalty landscape f is a quartic with a positive leading term and $f'(x) = %s%s%s$. Gradient descent with a small step starts at $x_0 = %s$. Where does it converge?", lin(r1), lin(r2), lin(r3), f(x0, 1))
	return edakit.Choice(q, ans,
		edakit.Sprintf("f' is negative below %d, positive between %d and %d, negative between %d and %d, and positive above %d, so %d and %d are local minima and %d is a local maximum. Descent never crosses the maximum: from %s it slides to %s.", r1, r1, r2, r2, r3, r3, r1, r3, r2, f(x0, 1), ans),
		dis...)
}

// lin renders the factor (x - r) with a tidy sign.
func lin(r int) string {
	switch {
	case r == 0:
		return "x"
	case r > 0:
		return edakit.Sprintf("(x - %d)", r)
	default:
		return edakit.Sprintf("(x + %d)", -r)
	}
}

// poly renders sum c_i m_i with tidy signs, bare unit coefficients and zero terms dropped.
// An empty monomial is a constant term.
func poly(cs []int, ms []string) string {
	var b strings.Builder
	for i, c := range cs {
		if c == 0 {
			continue
		}
		a := abs(c)
		coef := edakit.IntStr(a)
		if a == 1 && ms[i] != "" {
			coef = ""
		}
		switch {
		case b.Len() == 0 && c < 0:
			b.WriteString("-")
		case b.Len() > 0 && c < 0:
			b.WriteString(" - ")
		case b.Len() > 0:
			b.WriteString(" + ")
		}
		b.WriteString(coef + ms[i])
	}
	if b.Len() == 0 {
		return "0"
	}
	return b.String()
}

// coefStr renders a coefficient in front of a variable, dropping a unit coefficient.
func coefStr[T int | float64](c T) string {
	if float64(c) == 1 {
		return ""
	}
	return f(float64(c), 1)
}
