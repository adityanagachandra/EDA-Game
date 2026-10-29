package dev

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// soiBackBias computes an FD-SOI threshold under back-gate bias.
func soiBackBias(ctx generator.GeneratorContext) generator.Problem {
	vt := 10 * edakit.Int(30, 45)
	r := 5 * edakit.Int(10, 20)
	vbb := 0.1 * float64(edakit.Int(3, 15))
	forward := edakit.Int(0, 1) == 0
	shift := float64(r) * vbb
	var ans float64
	kind := "forward"
	if forward {
		ans = float64(vt) - shift
	} else {
		kind = "reverse"
		ans = float64(vt) + shift
	}
	q := edakit.Sprintf("An FD-SOI NMOS has $V_t = %d$ mV with its back gate at 0 V and a back-gate coupling of %d mV per volt of back bias. A body-bias generator applies %s V of %s back bias. What is the new threshold, in mV?",
		vt, r, edakit.Num(vbb, 1), kind)
	wrong := float64(vt) + shift
	if !forward {
		wrong = float64(vt) - shift
	}
	return edakit.Choice(q, edakit.Num(ans, 0),
		edakit.Sprintf("The thin buried oxide lets the substrate act as a second gate. The shift is $%d \\times %s = %s$ mV. Forward back bias lowers an NMOS threshold and reverse raises it, so $V_t = %s$ mV.", r, edakit.Num(vbb, 1), edakit.Num(shift, 0), edakit.Num(ans, 0)),
		edakit.Num(wrong, 0), edakit.Num(float64(vt), 0), edakit.Num(ans-shift*map[bool]float64{true: 1, false: -1}[forward], 0))
}

var soiItems = []edakit.Item{
	{KP: 1, Q: "What distinguishes fully depleted SOI from partially depleted SOI?",
		A: "The silicon film is thin enough that the gate depletes it completely, leaving no neutral floating body", D: []string{"FD-SOI has no buried oxide", "PD-SOI uses a thinner film than FD-SOI", "FD-SOI devices have no threshold voltage"},
		Why: "In PD-SOI the film is thicker than the maximum depletion depth, so a neutral body region floats. FD-SOI films are a few nm to tens of nm thick."},
	{KP: 1, Q: "A PD-SOI circuit shows delays that depend on how long the gate has been idle. What causes this history effect?",
		A: "Charge collects in the floating body and shifts the threshold over time", D: []string{"The buried oxide charges like a capacitor in the supply network", "Electromigration in the contacts", "The clock tree heats up"},
		Why: "Impact ionization and junction leakage charge the floating body, and the body potential sets $V_t$ through the body effect. Its value depends on switching history."},
	{KP: 1, Q: "What does the buried oxide in SOI remove compared with bulk CMOS?",
		A: "Most source and drain junction capacitance to the substrate and the parasitic path that enables latch-up", D: []string{"All gate capacitance", "The need for a gate dielectric", "Subthreshold leakage entirely"},
		Why: "The junctions sit on oxide, so their bottom-wall capacitance and the well-substrate thyristor largely disappear."},
	{KP: 1, Q: "Why does undoped-channel FD-SOI show less random threshold variation than heavily doped bulk devices?",
		A: "Short-channel control comes from the thin film, so the channel needs few dopants that could fluctuate", D: []string{"The buried oxide averages out the dopants", "FD-SOI devices are always wider", "Variation is set only by the metal stack"},
		Why: "Random dopant fluctuation scales with channel doping. A thin film suppresses short-channel effects without heavy doping."},
	{KP: 3, Q: "A library team characterizes an FD-SOI process with several back-bias values. Why?",
		A: "Back bias tunes threshold, speed and leakage after fabrication, so timing and power depend on it", D: []string{"Back bias changes the cell area", "Back bias only matters for analog blocks", "Libraries must contain one table per metal layer"},
		Why: "Adaptive body bias trades speed for leakage at run time. Implementation tools need tables at each bias to sign off the chosen operating points."},
	{KP: 3, Q: "The silicon film in an SOI device sits on an oxide that conducts heat poorly. What must a compact model add?",
		A: "Self-heating: a thermal node that raises channel temperature with dissipated power", D: []string{"Nothing: SOI runs cooler than bulk", "A larger gate resistance", "A second threshold voltage for each corner"},
		Why: "The buried oxide's thermal conductivity is about a hundredth of silicon's. SOI compact models include a thermal network so current falls under sustained power."},
	{KP: 3, Q: "In FD-SOI, what limits how thin the silicon film can usefully be made?",
		A: "Thickness variation and quantum confinement shift the threshold, and mobility falls in very thin films", D: []string{"The buried oxide stops working below 50 nm", "Gate capacitance rises without limit", "Thin films cannot be doped p-type"},
		Why: "Confinement energy grows as $1/t^2$, so a nanometre of film variation becomes a threshold shift, and surface scattering lowers mobility."},
	{KP: 3, Q: "A circuit simulator runs a PD-SOI netlist transient and sees a slowly drifting output. Which model element produces this?",
		A: "The floating-body node with its slow charging and discharging time constants", D: []string{"The overlap capacitance", "The series resistance", "The gate leakage current only"},
		Why: "PD-SOI compact models carry an internal body node. Its long time constants require transient runs long enough to settle the body state."},
}

