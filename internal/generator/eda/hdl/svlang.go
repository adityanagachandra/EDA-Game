package hdl

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// logicTypes asks for bit counts of packed arrays (KP 1), packed structs (KP 2) and the
// width of selects into unpacked arrays of packed arrays (KP 3).
func logicTypes(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	switch k {
	case 1:
		dims := []int{edakit.Pick([]int{2, 3, 4, 8}), edakit.Pick([]int{4, 8, 16})}
		if ctx.Difficulty > 0.5 && edakit.Int(0, 1) == 0 {
			dims = append([]int{edakit.Int(2, 4)}, dims...)
		}
		decl := "logic "
		total := 1
		for _, d := range dims {
			decl += edakit.Sprintf("[%d:0]", d-1)
			total *= d
		}
		decl += " bus;"
		q := "In SystemVerilog:" + code(decl) + "How many bits does bus hold in total?"
		return edakit.Numeric(q, edakit.IntStr(total),
			edakit.Sprintf("Packed dimensions multiply into one contiguous vector: %s = %d bits.", joinInts(dims, " x "), total),
			map[string]string{"bits": edakit.IntStr(total)})
	case 2:
		fields := []struct {
			name string
			w    int
		}{
			{"valid", 1}, {"tag", edakit.Pick([]int{3, 4, 6})}, {"data", edakit.Pick([]int{8, 16, 32})}, {"parity", 1}, {"id", edakit.Pick([]int{2, 5})},
		}
		nf := edakit.Int(3, 5)
		chosen := edakit.Shuffle(fields)[:nf]
		sum := 0
		lines := []string{"typedef struct packed {"}
		var ws []int
		for _, f := range chosen {
			if f.w == 1 {
				lines = append(lines, "  logic "+f.name+";")
			} else {
				lines = append(lines, edakit.Sprintf("  logic [%d:0] %s;", f.w-1, f.name))
			}
			sum += f.w
			ws = append(ws, f.w)
		}
		lines = append(lines, "} pkt_t;")
		n := edakit.Pick([]int{1, 2, 4, 8})
		total := sum * n
		ask := "How many bits wide is pkt_t?"
		if n > 1 {
			lines = append(lines, edakit.Sprintf("pkt_t [%d:0] q;", n-1))
			ask = "How many bits does q hold in total?"
		}
		q := "In SystemVerilog:" + code(lines...) + ask
		return edakit.Numeric(q, edakit.IntStr(total),
			edakit.Sprintf("A packed struct is a vector whose width is the sum of its fields, %s = %d bits; a packed array of %d of them holds %d bits.", joinInts(ws, " + "), sum, n, total),
			map[string]string{"struct_bits": edakit.IntStr(sum)})
	}
	a, b := edakit.Pick([]int{2, 4, 8}), edakit.Pick([]int{4, 8, 16})
	depth := edakit.Pick([]int{16, 32, 64, 128})
	decl := edakit.Sprintf("logic [%d:0][%d:0] mem [%d];", a-1, b-1, depth)
	var ask string
	var ans int
	var why string
	switch edakit.Int(0, 2) {
	case 0:
		ask = "How many bits of storage does mem declare in total?"
		ans = a * b * depth
		why = edakit.Sprintf("Each unpacked element is a packed %d x %d = %d-bit vector, and there are %d elements: %d bits.", a, b, a*b, depth, ans)
	case 1:
		i := edakit.Int(0, depth-1)
		ask = edakit.Sprintf("How many bits wide is the expression mem[%d]?", i)
		ans = a * b
		why = edakit.Sprintf("Selecting one unpacked element gives the whole packed vector, %d x %d = %d bits.", a, b, ans)
	default:
		i, j := edakit.Int(0, depth-1), edakit.Int(0, a-1)
		ask = edakit.Sprintf("How many bits wide is the expression mem[%d][%d]?", i, j)
		ans = b
		why = edakit.Sprintf("The first index selects an unpacked element; the second selects along the leftmost packed dimension, leaving one %d-bit row.", b)
	}
	q := "In SystemVerilog:" + code(decl) + ask
	return edakit.Numeric(q, edakit.IntStr(ans), why, map[string]string{"answer": edakit.IntStr(ans)})
}

