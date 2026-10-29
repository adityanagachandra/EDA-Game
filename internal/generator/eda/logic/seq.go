package logic

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func registerSeq(reg *generator.Registry) {
	reg.Register("logic.seq.latch_ff", retryTies(latchFF))
	reg.Register("logic.seq.registers", retryTies(registers))
	reg.Register("logic.seq.flip_flop_timing", retryTies(ffTiming))
	reg.Register("logic.seq.max_frequency", retryTies(maxFrequency))
	reg.Register("logic.seq.metastability", retryTies(metastability))
	reg.Register("logic.seq.synchronizer_mtbf", retryTies(synchronizerMTBF))
	reg.Register("logic.seq.synchronous_design", &edakit.Bank{Items: syncDesignItems})
	reg.Register("logic.seq.cdc", &edakit.Bank{Items: cdcItems})
	reg.Register("logic.seq.counters", retryTies(counters))
}

func jk(q, j, k int) int { return (j & (1 - q)) | ((1 - k) & q) }

func latchFF(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		q := edakit.Int(0, 1)
		switch edakit.Int(0, 2) {
		case 0:
			j, k := edakit.Int(0, 1), edakit.Int(0, 1)
			ans := jk(q, j, k)
			return edakit.Numeric(edakit.Sprintf("A JK flip-flop holds $Q = %d$. At the next clock edge $J = %d$ and $K = %d$. What is $Q$ after the edge?", q, j, k), edakit.IntStr(ans),
				edakit.Sprintf("$Q^+ = J\\bar{Q} + \\bar{K}Q = %d \\cdot %d + %d \\cdot %d = %d$.", j, 1-q, 1-k, q, ans), map[string]string{"next": edakit.IntStr(ans)})
		case 1:
			t := edakit.Int(0, 1)
			ans := t ^ q
			return edakit.Numeric(edakit.Sprintf("A T flip-flop holds $Q = %d$. At the next clock edge $T = %d$. What is $Q$ after the edge?", q, t), edakit.IntStr(ans),
				edakit.Sprintf("$Q^+ = T \\oplus Q = %d \\oplus %d = %d$.", t, q, ans), map[string]string{"next": edakit.IntStr(ans)})
		default:
			s, r := edakit.Int(0, 1), edakit.Int(0, 1)
			if s == 1 && r == 1 {
				r = 0
			}
			ans := s | ((1 - r) & q)
			return edakit.Numeric(edakit.Sprintf("A clocked SR flip-flop holds $Q = %d$. At the next clock edge $S = %d$ and $R = %d$. What is $Q$ after the edge?", q, s, r), edakit.IntStr(ans),
				edakit.Sprintf("$Q^+ = S + \\bar{R}Q = %d + %d \\cdot %d = %d$.", s, 1-r, q, ans), map[string]string{"next": edakit.IntStr(ans)})
		}
	default:
		count := ctx.KPOrdinal == 3 || (ctx.KPOrdinal == 0 && edakit.Int(0, 1) == 0)
		L := edakit.Scaled(ctx, 4, 8)
		q := edakit.Int(0, 1)
		q0 := q
		var steps []string
		ones := 0
		var trace []int
		for i := 0; i < L; i++ {
			j, k := edakit.Int(0, 1), edakit.Int(0, 1)
			q = jk(q, j, k)
			trace = append(trace, q)
			ones += q
			steps = append(steps, edakit.Sprintf("(%d, %d)", j, k))
		}
		seq := strings.Join(steps, ", ")
		if !count {
			return edakit.Numeric(edakit.Sprintf("A JK flip-flop starts at $Q = %d$. On %d successive clock edges it sees $(J, K)$ = %s. What is $Q$ after the last edge?", q0, L, seq), edakit.IntStr(q),
				"Applying the characteristic equation once per edge gives the states "+joinInts(trace, ", ")+".", map[string]string{"final": edakit.IntStr(q)})
		}
		return edakit.Numeric(edakit.Sprintf("A JK flip-flop starts at $Q = %d$. On %d successive clock edges it sees $(J, K)$ = %s. After how many of those edges is $Q = 1$?", q0, L, seq), edakit.IntStr(ones),
			"The states after each edge are "+joinInts(trace, ", ")+edakit.Sprintf(", so Q is 1 after %d edges.", ones), map[string]string{"final": edakit.IntStr(q)})
	}
}

