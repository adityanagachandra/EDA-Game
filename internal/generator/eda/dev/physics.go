package dev

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// fermiOccupancy asks for a Fermi-Dirac occupation probability, with the classic wrong
// computations (the hole probability, a halved or doubled exponent) as distractors.
func fermiOccupancy(ctx generator.GeneratorContext) generator.Problem {
	T := edakit.Pick([]int{250, 300, 350, 400})
	kT := kB * float64(T) * 1000 // meV
	dE := 5 * edakit.Int(1, 14)
	if edakit.Int(0, 1) == 0 {
		dE = -dE
	}
	x := float64(dE) / kT
	f := 1 / (1 + math.Exp(x))
	where := "above"
	if dE < 0 {
		where = "below"
	}
	q := edakit.Sprintf("At $T = %d$ K ($kT = %.2f$ meV), an allowed electron state sits %d meV %s the Fermi level. What is the probability that the state is occupied by an electron?",
		T, kT, int(math.Abs(float64(dE))), where)
	return edakit.Choice(q, fixed(f, 3),
		edakit.Sprintf("Fermi-Dirac: $f = 1/(1 + e^{(E - E_F)/kT})$ with $(E-E_F)/kT = %d/%.2f = %.3f$, so $f = %.3f$. The value $1 - f$ is the probability that the state is empty (holds a hole).", dE, kT, x, f),
		fixed(1-f, 3), fixed(1/(1+math.Exp(x/2)), 3), fixed(1/(1+math.Exp(2*x)), 3))
}

var bandsItems = []edakit.Item{
	{KP: 1, Q: "Which band picture separates a semiconductor such as silicon from a metal?",
		A: "A filled valence band and an empty conduction band separated by a gap of about 1 eV", D: []string{"A partly filled band crossing the Fermi level", "A gap of about 9 eV that no thermal carrier can cross", "Overlapping valence and conduction bands with no gap"},
		Why: "Silicon has a 1.12 eV gap at room temperature: small enough that a few carriers are thermally excited, large enough that conduction is controlled by doping. A metal has a partly filled band."},
	{KP: 1, Q: "An electron is excited from the valence band to the conduction band. What carriers result?",
		A: "One conduction electron and one hole in the valence band", D: []string{"Two holes in the valence band", "One conduction electron only, since holes are not real particles", "One ionized donor and one electron"},
		Why: "Band-to-band generation creates an electron-hole pair. The hole is the missing electron state, and it moves like a positive carrier."},
	{KP: 1, Q: "Why does a TCAD tool describe a conduction electron with an effective mass instead of the free-electron mass?",
		A: "The periodic lattice potential changes how the electron accelerates, which the band curvature captures", D: []string{"Electrons gain mass from the doping atoms they pass", "The effective mass corrects for relativity at high fields", "It is a fitting constant with no physical meaning"},
		Why: "Near a band edge, $E \\approx \\hbar^2 k^2/(2m^*)$. The curvature of the band sets $m^*$, which then enters mobility, density of states and tunneling."},
	{KP: 1, Q: "Silicon's band gap is indirect. What practical consequence does that have?",
		A: "Radiative recombination is weak, so silicon is a poor light emitter", D: []string{"Silicon cannot conduct current at room temperature", "Silicon has no holes", "Silicon's band gap grows with temperature"},
		Why: "In an indirect-gap material the conduction minimum and valence maximum sit at different crystal momenta, so recombination needs a phonon and is slow and non-radiative in practice."},
	{KP: 3, Q: "A p-n junction sits in thermal equilibrium with no applied bias. What does the Fermi level do across the junction?",
		A: "It is flat: the same energy everywhere", D: []string{"It steps up by the built-in potential at the junction", "It follows the conduction band edge", "It is undefined inside the depletion region"},
		Why: "In equilibrium there is no net current, which requires a constant Fermi level. The bands bend instead, and the bending is the built-in potential."},
	{KP: 3, Q: "Under forward bias, a device simulator reports two separate Fermi levels in the junction region. What are they?",
		A: "Quasi-Fermi levels for electrons and holes, which split by the applied voltage", D: []string{"The Fermi levels of the two contacts, which never meet", "The intrinsic level and the Fermi level", "A numerical artifact of the mesh"},
		Why: "Out of equilibrium, each carrier population has its own quasi-Fermi level. Their separation is $qV$ across the junction, and their gradients drive the currents."},
	{KP: 3, Q: "Where does the Fermi level sit in undoped (intrinsic) silicon at room temperature?",
		A: "Close to midgap, offset slightly by the unequal densities of states", D: []string{"At the conduction band edge", "At the valence band edge", "Exactly at the middle with no offset under any condition"},
		Why: "With $n = p$, $E_F = (E_c + E_v)/2 + (kT/2)\\ln(N_v/N_c)$. The correction is a few meV in silicon."},
	{KP: 3, Q: "A source region is doped to $10^{20}$ cm$^{-3}$, above the conduction-band density of states. What changes in the model?",
		A: "The Fermi level enters the band, so Boltzmann statistics fail and Fermi-Dirac statistics are needed", D: []string{"Nothing: Boltzmann statistics hold at any doping", "The material becomes an insulator", "The band gap doubles"},
		Why: "Degenerate doping puts $E_F$ inside the conduction band. The Boltzmann approximation $f \\approx e^{-(E-E_F)/kT}$ needs $E - E_F$ of several $kT$, which no longer holds. Heavy doping also narrows the gap."},
}

