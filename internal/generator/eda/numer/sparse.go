package numer

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// netlist draws a resistor netlist on n non-ground nodes: a connected set of node-to-node
// resistors, some resistors to ground, and sometimes a parallel duplicate.
func netlist(n int) (*graph, []string, bool) {
	g := randomGraph(n, 0.25)
	var lines []string
	k := 1
	for _, e := range g.edges() {
		var u, v int
		_, _ = sscan2(e, &u, &v)
		lines = append(lines, sp("R%d between nodes %d and %d", k, u, v))
		k++
	}
	dup := false
	if edakit.Int(0, 1) == 1 {
		es := g.edges()
		var u, v int
		_, _ = sscan2(edakit.Pick(es), &u, &v)
		lines = append(lines, sp("R%d between nodes %d and %d", k, v, u))
		k++
		dup = true
	}
	for _, v := range edakit.Shuffle(seq(1, n))[:edakit.Int(1, 3)] {
		lines = append(lines, sp("R%d from node %d to ground", k, v))
		k++
	}
	return g, edakit.Shuffle(lines), dup
}

// sscan2 parses "u-v".
func sscan2(e string, u, v *int) (int, error) {
	parts := strings.Split(e, "-")
	*u, *v = atoi(parts[0]), atoi(parts[1])
	return 2, nil
}

func atoi(s string) int {
	n := 0
	for _, r := range s {
		n = n*10 + int(r-'0')
	}
	return n
}

func mnaPattern(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 7)
		g, lines, _ := netlist(n)
		v := edakit.Int(1, n)
		ans := 1 + len(g.adj[v])
		q := sp("A resistor network has %d non-ground nodes: %s. In the nodal conductance matrix $G$ (one row and column per non-ground node), how many nonzero entries does row %d have?",
			n, strings.Join(lines, "; "), v)
		return numeric(q, itoa(ans), sp("Row %d has its diagonal entry plus one off-diagonal entry for each distinct non-ground neighbour of node %d (%d of them). Resistors to ground only add to the diagonal, and parallel resistors share one entry: %d nonzeros.", v, v, len(g.adj[v]), ans),
			map[string]string{"neighbours": itoa(len(g.adj[v]))})
	case 2:
		n := edakit.Scaled(ctx, 4, 8)
		g, lines, _ := netlist(n)
		m := g.edgeCount()
		ans := n + 2*m
		q := sp("A resistor network has %d non-ground nodes: %s. How many nonzero entries does the full nodal conductance matrix $G$ have?",
			n, strings.Join(lines, "; "))
		return numeric(q, itoa(ans), sp("There are %d diagonal entries and %d distinct pairs of non-ground nodes joined by a resistor; each pair puts a symmetric pair of off-diagonal entries in $G$. Total %d + 2 · %d = %d.", n, m, n, m, ans),
			map[string]string{"pairs": itoa(m), "n": itoa(n)})
	default:
		k := edakit.Scaled(ctx, 10, 60)
		n := k * k
		edges := 2 * k * (k - 1)
		ans := n + 2*edges
		q := sp("A power-grid mesh is a %d-by-%d array of nodes; each node connects through a resistor to its left, right, upper and lower neighbours (where they exist), and some nodes connect to supply pads (to ground in the nodal formulation). How many nonzero entries does the %d-by-%d nodal conductance matrix have?", k, k, n, n)
		return numeric(q, itoa(ans), sp("The grid has %d horizontal plus vertical resistors (2 · %d · %d). The matrix has %d diagonal entries plus 2 · %d off-diagonal ones: %d, under 5 per row against %d dense entries per row.", edges, k, k-1, n, edges, ans, n),
			map[string]string{"edges": itoa(edges)})
	}
}

// sparseMatrix draws an n-by-n matrix with a nonzero diagonal and random off-diagonal
// nonzeros (small integers).
func sparseMatrix(n int, density float64) [][]float64 {
	a := make([][]float64, n)
	for i := range a {
		a[i] = make([]float64, n)
		for j := range a[i] {
			if i == j {
				a[i][j] = float64(edakit.Int(2, 9))
			} else if randFloat() < density {
				a[i][j] = float64(edakit.Pick([]int{-4, -3, -2, -1, 1, 2, 3, 5}))
			}
		}
	}
	return a
}

