package pkg

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type material struct {
	name string
	k    float64
}

var conductors = []material{
	{"silicon", 150}, {"copper", 400}, {"SAC solder", 58}, {"thermal grease", 4},
	{"indium", 82}, {"mold compound", 0.9}, {"underfill", 0.55}, {"thermal gel", 3},
	{"silicon dioxide", 1.4}, {"aluminium nitride", 170},
}

// conduction: layer resistance, temperature drop, required conductivity.
func conduction(ctx generator.GeneratorContext) generator.Problem {
	m := edakit.Pick(conductors)
	L := float64(10 * edakit.Int(2, 80))
	A := float64(25 * edakit.Int(2, 32))
	R := L / (m.k * A)
	switch kpOf(ctx) {
	case 1:
		ans := f(R, 4)
		q := edakit.Sprintf("A layer of %s (k = %s W/(m K)) is %.0f um thick and %.0f mm^2 in area. What is its one-dimensional thermal resistance $L/(kA)$, in K/W? Round to 4 decimal places.", m.name, f(m.k, 2), L, A)
		return num(q, ans, edakit.Sprintf("With $L$ in um and $A$ in mm^2 the factors of $10^{-6}$ cancel: $%.0f/(%s \\times %.0f) = %s$ K/W.", L, f(m.k, 2), A, ans), "r_kw")
	case 2:
		P := float64(5 * edakit.Int(1, 60))
		ans := f(P*R, 2)
		q := edakit.Sprintf("%.0f W flows straight through a %.0f um layer of %s (k = %s W/(m K)) of area %.0f mm^2. What temperature drop does the layer cause, in kelvin? Round to 2 decimal places.", P, L, m.name, f(m.k, 2), A)
		return num(q, ans, edakit.Sprintf("$R = %.0f/(%s \\times %.0f) = %s$ K/W; $\\Delta T = %.0f \\times %s = %s$ K.", L, f(m.k, 2), A, f(R, 5), P, f(R, 5), ans), "dt_k")
	default:
		P := float64(10 * edakit.Int(5, 40))
		BLT := float64(5 * edakit.Int(2, 20))
		Ab := float64(50 * edakit.Int(2, 16))
		dT := step(1, 8, 0.5)
		k := P * BLT / (Ab * dT)
		ans := f(k, 2)
		q := edakit.Sprintf("A %.0f um interface layer over %.0f mm^2 must carry %.0f W with at most a %s K drop. What minimum thermal conductivity does the material need, in W/(m K)? Round to 2 decimal places.", BLT, Ab, P, f(dT, 1))
		return num(q, ans, edakit.Sprintf("$k \\ge QL/(A\\Delta T) = %.0f \\times %.0f/(%.0f \\times %s) = %s$ W/(m K).", P, BLT, Ab, f(dT, 1), ans), "k_min")
	}
}

// resistanceNetwork: series stack, parallel paths, heat split.
func resistanceNetwork(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		names := []string{"die", "TIM1", "lid", "TIM2", "heat sink", "spreader"}
		n := edakit.Scaled(ctx, 3, 6)
		sum := 0.0
		var parts []string
		for i := 0; i < n; i++ {
			r := step(0.01, 0.5, 0.01)
			sum += r
			parts = append(parts, edakit.Sprintf("%s %s K/W", names[i], f(r, 2)))
		}
		P := float64(10 * edakit.Int(2, 20))
		if P*sum > 90 {
			P = math.Max(10, 10*math.Floor(9/sum))
		}
		Ta := float64(edakit.Int(20, 50))
		ans := f(Ta+P*sum, 2)
		q := edakit.Sprintf("All the heat of a %.0f W die flows in series through: %s. Ambient is %.0f °C. What is the junction temperature, in °C? Round to 2 decimal places.", P, strings.Join(parts, ", "), Ta)
		return num(q, ans, edakit.Sprintf("Series sum %s K/W; $T_j = %.0f + %.0f \\times %s = %s$ °C.", f(sum, 2), Ta, P, f(sum, 2), ans), "tj")
	case 2:
		n := edakit.Int(2, 3)
		g := 0.0
		var parts []string
		for i := 0; i < n; i++ {
			r := step(0.2, 5, 0.1)
			g += 1 / r
			parts = append(parts, f(r, 1))
		}
		ans := f(1/g, 3)
		q := edakit.Sprintf("Heat leaves a die through %d parallel paths to ambient with resistances %s K/W. What is the combined junction-to-ambient resistance, in K/W? Round to 3 decimal places.", n, strings.Join(parts, ", "))
		return num(q, ans, edakit.Sprintf("The conductances $1/R_i$ of paths %s K/W sum to %s W/K; inverting gives %s K/W.", strings.Join(parts, ", "), f(g, 4), ans), "r_eq")
	default:
		Rt := step(0.2, 1.5, 0.05)
		Rb := step(1, 8, 0.1)
		ans := f(Rt/(Rt+Rb)*100, 1)
		q := edakit.Sprintf("A die loses heat upward through lid and heat sink with %s K/W to ambient, and downward through bumps, substrate and board with %s K/W to ambient. What percentage of the heat leaves through the board path? Round to 1 decimal place.", f(Rt, 2), f(Rb, 1))
		return num(q, ans, edakit.Sprintf("Heat divides inversely to resistance: $R_t/(R_t + R_b) = %s/%s = %s\\%%$.", f(Rt, 2), f(Rt+Rb, 2), ans), "board_pct")
	}
}