// intrinsic draws ni at 300 K (KP 1), ni at another temperature (KP 2) or the ratio of ni
// across a temperature step (KP 3).
func intrinsic(ctx generator.GeneratorContext) generator.Problem {
	nc := pickF(2.5, 2.8, 3.0, 3.2)
	nv := pickF(1.0, 1.04, 1.2, 1.5, 1.8)
	eg := pickF(1.08, 1.10, 1.12, 1.14)
	switch kpOf(ctx) {
	case 1:
		ni := math.Sqrt(nc*1e19*nv*1e19) * math.Exp(-eg/(2*vT300))
		m, e := mantissa(ni)
		ans := edakit.Num(m, 2)
		q := edakit.Sprintf("A semiconductor at 300 K has $N_c = %s$ cm$^{-3}$, $N_v = %s$ cm$^{-3}$ and band gap $E_g = %.2f$ eV. Take $kT = 0.02585$ eV. Compute the intrinsic carrier concentration $n_i = \\sqrt{N_c N_v}\\, e^{-E_g/(2kT)}$ and give it in units of $10^{%d}$ cm$^{-3}$. Round to 2 decimal places.",
			sci(nc, 19), sci(nv, 19), eg, e)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$\\sqrt{N_c N_v} = %s$ cm$^{-3}$ and $e^{-E_g/2kT} = e^{-%.2f} = %s$, so $n_i = %s$ cm$^{-3}$.",
				sciOf(math.Sqrt(nc*nv)*1e19), eg/(2*vT300), sciOf(math.Exp(-eg/(2*vT300))), sciOf(ni)),
			map[string]string{"ni": sciOf(ni)})
	case 2:
		T := float64(edakit.Pick([]int{250, 350, 400, 450}))
		kT := kB * T
		s := math.Pow(T/300, 1.5)
		ni := math.Sqrt(nc*nv*1e38) * s * math.Exp(-eg/(2*kT))
		m, e := mantissa(ni)
		ans := edakit.Num(m, 2)
		q := edakit.Sprintf("At 300 K a semiconductor has $N_c = %s$ cm$^{-3}$ and $N_v = %s$ cm$^{-3}$. Both scale as $T^{3/2}$. Treat the band gap as fixed at $E_g = %.2f$ eV and take $k = 8.617 \\times 10^{-5}$ eV/K. What is $n_i$ at $T = %d$ K? Give it in units of $10^{%d}$ cm$^{-3}$. Round to 2 decimal places.",
			sci(nc, 19), sci(nv, 19), eg, int(T), e)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$kT = %.5f$ eV. $\\sqrt{N_c N_v}$ grows by $(T/300)^{3/2} = %.4f$, and $e^{-E_g/2kT} = e^{-%.3f}$. The product is $n_i = %s$ cm$^{-3}$.", kT, s, eg/(2*kT), sciOf(ni)),
			map[string]string{"ni": sciOf(ni)})
	default:
		t1 := float64(edakit.Pick([]int{300, 320, 350}))
		t2 := t1 + float64(edakit.Pick([]int{10, 20, 25, 40, 50}))
		r := math.Pow(t2/t1, 1.5) * math.Exp(eg/(2*kB)*(1/t1-1/t2))
		ans, rt := sigAnswer(r, 3)
		q := edakit.Sprintf("A leakage model needs to know how fast $n_i$ grows with temperature. With $N_c N_v \\propto T^3$, a fixed band gap $E_g = %.2f$ eV and $k = 8.617 \\times 10^{-5}$ eV/K, by what factor does $n_i$ increase from %d K to %d K? %s",
			eg, int(t1), int(t2), rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$n_i(T_2)/n_i(T_1) = (T_2/T_1)^{3/2} \\exp\\left[\\frac{E_g}{2k}\\left(\\frac{1}{T_1} - \\frac{1}{T_2}\\right)\\right] = %.4f \\times %.3f = %s$. The exponential dominates.",
				math.Pow(t2/t1, 1.5), math.Exp(eg/(2*kB)*(1/t1-1/t2)), ans),
			map[string]string{"ratio": ans})
	}
}

