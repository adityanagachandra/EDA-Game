package sim

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func codeCoverage(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		nf := edakit.Scaled(ctx, 2, 4)
		var parts []string
		tot, hit := 0, 0
		for i := 0; i < nf; i++ {
			s := edakit.Int(4, 40) * 5
			h := edakit.Int(s/3, s)
			tot += s
			hit += h
			parts = append(parts, edakit.Sprintf("module %c has %d statements of which %d executed", 'A'+i, s, h))
		}
		pct := 100 * float64(hit) / float64(tot)
		ans := edakit.Num(pct, 1)
		q := "A coverage report says: " + strings.Join(parts, "; ") + ". What is the overall statement coverage of the design, in percent? Round to 1 decimal place."
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Pool the counts: %d executed of %d statements is %s%%. Averaging the per-module percentages would weight a small module like a large one.", hit, tot, ans),
			map[string]string{"hit": edakit.IntStr(hit), "total": edakit.IntStr(tot)})
	case 2:
		m := edakit.Int(1, 3)
		k := edakit.Int(2, 5)
		ifOut := m + 1
		caseOut := k + 1
		tot := ifOut + caseOut
		hitIf := edakit.Int(1, ifOut)
		hitCase := edakit.Int(1, caseOut)
		hit := hitIf + hitCase
		pct := 100 * float64(hit) / float64(tot)
		ans := edakit.Num(pct, 1)
		q := edakit.Sprintf("A block has an if / else-if chain with %d condition(s) and no final else, and a case statement with %d items plus a default. Branch coverage counts every outcome, including the implicit fall-through of the if chain. The tests reached %d of the if-chain outcomes and %d of the case outcomes. What is the branch coverage, in percent? Round to 1 decimal place.", m, k, hitIf, hitCase)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The if chain has %d + 1 = %d outcomes (the missing else is still a branch), and the case has %d + 1 = %d. Coverage is %d/%d = %s%%.", m, ifOut, k, caseOut, hit, tot, ans),
			map[string]string{"outcomes": edakit.IntStr(tot)})
	default:
		w := edakit.Int(3, 5)
		n := edakit.Scaled(ctx, 3, 6)
		vals := make([]int, n)
		for i := range vals {
			vals[i] = edakit.Int(0, 1<<w-1)
		}
		rise := make([]bool, w)
		fall := make([]bool, w)
		for i := 1; i < n; i++ {
			for b := 0; b < w; b++ {
				p := (vals[i-1] >> b) & 1
				c := (vals[i] >> b) & 1
				if p == 0 && c == 1 {
					rise[b] = true
				}
				if p == 1 && c == 0 {
					fall[b] = true
				}
			}
		}
		cov := 0
		for b := 0; b < w; b++ {
			if rise[b] {
				cov++
			}
			if fall[b] {
				cov++
			}
		}
		strs := make([]string, n)
		for i, v := range vals {
			s := ""
			for b := w - 1; b >= 0; b-- {
				s += edakit.IntStr((v >> b) & 1)
			}
			strs[i] = s
		}
		pct := 100 * float64(cov) / float64(2*w)
		ans := edakit.Num(pct, 1)
		q := edakit.Sprintf("A %d-bit signal takes the values %s (binary) on successive cycles. Toggle coverage counts each bit's 0->1 and 1->0 transitions as separate items. What is its toggle coverage, in percent? Round to 1 decimal place.", w, strings.Join(strs, ", "))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Comparing consecutive values bit by bit, %d of the %d transition items were observed, which is %s%%.", cov, 2*w, ans),
			map[string]string{"items": edakit.IntStr(cov)})
	}
}

// boolean expression over named conditions
type bexpr struct {
	text string
	n    int
	f    func(v []int) int
}

func drawExpr() bexpr {
	opts := []bexpr{
		{"a && b", 2, func(v []int) int { return v[0] & v[1] }},
		{"a || b", 2, func(v []int) int { return v[0] | v[1] }},
		{"a && b && c", 3, func(v []int) int { return v[0] & v[1] & v[2] }},
		{"a || b || c", 3, func(v []int) int { return v[0] | v[1] | v[2] }},
		{"(a && b) || c", 3, func(v []int) int { return (v[0] & v[1]) | v[2] }},
		{"a && (b || c)", 3, func(v []int) int { return v[0] & (v[1] | v[2]) }},
		{"(a || b) && (c || d)", 4, func(v []int) int { return (v[0] | v[1]) & (v[2] | v[3]) }},
		{"(a && b) || (c && d)", 4, func(v []int) int { return (v[0] & v[1]) | (v[2] & v[3]) }},
	}
	return edakit.Pick(opts)
}

