package techmap

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Truth tables over n inputs: bit m is f at the assignment whose binary code is m, with x1 the
// least significant bit.

var perms3 = [][]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
var perms2 = [][]int{{0, 1}, {1, 0}}

// transform returns g(x) = out XOR f(y) with y_j = x_{p[j]} XOR neg_j.
func transform(tt uint32, n int, p []int, neg int, out bool) uint32 {
	var g uint32
	for x := 0; x < 1<<uint(n); x++ {
		y := 0
		for j := 0; j < n; j++ {
			b := (x >> uint(p[j])) & 1
			b ^= (neg >> uint(j)) & 1
			y |= b << uint(j)
		}
		v := (tt >> uint(y)) & 1
		if out {
			v ^= 1
		}
		g |= v << uint(x)
	}
	return g
}

func permsFor(n int) [][]int {
	if n == 2 {
		return perms2
	}
	return perms3
}

// orbit returns the set of functions reachable from tt under the chosen group.
func orbit(tt uint32, n int, usePerm, useNeg, useOut bool) map[uint32]bool {
	set := map[uint32]bool{}
	ps := [][]int{[]int{0, 1, 2}[:n]}
	if usePerm {
		ps = permsFor(n)
	}
	negs := []int{0}
	if useNeg {
		negs = nil
		for i := 0; i < 1<<uint(n); i++ {
			negs = append(negs, i)
		}
	}
	outs := []bool{false}
	if useOut {
		outs = []bool{false, true}
	}
	for _, p := range ps {
		for _, ng := range negs {
			for _, o := range outs {
				set[transform(tt, n, p, ng, o)] = true
			}
		}
	}
	return set
}

func minOf(s map[uint32]bool) uint32 {
	m := uint32(1 << 31)
	for x := range s {
		if x < m {
			m = x
		}
	}
	return m
}

func ttFrom(n int, f func(x []bool) bool) uint32 {
	var tt uint32
	for m := 0; m < 1<<uint(n); m++ {
		x := make([]bool, n)
		for j := 0; j < n; j++ {
			x[j] = (m>>uint(j))&1 == 1
		}
		if f(x) {
			tt |= 1 << uint(m)
		}
	}
	return tt
}

func minterms(tt uint32, n int) string {
	var ms []int
	for m := 0; m < 1<<uint(n); m++ {
		if (tt>>uint(m))&1 == 1 {
			ms = append(ms, m)
		}
	}
	return "{" + joinInts(ms) + "}"
}

type named3 struct {
	name string
	f    func(x []bool) bool
}

var named3Funcs = []named3{
	{"x1 AND x2 AND x3", func(x []bool) bool { return x[0] && x[1] && x[2] }},
	{"x1 OR x2 OR x3", func(x []bool) bool { return x[0] || x[1] || x[2] }},
	{"x1 XOR x2 XOR x3", func(x []bool) bool { return x[0] != x[1] != x[2] }},
	{"MAJ(x1, x2, x3)", func(x []bool) bool { return (x[0] && x[1]) || (x[0] && x[2]) || (x[1] && x[2]) }},
	{"the multiplexer x3 ? x2 : x1", func(x []bool) bool {
		if x[2] {
			return x[1]
		}
		return x[0]
	}},
	{"x1 AND x2 (x3 unused)", func(x []bool) bool { return x[0] && x[1] }},
	{"x1 XOR x2 (x3 unused)", func(x []bool) bool { return x[0] != x[1] }},
	{"x1 AND (x2 OR x3)", func(x []bool) bool { return x[0] && (x[1] || x[2]) }},
	{"x1 OR (x2 AND x3)", func(x []bool) bool { return x[0] || (x[1] && x[2]) }},
	{"x1 XOR (x2 AND x3)", func(x []bool) bool { return x[0] != (x[1] && x[2]) }},
	{"x1 (x2 and x3 unused)", func(x []bool) bool { return x[0] }},
	{"(x1 AND x2) OR (NOT x1 AND x3)", func(x []bool) bool { return (x[0] && x[1]) || (!x[0] && x[2]) }},
}