// massAction asks for the minority carrier concentration in an n-type, p-type or compensated sample.
func massAction(ctx generator.GeneratorContext) generator.Problem {
	niM := pickF(1.0, 1.5)
	ni := niM * 1e10
	m := pickF(1, 2, 4, 5, 8)
	e := edakit.Int(15, 18)
	N := m * math.Pow(10, float64(e))
	switch kpOf(ctx) {
	case 1, 2:
		donor := kpOf(ctx) == 1
		if ctx.KPOrdinal == 0 {
			donor = edakit.Int(0, 1) == 0
		}
		kind, maj, mino, dop := "donors", "electrons", "holes", "N_D"
		if !donor {
			kind, maj, mino, dop = "acceptors", "holes", "electrons", "N_A"
		}
		minor := ni * ni / N
		ans := edakit.Num(minor, 0)
		q := edakit.Sprintf("Silicon at 300 K ($n_i = %s$ cm$^{-3}$) is doped with $%s = %s$ cm$^{-3}$ %s, all ionized. What is the concentration of minority %s, in cm$^{-3}$? Round to the nearest whole number.",
			sci(niM, 10), dop, sci(m, e), kind, mino)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The majority %s equal the doping, $%s$ cm$^{-3}$. The mass-action law $np = n_i^2 = %s$ gives the minority %s as $n_i^2/%s = %s$ cm$^{-3}$.", maj, sci(m, e), sciOf(ni*ni), mino, dop, ans),
			map[string]string{"majority": sci(m, e), "minority": ans})
	default:
		na := N
		k := pickF(2, 3, 5, 10)
		nd := na * k
		net := nd - na
		minor := ni * ni / net
		ans := edakit.Num(minor, 0)
		q := edakit.Sprintf("A silicon region at 300 K ($n_i = %s$ cm$^{-3}$) has an implant of $N_A = %s$ cm$^{-3}$ acceptors on top of $N_D = %s$ cm$^{-3}$ donors, all ionized. What is the hole concentration, in cm$^{-3}$? Round to the nearest whole number.",
			sci(niM, 10), sci(m, e), sciOf(nd))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The donors win, so the region is n-type with $n \\approx N_D - N_A = %s$ cm$^{-3}$. Then $p = n_i^2/n = %s/%s = %s$ cm$^{-3}$.", sciOf(net), sciOf(ni*ni), sciOf(net), ans),
			map[string]string{"n": sciOf(net), "p": ans})
	}
}