func vecText(v []int) string {
	s := ""
	for _, x := range v {
		s += edakit.IntStr(x)
	}
	return s
}

func conditionCoverage(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	e := drawExpr()
	nv := edakit.Scaled(ctx, 2, 5)
	if nv > 1<<e.n {
		nv = 1 << e.n
	}
	if kp == 3 && nv > e.n+1 {
		nv = e.n + 1
	}
	seen := map[string]bool{}
	var vecs [][]int
	base := make([]int, e.n)
	for i := range base {
		base[i] = edakit.Int(0, 1)
	}
	for len(vecs) < nv {
		v := make([]int, e.n)
		for i := range v {
			v[i] = edakit.Int(0, 1)
		}
		if kp == 3 {
			// MC/DC sets are built from single-condition flips of a base vector.
			copy(v, base)
			if len(vecs) > 0 {
				v[edakit.Int(0, e.n-1)] ^= 1
			}
		}
		if !seen[vecText(v)] {
			seen[vecText(v)] = true
			vecs = append(vecs, v)
		}
	}
	names := []string{"a", "b", "c", "d"}[:e.n]
	var vt []string
	for _, v := range vecs {
		vt = append(vt, vecText(v))
	}
	setup := edakit.Sprintf("The condition of an if statement is %s. The tests apply the vectors %s, each listing %s in that order. Treat every condition as fully evaluated (no short-circuit masking).", e.text, strings.Join(vt, ", "), strings.Join(names, ""))
	pairs := 0
	for i := 0; i < e.n; i++ {
		has0, has1 := false, false
		for _, v := range vecs {
			if v[i] == 0 {
				has0 = true
			} else {
				has1 = true
			}
		}
		if has0 {
			pairs++
		}
		if has1 {
			pairs++
		}
	}
	facts := map[string]string{"pairs": edakit.IntStr(pairs)}
	switch kp {
	case 1:
		return edakit.Numeric(setup+edakit.Sprintf(" Condition coverage needs each of the %d conditions to be observed both true and false. How many of the %d condition-value pairs are observed?", e.n, 2*e.n),
			edakit.IntStr(pairs),
			edakit.Sprintf("Scanning each column of the vectors for a 0 and a 1 finds %d of the %d pairs.", pairs, 2*e.n), facts)
	case 2:
		pct := 100 * float64(pairs) / float64(2*e.n)
		ans := edakit.Num(pct, 1)
		return edakit.Numeric(setup+" What is the condition coverage, in percent? Round to 1 decimal place.", ans,
			edakit.Sprintf("%d of %d condition-value pairs are observed: %s%%.", pairs, 2*e.n, ans), facts)
	default:
		// MC/DC: condition i is demonstrated if two vectors differ only in i and give different outcomes.
		shown := 0
		for i := 0; i < e.n; i++ {
			ok := false
			for x := 0; x < len(vecs) && !ok; x++ {
				for y := x + 1; y < len(vecs); y++ {
					diff := 0
					for j := 0; j < e.n; j++ {
						if vecs[x][j] != vecs[y][j] {
							diff++
						}
					}
					if diff == 1 && vecs[x][i] != vecs[y][i] && e.f(vecs[x]) != e.f(vecs[y]) {
						ok = true
						break
					}
				}
			}
			if ok {
				shown++
			}
		}
		if shown == 0 && edakit.Int(0, 3) > 0 {
			// Redraw most instances with nothing demonstrated, so the answer is rarely a trivial 0.
			return conditionCoverage(ctx)
		}
		return edakit.Numeric(setup+" Under MC/DC with unique-cause pairs, a condition is demonstrated when two of the vectors differ only in that condition and produce different outcomes. How many conditions are demonstrated?",
			edakit.IntStr(shown),
			edakit.Sprintf("Checking every pair of vectors that differ in exactly one position and flip the outcome demonstrates %d of the %d conditions. Full MC/DC on %d conditions needs at least %d vectors.", shown, e.n, e.n, e.n+1),
			map[string]string{"demonstrated": edakit.IntStr(shown)})
	}
}

