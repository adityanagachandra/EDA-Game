package boolalg

import (
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func cubeList(cs []cube) string {
	parts := make([]string, len(cs))
	for i, c := range cs {
		parts[i] = c.positional()
	}
	return strings.Join(parts, ", ")
}

func cubes(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 7)
		if rand.Intn(2) == 0 {
			c := randCube(n, 1, n-1)
			ans := edakit.IntStr(c.size())
			q := edakit.Sprintf("In positional notation over %d variables (one position per variable, 0 for a complemented literal, 1 for a true literal, - for an absent variable), how many minterms does the cube %s contain?", n, c.positional())
			return edakit.Numeric(q, ans, edakit.Sprintf("The cube has %d dashes, each free to be 0 or 1, so it contains 2^%d = %s minterms.", c.dashes(), c.dashes(), ans), nil)
		}
		k := edakit.Int(1, n-1)
		ans := edakit.IntStr(binom(n, k) * pow2(k))
		q := edakit.Sprintf("How many distinct cubes over %d variables have exactly %d literals?", n, k)
		return edakit.Numeric(q, ans, edakit.Sprintf("Choose which %d variables appear, C(%d, %d) = %d ways, and a polarity for each, 2^%d ways: %s cubes.", k, n, k, binom(n, k), k, ans), nil)
	case 2:
		n := edakit.Scaled(ctx, 4, 5)
		k := edakit.Int(2, 3)
		var cs []cube
		u := tt{n, 0}
		for len(cs) < k {
			c := randCube(n, 1, n-2)
			cs = append(cs, c)
			u = u.or(c.tt(n))
		}
		ans := edakit.IntStr(u.ones())
		q := edakit.Sprintf("Over %d variables, how many minterms lie in the union of the cubes %s?", n, cubeList(cs))
		return edakit.Numeric(q, ans, "Add the cube sizes and subtract the sizes of the pairwise intersections (adding back the triple intersection when there are three cubes). An intersection is empty when two cubes conflict, 0 against 1, in some position. The union has "+ans+" minterms.", nil)
	default:
		n := edakit.Scaled(ctx, 4, 6)
		for {
			a := randCube(n, 2, n-1)
			b := randCube(n, 2, n-1)
			dist := 0
			pos := -1
			for v := 0; v < n; v++ {
				if a[v] != 2 && b[v] != 2 && a[v] != b[v] {
					dist++
					pos = v
				}
			}
			if rand.Intn(3) == 0 {
				ans := edakit.IntStr(dist)
				q := edakit.Sprintf("What is the distance between the cubes %s and %s, the number of positions where one has 0 and the other 1?", a.positional(), b.positional())
				return edakit.Numeric(q, ans, "Count the conflicting positions; a dash conflicts with nothing. The distance is "+ans+".", nil)
			}
			if dist != 1 {
				continue
			}
			c := make(cube, n)
			for v := 0; v < n; v++ {
				switch {
				case v == pos:
					c[v] = 2
				case a[v] == 2:
					c[v] = b[v]
				default:
					c[v] = a[v]
				}
			}
			ans := edakit.IntStr(c.size())
			q := edakit.Sprintf("The cubes %s and %s are at distance 1. How many minterms does their consensus cube contain?", a.positional(), b.positional())
			return edakit.Numeric(q, ans, "The consensus drops the one conflicting variable and intersects the other positions: "+c.positional()+", which has "+ans+" minterms.", map[string]string{"consensus": c.positional()})
		}
	}
}

