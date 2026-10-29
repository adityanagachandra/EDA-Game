package pcb

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var viaTypeItems = []edakit.Item{
	{KP: 1, Q: "Which via joins an outer layer to an inner layer without passing through the whole board?",
		A: "A blind via", D: []string{"A buried via", "A through via", "A thermal via"},
		Why: "A blind via is visible on one surface only and ends at an inner layer."},
	{KP: 1, Q: "Which via joins only inner layers and is invisible from both surfaces?",
		A: "A buried via", D: []string{"A blind via", "A through via", "A microvia on the top layer"},
		Why: "A buried via is drilled and plated in an inner sub-stack before the outer layers are laminated over it."},
	{KP: 1, Q: "What defines a microvia in current IPC usage?",
		A:   "A small laser-drilled blind structure with an aspect ratio of at most 1:1, usually spanning one dielectric",
		D:   []string{"Any via narrower than the trace it serves", "A mechanically drilled through via under 0.3 mm", "A buried via between two cores"},
		Why: "IPC defines microvias by aspect ratio and depth rather than by diameter. Laser drilling forms them one dielectric layer at a time."},
	{KP: 1, Q: "Why is a through via the cheapest kind of via?",
		A:   "It is drilled once through the finished laminated stack and plated in one step",
		D:   []string{"It contains no copper", "It needs no drilling", "It is drilled by laser"},
		Why: "Every through via on the board shares one drill program and one plating step after the single lamination."},
	{KP: 1, Q: "A buried via joins L3 to L6 of a 10-layer board. What does fabrication have to do?",
		A:   "Drill and plate the L3 to L6 sub-stack before the final lamination",
		D:   []string{"Laser-drill from the top surface after lamination", "Backdrill from the bottom surface", "Nothing special, since it is a through via"},
		Why: "Each buried span is a separate sub-assembly that is drilled, plated and then laminated into the stack, adding process steps and cost."},
	{KP: 2, Q: "A 0.5 mm pitch BGA leaves no room for dog-bone fanout between its pads. Which via structure is typical?",
		A:   "Via-in-pad: microvias or filled and plated-over vias placed in the BGA lands",
		D:   []string{"Larger through vias between the pads", "Buried vias only, with no surface connection", "Thermal vias around the package"},
		Why: "At fine pitch the via must sit in the land itself. It is filled and capped so solder does not wick into the barrel."},
	{KP: 2, Q: "Why are blind and buried vias valued for dense routing?",
		A:   "They occupy only the layers they connect, leaving the other layers free to route over or under them",
		D:   []string{"They lower the board cost", "They are larger and carry more current", "They remove the need for planes"},
		Why: "A through via blocks routing on every layer. A blind or buried via frees the layers it does not span."},
	{KP: 2, Q: "A through via passes a plane it must not connect to. What does that plane need around the barrel?",
		A: "An antipad, a clearance opening", D: []string{"A thermal relief", "A teardrop", "A solder mask opening"},
		Why: "The antipad isolates the barrel from the plane. Its size also sets the via's capacitance and hence its impedance."},
	{KP: 2, Q: "What is the main trade-off between stacked and staggered microvias?",
		A:   "Stacked saves area but needs copper-filled vias and carries more reliability risk; staggered uses more area",
		D:   []string{"Staggered microvias need no laser drilling", "Stacked microvias are cheaper in every case", "Staggered microvias cannot carry signals"},
		Why: "Stacking puts each microvia on the filled one below it, which is compact but concentrates thermal-cycling stress at the interfaces."},
	{KP: 2, Q: "A through via carries a signal from L1 to L3 of a 12-layer board. What is left over?",
		A:   "A stub from L3 to the bottom that behaves like an open-ended transmission line",
		D:   []string{"Nothing, because the whole barrel carries current", "A short to every plane it passes", "A stub above L1"},
		Why: "The unused barrel below L3 is a capacitive open stub. Its quarter-wave resonance notches the channel."},
	{KP: 2, Q: "How does backdrilling remove a via stub?",
		A:   "A slightly larger drill removes the unused plated barrel from the far side to a controlled depth",
		D:   []string{"A laser vaporizes the pad on the far side", "The stub is filled with epoxy", "The barrel is plated thinner"},
		Why: "Controlled-depth drilling cuts away the copper below the last used layer, leaving a short residual stub."},
}

