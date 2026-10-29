package csim

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// shooting covers the periodicity condition, a shooting Newton step and settling cost.
func shooting(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		a := stepF(0.3, 0.05, 13)
		b := stepF(0.1, 0.05, 18)
		x := b / (1 - a)
		return num(edakit.Sprintf("Integrating a linear circuit driven by a periodic source over one period T maps its capacitor voltage as $x(T) = %s\\, x(0) + %s$ V. Periodic steady state requires $x(T) = x(0)$. What is that initial voltage, in volts?", f(a, 2), f(b, 2)), x, 4,
			edakit.Sprintf("Solve $x = %s x + %s$: $x = %s/(1 - %s) = %s$ V. Starting there, one period returns to the same state.", f(a, 2), f(b, 2), f(b, 2), f(a, 2), f(x, 4)), nil)
	case 2:
		x0 := stepF(0.2, 0.1, 15)
		phi := x0 + stepF(-0.5, 0.05, 20)
		if phi == x0 {
			phi += 0.15
		}
		dphi := stepF(0.1, 0.05, 15)
		x1 := x0 - (phi-x0)/(dphi-1)
		return num(edakit.Sprintf("Shooting Newton guesses the initial state $x_0 = %s$ V, integrates one period and gets $\\phi(x_0) = %s$ V; the transient sensitivity gives $\\phi'(x_0) = %s$. Newton on $r(x) = \\phi(x) - x = 0$ uses $r' = \\phi' - 1$. What is the next guess $x_1$, in volts?", f(x0, 2), f(phi, 2), f(dphi, 2)), x1, 4,
			edakit.Sprintf("$r = %s$, $r' = %s$, so $x_1 = %s - (%s)/(%s) = %s$ V. The sensitivity $\\phi'$ comes from the Jacobians the transient already factored at every step.", f(phi-x0, 2), f(dphi-1, 2), f(x0, 2), f(phi-x0, 2), f(dphi-1, 2), f(x1, 4)), nil)
	}
	a := pickF(0.9, 0.95, 0.98, 0.99, 0.995, 0.999)
	e0 := pickF(0.5, 1, 2)
	tol := pickF(1e-3, 1e-6)
	n := int(math.Ceil(math.Log(tol/e0) / math.Log(a)))
	return edakit.Numeric(edakit.Sprintf("A lightly damped resonant circuit loses only a fraction of its initial-condition error per period: the error is multiplied by %s each period. Starting %s V away from steady state, how many full periods must a plain transient run before the error falls to %s V or below?", f(a, 3), f(e0, 2), sci(tol)),
		edakit.IntStr(n),
		edakit.Sprintf("Need $%s \\cdot %s^N \\le %s$, so $N \\ge \\ln(%s)/\\ln(%s) = %s$, giving %d periods. Shooting finds the periodic state with a few Newton iterations, each one period long.", f(e0, 2), f(a, 3), sciM(tol), sciM(tol/e0), f(a, 3), f(math.Log(tol/e0)/math.Log(a), 2), n), nil)
}

