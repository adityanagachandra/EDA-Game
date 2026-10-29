package extract

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// fieldSolverCalc reads one column of a solver run (one conductor at 1 V, the rest grounded)
// and asks for a coupling, a total or a ground capacitance.
func fieldSolverCalc(ctx generator.GeneratorContext) generator.Problem {
	j := edakit.Int(1, 3)
	cg := float64(edakit.Int(3, 15)) / 2
	ca := float64(edakit.Int(2, 12)) / 4
	cb := float64(edakit.Int(1, 12)) / 4
	others := []int{}
	for k := 1; k <= 3; k++ {
		if k != j {
			others = append(others, k)
		}
	}
	qj := cg + ca + cb
	q := edakit.Sprintf("A 3D field solver holds conductor %d at 1 V and conductors %d and %d at 0 V, together with the ground plane. It reports induced charges of %s fC on conductor %d, %s fC on conductor %d and %s fC on conductor %d. ",
		j, others[0], others[1], f(qj, 2), j, f(-ca, 2), others[0], f(-cb, 2), others[1])
	switch edakit.Int(1, 3) {
	case 1:
		return numChoice(q+edakit.Sprintf("What is the coupling capacitance between conductors %d and %d?", j, others[0]), ca, 2, " fF",
			edakit.Sprintf("With conductor %d at 1 V, the charge on grounded conductor %d is $-C_c \\cdot 1$ V, so the coupling is %s fF.", j, others[0], f(ca, 2)),
			-ca, qj, qj-ca, ca+cb)
	case 2:
		return numChoice(q+edakit.Sprintf("What is the total capacitance of conductor %d?", j), qj, 2, " fF",
			edakit.Sprintf("The charge on the driven conductor per volt is its total capacitance, the diagonal entry $C_{%d%d}$ = %s fF.", j, j, f(qj, 2)),
			cg, ca+cb, qj+ca+cb, qj-ca)
	}
	return numChoice(q+edakit.Sprintf("What is the capacitance from conductor %d to the ground plane?", j), cg, 2, " fF",
		edakit.Sprintf("Total charge is conserved: the ground plane carries the rest. Ground capacitance is total minus couplings, %s - %s - %s = %s fF.", f(qj, 2), f(ca, 2), f(cb, 2), f(cg, 2)),
		qj, ca+cb, qj+ca+cb, qj-ca)
}

