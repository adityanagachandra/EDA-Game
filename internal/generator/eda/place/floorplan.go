package place

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type blk struct{ w, h int }

func randBlocks(n, lo, hi int) []blk {
	out := make([]blk, n)
	for i := range out {
		out[i] = blk{edakit.Int(lo, hi), edakit.Int(lo, hi)}
	}
	return out
}

func fmtBlocks(names []string, bs []blk) string {
	parts := make([]string, len(bs))
	for i, b := range bs {
		parts[i] = edakit.Sprintf("%s is %d x %d", names[i], b.w, b.h)
	}
	return strings.Join(parts, ", ")
}

// whitespace computes whitespace fractions and outline sizes.
func whitespace(ctx generator.GeneratorContext) generator.Problem {
	switch ctx.KPOrdinal {
	case 2:
		a := blk{edakit.Int(10, 60), edakit.Int(10, 60)}
		b := blk{edakit.Int(10, 60), edakit.Int(10, 60)}
		W, H := a.w+b.w, maxI(a.h, b.h)
		ws := 100 * float64(W*H-a.w*a.h-b.w*b.h) / float64(W*H)
		ans := edakit.Num(ws, 1)
		return edakit.Numeric(edakit.Sprintf("Block A is %d x %d um and block B is %d x %d um (width x height). They are placed side by side, bottoms aligned, with no gap, and the floorplan outline is their bounding box. What percentage of the outline is whitespace? Round to 1 decimal place.", a.w, a.h, b.w, b.h), ans,
			edakit.Sprintf("The outline is %d x %d = %d um^2 and the blocks cover %d um^2, so whitespace is %s%%.", W, H, W*H, a.w*a.h+b.w*b.h, ans),
			map[string]string{"outline": edakit.IntStr(W * H)})
	case 3:
		area := 100 * edakit.Scaled(ctx, 20, 200)
		ws := edakit.Pick([]float64{0.05, 0.1, 0.15, 0.2, 0.25})
		H := edakit.Int(40, 160)
		W := float64(area) / ((1 - ws) * float64(H))
		ans := edakit.Num(W, 2)
		return edakit.Numeric(edakit.Sprintf("The blocks of a floorplan have total area %d um^2. A fixed-outline flow fixes the outline height at %d um and allows %s%% whitespace. What outline width is required, in um? Round to 2 decimal places.", area, H, edakit.Num(100*ws, 0)), ans,
			edakit.Sprintf("Block area must be %s of the outline area, so the outline is %d / %s = %s um^2, and the width is that over %d: %s um.", edakit.Num(1-ws, 2), area, edakit.Num(1-ws, 2), edakit.Num(float64(area)/(1-ws), 2), H, ans),
			map[string]string{"outline_area": edakit.Num(float64(area)/(1-ws), 2)})
	}
	n := edakit.Scaled(ctx, 3, 6)
	bs := randBlocks(n, 5, 40)
	sum := 0
	for _, b := range bs {
		sum += b.w * b.h
	}
	W := int(math.Ceil(math.Sqrt(float64(sum)*1.3))) + edakit.Int(0, 10)
	H := int(math.Ceil(float64(sum)*1.1/float64(W))) + edakit.Int(1, 10)
	ws := 100 * float64(W*H-sum) / float64(W*H)
	ans := edakit.Num(ws, 1)
	names := make([]string, n)
	for i := range names {
		names[i] = letter(i)
	}
	return edakit.Numeric(edakit.Sprintf("A floorplan packs %d blocks (%s, width x height in um) into an outline %d um wide and %d um tall without overlap. What percentage of the outline is whitespace? Round to 1 decimal place.", n, fmtBlocks(names, bs), W, H), ans,
		edakit.Sprintf("The blocks cover %d um^2 of the %d um^2 outline, leaving %d um^2, which is %s%%.", sum, W*H, W*H-sum, ans),
		map[string]string{"block_area": edakit.IntStr(sum)})
}

