package csim

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// rcChain draws an RC chain: resistance r[i] in ohms from node i-1 to node i, and c[i] in fF
// from node i to ground. Node 0 is the driven root.
func rcChain(n int) (r, c []float64, desc string) {
	var parts []string
	for i := 0; i < n; i++ {
		r = append(r, 50*float64(edakit.Int(1, 8)))
		c = append(c, 5*float64(edakit.Int(1, 8)))
		parts = append(parts, edakit.Sprintf("$R_%d = %s\\ \\Omega$ into node %d with $C_%d = %s$ fF", i+1, f(r[i], 0), i+1, i+1, f(c[i], 0)))
	}
	return r, c, strings.Join(parts, ", ")
}

// chainM1 returns the first moment (negative Elmore delay) at each node, in ps.
func chainM1(r, c []float64) []float64 {
	n := len(r)
	m := make([]float64, n)
	for i := 0; i < n; i++ {
		s := 0.0
		for k := 0; k < n; k++ {
			shared := 0.0
			for j := 0; j <= i && j <= k; j++ {
				shared += r[j]
			}
			s += shared * c[k]
		}
		m[i] = -s / 1000
	}
	return m
}

// moments asks for the first or second moment of an RC chain.
func moments(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	n := edakit.Scaled(ctx, 2, 4)
	if k == 3 {
		n = edakit.Int(2, 3)
	}
	r, c, desc := rcChain(n)
	m1 := chainM1(r, c)
	stem := "An ideal voltage source drives the root of an RC chain: " + desc + ", each capacitor to ground. The transfer function to each node expands as $H(s) = m_0 + m_1 s + m_2 s^2 + \\cdots$ with $m_0 = 1$."
	switch k {
	case 1:
		return num(stem+edakit.Sprintf(" What is the first moment $m_1$ at the far end, node %d, in picoseconds?", n), m1[n-1], 3,
			edakit.Sprintf("$m_1 = -\\sum_k R_{%dk} C_k$, minus the Elmore delay: %s ps (ohm times fF is fs).", n, f(m1[n-1], 3)), nil)
	case 2:
		i := edakit.Int(1, n-1)
		shared := "$R_1$"
		if i > 1 {
			shared = edakit.Sprintf("$R_1 + \\cdots + R_%d$", i)
		}
		return num(stem+edakit.Sprintf(" What is the first moment $m_1$ at node %d, in picoseconds?", i), m1[i-1], 3,
			edakit.Sprintf("Each capacitor contributes its value times the resistance its root path shares with node %d's path. Capacitors at or beyond node %d share only %s. The sum gives $m_1 = %s$ ps.", i, i, shared, f(m1[i-1], 3)), nil)
	}
	// m2_i = -sum_k R_ik C_k m1_k, in ps^2 (ohm*fF = fs = 1e-3 ps).
	m2 := 0.0
	for kk := 0; kk < n; kk++ {
		shared := 0.0
		for j := 0; j <= kk; j++ {
			shared += r[j]
		}
		m2 += -shared * c[kk] * m1[kk] / 1000
	}
	var m1s []string
	for i := range m1 {
		m1s = append(m1s, edakit.Sprintf("$m_1^{(%d)} = %s$", i+1, f(m1[i], 3)))
	}
	return num(stem+" The first moments are "+strings.Join(m1s, ", ")+edakit.Sprintf(" ps. Using $m_2^{(i)} = -\\sum_k R_{ik} C_k m_1^{(k)}$, what is $m_2$ at node %d, in ps$^2$?", n), m2, 3,
		edakit.Sprintf("Each capacitor's charging current is weighted by the previous moment at its own node: the sum gives $m_2 = %s$ ps$^2$, positive as for any RC tree.", f(m2, 3)), nil)
}