// npnClassStats enumerates the NPN classes of n-input functions (n = 2 or 3).
func npnClassStats(n int) (classes int, byMinOnes map[int]int, fullSupport int) {
	seen := map[uint32]bool{}
	byMinOnes = map[int]int{}
	total := 1 << uint(1<<uint(n))
	for tt := 0; tt < total; tt++ {
		t := uint32(tt)
		if seen[t] {
			continue
		}
		orb := orbit(t, n, true, true, true)
		for x := range orb {
			seen[x] = true
		}
		classes++
		ones := popcount(t)
		if c := (1 << uint(n)) - ones; c < ones {
			ones = c
		}
		byMinOnes[ones]++
		if dependsOnAll(t, n) {
			fullSupport++
		}
	}
	return
}

func dependsOnAll(tt uint32, n int) bool {
	for j := 0; j < n; j++ {
		dep := false
		for m := 0; m < 1<<uint(n); m++ {
			if (m>>uint(j))&1 == 0 {
				if (tt>>uint(m))&1 != (tt>>uint(m|1<<uint(j)))&1 {
					dep = true
				}
			}
		}
		if !dep {
			return false
		}
	}
	return true
}

func factorial(n int) int {
	r := 1
	for i := 2; i <= n; i++ {
		r *= i
	}
	return r
}

func npnClasses(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		if edakit.Int(0, 1) == 0 {
			n := edakit.Int(2, 6)
			g := intPow(2, n+1) * factorial(n)
			return edakit.Numeric(edakit.Sprintf("How many NPN transformations act on %d-input Boolean functions, counting every combination of input permutation, input negations and output negation?", n),
				edakit.IntStr(g),
				edakit.Sprintf("There are %d! = %d input permutations, 2^%d = %d input-negation patterns and 2 output polarities: %d x %d x 2 = %d.", n, factorial(n), n, intPow(2, n), factorial(n), intPow(2, n), g),
				map[string]string{"group_order": edakit.IntStr(g)})
		}
		n := edakit.Int(1, 4)
		cnt := intPow(2, intPow(2, n))
		return edakit.Numeric(edakit.Sprintf("How many distinct Boolean functions of %d inputs exist?", n),
			edakit.IntStr(cnt),
			edakit.Sprintf("A function is a truth table of 2^%d = %d rows, each 0 or 1, so there are 2^%d = %d functions.", n, intPow(2, n), intPow(2, n), cnt),
			map[string]string{"functions": edakit.IntStr(cnt)})
	case 2:
		f := edakit.Pick(named3Funcs)
		tt := ttFrom(3, f.f)
		type grp struct {
			name    string
			p, n, o bool
			size    int
		}
		g := edakit.Pick([]grp{
			{"input negations only (8 transformations)", false, true, false, 8},
			{"input permutations only (6 transformations)", true, false, false, 6},
			{"input permutations and input negations (48 transformations)", true, true, false, 48},
			{"the full NPN group (96 transformations)", true, true, true, 96},
		})
		orb := orbit(tt, 3, g.p, g.n, g.o)
		return edakit.Numeric("Consider the 3-input function f = "+f.name+" (truth table "+minterms(tt, 3)+" as minterms, with x1 the least significant bit). Under "+g.name+", how many distinct functions are in the orbit of f, f itself included?",
			edakit.IntStr(len(orb)),
			edakit.Sprintf("Applying all %d transformations and removing duplicates leaves %d distinct functions. By orbit-stabilizer, %d divided by the number of transformations that fix f gives the same count, so %d transformations fix f.", g.size, len(orb), g.size, g.size/len(orb)),
			map[string]string{"orbit": edakit.IntStr(len(orb))})
	}
	switch edakit.Int(0, 2) {
	case 0:
		n := edakit.Int(1, 4)
		known := map[int]int{1: 2, 2: 4, 3: 14, 4: 222}
		return edakit.Numeric(edakit.Sprintf("Into how many NPN equivalence classes do the %d Boolean functions of %d inputs fall?", intPow(2, intPow(2, n)), n),
			edakit.IntStr(known[n]),
			edakit.Sprintf("The NPN class counts for n = 1, 2, 3, 4 are 2, 4, 14 and 222. For n = %d the answer is %d, which is why a mapper can index a whole small-input library by canonical form.", n, known[n]),
			map[string]string{"classes": edakit.IntStr(known[n])})
	case 1:
		n := edakit.Int(2, 3)
		_, _, full := npnClassStats(n)
		return edakit.Numeric(edakit.Sprintf("How many NPN classes of %d-input functions contain functions that depend on all %d inputs?", n, n),
			edakit.IntStr(full),
			edakit.Sprintf("Of the NPN classes of %d-input functions, the others are degenerate: constants and functions of fewer inputs. Enumerating classes and testing support leaves %d that use every input.", n, full),
			map[string]string{"full_support_classes": edakit.IntStr(full)})
	}
	_, by, _ := npnClassStats(3)
	t := edakit.Int(0, 4)
	stem := edakit.Sprintf("Output negation maps a function with j ones to one with 8 - j ones, so every NPN class of 3-input functions has members with j and with 8 - j ones for some j. How many classes have members with exactly %d ones (and so also %d)?", t, 8-t)
	if t == 4 {
		stem = "Output negation maps a function with j ones to one with 8 - j ones. How many NPN classes of 3-input functions consist of functions with exactly 4 ones?"
	}
	return edakit.Numeric(stem,
		edakit.IntStr(by[t]),
		edakit.Sprintf("Grouping the 14 classes by the smaller of their ones count and its complement gives %d classes for %d ones.", by[t], t),
		map[string]string{"classes": edakit.IntStr(by[t])})
}

