package pkg

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// sipMCM: compound module yield, known-good-die yield and cost per good module.
func sipMCM(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		k := edakit.Scaled(ctx, 2, 5)
		ya := step(0.97, 0.999, 0.001)
		Y := math.Pow(ya, float64(k))
		var parts []string
		for i := 0; i < k; i++ {
			y := step(0.85, 0.99, 0.01)
			Y *= y
			parts = append(parts, f(y, 2))
		}
		ans := f(Y*100, 2)
		q := edakit.Sprintf("A multi-chip module holds %d untested dies with yields %s. Each die attach succeeds with probability %s, independently. What is the module yield, as a percentage? Round to 2 decimal places.", k, strings.Join(parts, ", "), f(ya, 3))
		return num(q, ans, edakit.Sprintf("$Y = %s^{%d} \\times %s = %s\\%%$. Every part must be good, so yields multiply.", f(ya, 3), k, strings.Join(parts, " \\times "), ans), "yield_pct")
	case 2:
		k := edakit.Scaled(ctx, 2, 8)
		e := step(0.5, 5, 0.5)
		ya := step(0.98, 0.999, 0.001)
		Y := math.Pow((1-e/100)*ya, float64(k))
		ans := f(Y*100, 2)
		q := edakit.Sprintf("A package holds %d dies that passed wafer test, but test lets bad dies escape so that %s%% of dies shipped as good are actually bad. Each attach succeeds with probability %s. What is the package yield, as a percentage? Round to 2 decimal places.", k, f(e, 1), f(ya, 3))
		return num(q, ans, edakit.Sprintf("Each die is good with probability $1 - %s = %s$. $Y = (%s \\times %s)^{%d} = %s\\%%$.", f(e/100, 3), f(1-e/100, 3), f(1-e/100, 3), f(ya, 3), k, ans), "yield_pct")
	default:
		k := edakit.Int(2, 5)
		sum := 0
		var parts []string
		for i := 0; i < k; i++ {
			c := edakit.Int(3, 40)
			sum += c
			parts = append(parts, edakit.IntStr(c))
		}
		asm := edakit.Int(2, 10)
		Y := step(0.6, 0.95, 0.01)
		ans := f(float64(sum+asm)/Y, 2)
		q := edakit.Sprintf("A SiP uses dies costing %s cost units, plus %d units of substrate and assembly per module. The module yield is %s. What is the cost per good module, in cost units? Round to 2 decimal places.", strings.Join(parts, ", "), asm, f(Y, 2))
		return num(q, ans, edakit.Sprintf("All modules cost $%d + %d = %d$ units, and only a fraction %s are good: $%d/%s = %s$.", sum, asm, sum+asm, f(Y, 2), sum+asm, f(Y, 2), ans), "cost_good")
	}
}

// interposer: reticle fields, wires across an edge, bump-limited wire count.
func interposer(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		k := edakit.Int(2, 4)
		W := step(8, 30, 0.5)
		H := step(10, 30, 0.5)
		g := step(0.1, 0.5, 0.1)
		m := step(0.5, 2, 0.5)
		width := float64(k)*W + float64(k-1)*g + 2*m
		height := H + 2*m
		area := width * height
		ans := edakit.IntStr(cl(area, 858))
		q := edakit.Sprintf("A silicon interposer carries %d dies of %s mm by %s mm in a row, %s mm apart, with a %s mm margin on every side. A reticle field is 26 mm by 33 mm. Counting by area alone, at least how many reticle fields must be stitched?", k, f(W, 1), f(H, 1), f(g, 1), f(m, 1))
		return num(q, ans, edakit.Sprintf("Width $%d \\times %s + %d \\times %s + 2 \\times %s = %s$ mm, height %s mm, area %s mm^2. $\\lceil %s/858 \\rceil = %s$.", k, f(W, 1), k-1, f(g, 1), f(m, 1), f(width, 2), f(height, 2), f(area, 2), f(area, 2), ans), "reticles")
	case 2:
		E := step(3, 12, 0.5)
		p := edakit.Pick([]float64{0.4, 0.6, 0.8, 1.0, 1.2, 1.6, 2.0})
		N := edakit.Int(2, 4)
		per := fl(E*1000*10, p*10)
		ans := edakit.IntStr(N * per)
		q := edakit.Sprintf("Interposer routing along a %s mm die edge uses %d layers at a %s um wiring pitch. How many wires can cross that edge?", f(E, 1), N, f(p, 1))
		return num(q, ans, edakit.Sprintf("$\\lfloor %.0f/%s \\rfloor = %d$ per layer; $%d \\times %d = %s$ wires.", E*1000, f(p, 1), per, N, per, ans), "wires")
	default:
		E := step(3, 12, 0.5)
		wires := 1000 * edakit.Int(5, 40)
		pb := float64(edakit.Pick([]int{25, 30, 36, 40, 45, 50, 55}))
		D := step(0.5, 2, 0.25)
		b := fl(E*1000, pb) * fl(D*1000, pb)
		use := b
		if wires < b {
			use = wires
		}
		ans := edakit.IntStr(use)
		q := edakit.Sprintf("An interposer can route %d wires across a %s mm die edge. The die lands them on signal microbumps at %.0f um pitch in a band %s mm deep along that edge, one wire per bump. How many die-to-die signals can actually cross?", wires, f(E, 1), pb, f(D, 2))
		return num(q, ans, edakit.Sprintf("Bumps: $\\lfloor %.0f/%.0f \\rfloor \\times \\lfloor %.0f/%.0f \\rfloor = %d$. Usable: $\\min(%d, %d) = %s$.", E*1000, pb, D*1000, pb, b, wires, b, ans), "signals")
	}
}