var fieldSolverItems = []edakit.Item{
	{KP: 1, Q: "What equation does an electrostatic capacitance field solver solve in the dielectric between conductors?",
		A: "Laplace's equation for the potential, with conductor potentials as boundary conditions", D: []string{"The telegrapher's equations along each wire", "Kirchhoff's current law at every node", "The heat equation with the wires as sources"},
		Why: "With no free charge in the dielectric, the potential satisfies the Laplace (or, with varying permittivity, the generalized Poisson) equation. The charges on conductors follow from the field at their surfaces."},
	{KP: 1, Q: "How does a solver build the whole capacitance matrix of N conductors?",
		A: "Run N solves, each with one conductor at 1 V and the rest at 0 V, and record the induced charges as one column", D: []string{"Run one solve with every conductor at 1 V", "Run N squared solves, one per pair of conductors", "Measure the current into each conductor under a step input"},
		Why: "Charge is linear in the conductor voltages, $Q = C V$. Exciting one conductor at a time reads out C one column per solve."},
	{KP: 1, Q: "Why does a design flow call field solvers golden rather than using them on every net?",
		A: "They are accurate to the geometry but far too slow for billions of wire segments, so they calibrate faster methods", D: []string{"They ignore fringe fields and are only approximate", "They cannot compute coupling capacitance", "They only handle a single conductor at a time"},
		Why: "A 3D solve is expensive. Production extraction uses rule-based models whose tables were built and checked against solver results."},
	{KP: 1, Q: "A solver reports a capacitance matrix with a positive off-diagonal entry. What does this indicate?",
		A: "An error: in the Maxwell matrix off-diagonal entries are non-positive", D: []string{"Strong coupling between those two conductors", "A grounded shield between them", "That the conductors are on different layers"},
		Why: "Raising one conductor's voltage induces negative charge on its grounded neighbours, so $C_{ij} \\le 0$ for $i \\ne j$. A positive entry signals a meshing or convergence problem."},
	{KP: 1, Q: "Which property of the Maxwell capacitance matrix is a cheap sanity check on solver output?",
		A: "It is symmetric, so C_ij and C_ji must agree within the solver's accuracy", D: []string{"Its diagonal entries are all equal", "Its determinant is always 1", "Each row sums to zero"},
		Why: "Reciprocity makes the matrix symmetric. Rows sum to the ground capacitance, which is zero only if there is no ground conductor."},
	{KP: 1, Q: "Why does a capacitance solver only need the geometry and permittivities, and not the frequency?",
		A: "At on-chip dimensions far below the wavelength, the quasi-static electric field decouples from the magnetic field", D: []string{"Because capacitance is defined only at DC by convention", "Because dielectrics have no frequency dependence", "Because the solver averages over all frequencies"},
		Why: "When structures are electrically small, capacitance is an electrostatic quantity. Frequency effects enter through inductance and resistance (skin effect), extracted separately."},
	{KP: 3, Q: "A 3D solver's coupling result changes by 4% when the mesh is refined once and by 0.3% on the next refinement. What should the extraction engineer conclude?",
		A: "The first mesh was under-resolved; the refined value has converged to about 0.3%", D: []string{"The solver is unstable and cannot be trusted", "The 4% change is the physical variation of the process", "Mesh refinement never changes capacitance"},
		Why: "Discretization error falls as the mesh is refined. A golden reference must be shown converged by refinement before it is used to calibrate anything."},
	{KP: 3, Q: "Why do production extractors use 2.5D pattern tables rather than 3D solves for most nets?",
		A: "Most wiring is locally a cross-section that repeats along the wire, so 2D solves tabulated by width and spacing capture it cheaply", D: []string{"3D effects do not exist on chips", "2D solvers are more accurate than 3D", "SPEF cannot record 3D results"},
		Why: "Long parallel runs look the same in every cross-section. Corners, line ends and vias are the 3D exceptions that need special patterns or a 3D solve."},
	{KP: 3, Q: "A critical analog net needs extraction within 1% of silicon. The production flow's rule-based extractor is within 5% on average. What is the standard response?",
		A: "Extract that net with a 3D field solver while the rest of the chip stays rule-based", D: []string{"Re-extract the whole chip with a 3D solver", "Add 5% margin to every net", "Use the rule-based value, since averages cancel"},
		Why: "Field solvers are reserved for the nets where accuracy matters most. Mixed-mode extraction applies them selectively."},
	{KP: 3, Q: "A field solver is given a dielectric stack that omits a thin etch-stop layer of higher permittivity. Which way will its coupling and ground capacitances err?",
		A: "Low, because a higher-permittivity layer in the field path is missing", D: []string{"High, because thinner stacks give more capacitance", "Not at all, because thin layers do not matter", "Coupling high and ground low"},
		Why: "Capacitance grows with the permittivity of the material the field passes through. The solver is only golden if its process description matches the process."},
	{KP: 3, Q: "Which representation of the solver's result lets a timer compute crosstalk on a victim net?",
		A: "Separate coupling capacitances to each neighbour, not only the victim's total capacitance", D: []string{"The victim's total capacitance alone", "The ground capacitance with coupling dropped", "The average capacitance per micron of the layer"},
		Why: "Crosstalk depends on which neighbour switches and when, so the coupling to each aggressor must be kept as a distinct element."},
	{KP: 3, Q: "Two adjacent conductors each have 10 fF to ground and 5 fF to each other. A tool reports C_11 = 15 fF and C_12 = 5 fF. What went wrong?",
		A: "The sign convention: C_12 should be -5 fF in the Maxwell matrix", D: []string{"C_11 should be 10 fF", "C_11 should be 20 fF", "Nothing; both are correct"},
		Why: "The diagonal is total capacitance, 15 fF, and the off-diagonal is the negative of the coupling. Mixing up Maxwell and mutual conventions is a common integration bug."},
}