func litName(j, neg int) string {
	s := edakit.Sprintf("x%d", j+1)
	if neg != 0 {
		return "NOT " + s
	}
	return s
}

func randTT3() uint32 {
	for {
		t := uint32(edakit.Int(1, 254))
		if dependsOnAll(t, 3) {
			return t
		}
	}
}

func canonicalForms(ctx generator.GeneratorContext) generator.Problem {
	conv := "Truth tables are written as integers: bit m is the value of the function on the assignment whose binary code is m, with x1 the least significant bit."
	switch kp(ctx) {
	case 1:
		tt := randTT3()
		p := edakit.Pick(perms3)
		neg := edakit.Int(0, 7)
		out := edakit.Int(0, 3) == 0
		if p[0] == 0 && p[1] == 1 && neg == 0 && !out {
			neg = 1
		}
		g := transform(tt, 3, p, neg, out)
		var lits []string
		for j := 0; j < 3; j++ {
			lits = append(lits, litName(p[j], (neg>>uint(j))&1))
		}
		expr := "f(" + strings.Join(lits, ", ") + ")"
		if out {
			expr = "NOT " + expr
		}
		return edakit.Numeric(conv+edakit.Sprintf(" The function f has truth table %d (ones at minterms %s). Define g(x1, x2, x3) = %s. What is the truth table of g as an integer?", tt, minterms(tt, 3), expr),
			edakit.IntStr(int(g)),
			edakit.Sprintf("For each assignment of x1, x2, x3, evaluate the arguments, look up f, and set the bit of g. The ones of g are at minterms %s, so g = %d.", minterms(g, 3), g),
			map[string]string{"tt_f": edakit.IntStr(int(tt)), "tt_g": edakit.IntStr(int(g))})
	case 2:
		if edakit.Int(0, 1) == 0 {
			tt := randTT3()
			orb := orbit(tt, 3, true, false, false)
			c := minOf(orb)
			var all []int
			for x := range orb {
				all = append(all, int(x))
			}
			sort.Ints(all)
			return edakit.Numeric(conv+edakit.Sprintf(" The P-canonical form of a function is the smallest truth-table integer over all permutations of its inputs. f has truth table %d (ones at %s). What is its P-canonical form?", tt, minterms(tt, 3)),
				edakit.IntStr(int(c)),
				edakit.Sprintf("The six permutations give the distinct tables %s. The smallest is %d.", joinInts(all), c),
				map[string]string{"canonical": edakit.IntStr(int(c))})
		}
		tt := uint32(edakit.Int(1, 14))
		orb := orbit(tt, 2, true, true, true)
		c := minOf(orb)
		var all []int
		for x := range orb {
			all = append(all, int(x))
		}
		sort.Ints(all)
		return edakit.Numeric(conv+edakit.Sprintf(" The NPN-canonical form is the smallest truth-table integer over all 16 NPN transformations of a 2-input function. f(x1, x2) has truth table %d (ones at %s). What is its NPN-canonical form?", tt, minterms(tt, 2)),
			edakit.IntStr(int(c)),
			edakit.Sprintf("The orbit of f is %s, and the smallest member is %d. Two functions match exactly when their canonical forms are equal.", joinInts(all), c),
			map[string]string{"canonical": edakit.IntStr(int(c))})
	}
	tt := randTT3()
	j := edakit.Int(0, 2)
	pos := edakit.Int(0, 1) == 1
	ones := 0
	for m := 0; m < 8; m++ {
		bit := (m >> uint(j)) & 1
		if (bit == 1) == pos && (tt>>uint(m))&1 == 1 {
			ones++
		}
	}
	val := "0"
	if pos {
		val = "1"
	}
	return edakit.Numeric(conv+edakit.Sprintf(" A Boolean matcher prunes candidates with cofactor signatures. f has truth table %d (ones at %s). How many ones does the cofactor of f with x%d = %s have, as a function of the other two inputs?", tt, minterms(tt, 3), j+1, val),
		edakit.IntStr(ones),
		edakit.Sprintf("Keep the minterms of f whose bit for x%d is %s: there are %d. The multiset of such counts over all inputs and polarities is invariant under permutation, so two functions with different signatures cannot be P-equivalent.", j+1, val, ones),
		map[string]string{"cofactor_ones": edakit.IntStr(ones)})
}