func sopPos(ctx generator.GeneratorContext) generator.Problem {
	n := nVars(ctx)
	switch kpOf(ctx) {
	case 1:
		f, fs := drawFunc(n, 0.6)
		if rand.Intn(2) == 0 {
			ans := edakit.IntStr(f.ones())
			return edakit.Numeric("Let "+fs+". How many product terms does its canonical sum of minterms have?", ans,
				"One minterm per ON row: "+ans+" terms.", nil)
		}
		ans := edakit.IntStr(n * f.ones())
		return edakit.Numeric("Let "+fs+". How many literals does its canonical sum of minterms contain in total?", ans,
			edakit.Sprintf("There are %d ON rows, each a minterm of %d literals: %s literals.", f.ones(), n, ans), nil)
	case 2:
		f, fs := drawFunc(n, 0.6)
		off := rows(n) - f.ones()
		if rand.Intn(2) == 0 {
			ans := edakit.IntStr(off)
			return edakit.Numeric("Let "+fs+". How many sum terms (maxterms) does its canonical product of maxterms have?", ans,
				edakit.Sprintf("One maxterm per OFF row: %d - %d = %s.", rows(n), f.ones(), ans), nil)
		}
		ans := edakit.IntStr(n * off)
		return edakit.Numeric("Let "+fs+". How many literals does its canonical product of maxterms contain in total?", ans,
			edakit.Sprintf("There are %d OFF rows, each a maxterm of %d literals: %s literals.", off, n, ans), nil)
	default:
		k := edakit.Int(2, 3)
		var clauses []string
		off := tt{n, 0}
		var used []cube
		for len(used) < k {
			c := randCube(n, 2, n-1) // the falsifying cube of the clause
			dup := false
			for _, u := range used {
				if u.equal(c) {
					dup = true
				}
			}
			if dup {
				continue
			}
			used = append(used, c)
			off = off.or(c.tt(n))
			var lits []string
			for v := 0; v < n; v++ {
				if c[v] != 2 {
					lits = append(lits, lit(v, 1-c[v]))
				}
			}
			clauses = append(clauses, "("+strings.Join(lits, " + ")+")")
		}
		on := rows(n) - off.ones()
		ans := edakit.IntStr(on)
		q := edakit.Sprintf("Let $f(%s) = %s$, a product of sums. On how many of the %d rows is $f = 1$?", varList(n), strings.Join(clauses, ""), rows(n))
		return edakit.Numeric(q, ans, edakit.Sprintf("Each sum term is 0 on exactly one cube of rows. Those cubes together cover %d rows, so f is 1 on %d - %d = %s rows.", off.ones(), rows(n), off.ones(), ans), nil)
	}
}

// classifyCube returns 0 (not an implicant), 1 (implicant, not prime) or 2 (prime).
func classifyCube(c cube, f tt) int {
	if !isImplicant(c, f) {
		return 0
	}
	if isPrime(c, f) {
		return 2
	}
	return 1
}

func implicants(ctx generator.GeneratorContext) generator.Problem {
	n := nVars(ctx)
	switch kpOf(ctx) {
	case 1:
		want := rand.Intn(3)
		for {
			f, fs := drawFunc(n, 0.5)
			c := randCube(n, 1, n)
			if classifyCube(c, f) != want {
				continue
			}
			labels := []string{"Not an implicant of $f$", "An implicant of $f$ but not prime", "A prime implicant of $f$"}
			why := []string{
				"Some minterm of the cube is an OFF row of f, so the cube is not contained in the ON-set.",
				"Every minterm of the cube is ON, but dropping one of its literals still gives an implicant, so it is not maximal.",
				"Every minterm is ON, and dropping any single literal brings in an OFF row, so the cube is maximal: prime.",
			}[want]
			var ds []string
			for i, l := range labels {
				if i != want {
					ds = append(ds, l)
				}
			}
			ds = append(ds, "An essential prime of the complement of $f$")
			if want == 0 && isImplicant(c, f.not()) {
				ds = ds[:len(ds)-1]
			}
			q := "Let " + fs + ". Classify the cube $" + c.latex() + "$."
			return edakit.Choice(q, labels[want], why, ds...)
		}
	case 2:
		f, fs := drawFunc(n, 0.5)
		ps := primes(f)
		ans := len(ps)
		var lst []string
		for _, p := range ps {
			lst = append(lst, p.latex())
		}
		q := "Let " + fs + ". How many prime implicants does $f$ have?"
		return edakit.Choice(q, edakit.IntStr(ans), "The primes are $"+strings.Join(lst, "$, $")+"$.", edakit.Distinct(ans, 3, ans+1, ans-1, ans+2)...)
	default:
		for {
			f, fs := drawFunc(n, 0.5)
			ps := primes(f)
			if len(ps) == 0 {
				continue
			}
			p := edakit.Pick(ps)
			// Distractors: a non-prime implicant inside p, and cubes that are not implicants.
			var ds []string
			for _, c := range allCubes(n) {
				if len(ds) >= 6 {
					break
				}
				cls := classifyCube(c, f)
				if cls == 2 || c.lits() == 0 {
					continue
				}
				if (cls == 1 && c.within(p)) || (cls == 0 && c.lits() <= p.lits()+1 && rand.Intn(3) == 0) {
					ds = append(ds, "$"+c.latex()+"$")
				}
			}
			ds = edakit.Shuffle(ds)
			if len(ds) < 3 {
				continue
			}
			q := "Let " + fs + ". Which cube is a prime implicant of $f$?"
			return edakit.Choice(q, "$"+p.latex()+"$", "Only $"+p.latex()+"$ lies inside the ON-set and cannot lose a literal without covering an OFF row. The other cubes either cover an OFF row or can still be enlarged.", ds[:3]...)
		}
	}
}

