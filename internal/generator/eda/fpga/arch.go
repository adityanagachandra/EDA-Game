package fpga

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func lutBits(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		switch edakit.Int(0, 2) {
		case 0:
			k := edakit.Int(3, 7)
			n := edakit.Int(2, 400)
			q := edakit.Sprintf("A design maps to %d LUTs, each a %d-input LUT. How many SRAM configuration bits do those LUTs hold in total?", n, k)
			return num(q, edakit.IntStr(n*pow2(k)),
				edakit.Sprintf("A %d-LUT stores one bit per input combination, 2^%d = %d bits, so %d LUTs hold %d × %d = %d bits.", k, k, pow2(k), n, n, pow2(k), n*pow2(k)),
				map[string]string{"bits_per_lut": edakit.IntStr(pow2(k))})
		case 1:
			k := edakit.Int(2, 7)
			q := edakit.Sprintf("A %d-input LUT is built as a tree of 2:1 multiplexers whose leaves are the configuration bits and whose select lines are the LUT inputs. How many 2:1 multiplexers does the tree contain?", k)
			return num(q, edakit.IntStr(pow2(k)-1),
				edakit.Sprintf("The tree has 2^%d = %d leaves. A binary tree with L leaves has L - 1 internal nodes, so it needs %d multiplexers.", k, pow2(k), pow2(k)-1),
				map[string]string{"leaves": edakit.IntStr(pow2(k))})
		default:
			bits := make([]int, 8)
			var s []string
			for i := range bits {
				bits[i] = edakit.Int(0, 1)
				s = append(s, edakit.IntStr(bits[i]))
			}
			a, b, c := edakit.Int(0, 1), edakit.Int(0, 1), edakit.Int(0, 1)
			addr := 4*a + 2*b + c
			q := edakit.Sprintf("A 3-input LUT holds the configuration bits %s, listed from address 0 to address 7. Input a is the most significant select bit and c the least, so the address is 4a + 2b + c. Which bit does the LUT output when (a, b, c) = (%d, %d, %d)?", strings.Join(s, " "), a, b, c)
			return num(q, edakit.IntStr(bits[addr]),
				edakit.Sprintf("The inputs select address 4×%d + 2×%d + %d = %d, and the bit stored there is %d.", a, b, c, addr, bits[addr]),
				map[string]string{"address": edakit.IntStr(addr)})
		}
	case 2:
		if edakit.Int(0, 1) == 0 {
			k := edakit.Int(2, 5)
			f := uint64(1) << uint(pow2(k))
			q := edakit.Sprintf("How many distinct Boolean functions of %d inputs can a single %d-input LUT implement?", k, k)
			return num(q, edakit.Sprintf("%d", f),
				edakit.Sprintf("The LUT holds 2^%d = %d bits, and every one of the 2^%d assignments of those bits is a different truth table, so it implements 2^(2^%d) = %d functions.", k, pow2(k), pow2(k), k, f),
				map[string]string{"bits": edakit.IntStr(pow2(k))})
		}
		k := edakit.Int(4, 6)
		j := edakit.Int(2, k-1)
		q := edakit.Sprintf("A %d-input function is placed in a %d-input LUT, and the unused LUT inputs are tied to constants by the router. The tool must still write all %d configuration bits. How many of those bits are actually read for some input combination of the function?", j, k, pow2(k))
		return num(q, edakit.IntStr(pow2(j)),
			edakit.Sprintf("With the %d unused inputs held constant, only the 2^%d = %d addresses that agree with those constants are ever selected. Tools usually replicate the %d-bit truth table into the rest so that the choice of constant does not matter.", k-j, j, pow2(j), pow2(j)),
			map[string]string{"used": edakit.IntStr(pow2(j))})
	default:
		pool := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
		n := edakit.Scaled(ctx, 4, 6)
		type fn struct{ in []string }
		fs := make([]fn, n)
		var parts []string
		for i := range fs {
			sz := edakit.Int(2, 5)
			in := pickDistinct(pool[:edakit.Int(5, 8)], sz)
			sortStrings(in)
			fs[i] = fn{in}
			parts = append(parts, edakit.Sprintf("F%d = %s", i+1, setText(in)))
		}
		used := make([]bool, n)
		luts := 0
		var pairs []string
		for i := 0; i < n; i++ {
			if used[i] {
				continue
			}
			used[i] = true
			luts++
			for j := i + 1; j < n; j++ {
				if !used[j] && len(union(fs[i].in, fs[j].in)) <= 5 {
					used[j] = true
					pairs = append(pairs, edakit.Sprintf("F%d with F%d", i+1, j+1))
					break
				}
			}
		}
		ps := "no pair fits"
		if len(pairs) > 0 {
			ps = "the pairs are " + strings.Join(pairs, ", ")
		}
		q := "A fracturable 6-input LUT can serve either one function of up to 6 inputs or two functions whose combined set of distinct inputs has at most 5 signals. A packer walks the functions in order and pairs each unpaired function with the first later unpaired function that fits with it. Functions and their inputs: " + strings.Join(parts, ", ") + ". How many fracturable LUTs does it use?"
		return num(q, edakit.IntStr(luts),
			edakit.Sprintf("Check each candidate pair's input union against the 5-signal limit in the stated order: %s. That leaves %d physical LUTs for %d functions.", ps, luts, n),
			map[string]string{"luts": edakit.IntStr(luts), "pairs": edakit.IntStr(len(pairs))})
	}
}

