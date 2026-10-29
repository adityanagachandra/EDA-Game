package hls

import (
	"math/bits"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func initiationInterval(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		f := edakit.Pick([]int{100, 125, 150, 200, 250, 300, 400})
		ii := edakit.Int(1, 9)
		t := float64(f) / float64(ii)
		q := edakit.Sprintf("A pipelined loop runs at %d MHz with an initiation interval of %d cycles, consuming one input sample per iteration. What is its steady-state throughput in millions of samples per second? Round to 2 decimal places.", f, ii)
		return num(q, edakit.Num(t, 2),
			edakit.Sprintf("A new iteration starts every %d cycles, so throughput is %d / %d = %s million samples per second, independent of the iteration latency.", ii, f, ii, edakit.Num(t, 2)),
			map[string]string{"throughput": edakit.Num(t, 2)})
	case 2:
		f := edakit.Pick([]int{100, 150, 200, 250, 300, 400})
		r := edakit.Int(f/16+1, f/2)
		ii := f / r
		q := edakit.Sprintf("A filter must process %d million samples per second, one sample per loop iteration, on a %d MHz clock. What is the largest initiation interval, in whole cycles, that still meets the rate?", r, f)
		return num(q, edakit.IntStr(ii),
			edakit.Sprintf("The loop may spend at most %d / %d cycles per sample; rounding down to a whole cycle gives II = %d.", f, r, ii),
			map[string]string{"ii": edakit.IntStr(ii)})
	default:
		k := edakit.Scaled(ctx, 2, 4)
		set := map[int]bool{0: true}
		for len(set) < k {
			set[edakit.Int(1, 8)] = true
		}
		var offs []int
		for o := range set {
			offs = append(offs, o)
		}
		sort.Ints(offs)
		ii := len(offs)
		for ; ; ii++ {
			seen := map[int]bool{}
			ok := true
			for _, o := range offs {
				if seen[o%ii] {
					ok = false
					break
				}
				seen[o%ii] = true
			}
			if ok {
				break
			}
		}
		q := edakit.Sprintf("One iteration of a loop uses the single multiplier at cycle offsets %s relative to the iteration's start. Iterations start every II cycles. What is the smallest II at which no two iterations ever need the multiplier in the same cycle?", joinInts(offs))
		return num(q, edakit.IntStr(ii),
			edakit.Sprintf("Two uses collide exactly when their offsets are congruent modulo II, so II must give %d distinct residues. Trying II = %d, %d, ... the first that works is %d.", len(offs), len(offs), len(offs)+1, ii),
			map[string]string{"ii": edakit.IntStr(ii)})
	}
}

func loopPipelining(ctx generator.GeneratorContext) generator.Problem {
	N := edakit.Scaled(ctx, 8, 1000)
	L := edakit.Int(4, 30)
	ii := edakit.Int(1, 4)
	if ii > L {
		ii = 1
	}
	tot := (N-1)*ii + L
	switch whichKP(ctx) {
	case 1:
		q := edakit.Sprintf("A loop of %d iterations is pipelined with II = %d. One iteration takes %d cycles from start to finish. How many cycles pass from the first iteration's start to the last iteration's finish?", N, ii, L)
		return num(q, edakit.IntStr(tot),
			edakit.Sprintf("The last iteration starts %d × %d cycles after the first and needs %d more, so the total is (%d - 1) × %d + %d = %d.", N-1, ii, L, N, ii, L, tot),
			map[string]string{"total": edakit.IntStr(tot)})
	case 2:
		f := ceilDiv(L, ii)
		q := edakit.Sprintf("A loop is pipelined with II = %d and an iteration latency of %d cycles. In steady state, how many iterations are in flight at once (at most)?", ii, L)
		return num(q, edakit.IntStr(f),
			edakit.Sprintf("An iteration lives %d cycles and a new one starts every %d, so the ceiling of %d / %d = %d overlap. That is the number of pipeline stages.", L, ii, L, ii, f),
			map[string]string{"in_flight": edakit.IntStr(f)})
	default:
		seq := N * L
		sp := float64(seq) / float64(tot)
		q := edakit.Sprintf("A loop of %d iterations with an iteration latency of %d cycles is pipelined with II = %d. What is its speedup over running the iterations back to back? Round to 2 decimal places.", N, L, ii)
		return num(q, edakit.Num(sp, 2),
			edakit.Sprintf("Sequential: %d × %d = %d cycles. Pipelined: (%d - 1) × %d + %d = %d. Speedup %d / %d = %s, which approaches L / II = %s for long loops.", N, L, seq, N, ii, L, tot, seq, tot, edakit.Num(sp, 2), edakit.Num(float64(L)/float64(ii), 2)),
			map[string]string{"pipelined": edakit.IntStr(tot), "sequential": edakit.IntStr(seq)})
	}
}