// junctionTemp: Tj from theta_ja, maximum power, series case and sink.
func junctionTemp(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		th := step(0.5, 30, 0.1)
		P := step(0.5, math.Max(1, math.Floor(140/th)/2), 0.5)
		Ta := float64(edakit.Int(20, 70))
		ans := f(Ta+P*th, 1)
		q := edakit.Sprintf("A part dissipates %s W with $\\theta_{ja} = %s$ K/W in %.0f °C ambient. What is its junction temperature, in °C? Round to 1 decimal place.", f(P, 1), f(th, 1), Ta)
		return num(q, ans, edakit.Sprintf("$T_j = %.0f + %s \\times %s = %s$ °C.", Ta, f(P, 1), f(th, 1), ans), "tj")
	case 2:
		Tj := float64(5 * edakit.Int(17, 25))
		Ta := float64(5 * edakit.Int(4, 12))
		th := step(0.2, 10, 0.1)
		ans := f((Tj-Ta)/th, 2)
		q := edakit.Sprintf("A part has $T_{j,max} = %.0f$ °C and $\\theta_{ja} = %s$ K/W, and must work in ambient up to %.0f °C. What is the maximum power it may dissipate, in watts? Round to 2 decimal places.", Tj, f(th, 1), Ta)
		return num(q, ans, edakit.Sprintf("$P_{max} = (%.0f - %.0f)/%s = %s$ W.", Tj, Ta, f(th, 1), ans), "pmax")
	default:
		P := float64(5 * edakit.Int(4, 60))
		if edakit.Int(0, 1) == 0 {
			jc := step(0.05, 0.4, 0.01)
			cs := step(0.02, 0.2, 0.01)
			sa := step(0.1, 0.8, 0.01)
			Ta := float64(edakit.Int(20, 45))
			ans := f(Ta+P*(jc+cs+sa), 1)
			q := edakit.Sprintf("A %.0f W processor has $\\theta_{jc} = %s$ K/W, an interface with $\\theta_{cs} = %s$ K/W and a heat sink with $\\theta_{sa} = %s$ K/W, in %.0f °C air. What is the junction temperature, in °C? Round to 1 decimal place.", P, f(jc, 2), f(cs, 2), f(sa, 2), Ta)
			return num(q, ans, edakit.Sprintf("$T_j = %.0f + %.0f \\times (%s + %s + %s) = %s$ °C.", Ta, P, f(jc, 2), f(cs, 2), f(sa, 2), ans), "tj")
		}
		psi := step(0.05, 1.0, 0.05)
		Tt := float64(edakit.Int(40, 90))
		ans := f(Tt+P*psi, 1)
		q := edakit.Sprintf("In a working system, a thermocouple reads %.0f °C on the top of a package that dissipates %.0f W. Its datasheet gives $\\Psi_{jt} = %s$ K/W. What is the estimated junction temperature, in °C? Round to 1 decimal place.", Tt, P, f(psi, 2))
		return num(q, ans, edakit.Sprintf("$\\Psi_{jt}$ relates junction to measured top temperature: $T_j = %.0f + %.0f \\times %s = %s$ °C.", Tt, P, f(psi, 2), ans), "tj")
	}
}