// snode is a slicing-tree node: a leaf block or a V or H cut.
type snode struct {
	op   byte
	l, r *snode
	b    int
	name string
}

func buildSlicing(lo, hi int, k *int) *snode {
	if hi-lo == 1 {
		return &snode{b: lo}
	}
	m := edakit.Int(lo+1, hi-1)
	n := &snode{l: buildSlicing(lo, m, k), r: buildSlicing(m, hi, k)}
	n.op = edakit.Pick([]byte{'V', 'H'})
	if n.r.op != 0 {
		if n.r.op == 'V' {
			n.op = 'H'
		} else {
			n.op = 'V'
		}
	}
	return n
}

func (n *snode) dims(bs []blk) blk {
	if n.op == 0 {
		return bs[n.b]
	}
	a, b := n.l.dims(bs), n.r.dims(bs)
	if n.op == 'V' {
		return blk{a.w + b.w, maxI(a.h, b.h)}
	}
	return blk{maxI(a.w, b.w), a.h + b.h}
}

func (n *snode) postfix(labels []string) []string {
	if n.op == 0 {
		return []string{labels[n.b]}
	}
	return append(append(n.l.postfix(labels), n.r.postfix(labels)...), string(n.op))
}

// describe names the internal nodes r1, r2, ... in pre-order and lists them.
func (n *snode) describe(labels []string) string {
	var lines []string
	cnt := 0
	var name func(m *snode)
	name = func(m *snode) {
		if m.op == 0 {
			m.name = "block " + labels[m.b]
			return
		}
		cnt++
		m.name = edakit.Sprintf("node r%d", cnt)
		name(m.l)
		name(m.r)
	}
	name(n)
	var walk func(m *snode)
	walk = func(m *snode) {
		if m.op == 0 {
			return
		}
		if m.op == 'V' {
			lines = append(lines, edakit.Sprintf("%s is a vertical cut with %s on the left and %s on the right", m.name, m.l.name, m.r.name))
		} else {
			lines = append(lines, edakit.Sprintf("%s is a horizontal cut with %s below and %s above", m.name, m.l.name, m.r.name))
		}
		walk(m.l)
		walk(m.r)
	}
	walk(n)
	return strings.Join(lines, "; ")
}

func shuffledLabels(n int, numeric bool) []string {
	out := make([]string, n)
	for i := range out {
		if numeric {
			out[i] = edakit.IntStr(i + 1)
		} else {
			out[i] = letter(i)
		}
	}
	return edakit.Shuffle(out)
}

func labelBlocks(labels []string, bs []blk) string {
	type pair struct {
		l string
		b blk
	}
	ps := make([]pair, len(bs))
	for i := range bs {
		ps[i] = pair{labels[i], bs[i]}
	}
	// list in label order so the description does not reveal the tree.
	for i := 0; i < len(ps); i++ {
		for j := i + 1; j < len(ps); j++ {
			if ps[j].l < ps[i].l {
				ps[i], ps[j] = ps[j], ps[i]
			}
		}
	}
	parts := make([]string, len(ps))
	for i, p := range ps {
		parts[i] = edakit.Sprintf("%s is %d x %d", p.l, p.b.w, p.b.h)
	}
	return strings.Join(parts, ", ")
}