// bridges: silicon saved, wires per bridge set, bandwidth through bridges.
func bridges(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		A := float64(50 * edakit.Int(16, 50))
		n := edakit.Int(2, 8)
		bw := step(2, 10, 0.5)
		bh := step(1, 4, 0.5)
		ans := f(A/(float64(n)*bw*bh), 1)
		q := edakit.Sprintf("A full silicon interposer for a multi-die package would be %.0f mm^2. Instead, %d embedded bridges of %s mm by %s mm connect the facing die edges. By what factor does the silicon area fall? Round to 1 decimal place.", A, n, f(bw, 1), f(bh, 1))
		return num(q, ans, edakit.Sprintf("Bridge silicon: $%d \\times %s \\times %s = %s$ mm^2. Factor $%.0f/%s = %s$.", n, f(bw, 1), f(bh, 1), f(float64(n)*bw*bh, 2), A, f(float64(n)*bw*bh, 2), ans), "factor")
	case 2:
		Eb := step(2, 10, 0.5)
		p := edakit.Pick([]float64{1, 1.5, 2, 2.5, 3})
		N := edakit.Int(2, 3)
		n := edakit.Int(2, 6)
		per := fl(Eb*1000*10, p*10)
		ans := edakit.IntStr(n * N * per)
		q := edakit.Sprintf("Each of %d bridges spans %s mm of die edge with %d wiring layers at %s um pitch. How many die-to-die wires do the bridges carry in total?", n, f(Eb, 1), N, f(p, 1))
		return num(q, ans, edakit.Sprintf("Per layer: $\\lfloor %.0f/%s \\rfloor = %d$. Total: $%d \\times %d \\times %d = %s$.", Eb*1000, f(p, 1), per, n, N, per, ans), "wires")
	default:
		w := 100 * edakit.Int(10, 80)
		r := edakit.Pick([]float64{2, 4, 6, 8, 12, 16})
		n := edakit.Int(2, 6)
		ans := f(float64(w)*r*float64(n)/1000, 2)
		q := edakit.Sprintf("Each of %d bridges carries %d signal wires, each running at %s Gb/s. What is the total bandwidth through the bridges, in Tb/s? Round to 2 decimal places.", n, w, f(r, 0))
		return num(q, ans, edakit.Sprintf("$%d \\times %d \\times %s = %s$ Gb/s, or %s Tb/s.", n, w, f(r, 0), f(float64(w)*r*float64(n), 0), ans), "tbps")
	}
}

