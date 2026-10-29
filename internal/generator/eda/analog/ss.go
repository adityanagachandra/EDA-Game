package analog

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// smallSignal: gm from Id and Vov, ro from lambda, gm from the square law.
func smallSignal(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		id := float64(10 * edakit.Int(2, 80))
		vov := unif(0.08, 0.40, 0.01)
		gm := 2 * id / vov
		q := edakit.Sprintf("An NMOS transistor in saturation carries $I_D = %s\\ \\mu$A at an overdrive voltage $V_{OV} = %s$ V. Using the square-law model, what is its transconductance $g_m$ in microsiemens? Round to the nearest integer.", f(id, 0), f(vov, 2))
		return num(q, f(gm, 0),
			edakit.Sprintf("In the square law $I_D = \\frac{k}{2}V_{OV}^2$, so $g_m = \\partial I_D/\\partial V_{GS} = kV_{OV} = 2I_D/V_{OV} = 2 \\cdot %s / %s = %s\\ \\mu$S.", f(id, 0), f(vov, 2), f(gm, 0)),
			map[string]string{"gm_uS": f(gm, 0)})
	case 2:
		id := float64(10 * edakit.Int(2, 60))
		lam := unif(0.04, 0.30, 0.01)
		ro := 1 / (lam * id * 1e-6) / 1000
		q := edakit.Sprintf("A saturated MOSFET has channel-length modulation $\\lambda = %s\\ \\mathrm{V}^{-1}$ and carries $I_D = %s\\ \\mu$A. What is its small-signal output resistance $r_o$ in kilohms? Round to 1 decimal place.", f(lam, 2), f(id, 0))
		return num(q, f(ro, 1),
			edakit.Sprintf("$r_o = 1/(\\lambda I_D) = 1/(%s \\cdot %s \\times 10^{-6}) = %s$ kilohms.", f(lam, 2), f(id, 0), f(ro, 1)),
			map[string]string{"ro_k": f(ro, 1)})
	default:
		k := float64(10 * edakit.Int(10, 40))
		wl := float64(edakit.Int(2, 40))
		id := float64(10 * edakit.Int(2, 50))
		gm := math.Sqrt(2 * k * wl * id)
		q := edakit.Sprintf("An NMOS device has $\\mu_n C_{ox} = %s\\ \\mu\\mathrm{A/V}^2$ and $W/L = %s$, biased in saturation at $I_D = %s\\ \\mu$A. Using $g_m = \\sqrt{2\\mu_n C_{ox}(W/L)I_D}$, what is $g_m$ in microsiemens? Round to 1 decimal place.", f(k, 0), f(wl, 0), f(id, 0))
		return num(q, f(gm, 1),
			edakit.Sprintf("$g_m = \\sqrt{2 \\cdot %s \\cdot %s \\cdot %s} = %s\\ \\mu$S. Doubling $W/L$ at fixed current raises $g_m$ only by $\\sqrt{2}$.", f(k, 0), f(wl, 0), f(id, 0), f(gm, 1)),
			map[string]string{"gm_uS": f(gm, 1)})
	}
}