// fermiPosition asks for the Fermi level relative to the intrinsic level (KP 1), the doping
// that gives a stated offset (KP 2), or the offset below the conduction band edge (KP 3).
func fermiPosition(ctx generator.GeneratorContext) generator.Problem {
	m := pickF(1, 2, 3, 5)
	e := edakit.Int(15, 18)
	N := m * math.Pow(10, float64(e))
	kT := 25.85
	switch kpOf(ctx) {
	case 1:
		ntype := edakit.Int(0, 1) == 0
		d := kT * math.Log(N/1e10)
		ans := edakit.Num(d, 1)
		var q string
		if ntype {
			q = edakit.Sprintf("Silicon at 300 K ($kT = 25.85$ meV, $n_i = 10^{10}$ cm$^{-3}$) is doped n-type with $N_D = %s$ cm$^{-3}$. How far above the intrinsic level $E_i$ is the Fermi level, in meV? Round to 1 decimal place.", sci(m, e))
		} else {
			q = edakit.Sprintf("Silicon at 300 K ($kT = 25.85$ meV, $n_i = 10^{10}$ cm$^{-3}$) is doped p-type with $N_A = %s$ cm$^{-3}$. How far below the intrinsic level $E_i$ is the Fermi level, in meV? Round to 1 decimal place.", sci(m, e))
		}
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$|E_F - E_i| = kT \\ln(N/n_i) = 25.85 \\times \\ln(%s) = 25.85 \\times %.3f = %s$ meV.", sciOf(N/1e10), math.Log(N/1e10), ans),
			map[string]string{"offset_meV": ans})
	case 2:
		d := float64(10 * edakit.Int(25, 50))
		n := 1e10 * math.Exp(d/kT)
		mm, ee := mantissa(n)
		ans := edakit.Num(mm, 2)
		q := edakit.Sprintf("A TCAD run reports that in an n-type region of silicon at 300 K the Fermi level sits %d meV above $E_i$. Taking $kT = 25.85$ meV and $n_i = 10^{10}$ cm$^{-3}$, what is the electron concentration? Give it in units of $10^{%d}$ cm$^{-3}$. Round to 2 decimal places.", int(d), ee)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$n = n_i e^{(E_F - E_i)/kT} = 10^{10} e^{%d/25.85} = 10^{10} e^{%.3f} = %s$ cm$^{-3}$.", int(d), d/kT, sciOf(n)),
			map[string]string{"n": sciOf(n)})
	default:
		nc := pickF(2.8, 3.2)
		d := kT * math.Log(nc*1e19/N)
		ans := edakit.Num(d, 1)
		q := edakit.Sprintf("n-type silicon at 300 K has $N_D = %s$ cm$^{-3}$ and an effective conduction-band density of states $N_c = %s$ cm$^{-3}$. With $kT = 25.85$ meV and non-degenerate statistics, how far below $E_c$ is the Fermi level, in meV? Round to 1 decimal place.", sci(m, e), sci(nc, 19))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$n = N_c e^{-(E_c - E_F)/kT}$, so $E_c - E_F = kT \\ln(N_c/N_D) = 25.85 \\times \\ln(%.1f) = %s$ meV.", nc*1e19/N, ans),
			map[string]string{"offset_meV": ans})
	}
}

// driftMobility asks for a drift velocity (KP 1), a resistivity (KP 2) or a combined mobility
// by Matthiessen's rule (KP 3).
func driftMobility(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		mu := float64(50 * edakit.Int(4, 28))
		E := float64(100 * edakit.Int(2, 30))
		v := mu * E
		ans := edakit.Num(v/1e5, 2)
		q := edakit.Sprintf("Electrons in a silicon resistor have low-field mobility $\\mu_n = %d$ cm$^2$/(V s). A lateral field of %d V/cm is applied, well below velocity saturation. What is the drift velocity, in units of $10^5$ cm/s? Round to 2 decimal places.", int(mu), int(E))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$v_d = \\mu E = %d \\times %d = %s$ cm/s, which is %s in units of $10^5$ cm/s.", int(mu), int(E), sciOf(v), ans),
			map[string]string{"v": sciOf(v)})
	case 2:
		m := pickF(1, 2, 5)
		e := edakit.Int(15, 18)
		n := m * math.Pow(10, float64(e))
		mu := float64(edakit.Pick([]int{1350, 1200, 1000, 800, 500, 300, 150}))
		rho := 1 / (qC * n * mu)
		ans, rt := sigAnswer(rho, 3)
		q := edakit.Sprintf("An n-type silicon layer has $n = %s$ cm$^{-3}$ and electron mobility %d cm$^2$/(V s) at that doping. Ignore holes and take $q = 1.602 \\times 10^{-19}$ C. What is its resistivity, in ohm cm? %s", sci(m, e), int(mu), rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$\\sigma = q n \\mu = 1.602 \\times 10^{-19} \\times %s \\times %d = %.4g$ S/cm, so $\\rho = 1/\\sigma = %s$ ohm cm.", sci(m, e), int(mu), 1/rho, ans),
			map[string]string{"sigma": edakit.Sprintf("%.4g", 1/rho)})
	default:
		muL := float64(edakit.Pick([]int{1400, 1350, 1300}))
		muI := float64(100 * edakit.Int(3, 30))
		mu := 1 / (1/muL + 1/muI)
		ans := edakit.Num(mu, 0)
		q := edakit.Sprintf("In a doped channel, lattice (phonon) scattering alone would give a mobility of %d cm$^2$/(V s) and ionized-impurity scattering alone %d cm$^2$/(V s). Using Matthiessen's rule, what is the combined mobility, in cm$^2$/(V s)? Round to the nearest whole number.", int(muL), int(muI))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Scattering rates add, so $1/\\mu = 1/%d + 1/%d$, giving $\\mu = %s$ cm$^2$/(V s). The combined value is below both, and closer to the smaller one.", int(muL), int(muI), ans),
			map[string]string{"mu": ans})
	}
}

