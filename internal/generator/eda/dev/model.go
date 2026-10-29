package dev

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var tcadItems = []edakit.Item{
	{KP: 1, Q: "What separates process TCAD from device TCAD?",
		A: "Process TCAD simulates fabrication steps to predict structure and doping; device TCAD solves carrier transport in that structure", D: []string{"Process TCAD is for analog, device TCAD for digital", "Process TCAD runs SPICE; device TCAD runs timing", "They are two names for the same solver"},
		Why: "Implant, diffusion, oxidation and etch simulation produce the geometry and doping that device simulation then biases to get I-V and C-V curves."},
	{KP: 1, Q: "A device TCAD run produces a family of $I_D$-$V_{GS}$ curves. What does an EDA flow mainly use them for?",
		A: "To calibrate and extract compact-model parameters before silicon is available, and to explore device options", D: []string{"To run static timing directly", "As the transistor model inside SPICE for full chips", "To generate the layout"},
		Why: "TCAD is far too slow to call inside circuit simulation. Its curves feed compact-model extraction and design-technology co-optimization."},
	{KP: 1, Q: "Why must a TCAD tool refine the mesh near junctions and the gate-oxide interface?",
		A: "Potential and carrier densities change over nanometres there, and a coarse mesh gives wrong currents", D: []string{"The solver requires a uniform mesh", "The oxide has the highest doping", "Refinement is only for visualization"},
		Why: "Depletion edges and inversion layers have steep gradients. Mesh error there changes threshold and current directly."},
	{KP: 1, Q: "Which discretization do most device simulators use for the semiconductor equations on unstructured meshes?",
		A: "A finite-volume (box) method with Scharfetter-Gummel fluxes along edges", D: []string{"Spectral methods with Fourier modes", "Explicit time stepping on a uniform grid only", "Monte Carlo sampling of the potential"},
		Why: "Box integration conserves current exactly across control volumes, and Scharfetter-Gummel fluxes stay stable for large potential steps between nodes."},
	{KP: 2, Q: "A TCAD solve fails to converge when the gate is set straight to 1 V. What is the standard remedy?",
		A: "Ramp the bias in small steps, using each converged solution as the next initial guess", D: []string{"Coarsen the mesh until it converges", "Remove Poisson's equation", "Lower the temperature to 0 K"},
		Why: "Newton converges only from a nearby guess. Bias continuation keeps every step inside the basin of attraction."},
	{KP: 2, Q: "A fully coupled Newton solve of a 3D device has hundreds of thousands of unknowns per carrier equation. What dominates the run time?",
		A: "Factoring or iteratively solving the sparse Jacobian at each Newton step", D: []string{"Writing the output files", "Evaluating the mobility model", "Meshing the device once"},
		Why: "Each Newton iteration needs a linear solve with the sparse Jacobian of the coupled Poisson and continuity equations. Fill-in makes 3D direct solves expensive."},
	{KP: 2, Q: "A TCAD engineer calibrates a simulator against measured silicon. What is being tuned?",
		A: "Physical model parameters such as mobility, recombination and doping-profile details, until simulated curves match measurements", D: []string{"The measured data, until it matches simulation", "The mesh, until it matches the layout", "The supply voltage of the chip"},
		Why: "Calibration fixes the physics models to a process. A calibrated deck then predicts the effect of changes before they are fabricated."},
	{KP: 2, Q: "How does a TCAD tool obtain small-signal capacitances such as $C_{gg}$?",
		A: "By an AC analysis that linearizes the equations around a DC solution, or by differentiating charge with respect to terminal voltage", D: []string{"From the layout geometry alone", "From the drain current only", "By assuming the Meyer model"},
		Why: "Linearizing the coupled equations at an operating point gives the admittance matrix, whose imaginary part yields the capacitances."},
	{KP: 3, Q: "Why can a single 3D TCAD I-V sweep take hours while a compact model evaluates in microseconds?",
		A: "TCAD solves coupled nonlinear PDEs over a mesh, while a compact model evaluates closed-form expressions", D: []string{"TCAD runs on slower computers", "Compact models skip the physics entirely", "TCAD must also route the device"},
		Why: "This gap is why TCAD feeds models and design exploration, never full-chip simulation."},
	{KP: 3, Q: "A team simulates many statistical TCAD samples with random dopant placements. What are they producing?",
		A: "Variability statistics, such as threshold sigma, to inform statistical compact models", D: []string{"A single nominal device", "A timing library", "A layout"},
		Why: "Atomistic sampling of dopants, edge roughness or grains yields distributions that statistical compact models must reproduce."},
	{KP: 3, Q: "Which open-source project provides a programmable device simulator for semiconductor equations?",
		A: "DEVSIM", D: []string{"OpenSTA", "Yosys", "KLayout"},
		Why: "DEVSIM solves the semiconductor PDEs with a finite-volume method and a scripting interface. The other three are timing, synthesis and layout tools."},
	{KP: 3, Q: "What does design-technology co-optimization (DTCO) use TCAD for?",
		A: "Evaluating device and process options by their effect on cell and circuit metrics, before committing to a process", D: []string{"Only for final signoff of a chip", "To replace place and route", "To measure finished wafers"},
		Why: "DTCO loops device choices through compact models and standard cells to compare power, performance and area, with TCAD at the device end."},
}

