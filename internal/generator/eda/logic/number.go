package logic

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func registerNumber(reg *generator.Registry) {
	reg.Register("logic.number.twos_complement", retryTies(twosComplement))
	reg.Register("logic.number.overflow", retryTies(overflow))
	reg.Register("logic.number.fixed_point", retryTies(fixedPoint))
	reg.Register("logic.number.floating_point", retryTies(floatingPoint))
	reg.Register("logic.number.gray_code", retryTies(grayCode))
}

func widthPick(ctx generator.GeneratorContext, ws []int) int {
	return ws[edakit.Scaled(ctx, 0, len(ws)-1)]
}

func twosComplement(ctx generator.GeneratorContext) generator.Problem {
	n := widthPick(ctx, []int{4, 6, 8, 8, 10, 12, 16})
	switch variant(ctx) {
	case 1:
		switch edakit.Int(0, 2) {
		case 0:
			ans := -(int64(1) << uint(n-1))
			return edakit.Numeric(edakit.Sprintf("What is the most negative value a two's-complement word of %d bits can hold?", n), i64(ans),
				edakit.Sprintf("Only the top bit set gives $-2^{%d} = %d$.", n-1, ans), map[string]string{"n": edakit.IntStr(n)})
		case 1:
			ans := (int64(1) << uint(n-1)) - 1
			return edakit.Numeric(edakit.Sprintf("What is the largest value a two's-complement word of %d bits can hold?", n), i64(ans),
				edakit.Sprintf("Every bit except the top one set gives $2^{%d} - 1 = %d$.", n-1, ans), map[string]string{"n": edakit.IntStr(n)})
		default:
			ans := int64(1) << uint(n-1)
			return edakit.Numeric(edakit.Sprintf("How many of the $2^{%d}$ bit patterns of a two's-complement word of %d bits represent negative values?", n, n), i64(ans),
				edakit.Sprintf("Exactly the patterns with the top bit set are negative: $2^{%d} = %d$ of them.", n-1, ans), map[string]string{"n": edakit.IntStr(n)})
		}
	case 2:
		p := uint64(edakit.Int(0, (1<<uint(n))-1))
		if edakit.Int(0, 9) < 7 {
			p |= 1 << uint(n-1)
		}
		u := int64(p)
		ans := signed(p, n)
		return edakit.Numeric(edakit.Sprintf("What signed value does the %d-bit two's-complement pattern %s represent?", n, bin(p, n)), i64(ans),
			edakit.Sprintf("As unsigned the pattern is %d. The top bit is %d, so the signed value is %d.", u, p>>uint(n-1)&1, ans),
			map[string]string{"unsigned": i64(u)})
	default:
		if edakit.Int(0, 4) == 0 {
			mn := -(int64(1) << uint(n-1))
			return edakit.Numeric(edakit.Sprintf("A negator of width %d bits computes invert-and-add-one on the pattern for %d. What signed value does its output pattern represent?", n, mn), i64(mn),
				edakit.Sprintf("The pattern 1 followed by zeros inverts to 0 followed by ones; adding one restores 1 followed by zeros. The true answer %d is out of range, so the hardware returns %d again.", -mn, mn),
				map[string]string{"n": edakit.IntStr(n)})
		}
		x := int64(edakit.Int(1, (1<<uint(n-1))-1))
		ans := (int64(1) << uint(n)) - x
		return edakit.Numeric(edakit.Sprintf("Negate %d in %d-bit two's complement by inverting every bit and adding one. Read the resulting pattern as an unsigned integer. What is it?", x, n), i64(ans),
			edakit.Sprintf("%d is %s; inverting gives %s and adding one gives %s, which is $2^{%d} - %d = %d$.", x, bin(uint64(x), n), bin(^uint64(x), n), bin(uint64(ans), n), n, x, ans),
			map[string]string{"pattern": bin(uint64(ans), n)})
	}
}

