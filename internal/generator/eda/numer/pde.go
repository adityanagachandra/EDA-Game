package numer

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var materials = []struct {
	name string
	k    float64
}{{"silicon", 150}, {"copper", 400}, {"a silicon-copper composite", 250}, {"a bonded silicon layer", 120}}

func finiteDifference(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		h := float64(edakit.Int(1, 10)) / 10
		u0 := float64(edakit.Int(300, 400)) / 10
		u1 := u0 + float64(edakit.Int(-30, 30))/10
		u2 := u1 + float64(edakit.Int(-30, 30))/10
		d2 := (u0 - 2*u1 + u2) / (h * h)
		ans := f(d2, 4)
		q := sp("Temperatures along a line at spacing $h = %s$ mm are %s, %s and %s °C. What is the central second-difference estimate of $T''$ at the middle point, in K/mm$^2$? Round to 4 decimal places.", f(h, 1), f(u0, 1), f(u1, 1), f(u2, 1))
		return numeric(q, ans, sp("$(T_{i-1} - 2T_i + T_{i+1})/h^2 = (%s)/%s = %s$ K/mm$^2$.", f(u0-2*u1+u2, 2), f(h*h, 2), ans), nil)
	case 2:
		m := edakit.Pick(materials)
		h := float64(edakit.Int(5, 20)) * 10 // um
		qv := float64(edakit.Int(1, 20))     // 1e9 W/m^3
		tn := []float64{}
		sum := 0.0
		for i := 0; i < 4; i++ {
			t := float64(edakit.Int(500, 900)) / 10
			tn = append(tn, t)
			sum += t
		}
		src := qv * 1e9 * (h * 1e-6) * (h * 1e-6) / m.k
		tc := (sum + src) / 4
		ans := f(tc, 4)
		q := sp("A steady thermal grid of %s ($k = %s$ W/(m K)) has spacing $h = %s$ µm. A node's four neighbours sit at %s, %s, %s and %s °C, and the node's cell dissipates $q = %s \\times 10^9$ W/m$^3$. The five-point stencil gives $T_c = (T_N + T_S + T_E + T_W + qh^2/k)/4$. What is $T_c$, in °C? Round to 4 decimal places.", m.name, f(m.k, 0), f(h, 0), f(tn[0], 1), f(tn[1], 1), f(tn[2], 1), f(tn[3], 1), f(qv, 0))
		return numeric(q, ans, sp("$qh^2/k = %s$ K. $T_c = (%s + %s)/4 = %s$ °C.", f(src, 4), f(sum, 1), f(src, 4), ans), map[string]string{"source_K": f(src, 4)})
	default:
		e := float64(edakit.Int(10, 900)) / 100
		r := edakit.Int(2, 4)
		ans := f(e/float64(r*r), 4)
		q := sp("A thermal solver with a second-order stencil shows a discretization error of %s K at the hot spot. The mesh spacing is divided by %d in every direction. What error should the refined solution show, in K? Round to 4 decimal places.", f(e, 2), r)
		return numeric(q, ans, sp("The truncation error of the second difference is $O(h^2)$, so the error falls by $%d^2$: %s K.", r, ans), nil)
	}
}

