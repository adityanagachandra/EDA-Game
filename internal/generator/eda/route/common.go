package route

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// kp returns the requested knowledge point, or a uniform one for an unguided request.
func kp(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

type pt struct{ x, y int }

var dirs4 = []pt{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

func absI(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func minI(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxI(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func manh(a, b pt) int { return absI(a.x-b.x) + absI(a.y-b.y) }

// grid is a routing grid with blocked cells, indexed (column, row) from (1, 1).
type grid struct {
	w, h    int
	blocked map[pt]bool
}

func (g grid) free(p pt) bool {
	return p.x >= 1 && p.y >= 1 && p.x <= g.w && p.y <= g.h && !g.blocked[p]
}

// bfs returns the Lee wavefront label of every reachable cell.
func (g grid) bfs(s pt) map[pt]int {
	dist := map[pt]int{s: 0}
	q := []pt{s}
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, d := range dirs4 {
			v := pt{u.x + d.x, u.y + d.y}
			if !g.free(v) {
				continue
			}
			if _, ok := dist[v]; !ok {
				dist[v] = dist[u] + 1
				q = append(q, v)
			}
		}
	}
	return dist
}

// minBends returns the fewest direction changes on any route from s to t (or -1), by 0-1 BFS
// over (cell, direction) states.
func (g grid) minBends(s, t pt) int {
	type st struct {
		p pt
		d int
	}
	const inf = 1 << 30
	best := map[st]int{}
	var dq []st
	for d := range dirs4 {
		best[st{s, d}] = 0
		dq = append(dq, st{s, d})
	}
	get := func(k st) int {
		if v, ok := best[k]; ok {
			return v
		}
		return inf
	}
	for len(dq) > 0 {
		u := dq[0]
		dq = dq[1:]
		cu := get(u)
		for d, dv := range dirs4 {
			v := pt{u.p.x + dv.x, u.p.y + dv.y}
			if !g.free(v) {
				continue
			}
			c := cu
			if d != u.d {
				c++
			}
			k := st{v, d}
			if c < get(k) {
				best[k] = c
				if d == u.d {
					dq = append([]st{k}, dq...)
				} else {
					dq = append(dq, k)
				}
			}
		}
	}
	ans := inf
	for d := range dirs4 {
		if v := get(st{t, d}); v < ans {
			ans = v
		}
	}
	if ans == inf {
		return -1
	}
	return ans
}

// randGrid draws a grid whose blockages force the shortest source-target route to detour.
func randGrid(ctx generator.GeneratorContext) (grid, pt, pt, map[pt]int) {
	for try := 0; try < generator.MaxDrawAttempts; try++ {
		w, h := edakit.Scaled(ctx, 6, 8), edakit.Scaled(ctx, 5, 7)
		g := grid{w: w, h: h, blocked: map[pt]bool{}}
		s := pt{edakit.Int(1, 2), edakit.Int(1, h)}
		t := pt{edakit.Int(w-1, w), edakit.Int(1, h)}
		wx := edakit.Int(3, w-2)
		gap := edakit.Pick([]int{1, h})
		if edakit.Int(0, 2) == 0 {
			gap = edakit.Int(1, h)
		}
		for y := 1; y <= h; y++ {
			if y != gap {
				g.blocked[pt{wx, y}] = true
			}
		}
		extra := edakit.Scaled(ctx, 2, 6)
		for i := 0; i < extra; i++ {
			p := pt{edakit.Int(2, w-1), edakit.Int(1, h)}
			if p != s && p != t {
				g.blocked[p] = true
			}
		}
		delete(g.blocked, s)
		delete(g.blocked, t)
		dist := g.bfs(s)
		if d, ok := dist[t]; ok && d > manh(s, t) && d <= manh(s, t)+10 {
			return g, s, t, dist
		}
	}
	// Known-good fallback: a wall in column 3 with its gap in the top row.
	g := grid{w: 6, h: 5, blocked: map[pt]bool{{3, 1}: true, {3, 2}: true, {3, 3}: true, {3, 4}: true, {5, 2}: true}}
	s, t := pt{1, 1}, pt{6, 1}
	return g, s, t, g.bfs(s)
}

func sortPts(ps []pt) {
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].x != ps[j].x {
			return ps[i].x < ps[j].x
		}
		return ps[i].y < ps[j].y
	})
}

func ptText(p pt) string { return edakit.Sprintf("(%d, %d)", p.x, p.y) }

func blockedText(g grid) string {
	var ps []pt
	for p := range g.blocked {
		ps = append(ps, p)
	}
	sortPts(ps)
	var parts []string
	for _, p := range ps {
		parts = append(parts, ptText(p))
	}
	return strings.Join(parts, ", ")
}

func gridIntro(g grid, s, t pt) string {
	return edakit.Sprintf("A routing grid has columns 1 to %d and rows 1 to %d, and cells are written (column, row). Blocked cells: %s. The source pin is at %s and the target pin at %s. A route moves between the 4 neighbours of a cell. ",
		g.w, g.h, blockedText(g), ptText(s), ptText(t))
}

func intsText(xs []int) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = edakit.IntStr(x)
	}
	return strings.Join(parts, ", ")
}

func ceilDiv(a, b int) int { return (a + b - 1) / b }

// plural writes "1 cell" or "3 cells".
func plural(n int, word string) string {
	if n == 1 {
		return edakit.Sprintf("%d %s", n, word)
	}
	return edakit.Sprintf("%d %ss", n, word)
}
