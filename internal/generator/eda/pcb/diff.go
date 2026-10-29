package pcb

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var oddEvenItems = []edakit.Item{
	{KP: 1, Q: "Two coupled traces are driven with equal and opposite voltages. Which mode is excited?",
		A: "The odd mode", D: []string{"The even mode", "The common mode", "No mode, since the voltages cancel"},
		Why: "Opposite polarity on the two lines is the odd mode; equal polarity is the even mode. Any drive is a sum of the two."},
	{KP: 1, Q: "Two traces of a pair are brought closer together. What happens to the odd-mode impedance?",
		A: "It falls below the impedance of one uncoupled trace", D: []string{"It rises above the uncoupled impedance", "It stays equal to the uncoupled impedance", "It drops to zero"},
		Why: "In the odd mode the neighbour is at the opposite voltage, so mutual capacitance adds to each line's capacitance and the impedance falls."},
	{KP: 1, Q: "Two traces of a pair are brought closer together. What happens to the even-mode impedance?",
		A: "It rises above the impedance of one uncoupled trace", D: []string{"It falls below the uncoupled impedance", "It stays equal to the uncoupled impedance", "It becomes infinite"},
		Why: "In the even mode both lines swing together, so the mutual capacitance carries no charge and the mutual inductance adds to each line's inductance."},
	{KP: 1, Q: "Why does odd-mode impedance fall as coupling tightens?",
		A:   "Mutual capacitance adds to each line's capacitance and mutual inductance subtracts from its inductance",
		D:   []string{"Mutual capacitance subtracts and mutual inductance adds", "The dielectric constant rises between close traces", "Copper resistance falls between close traces"},
		Why: "With opposite currents, $L_{odd} = L - L_m$ and $C_{odd} = C + C_m$, so $\\sqrt{L/C}$ drops."},
	{KP: 2, Q: "How does differential impedance relate to odd-mode impedance?",
		A: "$Z_{diff} = 2 Z_{odd}$", D: []string{"$Z_{diff} = Z_{odd} / 2$", "$Z_{diff} = Z_{odd}$", "$Z_{diff} = \\sqrt{2}\\, Z_{odd}$"},
		Why: "The differential signal sees the two odd-mode impedances in series, from one line to the other."},
	{KP: 2, Q: "How does common-mode impedance relate to even-mode impedance?",
		A: "$Z_{cm} = Z_{even} / 2$", D: []string{"$Z_{cm} = 2 Z_{even}$", "$Z_{cm} = Z_{even}$", "$Z_{cm} = Z_{odd} / 2$"},
		Why: "A common-mode signal drives both lines in parallel against the reference, so the two even-mode impedances combine in parallel."},
	{KP: 2, Q: "A receiver places one 100 ohm resistor across a pair whose differential impedance is 100 ohms. Which mode does that terminate?",
		A:   "Only the odd (differential) mode; the even mode is left unterminated",
		D:   []string{"Both modes", "Only the even mode", "Neither mode"},
		Why: "Common-mode current flows through both lines in the same direction and does not pass through a resistor connected across them."},
	{KP: 2, Q: "Which network terminates both modes of a pair?",
		A:   "A T network: $Z_{odd}$ from each line to a center node and $(Z_{even} - Z_{odd})/2$ from that node to ground",
		D:   []string{"One resistor equal to $Z_{even}$ across the lines", "One $Z_{odd}$ resistor from each line straight to ground", "One $2 Z_{even}$ resistor from the lines to ground"},
		Why: "In the odd mode the center node is a virtual ground, so each line sees $Z_{odd}$. In the even mode each line sees $Z_{odd} + 2 \\cdot (Z_{even} - Z_{odd})/2 = Z_{even}$."},
	{KP: 3, Q: "In surface microstrip, which mode travels faster?",
		A:   "The odd mode, because more of its field lies in the air above the board",
		D:   []string{"The even mode, because its field stays in the laminate", "Both modes travel at the same speed", "It depends only on the trace length"},
		Why: "Odd-mode fields fringe strongly between the traces through air, lowering its effective dielectric constant."},
	{KP: 3, Q: "Why is far-end crosstalk nearly zero between striplines in a uniform dielectric?",
		A:   "Both modes see the same dielectric and travel at one speed, so the capacitive and inductive far-end terms cancel",
		D:   []string{"Striplines have no mutual capacitance", "The planes absorb all crosstalk", "Far-end crosstalk only exists on flex"},
		Why: "Far-end crosstalk grows with the difference between odd- and even-mode velocities, which is zero in a homogeneous medium."},
	{KP: 3, Q: "A pair targeting 100 ohms differential is routed loosely coupled with a wide gap. What follows?",
		A:   "Each trace is designed near 50 ohms single-ended, and the pair is less sensitive to gap variation",
		D:   []string{"Each trace must be 100 ohms single-ended", "The pair no longer rejects common-mode noise at the receiver", "Loose coupling makes the differential impedance zero"},
		Why: "With weak coupling $Z_{odd} \\approx Z_0$, so $Z_{diff} \\approx 2 Z_0$. Tight coupling saves area but makes impedance depend strongly on spacing."},
	{KP: 3, Q: "A tightly coupled pair crosses a split in its reference plane. What happens?",
		A:   "Part of the differential signal converts to common mode, which radiates and closes the eye",
		D:   []string{"Nothing, since the pair carries its own return", "The pair becomes purely odd mode", "Its impedance stays exactly constant"},
		Why: "Much of each trace's field still terminates on the plane. Losing the plane breaks the symmetry between the lines and couples the modes."},
	{KP: 3, Q: "Does a board differential pair need a reference plane?",
		A:   "In practice yes: most of each trace's field couples to the plane, which sets the impedance and returns common-mode current",
		D:   []string{"No, the return current flows only in the partner trace", "Only when the pair is longer than one meter", "No, planes only matter for power nets"},
		Why: "Board pairs are usually edge-coupled over a plane, so the plane is part of the transmission-line structure."},
}

