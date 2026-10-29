# Track assignment

Between global and detailed routing, many flows run track assignment: the long horizontal or vertical runs that global routing put in a row or column of gcells (a panel) are assigned to specific tracks before any maze search. TritonRoute does this, and the problem is interval packing, the same mathematics as channel routing, with overlap as the cost of a bad assignment.

## Interval packing in a panel

In a horizontal panel each global segment covers an interval of gcells $[l, r]$. Two segments conflict if their intervals share a gcell: they cannot share a track there. The conflict graph is an interval graph, and interval graphs are perfect: the fewest tracks that hold all segments without conflict equals the largest clique, which is the maximum number of segments covering any one gcell.

Example: segments [17, 18], [12, 19], [18, 19], [1, 6] and [9, 14] all cover gcell 18 three at a time ([17, 18], [12, 19], [18, 19]), so the panel needs 3 tracks. A left-edge sweep, sorting segments by left end and placing each on the first track whose last segment ended before it, achieves the bound.

## Conflicts between segments

The conflict graph has an edge for every pair of segments that share a gcell. With $n$ segments there are $n(n-1)/2$ pairs to check, or a sort by left end finds the overlaps in $O(n \log n + k)$ for $k$ conflicts. Example: segments [16, 20], [4, 14], [5, 6], [13, 19] and [18, 20] have five conflicting pairs.

The edges are what a track assigner must keep apart, and they carry more than a yes or no: the overlap length matters, and so do soft conflicts, such as two segments on adjacent tracks running alongside an aggressor net, or a segment ending near a pin that its neighbour must reach. Real track assigners weigh all of these in one cost.

## Overlap cost when tracks run out

A panel often has fewer tracks than its density, so some overlap is unavoidable. A track assigner then minimises a cost such as total overlap: for every pair of segments on the same track, the number of gcells they share. Example: with three tracks, an assignment that puts [11, 18] and [14, 20] on one track pays 5, the length of [14, 18].

Every unit of overlap is work handed to detailed routing, which must resolve it by moving one segment off its track with a detour. Assigners therefore also consider pin positions, blockages and previously assigned nets, and solve greedily, with matching, or with iterative improvement. A good assignment lets the detailed router spend its time on pin access and rule details rather than on long rearrangements.

## Going deeper

A. B. Kahng, L. Wang and B. Xu, "TritonRoute: The Open-Source Detailed Router" (IEEE TCAD, 2021); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.).
