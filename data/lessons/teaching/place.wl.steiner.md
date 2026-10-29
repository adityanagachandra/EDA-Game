# Steiner and spanning-tree wirelength estimates

A routed net is a rectilinear tree, and its shortest form is the rectilinear Steiner minimum tree (RSMT). Placers usually optimise HPWL, a lower bound, but timing estimation, congestion models and final reporting need a tree length. Spanning trees and Steiner trees bracket the truth, and knowing when each is exact tells a tool designer which estimate to use where.

## Three pins: HPWL is exact

For a three-pin net, the RSMT connects the pins through a single Steiner point at (median x, median y), and its length equals the HPWL. Pins at (0, 0), (10, 4) and (6, 12) have median point (6, 4); the three branches have lengths 10, 4 and 8, total 22, while the bounding box is 10 by 12, HPWL 22. Two-pin nets are trivially exact as well. Since most nets in real designs have two or three pins, HPWL is an exact routed-length estimate for the majority of nets, ignoring detours. The gap opens only for larger nets.

## Spanning trees as an upper estimate

A rectilinear minimum spanning tree (RMST) connects pins only to each other with Manhattan edges. Prim's algorithm builds it in $O(n^2)$ for $n$ pins, or $O(n \log n)$ with geometric tricks. Hwang's theorem bounds its quality: the RMST is at most $3/2$ times the RSMT. For pins at (0, 0), (10, 0), (0, 10) and (10, 10), the RMST uses three edges of length 10, total 30, while the HPWL is 20. The RMST is easy and gives a valid routing topology, which is why it is used as a starting point for Steiner heuristics and for Elmore-delay estimation of nets.

## Steiner points on the Hanan grid

Hanan's theorem says some RSMT has all its Steiner points on the Hanan grid, the intersections of the horizontal and vertical lines through the pins, and at most $n - 2$ Steiner points are ever needed. That turns the search into a finite one. Take four pins at (0, 5), (5, 0), (10, 5) and (5, 10). Their Hanan grid has lines at x = 0, 5, 10 and y = 0, 5, 10, so (5, 5) is a candidate. A star from (5, 5) has four arms of length 5, total 20, equal to the HPWL; the best spanning tree needs three edges of length 10, total 30. This instance meets Hwang's 3/2 bound exactly. Exact RSMT is NP-hard in general, so tools use fast methods such as the FLUTE lookup tables, which are optimal for small nets (up to nine pins) and heuristic beyond.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), sections 4.2 and 5.6; Chu and Wong, "FLUTE: fast lookup table based rectilinear Steiner minimal tree algorithm" (IEEE TCAD, 2008).
