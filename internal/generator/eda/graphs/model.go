package graphs

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// hypergraphGen draws a small netlist and asks for its pin count, a cell degree or a cell's
// neighbourhood size.
func hypergraphGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 9)
	h := randHyper(n, edakit.Scaled(ctx, 4, 7), 4)
	intro := "A netlist has cells c1 to c" + edakit.IntStr(n) + " and nets " + netText(h, false) + ". "
	switch kpOf(ctx) {
	case 1:
		pins := 0
		var ks []int
		for _, net := range h.nets {
			pins += len(net)
			ks = append(ks, len(net))
		}
		return edakit.Numeric(intro+"Seen as a hypergraph, how many pins (vertex-hyperedge incidences) does it have?",
			edakit.IntStr(pins),
			"Each net contributes one pin per cell it touches. The net sizes are "+intsText(ks)+", which sum to "+edakit.IntStr(pins)+".",
			map[string]string{"pins": edakit.IntStr(pins)})
	case 2:
		c := edakit.Int(0, n-1)
		deg := 0
		var on []string
		for j, net := range h.nets {
			for _, p := range net {
				if p == c {
					deg++
					on = append(on, edakit.Sprintf("n%d", j+1))
				}
			}
		}
		return edakit.Numeric(intro+"What is the degree of cell "+cellName(c)+" in the hypergraph, the number of nets incident to it?",
			edakit.IntStr(deg),
			cellName(c)+" appears on "+strings.Join(on, ", ")+", so its degree is "+edakit.IntStr(deg)+".",
			map[string]string{"degree": edakit.IntStr(deg)})
	default:
		c := edakit.Int(0, n-1)
		nb := map[int]bool{}
		for _, net := range h.nets {
			has := false
			for _, p := range net {
				if p == c {
					has = true
				}
			}
			if has {
				for _, p := range net {
					if p != c {
						nb[p] = true
					}
				}
			}
		}
		var names []string
		for i := 0; i < n; i++ {
			if nb[i] {
				names = append(names, cellName(i))
			}
		}
		return edakit.Numeric(intro+"How many distinct other cells share at least one net with "+cellName(c)+"?",
			edakit.IntStr(len(nb)),
			"Take the union of the nets on "+cellName(c)+" and drop "+cellName(c)+" itself: {"+strings.Join(names, ", ")+"}, which is "+edakit.IntStr(len(nb))+" cells. A cell on two nets with the same neighbour counts that neighbour once.",
			map[string]string{"neighbours": edakit.IntStr(len(nb))})
	}
}

// pinGraph draws nets each with one driver and sinks of higher index, so the cell graph is acyclic.
type drivenNet struct {
	drv   int
	sinks []int
}

func randDriven(n, m int) []drivenNet {
	drivers := edakit.Shuffle(seq(n - 1))[:m]
	sort.Ints(drivers)
	var out []drivenNet
	for _, d := range drivers {
		cands := edakit.Shuffle(seq(n - d - 1))
		k := edakit.Int(1, min(3, len(cands)))
		var s []int
		for _, c := range cands[:k] {
			s = append(s, d+1+c)
		}
		sort.Ints(s)
		out = append(out, drivenNet{d, s})
	}
	return out
}

func drivenText(nets []drivenNet) string {
	var parts []string
	for j, nt := range nets {
		var sn []string
		for _, s := range nt.sinks {
			sn = append(sn, cellName(s))
		}
		parts = append(parts, edakit.Sprintf("n%d driven by %s, sinks %s", j+1, cellName(nt.drv), strings.Join(sn, ", ")))
	}
	return strings.Join(parts, "; ")
}

func pinGraphGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 6, 9)
	m := edakit.Int(3, min(n-1, 3+int(ctx.Difficulty*3)))
	nets := randDriven(n, m)
	intro := "Each cell has one output pin. The nets are: " + drivenText(nets) + ". Build the directed pin graph with one arc from each net's driver to each of its sinks. "
	switch kpOf(ctx) {
	case 1:
		arcs := 0
		for _, nt := range nets {
			arcs += len(nt.sinks)
		}
		return edakit.Numeric(intro+"How many arcs does the graph have?", edakit.IntStr(arcs),
			"A net with k pins has one driver and k - 1 sinks, so it contributes k - 1 arcs. Summing the sink counts gives "+edakit.IntStr(arcs)+".",
			map[string]string{"arcs": edakit.IntStr(arcs)})
	case 2:
		var sinks []int
		for _, nt := range nets {
			sinks = append(sinks, nt.sinks...)
		}
		v := edakit.Pick(sinks)
		fanin := 0
		for _, nt := range nets {
			for _, s := range nt.sinks {
				if s == v {
					fanin++
				}
			}
		}
		return edakit.Numeric(intro+"What is the in-degree of "+cellName(v)+", the number of its input pins that are driven?", edakit.IntStr(fanin),
			"Count the nets that list "+cellName(v)+" as a sink: "+edakit.IntStr(fanin)+". Each such net is one arc into "+cellName(v)+".",
			map[string]string{"fanin": edakit.IntStr(fanin)})
	default:
		start := nets[0].drv
		reach := map[int]bool{}
		stack := []int{start}
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, nt := range nets {
				if nt.drv == u {
					for _, s := range nt.sinks {
						if !reach[s] {
							reach[s] = true
							stack = append(stack, s)
						}
					}
				}
			}
		}
		var names []string
		for i := 0; i < n; i++ {
			if reach[i] {
				names = append(names, cellName(i))
			}
		}
		return edakit.Numeric(intro+"How many cells lie in the transitive fanout cone of "+cellName(start)+" (cells reachable from it, not counting itself)?",
			edakit.IntStr(len(reach)),
			"Follow arcs forward from "+cellName(start)+": the reachable set is {"+strings.Join(names, ", ")+"}, "+edakit.IntStr(len(reach))+" cells.",
			map[string]string{"cone": edakit.IntStr(len(reach))})
	}
}

