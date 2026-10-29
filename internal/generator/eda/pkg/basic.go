package pkg

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// wireBonding: perimeter pad count, bond wire resistance and ground bounce.
func wireBonding(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		S := step(3, 12, 0.5)
		p := float64(5 * edakit.Int(8, 20))
		rows := edakit.Int(1, 2)
		n := fl(S*1000, p)
		total := 4 * rows * n
		ans := edakit.IntStr(total)
		rowWord := "one row"
		if rows == 2 {
			rowWord = "two staggered rows"
		}
		q := edakit.Sprintf("A square die of side %s mm has wire-bond pads in %s along every edge, each row at a pad pitch of %.0f um. Each row holds $\\lfloor$side/pitch$\\rfloor$ pads per side. How many pads does the die hold in total?", f(S, 1), rowWord, p)
		return num(q, ans, edakit.Sprintf("Per side and row: $\\lfloor %.0f/%.0f \\rfloor = %d$. Total: $4 \\times %d \\times %d = %s$ pads. Pad count grows with perimeter, not area.", S*1000, p, n, rows, n, ans), "pads")
	case 2:
		metals := []struct {
			name string
			rho  float64
		}{{"gold", 2.44}, {"copper", 1.7}, {"aluminium", 2.8}}
		m := edakit.Pick(metals)
		L := step(1, 5, 0.1)
		d := float64(edakit.Pick([]int{18, 20, 23, 25, 30, 33}))
		area := math.Pi * d * d / 4 * 1e-12
		R := m.rho * 1e-8 * L * 1e-3 / area * 1000
		ans := f(R, 1)
		q := edakit.Sprintf("A %s bond wire is %s mm long with a diameter of %.0f um. Taking the resistivity of %s as %s micro-ohm cm, what is the wire's resistance, in milliohms? Round to 1 decimal place.", m.name, f(L, 1), d, m.name, f(m.rho, 2))
		return num(q, ans, edakit.Sprintf("$A = \\pi d^2/4 = %.4g$ m^2. $R = \\rho L/A = %s \\times 10^{-8} \\times %s \\times 10^{-3} / %.4g = %s$ milliohm.", area, f(m.rho, 2), f(L, 1), area, ans), "r_mohm")
	default:
		L := step(1, 5, 0.5)
		N := edakit.Int(1, 8)
		dI := step(0.1, 1.0, 0.05)
		tr := step(0.5, 2, 0.1)
		V := L * dI / tr / float64(N) * 1000
		ans := f(V, 0)
		wires := "one ground wire"
		if N > 1 {
			wires = edakit.Sprintf("%d identical ground wires in parallel (ignore mutual inductance)", N)
		}
		q := edakit.Sprintf("Simultaneously switching outputs draw an extra %s A of ground current in %s ns. The ground return is %s, each %s mm long with 1 nH/mm of inductance. What is the ground bounce $L\\,\\Delta I/t_r$, in millivolts? Round to the nearest whole number.", f(dI, 2), f(tr, 1), wires, f(L, 1))
		return num(q, ans, edakit.Sprintf("One wire has %s nH; %d in parallel give %s nH. $V = %s \\times %s / %s = %s$ mV.", f(L, 1), N, f(L/float64(N), 3), f(L/float64(N), 3), f(dI, 2), f(tr, 1), ans), "bounce_mv")
	}
}