// patternMatching interpolates a capacitance table by width and spacing.
func patternMatching(ctx generator.GeneratorContext) generator.Problem {
	ws := []float64{0.02, 0.03, 0.04, 0.06, 0.08}
	ss := []float64{0.02, 0.03, 0.04, 0.06, 0.08, 0.12}
	i := edakit.Int(0, len(ws)-2)
	j := edakit.Int(0, len(ss)-2)
	w1, w2, s1, s2 := ws[i], ws[i+1], ss[j], ss[j+1]
	// Coupling per length falls with spacing and rises a little with width.
	v11 := float64(edakit.Int(60, 120))
	v12 := v11 - float64(edakit.Int(15, 40))
	v21 := v11 + float64(edakit.Int(3, 12))
	v22 := v12 + (v21 - v11) + float64(edakit.Int(-2, 2))
	interp := func(w, s float64) float64 {
		tw := (w - w1) / (w2 - w1)
		ts := (s - s1) / (s2 - s1)
		return v11*(1-tw)*(1-ts) + v12*(1-tw)*ts + v21*tw*(1-ts) + v22*tw*ts
	}
	table := edakit.Sprintf("A rule-based extractor's coupling table (aF/um of parallel run) has, at width %s um: %s aF/um at spacing %s um and %s aF/um at spacing %s um; at width %s um: %s aF/um at spacing %s um and %s aF/um at spacing %s um. It interpolates linearly in each index. ",
		f(w1, 2), f(v11, 0), f(s1, 2), f(v12, 0), f(s2, 2), f(w2, 2), f(v21, 0), f(s1, 2), f(v22, 0), f(s2, 2))
	frac := func() float64 { return float64(edakit.Int(1, 3)) / 4 }
	l := float64(10 * edakit.Scaled(ctx, 3, 30))
	switch kp(ctx) {
	case 1:
		s := s1 + (s2-s1)*frac()
		v := interp(w1, s)
		return num(table+edakit.Sprintf("A wire of width %s um runs %s um beside a neighbour at spacing %s um. What coupling per length does the extractor use, in aF/um?", f(w1, 2), f(l, 0), f(s, 4)),
			v, 2, edakit.Sprintf("On the width-%s row, interpolate in spacing: $t = %s$, giving %s + %s x (%s - %s) = %s aF/um.", f(w1, 2), f((s-s1)/(s2-s1), 2), f(v11, 0), f((s-s1)/(s2-s1), 2), f(v12, 0), f(v11, 0), f(v, 2)), nil)
	case 2:
		w := w1 + (w2-w1)*frac()
		s := s1 + (s2-s1)*frac()
		v := interp(w, s)
		c := v * l / 1000
		return num(table+edakit.Sprintf("A wire of width %s um runs %s um beside a neighbour at spacing %s um. What coupling capacitance does the extractor report, in fF?", f(w, 4), f(l, 0), f(s, 4)),
			c, 3, edakit.Sprintf("Bilinear interpolation gives %s aF/um at width %s and spacing %s. Times %s um of run: %s fF.", f(v, 2), f(w, 4), f(s, 4), f(l, 0), f(c, 3)),
			map[string]string{"per_len": f(v, 2)})
	}
	la := float64(10 * edakit.Scaled(ctx, 2, 20))
	lb := float64(10 * edakit.Scaled(ctx, 2, 20))
	sa := s1 + (s2-s1)*frac()
	va, vb := interp(w1, sa), interp(w1, s2)
	c := (va*la + vb*lb) / 1000
	return num(table+edakit.Sprintf("A wire of width %s um runs %s um beside a neighbour at spacing %s um, then the neighbour jogs away and the next %s um run at spacing %s um. The extractor fractures the wire into the two segments and looks up each. What total coupling to that neighbour does it report, in fF?", f(w1, 2), f(la, 0), f(sa, 4), f(lb, 0), f(s2, 2)),
		c, 3, edakit.Sprintf("Segment 1: %s aF/um x %s um. Segment 2 sits on a table point: %s aF/um x %s um. Sum: %s fF. Pattern extraction fractures geometry into pieces whose cross-section matches a table entry.", f(va, 2), f(la, 0), f(vb, 0), f(lb, 0), f(c, 3)), nil)
}