// slicing evaluates the area of a slicing tree of fixed-shape blocks.
func slicing(ctx generator.GeneratorContext) generator.Problem {
	n := 3
	if ctx.KPOrdinal != 1 {
		n = edakit.Scaled(ctx, 4, 6)
	}
	bs := randBlocks(n, 2, 12)
	k := 0
	root := buildSlicing(0, n, &k)
	labels := shuffledLabels(n, false)
	d := root.dims(bs)
	sum := 0
	for _, b := range bs {
		sum += b.w * b.h
	}
	q := edakit.Sprintf("Blocks have fixed shapes, width x height in um: %s. A slicing tree has root r1: %s. A vertical cut puts its children side by side (widths add, height is the larger), and a horizontal cut stacks them (heights add, width is the larger). ",
		labelBlocks(labels, bs), root.describe(labels))
	if ctx.KPOrdinal == 3 {
		ws := d.w*d.h - sum
		return edakit.Numeric(q+"How much whitespace, in um^2, does the resulting floorplan contain inside its bounding box?", edakit.IntStr(ws),
			edakit.Sprintf("Evaluating bottom-up gives a %d x %d floorplan of area %d um^2. The blocks cover %d um^2, so whitespace is %d um^2.", d.w, d.h, d.w*d.h, sum, ws),
			map[string]string{"w": edakit.IntStr(d.w), "h": edakit.IntStr(d.h)})
	}
	return edakit.Numeric(q+"What is the area of the floorplan's bounding box, in um^2?", edakit.IntStr(d.w*d.h),
		edakit.Sprintf("Evaluate the tree bottom-up: the root comes out %d um wide and %d um tall, so the area is %d um^2.", d.w, d.h, d.w*d.h),
		map[string]string{"w": edakit.IntStr(d.w), "h": edakit.IntStr(d.h)})
}

const polishRule = "In a Polish expression, operands are blocks and operators are cuts: \"P Q V\" places P to the left of Q (widths add, height is the larger) and \"P Q H\" places P below Q (heights add, width is the larger). "

// polish evaluates a normalized Polish expression with a stack.
func polish(ctx generator.GeneratorContext) generator.Problem {
	n := 3
	if ctx.KPOrdinal != 1 {
		n = edakit.Scaled(ctx, 4, 6)
	}
	bs := randBlocks(n, 2, 12)
	k := 0
	root := buildSlicing(0, n, &k)
	labels := shuffledLabels(n, true)
	expr := strings.Join(root.postfix(labels), " ")
	d := root.dims(bs)
	sum := 0
	for _, b := range bs {
		sum += b.w * b.h
	}
	head := polishRule + edakit.Sprintf("Blocks, width x height in um: %s. The expression is %s. ", labelBlocks(labels, bs), expr)
	switch ctx.KPOrdinal {
	case 1:
		dim, val := "width", d.w
		if edakit.Int(0, 1) == 1 {
			dim, val = "height", d.h
		}
		return edakit.Numeric(head+edakit.Sprintf("What is the %s of the floorplan, in um?", dim), edakit.IntStr(val),
			edakit.Sprintf("Scan left to right with a stack: push blocks, and on an operator pop two shapes and push their combination. The final shape is %d x %d.", d.w, d.h),
			map[string]string{"w": edakit.IntStr(d.w), "h": edakit.IntStr(d.h)})
	case 3:
		return edakit.Numeric(head+"How much whitespace, in um^2, lies inside the floorplan's bounding box?", edakit.IntStr(d.w*d.h-sum),
			edakit.Sprintf("The stack evaluation gives %d x %d = %d um^2 against %d um^2 of blocks, so whitespace is %d um^2.", d.w, d.h, d.w*d.h, sum, d.w*d.h-sum),
			map[string]string{"area": edakit.IntStr(d.w * d.h)})
	}
	return edakit.Numeric(head+"What is the floorplan's area, in um^2?", edakit.IntStr(d.w*d.h),
		edakit.Sprintf("Evaluating with a stack gives a %d x %d floorplan, area %d um^2.", d.w, d.h, d.w*d.h),
		map[string]string{"w": edakit.IntStr(d.w), "h": edakit.IntStr(d.h)})
}

func isOp(t string) bool { return t == "V" || t == "H" }

func flip(t string) string {
	if t == "V" {
		return "H"
	}
	return "V"
}

// legalPolish checks the balloting property and normalization of a token list.
func legalPolish(ts []string) (bool, string) {
	ops, opd := 0, 0
	for i, t := range ts {
		if isOp(t) {
			ops++
			if i > 0 && ts[i-1] == t {
				return false, "it has two identical operators in a row, so it is not normalized"
			}
		} else {
			opd++
		}
		if ops >= opd {
			return false, "a prefix has as many operators as operands, which breaks the balloting property"
		}
	}
	return true, ""
}