// flipChipBGA: area-array bump count, area versus perimeter, BGA ball count.
func flipChipBGA(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		W := step(5, 25, 0.5)
		H := step(5, 25, 0.5)
		p := float64(edakit.Pick([]int{130, 140, 150, 160, 180, 200}))
		nx, ny := fl(W*1000, p), fl(H*1000, p)
		ans := edakit.IntStr(nx * ny)
		q := edakit.Sprintf("A flip-chip die is %s mm by %s mm. Bumps sit on a full square grid at %.0f um pitch, the first bump half a pitch in from each edge. How many bumps does the die carry?", f(W, 1), f(H, 1), p)
		return num(q, ans, edakit.Sprintf("$\\lfloor %.0f/%.0f \\rfloor \\times \\lfloor %.0f/%.0f \\rfloor = %d \\times %d = %s$ bumps.", W*1000, p, H*1000, p, nx, ny, ans), "bumps")
	case 2:
		n := edakit.Scaled(ctx, 10, 80)
		p := float64(edakit.Pick([]int{100, 130, 150, 180, 200}))
		S := float64(n) * p / 1000
		r := float64(n*n) / float64(4*(n-1))
		ans := f(r, 2)
		q := edakit.Sprintf("A square die of side %s mm uses %.0f um pitch, so a full area array would be %d by %d sites. A perimeter design uses only the outer ring of the same grid. How many times more connections does the area array provide? Round to 2 decimal places.", f(S, 2), p, n, n)
		return num(q, ans, edakit.Sprintf("Area array: $%d^2 = %d$. Outer ring: $4(%d - 1) = %d$. Ratio $%d/%d = %s$, growing linearly with $n$.", n, n*n, n, 4*(n-1), n*n, 4*(n-1), ans), "ratio")
	default:
		L := float64(edakit.Int(15, 45))
		p := edakit.Pick([]float64{1.0, 0.8, 0.65, 0.5})
		n := fl(L*1000, p*1000)
		k := edakit.Int(0, n/3)
		ans := edakit.IntStr(n*n - k*k)
		hole := "with no depopulated balls"
		if k > 0 {
			hole = edakit.Sprintf("with a central %d by %d block of sites depopulated", k, k)
		}
		q := edakit.Sprintf("A %.0f mm square BGA package has a full ball grid at %s mm pitch, %s. Count the grid as $\\lfloor$side/pitch$\\rfloor$ balls per side. How many balls does it have?", L, f(p, 2), hole)
		return num(q, ans, edakit.Sprintf("$\\lfloor %.0f/%s \\rfloor = %d$ per side; $%d^2 - %d^2 = %s$ balls.", L, f(p, 2), n, n, k, ans), "balls")
	}
}

// substrates: wiring capacity, layer count and CTE displacement.
func substrates(ctx generator.GeneratorContext) generator.Problem {
	w := float64(edakit.Pick([]int{5, 6, 8, 10, 12, 15, 20}))
	s := float64(edakit.Pick([]int{5, 6, 8, 10, 12, 15, 20}))
	per := fl(1000, w+s)
	switch kpOf(ctx) {
	case 1:
		E := float64(edakit.Int(5, 20))
		N := edakit.Int(2, 8)
		ans := edakit.IntStr(int(E) * N * per)
		q := edakit.Sprintf("A package substrate routes at %.0f um line and %.0f um space. How many traces can cross a %.0f mm die edge using %d signal layers?", w, s, E, N)
		return num(q, ans, edakit.Sprintf("Per millimetre per layer: $\\lfloor 1000/%.0f \\rfloor = %d$. Total: $%.0f \\times %d \\times %d = %s$.", w+s, per, E, N, per, ans), "traces")
	case 2:
		E := float64(edakit.Int(5, 20))
		S := 100 * edakit.Int(10, 80)
		cap1 := int(E) * per
		ans := edakit.IntStr(int(math.Ceil(float64(S) / float64(cap1))))
		q := edakit.Sprintf("%d signals must leave a die across one %.0f mm edge on a substrate with %.0f/%.0f um line and space. How many signal layers are needed at minimum?", S, E, w, s)
		return num(q, ans, edakit.Sprintf("One layer carries $%.0f \\times \\lfloor 1000/%.0f \\rfloor = %d$ traces. Layers: $\\lceil %d/%d \\rceil = %s$.", E, w+s, cap1, S, cap1, ans), "layers")
	default:
		glass := edakit.Int(0, 1) == 1
		var a float64
		name := "an organic core"
		if glass {
			a = step(3.5, 9, 0.5)
			name = "a glass core"
		} else {
			a = step(12, 18, 0.5)
		}
		D := step(5, 25, 0.5)
		dT := float64(10 * edakit.Int(6, 20))
		d := (a - 2.6) * 1e-6 * D * 1000 * dT
		ans := f(d, 2)
		q := edakit.Sprintf("A silicon die (2.6 ppm/K) sits on a substrate with %s of %s ppm/K. For a point %s mm from the neutral point and a %.0f K temperature swing, by how many micrometres do die and substrate try to move apart? Round to 2 decimal places.", name, f(a, 1), f(D, 1), dT)
		return num(q, ans, edakit.Sprintf("$|\\alpha_{sub} - \\alpha_{Si}| D \\Delta T = %s \\times 10^{-6} \\times %.0f \\times %.0f = %s$ um.", f(a-2.6, 1), D*1000, dT, ans), "mismatch_um")
	}
}

