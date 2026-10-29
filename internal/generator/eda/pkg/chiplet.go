package pkg

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// disaggregation: monolithic yield, chiplet yield with overhead, silicon per good system.
func disaggregation(ctx generator.GeneratorContext) generator.Problem {
	A := float64(10 * edakit.Int(10, 80))
	D := step(0.05, 0.3, 0.01)
	k := edakit.Int(2, 8)
	o := step(0.03, 0.2, 0.01)
	Ym := math.Exp(-A / 100 * D)
	a := A * (1 + o) / float64(k)
	Yc := math.Exp(-a / 100 * D)
	switch kpOf(ctx) {
	case 1:
		ans := f(Ym*100, 2)
		q := edakit.Sprintf("A %.0f mm^2 monolithic die is made in a process with %s killer defects per cm^2, randomly scattered. Using the Poisson model $Y = e^{-AD_0}$, what is its yield, as a percentage? Round to 2 decimal places.", A, f(D, 2))
		return num(q, ans, edakit.Sprintf("$A = %s$ cm^2, $AD_0 = %s$, $Y = e^{-%s} = %s\\%%$.", f(A/100, 2), f(A/100*D, 4), f(A/100*D, 4), ans), "yield_pct")
	case 2:
		ans := f(Yc*100, 2)
		q := edakit.Sprintf("A %.0f mm^2 design is split into %d equal chiplets. Each chiplet needs an extra %.0f%% of its logic area for die-to-die interfaces. With %s defects per cm^2 and the Poisson model, what is the yield of one chiplet, as a percentage? Round to 2 decimal places.", A, k, o*100, f(D, 2))
		return num(q, ans, edakit.Sprintf("Chiplet area $%.0f \\times %s/%d = %s$ mm^2. $Y = e^{-%s \\times %s} = %s\\%%$.", A, f(1+o, 2), k, f(a, 2), f(a/100, 4), f(D, 2), ans), "yield_pct")
	default:
		ratio := Yc / (Ym * (1 + o))
		ans := f(ratio, 2)
		q := edakit.Sprintf("A %.0f mm^2 design at %s defects per cm^2 can be built monolithically or as %d tested chiplets, each carrying %.0f%% interface overhead. Using Poisson yields, what is the ratio of manufactured silicon per good system, monolithic over chiplet? Ignore assembly loss. Round to 2 decimal places.", A, f(D, 2), k, o*100)
		return num(q, ans, edakit.Sprintf("Monolithic: $%.0f/%s = %s$ mm^2. Chiplets: $%s/%s = %s$ mm^2. Ratio %s.", A, f(Ym, 4), f(A/Ym, 1), f(A*(1+o), 1), f(Yc, 4), f(A*(1+o)/Yc, 1), ans), "ratio")
	}
}

// ucie: module bandwidth, shoreline bandwidth, module count for a byte target.
func ucie(ctx generator.GeneratorContext) generator.Problem {
	adv := edakit.Int(0, 1) == 1
	lanes := 16
	kind := "standard-package"
	if adv {
		lanes = 64
		kind = "advanced-package"
	}
	rate := float64(edakit.Pick([]int{4, 8, 12, 16, 24, 32}))
	B := float64(lanes) * rate
	switch kpOf(ctx) {
	case 1:
		ans := f(B, 0)
		q := edakit.Sprintf("A UCIe %s module has %d data lanes per direction, each at %.0f GT/s (one bit per transfer). What is its raw bandwidth per direction, in Gb/s?", kind, lanes, rate)
		return num(q, ans, edakit.Sprintf("$%d \\times %.0f = %s$ Gb/s each way, before protocol overhead.", lanes, rate, ans), "gbps")
	case 2:
		w := float64(50 * edakit.Int(6, 24))
		E := step(3, 15, 0.5)
		n := fl(E*1000, w)
		ans := f(float64(n)*B/1000, 2)
		q := edakit.Sprintf("Suppose each %s module of %d lanes at %.0f GT/s occupies %.0f um of die edge, with modules placed side by side in one row. What raw bandwidth per direction can a %s mm edge deliver, in Tb/s? Round to 2 decimal places.", kind, lanes, rate, w, f(E, 1))
		return num(q, ans, edakit.Sprintf("Modules: $\\lfloor %.0f/%.0f \\rfloor = %d$, each %s Gb/s. Total $%d \\times %s = %s$ Gb/s, or %s Tb/s.", E*1000, w, n, f(B, 0), n, f(B, 0), f(float64(n)*B, 0), ans), "tbps")
	default:
		T := step(0.2, 4, 0.1)
		ans := edakit.IntStr(cl(8000*T, B))
		q := edakit.Sprintf("An accelerator needs %s TB/s (terabytes per second) per direction to its neighbour. Using UCIe %s modules of %d lanes at %.0f GT/s, and ignoring overhead, how many modules are needed?", f(T, 1), kind, lanes, rate)
		return num(q, ans, edakit.Sprintf("$%s$ TB/s is %s Gb/s. Each module carries %s Gb/s. $\\lceil %s/%s \\rceil = %s$ modules.", f(T, 1), f(8000*T, 0), f(B, 0), f(8000*T, 0), f(B, 0), ans), "modules")
	}
}