func functionalCoverage(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		type decl struct {
			text string
			n    int
		}
		lo := 0
		var ds []decl
		k := edakit.Scaled(ctx, 2, 4)
		for i := 0; i < k; i++ {
			w := edakit.Int(2, 8)
			hi := lo + w - 1
			switch edakit.Int(0, 2) {
			case 0:
				ds = append(ds, decl{edakit.Sprintf("bins r%d = {[%d:%d]};", i, lo, hi), 1})
			case 1:
				ds = append(ds, decl{edakit.Sprintf("bins r%d[] = {[%d:%d]};", i, lo, hi), w})
			default:
				nb := edakit.Int(2, w)
				ds = append(ds, decl{edakit.Sprintf("bins r%d[%d] = {[%d:%d]};", i, nb, lo, hi), nb})
			}
			lo = hi + 1
		}
		if edakit.Int(0, 3) == 0 {
			wbits := edakit.Int(3, 10)
			nb := 1 << wbits
			if nb > 64 {
				nb = 64
			}
			q := edakit.Sprintf("A coverpoint samples an unsigned %d-bit signal and declares no bins, so the tool creates automatic bins with the default auto_bin_max of 64. How many bins does the coverpoint have?", wbits)
			return edakit.Numeric(q, edakit.IntStr(nb),
				edakit.Sprintf("There are %d possible values; automatic binning creates one bin per value up to auto_bin_max, so %d bins.", 1<<wbits, nb),
				map[string]string{"bins": edakit.IntStr(nb)})
		}
		var txt []string
		tot := 0
		for _, d := range ds {
			txt = append(txt, d.text)
			tot += d.n
		}
		q := "A coverpoint declares " + strings.Join(txt, " ") + " How many bins does it have?"
		return edakit.Numeric(q, edakit.IntStr(tot),
			"A scalar bin covers its whole range as one bin, bins name[] makes one bin per value, and bins name[N] splits the values into N bins. The total is "+edakit.IntStr(tot)+".",
			map[string]string{"bins": edakit.IntStr(tot)})
	case 2:
		nb := edakit.Scaled(ctx, 5, 10)
		atLeast := edakit.Int(1, 3)
		var hits []string
		cov := 0
		for i := 0; i < nb; i++ {
			h := edakit.Int(0, 6)
			if h >= atLeast {
				cov++
			}
			hits = append(hits, edakit.IntStr(h))
		}
		pct := 100 * float64(cov) / float64(nb)
		ans := edakit.Num(pct, 1)
		q := edakit.Sprintf("A coverpoint has %d bins with option.at_least = %d. After the regression the bins have hit counts %s. What is the coverpoint's coverage, in percent? Round to 1 decimal place.", nb, atLeast, strings.Join(hits, ", "))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("A bin is covered only when its count reaches %d. %d of %d bins qualify: %s%%.", atLeast, cov, nb, ans),
			map[string]string{"covered": edakit.IntStr(cov)})
	default:
		n := edakit.Int(2, 4)
		var parts []string
		num, den := 0.0, 0.0
		for i := 0; i < n; i++ {
			c := float64(edakit.Int(4, 20) * 5)
			w := edakit.Int(1, 4)
			num += c * float64(w)
			den += float64(w)
			parts = append(parts, edakit.Sprintf("cp%d at %s%% with weight %d", i+1, edakit.Num(c, 0), w))
		}
		ans := edakit.Num(num/den, 1)
		q := "A covergroup contains coverpoints " + strings.Join(parts, ", ") + ". What is the covergroup's coverage, in percent? Round to 1 decimal place."
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Covergroup coverage is the weighted average of its items: %s / %s = %s%%.", edakit.Num(num, 0), edakit.Num(den, 0), ans),
			map[string]string{"weight_sum": edakit.Num(den, 0)})
	}
}