// harmonicBalance covers unknown counts and harmonic content of polynomial nonlinearities.
func harmonicBalance(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 5, 200)
		k := edakit.Pick([]int{3, 5, 7, 10, 15, 31})
		u := n * (2*k + 1)
		return edakit.Numeric(edakit.Sprintf("Harmonic balance represents each of a circuit's %d MNA unknowns by a DC term and the cosine and sine amplitudes of harmonics 1 to %d. How many real unknowns does the harmonic-balance system have?", n, k),
			edakit.IntStr(u),
			edakit.Sprintf("Each unknown carries $2K + 1 = %d$ real coefficients, so $%d \\times %d = %d$. The Jacobian couples harmonics through every nonlinear device, which is why large HB uses matrix-free Krylov solvers.", 2*k+1, n, 2*k+1, u), nil)
	case 2:
		a1 := stepF(0.5, 0.5, 9)
		a2 := stepF(0.5, 0.5, 9)
		amp := stepF(0.1, 0.1, 9)
		which := edakit.Pick([]string{"DC", "fundamental", "second-harmonic"})
		var v float64
		switch which {
		case "DC":
			v = a2 * amp * amp / 2
		case "fundamental":
			v = a1 * amp
		default:
			v = a2 * amp * amp / 2
		}
		return num(edakit.Sprintf("A nonlinear conductance draws $i = a_1 v + a_2 v^2$ with $a_1 = %s$ mS and $a_2 = %s$ mA/V$^2$. Harmonic balance evaluates it on $v = %s\\cos\\omega t$ V. What is the %s amplitude of the current, in mA?", f(a1, 2), f(a2, 2), f(amp, 2), which), v, 4,
			edakit.Sprintf("$\\cos^2 = \\tfrac12 + \\tfrac12\\cos 2\\omega t$, so $i = a_2 A^2/2 + a_1 A \\cos\\omega t + (a_2 A^2/2)\\cos 2\\omega t$. The %s term is %s mA.", which, f(v, 4)), nil)
	}
	a1 := stepF(0.5, 0.5, 9)
	a3 := -stepF(0.1, 0.1, 9)
	amp := stepF(0.2, 0.1, 8)
	which := edakit.Pick([]string{"fundamental", "third-harmonic"})
	v := a3 * amp * amp * amp / 4
	if which == "fundamental" {
		v = a1*amp + 3*a3*amp*amp*amp/4
	}
	return num(edakit.Sprintf("A device current is $i = a_1 v + a_3 v^3$ with $a_1 = %s$ mS and $a_3 = %s$ mA/V$^3$, driven by $v = %s\\cos\\omega t$ V. What is the %s amplitude of the current, in mA? (Use $\\cos^3 = \\tfrac34\\cos + \\tfrac14\\cos 3$.)", f(a1, 2), f(a3, 2), f(amp, 2), which), v, 4,
		edakit.Sprintf("The fundamental is $a_1 A + \\tfrac34 a_3 A^3 = %s$ mA (compressed by the negative cubic term) and the third harmonic is $\\tfrac14 a_3 A^3 = %s$ mA.", f(a1*amp+3*a3*amp*amp*amp/4, 4), f(a3*amp*amp*amp/4, 4)), nil)
}