func registers(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		w := 8
		mask := uint64(1)<<uint(w) - 1
		r := uint64(edakit.Int(1, 255))
		r0 := r
		k := edakit.Scaled(ctx, 3, 5)
		var bitsIn []int
		left := edakit.Int(0, 1) == 0
		for i := 0; i < k; i++ {
			b := edakit.Int(0, 1)
			bitsIn = append(bitsIn, b)
			if left {
				r = ((r << 1) | uint64(b)) & mask
			} else {
				r = (r >> 1) | uint64(b)<<uint(w-1)
			}
		}
		dir := "left, with the serial input entering at the least significant bit"
		if !left {
			dir = "right, with the serial input entering at the most significant bit"
		}
		return edakit.Numeric(edakit.Sprintf("An 8-bit shift register holds %s (%d). It shifts %s, and on %d successive clock edges the serial input is %s. What value does it hold afterwards, as an unsigned decimal integer?", bin(r0, w), r0, dir, k, joinInts(bitsIn, ", ")), edakit.IntStr(int(r)),
			edakit.Sprintf("Each edge moves every bit one place and inserts the serial bit; the final pattern is %s = %d.", bin(r, w), r), map[string]string{"pattern": bin(r, w)})
	case 2:
		type rtl struct {
			code string
			f    func(a, b int) (int, int)
		}
		forms := []rtl{
			{"A <= A + B; B <= A;", func(a, b int) (int, int) { return (a + b) % 256, a }},
			{"A <= B; B <= A + B;", func(a, b int) (int, int) { return b, (a + b) % 256 }},
			{"A <= B; B <= A;", func(a, b int) (int, int) { return b, a }},
			{"A <= A + 2*B; B <= A - B;", func(a, b int) (int, int) { return (a + 2*b) % 256, ((a-b)%256 + 256) % 256 }},
			{"A <= B + 1; B <= A + B;", func(a, b int) (int, int) { return (b + 1) % 256, (a + b) % 256 }},
		}
		fm := edakit.Pick(forms)
		a, b := edakit.Int(0, 9), edakit.Int(1, 9)
		a0, b0 := a, b
		k := edakit.Scaled(ctx, 3, 6)
		var tr []string
		for i := 0; i < k; i++ {
			a, b = fm.f(a, b)
			tr = append(tr, edakit.Sprintf("(%d, %d)", a, b))
		}
		which, val := "A", a
		if edakit.Int(0, 1) == 0 {
			which, val = "B", b
		}
		return edakit.Numeric(edakit.Sprintf("Two 8-bit unsigned registers start at A = %d and B = %d. Every clock edge executes the nonblocking assignments `%s` (arithmetic modulo 256). What is %s after %d edges?", a0, b0, fm.code, which, k), edakit.IntStr(val),
			"Every right-hand side uses the values from before the edge. The pairs (A, B) after each edge are "+strings.Join(tr, ", ")+".", map[string]string{"A": edakit.IntStr(a), "B": edakit.IntStr(b)})
	default:
		r := edakit.Int(0, 15)
		r0 := r
		k := edakit.Scaled(ctx, 5, 7)
		var ins []string
		var tr []int
		for i := 0; i < k; i++ {
			en := bit(edakit.Int(0, 3) != 0)
			clr := bit(edakit.Int(0, 4) == 0)
			ins = append(ins, edakit.Sprintf("(%d, %d)", en, clr))
			switch {
			case clr == 1:
				r = 0
			case en == 1:
				r = (r + 1) % 16
			}
			tr = append(tr, r)
		}
		return edakit.Numeric(edakit.Sprintf("A 4-bit counter computes R <= R + 1 (mod 16) when its enable is 1, holds when the enable is 0, and has a synchronous clear that overrides the enable. It starts at R = %d and sees (enable, clear) = %s on successive edges. What is R after the last edge?", r0, strings.Join(ins, ", ")), edakit.IntStr(r),
			"Clear first, then enable, else hold: the values after each edge are "+joinInts(tr, ", ")+".", map[string]string{"final": edakit.IntStr(r)})
	}
}