// awe fits poles to moments.
func awe(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		m0 := pickF(1, 1, 0.5, 0.8, 0.9)
		m1 := -stepF(0.2, 0.1, 30)
		p := m0 / m1
		return num(edakit.Sprintf("A one-pole AWE model $H(s) \\approx k/(1 - s/p)$ is fitted to the moments $m_0 = %s$ and $m_1 = %s$ ns. Matching the expansion $k(1 + s/p + \\cdots)$ term by term, what is the pole $p$ in 1/ns?", f(m0, 2), f(m1, 2)), p, 4,
			edakit.Sprintf("$k = m_0$ and $k/p = m_1$, so $p = m_0/m_1 = %s$ 1/ns, a time constant of %s ns.", f(p, 4), f(-1/p, 4)), nil)
	case 2:
		m1 := -stepF(0.2, 0.1, 30)
		t50 := math.Ln2 * -m1
		return num(edakit.Sprintf("A node's step response is approximated by the dominant pole fitted from $m_0 = 1$ and $m_1 = %s$ ns. What 50%% delay does the one-pole model predict, in ns?", f(m1, 2)), t50, 4,
			edakit.Sprintf("The time constant is $-m_1/m_0 = %s$ ns, and a single-pole step reaches 50%% at $\\tau \\ln 2 = %s$ ns.", f(-m1, 2), f(t50, 4)), nil)
	}
	t1 := pickF(1, 2, 3, 4)
	t2 := pickF(0.1, 0.2, 0.5)
	w1 := pickF(0.6, 0.7, 0.8, 0.9)
	w2 := 1 - w1
	m1 := -(w1*t1 + w2*t2)
	m2 := w1*t1*t1 + w2*t2*t2
	m3 := -(w1*t1*t1*t1 + w2*t2*t2*t2)
	det := m2 - m1*m1
	b1 := (m1*m2 - m3) / det
	return num(edakit.Sprintf("A two-pole Pade (AWE) model uses the denominator $1 + b_1 s + b_2 s^2$, found from $\\begin{pmatrix} m_0 & m_1 \\\\ m_1 & m_2 \\end{pmatrix} \\begin{pmatrix} b_2 \\\\ b_1 \\end{pmatrix} = -\\begin{pmatrix} m_2 \\\\ m_3 \\end{pmatrix}$. The moments are $m_0 = 1$, $m_1 = %s$, $m_2 = %s$, $m_3 = %s$ (in powers of ns). What is $b_1$ in ns?", f(m1, 6), f(m2, 6), f(m3, 6)), b1, 4,
		edakit.Sprintf("By Cramer's rule $b_1 = (m_1 m_2 - m_0 m_3)/(m_0 m_2 - m_1^2) = %s/%s = %s$ ns. This equals the sum of the two time constants, %s and %s ns. With more moments AWE can produce unstable poles, which PRIMA avoids.", f(m1*m2-m3, 6), f(det, 6), f(b1, 4), f(t1, 2), f(t2, 2)), nil)
}