var rfMethodItems = []edakit.Item{
	{KP: 1, Q: "A mixer is driven by a large local oscillator and a tiny RF input. Why is ordinary AC analysis at the DC operating point wrong for the RF path?",
		A: "The small signal sees a circuit whose operating point varies periodically with the LO, so it must be linearized about the periodic steady state", D: []string{"AC analysis cannot handle frequencies above 1 GHz", "The mixer has no DC operating point", "AC analysis ignores capacitors"},
		Why: "Periodic small-signal analysis (PAC) linearizes about the large-signal periodic solution, producing a linear time-varying system that translates frequencies."},
	{KP: 1, Q: "Periodic noise analysis of a mixer reports noise at the IF output that originated near the image frequency. What phenomenon is this?",
		A: "Noise folding: the time-varying operating point converts noise from sidebands of the LO harmonics to the output frequency", D: []string{"A numerical artifact of harmonic balance", "Thermal noise that is independent of the LO", "Aliasing from the transient timestep"},
		Why: "Cyclostationary analysis tracks how a periodically modulated circuit mixes noise from every sideband onto the output, which stationary noise analysis cannot."},
	{KP: 1, Q: "Which quantity is computed by periodic transfer-function analysis (PXF) after a PSS solve?",
		A: "Transfer functions from many inputs, at many sidebands, to one output, using an adjoint solve", D: []string{"The large-signal waveform of the LO", "The DC operating point at every harmonic", "The Monte Carlo spread of conversion gain"},
		Why: "Like adjoint sensitivity, PXF solves one adjoint problem for a chosen output and reads off the transfer from every source and sideband at once."},
	{KP: 1, Q: "Why must periodic small-signal analyses run after a converged periodic steady-state solve?",
		A: "They linearize around that periodic solution, so its accuracy bounds theirs", D: []string{"They reuse its random seed", "They need its DC operating point only", "They do not; they are independent"},
		Why: "PAC, PXF and PNoise all use the time-varying Jacobian of the periodic solution. A poorly converged PSS gives wrong conversion gains and noise."},
	{KP: 2, Q: "A two-tone intermodulation test uses tones at 1.000 GHz and 1.001 GHz. Why is a plain transient a poor way to get the steady state?",
		A: "The common period is 1 us while the carrier period is 1 ns, so thousands of fast cycles must be simulated just to see one beat", D: []string{"Transient cannot simulate two sources", "The intermodulation products are below machine precision in any method", "Transient analysis rejects sinusoidal sources"},
		Why: "Closely spaced tones make the combined waveform quasi-periodic with a long period. Multitone harmonic balance represents both tones directly in the frequency domain."},
	{KP: 2, Q: "A power amplifier is driven by a modulated carrier whose envelope changes over microseconds. Which method suits it?",
		A: "Envelope following: solve a periodic problem at the carrier for each slow time step of the envelope", D: []string{"A single harmonic-balance solve with one tone", "DC sweep of the envelope", "AC analysis at the carrier frequency"},
		Why: "Envelope methods separate the fast carrier from the slow modulation, stepping a periodic-steady-state problem along the slow time axis."},
	{KP: 2, Q: "Multitone harmonic balance with two tones keeps all mixing products up to order 7. What grows quickly as the order rises?",
		A: "The number of frequencies, and with it the size of the nonlinear system", D: []string{"The number of circuit nodes", "The timestep", "The DC gain"},
		Why: "Each retained mixing product $m f_1 + n f_2$ is another frequency per unknown. Truncation by total order keeps the count manageable but still grows roughly quadratically."},
	{KP: 2, Q: "Harmonic balance evaluates nonlinear devices where?",
		A: "In the time domain at sample points of the period, moving between time and frequency with FFTs", D: []string{"Only in the frequency domain, using closed-form device spectra", "At the DC operating point only", "At each Newton iteration of a transient run"},
		Why: "Device models are written in time. HB transforms the current frequency-domain guess to time samples, evaluates currents and charges, and transforms back."},
	{KP: 3, Q: "A switched-capacitor filter with hard clock edges is to be solved for its periodic steady state. Which method is the more natural choice?",
		A: "Shooting, because time-domain integration handles sharp switching without many harmonics", D: []string{"Harmonic balance with five harmonics", "AC analysis", "DC sweep"},
		Why: "Strongly nonlinear, switching waveforms need very many harmonics in HB, while shooting integrates the edges directly with adaptive steps."},
	{KP: 3, Q: "A low-noise amplifier with S-parameter models of its package and board is to be analyzed with mild nonlinearity. Which method is the more natural choice?",
		A: "Harmonic balance, because frequency-domain models stamp directly at each harmonic", D: []string{"Shooting, because S-parameters are time-domain models", "Pseudo-transient analysis", "Fast SPICE"},
		Why: "Frequency-domain data such as S-parameters are awkward in transient (they need convolution or fitting) but natural in HB, which works harmonic by harmonic."},
	{KP: 3, Q: "Large harmonic-balance problems are solved with GMRES and a preconditioner rather than a direct factorization. Why?",
		A: "The Jacobian is large and dense in the harmonic dimension, but its product with a vector is cheap through FFTs", D: []string{"Direct solvers cannot handle complex numbers", "GMRES needs no preconditioner", "The Jacobian is diagonal"},
		Why: "Matrix-implicit Krylov methods apply the HB Jacobian through device evaluations and FFTs without forming it, and use block-diagonal (linear-part) preconditioners."},
	{KP: 3, Q: "An oscillator has no external periodic drive. What extra unknown does a PSS solver need?",
		A: "The period itself, fixed by a phase condition that pins one point of the waveform", D: []string{"None; the period is the timestep", "The DC gain", "A second voltage source"},
		Why: "An autonomous oscillator's frequency is unknown. Shooting or HB adds the period as an unknown and an extra equation removing the arbitrary time shift."},
}
