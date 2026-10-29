package pkg

import (
	"math"
	"math/rand"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// netlist draws m nets over n cells, each net joining 2 to 4 distinct cells.
func netlist(n, m int) [][]int {
	nets := make([][]int, 0, m)
	for len(nets) < m {
		k := edakit.Int(2, 4)
		if k > n {
			k = n
		}
		perm := rand.Perm(n)[:k]
		sort.Ints(perm)
		nets = append(nets, perm)
	}
	return nets
}

func cellName(i int) string { return string(rune('a' + i)) }

func describe(nets [][]int, die []int, dies int) string {
	var ns []string
	for i, net := range nets {
		var cs []string
		for _, c := range net {
			cs = append(cs, cellName(c))
		}
		ns = append(ns, edakit.Sprintf("N%d = {%s}", i+1, strings.Join(cs, ", ")))
	}
	var as []string
	for d := 0; d < dies; d++ {
		var cs []string
		for c, x := range die {
			if x == d {
				cs = append(cs, cellName(c))
			}
		}
		as = append(as, edakit.Sprintf("die %d holds %s", d, strings.Join(cs, ", ")))
	}
	return "Nets: " + strings.Join(ns, "; ") + ". Assignment: " + strings.Join(as, "; ") + "."
}

func cut(nets [][]int, die []int) int {
	c := 0
	for _, net := range nets {
		for _, x := range net[1:] {
			if die[x] != die[net[0]] {
				c++
				break
			}
		}
	}
	return c
}

func assign(n, dies int) []int {
	for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
		die := make([]int, n)
		used := map[int]bool{}
		for i := range die {
			die[i] = rand.Intn(dies)
			used[die[i]] = true
		}
		if len(used) == dies {
			return die
		}
	}
	die := make([]int, n)
	for i := range die {
		die[i] = i % dies
	}
	return die
}

// partitioning: cut nets, vertical connections in a stack, move gain.
func partitioning(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 6, 12)
	m := edakit.Scaled(ctx, 5, 10)
	nets := netlist(n, m)
	switch kpOf(ctx) {
	case 1:
		die := assign(n, 2)
		ans := edakit.IntStr(cut(nets, die))
		q := "A design is split across two dies. " + describe(nets, die, 2) + " Each net with cells on both dies needs a vertical connection. How many nets are cut?"
		return num(q, ans, edakit.Sprintf("A net is cut when its cells are not all on one die; counting them gives %s.", ans), "cut")
	case 2:
		die := assign(n, 3)
		total := 0
		for _, net := range nets {
			lo, hi := 3, -1
			for _, c := range net {
				if die[c] < lo {
					lo = die[c]
				}
				if die[c] > hi {
					hi = die[c]
				}
			}
			total += hi - lo
		}
		ans := edakit.IntStr(total)
		q := "Three dies are stacked as die 0 at the bottom, die 1 in the middle and die 2 on top. " + describe(nets, die, 3) + " A net needs one vertical connection at each die interface it crosses between its lowest and highest die. How many vertical connections are needed in total?"
		return num(q, ans, edakit.Sprintf("Sum over nets of (highest die minus lowest die) gives %s.", ans), "vertical")
	default:
		die := assign(n, 2)
		c := rand.Intn(n)
		before := cut(nets, die)
		moved := append([]int(nil), die...)
		moved[c] = 1 - moved[c]
		after := cut(nets, moved)
		ans := edakit.IntStr(before - after)
		q := "A design is split across two dies. " + describe(nets, die, 2) + edakit.Sprintf(" What is the gain (cut before minus cut after) of moving cell %s to the other die? The answer may be negative.", cellName(c))
		return num(q, ans, edakit.Sprintf("The cut is %d before and %d after moving %s, so the gain is %s.", before, after, cellName(c), ans), "gain")
	}
}