// powerDensity: average density, hot-spot ratio, hot-spot temperature excess.
func powerDensity(ctx generator.GeneratorContext) generator.Problem {
	P := float64(10 * edakit.Int(2, 40))
	A := float64(10 * edakit.Int(10, 80))
	switch kpOf(ctx) {
	case 1:
		ans := f(P/(A/100), 2)
		q := edakit.Sprintf("A %.0f W die has an area of %.0f mm^2. What is its average power density, in W/cm^2? Round to 2 decimal places.", P, A)
		return num(q, ans, edakit.Sprintf("$%.0f$ mm^2 $= %s$ cm^2; $%.0f/%s = %s$ W/cm^2.", A, f(A/100, 2), P, f(A/100, 2), ans), "wcm2")
	case 2:
		Ab := float64(edakit.Int(5, 60))
		Pb := step(0.15, 0.5, 0.01) * P
		if Ab >= A {
			Ab = A / 4
		}
		ratio := (Pb / Ab) / (P / A)
		ans := f(ratio, 2)
		q := edakit.Sprintf("A %.0f W, %.0f mm^2 die contains a %.0f mm^2 block that dissipates %s W. How many times the die-average power density is the block's density? Round to 2 decimal places.", P, A, Ab, f(Pb, 1))
		return num(q, ans, edakit.Sprintf("Block $%s/%.0f = %s$ W/mm^2, average $%.0f/%.0f = %s$ W/mm^2, ratio %s.", f(Pb, 1), Ab, f(Pb/Ab, 4), P, A, f(P/A, 4), ans), "ratio")
	default:
		qh := step(0.5, 3, 0.05)
		qb := step(0.1, 0.45, 0.05)
		Rpp := float64(edakit.Int(5, 40))
		ans := f((qh-qb)*Rpp, 2)
		q := edakit.Sprintf("In a vertical thermal model, each region rises above the coolant by $q'' R''$ with $R'' = %.0f$ K mm^2/W. A hot spot dissipates %s W/mm^2 and its surroundings %s W/mm^2. Ignoring lateral spreading, how many kelvin hotter is the hot spot than its surroundings? Round to 2 decimal places.", Rpp, f(qh, 2), f(qb, 2))
		return num(q, ans, edakit.Sprintf("$(%s - %s) \\times %.0f = %s$ K. Lateral spreading would lower this.", f(qh, 2), f(qb, 2), Rpp, ans), "dt_k")
	}
}

// jouleHeating: I^2 R, wire temperature rise, transistor self-heating.
func jouleHeating(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		I := float64(10 * edakit.Int(1, 50))
		R := step(2, 50, 0.5)
		P := (I / 1000) * (I / 1000) * (R / 1000) * 1000
		ans := f(P, 4)
		q := edakit.Sprintf("A power bump of %s milliohm carries %.0f mA. How much power does it dissipate, in milliwatts? Round to 4 decimal places.", f(R, 1), I)
		return num(q, ans, edakit.Sprintf("$P = I^2 R = %s^2 \\times %s = %s$ mW.", f(I/1000, 3), f(R/1000, 4), ans), "p_mw")
	case 2:
		I := step(0.5, 5, 0.1)
		R := float64(edakit.Int(10, 200))
		Rth := float64(500 * edakit.Int(1, 40))
		P := (I / 1000) * (I / 1000) * R
		ans := f(P*Rth, 2)
		q := edakit.Sprintf("An upper-metal wire segment of %.0f ohm carries %s mA DC. Its thermal resistance to the substrate is %.0f K/W. How far above the substrate does it heat, in kelvin? Round to 2 decimal places.", R, f(I, 1), Rth)
		return num(q, ans, edakit.Sprintf("$P = (%s \\times 10^{-3})^2 \\times %.0f = %s$ mW; $\\Delta T = %s \\times 10^{-3} \\times %.0f = %s$ K.", f(I, 1), R, f(P*1000, 4), f(P*1000, 4), Rth, ans), "dt_k")
	default:
		Id := float64(10 * edakit.Int(2, 50))
		V := step(0.3, 1.0, 0.05)
		Rth := float64(edakit.Int(5, 60))
		P := Id * V / 1000
		ans := f(P*Rth, 2)
		q := edakit.Sprintf("A FinFET conducts %.0f uA at $V_{DS} = %s$ V. Its self-heating thermal resistance is %.0f K/mW. How much does the channel heat above the substrate, in kelvin? Round to 2 decimal places.", Id, f(V, 2), Rth)
		return num(q, ans, edakit.Sprintf("$P = %.0f \\times %s = %s$ uW $= %s$ mW; $\\Delta T = %s \\times %.0f = %s$ K.", Id, f(V, 2), f(Id*V, 1), f(P, 4), f(P, 4), Rth, ans), "dt_k")
	}
}

