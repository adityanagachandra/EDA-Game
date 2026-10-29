# Channel routing by the left-edge algorithm

The left-edge algorithm of Hashimoto and Stevens (1971) is the first and simplest channel router: sort trunks by left end and pack them greedily onto tracks. Without vertical constraints it is optimal; with them, a constrained version respects the vertical constraint graph by filling tracks from the top. It is still the template for interval packing in detailed routers.

## Unconstrained left edge

Sort the nets by the left end of their spans. Fill track 1: scan the sorted list and place each net whose left end lies to the right of the last net placed on the track. When the scan ends, start track 2 with the remaining nets, and so on.

Example: spans [9, 16], [4, 7], [3, 10], [2, 3] and [7, 10]. Track 1 takes [2, 3], then [4, 7], then [9, 16]; track 2 takes [3, 10]; track 3 takes [7, 10]. Three tracks, which equals the density: without vertical constraints the greedy is optimal, because when it opens a new track every existing track is busy at the current net's left end, so that many spans overlap there.

## Constrained left edge

With vertical constraints, fill tracks from the top. For each track, scan the unplaced nets by left end and add a net only if every net that must lie above it (its VCG predecessors) is already on an earlier track, and it fits horizontally. A net placed on this track blocks its VCG successors until the next track.

Example: spans net 1 [2, 7], net 2 [2, 6], net 3 [5, 8], net 4 [5, 7] and net 5 [3, 6] all overlap at columns 5 and 6, so five tracks are needed in any case, and the VCG only decides their top-to-bottom order. The algorithm is fast, $O(n \log n)$ plus the constraint checks, and it is guaranteed to finish whenever the VCG is acyclic, but it can use more tracks than the best routing.

## The top track and its limits

Only nets with no VCG predecessor can go on the top track, so the first pass is the unconstrained greedy restricted to those nets. Example: if nets 4 and 1 are the only predecessor-free nets and their spans do not overlap, the top track holds both.

The constrained greedy has two weaknesses. It ignores how long a net's chain of successors is, so it can delay a net that heads a long chain and waste tracks; refinements order candidates by chain length. And it fails outright on a cyclic VCG, where no net in the cycle is ever admitted. Doglegs, splitting a net's trunk at intermediate pins, fix both problems, and are the next step up.

## Going deeper

A. Hashimoto and J. Stevens (1971); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.), on channel routing.