func ffTiming(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		T := 50 * edakit.Int(10, 40)
		tcq := 5 * edakit.Int(6, 18)
		tsu := 5 * edakit.Int(4, 12)
		slack := 5 * edakit.Int(-20, 30)
		tmax := T - tcq - tsu - slack
		return edakit.Numeric(edakit.Sprintf("A register-to-register path has clock period %d ps, clock-to-q %d ps, longest logic delay %d ps and capture setup time %d ps, with no clock skew. What is the setup slack, in ps?", T, tcq, tmax, tsu), edakit.IntStr(slack),
			edakit.Sprintf("Setup slack is $%d - (%d + %d + %d) = %d$ ps.", T, tcq, tmax, tsu, slack), map[string]string{"arrival": edakit.IntStr(tcq + tmax)})
	case 2:
		tcq := 5 * edakit.Int(4, 12)
		tmin := 5 * edakit.Int(0, 16)
		th := 5 * edakit.Int(2, 16)
		slack := tcq + tmin - th
		return edakit.Numeric(edakit.Sprintf("A short path has minimum clock-to-q %d ps and minimum logic delay %d ps. The capture flip-flop's hold time is %d ps, with no clock skew. What is the hold slack, in ps?", tcq, tmin, th), edakit.IntStr(slack),
			edakit.Sprintf("Hold slack is $(%d + %d) - %d = %d$ ps. The clock period does not enter.", tcq, tmin, th, slack), map[string]string{"early_arrival": edakit.IntStr(tcq + tmin)})
	default:
		s := 5 * edakit.Int(-8, 10)
		tcq := 5 * edakit.Int(6, 16)
		if edakit.Int(0, 1) == 0 {
			T := 50 * edakit.Int(10, 30)
			tsu := 5 * edakit.Int(4, 10)
			tmax := T - tcq - tsu - 5*edakit.Int(-10, 20)
			slack := T + s - (tcq + tmax + tsu)
			return edakit.Numeric(edakit.Sprintf("A path has clock period %d ps, clock-to-q %d ps, longest logic delay %d ps and setup time %d ps. The capture clock arrives %d ps later than the launch clock (a negative value means earlier). What is the setup slack, in ps?", T, tcq, tmax, tsu, s), edakit.IntStr(slack),
				edakit.Sprintf("The capture edge moves by the skew: $%d + (%d) - (%d + %d + %d) = %d$ ps.", T, s, tcq, tmax, tsu, slack), map[string]string{"skew": edakit.IntStr(s)})
		}
		tmin := 5 * edakit.Int(0, 16)
		th := 5 * edakit.Int(2, 12)
		slack := tcq + tmin - (th + s)
		return edakit.Numeric(edakit.Sprintf("A short path has minimum clock-to-q %d ps, minimum logic delay %d ps and hold time %d ps. The capture clock arrives %d ps later than the launch clock (a negative value means earlier). What is the hold slack, in ps?", tcq, tmin, th, s), edakit.IntStr(slack),
			edakit.Sprintf("A later capture edge raises the hold requirement: $(%d + %d) - (%d + (%d)) = %d$ ps.", tcq, tmin, th, s, slack), map[string]string{"skew": edakit.IntStr(s)})
	}
}

func mhz(periodPs float64) string { return edakit.Num(1e6/periodPs, 1) }