// polishMoves asks for the result of M1 or M2, or which M3 candidate is legal.
func polishMoves(ctx generator.GeneratorContext) generator.Problem {
	mode := ctx.KPOrdinal
	if mode == 0 {
		mode = edakit.Int(1, 3)
	}
	for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
		n := edakit.Int(4, 6)
		k := 0
		root := buildSlicing(0, n, &k)
		ts := root.postfix(shuffledLabels(n, true))
		expr := strings.Join(ts, " ")
		join := func(x []string) string { return strings.Join(x, " ") }
		var opdIdx []int
		for i, t := range ts {
			if !isOp(t) {
				opdIdx = append(opdIdx, i)
			}
		}
		switch mode {
		case 1:
			i := edakit.Int(0, len(opdIdx)-2)
			swapOps := func(a, b int) []string {
				c := append([]string(nil), ts...)
				c[opdIdx[a]], c[opdIdx[b]] = c[opdIdx[b]], c[opdIdx[a]]
				return c
			}
			ans := join(swapOps(i, i+1))
			var d []string
			for j := 0; j < len(opdIdx)-1; j++ {
				if j != i {
					d = append(d, join(swapOps(j, j+1)))
				}
			}
			if i+2 < len(opdIdx) {
				d = append(d, join(swapOps(i, i+2)))
			}
			if i >= 1 {
				d = append(d, join(swapOps(i-1, i+1)))
			}
			d = edakit.Shuffle(d)
			if len(d) < 3 {
				continue
			}
			q := polishRule + edakit.Sprintf("A simulated-annealing floorplanner holds the normalized expression %s. It applies move M1 to operands %s and %s, which are adjacent in the operand sequence. What is the new expression?", expr, ts[opdIdx[i]], ts[opdIdx[i+1]])
			return edakit.Choice(q, ans, "M1 swaps two operands that are adjacent once operators are ignored, leaving every operator where it was. The result is "+ans+".", d[:3]...)
		case 2:
			type chain struct{ s, e int }
			var chains []chain
			for i := 0; i < len(ts); i++ {
				if isOp(ts[i]) {
					j := i
					for j+1 < len(ts) && isOp(ts[j+1]) {
						j++
					}
					chains = append(chains, chain{i, j})
					i = j
				}
			}
			if len(chains) < 2 {
				continue
			}
			ci := edakit.Int(0, len(chains)-1)
			c := chains[ci]
			comp := func(s, e int) []string {
				x := append([]string(nil), ts...)
				for t := s; t <= e; t++ {
					x[t] = flip(x[t])
				}
				return x
			}
			ans := join(comp(c.s, c.e))
			var d []string
			for cj, o := range chains {
				if cj != ci {
					d = append(d, join(comp(o.s, o.e)))
				}
			}
			if c.e > c.s {
				d = append(d, join(comp(c.s, c.s)), join(comp(c.e, c.e)))
			}
			all := append([]string(nil), ts...)
			for t := range all {
				if isOp(all[t]) {
					all[t] = flip(all[t])
				}
			}
			d = append(d, join(all))
			seen := map[string]bool{ans: true}
			var dd []string
			for _, x := range edakit.Shuffle(d) {
				if !seen[x] {
					seen[x] = true
					dd = append(dd, x)
				}
			}
			if len(dd) < 3 {
				continue
			}
			ord := []string{"first", "second", "third", "fourth", "fifth"}[ci]
			q := polishRule + edakit.Sprintf("A simulated-annealing floorplanner holds the normalized expression %s. A chain is a maximal run of consecutive operators. It applies move M2 to the %s chain from the left. What is the new expression?", expr, ord)
			return edakit.Choice(q, ans, "M2 complements every operator in one chain, V to H and H to V. Complementing the whole chain keeps the expression normalized. The result is "+ans+".", dd[:3]...)
		default:
			var good, bad []string
			var badWhy []string
			for i := 0; i+1 < len(ts); i++ {
				if isOp(ts[i]) == isOp(ts[i+1]) {
					continue
				}
				x := append([]string(nil), ts...)
				x[i], x[i+1] = x[i+1], x[i]
				if ok, why := legalPolish(x); ok {
					good = append(good, join(x))
				} else {
					bad = append(bad, join(x))
					badWhy = append(badWhy, why)
				}
			}
			if len(good) == 0 || len(bad) < 3 {
				continue
			}
			ans := edakit.Pick(good)
			q := polishRule + edakit.Sprintf("A simulated-annealing floorplanner holds the normalized expression %s. Move M3 swaps an adjacent operand and operator. Which of these candidate results is a legal normalized Polish expression the annealer may accept?", expr)
			ids := make([]int, len(bad))
			for t := range ids {
				ids[t] = t
			}
			ord := edakit.Shuffle(ids)
			shown := []string{bad[ord[0]], bad[ord[1]], bad[ord[2]]}
			return edakit.Choice(q, ans, "M3 must be checked: the result has to keep the balloting property (every prefix has more operands than operators) and stay normalized (no two identical operators in a row). "+ans+" passes both. Each rejected candidate fails one of them; "+shown[0]+" fails because "+badWhy[ord[0]]+".", shown...)
		}
	}
	// Known-good fallback: an M1 move on a fixed normalized expression.
	q := polishRule + "A simulated-annealing floorplanner holds the normalized expression 1 2 V 3 4 H V. It applies move M1 to operands 2 and 3, which are adjacent in the operand sequence. What is the new expression?"
	return edakit.Choice(q, "1 3 V 2 4 H V", "M1 swaps two operands that are adjacent once operators are ignored, leaving every operator where it was. The result is 1 3 V 2 4 H V.",
		"2 1 V 3 4 H V", "1 2 V 4 3 H V", "1 4 V 3 2 H V")
}