func npnEquivalence(ctx generator.GeneratorContext) generator.Problem {
	for {
		f := randTT3()
		orb := orbit(f, 3, true, true, true)
		canon := minOf(orb)
		var members []uint32
		for x := range orb {
			if x != f {
				members = append(members, x)
			}
		}
		if len(members) == 0 {
			continue
		}
		sort.Slice(members, func(i, j int) bool { return members[i] < members[j] })
		g := edakit.Pick(members)
		ones := popcount(f)
		var dis []string
		seen := map[uint32]bool{}
		for tries := 0; tries < 400 && len(dis) < 3; tries++ {
			c := uint32(edakit.Int(1, 254))
			pc := popcount(c)
			if tries < 300 && pc != ones && pc != 8-ones {
				continue
			}
			if seen[c] || minOf(orbit(c, 3, true, true, true)) == canon {
				continue
			}
			seen[c] = true
			dis = append(dis, minterms(c, 3))
		}
		if len(dis) < 3 {
			continue
		}
		q := "Minterms are numbered with x1 as the least significant bit. The 3-input function f has ones at minterms " + minterms(f, 3) + ". Which of these functions, given by their ones, is NPN-equivalent to f?"
		return edakit.Choice(q, minterms(g, 3),
			edakit.Sprintf("The NPN-canonical form (smallest truth-table integer over all 96 transformations) of f is %d, and the correct option has the same canonical form. The others share f's ones count or its complement's but fall in different classes.", canon),
			dis...)
	}
}

type sgGate struct {
	name  string
	pins  int
	area  int
	delay []int
}

func sgLibrary() []sgGate {
	base := []sgGate{
		{"INV", 1, 2, []int{8}},
		{"NAND2", 2, 3, []int{10, 12}},
		{"NOR2", 2, 3, []int{12, 15}},
		{"NAND3", 3, 4, []int{12, 14, 16}},
		{"AOI21", 3, 4, []int{14, 16, 11}},
	}
	for i := range base {
		base[i].area += edakit.Int(0, 1)
		for j := range base[i].delay {
			base[i].delay[j] += edakit.Int(-2, 2)
		}
	}
	return base
}

