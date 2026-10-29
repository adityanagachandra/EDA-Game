package fab

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// wafers: wafer area, edge-exclusion loss and the area gain of a larger wafer.
func wafers(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		d := edakit.Pick([]int{100, 150, 200, 300})
		n := edakit.Scaled(ctx, 1, 25)
		r := float64(d) / 20 // radius in cm
		area := math.Pi * r * r * float64(n)
		ans := edakit.Num(area, 0)
		q := edakit.Sprintf("A lot holds %d silicon wafers of %d mm diameter. What is the total silicon area of the lot, in square centimetres? Round to the nearest whole number.", n, d)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Each wafer has radius %.1f cm, so area $\\pi r^2 = %.2f$ square cm. Times %d wafers gives %s square cm.", r, math.Pi*r*r, n, ans),
			map[string]string{"wafer_area_cm2": edakit.Num(math.Pi*r*r, 2)})
	case 2:
		d := edakit.Pick([]int{150, 200, 300, 450})
		e := edakit.Int(1, 6)
		R := float64(d) / 2
		lost := (1 - math.Pow((R-float64(e))/R, 2)) * 100
		ans := edakit.Num(lost, 2)
		q := edakit.Sprintf("A %d mm wafer carries an edge exclusion of %d mm: no die may sit closer than that to the edge. What percentage of the wafer area does the exclusion ring remove? Round to 2 decimal places.", d, e)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The usable disc has radius %.0f - %d = %.0f mm. The lost fraction is $1 - (%.0f/%.0f)^2$, which is %s%%.", R, e, R-float64(e), R-float64(e), R, ans),
			map[string]string{"usable_radius_mm": edakit.Num(R-float64(e), 0)})
	default:
		pairs := [][2]int{{100, 150}, {100, 200}, {150, 200}, {150, 300}, {200, 300}, {300, 450}, {200, 450}, {125, 200}}
		p := edakit.Pick(pairs)
		e := edakit.Int(0, 5)
		r1, r2 := float64(p[0])/2-float64(e), float64(p[1])/2-float64(e)
		ratio := (r2 * r2) / (r1 * r1)
		ans := edakit.Num(ratio, 3)
		q := edakit.Sprintf("A fab moves from %d mm to %d mm wafers and keeps a %d mm edge exclusion on both. By what factor does the usable area per wafer grow? Round to 3 decimal places.", p[0], p[1], e)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Usable radii are %.1f mm and %.1f mm. Area scales with radius squared: $(%.1f/%.1f)^2 = %s$. A fixed exclusion ring costs the smaller wafer proportionally more.", r1, r2, r2, r1, ans),
			map[string]string{"ratio": ans})
	}
}