func bleClb(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		k := edakit.Int(4, 6)
		n := edakit.Int(2, 10)
		b := n * (pow2(k) + 1)
		q := edakit.Sprintf("A cluster holds %d basic logic elements. Each BLE is a %d-input LUT, a flip-flop and a 2:1 output multiplexer that chooses the LUT output or the registered output. Counting only the LUT bits and the output-multiplexer select bits, how many configuration bits does the cluster hold?", n, k)
		return num(q, edakit.IntStr(b),
			edakit.Sprintf("Each BLE stores 2^%d = %d LUT bits plus 1 select bit, %d bits; %d BLEs hold %d bits.", k, pow2(k), pow2(k)+1, n, b),
			map[string]string{"per_ble": edakit.IntStr(pow2(k) + 1)})
	case 2:
		k := edakit.Pick([]int{4, 6})
		n := edakit.Int(2, 12)
		in := k * (n + 1) / 2
		if edakit.Int(0, 1) == 0 {
			q := edakit.Sprintf("Betz and Rose found that a cluster of N BLEs built from k-input LUTs needs only about I = (k/2)(N + 1) distinct cluster inputs to reach near-full logic utilization. How many cluster inputs does that rule give for k = %d and N = %d?", k, n)
			return num(q, edakit.IntStr(in),
				edakit.Sprintf("I = (%d/2)(%d + 1) = %d. A fully used cluster would need k·N = %d inputs, but BLEs share inputs and feed one another locally.", k, n, in, k*n),
				map[string]string{"inputs": edakit.IntStr(in)})
		}
		sw := n * k * (in + n)
		q := edakit.Sprintf("A cluster has N = %d BLEs of %d-input LUTs and I = %d cluster inputs. Its local interconnect is fully connected: every LUT input can select any cluster input or any BLE output of the cluster. How many programmable switches (one per possible source per LUT input) does that crossbar contain?", n, k, in)
		return num(q, edakit.IntStr(sw),
			edakit.Sprintf("There are N·k = %d LUT inputs, each choosing among I + N = %d sources, so %d × %d = %d switches.", n*k, in+n, n*k, in+n, sw),
			map[string]string{"lut_inputs": edakit.IntStr(n * k), "sources": edakit.IntStr(in + n)})
	default:
		n := edakit.Pick([]int{4, 8, 10})
		l := edakit.Scaled(ctx, 50, 5000)
		cl := ceilDiv(l, n)
		side := int(math.Ceil(math.Sqrt(float64(cl))))
		for side*side < cl {
			side++
		}
		for (side-1)*(side-1) >= cl {
			side--
		}
		q := edakit.Sprintf("A circuit maps to %d LUTs, each with its flip-flop absorbed into the same BLE. Clusters hold N = %d BLEs, and the tool sizes the smallest square array of clusters that holds the design, assuming every cluster can be filled. What is the side length of that array, in clusters?", l, n)
		return num(q, edakit.IntStr(side),
			edakit.Sprintf("The design needs at least ceil(%d / %d) = %d clusters. The smallest n with n² ≥ %d is %d.", l, n, cl, cl, side),
			map[string]string{"clusters": edakit.IntStr(cl)})
	}
}