// prima covers moment matching counts, congruence projection and the reduced pole.
func prima(ctx generator.GeneratorContext) generator.Problem {
	vs := [][2]float64{{1, 0}, {0, 1}, {0.6, 0.8}, {0.8, 0.6}, {-0.6, 0.8}}
	gd := []float64{pickF(1, 2, 3, 4), pickF(1, 2, 3, 4)}
	go_ := -pickF(0.5, 1)
	switch kp(ctx) {
	case 1:
		p := edakit.Int(1, 8)
		m := edakit.Int(2, 10)
		q := p * m
		return edakit.Numeric(edakit.Sprintf("PRIMA reduces a %d-port RC interconnect model with a block Arnoldi basis of $q = %d$ columns. How many block moments of the port admittance does the reduced model match?", p, q),
			edakit.IntStr(m),
			edakit.Sprintf("Each block Krylov step adds one column per port, and q columns match $q/p = %d/%d = %d$ block moments.", q, p, m), nil)
	case 2:
		v := vs[edakit.Int(0, len(vs)-1)]
		gg := gd[0]*v[0]*v[0] + 2*go_*v[0]*v[1] + gd[1]*v[1]*v[1]
		return num(edakit.Sprintf("A two-node conductance matrix is $G = \\begin{pmatrix} %s & %s \\\\ %s & %s \\end{pmatrix}$ mS. PRIMA projects it with an orthonormal basis, here the single vector $v = (%s, %s)$, as $\\hat{G} = v^T G v$. What is $\\hat{G}$ in mS?", f(gd[0], 2), f(go_, 2), f(go_, 2), f(gd[1], 2), f(v[0], 2), f(v[1], 2)), gg, 4,
			edakit.Sprintf("$v^T G v = %s \\cdot (%s)^2 + 2(%s)(%s)(%s) + %s \\cdot (%s)^2 = %s$ mS. Congruence keeps $\\hat{G}$ positive definite when $G$ is, which preserves passivity.", f(gd[0], 2), f(v[0], 2), f(go_, 2), f(v[0], 2), f(v[1], 2), f(gd[1], 2), f(v[1], 2), f(gg, 4)), nil)
	}
	v := vs[edakit.Int(2, len(vs)-1)]
	cd := []float64{pickF(0.5, 1, 2), pickF(0.5, 1, 2)}
	gg := gd[0]*v[0]*v[0] + 2*go_*v[0]*v[1] + gd[1]*v[1]*v[1]
	cc := cd[0]*v[0]*v[0] + cd[1]*v[1]*v[1]
	p := -gg / cc
	return num(edakit.Sprintf("A two-node RC model has $G = \\begin{pmatrix} %s & %s \\\\ %s & %s \\end{pmatrix}$ mS and $C = \\mathrm{diag}(%s, %s)$ pF. It is reduced to one state with $v = (%s, %s)$: $\\hat{G} = v^T G v$, $\\hat{C} = v^T C v$. What is the pole $-\\hat{G}/\\hat{C}$ of the reduced model, in 1/ns?", f(gd[0], 2), f(go_, 2), f(go_, 2), f(gd[1], 2), f(cd[0], 2), f(cd[1], 2), f(v[0], 2), f(v[1], 2)), p, 4,
		edakit.Sprintf("$\\hat{G} = %s$ mS and $\\hat{C} = %s$ pF, so the pole is $-%s/%s = %s$ 1/ns (mS per pF is 1/ns). Both reduced matrices stay positive, so the pole stays in the left half plane.", f(gg, 4), f(cc, 4), f(gg, 4), f(cc, 4), f(p, 4)), nil)
}