func boundaryConditions(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		L := float64(edakit.Int(2, 10))
		n := edakit.Int(4, 10)
		i := edakit.Int(1, n-1)
		ta := float64(edakit.Int(250, 450)) / 10
		tb := float64(edakit.Int(250, 450)) / 10
		Q := float64(edakit.Int(1, 40)) / 10
		x := L * float64(i) / float64(n)
		t := ta + (tb-ta)*x/L + Q*x*(L-x)/2
		ans := f(t, 4)
		q := sp("A rod of length %s mm with uniform heating obeys $-T'' = Q$ with $Q = q/k = %s$ K/mm$^2$, and its ends are held at %s °C and %s °C (Dirichlet). A finite-difference grid has %d equal intervals. What temperature does it compute at node %d (node 0 is the first end), in °C? Round to 4 decimal places.", f(L, 0), f(Q, 1), f(ta, 1), f(tb, 1), n, i)
		return numeric(q, ans, sp("The exact solution $T = T_a + (T_b - T_a)x/L + Qx(L - x)/2$ is quadratic, so the second difference reproduces it exactly at the nodes. At $x = %s$ mm: %s °C.", f(x, 4), ans), map[string]string{"x": f(x, 4)})
	case 2:
		h := float64(edakit.Int(1, 10)) / 10
		t := float64(edakit.Int(400, 900)) / 10
		Q := float64(edakit.Int(1, 50)) / 10
		tn := t + Q*h*h/2
		ans := f(tn, 4)
		q := sp("The end node $N$ of a 1D thermal grid is insulated ($T' = 0$). With a ghost node $T_{N+1} = T_{N-1}$, the stencil at $N$ reads $2(T_{N-1} - T_N)/h^2 + Q = 0$. With $h = %s$ mm, $Q = %s$ K/mm$^2$ and $T_{N-1} = %s$ °C, what is $T_N$, in °C? Round to 4 decimal places.", f(h, 1), f(Q, 1), f(t, 1))
		return numeric(q, ans, sp("$T_N = T_{N-1} + Qh^2/2 = %s + %s = %s$ °C.", f(t, 1), f(Q*h*h/2, 4), ans), nil)
	default:
		k := float64(edakit.Pick([]int{1, 2, 5, 20, 150}))
		dx := float64(edakit.Int(1, 20)) / 10 // mm
		hc := float64(edakit.Int(1, 50)) * 100
		t1 := float64(edakit.Int(500, 1000)) / 10
		tamb := float64(edakit.Int(200, 400)) / 10
		g := k / (dx * 1e-3)
		t0 := (g*t1 + hc*tamb) / (g + hc)
		ans := f(t0, 4)
		q := sp("The surface of a package layer ($k = %s$ W/(m K)) is cooled by convection to %s °C with heat-transfer coefficient $h_c = %s$ W/(m$^2$ K). The first interior node, %s mm below the surface node, is at %s °C. The one-sided Robin condition $k(T_1 - T_0)/\\Delta x = h_c(T_0 - T_{amb})$ fixes the surface temperature. What is $T_0$, in °C? Round to 4 decimal places.", f(k, 0), f(tamb, 1), f(hc, 0), f(dx, 1), f(t1, 1))
		return numeric(q, ans, sp("$k/\\Delta x = %s$ W/(m$^2$ K). $T_0 = ((k/\\Delta x)T_1 + h_cT_{amb})/(k/\\Delta x + h_c) = %s$ °C.", f(g, 1), ans), map[string]string{"g": f(g, 1)})
	}
}

func explicitHeat(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		alpha := float64(edakit.Int(10, 120))
		h := float64(edakit.Int(1, 10)) / 10
		t0 := float64(edakit.Int(400, 800)) / 10
		t1 := t0 + float64(edakit.Int(-50, 50))/10
		t2 := t1 + float64(edakit.Int(-50, 50))/10
		r := alpha * (t0 - 2*t1 + t2) / (h * h)
		ans := f(r, 4)
		q := sp("The method of lines turns $T_t = \\alpha T_{xx}$ into $dT_i/dt = \\alpha(T_{i-1} - 2T_i + T_{i+1})/h^2$. With $\\alpha = %s$ mm$^2$/s, $h = %s$ mm and neighbouring temperatures %s, %s, %s °C, what is $dT_i/dt$ at the middle node, in K/s? Round to 4 decimal places.", f(alpha, 0), f(h, 1), f(t0, 1), f(t1, 1), f(t2, 1))
		return numeric(q, ans, sp("$%s \\cdot (%s)/%s = %s$ K/s.", f(alpha, 0), f(t0-2*t1+t2, 2), f(h*h, 2), ans), nil)
	case 2:
		alpha := float64(edakit.Int(10, 120))
		h := float64(edakit.Int(5, 100)) // um
		d := edakit.Int(1, 3)
		dt := (h * 1e-3) * (h * 1e-3) / (2 * float64(d) * alpha) * 1e6
		ans := f(dt, 4)
		q := sp("Forward Euler in time on a %dD heat grid with the standard second-difference stencil is stable only if $\\alpha\\,\\Delta t/h^2 \\le 1/(2d)$. For $\\alpha = %s$ mm$^2$/s and $h = %s$ µm, what is the largest stable $\\Delta t$, in µs? Round to 4 decimal places.", d, f(alpha, 0), f(h, 0))
		return numeric(q, ans, sp("$\\Delta t \\le h^2/(2d\\alpha) = (%s\\text{ mm})^2/(%d \\cdot %s) = %s$ µs. Halving $h$ quarters it.", f(h/1000, 3), 2*d, f(alpha, 0), ans), nil)
	default:
		r := float64(edakit.Int(5, 50)) / 100
		t0 := float64(edakit.Int(400, 800)) / 10
		t1 := t0 + float64(edakit.Int(-50, 50))/10
		t2 := t1 + float64(edakit.Int(-50, 50))/10
		tn := t1 + r*(t0-2*t1+t2)
		ans := f(tn, 4)
		q := sp("One forward-time centred-space step updates $T_i \\leftarrow T_i + r(T_{i-1} - 2T_i + T_{i+1})$ with $r = \\alpha\\Delta t/h^2 = %s$. The node and its neighbours are at $T_{i-1} = %s$, $T_i = %s$, $T_{i+1} = %s$ °C. What is the new $T_i$, in °C? Round to 4 decimal places.", f(r, 2), f(t0, 1), f(t1, 1), f(t2, 1))
		return numeric(q, ans, sp("$%s + %s(%s) = %s$ °C. Keeping $r \\le 1/2$ makes the new value a convex combination of old ones, which is why the step is stable.", f(t1, 1), f(r, 2), f(t0-2*t1+t2, 2), ans), nil)
	}
}