func maxFrequency(ctx generator.GeneratorContext) generator.Problem {
	tcq := 5 * edakit.Int(6, 16)
	tsu := 5 * edakit.Int(4, 10)
	switch variant(ctx) {
	case 1:
		tpd := 10 * edakit.Int(20, 150)
		T := tcq + tpd + tsu
		return withExact(edakit.Numeric(edakit.Sprintf("A register-to-register path has clock-to-q %d ps, logic delay %d ps and setup time %d ps. What is the maximum clock frequency, in MHz? Round to 1 decimal place.", tcq, tpd, tsu), mhz(float64(T)),
			edakit.Sprintf("$T_{min} = %d + %d + %d = %d$ ps, so $f_{max} = 10^6 / %d = %s$ MHz.", tcq, tpd, tsu, T, T, mhz(float64(T))), map[string]string{"period": edakit.IntStr(T)}), 1e6/float64(T))
	case 2:
		k := edakit.Scaled(ctx, 3, 5)
		var ds []int
		worst := 0
		for i := 0; i < k; i++ {
			d := 10 * edakit.Int(20, 120)
			ds = append(ds, d)
			worst = max(worst, d)
		}
		T := tcq + worst + tsu
		return withExact(edakit.Numeric(edakit.Sprintf("A design's register-to-register paths have logic delays of %s ps. Every flip-flop has clock-to-q %d ps and setup %d ps, with no skew. What is the maximum clock frequency, in MHz? Round to 1 decimal place.", joinInts(ds, ", "), tcq, tsu), mhz(float64(T)),
			edakit.Sprintf("The %d ps path is critical: $T_{min} = %d + %d + %d = %d$ ps, giving %s MHz.", worst, tcq, worst, tsu, T, mhz(float64(T))), map[string]string{"critical": edakit.IntStr(worst)}), 1e6/float64(T))
	default:
		tpd := 10 * edakit.Int(30, 120)
		s := 5 * edakit.Int(-10, 12)
		T := tcq + tpd + tsu - s
		return withExact(edakit.Numeric(edakit.Sprintf("A critical path has clock-to-q %d ps, logic delay %d ps and setup %d ps. Its capture clock arrives %d ps after its launch clock (negative means before). What is the maximum clock frequency for this path, in MHz? Round to 1 decimal place.", tcq, tpd, tsu, s), mhz(float64(T)),
			edakit.Sprintf("$T_{min} = %d + %d + %d - (%d) = %d$ ps, giving %s MHz.", tcq, tpd, tsu, s, T, mhz(float64(T))), map[string]string{"period": edakit.IntStr(T)}), 1e6/float64(T))
	}
}

func metastability(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		fc := edakit.Pick([]int{100, 200, 250, 400, 500, 800, 1000})
		fd := edakit.Pick([]int{1, 2, 5, 10, 20, 50})
		tw := edakit.Pick([]int{10, 15, 20, 25, 30, 50})
		r := fc * fd * tw
		return edakit.Numeric(edakit.Sprintf("A flip-flop clocked at %d MHz samples an asynchronous signal that toggles %d million times per second. Its metastability aperture is %d ps. About how many metastable events occur per second?", fc, fd, tw), edakit.IntStr(r),
			edakit.Sprintf("$R = f_c f_d T_w = (%d \\times 10^6)(%d \\times 10^6)(%d \\times 10^{-12}) = %d$ per second.", fc, fd, tw, r), map[string]string{"rate": edakit.IntStr(r)})
	case 2:
		tau := edakit.Pick([]int{8, 10, 12, 15, 20, 25})
		t := edakit.Int(2*tau, 7*tau)
		p := math.Exp(-float64(t) / float64(tau))
		return withExact(edakit.Numeric(edakit.Sprintf("A flip-flop with regeneration time constant $\\tau = %d$ ps has gone metastable. What is the probability that it is still unresolved %d ps later? Round to 5 decimal places.", tau, t), edakit.Num(p, 5),
			edakit.Sprintf("$P = e^{-t/\\tau} = e^{-%d/%d} = %s$.", t, tau, edakit.Num(p, 5)), map[string]string{"ratio": edakit.Num(float64(t)/float64(tau), 3)}), p)
	default:
		tau := edakit.Pick([]int{8, 10, 12, 15, 20, 25, 30})
		k := edakit.Int(3, 12)
		dt := float64(tau) * float64(k) * math.Ln10
		return withExact(edakit.Numeric(edakit.Sprintf("A synchronizer's flip-flops have $\\tau = %d$ ps. How much extra resolution time, in ps, improves its failure rate by a factor of $10^{%d}$? Round to 1 decimal place.", tau, k), edakit.Num(dt, 1),
			edakit.Sprintf("$\\Delta t = \\tau \\ln 10^{%d} = %d \\cdot %d \\ln 10 = %s$ ps.", k, tau, k, edakit.Num(dt, 1)), map[string]string{"ln_factor": edakit.Num(float64(k)*math.Ln10, 3)}), dt)
	}
}

// log10MTBF returns log10 of the MTBF in seconds for resolution time tr (ps), tau (ps), aperture tw
// (ps) and clock and data rates in MHz.
func log10MTBF(tr, tau float64, tw, fc, fd int) float64 {
	return tr/(tau*math.Ln10) - math.Log10(float64(tw*fc*fd))
}