func essentialPrimes(ctx generator.GeneratorContext) generator.Problem {
	n := nVars(ctx)
	f, fs := drawFunc(n, 0.4)
	ps := primes(f)
	es := essentials(f, ps)
	var pl, el []string
	for _, p := range ps {
		pl = append(pl, p.latex())
	}
	for _, e := range es {
		el = append(el, e.latex())
	}
	elist := "none"
	if len(el) > 0 {
		elist = "$" + strings.Join(el, "$, $") + "$"
	}
	switch kpOf(ctx) {
	case 1:
		ans := edakit.IntStr(len(ps))
		return edakit.Numeric("Let "+fs+". How many prime implicants does $f$ have?", ans,
			"The primes are $"+strings.Join(pl, "$, $")+"$: "+ans+".", map[string]string{"primes": strings.Join(pl, ",")})
	case 2:
		ans := edakit.IntStr(len(es))
		return edakit.Numeric("Let "+fs+". How many of its prime implicants are essential?", ans,
			"The primes are $"+strings.Join(pl, "$, $")+"$. Those that alone cover some ON minterm are essential: "+elist+".", map[string]string{"essentials": elist})
	default:
		cov := tt{n, 0}
		for _, e := range es {
			cov = cov.or(e.tt(n))
		}
		left := f.ones() - cov.and(f).ones()
		ans := edakit.IntStr(left)
		return edakit.Numeric("Let "+fs+". After selecting all essential prime implicants, how many ON minterms remain uncovered?", ans,
			edakit.Sprintf("The essential primes are %s. They cover %d of the %d ON minterms, leaving %s.", elist, cov.and(f).ones(), f.ones(), ans), nil)
	}
}

func drawCoverable(ctx generator.GeneratorContext, exprProb float64) (tt, string, []cube) {
	for {
		f, fs := drawFunc(nVars(ctx), exprProb)
		ps := primes(f)
		if len(ps) <= 16 {
			return f, fs, ps
		}
	}
}

func minCoverGen(ctx generator.GeneratorContext) generator.Problem {
	f, fs, ps := drawCoverable(ctx, 0.4)
	k, lits := minCover(f, ps)
	es := essentials(f, ps)
	switch kpOf(ctx) {
	case 1:
		ans := edakit.IntStr(k)
		return edakit.Numeric("Let "+fs+". What is the smallest number of product terms in any sum-of-products expression for $f$?", ans,
			edakit.Sprintf("A minimum SOP uses primes only. Of the %d primes, %d are essential, and the cheapest set covering every ON minterm has %s cubes.", len(ps), len(es), ans), map[string]string{"min_cubes": ans})
	case 2:
		ans := edakit.IntStr(k - len(es))
		return edakit.Numeric("Let "+fs+". A minimum sum-of-products cover must contain every essential prime. How many more primes, beyond the essential ones, does a minimum cover need?", ans,
			edakit.Sprintf("There are %d essential primes and a minimum cover has %d cubes, so the cyclic core needs %s more.", len(es), k, ans), map[string]string{"min_cubes": edakit.IntStr(k)})
	default:
		ans := edakit.IntStr(lits)
		return edakit.Numeric("Let "+fs+". Among the sum-of-products covers of $f$ with the fewest cubes, what is the smallest total number of literals?", ans,
			edakit.Sprintf("The minimum cover has %d cubes; the cheapest such cover has %s literals.", k, ans), map[string]string{"min_cubes": edakit.IntStr(k)})
	}
}