// tsvBasics: aspect ratio, TSV resistance, TSV area overhead.
func tsvBasics(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		ar := step(5, 15, 1)
		if edakit.Int(0, 1) == 0 {
			d := edakit.Pick([]float64{2, 3, 4, 5, 6, 8, 10})
			ans := f(ar*d, 1)
			q := edakit.Sprintf("A TSV process fills holes reliably up to an aspect ratio (depth over diameter) of %s. For %s um diameter TSVs, what is the thickest die, in micrometres, that the TSVs can pass through? Round to 1 decimal place.", f(ar, 0), f(d, 0))
			return num(q, ans, edakit.Sprintf("$t = AR \\times d = %s \\times %s = %s$ um.", f(ar, 0), f(d, 0), ans), "t_um")
		}
		t := float64(5 * edakit.Int(6, 40))
		ans := f(t/ar, 2)
		q := edakit.Sprintf("A die is thinned to %.0f um. The TSV process allows an aspect ratio of at most %s. What is the smallest TSV diameter, in micrometres? Round to 2 decimal places.", t, f(ar, 0))
		return num(q, ans, edakit.Sprintf("$d = t/AR = %.0f/%s = %s$ um.", t, f(ar, 0), ans), "d_um")
	case 2:
		d := edakit.Pick([]float64{2, 3, 4, 5, 6, 8, 10})
		t := float64(5 * edakit.Int(4, 24))
		rho := step(1.7, 2.2, 0.1)
		A := math.Pi * d * d / 4
		R := rho * 0.01 * t / A * 1000
		ans := f(R, 1)
		q := edakit.Sprintf("A copper TSV is %.0f um long and %s um in diameter, with resistivity %s micro-ohm cm. What is its resistance, in milliohms? Round to 1 decimal place.", t, f(d, 0), f(rho, 1))
		return num(q, ans, edakit.Sprintf("$A = \\pi %s^2/4 = %s$ um^2. $R = %s \\times %.0f / %s = %s$ milliohm.", f(d, 0), f(A, 3), f(rho*0.01, 3), t, f(A, 3), ans), "r_mohm")
	default:
		for attempt := 0; ; attempt++ {
			N := 1000 * edakit.Int(2, 100)
			s := float64(edakit.Pick([]int{4, 5, 6, 8, 10, 12, 15, 20}))
			A := float64(5 * edakit.Int(5, 80))
			if attempt >= generator.MaxDrawAttempts {
				N, s, A = 20000, 10, 100
			}
			frac := float64(N) * s * s / (A * 1e6) * 100
			if frac > 40 {
				continue
			}
			ans := f(frac, 2)
			q := edakit.Sprintf("A %.0f mm^2 die is crossed by %d TSVs, each with a %.0f um square landing pad. What percentage of the die area do the pads occupy? Round to 2 decimal places.", A, N, s)
			return num(q, ans, edakit.Sprintf("$%d \\times %.0f^2 = %.0f$ um^2 $= %s$ mm^2, which is %s%% of %.0f mm^2.", N, s, float64(N)*s*s, f(float64(N)*s*s/1e6, 4), ans, A), "overhead_pct")
		}
	}
}

// tsvKeepOut: blocked area per TSV, fraction of a block, KOZ from stress decay.
func tsvKeepOut(ctx generator.GeneratorContext) generator.Problem {
	d := edakit.Pick([]float64{2, 3, 4, 5, 6, 8, 10})
	k := step(1, 10, 0.5)
	side := d + 2*k
	switch kpOf(ctx) {
	case 1:
		ans := f(side*side, 2)
		q := edakit.Sprintf("A placer blocks a square keep-out around each TSV. The TSV is %s um in diameter and the keep-out extends %s um beyond the TSV edge on every side. What area does one TSV block, in square micrometres? Round to 2 decimal places.", f(d, 0), f(k, 1))
		return num(q, ans, edakit.Sprintf("Side $%s + 2 \\times %s = %s$ um; area $%s^2 = %s$ um^2.", f(d, 0), f(k, 1), f(side, 1), f(side, 1), ans), "blocked_um2")
	case 2:
		for attempt := 0; ; attempt++ {
			N := 100 * edakit.Int(1, 50)
			A := step(0.5, 10, 0.5)
			if attempt >= generator.MaxDrawAttempts {
				N, A = 100, 10
			}
			frac := float64(N) * side * side / (A * 1e6) * 100
			if frac > 60 {
				continue
			}
			ans := f(frac, 2)
			q := edakit.Sprintf("A %s mm^2 logic block contains %d isolated TSVs of %s um diameter, each with a square keep-out reaching %s um beyond the TSV edge. What percentage of the block is lost to keep-out squares? Round to 2 decimal places.", f(A, 1), N, f(d, 0), f(k, 1))
			return num(q, ans, edakit.Sprintf("Each blocks $(%s + 2 \\times %s)^2 = %s$ um^2; $%d$ of them block %s mm^2, which is %s%% of the block.", f(d, 0), f(k, 1), f(side*side, 2), N, f(float64(N)*side*side/1e6, 4), ans), "lost_pct")
		}
	default:
		R := d / 2
		s0 := float64(10 * edakit.Int(10, 40))
		sth := float64(5 * edakit.Int(2, 16))
		if sth >= s0 {
			sth = s0 / 4
		}
		koz := R * (math.Sqrt(s0/sth) - 1)
		ans := f(koz, 2)
		q := edakit.Sprintf("Thermal-mismatch stress around a TSV of %s um diameter falls as $(R/r)^2$ from %.0f MPa at the TSV edge, where $R$ is the TSV radius and $r$ the distance from its axis. Transistors tolerate %s MPa. How far from the TSV edge must the keep-out extend, in micrometres? Round to 2 decimal places.", f(d, 0), s0, f(sth, 0))
		return num(q, ans, edakit.Sprintf("$r = R\\sqrt{\\sigma_0/\\sigma_{th}} = %s \\sqrt{%.0f/%s} = %s$ um; keep-out $r - R = %s$ um.", f(R, 1), s0, f(sth, 0), f(R*math.Sqrt(s0/sth), 3), ans), "koz_um")
	}
}