func synchronizerMTBF(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		tau := edakit.Pick([]int{8, 10, 12, 15, 20})
		tr := 10 * edakit.Int(30, 100)
		tw := edakit.Pick([]int{10, 20, 30, 50})
		fc := edakit.Pick([]int{200, 500, 1000})
		fd := edakit.Pick([]int{1, 10, 50, 100})
		v := log10MTBF(float64(tr), float64(tau), tw, fc, fd)
		return withExact(edakit.Numeric(edakit.Sprintf("A synchronizer allows %d ps of resolution time. Its flip-flops have $\\tau = %d$ ps and aperture $T_w = %d$ ps; the clock is %d MHz and the data toggles at %d MHz. What is $\\log_{10}$ of the MTBF in seconds? Round to 2 decimal places.", tr, tau, tw, fc, fd), edakit.Num(v, 2),
			edakit.Sprintf("$\\log_{10} \\text{MTBF} = \\frac{%d}{%d \\ln 10} - \\log_{10}(%d) = %s$, where %d per second is the metastable event rate $T_w f_c f_d$.", tr, tau, tw*fc*fd, edakit.Num(v, 2), tw*fc*fd),
			map[string]string{"rate": edakit.IntStr(tw * fc * fd)}), v)
	case 2:
		fc := edakit.Pick([]int{250, 400, 500, 800, 1000, 1250})
		Tc := 1000000 / fc
		tcq := 5 * edakit.Int(6, 16)
		tsu := 5 * edakit.Int(4, 10)
		tau := edakit.Pick([]int{10, 15, 20, 25})
		tw := edakit.Pick([]int{10, 20, 30})
		fd := edakit.Pick([]int{1, 10, 50})
		tr := Tc - tcq - tsu
		v := log10MTBF(float64(tr), float64(tau), tw, fc, fd)
		return withExact(edakit.Numeric(edakit.Sprintf("A two-flip-flop synchronizer runs at %d MHz with no logic between the flops. Clock-to-q is %d ps, setup is %d ps, $\\tau = %d$ ps and $T_w = %d$ ps; the asynchronous input toggles at %d MHz. What is $\\log_{10}$ of the MTBF in seconds? Round to 2 decimal places.", fc, tcq, tsu, tau, tw, fd), edakit.Num(v, 2),
			edakit.Sprintf("The period is %d ps, so $t_r = %d - %d - %d = %d$ ps and $\\log_{10}\\text{MTBF} = %d/(%d \\ln 10) - \\log_{10}(%d) = %s$.", Tc, Tc, tcq, tsu, tr, tr, tau, tw*fc*fd, edakit.Num(v, 2)),
			map[string]string{"t_r": edakit.IntStr(tr)}), v)
	default:
		for {
			fc := edakit.Pick([]int{1000, 1250, 2000, 2500})
			Tc := 1000000 / fc
			tcq := 5 * edakit.Int(8, 14)
			tsu := 5 * edakit.Int(4, 8)
			tau := edakit.Pick([]int{20, 25, 30, 40})
			tw := edakit.Pick([]int{20, 30, 50})
			fd := edakit.Pick([]int{50, 100, 200})
			target := edakit.Pick([]int{9, 10, 12, 14, 16})
			r := Tc - tcq - tsu
			if r <= 0 {
				continue
			}
			N := 2
			for ; N <= 8; N++ {
				if log10MTBF(float64((N-1)*r), float64(tau), tw, fc, fd) >= float64(target) {
					break
				}
			}
			if N > 6 {
				continue
			}
			hi := log10MTBF(float64((N-1)*r), float64(tau), tw, fc, fd)
			lo := log10MTBF(float64((N-2)*r), float64(tau), tw, fc, fd)
			if hi-float64(target) < 0.1 || (N > 2 && float64(target)-lo < 0.1) {
				continue
			}
			return edakit.Numeric(edakit.Sprintf("A synchronizer chain runs at %d MHz. Each extra flip-flop adds one period minus clock-to-q (%d ps) minus setup (%d ps) of resolution time, so an N-flop chain resolves for $(N-1)$ such intervals. With $\\tau = %d$ ps, $T_w = %d$ ps and a %d MHz asynchronous input, what is the smallest N (at least 2) whose MTBF reaches $10^{%d}$ seconds?", fc, tcq, tsu, tau, tw, fd, target), edakit.IntStr(N),
				edakit.Sprintf("Each interval is %d ps. With N = %d, $\\log_{10}\\text{MTBF} \\approx %s$, which meets the target; one flop fewer falls short.", r, N, edakit.Num(hi, 2)),
				map[string]string{"interval": edakit.IntStr(r)})
		}
	}
}

