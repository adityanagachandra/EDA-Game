package synth

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// fsm is a completely specified machine over one input bit. For a Moore machine out[s][0]
// equals out[s][1] and is the state's output.
type fsm struct {
	n     int
	next  [][2]int
	out   [][2]int
	moore bool
}

// refine runs partition refinement and returns the block of each state per round.
func (m *fsm) refine() (rounds [][]int) {
	cls := make([]int, m.n)
	key := map[[2]int]int{}
	for s := 0; s < m.n; s++ {
		k := m.out[s]
		if _, ok := key[k]; !ok {
			key[k] = len(key)
		}
		cls[s] = key[k]
	}
	rounds = append(rounds, cls)
	for {
		nk := map[[3]int]int{}
		nc := make([]int, m.n)
		for s := 0; s < m.n; s++ {
			k := [3]int{cls[s], cls[m.next[s][0]], cls[m.next[s][1]]}
			if _, ok := nk[k]; !ok {
				nk[k] = len(nk)
			}
			nc[s] = nk[k]
		}
		if len(nk) == count(cls) {
			return rounds
		}
		cls = nc
		rounds = append(rounds, cls)
	}
}

func count(cls []int) int {
	m := map[int]bool{}
	for _, c := range cls {
		m[c] = true
	}
	return len(m)
}

func (m *fsm) desc() string {
	var parts []string
	for s := 0; s < m.n; s++ {
		if m.moore {
			parts = append(parts, edakit.Sprintf("S%d (output %d): on 0 go to S%d, on 1 go to S%d", s, m.out[s][0], m.next[s][0], m.next[s][1]))
		} else {
			parts = append(parts, edakit.Sprintf("S%d: 0/%d to S%d, 1/%d to S%d", s, m.out[s][0], m.next[s][0], m.out[s][1], m.next[s][1]))
		}
	}
	return strings.Join(parts, "; ")
}

// randFSM builds a minimal machine of k states and pads it with n-k redundant copies.
func randFSM(k, n int, moore bool) *fsm {
	for {
		base := &fsm{n: k, moore: moore, next: make([][2]int, k), out: make([][2]int, k)}
		for s := 0; s < k; s++ {
			base.next[s] = [2]int{edakit.Int(0, k-1), edakit.Int(0, k-1)}
			if moore {
				o := edakit.Int(0, 1)
				base.out[s] = [2]int{o, o}
			} else {
				base.out[s] = [2]int{edakit.Int(0, 1), edakit.Int(0, 1)}
			}
		}
		r := base.refine()
		if count(r[len(r)-1]) != k {
			continue
		}
		// orig[i] is the base state that state i copies.
		orig := make([]int, n)
		for i := 0; i < n; i++ {
			if i < k {
				orig[i] = i
			} else {
				orig[i] = edakit.Int(0, k-1)
			}
		}
		perm := edakit.Shuffle(seq(n))
		members := make([][]int, k)
		for i := 0; i < n; i++ {
			members[orig[i]] = append(members[orig[i]], perm[i])
		}
		m := &fsm{n: n, moore: moore, next: make([][2]int, n), out: make([][2]int, n)}
		for i := 0; i < n; i++ {
			s := perm[i]
			b := orig[i]
			m.out[s] = base.out[b]
			for x := 0; x < 2; x++ {
				m.next[s][x] = edakit.Pick(members[base.next[b][x]])
			}
		}
		return m
	}
}