func unrolling(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		U := edakit.Pick([]int{2, 3, 4, 5, 8})
		N := edakit.Int(10, 200)
		for N%U == 0 {
			N++
		}
		q := edakit.Sprintf("A loop of %d iterations is unrolled by a factor of %d without knowing that the trip count is a multiple of %d, so the compiler emits a main unrolled loop and a remainder loop. How many original iterations run in the remainder loop?", N, U, U)
		return num(q, edakit.IntStr(N%U),
			edakit.Sprintf("The main loop runs %d times covering %d iterations; %d mod %d = %d remain.", N/U, (N/U)*U, N, U, N%U),
			map[string]string{"main_trips": edakit.IntStr(N / U)})
	case 2:
		a, b := edakit.Int(1, 3), edakit.Int(1, 3)
		D := edakit.Int(2, 5)
		km, ka := edakit.Int(1, 4), edakit.Int(1, 4)
		U := edakit.Pick([]int{2, 4, 8})
		N := U * edakit.Int(4, 32)
		body := maxInt(D, ceilDiv(U*a, km), ceilDiv(U*b, ka))
		tot := (N / U) * body
		q := edakit.Sprintf("Each iteration of a loop has %d multiplication(s) and %d addition(s), all single-cycle, with a dependence depth of %d steps, and iterations are independent. The loop runs %d iterations, is unrolled by %d and not pipelined. There are %d multipliers and %d adders. Assume each unrolled body takes exactly the larger of its depth and its resource bounds. How many cycles does the loop take?", a, b, D, N, U, km, ka)
		return num(q, edakit.IntStr(tot),
			edakit.Sprintf("One unrolled body needs max(%d, ceil(%d/%d), ceil(%d/%d)) = %d cycles and runs %d times, %d cycles in all.", D, U*a, km, U*b, ka, body, N/U, tot),
			map[string]string{"body": edakit.IntStr(body)})
	default:
		D := edakit.Int(2, 6)
		km, ka := edakit.Int(1, 4), edakit.Int(1, 6)
		a := edakit.Int(1, minInt(4, D*km))
		b := edakit.Int(1, minInt(4, D*ka))
		u := minInt(D*km/a, D*ka/b)
		q := edakit.Sprintf("Each iteration of a loop has %d multiplication(s) and %d addition(s), all single-cycle, with a dependence depth of %d steps; iterations are independent. There are %d multipliers and %d adders. When the loop is unrolled by U, one unrolled body takes the larger of %d and its resource bounds. What is the largest U for which the body still takes only %d cycles?", a, b, D, km, ka, D, D)
		return num(q, edakit.IntStr(u),
			edakit.Sprintf("The multipliers allow U × %d ≤ %d × %d, so U ≤ %d; the adders allow U × %d ≤ %d × %d, so U ≤ %d. The binding limit is %d; beyond it unrolling adds area without cutting cycles.", a, D, km, D*km/a, b, D, ka, D*ka/b, u),
			map[string]string{"u_max": edakit.IntStr(u)})
	}
}