// oxidation: Deal-Grove thermal oxidation.
func oxidation(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		x := 10 * edakit.Int(2, 60)
		if edakit.Int(0, 1) == 0 {
			ans := edakit.Num(0.44*float64(x), 1)
			q := edakit.Sprintf("A thermal oxide of %d nm is grown on bare silicon. How much silicon thickness does the growth consume, in nanometres? Use the consumption ratio 0.44. Round to 1 decimal place.", x)
			return edakit.Numeric(q, ans,
				edakit.Sprintf("Oxide grows partly into the silicon: each nanometre of oxide consumes 0.44 nm of silicon, so %d nm consumes %s nm.", x, ans),
				map[string]string{"consumed_nm": ans})
		}
		si := float64(x) * 0.44
		ans := edakit.Num(float64(x), 1)
		q := edakit.Sprintf("Thermal oxidation consumed %s nm of silicon. How thick is the oxide, in nanometres? Use the consumption ratio 0.44. Round to 1 decimal place.", edakit.Num(si, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Oxide thickness is consumed silicon divided by 0.44: %s / 0.44 = %s nm.", edakit.Num(si, 2), ans),
			map[string]string{"oxide_nm": ans})
	case 2:
		A := step(0.05, 0.30, 0.01)
		B := step(0.01, 0.30, 0.01)
		tau := step(0, 0.4, 0.05)
		t := step(0.5, 6, 0.5)
		x := (-A + math.Sqrt(A*A+4*B*(t+tau))) / 2
		ans := edakit.Num(x*1000, 0)
		q := edakit.Sprintf("Deal-Grove oxidation: $x^2 + A x = B (t + \\tau)$ with $A = %s$ um, $B = %s$ um^2/h and $\\tau = %s$ h. What oxide thickness grows in %s h, in nanometres? Round to the nearest nanometre.",
			edakit.Num(A, 2), edakit.Num(B, 2), edakit.Num(tau, 2), edakit.Num(t, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Solve the quadratic: $x = (-A + \\sqrt{A^2 + 4B(t+\\tau)})/2 = %s$ um, which is %s nm.", edakit.Num(x, 4), ans),
			map[string]string{"x_um": edakit.Num(x, 4)})
	default:
		A := step(0.05, 0.30, 0.01)
		B := step(0.02, 0.30, 0.01)
		tau := step(0, 0.3, 0.05)
		x := step(0.1, 0.8, 0.05) // um
		t := (x*x+A*x)/B - tau
		if t <= 0 {
			tau = 0
			t = (x*x + A*x) / B
		}
		ans := edakit.Num(t, 2)
		q := edakit.Sprintf("Deal-Grove oxidation with $A = %s$ um, $B = %s$ um^2/h and $\\tau = %s$ h. How long must the wafer oxidize to reach %s nm of oxide, in hours? Round to 2 decimal places.",
			edakit.Num(A, 2), edakit.Num(B, 2), edakit.Num(tau, 2), edakit.Num(x*1000, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Rearrange: $t = (x^2 + A x)/B - \\tau = (%s + %s)/%s - %s = %s$ h. The $x^2/B$ term is the parabolic, diffusion-limited part; $A x / B$ is the linear, reaction-limited part.",
				edakit.Num(x*x, 4), edakit.Num(A*x, 4), edakit.Num(B, 2), edakit.Num(tau, 2), ans),
			map[string]string{"t_h": ans})
	}
}

// diffusion: diffusion length, Arrhenius temperature dependence and thermal budget.
func diffusion(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		m, e := edakit.Int(1, 9), edakit.Pick([]int{14, 15})
		D := float64(m) * math.Pow(10, -float64(e))
		mins := 5 * edakit.Int(2, 24)
		L := 2 * math.Sqrt(D*float64(mins)*60) * 1e7
		ans := edakit.Num(L, 1)
		q := edakit.Sprintf("A dopant has diffusivity $D = %d \\times 10^{-%d}$ cm^2/s at the anneal temperature. What is its characteristic diffusion length $2\\sqrt{Dt}$ after a %d-minute anneal, in nanometres? Round to 1 decimal place.", m, e, mins)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$Dt = %.3g$ cm^2, $\\sqrt{Dt} = %.4g$ cm, and $2\\sqrt{Dt} = %s$ nm (1 cm = $10^7$ nm).", D*float64(mins)*60, math.Sqrt(D*float64(mins)*60), ans),
			map[string]string{"Dt_cm2": edakit.Sprintf("%.3g", D*float64(mins)*60)})
	case 2:
		Ea := step(3.0, 4.0, 0.1)
		t1 := 25 * edakit.Int(36, 42) // 900 to 1050 C
		t2 := t1 + 25*edakit.Int(1, 4)
		k := 8.617e-5
		T1, T2 := float64(t1)+273, float64(t2)+273
		r := math.Exp(Ea / k * (1/T1 - 1/T2))
		ans := edakit.Num(r, 2)
		q := edakit.Sprintf("A dopant's diffusivity follows $D = D_0 e^{-E_a/kT}$ with $E_a = %s$ eV ($k = 8.617 \\times 10^{-5}$ eV/K). By what factor does $D$ grow when the anneal moves from %d C to %d C? Use kelvin = Celsius + 273. Round to 2 decimal places.", edakit.Num(Ea, 1), t1, t2)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The ratio is $\\exp[(E_a/k)(1/T_1 - 1/T_2)]$ with $T_1 = %.0f$ K and $T_2 = %.0f$ K, giving %s.", T1, T2, ans),
			map[string]string{"ratio": ans})
	default:
		m1, m2 := edakit.Int(1, 9), edakit.Int(1, 9)
		D1 := float64(m1) * 1e-14
		D2 := float64(m2) * 1e-15
		s1, s2 := 10*edakit.Int(1, 9), 60*edakit.Int(1, 10)
		dt := D1*float64(s1) + D2*float64(s2)
		L := 2 * math.Sqrt(dt) * 1e7
		ans := edakit.Num(L, 1)
		q := edakit.Sprintf("A wafer sees a spike anneal of %d s at $D = %d \\times 10^{-14}$ cm^2/s and later a %d s anneal at $D = %d \\times 10^{-15}$ cm^2/s for the same dopant. What is the combined diffusion length $2\\sqrt{\\sum D_i t_i}$, in nanometres? Round to 1 decimal place.", s1, m1, s2, m2)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Thermal budgets add as $Dt$ products: $%.3g + %.3g = %.3g$ cm^2. Then $2\\sqrt{Dt} = %s$ nm.", D1*float64(s1), D2*float64(s2), dt, ans),
			map[string]string{"Dt_cm2": edakit.Sprintf("%.3g", dt)})
	}
}