// randomWalk estimates a capacitance from walk scores, its standard error, and the walk count
// needed for a target accuracy.
func randomWalk(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 6, 10)
		var xs []string
		sum := 0.0
		for k := 0; k < n; k++ {
			v := float64(edakit.Int(0, 40)) / 10
			sum += v
			xs = append(xs, f(v, 1))
		}
		m := sum / float64(n)
		return num(edakit.Sprintf("A floating random-walk solver estimates the coupling between conductors i and j. Each walk starts on a surface around conductor i, hops through the dielectric, and contributes a score (in fF) when it lands on conductor j, and zero otherwise. %d walks scored %s. What is the coupling estimate, in fF?", n, strings.Join(xs, ", ")),
			m, 3, edakit.Sprintf("The estimate is the sample mean of the walk scores: %s / %d = %s fF. Each walk is independent, so the estimate is unbiased and its error shrinks as walks are added.", f(sum, 1), n, f(m, 3)), nil)
	case 2:
		mean := float64(edakit.Int(5, 40)) / 2
		sd := mean * float64(edakit.Int(5, 30)) / 10
		n := edakit.Pick([]int{1000, 2500, 4000, 10000, 40000})
		rel := 100 * sd / math.Sqrt(float64(n)) / mean
		return num(edakit.Sprintf("After %d walks, a random-walk capacitance solver has a sample mean of %s fF and a sample standard deviation of the per-walk score of %s fF. What is the relative standard error of the estimate, in percent?", n, f(mean, 1), f(sd, 2)),
			rel, 3, edakit.Sprintf("The standard error is $\\sigma/\\sqrt{N}$ = %s / %s = %s fF, which is %s%% of the mean.", f(sd, 2), f(math.Sqrt(float64(n)), 2), f(sd/math.Sqrt(float64(n)), 4), f(rel, 3)), nil)
	}
	mean := float64(edakit.Int(5, 40)) / 2
	sd := mean * float64(edakit.Int(5, 30)) / 10
	tgt := edakit.Pick([]float64{0.5, 1, 2})
	ratio := sd / (tgt / 100 * mean)
	need := int(math.Ceil(ratio*ratio - 1e-9))
	return num(edakit.Sprintf("A random-walk solver's per-walk scores for one coupling have mean %s fF and standard deviation %s fF. The flow requires a relative standard error of at most %s%%. How many walks are needed, at minimum?", f(mean, 1), f(sd, 2), f(tgt, 1)),
		float64(need), 0, edakit.Sprintf("Need $\\sigma/\\sqrt{N} \\le %s\\%%$ of the mean, so $N \\ge (\\sigma / (\\varepsilon \\mu))^2$ = (%s / %s)^2 = %s; round up to %d. Halving the error costs four times the walks, which is why random-walk solvers expose accuracy as a parameter.", f(tgt, 1), f(sd, 2), f(tgt/100*mean, 4), f(ratio*ratio, 2), need), nil)
}