// shapeCurves combines two blocks under a vertical cut, with rotation.
func shapeCurves(ctx generator.GeneratorContext) generator.Problem {
	a, b := blk{3, 5}, blk{4, 7}
	for try := 0; try < generator.MaxDrawAttempts; try++ {
		ca := blk{edakit.Int(2, 15), edakit.Int(2, 15)}
		cb := blk{edakit.Int(2, 15), edakit.Int(2, 15)}
		if ca.w != ca.h && cb.w != cb.h {
			a, b = ca, cb
			break
		}
	}
	oa := []blk{a, {a.h, a.w}}
	ob := []blk{b, {b.h, b.w}}
	if ctx.KPOrdinal == 1 {
		ob = ob[:1]
	}
	type sh struct{ w, h int }
	var cands []sh
	for _, x := range oa {
		for _, y := range ob {
			cands = append(cands, sh{x.w + y.w, maxI(x.h, y.h)})
		}
	}
	best := cands[0]
	for _, c := range cands {
		if c.w*c.h < best.w*best.h {
			best = c
		}
	}
	rot := "Both blocks may be rotated by 90 degrees."
	if ctx.KPOrdinal == 1 {
		rot = "Block A may be rotated by 90 degrees; block B is fixed."
	}
	head := edakit.Sprintf("Block A is %d x %d um and block B is %d x %d um (width x height). A vertical cut places A to the left of B. %s ", a.w, a.h, b.w, b.h, rot)
	if ctx.KPOrdinal == 3 {
		uniq := map[sh]bool{}
		for _, c := range cands {
			uniq[c] = true
		}
		cnt := 0
		var list []string
		for c := range uniq {
			dom := false
			for o := range uniq {
				if o != c && o.w <= c.w && o.h <= c.h {
					dom = true
				}
			}
			if !dom {
				cnt++
			}
			list = append(list, edakit.Sprintf("%d x %d", c.w, c.h))
		}
		return edakit.Numeric(head+"The shape function of the combined block keeps only non-dominated (width, height) pairs: a pair is dropped if another pair is no wider and no taller. How many distinct pairs remain on the combined shape curve?", edakit.IntStr(cnt),
			edakit.Sprintf("Width adds and height takes the maximum for every orientation choice, giving the distinct shapes %s. Removing dominated shapes leaves %d.", strings.Join(sortedStrings(list), ", "), cnt),
			map[string]string{"candidates": edakit.IntStr(len(uniq))})
	}
	return edakit.Numeric(head+"What is the smallest bounding-box area, in um^2, over the allowed orientations?", edakit.IntStr(best.w*best.h),
		edakit.Sprintf("Each orientation choice gives width = sum of widths and height = max of heights. The best is %d x %d = %d um^2.", best.w, best.h, best.w*best.h),
		map[string]string{"w": edakit.IntStr(best.w), "h": edakit.IntStr(best.h)})
}

