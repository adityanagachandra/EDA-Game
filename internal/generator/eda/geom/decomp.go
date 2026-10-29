package geom

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// shape is a rectilinear region rasterised on the grid compressed from its defining
// rectangles: cell (i, j) spans [xs[i], xs[i+1]] x [ys[j], ys[j+1]].
type shape struct {
	xs, ys []int
	fill   [][]bool // fill[i][j]
}

func (s shape) at(i, j int) bool {
	if i < 0 || j < 0 || i >= len(s.xs)-1 || j >= len(s.ys)-1 {
		return false
	}
	return s.fill[i][j]
}

// buildShape rasterises the union of adds minus the union of subs.
func buildShape(adds, subs []rect) shape {
	var xs, ys []int
	for _, r := range append(append([]rect(nil), adds...), subs...) {
		xs = append(xs, r.x1, r.x2)
		ys = append(ys, r.y1, r.y2)
	}
	s := shape{xs: uniq(xs), ys: uniq(ys)}
	in := func(rs []rect, i, j int) bool {
		for _, r := range rs {
			if r.x1 <= s.xs[i] && s.xs[i+1] <= r.x2 && r.y1 <= s.ys[j] && s.ys[j+1] <= r.y2 {
				return true
			}
		}
		return false
	}
	s.fill = make([][]bool, len(s.xs)-1)
	for i := range s.fill {
		s.fill[i] = make([]bool, len(s.ys)-1)
		for j := range s.fill[i] {
			s.fill[i][j] = in(adds, i, j) && !in(subs, i, j)
		}
	}
	return s.trim()
}

// trim drops grid lines that no boundary uses, so every grid line carries a vertex.
func (s shape) trim() shape {
	nx, ny := len(s.xs)-1, len(s.ys)-1
	keepX := []int{s.xs[0]}
	colKeep := []int{0}
	for i := 1; i < nx; i++ {
		same := true
		for j := 0; j < ny; j++ {
			if s.fill[i][j] != s.fill[i-1][j] {
				same = false
			}
		}
		if !same {
			keepX = append(keepX, s.xs[i])
			colKeep = append(colKeep, i)
		}
	}
	keepX = append(keepX, s.xs[nx])
	t := shape{xs: keepX, ys: s.ys}
	t.fill = make([][]bool, len(colKeep))
	for k, i := range colKeep {
		t.fill[k] = s.fill[i]
	}
	// rows
	nx = len(t.xs) - 1
	keepY := []int{t.ys[0]}
	rowKeep := []int{0}
	for j := 1; j < ny; j++ {
		same := true
		for i := 0; i < nx; i++ {
			if t.fill[i][j] != t.fill[i][j-1] {
				same = false
			}
		}
		if !same {
			keepY = append(keepY, t.ys[j])
			rowKeep = append(rowKeep, j)
		}
	}
	keepY = append(keepY, t.ys[ny])
	u := shape{xs: t.xs, ys: keepY, fill: make([][]bool, nx)}
	for i := 0; i < nx; i++ {
		for _, j := range rowKeep {
			u.fill[i] = append(u.fill[i], t.fill[i][j])
		}
	}
	return u
}

func (s shape) transpose() shape {
	t := shape{xs: s.ys, ys: s.xs, fill: make([][]bool, len(s.ys)-1)}
	for j := range t.fill {
		t.fill[j] = make([]bool, len(s.xs)-1)
		for i := range t.fill[j] {
			t.fill[j][i] = s.fill[i][j]
		}
	}
	return t
}

// corner classifies grid point (i, j) by its four surrounding cells.
func (s shape) corner(i, j int) (filled int, diagonal bool) {
	a, b, c, d := s.at(i-1, j-1), s.at(i, j-1), s.at(i-1, j), s.at(i, j)
	for _, x := range []bool{a, b, c, d} {
		if x {
			filled++
		}
	}
	return filled, filled == 2 && a == d
}

type analysis struct {
	ok       bool
	vertices int
	reflex   int
	holes    int
	hChords  int
	vChords  int
	minRects int
	rings    [][]pt
	slabH    int
	slabV    int
}