// finfet asks for an effective width (KP 1), a drive current (KP 2) or the fin count for a
// target current (KP 3).
func finfet(ctx generator.GeneratorContext) generator.Problem {
	h := float64(edakit.Int(30, 55))
	wf := float64(edakit.Int(5, 8))
	n := float64(edakit.Int(1, 4))
	weff := n * (2*h + wf)
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(weff, 0)
		q := edakit.Sprintf("A FinFET has %d fin(s), each %d nm tall and %d nm wide, with the gate wrapping both sidewalls and the top. What is its effective electrical width, in nm?", int(n), int(h), int(wf))
		return edakit.Numeric(q, ans, edakit.Sprintf("Each fin contributes $2H_{fin} + W_{fin} = %d$ nm, so $W_{eff} = %d \\times %d = %s$ nm.", int(2*h+wf), int(n), int(2*h+wf), ans), map[string]string{"per_fin": edakit.IntStr(int(2*h + wf))})
	case 2:
		ion := float64(50 * edakit.Int(12, 36))
		i := weff / 1000 * ion
		ans := edakit.Num(i, 1)
		fins := "fins"
		if n == 1 {
			fins = "fin"
		}
		q := edakit.Sprintf("A process quotes NMOS drive current as %d µA per µm of effective width. A FinFET has %d %s of height %d nm and width %d nm. What is its drive current, in µA? Round to 1 decimal place.", int(ion), int(n), fins, int(h), int(wf))
		return edakit.Numeric(q, ans, edakit.Sprintf("$W_{eff} = %d \\times (2 \\times %d + %d) = %d$ nm $= %s$ µm, so $I = %s \\times %d = %s$ µA.", int(n), int(h), int(wf), int(weff), edakit.Num(weff/1000, 3), edakit.Num(weff/1000, 3), int(ion), ans), map[string]string{"weff": edakit.IntStr(int(weff))})
	default:
		ifin := float64(5 * edakit.Int(6, 16))
		target := float64(10 * edakit.Int(5, 60))
		k := int(math.Ceil(target/ifin - 1e-9))
		ans := edakit.IntStr(k)
		q := edakit.Sprintf("In a FinFET process each fin of a given device delivers %d µA. A driver needs at least %d µA. Width can only grow in whole fins. What is the smallest number of fins that meets the requirement?", int(ifin), int(target))
		return edakit.Numeric(q, ans, edakit.Sprintf("$\\lceil %d/%d \\rceil = \\lceil %.3f \\rceil = %d$ fins. Width is quantized, so the delivered current %d µA overshoots the need.", int(target), int(ifin), target/ifin, k, k*int(ifin)), map[string]string{"fins": ans})
	}
}