func sortedStrings(xs []string) []string {
	out := append([]string(nil), xs...)
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j] < out[i] {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}

// seqPairPack computes block coordinates from a sequence pair by longest paths.
func seqPairPack(gp, gm []int, bs []blk) (xs, ys []int) {
	n := len(bs)
	pp, pm := make([]int, n), make([]int, n)
	for i, b := range gp {
		pp[b] = i
	}
	for i, b := range gm {
		pm[b] = i
	}
	xs, ys = make([]int, n), make([]int, n)
	for _, b := range gm {
		for _, a := range gm {
			if a == b {
				break
			}
			if pp[a] < pp[b] && pm[a] < pm[b] {
				xs[b] = maxI(xs[b], xs[a]+bs[a].w)
			}
			if pp[a] > pp[b] && pm[a] < pm[b] {
				ys[b] = maxI(ys[b], ys[a]+bs[a].h)
			}
		}
	}
	return
}

func seqString(g []int) string {
	parts := make([]string, len(g))
	for i, b := range g {
		parts[i] = letter(b)
	}
	return "(" + strings.Join(parts, " ") + ")"
}

// sequencePair packs a sequence pair and asks for a coordinate, the width or the area.
func sequencePair(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 6)
	bs := randBlocks(n, 1, 8)
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	gp, gm := edakit.Shuffle(idx), edakit.Shuffle(idx)
	xs, ys := seqPairPack(gp, gm, bs)
	W, H := 0, 0
	for i := range bs {
		W = maxI(W, xs[i]+bs[i].w)
		H = maxI(H, ys[i]+bs[i].h)
	}
	names := make([]string, n)
	for i := range names {
		names[i] = letter(i)
	}
	head := edakit.Sprintf("Blocks (width x height): %s. A floorplan is encoded by the sequence pair $\\Gamma_+$ = %s, $\\Gamma_-$ = %s. Block a is left of b when a precedes b in both sequences, and a is below b when a follows b in $\\Gamma_+$ but precedes it in $\\Gamma_-$. Each block is packed as far left and as far down as these constraints allow, starting from (0, 0). ",
		fmtBlocks(names, bs), seqString(gp), seqString(gm))
	facts := map[string]string{"width": edakit.IntStr(W), "height": edakit.IntStr(H)}
	switch ctx.KPOrdinal {
	case 1:
		b := edakit.Int(0, n-1)
		for tries := 0; xs[b] == 0 && tries < 10; tries++ {
			b = edakit.Int(0, n-1)
		}
		axis, v := "x", xs[b]
		if edakit.Int(0, 1) == 1 {
			axis, v = "y", ys[b]
		}
		return edakit.Numeric(head+edakit.Sprintf("What is the %s coordinate of the lower-left corner of block %s?", axis, letter(b)), edakit.IntStr(v),
			edakit.Sprintf("The coordinate is the longest path to %s in the %s constraint graph. Packing gives x = %s and y = %s for blocks %s in order.", letter(b), map[string]string{"x": "horizontal", "y": "vertical"}[axis], fmtInts(xs), fmtInts(ys), strings.Join(names, ", ")), facts)
	case 3:
		return edakit.Numeric(head+"What is the area of the packed floorplan's bounding box?", edakit.IntStr(W*H),
			edakit.Sprintf("The width is the longest path in the horizontal graph, %d, and the height the longest path in the vertical graph, %d. Area = %d.", W, H, W*H), facts)
	}
	return edakit.Numeric(head+"What is the width of the packed floorplan?", edakit.IntStr(W),
		edakit.Sprintf("Longest paths in the horizontal constraint graph give x = %s for blocks %s; the widest extent x + w is %d.", fmtInts(xs), strings.Join(names, ", "), W), facts)
}