// thermalSim: a Jacobi update, an exact chain, and 3-D FDM sparsity.
func thermalSim(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		Tl := step(40, 90, 0.5)
		Tr := step(40, 90, 0.5)
		R := step(0.5, 5, 0.5)
		Q := step(0, 4, 0.25)
		ans := f((Tl+Tr+Q*R)/2, 2)
		q := edakit.Sprintf("In a 1-D finite-difference thermal model, each node joins its two neighbours through %s K/W. An interior node receives %s W of power. Its neighbours currently hold %s °C and %s °C. What is the node's new temperature after one Jacobi update, in °C? Round to 2 decimal places.", f(R, 1), f(Q, 2), f(Tl, 1), f(Tr, 1))
		return num(q, ans, edakit.Sprintf("Energy balance gives $T_i = (T_{i-1} + T_{i+1} + Q R)/2 = (%s + %s + %s)/2 = %s$ °C.", f(Tl, 1), f(Tr, 1), f(Q*R, 2), ans), "t_new")
	case 2:
		N := edakit.Scaled(ctx, 3, 10)
		R := step(0.1, 1, 0.05)
		qq := step(0.5, 5, 0.5)
		Ts := float64(edakit.Int(25, 60))
		ans := f(Ts+R*qq*float64(N*(N+1))/2, 2)
		q := edakit.Sprintf("A chain of %d nodes runs from a heat sink held at %.0f °C: the sink joins node 1 through %s K/W, and each node joins the next through %s K/W. Every node receives %s W, and the far end is adiabatic. What is the steady temperature of the far node, in °C? Round to 2 decimal places.", N, Ts, f(R, 2), f(R, 2), f(qq, 1))
		return num(q, ans, edakit.Sprintf("Link $k$ carries $(N-k+1)q$, so the rise is $Rq N(N+1)/2 = %s \\times %s \\times %d = %s$ K above %.0f °C.", f(R, 2), f(qq, 1), N*(N+1)/2, f(R*qq*float64(N*(N+1))/2, 3), Ts), "t_far")
	default:
		nx := 10 * edakit.Int(2, 20)
		ny := 10 * edakit.Int(2, 20)
		nz := edakit.Int(2, 20)
		N := nx * ny * nz
		links := (nx-1)*ny*nz + nx*(ny-1)*nz + nx*ny*(nz-1)
		ans := edakit.IntStr(N + 2*links)
		q := edakit.Sprintf("A 3-D finite-difference thermal grid has %d by %d by %d cells with a 7-point stencil, one unknown per cell. How many nonzero entries does the system matrix have?", nx, ny, nz)
		return num(q, ans, edakit.Sprintf("Diagonal: %d. Neighbour links: $%d + %d + %d = %d$, each giving two off-diagonal entries. Total $%d + 2 \\times %d = %s$.", N, (nx-1)*ny*nz, nx*(ny-1)*nz, nx*ny*(nz-1), links, N, links, ans), "nnz")
	}
}

