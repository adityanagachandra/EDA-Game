package hls

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// dfgProblem draws a straight-line block and asks for a structural count of its DFG.
func dfgProblem(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 10)
	g := newDFG(n, 0.4)
	code := g.text()
	switch whichKP(ctx) {
	case 1:
		e := g.edges()
		q := "A straight-line block is written in three-address form: " + code +
			". Letters are primary inputs. Build its data-flow graph with one node per operation. How many data edges run from one operation to another (edges from primary inputs do not count)?"
		return num(q, edakit.IntStr(e),
			edakit.Sprintf("Each operand that names an earlier result n_k is one edge from n_k to the reading operation. Counting those operands gives %d edges.", e),
			map[string]string{"edges": edakit.IntStr(e)})
	case 2:
		h := g.heights(g.lat(1, 1))
		d := maxInt(h...)
		q := "A straight-line block is written in three-address form: " + code +
			". Letters are primary inputs. How many operations lie on the longest chain of dependent operations in its data-flow graph?"
		return num(q, edakit.IntStr(d),
			edakit.Sprintf("Walking back from each operation to the inputs, the longest dependence chain has %d operations. With unit-latency operations this is also the minimum number of control steps.", d),
			map[string]string{"depth": edakit.IntStr(d)})
	default:
		su := g.succs()
		var sinks []string
		for i := range g.ops {
			if len(su[i]) == 0 {
				sinks = append(sinks, opName(i))
			}
		}
		q := "A straight-line block is written in three-address form: " + code +
			". Letters are primary inputs. In its data-flow graph, how many operations are sinks, whose result no later operation reads and which must therefore leave the block as outputs?"
		return num(q, edakit.IntStr(len(sinks)),
			edakit.Sprintf("The operations whose names never appear as an operand are %s, so there are %d sinks.", strings.Join(sinks, ", "), len(sinks)),
			map[string]string{"sinks": edakit.IntStr(len(sinks))})
	}
}

// opLatency turns combinational delays into cycle counts at a clock period.
func opLatency(ctx generator.GeneratorContext) generator.Problem {
	periods := []int{20, 25, 30, 40, 50} // tenths of a ns
	tt := edakit.Pick(periods)
	tenths := func(x int) string { return edakit.Num(float64(x)/10, 1) }
	switch whichKP(ctx) {
	case 1:
		d := edakit.Int(8, 140)
		c := ceilDiv(d, tt)
		q := edakit.Sprintf("A combinational multiplier has a delay of %s ns. The design clock period is %s ns and register overhead is ignored. How many clock cycles of latency must the scheduler give one multiplication if its result must be registered before use?", tenths(d), tenths(tt))
		return num(q, edakit.IntStr(c),
			edakit.Sprintf("The operation needs the ceiling of %s / %s, which is %d cycles. A partial cycle still costs a whole cycle.", tenths(d), tenths(tt), c),
			map[string]string{"cycles": edakit.IntStr(c)})
	case 2:
		m := edakit.Scaled(ctx, 3, 6)
		total := 0
		var ds []string
		for i := 0; i < m; i++ {
			d := edakit.Int(6, 110)
			total += ceilDiv(d, tt)
			ds = append(ds, tenths(d))
		}
		q := edakit.Sprintf("A chain of %d dependent operations has delays %s ns. The clock period is %s ns, register overhead is ignored, and every result is registered (no chaining). How many cycles does the chain take?", m, strings.Join(ds, ", "), tenths(tt))
		return num(q, edakit.IntStr(total),
			edakit.Sprintf("Each operation costs the ceiling of its delay over %s ns; the sum over the chain is %d cycles.", tenths(tt), total),
			map[string]string{"cycles": edakit.IntStr(total)})
	default:
		m := edakit.Scaled(ctx, 3, 5)
		var ds []int
		var dstr []string
		for i := 0; i < m; i++ {
			d := edakit.Int(8, 90)
			ds = append(ds, d)
			dstr = append(dstr, tenths(d))
		}
		t1 := edakit.Pick([]int{20, 25, 30})
		t2 := edakit.Pick([]int{40, 45, 50, 60})
		cyc := func(t int) int {
			s := 0
			for _, d := range ds {
				s += ceilDiv(d, t)
			}
			return s
		}
		c1, c2 := cyc(t1), cyc(t2)
		x1, x2 := c1*t1, c2*t2
		best := minInt(x1, x2)
		q := edakit.Sprintf("A chain of dependent operations has delays %s ns, and every result is registered (no chaining, register overhead ignored). Clock A has a period of %s ns and clock B a period of %s ns. What is the smaller of the two total execution times, in ns? Round to 1 decimal place.",
			strings.Join(dstr, ", "), tenths(t1), tenths(t2))
		return num(q, tenths(best),
			edakit.Sprintf("Clock A needs %d cycles, %s ns; clock B needs %d cycles, %s ns. Quantizing each delay up to whole cycles wastes time differently at each period, so the smaller total is %s ns.",
				c1, tenths(x1), c2, tenths(x2), tenths(best)),
			map[string]string{"cycles_a": edakit.IntStr(c1), "cycles_b": edakit.IntStr(c2)})
	}
}