const eps0fFum = 8.854e-3 // vacuum permittivity in fF/um

// tsvCapacitance: coaxial C, RC delay, liner sensitivity.
func tsvCapacitance(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		r := step(1, 5, 0.5)
		t := step(0.1, 0.5, 0.05)
		L := float64(5 * edakit.Int(4, 20))
		C := 2 * math.Pi * eps0fFum * 3.9 * L / math.Log((r+t)/r)
		ans := f(C, 1)
		q := edakit.Sprintf("A TSV has a copper core of radius %s um, an oxide liner %s um thick ($\\varepsilon_r = 3.9$) and a length of %.0f um. Using the coaxial formula $C = 2\\pi\\varepsilon_0\\varepsilon_r L/\\ln(r_{out}/r_{in})$ with $\\varepsilon_0 = 8.854 \\times 10^{-3}$ fF/um, what is its oxide capacitance, in fF? Round to 1 decimal place.", f(r, 1), f(t, 2), L)
		return num(q, ans, edakit.Sprintf("$\\ln(%s/%s) = %s$. $C = 2\\pi \\times 8.854 \\times 10^{-3} \\times 3.9 \\times %.0f / %s = %s$ fF.", f(r+t, 2), f(r, 1), f(math.Log((r+t)/r), 4), L, f(math.Log((r+t)/r), 4), ans), "c_ff")
	case 2:
		Rd := step(0.5, 5, 0.5)
		Ct := float64(edakit.Int(20, 200))
		Cl := float64(edakit.Int(2, 30))
		d := 0.69 * Rd * (Ct + Cl)
		ans := f(d, 1)
		q := edakit.Sprintf("A driver with %s kilo-ohm output resistance drives a %.0f fF TSV into a receiver of %.0f fF. Neglecting the TSV's own resistance, what is the 50%% delay $0.69 R (C_{TSV} + C_L)$, in picoseconds? Round to 1 decimal place.", f(Rd, 1), Ct, Cl)
		return num(q, ans, edakit.Sprintf("One kilo-ohm times one fF is 1 ps. $0.69 \\times %s \\times %.0f = %s$ ps.", f(Rd, 1), Ct+Cl, ans), "delay_ps")
	default:
		r := step(1, 5, 0.5)
		t1 := step(0.1, 0.3, 0.05)
		t2 := t1 + step(0.05, 0.4, 0.05)
		ratio := math.Log((r+t1)/r) / math.Log((r+t2)/r)
		ans := f(ratio, 3)
		q := edakit.Sprintf("A TSV core of radius %s um has its oxide liner thickened from %s um to %s um, with the core radius and length unchanged. By what factor does the coaxial capacitance change (new over old)? Round to 3 decimal places.", f(r, 1), f(t1, 2), f(t2, 2))
		return num(q, ans, edakit.Sprintf("$C \\propto 1/\\ln(1 + t/r)$: $\\ln(%s/%s)/\\ln(%s/%s) = %s/%s = %s$.", f(r+t1, 2), f(r, 1), f(r+t2, 2), f(r, 1), f(math.Log((r+t1)/r), 4), f(math.Log((r+t2)/r), 4), ans), "ratio")
	}
}

