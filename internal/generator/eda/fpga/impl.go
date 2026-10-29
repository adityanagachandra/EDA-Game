package fpga

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type timedPath struct {
	table [4][4]int
	locs  []pt
	tcq   int
	tlut  int
	tsu   int
}

func newTimedPath(ctx generator.GeneratorContext) *timedPath {
	p := &timedPath{tcq: 10 * edakit.Int(5, 15), tlut: 10 * edakit.Int(10, 30), tsu: 10 * edakit.Int(3, 10)}
	base := 10 * edakit.Int(3, 8)
	ax, ay := 10*edakit.Int(2, 6), 10*edakit.Int(2, 6)
	for dx := 0; dx < 4; dx++ {
		for dy := 0; dy < 4; dy++ {
			p.table[dx][dy] = base + ax*dx + ay*dy + 10*edakit.Int(0, 2)
		}
	}
	n := edakit.Scaled(ctx, 2, 4)
	cur := pt{edakit.Int(3, 6), edakit.Int(3, 6)}
	p.locs = []pt{cur}
	for i := 0; i <= n; i++ {
		nx := cur
		for nx == cur || nx.x < 1 || nx.y < 1 || nx.x > 12 || nx.y > 12 {
			nx = pt{cur.x + edakit.Int(-3, 3), cur.y + edakit.Int(-3, 3)}
		}
		p.locs = append(p.locs, nx)
		cur = nx
	}
	return p
}

func (p *timedPath) text() string {
	var rows []string
	for dx := 0; dx < 4; dx++ {
		var r []string
		for dy := 0; dy < 4; dy++ {
			r = append(r, edakit.IntStr(p.table[dx][dy]))
		}
		rows = append(rows, edakit.Sprintf("dx = %d: %s", dx, strings.Join(r, ", ")))
	}
	var hops []string
	for i, l := range p.locs {
		name := edakit.Sprintf("LUT%d", i)
		if i == 0 {
			name = "launch FF"
		} else if i == len(p.locs)-1 {
			name = "capture FF"
		}
		hops = append(hops, edakit.Sprintf("%s at (%d, %d)", name, l.x, l.y))
	}
	return edakit.Sprintf("The placement delay model gives the routing delay of a connection in ps from its |dx| and |dy|, listed for dy = 0, 1, 2, 3: %s. The path runs %s. Clock-to-q is %d ps, every LUT adds %d ps and setup is %d ps; the clock is ideal.", strings.Join(rows, "; "), strings.Join(hops, " → "), p.tcq, p.tlut, p.tsu)
}

func (p *timedPath) arrival() (int, []int) {
	t := p.tcq
	var ds []int
	for i := 1; i < len(p.locs); i++ {
		d := p.table[absInt(p.locs[i].x-p.locs[i-1].x)][absInt(p.locs[i].y-p.locs[i-1].y)]
		ds = append(ds, d)
		t += d
		if i < len(p.locs)-1 {
			t += p.tlut
		}
	}
	return t, ds
}

func delayModel(ctx generator.GeneratorContext) generator.Problem {
	p := newTimedPath(ctx)
	a, ds := p.arrival()
	nl := len(p.locs) - 2
	base := edakit.Sprintf("Arrival = %d (clock-to-q) + routing %s + %d LUTs × %d = %d ps.", p.tcq, joinInts(ds), nl, p.tlut, a)
	switch whichKP(ctx) {
	case 1:
		return num(p.text()+" What is the arrival time at the capture flip-flop's D pin, in ps?", edakit.IntStr(a),
			base, map[string]string{"arrival": edakit.IntStr(a)})
	case 2:
		per := 10 * edakit.Int((a+p.tsu)/10-15, (a+p.tsu)/10+20)
		s := per - p.tsu - a
		return num(p.text()+edakit.Sprintf(" The clock period is %d ps. What is the setup slack, in ps?", per), edakit.IntStr(s),
			base+edakit.Sprintf(" Slack = %d - %d - %d = %d ps.", per, p.tsu, a, s),
			map[string]string{"arrival": edakit.IntStr(a)})
	default:
		f := 1e6 / float64(a+p.tsu)
		return num(p.text()+" If this is the critical path, what is the maximum clock frequency, in MHz? Round to 1 decimal place.", edakit.Num(f, 1),
			base+edakit.Sprintf(" The period must cover %d + %d = %d ps, so Fmax = 10^6 / %d = %s MHz.", a, p.tsu, a+p.tsu, a+p.tsu, edakit.Num(f, 1)),
			map[string]string{"period": edakit.IntStr(a + p.tsu)})
	}
}