// interDieTiming: setup slack across dies, corner pessimism, hold slack.
func interDieTiming(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		T := float64(50 * edakit.Int(10, 30))
		La := float64(10 * edakit.Int(15, 40))
		Lb := float64(10 * edakit.Int(15, 40))
		dA := float64(10 * edakit.Int(5, 50))
		dX := float64(5 * edakit.Int(6, 40))
		dB := float64(10 * edakit.Int(5, 50))
		su := float64(5 * edakit.Int(2, 10))
		s := T + Lb - su - (La + dA + dX + dB)
		ans := f(s, 0)
		q := edakit.Sprintf("A path launches on die A and captures on die B. Clock period %.0f ps; launch clock latency on die A %.0f ps; capture clock latency on die B %.0f ps. Data delay: %.0f ps of logic on die A, %.0f ps for the die crossing, %.0f ps of logic on die B. Setup time %.0f ps. What is the setup slack, in ps?", T, La, Lb, dA, dX, dB, su)
		return num(q, ans, edakit.Sprintf("Required $%.0f + %.0f - %.0f = %.0f$; arrival $%.0f + %.0f + %.0f + %.0f = %.0f$; slack %s ps.", T, Lb, su, T+Lb-su, La, dA, dX, dB, La+dA+dX+dB, ans), "slack_ps")
	case 2:
		dB := float64(10 * edakit.Int(5, 50))
		ks := step(1.08, 1.25, 0.01)
		kb := step(0.85, 1.05, 0.01)
		ans := f((ks-kb)*dB, 1)
		q := edakit.Sprintf("A cross-die path has %.0f ps of nominal logic on the capture die. A single-corner analysis puts both dies at the slow corner, scaling delays by %s. The capture die is actually from a wafer at a corner that scales delays by %s. How many picoseconds of pessimism does the single-corner analysis add on that die's logic? Round to 1 decimal place.", dB, f(ks, 2), f(kb, 2))
		return num(q, ans, edakit.Sprintf("$(%s - %s) \\times %.0f = %s$ ps.", f(ks, 2), f(kb, 2), dB, ans), "pessimism_ps")
	default:
		La := float64(10 * edakit.Int(15, 40))
		Lb := float64(10 * edakit.Int(15, 40))
		dmin := float64(5 * edakit.Int(6, 40))
		th := float64(5 * edakit.Int(1, 8))
		s := La + dmin - Lb - th
		ans := f(s, 0)
		q := edakit.Sprintf("For a hold check on a cross-die path, the launch clock reaches die A's flop at %.0f ps and the capture clock reaches die B's flop at %.0f ps. The minimum data delay, including the die crossing, is %.0f ps, and the hold time is %.0f ps. What is the hold slack, in ps?", La, Lb, dmin, th)
		return num(q, ans, edakit.Sprintf("$(%.0f + %.0f) - (%.0f + %.0f) = %s$ ps; a negative value is a violation.", La, dmin, Lb, th, ans), "hold_ps")
	}
}

// thermalFloorplan: stack peak, die-order gain, lateral spreading gain.
func thermalFloorplan(ctx generator.GeneratorContext) generator.Problem {
	Rs := step(0.1, 0.6, 0.05)
	Rb := step(0.1, 0.6, 0.05)
	Ta := float64(edakit.Int(25, 50))
	Pa := float64(5 * edakit.Int(1, 30))
	Pb := float64(5 * edakit.Int(1, 30))
	if Pb == Pa {
		Pb = Pa + 5
	}
	stack := edakit.Sprintf("Two dies are stacked under a heat sink. The die next to the sink connects to ambient (%.0f °C) through %s K/W, and the far die connects to the near die through %s K/W. All heat leaves through the sink.", Ta, f(Rs, 2), f(Rb, 2))
	switch kpOf(ctx) {
	case 1:
		Tn := Ta + (Pa+Pb)*Rs
		Tf := Tn + Pb*Rb
		ans := f(Tf, 1)
		q := stack + edakit.Sprintf(" The near die dissipates %.0f W and the far die %.0f W. What is the peak temperature in the stack, in °C? Round to 1 decimal place.", Pa, Pb)
		return num(q, ans, edakit.Sprintf("Near die $%.0f + %.0f \\times %s = %s$ °C; far die adds $%.0f \\times %s$ to reach %s °C.", Ta, Pa+Pb, f(Rs, 2), f(Tn, 2), Pb, f(Rb, 2), ans), "t_peak")
	case 2:
		ans := f(math.Abs(Pa-Pb)*Rb, 1)
		q := stack + edakit.Sprintf(" One die dissipates %.0f W and the other %.0f W, and either may go next to the sink. How many kelvin lower is the peak temperature with the better order than with the worse? Round to 1 decimal place.", Pa, Pb)
		return num(q, ans, edakit.Sprintf("The peak is $T_a + (P_1 + P_2)R_s + P_{far}R_b$; only the far die's power changes, so the gain is $|%.0f - %.0f| \\times %s = %s$ K.", Pa, Pb, f(Rb, 2), ans), "gain_k")
	default:
		P1 := float64(edakit.Int(3, 25))
		P2 := float64(edakit.Int(3, 25))
		Rii := step(0.5, 4, 0.1)
		Ra := edakit.Round(Rii*step(0.2, 0.6, 0.05), 2)
		Rf := edakit.Round(Ra*step(0.1, 0.5, 0.05), 3)
		pk := func(rm float64) float64 { return math.Max(Rii*P1+rm*P2, Rii*P2+rm*P1) }
		ans := f(pk(Ra)-pk(Rf), 2)
		q := edakit.Sprintf("Two hot blocks dissipate %.0f W and %.0f W. Each block's temperature rise is $R_{ii}P_i + R_{ij}P_j$ with self resistance $R_{ii} = %s$ K/W. The mutual resistance is %s K/W when the blocks are adjacent and %s K/W when they are placed apart. By how many kelvin does separating them lower the peak rise? Round to 2 decimal places.", P1, P2, f(Rii, 1), f(Ra, 2), f(Rf, 3))
		return num(q, ans, edakit.Sprintf("Adjacent peak rise %s K, separated %s K, gain %s K.", f(pk(Ra), 3), f(pk(Rf), 3), ans), "gain_k")
	}
}