// implant: Gaussian implanted profile.
func implant(ctx generator.GeneratorContext) generator.Problem {
	m, e := edakit.Int(1, 9), edakit.Pick([]int{13, 14, 15})
	Q := float64(m) * math.Pow(10, float64(e))
	dRp := float64(edakit.Int(5, 30))
	Rp := 5 * math.Round(dRp*step(2.5, 5, 0.5)/5)
	Np := Q / (math.Sqrt(2*math.Pi) * dRp * 1e-7)
	base := edakit.Sprintf("An implant delivers a dose of $%d \\times 10^{%d}$ cm^-2 with projected range $R_p = %.0f$ nm and straggle $\\Delta R_p = %.0f$ nm (Gaussian profile).", m, e, Rp, dRp)
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(Np/1e18, 2)
		q := base + " What is the peak concentration, in units of $10^{18}$ cm^-3? Round to 2 decimal places."
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$N_p = Q/(\\sqrt{2\\pi}\\,\\Delta R_p)$ with $\\Delta R_p = %.0f \\times 10^{-7}$ cm gives %.4g cm^-3, which is %s in units of $10^{18}$.", dRp, Np, ans),
			map[string]string{"Np": edakit.Sprintf("%.4g", Np)})
	case 2:
		kk := step(0.5, 2.5, 0.25)
		sign := edakit.Pick([]float64{-1, 1})
		x := Rp + sign*kk*dRp
		if x < 0 {
			x = Rp + kk*dRp
		}
		N := Np * math.Exp(-(x-Rp)*(x-Rp)/(2*dRp*dRp))
		ans := edakit.Num(N/1e18, 2)
		q := base + edakit.Sprintf(" What is the concentration at a depth of %s nm, in units of $10^{18}$ cm^-3? Round to 2 decimal places.", edakit.Num(x, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Peak $N_p = %.4g$ cm^-3. The depth is %s straggles from $R_p$, so $N = N_p e^{-(x-R_p)^2/2\\Delta R_p^2} = %.4g$ cm^-3, which is %s.", Np, edakit.Num(math.Abs(x-Rp)/dRp, 2), N, ans),
			map[string]string{"Np": edakit.Sprintf("%.4g", Np)})
	default:
		bm, be := edakit.Int(1, 9), edakit.Pick([]int{15, 16, 17})
		for float64(bm)*math.Pow(10, float64(be)) >= Np/2 {
			be--
		}
		NB := float64(bm) * math.Pow(10, float64(be))
		xj := Rp + dRp*math.Sqrt(2*math.Log(Np/NB))
		ans := edakit.Num(xj, 1)
		q := base + edakit.Sprintf(" The substrate has a background doping of $%d \\times 10^{%d}$ cm^-3 of the opposite type. The junction lies on the deep side of the peak. At what depth below the surface is it, in nanometres? Round to 1 decimal place.", bm, be)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Set $N(x_j) = N_B$: $x_j = R_p + \\Delta R_p\\sqrt{2\\ln(N_p/N_B)}$ with $N_p = %.4g$ cm^-3, so $x_j = %s$ nm.", Np, ans),
			map[string]string{"Np": edakit.Sprintf("%.4g", Np)})
	}
}

