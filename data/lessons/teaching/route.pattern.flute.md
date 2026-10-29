# FLUTE lookup-table Steiner trees

FLUTE (Chu and Wong) builds rectilinear Steiner minimal trees by table lookup. For low-degree nets, which dominate real designs, it is optimal and runs in time comparable to computing a bounding box. Nearly every academic and open-source global router, including FastRoute and OpenROAD's grt, starts its net topologies from FLUTE.

## Position sequences

Sort a net's $n$ pins by x coordinate, and record for each pin the rank of its y coordinate. The resulting permutation of $1, \dots, n$ is the net's position sequence. Two nets with the same sequence have the same relative arrangement of pins, so they share the same family of candidate optimal topologies, regardless of the actual distances.

There are $n!$ possible sequences: 24 for degree 4, 720 for degree 6, 40320 for degree 8. FLUTE precomputes a table entry for each sequence of low degree (the published tables cover degrees up to 9, with symmetries reducing storage). Larger nets are broken recursively into smaller ones that the table covers, at a small cost in optimality.

## Wirelength vectors

Let $h_1, \dots, h_{n-1}$ be the gaps between adjacent distinct x coordinates of the pins and $v_1, \dots, v_{n-1}$ the gaps between adjacent y coordinates. Any rectilinear Steiner tree on the Hanan grid has a length $\sum a_i h_i + \sum b_i v_i$ with small non-negative integer coefficients. FLUTE's table stores, for each sequence, only the coefficient vectors that can be optimal for some gaps: the potentially optimal wirelength vectors.

At query time FLUTE computes the gaps, evaluates each stored vector and returns the minimum. Example: for a degree-4 sequence with vectors (1, 1, 1 | 1, 2, 1) and (1, 2, 1 | 1, 1, 1) and gaps h = (5, 8, 3), v = (3, 9, 6), the values are 43 and 42, so the tree length is 42, the half-perimeter 34 plus the smaller middle gap 8.

## FLUTE inside a global router

FLUTE optimises wirelength only. It knows nothing about congestion, so a global router treats its tree as a starting topology. FastRoute, for example, constructs congestion-driven Steiner trees by shifting tree edges and Steiner points toward less congested gcells, then decomposes the tree into two-pin connections for pattern and maze routing.

The design lesson is the split between a fast, exact-enough geometric kernel and a cost-aware policy layered on top. FLUTE is a pure function from pin coordinates to a tree, deterministic and cacheable, which is exactly what a router wants in its inner loop. Its accuracy is excellent for the low degrees where it is exact and degrades slowly for large nets because of net breaking.

## Going deeper

C. Chu and Y.-C. Wong, "FLUTE: Fast Lookup Table Based Rectilinear Steiner Minimal Tree Algorithm for VLSI Design" (IEEE TCAD, 2008); M. Pan and C. Chu, "FastRoute" papers.
