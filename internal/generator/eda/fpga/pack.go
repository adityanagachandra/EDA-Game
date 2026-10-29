package fpga

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type lutFF struct {
	nLUT, nFF int
	driver    []int // per FF: index of the driving LUT, or -1 for an input pin
	ffCount   []int // per LUT: flip-flops it drives
	other     []int // per LUT: sinks that are not flip-flops
}

func newLutFF(ctx generator.GeneratorContext) *lutFF {
	l := edakit.Scaled(ctx, 4, 8)
	f := edakit.Int(2, l)
	x := &lutFF{nLUT: l, nFF: f, driver: make([]int, f), ffCount: make([]int, l), other: make([]int, l)}
	for i := 0; i < f; i++ {
		if edakit.Int(0, 3) > 0 {
			x.driver[i] = edakit.Int(0, l-1)
			x.ffCount[x.driver[i]]++
		} else {
			x.driver[i] = -1
		}
	}
	for j := 0; j < l; j++ {
		if x.ffCount[j] == 0 {
			x.other[j] = edakit.Int(1, 3)
		} else {
			x.other[j] = edakit.Pick([]int{0, 0, 0, 1, 2})
		}
	}
	return x
}

func (x *lutFF) text() string {
	var p []string
	for j := 0; j < x.nLUT; j++ {
		var ffs []string
		for i, d := range x.driver {
			if d == j {
				ffs = append(ffs, edakit.Sprintf("F%d", i+1))
			}
		}
		s := edakit.Sprintf("L%d drives ", j+1)
		switch {
		case len(ffs) > 0 && x.other[j] > 0:
			s += strings.Join(ffs, " and ") + edakit.Sprintf(" plus %d other sink(s)", x.other[j])
		case len(ffs) > 0:
			s += "only " + strings.Join(ffs, " and ")
		default:
			s += edakit.Sprintf("%d non-flip-flop sink(s)", x.other[j])
		}
		p = append(p, s)
	}
	for i, d := range x.driver {
		if d < 0 {
			p = append(p, edakit.Sprintf("F%d is fed directly from an input pin", i+1))
		}
	}
	return strings.Join(p, "; ")
}

// pairs counts LUT-FF pairs. single: the BLE has one output, so the LUT may feed only its FF.
func (x *lutFF) pairs(single bool) int {
	n := 0
	for j := 0; j < x.nLUT; j++ {
		if single {
			if x.ffCount[j] == 1 && x.other[j] == 0 {
				n++
			}
		} else if x.ffCount[j] >= 1 {
			n++
		}
	}
	return n
}

func bleFormation(ctx generator.GeneratorContext) generator.Problem {
	x := newLutFF(ctx)
	head := edakit.Sprintf("A mapped netlist has %d LUTs (L1 to L%d) and %d flip-flops (F1 to F%d): %s.", x.nLUT, x.nLUT, x.nFF, x.nFF, x.text())
	single := "Each BLE holds one LUT and one flip-flop and has a single output, either the LUT output or the flip-flop output, so a LUT and a flip-flop share a BLE only when the LUT's sole sink is that flip-flop."
	switch whichKP(ctx) {
	case 1:
		p := x.pairs(true)
		b := x.nLUT + x.nFF - p
		return num(head+" "+single+" How many BLEs does the netlist need?", edakit.IntStr(b),
			edakit.Sprintf("LUTs that satisfy the pairing rule: %d. BLEs = %d + %d - %d = %d.", p, x.nLUT, x.nFF, p, b),
			map[string]string{"pairs": edakit.IntStr(p)})
	case 2:
		p := x.pairs(true)
		b := x.nLUT + x.nFF - p
		n := edakit.Int(2, 4)
		c := ceilDiv(b, n)
		return num(head+" "+single+edakit.Sprintf(" Clusters hold N = %d BLEs. Assuming input pins never limit packing, how many clusters are needed?", n), edakit.IntStr(c),
			edakit.Sprintf("%d pairs form, giving %d + %d - %d = %d BLEs, and ceil(%d / %d) = %d clusters.", p, x.nLUT, x.nFF, p, b, b, n, c),
			map[string]string{"bles": edakit.IntStr(b)})
	default:
		p := x.pairs(false)
		b := x.nLUT + x.nFF - p
		return num(head+" Now each BLE exposes both the LUT output and the flip-flop output, so a LUT can share a BLE with one flip-flop it drives even when the LUT has other sinks. How many BLEs does the netlist need?", edakit.IntStr(b),
			edakit.Sprintf("Every LUT that drives at least one flip-flop can absorb one of them: %d pairs, so BLEs = %d + %d - %d = %d. The second output is why modern logic elements expose both signals.", p, x.nLUT, x.nFF, p, b),
			map[string]string{"pairs": edakit.IntStr(p)})
	}
}