var macromodelItems = []edakit.Item{
	{KP: 1, Q: "A package is characterized by measured S-parameters from 10 MHz to 20 GHz. What does vector fitting produce for the simulator?",
		A: "A rational function in s, a sum of pole-residue terms, that can be realized as an equivalent circuit or state-space model", D: []string{"A lookup table indexed by time", "A single resistor and capacitor", "A Monte Carlo distribution of the S-parameters"},
		Why: "Vector fitting iteratively relocates poles to fit frequency data with a rational model, which transient simulators can integrate directly."},
	{KP: 1, Q: "Why can't a transient simulator use tabulated S-parameter data directly without a fitted model?",
		A: "Frequency samples define a time response only through a convolution, which is expensive and sensitive to band limits", D: []string{"S-parameters are not linear", "Transient analysis only accepts resistors", "The data have no phase"},
		Why: "Without a rational fit, the simulator must convolve with the impulse response at every step, and truncation of the measured band causes causality artifacts."},
	{KP: 1, Q: "A rational fit reproduces the data well but has a pole with positive real part. What happens in transient?",
		A: "The model has a growing exponential mode and the simulation can diverge", D: []string{"Nothing; poles affect only AC", "The fit becomes passive automatically", "The simulator ignores that pole"},
		Why: "A right-half-plane pole is an unstable mode. Fitting algorithms flip or constrain poles to the left half plane to keep the model stable."},
	{KP: 1, Q: "Why is a reduced or fitted model typically stamped as a small dense block in the MNA matrix?",
		A: "Its state variables couple all ports to each other, unlike the sparse local stamps of physical elements", D: []string{"Dense blocks are required by Newton", "SPICE cannot store sparse blocks", "Reduction always produces diagonal matrices"},
		Why: "Projection mixes the original states, so the reduced G and C are dense. Keeping the order small keeps that dense block cheap."},
	{KP: 2, Q: "A fitted model of a passive interconnect is stable but generates energy at some frequencies. Why is that a problem when it is connected to other circuits?",
		A: "A non-passive model can make an otherwise stable interconnection oscillate or blow up", D: []string{"It only changes the DC operating point", "It makes the matrix symmetric", "It has no effect if each part is stable alone"},
		Why: "Stability of each block does not guarantee stability of the connection. Passivity does: a network of passive parts cannot create energy."},
	{KP: 2, Q: "What property of PRIMA guarantees a passive reduced model for RLC interconnect?",
		A: "A congruence projection of the MNA matrices in their passive form, which preserves positive (semi)definiteness", D: []string{"Matching more moments than AWE", "Using a Pade approximation", "Fitting the poles by least squares"},
		Why: "PRIMA reduces $V^T G V$ and $V^T C V$ directly. With the MNA matrices arranged so their symmetric parts are positive semidefinite, the projections inherit that and stay passive."},
	{KP: 2, Q: "A fitted model's impulse response is nonzero before t = 0. What property is violated?",
		A: "Causality", D: []string{"Passivity", "Reciprocity", "Linearity"},
		Why: "A causal system cannot respond before its input. Band-limited or inconsistent data can produce acausal fits, which must be checked and corrected."},
	{KP: 2, Q: "After vector fitting, an engineer checks passivity by testing the eigenvalues of the real part of the admittance at many frequencies. What fix follows if some are negative?",
		A: "Passivity enforcement: perturb the residues minimally until the violations disappear", D: []string{"Increase the simulator timestep", "Delete the offending ports", "Switch to backward Euler"},
		Why: "Passivity enforcement adjusts the model, usually the residues, with a small change that removes negative-resistance regions while keeping the fit accurate."},
	{KP: 3, Q: "A system-level testbench replaces a transistor-level op-amp with a Verilog-A model containing gain, one pole and output limits. What does this trade?",
		A: "Simulation speed for fidelity outside the behaviors the model was written to capture", D: []string{"Nothing: behavioral models are exact", "Accuracy for robustness only in DC", "It slows simulation because Verilog-A is interpreted"},
		Why: "A behavioral macromodel evaluates a few equations instead of dozens of transistors, but misses effects it does not encode, such as supply sensitivity or slewing details."},
	{KP: 3, Q: "A Verilog-A model computes an output current with an if-statement that jumps abruptly at a threshold. What numerical problem can that cause?",
		A: "A discontinuous current and derivative that make Newton oscillate and force tiny timesteps", D: []string{"A singular capacitance matrix", "Loss of passivity in AC", "Nothing, because Verilog-A smooths it"},
		Why: "Simulators need continuous currents and charges with continuous derivatives. Behavioral models should use smooth transitions such as tanh or limited slopes."},
	{KP: 3, Q: "A behavioral model of a voltage regulator is validated only at 27 C and nominal load. What must the model's documentation state?",
		A: "Its validity range, so users do not trust it at temperatures and loads it was never fitted to", D: []string{"Nothing, since behavioral models generalize", "Its LU fill-in", "Its compiled binary size"},
		Why: "A macromodel is a fit. Outside its fitted conditions its errors are unknown, so the conditions are part of the model's contract."},
	{KP: 3, Q: "A large digital block is replaced by a table of its supply current versus time for a power-grid simulation. What kind of model is this?",
		A: "A behavioral current-source macromodel that ignores the block's internal response to supply droop", D: []string{"A PRIMA-reduced model of the transistors", "A full transistor model", "A passive RC model"},
		Why: "Current signatures make grid simulation tractable, but because they are fixed, they do not change when the droop slows the logic."},
}
