package fab

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// interconnectStack: track supply per layer and across the stack, and the R penalty of thin layers.
func interconnectStack(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		W := float64(edakit.Int(5, 200))
		p := float64(2 * edakit.Int(12, 40))
		n := int(math.Floor(W * 1000 / p))
		ans := edakit.IntStr(n)
		q := edakit.Sprintf("A block is %.0f um wide. A vertical routing layer has a %.0f nm pitch. How many vertical tracks cross the block?", W, p)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$\\lfloor %.0f \\times 1000 / %.0f \\rfloor = %s$ tracks.", W, p, ans),
			map[string]string{"tracks": ans})
	case 2:
		W := float64(edakit.Int(10, 60))
		H := float64(edakit.Int(10, 60))
		nl := edakit.Scaled(ctx, 3, 6)
		total := 0
		var parts []string
		for i := 0; i < nl; i++ {
			p := float64(edakit.Pick([]int{28, 32, 36, 40, 48, 64, 80}))
			if i%2 == 0 {
				parts = append(parts, edakit.Sprintf("M%d horizontal at %.0f nm", i+2, p))
				total += int(math.Floor(H * 1000 / p))
			} else {
				parts = append(parts, edakit.Sprintf("M%d vertical at %.0f nm", i+2, p))
				total += int(math.Floor(W * 1000 / p))
			}
		}
		ans := edakit.IntStr(total)
		q := edakit.Sprintf("A block is %.0f um wide and %.0f um tall. Its routing layers are %s. Horizontal tracks span the width and are counted across the height; vertical tracks the reverse. What is the total number of routing tracks?", W, H, strings.Join(parts, ", "))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Horizontal layers give height/pitch tracks, vertical layers width/pitch, each floored; the sum is %s tracks.", ans),
			map[string]string{"tracks": ans})
	default:
		p1 := float64(2 * edakit.Int(12, 24))
		mult := edakit.Pick([]float64{1.5, 2, 3, 4, 8})
		p2 := p1 * mult
		ans := edakit.Num(mult*mult, 2)
		q := edakit.Sprintf("A thin local layer has a %.0f nm pitch and a semi-global layer a %.0f nm pitch. Both have width equal to half the pitch and the same aspect ratio, and the same resistivity. How many times larger is the resistance per unit length of the thin layer? Round to 2 decimal places.", p1, p2)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Cross-section area scales with pitch squared at fixed aspect ratio, so $R/L$ scales as $(%.0f/%.0f)^2 = %s$. Long routes are promoted to upper layers for this reason.", p2, p1, ans),
			map[string]string{"ratio": ans})
	}
}

// damascene: barrier-limited copper cross-section and resistance.
func damascene(ctx generator.GeneratorContext) generator.Problem {
	w := float64(2 * edakit.Int(6, 30))
	ar := step(1.5, 2.5, 0.5)
	h := w * ar
	b := float64(edakit.Int(1, 4))
	acu := (w - 2*b) * (h - b)
	base := edakit.Sprintf("A dual-damascene copper trench is %.0f nm wide and %.0f nm tall. A barrier and liner of %.0f nm lines both sidewalls and the bottom.", w, h, b)
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(acu/(w*h)*100, 1)
		q := base + " What percentage of the trench cross-section is copper? Round to 1 decimal place."
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Copper fills $(%.0f - 2 \\times %.0f) \\times (%.0f - %.0f) = %.0f$ nm^2 of %.0f nm^2, which is %s%%.", w, b, h, b, acu, w*h, ans),
			map[string]string{"acu": edakit.Num(acu, 0)})
	case 2:
		rho := step(2.5, 6.0, 0.1)
		r := rho * 10 * 1000 / acu
		ans := edakit.Num(r, 1)
		q := base + edakit.Sprintf(" The copper's effective resistivity at this width is %s micro-ohm cm, and the barrier carries negligible current. What is the wire resistance per micrometre, in ohms? Round to 1 decimal place.", edakit.Num(rho, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("1 micro-ohm cm is 10 ohm nm. $R/L = \\rho / A = %s \\times 1000 / %.0f = %s$ ohm per um.", edakit.Num(rho*10, 0), acu, ans),
			map[string]string{"acu": edakit.Num(acu, 0)})
	default:
		ans := edakit.Num(w*h/acu, 3)
		q := base + " At the same copper resistivity, by what factor is the resistance higher than if copper filled the whole trench? Round to 3 decimal places."
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The factor is $A_{trench}/A_{Cu} = %.0f / %.0f = %s$. Barriers do not shrink with the wire, so the penalty grows each node.", w*h, acu, ans),
			map[string]string{"acu": edakit.Num(acu, 0)})
	}
}