func resMII(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		uses := edakit.Int(2, 16)
		units := edakit.Int(1, 4)
		r := ceilDiv(uses, units)
		q := edakit.Sprintf("One iteration of a loop body performs %d multiplications and the allocation has %d fully pipelined multipliers. What is the resource-constrained minimum initiation interval ResMII from the multipliers?", uses, units)
		return num(q, edakit.IntStr(r),
			edakit.Sprintf("Every II cycles the multipliers offer %d × II slots, which must cover %d uses, so ResMII = ceil(%d / %d) = %d.", units, uses, uses, units, r),
			map[string]string{"resmii": edakit.IntStr(r)})
	case 2:
		nm, km := edakit.Int(2, 12), edakit.Int(1, 4)
		na, ka := edakit.Int(2, 12), edakit.Int(1, 4)
		nd, kd := edakit.Int(1, 3), 1
		bm, ba, bd := ceilDiv(nm, km), ceilDiv(na, ka), ceilDiv(nd, kd)
		r := maxInt(bm, ba, bd, 1)
		q := edakit.Sprintf("One iteration of a loop body performs %d multiplications, %d additions and %d divisions. The allocation has %d multipliers, %d adders and %d divider, all fully pipelined. What is ResMII?", nm, na, nd, km, ka, kd)
		return num(q, edakit.IntStr(r),
			edakit.Sprintf("Per resource: multipliers ceil(%d/%d) = %d, adders ceil(%d/%d) = %d, divider ceil(%d/%d) = %d. ResMII is the maximum, %d.", nm, km, bm, na, ka, ba, nd, kd, bd, r),
			map[string]string{"resmii": edakit.IntStr(r)})
	default:
		rd, wr := edakit.Int(1, 6), edakit.Int(1, 3)
		ports := edakit.Int(1, 2)
		nm, km := edakit.Int(1, 8), edakit.Int(1, 4)
		bmem, bm := ceilDiv(rd+wr, ports), ceilDiv(nm, km)
		r := maxInt(bmem, bm)
		q := edakit.Sprintf("One iteration reads array A %d times and writes it %d times. A sits in a single block RAM with %d port(s), each able to perform one read or one write per cycle. The iteration also performs %d multiplications on %d pipelined multipliers. What is ResMII?", rd, wr, ports, nm, km)
		return num(q, edakit.IntStr(r),
			edakit.Sprintf("Memory: %d accesses over %d port(s) gives ceil(%d/%d) = %d. Multipliers: ceil(%d/%d) = %d. ResMII = %d.", rd+wr, ports, rd+wr, ports, bmem, nm, km, bm, r),
			map[string]string{"mem_bound": edakit.IntStr(bmem), "mul_bound": edakit.IntStr(bm)})
	}
}