type ble struct {
	in      []string
	out     string
	outside bool
}

func newCluster(nb int) []ble {
	var bs []ble
	ext := []string{"e1", "e2", "e3", "e4", "e5", "e6", "e7", "e8", "e9"}
	for i := 0; i < nb; i++ {
		seen := map[string]bool{}
		var in []string
		for len(in) < 4 {
			var s string
			if i > 0 && edakit.Int(0, 9) < 3 {
				s = bs[edakit.Int(0, i-1)].out
			} else {
				s = edakit.Pick(ext)
			}
			if !seen[s] {
				seen[s] = true
				in = append(in, s)
			}
		}
		sortStrings(in)
		bs = append(bs, ble{in: in, out: edakit.Sprintf("o%d", i+1), outside: edakit.Int(0, 1) == 0})
	}
	bs[nb-1].outside = true
	return bs
}

func clusterText(bs []ble) string {
	var p []string
	for i, b := range bs {
		s := edakit.Sprintf("B%d has inputs %s and output %s", i+1, setText(b.in), b.out)
		if b.outside {
			s += " (also used outside the cluster)"
		} else {
			s += " (used only inside the cluster)"
		}
		p = append(p, s)
	}
	return strings.Join(p, "; ")
}

func extInputs(bs []ble, extra ...ble) []string {
	outs := map[string]bool{}
	all := append(append([]ble(nil), bs...), extra...)
	var ins [][]string
	for _, b := range all {
		outs[b.out] = true
		ins = append(ins, b.in)
	}
	var res []string
	for _, s := range union(ins...) {
		if !outs[s] {
			res = append(res, s)
		}
	}
	return res
}

func clustering(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		bs := newCluster(edakit.Scaled(ctx, 3, 5))
		ex := extInputs(bs)
		q := "A cluster of 4-LUT BLEs holds: " + clusterText(bs) + ". Signals produced inside the cluster reach other BLEs of the cluster through local feedback. How many distinct cluster input pins does the cluster need?"
		return num(q, edakit.IntStr(len(ex)),
			edakit.Sprintf("Take the union of all BLE inputs and remove the outputs generated inside the cluster. The remaining nets %s each need one cluster input: %d.", setText(ex), len(ex)),
			map[string]string{"inputs": edakit.IntStr(len(ex))})
	case 2:
		bs := newCluster(3)
		cur := extInputs(bs)
		lim := len(cur) + edakit.Int(0, 1)
		var cands []string
		ok := 0
		var why []string
		for c := 0; c < 4; c++ {
			sz := edakit.Int(2, 4)
			pool := []string{"e1", "e2", "e3", "e4", "e5", "e6", "e7", "e8", "e9", "e10", "e11", "e12", "o1", "o2", "o3"}
			in := pickDistinct(pool, sz)
			sortStrings(in)
			out := edakit.Pick([]string{"x", cur[edakit.Int(0, len(cur)-1)]})
			name := edakit.Sprintf("X%d", c+1)
			for _, s := range in {
				if s == out {
					out = "x"
				}
			}
			if out == "x" {
				out = edakit.Sprintf("x%d", c+1)
			}
			n := len(extInputs(bs, ble{in: in, out: out}))
			if n <= lim {
				ok++
			}
			cands = append(cands, edakit.Sprintf("%s with inputs %s and output %s", name, setText(in), out))
			why = append(why, edakit.Sprintf("%s → %d", name, n))
		}
		q := edakit.Sprintf("A partly filled cluster holds: %s. The cluster has room for more BLEs and I = %d input pins. Candidates, each considered alone (a candidate whose output is one of the e nets is that net's driver): %s. A net produced by a BLE inside the cluster needs no input pin. How many candidates can legally be added?", clusterText(bs), lim, strings.Join(cands, "; "))
		return num(q, edakit.IntStr(ok),
			edakit.Sprintf("The cluster now uses %d inputs %s. Recount the inputs with each candidate added, remembering that a candidate whose output is a net the cluster already imports frees that pin: %s. %d of them fit within I = %d.", len(cur), setText(cur), strings.Join(why, ", "), ok, lim),
			map[string]string{"current_inputs": edakit.IntStr(len(cur))})
	default:
		bs := newCluster(edakit.Scaled(ctx, 3, 5))
		used := map[string]bool{}
		for _, b := range bs {
			for _, s := range b.in {
				used[s] = true
			}
		}
		abs := 0
		var names []string
		for _, b := range bs {
			if !b.outside && used[b.out] {
				abs++
				names = append(names, b.out)
			}
		}
		q := "A cluster of 4-LUT BLEs holds: " + clusterText(bs) + ". A net is absorbed when its driver and all of its sinks lie inside the cluster, so it needs no inter-cluster routing. How many nets does this cluster absorb?"
		ex := "none"
		if abs > 0 {
			ex = strings.Join(names, ", ")
		}
		return num(q, edakit.IntStr(abs),
			edakit.Sprintf("Look for outputs used only inside the cluster that some BLE of the cluster actually reads: %s, %d in all. Absorbing nets is what makes clustering reduce the routing problem.", ex, abs),
			map[string]string{"absorbed": edakit.IntStr(abs)})
	}
}