func fem(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		k := float64(edakit.Pick([]int{1, 2, 20, 120, 150, 400}))
		a := float64(edakit.Int(1, 20)) / 10 // mm^2
		h := float64(edakit.Int(1, 20)) / 10 // mm
		g := k * a / h
		ans := f(-g, 4)
		q := sp("A 1D linear finite element of length %s mm, cross-section %s mm$^2$ and conductivity %s W/(m K) has element matrix $(kA/h)\\begin{pmatrix} 1 & -1 \\\\ -1 & 1 \\end{pmatrix}$. What is its off-diagonal entry, in mW/K? Round to 4 decimal places.", f(h, 1), f(a, 1), f(k, 0))
		return numeric(q, ans, sp("$kA/h = %s \\cdot %s\\text{ mm}^2/%s\\text{ mm} = %s$ mW/K, so the off-diagonal entry is %s.", f(k, 0), f(a, 1), f(h, 1), f(g, 4), ans), nil)
	case 2:
		g := []float64{float64(edakit.Int(1, 30)), float64(edakit.Int(1, 30)), float64(edakit.Int(1, 30))}
		// K is 4x4 tridiagonal: K11=g1, K22=g1+g2, K33=g2+g3, K44=g3, K12=-g1, K23=-g2, K34=-g3.
		type ent struct {
			name string
			v    float64
		}
		es := []ent{{"K_{22}", g[0] + g[1]}, {"K_{33}", g[1] + g[2]}, {"K_{11}", g[0]}, {"K_{23}", -g[1]}, {"K_{34}", -g[2]}, {"K_{13}", 0}}
		e := edakit.Pick(es)
		ans := f(e.v, 4)
		q := sp("Three linear elements in series join nodes 1-2, 2-3 and 3-4, with conductances $kA/h$ of %s, %s and %s mW/K. After assembly of the global stiffness matrix, what is $%s$, in mW/K?", f(g[0], 0), f(g[1], 0), f(g[2], 0), e.name)
		return numeric(q, ans, sp("Each element adds $g$ to its two diagonal entries and $-g$ to its off-diagonal pair. The global matrix is tridiagonal with diagonal (%s, %s, %s, %s); $%s = %s$.", f(g[0], 0), f(g[0]+g[1], 0), f(g[1]+g[2], 0), f(g[2], 0), e.name, ans), nil)
	default:
		g1, g2 := float64(edakit.Int(1, 30)), float64(edakit.Int(1, 30))
		tl, tr := float64(edakit.Int(250, 450))/10, float64(edakit.Int(250, 450))/10
		p := float64(edakit.Int(1, 500))
		t := (g1*tl + g2*tr + p) / (g1 + g2)
		ans := f(t, 4)
		q := sp("Two linear elements with $kA/h$ of %s and %s mW/K join a left end held at %s °C, a middle node, and a right end held at %s °C. A via dissipates %s mW at the middle node (a point load). Assemble and solve: what is the middle temperature, in °C? Round to 4 decimal places.", f(g1, 0), f(g2, 0), f(tl, 1), f(tr, 1), f(p, 0))
		return numeric(q, ans, sp("The middle row reads $(g_1 + g_2)T_m - g_1T_L - g_2T_R = P$, so $T_m = (%s + %s + %s)/%s = %s$ °C.", f(g1*tl, 1), f(g2*tr, 1), f(p, 0), f(g1+g2, 0), ans), nil)
	}
}

