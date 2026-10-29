package place

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// rulePluginCalc evaluates the cost a Tier 2 "abut regular cells" plug-in adds.
func rulePluginCalc(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Int(5, 7)
	type rc struct {
		name string
		x, w int
		g    bool
	}
	cells := make([]rc, n)
	x := edakit.Int(0, 4)
	for i := range cells {
		cells[i] = rc{name: edakit.Sprintf("u%d", i+1), x: x, w: edakit.Int(1, 4), g: edakit.Int(0, 1) == 0}
		x += cells[i].w + edakit.Int(0, 3)
	}
	cnt := 0
	for _, c := range cells {
		if c.g {
			cnt++
		}
	}
	if cnt < 3 {
		cells[0].g, cells[n-1].g, cells[n/2].g = true, true, true
	}
	lam := edakit.Pick([]float64{0.5, 1, 2, 5})
	gapSum, allGaps, whiteOnly := 0, 0, 0
	prev := -1
	for i, c := range cells {
		if i > 0 {
			allGaps += c.x - cells[i-1].x - cells[i-1].w
		}
		if !c.g {
			continue
		}
		if prev >= 0 {
			p := cells[prev]
			gapSum += c.x - p.x - p.w
			for k := prev + 1; k <= i; k++ {
				whiteOnly += cells[k].x - cells[k-1].x - cells[k-1].w
			}
		}
		prev = i
	}
	var desc []string
	for _, c := range cells {
		tag := ""
		if c.g {
			tag = ", group R"
		}
		desc = append(desc, edakit.Sprintf("%s at x = %d, width %d%s", c.name, c.x, c.w, tag))
	}
	q := edakit.Sprintf("A team adds an \"abut regular cells\" rule to a placer as a Tier 2 plug-in at the detailed-placement cost seam. For each row it adds $\\lambda$ times the total distance between consecutive members of the regular group R (from one member's right edge to the next member's left edge, whatever lies between), with $\\lambda$ = %s. One row holds, left to right (site units): %s. What cost does the plug-in add for this row?", edakit.Num(lam, 1), strings.Join(desc, "; "))
	ans, d := numChoices(lam*float64(gapSum), 1, "",
		float64(gapSum), lam*float64(allGaps), lam*float64(whiteOnly), lam*float64(gapSum+1))
	return edakit.Choice(q, ans, edakit.Sprintf("Walking the R members in x order, the distances between consecutive members sum to %d sites, so the plug-in adds %s x %d = %s. The core still owns legality; the plug-in only scores the arrangement.", gapSum, edakit.Num(lam, 1), gapSum, ans), d...)
}