func nnzOf(a [][]float64) int {
	c := 0
	for _, r := range a {
		for _, v := range r {
			if v != 0 {
				c++
			}
		}
	}
	return c
}

// csrOf converts a dense matrix to CSR arrays (0-based).
func csrOf(a [][]float64) (vals []float64, cols []int, rowPtr []int) {
	rowPtr = []int{0}
	for _, r := range a {
		for j, v := range r {
			if v != 0 {
				vals = append(vals, v)
				cols = append(cols, j)
			}
		}
		rowPtr = append(rowPtr, len(vals))
	}
	return
}

func floatsStr(xs []float64) string {
	var s []string
	for _, x := range xs {
		s = append(s, f(x, 3))
	}
	return strings.Join(s, ", ")
}

func storage(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 6)
		a := sparseMatrix(n, 0.3)
		nnz := nnzOf(a)
		if edakit.Int(0, 1) == 0 {
			ans := 2*nnz + n + 1
			q := "Store the matrix " + pmat(a, 0) + " in compressed sparse row (CSR) format with arrays val, col_idx and row_ptr. What is the total number of array entries across the three arrays?"
			return numeric(q, itoa(ans), sp("The matrix has %d nonzeros, so val and col_idx each hold %d entries and row_ptr holds n + 1 = %d: total %d.", nnz, nnz, n+1, ans),
				map[string]string{"nnz": itoa(nnz)})
		}
		ans := 3 * nnz
		q := "Store the matrix " + pmat(a, 0) + " in coordinate (COO) format as (row, column, value) triples. How many numbers are stored in total?"
		return numeric(q, itoa(ans), sp("COO keeps one triple per nonzero. There are %d nonzeros, so %d numbers.", nnz, ans),
			map[string]string{"nnz": itoa(nnz)})
	case 2:
		n := edakit.Scaled(ctx, 4, 6)
		a := sparseMatrix(n, 0.35)
		_, _, rp := csrOf(a)
		k := edakit.Int(2, n-1)
		q := "Rows are numbered from 0. In CSR storage of " + pmat(a, 0) + sp(", the array row_ptr starts with row_ptr[0] = 0. What is row_ptr[%d]?", k)
		return numeric(q, itoa(rp[k]), sp("row_ptr[%d] is the number of nonzeros in rows 0 to %d, the position in val where row %d starts: %d.", k, k-1, k, rp[k]),
			map[string]string{"row_ptr": ints(rp)})
	default:
		n := 1000 * edakit.Scaled(ctx, 100, 5000)
		d := edakit.Int(3, 7)
		nnz := n * d
		if edakit.Int(0, 1) == 0 {
			bytes := 12*float64(nnz) + 4*float64(n+1)
			ans := f(bytes/1e6, 2)
			q := sp("A nodal matrix has n = %d rows and an average of %d nonzeros per row. Values are 8-byte doubles and indices are 4-byte integers. How much memory does CSR storage take, in megabytes ($10^6$ bytes)? Round to 2 decimal places.", n, d)
			return numeric(q, ans, sp("nnz = %d. CSR stores 8 + 4 = 12 bytes per nonzero plus 4(n + 1) bytes of row pointers: %s MB.", nnz, ans),
				map[string]string{"nnz": itoa(nnz)})
		}
		dense := 8 * float64(n) * float64(n)
		csr := 12*float64(nnz) + 4*float64(n+1)
		ans := f(dense/csr, 0)
		q := sp("A nodal matrix has n = %d rows and an average of %d nonzeros per row. With 8-byte values and 4-byte indices, how many times more memory does dense storage need than CSR? Round to the nearest integer.", n, d)
		return numeric(q, ans, sp("Dense: 8n^2 bytes. CSR: 12 nnz + 4(n + 1) bytes with nnz = %d. The ratio is %s.", nnz, ans),
			map[string]string{"nnz": itoa(nnz)})
	}
}