func carryChain(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 8, 64)
		t0 := 10 * edakit.Int(10, 40)
		tc := edakit.Int(10, 40)
		tl := 10 * edakit.Int(10, 30)
		tr := 10 * edakit.Int(15, 50)
		if edakit.Int(0, 1) == 0 {
			d := t0 + n*tc
			q := edakit.Sprintf("A %d-bit adder uses the dedicated carry chain. Getting onto and off the chain costs %d ps in total, and each bit adds %d ps of carry propagation. What is the delay from the carry-in to the last sum bit, in ps?", n, t0, tc)
			return num(q, edakit.IntStr(d),
				edakit.Sprintf("Delay = %d + %d × %d = %d ps.", t0, n, tc, d),
				map[string]string{"delay": edakit.IntStr(d)})
		}
		d := n * (tl + tr)
		q := edakit.Sprintf("Without a carry chain, a ripple adder of %d bits puts each bit's carry in its own LUT, and every carry passes through one LUT (%d ps) and one general-routing hop (%d ps). What is the carry-in to last-bit delay, in ps?", n, tl, tr)
		return num(q, edakit.IntStr(d),
			edakit.Sprintf("Each bit costs %d + %d = %d ps, so %d bits cost %d ps. The routing hop is what the dedicated chain removes.", tl, tr, tl+tr, n, d),
			map[string]string{"per_bit": edakit.IntStr(tl + tr)})
	case 2:
		n := edakit.Scaled(ctx, 6, 96)
		b := edakit.Pick([]int{2, 4, 8})
		if edakit.Int(0, 1) == 0 {
			q := edakit.Sprintf("Each logic block's carry logic handles %d bits, and the chain continues into the block directly above. How many vertically adjacent logic blocks does a %d-bit adder occupy?", b, n)
			return num(q, edakit.IntStr(ceilDiv(n, b)),
				edakit.Sprintf("ceil(%d / %d) = %d blocks, which must be stacked in one column for the chain to connect, a placement constraint the placer must honor.", n, b, ceilDiv(n, b)),
				map[string]string{"blocks": edakit.IntStr(ceilDiv(n, b))})
		}
		tc := edakit.Int(10, 40)
		tl := 10 * edakit.Int(10, 30)
		tr := 10 * edakit.Int(15, 50)
		r := float64(tl+tr) / float64(tc)
		q := edakit.Sprintf("A carry costs %d ps per bit on the dedicated chain, and %d ps (LUT) plus %d ps (routing hop) per bit when rippled through soft logic. Ignoring entry and exit delay, how many times faster is the dedicated chain per bit? Round to 2 decimal places.", tc, tl, tr)
		return num(q, edakit.Num(r, 2),
			edakit.Sprintf("(%d + %d) / %d = %s.", tl, tr, tc, edakit.Num(r, 2)),
			map[string]string{"ratio": edakit.Num(r, 2)})
	default:
		t0 := 10 * edakit.Int(5, 20)
		tc := edakit.Int(15, 40)
		tl := 10 * edakit.Int(25, 45)
		chain := func(n int) int { return t0 + n*tc }
		pre := func(n int) int { return (bitsFor(n) + 1) * tl }
		ans := -1
		for n := 2; n <= 8192; n++ {
			if pre(n) < chain(n) {
				ans = n
				break
			}
		}
		if ans < 0 {
			tc += 20
			for n := 2; n <= 8192; n++ {
				if t0+n*tc > pre(n) {
					ans = n
					break
				}
			}
		}
		q := edakit.Sprintf("Model an n-bit addition two ways. On the dedicated carry chain it takes %d + %d·n ps. As a parallel-prefix adder in soft logic it takes (ceil(log2 n) + 1) levels of LUT plus routing at %d ps per level. What is the smallest n for which the soft prefix adder is strictly faster?", t0, tc, tl)
		return num(q, edakit.IntStr(ans),
			edakit.Sprintf("At n = %d the chain takes %d ps and the prefix adder %d ps; at n = %d they were %d and %d ps. Because the chain is so fast per bit, soft prefix adders only win for very wide additions.", ans, t0+ans*tc, pre(ans), ans-1, t0+(ans-1)*tc, pre(ans-1)),
			map[string]string{"crossover": edakit.IntStr(ans)})
	}
}