// escapeRouting: traces per channel, layers for a row count, balls escaped on one layer.
func escapeRouting(ctx generator.GeneratorContext) generator.Problem {
	var p, d, w, s float64
	n := 0
	for attempt := 0; attempt < generator.MaxDrawAttempts && (n < 1 || n > 5); attempt++ {
		p = float64(edakit.Pick([]int{1270, 1000, 800, 650, 500}))
		d = 5 * math.Round(p*step(0.4, 0.6, 0.05)/5)
		w = float64(edakit.Pick([]int{25, 30, 40, 50, 60, 75, 100}))
		s = float64(edakit.Pick([]int{25, 30, 40, 50, 60, 75, 100}))
		n = fl(p-d-s, w+s)
	}
	if n < 1 || n > 5 {
		p, d, w, s = 1000, 500, 75, 75
		n = fl(p-d-s, w+s)
	}
	geo := edakit.Sprintf("A BGA has %.0f um ball pitch and %.0f um pads. The board routes at %.0f um trace width and %.0f um space.", p, d, w, s)
	switch kpOf(ctx) {
	case 1:
		ans := edakit.IntStr(n)
		q := geo + " How many traces fit through one channel between adjacent pads, if $n$ traces need $n$ widths and $n+1$ spaces?"
		return num(q, ans, edakit.Sprintf("$n = \\lfloor (%.0f - %.0f - %.0f)/(%.0f + %.0f) \\rfloor = \\lfloor %.0f/%.0f \\rfloor = %s$.", p, d, s, w, s, p-d-s, w+s, ans), "traces_per_channel")
	case 2:
		R := edakit.Scaled(ctx, 4, 20)
		L := int(math.Ceil(float64(R) / float64(n+1)))
		ans := edakit.IntStr(L)
		q := geo + edakit.Sprintf(" On each layer the outer row escapes directly and each channel passes its traces from inner rows. The array has %d rows from edge to centre, all signals. How many routing layers are needed?", R)
		return num(q, ans, edakit.Sprintf("Traces per channel: %d, so one layer escapes %d rows. Layers: $\\lceil %d/%d \\rceil = %s$.", n, n+1, R, n+1, ans), "layers")
	default:
		N := edakit.Scaled(ctx, 16, 40)
		k := n + 1
		inner := N - 2*k
		if inner < 0 {
			inner = 0
		}
		ans := edakit.IntStr(N*N - inner*inner)
		q := edakit.Sprintf("A full %d by %d BGA array is routed on a board that passes %d traces per channel, so the first layer escapes the outer %d rings of balls. How many balls escape on the first layer?", N, N, n, k)
		return num(q, ans, edakit.Sprintf("Inner square: $(%d - 2 \\times %d)^2 = %d$. Escaped: $%d - %d = %s$.", N, k, inner*inner, N*N, inner*inner, ans), "escaped")
	}
}