func counters(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 16)
	tcq := 5 * edakit.Int(8, 20)
	switch variant(ctx) {
	case 1:
		return edakit.Numeric(edakit.Sprintf("A ripple counter of %d bits is built from toggle flip-flops, each clocked by the previous stage's output, with clock-to-q %d ps. In the worst case, how long after the clock edge does the most significant bit settle, in ps?", n, tcq), edakit.IntStr(n*tcq),
			edakit.Sprintf("On a step like 0111...1 to 1000...0 every stage toggles in turn: $%d \\times %d = %d$ ps.", n, tcq, n*tcq), map[string]string{"stages": edakit.IntStr(n)})
	case 2:
		tsu := 5 * edakit.Int(4, 10)
		T := n*tcq + tsu
		return withExact(edakit.Numeric(edakit.Sprintf("Logic clocked by the main clock samples the outputs of a ripple counter of %d bits. Each counter stage has clock-to-q %d ps, and the sampling flip-flops need %d ps of setup. What is the maximum clock frequency, in MHz? Round to 1 decimal place.", n, tcq, tsu), mhz(float64(T)),
			edakit.Sprintf("The outputs settle after $%d \\times %d = %d$ ps; adding setup gives %d ps, so %s MHz.", n, tcq, n*tcq, T, mhz(float64(T))), map[string]string{"period": edakit.IntStr(T)}), 1e6/float64(T))
	default:
		var N int
		for {
			N = edakit.Int(5, 200)
			if N&(N-1) != 0 {
				break
			}
		}
		b := clog2(N)
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(edakit.Sprintf("How many flip-flops does a binary-encoded mod-%d counter need?", N), edakit.IntStr(b),
				edakit.Sprintf("$\\lceil \\log_2 %d \\rceil = %d$.", N, b), map[string]string{"unused": edakit.IntStr((1 << uint(b)) - N)})
		}
		return edakit.Numeric(edakit.Sprintf("A binary-encoded mod-%d counter uses the minimum number of flip-flops. How many of their codes are never used in normal counting?", N), edakit.IntStr((1<<uint(b))-N),
			edakit.Sprintf("It needs %d flip-flops, which have %d codes; %d are used, leaving %d.", b, 1<<uint(b), N, (1<<uint(b))-N), map[string]string{"flops": edakit.IntStr(b)})
	}
}