func overflow(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		n := widthPick(ctx, []int{6, 8, 8, 10, 12})
		lim := int64(1) << uint(n-1)
		var a, b int64
		for {
			a = int64(edakit.Int(int(-lim), int(lim-1)))
			b = int64(edakit.Int(int(-lim), int(lim-1)))
			if (a >= 0) == (b >= 0) && a != 0 && b != 0 {
				break
			}
		}
		if edakit.Int(0, 3) == 0 {
			b = -b
			if b >= lim {
				b = lim - 1
			}
		}
		ans := wrap(a+b, n)
		ov := ans != a+b
		why := "The true sum fits in range, so no overflow occurs."
		if ov {
			why = edakit.Sprintf("The true sum %d is outside $[%d, %d]$, so it wraps by $2^{%d}$: signed overflow.", a+b, -lim, lim-1, n)
		}
		return edakit.Numeric(edakit.Sprintf("A two's-complement adder of width %d bits adds %d and %d. What signed value does its %d-bit result represent?", n, a, b, n), i64(ans),
			why+edakit.Sprintf(" The stored result is %d.", ans), map[string]string{"true_sum": i64(a + b), "overflow": edakit.IntStr(bit(ov))})
	case 2:
		n := edakit.Pick([]int{4, 5, 6, 8})
		m := edakit.Pick([]int{8, 12, 16})
		for m <= n {
			m += 4
		}
		p := uint64(edakit.Int(0, (1<<uint(n))-1))
		if edakit.Int(0, 4) != 0 {
			p |= 1 << uint(n-1)
		}
		ext := p
		if p>>uint(n-1)&1 == 1 {
			ext = p | (((1 << uint(m)) - 1) &^ ((1 << uint(n)) - 1))
		}
		return edakit.Numeric(edakit.Sprintf("The %d-bit two's-complement pattern %s is sign-extended to %d bits. Read the %d-bit result as an unsigned integer. What is it?", n, bin(p, n), m, m), i64(int64(ext)),
			edakit.Sprintf("Copy the sign bit into the %d new upper positions: %s, which is %d unsigned and still %d signed.", m-n, bin(ext, m), ext, signed(p, n)),
			map[string]string{"extended": bin(ext, m)})
	default:
		n := edakit.Pick([]int{8, 8, 16})
		lim := 1 << uint(n-1)
		if n == 16 {
			lim = 2000
		}
		x := int64(edakit.Int(1, lim-1))
		if edakit.Int(0, 4) != 0 {
			x = -x
		}
		s := edakit.Int(1, 4)
		ans := int64(math.Floor(float64(x) / float64(int64(1)<<uint(s))))
		pat := uint64(x) & ((1 << uint(n)) - 1)
		return edakit.Numeric(edakit.Sprintf("The %d-bit two's-complement value %d (pattern %s) is shifted right arithmetically by %d. What signed value results?", n, x, bin(pat, n), s), i64(ans),
			edakit.Sprintf("An arithmetic shift computes the floor: $\\lfloor %d / %d \\rfloor = %d$. Truncation toward zero would give %d for a negative value with bits shifted out.", x, int64(1)<<uint(s), ans, x/(int64(1)<<uint(s))),
			map[string]string{"shift": edakit.IntStr(s)})
	}
}