// gaa asks for a nanosheet effective width (KP 1), its ratio to a FinFET in the same footprint
// (KP 2) or the sheet width needed for a target effective width (KP 3).
func gaa(ctx generator.GeneratorContext) generator.Problem {
	n := float64(edakit.Int(2, 4))
	w := float64(5 * edakit.Int(3, 10))
	t := float64(edakit.Int(5, 7))
	weff := n * 2 * (w + t)
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(weff, 0)
		q := edakit.Sprintf("A gate-all-around device stacks %d nanosheets, each %d nm wide and %d nm thick, with the gate wrapping every sheet completely. What is its effective electrical width, in nm?", int(n), int(w), int(t))
		return edakit.Numeric(q, ans, edakit.Sprintf("Each sheet's perimeter is $2(W + T) = %d$ nm, so $W_{eff} = %d \\times %d = %s$ nm.", int(2*(w+t)), int(n), int(2*(w+t)), ans), map[string]string{"perimeter": edakit.IntStr(int(2 * (w + t)))})
	case 2:
		nf := float64(edakit.Int(2, 3))
		h := float64(edakit.Int(40, 55))
		wf := float64(edakit.Int(5, 7))
		wfin := nf * (2*h + wf)
		r := weff / wfin
		ans := edakit.Num(r, 2)
		q := edakit.Sprintf("In one cell footprint, a FinFET uses %d fins (height %d nm, width %d nm) and a nanosheet alternative uses %d sheets (width %d nm, thickness %d nm). What is the ratio of nanosheet effective width to FinFET effective width? Round to 2 decimal places.", int(nf), int(h), int(wf), int(n), int(w), int(t))
		return edakit.Numeric(q, ans, edakit.Sprintf("FinFET: $%d \\times (2 \\times %d + %d) = %d$ nm. Nanosheet: $%d \\times 2 \\times (%d + %d) = %d$ nm. Ratio $= %s$.", int(nf), int(h), int(wf), int(wfin), int(n), int(w), int(t), int(weff), ans), map[string]string{"weff_fin": edakit.IntStr(int(wfin)), "weff_ns": edakit.IntStr(int(weff))})
	default:
		target := float64(10 * edakit.Int(15, 40))
		need := target/(2*n) - t
		for need < 5 {
			target += 50
			need = target/(2*n) - t
		}
		ans := edakit.Num(need, 1)
		q := edakit.Sprintf("A nanosheet process stacks %d sheets of thickness %d nm, and the sheet width is a layout choice. What sheet width gives an effective width of %d nm, in nm? Round to 1 decimal place.", int(n), int(t), int(target))
		return edakit.Numeric(q, ans, edakit.Sprintf("$W_{eff} = 2N(W + T)$, so $W = W_{eff}/(2N) - T = %d/%d - %d = %s$ nm. Unlike fins, sheet width varies continuously.", int(target), int(2*n), int(t), ans), map[string]string{"W": ans})
	}
}