func recMII(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		l, d := edakit.Int(2, 20), edakit.Int(1, 3)
		r := ceilDiv(l, d)
		q := edakit.Sprintf("A loop has one recurrence: a cycle in its dependence graph whose operations have a total latency of %d cycles and whose loop-carried edges have a total iteration distance of %d. What is RecMII?", l, d)
		return num(q, edakit.IntStr(r),
			edakit.Sprintf("Around the cycle the schedule needs %d ≤ II × %d, so II ≥ ceil(%d / %d) = %d.", l, d, l, d, r),
			map[string]string{"recmii": edakit.IntStr(r)})
	case 2:
		k := edakit.Int(2, 4)
		best := 0
		var p, w []string
		for i := 0; i < k; i++ {
			l, d := edakit.Int(2, 18), edakit.Int(1, 3)
			c := ceilDiv(l, d)
			best = maxInt(best, c)
			p = append(p, edakit.Sprintf("latency %d with distance %d", l, d))
			w = append(w, edakit.Sprintf("ceil(%d/%d) = %d", l, d, c))
		}
		q := "A loop's dependence graph has these elementary cycles: " + strings.Join(p, "; ") + ". What is RecMII?"
		return num(q, edakit.IntStr(best),
			"Each cycle gives a bound: "+strings.Join(w, ", ")+edakit.Sprintf(". RecMII is the largest, %d.", best),
			map[string]string{"recmii": edakit.IntStr(best)})
	default:
		n := edakit.Int(3, 5)
		lat := make([]int, n)
		for i := range lat {
			lat[i] = edakit.Pick([]int{1, 2, 3, 4, 5})
		}
		type be struct{ from, to, d int }
		var bes []be
		var bt []string
		for len(bes) < 2 {
			j := edakit.Int(1, n-1)
			i := edakit.Int(0, j-1)
			if edakit.Int(0, 5) == 0 {
				i = j
			}
			dup := false
			for _, b := range bes {
				if b.from == j && b.to == i {
					dup = true
				}
			}
			if dup {
				continue
			}
			d := edakit.Int(1, 2)
			bes = append(bes, be{j, i, d})
			bt = append(bt, edakit.Sprintf("o%d → o%d with distance %d", j+1, i+1, d))
		}
		var lt []string
		for i, l := range lat {
			lt = append(lt, edakit.Sprintf("o%d: %d", i+1, l))
		}
		best := 0
		var w []string
		for _, b := range bes {
			s := 0
			for k := b.to; k <= b.from; k++ {
				s += lat[k]
			}
			c := ceilDiv(s, b.d)
			best = maxInt(best, c)
			w = append(w, edakit.Sprintf("o%d..o%d: ceil(%d/%d) = %d", b.to+1, b.from+1, s, b.d, c))
		}
		var chain []string
		for i := 0; i < n; i++ {
			chain = append(chain, edakit.Sprintf("o%d", i+1))
		}
		q := edakit.Sprintf("A loop body is a chain %s of dependent operations with latencies %s cycles. Loop-carried edges feed results into a later iteration: %s. Each loop-carried edge closes a cycle through the chain. What is RecMII?", strings.Join(chain, " → "), strings.Join(lt, ", "), strings.Join(bt, "; "))
		return num(q, edakit.IntStr(best),
			"Each loop-carried edge from oj back to oi forms a cycle with latency equal to the sum over oi..oj: "+strings.Join(w, ", ")+edakit.Sprintf(". RecMII = %d.", best),
			map[string]string{"recmii": edakit.IntStr(best)})
	}
}

func moduloSched(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		ii := edakit.Int(2, 5)
		k := edakit.Scaled(ctx, 3, 6)
		var ts []int
		rows := map[int]int{}
		m := 0
		for i := 0; i < k; i++ {
			t := edakit.Int(0, 14)
			ts = append(ts, t)
			rows[t%ii]++
			m = maxInt(m, rows[t%ii])
		}
		sort.Ints(ts)
		q := edakit.Sprintf("A modulo schedule with II = %d places one iteration's multiplications at start times %s (cycles from the iteration's start). Each multiplication occupies its multiplier for one cycle. In the modulo reservation table, with rows 0 to %d, how many multiplications fall in the most crowded row, and hence how many multipliers does this schedule need?", ii, joinInts(ts), ii-1)
		var rs []int
		for _, t := range ts {
			rs = append(rs, t%ii)
		}
		return num(q, edakit.IntStr(m),
			edakit.Sprintf("Row = start time mod %d gives %s. The fullest row holds %d multiplications, because in steady state all of them execute in the same cycle, from different iterations.", ii, joinInts(rs), m),
			map[string]string{"rows": joinInts(rs)})
	case 2:
		ii := edakit.Int(3, 6)
		occ := map[int]bool{}
		k := edakit.Int(1, ii-1)
		for len(occ) < k {
			occ[edakit.Int(0, ii-1)] = true
		}
		est := edakit.Int(0, 12)
		t := est
		for occ[t%ii] {
			t++
		}
		var o []int
		for r := range occ {
			o = append(o, r)
		}
		sort.Ints(o)
		q := edakit.Sprintf("Iterative modulo scheduling with II = %d and a single multiplier. Earlier multiplications already occupy MRT rows %s. The next multiplication's dependences allow it to start no earlier than cycle %d. What is the earliest start time at which it can be placed without a resource conflict?", ii, joinInts(o), est)
		return num(q, edakit.IntStr(t),
			edakit.Sprintf("Try t = %d, %d, ... and test row t mod %d. The first free row is reached at t = %d (row %d).", est, est+1, ii, t, t%ii),
			map[string]string{"start": edakit.IntStr(t)})
	default:
		ii := edakit.Int(1, 5)
		sl := edakit.Int(ii+1, 30)
		st := ceilDiv(sl, ii)
		pro := (st - 1) * ii
		q := edakit.Sprintf("A modulo-scheduled loop has II = %d and a single-iteration schedule length of %d cycles. The kernel overlaps all stages, so the pipeline has ceil(%d / %d) stages. How many cycles does the prologue take before the kernel runs with every stage full?", ii, sl, sl, ii)
		return num(q, edakit.IntStr(pro),
			edakit.Sprintf("There are %d stages of %d cycles each. The prologue fills stages one by one for %d stages, which is (%d - 1) × %d = %d cycles; the epilogue drains the same amount.", st, ii, st-1, st, ii, pro),
			map[string]string{"stages": edakit.IntStr(st)})
	}
}

