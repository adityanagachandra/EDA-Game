package dev

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// dennard asks for a constant-field scaling factor (KP 1), a Moore's-law density projection
// (KP 2) or the power-density growth once voltage stops scaling (KP 3).
func dennard(ctx generator.GeneratorContext) generator.Problem {
	s := pickF(0.7, 0.7, 0.75, 0.8)
	g := edakit.Int(1, 4)
	switch kpOf(ctx) {
	case 1:
		type qty struct {
			name string
			pow  float64
			why  string
		}
		qs := []qty{
			{"supply voltage", 1, "voltage scales with dimensions to keep the field constant"},
			{"gate delay", 1, "delay $CV/I$ scales as $(s)(s)/(s) = s$"},
			{"power per circuit", 2, "power $VI$ scales as $s \\cdot s = s^2$"},
			{"area per circuit", 2, "area is a product of two linear dimensions, $s^2$"},
			{"power density", 0, "power per circuit $s^2$ divided by area per circuit $s^2$ is constant"},
			{"switching energy $CV^2$ per gate", 3, "energy $CV^2$ scales as $s \\cdot s^2 = s^3$"},
		}
		x := edakit.Pick(qs)
		f := math.Pow(s, x.pow*float64(g))
		ans, rt := sigAnswer(f, 3)
		q := edakit.Sprintf("Under ideal constant-field (Dennard) scaling, every linear dimension and the supply voltage shrink by a factor %s per generation, and doping rises to match. After %d generation(s), what is the new %s divided by the old? %s", edakit.Num(s, 2), g, x.name, rt)
		return edakit.Numeric(q, ans, edakit.Sprintf("Per generation, %s. Over %d generation(s) the factor is $%s^{%s} = %s$.", x.why, g, edakit.Num(s, 2), edakit.Num(x.pow*float64(g), 0), ans), map[string]string{"factor": ans})
	case 2:
		if edakit.Int(0, 1) == 0 {
			n0 := pickF(0.5, 1, 2, 5, 10)
			per := pickF(2, 2, 2.5, 3)
			yrs := float64(edakit.Int(2, 10))
			n := n0 * math.Pow(2, yrs/per)
			ans, rt := sigAnswer(n, 3)
			q := edakit.Sprintf("A product line ships a chip with %s billion transistors. If transistor count doubles every %s years, how many billion transistors does the line's chip have %d years later? %s", edakit.Num(n0, 1), edakit.Num(per, 1), int(yrs), rt)
			return edakit.Numeric(q, ans, edakit.Sprintf("$N = N_0 \\cdot 2^{t/T} = %s \\times 2^{%d/%s} = %s$ billion.", edakit.Num(n0, 1), int(yrs), edakit.Num(per, 1), ans), map[string]string{"doublings": edakit.Num(yrs/per, 3)})
		}
		d := math.Pow(1/(s*s), float64(g))
		ans, rt := sigAnswer(d, 3)
		q := edakit.Sprintf("Each process generation shrinks linear dimensions by a factor %s. By what factor does transistor density grow over %d generation(s)? %s", edakit.Num(s, 2), g, rt)
		return edakit.Numeric(q, ans, edakit.Sprintf("Density scales as $1/s^2 = %.4f$ per generation, so over %d generation(s) it grows by %s. A factor of 0.7 gives about 2x per generation, the scaling behind Moore's law.", 1/(s*s), g, ans), map[string]string{"per_gen": edakit.Num(1/(s*s), 4)})
	default:
		pd := math.Pow(1/(s*s), float64(g))
		ans, rt := sigAnswer(pd, 3)
		q := edakit.Sprintf("Dimensions shrink by %s per generation, but the supply voltage stays fixed because the threshold cannot fall further. Assume capacitance per gate scales by %s, clock frequency rises by $1/%s$, and density by $1/%s^2$. By what factor does power density grow after %d generation(s)? %s", edakit.Num(s, 2), edakit.Num(s, 2), edakit.Num(s, 2), edakit.Num(s, 2), g, rt)
		return edakit.Numeric(q, ans, edakit.Sprintf("Power per gate $CV^2 f$ scales as $s \\cdot 1 \\cdot (1/s) = 1$. Density grows by $1/s^2$, so power density grows by $(1/%s^2)^{%d} = %s$. This is the power wall.", edakit.Num(s, 2), g, ans), map[string]string{"factor": ans})
	}
}