func joinInts(xs []int, sep string) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = edakit.IntStr(x)
	}
	return strings.Join(parts, sep)
}

type vStmt struct {
	target   string
	variable bool
	a, b     string // b may be a decimal constant
}

// vhdlSignalsVariables simulates a clocked VHDL process (KP 1 one edge, KP 2 two edges) or
// counts delta cycles through a chain of concurrent assignments (KP 3).
func vhdlSignalsVariables(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	if k == 3 {
		return vhdlDeltas(ctx)
	}
	edges := k
	sigs := []string{"a", "b", "c"}
	for tries := 0; ; tries++ {
		init := map[string]int{}
		for _, s := range sigs {
			init[s] = edakit.Int(1, 9)
		}
		operand := func(pool []string) string {
			if edakit.Int(0, 2) == 0 {
				return edakit.IntStr(edakit.Int(1, 4))
			}
			return edakit.Pick(pool)
		}
		var stmts []vStmt
		stmts = append(stmts, vStmt{target: "v", variable: true, a: edakit.Pick(sigs), b: operand(sigs)})
		pool := append([]string{"v"}, sigs...)
		order := edakit.Shuffle(sigs)
		for i, s := range order {
			stmts = append(stmts, vStmt{target: s, a: edakit.Pick(pool), b: operand(pool)})
			if i == 0 && edakit.Int(0, 1) == 0 {
				stmts = append(stmts, vStmt{target: "v", variable: true, a: "v", b: operand(sigs)})
			}
		}
		run := func(immediate bool) map[string]int {
			cur := map[string]int{}
			for kk, vv := range init {
				cur[kk] = vv
			}
			for e := 0; e < edges; e++ {
				pending := map[string]int{}
				vars := map[string]int{"v": cur["v"]}
				val := func(x string) int {
					if x == "v" {
						return vars["v"]
					}
					if x[0] >= '0' && x[0] <= '9' {
						n := 0
						for _, ch := range x {
							n = n*10 + int(ch-'0')
						}
						return n
					}
					return cur[x]
				}
				for _, s := range stmts {
					r := val(s.a) + val(s.b)
					if s.variable {
						vars["v"] = r
					} else if immediate {
						cur[s.target] = r
					} else {
						pending[s.target] = r
					}
				}
				for kk, vv := range pending {
					cur[kk] = vv
				}
				cur["v"] = vars["v"]
			}
			return cur
		}
		ask := edakit.Pick(sigs)
		got := run(false)[ask]
		alt := run(true)[ask]
		if got == alt && tries < 40 {
			continue
		}
		lines := []string{"process (clk)", "  variable v : integer := 0;", "begin", "  if rising_edge(clk) then"}
		for _, s := range stmts {
			op := "<="
			if s.variable {
				op = ":="
			}
			lines = append(lines, edakit.Sprintf("    %s %s %s + %s;", s.target, op, s.a, s.b))
		}
		lines = append(lines, "  end if;", "end process;")
		edgeWord := "one rising edge"
		if edges == 2 {
			edgeWord = "two rising edges"
		}
		q := edakit.Sprintf("Signals a, b and c are integers starting at a = %d, b = %d, c = %d, and only this process drives them.", init["a"], init["b"], init["c"]) +
			code(lines...) + edakit.Sprintf("What is the value of signal %s after %s, once all delta cycles have settled?", ask, edgeWord)
		why := edakit.Sprintf("A variable assignment (:=) takes effect at once; a signal assignment (<=) is scheduled and lands only when the process suspends, so every read of a signal inside the process sees the value from before the edge. Following that for %s gives %s = %d. Updating signals immediately would give %d.", edgeWord, ask, got, alt)
		return edakit.Numeric(q, edakit.IntStr(got), why, map[string]string{"naive": edakit.IntStr(alt)})
	}
}