var compactItems = []edakit.Item{
	{KP: 1, Q: "A circuit simulator calls a compact model at every Newton iteration for every transistor. What must the model return?",
		A: "Terminal currents and charges, plus their derivatives with respect to terminal voltages", D: []string{"Only the drain current", "The threshold voltage only", "A mesh of the device"},
		Why: "Newton needs the Jacobian: conductances $\\partial I/\\partial V$ and capacitances $\\partial Q/\\partial V$. Missing or wrong derivatives break convergence."},
	{KP: 1, Q: "Why must a compact model be smooth, with continuous derivatives, across all operating regions?",
		A: "Kinks in current or its derivatives make Newton iterations oscillate or fail to converge", D: []string{"Smoothness is only cosmetic", "Discontinuities speed up simulation", "Only the subthreshold region matters for convergence"},
		Why: "A region switch with a jump in $g_m$ or $g_{ds}$ makes the Jacobian inconsistent, so the solver bounces between regions. Analog circuits also need smooth derivatives for distortion analysis."},
	{KP: 1, Q: "What is the main difference between a compact model and TCAD?",
		A: "A compact model is a closed-form approximation evaluated in microseconds; TCAD solves the transport PDEs numerically", D: []string{"A compact model is more accurate than TCAD", "TCAD uses closed-form equations", "Compact models need a mesh"},
		Why: "Circuit simulation evaluates millions of device instances repeatedly, so only closed-form physics-based expressions with fitted parameters are fast enough."},
	{KP: 1, Q: "Why did charge-based capacitance models replace the Meyer model in transient simulation?",
		A: "The Meyer model is not charge conserving, which can create or destroy charge in switched circuits", D: []string{"The Meyer model is too slow", "The Meyer model has no saturation region", "Charge-based models ignore the gate"},
		Why: "Meyer capacitances are not derivatives of one charge function, so integrating them around a cycle need not return to zero. Charge-based models define terminal charges first."},
	{KP: 2, Q: "A compact model passes DC tests but a mixer simulation shows wrong harmonics near $V_{DS} = 0$. Which check would have caught it?",
		A: "The Gummel symmetry test, which checks source-drain symmetry and smooth derivatives at zero drain bias", D: []string{"A check of the threshold at high drain bias", "A layout-versus-schematic check", "A leakage check at the hot corner"},
		Why: "Swapping source and drain must not change physics. The test drives the device symmetrically around $V_{DS} = 0$ and examines higher derivatives of current."},
	{KP: 2, Q: "Which three families describe how compact MOSFET models compute their core quantity?",
		A: "Threshold-voltage based, surface-potential based and inversion-charge based", D: []string{"Linear, quadratic and cubic", "Bulk, SOI and FinFET", "Digital, analog and RF"},
		Why: "BSIM3 and BSIM4 are threshold-voltage based, PSP and BSIM-CMG solve for surface potential, and EKV-style models are built on inversion charge."},
	{KP: 2, Q: "A compact model must scale with geometry. What does that mean in practice?",
		A: "One parameter set must predict devices across the range of drawn widths and lengths, not one device at a time", D: []string{"The model file grows with the number of instances", "The model must be rewritten per layout", "Only minimum-length devices are modeled"},
		Why: "Circuit designers pick any legal W and L. Length and width dependences of threshold, mobility and parasitics are built into the model or handled by binning."},
	{KP: 2, Q: "Why do compact models include temperature as an input?",
		A: "Threshold, mobility, saturation velocity and leakage all change with temperature, and corners span a wide range", D: []string{"Temperature only matters for TCAD", "To compute the package thermal resistance", "Temperature only affects resistors"},
		Why: "Signoff corners span cold to hot. A model that is right only at room temperature mispredicts delay and leakage at the corners that matter."},
	{KP: 3, Q: "A circuit simulator spends most of its time in model evaluation on a large post-layout netlist. What model property matters most?",
		A: "Evaluation cost per call, since the model runs per device per Newton iteration per time step", D: []string{"The number of fitting parameters stored", "The length of the model's documentation", "Whether the model supports binning"},
		Why: "Model evaluation can be a large fraction of SPICE time. Compact-model developers trade expressions for speed while keeping accuracy and smoothness."},
	{KP: 3, Q: "Why is a physics-based compact model preferred over a pure table of measured I-V points inside SPICE?",
		A: "It extrapolates sensibly, stays smooth and conserves charge, while tables need interpolation that can break derivatives", D: []string{"Tables are always less accurate at measured points", "SPICE cannot read tables", "Physics models need no parameters"},
		Why: "Tables are exact at their points but derivatives between points depend on the interpolant, and behaviour outside the table is undefined."},
	{KP: 3, Q: "A foundry ships a model with parameters valid only for $V_{DS} \\le 1.0$ V. A designer simulates at 1.3 V. What should a robust model and simulator do?",
		A: "Warn about operation outside the qualified range, since accuracy is not guaranteed there", D: []string{"Silently clamp the voltage to 1.0 V", "Abort every simulation", "Switch to TCAD automatically"},
		Why: "Models are qualified over a bias, geometry and temperature range. Outside it, equations still evaluate, but nobody has verified them."},
	{KP: 3, Q: "Why are noise models part of a compact model?",
		A: "Thermal and flicker noise depend on bias and geometry, and analog and RF design needs them", D: []string{"Noise is added by the simulator independently of the device", "Digital timing depends only on noise", "Noise models replace the I-V equations"},
		Why: "Channel thermal noise and 1/f noise are bias-dependent device properties. A model that omits them cannot be used for analog signoff."},
}