// stateMinimization: initial partition, minimal state count, refinement rounds.
func stateMinimization(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 8)
	k := edakit.Int(2, n-1)
	moore := ctx.Difficulty < 0.6
	m := randFSM(k, n, moore)
	r := m.refine()
	kind := "Moore machine"
	if !moore {
		kind = "Mealy machine (entries are input/output to next state)"
	}
	stem := "A completely specified " + kind + " has states " + m.desc() + "."
	switch kp(ctx) {
	case 1:
		p0 := count(r[0])
		return edakit.Numeric(stem+" Partition refinement starts from the partition by outputs. How many blocks does that initial partition $P_0$ have?",
			edakit.IntStr(p0), edakit.Sprintf("States with the same output response to every input share a block of $P_0$. There are %d blocks.", p0),
			map[string]string{"p0": edakit.IntStr(p0)})
	case 2:
		f := count(r[len(r)-1])
		return edakit.Numeric(stem+" How many states does the equivalent minimal machine have?",
			edakit.IntStr(f), edakit.Sprintf("Refine blocks by the blocks of the next states until nothing splits. The stable partition has %d blocks, one state each in the minimal machine.", f),
			map[string]string{"states": edakit.IntStr(f)})
	}
	splits := len(r) - 1
	if splits == 0 && edakit.Int(0, 2) > 0 {
		return stateMinimization(ctx)
	}
	return edakit.Numeric(stem+" Starting from $P_0$ (blocks by output), each refinement pass splits blocks whose states go to different blocks on some input. How many passes split at least one block before the partition is stable?",
		edakit.IntStr(splits), edakit.Sprintf("The block counts per pass are %s. %d passes changed the partition; the next one changes nothing.", roundCounts(r), splits),
		map[string]string{"passes": edakit.IntStr(splits)})
}

func roundCounts(r [][]int) string {
	var xs []int
	for _, c := range r {
		xs = append(xs, count(c))
	}
	return intsCSV(xs)
}

// distinguish returns the length of the shortest input sequence whose output sequences differ
// from states p and q of a Mealy machine, or 0 when they are equivalent.
func (m *fsm) distinguish(p, q int) int {
	type pr struct{ a, b int }
	dist := map[pr]int{{p, q}: 0}
	queue := []pr{{p, q}}
	for len(queue) > 0 {
		c := queue[0]
		queue = queue[1:]
		for x := 0; x < 2; x++ {
			if m.out[c.a][x] != m.out[c.b][x] {
				return dist[c] + 1
			}
		}
		for x := 0; x < 2; x++ {
			nx := pr{m.next[c.a][x], m.next[c.b][x]}
			if _, ok := dist[nx]; !ok {
				dist[nx] = dist[c] + 1
				queue = append(queue, nx)
			}
		}
	}
	return 0
}

// fsmEquivalence: one-step distinguishable pairs, shortest distinguishing sequences.
func fsmEquivalence(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 6)
	m := randFSM(edakit.Int(n-1, n), n, false)
	stem := "A Mealy machine (entries are input/output to next state) has states " + m.desc() + "."
	switch kp(ctx) {
	case 1:
		c := 0
		for p := 0; p < n; p++ {
			for q := p + 1; q < n; q++ {
				if m.out[p] != m.out[q] {
					c++
				}
			}
		}
		return edakit.Numeric(stem+" How many unordered pairs of states are distinguished by a single input symbol (their outputs differ for input 0 or for input 1)?",
			edakit.IntStr(c), edakit.Sprintf("Compare the output pairs of every two states: %d pairs differ immediately. These are the pairs split by $P_0$.", c),
			map[string]string{"pairs": edakit.IntStr(c)})
	case 2:
		for {
			p, q := edakit.Int(0, n-1), edakit.Int(0, n-1)
			if p == q {
				continue
			}
			if p > q {
				p, q = q, p
			}
			d := m.distinguish(p, q)
			if d == 0 {
				continue
			}
			return edakit.Numeric(stem+edakit.Sprintf(" What is the length of the shortest input sequence that makes S%d and S%d produce different output sequences?", p, q),
				edakit.IntStr(d), edakit.Sprintf("Search breadth-first over pairs of states reached by the same input sequence; the first pair with differing outputs appears after %d inputs.", d),
				map[string]string{"length": edakit.IntStr(d)})
		}
	}
	best := 0
	for p := 0; p < n; p++ {
		for q := p + 1; q < n; q++ {
			if d := m.distinguish(p, q); d > best {
				best = d
			}
		}
	}
	return edakit.Numeric(stem+" Over all pairs of inequivalent states, what is the longest shortest-distinguishing-sequence length? (This is the number of refinement steps needed, counting $P_0$ as step 1.)",
		edakit.IntStr(best), edakit.Sprintf("States separated first in partition $P_j$ need a sequence of length j + 1. The hardest pair needs %d inputs.", best),
		map[string]string{"max_length": edakit.IntStr(best)})
}