// stubResonance is the computed KP 3 of via.types: the quarter-wave resonance of a via stub.
func stubResonance(ctx generator.GeneratorContext) generator.Problem {
	er := edakit.Pick([]float64{3.4, 3.6, 3.8, 4.0, 4.2})
	L := float64(edakit.Int(3, 24) * 5)
	f := 11.8 / (4 * L / 1000 * sqrt(er))
	ans, ds := options(f, 1, "GHz", 2*f, f/2, f*sqrt(er), f/sqrt(er))
	q := edakit.Sprintf("A via stub is an open-ended line whose first resonance (a notch in the channel) falls where the stub is a quarter wavelength long: $f_0 = \\frac{c}{4 L \\sqrt{\\varepsilon_r}}$, with $c$ = 11.8 inches per ns. A through via leaves a %s mil stub in a dielectric with $\\varepsilon_r = %s$. What is $f_0$?", edakit.Num(L, 0), edakit.Num(er, 1))
	return edakit.Choice(q, ans, edakit.Sprintf("$f_0 = 11.8 / (4 \\cdot %s \\cdot \\sqrt{%s})$ = %s, with the length in inches. A shorter stub, from backdrilling or a blind via, pushes the notch above the signal's band.", edakit.Num(L/1000, 3), edakit.Num(er, 1), ans), ds...)
}

