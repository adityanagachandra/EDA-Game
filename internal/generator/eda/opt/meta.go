package opt

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func localSearch(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 8, 14)
	cost := distinctInts(n, 10, 99)
	land := "Configurations 0 to " + edakit.IntStr(n-1) + " of a partitioner have costs " + ints(cost, ", ") + "."
	switch kp(ctx) {
	case 1:
		cnt := 0
		var at []string
		for i := 0; i < n; i++ {
			l := i == 0 || cost[i] < cost[i-1]
			r := i == n-1 || cost[i] < cost[i+1]
			if l && r {
				cnt++
				at = append(at, edakit.IntStr(i))
			}
		}
		return num(land+" The move set changes the configuration index by plus or minus 1. How many local minima (no neighbour is cheaper) are there?"+roundNote(0), edakit.IntStr(cnt),
			edakit.Sprintf("A configuration is a local minimum when every neighbour costs more; the endpoints have one neighbour. The local minima are at configurations %s: %d in all. Only one of them is global, which is why plain descent is not enough.", strings.Join(at, ", "), cnt),
			"local_minima", edakit.IntStr(cnt))
	case 2:
		s := edakit.Int(0, n-1)
		i := s
		for {
			best := i
			for _, j := range []int{i - 1, i + 1} {
				if j >= 0 && j < n && cost[j] < cost[best] {
					best = j
				}
			}
			if best == i {
				break
			}
			i = best
		}
		return num(land+edakit.Sprintf(" Steepest descent starts at configuration %d, moves to the cheaper of its plus-or-minus-1 neighbours while that improves, and stops otherwise. At what cost does it stop?", s)+roundNote(0), edakit.IntStr(cost[i]),
			edakit.Sprintf("Following strictly improving moves from %d ends at configuration %d with cost %d, a local minimum of this move set.", s, i, cost[i]),
			"stop", edakit.IntStr(i), "cost", edakit.IntStr(cost[i]))
	default:
		s := edakit.Int(0, n-1)
		tenure := edakit.Int(2, 3)
		steps := edakit.Int(4, 7)
		i := s
		best := cost[s]
		recent := []int{s}
		var path []string
		path = append(path, edakit.IntStr(s))
		for k := 0; k < steps; k++ {
			cand := -1
			for _, j := range []int{i - 2, i - 1, i + 1, i + 2} {
				if j < 0 || j >= n {
					continue
				}
				tabu := false
				for _, r := range recent {
					if r == j {
						tabu = true
					}
				}
				if tabu {
					continue
				}
				if cand < 0 || cost[j] < cost[cand] {
					cand = j
				}
			}
			if cand < 0 {
				break
			}
			i = cand
			path = append(path, edakit.IntStr(i))
			if cost[i] < best {
				best = cost[i]
			}
			recent = append(recent, i)
			if len(recent) > tenure+1 {
				recent = recent[1:]
			}
		}
		return num(land+edakit.Sprintf(" Tabu search starts at configuration %d. Each move goes to the cheapest neighbour within plus or minus 2 that is not tabu, even if it is worse. The current configuration and the last %d visited ones are tabu. After %d moves, what is the lowest cost seen?", s, tenure, steps)+roundNote(0), edakit.IntStr(best),
			edakit.Sprintf("The walk visits %s. Forced moves to worse neighbours let it leave local minima, and the tabu list stops it cycling straight back. The best cost seen is %d.", strings.Join(path, ", "), best),
			"best", edakit.IntStr(best))
	}
}