func csrMatvec(ctx generator.GeneratorContext) generator.Problem {
	k := kpOf(ctx)
	if k == 3 {
		n := 1000 * edakit.Scaled(ctx, 50, 2000)
		d := edakit.Int(3, 9)
		nnz := n * d
		bytes := 12*float64(nnz) + 4*float64(n+1) + 16*float64(n)
		if edakit.Int(0, 1) == 0 {
			ans := f(bytes/1e6, 2)
			q := sp("A CSR matrix has n = %d rows and nnz = %d. One SpMV $y = Ax$ reads val (8 bytes per entry), col_idx (4 bytes per entry), row_ptr (4 bytes per entry, n + 1 entries), reads $x$ once (8n bytes) and writes $y$ (8n bytes). How many megabytes ($10^6$ bytes) move? Round to 2 decimal places.", n, nnz)
			return numeric(q, ans, sp("12 nnz + 4(n + 1) + 16n = %s MB.", ans), map[string]string{"bytes": f(bytes, 0)})
		}
		ans := itoa(2 * nnz)
		q := sp("A CSR matrix has n = %d rows and nnz = %d. How many floating-point operations (count multiplies and adds separately) does one SpMV $y = Ax$ perform?", n, nnz)
		return numeric(q, ans, "Each stored entry contributes one multiply and one add: 2 nnz = "+ans+".", nil)
	}
	n := 4
	if k == 2 {
		n = edakit.Scaled(ctx, 4, 5)
	}
	a := sparseMatrix(n, 0.35)
	vals, cols, rp := csrOf(a)
	if k == 1 {
		// Position of entry (i, j) in val.
		type pos struct{ i, j int }
		var ps []pos
		for i := range a {
			for j := range a[i] {
				if a[i][j] != 0 && rp[i] > 0 {
					ps = append(ps, pos{i, j})
				}
			}
		}
		p := edakit.Pick(ps)
		idx := rp[p.i]
		for j := 0; j < p.j; j++ {
			if a[p.i][j] != 0 {
				idx++
			}
		}
		q := "Rows and columns are numbered from 0. The matrix " + pmat(a, 0) + sp(" is stored in CSR with nonzeros of each row kept in increasing column order. At which position (0-based) of the val array is the entry in row %d, column %d?", p.i, p.j)
		return numeric(q, itoa(idx), sp("Row %d starts at row_ptr[%d] = %d. Count the nonzeros of row %d to the left of column %d and add: position %d.", p.i, p.i, rp[p.i], p.i, p.j, idx),
			map[string]string{"row_ptr": ints(rp)})
	}
	x := make([]float64, n)
	for i := range x {
		x[i] = float64(edakit.Int(-3, 4))
	}
	i := edakit.Int(0, n-1)
	y := 0.0
	for kk := rp[i]; kk < rp[i+1]; kk++ {
		y += vals[kk] * x[cols[kk]]
	}
	q := sp("A matrix is stored in CSR (0-based): val = [%s], col_idx = [%s], row_ptr = [%s]. With $x = (%s)$, what is $y_%d$ in $y = Ax$?",
		floatsStr(vals), ints(cols), ints(rp), floatsStr(x), i)
	return numeric(q, f(y, 3), sp("Row %d occupies positions %d to %d of val. Summing val[k] times x[col_idx[k]] over them gives %s.", i, rp[i], rp[i+1]-1, f(y, 3)),
		map[string]string{"y": f(y, 3)})
}