// energyPerBit: link power, on- versus off-package saving, energy from capacitance.
func energyPerBit(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		Eb := step(0.1, 2, 0.05)
		BW := step(1, 40, 0.5)
		both := edakit.Int(0, 1) == 1
		P := Eb * BW
		dir := "in one direction"
		if both {
			P *= 2
			dir = "in each of the two directions, both active"
		}
		ans := f(P, 2)
		q := edakit.Sprintf("A die-to-die link spends %s pJ per bit and carries %s Tb/s %s. What power does it draw, in watts? Round to 2 decimal places.", f(Eb, 2), f(BW, 1), dir)
		return num(q, ans, edakit.Sprintf("1 pJ/bit at 1 Tb/s is 1 W. $%s \\times %s$ per direction gives %s W in total.", f(Eb, 2), f(BW, 1), ans), "watts")
	case 2:
		E1 := step(3, 10, 0.5)
		E2 := step(0.2, 0.8, 0.05)
		BW := step(1, 20, 0.5)
		ans := f((E1-E2)*BW, 2)
		q := edakit.Sprintf("Moving %s Tb/s between two chips costs %s pJ/bit over off-package SerDes, or %s pJ/bit if the chips become chiplets on one package. How many watts does the package option save? Round to 2 decimal places.", f(BW, 1), f(E1, 1), f(E2, 2))
		return num(q, ans, edakit.Sprintf("$(%s - %s) \\times %s = %s$ W.", f(E1, 1), f(E2, 2), f(BW, 1), ans), "saving_w")
	default:
		C := float64(edakit.Int(30, 400))
		V := step(0.4, 1.0, 0.05)
		e := C * V * V / 4
		ans := f(e, 2)
		q := edakit.Sprintf("A full-swing, unterminated die-to-die wire has %.0f fF of total capacitance and swings %s V. Each 0-to-1 transition draws $CV^2$ from the supply, and random data makes a 0-to-1 transition on one quarter of bits. What is the average energy per bit, in fJ? Round to 2 decimal places.", C, f(V, 2))
		return num(q, ans, edakit.Sprintf("$CV^2/4 = %.0f \\times %s^2/4 = %s$ fJ per bit.", C, f(V, 2), ans), "fj_per_bit")
	}
}

// packageBoard: power ball count, target impedance, parallel paths for inductance.
func packageBoard(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		I := float64(edakit.Int(10, 300))
		Ib := step(0.3, 1.0, 0.1)
		ans := edakit.IntStr(cl(I*10, Ib*10))
		q := edakit.Sprintf("A core rail draws %.0f A, and each package ball may carry at most %s A. How many supply balls are needed for this rail, not counting the ground return?", I, f(Ib, 1))
		return num(q, ans, edakit.Sprintf("$\\lceil %.0f/%s \\rceil = %s$ balls, with as many ground balls again.", I, f(Ib, 1), ans), "balls")
	case 2:
		V := step(0.6, 1.2, 0.05)
		r := float64(edakit.Pick([]int{1, 2, 3, 5}))
		dI := float64(edakit.Int(5, 150))
		Z := V * r / 100 / dI * 1000
		ans := f(Z, 3)
		q := edakit.Sprintf("A %s V rail may ripple by %.0f%% during a %.0f A load step. What is the target impedance of its power delivery network, in milliohms? Round to 3 decimal places.", f(V, 2), r, dI)
		return num(q, ans, edakit.Sprintf("$Z_t = V r/\\Delta I = %s \\times %s/%.0f = %s$ milliohm.", f(V, 2), f(r/100, 2), dI, ans), "z_mohm")
	default:
		L1 := float64(10 * edakit.Int(10, 80))
		Lt := step(1, 20, 0.5)
		ans := edakit.IntStr(cl(L1*10, Lt*10))
		q := edakit.Sprintf("One ball-and-via power path has %.0f pH of inductance. Ignoring mutual inductance, how many identical paths in parallel are needed to bring the inductance to at most %s pH?", L1, f(Lt, 1))
		return num(q, ans, edakit.Sprintf("$N = \\lceil %.0f/%s \\rceil = %s$ paths.", L1, f(Lt, 1), ans), "paths")
	}
}