func fixedPoint(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		w := widthPick(ctx, []int{8, 8, 10, 12, 16})
		f := edakit.Int(2, 6)
		p := uint64(edakit.Int(0, (1<<uint(w))-1))
		k := signed(p, w)
		val := float64(k) / float64(int64(1)<<uint(f))
		ans := edakit.Num(val, f)
		return edakit.Numeric(edakit.Sprintf("A signed fixed-point format has %d bits in total, %d of them fractional, in two's complement. What value does the pattern %s represent? Give the exact decimal value.", w, f, bin(p, w)), ans,
			edakit.Sprintf("As a two's-complement integer the pattern is %d. Dividing by $2^{%d} = %d$ gives %s.", k, f, int64(1)<<uint(f), ans),
			map[string]string{"integer": i64(k)})
	case 2:
		w := edakit.Pick([]int{8, 10, 12, 16})
		f := edakit.Int(2, 7)
		step := 1 / float64(int64(1)<<uint(f))
		switch edakit.Int(0, 3) {
		case 0:
			v := math.Pow(2, float64(w-1-f)) - step
			ans := edakit.Num(v, f)
			return edakit.Numeric(edakit.Sprintf("A signed two's-complement fixed-point format has %d total bits and %d fractional bits. What is the largest value it can represent? Give the exact decimal value.", w, f), ans,
				edakit.Sprintf("The largest integer is $2^{%d} - 1$; scaled by $2^{-%d}$ it is $2^{%d} - 2^{-%d} = %s$.", w-1, f, w-1-f, f, ans), map[string]string{"step": edakit.Num(step, f)})
		case 1:
			v := -math.Pow(2, float64(w-1-f))
			ans := edakit.Num(v, 0)
			return edakit.Numeric(edakit.Sprintf("A signed two's-complement fixed-point format has %d total bits and %d fractional bits. What is the most negative value it can represent?", w, f), ans,
				edakit.Sprintf("The most negative integer $-2^{%d}$ scaled by $2^{-%d}$ is $-2^{%d} = %s$.", w-1, f, w-1-f, ans), map[string]string{"step": edakit.Num(step, f)})
		case 2:
			ans := edakit.Num(step, f)
			return edakit.Numeric(edakit.Sprintf("A fixed-point format has %d total bits and %d fractional bits. What is its resolution (the step between adjacent values)? Give the exact decimal value.", w, f), ans,
				edakit.Sprintf("The step is $2^{-%d} = %s$, independent of the total width.", f, ans), map[string]string{"step": ans})
		default:
			ans := edakit.Num(step/2, f+1)
			return edakit.Numeric(edakit.Sprintf("A real number is rounded to the nearest value of a fixed-point format with %d fractional bits. What is the largest possible rounding error? Give the exact decimal value.", f), ans,
				edakit.Sprintf("Rounding to nearest errs by at most half a step: $2^{-%d} = %s$.", f+1, ans), map[string]string{"step": edakit.Num(step, f)})
		}
	default:
		switch edakit.Int(0, 2) {
		case 0, 1:
			w1, f1 := edakit.Pick([]int{8, 12, 16}), 0
			w2, f2 := edakit.Pick([]int{8, 10, 12, 16}), 0
			f1 = edakit.Int(2, w1-1)
			f2 = edakit.Int(2, w2-1)
			if edakit.Int(0, 1) == 0 {
				return edakit.Numeric(edakit.Sprintf("A signed value with %d total bits (%d fractional) is multiplied by a signed value with %d total bits (%d fractional). How many total bits does the exact product need?", w1, f1, w2, f2), edakit.IntStr(w1+w2),
					edakit.Sprintf("The integers multiply into %d + %d = %d bits.", w1, w2, w1+w2), map[string]string{"frac": edakit.IntStr(f1 + f2)})
			}
			return edakit.Numeric(edakit.Sprintf("A signed value with %d total bits (%d fractional) is multiplied by a signed value with %d total bits (%d fractional). How many fractional bits does the exact product have?", w1, f1, w2, f2), edakit.IntStr(f1+f2),
				edakit.Sprintf("The scale factors $2^{-%d}$ and $2^{-%d}$ multiply to $2^{-%d}$.", f1, f2, f1+f2), map[string]string{"width": edakit.IntStr(w1 + w2)})
		default:
			nTerms := edakit.Int(3, 300)
			pw := edakit.Pick([]int{16, 24, 28, 32})
			g := clog2(nTerms)
			return edakit.Numeric(edakit.Sprintf("An accumulator sums %d products, each a %d-bit signed value. How many bits wide must the accumulator be so the sum can never overflow?", nTerms, pw), edakit.IntStr(pw+g),
				edakit.Sprintf("Summing %d terms needs $\\lceil \\log_2 %d \\rceil = %d$ guard bits on top of the %d product bits: %d bits.", nTerms, nTerms, g, pw, pw+g),
				map[string]string{"guard_bits": edakit.IntStr(g)})
		}
	}
}

