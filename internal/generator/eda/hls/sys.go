package hls

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func codesign(ctx generator.GeneratorContext) generator.Problem {
	T := float64(10 * edakit.Int(5, 100))
	f := float64(edakit.Int(40, 95)) / 100
	s := float64(edakit.Int(4, 60))
	c := edakit.Round(T*f*float64(edakit.Int(2, 25))/100, 1)
	newT := (1-f)*T + f*T/s + c
	head := edakit.Sprintf("A program takes %s ms in software. A kernel accounting for %s%% of that time is moved to an HLS accelerator that runs it %s times faster, and moving data to and from the accelerator adds %s ms.",
		edakit.Num(T, 0), edakit.Num(f*100, 0), edakit.Num(s, 0), edakit.Num(c, 1))
	switch whichKP(ctx) {
	case 1:
		return num(head+" What is the new run time in ms? Round to 2 decimal places.", edakit.Num(newT, 2),
			edakit.Sprintf("Software part (1 - %s) × %s = %s ms, accelerated part %s / %s = %s ms, plus %s ms of transfer: %s ms.", edakit.Num(f, 2), edakit.Num(T, 0), edakit.Num((1-f)*T, 2), edakit.Num(f*T, 2), edakit.Num(s, 0), edakit.Num(f*T/s, 2), edakit.Num(c, 1), edakit.Num(newT, 2)),
			map[string]string{"new_time": edakit.Num(newT, 2)})
	case 2:
		sp := T / newT
		return num(head+" What is the overall speedup? Round to 2 decimal places.", edakit.Num(sp, 2),
			edakit.Sprintf("The new time is %s ms, so the speedup is %s / %s = %s. Amdahl's law caps it below 1 / (1 - %s) = %s however fast the accelerator is.", edakit.Num(newT, 2), edakit.Num(T, 0), edakit.Num(newT, 2), edakit.Num(sp, 2), edakit.Num(f, 2), edakit.Num(1/(1-f), 2)),
			map[string]string{"speedup": edakit.Num(sp, 2)})
	default:
		be := f * T / (f*T - c)
		q := edakit.Sprintf("A program takes %s ms in software. A kernel accounting for %s%% of that time is a candidate for an HLS accelerator, and moving data to and from it would add %s ms. What accelerator speedup on the kernel is needed just to break even with software? Round to 2 decimal places.", edakit.Num(T, 0), edakit.Num(f*100, 0), edakit.Num(c, 1))
		return num(q, edakit.Num(be, 2),
			edakit.Sprintf("Break-even needs %s / s + %s = %s, so s = %s / (%s - %s) = %s.", edakit.Num(f*T, 2), edakit.Num(c, 1), edakit.Num(f*T, 2), edakit.Num(f*T, 2), edakit.Num(f*T, 2), edakit.Num(c, 1), edakit.Num(be, 2)),
			map[string]string{"break_even": edakit.Num(be, 2)})
	}
}

type pt struct{ lat, area int }

func paretoFront(ps []pt) []pt {
	var out []pt
	for i, p := range ps {
		dom := false
		for j, q := range ps {
			if i != j && q.lat <= p.lat && q.area <= p.area && (q.lat < p.lat || q.area < p.area) {
				dom = true
				break
			}
		}
		if !dom {
			out = append(out, p)
		}
	}
	sort.Slice(out, func(a, b int) bool { return out[a].lat < out[b].lat })
	return out
}

func randPoints(n int) []pt {
	seen := map[pt]bool{}
	var ps []pt
	for len(ps) < n {
		l := edakit.Int(1, 20)
		a := 21 - l + edakit.Int(-5, 5)
		if a < 1 {
			a = 1
		}
		if a > 20 {
			a = 20
		}
		p := pt{10 * l, 100 * a}
		if !seen[p] {
			seen[p] = true
			ps = append(ps, p)
		}
	}
	return ps
}

func ptText(ps []pt) string {
	var s []string
	for i, p := range ps {
		s = append(s, edakit.Sprintf("D%d (%d cycles, %d LUTs)", i+1, p.lat, p.area))
	}
	return strings.Join(s, ", ")
}

func dse(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		loops := edakit.Int(2, 4)
		uo := edakit.Pick([]int{3, 4, 5})
		arrays := edakit.Int(1, 3)
		po := edakit.Pick([]int{3, 4})
		size := int(math.Pow(float64(uo*2), float64(loops)) * math.Pow(float64(po), float64(arrays)))
		q := edakit.Sprintf("A kernel has %d loops; each can take one of %d unroll factors and may be pipelined or not. It has %d arrays; each can take one of %d partitioning choices. All choices are independent. How many configurations does the pragma design space contain?", loops, uo, arrays, po)
		return num(q, edakit.IntStr(size),
			edakit.Sprintf("Each loop has %d × 2 = %d settings and each array %d, so the space has %d^%d × %d^%d = %d points.", uo, uo*2, po, uo*2, loops, po, arrays, size),
			map[string]string{"size": edakit.IntStr(size)})
	case 2:
		ps := randPoints(edakit.Scaled(ctx, 5, 9))
		fr := paretoFront(ps)
		q := "HLS runs of one kernel under different pragmas report " + ptText(ps) + ". Lower is better in both. How many of these designs are Pareto-optimal?"
		return num(q, edakit.IntStr(len(fr)),
			edakit.Sprintf("A design is Pareto-optimal when no other is at least as good in both metrics and better in one. The front is %s: %d designs.", ptText2(fr), len(fr)),
			map[string]string{"front": edakit.IntStr(len(fr))})
	default:
		ps := randPoints(edakit.Scaled(ctx, 4, 7))
		fr := paretoFront(ps)
		rl, ra := 220, 2200
		hv := 0
		for i, p := range fr {
			next := rl
			if i+1 < len(fr) {
				next = fr[i+1].lat
			}
			hv += (next - p.lat) * (ra - p.area)
		}
		q := "HLS runs report " + ptText(ps) + edakit.Sprintf(". Lower is better in both. The hypervolume indicator is the area of the region dominated by the Pareto front and bounded by the reference point (%d cycles, %d LUTs). What is it, in cycle-LUTs?", rl, ra)
		return num(q, edakit.IntStr(hv),
			edakit.Sprintf("The front, sorted by latency, is %s. Summing the rectangles from each front point to the next point's latency and up to the reference area gives %d.", ptText2(fr), hv),
			map[string]string{"hypervolume": edakit.IntStr(hv)})
	}
}

func ptText2(ps []pt) string {
	var s []string
	for _, p := range ps {
		s = append(s, edakit.Sprintf("(%d, %d)", p.lat, p.area))
	}
	return strings.Join(s, ", ")
}