func crossCoverage(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	na, nb := edakit.Int(2, 6), edakit.Int(2, 6)
	switch kp {
	case 1:
		if edakit.Int(0, 1) == 0 {
			nc := edakit.Int(2, 5)
			q := edakit.Sprintf("Coverpoints cp_op, cp_size and cp_mode have %d, %d and %d bins. How many bins does cross cp_op, cp_size, cp_mode have?", na, nb, nc)
			return edakit.Numeric(q, edakit.IntStr(na*nb*nc),
				edakit.Sprintf("A cross has one bin per combination of its coverpoints' bins: %d x %d x %d = %d.", na, nb, nc, na*nb*nc),
				map[string]string{"bins": edakit.IntStr(na * nb * nc)})
		}
		q := edakit.Sprintf("Coverpoint cp_a has %d bins and cp_b has %d bins. How many bins does cross cp_a, cp_b have?", na, nb)
		return edakit.Numeric(q, edakit.IntStr(na*nb),
			edakit.Sprintf("The cross is the Cartesian product of the bins: %d x %d = %d.", na, nb, na*nb),
			map[string]string{"bins": edakit.IntStr(na * nb)})
	default:
		// cp_a bins a0..a{na-1}, cp_b bins b0..b{nb-1}.
		ia := edakit.Int(0, na-1)
		jb := edakit.Int(0, nb-1)
		ia2, jb2 := edakit.Int(0, na-1), edakit.Int(0, nb-1)
		excluded := map[[2]int]bool{}
		for j := 0; j < nb; j++ {
			excluded[[2]int{ia, j}] = true
		}
		for i := 0; i < na; i++ {
			excluded[[2]int{i, jb}] = true
		}
		illegal := [2]int{ia2, jb2}
		excluded[illegal] = true
		illegalNote := "The illegal combination removes one more."
		if illegal[0] == ia || illegal[1] == jb {
			illegalNote = "The illegal combination was already ignored."
		}
		valid := na*nb - len(excluded)
		if valid <= 0 {
			valid = 0
		}
		decl := edakit.Sprintf("cp_a has bins a0 to a%d and cp_b has bins b0 to b%d. Their cross declares ignore_bins i1 = binsof(cp_a.a%d); ignore_bins i2 = binsof(cp_b.b%d); illegal_bins x = binsof(cp_a.a%d) && binsof(cp_b.b%d);", na-1, nb-1, ia, jb, ia2, jb2)
		if kp == 2 || valid < 2 {
			return edakit.Numeric(decl+" How many cross bins count toward coverage?", edakit.IntStr(valid),
				edakit.Sprintf("Of the %d combinations, i1 removes %d and i2 removes %d, overlapping in one, so %d go. %s That leaves %d bins.", na*nb, nb, na, na+nb-1, illegalNote, valid),
				map[string]string{"valid": edakit.IntStr(valid)})
		}
		var validList [][2]int
		for i := 0; i < na; i++ {
			for j := 0; j < nb; j++ {
				if !excluded[[2]int{i, j}] {
					validList = append(validList, [2]int{i, j})
				}
			}
		}
		h := edakit.Int(1, len(validList))
		hitSet := edakit.Shuffle(validList)[:h]
		// Add some samples that land in ignored bins and some repeats.
		samples := append([][2]int(nil), hitSet...)
		samples = append(samples, [2]int{ia, edakit.Int(0, nb-1)}, hitSet[0])
		samples = edakit.Shuffle(samples)
		var st []string
		for _, s := range samples {
			st = append(st, edakit.Sprintf("(a%d, b%d)", s[0], s[1]))
		}
		pct := 100 * float64(h) / float64(valid)
		ans := edakit.Num(pct, 1)
		return edakit.Numeric(decl+" The sampled (cp_a, cp_b) pairs are "+strings.Join(st, ", ")+". What is the cross coverage, in percent? Round to 1 decimal place.", ans,
			edakit.Sprintf("%d bins remain after ignore and illegal bins. The samples hit %d distinct remaining bins (repeats and ignored samples add nothing), so coverage is %s%%.", valid, h, ans),
			map[string]string{"valid": edakit.IntStr(valid), "hit": edakit.IntStr(h)})
	}
}

func harmonic(n int) float64 {
	s := 0.0
	for i := 1; i <= n; i++ {
		s += 1 / float64(i)
	}
	return s
}