func dbConversions(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		if edakit.Int(0, 1) == 0 {
			a := float64(edakit.Int(3, 9000))
			g := db20(a)
			q := edakit.Sprintf("An amplifier has a voltage gain of %s V/V. What is the gain in decibels? Round to 2 decimal places.", f(a, 0))
			return num(q, f(g, 2), edakit.Sprintf("Voltage ratios use $20\\log_{10}$: $20\\log_{10}(%s) = %s$ dB.", f(a, 0), f(g, 2)), map[string]string{"db": f(g, 2)})
		}
		p := float64(edakit.Int(2, 5000))
		g := db10(p)
		q := edakit.Sprintf("A stage delivers %s times more power to its load than it draws from its source. What is its power gain in decibels? Round to 2 decimal places.", f(p, 0))
		return num(q, f(g, 2), edakit.Sprintf("Power ratios use $10\\log_{10}$: $10\\log_{10}(%s) = %s$ dB.", f(p, 0), f(g, 2)), map[string]string{"db": f(g, 2)})
	case 2:
		if edakit.Int(0, 1) == 0 {
			p := unif(0.01, 200, 0.01)
			d := db10(p / 1)
			q := edakit.Sprintf("A signal generator outputs %s mW. What is this power in dBm? Round to 2 decimal places.", f(p, 2))
			return num(q, f(d, 2), edakit.Sprintf("dBm is power relative to 1 mW: $10\\log_{10}(%s / 1) = %s$ dBm.", f(p, 2), f(d, 2)), map[string]string{"dbm": f(d, 2)})
		}
		d := unif(-40, 20, 0.5)
		p := undb10(d)
		q := edakit.Sprintf("A receiver input sees %s dBm. What is this power in milliwatts? Round to 4 decimal places.", f(d, 1))
		return num(q, f(p, 4), edakit.Sprintf("$P = 10^{%s/10}$ mW $= %s$ mW.", f(d, 1), f(p, 4)), map[string]string{"mw": f(p, 4)})
	default:
		g1 := unif(6, 30, 0.5)
		a2 := float64(edakit.Int(2, 40))
		att := unif(1, 10, 0.5)
		tot := g1 + db20(a2) - att
		q := edakit.Sprintf("A chain has a %s dB voltage-gain stage, then a stage with voltage gain %s V/V, then a %s dB attenuator. What is the total voltage gain in dB? Round to 2 decimal places.", f(g1, 1), f(a2, 0), f(att, 1))
		return num(q, f(tot, 2), edakit.Sprintf("Convert the middle stage, $20\\log_{10}(%s) = %s$ dB, then add the cascade in dB: $%s + %s - %s = %s$ dB.", f(a2, 0), f(db20(a2), 2), f(g1, 1), f(db20(a2), 2), f(att, 1), f(tot, 2)), map[string]string{"db": f(tot, 2)})
	}
}

func intrinsicGain(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		gm := unif(0.2, 5, 0.1)
		ro := float64(edakit.Int(10, 200))
		a := gm * ro
		q := edakit.Sprintf("A transistor has $g_m = %s$ mS and $r_o = %s$ kilohms. What is its intrinsic gain $g_m r_o$ in V/V? Round to 1 decimal place.", f(gm, 1), f(ro, 0))
		return num(q, f(a, 1), edakit.Sprintf("Millisiemens times kilohms is dimensionless: $%s \\cdot %s = %s$ V/V.", f(gm, 1), f(ro, 0), f(a, 1)), map[string]string{"gain": f(a, 1)})
	case 2:
		lam := unif(0.05, 0.4, 0.01)
		vov := unif(0.1, 0.4, 0.01)
		a := 2 / (lam * vov)
		q := edakit.Sprintf("In the square-law model with channel-length modulation $\\lambda = %s\\ \\mathrm{V}^{-1}$, a device runs at $V_{OV} = %s$ V. What is its intrinsic gain in V/V? Round to 1 decimal place.", f(lam, 2), f(vov, 2))
		return num(q, f(a, 1), edakit.Sprintf("$g_m = 2I_D/V_{OV}$ and $r_o = 1/(\\lambda I_D)$, so $g_m r_o = 2/(\\lambda V_{OV}) = 2/(%s \\cdot %s) = %s$. The current cancels.", f(lam, 2), f(vov, 2), f(a, 1)), map[string]string{"gain": f(a, 1)})
	default:
		gmid := float64(edakit.Int(5, 24))
		va := unif(2, 20, 0.5)
		a := gmid * va
		q := edakit.Sprintf("A device is biased at $g_m/I_D = %s\\ \\mathrm{V}^{-1}$ and has an Early voltage $V_A = 1/\\lambda = %s$ V at its channel length. What is its intrinsic gain in dB? Round to 1 decimal place.", f(gmid, 0), f(va, 1))
		return num(q, f(db20(a), 1), edakit.Sprintf("$g_m r_o = (g_m/I_D)(I_D r_o) = (g_m/I_D) V_A = %s \\cdot %s = %s$ V/V, which is %s dB.", f(gmid, 0), f(va, 1), f(a, 1), f(db20(a), 1)), map[string]string{"gain": f(a, 1)})
	}
}