// fsmEncoding: flip-flop counts and weighted-adjacency cost of a state assignment.
func fsmEncoding(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 5, 40)
		b := ceilLog2(n)
		return edakit.Numeric(edakit.Sprintf("A controller has %d states. How many more flip-flops does a one-hot encoding use than a minimum-length binary encoding?", n),
			edakit.IntStr(n-b), edakit.Sprintf("One-hot uses one flip-flop per state, %d. Binary uses $\\lceil \\log_2 %d \\rceil = %d$. Difference %d.", n, n, b, n-b),
			map[string]string{"onehot": edakit.IntStr(n), "binary": edakit.IntStr(b)})
	case 2:
		n := edakit.Int(4, 6)
		bitsN := ceilLog2(n)
		codes := edakit.Shuffle(seq(1 << bitsN))[:n]
		cost, desc, wdesc := 0, []string{}, []string{}
		for i := 0; i < n; i++ {
			desc = append(desc, edakit.Sprintf("S%d = %0*b", i, bitsN, codes[i]))
		}
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if edakit.Int(0, 2) == 0 {
					continue
				}
				w := edakit.Int(1, 4)
				cost += w * popcount(uint64(codes[i]^codes[j]))
				wdesc = append(wdesc, edakit.Sprintf("w(S%d,S%d) = %d", i, j, w))
			}
		}
		if len(wdesc) == 0 {
			w := 2
			cost += w * popcount(uint64(codes[0]^codes[1]))
			wdesc = append(wdesc, "w(S0,S1) = 2")
		}
		return edakit.Numeric("An encoder wants state pairs with high affinity weight to get adjacent codes. Weights: "+strings.Join(wdesc, ", ")+" (all others 0). The assignment is "+strings.Join(desc, ", ")+". What is the cost $\\sum w_{ij} \\cdot \\mathrm{dist}(code_i, code_j)$, with Hamming distance?",
			edakit.IntStr(cost), edakit.Sprintf("Multiply each weight by the Hamming distance of the two codes and sum: %d.", cost),
			map[string]string{"cost": edakit.IntStr(cost)})
	}
	// Four states into two bits: only which pairs end up diagonal matters.
	w := map[[2]int]int{}
	var wdesc []string
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			w[[2]int{i, j}] = edakit.Int(0, 5)
			wdesc = append(wdesc, edakit.Sprintf("w(S%d,S%d) = %d", i, j, w[[2]int{i, j}]))
		}
	}
	best := 1 << 30
	for _, p := range perms4() {
		c := 0
		for i := 0; i < 4; i++ {
			for j := i + 1; j < 4; j++ {
				c += w[[2]int{i, j}] * popcount(uint64(p[i]^p[j]))
			}
		}
		if c < best {
			best = c
		}
	}
	return edakit.Numeric("Four states S0 to S3 get the 2-bit codes 00, 01, 10, 11 in some order. Affinity weights: "+strings.Join(wdesc, ", ")+". What is the minimum over all assignments of $\\sum w_{ij} \\cdot \\mathrm{dist}(code_i, code_j)$ with Hamming distance?",
		edakit.IntStr(best), edakit.Sprintf("Every pair is at distance 1 except the two diagonal pairs at distance 2, so the cost is the total weight plus the weight of the diagonal pairs. Choose the diagonal split with the least weight: minimum %d.", best),
		map[string]string{"min_cost": edakit.IntStr(best)})
}

func perms4() [][]int {
	var out [][]int
	var rec func(cur []int, used int)
	rec = func(cur []int, used int) {
		if len(cur) == 4 {
			out = append(out, append([]int(nil), cur...))
			return
		}
		for i := 0; i < 4; i++ {
			if used>>i&1 == 0 {
				rec(append(cur, i), used|1<<i)
			}
		}
	}
	rec(nil, 0)
	return out
}