func annealing(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		d := float64(edakit.Int(1, 40))
		T := float64(edakit.Int(2, 60))
		for d/T > 4 {
			T = float64(edakit.Int(2, 60))
		}
		p := math.Exp(-d / T)
		return num(edakit.Sprintf("A simulated-annealing placer proposes a swap that raises HPWL by $\\Delta = %s$ um at temperature T = %s um. With the Metropolis rule, what is the probability the swap is accepted?", f(d, 0), f(T, 0))+roundNote(3), f(p, 3),
			edakit.Sprintf("Uphill moves are accepted with probability $e^{-\\Delta/T} = e^{-%s/%s} = %s$. Downhill moves are always accepted.", f(d, 0), f(T, 0), f(p, 3)),
			"p", f(p, 3))
	case 2:
		d := float64(edakit.Int(2, 50))
		p := edakit.Pick([]float64{0.1, 0.2, 0.25, 0.3, 0.4, 0.5, 0.6, 0.8, 0.9})
		T := -d / math.Log(p)
		return num(edakit.Sprintf("A floorplanner wants a typical uphill move of $\\Delta = %s$ area units to be accepted with probability %s. What temperature T gives that?", f(d, 0), f(p, 2))+roundNote(2), f(T, 2),
			edakit.Sprintf("Solve $e^{-\\Delta/T} = p$: $T = -\\Delta/\\ln p = -%s/\\ln %s = %s$.", f(d, 0), f(p, 2), f(T, 2)),
			"t", f(T, 2))
	default:
		n := edakit.Scaled(ctx, 3, 6)
		T := float64(edakit.Int(5, 30))
		ds := make([]int, n)
		exp := 0.0
		var terms []string
		for i := range ds {
			ds[i] = edakit.Int(-10, 30)
			pr := 1.0
			if ds[i] > 0 {
				pr = math.Exp(-float64(ds[i]) / T)
			}
			exp += pr
			terms = append(terms, f(pr, 3))
		}
		return num(edakit.Sprintf("At T = %s, an annealer proposes moves with cost changes %s. What is the expected number of accepted moves?", f(T, 0), ints(ds, ", "))+roundNote(2), f(exp, 2),
			edakit.Sprintf("Each move contributes its acceptance probability, 1 if $\\Delta \\le 0$ and $e^{-\\Delta/T}$ otherwise: %s. Their sum is %s.", strings.Join(terms, ", "), f(exp, 2)),
			"expected", f(exp, 2))
	}
}

func cooling(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		T0 := float64(10 * edakit.Int(5, 100))
		a := edakit.Pick([]float64{0.8, 0.85, 0.9, 0.95, 0.98})
		k := edakit.Scaled(ctx, 3, 20)
		Tk := T0 * math.Pow(a, float64(k))
		return num(edakit.Sprintf("A geometric cooling schedule sets $T_{k+1} = \\alpha T_k$ with $T_0 = %s$ and $\\alpha = %s$. What is $T_{%d}$?", f(T0, 0), f(a, 2), k)+roundNote(3), f(Tk, 3),
			edakit.Sprintf("$T_k = T_0 \\alpha^k = %s \\cdot %s^{%d} = %s$.", f(T0, 0), f(a, 2), k, f(Tk, 3)),
			"tk", f(Tk, 3))
	case 2:
		for {
			T0 := float64(10 * edakit.Int(5, 100))
			Tf := edakit.Pick([]float64{0.01, 0.05, 0.1, 0.5, 1})
			a := edakit.Pick([]float64{0.8, 0.85, 0.9, 0.95})
			x := math.Log(Tf/T0) / math.Log(a)
			if math.Abs(x-math.Round(x)) < 1e-6 {
				continue
			}
			k := int(math.Ceil(x))
			return num(edakit.Sprintf("An annealer cools geometrically from $T_0 = %s$ with $\\alpha = %s$ and stops once $T \\le %s$. How many temperature steps does it take?", f(T0, 0), f(a, 2), f(Tf, 2))+roundNote(0), edakit.IntStr(k),
				edakit.Sprintf("Need $T_0\\alpha^k \\le T_f$, so $k \\ge \\ln(T_f/T_0)/\\ln\\alpha = %s$; the smallest integer is %d. Each step also runs many moves, so the move budget is steps times moves per temperature.", f(x, 3), k),
				"steps", edakit.IntStr(k))
		}
	default:
		d := float64(edakit.Int(2, 40))
		chi := edakit.Pick([]float64{0.8, 0.85, 0.9, 0.95})
		T0 := -d / math.Log(chi)
		return num(edakit.Sprintf("Before annealing, a placer samples random moves; the uphill ones average $\\overline{\\Delta^+} = %s$ um. It wants an initial acceptance ratio of about %s for uphill moves. Using $T_0 = -\\overline{\\Delta^+}/\\ln \\chi_0$, what is $T_0$?", f(d, 0), f(chi, 2))+roundNote(2), f(T0, 2),
			edakit.Sprintf("$T_0 = -%s/\\ln %s = %s$. Starting hot enough that most moves are accepted lets the early search roam; the schedule then cools toward a greedy descent.", f(d, 0), f(chi, 2), f(T0, 2)),
			"t0", f(T0, 2))
	}
}