func gmID(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		gm := float64(50 * edakit.Int(4, 80))
		gmid := float64(edakit.Int(5, 25))
		id := gm / gmid
		q := edakit.Sprintf("A gain stage needs $g_m = %s\\ \\mu$S. The designer picks $g_m/I_D = %s\\ \\mathrm{V}^{-1}$. What drain current does the device need, in microamperes? Round to 1 decimal place.", f(gm, 0), f(gmid, 0))
		return num(q, f(id, 1), edakit.Sprintf("$I_D = g_m/(g_m/I_D) = %s / %s = %s\\ \\mu$A. A higher $g_m/I_D$ buys the same $g_m$ for less current.", f(gm, 0), f(gmid, 0), f(id, 1)), map[string]string{"id_uA": f(id, 1)})
	case 2:
		if edakit.Int(0, 1) == 0 {
			gmid := float64(edakit.Int(4, 16))
			vov := 2 / gmid * 1000
			q := edakit.Sprintf("In strong inversion the square law gives $g_m/I_D = 2/V_{OV}$. What overdrive voltage corresponds to $g_m/I_D = %s\\ \\mathrm{V}^{-1}$, in millivolts? Round to 1 decimal place.", f(gmid, 0))
			return num(q, f(vov, 1), edakit.Sprintf("$V_{OV} = 2/(g_m/I_D) = 2/%s$ V $= %s$ mV.", f(gmid, 0), f(vov, 1)), map[string]string{"vov_mV": f(vov, 1)})
		}
		n := unif(1.1, 1.6, 0.05)
		t := float64(edakit.Int(250, 400))
		ut := kBoltzmann * t / 1.602176634e-19
		mx := 1 / (n * ut)
		q := edakit.Sprintf("Deep in weak inversion, $g_m/I_D$ approaches $1/(n U_T)$ with $U_T = kT/q$. For a subthreshold slope factor $n = %s$ at $T = %s$ K, what is this ceiling on $g_m/I_D$ in $\\mathrm{V}^{-1}$? Round to 2 decimal places.", f(n, 2), f(t, 0))
		return num(q, f(mx, 2), edakit.Sprintf("$U_T = kT/q = %s$ mV, so $1/(nU_T) = 1/(%s \\cdot %s\\ \\mathrm{V}) = %s\\ \\mathrm{V}^{-1}$. No bias point does better.", f(ut*1000, 2), f(n, 2), f(ut, 5), f(mx, 2)), map[string]string{"max": f(mx, 2)})
	default:
		gmid := float64(edakit.Int(5, 22))
		id := float64(5 * edakit.Int(2, 60))
		cgg := float64(edakit.Int(5, 200))
		gm := gmid * id
		ft := gm / (2 * math.Pi * cgg)
		q := edakit.Sprintf("A device runs at $g_m/I_D = %s\\ \\mathrm{V}^{-1}$ with $I_D = %s\\ \\mu$A and total gate capacitance $C_{gg} = %s$ fF. What is its transit frequency $f_T \\approx g_m/(2\\pi C_{gg})$ in GHz? Round to 2 decimal places.", f(gmid, 0), f(id, 0), f(cgg, 0))
		return num(q, f(ft, 2), edakit.Sprintf("$g_m = %s \\cdot %s = %s\\ \\mu$S. One microsiemens over one femtofarad is $10^9$ per second, so $f_T = %s/(2\\pi \\cdot %s) = %s$ GHz.", f(gmid, 0), f(id, 0), f(gm, 0), f(gm, 0), f(cgg, 0), f(ft, 2)), map[string]string{"ft": f(ft, 2)})
	}
}