func (s shape) slabs() int {
	n := 0
	ny := len(s.ys) - 1
	var prev [][2]int
	for j := 0; j < ny; j++ {
		var runs [][2]int
		for i := 0; i < len(s.xs)-1; i++ {
			if s.at(i, j) && !s.at(i-1, j) {
				k := i
				for s.at(k+1, j) {
					k++
				}
				runs = append(runs, [2]int{i, k})
			}
		}
		for _, r := range runs {
			found := false
			for _, p := range prev {
				if p == r {
					found = true
				}
			}
			if !found {
				n++
			}
		}
		prev = runs
	}
	return n
}

// analyse validates the shape (connected, no diagonal touches, no collinear chord chains) and
// computes its vertex counts, chords, minimum rectangle partition and slab counts.
func (s shape) analyse() analysis {
	var A analysis
	nx, ny := len(s.xs)-1, len(s.ys)-1
	// connectivity of filled cells and count of enclosed empty components
	seen := make([][]bool, nx+2)
	for i := range seen {
		seen[i] = make([]bool, ny+2)
	}
	get := func(i, j int) bool { return s.at(i-1, j-1) }
	var flood func(i, j int, want bool)
	flood = func(i, j int, want bool) {
		if i < 0 || j < 0 || i >= nx+2 || j >= ny+2 || seen[i][j] || get(i, j) != want {
			return
		}
		seen[i][j] = true
		flood(i+1, j, want)
		flood(i-1, j, want)
		flood(i, j+1, want)
		flood(i, j-1, want)
	}
	comps := 0
	for i := 0; i < nx+2; i++ {
		for j := 0; j < ny+2; j++ {
			if !seen[i][j] && get(i, j) {
				comps++
				flood(i, j, true)
			}
		}
	}
	if comps != 1 {
		return A
	}
	flood(0, 0, false)
	for i := 0; i < nx+2; i++ {
		for j := 0; j < ny+2; j++ {
			if !seen[i][j] && !get(i, j) {
				A.holes++
				flood(i, j, false)
			}
		}
	}
	// vertices
	type gp struct{ i, j int }
	isReflex := map[gp]bool{}
	for i := 0; i <= nx; i++ {
		for j := 0; j <= ny; j++ {
			f, diag := s.corner(i, j)
			if diag {
				return A
			}
			if f == 1 || f == 3 {
				A.vertices++
			}
			if f == 3 {
				A.reflex++
				isReflex[gp{i, j}] = true
			}
		}
	}
	// chords between consecutive reflex vertices on a grid line, through the interior
	type chord struct{ a, b gp }
	var hc, vc []chord
	ends := map[gp]int{}
	for j := 0; j <= ny; j++ {
		var on []int
		for i := 0; i <= nx; i++ {
			if isReflex[gp{i, j}] {
				on = append(on, i)
			}
		}
		for k := 0; k+1 < len(on); k++ {
			ok := true
			for i := on[k]; i < on[k+1]; i++ {
				if !s.at(i, j-1) || !s.at(i, j) {
					ok = false
				}
			}
			if ok {
				hc = append(hc, chord{gp{on[k], j}, gp{on[k+1], j}})
				ends[gp{on[k], j}]++
				ends[gp{on[k+1], j}]++
			}
		}
	}
	for _, c := range ends {
		if c > 1 {
			return A
		}
	}
	ends = map[gp]int{}
	for i := 0; i <= nx; i++ {
		var on []int
		for j := 0; j <= ny; j++ {
			if isReflex[gp{i, j}] {
				on = append(on, j)
			}
		}
		for k := 0; k+1 < len(on); k++ {
			ok := true
			for j := on[k]; j < on[k+1]; j++ {
				if !s.at(i-1, j) || !s.at(i, j) {
					ok = false
				}
			}
			if ok {
				vc = append(vc, chord{gp{i, on[k]}, gp{i, on[k+1]}})
				ends[gp{i, on[k]}]++
				ends[gp{i, on[k+1]}]++
			}
		}
	}
	for _, c := range ends {
		if c > 1 {
			return A
		}
	}
	A.hChords, A.vChords = len(hc), len(vc)
	// bipartite intersection graph and maximum matching (Kuhn)
	adj := make([][]int, len(hc))
	for a, h := range hc {
		for b, v := range vc {
			if h.a.i <= v.a.i && v.a.i <= h.b.i && v.a.j <= h.a.j && h.a.j <= v.b.j {
				adj[a] = append(adj[a], b)
			}
		}
	}
	match := make([]int, len(vc))
	for i := range match {
		match[i] = -1
	}
	var try func(a int, vis []bool) bool
	try = func(a int, vis []bool) bool {
		for _, b := range adj[a] {
			if vis[b] {
				continue
			}
			vis[b] = true
			if match[b] < 0 || try(match[b], vis) {
				match[b] = a
				return true
			}
		}
		return false
	}
	m := 0
	for a := range hc {
		if try(a, make([]bool, len(vc))) {
			m++
		}
	}
	g := len(hc) + len(vc) - m
	A.minRects = A.reflex - g - A.holes + 1
	A.rings = s.rings()
	A.slabH = s.slabs()
	A.slabV = s.transpose().slabs()
	A.ok = true
	return A
}