// diffusionEinstein asks for a diffusion coefficient (KP 1), a diffusion current density
// (KP 2) or a diffusion length (KP 3).
func diffusionEinstein(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		mu := float64(50 * edakit.Int(4, 28))
		T := float64(edakit.Pick([]int{250, 300, 350, 400}))
		D := mu * kB * T
		ans := edakit.Num(D, 2)
		q := edakit.Sprintf("Carriers at $T = %d$ K have mobility %d cm$^2$/(V s). Using the Einstein relation with $k/q = 8.617 \\times 10^{-5}$ V/K, what is the diffusion coefficient, in cm$^2$/s? Round to 2 decimal places.", int(T), int(mu))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$D = \\mu kT/q = %d \\times %.5f = %s$ cm$^2$/s.", int(mu), kB*T, ans),
			map[string]string{"kT_q": edakit.Sprintf("%.5f", kB*T)})
	case 2:
		D := float64(edakit.Int(5, 36))
		m := pickF(1, 2, 5)
		e := edakit.Int(15, 17)
		dn := m * math.Pow(10, float64(e))
		dx := pickF(0.5, 1, 2, 5)
		J := qC * D * dn / (dx * 1e-4)
		ans, rt := sigAnswer(J, 3)
		q := edakit.Sprintf("The electron concentration falls linearly by $%s$ cm$^{-3}$ over %s µm in a region with $D_n = %d$ cm$^2$/s. With $q = 1.602 \\times 10^{-19}$ C, what is the magnitude of the electron diffusion current density, in A/cm$^2$? %s",
			sci(m, e), edakit.Num(dx, 1), int(D), rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$|J| = q D_n\\, |dn/dx| = 1.602 \\times 10^{-19} \\times %d \\times %s / (%s \\times 10^{-4}\\ \\text{cm}) = %s$ A/cm$^2$.", int(D), sci(m, e), edakit.Num(dx, 1), ans),
			map[string]string{"gradient": sciOf(dn / (dx * 1e-4))})
	default:
		D := float64(edakit.Int(2, 36))
		tauNs := pickF(10, 50, 100, 500, 1000, 2000)
		L := math.Sqrt(D*tauNs*1e-9) * 1e4
		ans, rt := sigAnswer(L, 3)
		q := edakit.Sprintf("Minority carriers have $D = %d$ cm$^2$/s and lifetime $\\tau = %s$ ns. What is the diffusion length $L = \\sqrt{D\\tau}$, in µm? %s", int(D), edakit.Num(tauNs, 0), rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$L = \\sqrt{%d \\times %s \\times 10^{-9}}$ cm $= %.4g$ cm $= %s$ µm.", int(D), edakit.Num(tauNs, 0), L*1e-4, ans),
			map[string]string{"L_cm": edakit.Sprintf("%.4g", L*1e-4)})
	}
}

