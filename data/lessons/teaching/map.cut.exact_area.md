# Exact local area by reference counting

Area flow is an estimate. Exact local area asks a sharper question against the current mapping: if node n switched to this cut, how much area would actually be added, counting each newly needed gate once? Reference counts answer it exactly for the local change. It is the final area-recovery pass in ABC's mappers and the reason they reach low LUT and cell counts.

## Reference counts in a mapping

Given a current cover, the **reference count** of a node is the number of selected gates or outputs that use it as an input. A node with count 0 is not built; a node with count 1 or more is. When a cut of n is chosen, each leaf is referenced once more. If a leaf's count rises from 0 to 1, that leaf must now be built through its own best cut, whose leaves are referenced in turn. When n drops a cut, the reverse happens: counts fall and nodes reaching 0 disappear with their area.

## Computing the exact area of a cut

The exact area of cut C at node n is A(C) plus, for every leaf whose count rises from 0, the area of that leaf's best cut computed the same way, recursively; a node is added once even if several leaves reach it. Example: cut {c, q, r} with gate area 4. c is a primary input. q (area 2, cut {a, w}) and r (area 5, cut {a, u}) both have count 0, so both are added; w and u already have counts of at least 1, so they add nothing. Exact area = 4 + 2 + 5 = 11. The procedure references, sums, then dereferences to restore the counts.

## Why it improves on area flow

Exact local area sees the real sharing in the current cover. A leaf already used elsewhere costs nothing, whatever its area flow, and two leaves that need the same unbuilt node pay for it once. Visiting nodes in topological order and switching each to the cut with the smallest exact area, while respecting required times, never increases the total area of the cover, since each step compares actual increments. The recursion is cut off at a small depth or at nodes already referenced, so the pass is fast. Mappers typically run one or two area-flow passes and then one or two exact-area passes.

## Going deeper

Mishchenko, Chatterjee and Brayton, "Improvements to technology mapping for LUT-based FPGAs", IEEE TCAD 2007; the ABC source code for if and amap.