// objectiveWeights compares placements under a weighted objective.
func objectiveWeights(ctx generator.GeneratorContext) generator.Problem {
	type sol struct{ wl, of, tns float64 }
	p, q := sol{900000, 10, 100}, sol{950000, 4, 50}
	for try := 0; try < generator.MaxDrawAttempts; try++ {
		cp := sol{float64(1000 * edakit.Int(800, 1200)), float64(edakit.Int(1, 15)), float64(edakit.Int(0, 400))}
		cq := sol{cp.wl + float64(1000*edakit.Int(-60, 60)), float64(edakit.Int(1, 15)), float64(edakit.Int(0, 400))}
		if cp.of != cq.of && cp.wl != cq.wl {
			p, q = cp, cq
			break
		}
	}
	b := float64(1000 * edakit.Pick([]int{1, 2, 5, 10}))
	g := float64(edakit.Pick([]int{10, 50, 100, 200}))
	cost := func(s sol) float64 { return s.wl + b*s.of + g*s.tns }
	head := edakit.Sprintf("A placer exposes its objective weights as Tier 1 configuration: cost = HPWL + $\\beta$ x overflow + $\\gamma$ x |TNS|, with HPWL in um, overflow in percent and TNS in ps. Run P gives HPWL %s, overflow %s%%, TNS -%s ps; run Q gives HPWL %s, overflow %s%%, TNS -%s ps. ",
		edakit.Num(p.wl, 0), edakit.Num(p.of, 0), edakit.Num(p.tns, 0), edakit.Num(q.wl, 0), edakit.Num(q.of, 0), edakit.Num(q.tns, 0))
	switch ctx.KPOrdinal {
	case 1:
		ans := edakit.Num(cost(p), 0)
		return edakit.Numeric(head+edakit.Sprintf("With $\\beta$ = %s and $\\gamma$ = %s, what is the cost of run P?", edakit.Num(b, 0), edakit.Num(g, 0)), ans,
			edakit.Sprintf("%s + %s x %s + %s x %s = %s.", edakit.Num(p.wl, 0), edakit.Num(b, 0), edakit.Num(p.of, 0), edakit.Num(g, 0), edakit.Num(p.tns, 0), ans), map[string]string{"cost_p": ans})
	case 3:
		cc := func(s sol) float64 { return s.wl + g*s.tns }
		bad := func(p, q sol) bool { return p.of == q.of || (p.of > q.of) != (cc(q) > cc(p)) || cc(q) == cc(p) }
		for try := 0; bad(p, q) && try < generator.MaxDrawAttempts; try++ {
			p = sol{float64(1000 * edakit.Int(800, 1200)), float64(edakit.Int(1, 15)), float64(edakit.Int(0, 400))}
			q = sol{p.wl + float64(1000*edakit.Int(-60, 60)), float64(edakit.Int(1, 15)), float64(edakit.Int(0, 400))}
		}
		if bad(p, q) {
			p, q = sol{900000, 10, 100}, sol{950000, 4, 50}
		}
		head = edakit.Sprintf("A placer exposes its objective weights as Tier 1 configuration: cost = HPWL + $\\beta$ x overflow + $\\gamma$ x |TNS|, with HPWL in um, overflow in percent and TNS in ps. Run P gives HPWL %s, overflow %s%%, TNS -%s ps; run Q gives HPWL %s, overflow %s%%, TNS -%s ps. ",
			edakit.Num(p.wl, 0), edakit.Num(p.of, 0), edakit.Num(p.tns, 0), edakit.Num(q.wl, 0), edakit.Num(q.of, 0), edakit.Num(q.tns, 0))
		be := (cc(q) - cc(p)) / (p.of - q.of)
		ans := edakit.Num(be, 2)
		return edakit.Numeric(head+edakit.Sprintf("With $\\gamma$ = %s, at what value of $\\beta$ do the two runs cost the same? Round to 2 decimal places.", edakit.Num(g, 0)), ans,
			edakit.Sprintf("Without the overflow term, P costs %s and Q costs %s. Setting the full costs equal gives $\\beta$ = (%s - %s) / (%s - %s) = %s. Above it the run with less overflow wins. An agent can sweep this crossover without recompiling anything.", edakit.Num(cc(p), 0), edakit.Num(cc(q), 0), edakit.Num(cc(q), 0), edakit.Num(cc(p), 0), edakit.Num(p.of, 0), edakit.Num(q.of, 0), ans),
			map[string]string{"beta": ans})
	}
	ans := edakit.Num(cost(q)-cost(p), 0)
	return edakit.Numeric(head+edakit.Sprintf("With $\\beta$ = %s and $\\gamma$ = %s, what is cost(Q) - cost(P)?", edakit.Num(b, 0), edakit.Num(g, 0)), ans,
		edakit.Sprintf("cost(P) = %s and cost(Q) = %s, so the difference is %s; the lower-cost run is preferred.", edakit.Num(cost(p), 0), edakit.Num(cost(q), 0), ans), map[string]string{"cost_p": edakit.Num(cost(p), 0), "cost_q": edakit.Num(cost(q), 0)})
}