func fvm(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		k := float64(edakit.Pick([]int{1, 2, 20, 120, 150, 400}))
		a := float64(edakit.Int(1, 50)) / 100 // mm^2
		d := float64(edakit.Int(1, 20)) / 100 // mm
		ti := float64(edakit.Int(500, 900)) / 10
		tj := ti - float64(edakit.Int(1, 100))/10
		qf := k * a / d * (ti - tj) // mW
		ans := f(qf, 4)
		q := sp("Two control volumes of a material with $k = %s$ W/(m K) share a face of area %s mm$^2$; their centres are %s mm apart and sit at %s and %s °C. What heat flows through the face, in mW? Round to 4 decimal places.", f(k, 0), f(a, 2), f(d, 2), f(ti, 1), f(tj, 1))
		return numeric(q, ans, sp("$q = kA(T_i - T_j)/d = %s \\cdot %s/%s \\cdot %s = %s$ mW (with mm units, W/(m K) times mm gives mW/K).", f(k, 0), f(a, 2), f(d, 2), f(ti-tj, 1), ans), nil)
	case 2:
		k1 := float64(edakit.Pick([]int{1, 2, 5, 20}))
		k2 := float64(edakit.Pick([]int{120, 150, 400}))
		if edakit.Int(0, 1) == 0 {
			ke := 2 * k1 * k2 / (k1 + k2)
			ans := f(ke, 4)
			q := sp("A finite-volume face separates a cell of conductivity %s W/(m K) from a cell of %s W/(m K); both cell centres are the same distance from the face. What effective face conductivity makes the flux exact for this series path, in W/(m K)? Round to 4 decimal places.", f(k1, 0), f(k2, 0))
			return numeric(q, ans, sp("Series resistances give the harmonic mean $2k_1k_2/(k_1 + k_2) = %s$ W/(m K), close to the smaller value, not the arithmetic mean %s.", ans, f((k1+k2)/2, 1)), nil)
		}
		d1, d2 := float64(edakit.Int(1, 20))/100, float64(edakit.Int(1, 20))/100
		a := float64(edakit.Int(1, 50)) / 100
		g := a / (d1/k1 + d2/k2)
		ans := f(g, 4)
		q := sp("A face of area %s mm$^2$ separates a cell of conductivity %s W/(m K), whose centre is %s mm from the face, from a cell of %s W/(m K) whose centre is %s mm from it. What is the face conductance $A/(d_1/k_1 + d_2/k_2)$, in mW/K? Round to 4 decimal places.", f(a, 2), f(k1, 0), f(d1, 2), f(k2, 0), f(d2, 2))
		return numeric(q, ans, sp("$%s/(%s/%s + %s/%s) = %s$ mW/K.", f(a, 2), f(d1, 2), f(k1, 0), f(d2, 2), f(k2, 0), ans), nil)
	default:
		n := edakit.Int(3, 5)
		sg, sgt := 0.0, 0.0
		var parts []string
		for i := 0; i < n; i++ {
			g := float64(edakit.Int(1, 40)) / 10
			t := float64(edakit.Int(400, 800)) / 10
			sg += g
			sgt += g * t
			parts = append(parts, sp("%s mW/K to a cell at %s °C", f(g, 1), f(t, 1)))
		}
		p := float64(edakit.Int(1, 200))
		tp := (sgt + p) / sg
		ans := f(tp, 4)
		q := "A finite-volume cell on a nonuniform mesh has face conductances " + joinSemi(parts) + sp(". The cell dissipates %s mW. Balancing fluxes, what is the cell temperature, in °C? Round to 4 decimal places.", f(p, 0))
		return numeric(q, ans, sp("$\\sum G_f(T_P - T_f) = P$ gives $T_P = (\\sum G_fT_f + P)/\\sum G_f = (%s + %s)/%s = %s$ °C.", f(sgt, 2), f(p, 0), f(sg, 1), ans), nil)
	}
}

func joinSemi(xs []string) string {
	out := ""
	for i, x := range xs {
		if i > 0 {
			out += "; "
		}
		out += x
	}
	return out
}