// bitsFor returns ceil(log2 n) for n >= 1.
func bitsFor(n int) int {
	b := 0
	for (1 << uint(b)) < n {
		b++
	}
	return b
}

func hardBlocks(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		type cfg struct{ d, w int }
		cfgs := []cfg{{32768, 1}, {16384, 2}, {8192, 4}, {4096, 9}, {2048, 18}, {1024, 36}, {512, 72}}
		d := edakit.Scaled(ctx, 200, 40000)
		w := edakit.Int(1, 80)
		best, bi := 1<<30, 0
		var lines []string
		for i, c := range cfgs {
			n := ceilDiv(d, c.d) * ceilDiv(w, c.w)
			lines = append(lines, edakit.Sprintf("%d×%d needs %d", c.d, c.w, n))
			if n < best {
				best, bi = n, i
			}
		}
		q := edakit.Sprintf("A block RAM tile can be configured as 32768×1, 16384×2, 8192×4, 4096×9, 2048×18, 1024×36 or 512×72 (depth × width). A memory of %d words of %d bits is built from identically configured tiles, cascading tiles in depth and placing them side by side in width. What is the fewest tiles it needs?", d, w)
		return num(q, edakit.IntStr(best),
			edakit.Sprintf("For each shape, tiles = ceil(depth / tile depth) × ceil(width / tile width): %s. The best is %d×%d with %d tiles.", strings.Join(lines, "; "), cfgs[bi].d, cfgs[bi].w, best),
			map[string]string{"tiles": edakit.IntStr(best)})
	case 2:
		p := edakit.Pick([]int{18, 24, 27})
		qq := edakit.Pick([]int{17, 18})
		m := edakit.Int(8, 64)
		n := edakit.Int(8, 64)
		a := ceilDiv(m, p) * ceilDiv(n, qq)
		b := ceilDiv(m, qq) * ceilDiv(n, p)
		best := minInt(a, b)
		q := edakit.Sprintf("Each DSP block multiplies a %d-bit unsigned operand by a %d-bit unsigned operand. A %d × %d unsigned multiplication is split into partial products, each on one DSP, in whichever operand orientation is cheaper. Counting only the partial-product multipliers, how many DSP blocks are needed?", p, qq, m, n)
		return num(q, edakit.IntStr(best),
			edakit.Sprintf("Splitting the %d-bit operand by %d and the %d-bit operand by %d gives %d × %d = %d; the other orientation gives %d × %d = %d. The minimum is %d.", m, p, n, qq, ceilDiv(m, p), ceilDiv(n, qq), a, ceilDiv(m, qq), ceilDiv(n, p), b, best),
			map[string]string{"orient_a": edakit.IntStr(a), "orient_b": edakit.IntStr(b)})
	default:
		d := 32 * edakit.Int(1, 32)
		w := edakit.Int(1, 64)
		l := ceilDiv(d, 64) * w
		q := edakit.Sprintf("Each 6-input LUT can act as a 64×1 distributed RAM. Ignoring any LUTs for read-address multiplexing, how many LUTs does a %d-word by %d-bit distributed RAM need?", d, w)
		return num(q, edakit.IntStr(l),
			edakit.Sprintf("ceil(%d / 64) = %d LUTs per bit column, times %d bit columns = %d LUTs. Small, shallow memories fit LUT RAM; deep ones belong in block RAM.", d, ceilDiv(d, 64), w, l),
			map[string]string{"luts": edakit.IntStr(l)})
	}
}