func supergates(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	lib := sgLibrary()
	roots := []sgGate{lib[1], lib[2], lib[3], lib[4]}
	root := edakit.Pick(roots)
	nsub := 1
	if root.pins >= 2 && edakit.Int(0, 1) == 1 {
		nsub = 2
	}
	pinIdx := edakit.Shuffle([]int{0, 1, 2}[:root.pins])[:nsub]
	sort.Ints(pinIdx)
	subs := map[int]sgGate{}
	for _, p := range pinIdx {
		subs[p] = edakit.Pick(lib[:3])
	}
	pinName := []string{"A", "B", "C"}
	next := 0
	area := root.area
	var parts []string
	type input struct {
		name string
		path int
		arr  int
	}
	var ins []input
	for p := 0; p < root.pins; p++ {
		if s, ok := subs[p]; ok {
			area += s.area
			var nm []string
			for q := 0; q < s.pins; q++ {
				in := input{piName(next), s.delay[q] + root.delay[p], edakit.Int(0, 30)}
				next++
				ins = append(ins, in)
				nm = append(nm, in.name)
			}
			word := "inputs"
			if len(nm) == 1 {
				word = "input"
			}
			parts = append(parts, edakit.Sprintf("pin %s is driven by %s %s with %s %s", pinName[p], article(s.name), s.name, word, strings.Join(nm, ", ")))
		} else {
			in := input{piName(next), root.delay[p], edakit.Int(0, 30)}
			next++
			ins = append(ins, in)
			parts = append(parts, edakit.Sprintf("pin %s is input %s", pinName[p], in.name))
		}
	}
	var lib2 []string
	for _, g := range lib {
		var ds []string
		for j, d := range g.delay {
			ds = append(ds, edakit.Sprintf("%s %d", pinName[j], d))
		}
		lib2 = append(lib2, edakit.Sprintf("%s (area %d; pin-to-output delays in ps: %s)", g.name, g.area, strings.Join(ds, ", ")))
	}
	q := "Library: " + strings.Join(lib2, "; ") + ". A supergate has " + article(root.name) + " " + root.name + " at its output, where " + strings.Join(parts, "; ") + "."
	switch k {
	case 1:
		return edakit.Numeric(q+" What is the area of the supergate?", edakit.IntStr(area),
			edakit.Sprintf("A supergate is a fixed composition of library gates, so its area is the sum of its component areas: %d.", area),
			map[string]string{"area": edakit.IntStr(area)})
	case 2:
		var arr []string
		best, bestIn := 0, ""
		var why []string
		for _, in := range ins {
			arr = append(arr, edakit.Sprintf("%s at %d ps", in.name, in.arr))
			t := in.arr + in.path
			why = append(why, edakit.Sprintf("%s: %d + %d = %d", in.name, in.arr, in.path, t))
			if t > best {
				best, bestIn = t, in.name
			}
		}
		return edakit.Numeric(q+" Inputs arrive "+strings.Join(arr, ", ")+". With these load-independent pin delays, when does the supergate output settle, in ps?", edakit.IntStr(best),
			"Each supergate input's delay is the sum of pin delays along its path to the output. "+strings.Join(why, "; ")+". The latest is "+bestIn+edakit.Sprintf(", at %d ps.", best),
			map[string]string{"arrival": edakit.IntStr(best)})
	}
	return edakit.Numeric(q+" How many inputs does the supergate have?", edakit.IntStr(len(ins)),
		edakit.Sprintf("The root contributes its %d pins, each internal gate replaces the pin it drives with its own inputs: %d inputs in all. Supergate generators bound this count (and the number of levels) so the matching library stays small.", root.pins, len(ins)),
		map[string]string{"inputs": edakit.IntStr(len(ins))})
}

func article(name string) string {
	if strings.ContainsRune("AEIOU", rune(name[0])) {
		return "an"
	}
	return "a"
}