func couponCollector(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		m := edakit.Int(1, 40)
		den := edakit.Pick([]int{200, 500, 1000})
		e := float64(den) / float64(m)
		ans := edakit.Num(e, 2)
		q := edakit.Sprintf("A coverage bin is hit by %d out of every %d random tests on average, independently. What is the expected number of tests until it is first hit? Round to 2 decimal places.", m, den)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The wait is geometric with p = %d/%d, so the mean is 1/p = %s tests.", m, den, ans),
			map[string]string{"p": edakit.Num(float64(m)/float64(den), 4)})
	case 2:
		n := edakit.Scaled(ctx, 4, 30)
		e := float64(n) * harmonic(n)
		ans := edakit.Num(e, 2)
		q := edakit.Sprintf("Each random test hits exactly one of %d equally likely coverage bins. What is the expected number of tests to hit all %d bins? Round to 2 decimal places.", n, n)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("This is the coupon collector: n H_n = %d x %s = %s.", n, edakit.Num(harmonic(n), 4), ans),
			map[string]string{"H_n": edakit.Num(harmonic(n), 4)})
	default:
		n := edakit.Scaled(ctx, 8, 40)
		k := edakit.Int(n/2, n-1)
		e := float64(n) * harmonic(n-k)
		ans := edakit.Num(e, 2)
		q := edakit.Sprintf("Each random test hits exactly one of %d equally likely bins, and %d of them are already covered. What is the expected number of further tests to cover the rest? Round to 2 decimal places.", n, k)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("With j bins missing, the next new bin takes n/j tests on average. Summing j = 1..%d gives n H_%d = %d x %s = %s; the last few bins dominate.", n-k, n-k, n, edakit.Num(harmonic(n-k), 4), ans),
			map[string]string{"H": edakit.Num(harmonic(n-k), 4)})
	}
}

// greedyRanking asks what the greedy set-cover ranking of tests picks.
func greedyRanking(ctx generator.GeneratorContext) generator.Problem {
	for {
		nt := edakit.Int(4, 5)
		nbins := edakit.Int(8, 12)
		sets := make([]map[int]bool, nt)
		for i := range sets {
			sets[i] = map[int]bool{}
			sz := edakit.Int(2, nbins/2+1)
			for len(sets[i]) < sz {
				sets[i][edakit.Int(1, nbins)] = true
			}
		}
		covered := map[int]bool{}
		used := make([]bool, nt)
		var order []int
		unique := true
		for {
			best, bestGain, ties := -1, 0, 0
			for i, s := range sets {
				if used[i] {
					continue
				}
				g := 0
				for b := range s {
					if !covered[b] {
						g++
					}
				}
				if g > bestGain {
					best, bestGain, ties = i, g, 1
				} else if g == bestGain && g > 0 {
					ties++
				}
			}
			if best < 0 {
				break
			}
			if ties > 1 {
				unique = false
				break
			}
			used[best] = true
			order = append(order, best)
			for b := range sets[best] {
				covered[b] = true
			}
		}
		if !unique || len(order) < 2 {
			continue
		}
		var desc []string
		for i, s := range sets {
			var bs []int
			for b := 1; b <= nbins; b++ {
				if s[b] {
					bs = append(bs, b)
				}
			}
			desc = append(desc, edakit.Sprintf("T%d covers bins {%s}", i+1, intsText(bs)))
		}
		setup := "Test ranking keeps the fewest tests that preserve coverage. Greedy ranking repeatedly picks the test that adds the most not-yet-covered bins. " + strings.Join(desc, "; ") + "."
		if edakit.Int(0, 1) == 0 {
			ans := edakit.Sprintf("T%d", order[1]+1)
			var d []string
			for i := 0; i < nt; i++ {
				if i != order[1] {
					d = append(d, edakit.Sprintf("T%d", i+1))
				}
			}
			return edakit.Choice(setup+" Which test does greedy ranking pick second?", ans,
				edakit.Sprintf("Greedy first picks T%d, the largest set. Counting only bins it does not already cover, T%d adds the most next.", order[0]+1, order[1]+1), d...)
		}
		ans := edakit.IntStr(len(order))
		return edakit.Choice(setup+" How many tests does greedy ranking keep before no test adds a new bin?", ans,
			edakit.Sprintf("Greedy picks %d tests in turn until every bin covered by any test is covered.", len(order)),
			edakit.Distinct(len(order), 3, nt, len(order)-1, len(order)+1)...)
	}
}

