# Coloring-aware detailed routing

At advanced nodes the tightest metal layers are printed with multiple masks: features too close for one exposure are split between masks. That turns every routing decision into a coloring problem, because two nearby wires must land on different masks. A coloring-aware detailed router avoids creating conflicts that no mask assignment can satisfy, rather than leaving them for decomposition to fail on.

## The conflict graph

Build a graph with one vertex per wire shape on the layer and an edge between any two shapes closer than the same-mask spacing. A legal mask assignment is a proper coloring with as many colors as masks: two masks for litho-etch-litho-etch double patterning (LELE), three for triple patterning.

The layout is decomposable with two masks exactly when the conflict graph is bipartite, that is, has no odd cycle. Example: wires on adjacent tracks conflict where they run alongside each other; those conflicts alone are always two-colorable by track parity. But if two wires on the same track also end too close tip to tip, they conflict too, and a tip-to-tip edge on one track with paths through the neighbouring track can close an odd cycle, which forces a third mask.

## Counting conflicts from geometry

Conflicts come from geometry: adjacent tracks with overlapping x ranges, or the same track with a tip-to-tip gap below the same-mask limit. Example: wires [5, 63] and [103, 194] on track 1, [15, 89] and [99, 159] on track 2, with a tip-to-tip limit of 30 nm: the track-2 wires are 10 nm apart and conflict, the track-1 wires are 40 nm apart and do not, and each pair of overlapping wires across tracks 1 and 2 conflicts.

The edge count measures how constrained the layer is. A router that knows the coloring rules prices moves that add edges, especially edges that would close an odd cycle, and prefers line-end positions that keep tip-to-tip gaps above the limit.

## Coloring violations

Given a proposed mask assignment, a coloring violation is a conflict edge whose two shapes share a mask. Example: if a decomposer colors wires by track parity but two same-track wires conflict tip to tip, every such pair is a violation, because they necessarily share a color.

Violations are fixed by rerouting one of the shapes (moving a line end, changing track), or where the process allows, by a stitch that splits one shape across both masks at the cost of overlay risk. Coloring-aware routers either carry a color with each wire during search (correct by construction) or route first and repair coloring conflicts like any other violation. In self-aligned double patterning (SADP) the rules are different again, since line ends and cuts, not adjacency alone, define what prints.

## Going deeper

D. Z. Pan, B. Yu and J.-R. Gao, "Design for Manufacturing with Emerging Nanolithography" (IEEE TCAD, 2013); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.).