func routingArch(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		w := 8 * edakit.Int(2, 12)
		fin := edakit.Pick([]float64{0.25, 0.5, 0.75, 1})
		fout := edakit.Pick([]float64{0.125, 0.25, 0.5})
		i := edakit.Int(8, 40)
		o := edakit.Int(4, 10)
		ci := int(fin * float64(w))
		co := int(fout * float64(w))
		s := i*ci + o*co
		q := edakit.Sprintf("A logic block has %d input pins and %d output pins next to a channel of W = %d tracks. Each input pin connects to Fc_in·W tracks with Fc_in = %s, and each output pin to Fc_out·W tracks with Fc_out = %s. How many connection-box switches serve this block?", i, o, w, edakit.Num(fin, 3), edakit.Num(fout, 3))
		return num(q, edakit.IntStr(s),
			edakit.Sprintf("Inputs: %d × %d = %d switches. Outputs: %d × %d = %d switches. Total %d.", i, ci, i*ci, o, co, o*co, s),
			map[string]string{"in_switches": edakit.IntStr(i * ci), "out_switches": edakit.IntStr(o * co)})
	case 2:
		w := edakit.Int(4, 60)
		fs := 3
		m := edakit.Int(2, 12)
		per := 2 * fs * w
		if edakit.Int(0, 1) == 0 {
			q := edakit.Sprintf("An interior switch box joins four channel sides of W = %d tracks. Every track end that enters it can connect to Fs = %d track ends on the other three sides, and one switch joins two ends. How many switches does the switch box contain?", w, fs)
			return num(q, edakit.IntStr(per),
				edakit.Sprintf("There are 4W = %d track ends, each with Fs = %d connections, and each switch is counted from both of its ends: 4 × %d × %d / 2 = %d.", 4*w, fs, w, fs, per),
				map[string]string{"ends": edakit.IntStr(4 * w)})
		}
		tot := per * m * m
		q := edakit.Sprintf("A routing fabric has a %d × %d grid of interior switch boxes, each joining four channel sides of W = %d tracks with Fs = %d. How many switch-box switches does the fabric contain?", m, m, w, fs)
		return num(q, edakit.IntStr(tot),
			edakit.Sprintf("One box has 4 × %d × %d / 2 = %d switches; %d boxes have %d.", w, fs, per, m*m, tot),
			map[string]string{"per_box": edakit.IntStr(per)})
	default:
		w := edakit.Int(8, 16)
		na := edakit.Int(w/4, w/2)
		nb := edakit.Int(w/4, w/2)
		idx := make([]int, w)
		for i := range idx {
			idx[i] = i
		}
		a := pickDistinct(idx, na)
		b := pickDistinct(idx, nb)
		sortIntSlice(a)
		sortIntSlice(b)
		in := map[int]bool{}
		for _, x := range a {
			in[x] = true
		}
		common := 0
		var cs []int
		for _, x := range b {
			if in[x] {
				common++
				cs = append(cs, x)
			}
		}
		q := edakit.Sprintf("A channel has W = %d tracks numbered 0 to %d, and every switch box is disjoint: track i connects only to track i on the other sides. A source output pin connects to tracks %s; the destination input pin, several tiles away, connects to tracks %s. On how many track indices can a route from the source reach the destination pin?", w, w-1, joinInts(a), joinInts(b))
		ex := "none"
		if common > 0 {
			ex = joinInts(cs)
		}
		return num(q, edakit.IntStr(common),
			edakit.Sprintf("A disjoint switch box never changes the track index, so the fabric splits into W separate domains. The route must start and end on the same index, which leaves the common tracks (%s), %d in all. A Wilton switch box removes this restriction by rotating track indices.", ex, common),
			map[string]string{"common": edakit.IntStr(common)})
	}
}

