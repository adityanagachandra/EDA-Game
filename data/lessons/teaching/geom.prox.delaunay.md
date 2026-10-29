# Delaunay triangulation

The Delaunay triangulation connects sites whose Voronoi cells share an edge. It is the triangulation that avoids skinny triangles where it can: it maximizes the minimum angle over all triangulations. That makes it the standard mesh for field solvers and finite-element models in extraction and thermal analysis. Its counts follow from Euler's formula, and its defining test, the empty circle, is a single determinant.

## Counting triangles and edges

Any triangulation of $n$ points in the plane, $h$ of them on the convex hull, has exactly $2n - 2 - h$ triangles and $3n - 3 - h$ edges. Euler's formula gives it: with $V = n$ and $F = T + 1$ faces counting the outer one, each triangle has 3 edges and each edge borders two faces, except the $h$ hull edges, so $3T + h = 2E$. Example: 40 points with 6 on the hull give 72 triangles and 111 edges, whether the triangulation is Delaunay or not. The Delaunay triangulation is the dual of the Voronoi diagram: one triangle per Voronoi vertex and one edge per Voronoi edge, which is why the counts match. Building it takes $O(n \log n)$ by divide and conquer or by randomized incremental insertion.

## The incircle test

A triangulation is Delaunay when no point lies strictly inside the circumcircle of any triangle. For a counter-clockwise triangle $abc$ and a point $d$, the test is the sign of

$$\det \begin{pmatrix} a_x - d_x & a_y - d_y & (a_x - d_x)^2 + (a_y - d_y)^2 \\ b_x - d_x & b_y - d_y & (b_x - d_x)^2 + (b_y - d_y)^2 \\ c_x - d_x & c_y - d_y & (c_x - d_x)^2 + (c_y - d_y)^2 \end{pmatrix},$$

positive when $d$ is inside, negative outside, zero on the circle. It is the orientation test lifted onto the paraboloid $z = x^2 + y^2$. Example: $a = (0, 0)$, $b = (4, 0)$, $c = (0, 4)$ and $d = (1, 1)$ give a determinant of 96, so $d$ is inside. Like orientation, it must be evaluated exactly, and on integer input it needs roughly twice the bits of the orientation test. Incremental Delaunay algorithms call it once per candidate edge flip, so its cost and exactness dominate their robustness.

## Refining a mesh

Field solvers refine meshes by inserting Steiner points where triangles are too large or too skinny. Inserting a point strictly inside the hull leaves the hull unchanged, so by the counting formula each insertion adds 2 triangles and 3 edges; locally, the new point splits its triangle into three, and edge flips driven by the incircle test restore the Delaunay property. Example: a mesh of 200 points with 12 on the hull has 386 triangles; adding 30 interior points gives 446. Mesh size drives solver cost, since the linear system has one unknown per node, so quality-driven refinement (Ruppert's algorithm, Chew's algorithm) that bounds angles while adding few points directly sets extraction run time and accuracy.

## Going deeper

de Berg et al., "Computational Geometry: Algorithms and Applications" (3rd ed., 2008), chapter 9; Shewchuk, "Delaunay refinement algorithms for triangular mesh generation" (Computational Geometry, 2002); Guibas and Stolfi, "Primitives for the manipulation of general subdivisions and the computation of Voronoi diagrams" (ACM Transactions on Graphics, 1985).