// rings traces the boundary with the interior on the left: the outer ring comes out
// counter-clockwise and holes clockwise. Each ring starts at its lowest-leftmost vertex.
func (s shape) rings() [][]pt {
	type gp struct{ i, j int }
	next := map[gp]gp{}
	nx, ny := len(s.xs)-1, len(s.ys)-1
	for i := 0; i < nx; i++ {
		for j := 0; j < ny; j++ {
			if !s.fill[i][j] {
				continue
			}
			if !s.at(i, j-1) {
				next[gp{i, j}] = gp{i + 1, j}
			}
			if !s.at(i+1, j) {
				next[gp{i + 1, j}] = gp{i + 1, j + 1}
			}
			if !s.at(i, j+1) {
				next[gp{i + 1, j + 1}] = gp{i, j + 1}
			}
			if !s.at(i-1, j) {
				next[gp{i, j + 1}] = gp{i, j}
			}
		}
	}
	var keys []gp
	for k := range next {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(a, b int) bool {
		if keys[a].j != keys[b].j {
			return keys[a].j < keys[b].j
		}
		return keys[a].i < keys[b].i
	})
	used := map[gp]bool{}
	var out [][]pt
	for _, k := range keys {
		if used[k] {
			continue
		}
		var loop []gp
		for c := k; !used[c]; c = next[c] {
			used[c] = true
			loop = append(loop, c)
		}
		var ring []pt
		for idx, c := range loop {
			p, n := loop[(idx+len(loop)-1)%len(loop)], loop[(idx+1)%len(loop)]
			if (p.i == c.i && c.i == n.i) || (p.j == c.j && c.j == n.j) {
				continue
			}
			ring = append(ring, pt{s.xs[c.i], s.ys[c.j]})
		}
		// rotate to lowest-leftmost
		best := 0
		for idx, p := range ring {
			if p.y < ring[best].y || (p.y == ring[best].y && p.x < ring[best].x) {
				best = idx
			}
		}
		ring = append(ring[best:], ring[:best]...)
		out = append(out, ring)
	}
	return out
}

func (A analysis) text() string {
	var b strings.Builder
	for k, r := range A.rings {
		if k == 0 {
			b.WriteString("outer boundary " + ringText(r))
		} else {
			b.WriteString(edakit.Sprintf("; hole %d: %s", k, ringText(r)))
		}
	}
	return b.String()
}

// cuts draws k sorted distinct integers strictly between lo and hi.
func cuts(k, lo, hi int) []int {
	set := map[int]bool{}
	for len(set) < k {
		set[edakit.Int(lo+1, hi-1)] = true
	}
	var out []int
	for v := range set {
		out = append(out, v)
	}
	sort.Ints(out)
	return out
}