func sortIntSlice(xs []int) {
	for i := 1; i < len(xs); i++ {
		for j := i; j > 0 && xs[j] < xs[j-1]; j-- {
			xs[j], xs[j-1] = xs[j-1], xs[j]
		}
	}
}

func segments(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		l := edakit.Pick([]int{2, 4, 8})
		w := l * edakit.Int(2, 12)
		n := l * edakit.Int(2, 8)
		s := n * w / l
		q := edakit.Sprintf("A horizontal channel is %d tiles long and has W = %d tracks, all built from length-%d wire segments laid end to end (each segment spans %d tiles). How many wire segments does the channel contain?", n, w, l, l)
		return num(q, edakit.IntStr(s),
			edakit.Sprintf("Each track holds %d / %d = %d segments, so %d tracks hold %d. Staggering the start points changes where segments begin, not how many there are per track.", n, l, n/l, w, s),
			map[string]string{"per_track": edakit.IntStr(n / l)})
	case 2:
		l := edakit.Pick([]int{1, 2, 4, 8})
		d := edakit.Int(1, 24)
		r := float64(10 * edakit.Int(5, 20))
		c := float64(10 * edakit.Int(2, 8))
		tb := float64(10 * edakit.Int(3, 8))
		seg := tb + 0.5*r*c*float64(l*l)/1000
		ns := ceilDiv(d, l)
		tot := float64(ns) * seg
		q := edakit.Sprintf("A connection travels %d tiles in a straight line on buffered length-%d wires, starting at a segment's start. Each segment adds a %s ps switch-buffer delay plus its distributed RC delay of 0.5·r·c·L², where r = %s ohm and c = %s fF per tile. What is the total delay, in ps? Round to 1 decimal place.", d, l, edakit.Num(tb, 0), edakit.Num(r, 0), edakit.Num(c, 0))
		return num(q, edakit.Num(tot, 1),
			edakit.Sprintf("The connection uses ceil(%d / %d) = %d segments. One segment costs %s + 0.5 × %s × %s × %d / 1000 = %s ps (ohm × fF = 0.001 ps). Total %d × %s = %s ps.", d, l, ns, edakit.Num(tb, 0), edakit.Num(r, 0), edakit.Num(c, 0), l*l, edakit.Num(seg, 2), ns, edakit.Num(seg, 2), edakit.Num(tot, 1)),
			map[string]string{"segments": edakit.IntStr(ns), "per_segment": edakit.Num(seg, 2)})
	default:
		for {
			d := edakit.Int(3, 32)
			r := float64(10 * edakit.Int(5, 30))
			c := float64(10 * edakit.Int(2, 10))
			tb := float64(10 * edakit.Int(3, 10))
			bestL, best, second := 0, math.Inf(1), math.Inf(1)
			var parts []string
			for _, l := range []int{1, 2, 4, 8, 16} {
				t := float64(ceilDiv(d, l)) * (tb + 0.5*r*c*float64(l*l)/1000)
				parts = append(parts, edakit.Sprintf("L = %d: %s ps", l, edakit.Num(t, 1)))
				if t < best {
					second = best
					best, bestL = t, l
				} else if t < second {
					second = t
				}
			}
			if second-best < 1 {
				continue
			}
			q := edakit.Sprintf("An architect compares buffered segment lengths L in {1, 2, 4, 8, 16} for a %d-tile straight connection that starts at a segment start. Each segment costs %s ps of switch-buffer delay plus 0.5·r·c·L², with r = %s ohm and c = %s fF per tile (ohm × fF = 0.001 ps). Which L gives the smallest delay?", d, edakit.Num(tb, 0), edakit.Num(r, 0), edakit.Num(c, 0))
			return num(q, edakit.IntStr(bestL),
				edakit.Sprintf("Delay = ceil(d/L)·(t_b + 0.5·r·c·L²): %s. The best is L = %d. Short wires pay many buffers; long wires pay quadratic RC and waste track on short hops.", strings.Join(parts, "; "), bestL),
				map[string]string{"best_l": edakit.IntStr(bestL)})
		}
	}
}