// bemFem compares dense boundary-element and sparse finite-element solver costs.
func bemFem(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		p := edakit.Pick([]float64{0.01, 0.02, 0.05})
		lw := edakit.Scaled(ctx, 10, 60)
		ww := edakit.Int(1, 4)
		tw := edakit.Int(1, 4)
		l, w, t := float64(lw)*p, float64(ww)*p, float64(tw)*p
		n := 2 * (lw*ww + lw*tw + ww*tw)
		return num(edakit.Sprintf("A boundary-element solver meshes only conductor surfaces, into square panels of side %s um. A wire is a box %s um long, %s um wide and %s um thick. How many panels cover its whole surface?", f(p, 2), f(l, 2), f(w, 2), f(t, 2)),
			float64(n), 0, edakit.Sprintf("Surface area is $2(LW + LT + WT)$; in panel units that is 2 x (%d x %d + %d x %d + %d x %d) = %d panels. A finite-element mesh of the surrounding volume would need far more unknowns.", lw, ww, lw, tw, ww, tw, n), nil)
	case 2:
		n := 1000 * edakit.Scaled(ctx, 2, 40)
		mb := float64(n) * float64(n) * 8 / 1e6
		return num(edakit.Sprintf("A boundary-element capacitance solve has %d panels. Its interaction matrix is dense, stored in double precision (8 bytes per entry). How much memory does the matrix take, in MB (1 MB is one million bytes)?", n),
			mb, 1, edakit.Sprintf("$N^2 \\times 8$ = %d^2 x 8 = %s MB. Dense storage and an $O(N^3)$ direct solve are why practical BEM extractors use iterative solvers with fast multipole or similar acceleration.", n, f(mb, 1)), nil)
	}
	n := 1000 * edakit.Scaled(ctx, 2, 20)
	mult := edakit.Pick([]int{20, 30, 50})
	nnz := edakit.Pick([]int{7, 15, 27})
	m := n * mult
	ratio := float64(n) * float64(n) / (float64(m) * float64(nnz))
	return num(edakit.Sprintf("For the same structure, a boundary-element solver uses %d surface panels with a dense matrix, while a finite-element solver meshes the volume with %d times as many unknowns but a sparse matrix averaging %d nonzeros per row. What is the ratio of stored entries, BEM over FEM?", n, mult, nnz),
		ratio, 2, edakit.Sprintf("BEM stores $N^2$ = %s entries; FEM stores %d x %d = %s. Ratio %s. BEM has fewer unknowns but dense coupling; FEM has many unknowns but sparse, local coupling and handles complex dielectric stacks naturally.", f(float64(n)*float64(n), 0), m, nnz, f(float64(m)*float64(nnz), 0), f(ratio, 2)), nil)
}

// calibration fits a rule-based extractor's output to golden solver values.
func calibration(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 6)
	k := 1 + float64(edakit.Int(-8, 8))/100
	xs := make([]float64, n)
	gs := make([]float64, n)
	var rows []string
	for i := 0; i < n; i++ {
		xs[i] = float64(edakit.Int(20, 120)) / 10
		gs[i] = edakit.Round(xs[i]*k*(1+float64(edakit.Int(-3, 3))/100), 2)
		rows = append(rows, edakit.Sprintf("net %d: fast %s fF, golden %s fF", i+1, f(xs[i], 1), f(gs[i], 2)))
	}
	stem := "A fast rule-based extractor is compared with a 3D field solver on a set of nets: " + strings.Join(rows, "; ") + ". "
	sxy, sxx := 0.0, 0.0
	for i := range xs {
		sxy += xs[i] * gs[i]
		sxx += xs[i] * xs[i]
	}
	scale := sxy / sxx
	maxErr := func(s float64) float64 {
		m := 0.0
		for i := range xs {
			e := math.Abs(s*xs[i]-gs[i]) / gs[i] * 100
			if e > m {
				m = e
			}
		}
		return m
	}
	switch kp(ctx) {
	case 1:
		e := maxErr(1)
		return num(stem+"What is the largest relative error of the fast extractor against the golden values, in percent of the golden value?", e, 2,
			edakit.Sprintf("For each net compute |fast - golden| / golden. The largest is %s%%.", f(e, 2)), nil)
	case 2:
		return num(stem+"Calibrate a single scale factor k so that k times the fast value best matches the golden values in least squares. What is k?", scale, 4,
			edakit.Sprintf("Minimizing $\\sum (k x_i - g_i)^2$ gives $k = \\sum x_i g_i / \\sum x_i^2$ = %s / %s = %s.", f(sxy, 4), f(sxx, 4), f(scale, 4)),
			map[string]string{"k": f(scale, 4)})
	}
	e := maxErr(edakit.Round(scale, 4))
	return num(stem+edakit.Sprintf("After calibration with the least-squares scale factor k = %s, what is the largest remaining relative error, in percent?", f(scale, 4)), e, 2,
		edakit.Sprintf("Multiply each fast value by %s and compare with golden. The worst net is now off by %s%%. A single factor removes systematic bias but not the scatter, which needs better patterns.", f(scale, 4), f(e, 2)), nil)
}