// confinement asks for a ground-state confinement energy (KP 1), the first subband spacing
// (KP 2) or the effective band-gap widening (KP 3) in an infinite-well model.
func confinement(ctx generator.GeneratorContext) generator.Problem {
	t := float64(edakit.Int(3, 10))
	c := 0.0381 * math.Pi * math.Pi // eV nm^2, for E1 = c/(m t^2)
	stem := "Model the body as an infinite square well, $E_n = n^2 \\pi^2 \\hbar^2/(2 m^* t^2)$, with $\\hbar^2/(2m_0) = 0.0381$ eV nm$^2$."
	switch kpOf(ctx) {
	case 1:
		m := pickF(0.19, 0.26, 0.5, 0.916)
		e1 := c / (m * t * t) * 1000
		ans := edakit.Num(e1, 1)
		q := edakit.Sprintf("An electron is confined across a silicon body $t = %d$ nm thick, with confinement effective mass $%s\\, m_0$. %s What is the ground-state energy above the band edge, in meV? Round to 1 decimal place.", int(t), edakit.Num(m, 3), stem)
		return edakit.Numeric(q, ans, edakit.Sprintf("$E_1 = \\pi^2 \\times 0.0381 / (%s \\times %d^2)$ eV $= %s$ meV.", edakit.Num(m, 3), int(t), ans), map[string]string{"E1": ans})
	case 2:
		m := pickF(0.19, 0.26, 0.5, 0.916)
		e1 := c / (m * t * t) * 1000
		ans := edakit.Num(3*e1, 1)
		q := edakit.Sprintf("Electrons are confined in a nanosheet $t = %d$ nm thick, with confinement mass $%s\\, m_0$. %s What is the spacing between the first and second subbands, in meV? Round to 1 decimal place.", int(t), edakit.Num(m, 3), stem)
		return edakit.Numeric(q, ans, edakit.Sprintf("$E_1 = %.2f$ meV and $E_2 = 4E_1$, so $E_2 - E_1 = 3E_1 = %s$ meV. When this spacing is several $kT$, most carriers sit in the lowest subband.", e1, ans), map[string]string{"E1": edakit.Num(e1, 2)})
	default:
		me := pickF(0.19, 0.916)
		mh := pickF(0.29, 0.49)
		de := c / (t * t) * (1/me + 1/mh) * 1000
		ans := edakit.Num(de, 1)
		q := edakit.Sprintf("In a %d nm thick silicon body, electrons have confinement mass $%s\\, m_0$ and holes $%s\\, m_0$. %s By how much does confinement widen the effective band gap (the sum of both ground-state shifts), in meV? Round to 1 decimal place.", int(t), edakit.Num(me, 3), edakit.Num(mh, 2), stem)
		return edakit.Numeric(q, ans, edakit.Sprintf("The electron ground state rises by %.2f meV and the hole ground state falls by %.2f meV, so the gap widens by %s meV.", c/(me*t*t)*1000, c/(mh*t*t)*1000, ans), map[string]string{"dEg": ans})
	}
}