var bemItems = []edakit.Item{
	{KP: 3, Q: "A BEM capacitance solver uses the fast multipole method inside GMRES. What does the multipole method change?",
		A: "Each matrix-vector product with the dense panel matrix costs about O(N) instead of O(N^2)", D: []string{"It makes the panel matrix sparse by deleting far interactions", "It replaces GMRES by a direct factorization", "It removes the need for surface panels"},
		Why: "Far-field interactions are aggregated into multipole expansions, so the product is approximated to a set accuracy in near-linear time; the matrix is never formed."},
	{KP: 3, Q: "Why is BEM attractive for interconnect capacitance compared with FEM on a volume mesh?",
		A: "Only conductor surfaces are discretized, and the unbounded dielectric region needs no outer mesh", D: []string{"The BEM matrix is sparse", "BEM needs no Green's function", "BEM handles any inhomogeneous dielectric with no extra work"},
		Why: "The Green's function carries the field through the open region analytically; the cost is a dense system on the surface unknowns."},
	{KP: 3, Q: "Which compression idea lets a hierarchical-matrix BEM solver store the far-field blocks cheaply?",
		A: "Blocks coupling well-separated panel clusters are numerically low rank", D: []string{"Far-field blocks are exactly zero", "Near-field blocks are low rank", "All blocks are diagonal after reordering"},
		Why: "The kernel 1/r is smooth for separated clusters, so those blocks admit low-rank factorizations such as adaptive cross approximation."},
	{KP: 3, Q: "A multilayer dielectric stack complicates BEM. What is the standard remedy?",
		A: "Add panels on the dielectric interfaces, or use a layered-media Green's function", D: []string{"Ignore the layers and use the free-space kernel", "Switch to the pure Neumann problem", "Mesh only the top layer"},
		Why: "The free-space Green's function assumes one permittivity, so interfaces need their own unknowns or a Green's function that already includes the layers."},
	{KP: 3, Q: "Floating random-walk extraction is a third family besides BEM and FEM. What is its distinguishing property?",
		A: "It estimates capacitance statistically without assembling or solving a matrix, with error falling as one over the square root of the number of walks", D: []string{"It is exact after one walk", "It solves the dense BEM system by random pivoting", "It needs a volume mesh of the dielectric"},
		Why: "Random-walk methods sample the Green's-function representation of the potential; memory is tiny and accuracy is set by the number of walks."},
	{KP: 3, Q: "A BEM solver doubles the panel density on every surface of a 3D structure. Roughly how does the cost of a dense direct LU solve change?",
		A: "N grows by 4, so LU work grows by about 64", D: []string{"It doubles", "It grows by 4", "It grows by 16"},
		Why: "Doubling density in two surface dimensions multiplies N by 4, and dense LU costs O(N^3), so the work grows by 4^3 = 64."},
}

func bemPanels(ctx generator.GeneratorContext) generator.Problem {
	c := edakit.Int(2, 4)
	p, qq, r := edakit.Int(2, 10), edakit.Int(2, 10), edakit.Int(1, 4)
	per := 2 * (p*qq + p*r + qq*r)
	ans := c * per
	q := sp("A BEM capacitance extractor models %d identical box-shaped wires. Each box is meshed with %d-by-%d panels on its top and bottom faces, %d-by-%d on its front and back, and %d-by-%d on its two ends. How many unknowns (panel charges) does the system have?", c, p, qq, p, r, qq, r)
	return edakit.Choice(q, itoa(ans), sp("Each box has $2(%d + %d + %d) = %d$ panels; %d boxes give %d unknowns, one charge per panel.", p*qq, p*r, qq*r, per, c, ans),
		edakit.Distinct(ans, 3, c*p*qq*2, per, c*(p*qq+p*r+qq*r), c*6*p*qq)...)
}

func bemCost(ctx generator.GeneratorContext) generator.Problem {
	n := 1000 * edakit.Scaled(ctx, 2, 60)
	if edakit.Int(0, 1) == 0 {
		mb := 8 * float64(n) * float64(n) / 1e6
		ans := f(mb, 0) + " MB"
		q := sp("A BEM extraction has N = %d panels and forms the dense panel matrix in 8-byte doubles. How much memory does the matrix take?", n)
		return edakit.Choice(q, ans, sp("$8N^2 = %s$ bytes, about %s. Memory, not just time, forces fast methods for large N.", f(8*float64(n)*float64(n), 0), ans),
			f(mb/1000, 0)+" MB", f(8*float64(n)/1e6*1000, 0)+" MB", f(mb*8, 0)+" MB", f(mb/8, 0)+" MB")
	}
	k := edakit.Int(10, 60)
	lu := 2.0 / 3 * math.Pow(float64(n), 3)
	gm := float64(k) * 2 * float64(n) * float64(n)
	ratio := lu / gm
	ans := f(ratio, 0)
	q := sp("A BEM system with N = %d panels can be solved by dense LU ($2N^3/3$ flops) or by %d GMRES iterations, each a dense matrix-vector product ($2N^2$ flops). Ignoring orthogonalization, about how many times more flops does LU need? Round to the nearest integer.", n, k)
	return edakit.Choice(q, ans, sp("$(2N^3/3)/(%d \\cdot 2N^2) = N/(3 \\cdot %d) = %s$.", k, k, ans),
		distinctFloats(ans, 0, ratio*3, ratio/2, float64(n)/float64(k), ratio*2)...)
}