func bitstream(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		k := edakit.Int(4, 6)
		n := edakit.Pick([]int{4, 8, 10})
		s := 10 * edakit.Int(20, 150)
		m := edakit.Scaled(ctx, 4, 60)
		tile := n*(pow2(k)+1) + s
		tot := tile * m * m
		q := edakit.Sprintf("Each tile of a %d × %d FPGA holds one cluster of %d BLEs (a %d-LUT plus one output-select bit each) and %d routing-switch configuration bits. Ignoring I/O and clocking, how many configuration bits does the device hold?", m, m, n, k, s)
		return num(q, edakit.IntStr(tot),
			edakit.Sprintf("One tile: %d × (%d + 1) + %d = %d bits. %d tiles: %d bits. Routing bits usually outnumber logic bits, which is why the routing fabric dominates FPGA area.", n, pow2(k), s, tile, m*m, tot),
			map[string]string{"tile_bits": edakit.IntStr(tile)})
	case 2:
		bits := 1000 * edakit.Int(500, 30000)
		if edakit.Int(0, 1) == 0 {
			fw := 32 * edakit.Int(50, 130)
			fr := ceilDiv(bits, fw)
			q := edakit.Sprintf("A device holds %d configuration bits, written in frames of %d bits each (a frame is the smallest unit that can be addressed). How many frames does a full configuration contain?", bits, fw)
			return num(q, edakit.IntStr(fr),
				edakit.Sprintf("ceil(%d / %d) = %d frames.", bits, fw, fr),
				map[string]string{"frames": edakit.IntStr(fr)})
		}
		w := edakit.Pick([]int{1, 8, 16, 32})
		f := edakit.Pick([]int{25, 50, 66, 100})
		t := float64(bits) / float64(w*f) / 1000
		q := edakit.Sprintf("A %d-bit bitstream is loaded over a %d-bit-wide configuration port clocked at %d MHz, one word per cycle. How long does configuration take, in ms? Round to 2 decimal places.", bits, w, f)
		return num(q, edakit.Num(t, 2),
			edakit.Sprintf("%d / (%d × %d × 10^6) s = %s ms.", bits, w, f, edakit.Num(t, 2)),
			map[string]string{"cycles": edakit.IntStr(ceilDiv(bits, w))})
	default:
		mx := edakit.Scaled(ctx, 100, 5000)
		m := edakit.Int(5, 40)
		enc := mx * bitsFor(m)
		oh := mx * m
		if edakit.Int(0, 1) == 0 {
			q := edakit.Sprintf("A fabric has %d routing multiplexers, each choosing among %d inputs. If each multiplexer's select is stored in binary-encoded form, how many configuration bits do the multiplexers need?", mx, m)
			return num(q, edakit.IntStr(enc),
				edakit.Sprintf("ceil(log2 %d) = %d bits each, %d in total; one-hot storage would need %d.", m, bitsFor(m), enc, oh),
				map[string]string{"bits_each": edakit.IntStr(bitsFor(m))})
		}
		q := edakit.Sprintf("A fabric has %d routing multiplexers with %d inputs each. How many configuration bits does binary-encoded select storage save compared with one bit per multiplexer input?", mx, m)
		return num(q, edakit.IntStr(oh-enc),
			edakit.Sprintf("One-hot: %d × %d = %d. Encoded: %d × %d = %d. Saving %d bits, paid for with a decoder in front of each multiplexer.", mx, m, oh, mx, bitsFor(m), enc, oh-enc),
			map[string]string{"one_hot": edakit.IntStr(oh), "encoded": edakit.IntStr(enc)})
	}
}

func partialReconfig(ctx generator.GeneratorContext) generator.Problem {
	c := edakit.Scaled(ctx, 2, 30)
	r := edakit.Int(1, 4)
	f := edakit.Pick([]int{28, 36, 128})
	wd := edakit.Pick([]int{93, 101, 123})
	bytes := c * r * f * wd * 4
	switch whichKP(ctx) {
	case 1:
		kb := float64(bytes) / 1024
		q := edakit.Sprintf("A reconfigurable region spans %d logic columns over %d clock-region rows. In this device each such column-row needs %d configuration frames, and a frame is %d 32-bit words. Ignoring headers and commands, how large is the partial bitstream, in KB (1 KB = 1024 bytes)? Round to 1 decimal place.", c, r, f, wd)
		return num(q, edakit.Num(kb, 1),
			edakit.Sprintf("%d × %d × %d frames × %d words × 4 bytes = %d bytes = %s KB. A partial bitstream scales with the region, not with the device.", c, r, f, wd, bytes, edakit.Num(kb, 1)),
			map[string]string{"bytes": edakit.IntStr(bytes)})
	case 2:
		mhz := edakit.Pick([]int{50, 100, 125, 200})
		us := float64(bytes) / (4 * float64(mhz))
		q := edakit.Sprintf("A %d-byte partial bitstream is written through an internal configuration port that accepts one 32-bit word per cycle at %d MHz. How long does reconfiguration take, in microseconds? Round to 1 decimal place.", bytes, mhz)
		return num(q, edakit.Num(us, 1),
			edakit.Sprintf("Throughput is 4 bytes × %d × 10^6 per second, so %d bytes take %d / (4 × %d) = %s microseconds.", mhz, bytes, bytes, mhz, edakit.Num(us, 1)),
			map[string]string{"throughput_mb_s": edakit.IntStr(4 * mhz)})
	default:
		tr := float64(edakit.Int(1, 40)) / 2
		tt := float64(edakit.Int(1, 200)) / 2
		u := 100 * tt / (tt + tr)
		q := edakit.Sprintf("A region alternates between two accelerators: each swap reconfigures the region in %s ms, then the loaded accelerator runs for %s ms before the next swap. What percentage of the time is the region computing? Round to 1 decimal place.", edakit.Num(tr, 1), edakit.Num(tt, 1))
		return num(q, edakit.Num(u, 1),
			edakit.Sprintf("%s / (%s + %s) × 100 = %s%%. Partial reconfiguration pays off only when tasks run long compared with the swap, and the static region must hold the interface steady while the region is being rewritten.", edakit.Num(tt, 1), edakit.Num(tt, 1), edakit.Num(tr, 1), edakit.Num(u, 1)),
			map[string]string{"utilization": edakit.Num(u, 1)})
	}
}