// family draws a named rectilinear shape. withHoles selects the families that carry holes.
func family(withHoles bool) (string, []rect, []rect) {
	W, H := edakit.Int(10, 24), edakit.Int(10, 24)
	if withHoles {
		switch edakit.Int(0, 3) {
		case 0:
			a, b := cuts(2, 0, W), cuts(2, 0, H)
			return "frame", []rect{{0, 0, W, H}}, []rect{{a[0], b[0], a[1], b[1]}}
		case 1:
			a, b := cuts(4, 0, W), cuts(2, 0, H)
			return "two aligned holes", []rect{{0, 0, W, H}}, []rect{{a[0], b[0], a[1], b[1]}, {a[2], b[0], a[3], b[1]}}
		case 2:
			a, b := cuts(4, 0, W), cuts(4, 0, H)
			bb := edakit.Shuffle(b)
			return "two holes", []rect{{0, 0, W, H}}, []rect{{a[0], imin(bb[0], bb[1]), a[1], imax(bb[0], bb[1])}, {a[2], imin(bb[2], bb[3]), a[3], imax(bb[2], bb[3])}}
		default:
			adds := []rect{randRect(20, 8, 16), randRect(20, 8, 16)}
			m := mbr(adds)
			hx, hy := cuts(2, m.x1, m.x2), cuts(2, m.y1, m.y2)
			return "random with hole", adds, []rect{{hx[0], hy[0], hx[1], hy[1]}}
		}
	}
	switch edakit.Int(0, 8) {
	case 0:
		a, b := cuts(1, 0, W)[0], cuts(1, 0, H)[0]
		return "L", []rect{{0, 0, W, H}}, []rect{{a, b, W, H}}
	case 1:
		a, b := cuts(2, 0, W), cuts(1, 0, H)[0]
		return "T", []rect{{0, b, W, H}, {a[0], 0, a[1], b}}, nil
	case 2:
		a, b := cuts(2, 0, W), cuts(1, 0, H)[0]
		return "U", []rect{{0, 0, W, H}}, []rect{{a[0], b, a[1], H}}
	case 3:
		a, b := cuts(2, 0, W), cuts(2, 0, H)
		return "plus", []rect{{a[0], 0, a[1], H}, {0, b[0], W, b[1]}}, nil
	case 4:
		a, b := cuts(2, 0, W), cuts(2, 0, H)
		return "H", []rect{{0, 0, a[0], H}, {a[1], 0, W, H}, {a[0], b[0], a[1], b[1]}}, nil
	case 5:
		t := edakit.Int(3, 4)
		sp := cuts(1, 0, W)[0]
		ys := append(append([]int{0}, cuts(2*t-2, 0, H+6)...), H+6)
		adds := []rect{{0, 0, sp, H + 6}}
		for k := 0; k < t; k++ {
			adds = append(adds, rect{sp, ys[2*k], W, ys[2*k+1]})
		}
		return "comb", adds, nil
	case 6:
		k := edakit.Int(3, 5)
		xs := append(append([]int{0}, cuts(k-1, 0, W+6)...), W+6)
		hs := cuts(k, 0, H+8)
		var adds []rect
		for i := 0; i < k; i++ {
			adds = append(adds, rect{xs[i], 0, xs[i+1], hs[k-1-i]})
		}
		return "staircase", adds, nil
	default:
		return "random", []rect{randRect(20, 4, 14), randRect(20, 4, 14), randRect(20, 4, 14)}, nil
	}
}

func drawShape(withHoles bool) (string, analysis) {
	for {
		name, adds, subs := family(withHoles)
		A := buildShape(adds, subs).analyse()
		if !A.ok || (withHoles && A.holes == 0) || (!withHoles && A.holes > 0) || A.reflex == 0 {
			continue
		}
		return name, A
	}
}