var bsimItems = []edakit.Item{
	{KP: 1, Q: "Which BSIM model targets FinFETs and other common-gate multi-gate devices?",
		A: "BSIM-CMG", D: []string{"BSIM4", "BSIM-SOI", "BSIM3v3"},
		Why: "BSIM-CMG (common multi-gate) models devices whose gates are tied together around a fin or sheet. BSIM4 is for planar bulk devices."},
	{KP: 1, Q: "Which BSIM model covers planar bulk MOSFETs used through the planar CMOS era?",
		A: "BSIM4", D: []string{"BSIM-CMG", "BSIM-IMG", "BSIM-SOI"},
		Why: "BSIM4 extended BSIM3 for planar bulk devices with gate tunneling, stress and other effects."},
	{KP: 1, Q: "Which BSIM model handles independent front and back gates, as in UTBB FD-SOI with back biasing?",
		A: "BSIM-IMG", D: []string{"BSIM-CMG", "BSIM4", "BSIM3v3"},
		Why: "IMG stands for independent multi-gate. It treats the back gate as a separate terminal, which back-bias design needs."},
	{KP: 1, Q: "Who standardizes compact models such as BSIM for industry use?",
		A: "The Compact Model Coalition, which qualifies models and their releases", D: []string{"Each SPICE vendor separately with no shared standard", "The IEEE Verilog committee", "The foundry for each customer separately"},
		Why: "Standard models let one foundry model file run in many simulators. The coalition maintains the standard model list and versions."},
	{KP: 2, Q: "A foundry model card uses binning. What does that mean?",
		A: "Parameters are fitted separately over ranges of W and L, with size-dependent terms, and the simulator picks the bin for each device", D: []string{"Devices are grouped by threshold flavor", "The model is stored in compressed bins", "Each bin is a different temperature"},
		Why: "Binning, with $P = P_0 + P_L/L + P_W/W + P_P/(WL)$ inside each bin, improves fit across geometry. Bin edges must stay continuous."},
	{KP: 2, Q: "BSIM-CMG is described as surface-potential based. What does that mean for its core computation?",
		A: "It solves for the channel surface potential from the gate bias and electrostatics, then derives charge and current from it", D: []string{"It assumes a fixed threshold and applies the square law", "It interpolates measured tables", "It runs a TCAD solve internally"},
		Why: "Surface-potential models give one continuous description from weak to strong inversion without region switching."},
	{KP: 2, Q: "BSIM-CMG adds a model for the effect of a thin fin body on carriers. Which physics is it?",
		A: "Quantum confinement, which raises the threshold and moves the charge centroid away from the interface", D: []string{"Latch-up", "Electromigration", "Skin effect"},
		Why: "Thin fins confine carriers, which multi-gate models include through threshold and effective oxide corrections."},
	{KP: 2, Q: "Why do BSIM models carry many parameters, often hundreds?",
		A: "They capture many physical effects, each with geometry and temperature dependence, over a wide qualified range", D: []string{"Most parameters are unused padding", "Each parameter corresponds to one transistor instance", "SPICE requires a minimum count"},
		Why: "Short-channel effects, mobility degradation, parasitics, leakage, noise and temperature each need parameters. Only a subset is tuned for a given process."},
	{KP: 3, Q: "BSIM models are released as Verilog-A source as well as built into simulators. Why does that matter to a simulator developer?",
		A: "A Verilog-A compiler path lets the simulator load the exact reference model without hand-porting the equations", D: []string{"Verilog-A runs only on FPGAs", "Verilog-A removes the need for derivatives", "Built-in models are always wrong"},
		Why: "Hand-ported C implementations can diverge from the reference. Compiling the reference Verilog-A, with automatic derivatives, keeps simulators consistent."},
	{KP: 3, Q: "How do the open-source simulators ngspice and its model ecosystem typically handle BSIM-CMG?",
		A: "ngspice loads compiled Verilog-A models such as BSIM-CMG through its OSDI interface, for example compiled with OpenVAF", D: []string{"ngspice cannot simulate FinFET models at all", "ngspice requires TCAD for FinFETs", "BSIM-CMG is only available as a Liberty file"},
		Why: "OSDI is a binary interface for compiled device models. It lets ngspice run Verilog-A compact models at compiled speed."},
	{KP: 3, Q: "A foundry PDK for a FinFET node includes a BSIM-CMG model card. What does a standard-cell characterization flow do with it?",
		A: "Runs SPICE on each cell across slews, loads and corners to build Liberty timing and power tables", D: []string{"Converts it into a layout", "Uses it to compute wire resistance", "Ignores it in favour of measured chips"},
		Why: "Characterization is many SPICE runs. Liberty tables are the abstraction that static timing and power analysis consume."},
	{KP: 3, Q: "Why do BSIM releases include version numbers that PDKs pin exactly?",
		A: "Equations and default behaviour change between versions, so a parameter set is only valid for the version it was extracted with", D: []string{"Versions only change documentation", "Newer versions are always backward identical", "Simulators ignore the version"},
		Why: "Extraction fits parameters to one version's equations. Running them under another version can shift currents and charges."},
}