func fillIn(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 5, 7)
		g := randomGraph(n, 0.3)
		var cands []int
		for v := 1; v <= n; v++ {
			if len(g.adj[v]) >= 2 {
				cands = append(cands, v)
			}
		}
		v := edakit.Pick(cands)
		fill := g.clone().eliminate(v, map[int]bool{})
		q := sp("A symmetric sparse matrix has a nonzero diagonal and off-diagonal nonzeros at the vertex pairs %s (the graph of the matrix, vertices 1 to %d). Gaussian elimination eliminates vertex %d first. How many fill edges (new off-diagonal nonzero pairs) does that step create?",
			strings.Join(g.edges(), ", "), n, v)
		return numeric(q, itoa(fill), sp("Eliminating %d makes its neighbours a clique. Count the neighbour pairs that were not already joined: %d.", v, fill),
			map[string]string{"degree": itoa(len(g.adj[v]))})
	case 2:
		n := edakit.Scaled(ctx, 5, 7)
		g := randomGraph(n, 0.3)
		order := edakit.Shuffle(seq(1, n))
		fill := fillForOrder(g, order)
		q := sp("A symmetric sparse matrix has the graph with edges %s on vertices 1 to %d. Eliminate in the order %s. How many fill edges are created in total? (Each fill edge is one new nonzero in $L$ below the diagonal.)",
			strings.Join(g.edges(), ", "), n, ints(order))
		return numeric(q, itoa(fill), sp("Play the elimination game: each eliminated vertex joins its remaining neighbours into a clique. The new edges over the whole order total %d.", fill),
			map[string]string{"fill": itoa(fill)})
	default:
		leaves := edakit.Scaled(ctx, 4, 12)
		path := edakit.Int(0, 1) == 1
		fill := leaves * (leaves - 1) / 2
		desc := sp("a hub vertex coupled to each of %d other vertices, which are not coupled to each other", leaves)
		if path {
			fill -= leaves - 1
			desc = sp("a hub vertex coupled to each of %d other vertices, which in addition form a chain (each coupled to the next)", leaves)
		}
		q := sp("An arrow-shaped matrix has %s. If the hub is eliminated first, how many fill edges appear? (Eliminating the hub last creates none.)", desc)
		return numeric(q, itoa(fill), sp("Eliminating the hub first joins all %d neighbours into a clique of %d pairs; %s So the fill is %d.", leaves, leaves*(leaves-1)/2,
			map[bool]string{true: sp("%d of those pairs are chain edges already.", leaves-1), false: "none existed before."}[path], fill),
			map[string]string{"leaves": itoa(leaves)})
	}
}

// pattern is an unsymmetric 0/1 pattern with a full diagonal.
func randPattern(n int, density float64) [][]bool {
	p := make([][]bool, n)
	for i := range p {
		p[i] = make([]bool, n)
		for j := range p[i] {
			p[i][j] = i == j || randFloat() < density
		}
	}
	return p
}

func patternRows(p [][]bool) string {
	var rows []string
	for i, r := range p {
		var cs []int
		for j, b := range r {
			if b {
				cs = append(cs, j+1)
			}
		}
		rows = append(rows, sp("row %d: columns %s", i+1, ints(cs)))
	}
	return strings.Join(rows, "; ")
}

func rowCount(p [][]bool, i int) int {
	c := 0
	for _, b := range p[i] {
		if b {
			c++
		}
	}
	return c
}

func colCount(p [][]bool, j int) int {
	c := 0
	for i := range p {
		if p[i][j] {
			c++
		}
	}
	return c
}

func markowitz(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 6)
	p := randPattern(n, 0.35)
	type ent struct{ i, j, m int }
	var es []ent
	minM := 1 << 30
	for i := range p {
		for j := range p[i] {
			if p[i][j] {
				m := (rowCount(p, i) - 1) * (colCount(p, j) - 1)
				es = append(es, ent{i, j, m})
				if m < minM {
					minM = m
				}
			}
		}
	}
	stem := "An unsymmetric sparse matrix has this nonzero pattern: " + patternRows(p) + ". "
	switch kpOf(ctx) {
	case 1:
		e := edakit.Pick(es)
		q := stem + sp("What is the Markowitz count of the entry in row %d, column %d?", e.i+1, e.j+1)
		return numeric(q, itoa(e.m), sp("Row %d has %d nonzeros and column %d has %d, so the count is (%d - 1)(%d - 1) = %d.", e.i+1, rowCount(p, e.i), e.j+1, colCount(p, e.j), rowCount(p, e.i), colCount(p, e.j), e.m),
			map[string]string{"markowitz": itoa(e.m)})
	case 2:
		q := stem + "Over all nonzero entries, what is the smallest Markowitz count, the one a Markowitz pivot search would choose (ignoring stability)?"
		return numeric(q, itoa(minM), sp("Compute (r_i - 1)(c_j - 1) for every nonzero; the minimum is %d.", minM),
			map[string]string{"min": itoa(minM)})
	default:
		var cands []ent
		for _, e := range es {
			if e.m >= 2 {
				cands = append(cands, e)
			}
		}
		if len(cands) == 0 {
			for _, e := range es {
				if e.m > 0 {
					cands = append(cands, e)
				}
			}
		}
		if len(cands) == 0 {
			cands = es
		}
		e := edakit.Pick(cands)
		fill := 0
		for r := 0; r < n; r++ {
			if r == e.i || !p[r][e.j] {
				continue
			}
			for c := 0; c < n; c++ {
				if c == e.j || !p[e.i][c] {
					continue
				}
				if !p[r][c] {
					fill++
				}
			}
		}
		q := stem + sp("The entry in row %d, column %d has Markowitz count %d. If it is chosen as the first pivot, how many fill entries (zero positions that become nonzero) does that step actually create?", e.i+1, e.j+1, e.m)
		return numeric(q, itoa(fill), sp("The update touches the %d positions (r, c) with r in the pivot column and c in the pivot row. Of those, %d were zero, so the actual fill is %d, at most the Markowitz count.", e.m, fill, fill),
			map[string]string{"bound": itoa(e.m), "fill": itoa(fill)})
	}
}