// wireRC: resistance per length, distributed RC delay, and scaling.
func wireRC(ctx generator.GeneratorContext) generator.Problem {
	w := float64(2 * edakit.Int(6, 40))
	t := w * step(1.5, 2.5, 0.5)
	rho := step(2.0, 6.0, 0.1)
	r := rho * 10 * 1000 / (w * t)
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(r, 2)
		q := edakit.Sprintf("A wire is %.0f nm wide and %.0f nm thick with effective resistivity %s micro-ohm cm. What is its resistance per micrometre, in ohms? Round to 2 decimal places.", w, t, edakit.Num(rho, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$r = \\rho/(w t) = %s$ ohm nm $\\times 1000$ nm $/ %.0f$ nm^2 $= %s$ ohm per um.", edakit.Num(rho*10, 0), w*t, ans),
			map[string]string{"r": ans})
	case 2:
		rr := step(5, 100, 5)
		c := step(0.10, 0.25, 0.01)
		L := float64(10 * edakit.Int(2, 30))
		d := 0.5 * rr * c * L * L / 1000
		ans := edakit.Num(d, 2)
		q := edakit.Sprintf("An unbuffered wire of length %.0f um has $r = %s$ ohm/um and $c = %s$ fF/um. Ignoring driver and load, what is its distributed Elmore delay $r c L^2 / 2$, in picoseconds? Round to 2 decimal places.", L, edakit.Num(rr, 0), edakit.Num(c, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$0.5 \\times %s \\times %s \\times %.0f^2 = %s$ ohm fF, and 1000 ohm fF is 1 ps, giving %s ps. Delay grows with length squared, which is why long wires are buffered.", edakit.Num(rr, 0), edakit.Num(c, 2), L, edakit.Num(d*1000, 0), ans),
			map[string]string{"rcl2": edakit.Num(d*1000, 0)})
	default:
		d0 := step(5, 80, 0.5)
		s := step(0.6, 0.85, 0.05)
		rf := step(1.0, 1.5, 0.05)
		d1 := d0 * rf / (s * s)
		ans := edakit.Num(d1, 2)
		q := edakit.Sprintf("A fixed-length route has %s ps of wire delay. The next node shrinks wire width and thickness by a factor of %s each, capacitance per unit length stays the same, and size effects raise resistivity by a factor of %s. What is the new wire delay for the same route length, in picoseconds? Round to 2 decimal places.", edakit.Num(d0, 1), edakit.Num(s, 2), edakit.Num(rf, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Resistance per length scales by $%s / %s^2 = %s$ and $c$ is unchanged, so delay becomes $%s \\times %s = %s$ ps. Wires get slower as they shrink while gates get faster.", edakit.Num(rf, 2), edakit.Num(s, 2), edakit.Num(rf/(s*s), 3), edakit.Num(d0, 1), edakit.Num(rf/(s*s), 3), ans),
			map[string]string{"r_scale": edakit.Num(rf/(s*s), 3)})
	}
}

// lowK: capacitance with dielectric constant, effective k of a stack, and partial air gaps.
func lowK(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		c0 := step(0.12, 0.30, 0.01)
		k0 := edakit.Pick([]float64{4.0, 3.9, 3.0, 2.9})
		k1 := step(2.2, 2.8, 0.1)
		ans := edakit.Num(c0*k1/k0, 3)
		q := edakit.Sprintf("A wire has %s fF/um of capacitance in a dielectric of $k = %s$. If every dielectric around it is replaced by one of $k = %s$ with the same geometry, what is the new capacitance per micrometre, in fF? Round to 3 decimal places.", edakit.Num(c0, 2), edakit.Num(k0, 1), edakit.Num(k1, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Capacitance is proportional to permittivity at fixed geometry: $%s \\times %s / %s = %s$ fF/um.", edakit.Num(c0, 2), edakit.Num(k1, 1), edakit.Num(k0, 1), ans),
			map[string]string{"c": ans})
	case 2:
		t1 := float64(edakit.Int(5, 30))
		k1 := step(4.5, 7.0, 0.5)
		t2 := float64(10 * edakit.Int(5, 20))
		k2 := step(2.2, 3.0, 0.1)
		keff := (t1 + t2) / (t1/k1 + t2/k2)
		ans := edakit.Num(keff, 2)
		q := edakit.Sprintf("Between two metal layers sit a %.0f nm etch-stop cap with $k = %s$ and a %.0f nm low-k film with $k = %s$. Treating them as plates in series, what is the effective dielectric constant of the stack? Round to 2 decimal places.", t1, edakit.Num(k1, 1), t2, edakit.Num(k2, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Series capacitors: $k_{eff} = (t_1 + t_2)/(t_1/k_1 + t_2/k_2) = %.0f / %s = %s$. A thin high-k cap costs more than its thickness suggests.", t1+t2, edakit.Num(t1/k1+t2/k2, 3), ans),
			map[string]string{"keff": ans})
	default:
		f := step(0.4, 0.8, 0.05)
		k := step(2.4, 3.0, 0.1)
		kg := step(1.0, 1.8, 0.1)
		red := f * (1 - kg/k) * 100
		ans := edakit.Num(red, 1)
		q := edakit.Sprintf("Of a wire's total capacitance, a fraction %s is to its same-layer neighbours. Air gaps lower the dielectric between neighbours from $k = %s$ to an effective %s, and the rest is unchanged. By what percentage does total capacitance, and so RC delay, fall? Round to 1 decimal place.", edakit.Num(f, 2), edakit.Num(k, 1), edakit.Num(kg, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Only the lateral part scales: reduction $= %s \\times (1 - %s/%s) = %s$%%.", edakit.Num(f, 2), edakit.Num(kg, 1), edakit.Num(k, 1), ans),
			map[string]string{"reduction": ans})
	}
}

// viasContact: specific contact resistivity, via stacks, and the contact share of a path.
func viasContact(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		m := step(0.5, 5.0, 0.5)
		a := float64(edakit.Int(10, 40))
		bb := float64(edakit.Int(10, 40))
		R := m * 1e-9 / (a * bb * 1e-14)
		ans := edakit.Num(R, 1)
		q := edakit.Sprintf("A contact to silicide has specific contact resistivity $\\rho_c = %s \\times 10^{-9}$ ohm cm^2 and an interface of %.0f nm by %.0f nm. What is its interface resistance $\\rho_c / A$, in ohms? Round to 1 decimal place.", edakit.Num(m, 1), a, bb)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$A = %.0f$ nm^2 $= %.4g$ cm^2, so $R = %s \\times 10^{-9} / %.4g = %s$ ohm. Halving the area doubles it.", a*bb, a*bb*1e-14, edakit.Num(m, 1), a*bb*1e-14, ans),
			map[string]string{"area_cm2": edakit.Sprintf("%.4g", a*bb*1e-14)})
	case 2:
		nl := edakit.Scaled(ctx, 2, 6)
		total := 0.0
		var parts []string
		for i := 0; i < nl; i++ {
			rv := float64(5 * edakit.Int(2, 20))
			k := edakit.Int(1, 3)
			total += rv / float64(k)
			if k == 1 {
				parts = append(parts, edakit.Sprintf("V%d: one via of %.0f ohm", i, rv))
			} else {
				parts = append(parts, edakit.Sprintf("V%d: %d parallel vias of %.0f ohm each", i, k, rv))
			}
		}
		ans := edakit.Num(total, 2)
		q := "A signal climbs a via stack: " + strings.Join(parts, "; ") + ". Ignoring the short metal landings, what is the total stack resistance, in ohms? Round to 2 decimal places."
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Levels are in series; parallel vias at one level divide that level's resistance. The sum is %s ohm.", ans),
			map[string]string{"r_stack": ans})
	default:
		rc := float64(10 * edakit.Int(5, 40))
		rv := float64(10 * edakit.Int(2, 20))
		rw := float64(10 * edakit.Int(5, 100))
		frac := 2 * rc / (2*rc + rv + rw) * 100
		ans := edakit.Num(frac, 1)
		q := edakit.Sprintf("A short local connection runs from one transistor's contact, through vias of %.0f ohm total and a wire of %.0f ohm, into another transistor's contact. Each contact is %.0f ohm. What percentage of the path resistance is in the two contacts? Round to 1 decimal place.", rv, rw, rc)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Total is $2 \\times %.0f + %.0f + %.0f = %.0f$ ohm; contacts are $%.0f / %.0f = %s$%%. At short distances, MOL contact resistance dominates.", rc, rv, rw, 2*rc+rv+rw, 2*rc, 2*rc+rv+rw, ans),
			map[string]string{"total": edakit.Num(2*rc+rv+rw, 0)})
	}
}