var cfetItems = []edakit.Item{
	{KP: 1, Q: "What does a complementary FET (CFET) stack vertically?",
		A: "An NMOS and a PMOS device, one above the other, sharing a footprint", D: []string{"Two NMOS devices in series", "A transistor above its own metal interconnect", "A DRAM capacitor above a transistor"},
		Why: "A CFET places the n and p nanosheet stacks on top of each other instead of side by side, so one footprint carries both halves of a CMOS pair."},
	{KP: 1, Q: "Why does stacking the n and p devices promise a smaller standard cell?",
		A: "The n-to-p spacing that sets cell height in side-by-side layouts disappears", D: []string{"Each transistor becomes twice as fast", "Metal layers are no longer needed", "The supply voltage halves"},
		Why: "Planar, FinFET and nanosheet cells place n and p regions side by side with spacing between them. Stacking removes that lateral separation, which shrinks cell height."},
	{KP: 1, Q: "What are the two main fabrication approaches discussed for CFET?",
		A: "Monolithic, building both tiers on one wafer in sequence, and sequential, bonding a second device layer on top", D: []string{"Bulk and SOI", "Planar and FinFET", "Front-side and back-side lithography only"},
		Why: "Monolithic integration patterns both tiers together. Sequential integration builds the second tier after bonding, which loosens alignment but limits the thermal budget."},
	{KP: 1, Q: "How does a CFET relate to the nanosheet transistor?",
		A: "It is built from gate-all-around nanosheet stacks, with the n and p stacks placed vertically", D: []string{"It replaces nanosheets with planar devices", "It uses FinFETs with the gate on top only", "It is a bipolar transistor"},
		Why: "CFET is a layout of gate-all-around devices, not a new channel physics. It inherits nanosheet electrostatics."},
	{KP: 2, Q: "In a CFET cell, the input must reach both the top and bottom device gates. What does a layout tool now have to model?",
		A: "Vertical connections between tiers, and which tier each pin and contact lives on", D: []string{"Nothing new: the cell is still planar", "Only the top device, since the bottom one is hidden", "A separate netlist for each tier with no shared nodes"},
		Why: "A shared gate, and the output node joining top and bottom drains, need vertical contacts. Cell layout and extraction must represent tier-aware geometry."},
	{KP: 2, Q: "Why does CFET make backside power delivery attractive?",
		A: "Front-side routing tracks are scarce in a shorter cell, and the bottom device is close to the wafer back", D: []string{"Backside power removes the need for vias", "CFET devices cannot be powered from the front", "Backside power doubles mobility"},
		Why: "Shrinking cell height squeezes signal and power rails. Moving power to the back frees front tracks and gives the bottom tier a short path to its supply."},
	{KP: 2, Q: "A parasitic extractor for a CFET process must add a coupling that side-by-side devices did not have. Which one?",
		A: "Capacitance between the stacked top and bottom devices and their contacts", D: []string{"Coupling between metal layers that are not adjacent", "Substrate capacitance of the buried oxide only", "Inductance of the gate"},
		Why: "Stacking places the n and p devices and their local interconnect within nanometres vertically, so inter-tier capacitance appears in the netlist."},
	{KP: 2, Q: "A placer's cell library moves from nanosheet to CFET with the same cell width but fewer routing tracks. What problem should the router designer expect?",
		A: "Pin access gets harder, so pin-access-aware placement and routing matter more", D: []string{"Routing gets easier because cells are shorter", "No change, since pins are unchanged", "Timing analysis becomes unnecessary"},
		Why: "Fewer tracks per cell means fewer legal access points per pin. Pin access becomes a first-class constraint for detailed placement and routing."},
	{KP: 3, Q: "Why is thermal management a concern for stacked devices such as CFET?",
		A: "Stacking puts more power in the same footprint and buries one tier farther from the heat sink path", D: []string{"Stacked devices dissipate no heat", "Heat flows sideways only in bulk silicon", "The gate oxide insulates the whole stack perfectly"},
		Why: "Power density per footprint rises, and the inner tier is thermally isolated by surrounding layers. Device and system models need self-heating terms."},
	{KP: 3, Q: "A compact-model team is asked to support CFET. What can they largely reuse?",
		A: "Gate-all-around channel models, with added coupling and parasitics between tiers", D: []string{"Nothing: CFET needs an entirely new transport theory", "Bulk planar models with no changes", "Bipolar transistor models"},
		Why: "Each tier is a nanosheet device. The new work is in tier-to-tier coupling, separate tier thresholds and the parasitic network."},
	{KP: 3, Q: "Why should claims about CFET density gains be read with care?",
		A: "The gain depends on cell architecture, routing tracks and contact schemes, and production results are not yet established", D: []string{"CFET has already doubled density in every product", "Density gains are fixed by physics at exactly 2x", "CFET increases cell area"},
		Why: "Stacking removes the n-p space, but contacts, vertical interconnect and track count decide the real cell size. Treat early numbers as projections."},
	{KP: 3, Q: "In a CFET, the top and bottom devices may receive different thermal budgets during fabrication. What modeling consequence follows?",
		A: "The two tiers can have different thresholds and mobilities and need separate model parameters", D: []string{"Both tiers must share every parameter", "The top tier has no threshold", "Thermal budget only affects metal resistance"},
		Why: "Process differences between tiers show up as different device characteristics, so extraction and corners must treat tiers separately."},
}
