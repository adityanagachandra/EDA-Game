# Nested dissection and separators

Nested dissection orders a mesh-like graph by cutting it with small separators and numbering each separator after the pieces it separates. On power grids and thermal meshes it gives far less fill than band orderings and produces an elimination tree with natural parallelism, which is why large IR-drop and thermal direct solvers rely on it.

## Separators split the elimination

A vertex separator $S$ splits a graph into parts $A$ and $B$ with no edges between them. Order $A$ first, then $B$, then $S$. While eliminating vertices of $A$, every neighbour lies in $A$ or $S$, so fill can join $A$ to $S$ but never to $B$. The factor therefore has a zero block coupling $A$ and $B$, and the two parts can be factored independently.

Recursing on $A$ and $B$ gives nested dissection. On a $k$-by-$k$ grid, the middle column is a separator of $k$ vertices; for $k = 7$ it leaves two halves of $7 \cdot 3 = 21$ vertices each. Each half is then cut by its middle row, and so on until pieces are tiny.

## Ordering a grid by dissection

The cost concentrates in the separators. When the top-level separator is eliminated last, every pair of its vertices is joined through the parts already eliminated, so its block of $L$ becomes dense: for a separator of $s$ vertices, $s(s + 1)/2$ entries including the diagonal. For a 31-by-31 grid, $s = 31$ and the block holds 496 entries.

Summing over levels, a 2D grid with $n$ unknowns has separators of size about $\sqrt n$, and George showed that nested dissection gives $O(n \log n)$ fill and $O(n^{3/2})$ factorization work. A band ordering on the same grid gives $O(n^{3/2})$ fill and $O(n^2)$ work. In 3D the separators are planes of about $n^{2/3}$ vertices, and the work rises to $O(n^2)$, which is why large 3D thermal problems often switch to iterative solvers.

## Asymptotic fill for grids and why it matters

Separators also shape the elimination tree. Parts $A$ and $B$ form independent subtrees, so the tree is short and bushy, and the subtrees can be factored on different cores or nodes. A band ordering produces a long chain with little parallelism.

Practical tools compute separators with multilevel graph partitioning: METIS is the standard, and solvers such as CHOLMOD, MUMPS and PARDISO call it or something like it. The choice between nested dissection and minimum degree depends on the matrix. Mesh-like matrices (power grids, thermal grids, substrate meshes) have small separators and benefit; irregular analog netlists with hub nodes do not, and AMD usually wins there.

## Going deeper

George, "Nested dissection of a regular finite element mesh", SIAM J. Numer. Anal. (1973); Karypis and Kumar, "A fast and high quality multilevel scheme for partitioning irregular graphs", SIAM J. Sci. Comput. (1998).