// decompRectangles: reflex count (KP 1), minimum partition of a simple polygon (KP 2), minimum
// partition with holes (KP 3).
func decompRectangles(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		_, A := drawShape(false)
		q := edakit.Sprintf("A rectilinear polygon has vertices %s. How many of its vertices are reflex (interior angle 270 degrees)?", ringText(A.rings[0]))
		return edakit.Numeric(q, edakit.IntStr(A.reflex),
			edakit.Sprintf("The polygon has %d vertices. Walking counter-clockwise, every left turn is convex and every right turn reflex, and convex minus reflex is 4, so reflex = (n - 4)/2 = %d.", A.vertices, A.reflex),
			map[string]string{"vertices": edakit.IntStr(A.vertices), "reflex": edakit.IntStr(A.reflex)})
	case 2:
		_, A := drawShape(false)
		q := edakit.Sprintf("A rectilinear polygon has vertices %s. What is the minimum number of non-overlapping rectangles that partition it exactly?", ringText(A.rings[0]))
		g := A.reflex + 1 - A.minRects
		return edakit.Numeric(q, edakit.IntStr(A.minRects),
			edakit.Sprintf("It has r = %d reflex vertices. Chords join pairs of reflex vertices through the interior (%d horizontal, %d vertical); the largest set of pairwise non-touching chords has g = %d. The minimum is r - g + 1 = %d.", A.reflex, A.hChords, A.vChords, g, A.minRects),
			map[string]string{"reflex": edakit.IntStr(A.reflex), "g": edakit.IntStr(g)})
	default:
		_, A := drawShape(true)
		g := A.reflex - A.holes + 1 - A.minRects
		q := edakit.Sprintf("A layout polygon with holes has %s. What is the minimum number of non-overlapping rectangles that partition it exactly?", A.text())
		return edakit.Numeric(q, edakit.IntStr(A.minRects),
			edakit.Sprintf("Count reflex vertices on all rings: r = %d (each hole corner is reflex for the region). Non-touching chords: g = %d. With h = %d holes the minimum is r - g - h + 1 = %d.", A.reflex, g, A.holes, A.minRects),
			map[string]string{"reflex": edakit.IntStr(A.reflex), "g": edakit.IntStr(g), "holes": edakit.IntStr(A.holes)})
	}
}

// decompSlab: horizontal slab count (KP 1), its excess over the minimum (KP 2), the better of
// the two slab directions (KP 3).
func decompSlab(ctx generator.GeneratorContext) generator.Problem {
	withHoles := edakit.Int(0, 3) == 0
	_, A := drawShape(withHoles)
	desc := A.text()
	switch kp(ctx) {
	case 1:
		q := edakit.Sprintf("A fracturing step cuts a polygon with horizontal lines extended from every reflex vertex across the interior, and keeps the pieces. The polygon has %s. How many rectangles result?", desc)
		return edakit.Numeric(q, edakit.IntStr(A.slabH),
			edakit.Sprintf("Between consecutive vertex y-values the polygon is a set of horizontal runs; stacking identical runs gives %d rectangles.", A.slabH),
			map[string]string{"slab": edakit.IntStr(A.slabH)})
	case 2:
		q := edakit.Sprintf("A polygon has %s. How many more rectangles does the horizontal slab decomposition (cuts from every reflex vertex, horizontal only) produce than the minimum partition?", desc)
		return edakit.Numeric(q, edakit.IntStr(A.slabH-A.minRects),
			edakit.Sprintf("Horizontal slabs give %d rectangles. The minimum, r - g - h + 1, is %d. The slab method is simple and sweep-friendly but ignores vertical chords that would save pieces: %d extra.", A.slabH, A.minRects, A.slabH-A.minRects),
			map[string]string{"slab": edakit.IntStr(A.slabH), "min": edakit.IntStr(A.minRects)})
	default:
		best := imin(A.slabH, A.slabV)
		q := edakit.Sprintf("A mask fracturer may slab a polygon with horizontal cuts or with vertical cuts from every reflex vertex, and keeps the direction giving fewer rectangles. The polygon has %s. How many rectangles does it keep?", desc)
		return edakit.Numeric(q, edakit.IntStr(best),
			edakit.Sprintf("Horizontal slabs give %d rectangles and vertical slabs %d, so it keeps %d. The true minimum is %d.", A.slabH, A.slabV, best, A.minRects),
			map[string]string{"slab_h": edakit.IntStr(A.slabH), "slab_v": edakit.IntStr(A.slabV)})
	}
}