// bstar packs a B*-tree with a contour and asks for coordinates or area.
func bstar(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 6)
	bs := randBlocks(n, 1, 8)
	left, right := make([]int, n), make([]int, n)
	for i := range left {
		left[i], right[i] = -1, -1
	}
	var edges []string
	for i := 1; i < n; i++ {
		placedChild := false
		for try := 0; try < generator.MaxDrawAttempts && !placedChild; try++ {
			p := edakit.Int(0, i-1)
			if edakit.Int(0, 1) == 0 && left[p] < 0 {
				left[p] = i
				edges = append(edges, edakit.Sprintf("%s's left child is %s", letter(p), letter(i)))
				placedChild = true
			} else if right[p] < 0 {
				right[p] = i
				edges = append(edges, edakit.Sprintf("%s's right child is %s", letter(p), letter(i)))
				placedChild = true
			}
		}
		for p := 0; p < i && !placedChild; p++ {
			if left[p] < 0 {
				left[p] = i
				edges = append(edges, edakit.Sprintf("%s's left child is %s", letter(p), letter(i)))
				placedChild = true
			} else if right[p] < 0 {
				right[p] = i
				edges = append(edges, edakit.Sprintf("%s's right child is %s", letter(p), letter(i)))
				placedChild = true
			}
		}
	}
	xs, ys := make([]int, n), make([]int, n)
	var placed []int
	var dfs func(v int)
	dfs = func(v int) {
		y := 0
		for _, u := range placed {
			if xs[u] < xs[v]+bs[v].w && xs[v] < xs[u]+bs[u].w {
				y = maxI(y, ys[u]+bs[u].h)
			}
		}
		ys[v] = y
		placed = append(placed, v)
		if left[v] >= 0 {
			xs[left[v]] = xs[v] + bs[v].w
			dfs(left[v])
		}
		if right[v] >= 0 {
			xs[right[v]] = xs[v]
			dfs(right[v])
		}
	}
	dfs(0)
	W, H := 0, 0
	for i := range bs {
		W = maxI(W, xs[i]+bs[i].w)
		H = maxI(H, ys[i]+bs[i].h)
	}
	names := make([]string, n)
	for i := range names {
		names[i] = letter(i)
	}
	head := edakit.Sprintf("Blocks (width x height): %s. A B*-tree has root A at (0, 0); %s. A left child is placed immediately to the right of its parent (x = parent x + parent width), a right child at its parent's x. Blocks are placed in pre-order, and each is dropped onto the contour: its y is the highest top edge of the already placed blocks that overlap it in x. ",
		fmtBlocks(names, bs), strings.Join(edges, ", "))
	facts := map[string]string{"xs": fmtInts(xs), "ys": fmtInts(ys)}
	switch ctx.KPOrdinal {
	case 1:
		b := edakit.Int(1, n-1)
		return edakit.Numeric(head+edakit.Sprintf("What is the x coordinate of block %s?", letter(b)), edakit.IntStr(xs[b]),
			edakit.Sprintf("Follow the tree from A: every left edge adds the parent's width and every right edge keeps the parent's x. Block %s lands at x = %d.", letter(b), xs[b]), facts)
	case 2:
		b := edakit.Int(1, n-1)
		for tries := 0; ys[b] == 0 && tries < 12; tries++ {
			b = edakit.Int(1, n-1)
		}
		return edakit.Numeric(head+edakit.Sprintf("What is the y coordinate of block %s?", letter(b)), edakit.IntStr(ys[b]),
			edakit.Sprintf("Packing in pre-order gives x = %s and y = %s for blocks %s. Block %s sits on the contour at y = %d.", fmtInts(xs), fmtInts(ys), strings.Join(names, ", "), letter(b), ys[b]), facts)
	}
	return edakit.Numeric(head+"What is the area of the packed floorplan's bounding box?", edakit.IntStr(W*H),
		edakit.Sprintf("Packing gives x = %s and y = %s for blocks %s. The bounding box is %d x %d = %d.", fmtInts(xs), fmtInts(ys), strings.Join(names, ", "), W, H, W*H), facts)
}