// unifiedIV evaluates the EKV-style single-piece interpolation F(x) = ln^2(1 + e^x).
func unifiedIV(ctx generator.GeneratorContext) generator.Problem {
	n := pickF(1.2, 1.3, 1.4, 1.5)
	stem := edakit.Sprintf("A compact model uses the single-piece interpolation $I_D = I_{spec} F(x)$ with $F(x) = \\ln^2(1 + e^{x})$ and $x = (V_{GS} - V_t)/(2 n\\, kT/q)$. Take $n = %s$ and $kT/q = 25.85$ mV.", edakit.Num(n, 1))
	switch kpOf(ctx) {
	case 1:
		vov := 0.01 * float64(edakit.Int(-20, 40))
		x := vov / (2 * n * vT300)
		f := math.Pow(math.Log1p(math.Exp(x)), 2)
		ans, rt := sigAnswer(f, 3)
		q := stem + edakit.Sprintf(" What is $F$ at $V_{GS} - V_t = %d$ mV? %s", int(math.Round(vov*1000)), rt)
		return edakit.Numeric(q, ans, edakit.Sprintf("$x = %.4f/(2 \\times %s \\times 0.02585) = %.4f$, $\\ln(1 + e^{x}) = %.5f$, so $F = %s$.", vov, edakit.Num(n, 1), x, math.Log1p(math.Exp(x)), ans), map[string]string{"x": edakit.Num(x, 4)})
	case 2:
		vov := 0.05 * float64(edakit.Int(2, 10))
		x := vov / (2 * n * vT300)
		f := math.Pow(math.Log1p(math.Exp(x)), 2)
		r := f / (x * x)
		ans := edakit.Num(r, 3)
		q := stem + edakit.Sprintf(" In strong inversion $F$ approaches the square law $x^2$. At $V_{GS} - V_t = %d$ mV, what is the ratio $F(x)/x^2$? Round to 3 decimal places.", int(math.Round(vov*1000)))
		return edakit.Numeric(q, ans, edakit.Sprintf("$x = %.4f$, $F = %.5f$ and $x^2 = %.5f$, so the ratio is %s. It tends to 1 as overdrive grows.", x, f, x*x, ans), map[string]string{"x": edakit.Num(x, 4)})
	default:
		vov := -0.01 * float64(edakit.Int(3, 25))
		x := vov / (2 * n * vT300)
		f := math.Pow(math.Log1p(math.Exp(x)), 2)
		r := f / math.Exp(2*x)
		ans := edakit.Num(r, 3)
		q := stem + edakit.Sprintf(" Below threshold $F$ approaches the exponential $e^{2x}$. At $V_{GS} - V_t = %d$ mV, what is the ratio $F(x)/e^{2x}$? Round to 3 decimal places.", int(math.Round(vov*1000)))
		return edakit.Numeric(q, ans, edakit.Sprintf("$x = %.4f$, $F = %.6g$ and $e^{2x} = %.6g$, so the ratio is %s. Since $\\ln(1+e^x) \\approx e^x$ for $x \\ll 0$, the ratio tends to 1 deep in weak inversion, where current falls by a decade per $n (kT/q) \\ln 10$.", x, f, math.Exp(2*x), ans), map[string]string{"x": edakit.Num(x, 4)})
	}
}