func floatingPoint(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		s := edakit.Int(0, 1)
		e := edakit.Int(10, 21)
		f := 32 * edakit.Int(0, 31)
		val := math.Pow(2, float64(e-15)) * float64(1024+f) / 1024
		if s == 1 {
			val = -val
		}
		ans := edakit.Num(val, 10)
		return edakit.Numeric(edakit.Sprintf("A half-precision (binary16) number has sign bit %d, exponent field E = %d and fraction field F = %d (5 exponent bits, bias 15, 10 fraction bits). What is its value? Give the exact decimal value.", s, e, f), ans,
			edakit.Sprintf("The value is $(-1)^{%d} \\cdot 2^{%d - 15} \\cdot (1 + %d/1024) = %s$.", s, e, f, ans),
			map[string]string{"unbiased": edakit.IntStr(e - 15)})
	case 2:
		k := edakit.Int(-3, 9)
		j := edakit.Int(8, 15) // significand j/8 in [1, 2)
		val := float64(j) * math.Pow(2, float64(k-3))
		vs := edakit.Num(val, 8)
		fField := (j - 8) * 128
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(edakit.Sprintf("The value %s is encoded in IEEE half precision (5 exponent bits, bias 15, 10 fraction bits). What is the exponent field E, as an unsigned integer?", vs), edakit.IntStr(k+15),
				edakit.Sprintf("$%s = %s \\cdot 2^{%d}$ with the significand in $[1, 2)$, so $E = %d + 15 = %d$.", vs, edakit.Num(float64(j)/8, 3), k, k, k+15),
				map[string]string{"fraction": edakit.IntStr(fField)})
		}
		return edakit.Numeric(edakit.Sprintf("The value %s is encoded in IEEE half precision (5 exponent bits, bias 15, 10 fraction bits). What is the fraction field F, as an unsigned integer?", vs), edakit.IntStr(fField),
			edakit.Sprintf("$%s = %s \\cdot 2^{%d}$, so $F = (%s - 1) \\cdot 1024 = %d$.", vs, edakit.Num(float64(j)/8, 3), k, edakit.Num(float64(j)/8, 3), fField),
			map[string]string{"exponent": edakit.IntStr(k + 15)})
	default:
		switch edakit.Int(0, 5) {
		case 0:
			return edakit.Numeric("In IEEE half precision (5 exponent bits, bias 15, 10 fraction bits), the smallest positive subnormal number is $2^p$. What is p?", "-24",
				"Subnormals are $2^{-14} \\cdot F/2^{10}$; the smallest has F = 1, giving $2^{-24}$.", map[string]string{"min_normal_exp": "-14"})
		case 1:
			return edakit.Numeric("In IEEE half precision (5 exponent bits, bias 15, 10 fraction bits), the smallest positive normal number is $2^p$. What is p?", "-14",
				"The smallest normal exponent field is 1, so the value is $2^{1-15} = 2^{-14}$.", map[string]string{"subnormal_exp": "-24"})
		default:
			k := edakit.Int(-4, 15)
			lo := math.Pow(2, float64(k))
			j := edakit.Int(1, 7)
			x := lo * (1 + float64(j)/8)
			return edakit.Numeric(edakit.Sprintf("In IEEE half precision (5 exponent bits, bias 15, 10 fraction bits), the spacing between adjacent representable numbers near %s is $2^p$. What is p?", edakit.Num(x, 8)), edakit.IntStr(k-10),
				edakit.Sprintf("%s lies in the binade $[2^{%d}, 2^{%d})$, where the ulp is $2^{%d - 10} = 2^{%d}$.", edakit.Num(x, 8), k, k+1, k, k-10),
				map[string]string{"binade": edakit.IntStr(k)})
		}
	}
}

func grayCode(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 8)
	switch variant(ctx) {
	case 1:
		b := uint64(edakit.Int(1, (1<<uint(n))-1))
		g := b ^ (b >> 1)
		return edakit.Numeric(edakit.Sprintf("Convert the %d-bit binary number %s to reflected binary Gray code. Enter the Gray code word read as an unsigned decimal integer.", n, bin(b, n)), edakit.IntStr(int(g)),
			edakit.Sprintf("$g = b \\oplus (b \\gg 1) = %s \\oplus %s = %s$, which is %d.", bin(b, n), bin(b>>1, n), bin(g, n), g),
			map[string]string{"gray": bin(g, n)})
	case 2:
		g := uint64(edakit.Int(1, (1<<uint(n))-1))
		b := g
		for sh := uint64(1); sh < 64; sh <<= 1 {
			b ^= b >> sh
		}
		return edakit.Numeric(edakit.Sprintf("The %d-bit Gray code word %s encodes which unsigned integer?", n, bin(g, n)), edakit.IntStr(int(b)),
			edakit.Sprintf("Prefix-XOR from the top: $b_i$ is the XOR of all Gray bits at or above position i, giving %s = %d.", bin(b, n), b),
			map[string]string{"binary": bin(b, n)})
	default:
		a := edakit.Int(0, (1<<uint(n))-14)
		steps := edakit.Int(3, 12)
		bEnd := a + steps
		tog := 0
		for i := a; i < bEnd; i++ {
			tog += popcount(uint64(i ^ (i + 1)))
		}
		return edakit.Numeric(edakit.Sprintf("A plain binary counter counts from %d up to %d, one step per clock. How many bit toggles occur in total over those %d steps? (A Gray-code counter would make exactly %d.)", a, bEnd, steps, steps), edakit.IntStr(tog),
			edakit.Sprintf("Each step from i to i + 1 toggles popcount(i XOR (i + 1)) bits, one plus the number of trailing ones of i. Summed over the %d steps that is %d.", steps, tog),
			map[string]string{"steps": edakit.IntStr(steps)})
	}
}