// diffPairs: Zdiff from Zodd (KP 1), Zcm from Zeven (KP 2), coupling from measured Zdiff and Zcm (KP 3).
func diffPairs(ctx generator.GeneratorContext) generator.Problem {
	zo := float64(edakit.Int(380, 560)) / 10
	ze := zo + float64(edakit.Int(30, 250))/10
	switch kpOf(ctx) {
	case 1:
		if edakit.Int(1, 2) == 1 {
			ans := edakit.Num(2*zo, 1)
			q := edakit.Sprintf("A field solver reports an edge-coupled pair with odd-mode impedance %s ohms and even-mode impedance %s ohms. What is the differential impedance, in ohms? Round to 1 decimal place.", edakit.Num(zo, 1), edakit.Num(ze, 1))
			return edakit.Numeric(q, ans, edakit.Sprintf("$Z_{diff} = 2 Z_{odd} = 2 \\cdot %s = %s$ ohms. The even-mode value does not enter.", edakit.Num(zo, 1), ans), nil)
		}
		zd := float64(edakit.Pick([]int{85, 90, 100, 100, 110, 120}))
		ans := edakit.Num(zd/2, 1)
		q := edakit.Sprintf("An interface specifies %s ohms differential impedance. What odd-mode impedance must each trace of the pair present, in ohms? Round to 1 decimal place.", edakit.Num(zd, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("$Z_{odd} = Z_{diff}/2 = %s$ ohms.", ans), nil)
	case 2:
		ans := edakit.Num(ze/2, 2)
		q := edakit.Sprintf("A pair has odd-mode impedance %s ohms and even-mode impedance %s ohms. What is its common-mode impedance, in ohms? Round to 2 decimal places.", edakit.Num(zo, 1), edakit.Num(ze, 1))
		return edakit.Numeric(q, ans, edakit.Sprintf("Common-mode drive puts the two even-mode impedances in parallel: $Z_{cm} = Z_{even}/2 = %s$ ohms.", ans), nil)
	}
	zd := 2 * zo
	zc := ze / 2
	k := (ze - zo) / (ze + zo)
	ans := edakit.Num(k, 3)
	q := edakit.Sprintf("A TDR measures a pair at %s ohms differential and %s ohms common mode. The coupling coefficient is $k = (Z_{even} - Z_{odd}) / (Z_{even} + Z_{odd})$. What is $k$? Round to 3 decimal places.", edakit.Num(zd, 1), edakit.Num(zc, 2))
	return edakit.Numeric(q, ans, edakit.Sprintf("$Z_{odd} = %s/2 = %s$ and $Z_{even} = 2 \\cdot %s = %s$ ohms, so $k = %s / %s = %s$.", edakit.Num(zd, 1), edakit.Num(zo, 1), edakit.Num(zc, 2), edakit.Num(ze, 1), edakit.Num(ze-zo, 1), edakit.Num(ze+zo, 1), ans), map[string]string{"zodd": edakit.Num(zo, 1), "zeven": edakit.Num(ze, 1)})
}

// lengthMatching: skew from mismatch (KP 1), tolerance from a budget (KP 2), a byte-lane window (KP 3).
func lengthMatching(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		dl := float64(edakit.Int(5, 200))
		tpd := float64(edakit.Int(140, 185))
		ans := edakit.Num(dl/1000*tpd, 2)
		q := edakit.Sprintf("The two traces of a differential pair differ in length by %s mil, and the layer's delay is %s ps per inch. What is the intra-pair skew, in ps? Round to 2 decimal places.", edakit.Num(dl, 0), edakit.Num(tpd, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("%s mil is %s in, times %s ps per inch is %s ps. Skew converts part of the differential signal into common mode.", edakit.Num(dl, 0), edakit.Num(dl/1000, 3), edakit.Num(tpd, 0), ans), nil)
	case 2:
		er := edakit.Pick([]float64{3.2, 3.5, 3.7, 3.9, 4.1, 4.3})
		budget := float64(edakit.Int(2, 20)) / 2
		tpd := 1000 * sqrt(er) / 11.8
		tol := floorInt(budget / tpd * 1000)
		q := edakit.Sprintf("A link allows %s ps of skew between two signals routed as stripline in a laminate with $\\varepsilon_r = %s$. The delay is $t_{pd} = \\sqrt{\\varepsilon_r}/11.8$ ns per inch. What is the largest length mismatch allowed, in whole mil? Round down to a whole mil.", edakit.Num(budget, 1), edakit.Num(er, 1))
		return edakit.Numeric(q, edakit.IntStr(tol), edakit.Sprintf("$t_{pd} = %s$ ps per inch, so the mismatch may be $%s / %s = %s$ in, which rounds down to %d mil.", edakit.Num(tpd, 2), edakit.Num(budget, 1), edakit.Num(tpd, 2), edakit.Num(budget/tpd, 5), tol), map[string]string{"tpd": edakit.Num(tpd, 2)})
	}
	tpd := float64(edakit.Int(140, 180))
	win := float64(edakit.Pick([]int{2, 3, 5, 8, 10}))
	ls := edakit.Int(1500, 3500)
	n := edakit.Scaled(ctx, 5, 9)
	out := 0
	var sig []string
	for i := 0; i < n; i++ {
		l := ls + edakit.Int(-90, 90)
		sk := absf(float64(l-ls)) / 1000 * tpd
		if sk > win+1e-9 {
			out++
		}
		sig = append(sig, edakit.Sprintf("DQ%d %d mil", i, l))
	}
	q := edakit.Sprintf("In a memory byte lane every DQ must arrive within $\\pm %s$ ps of its strobe. The strobe pair is %d mil long, the layer delay is %s ps per inch, and the DQ lengths are %s. How many DQ signals fall outside the window and must be tuned?", edakit.Num(win, 0), ls, edakit.Num(tpd, 0), strings.Join(sig, ", "))
	return edakit.Numeric(q, edakit.IntStr(out), edakit.Sprintf("The window is $\\pm %s / %s$ in, about $\\pm %s$ mil. %s from the strobe by more than that.", edakit.Num(win, 0), edakit.Num(tpd, 0), edakit.Num(win/tpd*1000, 1), plural(out, "signal differs", "signals differ")), map[string]string{"window_mil": edakit.Num(win/tpd*1000, 2)})
}