func netModelsGen(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		m := edakit.Scaled(ctx, 3, 6)
		var ks []int
		total := 0
		for i := 0; i < m; i++ {
			k := edakit.Int(2, 8)
			ks = append(ks, k)
			total += k * (k - 1) / 2
		}
		return edakit.Numeric("A netlist has nets with "+intsText(ks)+" pins. Expanding every net as a clique, how many graph edges result (count parallel edges between the same pair separately)?",
			edakit.IntStr(total),
			"A k-pin clique has k(k - 1)/2 edges. Summing over the nets gives "+edakit.IntStr(total)+". Large nets dominate: the clique count grows quadratically.",
			map[string]string{"edges": edakit.IntStr(total)})
	case 2:
		k := edakit.Scaled(ctx, 6, 30)
		clique := k * (k - 1) / 2
		diff := clique - k
		return edakit.Numeric(edakit.Sprintf("A clock-enable net has %d pins. How many more edges does its clique model have than its star model (one new star vertex joined to every pin)?", k),
			edakit.IntStr(diff),
			edakit.Sprintf("Clique: %d·%d/2 = %d edges. Star: %d edges. The difference is %d.", k, k-1, clique, k, diff),
			map[string]string{"clique": edakit.IntStr(clique), "star": edakit.IntStr(k)})
	default:
		n := edakit.Scaled(ctx, 5, 8)
		for {
			h := randHyper(n, edakit.Scaled(ctx, 4, 7), 5)
			a, b := edakit.Int(0, n-1), edakit.Int(0, n-1)
			if a == b {
				continue
			}
			total := 0.0
			var terms []string
			for _, net := range h.nets {
				ha, hb := false, false
				for _, p := range net {
					ha = ha || p == a
					hb = hb || p == b
				}
				if ha && hb {
					total += 1.0 / float64(len(net)-1)
					terms = append(terms, edakit.Sprintf("1/%d", len(net)-1))
				}
			}
			if len(terms) == 0 {
				continue
			}
			ans := edakit.Num(total, 2)
			return edakit.Numeric("Nets: "+netText(h, false)+". Each net of k pins is expanded as a clique whose edges get weight $1/(k-1)$, and parallel edges are merged by adding weights. What is the weight of the merged edge between "+cellName(a)+" and "+cellName(b)+"? Round to 2 decimal places.",
				ans,
				"Only nets that contain both cells contribute, each 1/(k - 1): "+strings.Join(terms, " + ")+" = "+ans+".",
				map[string]string{"weight": ans})
		}
	}
}

func csrGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 9)
	h := randHyper(n, edakit.Scaled(ctx, 4, 7), 4)
	m := len(h.nets)
	pins := 0
	off := []int{0}
	for _, net := range h.nets {
		pins += len(net)
		off = append(off, pins)
	}
	intro := "Nets: " + netText(h, false) + ". The tool stores net-to-pin incidence in compressed form: an offset array netStart of length m + 1 and a pin array listing each net's cells in order, so net j's pins are pin[netStart[j]] up to pin[netStart[j+1] - 1], with nets indexed from 0 (n1 is index 0). "
	switch kpOf(ctx) {
	case 1:
		j := edakit.Int(1, m-1)
		return edakit.Numeric(intro+edakit.Sprintf("What is netStart[%d]?", j), edakit.IntStr(off[j]),
			edakit.Sprintf("netStart[%d] is the number of pins in nets 0 to %d, which is %d. The offsets are %s.", j, j-1, off[j], intsText(off)),
			map[string]string{"offset": edakit.IntStr(off[j])})
	case 2:
		c := edakit.Int(0, n-1)
		// Cell-to-net offsets: cellStart[c] = sum of degrees of cells before c.
		deg := make([]int, n)
		for _, net := range h.nets {
			for _, p := range net {
				deg[p]++
			}
		}
		s := 0
		for i := 0; i < c; i++ {
			s += deg[i]
		}
		return edakit.Numeric(intro+"The tool also builds the transpose, cell-to-net incidence, with an offset array cellStart indexed from 0 (c1 is index 0). What is cellStart["+edakit.IntStr(c)+"]?",
			edakit.IntStr(s),
			"cellStart[c] is the sum of the degrees of the cells before it. The degrees are "+intsText(deg)+", and the first "+edakit.IntStr(c)+" of them sum to "+edakit.IntStr(s)+".",
			map[string]string{"offset": edakit.IntStr(s)})
	default:
		bytes := 4 * ((m + 1) + pins + (n + 1) + pins)
		return edakit.Numeric(intro+edakit.Sprintf("Store both directions (net-to-pin and cell-to-net, each an offset array plus an index array) with 4-byte integers. There are %d cells. How many bytes in total?", n),
			edakit.IntStr(bytes),
			edakit.Sprintf("Net side: %d offsets + %d pins. Cell side: %d offsets + %d pins. Total %d integers times 4 bytes = %d bytes.", m+1, pins, n+1, pins, (m+1)+pins+(n+1)+pins, bytes),
			map[string]string{"pins": edakit.IntStr(pins), "bytes": edakit.IntStr(bytes)})
	}
}