// waferLevel: fan-in capacity, fan-out area expansion, units per panel.
func waferLevel(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		W := step(2, 8, 0.1)
		H := step(2, 8, 0.1)
		p := edakit.Pick([]float64{0.3, 0.35, 0.4, 0.5})
		nx, ny := fl(W*1000, p*1000), fl(H*1000, p*1000)
		ans := edakit.IntStr(nx * ny)
		q := edakit.Sprintf("A fan-in wafer-level package places balls on a full grid under a %s mm by %s mm die at %s mm pitch. How many balls fit, counting $\\lfloor$side/pitch$\\rfloor$ per side?", f(W, 1), f(H, 1), f(p, 2))
		return num(q, ans, edakit.Sprintf("$\\lfloor %s/%s \\rfloor \\times \\lfloor %s/%s \\rfloor = %d \\times %d = %s$ balls.", f(W, 1), f(p, 2), f(H, 1), f(p, 2), nx, ny, ans), "balls")
	case 2:
		a := step(3, 8, 0.5)
		p := edakit.Pick([]float64{0.35, 0.4, 0.5})
		nmin := fl(a*1000, p*1000) + 2
		n := edakit.Int(nmin, nmin+15)
		B := edakit.Int((n-1)*(n-1)+1, n*n)
		side := float64(n) * p
		r := side * side / (a * a)
		ans := f(r, 2)
		q := edakit.Sprintf("A %s mm square die needs %d balls at %s mm pitch in a square full array. A fan-out package provides the area. What is the ratio of package area to die area? Round to 2 decimal places.", f(a, 1), B, f(p, 2))
		return num(q, ans, edakit.Sprintf("Grid side $\\lceil \\sqrt{%d} \\rceil = %d$, package side $%d \\times %s = %s$ mm. Ratio $%s^2/%s^2 = %s$.", B, n, n, f(p, 2), f(side, 2), f(side, 2), f(a, 1), ans), "area_ratio")
	default:
		panels := []struct{ w, h float64 }{{510, 515}, {600, 600}, {300, 300}, {320, 320}}
		pn := edakit.Pick(panels)
		sz := step(5, 20, 0.5)
		g := edakit.Pick([]float64{0.1, 0.2, 0.3})
		st := sz + g
		nx, ny := fl(pn.w*10, st*10), fl(pn.h*10, st*10)
		ans := edakit.IntStr(nx * ny)
		q := edakit.Sprintf("A fan-out panel of %.0f mm by %.0f mm holds square packages %s mm on a side separated by %s mm saw streets, so the layout step is size plus street. How many units fit on the panel?", pn.w, pn.h, f(sz, 1), f(g, 1))
		return num(q, ans, edakit.Sprintf("Step %s mm. $\\lfloor %.0f/%s \\rfloor \\times \\lfloor %.0f/%s \\rfloor = %d \\times %d = %s$ units.", f(st, 1), pn.w, f(st, 1), pn.h, f(st, 1), nx, ny, ans), "units")
	}
}

// rdl: trace resistance, IR drop, density against substrate.
func rdl(ctx generator.GeneratorContext) generator.Problem {
	rho := step(1.8, 2.4, 0.1)
	L := float64(100 * edakit.Int(5, 100))
	w := float64(edakit.Pick([]int{2, 3, 4, 5, 8, 10, 15, 20}))
	t := float64(edakit.Pick([]int{1, 2, 3, 4, 5}))
	R := 0.01 * rho * L / (w * t)
	geo := edakit.Sprintf("An RDL copper trace is %.0f um long, %.0f um wide and %.0f um thick, with resistivity %s micro-ohm cm.", L, w, t, f(rho, 1))
	switch kpOf(ctx) {
	case 1:
		ans := f(R, 2)
		q := geo + " What is its resistance, in ohms? Round to 2 decimal places."
		return num(q, ans, edakit.Sprintf("%s micro-ohm cm is %s ohm um. $R = %s \\times %.0f/(%.0f \\times %.0f) = %s$ ohm.", f(rho, 1), f(rho*0.01, 3), f(rho*0.01, 3), L, w, t, ans), "r_ohm")
	case 2:
		I := float64(10 * edakit.Int(2, 30))
		V := I * R
		ans := f(V, 1)
		q := geo + edakit.Sprintf(" It carries %.0f mA of supply current. What is the IR drop along it, in millivolts? Round to 1 decimal place.", I)
		return num(q, ans, edakit.Sprintf("$R = 0.01 \\times %s \\times %.0f/(%.0f \\times %.0f) = %s$ ohm. $V = %.0f \\times %s = %s$ mV.", f(rho, 1), L, w, t, f(R, 4), I, f(R, 4), ans), "drop_mv")
	default:
		w1 := float64(edakit.Pick([]int{8, 10, 12, 15, 20}))
		s1 := float64(edakit.Pick([]int{8, 10, 12, 15, 20}))
		w2 := edakit.Pick([]float64{1, 1.5, 2, 3, 5})
		s2 := edakit.Pick([]float64{1, 1.5, 2, 3, 5})
		ans := f((w1+s1)/(w2+s2), 2)
		q := edakit.Sprintf("A laminate substrate routes at %.0f/%.0f um line and space, and a fan-out RDL at %s/%s um. How many times more wires per millimetre per layer does the RDL provide? Use the pitch ratio. Round to 2 decimal places.", w1, s1, f(w2, 1), f(s2, 1))
		return num(q, ans, edakit.Sprintf("Density is $1/(w+s)$, so the ratio is $(%.0f + %.0f)/(%s + %s) = %s$.", w1, s1, f(w2, 1), f(s2, 1), ans), "density_ratio")
	}
}