var driftDiffusionItems = []edakit.Item{
	{KP: 1, Q: "A device simulator solves the drift-diffusion model. Which set of unknowns does it solve for at every mesh node?",
		A: "Electrostatic potential, electron density and hole density", D: []string{"Drain current, gate charge and threshold voltage", "Temperature only", "Electric field and magnetic field"},
		Why: "The model couples Poisson's equation for the potential with continuity equations for $n$ and $p$. Terminal currents and thresholds are derived afterwards."},
	{KP: 1, Q: "Which equation links the electrostatic potential to the charge in the drift-diffusion model?",
		A: "Poisson's equation, $\\nabla \\cdot (\\varepsilon \\nabla \\psi) = -q(p - n + N_D^+ - N_A^-)$", D: []string{"Ohm's law, $J = \\sigma E$", "The continuity equation, $\\partial n/\\partial t = \\nabla \\cdot J_n / q + G - R$", "The diode equation"},
		Why: "Poisson's equation sets the potential from the net charge: holes and ionized donors positive, electrons and ionized acceptors negative."},
	{KP: 1, Q: "In the electron current equation $J_n = q\\mu_n n E + q D_n \\nabla n$, what do the two terms describe?",
		A: "Drift in the electric field and diffusion down the concentration gradient", D: []string{"Electron and hole currents", "Generation and recombination", "Tunneling and thermionic emission"},
		Why: "Drift-diffusion is named for these two terms. Both are linear in the carrier density, which keeps the model tractable."},
	{KP: 1, Q: "What does the electron continuity equation enforce?",
		A: "Electrons are conserved: the density changes only through current divergence and net generation", D: []string{"The potential is continuous across material interfaces", "The electric field is zero in the bulk", "The total charge is zero everywhere"},
		Why: "$\\partial n/\\partial t = (1/q)\\nabla \\cdot J_n + G - R$ is a conservation law. Every electron that appears either flowed in or was generated."},
	{KP: 2, Q: "In equilibrium the electron current is zero, yet the field inside a junction is large. What must be true?",
		A: "The drift and diffusion currents are equal and opposite at every point", D: []string{"Both drift and diffusion are zero separately", "The mobility is zero in the depletion region", "Only holes move in equilibrium"},
		Why: "Large fields and large gradients coexist in a junction. Zero net current means the two components cancel exactly, which is also how the Einstein relation is derived."},
	{KP: 2, Q: "A simulator discretizes the current equation between two mesh nodes with plain central differences and sees oscillating, negative carrier densities. What is the standard fix?",
		A: "The Scharfetter-Gummel discretization, which assumes an exponential carrier profile between nodes", D: []string{"Use a coarser mesh", "Drop the diffusion term", "Switch to single precision"},
		Why: "Carrier densities vary exponentially with potential. Central differences are unstable when the potential drop between nodes exceeds a few $kT/q$. Scharfetter-Gummel builds the exponential into the flux, which keeps densities positive."},
	{KP: 2, Q: "Why does a device simulator often solve for the potential and quasi-Fermi potentials, or scale densities by $n_i$, instead of solving raw densities?",
		A: "Carrier densities span around 20 orders of magnitude across one device, which wrecks conditioning", D: []string{"Densities cannot be represented in floating point at all", "The potential is the only measurable quantity", "Quasi-Fermi levels are always linear in position"},
		Why: "Densities range from about $10^{0}$ to $10^{20}$ cm$^{-3}$ in one device. Variable choice and scaling keep the Newton Jacobian well conditioned."},
	{KP: 2, Q: "Gummel's method solves Poisson, then the electron equation, then the hole equation, and repeats. When does it converge poorly?",
		A: "At high injection or strong coupling, where the equations depend tightly on each other", D: []string{"At zero bias, where nothing moves", "Whenever the mesh is uniform", "Only for p-type devices"},
		Why: "The decoupled iteration is cheap and robust at low bias, but strong coupling between carriers and potential makes it slow. Simulators then switch to a fully coupled Newton solve."},
	{KP: 3, Q: "A 15 nm channel shows currents above what drift-diffusion predicts even with a saturated-velocity model. Which physics is missing?",
		A: "Non-local transport such as velocity overshoot, which needs hydrodynamic, Monte Carlo or quantum models", D: []string{"Poisson's equation", "The continuity equation", "Ionized-impurity scattering"},
		Why: "Drift-diffusion assumes carriers relax instantly to the local field. In very short channels they cross before relaxing and overshoot the saturation velocity."},
	{KP: 3, Q: "What does a drift-diffusion simulator assume about the carriers that an energy-balance (hydrodynamic) simulator relaxes?",
		A: "Carriers stay at the lattice temperature", D: []string{"Carriers have no charge", "The potential is fixed", "Mobility is independent of doping"},
		Why: "Drift-diffusion uses the lattice temperature in $D = \\mu kT/q$. Hydrodynamic models add an energy equation and a carrier temperature, which matters for hot carriers."},
	{KP: 3, Q: "A TCAD engineer adds the density-gradient correction to a drift-diffusion run on a 1 nm oxide device. What effect is being added?",
		A: "Quantum confinement that pushes the inversion charge away from the oxide interface", D: []string{"Self-heating of the lattice", "Interface trap generation", "Velocity saturation"},
		Why: "Classical statistics put the charge peak at the interface. Density-gradient adds a quantum potential that reproduces the confinement offset, raising the effective oxide thickness and the threshold."},
	{KP: 3, Q: "The drift-diffusion equations are nonlinear. What makes them nonlinear even with constant mobility?",
		A: "Carrier densities depend exponentially on the potential, and the drift term multiplies density by field", D: []string{"Poisson's equation is nonlinear in the permittivity", "The mesh moves during the solve", "Diffusion is quadratic in the gradient"},
		Why: "The product $n E$ couples two unknowns, and in equilibrium $n \\propto e^{q\\psi/kT}$. Hence a Newton solve on a sparse Jacobian."},
}