// powerWall uses generalized scaling (dimensions 1/kappa, voltage 1/s, long-channel current)
// and asks for the delay factor (KP 1), the power-density factor (KP 2) or the power-density
// factor at fixed frequency (KP 3).
func powerWall(ctx generator.GeneratorContext) generator.Problem {
	k := pickF(1.4, 1.5, 2)
	s := pickF(1, 1.1, 1.2, 1.3)
	stem := edakit.Sprintf("In generalized scaling, all dimensions including $t_{ox}$ shrink by $1/\\kappa$ and the supply voltage by $1/s$, with $\\kappa = %s$ and $s = %s$. Use the long-channel current $I \\propto (W/L) C_{ox} V^2$ and gate capacitance $C \\propto WL/t_{ox}$.", edakit.Num(k, 1), edakit.Num(s, 1))
	switch kpOf(ctx) {
	case 1:
		f := s / (k * k)
		ans := edakit.Num(f, 3)
		q := stem + " By what factor does the gate delay $CV/I$ change? Round to 3 decimal places."
		return edakit.Numeric(q, ans, edakit.Sprintf("$C \\to C/\\kappa$, $V \\to V/s$, $I \\to I\\kappa/s^2$. Delay $CV/I$ scales as $(1/\\kappa)(1/s)/(\\kappa/s^2) = s/\\kappa^2 = %s$.", ans), map[string]string{"delay": ans})
	case 2:
		f := math.Pow(k/s, 3)
		ans, rt := sigAnswer(f, 3)
		q := stem + " Frequency follows the delay. By what factor does power density change? " + rt
		return edakit.Numeric(q, ans, edakit.Sprintf("Power per device $VI \\propto (1/s)(\\kappa/s^2) = \\kappa/s^3$, and density rises by $\\kappa^2$, so power density scales as $\\kappa^3/s^3 = %s$. Only $s = \\kappa$ (Dennard) keeps it constant.", ans), map[string]string{"pd": ans})
	default:
		f := k / (s * s)
		ans, rt := sigAnswer(f, 3)
		q := stem + " Now the clock frequency is held fixed instead of following the delay. By what factor does dynamic power density $C V^2 f$ per unit area change? " + rt
		return edakit.Numeric(q, ans, edakit.Sprintf("Per device $CV^2 f \\propto (1/\\kappa)(1/s^2)(1)$, and density rises by $\\kappa^2$, so power density scales as $\\kappa/s^2 = %s$. Without voltage scaling, even a frozen clock does not hold power density flat.", ans), map[string]string{"pd": ans})
	}
}

// vtFloor asks for the lowest threshold that meets an off-current target (KP 1), the same at a
// hot corner (KP 2) or the on/off ratio in decades (KP 3).
func vtFloor(ctx generator.GeneratorContext) generator.Problem {
	s := float64(edakit.Pick([]int{65, 70, 75, 80, 90}))
	ivt := float64(edakit.Pick([]int{100, 200, 300, 500}))      // nA/um at Vgs = Vt
	ioff := float64(edakit.Pick([]int{10, 30, 100, 300, 1000})) // pA/um
	switch kpOf(ctx) {
	case 1:
		v := s * math.Log10(ivt*1000/ioff)
		ans := edakit.Num(v, 1)
		q := edakit.Sprintf("A transistor carries %d nA/µm at $V_{GS} = V_t$ (a constant-current threshold definition), and its subthreshold swing is %d mV/decade. The off-current at $V_{GS} = 0$ must not exceed %d pA/µm. What is the lowest allowed threshold, in mV? Round to 1 decimal place.", int(ivt), int(s), int(ioff))
		return edakit.Numeric(q, ans, edakit.Sprintf("$I_{off} = I_{V_t} 10^{-V_t/S}$, so $V_t \\ge S\\log_{10}(I_{V_t}/I_{off}) = %d \\times \\log_{10}(%s) = %s$ mV.", int(s), edakit.Num(ivt*1000/ioff, 0), ans), map[string]string{"decades": edakit.Num(math.Log10(ivt*1000/ioff), 4)})
	case 2:
		th := float64(edakit.Pick([]int{358, 373, 398}))
		sh := s * th / 300
		v := sh * math.Log10(ivt*1000/ioff)
		ans := edakit.Num(v, 1)
		q := edakit.Sprintf("At 300 K a transistor has a swing of %d mV/decade and carries %d nA/µm at $V_{GS} = V_t$. Keeping $n$ and the current at threshold fixed, the swing scales with absolute temperature. What is the lowest threshold that keeps the off-current at or below %d pA/µm at %d K, in mV? Round to 1 decimal place.", int(s), int(ivt), int(ioff), int(th))
		return edakit.Numeric(q, ans, edakit.Sprintf("$S(%d\\ \\text{K}) = %d \\times %d/300 = %.3f$ mV/decade, so $V_t \\ge %.3f \\times \\log_{10}(%s) = %s$ mV. The hot corner sets the floor.", int(th), int(s), int(th), sh, sh, edakit.Num(ivt*1000/ioff, 0), ans), map[string]string{"S_hot": edakit.Num(sh, 3)})
	default:
		vt := float64(10 * edakit.Int(20, 45))
		ion := float64(100 * edakit.Int(5, 15)) // uA/um
		d := math.Log10(ion*1000/ivt) + vt/s
		ans := edakit.Num(d, 2)
		q := edakit.Sprintf("A transistor has $V_t = %d$ mV, a swing of %d mV/decade, %d nA/µm at $V_{GS} = V_t$ and an on-current of %d µA/µm at $V_{GS} = V_{DD}$. How many decades separate $I_{on}$ from $I_{off}$ (the current at $V_{GS} = 0$)? Round to 2 decimal places.", int(vt), int(s), int(ivt), int(ion))
		return edakit.Numeric(q, ans, edakit.Sprintf("$I_{off} = %d\\ \\text{nA/µm} \\times 10^{-%d/%d}$. So $\\log_{10}(I_{on}/I_{off}) = \\log_{10}(%d/%s) + %d/%d = %s$ decades.", int(ivt), int(vt), int(s), int(ion), edakit.Num(ivt/1000, 1), int(vt), int(s), ans), map[string]string{"decades": ans})
	}
}