// extraction asks for a threshold by maximum-transconductance extrapolation (KP 1), a threshold
// from two saturation points (KP 2) or the transconductance parameter from the same points (KP 3).
func extraction(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		vg := 0.05 * float64(edakit.Int(10, 16))
		id := float64(edakit.Int(10, 60))
		gm := float64(10 * edakit.Int(10, 40))
		vds := 0.05
		vt := vg - id/gm - vds/2
		ans := edakit.Num(vt, 3)
		q := edakit.Sprintf("A linear-region sweep at $V_{DS} = 50$ mV has its maximum transconductance %d µA/V at $V_{GS} = %s$ V, where $I_D = %d$ µA. Using linear extrapolation at maximum $g_m$, $V_t = V_{GS} - I_D/g_m - V_{DS}/2$, what threshold is extracted, in volts? Round to 3 decimal places.", int(gm), edakit.Num(vg, 2), int(id))
		return edakit.Numeric(q, ans, edakit.Sprintf("The tangent at maximum $g_m$ crosses $I_D = 0$ at $V_{GS} - I_D/g_m = %s - %.4f$; subtracting $V_{DS}/2 = 0.025$ V gives $V_t = %s$ V.", edakit.Num(vg, 2), id/gm, ans), map[string]string{"intercept": edakit.Num(vg-id/gm, 4)})
	default:
		vt := 0.01 * float64(edakit.Int(25, 50))
		k := float64(20 * edakit.Int(5, 40)) // uA/V^2, K = k' W/L
		v1 := 0.1 * float64(edakit.Int(7, 9))
		v2 := v1 + 0.1*float64(edakit.Int(2, 4))
		i1 := edakit.Round(0.5*k*(v1-vt)*(v1-vt), 1)
		i2 := edakit.Round(0.5*k*(v2-vt)*(v2-vt), 1)
		s1, s2 := math.Sqrt(i1), math.Sqrt(i2)
		stem := edakit.Sprintf("A saturation sweep of a long-channel NMOS gives $I_D = %s$ µA at $V_{GS} = %s$ V and $I_D = %s$ µA at $V_{GS} = %s$ V. Assume $I_D = \\tfrac12 K (V_{GS} - V_t)^2$, so $\\sqrt{I_D}$ is linear in $V_{GS}$.", edakit.Num(i1, 1), edakit.Num(v1, 1), edakit.Num(i2, 1), edakit.Num(v2, 1))
		if kpOf(ctx) == 2 || (ctx.KPOrdinal == 0 && edakit.Int(0, 1) == 0) {
			vte := (v2*s1 - v1*s2) / (s1 - s2)
			ans := edakit.Num(vte, 3)
			return edakit.Numeric(stem+" What threshold do the two points imply, in volts? Round to 3 decimal places.", ans,
				edakit.Sprintf("$\\sqrt{I_1} = %.4f$ and $\\sqrt{I_2} = %.4f$. The line through them hits zero at $V_t = (V_2\\sqrt{I_1} - V_1\\sqrt{I_2})/(\\sqrt{I_1} - \\sqrt{I_2}) = %s$ V.", s1, s2, ans),
				map[string]string{"sqrtI1": edakit.Num(s1, 4), "sqrtI2": edakit.Num(s2, 4)})
		}
		ke := 2 * math.Pow((s2-s1)/(v2-v1), 2)
		ans := edakit.Num(ke, 1)
		return edakit.Numeric(stem+" What value of $K$ do the two points imply, in µA/V$^2$? Round to 1 decimal place.", ans,
			edakit.Sprintf("The slope of $\\sqrt{I_D}$ is $(%.4f - %.4f)/%s = %.4f$ µA$^{1/2}$/V, which equals $\\sqrt{K/2}$. So $K = 2 \\times %.4f^2 = %s$ µA/V$^2$.", s2, s1, edakit.Num(v2-v1, 1), (s2-s1)/(v2-v1), (s2-s1)/(v2-v1), ans),
			map[string]string{"slope": edakit.Num((s2-s1)/(v2-v1), 4)})
	}
}
