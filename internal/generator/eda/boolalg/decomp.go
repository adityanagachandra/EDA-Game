package boolalg

import (
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// funcOnSupport draws a function of n variables that depends exactly on the variables in sup.
func funcOnSupport(n int, sup []int) tt {
	for {
		f := tt{n, 0}
		k := len(sup)
		g := rand.Uint32() & fullMask(k)
		for i := 0; i < rows(n); i++ {
			j := 0
			for _, v := range sup {
				j = j<<1 | bitOf(i, v, n)
			}
			if g>>uint(j)&1 == 1 {
				f.w |= 1 << uint(i)
			}
		}
		if f.supportSize() == k {
			return f
		}
	}
}

func support(ctx generator.GeneratorContext) generator.Problem {
	kp := kpOf(ctx)
	switch kp {
	case 1:
		n := 4
		k := edakit.Int(1, 4)
		sup := rand.Perm(n)[:k]
		f := funcOnSupport(n, sup)
		var names []string
		for v := 0; v < n; v++ {
			if f.dependsOn(v) {
				names = append(names, varNames[v])
			}
		}
		ans := edakit.IntStr(k)
		q := "Let " + sigma("f", f) + ". How many variables are in the true support of $f$, the variables whose Boolean difference is not identically 0?"
		return edakit.Numeric(q, ans, "Compare the two cofactors of each variable. They differ only for "+strings.Join(names, ", ")+", so the support has "+ans+" variables.", map[string]string{"support": strings.Join(names, ",")})
	default:
		if kp == 2 {
			n := nVars(ctx)
			for {
				s := randSOP(n, edakit.Int(2, 3))
				// Split one or two cubes on a variable they do not mention, which adds a
				// variable to the expression without adding it to the function.
				for t := 0; t < edakit.Int(1, 2); t++ {
					i := rand.Intn(len(s))
					var free []int
					for v := 0; v < n; v++ {
						if s[i][v] == 2 {
							free = append(free, v)
						}
					}
					if len(free) == 0 || s[i].lits() >= n-1 {
						continue
					}
					v := edakit.Pick(free)
					c0 := append(cube(nil), s[i]...)
					c1 := append(cube(nil), s[i]...)
					c0[v], c1[v] = 0, 1
					s = append(append(append(sop(nil), s[:i]...), c0, c1), s[i+1:]...)
				}
				f := s.tt(n)
				if f.isConst() {
					continue
				}
				mentioned, vac := 0, 0
				var vacNames []string
				for v := 0; v < n; v++ {
					if s.mentions(v) {
						mentioned++
						if !f.dependsOn(v) {
							vac++
							vacNames = append(vacNames, varNames[v])
						}
					}
				}
				ans := edakit.IntStr(vac)
				q := edakit.Sprintf("The expression $f = %s$ mentions %d variables. On how many of them does the function $f$ not actually depend?", s.latex(), mentioned)
				list := "none"
				if vac > 0 {
					list = strings.Join(vacNames, ", ")
				}
				return edakit.Numeric(q, ans, "A variable is vacuous when its two cofactors are equal. Here that holds for: "+list+". Terms like x y + x y-bar merge to x, so a variable can appear in a cover without being in the support.", map[string]string{"vacuous": list})
			}
		}
		// KP 3: functional dependence of g on (f1, f2).
		n := 3
		var f1, f2, g tt
		for {
			f1 = tt{n, rand.Uint32() & fullMask(n)}
			f2 = tt{n, rand.Uint32() & fullMask(n)}
			if f1.isConst() || f2.isConst() || f1 == f2 || f1 == f2.not() {
				continue
			}
			if rand.Intn(2) == 0 {
				h := rand.Intn(16)
				g = tt{n, 0}
				for i := 0; i < rows(n); i++ {
					j := 0
					if f1.at(i) {
						j |= 2
					}
					if f2.at(i) {
						j |= 1
					}
					if h>>uint(j)&1 == 1 {
						g.w |= 1 << uint(i)
					}
				}
			} else {
				g = tt{n, rand.Uint32() & fullMask(n)}
			}
			if !g.isConst() {
				break
			}
		}
		var seen0, seen1 [4]bool
		for i := 0; i < rows(n); i++ {
			j := 0
			if f1.at(i) {
				j |= 2
			}
			if f2.at(i) {
				j |= 1
			}
			if g.at(i) {
				seen1[j] = true
			} else {
				seen0[j] = true
			}
		}
		conflicts := 0
		for j := 0; j < 4; j++ {
			if seen0[j] && seen1[j] {
				conflicts++
			}
		}
		ans := edakit.IntStr(conflicts)
		q := "Over $a, b, c$: " + sigma("f_1", f1) + ", " + sigma("f_2", f2) + " and " + sigma("g", g) +
			". Group the 8 input rows by the value pair $(f_1, f_2)$. How many of the four value pairs occur with both $g = 0$ and $g = 1$? (Zero means $g$ is a function of $f_1$ and $f_2$.)"
		verdict := "g can be written as h(f1, f2)"
		if conflicts > 0 {
			verdict = "g is not a function of f1 and f2"
		}
		return edakit.Numeric(q, ans, edakit.Sprintf("Tabulate (f1, f2, g) per row. %s value pair(s) see both g values, so %s.", ans, verdict), map[string]string{"conflicts": ans})
	}
}

// multiplicity returns the column multiplicity of f (n = 4) with bound set {p, q}.
func multiplicity(f tt, p, q int) int {
	cols := map[uint32]bool{}
	for a := 0; a < 2; a++ {
		for b := 0; b < 2; b++ {
			lo, hi := p, q
			if lo > hi {
				lo, hi = hi, lo
			}
			// Cofactor the higher index first so the lower index stays valid. Every pair of
			// values is visited, so the set of columns does not depend on which value goes where.
			c := f.cof(hi, b).cof(lo, a)
			cols[c.w] = true
		}
	}
	return len(cols)
}

func structuredDecomposable() tt {
	n := 4
	perm := rand.Perm(n)
	x, y, z, w := perm[0], perm[1], perm[2], perm[3]
	gfun := edakit.Pick([]int{0x6, 0x8, 0xE, 0x9, 0x1, 0x7, 0x2, 0x4}) // g(x, y) as a 4-bit table
	h := rand.Intn(256)                                                // h(g, z, w)
	f := tt{n, 0}
	for i := 0; i < rows(n); i++ {
		gi := (gfun >> uint(bitOf(i, x, n)<<1|bitOf(i, y, n))) & 1
		j := gi<<2 | bitOf(i, z, n)<<1 | bitOf(i, w, n)
		if h>>uint(j)&1 == 1 {
			f.w |= 1 << uint(i)
		}
	}
	return f
}

func disjointDecomp(ctx generator.GeneratorContext) generator.Problem {
	n := 4
	var f tt
	for {
		if rand.Intn(2) == 0 {
			f = structuredDecomposable()
		} else {
			f, _ = drawFunc(n, 0)
		}
		if f.supportSize() >= 3 {
			break
		}
	}
	kp := kpOf(ctx)
	switch kp {
	case 1, 2:
		p, q := 0, 0
		for p == q {
			p, q = rand.Intn(n), rand.Intn(n)
		}
		if p > q {
			p, q = q, p
		}
		mu := multiplicity(f, p, q)
		bits := 0
		for (1 << uint(bits)) < mu {
			bits++
		}
		free := otherVars(n, p, q)
		if kp == 2 {
			ans := edakit.IntStr(bits)
			qs := edakit.Sprintf("Let %s. Decompose $f$ with bound set $\\{%s, %s\\}$ and free set $\\{%s\\}$, as $f = h(g_1, \\dots, g_k, %s)$ with each $g_i$ a function of $%s, %s$ only. What is the smallest possible $k$? (Answer 0 if $f$ does not depend on the bound set.)",
				sigma("f", f), varNames[p], varNames[q], free, free, varNames[p], varNames[q])
			return edakit.Numeric(qs, ans, edakit.Sprintf("The decomposition chart has column multiplicity %d: the four assignments of (%s, %s) give %d distinct functions of the free variables. Encoding %d classes needs the ceiling of log2 %d = %s intermediate signals.", mu, varNames[p], varNames[q], mu, mu, mu, ans), map[string]string{"multiplicity": edakit.IntStr(mu)})
		}
		ans := edakit.IntStr(mu)
		qs := edakit.Sprintf("Let %s. Draw the decomposition chart with bound set $\\{%s, %s\\}$ (one column per assignment of the bound set) and free set $\\{%s\\}$. What is the column multiplicity, the number of distinct columns?",
			sigma("f", f), varNames[p], varNames[q], free)
		return edakit.Numeric(qs, ans, edakit.Sprintf("Each column is the cofactor of f with respect to one assignment of (%s, %s), a function of %s. There are %s distinct cofactors. A simple disjoint decomposition f = h(g(%s, %s), %s) exists exactly when this is at most 2.", varNames[p], varNames[q], free, ans, varNames[p], varNames[q], free), map[string]string{"multiplicity": ans})
	default:
		count := 0
		var good []string
		for p := 0; p < n; p++ {
			for q := p + 1; q < n; q++ {
				if multiplicity(f, p, q) <= 2 {
					count++
					good = append(good, "{"+varNames[p]+","+varNames[q]+"}")
				}
			}
		}
		ans := edakit.IntStr(count)
		list := "none"
		if count > 0 {
			list = strings.Join(good, ", ")
		}
		qs := "Let " + sigma("f", f) + ". For how many of the 6 two-variable bound sets does $f$ have a simple disjoint decomposition $f = h(g(\\text{bound set}), \\text{free set})$?"
		return edakit.Numeric(qs, ans, "Build the chart for each bound set and keep those with column multiplicity at most 2: "+list+".", map[string]string{"bound_sets": list})
	}
}