// transient: thermal capacitance, time constant, step response.
func transient(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		W := step(5, 30, 0.5)
		H := step(5, 30, 0.5)
		t := float64(25 * edakit.Int(2, 32))
		V := W * H * t / 1e6 * 1e-6 // m^3
		C := 2330 * 700 * V * 1000  // mJ/K
		ans := f(C, 1)
		q := edakit.Sprintf("A silicon die is %s mm by %s mm and %.0f um thick. Silicon has density 2330 kg/m^3 and specific heat 700 J/(kg K). What is the die's thermal capacitance, in mJ/K? Round to 1 decimal place.", f(W, 1), f(H, 1), t)
		return num(q, ans, edakit.Sprintf("$V = %s$ mm^3 $= %.4g$ m^3. $C = \\rho c_p V = 2330 \\times 700 \\times %.4g = %s$ mJ/K.", f(W*H*t/1000, 2), V, V, ans), "c_mjk")
	case 2:
		C := float64(10 * edakit.Int(5, 200))
		R := step(0.05, 2, 0.05)
		ans := f(C*R, 2)
		q := edakit.Sprintf("A lumped thermal model has %.0f mJ/K of capacitance and %s K/W of resistance to ambient. What is its thermal time constant, in milliseconds? Round to 2 decimal places.", C, f(R, 2))
		return num(q, ans, edakit.Sprintf("$\\tau = RC = %s \\times %.0f$ mJ/K $= %s$ ms.", f(R, 2), C, ans), "tau_ms")
	default:
		P := float64(5 * edakit.Int(2, 40))
		R := step(0.05, 1, 0.05)
		tau := float64(10 * edakit.Int(2, 50))
		t := float64(5 * edakit.Int(1, 60))
		dT := P * R * (1 - math.Exp(-t/tau))
		ans := f(dT, 2)
		q := edakit.Sprintf("A power step of %.0f W is applied to a lumped thermal body with %s K/W to ambient and a time constant of %.0f ms. How far has its temperature risen after %.0f ms, in kelvin? Round to 2 decimal places.", P, f(R, 2), tau, t)
		return num(q, ans, edakit.Sprintf("Final rise $%.0f \\times %s = %s$ K. $\\Delta T = %s(1 - e^{-%.0f/%.0f}) = %s$ K.", P, f(R, 2), f(P*R, 2), f(P*R, 2), t, tau, ans), "dt_k")
	}
}

// management: TIM drop, required sink resistance, convective area.
func management(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		BLT := float64(5 * edakit.Int(2, 30))
		k := step(1, 10, 0.5)
		A := float64(50 * edakit.Int(2, 16))
		P := float64(10 * edakit.Int(3, 40))
		ans := f(P*BLT/(k*A), 2)
		q := edakit.Sprintf("A thermal interface material with conductivity %s W/(m K) has a %.0f um bond line over %.0f mm^2. Ignoring contact resistance, what temperature drop does it cause at %.0f W, in kelvin? Round to 2 decimal places.", f(k, 1), BLT, A, P)
		return num(q, ans, edakit.Sprintf("$R = %.0f/(%s \\times %.0f) = %s$ K/W; $\\Delta T = %.0f R = %s$ K.", BLT, f(k, 1), A, f(BLT/(k*A), 5), P, ans), "dt_k")
	case 2:
		for attempt := 0; ; attempt++ {
			P := float64(10 * edakit.Int(3, 30))
			Tj := float64(5 * edakit.Int(17, 22))
			Ta := float64(5 * edakit.Int(4, 10))
			jc := step(0.05, 0.3, 0.01)
			cs := step(0.02, 0.15, 0.01)
			if attempt >= generator.MaxDrawAttempts {
				P, Tj, Ta, jc, cs = 150, 95, 40, 0.1, 0.05
			}
			sa := (Tj-Ta)/P - jc - cs
			if sa <= 0.02 {
				continue
			}
			ans := f(sa, 3)
			q := edakit.Sprintf("A %.0f W part must keep its junction at or below %.0f °C in %.0f °C air. It has $\\theta_{jc} = %s$ K/W and the interface adds $\\theta_{cs} = %s$ K/W. What is the largest heat-sink resistance $\\theta_{sa}$ allowed, in K/W? Round to 3 decimal places.", P, Tj, Ta, f(jc, 2), f(cs, 2))
			return num(q, ans, edakit.Sprintf("$(%.0f - %.0f)/%.0f - %s - %s = %s$ K/W.", Tj, Ta, P, f(jc, 2), f(cs, 2), ans), "theta_sa")
		}
	default:
		h := float64(5 * edakit.Int(1, 40))
		th := step(0.1, 2, 0.05)
		ans := f(1/(h*th)*1e4, 1)
		q := edakit.Sprintf("A heat sink transfers heat to air with a heat-transfer coefficient of %.0f W/(m^2 K). Assuming fin efficiency of one, what wetted area gives a convective resistance of %s K/W, in cm^2? Round to 1 decimal place.", h, f(th, 2))
		return num(q, ans, edakit.Sprintf("$A = 1/(h\\theta) = 1/(%.0f \\times %s) = %s$ m^2 $= %s$ cm^2.", h, f(th, 2), f(1/(h*th), 5), ans), "area_cm2")
	}
}