type blk struct{ w, h int }

// evalPolish evaluates a postfix slicing expression: V puts operands side by side, H stacks them.
func evalPolish(tokens []string, dims map[string]blk) blk {
	var st []blk
	for _, t := range tokens {
		switch t {
		case "V":
			b, a := st[len(st)-1], st[len(st)-2]
			st = st[:len(st)-2]
			st = append(st, blk{a.w + b.w, max(a.h, b.h)})
		case "H":
			b, a := st[len(st)-1], st[len(st)-2]
			st = st[:len(st)-2]
			st = append(st, blk{max(a.w, b.w), a.h + b.h})
		default:
			st = append(st, dims[t])
		}
	}
	return st[0]
}

// randomPolish builds a random postfix expression over the given operands.
func randomPolish(ops []string) []string {
	if len(ops) == 1 {
		return []string{ops[0]}
	}
	k := edakit.Int(1, len(ops)-1)
	left := randomPolish(ops[:k])
	right := randomPolish(ops[k:])
	out := append(append(left, right...), edakit.Pick([]string{"H", "V"}))
	return out
}

func polishExpression(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 3, 6)
	names := edakit.Shuffle(rangeInts(1, n))
	ops := make([]string, n)
	dims := map[string]blk{}
	var dd []string
	for i := 1; i <= n; i++ {
		dims[edakit.IntStr(i)] = blk{edakit.Int(1, 6), edakit.Int(1, 6)}
		dd = append(dd, edakit.Sprintf("block %d is %d x %d", i, dims[edakit.IntStr(i)].w, dims[edakit.IntStr(i)].h))
	}
	for i, v := range names {
		ops[i] = edakit.IntStr(v)
	}
	expr := randomPolish(ops)
	b := evalPolish(expr, dims)
	area := b.w * b.h
	stem := "Blocks (width x height): " + strings.Join(dd, ", ") + ". In the Polish expression " + strings.Join(expr, " ") + ", operator V places its two operands side by side (widths add, height is the larger) and H stacks them (heights add, width is the larger)."
	// moves
	move := func() ([]string, string) {
		out := append([]string(nil), expr...)
		if edakit.Int(0, 1) == 0 {
			var idx []int
			for i, t := range out {
				if t != "H" && t != "V" {
					idx = append(idx, i)
				}
			}
			k := edakit.Int(0, len(idx)-2)
			i, j := idx[k], idx[k+1]
			out[i], out[j] = out[j], out[i]
			return out, edakit.Sprintf("move M1 swaps the adjacent operands %s and %s", expr[i], expr[j])
		}
		var runs [][2]int
		for i := 0; i < len(out); {
			if out[i] == "H" || out[i] == "V" {
				j := i
				for j < len(out) && (out[j] == "H" || out[j] == "V") {
					j++
				}
				runs = append(runs, [2]int{i, j})
				i = j
			} else {
				i++
			}
		}
		r := edakit.Pick(runs)
		for i := r[0]; i < r[1]; i++ {
			if out[i] == "H" {
				out[i] = "V"
			} else {
				out[i] = "H"
			}
		}
		return out, edakit.Sprintf("move M2 complements the operator chain %s at positions %d to %d", strings.Join(expr[r[0]:r[1]], " "), r[0]+1, r[1])
	}
	switch kp(ctx) {
	case 1:
		return num(stem+" What is the area of the enclosing rectangle?"+roundNote(0), edakit.IntStr(area),
			edakit.Sprintf("Evaluate the postfix expression with a stack, combining the top two sub-floorplans at each operator. The final rectangle is %d x %d, area %d.", b.w, b.h, area),
			"w", edakit.IntStr(b.w), "h", edakit.IntStr(b.h), "area", edakit.IntStr(area))
	case 2:
		nexpr, desc := move()
		nb := evalPolish(nexpr, dims)
		return num(stem+" An annealing "+desc+", giving "+strings.Join(nexpr, " ")+". What is the new enclosing area?"+roundNote(0), edakit.IntStr(nb.w*nb.h),
			edakit.Sprintf("The original is %d x %d (area %d). After the move the stack evaluation gives %d x %d, area %d.", b.w, b.h, area, nb.w, nb.h, nb.w*nb.h),
			"area_before", edakit.IntStr(area), "area_after", edakit.IntStr(nb.w*nb.h))
	default:
		T := float64(edakit.Int(2, 12))
		var nexpr []string
		var desc string
		var nb blk
		for t := 0; t < 60; t++ {
			nexpr, desc = move()
			nb = evalPolish(nexpr, dims)
			if nb.w*nb.h > area {
				break
			}
		}
		d := nb.w*nb.h - area
		p := 1.0
		if d > 0 {
			p = math.Exp(-float64(d) / T)
		}
		return num(stem+" An annealing "+desc+edakit.Sprintf(", giving %s. At temperature T = %s (area units), what is the probability this move is accepted?", strings.Join(nexpr, " "), f(T, 0))+roundNote(3), f(p, 3),
			edakit.Sprintf("Area goes from %d to %d, so $\\Delta = %d$. %s", area, nb.w*nb.h, d, map[bool]string{true: edakit.Sprintf("The acceptance probability is $e^{-%d/%s} = %s$.", d, f(T, 0), f(p, 3)), false: "A move that does not increase cost is always accepted."}[d > 0]),
			"delta", edakit.IntStr(d), "p", f(p, 3))
	}
}