// annealing computes cost changes, acceptance probabilities and cooling steps.
func annealing(ctx generator.GeneratorContext) generator.Problem {
	alpha := edakit.Pick([]float64{0.5, 1, 2})
	beta := edakit.Pick([]float64{0.5, 1, 2, 3})
	switch ctx.KPOrdinal {
	case 2:
		d := float64(edakit.Int(2, 60))
		T := float64(edakit.Pick([]int{10, 20, 25, 40, 50, 80, 100}))
		p := math.Exp(-d / T)
		ans := edakit.Num(p, 4)
		return edakit.Numeric(edakit.Sprintf("At temperature T = %s, a simulated-annealing floorplanner proposes an M3 move that raises the cost by %s. With the Metropolis rule, what is the probability that the move is accepted? Round to 4 decimal places.", edakit.Num(T, 0), edakit.Num(d, 0)), ans,
			edakit.Sprintf("An uphill move is accepted with probability $e^{-\\Delta/T}$ = exp(-%s/%s) = %s. A downhill move would always be accepted.", edakit.Num(d, 0), edakit.Num(T, 0), ans), map[string]string{"p": ans})
	case 3:
		T0 := float64(edakit.Pick([]int{100, 200, 500, 1000, 2000}))
		Tf := edakit.Pick([]float64{0.1, 0.5, 1, 2})
		r := edakit.Pick([]float64{0.8, 0.85, 0.9, 0.95})
		k, T := 0, T0
		for T >= Tf {
			T *= r
			k++
		}
		return edakit.Numeric(edakit.Sprintf("An annealing schedule starts at $T_0$ = %s and cools geometrically, $T_{k+1} = %s\\, T_k$, after each temperature step. It stops at the first step where T falls below %s. How many cooling steps does it take?", edakit.Num(T0, 0), edakit.Num(r, 2), edakit.Num(Tf, 1)), edakit.IntStr(k),
			edakit.Sprintf("We need the smallest k with %s x %s^k < %s, which is k > ln(%s/%s)/ln(%s) = %s, so k = %d.", edakit.Num(T0, 0), edakit.Num(r, 2), edakit.Num(Tf, 1), edakit.Num(Tf, 1), edakit.Num(T0, 0), edakit.Num(r, 2), edakit.Num(math.Log(Tf/T0)/math.Log(r), 3), k),
			map[string]string{"k": edakit.IntStr(k)})
	}
	a1, w1 := edakit.Int(400, 1200), edakit.Int(100, 600)
	a2, w2 := a1+edakit.Int(-120, 120), w1+edakit.Int(-80, 80)
	c1 := alpha*float64(a1) + beta*float64(w1)
	c2 := alpha*float64(a2) + beta*float64(w2)
	ans := edakit.Num(c2-c1, 1)
	return edakit.Numeric(edakit.Sprintf("A floorplanner's cost is $%s A + %s W$, where A is the bounding-box area in um^2 and W the total HPWL in um. The current floorplan has A = %d and W = %d. A proposed move gives A = %d and W = %d. What is the change in cost, new minus old? Round to 1 decimal place.", edakit.Num(alpha, 1), edakit.Num(beta, 1), a1, w1, a2, w2), ans,
		edakit.Sprintf("Old cost %s, new cost %s, change %s. A negative change is always accepted.", edakit.Num(c1, 1), edakit.Num(c2, 1), ans),
		map[string]string{"old": edakit.Num(c1, 1), "new": edakit.Num(c2, 1)})
}