var syncDesignItems = []edakit.Item{
	{KP: 1, Q: "Which structure violates the synchronous design discipline?",
		A: "A loop of gates feeding back to its own input with no flip-flop on it", D: []string{"A flip-flop whose output feeds its own input through an adder", "Two flip-flops in series with no logic between them", "A register whose input comes from a mux selecting its own output"},
		Why: "Every cycle must pass through a flip-flop. A loop broken by a register is fine; a purely combinational loop is not."},
	{KP: 1, Q: "Under the synchronous discipline, where may a flip-flop's clock pin be driven from?",
		A: "The clock network, possibly through a clock-gating cell", D: []string{"Any data signal that toggles once per cycle", "The output of another flip-flop's data path logic", "A combinational decode of the register outputs"},
		Why: "Clocks come from the clock network so their timing is known; deriving a clock from data logic creates skew and glitches the tools cannot reason about."},
	{KP: 1, Q: "Why are glitches in combinational logic between registers acceptable in a synchronous design?",
		A: "The registers sample only after the logic settles, within the setup and hold window", D: []string{"Glitches are always shorter than the hold time", "Flip-flops ignore any input change of less than one period", "Synthesis removes all glitches automatically"},
		Why: "Only the value present around the clock edge is stored, so transient values before settling do not matter."},
	{KP: 1, Q: "An RTL block uses a latch inferred from an incomplete `if` without an `else`. Which discipline rule does the inferred latch break, in a flop-based design?",
		A: "State is held in an element that is not an edge-triggered flip-flop", D: []string{"A clock is derived from data logic", "The design contains a register with an enable", "The combinational logic has more than two levels"},
		Why: "An unintended latch is a level-sensitive storage element outside the flop discipline, which complicates timing and is a standard lint error."},
	{KP: 2, Q: "Which property of synchronous design lets static timing analysis check a design without simulation vectors?",
		A: "Correctness depends only on every path settling between clock edges", D: []string{"Every flip-flop has the same setup time", "The logic is minimized before analysis", "All inputs change at the same time"},
		Why: "Because only path delays relative to clock edges matter, STA can bound every path's delay independently of data values."},
	{KP: 2, Q: "A cycle-based simulator such as Verilator evaluates each combinational block once per clock cycle. What would break that scheme?",
		A: "A combinational loop, which has no topological order", D: []string{"A register with a synchronous reset", "A module with a deep adder tree", "A design with two pipeline stages"},
		Why: "Single-pass evaluation needs a levelized order of the logic; a combinational loop has none and may oscillate."},
	{KP: 2, Q: "Retiming moves registers across combinational logic. Which assumption makes this a behavior-preserving transformation?",
		A: "The circuit is a synchronous network whose behavior is defined cycle by cycle at the registers", D: []string{"All registers have identical clock-to-q delays", "The logic contains no XOR gates", "Each register drives exactly one gate"},
		Why: "Retiming preserves the input-output behavior of a synchronous circuit (up to initial state), which is only defined when all state is in clocked registers."},
	{KP: 2, Q: "Synthesis reports a combinational loop. Why is it treated as an error rather than a warning?",
		A: "Timing, simulation and equivalence engines all assume an acyclic combinational network", D: []string{"Loops always consume excessive area", "Loops are illegal in the Verilog language", "Loops make the design impossible to place"},
		Why: "The loop can oscillate or hold state, and the timing graph, the levelized simulator and the equivalence model all lose their meaning."},
	{KP: 3, Q: "What is the safe way to stop a register from loading on some cycles?",
		A: "A recirculating mux on its D input, or an integrated clock-gating cell", D: []string{"An AND gate between the clock and the register's clock pin", "Disconnecting the register's reset", "Driving the D input with a tristate buffer"},
		Why: "Both keep the clock clean: the mux holds the value functionally, and the ICG cell latches the enable so the gated clock cannot glitch."},
	{KP: 3, Q: "What does an asynchronous-assert, synchronous-deassert reset synchronizer guarantee?",
		A: "Reset takes effect immediately, and all flip-flops leave reset on the same clock edge", D: []string{"Reset is delayed by two cycles in both directions", "Reset can never cause a timing check", "The reset signal can be driven by any combinational logic"},
		Why: "Assertion bypasses the clock, while the release is retimed through flip-flops so recovery and removal checks are met everywhere."},
	{KP: 3, Q: "Which timing checks apply to the release of an asynchronous reset at a flip-flop?",
		A: "Recovery and removal", D: []string{"Setup and hold on the D pin", "Clock-to-q and pulse width", "Minimum period and skew"},
		Why: "Recovery and removal are the reset-pin analogues of setup and hold: the reset must be released far enough from the clock edge."},
	{KP: 3, Q: "Gating a clock with a plain AND gate driven by an enable from a flip-flop clocked on the same edge risks what?",
		A: "A glitch or a shortened clock pulse when the enable changes while the clock is high", D: []string{"A permanent hold violation on every path", "The register losing its reset value", "The enable being optimized away by synthesis"},
		Why: "If the enable changes during the high phase, the AND output can chop or create a pulse; an ICG cell latches the enable to prevent it."},
}