func vhdlDeltas(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 6)
	names := make([]string, n)
	for i := range names {
		names[i] = edakit.Sprintf("s%d", i+1)
	}
	delta := map[string]int{"x": 0}
	var lines []string
	for i, s := range names {
		pool := append([]string{"x"}, names[:i]...)
		d1 := edakit.Pick(pool)
		var rhs string
		dl := delta[d1]
		if i > 0 && edakit.Int(0, 2) > 0 {
			d2 := edakit.Pick(pool)
			for d2 == d1 && len(pool) > 1 {
				d2 = edakit.Pick(pool)
			}
			rhs = d1 + " " + edakit.Pick([]string{"and", "or", "xor"}) + " " + d2
			dl = max(dl, delta[d2])
		} else {
			rhs = "not " + d1
		}
		delta[s] = dl + 1
		lines = append(lines, s+" <= "+rhs+";")
	}
	shuffled := edakit.Shuffle(lines)
	ask := names[edakit.Int(n/2, n-1)]
	ans := delta[ask]
	q := "These concurrent signal assignments sit in one VHDL architecture, written in arbitrary order. Input x changes in delta cycle 0, and assume every signal whose inputs change takes a new value." +
		code(shuffled...) + edakit.Sprintf("In which delta cycle does %s take its final value for this change?", ask)
	why := edakit.Sprintf("Each concurrent assignment is a process that wakes in the delta after one of its inputs changes, so a signal settles one delta after the latest of its inputs: delta(s) = 1 + max over its inputs. Along the longest chain from x this gives %s = %d. Textual order does not matter.", ask, ans)
	return edakit.Numeric(q, edakit.IntStr(ans), why, map[string]string{"delta": edakit.IntStr(ans)})
}

// std_logic resolution table from IEEE 1164, indexed by "UX01ZWLH-".
const stdLogic = "UX01ZWLH-"

var resolveTable = [9]string{
	"UUUUUUUUU", // U
	"UXXXXXXXX", // X
	"UX0X0000X", // 0
	"UXX11111X", // 1
	"UX01ZWLHX", // Z
	"UX01WWWWX", // W
	"UX01LWLWX", // L
	"UX01HWWHX", // H
	"UXXXXXXXX", // -
}

func resolve(a, b byte) byte {
	i := strings.IndexByte(stdLogic, a)
	j := strings.IndexByte(stdLogic, b)
	return resolveTable[i][j]
}

var stdLogicName = map[byte]string{
	'U': "'U' (uninitialized)", 'X': "'X' (forcing unknown)", '0': "'0' (forcing 0)", '1': "'1' (forcing 1)",
	'Z': "'Z' (high impedance)", 'W': "'W' (weak unknown)", 'L': "'L' (weak 0)", 'H': "'H' (weak 1)", '-': "'-' (don't care)",
}

// vhdlResolution asks for the resolved value of a std_logic signal with several drivers.
func vhdlResolution(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	var pool string
	var n int
	switch k {
	case 1:
		pool, n = "01Z", 2
	case 2:
		pool, n = "01ZLHW", 2
	default:
		pool, n = "01ZLHWX", 3
	}
	drivers := make([]byte, n)
	for {
		for i := range drivers {
			drivers[i] = pool[edakit.Int(0, len(pool)-1)]
		}
		if k == 1 || drivers[0] != drivers[1] {
			break
		}
	}
	res := drivers[0]
	for _, d := range drivers[1:] {
		res = resolve(res, d)
	}
	var shown []string
	for _, d := range drivers {
		shown = append(shown, "'"+string(d)+"'")
	}
	ans := stdLogicName[res]
	var ds []string
	for _, c := range []byte("01XZWLH") {
		if c != res {
			ds = append(ds, stdLogicName[c])
		}
	}
	ds = edakit.Shuffle(ds)[:4]
	q := edakit.Sprintf("A std_logic signal has %d drivers, currently driving %s. What value does the resolution function produce?", n, strings.Join(shown, ", "))
	why := "The IEEE 1164 resolution function is commutative and associative: forcing values beat weak ones, weak beats Z, two different values of the same strength give the unknown of that strength ('X' or 'W'), and 'X' dominates everything except 'U'. The result is " + ans + "."
	return edakit.Choice(q, ans, why, ds...)
}