// microbumps: density, power bump count, joint resistance.
func microbumps(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		p := float64(edakit.Pick([]int{20, 25, 30, 36, 40, 45, 50, 55}))
		if edakit.Int(0, 1) == 0 {
			ans := f(math.Pow(1000/p, 2), 0)
			q := edakit.Sprintf("Microbumps sit on a square grid at %.0f um pitch. How many bumps per square millimetre is that? Round to the nearest whole number.", p)
			return num(q, ans, edakit.Sprintf("$(1000/%.0f)^2 = %s$ per mm^2.", p, ans), "density")
		}
		A := float64(edakit.Int(2, 40))
		ans := f(A*math.Pow(1000/p, 2), 0)
		q := edakit.Sprintf("A %.0f mm^2 interface region is filled with microbumps on a square grid at %.0f um pitch. Using density times area, about how many bumps does it hold? Round to the nearest whole number.", A, p)
		return num(q, ans, edakit.Sprintf("Density $(1000/%.0f)^2 = %s$ per mm^2, times %.0f mm^2 gives %s.", p, f(math.Pow(1000/p, 2), 2), A, ans), "bumps")
	case 2:
		I := step(5, 150, 5)
		Im := float64(edakit.Pick([]int{50, 75, 100, 150, 200}))
		ans := edakit.IntStr(cl(I*1000, Im))
		q := edakit.Sprintf("A stacked die draws %s A from its supply through microbumps limited to %.0f mA each by electromigration. Assuming perfect current sharing, how many supply bumps are needed, not counting the ground return?", f(I, 0), Im)
		return num(q, ans, edakit.Sprintf("$\\lceil %.0f/%.0f \\rceil = %s$ bumps, and as many again for ground; uneven sharing needs margin on top.", I*1000, Im, ans), "bumps")
	default:
		rho := step(10, 16, 0.5)
		h := float64(edakit.Int(5, 25))
		d := float64(edakit.Int(10, 30))
		A := math.Pi * d * d / 4
		R := rho * 0.01 * h / A * 1000
		ans := f(R, 2)
		q := edakit.Sprintf("A microbump's solder cap is modelled as a cylinder %.0f um tall and %.0f um in diameter with resistivity %s micro-ohm cm. What is its resistance, in milliohms? Round to 2 decimal places.", h, d, f(rho, 1))
		return num(q, ans, edakit.Sprintf("$A = \\pi %.0f^2/4 = %s$ um^2. $R = %s \\times %.0f/%s = %s$ milliohm.", d, f(A, 2), f(rho*0.01, 3), h, f(A, 2), ans), "r_mohm")
	}
}

// hybridBonding: density, gain over microbumps, pad overlap under misalignment.
func hybridBonding(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		p := edakit.Pick([]float64{1, 1.5, 2, 3, 4, 5, 6, 9, 10})
		if edakit.Int(0, 1) == 0 {
			ans := f(math.Pow(1000/p, 2), 0)
			q := edakit.Sprintf("A hybrid-bonding interface places pads on a square grid at %s um pitch. How many connections per square millimetre does it offer? Round to the nearest whole number.", f(p, 1))
			return num(q, ans, edakit.Sprintf("$(1000/%s)^2 = %s$ per mm^2.", f(p, 1), ans), "density")
		}
		A := float64(edakit.Int(10, 200))
		ans := f(A*math.Pow(1000/p, 2)/1e6, 2)
		q := edakit.Sprintf("Two dies are hybrid bonded over a %.0f mm^2 overlap with pads at %s um pitch on a square grid. How many connections is that, in millions? Round to 2 decimal places.", A, f(p, 1))
		return num(q, ans, edakit.Sprintf("Density $(1000/%s)^2 = %s$ per mm^2; times %.0f gives %s million.", f(p, 1), f(math.Pow(1000/p, 2), 1), A, ans), "millions")
	case 2:
		pm := float64(edakit.Pick([]int{25, 30, 36, 40, 45, 50, 55}))
		ph := edakit.Pick([]float64{1, 1.5, 2, 3, 4, 5, 6, 9, 10})
		ans := f(math.Pow(pm/ph, 2), 1)
		q := edakit.Sprintf("An interface moves from microbumps at %.0f um pitch to hybrid bonding at %s um pitch. By what factor does the connection density rise? Round to 1 decimal place.", pm, f(ph, 1))
		return num(q, ans, edakit.Sprintf("Density scales as $1/p^2$: $(%.0f/%s)^2 = %s$.", pm, f(ph, 1), ans), "gain")
	default:
		a := step(1, 5, 0.5)
		ex := step(0, 0.5, 0.05)
		ey := step(0, 0.5, 0.05)
		ov := (a - ex) * (a - ey) / (a * a) * 100
		ans := f(ov, 1)
		q := edakit.Sprintf("Hybrid-bond pads are squares %s um on a side. After bonding, the top die is offset by %s um in x and %s um in y. What percentage of each pad's area still overlaps its partner? Round to 1 decimal place.", f(a, 1), f(ex, 2), f(ey, 2))
		return num(q, ans, edakit.Sprintf("$(%s - %s)(%s - %s)/%s^2 = %s/%s = %s\\%%$.", f(a, 1), f(ex, 2), f(a, 1), f(ey, 2), f(a, 1), f((a-ex)*(a-ey), 4), f(a*a, 2), ans), "overlap_pct")
	}
}