func genetic(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 7)
		fs := make([]int, n)
		sum := 0
		for i := range fs {
			fs[i] = edakit.Int(1, 40)
			sum += fs[i]
		}
		k := edakit.Int(1, n)
		p := float64(fs[k-1]) / float64(sum)
		return num(edakit.Sprintf("A genetic floorplanner's population has fitness values %s (individuals 1 to %d). With roulette-wheel (fitness-proportional) selection, what is the probability that one draw picks individual %d?", ints(fs, ", "), n, k)+roundNote(3), f(p, 3),
			edakit.Sprintf("Probability is fitness over total fitness: %d/%d = %s.", fs[k-1], sum, f(p, 3)),
			"p", f(p, 3))
	case 2:
		L := edakit.Scaled(ctx, 6, 10)
		A, B, w := make([]int, L), make([]int, L), make([]int, L)
		for i := 0; i < L; i++ {
			A[i], B[i], w[i] = edakit.Int(0, 1), edakit.Int(0, 1), edakit.Int(1, 9)
		}
		cut := edakit.Int(1, L-1)
		tot := 0
		for i := 0; i < L; i++ {
			g := B[i]
			if i < cut {
				g = A[i]
			}
			tot += g * w[i]
		}
		return num(edakit.Sprintf("Chromosomes assign %d cells to partition 1 (gene 1) or partition 0. Parent A is %s and parent B is %s. One-point crossover after gene %d takes genes 1 to %d from A and the rest from B. The cells' areas are %s. What total area does the child put in partition 1?", L, ints(A, ""), ints(B, ""), cut, cut, ints(w, ", "))+roundNote(0), edakit.IntStr(tot),
			edakit.Sprintf("The child is %s%s. Summing the areas where its gene is 1 gives %d.", ints(A[:cut], ""), ints(B[cut:], ""), tot),
			"area", edakit.IntStr(tot))
	default:
		N := edakit.Int(10, 100)
		k := edakit.Int(2, 7)
		p := 1 - math.Pow(1-1/float64(N), float64(k))
		return num(edakit.Sprintf("Tournament selection draws k = %d individuals uniformly with replacement from a population of N = %d and keeps the fittest. What is the probability that the single best individual wins a tournament?", k, N)+roundNote(4), f(p, 4),
			edakit.Sprintf("The best individual wins whenever it is drawn at least once: $1 - (1 - 1/N)^k = 1 - (%d/%d)^{%d} = %s$. Larger k raises selection pressure.", N-1, N, k, f(p, 4)),
			"p", f(p, 4))
	}
}