var cdcItems = []edakit.Item{
	{KP: 1, Q: "A single control bit crosses from a 400 MHz domain into an unrelated 100 MHz domain. What is the minimum safe structure?",
		A: "Drive it from a source flip-flop, then pass it through two destination flip-flops", D: []string{"Sample it with one destination flip-flop and use it directly", "Route it through a buffer chain to add delay", "Constrain it with a tighter clock period"},
		Why: "A two-flop synchronizer bounds the metastability risk, and launching from a flip-flop prevents glitches being sampled."},
	{KP: 1, Q: "A one-cycle pulse from a fast domain crosses into a much slower domain through a two-flop synchronizer. What can go wrong?",
		A: "The pulse can fall between two slow clock edges and never be seen", D: []string{"The pulse will always be seen twice", "The synchronizer turns the pulse into a level permanently", "The pulse violates the source domain's hold time"},
		Why: "A pulse shorter than the receiving clock's period may not overlap any sampling edge, so it must be stretched or converted to a toggle."},
	{KP: 1, Q: "Why must the signal entering a synchronizer come directly from a flip-flop rather than from combinational logic?",
		A: "Logic can glitch, and the receiving clock may sample the glitch as a real pulse", D: []string{"Combinational logic adds too much setup time", "Flip-flop outputs are never metastable", "Synthesis cannot place logic near a synchronizer"},
		Why: "In the source domain glitches are harmless because the source clock waits for settling; the asynchronous receiver does not."},
	{KP: 1, Q: "A toggle synchronizer sends an event across domains by flipping a source flip-flop. How does the receiver detect the event?",
		A: "It synchronizes the toggling level and detects a change between consecutive synchronized samples", D: []string{"It counts rising edges of the raw source signal", "It waits for the toggle to return to zero", "It samples the source flip-flop with its own clock once"},
		Why: "A level change is held until the next event, so it cannot be missed; an XOR of two successive synchronized samples recovers the event."},
	{KP: 2, Q: "An 8-bit binary counter value is passed to another clock domain by synchronizing each bit separately. Why is this unsafe?",
		A: "Bits can resolve in different cycles, so the receiver may see a value that never existed", D: []string{"Eight synchronizers use too much area", "Binary counters cannot be synchronized at all", "The counter's carry chain becomes metastable"},
		Why: "When several bits change together, independent synchronizers can split the change across cycles, producing a mixed value."},
	{KP: 2, Q: "Which value can be safely synchronized bit by bit?",
		A: "A Gray-coded pointer that changes by one step per source cycle", D: []string{"A binary address that increments by one", "A 32-bit data word written once per cycle", "Two independent control bits that must arrive together"},
		Why: "A Gray code changes one bit per step, so a sample is either the old or the new value, never a third value."},
	{KP: 2, Q: "A data bus crosses domains alongside a valid bit. Which scheme is correct?",
		A: "Hold the bus stable, synchronize only the valid bit, and sample the bus after valid arrives", D: []string{"Synchronize every bus bit and the valid bit with separate two-flop chains", "Sample the bus on the receiving clock without any qualifier", "Add equal delays on all bus wires"},
		Why: "The bus is stable by the time the synchronized qualifier arrives, so it can be sampled directly; only one bit needs synchronization."},
	{KP: 2, Q: "Two separately synchronized control signals are combined in an AND gate in the destination domain. What problem does a CDC tool report?",
		A: "Reconvergence: the two may arrive in different cycles, so the combination can be wrong for a cycle", D: []string{"A combinational loop through the AND gate", "A hold violation on the synchronizer", "A setup violation in the source domain"},
		Why: "Each synchronizer has a one-cycle uncertainty, so related signals synchronized separately can be misaligned when they reconverge."},
	{KP: 3, Q: "In an asynchronous FIFO, why are the read and write pointers Gray-coded before crossing?",
		A: "So each synchronized pointer is either its old or its new value, never an invalid one", D: []string{"So the FIFO needs fewer memory entries", "So the pointers can be compared without synchronization", "So the memory can be single-ported"},
		Why: "Only one bit changes per increment, so the synchronized copy is always a legal, possibly slightly old, pointer."},
	{KP: 3, Q: "An asynchronous FIFO's full flag is computed from the write pointer and a synchronized copy of the read pointer. What is true of that flag?",
		A: "It may stay asserted a few cycles after space frees up, which is safe", D: []string{"It may deassert before space is available, which is safe", "It is always exactly correct", "It can only be computed in the read domain"},
		Why: "The synchronized read pointer lags the real one, so full is pessimistic: the writer waits slightly longer but never overwrites."},
	{KP: 3, Q: "A writer sends bursts of 64 words at 200 MHz into an asynchronous FIFO drained at 150 MHz. About how many words accumulate during one burst, before synchronizer margin?",
		A: "16", D: []string{"48", "64", "8"},
		Why: "While 64 words are written, the reader removes $64 \\times 150/200 = 48$, so about $64 (1 - 150/200) = 16$ words build up."},
	{KP: 3, Q: "Why does static timing analysis not check paths between asynchronous clock domains?",
		A: "There is no fixed phase relation between the clocks, so no setup or hold requirement can be computed", D: []string{"Those paths are always too short to fail", "STA cannot handle more than one clock", "CDC paths are removed by synthesis"},
		Why: "With unrelated clocks the capture edge can fall anywhere relative to the launch edge, so the paths are declared false or asynchronous and checked structurally by CDC tools."},
}