func minimumDegree(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 6, 8)
	g := randomGraph(n, 0.3)
	stem := sp("A symmetric matrix has the graph with edges %s on vertices 1 to %d. ", strings.Join(g.edges(), ", "), n)
	order, fill := minDegreeOrder(g)
	switch kpOf(ctx) {
	case 1:
		idx := edakit.Int(0, 1)
		q := stem + sp("Minimum degree ordering always eliminates a vertex of least current degree, breaking ties by the lowest label. Which vertex is eliminated %s?", []string{"first", "second"}[idx])
		return numeric(q, itoa(order[idx]), sp("The minimum degree order is %s, so the answer is %d.", ints(order), order[idx]),
			map[string]string{"order": ints(order)})
	case 2:
		q := stem + "Run minimum degree ordering (ties broken by the lowest label), updating degrees after each elimination. How many fill edges does the whole ordering create?"
		return numeric(q, itoa(fill), sp("The order is %s and the elimination game along it creates %d fill edges.", ints(order), fill),
			map[string]string{"order": ints(order), "fill": itoa(fill)})
	default:
		var vs []int
		for v := 1; v <= n; v++ {
			if len(g.adj[v]) >= 2 {
				vs = append(vs, v)
			}
		}
		v := edakit.Pick(vs)
		var nb []int
		for u := range g.adj[v] {
			nb = append(nb, u)
		}
		sort.Ints(nb)
		u := edakit.Pick(nb)
		h := g.clone()
		removed := map[int]bool{}
		h.eliminate(v, removed)
		d := h.degree(u, removed)
		q := stem + sp("Vertex %d is eliminated. What is the degree of vertex %d in the resulting elimination graph?", v, u)
		return numeric(q, itoa(d), sp("After eliminating %d, vertex %d is adjacent to its old neighbours (minus %d) together with the other neighbours of %d: degree %d.", v, u, v, v, d),
			map[string]string{"degree": itoa(d)})
	}
}