// aspectRatio: board aspect ratio (KP 1), minimum drill for a fab limit (KP 2), microvia spans (KP 3).
func aspectRatio(ctx generator.GeneratorContext) generator.Problem {
	T := float64(edakit.Pick([]int{31, 47, 62, 62, 93, 125}))
	switch kpOf(ctx) {
	case 1:
		d := float64(edakit.Int(6, 20))
		ans := edakit.Num(T/d, 2)
		q := edakit.Sprintf("A %s mil thick board has through vias drilled at %s mil. The aspect ratio is board thickness divided by drilled hole diameter, written $x$:1. What is $x$? Round to 2 decimal places.", edakit.Num(T, 0), edakit.Num(d, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("$%s / %s = %s$. High ratios make it hard for plating solution to reach the middle of the barrel, so the copper there thins.", edakit.Num(T, 0), edakit.Num(d, 0), ans), map[string]string{"ar": ans})
	case 2:
		ar := float64(edakit.Pick([]int{6, 8, 10, 12}))
		T = float64(edakit.Pick([]int{62, 93, 125, 157}))
		d := ceilInt(T / ar)
		q := edakit.Sprintf("A fabricator plates through holes reliably up to an aspect ratio (board thickness to drilled diameter) of %s:1. The board is %s mil thick. What is the smallest drill diameter allowed, in whole mil? Round up to a whole mil.", edakit.Num(ar, 0), edakit.Num(T, 0))
		return edakit.Numeric(q, edakit.IntStr(d), edakit.Sprintf("$%s / %s = %s$ mil, rounded up to %d mil. A thicker board forces larger vias, which is one reason dense designs move to thinner boards or microvias.", edakit.Num(T, 0), edakit.Num(ar, 0), edakit.Num(T/ar, 3), d), map[string]string{"dmin": edakit.IntStr(d)})
	}
	dia := float64(edakit.Int(3, 6))
	lim := edakit.Pick([]float64{0.75, 0.8, 1.0})
	n := edakit.Scaled(ctx, 4, 7)
	ok := 0
	var ts []string
	for i := 0; i < n; i++ {
		t := float64(edakit.Int(4, 14)) / 2
		ts = append(ts, edakit.Num(t, 1))
		if t/dia <= lim+1e-9 {
			ok++
		}
	}
	q := edakit.Sprintf("A laser drills %s mil microvias, and the fabricator caps microvia aspect ratio (dielectric depth to diameter) at %s:1. Candidate build-up dielectrics are %s mil thick. How many of these candidates can a single microvia span within the limit?", edakit.Num(dia, 0), edakit.Num(lim, 2), joinAnd(ts))
	return edakit.Numeric(q, edakit.IntStr(ok), edakit.Sprintf("A dielectric qualifies when its thickness is at most %s x %s = %s mil. %d of the candidates do.", edakit.Num(lim, 2), edakit.Num(dia, 0), edakit.Num(lim*dia, 2), ok), map[string]string{"max_depth": edakit.Num(lim*dia, 2)})
}

// hdi works with i+N+i build-up notation.
func hdi(ctx generator.GeneratorContext) generator.Problem {
	i := edakit.Int(1, 3)
	N := edakit.Pick([]int{4, 6, 8, 10})
	L := 2*i + N
	name := edakit.Sprintf("%d+%d+%d", i, N, i)
	switch kpOf(ctx) {
	case 1:
		if edakit.Int(1, 2) == 1 {
			q := edakit.Sprintf("HDI build-ups are written $i + N + i$: an $N$-layer core with $i$ sequential build-up layers on each side. How many copper layers does a %s board have?", name)
			return edakit.Numeric(q, edakit.IntStr(L), edakit.Sprintf("$2 \\cdot %d + %d = %d$ layers.", i, N, L), nil)
		}
		q := edakit.Sprintf("HDI build-ups are written $i + N + i$: an $N$-layer core with $i$ sequential build-up layers on each side. %s %d-layer board uses %s on each side. How many layers does its core have?", strings.ToUpper(an(L)[:1])+an(L)[1:], L, plural(i, "build-up layer", "build-up layers"))
		return edakit.Numeric(q, edakit.IntStr(N), edakit.Sprintf("$%d - 2 \\cdot %d = %d$ core layers, so the board is %s.", L, i, N, name), nil)
	case 2:
		if edakit.Int(1, 2) == 1 {
			q := edakit.Sprintf("In a sequential build-up, the %d-layer core of a %s board is laminated once, and each build-up layer (added to both sides at the same time) needs one more lamination cycle. How many lamination cycles does the board need?", N, name)
			return edakit.Numeric(q, edakit.IntStr(i+1), edakit.Sprintf("One for the core plus %d for the build-up layers: %d cycles. Each cycle adds press time, registration error and cost, which is why HDI price climbs with $i$.", i, i+1), nil)
		}
		q := edakit.Sprintf("A %s board is made by sequential build-up. The core is mechanically drilled once for buried vias, each build-up layer on each side is laser-drilled once for its microvias, and through holes are mechanically drilled once after the last lamination. How many drilling operations are there in total?", name)
		return edakit.Numeric(q, edakit.IntStr(2*i+2), edakit.Sprintf("Core drill 1, laser drills $2 \\cdot %d = %d$, final drill 1: %d operations.", i, 2*i, 2*i+2), nil)
	}
	m := edakit.Int(2, L)
	micro, buried := 0, 0
	coreTop, coreBot := i+1, i+N
	switch {
	case m <= coreTop:
		micro = m - 1
	case m <= coreBot:
		micro, buried = i, 1
	default:
		micro, buried = i+(m-coreBot), 1
	}
	total := micro + buried
	span := func(a, b int) string {
		if a == b {
			return edakit.Sprintf("L%d", a)
		}
		return edakit.Sprintf("L%d to L%d", a, b)
	}
	layout := edakit.Sprintf("top build-up %s, core %s, bottom build-up %s", span(1, i), span(coreTop, coreBot), span(coreBot+1, L))
	q := edakit.Sprintf("On a %s board (%s), each stacked microvia spans one build-up dielectric, and the core is crossed only by a buried via spanning the whole core, which also connects the core's inner layers. A signal must go from L1 to L%d with stacked microvias. How many via structures (microvias plus buried vias) does the connection use?", name, layout, m)
	return edakit.Numeric(q, edakit.IntStr(total), edakit.Sprintf("%s and %s: %d structures. Every structure in the stack is another interface that must survive thermal cycling.", plural(micro, "microvia", "microvias"), plural(buried, "buried via", "buried vias"), total), map[string]string{"micro": edakit.IntStr(micro), "buried": edakit.IntStr(buried)})
}