// epitaxyCompute: lattice mismatch of SiGe on Si by Vegard's law.
func epitaxyCompute(ctx generator.GeneratorContext) generator.Problem {
	x := step(0.10, 0.60, 0.05)
	aSi, aGe := 5.431, 5.658
	a := aSi + (aGe-aSi)*x
	f := (a - aSi) / aSi * 100
	c, d := numOptions(f, 2, "%", (a-aSi)/aGe*100, (aGe-aSi)/aSi*100, f*2, f/2)
	q := edakit.Sprintf("A $\\mathrm{Si}_{1-x}\\mathrm{Ge}_x$ layer with $x = %s$ is grown on a silicon substrate. Take $a_{Si} = 5.431$ and $a_{Ge} = 5.658$ angstrom and Vegard's law. What is the lattice mismatch $(a_{film} - a_{sub})/a_{sub}$?", edakit.Num(x, 2))
	return edakit.Choice(q, c,
		edakit.Sprintf("Vegard's law interpolates linearly: $a = 5.431 + 0.227 \\times %s = %s$ angstrom. The mismatch is $%s/5.431$ = %s. A pseudomorphic film takes the substrate's in-plane spacing and is held in compression.", edakit.Num(x, 2), edakit.Num(a, 4), edakit.Num(a-aSi, 4), c),
		d...)
}

// depositionCompute: ALD cycle counts and conformal step coverage.
func depositionCompute(ctx generator.GeneratorContext) generator.Problem {
	if edakit.Int(0, 1) == 0 {
		T := step(1.0, 6.0, 0.1)
		g := step(0.05, 0.15, 0.01)
		n := int(math.Ceil(T/g - 1e-9))
		q := edakit.Sprintf("An ALD process grows %s nm per cycle. How many cycles are needed for at least %s nm of film?", edakit.Num(g, 2), edakit.Num(T, 1))
		c := edakit.IntStr(n)
		return edakit.Choice(q, c,
			edakit.Sprintf("ALD is self-limiting, so thickness is cycles times growth per cycle: $%s / %s = %s$, rounded up to %d whole cycles.", edakit.Num(T, 1), edakit.Num(g, 2), edakit.Num(T/g, 2), n),
			edakit.Distinct(n, 3, int(math.Round(T*g*100)), int(math.Floor(T/g))-1, n+edakit.Int(3, 9))...)
	}
	top := step(10, 40, 1)
	cov := step(0.3, 0.95, 0.05)
	bottom := top * cov
	q := edakit.Sprintf("A film deposits %s nm on the field above a trench and %s nm on the trench bottom. What is the bottom step coverage?", edakit.Num(top, 0), edakit.Num(bottom, 2))
	c, d := numOptions(cov*100, 0, "%", 100/cov, (1-cov)*100, cov*50)
	return edakit.Choice(q, c,
		edakit.Sprintf("Step coverage is bottom thickness over field thickness: %s / %s = %s. ALD approaches 100%%; directional PVD falls well short in deep trenches.", edakit.Num(bottom, 2), edakit.Num(top, 0), c),
		d...)
}