var nestedDissectionItems = []edakit.Item{
	{KP: 3, Q: "For a 2D grid with n unknowns, how does the fill in the Cholesky factor scale under nested dissection ordering?",
		A: "O(n log n)", D: []string{"O(n)", "O(n^1.5)", "O(n^2)"},
		Why: "George's result for planar grids: nested dissection gives O(n log n) fill and O(n^1.5) factorization work, against O(n^1.5) fill for a band ordering."},
	{KP: 3, Q: "Why does a 3D thermal mesh factor much more expensively than a 2D power grid with the same number of unknowns?",
		A: "Its separators are planes of about n^(2/3) unknowns, so the dense trailing blocks are much larger", D: []string{"3D meshes cannot be ordered by nested dissection", "Cholesky needs pivoting in 3D", "3D matrices are not symmetric"},
		Why: "A 3D grid separator has size n^(2/3) rather than n^(1/2), and its dense block costs the cube of that size, which gives O(n^2) work for direct factorization."},
	{KP: 3, Q: "A tool orders a power-grid matrix with nested dissection rather than minimum degree. Besides fill, what does it gain?",
		A: "A short, bushy elimination tree whose independent subtrees can be factored in parallel", D: []string{"A guarantee of zero fill", "Freedom from numerical pivoting", "A smaller number of nonzeros in the original matrix"},
		Why: "Separators split the problem into independent halves, so the elimination tree branches early and both halves factor concurrently."},
	{KP: 3, Q: "Which library is the usual source of nested dissection orderings for large sparse matrices?",
		A: "METIS, through its multilevel graph partitioning", D: []string{"BLAS level 3", "LAPACK dgetrf", "The Markowitz search in Sparse 1.3"},
		Why: "Nested dissection needs good small separators, and METIS computes them by multilevel partitioning; solvers such as CHOLMOD and MUMPS call it."},
	{KP: 3, Q: "Circuit matrices from a typical analog netlist (not a mesh) are often ordered with AMD rather than nested dissection. Why?",
		A: "They have no small geometric separators and are small and very sparse, so minimum degree gives low fill cheaply", D: []string{"Nested dissection only works on symmetric positive definite matrices", "AMD guarantees optimal fill", "Nested dissection produces dense factors on every graph"},
		Why: "Nested dissection pays off on mesh-like graphs with small separators. Irregular circuit graphs with hubs (supply nets) suit degree-based orderings, as in KLU."},
	{KP: 3, Q: "In nested dissection the separator vertices are numbered last. What would go wrong if they were eliminated first?",
		A: "Eliminating a separator vertex joins neighbours on both sides, creating fill that couples the two halves", D: []string{"Nothing: elimination order does not affect fill", "The matrix would lose symmetry", "The factorization would need pivoting"},
		Why: "Keeping the separator for last means no elimination step ever has neighbours in both halves, so no fill crosses between them."},
}

func ndGrid(ctx generator.GeneratorContext) generator.Problem {
	k := 2*edakit.Scaled(ctx, 2, 10) + 1
	ans := k * (k - 1) / 2
	q := sp("A %d-by-%d grid graph (a five-point power-grid mesh) is split by nested dissection with its middle column as the separator. How many vertices does each of the two remaining halves contain?", k, k)
	return edakit.Choice(q, itoa(ans), sp("The separator is the middle column of %d vertices. The other %d vertices split evenly: %d per half.", k, k*k-k, ans),
		edakit.Distinct(ans, 3, k*k/2, k*(k+1)/2, (k-1)*(k-1)/2, k*k-k)...)
}

func ndDense(ctx generator.GeneratorContext) generator.Problem {
	if edakit.Int(0, 1) == 0 {
		k := edakit.Scaled(ctx, 8, 60)
		ans := k * (k + 1) / 2
		q := sp("In a nested dissection ordering of a %d-by-%d grid, the top-level separator has %d vertices and is eliminated last; its block of the factor becomes completely dense. How many nonzeros does that block of $L$ hold (lower triangle including the diagonal)?", k, k, k)
		return edakit.Choice(q, itoa(ans), sp("A dense %d-by-%d lower triangle holds %d(%d + 1)/2 = %d entries.", k, k, k, k, ans),
			edakit.Distinct(ans, 3, k*k, k*(k-1)/2, 2*k, k*k*k/3)...)
	}
	k := edakit.Scaled(ctx, 6, 20)
	s := k * k
	ans := s * (s + 1) / 2
	q := sp("In a nested dissection ordering of a %d-by-%d-by-%d 3D thermal grid, the top-level separator is a plane of %d vertices eliminated last, and its block of the factor becomes dense. How many nonzeros does that block of $L$ hold (lower triangle including the diagonal)?", k, k, k, s)
	return edakit.Choice(q, itoa(ans), sp("The separator has %d vertices; a dense lower triangle of that size holds %d(%d + 1)/2 = %d entries.", s, s, s, ans),
		edakit.Distinct(ans, 3, s*s, k*(k+1)/2, s*(s-1)/2, 2*s)...)
}
