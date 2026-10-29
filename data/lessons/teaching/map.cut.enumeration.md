# Bottom-up cut enumeration and cut functions

Cut enumeration is the workhorse of every cut-based mapper and of many rewriting and resynthesis passes. Each node's cut set is built from its fanins' sets in one topological sweep, and each cut then yields a small truth table, its cut function, which is what LUT mapping stores and Boolean matching compares against the library.

## Merging fanin cut sets

Visit nodes in topological order. For n = AND(u, w), form every pair (one cut from cuts(u), one from cuts(w)), each set including its trivial cut, and take the union. If cuts(u) = {u}, {a, b} and cuts(w) = {w}, {b, c}, {c, d}, there are 2 x 3 = 6 pairs, giving {u, w}, {u, b, c}, {u, c, d}, {a, b, w}, {a, b, c}, {a, b, c, d}. Unions over k leaves are dropped immediately, and identical unions are merged. In practice cuts are stored as sorted leaf arrays or bit signatures so that size tests and unions are cheap.

## Filtering by size and dominance

After merging, add the trivial cut {n}, then remove every cut that strictly contains another cut of the same node. Such a cut is dominated: its extra leaves enlarge the LUT input count without making any new function available. With k = 3 in the example above, the five feasible unions are kept except for {a, b, c, d}, which is dropped for size; if both {a, b} and {a, b, c} appeared, the second would be dominated. A 64-bit signature per cut, with one bit per hashed leaf, gives a fast necessary test for containment before the exact check.

## Computing the cut function

A cut's function is the node's value as a function of the cut leaves. Treat the leaves as free variables, then evaluate the cone above them for every assignment, respecting complemented edges. For 3 leaves this is 8 evaluations, conveniently done with bit-parallel truth tables: the leaves get the patterns 0xAA, 0xCC, 0xF0, and each AND node combines its fanin tables with bitwise AND and complement. For n = AND(NOT x, y) over leaves x (least significant) and y, the only one is at row x = 0, y = 1, which is row 2, so the table is the integer 4. Mappers derive the cut function from the fanin cuts' tables rather than re-simulating the cone.

## Going deeper

Mishchenko, Cho, Chatterjee and Brayton, "Combinational and sequential mapping with priority cuts", ICCAD 2007; Knuth, "The Art of Computer Programming", volume 4A, on bitwise truth-table manipulation.