func recurrenceBreaking(ctx generator.GeneratorContext) generator.Problem {
	la := edakit.Int(3, 8)
	lm := edakit.Int(2, 6)
	switch whichKP(ctx) {
	case 1:
		forms := []struct {
			code string
			ii   int
			why  string
		}{
			{"s = s + x[i]", la, edakit.Sprintf("Only the addition is on the cycle, so II ≥ %d.", la)},
			{"s = s + a[i] × b[i]", la, edakit.Sprintf("The multiplication reads only array values and is off the cycle; only the %d-cycle addition is on it, so II ≥ %d.", la, la)},
			{"s = s × c + x[i]", la + lm, edakit.Sprintf("Both the multiplication and the addition sit on the cycle, so II ≥ %d + %d = %d.", lm, la, la+lm)},
		}
		f := edakit.Pick(forms)
		q := edakit.Sprintf("A loop computes %s for i = 0 to N-1, where s carries from one iteration to the next. The adder has a latency of %d cycles and the multiplier %d cycles, both pipelined, with resources unlimited. What is the minimum initiation interval?", f.code, la, lm)
		return num(q, edakit.IntStr(f.ii), "The recurrence through s has distance 1. "+f.why,
			map[string]string{"recmii": edakit.IntStr(f.ii)})
	case 2:
		target := edakit.Int(1, minInt(3, la-1))
		p := ceilDiv(la, target)
		q := edakit.Sprintf("A floating-point accumulation s = s + x[i] uses an adder with a latency of %d cycles. Rewriting it to keep P interleaved partial sums (iteration i updates partial sum i mod P) makes each recurrence distance P. What is the smallest P that achieves II = %d?", la, target)
		return num(q, edakit.IntStr(p),
			edakit.Sprintf("With distance P the recurrence bound is ceil(%d / P), which is at most %d once P ≥ ceil(%d / %d) = %d.", la, target, la, target, p),
			map[string]string{"p": edakit.IntStr(p)})
	default:
		p := la
		N := p * edakit.Int(8, 128)
		lv := bits.Len(uint(p - 1))
		tot := (N - 1) + la + lv*la
		q := edakit.Sprintf("An accumulation of N = %d values uses an adder with a latency of %d cycles. It is rewritten with P = %d interleaved partial sums so the loop runs at II = 1: its last addition completes %d cycles after it is issued at cycle %d. Then the P partial sums are combined by a balanced tree of additions, each tree level taking %d cycles. Counting from cycle 0, when does the final sum complete?", N, la, p, la, N-1, la)
		return num(q, edakit.IntStr(tot),
			edakit.Sprintf("Main loop: (%d - 1) + %d = %d cycles. A balanced tree over %d partial sums has ceil(log2 %d) = %d levels, adding %d × %d = %d. Total %d cycles, against %d × %d = %d for the original II = %d loop.", N, la, N-1+la, p, p, lv, lv, la, lv*la, tot, N, la, N*la, la),
			map[string]string{"levels": edakit.IntStr(lv), "total": edakit.IntStr(tot)})
	}
}