// lookupTables draws a gm/Id table at fixed L and asks for widths.
func lookupTables(ctx generator.GeneratorContext) generator.Problem {
	gmids := []float64{6, 10, 14, 18, 22}
	dens := make([]float64, len(gmids))
	gains := make([]float64, len(gmids))
	d := unif(40, 90, 1)
	g := unif(15, 30, 1)
	for i := range gmids {
		dens[i] = edakit.Round(d, 2)
		gains[i] = g
		d *= unif(0.35, 0.6, 0.01)
		g += unif(2, 8, 1)
	}
	var rows []string
	for i := range gmids {
		rows = append(rows, edakit.Sprintf("$g_m/I_D = %s$: $I_D/W = %s\\ \\mu$A/$\\mu$m, $g_m/g_{ds} = %s$", f(gmids[i], 0), f(dens[i], 2), f(gains[i], 0)))
	}
	table := "A gm/Id lookup table, simulated with the foundry model at the chosen channel length, gives: " + strings.Join(rows, "; ") + "."
	switch whichKP(ctx) {
	case 1:
		i := edakit.Int(0, len(gmids)-1)
		id := float64(5 * edakit.Int(2, 80))
		w := id / dens[i]
		q := table + edakit.Sprintf(" A device must carry $I_D = %s\\ \\mu$A at $g_m/I_D = %s\\ \\mathrm{V}^{-1}$. What width does it need, in micrometres? Round to 2 decimal places.", f(id, 0), f(gmids[i], 0))
		return num(q, f(w, 2), edakit.Sprintf("Width scales current at fixed current density: $W = I_D/(I_D/W) = %s/%s = %s\\ \\mu$m.", f(id, 0), f(dens[i], 2), f(w, 2)), map[string]string{"w": f(w, 2)})
	case 2:
		i := edakit.Int(0, len(gmids)-2)
		x := gmids[i] + float64(edakit.Int(1, 3))
		t := (x - gmids[i]) / (gmids[i+1] - gmids[i])
		dx := dens[i] + t*(dens[i+1]-dens[i])
		q := table + edakit.Sprintf(" Interpolating linearly between table rows, what current density $I_D/W$ corresponds to $g_m/I_D = %s\\ \\mathrm{V}^{-1}$, in $\\mu$A/$\\mu$m? Round to 3 decimal places.", f(x, 0))
		return num(q, f(dx, 3), edakit.Sprintf("The point lies a fraction %s of the way from %s to %s, so $I_D/W = %s + %s(%s - %s) = %s$. Real tools interpolate in log of density, since density falls roughly exponentially with $g_m/I_D$.", f(t, 2), f(gmids[i], 0), f(gmids[i+1], 0), f(dens[i], 2), f(t, 2), f(dens[i+1], 2), f(dens[i], 2), f(dx, 3)), map[string]string{"density": f(dx, 3)})
	default:
		i := edakit.Int(0, len(gmids)-1)
		gm := float64(50 * edakit.Int(4, 60))
		id := gm / gmids[i]
		w := id / dens[i]
		q := table + edakit.Sprintf(" A stage needs $g_m = %s\\ \\mu$S and the designer picks the row $g_m/I_D = %s\\ \\mathrm{V}^{-1}$. What width results, in micrometres? Round to 2 decimal places.", f(gm, 0), f(gmids[i], 0))
		return num(q, f(w, 2), edakit.Sprintf("First $I_D = g_m/(g_m/I_D) = %s/%s = %s\\ \\mu$A, then $W = I_D/(I_D/W) = %s/%s = %s\\ \\mu$m. The same row also gives the intrinsic gain, %s, with no extra simulation.", f(gm, 0), f(gmids[i], 0), f(id, 2), f(id, 2), f(dens[i], 2), f(w, 2), f(gains[i], 0)), map[string]string{"id": f(id, 2), "w": f(w, 2)})
	}
}