func vpackAttraction(ctx generator.GeneratorContext) generator.Problem {
	for {
		pool := []string{"n1", "n2", "n3", "n4", "n5", "n6", "n7", "n8", "n9", "n10", "n11", "n12"}
		cn := pickDistinct(pool, edakit.Int(5, 7))
		sortStrings(cn)
		inC := map[string]bool{}
		for _, s := range cn {
			inC[s] = true
		}
		type cand struct {
			nets   []string
			crit   float64
			shared int
		}
		var cs []cand
		var desc []string
		for i := 0; i < 4; i++ {
			ns := pickDistinct(pool, edakit.Int(3, 5))
			sortStrings(ns)
			sh := 0
			for _, s := range ns {
				if inC[s] {
					sh++
				}
			}
			cr := float64(edakit.Int(0, 20)) / 20
			cs = append(cs, cand{ns, cr, sh})
			desc = append(desc, edakit.Sprintf("X%d touches %s with criticality %s", i+1, setText(ns), edakit.Num(cr, 2)))
		}
		lam := edakit.Pick([]float64{0.5, 0.75, 0.9})
		g := 5.0
		attr := func(c cand) float64 { return lam*c.crit + (1-lam)*float64(c.shared)/g }
		head := edakit.Sprintf("A cluster being filled touches the nets %s. Candidate BLEs: %s.", setText(cn), strings.Join(desc, "; "))
		switch whichKP(ctx) {
		case 1:
			x := edakit.Int(0, 3)
			q := head + edakit.Sprintf(" What is VPack's attraction for X%d, the number of nets it shares with the cluster?", x+1)
			return num(q, edakit.IntStr(cs[x].shared),
				edakit.Sprintf("Intersect X%d's nets with the cluster's: %d shared nets. VPack ignores criticality.", x+1, cs[x].shared),
				map[string]string{"shared": edakit.IntStr(cs[x].shared)})
		case 2:
			x := edakit.Int(0, 3)
			a := attr(cs[x])
			q := head + edakit.Sprintf(" T-VPack scores a candidate as λ·Crit + (1 - λ)·shared / G with λ = %s and G = 5, the most nets a 4-LUT BLE can touch. What is X%d's attraction? Round to 3 decimal places.", edakit.Num(lam, 2), x+1)
			return num(q, edakit.Num(a, 3),
				edakit.Sprintf("X%d shares %d nets, so attraction = %s × %s + %s × %d / 5 = %s.", x+1, cs[x].shared, edakit.Num(lam, 2), edakit.Num(cs[x].crit, 2), edakit.Num(1-lam, 2), cs[x].shared, edakit.Num(a, 3)),
				map[string]string{"shared": edakit.IntStr(cs[x].shared)})
		default:
			best, bi, second := -1.0, 0, -1.0
			var sc []string
			for i, c := range cs {
				a := attr(c)
				sc = append(sc, edakit.Sprintf("X%d: %s", i+1, edakit.Num(a, 3)))
				if a > best {
					second = best
					best, bi = a, i
				} else if a > second {
					second = a
				}
			}
			if math.Abs(best-second) < 0.005 {
				continue
			}
			q := head + edakit.Sprintf(" T-VPack scores a candidate as λ·Crit + (1 - λ)·shared / G with λ = %s and G = 5, and adds the highest-scoring one. Which candidate does it add? Answer with its number.", edakit.Num(lam, 2))
			return num(q, edakit.IntStr(bi+1),
				edakit.Sprintf("Scores: %s. X%d wins. With a large λ, criticality dominates, so T-VPack pulls critical connections inside the cluster where local wiring is fast.", strings.Join(sc, ", "), bi+1),
				map[string]string{"winner": edakit.IntStr(bi + 1)})
		}
	}
}