// qorDifferential summarises a placer change over a regression suite.
func qorDifferential(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 7)
	base, nw, of := make([]float64, n), make([]float64, n), make([]float64, n)
	var rows []string
	for i := range base {
		base[i] = float64(1000 * edakit.Int(100, 900))
		nw[i] = base[i] * (1 + float64(edakit.Int(-40, 30))/1000)
		nw[i] = math.Round(nw[i])
		of[i] = float64(edakit.Int(4, 14)) / 100
		rows = append(rows, edakit.Sprintf("d%d: baseline HPWL %s, new HPWL %s, new overflow %s", i+1, edakit.Num(base[i], 0), edakit.Num(nw[i], 0), edakit.Num(of[i], 2)))
	}
	switch ctx.KPOrdinal {
	case 1:
		pc := 100 * (nw[0] - base[0]) / base[0]
		ans := edakit.Num(pc, 2)
		return edakit.Numeric(edakit.Sprintf("A placer change is run on a reference design: baseline HPWL %s um, new HPWL %s um. What is the HPWL change, in percent (negative is better)? Round to 2 decimal places.", edakit.Num(base[0], 0), edakit.Num(nw[0], 0)), ans,
			edakit.Sprintf("(%s - %s) / %s x 100 = %s%%.", edakit.Num(nw[0], 0), edakit.Num(base[0], 0), edakit.Num(base[0], 0), ans), map[string]string{"pct": ans})
	case 2:
		tol := edakit.Pick([]float64{0.5, 1, 2})
		ot := edakit.Pick([]float64{0.08, 0.1, 0.12})
		fail := 0
		for i := range base {
			if 100*(nw[i]-base[i])/base[i] > tol || of[i] > ot {
				fail++
			}
		}
		return edakit.Numeric(edakit.Sprintf("A placer plug-in is promoted only if, on every regression design, HPWL grows by no more than %s%% over baseline and final overflow is at most %s. Results: %s. On how many designs does the gate fail?", edakit.Num(tol, 1), edakit.Num(ot, 2), strings.Join(rows, "; ")), edakit.IntStr(fail),
			edakit.Sprintf("Check both conditions per design; %d designs break at least one, so the change is %s.", fail, map[bool]string{true: "blocked", false: "allowed to proceed to review"}[fail > 0]),
			map[string]string{"fails": edakit.IntStr(fail)})
	}
	lg := 0.0
	for i := range base {
		lg += math.Log(nw[i] / base[i])
	}
	gm := math.Exp(lg / float64(n))
	ans := edakit.Num(gm, 4)
	return edakit.Numeric(edakit.Sprintf("A placer change is summarised over a regression suite by the geometric mean of new-over-baseline HPWL ratios. Results: %s. What is the geometric mean ratio? Round to 4 decimal places.", strings.Join(rows, "; ")), ans,
		edakit.Sprintf("Average the logarithms of the %d ratios and exponentiate: %s. The geometric mean keeps one large design from dominating the summary.", n, ans),
		map[string]string{"geomean": ans})
}

// fastLoopCalc counts the candidates an agent must re-run with the accurate placer.
func fastLoopCalc(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Int(5, 8)
	e := edakit.Pick([]float64{0.02, 0.03, 0.05})
	est := make([]float64, n)
	for i := range est {
		est[i] = float64(edakit.Int(900, 1100))
	}
	s := append([]float64(nil), est...)
	sort.Float64s(s)
	bestHi := s[0] * (1 + e)
	cnt := 0
	for _, v := range est {
		if v*(1-e) <= bestHi {
			cnt++
		}
	}
	q := edakit.Sprintf("An agent explores %d floorplan variants with a fast placer whose HPWL is known to be within %s%% of the accurate placer's. The fast HPWL estimates (mm) are %s. A variant can be discarded only if even its most optimistic true HPWL is worse than the best variant's most pessimistic one. How many variants must the agent re-run with the accurate placer?", n, edakit.Num(100*e, 0), fmtFloats(est, 0))
	opts := edakit.Distinct(cnt, 3, 1, n, cnt+1, cnt-1)
	var d []string
	for _, o := range opts {
		if o != "0" && !strings.HasPrefix(o, "-") {
			d = append(d, plural(o))
		}
	}
	for k := cnt + 2; len(d) < 3; k++ {
		d = append(d, plural(edakit.IntStr(k)))
	}
	return edakit.Choice(q, plural(edakit.IntStr(cnt)), edakit.Sprintf("The best estimate is %s, whose pessimistic value is %s. A variant survives if its optimistic value, estimate x %s, is at most that: %d do.", edakit.Num(s[0], 0), edakit.Num(bestHi, 2), edakit.Num(1-e, 2), cnt), d...)
}

func plural(n string) string {
	if n == "1" {
		return "1 variant"
	}
	return n + " variants"
}