// etch: anisotropy, etch time with over-etch, and mask erosion through selectivity.
func etch(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		v := float64(10 * edakit.Int(5, 40))
		l := float64(edakit.Int(1, 60))
		ans := edakit.Num(1-l/v, 3)
		q := edakit.Sprintf("A plasma etch removes material vertically at %.0f nm/min and laterally at %.0f nm/min. What is its anisotropy $A = 1 - r_{lat}/r_{vert}$? Round to 3 decimal places.", v, l)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$A = 1 - %.0f/%.0f = %s$. A wet isotropic etch has $A = 0$; a perfectly directional etch has $A = 1$.", l, v, ans),
			map[string]string{"A": ans})
	case 2:
		h := float64(10 * edakit.Int(5, 60))
		r := float64(10 * edakit.Int(5, 50))
		oe := 5 * edakit.Int(2, 10)
		t := h / r * 60 * (1 + float64(oe)/100)
		ans := edakit.Num(t, 1)
		q := edakit.Sprintf("A %.0f nm film etches at %.0f nm/min. The recipe adds %d%% over-etch to clear thickness variation across the wafer. What is the total etch time, in seconds? Round to 1 decimal place.", h, r, oe)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Main etch takes $%.0f/%.0f = %s$ min, which is %s s. Over-etch multiplies by %s, giving %s s.", h, r, edakit.Num(h/r, 3), edakit.Num(h/r*60, 1), edakit.Num(1+float64(oe)/100, 2), ans),
			map[string]string{"main_s": edakit.Num(h/r*60, 1)})
	default:
		h := float64(10 * edakit.Int(5, 40))
		oe := 10 * edakit.Int(1, 5)
		S := float64(edakit.Int(3, 30))
		loss := h * (1 + float64(oe)/100) / S
		ans := edakit.Num(loss, 1)
		q := edakit.Sprintf("A %.0f nm film is etched with %d%% over-etch through a mask. The etch selectivity of film to mask is %.0f:1. How much mask thickness is consumed, in nanometres? Round to 1 decimal place.", h, oe, S)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The etch runs long enough to remove %s nm of film. The mask erodes %s times slower: %s / %.0f = %s nm.", edakit.Num(h*(1+float64(oe)/100), 1), edakit.Num(S, 0), edakit.Num(h*(1+float64(oe)/100), 1), S, ans),
			map[string]string{"film_equiv_nm": edakit.Num(h*(1+float64(oe)/100), 1)})
	}
}

// cmp: Preston scaling, density-dependent planarization and dummy fill.
func cmpGen(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		rr := float64(10 * edakit.Int(10, 50))
		p1, v1 := float64(edakit.Int(10, 30)), step(0.5, 2.0, 0.1)
		p2, v2 := float64(edakit.Int(10, 30)), step(0.5, 2.0, 0.1)
		r2 := rr * p2 * v2 / (p1 * v1)
		ans := edakit.Num(r2, 1)
		q := edakit.Sprintf("A CMP step removes %.0f nm/min at %.0f kPa down-force and %s m/s pad velocity. By Preston's law, what removal rate is expected at %.0f kPa and %s m/s, in nm/min? Round to 1 decimal place.", rr, p1, edakit.Num(v1, 1), p2, edakit.Num(v2, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Preston: rate $= K_p P v$, so the rate scales by $(%.0f \\times %s)/(%.0f \\times %s)$, giving %s nm/min.", p2, edakit.Num(v2, 1), p1, edakit.Num(v1, 1), ans),
			map[string]string{"ratio": edakit.Num(p2*v2/(p1*v1), 4)})
	case 2:
		K := float64(10 * edakit.Int(10, 40))
		rho := step(0.2, 0.9, 0.05)
		h := float64(50 * edakit.Int(2, 12))
		t := h * rho / K
		ans := edakit.Num(t, 2)
		q := edakit.Sprintf("In the density model of oxide CMP, raised features polish at the blanket rate divided by the local pattern density. The blanket rate is %.0f nm/min. How long does it take to remove a %.0f nm step in a region of pattern density %s, in minutes? Round to 2 decimal places.", K, h, edakit.Num(rho, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The raised area carries all the pad pressure, so its rate is $%.0f/%s = %s$ nm/min. The step clears in $%.0f/%s = %s$ min. Denser regions clear more slowly, which is why density rules exist.", K, edakit.Num(rho, 2), edakit.Num(K/rho, 1), h, edakit.Num(K/rho, 1), ans),
			map[string]string{"rate": edakit.Num(K/rho, 1)})
	default:
		W := 10 * edakit.Int(2, 10)
		dmin := 5 * edakit.Int(4, 8)
		have := float64(W*W) * float64(edakit.Int(2, dmin/5*5-5)) / 100
		need := float64(dmin)/100*float64(W*W) - have
		ans := edakit.Num(need, 0)
		q := edakit.Sprintf("A density rule requires at least %d%% metal in every %d um by %d um window. One window holds %s um^2 of drawn metal. How much dummy fill must be added to that window, in square micrometres? Round to the nearest whole number.", dmin, W, W, edakit.Num(have, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The window is %d um^2, so the rule needs %s um^2 of metal. Subtract the %s um^2 present: %s um^2 of fill.", W*W, edakit.Num(float64(dmin)/100*float64(W*W), 0), edakit.Num(have, 0), ans),
			map[string]string{"window_um2": edakit.IntStr(W * W)})
	}
}
