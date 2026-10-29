package hdl

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Four-state bit operations as IEEE 1364 defines them: z behaves as x on every logic
// operator input, a controlling 0 (for AND) or 1 (for OR) wins over an unknown.
func and4(a, b byte) byte {
	if a == '0' || b == '0' {
		return '0'
	}
	if a == '1' && b == '1' {
		return '1'
	}
	return 'x'
}

func or4(a, b byte) byte {
	if a == '1' || b == '1' {
		return '1'
	}
	if a == '0' && b == '0' {
		return '0'
	}
	return 'x'
}

func xor4(a, b byte) byte {
	if (a != '0' && a != '1') || (b != '0' && b != '1') {
		return 'x'
	}
	if a == b {
		return '0'
	}
	return '1'
}

func not4(a byte) byte {
	switch a {
	case '0':
		return '1'
	case '1':
		return '0'
	}
	return 'x'
}

func randBits4(n int, unknowns int) []byte {
	out := make([]byte, n)
	for i := range out {
		out[i] = edakit.Pick([]byte{'0', '1'})
	}
	for _, i := range edakit.Shuffle(seq(n))[:unknowns] {
		out[i] = edakit.Pick([]byte{'x', 'z'})
	}
	return out
}

func seq(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i
	}
	return out
}

func lit4(bits []byte) string { return edakit.Sprintf("%d'b%s", len(bits), string(bits)) }

// mapBits applies a two-input bit operation lane by lane under a chosen (possibly wrong) model.
func mapBits(a, b []byte, f func(x, y byte) byte) []byte {
	out := make([]byte, len(a))
	for i := range a {
		out[i] = f(a[i], b[i])
	}
	return out
}

// fourState asks for the result of a four-state operation. KP 1 is bitwise AND and OR, KP 2
// XOR, XNOR and inversion, KP 3 reductions and the == versus === comparison.
func fourState(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	n := edakit.Scaled(ctx, 4, 6)
	a := randBits4(n, edakit.Int(1, 2))
	b := randBits4(n, edakit.Int(0, 2))
	switch k {
	case 1, 2:
		var op string
		var f func(x, y byte) byte
		if k == 1 {
			op = edakit.Pick([]string{"&", "|"})
			if op == "&" {
				f = and4
			} else {
				f = or4
			}
		} else {
			op = edakit.Pick([]string{"^", "~^", "& ~"})
			switch op {
			case "^":
				f = xor4
			case "~^":
				f = func(x, y byte) byte { return not4(xor4(x, y)) }
			default:
				f = func(x, y byte) byte { return and4(x, not4(y)) }
			}
		}
		right := mapBits(a, b, f)
		// Misconceptions: z passes straight through; unknowns read as 0; any unknown
		// input poisons the bit even when the other input controls it.
		zThrough := mapBits(a, b, func(x, y byte) byte {
			r := f(x, y)
			if r == 'x' && (x == 'z' || y == 'z') {
				return 'z'
			}
			return r
		})
		asZero := mapBits(a, b, func(x, y byte) byte {
			fix := func(c byte) byte {
				if c == 'x' || c == 'z' {
					return '0'
				}
				return c
			}
			return f(fix(x), fix(y))
		})
		poison := mapBits(a, b, func(x, y byte) byte {
			if x == 'x' || x == 'z' || y == 'x' || y == 'z' {
				return 'x'
			}
			return f(x, y)
		})
		ds := []string{lit4(zThrough), lit4(asZero), lit4(poison)}
		ds = fillBitDistractors(lit4(right), ds, right)
		expr := "a " + op + " b"
		if op == "& ~" {
			expr = "a & ~b"
		}
		q := edakit.Sprintf("In Verilog, a = %s and b = %s. What is the value of %s?", lit4(a), lit4(b), expr)
		return edakit.Choice(q, lit4(right),
			edakit.Sprintf("Work lane by lane. z behaves as x at a logic-operator input, a 0 controls AND and a 1 controls OR even against an unknown, and any unknown into XOR gives x. The result is %s.", lit4(right)),
			ds...)
	}
	// KP 3: reductions and equality.
	kind := edakit.Pick([]string{"&", "|", "^", "==", "==="})
	var res byte
	switch kind {
	case "&":
		res = a[0]
		if res == 'z' {
			res = 'x'
		}
		for _, c := range a[1:] {
			res = and4(res, c)
		}
	case "|":
		res = a[0]
		if res == 'z' {
			res = 'x'
		}
		for _, c := range a[1:] {
			res = or4(res, c)
		}
	case "^":
		res = '0'
		for _, c := range a {
			res = xor4(res, c)
		}
	case "==":
		// Copy a into b in some lanes so equality is decided by the unknowns.
		for i := range b {
			if edakit.Int(0, 2) > 0 {
				b[i] = a[i]
			}
		}
		res = '1'
		known := true
		for i := range a {
			ka := a[i] == '0' || a[i] == '1'
			kb := b[i] == '0' || b[i] == '1'
			if ka && kb && a[i] != b[i] {
				res = '0'
				known = false
				break
			}
		}
		if known {
			for i := range a {
				if a[i] == 'x' || a[i] == 'z' || b[i] == 'x' || b[i] == 'z' {
					res = 'x'
					break
				}
			}
		}
	case "===":
		for i := range b {
			if edakit.Int(0, 2) > 0 {
				b[i] = a[i]
			}
		}
		res = '1'
		for i := range a {
			if a[i] != b[i] {
				res = '0'
				break
			}
		}
	}
	ans := "1'b" + string(res)
	var ds []string
	for _, c := range []string{"0", "1", "x", "z"} {
		if "1'b"+c != ans {
			ds = append(ds, "1'b"+c)
		}
	}
	var q string
	switch kind {
	case "==", "===":
		q = edakit.Sprintf("In Verilog, a = %s and b = %s. What is the value of (a %s b)?", lit4(a), lit4(b), kind)
	default:
		q = edakit.Sprintf("In Verilog, a = %s. What is the value of the reduction %sa?", lit4(a), kind)
	}
	why := map[string]string{
		"&":   "Reduction AND is 0 if any bit is 0, otherwise x if any bit is unknown, otherwise 1.",
		"|":   "Reduction OR is 1 if any bit is 1, otherwise x if any bit is unknown, otherwise 0.",
		"^":   "Reduction XOR is x as soon as any bit is x or z; with only known bits it is the parity.",
		"==":  "Logical equality is 0 if some lane has two known, different bits; otherwise it is x if any lane holds x or z; otherwise 1.",
		"===": "Case equality compares the four-state values literally, so x matches only x and z only z; the result is always 0 or 1.",
	}[kind]
	return edakit.Choice(q, ans, why+" Here the result is "+ans+".", ds...)
}

// fillBitDistractors removes duplicates of the answer and tops the list up with one-lane
// mutations of the correct vector until there are at least three distinct distractors.
func fillBitDistractors(ans string, ds []string, right []byte) []string {
	seen := map[string]bool{ans: true}
	var out []string
	for _, d := range ds {
		if !seen[d] {
			seen[d] = true
			out = append(out, d)
		}
	}
	for tries := 0; len(out) < 3 && tries < 50; tries++ {
		m := append([]byte(nil), right...)
		i := edakit.Int(0, len(m)-1)
		m[i] = edakit.Pick([]byte{'0', '1', 'x', 'z'})
		s := lit4(m)
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

// indent prefixes every line with n spaces.
func indent(n int, lines []string) []string {
	pad := strings.Repeat(" ", n)
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = pad + l
	}
	return out
}