func dontCares(ctx generator.GeneratorContext) generator.Problem {
	n := nVars(ctx)
	var on, dc tt
	for {
		w := rand.Uint32() & fullMask(n)
		d := rand.Uint32() & fullMask(n) &^ w
		// Keep a few don't cares only.
		for bitsSet(d) > 3 {
			d &= d - 1
		}
		on, dc = tt{n, w}, tt{n, d}
		if on.ones() >= 2 && dc.ones() >= 1 && on.or(dc).ones() < rows(n) {
			break
		}
	}
	care := on.or(dc)
	q0 := edakit.Sprintf("An incompletely specified function of %s has ON-set $\\sum m(%s)$ and don't-care set $d(%s)$; every other row is OFF.",
		varList(n), mintermList(on.minterms()), mintermList(dc.minterms()))
	switch kpOf(ctx) {
	case 1:
		ans := edakit.IntStr(pow2(dc.ones()))
		return edakit.Numeric(q0+" How many completely specified functions agree with it on every care row?", ans,
			edakit.Sprintf("Each of the %d don't-care rows may be 0 or 1: 2^%d = %s.", dc.ones(), dc.ones(), ans), nil)
	case 2:
		ps := primes(care)
		k, _ := minCover(on, ps)
		k0, _ := minCover(on, primes(on))
		ans := edakit.IntStr(k)
		return edakit.Numeric(q0+" What is the smallest number of cubes in a sum-of-products cover that is 1 on the ON-set and 0 on the OFF-set?", ans,
			edakit.Sprintf("Grow primes over ON plus don't-care rows, then cover only the ON rows: %s cubes. Without the don't cares the minimum would be %d.", ans, k0), map[string]string{"without_dc": edakit.IntStr(k0)})
	default:
		ps := primes(care)
		useful := 0
		for _, p := range ps {
			if p.tt(n).w&on.w != 0 {
				useful++
			}
		}
		ans := edakit.IntStr(useful)
		return edakit.Numeric(q0+" Compute the prime implicants of ON plus don't care. How many of them cover at least one ON minterm?", ans,
			edakit.Sprintf("ON plus don't care has %d primes; %s of them touch the ON-set, and only those can appear in a minimum cover.", len(ps), ans), nil)
	}
}

func bitsSet(w uint32) int {
	k := 0
	for w != 0 {
		w &= w - 1
		k++
	}
	return k
}

func tautology(ctx generator.GeneratorContext) generator.Problem {
	n := nVars(ctx)
	var F []cube
	k := edakit.Int(3, 5)
	for len(F) < k {
		F = append(F, randCube(n, 1, 2))
	}
	cov := sop(F).tt(n)
	q0 := edakit.Sprintf("Over the variables %s, in that order, a cover is $F = \\{%s\\}$ in positional notation.", varList(n), cubeList(F))
	switch kpOf(ctx) {
	case 1:
		// Sometimes complete the cover to a tautology.
		if rand.Intn(2) == 0 {
			for cov.ones() < rows(n) {
				missing := cov.not().minterms()[0]
				c := make(cube, n)
				for v := 0; v < n; v++ {
					c[v] = bitOf(missing, v, n)
				}
				// Enlarge the cube while it stays a single block; any cube works for coverage.
				c[rand.Intn(n)] = 2
				F = append(F, c)
				cov = sop(F).tt(n)
			}
			q0 = edakit.Sprintf("Over the variables %s, in that order, a cover is $F = \\{%s\\}$ in positional notation.", varList(n), cubeList(F))
		}
		miss := rows(n) - cov.ones()
		ans := edakit.IntStr(miss)
		return edakit.Numeric(q0+edakit.Sprintf(" How many of the %d minterms does $F$ miss? (Zero means $F$ is a tautology.)", rows(n)), ans,
			edakit.Sprintf("The cubes together cover %d of the %d minterms, so %s are missed.", cov.ones(), rows(n), ans), nil)
	case 2:
		v := rand.Intn(n)
		val := rand.Intn(2)
		cnt := 0
		for _, c := range F {
			if c[v] == 2 || c[v] == val {
				cnt++
			}
		}
		ans := edakit.IntStr(cnt)
		return edakit.Numeric(q0+edakit.Sprintf(" How many cubes does the cofactor cover $F_{%s}$ contain? Keep duplicate cubes; do not merge anything.", lit(v, val)), ans,
			edakit.Sprintf("Cubes with %s = %d conflict with the literal and drop out. The others survive with position %s set to dash: %s cubes.", varNames[v], 1-val, varNames[v], ans), nil)
	default:
		var c cube
		for {
			c = randCube(n, 1, 2)
			inF := false
			for _, d := range F {
				if c.within(d) {
					inF = true
				}
			}
			if !inF {
				break
			}
		}
		miss := c.tt(n).ones() - c.tt(n).and(cov).ones()
		ans := edakit.IntStr(miss)
		return edakit.Numeric(q0+edakit.Sprintf(" Is the cube %s contained in $F$? Answer with the number of minterms of %s that $F$ misses (zero means contained).", c.positional(), c.positional()), ans,
			edakit.Sprintf("The cube is contained exactly when the cofactor of F with respect to it is a tautology. Of its %d minterms, %s are not covered.", c.size(), ans), nil)
	}
}