var closureItems = []edakit.Item{
	{KP: 1, Q: "A coverage report shows 97% functional coverage after 2,000 seeds, unchanged over the last 1,500. What does the plateau suggest?",
		A: "The remaining bins are unlikely under the current constraints and need biased or directed stimulus", D: []string{"The remaining bins are already hit but not recorded", "More seeds of the same test will close it soon", "The design is fully verified"},
		Why: "A flat curve means new seeds are not reaching the holes; the stimulus distribution must change."},
	{KP: 1, Q: "A bin cannot be hit because the configuration that reaches it is not supported in this product. What is the right action?",
		A: "Exclude the bin with a documented, reviewed waiver", D: []string{"Write a test that forces the signal to that value", "Delete the coverpoint entirely", "Leave it as a hole and ignore the report"},
		Why: "Unreachable bins are excluded with justification so the coverage number means something."},
	{KP: 1, Q: "Why are coverage databases from many seeds merged?",
		A: "Coverage is a property of the whole regression, so bins hit by any test count", D: []string{"Each seed must reach 100% on its own", "Merging removes duplicate bugs", "Merging speeds up the simulator"},
		Why: "Closure is measured over the union of all runs."},
	{KP: 1, Q: "Code coverage is 100% but functional coverage is 70%. What does this mean?",
		A: "Every line ran, but many specified scenarios and combinations were never exercised", D: []string{"The functional coverage model is wrong by definition", "The design is verified", "The code coverage tool is broken"},
		Why: "Code coverage measures structure; functional coverage measures intent. They answer different questions."},
	{KP: 1, Q: "What is a coverage hole?",
		A: "A bin or scenario in the coverage model that no test has hit yet", D: []string{"A bug found by an assertion", "A line of code that never runs", "A test that always fails"},
		Why: "Holes are the gap between the coverage model and what the regression achieved."},
	{KP: 1, Q: "Functional coverage reaches 100% quickly on a new block. What should the team suspect first?",
		A: "The coverage model is too coarse to capture the interesting scenarios", D: []string{"The block has no bugs", "The simulator is skipping tests", "Code coverage is wrong"},
		Why: "Coverage is only as good as the model; fast closure often means missing crosses or corner bins."},
	{KP: 3, Q: "In coverage-driven verification, what drives the next round of work?",
		A: "Analysis of coverage holes, which leads to new constraints, sequences or directed tests", D: []string{"The number of tests written so far", "The simulator speed", "The order of bins in the report"},
		Why: "The loop is: run, measure coverage, analyze holes, adjust stimulus, repeat."},
	{KP: 3, Q: "A hole needs a FIFO full while a flush request arrives. Random traffic almost never fills the FIFO. What is the most efficient fix?",
		A: "Add a constraint mode or sequence that throttles the reader so the FIFO fills, then randomize around it", D: []string{"Run a hundred times more seeds", "Lower option.at_least to 0", "Remove the coverpoint"},
		Why: "Steering stimulus toward the precondition makes the corner likely instead of waiting for luck."},
	{KP: 3, Q: "Why rank tests by coverage contribution before a release regression?",
		A: "A small ranked subset reaches the same coverage, so each change can be checked quickly", D: []string{"Ranking finds more bugs than the full set", "Unranked tests cannot be merged", "Ranking makes tests deterministic"},
		Why: "Ranking is greedy set cover over bins; it shrinks regression time without losing reach."},
	{KP: 3, Q: "Greedy test ranking is not guaranteed to find the smallest covering subset. Why is it still used?",
		A: "Minimum set cover is NP-hard, and greedy is fast with a logarithmic approximation guarantee", D: []string{"Greedy always finds the optimum", "Exact ranking needs a simulator", "Coverage cannot be computed exactly"},
		Why: "Set cover is NP-hard; greedy achieves an H_n approximation, which is good enough in practice."},
	{KP: 3, Q: "A cross of opcode and address region has 400 bins, and 60 are illegal combinations. What should be done with them?",
		A: "Mark them as illegal_bins so a hit is flagged as an error and they leave the denominator", D: []string{"Leave them as holes", "Write directed tests to hit them", "Double their weight"},
		Why: "illegal_bins remove impossible combinations from coverage and turn an unexpected hit into a reported error."},
	{KP: 3, Q: "When is functional coverage closure a sign-off criterion that can be trusted?",
		A: "When the coverage model was reviewed against the specification and all holes are hit or waived", D: []string{"When the percentage reaches 90%", "When every test passes once", "When code coverage is 100%"},
		Why: "A reviewed model plus closed or waived holes is what makes the number meaningful."},
}