var feolMolBeolItems = []edakit.Item{
	{KP: 1, Q: "Which part of the process forms the transistors themselves: wells, isolation, gates, spacers and source-drain?",
		A:   "The front end of line (FEOL)",
		D:   []string{"The middle of line (MOL)", "The back end of line (BEOL)", "The packaging flow"},
		Why: "FEOL builds the devices in and on the silicon before any contact or metal is formed."},
	{KP: 1, Q: "Trench contacts to source and drain and the gate contact belong to which segment?",
		A:   "MOL",
		D:   []string{"FEOL", "BEOL", "Wafer sort"},
		Why: "The middle of line is the local contact layer that joins the devices to the first metal."},
	{KP: 1, Q: "Which segment contains the copper wiring layers and the vias between them?",
		A:   "BEOL",
		D:   []string{"FEOL", "MOL", "Epitaxy"},
		Why: "The back end of line is the multilevel interconnect stack, built after the devices and contacts."},
	{KP: 1, Q: "Why did foundries start naming a separate middle of line?",
		A:   "Local contacts became tight and resistive enough to need their own materials, rules and models",
		D:   []string{"To add more transistors", "Because FEOL was removed", "Because copper is no longer used"},
		Why: "At small pitches, contact and local interconnect resistance and capacitance became first-order, so the MOL got its own integration and design rules."},
	{KP: 2, Q: "A designer sees separate layer names for the gate contact, the trench contact and V0. Which tool most needs their resistances right?",
		A:   "Parasitic extraction feeding timing",
		D:   []string{"Logic synthesis's technology-independent pass", "The equivalence checker", "The HDL simulator"},
		Why: "Contact and V0 resistance are now a large share of a cell's output resistance; extraction must model them for accurate timing."},
	{KP: 2, Q: "Contact over active gate (COAG) lets the gate contact land above the active region rather than over isolation. What design benefit follows?",
		A:   "Smaller cell height or freer pin placement, because the gate contact no longer needs its own space",
		D:   []string{"Lower threshold voltage", "Fewer metal layers", "No need for vias"},
		Why: "Without COAG the gate contact sits on isolation between the n and p regions, which costs height. COAG removes that constraint at the price of a harder self-aligned process."},
	{KP: 2, Q: "Why are FEOL layers usually drawn only inside standard cells while BEOL layers are mostly drawn by the router?",
		A:   "Devices are fixed inside library cells, and the routing between cells is design specific",
		D:   []string{"FEOL cannot be drawn by any tool", "BEOL is generated by the foundry", "Routers cannot see metal layers"},
		Why: "The cell library fixes transistors and their local contacts; the router lays out the inter-cell wiring in the BEOL layers."},
	{KP: 2, Q: "A BEOL via must not land on a dielectric. Which process fact explains the self-aligned via processes at tight pitch?",
		A:   "Overlay error can shift a via onto the neighbouring line, so the via is confined by the line trench",
		D:   []string{"Vias are formed before the transistors", "Vias are printed with no lithography", "Copper cannot fill vias"},
		Why: "At small pitches, overlay error is a large fraction of spacing; self-aligned vias use the line pattern to clip the via and prevent shorts."},
	{KP: 3, Q: "A process change thins the MOL dielectric between the gate and the trench contact. What changes for the designer?",
		A:   "Gate-to-contact capacitance rises, slowing cells, and the short risk between them increases",
		D:   []string{"Only the top metal resistance changes", "Leakage through the substrate falls", "Nothing outside FEOL"},
		Why: "The gate and trench contacts run side by side; a thinner spacer raises their coupling capacitance and tightens the reliability margin."},
	{KP: 3, Q: "Why does the BEOL have a different thermal budget from the FEOL?",
		A:   "Copper and low-k films and completed devices cannot tolerate FEOL temperatures",
		D:   []string{"The BEOL is built on a separate wafer", "Lithography fails above room temperature", "Copper requires high heat to conduct"},
		Why: "Once silicides, devices and copper exist, process steps stay near 400 C to avoid diffusion, damage and stress."},
	{KP: 3, Q: "A new node replaces copper with another metal for the tightest-pitch layers. Why might a lower-bulk-conductivity metal win?",
		A:   "At small widths a metal needing a thin or no barrier and with less size-effect penalty can give lower line resistance",
		D:   []string{"Bulk conductivity always decides", "Copper cannot be patterned at any size", "The new metal has a lower dielectric constant"},
		Why: "Copper needs a barrier and suffers strong surface and grain-boundary scattering when narrow; metals such as ruthenium or molybdenum are considered for this reason."},
	{KP: 3, Q: "Which segment does backside power delivery move off the front side of the wafer?",
		A:   "The power distribution part of the BEOL",
		D:   []string{"The FEOL transistors", "The MOL contacts entirely", "The wafer sort step"},
		Why: "Backside power builds the power network on the wafer back, freeing front-side BEOL tracks for signals."},
}
